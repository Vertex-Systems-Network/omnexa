package events

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

func TestDurableConsumerResumesFromLastAcceptedCheckpoint(t *testing.T) {
	calls := 0
	registry := durableTestRegistry(t, false, func(context.Context, Envelope) error {
		calls++
		return nil
	})
	store := NewMemoryCheckpointStore()
	binding := durableTestBinding("")

	firstConsumer := durableTestConsumer(t, registry, store, binding)
	position, err := firstConsumer.ResumePosition(context.Background())
	if err != nil {
		t.Fatalf("ResumePosition() before first delivery error = %v", err)
	}
	if position != 1 {
		t.Fatalf("initial resume position = %d, want 1", position)
	}

	first := testEnvelope(t)
	processErr := firstConsumer.Process(context.Background(), Delivery{Position: 1, Envelope: first})
	if processErr != nil {
		t.Fatalf("first Process() error = %v", processErr)
	}
	checkpoint, exists, err := firstConsumer.LastCheckpoint(context.Background())
	if err != nil {
		t.Fatalf("LastCheckpoint() error = %v", err)
	}
	if !exists || checkpoint.Position != 1 || checkpoint.EventID != first.ID {
		t.Fatalf("checkpoint after first delivery = %#v exists=%v", checkpoint, exists)
	}

	restarted := durableTestConsumer(t, registry, store, binding)
	position, err = restarted.ResumePosition(context.Background())
	if err != nil {
		t.Fatalf("ResumePosition() after restart error = %v", err)
	}
	if position != 2 {
		t.Fatalf("resume position after restart = %d, want 2", position)
	}
	second := testEnvelope(t)
	if err := restarted.Process(context.Background(), Delivery{Position: 2, Envelope: second}); err != nil {
		t.Fatalf("second Process() error = %v", err)
	}
	if calls != 2 {
		t.Fatalf("handler calls = %d, want 2", calls)
	}
}

func TestDurableConsumerFailedOrCancelledWorkDoesNotAdvance(t *testing.T) {
	t.Run("handler failure", func(t *testing.T) {
		registry := durableTestRegistry(t, false, func(context.Context, Envelope) error {
			return errors.New("private handler diagnostic")
		})
		consumer := durableTestConsumer(t, registry, NewMemoryCheckpointStore(), durableTestBinding(""))

		err := consumer.Process(context.Background(), Delivery{Position: 1, Envelope: testEnvelope(t)})
		assertFailureCode(t, err, codeHandlerFailed)
		assertNoCheckpoint(t, consumer)
	})

	t.Run("cancellation after handler acknowledgement candidate", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		registry := durableTestRegistry(t, false, func(context.Context, Envelope) error {
			cancel()
			return nil
		})
		consumer := durableTestConsumer(t, registry, NewMemoryCheckpointStore(), durableTestBinding(""))

		err := consumer.Process(ctx, Delivery{Position: 1, Envelope: testEnvelope(t)})
		assertFailureCode(t, err, codeDurableInterrupted)
		assertNoCheckpoint(t, consumer)
	})
}

func TestDurableConsumerEnforcesContiguousMonotonicProgress(t *testing.T) {
	calls := 0
	registry := durableTestRegistry(t, false, func(context.Context, Envelope) error {
		calls++
		return nil
	})
	consumer := durableTestConsumer(t, registry, NewMemoryCheckpointStore(), durableTestBinding(""))
	first := testEnvelope(t)
	if err := consumer.Process(context.Background(), Delivery{Position: 1, Envelope: first}); err != nil {
		t.Fatalf("first Process() error = %v", err)
	}

	assertFailureCode(t, consumer.Process(context.Background(), Delivery{Position: 1, Envelope: first}), codeCheckpointStale)
	assertFailureCode(t, consumer.Process(context.Background(), Delivery{Position: 3, Envelope: testEnvelope(t)}), codeCheckpointConflict)
	if calls != 1 {
		t.Fatalf("stale/gapped deliveries reached handler: calls = %d, want 1", calls)
	}

	if err := consumer.Process(context.Background(), Delivery{Position: 2, Envelope: testEnvelope(t)}); err != nil {
		t.Fatalf("contiguous second Process() error = %v", err)
	}
	position, err := consumer.ResumePosition(context.Background())
	if err != nil {
		t.Fatalf("ResumePosition() error = %v", err)
	}
	if position != 3 {
		t.Fatalf("resume position = %d, want 3", position)
	}
}

func TestDurableConsumerRejectsConflictingOwnerAndScopeRebinding(t *testing.T) {
	store := NewMemoryCheckpointStore()
	registry := durableTestRegistryWithIdentity(t,
		Producer("urn:omnexa:module:commerce.orders"),
		"orders.durable_projection",
		false,
		func(context.Context, Envelope) error { return nil },
	)
	binding := durableTestBinding("")
	_ = durableTestConsumer(t, registry, store, binding)

	conflictingScope := binding
	conflictingScope.Scope.Stream = "orders.secondary"
	_, err := NewDurableConsumer(context.Background(), registry, store, conflictingScope)
	assertFailureCode(t, err, codeDurableBindingConflict)

	otherRegistry := durableTestRegistryWithIdentity(t,
		Producer("urn:omnexa:module:analytics.orders"),
		"orders.durable_projection",
		false,
		func(context.Context, Envelope) error { return nil },
	)
	otherBinding := binding
	otherBinding.Owner = Producer("urn:omnexa:module:analytics.orders")
	_, err = NewDurableConsumer(context.Background(), otherRegistry, store, otherBinding)
	assertFailureCode(t, err, codeDurableBindingConflict)
}

func TestDurableConsumerTenantScopesAreIsolatedAndMetadataIsPreserved(t *testing.T) {
	tenantOne := tenancy.TenantID(testUUIDv7(t))
	tenantTwo := tenancy.TenantID(testUUIDv7(t))
	var received []Envelope
	registry := durableTestRegistry(t, true, func(_ context.Context, envelope Envelope) error {
		received = append(received, envelope)
		return nil
	})
	store := NewMemoryCheckpointStore()
	firstBinding := durableTestBinding(tenantOne)
	secondBinding := durableTestBinding(tenantTwo)
	firstConsumer := durableTestConsumer(t, registry, store, firstBinding)
	secondConsumer := durableTestConsumer(t, registry, store, secondBinding)

	first := testEnvelopeForTenant(t, tenantOne)
	first.CausationID = CausationID(testUUIDv7(t))
	if err := first.ValidateForTenant(tenantOne); err != nil {
		t.Fatalf("first ValidateForTenant() error = %v", err)
	}
	second := testEnvelopeForTenant(t, tenantTwo)

	if err := firstConsumer.Process(context.Background(), Delivery{Position: 1, Envelope: first}); err != nil {
		t.Fatalf("tenant one Process() error = %v", err)
	}
	if err := secondConsumer.Process(context.Background(), Delivery{Position: 1, Envelope: second}); err != nil {
		t.Fatalf("tenant two Process() error = %v", err)
	}

	firstCheckpoint, exists, err := firstConsumer.LastCheckpoint(context.Background())
	if err != nil || !exists {
		t.Fatalf("tenant one checkpoint error=%v exists=%v", err, exists)
	}
	secondCheckpoint, exists, err := secondConsumer.LastCheckpoint(context.Background())
	if err != nil || !exists {
		t.Fatalf("tenant two checkpoint error=%v exists=%v", err, exists)
	}
	if firstCheckpoint.EventID != first.ID || secondCheckpoint.EventID != second.ID || firstCheckpoint.EventID == secondCheckpoint.EventID {
		t.Fatalf("tenant checkpoints collided: first=%#v second=%#v", firstCheckpoint, secondCheckpoint)
	}

	assertFailureCode(t, firstConsumer.Process(context.Background(), Delivery{Position: 2, Envelope: second}), codeTenantMismatch)
	unchanged, _, err := firstConsumer.LastCheckpoint(context.Background())
	if err != nil {
		t.Fatalf("tenant one checkpoint after mismatch error = %v", err)
	}
	if unchanged != firstCheckpoint {
		t.Fatalf("tenant mismatch advanced another tenant checkpoint: got %#v want %#v", unchanged, firstCheckpoint)
	}

	if len(received) != 2 {
		t.Fatalf("handler received %d envelopes, want 2", len(received))
	}
	got := received[0]
	if got.ID != first.ID || got.Type != first.Type || got.Source != first.Source || got.TenantID != first.TenantID || got.CorrelationID != first.CorrelationID || got.CausationID != first.CausationID || got.Classification != first.Classification || !bytes.Equal(got.Data, first.Data) {
		t.Fatal("durable boundary mutated canonical P04.01 metadata or payload")
	}
}

func TestDurableConsumerAllowsDuplicateDeliveryWhenCheckpointWriteFails(t *testing.T) {
	const privateStoreCause = "postgres://private-host/restricted-checkpoint-state"
	calls := 0
	registry := durableTestRegistry(t, false, func(context.Context, Envelope) error {
		calls++
		return nil
	})
	store := &failOnceAdvanceStore{
		base:         NewMemoryCheckpointStore(),
		privateCause: privateStoreCause,
	}
	binding := durableTestBinding("")
	consumer := durableTestConsumer(t, registry, store, binding)
	envelope := testEnvelope(t)

	firstErr := consumer.Process(context.Background(), Delivery{Position: 1, Envelope: envelope})
	assertFailureCode(t, firstErr, codeCheckpointWriteFailed)
	if strings.Contains(firstErr.Error(), privateStoreCause) {
		t.Fatalf("checkpoint failure leaked private provider diagnostic: %v", firstErr)
	}
	assertNoCheckpoint(t, consumer)

	restarted := durableTestConsumer(t, registry, store, binding)
	if err := restarted.Process(context.Background(), Delivery{Position: 1, Envelope: envelope}); err != nil {
		t.Fatalf("replayed Process() error = %v", err)
	}
	if calls != 2 {
		t.Fatalf("crash-window duplicate delivery was hidden: calls = %d, want 2", calls)
	}
	checkpoint, exists, err := restarted.LastCheckpoint(context.Background())
	if err != nil || !exists || checkpoint.Position != 1 || checkpoint.EventID != envelope.ID {
		t.Fatalf("checkpoint after replay = %#v exists=%v err=%v", checkpoint, exists, err)
	}
}

func TestDurableConsumerRejectsMalformedCheckpointAndSanitizesStoreFailure(t *testing.T) {
	t.Run("malformed checkpoint", func(t *testing.T) {
		called := false
		registry := durableTestRegistry(t, false, func(context.Context, Envelope) error {
			called = true
			return nil
		})
		store := &malformedCheckpointStore{}
		consumer := durableTestConsumer(t, registry, store, durableTestBinding(""))

		err := consumer.Process(context.Background(), Delivery{Position: 2, Envelope: testEnvelope(t)})
		assertFailureCode(t, err, codeCheckpointMalformed)
		if called {
			t.Fatal("malformed checkpoint state reached handler")
		}
	})

	t.Run("private read failure", func(t *testing.T) {
		const privateCause = "redis://private-topology/restricted"
		registry := durableTestRegistry(t, false, func(context.Context, Envelope) error { return nil })
		store := &loadFailureStore{privateCause: privateCause}
		consumer := durableTestConsumer(t, registry, store, durableTestBinding(""))

		_, _, err := consumer.LastCheckpoint(context.Background())
		assertFailureCode(t, err, codeCheckpointReadFailed)
		if strings.Contains(err.Error(), privateCause) {
			t.Fatalf("checkpoint read failure leaked private provider diagnostic: %v", err)
		}
	})
}

func TestDurableConsumerRejectsInvalidRouteAndScopeBeforeHandler(t *testing.T) {
	called := false
	registry := durableTestRegistry(t, false, func(context.Context, Envelope) error {
		called = true
		return nil
	})
	binding := durableTestBinding("")
	consumer := durableTestConsumer(t, registry, NewMemoryCheckpointStore(), binding)

	otherRoute := testEnvelope(t)
	otherRoute.Type = EventType("commerce.order.updated.v1")
	otherRoute.DataSchema = SchemaID("urn:omnexa:event-schema:commerce.order.updated:v1")
	assertFailureCode(t, consumer.Process(context.Background(), Delivery{Position: 1, Envelope: otherRoute}), codeDurableBindingConflict)
	if called {
		t.Fatal("route-conflicting delivery reached handler")
	}

	badBinding := binding
	badBinding.Scope.Stream = "orders/escape"
	_, err := NewDurableConsumer(context.Background(), registry, NewMemoryCheckpointStore(), badBinding)
	assertFailureCode(t, err, codeDurableInvalid)
}

func durableTestRegistry(t *testing.T, tenantScoped bool, handler Handler) *Registry {
	t.Helper()
	return durableTestRegistryWithIdentity(t,
		Producer("urn:omnexa:module:commerce.orders"),
		"orders.durable_projection",
		tenantScoped,
		handler,
	)
}

func durableTestRegistryWithIdentity(t *testing.T, owner Producer, consumerID string, tenantScoped bool, handler Handler) *Registry {
	t.Helper()
	registry := NewRegistry()
	if err := registry.Register(Registration{
		Owner:        owner,
		ConsumerID:   consumerID,
		EventTypes:   []EventType{"commerce.order.created.v1"},
		TenantScoped: tenantScoped,
		Handler:      handler,
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	return registry
}

func durableTestBinding(tenantID tenancy.TenantID) DurableBinding {
	return DurableBinding{
		Owner:      Producer("urn:omnexa:module:commerce.orders"),
		ConsumerID: "orders.durable_projection",
		EventType:  EventType("commerce.order.created.v1"),
		Scope: DurableScope{
			Stream:    "orders.events",
			Partition: "primary",
			TenantID:  tenantID,
		},
	}
}

func durableTestConsumer(t *testing.T, registry *Registry, store CheckpointStore, binding DurableBinding) *DurableConsumer {
	t.Helper()
	consumer, err := NewDurableConsumer(context.Background(), registry, store, binding)
	if err != nil {
		t.Fatalf("NewDurableConsumer() error = %v", err)
	}
	return consumer
}

func testEnvelopeForTenant(t *testing.T, tenantID tenancy.TenantID) Envelope {
	t.Helper()
	envelope := testEnvelope(t)
	envelope.TenantID = tenantID
	if err := envelope.ValidateForTenant(tenantID); err != nil {
		t.Fatalf("ValidateForTenant() error = %v", err)
	}
	return envelope
}

func assertNoCheckpoint(t *testing.T, consumer *DurableConsumer) {
	t.Helper()
	checkpoint, exists, err := consumer.LastCheckpoint(context.Background())
	if err != nil {
		t.Fatalf("LastCheckpoint() error = %v", err)
	}
	if exists || checkpoint != (Checkpoint{}) {
		t.Fatalf("unexpected checkpoint = %#v exists=%v", checkpoint, exists)
	}
}

type failOnceAdvanceStore struct {
	base         *MemoryCheckpointStore
	failed       bool
	privateCause string
}

func (store *failOnceAdvanceStore) Bind(ctx context.Context, binding DurableBinding) (BindingResult, error) {
	return store.base.Bind(ctx, binding)
}

func (store *failOnceAdvanceStore) Load(ctx context.Context, binding DurableBinding) (Checkpoint, bool, error) {
	return store.base.Load(ctx, binding)
}

func (store *failOnceAdvanceStore) Advance(ctx context.Context, binding DurableBinding, expected uint64, next Checkpoint) (AdvanceResult, error) {
	if !store.failed {
		store.failed = true
		return AdvanceResultUnknown, errors.New(store.privateCause)
	}
	return store.base.Advance(ctx, binding, expected, next)
}

type malformedCheckpointStore struct{}

func (*malformedCheckpointStore) Bind(context.Context, DurableBinding) (BindingResult, error) {
	return BindingAccepted, nil
}

func (*malformedCheckpointStore) Load(context.Context, DurableBinding) (Checkpoint, bool, error) {
	return Checkpoint{Position: 1, EventID: EventID("malformed-event-id")}, true, nil
}

func (*malformedCheckpointStore) Advance(context.Context, DurableBinding, uint64, Checkpoint) (AdvanceResult, error) {
	return CheckpointAdvanced, nil
}

type loadFailureStore struct {
	privateCause string
}

func (*loadFailureStore) Bind(context.Context, DurableBinding) (BindingResult, error) {
	return BindingAccepted, nil
}

func (store *loadFailureStore) Load(context.Context, DurableBinding) (Checkpoint, bool, error) {
	return Checkpoint{}, false, errors.New(store.privateCause)
}

func (*loadFailureStore) Advance(context.Context, DurableBinding, uint64, Checkpoint) (AdvanceResult, error) {
	return CheckpointAdvanced, nil
}

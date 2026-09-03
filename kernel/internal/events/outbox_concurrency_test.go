package events

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5"
)

func TestPostgresOutboxConcurrentRelaysPreservePublicationState(t *testing.T) {
	ctx, pool, store := setupP0404ReliabilityDatabase(t)
	scope := OutboxScope{
		Owner:    Producer("urn:omnexa:module:concurrency"),
		TenantID: tenancy.TenantID(mustP0404UUIDv7(t)),
	}
	envelope := mustP0404Envelope(t, scope, "concurrent-relay")
	if err := database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		result, enqueueErr := EnqueueOutbox(ctx, tx, store, scope, envelope)
		if enqueueErr != nil {
			return enqueueErr
		}
		if result != OutboxEnqueued {
			return errors.New("concurrency fixture enqueue did not report enqueued")
		}
		return nil
	}); err != nil {
		t.Fatalf("concurrency fixture enqueue error = %v", err)
	}

	arrived := make(chan struct{}, 2)
	release := make(chan struct{})
	barrierStore := &pendingBarrierOutboxStore{
		OutboxStore: store,
		arrived:     arrived,
		release:     release,
	}
	var publishes atomic.Int32
	publisher, err := NewPublisher(func(_ context.Context, candidate Envelope) error {
		if candidate.ID != envelope.ID {
			return errors.New("concurrent relay published unexpected event identity")
		}
		publishes.Add(1)
		return nil
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}

	type relayOutcome struct {
		report RelayReport
		err    error
	}
	outcomes := make(chan relayOutcome, 2)
	for range 2 {
		relay, relayErr := NewOutboxRelay(barrierStore, publisher)
		if relayErr != nil {
			t.Fatalf("NewOutboxRelay() error = %v", relayErr)
		}
		go func() {
			report, runErr := relay.RelayPending(ctx, scope, 1)
			outcomes <- relayOutcome{report: report, err: runErr}
		}()
	}

	// Both relays must have read the same committed pending revision before
	// either is allowed to publish/mark. This deterministically exercises the
	// accepted duplicate-publication window and PostgreSQL CAS mark behavior.
	<-arrived
	<-arrived
	close(release)

	first := <-outcomes
	second := <-outcomes
	if first.err != nil || second.err != nil {
		t.Fatalf("concurrent relay errors = %v / %v", first.err, second.err)
	}
	accepted := first.report.AcceptedPublications + second.report.AcceptedPublications
	marked := first.report.PublishedMarks + second.report.PublishedMarks
	duplicates := first.report.DuplicateMarkObserved + second.report.DuplicateMarkObserved
	if accepted != 2 || marked != 1 || duplicates != 1 || publishes.Load() != 2 {
		t.Fatalf("concurrent relay totals accepted=%d marked=%d duplicates=%d publishes=%d; first=%#v second=%#v", accepted, marked, duplicates, publishes.Load(), first.report, second.report)
	}
	pending, err := store.Pending(ctx, scope, 1)
	if err != nil {
		t.Fatalf("Pending() after concurrent relays error = %v", err)
	}
	if len(pending) != 0 {
		t.Fatalf("concurrent relays left corrupted pending state = %#v", pending)
	}
}

func TestPostgresOutboxConcurrentTenantRelaysRemainIsolated(t *testing.T) {
	ctx, pool, store := setupP0404ReliabilityDatabase(t)
	owner := Producer("urn:omnexa:module:tenant-isolation")
	scopeA := OutboxScope{Owner: owner, TenantID: tenancy.TenantID(mustP0404UUIDv7(t))}
	scopeB := OutboxScope{Owner: owner, TenantID: tenancy.TenantID(mustP0404UUIDv7(t))}
	envelopeA := mustP0404Envelope(t, scopeA, "tenant-a")
	envelopeB := mustP0404Envelope(t, scopeB, "tenant-b")

	for _, fixture := range []struct {
		scope    OutboxScope
		envelope Envelope
	}{{scopeA, envelopeA}, {scopeB, envelopeB}} {
		fixture := fixture
		if err := database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
			result, enqueueErr := EnqueueOutbox(ctx, tx, store, fixture.scope, fixture.envelope)
			if enqueueErr != nil {
				return enqueueErr
			}
			if result != OutboxEnqueued {
				return errors.New("tenant-isolation fixture enqueue did not report enqueued")
			}
			return nil
		}); err != nil {
			t.Fatalf("tenant-isolation fixture enqueue error = %v", err)
		}
	}

	var mutex sync.Mutex
	publishedByTenant := map[tenancy.TenantID][]EventID{}
	publisher, err := NewPublisher(func(_ context.Context, candidate Envelope) error {
		mutex.Lock()
		publishedByTenant[candidate.TenantID] = append(publishedByTenant[candidate.TenantID], candidate.ID)
		mutex.Unlock()
		return nil
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}

	type tenantOutcome struct {
		scope  OutboxScope
		report RelayReport
		err    error
	}
	outcomes := make(chan tenantOutcome, 2)
	for _, scope := range []OutboxScope{scopeA, scopeB} {
		scope := scope
		relay, relayErr := NewOutboxRelay(store, publisher)
		if relayErr != nil {
			t.Fatalf("NewOutboxRelay() error = %v", relayErr)
		}
		go func() {
			report, runErr := relay.RelayPending(ctx, scope, 8)
			outcomes <- tenantOutcome{scope: scope, report: report, err: runErr}
		}()
	}

	for range 2 {
		outcome := <-outcomes
		if outcome.err != nil {
			t.Fatalf("tenant %s RelayPending() error = %v", outcome.scope.TenantID, outcome.err)
		}
		if outcome.report.PendingObserved != 1 || outcome.report.AcceptedPublications != 1 || outcome.report.PublishedMarks != 1 || outcome.report.DuplicateMarkObserved != 0 {
			t.Fatalf("tenant %s relay report = %#v", outcome.scope.TenantID, outcome.report)
		}
	}

	mutex.Lock()
	publishedA := append([]EventID(nil), publishedByTenant[scopeA.TenantID]...)
	publishedB := append([]EventID(nil), publishedByTenant[scopeB.TenantID]...)
	mutex.Unlock()
	if len(publishedA) != 1 || publishedA[0] != envelopeA.ID {
		t.Fatalf("tenant A publications = %#v; want only %s", publishedA, envelopeA.ID)
	}
	if len(publishedB) != 1 || publishedB[0] != envelopeB.ID {
		t.Fatalf("tenant B publications = %#v; want only %s", publishedB, envelopeB.ID)
	}
	for _, scope := range []OutboxScope{scopeA, scopeB} {
		pending, pendingErr := store.Pending(ctx, scope, 8)
		if pendingErr != nil {
			t.Fatalf("tenant %s Pending() after relay error = %v", scope.TenantID, pendingErr)
		}
		if len(pending) != 0 {
			t.Fatalf("tenant %s retained pending records = %#v", scope.TenantID, pending)
		}
	}
}

type pendingBarrierOutboxStore struct {
	OutboxStore
	arrived chan<- struct{}
	release <-chan struct{}
}

func (store *pendingBarrierOutboxStore) Pending(ctx context.Context, scope OutboxScope, limit uint32) ([]OutboxRecord, error) {
	records, err := store.OutboxStore.Pending(ctx, scope, limit)
	if err != nil {
		return nil, err
	}
	store.arrived <- struct{}{}
	select {
	case <-store.release:
		return records, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

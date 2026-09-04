package events

import (
	"context"
	"errors"
	"testing"
)

type inboxTestStore struct {
	claimResult InboxClaimResult
	claimErr    error
	completeErr error

	claimCalls    int
	completeCalls int
	lastClaimTx   OutboxTransaction
	lastCompleteTx OutboxTransaction
	lastClaim     InboxRecord
	lastComplete  InboxRecord
}

func (store *inboxTestStore) Claim(_ context.Context, tx OutboxTransaction, record InboxRecord) (InboxClaimResult, error) {
	store.claimCalls++
	store.lastClaimTx = tx
	store.lastClaim = record
	return store.claimResult, store.claimErr
}

func (store *inboxTestStore) Complete(_ context.Context, tx OutboxTransaction, record InboxRecord) error {
	store.completeCalls++
	store.lastCompleteTx = tx
	store.lastComplete = record
	return store.completeErr
}

func TestNewInboxRecordScopesCanonicalEventToConsumerBinding(t *testing.T) {
	envelope := testEnvelope(t)
	binding := testInboxBinding(envelope, "billing.projection")

	record, err := NewInboxRecord(binding, envelope)
	if err != nil {
		t.Fatalf("NewInboxRecord() error = %v", err)
	}
	if record.Identity.EventID != envelope.ID || record.Identity.Owner != binding.Owner || record.Identity.ConsumerID != binding.ConsumerID {
		t.Fatalf("processing identity changed canonical event/consumer identity: %#v", record.Identity)
	}
	if record.Identity.EventType != binding.EventType || record.Identity.Stream != binding.Scope.Stream || record.Identity.Partition != binding.Scope.Partition || record.Identity.TenantID != binding.Scope.TenantID {
		t.Fatalf("processing identity dropped route/scope evidence: %#v", record.Identity)
	}

	secondBinding := testInboxBinding(envelope, "analytics.projection")
	second, err := NewInboxRecord(secondBinding, envelope)
	if err != nil {
		t.Fatalf("NewInboxRecord(second consumer) error = %v", err)
	}
	if second.Identity == record.Identity {
		t.Fatal("same EventID became a global cross-consumer processing lock")
	}
	if second.Fingerprint != record.Fingerprint {
		t.Fatal("consumer scope unexpectedly changed canonical envelope fingerprint")
	}
}

func TestNewInboxRecordFingerprintDetectsSameIdentityContentReuse(t *testing.T) {
	envelope := testEnvelope(t)
	binding := testInboxBinding(envelope, "billing.projection")
	first, err := NewInboxRecord(binding, envelope)
	if err != nil {
		t.Fatalf("NewInboxRecord(first) error = %v", err)
	}

	changed := envelope
	changed.Data = append([]byte(nil), envelope.Data...)
	changed.Subject = "same-id-different-canonical-content"
	second, err := NewInboxRecord(binding, changed)
	if err != nil {
		t.Fatalf("NewInboxRecord(changed) error = %v", err)
	}
	if first.Identity != second.Identity {
		t.Fatal("same canonical EventID/binding did not retain one processing identity")
	}
	if first.Fingerprint == second.Fingerprint {
		t.Fatal("same processing identity with conflicting canonical content had identical fingerprint")
	}
}

func TestNewInboxRecordRejectsRouteOrTenantRebinding(t *testing.T) {
	envelope := testEnvelope(t)
	binding := testInboxBinding(envelope, "billing.projection")
	binding.EventType = EventType("billing.other.changed.v1")
	_, err := NewInboxRecord(binding, envelope)
	assertFailureCode(t, err, codeInboxConflict)

	binding = testInboxBinding(envelope, "billing.projection")
	binding.Scope.TenantID = "01990f6e-1f30-7000-8000-000000000099"
	_, err = NewInboxRecord(binding, envelope)
	if err == nil {
		t.Fatal("tenant scope rebinding was accepted")
	}
}

func TestApplyInboxFirstDeliveryUsesExactCallerTransactionAndCompletesAfterMutation(t *testing.T) {
	envelope := testEnvelope(t)
	binding := testInboxBinding(envelope, "billing.projection")
	tx := &outboxTestTx{}
	store := &inboxTestStore{claimResult: InboxClaimed}
	mutationCalls := 0

	result, err := ApplyInbox(context.Background(), tx, store, binding, envelope, func(_ context.Context, mutationTx OutboxTransaction) error {
		mutationCalls++
		if mutationTx != tx {
			t.Fatal("protected mutation did not receive exact caller transaction")
		}
		if store.completeCalls != 0 {
			t.Fatal("inbox completion was recorded before protected mutation")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ApplyInbox() error = %v", err)
	}
	if result != InboxApplied || mutationCalls != 1 || store.claimCalls != 1 || store.completeCalls != 1 {
		t.Fatalf("first application = result %v mutation %d claim %d complete %d", result, mutationCalls, store.claimCalls, store.completeCalls)
	}
	if store.lastClaimTx != tx || store.lastCompleteTx != tx {
		t.Fatal("claim/completion did not use the exact protected-mutation transaction")
	}
	if store.lastClaim != store.lastComplete {
		t.Fatal("completion changed processing identity/fingerprint evidence")
	}
}

func TestApplyInboxDuplicateSkipsProtectedMutation(t *testing.T) {
	envelope := testEnvelope(t)
	store := &inboxTestStore{claimResult: InboxAlreadyCompleted}
	called := false

	result, err := ApplyInbox(
		context.Background(),
		&outboxTestTx{},
		store,
		testInboxBinding(envelope, "billing.projection"),
		envelope,
		func(context.Context, OutboxTransaction) error {
			called = true
			return nil
		},
	)
	if err != nil || result != InboxAlreadyApplied {
		t.Fatalf("duplicate = (%v, %v), want already-applied without error", result, err)
	}
	if called || store.completeCalls != 0 {
		t.Fatal("duplicate redelivery re-ran protected mutation or completion")
	}
}

func TestApplyInboxConflictAndConcurrentFailClosedWithoutMutation(t *testing.T) {
	envelope := testEnvelope(t)
	binding := testInboxBinding(envelope, "billing.projection")
	called := false
	mutation := func(context.Context, OutboxTransaction) error {
		called = true
		return nil
	}

	conflict := &inboxTestStore{claimResult: InboxIdentityConflict}
	result, err := ApplyInbox(context.Background(), &outboxTestTx{}, conflict, binding, envelope, mutation)
	if result != InboxConflict {
		t.Fatalf("identity conflict result = %v, want conflict", result)
	}
	assertFailureCode(t, err, codeInboxConflict)
	if called || conflict.completeCalls != 0 {
		t.Fatal("identity conflict reached mutation/completion")
	}

	concurrent := &inboxTestStore{claimResult: InboxConcurrentProcessing}
	result, err = ApplyInbox(context.Background(), &outboxTestTx{}, concurrent, binding, envelope, mutation)
	if result != InboxConcurrent {
		t.Fatalf("concurrent result = %v, want concurrent", result)
	}
	assertFailureCode(t, err, codeInboxConcurrent)
	if called || concurrent.completeCalls != 0 {
		t.Fatal("concurrent unresolved attempt reached mutation/completion")
	}
}

func TestApplyInboxMutationFailureLeavesCompletionUnwrittenAndPreservesError(t *testing.T) {
	envelope := testEnvelope(t)
	store := &inboxTestStore{claimResult: InboxClaimed}
	sentinel := errors.New("authorized mutation rejected")

	result, err := ApplyInbox(
		context.Background(),
		&outboxTestTx{},
		store,
		testInboxBinding(envelope, "billing.projection"),
		envelope,
		func(context.Context, OutboxTransaction) error { return sentinel },
	)
	if result != InboxApplyResultUnknown || !errors.Is(err, sentinel) {
		t.Fatalf("mutation failure = (%v, %v), want unchanged sentinel", result, err)
	}
	if store.completeCalls != 0 {
		t.Fatal("failed mutation created a completion record")
	}
}

func TestApplyInboxCompletionFailureForcesCallerTransactionFailure(t *testing.T) {
	const privateCause = "postgres private hostname and row detail"
	envelope := testEnvelope(t)
	store := &inboxTestStore{
		claimResult: InboxClaimed,
		completeErr: errors.New(privateCause),
	}

	result, err := ApplyInbox(
		context.Background(),
		&outboxTestTx{},
		store,
		testInboxBinding(envelope, "billing.projection"),
		envelope,
		func(context.Context, OutboxTransaction) error { return nil },
	)
	if result != InboxApplyResultUnknown {
		t.Fatalf("completion failure result = %v, want unknown so caller rolls back", result)
	}
	assertFailureCode(t, err, codeInboxStoreFailed)
	if err != nil && containsPrivateText(err.Error(), privateCause) {
		t.Fatalf("completion failure leaked provider details: %v", err)
	}
}

func TestApplyInboxRejectsInvalidStoreResultAndInterruptedContext(t *testing.T) {
	envelope := testEnvelope(t)
	binding := testInboxBinding(envelope, "billing.projection")

	invalid := &inboxTestStore{claimResult: InboxClaimResultUnknown}
	result, err := ApplyInbox(
		context.Background(),
		&outboxTestTx{},
		invalid,
		binding,
		envelope,
		func(context.Context, OutboxTransaction) error { return nil },
	)
	if result != InboxApplyResultUnknown {
		t.Fatalf("invalid store result = %v", result)
	}
	assertFailureCode(t, err, codeInboxStateMalformed)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	store := &inboxTestStore{claimResult: InboxClaimed}
	result, err = ApplyInbox(
		ctx,
		&outboxTestTx{},
		store,
		binding,
		envelope,
		func(context.Context, OutboxTransaction) error { return nil },
	)
	if result != InboxApplyResultUnknown || store.claimCalls != 0 {
		t.Fatalf("interrupted request crossed persistence seam: result=%v claims=%d", result, store.claimCalls)
	}
	assertFailureCode(t, err, codeInboxInterrupted)
}

func testInboxBinding(envelope Envelope, consumerID string) DurableBinding {
	return DurableBinding{
		Owner:      Producer("urn:omnexa:module:billing.consumer"),
		ConsumerID: consumerID,
		EventType:  envelope.Type,
		Scope: DurableScope{
			Stream:    "billing.events",
			Partition: "primary",
			TenantID:  envelope.TenantID,
		},
	}
}

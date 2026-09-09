package events

import (
	"context"
	"errors"
	"testing"
	"time"
)

type p0406T03RetryStore struct {
	record          RetryStateRecord
	transitionCalls int
	transitionErr   error
	order           *[]string
}

func (store *p0406T03RetryStore) Load(_ context.Context, identity InboxIdentity, position uint64) (RetryStateRecord, error) {
	if store == nil || store.record.Identity != identity || store.record.Position != position {
		return RetryStateRecord{}, errors.New("synthetic retry state not found")
	}
	return store.record, nil
}

func (store *p0406T03RetryStore) TransitionClaimed(
	_ context.Context,
	next RetryStateRecord,
	expectedRevision uint64,
	claimToken string,
) (RetryStateRecord, error) {
	if store.order != nil {
		*store.order = append(*store.order, "quarantine")
	}
	store.transitionCalls++
	if store.transitionErr != nil {
		return RetryStateRecord{}, store.transitionErr
	}
	if store.record.Revision != expectedRevision || store.record.ClaimToken != claimToken {
		return RetryStateRecord{}, errors.New("synthetic stale retry transition")
	}
	store.record = next
	return store.record, nil
}

type p0406T03CheckpointStore struct {
	checkpoint   Checkpoint
	exists       bool
	failNext     bool
	advanceCalls int
	order        *[]string
}

func (store *p0406T03CheckpointStore) Bind(context.Context, DurableBinding) (BindingResult, error) {
	return BindingAccepted, nil
}

func (store *p0406T03CheckpointStore) Load(context.Context, DurableBinding) (Checkpoint, bool, error) {
	if !store.exists {
		return Checkpoint{}, false, nil
	}
	return store.checkpoint, true, nil
}

func (store *p0406T03CheckpointStore) Advance(
	_ context.Context,
	_ DurableBinding,
	expectedPosition uint64,
	next Checkpoint,
) (AdvanceResult, error) {
	if store.order != nil {
		*store.order = append(*store.order, "checkpoint")
	}
	store.advanceCalls++
	if store.failNext {
		store.failNext = false
		return AdvanceResultUnknown, errors.New("synthetic checkpoint persistence gap")
	}
	if store.exists {
		if store.checkpoint.Position != expectedPosition || next.Position != store.checkpoint.Position+1 {
			return CheckpointConflict, nil
		}
	} else if expectedPosition != 0 || next.Position != 1 {
		return CheckpointConflict, nil
	}
	store.checkpoint = next
	store.exists = true
	return CheckpointAdvanced, nil
}

func TestRetryQuarantineCommitsBeforeCheckpointAndRecoversCrashGap(t *testing.T) {
	ctx := context.Background()
	claimed := p0406T03ClaimedRecord(t)
	binding := p0406T03Binding(claimed.Identity)
	transition, err := ComposeClaimedRetryResult(
		claimed,
		time.Date(2026, 9, 9, 10, 0, 0, 0, time.UTC),
		errors.New("synthetic terminal failure"),
		true,
	)
	if err != nil {
		t.Fatalf("ComposeClaimedRetryResult() error = %v", err)
	}
	if transition.Outcome != RetryRuntimeQuarantined {
		t.Fatalf("terminal transition outcome = %q, want quarantine", transition.Outcome)
	}

	order := []string{}
	retryStore := &p0406T03RetryStore{record: claimed, order: &order}
	checkpointStore := &p0406T03CheckpointStore{
		checkpoint: Checkpoint{Position: claimed.Position - 1, EventID: testEnvelope(t).ID},
		exists:     true,
		failNext:   true,
		order:      &order,
	}

	result, err := CommitRetryQuarantineBeforeCheckpoint(
		ctx,
		retryStore,
		checkpointStore,
		binding,
		claimed,
		transition,
	)
	assertFailureCode(t, err, codeRetryCheckpointWrite)
	if !result.QuarantineCommitted || result.CheckpointStatus != RetryCheckpointPending {
		t.Fatalf("crash-gap result = %+v", result)
	}
	if retryStore.record.State != RetryStateQuarantined || retryStore.transitionCalls != 1 {
		t.Fatalf("retry terminal state/calls = %q/%d", retryStore.record.State, retryStore.transitionCalls)
	}
	if len(order) != 2 || order[0] != "quarantine" || order[1] != "checkpoint" {
		t.Fatalf("commit order = %v, want [quarantine checkpoint]", order)
	}
	if checkpointStore.checkpoint.Position != claimed.Position-1 {
		t.Fatalf("checkpoint advanced despite synthetic persistence gap: %+v", checkpointStore.checkpoint)
	}

	order = order[:0]
	recovered, err := RecoverRetryQuarantineCheckpoint(
		ctx,
		retryStore,
		checkpointStore,
		binding,
		claimed.Identity,
		claimed.Position,
	)
	if err != nil {
		t.Fatalf("RecoverRetryQuarantineCheckpoint() error = %v", err)
	}
	if !recovered.QuarantineCommitted || recovered.CheckpointStatus != RetryCheckpointAdvanced {
		t.Fatalf("recovered result = %+v", recovered)
	}
	if retryStore.transitionCalls != 1 {
		t.Fatalf("recovery repeated terminal retry mutation: calls=%d", retryStore.transitionCalls)
	}
	if len(order) != 1 || order[0] != "checkpoint" {
		t.Fatalf("recovery order = %v, want checkpoint only", order)
	}
	if checkpointStore.checkpoint != (Checkpoint{Position: claimed.Position, EventID: claimed.Identity.EventID}) {
		t.Fatalf("recovered checkpoint = %+v", checkpointStore.checkpoint)
	}
}

func TestRetryQuarantineDoesNotCheckpointWhenTerminalCommitFails(t *testing.T) {
	ctx := context.Background()
	claimed := p0406T03ClaimedRecord(t)
	binding := p0406T03Binding(claimed.Identity)
	transition, err := ComposeClaimedRetryResult(
		claimed,
		time.Date(2026, 9, 9, 10, 0, 0, 0, time.UTC),
		errors.New("synthetic terminal failure"),
		true,
	)
	if err != nil {
		t.Fatalf("ComposeClaimedRetryResult() error = %v", err)
	}

	retryStore := &p0406T03RetryStore{
		record:        claimed,
		transitionErr: errors.New("synthetic quarantine commit failure"),
	}
	checkpointStore := &p0406T03CheckpointStore{
		checkpoint: Checkpoint{Position: claimed.Position - 1, EventID: testEnvelope(t).ID},
		exists:     true,
	}

	result, err := CommitRetryQuarantineBeforeCheckpoint(
		ctx,
		retryStore,
		checkpointStore,
		binding,
		claimed,
		transition,
	)
	assertFailureCode(t, err, codeRetryCheckpointState)
	if result.QuarantineCommitted {
		t.Fatalf("failed quarantine commit reported committed: %+v", result)
	}
	if checkpointStore.advanceCalls != 0 {
		t.Fatalf("checkpoint was attempted before terminal quarantine commit: calls=%d", checkpointStore.advanceCalls)
	}
}

func TestRetryQuarantineCheckpointRecoveryIsIdempotentForExactAcceptedProgress(t *testing.T) {
	ctx := context.Background()
	record := testRetryStateRecord(t, RetryStateQuarantined)
	binding := p0406T03Binding(record.Identity)
	retryStore := &p0406T03RetryStore{record: record}
	checkpointStore := &p0406T03CheckpointStore{
		checkpoint: Checkpoint{Position: record.Position, EventID: record.Identity.EventID},
		exists:     true,
	}

	result, err := RecoverRetryQuarantineCheckpoint(
		ctx,
		retryStore,
		checkpointStore,
		binding,
		record.Identity,
		record.Position,
	)
	if err != nil {
		t.Fatalf("RecoverRetryQuarantineCheckpoint() error = %v", err)
	}
	if result.CheckpointStatus != RetryCheckpointAlreadyAdvanced || !result.QuarantineCommitted {
		t.Fatalf("idempotent recovery result = %+v", result)
	}
	if checkpointStore.advanceCalls != 0 || retryStore.transitionCalls != 0 {
		t.Fatalf("idempotent recovery mutated state: checkpoint=%d retry=%d", checkpointStore.advanceCalls, retryStore.transitionCalls)
	}
}

func p0406T03ClaimedRecord(t *testing.T) RetryStateRecord {
	t.Helper()
	record := testRetryStateRecord(t, RetryStateScheduled)
	record.Revision = 2
	record.ClaimToken = "01990f6e-1f30-4000-8000-000000000903"
	record.ClaimExpiresAt = record.NextEligibleAt.Add(5 * time.Minute)
	if err := record.Validate(); err != nil {
		t.Fatalf("claimed retry fixture invalid: %v", err)
	}
	return record
}

func p0406T03Binding(identity InboxIdentity) DurableBinding {
	return DurableBinding{
		Owner:      identity.Owner,
		ConsumerID: identity.ConsumerID,
		EventType:  identity.EventType,
		Scope: DurableScope{
			Stream:    identity.Stream,
			Partition: identity.Partition,
			TenantID:  identity.TenantID,
		},
	}
}

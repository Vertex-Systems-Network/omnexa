package events

import (
	"context"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeRetryCheckpointInvalid   failure.Code = "events.retry.checkpoint_invalid"
	codeRetryCheckpointConflict  failure.Code = "events.retry.checkpoint_conflict"
	codeRetryCheckpointMalformed failure.Code = "events.retry.checkpoint_malformed"
	codeRetryCheckpointRead      failure.Code = "events.retry.checkpoint_read_failed"
	codeRetryCheckpointWrite     failure.Code = "events.retry.checkpoint_write_failed"
	codeRetryCheckpointState     failure.Code = "events.retry.checkpoint_state_failed"
)

// RetryCheckpointStateStore is the narrow accepted P04.06 persistence seam used
// by terminal retry/checkpoint composition. PostgresRetryStateStore satisfies it
// without adding a second schema, migration, transaction owner, or worker loop.
type RetryCheckpointStateStore interface {
	Load(context.Context, InboxIdentity, uint64) (RetryStateRecord, error)
	TransitionClaimed(context.Context, RetryStateRecord, uint64, string) (RetryStateRecord, error)
}

// RetryCheckpointStatus records only the bounded checkpoint side of terminal
// retry composition. It is not execution authority, acknowledgement evidence for
// another delivery, or an exactly-once guarantee.
type RetryCheckpointStatus string

const (
	RetryCheckpointPending         RetryCheckpointStatus = "pending"
	RetryCheckpointAdvanced        RetryCheckpointStatus = "advanced"
	RetryCheckpointAlreadyAdvanced RetryCheckpointStatus = "already_advanced"
)

// RetryQuarantineCheckpointResult makes the crash gap explicit. Once
// QuarantineCommitted is true, callers must never rerun the terminal retry
// transition merely because checkpoint persistence failed afterwards; recovery
// must load that already-quarantined row and repair only the exact checkpoint.
type RetryQuarantineCheckpointResult struct {
	QuarantineCommitted bool
	CheckpointStatus    RetryCheckpointStatus
}

// CommitRetryQuarantineBeforeCheckpoint commits one already-composed terminal
// transition through the existing claimed-retry CAS before attempting checkpoint
// advancement. A checkpoint failure after that commit intentionally returns a
// result with QuarantineCommitted=true so restart logic can recover without
// repeating handler execution or terminal-state mutation.
func CommitRetryQuarantineBeforeCheckpoint(
	ctx context.Context,
	retryStore RetryCheckpointStateStore,
	checkpointStore CheckpointStore,
	binding DurableBinding,
	claimed RetryStateRecord,
	transition RetryRuntimeTransition,
) (RetryQuarantineCheckpointResult, error) {
	if err := retryCheckpointContextError(ctx); err != nil {
		return RetryQuarantineCheckpointResult{}, err
	}
	if retryStore == nil || checkpointStore == nil {
		return RetryQuarantineCheckpointResult{}, classifiedFailure(codeRetryCheckpointInvalid, failure.CategoryValidation, "event retry checkpoint composition boundary is invalid")
	}
	if err := validateRetryCheckpointBinding(binding, claimed); err != nil {
		return RetryQuarantineCheckpointResult{}, err
	}
	if claimed.State != RetryStateScheduled || claimed.ClaimToken == "" || claimed.ClaimExpiresAt.IsZero() {
		return RetryQuarantineCheckpointResult{}, classifiedFailure(codeRetryCheckpointInvalid, failure.CategoryConflict, "event retry checkpoint composition requires an authoritative claimed retry")
	}
	if !transition.Persist || transition.Outcome != RetryRuntimeQuarantined || transition.Record.State != RetryStateQuarantined {
		return RetryQuarantineCheckpointResult{}, classifiedFailure(codeRetryCheckpointInvalid, failure.CategoryValidation, "event retry checkpoint composition requires a terminal quarantine transition")
	}
	if err := transition.Record.Validate(); err != nil {
		return RetryQuarantineCheckpointResult{}, err
	}
	if transition.Record.Identity != claimed.Identity ||
		transition.Record.Fingerprint != claimed.Fingerprint ||
		transition.Record.Position != claimed.Position ||
		transition.Record.Policy != claimed.Policy ||
		transition.Record.Revision != claimed.Revision+1 ||
		transition.Record.ClaimToken != "" || !transition.Record.ClaimExpiresAt.IsZero() {
		return RetryQuarantineCheckpointResult{}, classifiedFailure(codeRetryCheckpointMalformed, failure.CategoryInvariant, "event retry checkpoint terminal transition changed immutable or claim evidence")
	}

	stored, err := retryStore.TransitionClaimed(ctx, transition.Record, claimed.Revision, claimed.ClaimToken)
	if err != nil {
		return RetryQuarantineCheckpointResult{}, wrappedFailure(err, codeRetryCheckpointState, failure.CategoryUnavailable, "event retry quarantine could not be committed before checkpoint")
	}
	if stored != transition.Record || stored.State != RetryStateQuarantined {
		return RetryQuarantineCheckpointResult{}, classifiedFailure(codeRetryCheckpointMalformed, failure.CategoryInvariant, "event retry quarantine store returned unexpected committed state")
	}

	result := RetryQuarantineCheckpointResult{
		QuarantineCommitted: true,
		CheckpointStatus:    RetryCheckpointPending,
	}
	status, checkpointErr := advanceRetryQuarantineCheckpoint(ctx, checkpointStore, binding, stored)
	result.CheckpointStatus = status
	if checkpointErr != nil {
		return result, checkpointErr
	}
	return result, nil
}

// RecoverRetryQuarantineCheckpoint repairs only the checkpoint side of a
// previously committed quarantine/checkpoint crash gap. It never invokes a
// handler and never calls TransitionClaimed. Recovery fails closed unless the
// authoritative retry row is already terminal-quarantined for the exact binding
// and delivery position.
func RecoverRetryQuarantineCheckpoint(
	ctx context.Context,
	retryStore RetryCheckpointStateStore,
	checkpointStore CheckpointStore,
	binding DurableBinding,
	identity InboxIdentity,
	position uint64,
) (RetryQuarantineCheckpointResult, error) {
	if err := retryCheckpointContextError(ctx); err != nil {
		return RetryQuarantineCheckpointResult{}, err
	}
	if retryStore == nil || checkpointStore == nil || position == 0 {
		return RetryQuarantineCheckpointResult{}, classifiedFailure(codeRetryCheckpointInvalid, failure.CategoryValidation, "event retry checkpoint recovery boundary is invalid")
	}
	if err := identity.validate(); err != nil {
		return RetryQuarantineCheckpointResult{}, err
	}

	record, err := retryStore.Load(ctx, identity, position)
	if err != nil {
		return RetryQuarantineCheckpointResult{}, wrappedFailure(err, codeRetryCheckpointState, failure.CategoryUnavailable, "event retry quarantine state could not be loaded for checkpoint recovery")
	}
	if err := validateRetryCheckpointBinding(binding, record); err != nil {
		return RetryQuarantineCheckpointResult{}, err
	}
	if record.State != RetryStateQuarantined || record.ClaimToken != "" || !record.ClaimExpiresAt.IsZero() {
		return RetryQuarantineCheckpointResult{}, classifiedFailure(codeRetryCheckpointConflict, failure.CategoryConflict, "event retry checkpoint recovery requires an already-committed quarantine")
	}

	result := RetryQuarantineCheckpointResult{
		QuarantineCommitted: true,
		CheckpointStatus:    RetryCheckpointPending,
	}
	status, checkpointErr := advanceRetryQuarantineCheckpoint(ctx, checkpointStore, binding, record)
	result.CheckpointStatus = status
	if checkpointErr != nil {
		return result, checkpointErr
	}
	return result, nil
}

func advanceRetryQuarantineCheckpoint(
	ctx context.Context,
	store CheckpointStore,
	binding DurableBinding,
	record RetryStateRecord,
) (RetryCheckpointStatus, error) {
	checkpoint, exists, err := store.Load(ctx, binding)
	if err != nil {
		return RetryCheckpointPending, wrappedFailure(err, codeRetryCheckpointRead, failure.CategoryUnavailable, "event retry checkpoint could not be read")
	}

	expectedPosition := uint64(0)
	if exists {
		if !checkpoint.valid() {
			return RetryCheckpointPending, classifiedFailure(codeRetryCheckpointMalformed, failure.CategoryInvariant, "event retry checkpoint store returned malformed state")
		}
		if checkpoint.Position == record.Position {
			if checkpoint.EventID == record.Identity.EventID {
				return RetryCheckpointAlreadyAdvanced, nil
			}
			return RetryCheckpointPending, classifiedFailure(codeRetryCheckpointConflict, failure.CategoryConflict, "event retry checkpoint position is owned by another event")
		}
		if checkpoint.Position >= record.Position || checkpoint.Position+1 != record.Position {
			return RetryCheckpointPending, classifiedFailure(codeRetryCheckpointConflict, failure.CategoryConflict, "event retry quarantine checkpoint is not contiguous with accepted progress")
		}
		expectedPosition = checkpoint.Position
	} else if record.Position != 1 {
		return RetryCheckpointPending, classifiedFailure(codeRetryCheckpointConflict, failure.CategoryConflict, "event retry quarantine checkpoint cannot skip unacknowledged work")
	}

	next := Checkpoint{Position: record.Position, EventID: record.Identity.EventID}
	advance, err := store.Advance(ctx, binding, expectedPosition, next)
	if err != nil {
		return RetryCheckpointPending, wrappedFailure(err, codeRetryCheckpointWrite, failure.CategoryUnavailable, "event retry quarantine checkpoint could not be advanced")
	}
	switch advance {
	case CheckpointAdvanced:
		return RetryCheckpointAdvanced, nil
	case CheckpointStale, CheckpointConflict:
		// Another writer may have won the exact same checkpoint race. Re-read and
		// accept only byte-for-byte equivalent progress; every other result remains
		// a deterministic conflict.
		current, currentExists, loadErr := store.Load(ctx, binding)
		if loadErr != nil {
			return RetryCheckpointPending, wrappedFailure(loadErr, codeRetryCheckpointRead, failure.CategoryUnavailable, "event retry quarantine checkpoint could not be re-read after a race")
		}
		if currentExists && current == next {
			return RetryCheckpointAlreadyAdvanced, nil
		}
		return RetryCheckpointPending, classifiedFailure(codeRetryCheckpointConflict, failure.CategoryConflict, "event retry quarantine checkpoint conflicts with accepted progress")
	default:
		return RetryCheckpointPending, classifiedFailure(codeRetryCheckpointMalformed, failure.CategoryInvariant, "event retry checkpoint store returned an invalid advancement result")
	}
}

func validateRetryCheckpointBinding(binding DurableBinding, record RetryStateRecord) error {
	if err := binding.validateBasic(); err != nil {
		return classifiedFailure(codeRetryCheckpointInvalid, failure.CategoryValidation, "event retry checkpoint durable binding is invalid")
	}
	if err := record.Validate(); err != nil {
		return err
	}
	identity := record.Identity
	if identity.Owner != binding.Owner ||
		identity.ConsumerID != binding.ConsumerID ||
		identity.EventType != binding.EventType ||
		identity.Stream != binding.Scope.Stream ||
		identity.Partition != binding.Scope.Partition ||
		identity.TenantID != binding.Scope.TenantID {
		return classifiedFailure(codeRetryCheckpointConflict, failure.CategoryConflict, "event retry checkpoint state conflicts with the durable binding")
	}
	return nil
}

func retryCheckpointContextError(ctx context.Context) error {
	if ctx == nil {
		return classifiedFailure(codeRetryCheckpointInvalid, failure.CategoryValidation, "event retry checkpoint context is invalid")
	}
	if err := ctx.Err(); err != nil {
		return wrappedFailure(err, codeRetryCheckpointWrite, failure.CategoryUnavailable, "event retry checkpoint composition was interrupted")
	}
	return nil
}

package events

import (
	"context"
	"crypto/sha256"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

const (
	codeInboxInvalid        failure.Code = "events.inbox.invalid"
	codeInboxConflict       failure.Code = "events.inbox.conflict"
	codeInboxConcurrent     failure.Code = "events.inbox.concurrent"
	codeInboxStoreFailed    failure.Code = "events.inbox.store_failed"
	codeInboxStateMalformed failure.Code = "events.inbox.state_malformed"
	codeInboxInterrupted    failure.Code = "events.inbox.interrupted"
)

// InboxFingerprint is a fixed-size digest of the canonical P04.01 envelope.
// It is conflict-detection evidence only: it is not an event identifier,
// authorization credential, payload archive, or schema-registry substitute.
type InboxFingerprint [sha256.Size]byte

// InboxIdentity is the stable local processing identity for one protected
// consumer application. EventID is intentionally scoped by the accepted
// consumer owner/route/durable scope so the same canonical event may be applied
// independently by different legitimate consumers without becoming a global
// cross-consumer lock.
type InboxIdentity struct {
	EventID    EventID
	Owner      Producer
	ConsumerID string
	EventType  EventType
	Stream     string
	Partition  string
	TenantID   tenancy.TenantID
}

func (identity InboxIdentity) validate() error {
	binding := DurableBinding{
		Owner:      identity.Owner,
		ConsumerID: identity.ConsumerID,
		EventType:  identity.EventType,
		Scope: DurableScope{
			Stream:    identity.Stream,
			Partition: identity.Partition,
			TenantID:  identity.TenantID,
		},
	}
	if !identity.EventID.Valid() || binding.validateBasic() != nil {
		return classifiedFailure(codeInboxInvalid, failure.CategoryValidation, "event inbox processing identity is invalid")
	}
	return nil
}

// InboxRecord contains only the bounded identity/fingerprint evidence required
// to classify later redelivery. It deliberately does not retain the event
// payload or checkpoint position.
type InboxRecord struct {
	Identity    InboxIdentity
	Fingerprint InboxFingerprint
}

func (record InboxRecord) validate() error {
	if err := record.Identity.validate(); err != nil {
		return err
	}
	return nil
}

// InboxClaimResult separates a first transactional processing claim from an
// already-completed duplicate, conflicting canonical content, and a concurrent
// same-identity attempt whose authoritative state is not yet safely resolved.
type InboxClaimResult uint8

const (
	InboxClaimResultUnknown InboxClaimResult = iota
	InboxClaimed
	InboxAlreadyCompleted
	InboxIdentityConflict
	InboxConcurrentProcessing
)

// InboxApplyResult is the public core outcome for one attempted protected local
// application. Duplicate is a deterministic success-like outcome: the protected
// mutation is not called again. Conflict and concurrent processing fail closed.
type InboxApplyResult uint8

const (
	InboxApplyResultUnknown InboxApplyResult = iota
	InboxApplied
	InboxAlreadyApplied
	InboxConflict
	InboxConcurrent
)

// InboxStore is the provider-neutral P04.05 persistence seam. Claim and Complete
// receive the exact caller-provided transaction that also performs the protected
// local mutation. The store must not begin, commit, or hide a second transaction.
//
// A Claim returning InboxClaimed may create only transactional/incomplete state;
// that state must roll back with the caller transaction if the protected
// mutation or completion fails. Complete is called only after mutation success.
type InboxStore interface {
	Claim(context.Context, OutboxTransaction, InboxRecord) (InboxClaimResult, error)
	Complete(context.Context, OutboxTransaction, InboxRecord) error
}

// ProtectedMutation performs the already-authorized local mutation through the
// exact same transaction supplied to ApplyInbox. Receipt/dedup identity grants no
// authorization; callers remain responsible for their existing action boundary.
type ProtectedMutation func(context.Context, OutboxTransaction) error

// NewInboxRecord validates an accepted durable consumer binding against the
// canonical envelope and derives bounded conflict-detection evidence. Consumer
// owner identity is intentionally distinct from envelope.Source (the producer).
func NewInboxRecord(binding DurableBinding, envelope Envelope) (InboxRecord, error) {
	if err := binding.validateBasic(); err != nil {
		return InboxRecord{}, classifiedFailure(codeInboxInvalid, failure.CategoryValidation, "event inbox durable binding is invalid")
	}
	if binding.Scope.TenantID != "" {
		if err := envelope.ValidateForTenant(binding.Scope.TenantID); err != nil {
			return InboxRecord{}, err
		}
	} else if err := envelope.Validate(); err != nil {
		return InboxRecord{}, err
	}
	if envelope.Type != binding.EventType {
		return InboxRecord{}, classifiedFailure(codeInboxConflict, failure.CategoryConflict, "event inbox route conflicts with the canonical envelope")
	}

	serialized, err := envelope.Marshal()
	if err != nil {
		return InboxRecord{}, err
	}
	record := InboxRecord{
		Identity: InboxIdentity{
			EventID:    envelope.ID,
			Owner:      binding.Owner,
			ConsumerID: binding.ConsumerID,
			EventType:  binding.EventType,
			Stream:     binding.Scope.Stream,
			Partition:  binding.Scope.Partition,
			TenantID:   binding.Scope.TenantID,
		},
		Fingerprint: sha256.Sum256(serialized),
	}
	if err := record.validate(); err != nil {
		return InboxRecord{}, err
	}
	return record, nil
}

// ApplyInbox coordinates one duplicate-safe local application inside the exact
// caller transaction. It does not open or commit a transaction, advance a P04.03
// checkpoint, schedule retries, select a provider, or claim exactly-once external
// side effects. Callers compose this function inside the retained P01
// database.InTransaction boundary when PostgreSQL atomicity is required.
func ApplyInbox(
	ctx context.Context,
	tx OutboxTransaction,
	store InboxStore,
	binding DurableBinding,
	envelope Envelope,
	mutation ProtectedMutation,
) (InboxApplyResult, error) {
	if err := inboxContextError(ctx); err != nil {
		return InboxApplyResultUnknown, err
	}
	if tx == nil || store == nil || mutation == nil {
		return InboxApplyResultUnknown, classifiedFailure(codeInboxInvalid, failure.CategoryValidation, "event inbox application boundary is invalid")
	}

	record, err := NewInboxRecord(binding, envelope)
	if err != nil {
		return InboxApplyResultUnknown, err
	}
	claim, claimErr := store.Claim(ctx, tx, record)
	if claimErr != nil {
		return InboxApplyResultUnknown, wrappedFailure(claimErr, codeInboxStoreFailed, failure.CategoryUnavailable, "event inbox processing identity could not be claimed")
	}

	switch claim {
	case InboxAlreadyCompleted:
		return InboxAlreadyApplied, nil
	case InboxIdentityConflict:
		return InboxConflict, classifiedFailure(codeInboxConflict, failure.CategoryConflict, "event inbox processing identity conflicts with committed canonical evidence")
	case InboxConcurrentProcessing:
		return InboxConcurrent, classifiedFailure(codeInboxConcurrent, failure.CategoryConflict, "event inbox processing identity is being resolved concurrently")
	case InboxClaimed:
		// Continue below inside the exact caller transaction.
	default:
		return InboxApplyResultUnknown, classifiedFailure(codeInboxStateMalformed, failure.CategoryInvariant, "event inbox store returned an invalid claim result")
	}

	if mutationErr := mutation(ctx, tx); mutationErr != nil {
		// Preserve the already-authoritative mutation failure identity. The caller
		// transaction must roll back both mutation effects and any incomplete claim.
		return InboxApplyResultUnknown, mutationErr
	}
	if err := inboxContextError(ctx); err != nil {
		return InboxApplyResultUnknown, err
	}
	if completeErr := store.Complete(ctx, tx, record); completeErr != nil {
		return InboxApplyResultUnknown, wrappedFailure(completeErr, codeInboxStoreFailed, failure.CategoryUnavailable, "event inbox completion could not be recorded")
	}
	return InboxApplied, nil
}

func inboxContextError(ctx context.Context) error {
	if ctx == nil {
		return classifiedFailure(codeInboxInvalid, failure.CategoryValidation, "event inbox context is invalid")
	}
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return wrappedFailure(err, codeInboxInterrupted, failure.CategoryUnavailable, "event inbox application was interrupted")
		}
		return wrappedFailure(err, codeInboxInterrupted, failure.CategoryUnavailable, "event inbox application was interrupted")
	}
	return nil
}

package events

import (
	"context"
	"errors"
	"math"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	maxOutboxRelayBatch = uint32(256)

	codeOutboxInvalid        failure.Code = "events.outbox.invalid"
	codeOutboxConflict       failure.Code = "events.outbox.conflict"
	codeOutboxStoreFailed    failure.Code = "events.outbox.store_failed"
	codeOutboxStateMalformed failure.Code = "events.outbox.state_malformed"
	codeOutboxPublishFailed  failure.Code = "events.outbox.publish_failed"
	codeOutboxMarkFailed     failure.Code = "events.outbox.mark_failed"
	codeOutboxInterrupted    failure.Code = "events.outbox.interrupted"
)

// OutboxScope binds one producer-owned outbox record to the canonical producer
// and trusted tenant identity carried by the P04.01 envelope. It grants no
// authorization and contains no caller-selectable database namespace.
type OutboxScope struct {
	Owner    Producer
	TenantID tenancy.TenantID
}

func (scope OutboxScope) validate() error {
	if !scope.Owner.Valid() || (scope.TenantID != "" && !scope.TenantID.Valid()) {
		return classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox scope is invalid")
	}
	return nil
}

func (scope OutboxScope) validateEnvelope(envelope Envelope) error {
	if err := scope.validate(); err != nil {
		return err
	}
	if err := envelope.Validate(); err != nil {
		return err
	}
	if envelope.Source != scope.Owner || envelope.TenantID != scope.TenantID {
		return classifiedFailure(codeOutboxConflict, failure.CategoryConflict, "event outbox identity conflicts with its owner or tenant scope")
	}
	return nil
}

// OutboxState is producer-side publication progress only. Published never means
// consumer receipt, inbox deduplication or exactly-once protected mutation.
type OutboxState uint8

const (
	OutboxStateUnknown OutboxState = iota
	OutboxPending
	OutboxPublished
)

func (state OutboxState) valid() bool {
	return state == OutboxPending || state == OutboxPublished
}

// OutboxRecord retains one canonical P04.01 envelope unchanged with only the
// minimum producer-side state required for optimistic publication advancement.
// Revision is storage-owned compare-and-set state and is not event identity.
type OutboxRecord struct {
	Scope    OutboxScope
	Envelope Envelope
	State    OutboxState
	Revision uint64
}

// NewOutboxRecord validates and defensively copies a canonical envelope into the
// initial pending producer-side state. The envelope ID is never replaced.
func NewOutboxRecord(scope OutboxScope, envelope Envelope) (OutboxRecord, error) {
	if err := scope.validateEnvelope(envelope); err != nil {
		return OutboxRecord{}, err
	}
	cloned, err := cloneOutboxEnvelope(envelope)
	if err != nil {
		return OutboxRecord{}, err
	}
	return OutboxRecord{
		Scope:    scope,
		Envelope: cloned,
		State:    OutboxPending,
		Revision: 1,
	}, nil
}

func (record OutboxRecord) validate() error {
	if !record.State.valid() || record.Revision == 0 {
		return classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "stored event outbox state is malformed")
	}
	if err := record.Scope.validateEnvelope(record.Envelope); err != nil {
		return classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "stored event outbox identity is malformed")
	}
	return nil
}

// OutboxTransaction is the minimum SQL capability that a pgx.Tx already
// satisfies. Callers must pass the same transaction object that performs the
// authoritative owner mutation; the outbox core never begins a second
// transaction. T02 may use these operations in its owner-scoped PostgreSQL
// adapter without exposing a separate transaction framework.
type OutboxTransaction interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// EnqueueResult makes exact same-event retries distinguishable from conflicting
// EventID reuse. Conflict is producer-side identity integrity, not consumer
// deduplication.
type EnqueueResult uint8

const (
	EnqueueResultUnknown EnqueueResult = iota
	OutboxEnqueued
	OutboxUnchanged
	OutboxEnqueueConflict
)

// MarkResult is an optimistic publication-state transition result. A concurrent
// relay may have already marked the same record after both relays published it;
// this is an accepted duplicate-publication window rather than exactly-once.
type MarkResult uint8

const (
	MarkResultUnknown MarkResult = iota
	OutboxMarkedPublished
	OutboxAlreadyPublished
	OutboxMarkConflict
)

// OutboxStore is the provider-neutral P04.04 persistence seam. Enqueue receives
// the caller's existing local PostgreSQL transaction. Pending and MarkPublished
// operate only on committed producer-side state. Retry schedules, DLQ policy and
// background-job ownership are intentionally absent.
type OutboxStore interface {
	Enqueue(context.Context, OutboxTransaction, OutboxRecord) (EnqueueResult, error)
	Pending(context.Context, OutboxScope, uint32) ([]OutboxRecord, error)
	MarkPublished(context.Context, OutboxScope, EventID, uint64) (MarkResult, error)
}

// EnqueueOutbox writes the canonical event through the exact transaction handle
// supplied by the authorized owner mutation. It never opens or commits a second
// transaction itself.
func EnqueueOutbox(ctx context.Context, tx OutboxTransaction, store OutboxStore, scope OutboxScope, envelope Envelope) (EnqueueResult, error) {
	if err := outboxContextError(ctx); err != nil {
		return EnqueueResultUnknown, err
	}
	if tx == nil || store == nil {
		return EnqueueResultUnknown, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox enqueue boundary is invalid")
	}
	record, err := NewOutboxRecord(scope, envelope)
	if err != nil {
		return EnqueueResultUnknown, err
	}
	result, storeErr := store.Enqueue(ctx, tx, record)
	if storeErr != nil {
		return EnqueueResultUnknown, wrappedFailure(storeErr, codeOutboxStoreFailed, failure.CategoryUnavailable, "event outbox enqueue failed")
	}
	switch result {
	case OutboxEnqueued, OutboxUnchanged:
		return result, nil
	case OutboxEnqueueConflict:
		return result, classifiedFailure(codeOutboxConflict, failure.CategoryConflict, "event outbox identity conflicts with committed state")
	default:
		return EnqueueResultUnknown, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "event outbox store returned an invalid enqueue result")
	}
}

// OutboxRelay composes committed outbox state with the accepted P04.02 Publisher.
// It is an explicit bounded relay operation, not a scheduler, retry engine or
// provider worker.
type OutboxRelay struct {
	store     OutboxStore
	publisher *Publisher
}

func NewOutboxRelay(store OutboxStore, publisher *Publisher) (*OutboxRelay, error) {
	if store == nil || publisher == nil || publisher.transport == nil {
		return nil, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox relay configuration is invalid")
	}
	return &OutboxRelay{store: store, publisher: publisher}, nil
}

// RelayReport reports producer-side publish attempts only. AcceptedPublications
// means P04.02 transport acceptance; it is not consumer or business completion.
type RelayReport struct {
	PendingObserved       uint32
	AcceptedPublications  uint32
	PublishedMarks        uint32
	DuplicateMarkObserved uint32
}

// RelayPending publishes at most limit committed pending records for one exact
// owner/tenant scope. A publish failure or mark failure stops this explicit pass;
// the record remains recoverable in storage. If publish succeeds and marking then
// fails, a later pass may publish the same EventID again by design.
func (relay *OutboxRelay) RelayPending(ctx context.Context, scope OutboxScope, limit uint32) (RelayReport, error) {
	if relay == nil || relay.store == nil || relay.publisher == nil {
		return RelayReport{}, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox relay is invalid")
	}
	if err := outboxContextError(ctx); err != nil {
		return RelayReport{}, err
	}
	if err := scope.validate(); err != nil {
		return RelayReport{}, err
	}
	if limit == 0 || limit > maxOutboxRelayBatch {
		return RelayReport{}, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox relay batch is invalid")
	}

	records, err := relay.store.Pending(ctx, scope, limit)
	if err != nil {
		return RelayReport{}, wrappedFailure(err, codeOutboxStoreFailed, failure.CategoryUnavailable, "event outbox pending state could not be read")
	}
	if uint64(len(records)) > uint64(limit) || uint64(len(records)) > uint64(math.MaxUint32) {
		return RelayReport{}, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "event outbox store returned an invalid pending batch")
	}

	report := RelayReport{PendingObserved: uint32(len(records))}
	for _, record := range records {
		if err := outboxContextError(ctx); err != nil {
			return report, err
		}
		if err := record.validate(); err != nil || record.Scope != scope || record.State != OutboxPending {
			return report, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "event outbox store returned malformed or cross-scope pending state")
		}

		publishResult, publishErr := relay.publisher.Publish(ctx, record.Envelope)
		if publishErr != nil {
			return report, wrappedFailure(publishErr, codeOutboxPublishFailed, failure.CategoryUnavailable, "event outbox publication failed")
		}
		if !publishResult.Accepted {
			return report, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "event publisher returned invalid acceptance state")
		}
		report.AcceptedPublications++

		markResult, markErr := relay.store.MarkPublished(ctx, scope, record.Envelope.ID, record.Revision)
		if markErr != nil {
			return report, wrappedFailure(markErr, codeOutboxMarkFailed, failure.CategoryUnavailable, "event outbox published state could not be recorded")
		}
		switch markResult {
		case OutboxMarkedPublished:
			report.PublishedMarks++
		case OutboxAlreadyPublished:
			report.DuplicateMarkObserved++
		case OutboxMarkConflict:
			return report, classifiedFailure(codeOutboxConflict, failure.CategoryConflict, "event outbox publication state changed concurrently")
		default:
			return report, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "event outbox store returned an invalid publication result")
		}
	}
	return report, nil
}

func cloneOutboxEnvelope(envelope Envelope) (Envelope, error) {
	serialized, err := envelope.Marshal()
	if err != nil {
		return Envelope{}, err
	}
	cloned, err := Parse(serialized)
	if err != nil {
		return Envelope{}, err
	}
	return cloned, nil
}

func outboxContextError(ctx context.Context) error {
	if ctx == nil {
		return classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox context is invalid")
	}
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return wrappedFailure(err, codeOutboxInterrupted, failure.CategoryUnavailable, "event outbox operation was interrupted")
		}
		return wrappedFailure(err, codeOutboxInterrupted, failure.CategoryUnavailable, "event outbox operation was interrupted")
	}
	return nil
}

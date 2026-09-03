package events

import (
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxPostgresOutboxRevision uint64 = 1<<63 - 1

// PostgresOutboxStore persists P04.04 producer-side outbox state in the
// kernel.events owner schema. Enqueue uses the caller-provided transaction;
// committed-state reads and publication marks use the configured pool.
type PostgresOutboxStore struct {
	pool *pgxpool.Pool
}

// NewPostgresOutboxStore binds the outbox adapter to an existing PostgreSQL
// pool. It does not create a second migration or transaction authority.
func NewPostgresOutboxStore(pool *pgxpool.Pool) (*PostgresOutboxStore, error) {
	if pool == nil {
		return nil, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox postgres store is invalid")
	}
	return &PostgresOutboxStore{pool: pool}, nil
}

// Enqueue stores one canonical record through the exact transaction supplied by
// the owner mutation. Exact same-event replay is unchanged; conflicting EventID
// reuse fails closed through the core EnqueueResult contract.
func (store *PostgresOutboxStore) Enqueue(ctx context.Context, tx OutboxTransaction, record OutboxRecord) (EnqueueResult, error) {
	if store == nil || store.pool == nil || tx == nil {
		return EnqueueResultUnknown, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox postgres enqueue boundary is invalid")
	}
	if err := outboxContextError(ctx); err != nil {
		return EnqueueResultUnknown, err
	}
	if err := record.validate(); err != nil || record.State != OutboxPending || record.Revision > maxPostgresOutboxRevision {
		return EnqueueResultUnknown, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "event outbox postgres enqueue state is malformed")
	}

	serialized, err := record.Envelope.Marshal()
	if err != nil {
		return EnqueueResultUnknown, err
	}
	tenant := postgresOutboxTenant(record.Scope)
	tag, err := tx.Exec(
		ctx,
		`INSERT INTO omnexa_events.transactional_outbox
			(event_id, owner, tenant_id, envelope, publication_state, revision)
		 VALUES ($1::uuid, $2, $3::uuid, $4::jsonb, 'pending', $5)
		 ON CONFLICT (event_id) DO NOTHING`,
		string(record.Envelope.ID),
		string(record.Scope.Owner),
		tenant,
		string(serialized),
		int64(record.Revision),
	)
	if err != nil {
		return EnqueueResultUnknown, err
	}
	if tag.RowsAffected() == 1 {
		return OutboxEnqueued, nil
	}
	if tag.RowsAffected() != 0 {
		return EnqueueResultUnknown, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "event outbox postgres enqueue returned an invalid write count")
	}

	var unchanged bool
	err = tx.QueryRow(
		ctx,
		`SELECT owner = $2
			AND tenant_id IS NOT DISTINCT FROM $3::uuid
			AND envelope = $4::jsonb
		 FROM omnexa_events.transactional_outbox
		 WHERE event_id = $1::uuid`,
		string(record.Envelope.ID),
		string(record.Scope.Owner),
		tenant,
		string(serialized),
	).Scan(&unchanged)
	if errors.Is(err, pgx.ErrNoRows) {
		return EnqueueResultUnknown, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "event outbox conflicting state disappeared during enqueue")
	}
	if err != nil {
		return EnqueueResultUnknown, err
	}
	if unchanged {
		return OutboxUnchanged, nil
	}
	return OutboxEnqueueConflict, nil
}

// Pending reads only committed pending records for one exact producer/tenant
// scope. It does not lock, claim, schedule or infer global ordering.
func (store *PostgresOutboxStore) Pending(ctx context.Context, scope OutboxScope, limit uint32) ([]OutboxRecord, error) {
	if store == nil || store.pool == nil {
		return nil, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox postgres store is invalid")
	}
	if err := outboxContextError(ctx); err != nil {
		return nil, err
	}
	if err := scope.validate(); err != nil {
		return nil, err
	}
	if limit == 0 || limit > maxOutboxRelayBatch {
		return nil, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox postgres pending batch is invalid")
	}

	rows, err := store.pool.Query(
		ctx,
		`SELECT envelope, publication_state, revision
		 FROM omnexa_events.transactional_outbox
		 WHERE owner = $1
		   AND tenant_id IS NOT DISTINCT FROM $2::uuid
		   AND publication_state = 'pending'
		 ORDER BY created_at, event_id
		 LIMIT $3`,
		string(scope.Owner),
		postgresOutboxTenant(scope),
		int64(limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]OutboxRecord, 0, limit)
	for rows.Next() {
		var serialized []byte
		var persistedState string
		var revision int64
		if err = rows.Scan(&serialized, &persistedState, &revision); err != nil {
			return nil, err
		}
		if persistedState != "pending" || revision <= 0 {
			return nil, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "stored event outbox postgres state is malformed")
		}
		envelope, parseErr := Parse(serialized)
		if parseErr != nil {
			return nil, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "stored event outbox postgres envelope is malformed")
		}
		record := OutboxRecord{
			Scope:    scope,
			Envelope: envelope,
			State:    OutboxPending,
			Revision: uint64(revision),
		}
		if validateErr := record.validate(); validateErr != nil {
			return nil, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "stored event outbox postgres identity is malformed")
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

// MarkPublished advances one exact pending record with compare-and-set revision
// semantics. Concurrent relays may already have published the same event; that
// is reported as OutboxAlreadyPublished rather than rewritten as exactly-once.
func (store *PostgresOutboxStore) MarkPublished(ctx context.Context, scope OutboxScope, eventID EventID, revision uint64) (MarkResult, error) {
	if store == nil || store.pool == nil {
		return MarkResultUnknown, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox postgres store is invalid")
	}
	if err := outboxContextError(ctx); err != nil {
		return MarkResultUnknown, err
	}
	if err := scope.validate(); err != nil || !eventID.Valid() || revision == 0 || revision >= maxPostgresOutboxRevision {
		return MarkResultUnknown, classifiedFailure(codeOutboxInvalid, failure.CategoryValidation, "event outbox postgres mark boundary is invalid")
	}

	tag, err := store.pool.Exec(
		ctx,
		`UPDATE omnexa_events.transactional_outbox
		 SET publication_state = 'published',
		     revision = revision + 1,
		     published_at = CURRENT_TIMESTAMP
		 WHERE event_id = $1::uuid
		   AND owner = $2
		   AND tenant_id IS NOT DISTINCT FROM $3::uuid
		   AND publication_state = 'pending'
		   AND revision = $4`,
		string(eventID),
		string(scope.Owner),
		postgresOutboxTenant(scope),
		int64(revision),
	)
	if err != nil {
		return MarkResultUnknown, err
	}
	if tag.RowsAffected() == 1 {
		return OutboxMarkedPublished, nil
	}
	if tag.RowsAffected() != 0 {
		return MarkResultUnknown, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "event outbox postgres mark returned an invalid write count")
	}

	var state string
	var storedRevision int64
	err = store.pool.QueryRow(
		ctx,
		`SELECT publication_state, revision
		 FROM omnexa_events.transactional_outbox
		 WHERE event_id = $1::uuid
		   AND owner = $2
		   AND tenant_id IS NOT DISTINCT FROM $3::uuid`,
		string(eventID),
		string(scope.Owner),
		postgresOutboxTenant(scope),
	).Scan(&state, &storedRevision)
	if errors.Is(err, pgx.ErrNoRows) {
		return OutboxMarkConflict, nil
	}
	if err != nil {
		return MarkResultUnknown, err
	}
	if storedRevision <= 0 {
		return MarkResultUnknown, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "stored event outbox postgres publication state is malformed")
	}
	switch state {
	case "published":
		return OutboxAlreadyPublished, nil
	case "pending":
		return OutboxMarkConflict, nil
	default:
		return MarkResultUnknown, classifiedFailure(codeOutboxStateMalformed, failure.CategoryInvariant, "stored event outbox postgres publication state is malformed")
	}
}

func postgresOutboxTenant(scope OutboxScope) any {
	if scope.TenantID == "" {
		return nil
	}
	return string(scope.TenantID)
}

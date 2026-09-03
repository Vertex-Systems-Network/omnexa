package events

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresOutboxStoreIntegration(t *testing.T) {
	databaseURL := os.Getenv("P04_04_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P04_04_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		t.Fatalf("pool.Ping() error = %v", err)
	}

	foundation, err := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if err != nil {
		t.Fatalf("foundation migrator create error = %v", err)
	}
	if err = foundation.Run(ctx); err != nil {
		t.Fatalf("foundation migrator run error = %v", err)
	}
	defer resetP0404OutboxDatabase(t, context.Background(), pool)
	resetP0404OutboxDatabase(t, ctx, pool)

	migrationSQL, err := os.ReadFile("../../migrations/kernel.events/1_create_transactional_outbox.sql")
	if err != nil {
		t.Fatalf("read P04.04 migration error = %v", err)
	}
	migrator, err := database.NewMigrator(pool, "kernel.events", []database.Migration{{
		Version: 1,
		Name:    "create_transactional_outbox",
		SQL:     string(migrationSQL),
	}}, 5*time.Second)
	if err != nil {
		t.Fatalf("kernel.events migrator create error = %v", err)
	}
	if err = migrator.Run(ctx); err != nil {
		t.Fatalf("kernel.events migrator run error = %v", err)
	}
	if err = migrator.Run(ctx); err != nil {
		t.Fatalf("kernel.events migration replay error = %v", err)
	}

	if _, err = pool.Exec(ctx, `CREATE TABLE omnexa_events.p04_04_test_owner_state (id uuid PRIMARY KEY, note text NOT NULL)`); err != nil {
		t.Fatalf("create synthetic owner table error = %v", err)
	}
	store, err := NewPostgresOutboxStore(pool)
	if err != nil {
		t.Fatalf("NewPostgresOutboxStore() error = %v", err)
	}

	// This 19-character producer is the shortest canonical shape that exposed a
	// migration-only contract mismatch during Supervisor review.
	scope := OutboxScope{
		Owner:    Producer("urn:omnexa:module:x"),
		TenantID: tenancy.TenantID(mustP0404UUIDv7(t)),
	}
	envelope := mustP0404Envelope(t, scope, "atomic-commit")
	ownerMutationID := mustP0404UUIDv7(t)

	err = database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if _, execErr := tx.Exec(ctx, `INSERT INTO omnexa_events.p04_04_test_owner_state (id, note) VALUES ($1::uuid, 'committed')`, ownerMutationID); execErr != nil {
			return execErr
		}
		result, enqueueErr := EnqueueOutbox(ctx, tx, store, scope, envelope)
		if enqueueErr != nil {
			return enqueueErr
		}
		if result != OutboxEnqueued {
			return errors.New("initial outbox enqueue did not report enqueued")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("atomic owner mutation + enqueue error = %v", err)
	}
	assertP0404Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.p04_04_test_owner_state WHERE id=$1::uuid`, ownerMutationID, 1)
	assertP0404Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.transactional_outbox WHERE event_id=$1::uuid`, string(envelope.ID), 1)

	pending, err := store.Pending(ctx, scope, 8)
	if err != nil {
		t.Fatalf("Pending() after commit error = %v", err)
	}
	if len(pending) != 1 || pending[0].Envelope.ID != envelope.ID || pending[0].Revision != 1 {
		t.Fatalf("Pending() after commit = %#v", pending)
	}
	serializedOriginal, err := envelope.Marshal()
	if err != nil {
		t.Fatalf("original envelope Marshal() error = %v", err)
	}
	serializedRecovered, err := pending[0].Envelope.Marshal()
	if err != nil {
		t.Fatalf("recovered envelope Marshal() error = %v", err)
	}
	if string(serializedRecovered) != string(serializedOriginal) {
		t.Fatal("committed outbox envelope did not round-trip canonically")
	}

	if err = database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		result, enqueueErr := EnqueueOutbox(ctx, tx, store, scope, envelope)
		if enqueueErr != nil {
			return enqueueErr
		}
		if result != OutboxUnchanged {
			return errors.New("exact same-event retry did not report unchanged")
		}
		return nil
	}); err != nil {
		t.Fatalf("same-event replay error = %v", err)
	}

	conflicting := envelope
	conflicting.Data = json.RawMessage(`{"value":"conflicting-content"}`)
	if err = conflicting.Validate(); err == nil {
		t.Fatal("escaped invalid JSON fixture unexpectedly validated")
	}
	conflicting.Data = json.RawMessage(`{"value":"conflicting-content"}`[1:])
	if err = conflicting.Validate(); err == nil {
		t.Fatal("malformed conflict fixture unexpectedly validated")
	}
	conflicting.Data = json.RawMessage([]byte{'{', '"', 'v', 'a', 'l', 'u', 'e', '"', ':', '"', 'c', 'o', 'n', 'f', 'l', 'i', 'c', 't', 'i', 'n', 'g', '-', 'c', 'o', 'n', 't', 'e', 'n', 't', '"', '}'})
	if err = conflicting.Validate(); err != nil {
		t.Fatalf("conflicting test envelope validation error = %v", err)
	}
	rolledBackOwnerID := mustP0404UUIDv7(t)
	err = database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if _, execErr := tx.Exec(ctx, `INSERT INTO omnexa_events.p04_04_test_owner_state (id, note) VALUES ($1::uuid, 'must-rollback')`, rolledBackOwnerID); execErr != nil {
			return execErr
		}
		result, enqueueErr := EnqueueOutbox(ctx, tx, store, scope, conflicting)
		if result != OutboxEnqueueConflict || enqueueErr == nil {
			return errors.New("conflicting EventID reuse did not fail closed")
		}
		return enqueueErr
	})
	if err == nil {
		t.Fatal("conflicting enqueue transaction unexpectedly committed")
	}
	assertP0404Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.p04_04_test_owner_state WHERE id=$1::uuid`, rolledBackOwnerID, 0)

	uncommittedEnvelope := mustP0404Envelope(t, scope, "uncommitted")
	uncommittedOwnerID := mustP0404UUIDv7(t)
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("pool.Begin() error = %v", err)
	}
	if _, err = tx.Exec(ctx, `INSERT INTO omnexa_events.p04_04_test_owner_state (id, note) VALUES ($1::uuid, 'uncommitted')`, uncommittedOwnerID); err != nil {
		_ = tx.Rollback(ctx)
		t.Fatalf("uncommitted owner insert error = %v", err)
	}
	if result, enqueueErr := EnqueueOutbox(ctx, tx, store, scope, uncommittedEnvelope); enqueueErr != nil || result != OutboxEnqueued {
		_ = tx.Rollback(ctx)
		t.Fatalf("uncommitted EnqueueOutbox() = %v/%v", result, enqueueErr)
	}
	outsidePending, pendingErr := store.Pending(ctx, scope, 8)
	if pendingErr != nil {
		_ = tx.Rollback(ctx)
		t.Fatalf("Pending() while second transaction open error = %v", pendingErr)
	}
	if len(outsidePending) != 1 || outsidePending[0].Envelope.ID != envelope.ID {
		_ = tx.Rollback(ctx)
		t.Fatalf("relay-visible pending state included uncommitted event: %#v", outsidePending)
	}
	if err = tx.Rollback(ctx); err != nil {
		t.Fatalf("tx.Rollback() error = %v", err)
	}
	assertP0404Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.p04_04_test_owner_state WHERE id=$1::uuid`, uncommittedOwnerID, 0)
	assertP0404Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.transactional_outbox WHERE event_id=$1::uuid`, string(uncommittedEnvelope.ID), 0)

	wrongScope := OutboxScope{Owner: scope.Owner, TenantID: tenancy.TenantID(mustP0404UUIDv7(t))}
	wrongPending, err := store.Pending(ctx, wrongScope, 8)
	if err != nil {
		t.Fatalf("wrong-scope Pending() error = %v", err)
	}
	if len(wrongPending) != 0 {
		t.Fatalf("wrong-scope Pending() leaked records = %#v", wrongPending)
	}
	if result, markErr := store.MarkPublished(ctx, wrongScope, envelope.ID, 1); markErr != nil || result != OutboxMarkConflict {
		t.Fatalf("wrong-scope MarkPublished() = %v/%v", result, markErr)
	}

	markResult, err := store.MarkPublished(ctx, scope, envelope.ID, 1)
	if err != nil || markResult != OutboxMarkedPublished {
		t.Fatalf("MarkPublished() = %v/%v", markResult, err)
	}
	markResult, err = store.MarkPublished(ctx, scope, envelope.ID, 1)
	if err != nil || markResult != OutboxAlreadyPublished {
		t.Fatalf("duplicate MarkPublished() = %v/%v", markResult, err)
	}
	pending, err = store.Pending(ctx, scope, 8)
	if err != nil {
		t.Fatalf("Pending() after published mark error = %v", err)
	}
	if len(pending) != 0 {
		t.Fatalf("published event remained pending = %#v", pending)
	}

	assertP0404Count(t, ctx, pool, `SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner='kernel.events' AND version=1`, nil, 1)
}

func mustP0404Envelope(t *testing.T, scope OutboxScope, value string) Envelope {
	t.Helper()
	correlation := mustP0404UUIDv7(t)
	envelope, err := New(Params{
		Type:           EventType("kernel.events.outbox.v1"),
		Producer:       scope.Owner,
		OccurredAt:     time.Date(2026, time.September, 4, 0, 0, 0, 0, time.UTC),
		TenantID:       scope.TenantID,
		CorrelationID:  CorrelationID(correlation),
		Classification: DataClassInternal,
		DataSchema:     SchemaID("urn:omnexa:event-schema:kernel.events.outbox:v1"),
		Data:           map[string]any{"value": value},
	})
	if err != nil {
		t.Fatalf("New() envelope error = %v", err)
	}
	return envelope
}

func mustP0404UUIDv7(t *testing.T) string {
	t.Helper()
	identifier, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("uuid.NewV7() error = %v", err)
	}
	return identifier.String()
}

func assertP0404Count(t *testing.T, ctx context.Context, pool *pgxpool.Pool, query string, argument any, expected int) {
	t.Helper()
	var count int
	var err error
	if argument == nil {
		err = pool.QueryRow(ctx, query).Scan(&count)
	} else {
		err = pool.QueryRow(ctx, query, argument).Scan(&count)
	}
	if err != nil || count != expected {
		t.Fatalf("count query result/error = %d/%v; expected %d", count, err, expected)
	}
}

func resetP0404OutboxDatabase(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(ctx, `DROP SCHEMA IF EXISTS omnexa_events CASCADE`); err != nil {
		t.Fatalf("drop omnexa_events schema error = %v", err)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM omnexa_kernel.schema_migrations WHERE owner='kernel.events'`); err != nil {
		t.Fatalf("reset kernel.events migration ledger error = %v", err)
	}
}

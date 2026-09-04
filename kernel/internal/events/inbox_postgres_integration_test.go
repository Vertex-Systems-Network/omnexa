package events

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresInboxStoreIntegration(t *testing.T) {
	databaseURL := os.Getenv("P04_05_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P04_05_TEST_DATABASE_URL is not set")
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
	defer resetP0405InboxDatabase(t, context.Background(), pool)
	resetP0405InboxDatabase(t, ctx, pool)

	outboxSQL, err := os.ReadFile("../../migrations/kernel.events/1_create_transactional_outbox.sql")
	if err != nil {
		t.Fatalf("read P04.04 migration error = %v", err)
	}
	inboxSQL, err := os.ReadFile("../../migrations/kernel.events/2_create_consumer_inbox.sql")
	if err != nil {
		t.Fatalf("read P04.05 migration error = %v", err)
	}
	migrator, err := database.NewMigrator(pool, "kernel.events", []database.Migration{
		{Version: 1, Name: "create_transactional_outbox", SQL: string(outboxSQL)},
		{Version: 2, Name: "create_consumer_inbox", SQL: string(inboxSQL)},
	}, 5*time.Second)
	if err != nil {
		t.Fatalf("kernel.events migrator create error = %v", err)
	}
	if err = migrator.Run(ctx); err != nil {
		t.Fatalf("kernel.events migrator run error = %v", err)
	}
	if err = migrator.Run(ctx); err != nil {
		t.Fatalf("kernel.events migration replay error = %v", err)
	}

	if _, err = pool.Exec(ctx, `CREATE TABLE omnexa_events.p04_05_test_owner_state (
		event_id uuid NOT NULL,
		consumer_id text NOT NULL,
		note text NOT NULL,
		PRIMARY KEY (event_id, consumer_id)
	)`); err != nil {
		t.Fatalf("create synthetic owner table error = %v", err)
	}

	store := NewPostgresInboxStore()
	envelope := testEnvelope(t)
	binding := testInboxBinding(envelope, "billing.projection")

	result, err := applyP0405OwnerMutation(ctx, pool, store, binding, envelope, "committed")
	if err != nil || result != InboxApplied {
		t.Fatalf("first ApplyInbox() = %v/%v", result, err)
	}
	assertP0405Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.p04_05_test_owner_state WHERE event_id=$1::uuid AND consumer_id=$2`, []any{string(envelope.ID), binding.ConsumerID}, 1)
	assertP0405Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.consumer_inbox WHERE event_id=$1::uuid AND consumer_id=$2 AND processing_state='completed'`, []any{string(envelope.ID), binding.ConsumerID}, 1)

	// A new adapter instance represents restart of the process-local store object.
	restarted := NewPostgresInboxStore()
	duplicateMutationCalled := false
	err = database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		duplicate, applyErr := ApplyInbox(ctx, tx, restarted, binding, envelope, func(context.Context, OutboxTransaction) error {
			duplicateMutationCalled = true
			return errors.New("duplicate mutation must not execute")
		})
		if applyErr != nil {
			return applyErr
		}
		if duplicate != InboxAlreadyApplied {
			return errors.New("restart duplicate did not report already applied")
		}
		return nil
	})
	if err != nil || duplicateMutationCalled {
		t.Fatalf("restart duplicate result error/callback = %v/%v", err, duplicateMutationCalled)
	}

	// EventID is not a global lock: another accepted consumer scope may apply the
	// same canonical event independently.
	secondBinding := testInboxBinding(envelope, "analytics.projection")
	result, err = applyP0405OwnerMutation(ctx, pool, restarted, secondBinding, envelope, "second-consumer")
	if err != nil || result != InboxApplied {
		t.Fatalf("second consumer ApplyInbox() = %v/%v", result, err)
	}
	assertP0405Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.consumer_inbox WHERE event_id=$1::uuid AND processing_state='completed'`, []any{string(envelope.ID)}, 2)

	// Same processing identity with different canonical content is conflict, not
	// a harmless duplicate, and never reaches the protected mutation.
	conflicting := envelope
	conflicting.Subject = "same-identity-conflicting-content"
	if err = conflicting.Validate(); err != nil {
		t.Fatalf("conflicting envelope validation error = %v", err)
	}
	conflictMutationCalled := false
	err = database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		conflictResult, applyErr := ApplyInbox(ctx, tx, restarted, binding, conflicting, func(context.Context, OutboxTransaction) error {
			conflictMutationCalled = true
			return nil
		})
		if conflictResult != InboxConflict || applyErr == nil {
			return errors.New("conflicting canonical content did not fail closed")
		}
		return applyErr
	})
	if err == nil || conflictMutationCalled {
		t.Fatalf("conflict transaction/callback = %v/%v", err, conflictMutationCalled)
	}

	// A mutation failure rolls back both the synthetic owner row and the claimed
	// inbox identity, allowing a later legitimate redelivery to succeed.
	retryEnvelope := testEnvelope(t)
	retryBinding := testInboxBinding(retryEnvelope, "billing.retryable")
	sentinel := errors.New("synthetic protected mutation failed")
	err = database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, applyErr := ApplyInbox(ctx, tx, store, retryBinding, retryEnvelope, func(mutationCtx context.Context, mutationTx OutboxTransaction) error {
			if _, execErr := mutationTx.Exec(
				mutationCtx,
				`INSERT INTO omnexa_events.p04_05_test_owner_state (event_id, consumer_id, note) VALUES ($1::uuid, $2, 'must-rollback')`,
				string(retryEnvelope.ID),
				retryBinding.ConsumerID,
			); execErr != nil {
				return execErr
			}
			return sentinel
		})
		return applyErr
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("failed mutation transaction error = %v, want sentinel", err)
	}
	assertP0405Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.p04_05_test_owner_state WHERE event_id=$1::uuid`, []any{string(retryEnvelope.ID)}, 0)
	assertP0405Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.consumer_inbox WHERE event_id=$1::uuid AND consumer_id=$2`, []any{string(retryEnvelope.ID), retryBinding.ConsumerID}, 0)

	result, err = applyP0405OwnerMutation(ctx, pool, store, retryBinding, retryEnvelope, "retry-committed")
	if err != nil || result != InboxApplied {
		t.Fatalf("retry after rollback ApplyInbox() = %v/%v", result, err)
	}

	assertP0405Count(t, ctx, pool, `SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner='kernel.events' AND version=1`, nil, 1)
	assertP0405Count(t, ctx, pool, `SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner='kernel.events' AND version=2`, nil, 1)

	// Simulate inbox persistence unavailability after the owner mutation has been
	// attempted. ApplyInbox must return an error so InTransaction rolls it back.
	if _, err = pool.Exec(ctx, `DROP TABLE omnexa_events.consumer_inbox`); err != nil {
		t.Fatalf("drop inbox table for persistence failure error = %v", err)
	}
	failedPersistenceEnvelope := testEnvelope(t)
	failedPersistenceBinding := testInboxBinding(failedPersistenceEnvelope, "billing.persistence-failure")
	err = database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if _, execErr := tx.Exec(
			ctx,
			`INSERT INTO omnexa_events.p04_05_test_owner_state (event_id, consumer_id, note) VALUES ($1::uuid, $2, 'must-rollback')`,
			string(failedPersistenceEnvelope.ID),
			failedPersistenceBinding.ConsumerID,
		); execErr != nil {
			return execErr
		}
		_, applyErr := ApplyInbox(ctx, tx, store, failedPersistenceBinding, failedPersistenceEnvelope, func(context.Context, OutboxTransaction) error {
			return nil
		})
		return applyErr
	})
	if err == nil {
		t.Fatal("inbox persistence failure unexpectedly allowed transaction commit")
	}
	assertP0405Count(t, ctx, pool, `SELECT count(*) FROM omnexa_events.p04_05_test_owner_state WHERE event_id=$1::uuid`, []any{string(failedPersistenceEnvelope.ID)}, 0)
}

func applyP0405OwnerMutation(
	ctx context.Context,
	pool *pgxpool.Pool,
	store InboxStore,
	binding DurableBinding,
	envelope Envelope,
	note string,
) (InboxApplyResult, error) {
	result := InboxApplyResultUnknown
	err := database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		var applyErr error
		result, applyErr = ApplyInbox(ctx, tx, store, binding, envelope, func(mutationCtx context.Context, mutationTx OutboxTransaction) error {
			_, execErr := mutationTx.Exec(
				mutationCtx,
				`INSERT INTO omnexa_events.p04_05_test_owner_state (event_id, consumer_id, note) VALUES ($1::uuid, $2, $3)`,
				string(envelope.ID),
				binding.ConsumerID,
				note,
			)
			return execErr
		})
		return applyErr
	})
	return result, err
}

func assertP0405Count(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	query string,
	arguments []any,
	expected int,
) {
	t.Helper()
	var count int
	var err error
	if arguments == nil {
		err = pool.QueryRow(ctx, query).Scan(&count)
	} else {
		err = pool.QueryRow(ctx, query, arguments...).Scan(&count)
	}
	if err != nil || count != expected {
		t.Fatalf("count query result/error = %d/%v; expected %d", count, err, expected)
	}
}

func resetP0405InboxDatabase(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(ctx, `DROP SCHEMA IF EXISTS omnexa_events CASCADE`); err != nil {
		t.Fatalf("drop omnexa_events schema error = %v", err)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM omnexa_kernel.schema_migrations WHERE owner='kernel.events'`); err != nil {
		t.Fatalf("reset kernel.events migration ledger error = %v", err)
	}
}

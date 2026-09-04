package events

import (
	"context"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type p0405FailOnceCheckpointStore struct {
	delegate *MemoryCheckpointStore
	mu       sync.Mutex
	failNext bool
}

func (store *p0405FailOnceCheckpointStore) Bind(ctx context.Context, binding DurableBinding) (BindingResult, error) {
	return store.delegate.Bind(ctx, binding)
}

func (store *p0405FailOnceCheckpointStore) Load(ctx context.Context, binding DurableBinding) (Checkpoint, bool, error) {
	return store.delegate.Load(ctx, binding)
}

func (store *p0405FailOnceCheckpointStore) Advance(
	ctx context.Context,
	binding DurableBinding,
	expectedPosition uint64,
	next Checkpoint,
) (AdvanceResult, error) {
	store.mu.Lock()
	if store.failNext {
		store.failNext = false
		store.mu.Unlock()
		return AdvanceResultUnknown, errors.New("synthetic checkpoint persistence gap")
	}
	store.mu.Unlock()
	return store.delegate.Advance(ctx, binding, expectedPosition, next)
}

func TestInboxRestartAfterCheckpointGapDoesNotRepeatCommittedMutation(t *testing.T) {
	ctx, pool := setupP0405ReliabilityDatabase(t)
	store := NewPostgresInboxStore()
	envelope := testEnvelope(t)
	binding := testInboxBinding(envelope, "billing.checkpoint-gap")
	checkpointStore := &p0405FailOnceCheckpointStore{
		delegate: NewMemoryCheckpointStore(),
		failNext: true,
	}
	registry := NewRegistry()
	var mutationCalls atomic.Int32

	err := registry.Register(Registration{
		Owner:        binding.Owner,
		ConsumerID:   binding.ConsumerID,
		EventTypes:   []EventType{binding.EventType},
		TenantScoped: binding.Scope.TenantID != "",
		Handler: func(handlerCtx context.Context, received Envelope) error {
			return database.InTransaction(handlerCtx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
				result, applyErr := ApplyInbox(
					handlerCtx,
					tx,
					store,
					binding,
					received,
					func(mutationCtx context.Context, mutationTx OutboxTransaction) error {
						mutationCalls.Add(1)
						_, execErr := mutationTx.Exec(
							mutationCtx,
							`INSERT INTO omnexa_events.p04_05_reliability_owner_state (event_id, consumer_id, note) VALUES ($1::uuid, $2, 'checkpoint-gap')`,
							string(received.ID),
							binding.ConsumerID,
						)
						return execErr
					},
				)
				if applyErr != nil {
					return applyErr
				}
				if result != InboxApplied && result != InboxAlreadyApplied {
					return errors.New("inbox application returned unexpected reliability result")
				}
				return nil
			})
		},
	})
	if err != nil {
		t.Fatalf("registry.Register() error = %v", err)
	}

	consumer, err := NewDurableConsumer(ctx, registry, checkpointStore, binding)
	if err != nil {
		t.Fatalf("NewDurableConsumer() error = %v", err)
	}
	delivery := Delivery{Position: 1, Envelope: envelope}
	firstErr := consumer.Process(ctx, delivery)
	assertFailureCode(t, firstErr, codeCheckpointWriteFailed)

	assertP0405ReliabilityCount(
		t,
		ctx,
		pool,
		`SELECT count(*) FROM omnexa_events.p04_05_reliability_owner_state WHERE event_id=$1::uuid AND consumer_id=$2`,
		[]any{string(envelope.ID), binding.ConsumerID},
		1,
	)
	assertP0405ReliabilityCount(
		t,
		ctx,
		pool,
		`SELECT count(*) FROM omnexa_events.consumer_inbox WHERE event_id=$1::uuid AND consumer_id=$2 AND processing_state='completed'`,
		[]any{string(envelope.ID), binding.ConsumerID},
		1,
	)
	if _, exists, loadErr := consumer.LastCheckpoint(ctx); loadErr != nil || exists {
		t.Fatalf("checkpoint unexpectedly advanced after synthetic gap: exists=%v error=%v", exists, loadErr)
	}

	// Recreate the durable consumer to simulate process restart. The checkpoint is
	// still missing, so position 1 is redelivered. The committed inbox completion
	// must suppress the protected mutation while allowing checkpoint recovery.
	restarted, err := NewDurableConsumer(ctx, registry, checkpointStore, binding)
	if err != nil {
		t.Fatalf("restart NewDurableConsumer() error = %v", err)
	}
	if err = restarted.Process(ctx, delivery); err != nil {
		t.Fatalf("restart Process() error = %v", err)
	}
	if mutationCalls.Load() != 1 {
		t.Fatalf("protected mutation calls = %d, want exactly 1 across checkpoint gap", mutationCalls.Load())
	}
	checkpoint, exists, err := restarted.LastCheckpoint(ctx)
	if err != nil || !exists || checkpoint.Position != 1 || checkpoint.EventID != envelope.ID {
		t.Fatalf("recovered checkpoint = %#v exists=%v error=%v", checkpoint, exists, err)
	}
	assertP0405ReliabilityCount(
		t,
		ctx,
		pool,
		`SELECT count(*) FROM omnexa_events.p04_05_reliability_owner_state WHERE event_id=$1::uuid AND consumer_id=$2`,
		[]any{string(envelope.ID), binding.ConsumerID},
		1,
	)
}

func setupP0405ReliabilityDatabase(t *testing.T) (context.Context, *pgxpool.Pool) {
	t.Helper()
	databaseURL := os.Getenv("P04_05_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P04_05_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		cancel()
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		cancel()
		t.Fatalf("pool.Ping() error = %v", err)
	}

	foundation, err := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if err != nil {
		pool.Close()
		cancel()
		t.Fatalf("foundation migrator create error = %v", err)
	}
	if err = foundation.Run(ctx); err != nil {
		pool.Close()
		cancel()
		t.Fatalf("foundation migrator run error = %v", err)
	}
	resetP0405ReliabilityDatabase(t, ctx, pool)

	outboxSQL, err := os.ReadFile("../../migrations/kernel.events/1_create_transactional_outbox.sql")
	if err != nil {
		pool.Close()
		cancel()
		t.Fatalf("read P04.04 migration error = %v", err)
	}
	inboxSQL, err := os.ReadFile("../../migrations/kernel.events/2_create_consumer_inbox.sql")
	if err != nil {
		pool.Close()
		cancel()
		t.Fatalf("read P04.05 migration error = %v", err)
	}
	migrator, err := database.NewMigrator(pool, "kernel.events", []database.Migration{
		{Version: 1, Name: "create_transactional_outbox", SQL: string(outboxSQL)},
		{Version: 2, Name: "create_consumer_inbox", SQL: string(inboxSQL)},
	}, 5*time.Second)
	if err != nil {
		pool.Close()
		cancel()
		t.Fatalf("kernel.events migrator create error = %v", err)
	}
	if err = migrator.Run(ctx); err != nil {
		pool.Close()
		cancel()
		t.Fatalf("kernel.events migrator run error = %v", err)
	}
	if _, err = pool.Exec(ctx, `CREATE TABLE omnexa_events.p04_05_reliability_owner_state (
		event_id uuid NOT NULL,
		consumer_id text NOT NULL,
		note text NOT NULL,
		PRIMARY KEY (event_id, consumer_id)
	)`); err != nil {
		pool.Close()
		cancel()
		t.Fatalf("create reliability owner table error = %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cleanupCancel()
		resetP0405ReliabilityDatabase(t, cleanupCtx, pool)
		pool.Close()
		cancel()
	})
	return ctx, pool
}

func resetP0405ReliabilityDatabase(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(ctx, `DROP SCHEMA IF EXISTS omnexa_events CASCADE`); err != nil {
		t.Fatalf("drop omnexa_events reliability schema error = %v", err)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM omnexa_kernel.schema_migrations WHERE owner='kernel.events'`); err != nil {
		t.Fatalf("reset kernel.events reliability migration ledger error = %v", err)
	}
}

func assertP0405ReliabilityCount(
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
		t.Fatalf("reliability count result/error = %d/%v; expected %d", count, err, expected)
	}
}

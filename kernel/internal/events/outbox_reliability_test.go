package events

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresOutboxRestartRecoversAndRelaysCanonicalPendingEvent(t *testing.T) {
	ctx, pool, store := setupP0404ReliabilityDatabase(t)
	scope := OutboxScope{
		Owner:    Producer("urn:omnexa:module:reliability"),
		TenantID: tenancy.TenantID(mustP0404UUIDv7(t)),
	}
	envelope := mustP0404Envelope(t, scope, "restart-recovery")

	if err := database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		result, enqueueErr := EnqueueOutbox(ctx, tx, store, scope, envelope)
		if enqueueErr != nil {
			return enqueueErr
		}
		if result != OutboxEnqueued {
			return errors.New("restart fixture enqueue did not report enqueued")
		}
		return nil
	}); err != nil {
		t.Fatalf("restart fixture enqueue error = %v", err)
	}

	// Reconstruct the store/relay objects after commit. Recovery must come from
	// durable committed PostgreSQL state, not from process-local memory.
	restartedStore, err := NewPostgresOutboxStore(pool)
	if err != nil {
		t.Fatalf("NewPostgresOutboxStore() after restart error = %v", err)
	}
	var published []EventID
	publisher, err := NewPublisher(func(_ context.Context, candidate Envelope) error {
		published = append(published, candidate.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	relay, err := NewOutboxRelay(restartedStore, publisher)
	if err != nil {
		t.Fatalf("NewOutboxRelay() error = %v", err)
	}

	report, err := relay.RelayPending(ctx, scope, 8)
	if err != nil {
		t.Fatalf("RelayPending() after restart error = %v", err)
	}
	if len(published) != 1 || published[0] != envelope.ID {
		t.Fatalf("restart relay published IDs = %#v; want only %s", published, envelope.ID)
	}
	if report.PendingObserved != 1 || report.AcceptedPublications != 1 || report.PublishedMarks != 1 || report.DuplicateMarkObserved != 0 {
		t.Fatalf("restart relay report = %#v", report)
	}
	pending, err := restartedStore.Pending(ctx, scope, 8)
	if err != nil {
		t.Fatalf("Pending() after restart relay error = %v", err)
	}
	if len(pending) != 0 {
		t.Fatalf("restart relay left pending records = %#v", pending)
	}
}

func TestPostgresOutboxPublishSuccessMarkFailureRemainsRecoverableForDuplicateRelay(t *testing.T) {
	ctx, pool, store := setupP0404ReliabilityDatabase(t)
	scope := OutboxScope{
		Owner:    Producer("urn:omnexa:module:reliability"),
		TenantID: tenancy.TenantID(mustP0404UUIDv7(t)),
	}
	envelope := mustP0404Envelope(t, scope, "publish-mark-crash-window")

	if err := database.InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		result, enqueueErr := EnqueueOutbox(ctx, tx, store, scope, envelope)
		if enqueueErr != nil {
			return enqueueErr
		}
		if result != OutboxEnqueued {
			return errors.New("crash-window fixture enqueue did not report enqueued")
		}
		return nil
	}); err != nil {
		t.Fatalf("crash-window fixture enqueue error = %v", err)
	}

	failingStore := &failOnceOutboxMarkStore{OutboxStore: store}
	var mutex sync.Mutex
	published := make([]EventID, 0, 2)
	publisher, err := NewPublisher(func(_ context.Context, candidate Envelope) error {
		mutex.Lock()
		published = append(published, candidate.ID)
		mutex.Unlock()
		return nil
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	firstRelay, err := NewOutboxRelay(failingStore, publisher)
	if err != nil {
		t.Fatalf("NewOutboxRelay() first error = %v", err)
	}
	firstReport, firstErr := firstRelay.RelayPending(ctx, scope, 1)
	assertFailureCode(t, firstErr, codeOutboxMarkFailed)
	if firstReport.AcceptedPublications != 1 || firstReport.PublishedMarks != 0 {
		t.Fatalf("first crash-window report = %#v", firstReport)
	}

	stillPending, err := store.Pending(ctx, scope, 1)
	if err != nil {
		t.Fatalf("Pending() after synthetic mark failure error = %v", err)
	}
	if len(stillPending) != 1 || stillPending[0].Envelope.ID != envelope.ID {
		t.Fatalf("mark failure lost recoverable pending record = %#v", stillPending)
	}

	// Simulate a later process attempt with a newly constructed store/relay. The
	// same canonical EventID may be published again by design, then marked once.
	restartedStore, err := NewPostgresOutboxStore(pool)
	if err != nil {
		t.Fatalf("NewPostgresOutboxStore() second attempt error = %v", err)
	}
	secondRelay, err := NewOutboxRelay(restartedStore, publisher)
	if err != nil {
		t.Fatalf("NewOutboxRelay() second error = %v", err)
	}
	secondReport, secondErr := secondRelay.RelayPending(ctx, scope, 1)
	if secondErr != nil {
		t.Fatalf("second RelayPending() error = %v", secondErr)
	}
	if secondReport.AcceptedPublications != 1 || secondReport.PublishedMarks != 1 {
		t.Fatalf("second crash-window report = %#v", secondReport)
	}

	mutex.Lock()
	publishedCopy := append([]EventID(nil), published...)
	mutex.Unlock()
	if len(publishedCopy) != 2 || publishedCopy[0] != envelope.ID || publishedCopy[1] != envelope.ID {
		t.Fatalf("duplicate crash-window publications = %#v; want same EventID twice", publishedCopy)
	}
	pending, err := restartedStore.Pending(ctx, scope, 1)
	if err != nil {
		t.Fatalf("Pending() after recovery mark error = %v", err)
	}
	if len(pending) != 0 {
		t.Fatalf("recovered duplicate relay left pending records = %#v", pending)
	}
}

type failOnceOutboxMarkStore struct {
	OutboxStore
	failed bool
}

func (store *failOnceOutboxMarkStore) MarkPublished(ctx context.Context, scope OutboxScope, eventID EventID, revision uint64) (MarkResult, error) {
	if !store.failed {
		store.failed = true
		return MarkResultUnknown, errors.New("synthetic mark persistence interruption")
	}
	return store.OutboxStore.MarkPublished(ctx, scope, eventID, revision)
}

func setupP0404ReliabilityDatabase(t *testing.T) (context.Context, *pgxpool.Pool, *PostgresOutboxStore) {
	t.Helper()
	databaseURL := os.Getenv("P04_04_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P04_04_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
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
	t.Cleanup(func() {
		resetP0404OutboxDatabase(t, context.Background(), pool)
		pool.Close()
		cancel()
	})

	foundation, err := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if err != nil {
		t.Fatalf("foundation migrator create error = %v", err)
	}
	if err = foundation.Run(ctx); err != nil {
		t.Fatalf("foundation migrator run error = %v", err)
	}
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

	store, err := NewPostgresOutboxStore(pool)
	if err != nil {
		t.Fatalf("NewPostgresOutboxStore() error = %v", err)
	}
	return ctx, pool, store
}

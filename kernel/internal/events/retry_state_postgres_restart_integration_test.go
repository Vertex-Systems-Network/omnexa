package events

import (
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresRetryStateStorePersistsLifecycleAcrossFreshPool(t *testing.T) {
	ctx, _ := setupP0406RetryDatabase(t)
	databaseURL := p0406RestartDatabaseURL(t)

	writerPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("writer pgxpool.New() error = %v", err)
	}
	if err = writerPool.Ping(ctx); err != nil {
		writerPool.Close()
		t.Fatalf("writer pool.Ping() error = %v", err)
	}
	writerStore, err := NewPostgresRetryStateStore(writerPool)
	if err != nil {
		writerPool.Close()
		t.Fatalf("NewPostgresRetryStateStore(writer) error = %v", err)
	}

	scheduled := testRetryStateRecord(t, RetryStateScheduled)
	scheduled.Position = 31
	scheduled.NextEligibleAt = time.Now().UTC().Add(-time.Minute).Truncate(time.Microsecond)
	quarantined := testRetryStateRecord(t, RetryStateQuarantined)
	quarantined.Position = 32
	resolved := testRetryStateRecord(t, RetryStateResolved)
	resolved.Position = 33

	for _, record := range []RetryStateRecord{scheduled, quarantined, resolved} {
		created, createErr := writerStore.Create(ctx, record)
		if createErr != nil {
			writerPool.Close()
			t.Fatalf("Create(state=%s position=%d) error = %v", record.State, record.Position, createErr)
		}
		if !retryPostgresRecordsEqual(created, record) {
			writerPool.Close()
			t.Fatalf("Create(state=%s) round-trip drifted: got=%+v want=%+v", record.State, created, record)
		}
	}

	// Simulate process/store restart by closing the writer pool completely and
	// constructing a fresh pool/store against the same already-migrated database.
	writerPool.Close()

	readerPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("reader pgxpool.New() error = %v", err)
	}
	t.Cleanup(readerPool.Close)
	if err = readerPool.Ping(ctx); err != nil {
		t.Fatalf("reader pool.Ping() error = %v", err)
	}
	readerStore, err := NewPostgresRetryStateStore(readerPool)
	if err != nil {
		t.Fatalf("NewPostgresRetryStateStore(reader) error = %v", err)
	}

	loaded := make(map[RetryState]RetryStateRecord, 3)
	for _, want := range []RetryStateRecord{scheduled, quarantined, resolved} {
		got, loadErr := readerStore.Load(ctx, want.Identity, want.Position)
		if loadErr != nil {
			t.Fatalf("Load after restart(state=%s position=%d) error = %v", want.State, want.Position, loadErr)
		}
		if !retryPostgresRecordsEqual(got, want) {
			t.Fatalf("restart state drifted for %s: got=%+v want=%+v", want.State, got, want)
		}
		loaded[want.State] = got
	}

	scheduledDirective, err := RetryExecutionDirectiveFor(loaded[RetryStateScheduled])
	if err != nil {
		t.Fatalf("scheduled directive after restart error = %v", err)
	}
	if scheduledDirective.Eligibility != RetryExecutionAwaitingClaim {
		t.Fatalf("scheduled state became executable without authoritative claim after restart: %+v", scheduledDirective)
	}

	quarantineDirective, err := RetryExecutionDirectiveFor(loaded[RetryStateQuarantined])
	if err != nil {
		t.Fatalf("quarantined directive after restart error = %v", err)
	}
	if quarantineDirective.Eligibility != RetryExecutionQuarantined || quarantineDirective.NextAttempt != 0 {
		t.Fatalf("quarantined state became executable after restart: %+v", quarantineDirective)
	}

	resolvedDirective, err := RetryExecutionDirectiveFor(loaded[RetryStateResolved])
	if err != nil {
		t.Fatalf("resolved directive after restart error = %v", err)
	}
	if resolvedDirective.Eligibility != RetryExecutionResolved || resolvedDirective.NextAttempt != 0 {
		t.Fatalf("resolved state became executable after restart: %+v", resolvedDirective)
	}

	binding := retryDueBinding(scheduled.Identity)
	candidates, err := readerStore.ListDueCandidates(ctx, binding, maxPostgresRetryDueBatch)
	if err != nil {
		t.Fatalf("ListDueCandidates after restart error = %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("due candidates after restart = %d, want only scheduled state: %+v", len(candidates), candidates)
	}
	candidate := candidates[0]
	if candidate.Identity != scheduled.Identity || candidate.Position != scheduled.Position || candidate.Revision != scheduled.Revision {
		t.Fatalf("due candidate identity drifted after restart: %+v", candidate)
	}

	claimToken := testUUIDv7(t)
	claimed, err := readerStore.ClaimDue(
		ctx,
		candidate.Identity,
		candidate.Position,
		candidate.Revision,
		claimToken,
		time.Minute,
	)
	if err != nil {
		t.Fatalf("ClaimDue after restart error = %v", err)
	}
	claimedDirective, err := RetryExecutionDirectiveFor(claimed)
	if err != nil {
		t.Fatalf("claimed directive after restart error = %v", err)
	}
	if claimedDirective.Eligibility != RetryExecutionClaimed || claimed.ClaimToken != claimToken {
		t.Fatalf("scheduled state did not become executable only through authoritative post-restart claim: record=%+v directive=%+v", claimed, claimedDirective)
	}
}

func p0406RestartDatabaseURL(t *testing.T) string {
	t.Helper()
	databaseURL := os.Getenv("P04_06_TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("P04_05_TEST_DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("P04_06_TEST_DATABASE_URL or P04_05_TEST_DATABASE_URL is not set")
	}
	return databaseURL
}

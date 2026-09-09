package events

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type retryClaimAttempt struct {
	record RetryStateRecord
	err    error
	token  string
}

func TestPostgresRetryStateStoreCASLeaseIntegration(t *testing.T) {
	ctx, pool := setupP0406RetryDatabase(t)
	store, err := NewPostgresRetryStateStore(pool)
	if err != nil {
		t.Fatalf("NewPostgresRetryStateStore() error = %v", err)
	}

	record := testRetryStateRecord(t, RetryStateScheduled)
	record.NextEligibleAt = time.Now().UTC().Add(-time.Minute).Truncate(time.Microsecond)
	created, err := store.Create(ctx, record)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) {
			t.Fatalf("Create() failed with SQLSTATE=%s constraint=%s", postgresError.Code, postgresError.ConstraintName)
		}
		t.Fatalf("Create() error = %v", err)
	}
	if !retryPostgresRecordsEqual(created, record) {
		t.Fatalf("Create() round-trip mismatch\ncreated: %+v\nwant: %+v", created, record)
	}

	loaded, err := store.Load(ctx, record.Identity, record.Position)
	if err != nil || !retryPostgresRecordsEqual(loaded, record) {
		t.Fatalf("Load() record/error = %+v/%v", loaded, err)
	}
	replayed, err := store.Create(ctx, record)
	if err != nil || !retryPostgresRecordsEqual(replayed, record) {
		t.Fatalf("exact Create() replay record/error = %+v/%v", replayed, err)
	}

	conflicting := record
	conflicting.Fingerprint[0] ^= 0xff
	_, err = store.Create(ctx, conflicting)
	assertFailureCode(t, err, codeRetryPostgresConflict)

	malformed := record
	malformed.Revision = 0
	_, err = store.Create(ctx, malformed)
	assertFailureCode(t, err, codeRetryPostgresInvalid)

	preclaimed := record
	preclaimed.ClaimToken = testUUIDv7(t)
	preclaimed.ClaimExpiresAt = preclaimed.NextEligibleAt.Add(2 * time.Minute)
	_, err = store.Create(ctx, preclaimed)
	assertFailureCode(t, err, codeRetryPostgresInvalid)

	missing := testRetryStateRecord(t, RetryStateScheduled)
	_, err = store.Load(ctx, missing.Identity, missing.Position)
	assertFailureCode(t, err, codeRetryPostgresNotFound)

	firstToken := testUUIDv7(t)
	competingToken := testUUIDv7(t)
	attempts := make(chan retryClaimAttempt, 2)
	for _, token := range []string{firstToken, competingToken} {
		go func(candidate string) {
			claimed, claimErr := store.ClaimDue(ctx, record.Identity, record.Position, record.Revision, candidate, 2*time.Minute)
			attempts <- retryClaimAttempt{record: claimed, err: claimErr, token: candidate}
		}(token)
	}

	var winner retryClaimAttempt
	wins := 0
	conflicts := 0
	for range 2 {
		attempt := <-attempts
		if attempt.err == nil {
			winner = attempt
			wins++
			continue
		}
		if hasFailureCode(attempt.err, codeRetryPostgresConflict) {
			conflicts++
			continue
		}
		t.Fatalf("concurrent ClaimDue() unexpected error = %v", attempt.err)
	}
	if wins != 1 || conflicts != 1 {
		t.Fatalf("concurrent ClaimDue() wins/conflicts = %d/%d, want 1/1", wins, conflicts)
	}
	if winner.record.Revision != record.Revision+1 || winner.record.ClaimToken != winner.token {
		t.Fatalf("winning claim = %+v, token = %s", winner.record, winner.token)
	}

	if _, err = pool.Exec(
		ctx,
		`UPDATE omnexa_events.consumer_retry_state
		 SET claim_expires_at = CURRENT_TIMESTAMP - INTERVAL '1 second'
		 WHERE event_id = $1::uuid AND consumer_id = $2 AND delivery_position = $3::bigint`,
		string(record.Identity.EventID),
		record.Identity.ConsumerID,
		strconv.FormatUint(record.Position, 10),
	); err != nil {
		t.Fatalf("expire winning lease error = %v", err)
	}

	takeoverToken := testUUIDv7(t)
	takeover, err := store.ClaimDue(ctx, record.Identity, record.Position, winner.record.Revision, takeoverToken, 2*time.Minute)
	if err != nil {
		t.Fatalf("expired-lease ClaimDue() error = %v", err)
	}
	if takeover.Revision != winner.record.Revision+1 || takeover.ClaimToken != takeoverToken || takeover.AttemptsConsumed != winner.record.AttemptsConsumed {
		t.Fatalf("expired-lease takeover = %+v", takeover)
	}

	staleNext := winner.record
	staleNext.ClaimToken = ""
	staleNext.ClaimExpiresAt = time.Time{}
	staleNext.Revision = winner.record.Revision + 1
	staleNext.AttemptsConsumed++
	staleNext.NextEligibleAt = time.Now().UTC().Add(time.Minute).Truncate(time.Microsecond)
	_, err = store.TransitionClaimed(ctx, staleNext, winner.record.Revision, winner.token)
	assertFailureCode(t, err, codeRetryPostgresConflict)

	nextScheduled := takeover
	nextScheduled.ClaimToken = ""
	nextScheduled.ClaimExpiresAt = time.Time{}
	nextScheduled.Revision = takeover.Revision + 1
	nextScheduled.AttemptsConsumed++
	nextScheduled.NextEligibleAt = time.Now().UTC().Add(time.Minute).Truncate(time.Microsecond)
	scheduled, err := store.TransitionClaimed(ctx, nextScheduled, takeover.Revision, takeoverToken)
	if err != nil || !retryPostgresRecordsEqual(scheduled, nextScheduled) {
		t.Fatalf("scheduled transition record/error = %+v/%v", scheduled, err)
	}
	_, err = store.TransitionClaimed(ctx, nextScheduled, takeover.Revision, takeoverToken)
	assertFailureCode(t, err, codeRetryPostgresConflict)

	quarantine := testRetryStateRecord(t, RetryStateScheduled)
	quarantine.Position = 8
	quarantine.NextEligibleAt = time.Now().UTC().Add(-time.Minute).Truncate(time.Microsecond)
	if _, err = store.Create(ctx, quarantine); err != nil {
		t.Fatalf("create quarantine candidate error = %v", err)
	}
	quarantineToken := testUUIDv7(t)
	quarantineClaim, err := store.ClaimDue(ctx, quarantine.Identity, quarantine.Position, 1, quarantineToken, 2*time.Minute)
	if err != nil {
		t.Fatalf("quarantine ClaimDue() error = %v", err)
	}
	quarantinedNext := quarantineClaim
	quarantinedNext.ClaimToken = ""
	quarantinedNext.ClaimExpiresAt = time.Time{}
	quarantinedNext.Revision = quarantineClaim.Revision + 1
	quarantinedNext.AttemptsConsumed++
	quarantinedNext.State = RetryStateQuarantined
	quarantinedNext.NextEligibleAt = time.Time{}
	quarantinedNext.TerminalReason = RetryTerminalReasonNonRetryable
	quarantinedNext.Failure.Retryable = false
	quarantined, err := store.TransitionClaimed(ctx, quarantinedNext, quarantineClaim.Revision, quarantineToken)
	if err != nil || !retryPostgresRecordsEqual(quarantined, quarantinedNext) {
		t.Fatalf("quarantine transition record/error = %+v/%v", quarantined, err)
	}
	_, err = store.ClaimDue(ctx, quarantine.Identity, quarantine.Position, quarantined.Revision, testUUIDv7(t), time.Minute)
	assertFailureCode(t, err, codeRetryPostgresConflict)

	resolution := testRetryStateRecord(t, RetryStateScheduled)
	resolution.Position = 9
	resolution.NextEligibleAt = time.Now().UTC().Add(-time.Minute).Truncate(time.Microsecond)
	if _, err = store.Create(ctx, resolution); err != nil {
		t.Fatalf("create resolution candidate error = %v", err)
	}
	resolutionToken := testUUIDv7(t)
	resolutionClaim, err := store.ClaimDue(ctx, resolution.Identity, resolution.Position, 1, resolutionToken, 2*time.Minute)
	if err != nil {
		t.Fatalf("resolution ClaimDue() error = %v", err)
	}
	resolvedNext := resolutionClaim
	resolvedNext.ClaimToken = ""
	resolvedNext.ClaimExpiresAt = time.Time{}
	resolvedNext.Revision = resolutionClaim.Revision + 1
	resolvedNext.State = RetryStateResolved
	resolvedNext.NextEligibleAt = time.Time{}
	resolvedNext.TerminalReason = RetryTerminalReasonNone
	resolved, err := store.TransitionClaimed(ctx, resolvedNext, resolutionClaim.Revision, resolutionToken)
	if err != nil || !retryPostgresRecordsEqual(resolved, resolvedNext) {
		t.Fatalf("resolved transition record/error = %+v/%v", resolved, err)
	}

	if _, err = pool.Exec(ctx, `DROP TABLE omnexa_events.consumer_retry_state`); err != nil {
		t.Fatalf("drop retry table for safe-error test = %v", err)
	}
	_, err = store.Load(ctx, record.Identity, record.Position)
	assertFailureCode(t, err, codeRetryPostgresFailed)
	publicText := strings.ToLower(err.Error())
	for _, forbidden := range []string{"consumer_retry_state", "relation", "select", "sql"} {
		if strings.Contains(publicText, forbidden) {
			t.Fatalf("safe retry postgres error leaked provider detail %q in %q", forbidden, err.Error())
		}
	}
}

func TestPostgresRetryStateStoreDueCandidatesIntegration(t *testing.T) {
	ctx, pool := setupP0406RetryDatabase(t)
	store, err := NewPostgresRetryStateStore(pool)
	if err != nil {
		t.Fatalf("NewPostgresRetryStateStore() error = %v", err)
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	first := testRetryStateRecord(t, RetryStateScheduled)
	first.Position = 21
	first.NextEligibleAt = now.Add(-3 * time.Minute)
	second := testRetryStateRecord(t, RetryStateScheduled)
	second.Position = 22
	second.NextEligibleAt = now.Add(-2 * time.Minute)
	future := testRetryStateRecord(t, RetryStateScheduled)
	future.Position = 23
	future.NextEligibleAt = now.Add(time.Hour)
	active := testRetryStateRecord(t, RetryStateScheduled)
	active.Position = 24
	active.NextEligibleAt = now.Add(-time.Minute)
	expired := testRetryStateRecord(t, RetryStateScheduled)
	expired.Position = 25
	expired.NextEligibleAt = now.Add(-30 * time.Second)
	otherConsumer := testRetryStateRecord(t, RetryStateScheduled)
	otherConsumer.Position = 26
	otherConsumer.NextEligibleAt = now.Add(-4 * time.Minute)
	otherConsumer.Identity.ConsumerID = "other.projection"

	for _, candidate := range []RetryStateRecord{first, second, future, active, expired, otherConsumer} {
		if _, err = store.Create(ctx, candidate); err != nil {
			t.Fatalf("Create(position=%d) error = %v", candidate.Position, err)
		}
	}

	activeClaim, err := store.ClaimDue(ctx, active.Identity, active.Position, active.Revision, testUUIDv7(t), 2*time.Minute)
	if err != nil {
		t.Fatalf("active ClaimDue() error = %v", err)
	}
	if activeClaim.ClaimToken == "" {
		t.Fatal("active claim did not obtain a lease")
	}

	expiredToken := testUUIDv7(t)
	expiredClaim, err := store.ClaimDue(ctx, expired.Identity, expired.Position, expired.Revision, expiredToken, 2*time.Minute)
	if err != nil {
		t.Fatalf("expired candidate ClaimDue() error = %v", err)
	}
	if _, err = pool.Exec(
		ctx,
		`UPDATE omnexa_events.consumer_retry_state
		 SET claim_expires_at = CURRENT_TIMESTAMP - INTERVAL '1 second'
		 WHERE event_id = $1::uuid AND consumer_id = $2 AND delivery_position = $3::bigint`,
		string(expired.Identity.EventID),
		expired.Identity.ConsumerID,
		strconv.FormatUint(expired.Position, 10),
	); err != nil {
		t.Fatalf("expire due-scan lease error = %v", err)
	}

	binding := retryDueBinding(first.Identity)
	if _, err = store.ListDueCandidates(ctx, binding, 0); !hasFailureCode(err, codeRetryPostgresInvalid) {
		t.Fatalf("zero due-scan limit error = %v", err)
	}
	if _, err = store.ListDueCandidates(ctx, binding, maxPostgresRetryDueBatch+1); !hasFailureCode(err, codeRetryPostgresInvalid) {
		t.Fatalf("oversized due-scan limit error = %v", err)
	}

	limited, err := store.ListDueCandidates(ctx, binding, 2)
	if err != nil {
		t.Fatalf("ListDueCandidates(limit=2) error = %v", err)
	}
	if len(limited) != 2 || limited[0].Position != first.Position || limited[1].Position != second.Position {
		t.Fatalf("limited due candidates = %+v", limited)
	}

	candidates, err := store.ListDueCandidates(ctx, binding, maxPostgresRetryDueBatch)
	if err != nil {
		t.Fatalf("ListDueCandidates() error = %v", err)
	}
	if len(candidates) != 3 {
		t.Fatalf("due candidate count = %d, want 3: %+v", len(candidates), candidates)
	}
	wantPositions := []uint64{first.Position, second.Position, expired.Position}
	for index, candidate := range candidates {
		if candidate.Position != wantPositions[index] {
			t.Fatalf("candidate[%d].Position = %d, want %d", index, candidate.Position, wantPositions[index])
		}
		if candidate.Identity.ConsumerID != binding.ConsumerID || candidate.NextEligibleAt.IsZero() {
			t.Fatalf("candidate[%d] escaped binding or eligibility evidence: %+v", index, candidate)
		}
	}
	if candidates[2].Revision != expiredClaim.Revision {
		t.Fatalf("expired lease candidate revision = %d, want %d", candidates[2].Revision, expiredClaim.Revision)
	}

	claimToken := testUUIDv7(t)
	claimed, err := store.ClaimDue(
		ctx,
		candidates[0].Identity,
		candidates[0].Position,
		candidates[0].Revision,
		claimToken,
		time.Minute,
	)
	if err != nil {
		t.Fatalf("ClaimDue(scanned candidate) error = %v", err)
	}
	if claimed.ClaimToken != claimToken || claimed.Revision != candidates[0].Revision+1 {
		t.Fatalf("claimed scanned candidate = %+v", claimed)
	}
	_, err = store.ClaimDue(
		ctx,
		candidates[0].Identity,
		candidates[0].Position,
		candidates[0].Revision,
		testUUIDv7(t),
		time.Minute,
	)
	assertFailureCode(t, err, codeRetryPostgresConflict)
}

func retryDueBinding(identity InboxIdentity) DurableBinding {
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

func hasFailureCode(err error, code failure.Code) bool {
	return failure.IsCode(err, code)
}

func setupP0406RetryDatabase(t *testing.T) (context.Context, *pgxpool.Pool) {
	t.Helper()
	databaseURL := os.Getenv("P04_06_TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("P04_05_TEST_DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("P04_06_TEST_DATABASE_URL or P04_05_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	t.Cleanup(cancel)
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	t.Cleanup(pool.Close)
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
	resetP0406RetryDatabase(t, ctx, pool)
	t.Cleanup(func() {
		resetP0406RetryDatabase(t, context.Background(), pool)
	})

	outboxSQL, readErr := os.ReadFile("../../migrations/kernel.events/1_create_transactional_outbox.sql")
	if readErr != nil {
		t.Fatalf("read kernel.events migration 1 error = %v", readErr)
	}
	inboxSQL, readErr := os.ReadFile("../../migrations/kernel.events/2_create_consumer_inbox.sql")
	if readErr != nil {
		t.Fatalf("read kernel.events migration 2 error = %v", readErr)
	}
	retrySQL, readErr := os.ReadFile("../../migrations/kernel.events/3_create_retry_quarantine_state.sql")
	if readErr != nil {
		t.Fatalf("read kernel.events migration 3 error = %v", readErr)
	}
	migrations := []database.Migration{
		{Version: 1, Name: "create_transactional_outbox", SQL: string(outboxSQL)},
		{Version: 2, Name: "create_consumer_inbox", SQL: string(inboxSQL)},
		{Version: 3, Name: "create_retry_quarantine_state", SQL: string(retrySQL)},
	}
	migrator, err := database.NewMigrator(pool, "kernel.events", migrations, 5*time.Second)
	if err != nil {
		t.Fatalf("kernel.events migrator create error = %v", err)
	}
	if err = migrator.Run(ctx); err != nil {
		t.Fatalf("kernel.events migrator run error = %v", err)
	}
	if err = migrator.Run(ctx); err != nil {
		t.Fatalf("kernel.events migration replay error = %v", err)
	}
	return ctx, pool
}

func resetP0406RetryDatabase(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(ctx, `DROP SCHEMA IF EXISTS omnexa_events CASCADE`); err != nil {
		t.Fatalf("drop omnexa_events schema error = %v", err)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM omnexa_kernel.schema_migrations WHERE owner='kernel.events'`); err != nil {
		t.Fatalf("reset kernel.events migration ledger error = %v", err)
	}
}

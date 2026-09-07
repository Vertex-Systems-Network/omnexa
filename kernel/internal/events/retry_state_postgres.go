package events

import (
	"context"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxPostgresRetryRevision   uint64        = 1<<63 - 1
	maxPostgresRetryClaimLease time.Duration = 15 * time.Minute

	codeRetryPostgresInvalid   failure.Code = "events.retry.postgres_invalid"
	codeRetryPostgresNotFound  failure.Code = "events.retry.postgres_not_found"
	codeRetryPostgresConflict  failure.Code = "events.retry.postgres_conflict"
	codeRetryPostgresMalformed failure.Code = "events.retry.postgres_malformed"
	codeRetryPostgresFailed    failure.Code = "events.retry.postgres_failed"
	codeRetryPostgresCanceled  failure.Code = "events.retry.postgres_interrupted"
)

const retryPostgresSelectColumns = `
	event_id::text,
	owner,
	consumer_id,
	event_type,
	stream,
	partition_key,
	COALESCE(tenant_id::text, ''),
	delivery_position,
	canonical_fingerprint,
	policy_id,
	policy_version,
	max_attempts,
	initial_backoff_ms,
	max_backoff_ms,
	attempts_consumed,
	retry_state,
	next_eligible_at,
	terminal_reason,
	failure_code,
	failure_category,
	failure_retryable,
	claim_token::text,
	claim_expires_at,
	revision`

// PostgresRetryStateStore persists the already-accepted P04.06 retry/quarantine
// contract in kernel.events. It reuses an existing pool and does not create a
// second migration, scheduler, worker, or transaction authority.
type PostgresRetryStateStore struct {
	pool *pgxpool.Pool
}

func NewPostgresRetryStateStore(pool *pgxpool.Pool) (*PostgresRetryStateStore, error) {
	if pool == nil {
		return nil, classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres store is invalid")
	}
	return &PostgresRetryStateStore{pool: pool}, nil
}

// Create stores revision one for one exact processing identity and delivery.
// Exact replay is idempotent; any differing fingerprint, policy, or state under
// the same identity is a conflict rather than an overwrite.
func (store *PostgresRetryStateStore) Create(ctx context.Context, record RetryStateRecord) (RetryStateRecord, error) {
	if store == nil || store.pool == nil {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres store is invalid")
	}
	if err := retryPostgresContextError(ctx); err != nil {
		return RetryStateRecord{}, err
	}
	if err := validatePostgresRetryRecord(record); err != nil || record.Revision != 1 || record.ClaimToken != "" || !record.ClaimExpiresAt.IsZero() {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres initial state is invalid")
	}

	identity := record.Identity
	tag, err := store.pool.Exec(
		ctx,
		`INSERT INTO omnexa_events.consumer_retry_state (
			event_id, owner, consumer_id, event_type, stream, partition_key, tenant_id,
			delivery_position, canonical_fingerprint,
			policy_id, policy_version, max_attempts, initial_backoff_ms, max_backoff_ms,
			attempts_consumed, retry_state, next_eligible_at, terminal_reason,
			failure_code, failure_category, failure_retryable,
			claim_token, claim_expires_at, revision, quarantined_at, resolved_at
		) VALUES (
			$1::uuid, $2, $3, $4, $5, $6, $7::uuid,
			$8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17, $18,
			$19, $20, $21,
			$22::uuid, $23, $24,
			CASE WHEN $16 = 'quarantined' THEN CURRENT_TIMESTAMP ELSE NULL END,
			CASE WHEN $16 = 'resolved' THEN CURRENT_TIMESTAMP ELSE NULL END
		)
		ON CONFLICT DO NOTHING`,
		string(identity.EventID),
		string(identity.Owner),
		identity.ConsumerID,
		string(identity.EventType),
		identity.Stream,
		identity.Partition,
		retryPostgresTenant(identity),
		int64(record.Position),
		record.Fingerprint[:],
		record.Policy.ID,
		int64(record.Policy.Version),
		int64(record.Policy.MaxAttempts),
		record.Policy.InitialBackoff.Milliseconds(),
		record.Policy.MaxBackoff.Milliseconds(),
		int64(record.AttemptsConsumed),
		string(record.State),
		retryPostgresTime(record.NextEligibleAt),
		retryPostgresTerminalReason(record.TerminalReason),
		retryPostgresFailureCode(record.Failure),
		retryPostgresFailureCategory(record.Failure),
		retryPostgresFailureRetryable(record.Failure),
		retryPostgresClaimToken(record.ClaimToken),
		retryPostgresTime(record.ClaimExpiresAt),
		int64(record.Revision),
	)
	if err != nil {
		return RetryStateRecord{}, wrappedFailure(err, codeRetryPostgresFailed, failure.CategoryUnavailable, "event retry postgres initial state could not be stored")
	}
	if tag.RowsAffected() != 0 && tag.RowsAffected() != 1 {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresMalformed, failure.CategoryInvariant, "event retry postgres create returned an invalid write count")
	}

	stored, loadErr := store.Load(ctx, identity, record.Position)
	if loadErr != nil {
		return RetryStateRecord{}, loadErr
	}
	if tag.RowsAffected() == 0 && !retryPostgresRecordsEqual(stored, record) {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres identity conflicts with committed state")
	}
	return stored, nil
}

// Load returns one authoritative committed retry-state record for the exact
// processing identity and ordered delivery position.
func (store *PostgresRetryStateStore) Load(ctx context.Context, identity InboxIdentity, position uint64) (RetryStateRecord, error) {
	if store == nil || store.pool == nil {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres store is invalid")
	}
	if err := retryPostgresContextError(ctx); err != nil {
		return RetryStateRecord{}, err
	}
	if err := identity.validate(); err != nil || position == 0 || position > maxPostgresRetryRevision {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres load identity is invalid")
	}

	row := store.pool.QueryRow(
		ctx,
		`SELECT `+retryPostgresSelectColumns+`
		 FROM omnexa_events.consumer_retry_state
		 WHERE event_id = $1::uuid
		   AND owner = $2
		   AND consumer_id = $3
		   AND event_type = $4
		   AND stream = $5
		   AND partition_key = $6
		   AND tenant_id IS NOT DISTINCT FROM $7::uuid
		   AND delivery_position = $8`,
		string(identity.EventID),
		string(identity.Owner),
		identity.ConsumerID,
		string(identity.EventType),
		identity.Stream,
		identity.Partition,
		retryPostgresTenant(identity),
		int64(position),
	)
	record, err := scanPostgresRetryRecord(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresNotFound, failure.CategoryNotFound, "event retry postgres state was not found")
	}
	if err != nil {
		var structured *failure.Error
		if errors.As(err, &structured) {
			return RetryStateRecord{}, err
		}
		return RetryStateRecord{}, wrappedFailure(err, codeRetryPostgresFailed, failure.CategoryUnavailable, "event retry postgres state could not be loaded")
	}
	return record, nil
}

// ClaimDue obtains one bounded active lease using PostgreSQL authoritative time.
// Claim expiry is derived in PostgreSQL from a bounded duration; caller-provided
// wall-clock timestamps never decide whether a row is due or an old lease expired.
func (store *PostgresRetryStateStore) ClaimDue(
	ctx context.Context,
	identity InboxIdentity,
	position uint64,
	expectedRevision uint64,
	claimToken string,
	lease time.Duration,
) (RetryStateRecord, error) {
	if store == nil || store.pool == nil {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres store is invalid")
	}
	if err := retryPostgresContextError(ctx); err != nil {
		return RetryStateRecord{}, err
	}
	if err := identity.validate(); err != nil ||
		position == 0 || position > maxPostgresRetryRevision ||
		expectedRevision == 0 || expectedRevision >= maxPostgresRetryRevision ||
		!retryClaimTokenPattern.MatchString(claimToken) ||
		lease < time.Millisecond || lease > maxPostgresRetryClaimLease || lease%time.Millisecond != 0 {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres claim boundary is invalid")
	}

	row := store.pool.QueryRow(
		ctx,
		`UPDATE omnexa_events.consumer_retry_state
		 SET claim_token = $10::uuid,
		     claim_expires_at = CURRENT_TIMESTAMP + ($11::bigint * INTERVAL '1 millisecond'),
		     revision = revision + 1,
		     updated_at = CURRENT_TIMESTAMP
		 WHERE event_id = $1::uuid
		   AND owner = $2
		   AND consumer_id = $3
		   AND event_type = $4
		   AND stream = $5
		   AND partition_key = $6
		   AND tenant_id IS NOT DISTINCT FROM $7::uuid
		   AND delivery_position = $8
		   AND revision = $9
		   AND retry_state = 'retry_scheduled'
		   AND next_eligible_at <= CURRENT_TIMESTAMP
		   AND (claim_token IS NULL OR claim_expires_at <= CURRENT_TIMESTAMP)
		 RETURNING `+retryPostgresSelectColumns,
		string(identity.EventID),
		string(identity.Owner),
		identity.ConsumerID,
		string(identity.EventType),
		identity.Stream,
		identity.Partition,
		retryPostgresTenant(identity),
		int64(position),
		int64(expectedRevision),
		claimToken,
		lease.Milliseconds(),
	)
	record, err := scanPostgresRetryRecord(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return RetryStateRecord{}, store.retryPostgresClaimConflict(ctx, identity, position, expectedRevision)
	}
	if err != nil {
		var structured *failure.Error
		if errors.As(err, &structured) {
			return RetryStateRecord{}, err
		}
		return RetryStateRecord{}, wrappedFailure(err, codeRetryPostgresFailed, failure.CategoryUnavailable, "event retry postgres claim could not be stored")
	}
	if record.Revision != expectedRevision+1 || record.ClaimToken != claimToken || record.State != RetryStateScheduled {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresMalformed, failure.CategoryInvariant, "event retry postgres claim returned malformed state")
	}
	return record, nil
}

// TransitionClaimed persists the outcome of one actively leased retry attempt.
// The expected revision and exact claim token must both still own a non-expired
// lease, so an expired-lease takeover makes the old worker unable to commit.
func (store *PostgresRetryStateStore) TransitionClaimed(
	ctx context.Context,
	next RetryStateRecord,
	expectedRevision uint64,
	claimToken string,
) (RetryStateRecord, error) {
	if store == nil || store.pool == nil {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres store is invalid")
	}
	if err := retryPostgresContextError(ctx); err != nil {
		return RetryStateRecord{}, err
	}
	if err := validatePostgresRetryRecord(next); err != nil ||
		expectedRevision == 0 || expectedRevision >= maxPostgresRetryRevision ||
		next.Revision != expectedRevision+1 ||
		!retryClaimTokenPattern.MatchString(claimToken) ||
		next.ClaimToken != "" || !next.ClaimExpiresAt.IsZero() {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres transition boundary is invalid")
	}

	identity := next.Identity
	row := store.pool.QueryRow(
		ctx,
		`UPDATE omnexa_events.consumer_retry_state
		 SET attempts_consumed = $17,
		     retry_state = $18,
		     next_eligible_at = $19,
		     terminal_reason = $20,
		     failure_code = $21,
		     failure_category = $22,
		     failure_retryable = $23,
		     claim_token = NULL,
		     claim_expires_at = NULL,
		     revision = revision + 1,
		     updated_at = CURRENT_TIMESTAMP,
		     quarantined_at = CASE WHEN $18 = 'quarantined' THEN CURRENT_TIMESTAMP ELSE NULL END,
		     resolved_at = CASE WHEN $18 = 'resolved' THEN CURRENT_TIMESTAMP ELSE NULL END
		 WHERE event_id = $1::uuid
		   AND owner = $2
		   AND consumer_id = $3
		   AND event_type = $4
		   AND stream = $5
		   AND partition_key = $6
		   AND tenant_id IS NOT DISTINCT FROM $7::uuid
		   AND delivery_position = $8
		   AND canonical_fingerprint = $9
		   AND policy_id = $10
		   AND policy_version = $11
		   AND max_attempts = $12
		   AND initial_backoff_ms = $13
		   AND max_backoff_ms = $14
		   AND revision = $15
		   AND claim_token = $16::uuid
		   AND claim_expires_at > CURRENT_TIMESTAMP
		   AND retry_state = 'retry_scheduled'
		   AND ($17 = attempts_consumed OR $17 = attempts_consumed + 1)
		 RETURNING `+retryPostgresSelectColumns,
		string(identity.EventID),
		string(identity.Owner),
		identity.ConsumerID,
		string(identity.EventType),
		identity.Stream,
		identity.Partition,
		retryPostgresTenant(identity),
		int64(next.Position),
		next.Fingerprint[:],
		next.Policy.ID,
		int64(next.Policy.Version),
		int64(next.Policy.MaxAttempts),
		next.Policy.InitialBackoff.Milliseconds(),
		next.Policy.MaxBackoff.Milliseconds(),
		int64(expectedRevision),
		claimToken,
		int64(next.AttemptsConsumed),
		string(next.State),
		retryPostgresTime(next.NextEligibleAt),
		retryPostgresTerminalReason(next.TerminalReason),
		retryPostgresFailureCode(next.Failure),
		retryPostgresFailureCategory(next.Failure),
		retryPostgresFailureRetryable(next.Failure),
	)
	stored, err := scanPostgresRetryRecord(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return RetryStateRecord{}, store.retryPostgresTransitionConflict(ctx, next, expectedRevision, claimToken)
	}
	if err != nil {
		var structured *failure.Error
		if errors.As(err, &structured) {
			return RetryStateRecord{}, err
		}
		return RetryStateRecord{}, wrappedFailure(err, codeRetryPostgresFailed, failure.CategoryUnavailable, "event retry postgres transition could not be stored")
	}
	if !retryPostgresRecordsEqual(stored, next) {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresMalformed, failure.CategoryInvariant, "event retry postgres transition returned malformed state")
	}
	return stored, nil
}

func (store *PostgresRetryStateStore) retryPostgresClaimConflict(ctx context.Context, identity InboxIdentity, position, expectedRevision uint64) error {
	current, err := store.Load(ctx, identity, position)
	if err != nil {
		return err
	}
	if current.Revision != expectedRevision {
		return classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres claim revision is stale")
	}
	if current.State != RetryStateScheduled {
		return classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres state is terminal and cannot be claimed")
	}
	if current.ClaimToken != "" {
		return classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres state already has an active or unresolved lease")
	}
	return classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres state is not yet claimable")
}

func (store *PostgresRetryStateStore) retryPostgresTransitionConflict(ctx context.Context, next RetryStateRecord, expectedRevision uint64, claimToken string) error {
	current, err := store.Load(ctx, next.Identity, next.Position)
	if err != nil {
		return err
	}
	if !retryPostgresImmutableEqual(current, next) {
		return classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres transition conflicts with immutable evidence")
	}
	if current.Revision != expectedRevision {
		return classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres transition revision is stale")
	}
	if current.State != RetryStateScheduled {
		return classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres terminal state cannot be transitioned")
	}
	if current.ClaimToken != claimToken {
		return classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres transition does not own the active lease")
	}
	return classifiedFailure(codeRetryPostgresConflict, failure.CategoryConflict, "event retry postgres transition lease is expired or no longer claimable")
}

func scanPostgresRetryRecord(row pgx.Row) (RetryStateRecord, error) {
	var (
		eventID, owner, consumerID, eventType, stream, partition, tenant string
		fingerprint                                                      []byte
		policyID                                                         string
		position, policyVersion, maxAttempts                             int64
		initialBackoffMS, maxBackoffMS, attempts, revision               int64
		state                                                            string
		nextEligible, claimExpires                                       pgtype.Timestamptz
		terminalReason, failureCode, failureCategory, claimToken         pgtype.Text
		failureRetryable                                                 pgtype.Bool
	)
	if err := row.Scan(
		&eventID,
		&owner,
		&consumerID,
		&eventType,
		&stream,
		&partition,
		&tenant,
		&position,
		&fingerprint,
		&policyID,
		&policyVersion,
		&maxAttempts,
		&initialBackoffMS,
		&maxBackoffMS,
		&attempts,
		&state,
		&nextEligible,
		&terminalReason,
		&failureCode,
		&failureCategory,
		&failureRetryable,
		&claimToken,
		&claimExpires,
		&revision,
	); err != nil {
		return RetryStateRecord{}, err
	}
	if position <= 0 || policyVersion <= 0 || maxAttempts <= 0 || initialBackoffMS <= 0 || maxBackoffMS <= 0 || attempts <= 0 || revision <= 0 ||
		uint64(position) > maxPostgresRetryRevision || uint64(revision) > maxPostgresRetryRevision || len(fingerprint) != len(InboxFingerprint{}) {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresMalformed, failure.CategoryInvariant, "stored event retry postgres numeric or fingerprint state is malformed")
	}
	if policyVersion > int64(^uint32(0)) || maxAttempts > int64(^uint32(0)) || attempts > int64(^uint32(0)) {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresMalformed, failure.CategoryInvariant, "stored event retry postgres bounded state is malformed")
	}

	record := RetryStateRecord{
		Identity: InboxIdentity{
			EventID:    EventID(eventID),
			Owner:      Producer(owner),
			ConsumerID: consumerID,
			EventType:  EventType(eventType),
			Stream:     stream,
			Partition:  partition,
			TenantID:   tenancy.TenantID(tenant),
		},
		Position: uint64(position),
		Policy: RetryPolicy{
			ID:             policyID,
			Version:        uint32(policyVersion),
			MaxAttempts:    uint32(maxAttempts),
			InitialBackoff: time.Duration(initialBackoffMS) * time.Millisecond,
			MaxBackoff:     time.Duration(maxBackoffMS) * time.Millisecond,
		},
		AttemptsConsumed: uint32(attempts),
		State:            RetryState(state),
		Revision:         uint64(revision),
	}
	copy(record.Fingerprint[:], fingerprint)

	if nextEligible.Valid {
		record.NextEligibleAt = nextEligible.Time.UTC()
	}
	if terminalReason.Valid {
		record.TerminalReason = RetryTerminalReason(terminalReason.String)
	}
	if failureCode.Valid || failureCategory.Valid || failureRetryable.Valid {
		if !failureCode.Valid || !failureCategory.Valid || !failureRetryable.Valid {
			return RetryStateRecord{}, classifiedFailure(codeRetryPostgresMalformed, failure.CategoryInvariant, "stored event retry postgres failure evidence is malformed")
		}
		record.Failure = RetryFailureEvidence{
			Code:      failure.Code(failureCode.String),
			Category:  failure.Category(failureCategory.String),
			Retryable: failureRetryable.Bool,
		}
	}
	if claimToken.Valid {
		record.ClaimToken = claimToken.String
	}
	if claimExpires.Valid {
		record.ClaimExpiresAt = claimExpires.Time.UTC()
	}
	if err := record.Validate(); err != nil {
		return RetryStateRecord{}, classifiedFailure(codeRetryPostgresMalformed, failure.CategoryInvariant, "stored event retry postgres state is malformed")
	}
	return record, nil
}

func validatePostgresRetryRecord(record RetryStateRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	if record.Position > maxPostgresRetryRevision || record.Revision > maxPostgresRetryRevision {
		return classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres revision or position exceeds storage bounds")
	}
	if record.Policy.InitialBackoff%time.Millisecond != 0 || record.Policy.MaxBackoff%time.Millisecond != 0 {
		return classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres policy backoff must preserve millisecond precision")
	}
	if !retryPostgresTimeExact(record.NextEligibleAt) || !retryPostgresTimeExact(record.ClaimExpiresAt) {
		return classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres timestamp exceeds storage precision")
	}
	return nil
}

func retryPostgresTimeExact(value time.Time) bool {
	return value.IsZero() || value.Nanosecond()%1000 == 0
}

func retryPostgresRecordsEqual(left, right RetryStateRecord) bool {
	return retryPostgresImmutableEqual(left, right) &&
		left.AttemptsConsumed == right.AttemptsConsumed &&
		left.State == right.State &&
		retryPostgresTimesEqual(left.NextEligibleAt, right.NextEligibleAt) &&
		left.TerminalReason == right.TerminalReason &&
		left.Failure == right.Failure &&
		left.ClaimToken == right.ClaimToken &&
		retryPostgresTimesEqual(left.ClaimExpiresAt, right.ClaimExpiresAt) &&
		left.Revision == right.Revision
}

func retryPostgresImmutableEqual(left, right RetryStateRecord) bool {
	return left.Identity == right.Identity &&
		left.Fingerprint == right.Fingerprint &&
		left.Position == right.Position &&
		left.Policy == right.Policy
}

func retryPostgresTimesEqual(left, right time.Time) bool {
	return (left.IsZero() && right.IsZero()) || (!left.IsZero() && !right.IsZero() && left.Equal(right))
}

func retryPostgresTenant(identity InboxIdentity) any {
	if identity.TenantID == "" {
		return nil
	}
	return string(identity.TenantID)
}

func retryPostgresTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

func retryPostgresTerminalReason(reason RetryTerminalReason) any {
	if reason == RetryTerminalReasonNone {
		return nil
	}
	return string(reason)
}

func retryPostgresFailureCode(evidence RetryFailureEvidence) any {
	if evidence.Code == "" {
		return nil
	}
	return string(evidence.Code)
}

func retryPostgresFailureCategory(evidence RetryFailureEvidence) any {
	if evidence.Category == "" {
		return nil
	}
	return string(evidence.Category)
}

func retryPostgresFailureRetryable(evidence RetryFailureEvidence) any {
	if evidence.Code == "" && evidence.Category == "" {
		return nil
	}
	return evidence.Retryable
}

func retryPostgresClaimToken(token string) any {
	if token == "" {
		return nil
	}
	return token
}

func retryPostgresContextError(ctx context.Context) error {
	if ctx == nil {
		return classifiedFailure(codeRetryPostgresInvalid, failure.CategoryValidation, "event retry postgres context is invalid")
	}
	if err := ctx.Err(); err != nil {
		return wrappedFailure(err, codeRetryPostgresCanceled, failure.CategoryUnavailable, "event retry postgres operation was interrupted")
	}
	return nil
}

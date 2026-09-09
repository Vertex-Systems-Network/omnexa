#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

python scripts/validate_governance.py
python scripts/validate_p04_activation.py

expected_go="$(tr -d '[:space:]' < .go-version)"
actual_go="$(go env GOVERSION)"
if [[ "$actual_go" != "go${expected_go}" ]]; then
  echo "ERROR: Go toolchain mismatch: got ${actual_go}, want go${expected_go}" >&2
  exit 1
fi

required_files=(
  kernel/internal/events/retry.go
  kernel/internal/events/retry_test.go
  kernel/internal/events/retry_state.go
  kernel/internal/events/retry_state_test.go
  kernel/internal/events/retry_runtime.go
  kernel/internal/events/retry_runtime_test.go
  kernel/internal/events/retry_state_postgres.go
  kernel/internal/events/retry_state_postgres_integration_test.go
  kernel/internal/events/retry_state_postgres_restart_integration_test.go
  kernel/internal/events/retry_checkpoint.go
  kernel/internal/events/retry_checkpoint_integration_test.go
  kernel/internal/events/retry_inbox.go
  kernel/internal/events/retry_inbox_test.go
  kernel/internal/events/durable.go
  kernel/internal/events/inbox.go
  kernel/migrations/kernel.events/3_create_retry_quarantine_state.sql
  docs/roadmap/work-packages/P04.06.md
)
for file in "${required_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P04.06 required source missing: ${file}" >&2
    exit 1
  fi
done

if [[ -z "${P04_06_TEST_DATABASE_URL:-}" ]]; then
  echo "ERROR: P04_06_TEST_DATABASE_URL is required for P04.06 acceptance" >&2
  exit 1
fi

unformatted="$(gofmt -l kernel/internal/events/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P04 event files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P04.06 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

# P04.06 is provider-neutral. Retry/quarantine source must not select a broker,
# external workflow runtime, or provider-specific retry/DLQ dependency.
if grep -nE '"(github.com/segmentio/kafka-go|github.com/IBM/sarama|github.com/nats-io/|github.com/rabbitmq/|github.com/redis/|github.com/confluentinc/|go\.uber\.org/cadence|go\.temporal\.io)' \
  kernel/internal/events/retry*.go; then
  echo "ERROR: P04.06 introduced an unauthorized broker/provider/workflow dependency" >&2
  exit 1
fi

for marker in \
  'MaxRetryAttempts uint32 = 100' \
  'MaxRetryBackoff       = 24 * time.Hour' \
  'func DecideRetry(' \
  'RetryDispositionRetryable' \
  'RetryDispositionTerminal' \
  'RetryDispositionInterrupted'; do
  grep -Fq "$marker" kernel/internal/events/retry.go || { echo "ERROR: P04.06 retry marker missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'type RetryStateRecord struct' \
  'RetryStateScheduled' \
  'RetryStateQuarantined' \
  'RetryStateResolved' \
  'ClaimToken' \
  'Revision'; do
  grep -Fq "$marker" kernel/internal/events/retry_state.go || { echo "ERROR: P04.06 durable-state marker missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'func ComposeInitialRetryFailure(' \
  'func RetryExecutionDirectiveFor(' \
  'func ComposeClaimedRetryResult(' \
  'RetryExecutionAwaitingClaim' \
  'RetryExecutionClaimed'; do
  grep -Fq "$marker" kernel/internal/events/retry_runtime.go || { echo "ERROR: P04.06 runtime marker missing: ${marker}" >&2; exit 1; }
done

# Wave 2D T05: accepted P04.05 applied/already-applied outcomes must resolve an
# authoritative claimed retry without gaining authority to re-run the mutation.
for marker in \
  'func ComposeClaimedRetryInboxResult(' \
  'InboxApplied, InboxAlreadyApplied' \
  'RetryExecutionClaimed'; do
  grep -Fq "$marker" kernel/internal/events/retry_inbox.go || { echo "ERROR: P04.06 retry/inbox composition marker missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'maxPostgresRetryClaimLease time.Duration = 15 * time.Minute' \
  'maxPostgresRetryDueBatch   uint32        = 100' \
  'func (store *PostgresRetryStateStore) ListDueCandidates(' \
  'func (store *PostgresRetryStateStore) ClaimDue(' \
  'func (store *PostgresRetryStateStore) TransitionClaimed('; do
  grep -Fq "$marker" kernel/internal/events/retry_state_postgres.go || { echo "ERROR: P04.06 PostgreSQL marker missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'type RetryCheckpointStateStore interface' \
  'type RetryQuarantineCheckpointResult struct' \
  'QuarantineCommitted bool' \
  'func CommitRetryQuarantineBeforeCheckpoint(' \
  'func RecoverRetryQuarantineCheckpoint('; do
  grep -Fq "$marker" kernel/internal/events/retry_checkpoint.go || { echo "ERROR: P04.06 checkpoint marker missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'TestComposeInitialRetryFailureSchedulesDurableStateBeforeDeferral' \
  'TestComposeInitialRetryFailureInterruptionCreatesNoDurableState' \
  'TestComposeInitialRetryFailureUnknownErrorFailsClosedToQuarantine' \
  'TestRetryExecutionDirectiveRequiresAuthoritativeClaimBeforeExecution' \
  'TestComposeClaimedRetryResultInterruptionReleasesClaimWithoutConsumingBudget' \
  'TestComposeClaimedRetryResultExhaustionQuarantines'; do
  grep -Fq "$marker" kernel/internal/events/retry_runtime_test.go || { echo "ERROR: P04.06 runtime acceptance test missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'TestComposeClaimedRetryInboxResultAppliedAndAlreadyAppliedResolve' \
  'TestComposeClaimedRetryInboxResultFailsClosedForNonSuccessOutcomes' \
  'TestComposeClaimedRetryInboxResultRequiresAuthoritativeClaim'; do
  grep -Fq "$marker" kernel/internal/events/retry_inbox_test.go || { echo "ERROR: P04.06 retry/inbox acceptance test missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'TestPostgresRetryStateStoreCASLeaseIntegration' \
  'TestPostgresRetryStateStoreDueCandidatesIntegration'; do
  grep -Fq "$marker" kernel/internal/events/retry_state_postgres_integration_test.go || { echo "ERROR: P04.06 PostgreSQL acceptance test missing: ${marker}" >&2; exit 1; }
done

grep -Fq 'TestPostgresRetryStateStorePersistsLifecycleAcrossFreshPool' \
  kernel/internal/events/retry_state_postgres_restart_integration_test.go || {
    echo "ERROR: P04.06 restart persistence acceptance test missing" >&2
    exit 1
  }

for marker in \
  'TestRetryQuarantineCommitsBeforeCheckpointAndRecoversCrashGap' \
  'TestRetryQuarantineDoesNotCheckpointWhenTerminalCommitFails' \
  'TestRetryQuarantineCheckpointRecoveryIsIdempotentForExactAcceptedProgress'; do
  grep -Fq "$marker" kernel/internal/events/retry_checkpoint_integration_test.go || { echo "ERROR: P04.06 checkpoint acceptance test missing: ${marker}" >&2; exit 1; }
done

migration="kernel/migrations/kernel.events/3_create_retry_quarantine_state.sql"
for marker in \
  'CREATE TABLE IF NOT EXISTS omnexa_events.consumer_retry_state' \
  'CONSTRAINT events_retry_attempt_policy_bounded' \
  'CONSTRAINT events_retry_backoff_bounded' \
  'CONSTRAINT events_retry_claim_consistent' \
  'CONSTRAINT events_retry_state_evidence_consistent' \
  'CREATE UNIQUE INDEX IF NOT EXISTS consumer_retry_processing_identity_uq' \
  'CREATE INDEX IF NOT EXISTS consumer_retry_due_idx' \
  'CREATE INDEX IF NOT EXISTS consumer_retry_claim_expiry_idx' \
  'CREATE INDEX IF NOT EXISTS consumer_retry_quarantine_idx'; do
  grep -Fq "$marker" "$migration" || { echo "ERROR: P04.06 migration-3 invariant missing: ${marker}" >&2; exit 1; }
done

# Exact P04.06 family runs against real PostgreSQL first. The normal event-package
# run retains P04.01-P04.05 regression coverage; race evidence remains retry-focused.
go test ./kernel/internal/events -run 'Retry' -count=1
go vet ./kernel/internal/events
go test ./kernel/internal/events -count=1
go test -race ./kernel/internal/events -run 'Retry' -count=1
go build ./kernel/...

echo "P04.06 G0 active-package governance + exact bounded path authority: PASS"
echo "P04.06 G1 structured failure classification + unknown-error fail-closed + interruption law: PASS"
echo "P04.06 G2 finite attempt/backoff policy + UTC deterministic eligibility + hard ceilings: PASS"
echo "P04.06 G3 durable scheduled/quarantined/resolved state + bounded safe failure evidence: PASS"
echo "P04.06 G4 authoritative PostgreSQL due discovery + bounded batch + one-at-a-time CAS/lease claim: PASS"
echo "P04.06 G5 lease expiry/takeover + stale transition conflict + monotonic revision/attempt evidence: PASS"
echo "P04.06 G6 retry execution requires authoritative claim and clears/releases claim deterministically: PASS"
echo "P04.06 G7 terminal/exhausted retry commits quarantine before checkpoint advancement: PASS"
echo "P04.06 G8 quarantine-commit/checkpoint-failure crash gap recovers without handler/terminal mutation replay: PASS"
echo "P04.06 G9 exact owner/consumer/route/stream/partition/tenant processing scope remains isolated: PASS"
echo "P04.06 G10 immutable kernel.events migration 3 identity/state/due/claim/quarantine constraints retained: PASS"
echo "P04.06 G11 retained P04.01-P04.05 event-package regression + race + build boundary: PASS"
echo "P04.06 G12 no broker/provider DLQ, P04.07+ runtime, business feature, AI runtime, or exactly-once/global-ordering expansion: PASS"
echo "P04.06 G13 P04.05 applied/already-applied outcomes resolve only an authoritative claimed retry and do not replay protected mutation: PASS"
echo "P04.06 G14 scheduled/quarantined/resolved retry state survives a fresh PostgreSQL pool/store boundary and execution remains claim-gated: PASS"

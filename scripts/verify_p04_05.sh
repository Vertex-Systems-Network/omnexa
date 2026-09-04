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
  kernel/internal/events/envelope.go
  kernel/internal/events/durable.go
  kernel/internal/events/inbox.go
  kernel/internal/events/inbox_core_test.go
  kernel/internal/events/inbox_postgres.go
  kernel/internal/events/inbox_postgres_integration_test.go
  kernel/internal/events/inbox_reliability_test.go
  kernel/internal/events/inbox_concurrency_test.go
  kernel/migrations/kernel.events/1_create_transactional_outbox.sql
  kernel/migrations/kernel.events/2_create_consumer_inbox.sql
  docs/roadmap/work-packages/P04.05.md
)
for file in "${required_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P04.05 required source missing: ${file}" >&2
    exit 1
  fi
done

if [[ -z "${P04_05_TEST_DATABASE_URL:-}" ]]; then
  echo "ERROR: P04_05_TEST_DATABASE_URL is required for P04.05 acceptance" >&2
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
  echo "ERROR: P04.05 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

# P04.05 may use the retained pgx/PostgreSQL boundary but must remain broker-
# neutral and must not pull P04.06+ reliability/runtime packages forward.
if grep -nE '"(github.com/segmentio/kafka-go|github.com/IBM/sarama|github.com/nats-io/|github.com/rabbitmq/|github.com/redis/|github.com/confluentinc/|go\.uber\.org/cadence|go\.temporal\.io)' \
  kernel/internal/events/inbox*.go; then
  echo "ERROR: P04.05 introduced an unauthorized broker/provider/workflow dependency" >&2
  exit 1
fi
if grep -nE 'type[[:space:]]+(DeadLetter|RetryPolicy|RetrySchedule|SchemaRegistry|BackgroundJob|Quarantine)|func[[:space:]].*(DeadLetter|RetryBackoff|RetrySchedule|SchemaRegistry|Quarantine|BackgroundJob)' \
  kernel/internal/events/inbox*.go; then
  echo "ERROR: P04.05 introduced a future-package runtime surface" >&2
  exit 1
fi

for marker in \
  'type InboxIdentity struct' \
  'type InboxRecord struct' \
  'type InboxClaimResult uint8' \
  'type InboxApplyResult uint8' \
  'type InboxStore interface' \
  'type ProtectedMutation func' \
  'func NewInboxRecord(' \
  'func ApplyInbox(' \
  'InboxClaimed' \
  'InboxAlreadyCompleted' \
  'InboxIdentityConflict' \
  'InboxConcurrentProcessing' \
  'InboxApplied' \
  'InboxAlreadyApplied' \
  'InboxConflict' \
  'InboxConcurrent'; do
  if ! grep -Fq "$marker" kernel/internal/events/inbox.go; then
    echo "ERROR: P04.05 required core contract marker missing: ${marker}" >&2
    exit 1
  fi
done

for marker in \
  'type PostgresInboxStore struct' \
  'func NewPostgresInboxStore(' \
  'func (store *PostgresInboxStore) Claim(' \
  'func (store *PostgresInboxStore) Complete('; do
  if ! grep -Fq "$marker" kernel/internal/events/inbox_postgres.go; then
    echo "ERROR: P04.05 required PostgreSQL adapter marker missing: ${marker}" >&2
    exit 1
  fi
done

for marker in \
  'TestNewInboxRecordScopesCanonicalEventToConsumerBinding' \
  'TestNewInboxRecordFingerprintDetectsSameIdentityContentReuse' \
  'TestNewInboxRecordRejectsRouteOrTenantRebinding' \
  'TestApplyInboxFirstDeliveryUsesExactCallerTransactionAndCompletesAfterMutation' \
  'TestApplyInboxDuplicateSkipsProtectedMutation' \
  'TestApplyInboxConflictAndConcurrentFailClosedWithoutMutation' \
  'TestApplyInboxMutationFailureLeavesCompletionUnwrittenAndPreservesError' \
  'TestApplyInboxCompletionFailureForcesCallerTransactionFailure' \
  'TestApplyInboxRejectsInvalidStoreResultAndInterruptedContext'; do
  if ! grep -Fq "$marker" kernel/internal/events/inbox_core_test.go; then
    echo "ERROR: P04.05 required core acceptance test missing: ${marker}" >&2
    exit 1
  fi
done

if ! grep -Fq 'TestPostgresInboxStoreIntegration' kernel/internal/events/inbox_postgres_integration_test.go; then
  echo "ERROR: P04.05 PostgreSQL atomicity/integration evidence is missing" >&2
  exit 1
fi
if ! grep -Fq 'TestInboxRestartAfterCheckpointGapDoesNotRepeatCommittedMutation' kernel/internal/events/inbox_reliability_test.go; then
  echo "ERROR: P04.05 checkpoint-gap/restart acceptance evidence is missing" >&2
  exit 1
fi
for marker in \
  'TestPostgresInboxSameScopeConcurrentAttemptsCannotBothMutate' \
  'TestPostgresInboxSameEventDifferentConsumersRemainIndependentConcurrently' \
  'TestInboxReliabilityRejectsRouteAndTenantRebindingBeforeClaim'; do
  if ! grep -Fq "$marker" kernel/internal/events/inbox_concurrency_test.go; then
    echo "ERROR: P04.05 required concurrency/isolation acceptance test missing: ${marker}" >&2
    exit 1
  fi
done

migration="kernel/migrations/kernel.events/2_create_consumer_inbox.sql"
for marker in \
  'CREATE SCHEMA IF NOT EXISTS omnexa_events' \
  'CREATE TABLE IF NOT EXISTS omnexa_events.consumer_inbox' \
  'canonical_fingerprint bytea NOT NULL' \
  "processing_state text NOT NULL DEFAULT 'claimed'" \
  'CONSTRAINT events_inbox_fingerprint_size' \
  'CONSTRAINT events_inbox_state_valid' \
  'CONSTRAINT events_inbox_completion_consistent' \
  'CREATE UNIQUE INDEX IF NOT EXISTS consumer_inbox_processing_identity_uq' \
  "COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid)"; do
  if ! grep -Fq "$marker" "$migration"; then
    echo "ERROR: P04.05 migration invariant missing: ${marker}" >&2
    exit 1
  fi
done

# Prove the package against real PostgreSQL. The focused run exercises first
# application, rollback, persistence failure, restart/checkpoint-gap handling,
# duplicate suppression, cross-consumer independence and same-scope concurrency.
go test ./kernel/internal/events -run 'Inbox' -count=1
go vet ./kernel/internal/events
go test ./kernel/internal/events -count=1
go test -race ./kernel/internal/events -run 'Inbox' -count=1
go build ./kernel/...

echo "P04.05 G0 active-package governance + exact P04.05 path boundary: PASS"
echo "P04.05 G1 scoped canonical EventID/owner/consumer/tenant/route identity and conflict detection: PASS"
echo "P04.05 G2 same caller PostgreSQL transaction for protected mutation + inbox completion: PASS"
echo "P04.05 G3 mutation/persistence failure rollback and later legitimate retry eligibility: PASS"
echo "P04.05 G4 committed duplicate/restart redelivery returns already-applied without mutation replay: PASS"
echo "P04.05 G5 checkpoint-gap recovery keeps checkpoint and inbox facts distinct: PASS"
echo "P04.05 G6 same-scope concurrency excludes double mutation while cross-consumer EventID remains independent: PASS"
echo "P04.05 G7 route/tenant/content rebinding fails closed and authorization is not inferred from inbox identity: PASS"
echo "P04.05 G8 owner-scoped migration v2 replay/ledger and retained P01-P04.04 regressions: PASS"
echo "P04.05 G9 no broker/provider, retry/DLQ, schema-registry, job, business, AI or end-to-end exactly-once expansion: PASS"

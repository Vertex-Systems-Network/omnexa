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
  kernel/internal/events/bus.go
  kernel/internal/events/durable.go
  kernel/internal/events/outbox.go
  kernel/internal/events/outbox_core_test.go
  kernel/internal/events/outbox_postgres.go
  kernel/internal/events/outbox_postgres_integration_test.go
  kernel/internal/events/outbox_reliability_test.go
  kernel/internal/events/outbox_concurrency_test.go
  kernel/migrations/kernel.events/1_create_transactional_outbox.sql
  docs/roadmap/work-packages/P04.04.md
)
for file in "${required_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P04.04 required source missing: ${file}" >&2
    exit 1
  fi
done

if [[ -z "${P04_04_TEST_DATABASE_URL:-}" ]]; then
  echo "ERROR: P04_04_TEST_DATABASE_URL is required for P04.04 acceptance" >&2
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
  echo "ERROR: P04.04 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

# P04.04 may use the retained pgx/PostgreSQL boundary but must remain transport
# provider-neutral and must not pull future reliability packages forward.
if grep -nE '"(github.com/segmentio/kafka-go|github.com/IBM/sarama|github.com/nats-io/|github.com/rabbitmq/|github.com/redis/|github.com/confluentinc/|go\.uber\.org/cadence|go\.temporal\.io)' \
  kernel/internal/events/outbox*.go; then
  echo "ERROR: P04.04 introduced an unauthorized broker/provider/workflow dependency" >&2
  exit 1
fi
if grep -nE 'type[[:space:]]+(Inbox|DeadLetter|RetryPolicy|SchemaRegistry|BackgroundJob)|func[[:space:]].*(DeadLetter|RetryBackoff|InboxDedup|SchemaRegistry)' \
  kernel/internal/events/outbox*.go; then
  echo "ERROR: P04.04 introduced a future-package runtime surface" >&2
  exit 1
fi

for marker in \
  'type OutboxScope struct' \
  'type OutboxRecord struct' \
  'type OutboxTransaction interface' \
  'type OutboxStore interface' \
  'func EnqueueOutbox(' \
  'type OutboxRelay struct' \
  'func NewOutboxRelay(' \
  'func (relay *OutboxRelay) RelayPending(' \
  'OutboxEnqueued' \
  'OutboxUnchanged' \
  'OutboxEnqueueConflict' \
  'OutboxMarkedPublished' \
  'OutboxAlreadyPublished' \
  'OutboxMarkConflict'; do
  if ! grep -Fq "$marker" kernel/internal/events/outbox.go; then
    echo "ERROR: P04.04 required core contract marker missing: ${marker}" >&2
    exit 1
  fi
done

for marker in \
  'type PostgresOutboxStore struct' \
  'func NewPostgresOutboxStore(' \
  'func (store *PostgresOutboxStore) Enqueue(' \
  'func (store *PostgresOutboxStore) Pending(' \
  'func (store *PostgresOutboxStore) MarkPublished('; do
  if ! grep -Fq "$marker" kernel/internal/events/outbox_postgres.go; then
    echo "ERROR: P04.04 required PostgreSQL adapter marker missing: ${marker}" >&2
    exit 1
  fi
done

for marker in \
  'TestEnqueueOutboxUsesCallerTransactionWithoutOpeningAnother' \
  'TestEnqueueOutboxSeparatesIdempotentRetryFromConflictingIdentity' \
  'TestRelayPublishFailureLeavesRecordUnmarked' \
  'TestRelayMarkFailureKeepsDuplicatePublicationWindowExplicit' \
  'TestRelayRejectsCrossScopeOrMalformedPendingStateBeforePublish' \
  'TestRelayTreatsConcurrentAlreadyPublishedMarkAsDuplicateNotExactlyOnce'; do
  if ! grep -Fq "$marker" kernel/internal/events/outbox_core_test.go; then
    echo "ERROR: P04.04 required core acceptance test missing: ${marker}" >&2
    exit 1
  fi
done

if ! grep -Fq 'TestPostgresOutboxStoreIntegration' kernel/internal/events/outbox_postgres_integration_test.go; then
  echo "ERROR: P04.04 PostgreSQL atomicity/integration evidence is missing" >&2
  exit 1
fi
for marker in \
  'TestPostgresOutboxRestartRecoversAndRelaysCanonicalPendingEvent' \
  'TestPostgresOutboxPublishSuccessMarkFailureRemainsRecoverableForDuplicateRelay'; do
  if ! grep -Fq "$marker" kernel/internal/events/outbox_reliability_test.go; then
    echo "ERROR: P04.04 required restart/crash-window acceptance test missing: ${marker}" >&2
    exit 1
  fi
done
for marker in \
  'TestPostgresOutboxConcurrentRelaysPreservePublicationState' \
  'TestPostgresOutboxConcurrentTenantRelaysRemainIsolated'; do
  if ! grep -Fq "$marker" kernel/internal/events/outbox_concurrency_test.go; then
    echo "ERROR: P04.04 required concurrency/tenant-isolation acceptance test missing: ${marker}" >&2
    exit 1
  fi
done

migration="kernel/migrations/kernel.events/1_create_transactional_outbox.sql"
for marker in \
  'CREATE SCHEMA IF NOT EXISTS omnexa_events' \
  'CREATE TABLE IF NOT EXISTS omnexa_events.transactional_outbox' \
  'event_id uuid PRIMARY KEY' \
  'publication_state' \
  'revision bigint' \
  'CHECK'; do
  if ! grep -Fq "$marker" "$migration"; then
    echo "ERROR: P04.04 migration invariant missing: ${marker}" >&2
    exit 1
  fi
done

# The integration/reliability tests use the retained P01 migration runner and
# the exact caller transaction. They prove rollback, committed-only visibility,
# restart recovery, duplicate-relay crash windows, optimistic CAS and tenant
# isolation against real PostgreSQL state.
go test ./kernel/internal/events -run 'Outbox' -count=1
go vet ./kernel/internal/events
go test ./kernel/internal/events -count=1
go test -race ./kernel/internal/events -run 'Outbox' -count=1
go build ./kernel/...

echo "P04.04 G0 active-package governance + exact P04.04 path boundary: PASS"
echo "P04.04 G1 same-local-PostgreSQL-transaction owner mutation + outbox atomicity: PASS"
echo "P04.04 G2 canonical envelope identity and conflicting EventID fail-closed semantics: PASS"
echo "P04.04 G3 committed-only pending relay and publication-failure recoverability: PASS"
echo "P04.04 G4 restart recovery and publish-success/mark-failure duplicate window: PASS"
echo "P04.04 G5 optimistic publication-state CAS and concurrent duplicate acceptance: PASS"
echo "P04.04 G6 owner/tenant isolation and malformed-state fail-closed behavior: PASS"
echo "P04.04 G7 migration replay/ledger integrity and retained P01-P04.03 regressions: PASS"
echo "P04.04 G8 no broker/provider, inbox/dedup, retry/DLQ, schema-registry, job, business or AI runtime expansion: PASS"

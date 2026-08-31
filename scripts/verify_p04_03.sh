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
  kernel/internal/events/envelope_test.go
  kernel/internal/events/bus.go
  kernel/internal/events/bus_test.go
  kernel/internal/events/durable.go
  kernel/internal/events/durable_test.go
  kernel/internal/events/durable_concurrency_test.go
  docs/roadmap/work-packages/P04.03.md
)
for file in "${required_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P04.03 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/events/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P04 event files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P04.03 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

# P04.03 remains provider-neutral and must not pull future reliability packages,
# concrete brokers, database persistence or job runtime into the active package.
if grep -nE '"(database/sql|github.com/segmentio/kafka-go|github.com/IBM/sarama|github.com/nats-io/|github.com/rabbitmq/|github.com/redis/|github.com/confluentinc/|go\.uber\.org/cadence|go\.temporal\.io)' \
  kernel/internal/events/durable.go \
  kernel/internal/events/durable_test.go \
  kernel/internal/events/durable_concurrency_test.go; then
  echo "ERROR: P04.03 introduced an unauthorized provider/database/workflow dependency" >&2
  exit 1
fi
if grep -nE 'type[[:space:]]+(Outbox|Inbox|DeadLetter|RetryPolicy|SchemaRegistry|BackgroundJob)|func[[:space:]].*(ExactlyOnce|Retry|DeadLetter|Outbox|Inbox)' \
  kernel/internal/events/durable.go; then
  echo "ERROR: P04.03 introduced future-package or exactly-once runtime surface" >&2
  exit 1
fi

for marker in \
  'type DurableScope struct' \
  'type DurableBinding struct' \
  'type Checkpoint struct' \
  'type CheckpointStore interface' \
  'type DurableConsumer struct' \
  'func NewDurableConsumer(' \
  'func (consumer *DurableConsumer) LastCheckpoint(' \
  'func (consumer *DurableConsumer) ResumePosition(' \
  'func (consumer *DurableConsumer) Process(' \
  'type MemoryCheckpointStore struct' \
  'BindingConflict' \
  'CheckpointStale' \
  'CheckpointConflict'; do
  if ! grep -Fq "$marker" kernel/internal/events/durable.go; then
    echo "ERROR: P04.03 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

for marker in \
  'TestDurableConsumerResumesFromLastAcceptedCheckpoint' \
  'TestDurableConsumerFailedOrCancelledWorkDoesNotAdvance' \
  'TestDurableConsumerEnforcesContiguousMonotonicProgress' \
  'TestDurableConsumerRejectsConflictingOwnerAndScopeRebinding' \
  'TestDurableConsumerTenantScopesAreIsolatedAndMetadataIsPreserved' \
  'TestDurableConsumerAllowsDuplicateDeliveryWhenCheckpointWriteFails' \
  'TestDurableConsumerRejectsMalformedCheckpointAndSanitizesStoreFailure' \
  'TestDurableConsumerRejectsInvalidRouteAndScopeBeforeHandler'; do
  if ! grep -Fq "$marker" kernel/internal/events/durable_test.go; then
    echo "ERROR: P04.03 required acceptance test missing: ${marker}" >&2
    exit 1
  fi
done
if ! grep -Fq 'TestDurableConsumerConcurrentSameScopeCheckpointRaceIsFailClosed' kernel/internal/events/durable_concurrency_test.go; then
  echo "ERROR: P04.03 concurrent same-scope checkpoint race evidence is missing" >&2
  exit 1
fi

# Checkpoint state must remain progress metadata only; raw event data does not
# belong in the checkpoint record.
checkpoint_block="$(awk '/type Checkpoint struct/{capture=1} capture{print} capture && /^}/{exit}' kernel/internal/events/durable.go)"
if grep -qE 'Data|Payload|Secret|Credential|Token' <<<"$checkpoint_block"; then
  echo "ERROR: P04.03 checkpoint state contains payload/secret-like fields" >&2
  printf '%s\n' "$checkpoint_block" >&2
  exit 1
fi

# Dedicated P04.03 proof plus the full events package retain P04.01 envelope and
# P04.02 publish/subscribe regressions on the same exact head.
go test ./kernel/internal/events -run '^TestDurableConsumer' -count=1
go vet ./kernel/internal/events
go test ./kernel/internal/events -count=1
go test -race ./kernel/internal/events -count=1
go build ./kernel/...

echo "P04.03 G0 active-package governance + exact P04.03 contract boundary: PASS"
echo "P04.03 G1 provider-neutral durable owner/consumer/route/scope binding: PASS"
echo "P04.03 G2 contiguous monotonic checkpoint advancement, stale/conflict rejection and concurrent same-scope CAS race: PASS"
echo "P04.03 G3 restart/resume from last accepted checkpoint with explicit position-1 convention: PASS"
echo "P04.03 G4 handler failure/cancellation cannot advance unacknowledged work: PASS"
echo "P04.03 G5 tenant/owner/scope isolation and malformed-state fail-closed behavior: PASS"
echo "P04.03 G6 duplicate delivery remains possible across handler-success/checkpoint-failure crash window: PASS"
echo "P04.03 G7 P04.01 envelope metadata + P04.02 registration/handler regressions retained: PASS"
echo "P04.03 G8 no broker/provider, migration, outbox/inbox, retry/DLQ, job, business or AI runtime expansion: PASS"

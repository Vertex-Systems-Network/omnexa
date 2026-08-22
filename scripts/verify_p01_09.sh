#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

expected_go="$(tr -d '[:space:]' < .go-version)"
actual_go="$(go env GOVERSION)"
if [[ "$actual_go" != "go${expected_go}" ]]; then
  echo "ERROR: Go toolchain mismatch: got ${actual_go}, want go${expected_go}" >&2
  exit 1
fi

echo "P01.09 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P01.09 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
echo "P01.09 google/uuid version: ${uuid_version}"

mapfile -t go_files < <(find kernel -type f -name '*.go' -print | sort)
if [[ ${#go_files[@]} -eq 0 ]]; then
  echo "ERROR: no Go source files found under kernel/" >&2
  exit 1
fi

unformatted="$(gofmt -l "${go_files[@]}")"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: Go module/workspace metadata is not canonical" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

go vet ./kernel/...
go test ./kernel/...
go test -race ./kernel/internal/jobs -count=1
go test -v ./kernel/internal/jobs -run 'Test(Execution|Registry|Retry|Idempotency|ConcurrentDuplicate|Cancellation|Worker|Graceful|Shutdown|Observability|Schedules|Panicking|ConcurrentExecutions)' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/jobs"|"$module_prefix/kernel/internal/failure"|"$module_prefix/kernel/internal/observability"|"$module_prefix/kernel/internal/config"|"$module_prefix/kernel/internal/buildinfo") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P01.09 jobs package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/jobs)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(cache|database|storage|operations)' kernel/internal/jobs --include='*.go'; then
  echo "ERROR: P01.09 jobs source contains provider/module/later-kernel coupling" >&2
  exit 1
fi

if go list -f '{{join .Imports "\n"}}' ./kernel/internal/jobs | grep -E '(^|/)(database/sql|nats|jetstream)(/|$)|github\.com/nats-io|go\.temporal\.io|robfig/cron'; then
  echo "ERROR: P01.09 imports durable messaging/workflow/cron persistence machinery" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Tenant|Organization|User|Customer|Order|Invoice|Product|Payment|Workflow|Event|Outbox|Inbox|Agent|Planner|Model|FeatureFlag|AuditRecord)([[:space:]]|$)|[[:space:]](TenantID|OrganizationID|UserID|CustomerID|OrderID|PaymentID|AgentID|GoalID|PlanID|TaskID)[[:space:]]' kernel/internal/jobs --include='*.go'; then
  echo "ERROR: P01.09 declares unauthorized tenant/business/event/workflow/AI/later-package concepts" >&2
  exit 1
fi

if grep -R -nE 'nats\.Conn|JetStream|TransactionalOutbox|WorkflowTimer|cron\.New' kernel/internal/jobs --include='*.go' | grep -v '_test.go'; then
  echo "ERROR: P01.09 implements unauthorized durable messaging/workflow scheduling machinery" >&2
  exit 1
fi

go build ./kernel/...

echo "P01.09 G1 format/static/dependency/portable-boundary: PASS"
echo "P01.09 G2 unit/race registry/retry/idempotency/schedule: PASS"
echo "P01.09 G3 execution/duplicate/observability contract integration: PASS"
echo "P01.09 G5 scheduler non-authority/safe failure/scope guard: PASS"
echo "P01.09 G6 bounded concurrency/retry/cancellation/drain/shutdown resilience: PASS"
echo "P01.09 G7 build/package: PASS"

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

echo "P01.11 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P01.11 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
echo "P01.11 google/uuid version: ${uuid_version}"

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
go test -race ./kernel/internal/audit -count=1
go test -v ./kernel/internal/audit -run 'Test(Record|Memory|Required|Impersonation|Audit|Caller|Concurrent)' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/audit"|"$module_prefix/kernel/internal/failure"|"$module_prefix/kernel/internal/observability"|"$module_prefix/kernel/internal/config"|"$module_prefix/kernel/internal/buildinfo") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P01.11 audit package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/audit)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(cache|database|storage|operations|jobs|configuration)' kernel/internal/audit --include='*.go'; then
  echo "ERROR: P01.11 audit source contains provider/module/later-kernel coupling" >&2
  exit 1
fi

if go list -f '{{join .Imports "\n"}}' ./kernel/internal/audit | grep -E '(^|/)(database/sql|net/http)(/|$)|github\.com/nats-io|go\.temporal\.io|robfig/cron'; then
  echo "ERROR: P01.11 imports unauthorized persistence, HTTP, messaging, workflow, or scheduler machinery" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Tenant|Organization|User|Customer|Order|Invoice|Product|Payment|Workflow|Event|Outbox|Inbox|Agent|Planner|Model|Permission|Role|Entitlement|Experiment)([[:space:]]|$)' kernel/internal/audit --include='*.go'; then
  echo "ERROR: P01.11 declares unauthorized identity, business, workflow, entitlement, experimentation, or AI concepts" >&2
  exit 1
fi

if grep -R -nE 'func[[:space:]].*(Update|Delete|Export|ReadAll|Replace)|type[[:space:]].*(Repository|Retention|LegalHold|Outbox|Inbox)' kernel/internal/audit --include='*.go' | grep -v '_test.go'; then
  echo "ERROR: P01.11 implements unauthorized audit mutation/export/retention or durable messaging machinery" >&2
  exit 1
fi

go build ./kernel/...

echo "P01.11 G1 format/static/dependency/ownership boundary: PASS"
echo "P01.11 G2 unit/race envelope/integrity/sink/health: PASS"
echo "P01.11 G3 append/required-vs-degraded/metadata contract integration: PASS"
echo "P01.11 G5 classification/prohibited-secret/non-authority/security negatives: PASS"
echo "P01.11 G6 cancellation/panic/capacity/tamper/concurrency resilience: PASS"
echo "P01.11 G7 build/package: PASS"

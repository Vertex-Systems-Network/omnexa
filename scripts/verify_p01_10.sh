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

echo "P01.10 Go toolchain: ${actual_go}"

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
go test -race ./kernel/internal/configuration -count=1
go test -v ./kernel/internal/configuration -run 'Test(Registry|Evaluation|KillSwitch|Cache|Refresh|Caller|Provider|Concurrent)' -count=1

# P01.10 originally owned the whole configuration package. P02.09 now adds an
# authorized scoped layer in that same owner. Keep this historical regression
# guard on direct runtime dependencies so the new layer may compose only the
# accepted P02 audit/authorization/tenancy/organization capabilities without
# turning transitive prerequisite imports into false coupling failures.
module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/audit"|\
    "$module_prefix/kernel/internal/authorization"|\
    "$module_prefix/kernel/internal/failure"|\
    "$module_prefix/kernel/internal/organization"|\
    "$module_prefix/kernel/internal/tenancy") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: configuration runtime directly imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -f '{{join .Imports "\n"}}' ./kernel/internal/configuration)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(config|cache|database|storage|operations|jobs|observability)' kernel/internal/configuration --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: configuration runtime source contains static/provider/later-kernel coupling" >&2
  exit 1
fi

if go list -f '{{join .Imports "\n"}}' ./kernel/internal/configuration | grep -E '(^|/)(database/sql|net/http)(/|$)|github\.com/nats-io|go\.temporal\.io|robfig/cron'; then
  echo "ERROR: configuration runtime imports unauthorized transport, messaging, or workflow machinery" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Tenant|Organization|User|Customer|Order|Invoice|Product|Payment|Workflow|Event|Outbox|Inbox|Agent|Planner|Model|Permission|Role|Entitlement|Experiment|AuditRecord|Secret|Credential|Token)([[:space:]]|$)' kernel/internal/configuration --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: configuration runtime declares unauthorized identity, business, authorization-owner, entitlement, secret, or AI concepts" >&2
  exit 1
fi

go build ./kernel/...

echo "P01.10 G1 format/static/direct-runtime-dependency/ownership boundary: PASS"
echo "P01.10 G2 unit/race definitions/evaluation/cache/change metadata: PASS"
echo "P01.10 G3 provider/default/fallback/context contract integration: PASS"
echo "P01.10 G5 kill-switch fail-closed/non-authority/scope/security negatives: PASS"
echo "P01.10 G6 bounded refresh/provider timeout/panic/cache resilience: PASS"
echo "P01.10 G7 build/package: PASS"

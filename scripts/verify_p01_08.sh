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

echo "P01.08 Go toolchain: ${actual_go}"

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
go test -race ./kernel/internal/operations -count=1
go test -v ./kernel/internal/operations -run 'Test(Liveness|Optional|Required|Dependency|Caller|Diagnostic|Build|Registry|Lifecycle|Concurrent|Panicking)' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/operations"|"$module_prefix/kernel/internal/buildinfo"|"$module_prefix/kernel/internal/observability"|"$module_prefix/kernel/internal/config"|"$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P01.08 operations package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/operations)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(cache|database|storage)' kernel/internal/operations --include='*.go'; then
  echo "ERROR: P01.08 operations source contains provider/module coupling" >&2
  exit 1
fi

if grep -R -nE 'net/http|k8s\.io/|kubernetes|http\.Handle|http\.Server' kernel/internal/operations --include='*.go'; then
  echo "ERROR: P01.08 must remain a portable health primitive and must not implement a public/Kubernetes-specific status surface" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Tenant|User|Customer|Order|Invoice|Product|AuditEvent|Job|Scheduler|FeatureFlag|Agent|Planner|Model|Embedding|Vector|Workflow|Event|ModuleHealth)([[:space:]]|$)|[[:space:]](TenantID|UserID|CustomerID|OrderID|AuditID|AgentID|GoalID|PlanID|TaskID|ExecutionID)[[:space:]]' kernel/internal/operations --include='*.go'; then
  echo "ERROR: P01.08 declares unauthorized tenant/business/audit/scheduler/AI/module-health concepts" >&2
  exit 1
fi

if grep -R -nE 'ConnectionString|DSN|SQL|ObjectKey|Credentials|Password|Secret|Token' kernel/internal/operations --include='*.go' | grep -v '_test.go'; then
  echo "ERROR: P01.08 production diagnostic model must not expose secret/provider payload fields" >&2
  exit 1
fi

go build ./kernel/...

echo "P01.08 G1 format/static/dependency/portable-boundary: PASS"
echo "P01.08 G2 unit/race lifecycle/readiness/registry: PASS"
echo "P01.08 G3 health contract/build-identity/observability integration: PASS"
echo "P01.08 G5 safe diagnostic projection/security-critical fail-closed/scope guard: PASS"
echo "P01.08 G6 bounded timeout/cancellation/panic/stopping resilience: PASS"
echo "P01.08 G7 build/package: PASS"

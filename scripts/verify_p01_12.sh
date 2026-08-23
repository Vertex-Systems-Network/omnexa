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

echo "P01.12 Go toolchain: ${actual_go}"

if [[ -z "${P01_12_TEST_DATABASE_URL:-}" ]]; then
  echo "ERROR: P01_12_TEST_DATABASE_URL is required for canonical P01.12 migration evidence" >&2
  exit 1
fi

python scripts/validate_governance.py
python scripts/validate_p01_preparation.py
python scripts/validate_p01_package_specs.py

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
go test -race ./kernel/internal/developer ./kernel/cmd/omnexa -count=1
go test -v ./kernel/internal/developer -run 'Test(Help|Health|Database|Verify|Module|Unknown|OS)' -count=1

direct_imports="$(go list -f '{{join .Imports "\n"}}' ./kernel/internal/developer)"
if grep -E 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform|apps|shared)|kernel/internal/(cache|storage|jobs|configuration|audit)' <<<"$direct_imports"; then
  echo "ERROR: P01.12 developer CLI imports unauthorized business/module/kernel implementation" >&2
  exit 1
fi

if grep -R -nE '(^|[[:space:]])(kubectl|helm|terraform|docker|sudo)([[:space:]]|$)|sh[[:space:]]+-c|bash[[:space:]]+-c' kernel/internal/developer kernel/cmd/omnexa --include='*.go'; then
  echo "ERROR: P01.12 developer CLI contains unauthorized deployment/shell-bypass machinery" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Tenant|Organization|User|Customer|Order|Invoice|Product|Payment|Workflow|Event|Outbox|Inbox|Agent|Planner|Model|Permission|Role|Entitlement)([[:space:]]|$)' kernel/internal/developer kernel/cmd/omnexa --include='*.go'; then
  echo "ERROR: P01.12 declares unauthorized identity, business, workflow, entitlement, or AI concepts" >&2
  exit 1
fi

help_output="$(go run ./kernel/cmd/omnexa help)"
for marker in 'version' 'health' 'db migrate' 'verify <target>'; do
  if [[ "$help_output" != *"$marker"* ]]; then
    echo "ERROR: CLI help missing marker: ${marker}" >&2
    exit 1
  fi
done

version_output="$(go run ./kernel/cmd/omnexa version)"
if [[ "$version_output" != omnexa\ version=*\ commit=* ]]; then
  echo "ERROR: CLI version output is not deterministic" >&2
  exit 1
fi

health_output="$(OMNEXA_ENVIRONMENT=ci OMNEXA_DATABASE_URL="$P01_12_TEST_DATABASE_URL" go run ./kernel/cmd/omnexa health)"
if [[ "$health_output" != *'"readiness":"healthy"'* ]]; then
  echo "ERROR: CLI health output is not healthy JSON" >&2
  exit 1
fi
if [[ "$health_output" == *"postgres://"* || "$health_output" == *"synthetic-p01-04"* ]]; then
  echo "ERROR: CLI health output leaked restricted database configuration" >&2
  exit 1
fi

migration_output="$(OMNEXA_ENVIRONMENT=ci OMNEXA_DATABASE_URL="$P01_12_TEST_DATABASE_URL" go run ./kernel/cmd/omnexa db migrate)"
if [[ "$migration_output" != "omnexa db migrate: PASS environment=ci" ]]; then
  echo "ERROR: CLI migration boundary did not return the expected safe result" >&2
  exit 1
fi

negative_output="$(mktemp)"
trap 'rm -f "$negative_output"' EXIT
if OMNEXA_ENVIRONMENT=production OMNEXA_DATABASE_URL="$P01_12_TEST_DATABASE_URL" go run ./kernel/cmd/omnexa db migrate >"$negative_output" 2>&1; then
  echo "ERROR: production developer migration unexpectedly succeeded" >&2
  exit 1
fi
if grep -Fq "$P01_12_TEST_DATABASE_URL" "$negative_output" || grep -Fq 'postgres://' "$negative_output"; then
  echo "ERROR: failed migration output leaked restricted database configuration" >&2
  exit 1
fi

if go run ./kernel/cmd/omnexa verify definitely-not-a-gate >"$negative_output" 2>&1; then
  echo "ERROR: unknown verification target unexpectedly succeeded" >&2
  exit 1
fi

go mod verify
go build ./kernel/...

echo "P01.12 G0 governance/active-package/readiness: PASS"
echo "P01.12 G1 format/static/dependency/command-boundary: PASS"
echo "P01.12 G2 unit/race CLI parsing/output/orchestration: PASS"
echo "P01.12 G3 help/version/health/verification command contracts: PASS"
echo "P01.12 G4 guarded PostgreSQL migration command boundary: PASS"
echo "P01.12 G5 secret-safe output/non-production/non-authority negatives: PASS"
echo "P01.12 G6 fail-closed command/runner/verification resilience: PASS"
echo "P01.12 G7 build/package: PASS"
echo "P01.12 G8 module checksum verification: PASS"

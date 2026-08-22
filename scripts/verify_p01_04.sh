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

echo "P01.04 Go toolchain: ${actual_go}"

if [[ -z "${P01_04_TEST_DATABASE_URL:-}" ]]; then
  echo "ERROR: P01_04_TEST_DATABASE_URL is required for P01.04 integration evidence" >&2
  exit 1
fi

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
  exit 1
fi

pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: pgx version = ${pgx_version}, want v5.10.0" >&2
  exit 1
fi

echo "P01.04 pgx version: ${pgx_version}"

go vet ./kernel/...
go test ./kernel/...
go test -race ./kernel/internal/database -run 'Test(LoadConfiguration|Settings|MigrationChecksum|NewMigrator|AdvisoryLockKey|UnavailableConnection)' -count=1
go test -v ./kernel/internal/database -run 'TestPostgreSQLFoundationIntegration|TestUnavailableConnectionIsBoundedAndSafe' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/database"|"$module_prefix/kernel/internal/config"|"$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P01.04 database package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/database)

if go list -m all | grep -E '^(gorm\.io/|entgo\.io/ent($|/)|github\.com/uptrace/bun($|/))' >/dev/null; then
  echo "ERROR: P01.04 must not introduce an ORM" >&2
  exit 1
fi

if [[ ! -f kernel/migrations/README.md ]]; then
  echo "ERROR: owner-scoped migration repository convention is missing" >&2
  exit 1
fi
if ! grep -Fq 'kernel/migrations/<owner>/<version>_<name>.sql' kernel/migrations/README.md; then
  echo "ERROR: owner-scoped migration path convention is incomplete" >&2
  exit 1
fi

if find kernel/internal/database -type f -name '*.go' ! -name '*_test.go' -print0 | xargs -0 grep -nEi '\b(tenant|organization|outbox|inbox|cache|storage|telemetry)\b'; then
  echo "ERROR: P01.04 runtime source contains pulled-forward schema/capability markers" >&2
  exit 1
fi

go build ./kernel/...

echo "P01.04 G1 format/static/dependency boundary: PASS"
echo "P01.04 G2 unit/race: PASS"
echo "P01.04 G3 PostgreSQL connection/transaction integration: PASS"
echo "P01.04 G4 fresh/idempotent/upgrade/failure/ledger migration: PASS"
echo "P01.04 G5 secret redaction and scope guard: PASS"
echo "P01.04 G7 build/package: PASS"

#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

python scripts/validate_governance.py
python scripts/validate_p02_preparation.py
python scripts/validate_p02_package_specs.py

expected_go="$(tr -d '[:space:]' < .go-version)"
actual_go="$(go env GOVERSION)"
if [[ "$actual_go" != "go${expected_go}" ]]; then
  echo "ERROR: Go toolchain mismatch: got ${actual_go}, want go${expected_go}" >&2
  exit 1
fi

echo "P02.01 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P02.01 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.01 requires pinned github.com/jackc/pgx/v5 v5.10.0, got ${pgx_version}" >&2
  exit 1
fi
echo "P02.01 google/uuid version: ${uuid_version}"
echo "P02.01 pgx version: ${pgx_version}"

migration="kernel/migrations/kernel.identity/1_create_identity_foundation.sql"
if [[ ! -f "$migration" ]]; then
  echo "ERROR: missing P02.01 kernel.identity migration: ${migration}" >&2
  exit 1
fi

mapfile -t go_files < <(find kernel -type f -name '*.go' -print | sort)
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
go test -race ./kernel/internal/identity -count=1

if [[ -z "${P02_01_TEST_DATABASE_URL:-}" ]]; then
  echo "ERROR: P02_01_TEST_DATABASE_URL is required for canonical P02.01 migration/repository evidence" >&2
  exit 1
fi
go test -v ./kernel/internal/identity -run '^TestPostgresIdentityFoundationIntegration$' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/identity"|"$module_prefix/kernel/internal/database"|"$module_prefix/kernel/internal/failure"|"$module_prefix/kernel/internal/config") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P02.01 identity package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/identity)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(audit|cache|configuration|developer|jobs|observability|operations|storage)' kernel/internal/identity --include='*.go'; then
  echo "ERROR: P02.01 identity source contains later-domain or unrelated kernel coupling" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Tenant|Organization|Membership|Session|Credential|Role|Permission|Policy|MFA|Passkey|ServiceAccount|APIKey|Person)([[:space:]]|$)' kernel/internal/identity --include='*.go'; then
  echo "ERROR: P02.01 declares unauthorized tenancy/authentication/authorization/service-account/business concepts" >&2
  exit 1
fi

if grep -nEi '\b(tenant_id|organization_id|password|password_hash|session_id|role_id|permission_id|mfa_secret|passkey|api_key)\b' "$migration"; then
  echo "ERROR: P02.01 migration contains premature tenancy/credential/session/authorization fields" >&2
  exit 1
fi

if grep -R -nE '(fmt\.(Print|Printf|Println)|log\.|slog\.)' kernel/internal/identity --include='*.go' | grep -v '_test.go'; then
  echo "ERROR: P02.01 runtime identity code must not log CONFIDENTIAL/PII values directly" >&2
  exit 1
fi

go build ./kernel/...

echo "P02.01 G0 governance/active-package/owner boundary: PASS"
echo "P02.01 G1 format/static/dependency/migration ownership: PASS"
echo "P02.01 G2 unit/race UUIDv7/lifecycle/safe projection: PASS"
echo "P02.01 G3 PostgreSQL repository/optimistic transition contract: PASS"
echo "P02.01 G4 fresh/idempotent/immutable-ledger migration evidence: PASS"
echo "P02.01 G5 CONFIDENTIAL-PII safe output/no-premature-authority negatives: PASS"
echo "P02.01 G6 stale transition/disabled terminal lifecycle resilience: PASS"
echo "P02.01 G7 build/package: PASS"
echo "P02.01 G8 pinned identity/database dependencies: PASS"

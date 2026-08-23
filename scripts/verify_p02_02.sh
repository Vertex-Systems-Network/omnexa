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

echo "P02.02 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P02.02 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.02 requires pinned github.com/jackc/pgx/v5 v5.10.0, got ${pgx_version}" >&2
  exit 1
fi
echo "P02.02 google/uuid version: ${uuid_version}"
echo "P02.02 pgx version: ${pgx_version}"

migration="kernel/migrations/kernel.tenancy/1_create_tenancy_foundation.sql"
if [[ ! -f "$migration" ]]; then
  echo "ERROR: missing P02.02 kernel.tenancy migration: ${migration}" >&2
  exit 1
fi

for marker in 'omnexa_tenancy.tenants' 'tenant_id' 'principal_id' 'relationship_state'; do
  if ! grep -Fq "$marker" "$migration"; then
    echo "ERROR: P02.02 tenancy migration missing required isolation marker: ${marker}" >&2
    exit 1
  fi
done

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

integration_url="${P02_02_TEST_DATABASE_URL:-}"
P02_02_TEST_DATABASE_URL= go test ./kernel/...
P02_02_TEST_DATABASE_URL= go test -race ./kernel/internal/tenancy -count=1

if [[ -z "$integration_url" ]]; then
  echo "ERROR: P02_02_TEST_DATABASE_URL is required for canonical P02.02 migration/tenant-isolation evidence" >&2
  exit 1
fi
P02_02_TEST_DATABASE_URL="$integration_url" go test -v ./kernel/internal/tenancy -run '^TestPostgresTenancyFoundationIntegration$' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/tenancy"|"$module_prefix/kernel/internal/identity"|"$module_prefix/kernel/internal/database"|"$module_prefix/kernel/internal/failure"|"$module_prefix/kernel/internal/config") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P02.02 tenancy package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/tenancy)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(audit|cache|configuration|developer|jobs|observability|operations|storage)' kernel/internal/tenancy --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.02 tenancy runtime source contains later-domain or unrelated kernel coupling" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Organization|Session|Credential|Role|Permission|Policy|MFA|Passkey|ServiceAccount|APIKey|Person|TenantSetting)([[:space:]]|$)' kernel/internal/tenancy --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.02 declares unauthorized organization/authentication/authorization/service-account/settings/business concepts" >&2
  exit 1
fi

if grep -nEi '\b(organization_id|role_id|permission_id|policy_id|password|password_hash|session_id|mfa_secret|passkey|api_key|setting_key|global_tenant_id)\b' "$migration"; then
  echo "ERROR: P02.02 migration contains premature organization/credential/session/authorization/settings/global-tenant fields" >&2
  exit 1
fi

if grep -R -nE '(fmt\.(Print|Printf|Println)|log\.|slog\.)' kernel/internal/tenancy --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.02 runtime tenancy code must not emit tenant/principal context directly to ordinary logs" >&2
  exit 1
fi

if grep -R -nE 'func[[:space:]]+NewTrustedContext' kernel/internal/tenancy --include='*.go'; then
  echo "ERROR: P02.02 TrustedContext must not have a caller-constructible exported constructor" >&2
  exit 1
fi

for marker in 'type TrustedContext struct' 'func (trusted TrustedContext) ScopeFor' 'target != trusted.tenantID'; do
  if ! grep -Fq "$marker" kernel/internal/tenancy/context.go; then
    echo "ERROR: P02.02 trusted-context fail-closed marker missing: ${marker}" >&2
    exit 1
  fi
done

go build ./kernel/...

echo "P02.02 G0 governance/active-package/kernel.tenancy owner boundary: PASS"
echo "P02.02 G1 format/static/dependency/migration ownership: PASS"
echo "P02.02 G2 unit/race Tenant/membership/context/scope semantics: PASS"
echo "P02.02 G3 PostgreSQL tenant repository/trusted-context contract: PASS"
echo "P02.02 G4 fresh/idempotent/P02.01-upgrade/immutable-ledger migration evidence: PASS"
echo "P02.02 G5 same-tenant allow/cross-tenant forged-selector deny/no-global-fallback: PASS"
echo "P02.02 G6 tenant suspension/membership revocation/stale-transition resilience: PASS"
echo "P02.02 G7 build/package: PASS"
echo "P02.02 G8 pinned identity/database dependencies: PASS"

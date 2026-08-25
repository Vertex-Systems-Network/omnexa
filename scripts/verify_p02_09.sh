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
echo "P02.09 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$uuid_version" != "v1.6.0" || "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.09 pinned dependency mismatch: uuid=${uuid_version} pgx=${pgx_version}" >&2
  exit 1
fi

configuration_migration="kernel/migrations/kernel.configuration/1_create_scoped_settings.sql"
authorization_migration="kernel/migrations/kernel.authorization/3_add_configuration_permissions.sql"
for migration in "$configuration_migration" "$authorization_migration"; do
  if [[ ! -f "$migration" ]]; then
    echo "ERROR: missing P02.09 migration: ${migration}" >&2
    exit 1
  fi
done

for marker in \
  'omnexa_configuration.setting_overrides' \
  'tenant_id uuid NOT NULL' \
  'organization_id uuid NULL' \
  'configuration_override_organization_same_tenant' \
  'configuration_one_tenant_override_per_key' \
  'configuration_one_organization_override_per_key' \
  'revision bigint NOT NULL'; do
  if ! grep -Fq "$marker" "$configuration_migration"; then
    echo "ERROR: P02.09 configuration migration missing marker: ${marker}" >&2
    exit 1
  fi
done
for marker in 'configuration.setting.read' 'configuration.setting.manage'; do
  if ! grep -Fq "$marker" "$authorization_migration"; then
    echo "ERROR: P02.09 authorization permission migration missing marker: ${marker}" >&2
    exit 1
  fi
done

if grep -nEi '\b(user_id|global_value|global_override|role_id|permission_id|secret_value|password|credential|api_key|token_value)\b' "$configuration_migration"; then
  echo "ERROR: P02.09 configuration schema contains user/global/authorization/secret scope" >&2
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

integration_url="${P02_09_TEST_DATABASE_URL:-}"
P02_01_TEST_DATABASE_URL= P02_02_TEST_DATABASE_URL= P02_03_TEST_DATABASE_URL= P02_04_TEST_DATABASE_URL= P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= P02_07_TEST_DATABASE_URL= P02_08_TEST_DATABASE_URL= P02_09_TEST_DATABASE_URL= \
  go test ./kernel/...
P02_09_TEST_DATABASE_URL= go test -race ./kernel/internal/configuration ./kernel/internal/authorization -count=1

if [[ -z "$integration_url" ]]; then
  echo "ERROR: P02_09_TEST_DATABASE_URL is required for canonical P02.09 migration/scope evidence" >&2
  exit 1
fi
P02_01_TEST_DATABASE_URL= P02_02_TEST_DATABASE_URL= P02_03_TEST_DATABASE_URL= P02_04_TEST_DATABASE_URL= P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= P02_07_TEST_DATABASE_URL= P02_08_TEST_DATABASE_URL= P02_09_TEST_DATABASE_URL="$integration_url" \
  go test -v ./kernel/internal/configuration -run '^TestPostgresTenantScopedSettingsIntegration$' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/audit"|\
    "$module_prefix/kernel/internal/authorization"|\
    "$module_prefix/kernel/internal/failure"|\
    "$module_prefix/kernel/internal/organization"|\
    "$module_prefix/kernel/internal/tenancy") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P02.09 configuration runtime directly imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -f '{{join .Imports "\n"}}' ./kernel/internal/configuration)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(developer|jobs|operations|storage)' \
  kernel/internal/configuration --include='scoped*.go' --include='postgres_scoped_repository.go'; then
  echo "ERROR: P02.09 scoped runtime pulled business/future/deployment coupling forward" >&2
  exit 1
fi
if grep -R -nE '(fmt\.(Print|Printf|Println)|log\.|slog\.)' \
  kernel/internal/configuration/scoped*.go kernel/internal/configuration/postgres_scoped_repository.go; then
  echo "ERROR: P02.09 scoped runtime must not log classified setting values" >&2
  exit 1
fi

for marker in \
  'type TrustedSettingScope struct' \
  'func ScopeFromTenantContext' \
  'func ScopeFromOrganizationContext' \
  'scope.UserID != ""' \
  'policy.AllowOrganizationOverride' \
  'PermissionSettingRead' \
  'PermissionSettingManage' \
  'service.authorizer.Require' \
  'audit.RequirementRequired' \
  'service.evaluator.InvalidateKey' \
  'type PostgresScopedRepository struct'; do
  if ! grep -RFq "$marker" kernel/internal/configuration; then
    echo "ERROR: P02.09 trusted-scope/authorization/audit marker missing: ${marker}" >&2
    exit 1
  fi
done

if grep -R -nE 'func[[:space:]]+(New|Create|From)[A-Za-z0-9_]*Scope\([^)]*(string|TenantID|NodeID)' \
  kernel/internal/configuration/scoped_types.go; then
  echo "ERROR: P02.09 exposes a raw tenant/org setting-scope constructor" >&2
  exit 1
fi
if grep -R -nE 'type[[:space:]]+(BusinessSetting|ModuleSetting|DeploymentSetting|SecretStore|Entitlement|AdminBypass|SuperAdmin|AIAgent)([[:space:]]|$)' \
  kernel/internal/configuration --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.09 pulled business/deployment/secrets/authority/AI scope forward" >&2
  exit 1
fi

go build ./kernel/...

echo "P02.09 G0 governance/active-package/kernel.configuration owner boundary: PASS"
echo "P02.09 G1 format/static/direct-dependency/no-raw-scope/no-future-scope: PASS"
echo "P02.09 G2 unit/race classification/precedence/value semantics: PASS"
echo "P02.09 G3 PostgreSQL trusted tenant/org resolution plus authorization/audit: PASS"
echo "P02.09 G4 fresh/idempotent/P02.08 supported-upgrade migration evidence: PASS"
echo "P02.09 G5 wrong-tenant/wrong-org/protected-read/manage/no-global-user-secret negatives: PASS"
echo "P02.09 G6 exact-scope fallback/cache-invalidation/current-authorization resilience: PASS"
echo "P02.09 G7 build/package: PASS"
echo "P02.09 G8 pinned configuration/tenancy/organization/authorization/audit dependency chain: PASS"

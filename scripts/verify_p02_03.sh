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
echo "P02.03 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P02.03 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.03 requires pinned github.com/jackc/pgx/v5 v5.10.0, got ${pgx_version}" >&2
  exit 1
fi
echo "P02.03 google/uuid version: ${uuid_version}"
echo "P02.03 pgx version: ${pgx_version}"

migration="kernel/migrations/kernel.organization/1_create_organization_foundation.sql"
if [[ ! -f "$migration" ]]; then
  echo "ERROR: missing P02.03 kernel.organization migration: ${migration}" >&2
  exit 1
fi

for marker in \
  'omnexa_organization.nodes' \
  'omnexa_organization.scoped_memberships' \
  'organization_parent_same_tenant' \
  'organization_membership_scope_same_tenant' \
  "node_kind IN ('organization', 'legal_entity', 'business_unit', 'branch', 'team', 'location')"; do
  if ! grep -Fq "$marker" "$migration"; then
    echo "ERROR: P02.03 organization migration missing required hierarchy/isolation marker: ${marker}" >&2
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

integration_url="${P02_03_TEST_DATABASE_URL:-}"
P02_03_TEST_DATABASE_URL= go test ./kernel/...
P02_03_TEST_DATABASE_URL= go test -race ./kernel/internal/organization -count=1

if [[ -z "$integration_url" ]]; then
  echo "ERROR: P02_03_TEST_DATABASE_URL is required for canonical P02.03 migration/hierarchy-isolation evidence" >&2
  exit 1
fi
P02_03_TEST_DATABASE_URL="$integration_url" go test -v ./kernel/internal/organization -run '^TestPostgresOrganizationFoundationIntegration$' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/organization"|"$module_prefix/kernel/internal/tenancy"|"$module_prefix/kernel/internal/identity"|"$module_prefix/kernel/internal/database"|"$module_prefix/kernel/internal/config"|"$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P02.03 organization package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/organization)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(audit|cache|configuration|developer|jobs|observability|operations|storage)' kernel/internal/organization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.03 organization runtime source contains later-domain or unrelated kernel coupling" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Party|Person|Customer|Supplier|Employee|Role|Permission|Policy|Session|Credential|MFA|Passkey|ServiceAccount|APIKey|TenantSetting|Company|Workspace|Warehouse|Department|Store|Brand)([[:space:]]|$)' kernel/internal/organization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.03 declares unauthorized business/authentication/authorization/future hierarchy concepts" >&2
  exit 1
fi

if grep -nEi '\b(company_id|workspace_id|party_id|person_id|employee_id|role_id|permission_id|policy_id|session_id|password|password_hash|mfa_secret|passkey|api_key|setting_key)\b' "$migration"; then
  echo "ERROR: P02.03 migration contains unauthorized business/auth/authz/settings fields" >&2
  exit 1
fi

if grep -R -nE '(fmt\.(Print|Printf|Println)|log\.|slog\.)' kernel/internal/organization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.03 runtime organization code must not emit tenant/principal/scope context directly to ordinary logs" >&2
  exit 1
fi

if grep -R -nE 'func[[:space:]]+New(ScopedContext|Scope)' kernel/internal/organization --include='*.go'; then
  echo "ERROR: P02.03 relationship scope/context must not have caller-constructible exported constructors" >&2
  exit 1
fi

for marker in \
  'type ScopedContext struct' \
  'func (scoped ScopedContext) ScopeFor' \
  'target != scoped.scopeID' \
  'TenantContextResolver interface'; do
  if ! grep -RFq "$marker" kernel/internal/organization; then
    echo "ERROR: P02.03 scoped-context fail-closed marker missing: ${marker}" >&2
    exit 1
  fi
done

go build ./kernel/...

echo "P02.03 G0 governance/active-package/kernel.organization owner boundary: PASS"
echo "P02.03 G1 format/static/dependency/schema ownership: PASS"
echo "P02.03 G2 unit/race hierarchy-kind/transition/membership semantics: PASS"
echo "P02.03 G3 PostgreSQL hierarchy/traversal/scoped-context contract: PASS"
echo "P02.03 G4 fresh/idempotent/P02.02-upgrade/immutable-ledger migration evidence: PASS"
echo "P02.03 G5 same-tenant allow/cross-tenant parent+membership deny/non-authorizing scope context: PASS"
echo "P02.03 G6 cycle/stale-context/revoked-membership hierarchy resilience: PASS"
echo "P02.03 G7 build/package: PASS"
echo "P02.03 G8 pinned identity/tenancy/database dependency chain: PASS"

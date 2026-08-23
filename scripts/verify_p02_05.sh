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
echo "P02.05 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P02.05 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.05 requires pinned github.com/jackc/pgx/v5 v5.10.0, got ${pgx_version}" >&2
  exit 1
fi
echo "P02.05 google/uuid version: ${uuid_version}"
echo "P02.05 pgx version: ${pgx_version}"

migration="kernel/migrations/kernel.authorization/1_create_rbac_foundation.sql"
if [[ ! -f "$migration" ]]; then
  echo "ERROR: missing P02.05 kernel.authorization migration: ${migration}" >&2
  exit 1
fi
for marker in \
  'omnexa_authorization.permissions' \
  'omnexa_authorization.roles' \
  'omnexa_authorization.role_permissions' \
  'omnexa_authorization.role_assignments' \
  'authorization.role.read' \
  'authorization.role.manage' \
  'authorization.assignment.read' \
  'authorization.assignment.manage' \
  'authorization_role_organization_same_tenant' \
  'REFERENCES omnexa_identity.users(principal_id)'; do
  if ! grep -Fq "$marker" "$migration"; then
    echo "ERROR: P02.05 migration missing direct-RBAC marker: ${marker}" >&2
    exit 1
  fi
done

if grep -nEi 'CREATE TABLE[^;]*(policy|relationship|object|module|mfa|passkey|service_account|api_key)' "$migration"; then
  echo "ERROR: P02.05 migration introduces future relationship/policy/module/MFA/service-account scope" >&2
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

integration_url="${P02_05_TEST_DATABASE_URL:-}"
P02_05_TEST_DATABASE_URL= go test ./kernel/...
P02_05_TEST_DATABASE_URL= go test -race ./kernel/internal/authorization -count=1

if [[ -z "$integration_url" ]]; then
  echo "ERROR: P02_05_TEST_DATABASE_URL is required for canonical P02.05 migration/RBAC evidence" >&2
  exit 1
fi
P02_05_TEST_DATABASE_URL="$integration_url" go test -v ./kernel/internal/authorization -run '^TestPostgresRBACFoundationIntegration$' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/authorization"|\
    "$module_prefix/kernel/internal/audit"|\
    "$module_prefix/kernel/internal/failure"|\
    "$module_prefix/kernel/internal/identity"|\
    "$module_prefix/kernel/internal/tenancy"|\
    "$module_prefix/kernel/internal/organization") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P02.05 authorization package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/authorization)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(configuration|developer|jobs|observability|operations|storage)' kernel/internal/authorization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.05 runtime authorization source contains unrelated or future-domain coupling" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Policy|Relationship|ObjectPolicy|ContextRule|MFA|Passkey|ServiceAccount|APIKey|ModulePermission|Customer|Supplier|Employee)([[:space:]]|$)' kernel/internal/authorization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.05 declares unauthorized P02.06+/P03/business authorization concepts" >&2
  exit 1
fi

if grep -R -nE 'func \(repository \*PostgresRepository\) (CreateRole|ReplaceRolePermissions|CreateAssignment|RevokeAssignment)' kernel/internal/authorization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.05 persistence mutation bypass is exported outside authorization service" >&2
  exit 1
fi

for marker in \
  'type repository interface' \
  'func SubjectFromTenantContext' \
  'func SubjectFromOrganizationContext' \
  'func (scope Scope) Equal' \
  'func (service *Service) Check' \
  'func (service *Service) Require' \
  'func (service *Service) CreateRole' \
  'func (service *Service) ReplaceRolePermissions' \
  'func (service *Service) AssignRole' \
  'func (service *Service) RevokeAssignment' \
  'service.ensureGrantable' \
  'audit.RequirementRequired' \
  'DecisionDeny'; do
  if ! grep -RFq "$marker" kernel/internal/authorization; then
    echo "ERROR: P02.05 fail-closed direct-RBAC marker missing: ${marker}" >&2
    exit 1
  fi
done

if ! grep -Fq "r.organization_id IS NOT DISTINCT FROM \$3::uuid" kernel/internal/authorization/postgres_repository.go; then
  echo "ERROR: P02.05 exact tenant/organization scope predicate is missing" >&2
  exit 1
fi

if grep -R -nE '(role\.Name\(\)[[:space:]]*==|strings\.(EqualFold|Contains)\([^\n]*role)' kernel/internal/authorization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.05 role-name shortcut detected" >&2
  exit 1
fi

go build ./kernel/...

echo "P02.05 G0 governance/active-package/kernel.authorization owner boundary: PASS"
echo "P02.05 G1 format/static/dependency/schema/service-boundary ownership: PASS"
echo "P02.05 G2 unit/race permission/role/direct-decision semantics: PASS"
echo "P02.05 G3 PostgreSQL direct-RBAC role/assignment lifecycle contract: PASS"
echo "P02.05 G4 fresh/idempotent/P02.04-prerequisite/immutable-ledger migration evidence: PASS"
echo "P02.05 G5 deny-by-default/anti-escalation/exact-scope/audit-safety negatives: PASS"
echo "P02.05 G6 assignment revocation/unassigned/role-name non-bypass resilience: PASS"
echo "P02.05 G7 build/package: PASS"
echo "P02.05 G8 pinned identity/tenancy/organization/audit/database dependency chain: PASS"

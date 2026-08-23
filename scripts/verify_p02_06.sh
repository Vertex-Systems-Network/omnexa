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
echo "P02.06 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P02.06 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.06 requires pinned github.com/jackc/pgx/v5 v5.10.0, got ${pgx_version}" >&2
  exit 1
fi
echo "P02.06 google/uuid version: ${uuid_version}"
echo "P02.06 pgx version: ${pgx_version}"

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

integration_url="${P02_06_TEST_DATABASE_URL:-}"
P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= go test ./kernel/...
P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= go test -race ./kernel/internal/authorization -count=1

if [[ -z "$integration_url" ]]; then
  echo "ERROR: P02_06_TEST_DATABASE_URL is required for canonical P02.06 authorization evidence" >&2
  exit 1
fi
P02_05_TEST_DATABASE_URL="$integration_url" P02_06_TEST_DATABASE_URL= \
  go test -v ./kernel/internal/authorization -run '^TestPostgresRBACFoundationIntegration$' -count=1
P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL="$integration_url" \
  go test -v ./kernel/internal/authorization -run '^TestPostgresContextAwareAuthorizationIntegration$' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/audit"|\
    "$module_prefix/kernel/internal/failure"|\
    "$module_prefix/kernel/internal/identity"|\
    "$module_prefix/kernel/internal/tenancy"|\
    "$module_prefix/kernel/internal/organization") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P02.06 authorization package directly imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -f '{{join .Imports "\n"}}' ./kernel/internal/authorization)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(configuration|developer|jobs|observability|operations|storage)' kernel/internal/authorization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.06 runtime authorization source contains unrelated/future-domain coupling" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(MFA|Passkey|ServiceAccount|APIKey|ModulePermission|Customer|Supplier|Employee)([[:space:]]|$)' kernel/internal/authorization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.06 declares unauthorized P02.07+/P03/business concepts" >&2
  exit 1
fi

for marker in \
  'type RelationshipResolver interface' \
  'type ContextConstraintEvaluator interface' \
  'type ContextService struct' \
  'func NewContextService' \
  'func (service *ContextService) Check' \
  'service.rbac.Check' \
  'AccessSensitiveField' \
  'AccessExport' \
  'CallerInternal' \
  'CallerBackground' \
  'RelationshipQuery' \
  'RelationshipEvidence' \
  'audit.RequirementRequired' \
  'DecisionDeny'; do
  if ! grep -RFq "$marker" kernel/internal/authorization; then
    echo "ERROR: P02.06 fail-closed contextual authorization marker missing: ${marker}" >&2
    exit 1
  fi
done

if grep -R -nE '(role\.Name\(\)[[:space:]]*==|strings\.(EqualFold|Contains)\([^\n]*role)' kernel/internal/authorization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.06 role-name shortcut detected" >&2
  exit 1
fi

if grep -R -nE 'if[[:space:]].*(CallerInternal|CallerBackground).*(DecisionAllow|return[[:space:]]+DecisionAllow)' kernel/internal/authorization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.06 internal/background caller bypass detected" >&2
  exit 1
fi

if find kernel/migrations/kernel.authorization -maxdepth 1 -type f -name '2_*' -print | grep -q .; then
  echo "ERROR: P02.06 unexpectedly added authorization persistence; this implementation is a dependency-inverted policy layer" >&2
  exit 1
fi

go build ./kernel/...

echo "P02.06 G0 governance/active-package/kernel.authorization owner boundary: PASS"
echo "P02.06 G1 format/static/direct-dependency/no-future-scope/service-layer ownership: PASS"
echo "P02.06 G2 unit/race RBAC+relationship+context/caller-origin/field-export semantics: PASS"
echo "P02.06 G3 PostgreSQL-backed RBAC plus contextual authorization integration: PASS"
echo "P02.06 G4 data/migration: N/A — no P02.06 persistence change; retained P02.05 migration regression PASS"
echo "P02.06 G5 same-scope allow/wrong-tenant/wrong-org/wrong-object/missing-permission/audit negatives: PASS"
echo "P02.06 G6 resolver/constraint fail-closed and internal/background non-bypass resilience: PASS"
echo "P02.06 G7 build/package: PASS"
echo "P02.06 G8 pinned identity/tenancy/organization/authorization/audit dependency chain: PASS"

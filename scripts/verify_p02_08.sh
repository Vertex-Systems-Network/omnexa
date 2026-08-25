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

echo "P02.08 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P02.08 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
if [[ "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.08 requires pinned github.com/jackc/pgx/v5 v5.10.0, got ${pgx_version}" >&2
  exit 1
fi

identity_migration="kernel/migrations/kernel.identity/4_create_service_accounts_api_credentials.sql"
authorization_migration="kernel/migrations/kernel.authorization/2_allow_service_account_assignments.sql"
for migration in "$identity_migration" "$authorization_migration"; do
  if [[ ! -f "$migration" ]]; then
    echo "ERROR: missing P02.08 migration: ${migration}" >&2
    exit 1
  fi
done

for marker in \
  'omnexa_identity.service_accounts' \
  'omnexa_identity.api_credentials' \
  'secret_digest' \
  'superseded_at' \
  'revoked_at' \
  'tenant_id' \
  'organization_id'; do
  if ! grep -Fq "$marker" "$identity_migration"; then
    echo "ERROR: P02.08 identity migration missing marker: ${marker}" >&2
    exit 1
  fi
done
for marker in \
  'authorization_assignment_principal_fk' \
  'REFERENCES omnexa_identity.principals(id)'; do
  if ! grep -Fq "$marker" "$authorization_migration"; then
    echo "ERROR: P02.08 authorization migration missing principal-composition marker: ${marker}" >&2
    exit 1
  fi
done

if grep -nEi '\b(raw_secret|credential_value|api_secret|private_key|oauth_client_secret|master_key|superkey)\b' "$identity_migration"; then
  echo "ERROR: P02.08 migration contains reversible/raw credential material or master-key scope" >&2
  exit 1
fi
if grep -R -nE 'type[[:space:]]+(OAuth|OAuthApplication|DeveloperApplication|Connector|DeviceIdentity|AIAgent|Customer|Supplier|Employee)([[:space:]]|$)' \
  kernel/internal/identity kernel/internal/authorization --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.08 pulled OAuth/developer/integration/device/AI-agent/business scope forward" >&2
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

integration_url="${P02_08_TEST_DATABASE_URL:-}"
P02_01_TEST_DATABASE_URL= P02_02_TEST_DATABASE_URL= P02_03_TEST_DATABASE_URL= P02_04_TEST_DATABASE_URL= P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= P02_07_TEST_DATABASE_URL= P02_08_TEST_DATABASE_URL= \
  go test ./kernel/...
P02_08_TEST_DATABASE_URL= go test -race ./kernel/internal/identity ./kernel/internal/authorization -count=1

if [[ -z "$integration_url" ]]; then
  echo "ERROR: P02_08_TEST_DATABASE_URL is required for canonical P02.08 evidence" >&2
  exit 1
fi
P02_01_TEST_DATABASE_URL= P02_02_TEST_DATABASE_URL= P02_03_TEST_DATABASE_URL= P02_04_TEST_DATABASE_URL= P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= P02_07_TEST_DATABASE_URL= P02_08_TEST_DATABASE_URL="$integration_url" \
  go test -v ./kernel/internal/authorization -run '^TestPostgresServiceAccountCredentialIntegration$' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/database"|\
    "$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P02.08 identity package directly imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -f '{{join .Imports "\n"}}' ./kernel/internal/identity)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(authorization|tenancy|organization|audit|cache|configuration|developer|jobs|observability|operations|storage)' \
  kernel/internal/identity --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.08 identity runtime crossed its owner/dependency boundary" >&2
  exit 1
fi

if grep -R -nE '(fmt\.(Print|Printf|Println)|log\.|slog\.)' \
  kernel/internal/identity/service_account_*.go kernel/internal/identity/postgres_service_account_repository.go; then
  echo "ERROR: P02.08 service-account runtime must not log RESTRICTED credential material" >&2
  exit 1
fi
if grep -nE '\.Reveal\(\)' \
  kernel/internal/identity/service_account_service.go \
  kernel/internal/identity/service_account_repository.go \
  kernel/internal/identity/postgres_service_account_repository.go \
  kernel/internal/identity/service_account_errors.go; then
  echo "ERROR: P02.08 persistence/service/error path reveals raw API credentials" >&2
  exit 1
fi

for marker in \
  'type ServiceAccount struct' \
  'type APICredentialSecret struct' \
  'func (secret APICredentialSecret) MarshalJSON' \
  'func (service *ServiceAccountService) IssueCredential' \
  'func (service *ServiceAccountService) VerifyCredential' \
  'func (service *ServiceAccountService) RotateCredential' \
  'func (service *ServiceAccountService) RevokeCredential' \
  'func ServiceAccountSubjectFromAuthentication' \
  'func (service *Service) CheckServiceAccount' \
  'func (service *Service) AssignRoleToServiceAccount'; do
  if ! grep -RFq "$marker" kernel/internal; then
    echo "ERROR: P02.08 fail-closed service-account marker missing: ${marker}" >&2
    exit 1
  fi
done

if grep -R -nE '(oauth|client_secret|master[_-]?token|super[_-]?key)' \
  kernel/internal/identity/service_account_*.go kernel/internal/authorization/service_account.go --exclude='*_test.go'; then
  echo "ERROR: P02.08 runtime contains unauthorized OAuth/master-token scope" >&2
  exit 1
fi

go build ./kernel/...

echo "P02.08 G0 governance/active-package/service-principal owner boundary: PASS"
echo "P02.08 G1 format/static/direct-dependency/no-future-scope/secret-redaction: PASS"
echo "P02.08 G2 unit/race service lifecycle/credential/redaction semantics: PASS"
echo "P02.08 G3 PostgreSQL credential lifecycle plus direct-RBAC composition: PASS"
echo "P02.08 G4 fresh/idempotent/P02.07+P02.05 supported-upgrade migration evidence: PASS"
echo "P02.08 G5 wrong-tenant/wrong-org/wrong-permission/revoked/expired/superseded negatives: PASS"
echo "P02.08 G6 current-principal/current-tenant/revocation/rotation fail-closed resilience: PASS"
echo "P02.08 G7 build/package: PASS"
echo "P02.08 G8 pinned identity/database/authorization dependency chain: PASS"

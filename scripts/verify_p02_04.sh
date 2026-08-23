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
echo "P02.04 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P02.04 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.04 requires pinned github.com/jackc/pgx/v5 v5.10.0, got ${pgx_version}" >&2
  exit 1
fi
echo "P02.04 google/uuid version: ${uuid_version}"
echo "P02.04 pgx version: ${pgx_version}"

migration="kernel/migrations/kernel.identity/2_create_authentication_sessions.sql"
if [[ ! -f "$migration" ]]; then
  echo "ERROR: missing P02.04 kernel.identity migration: ${migration}" >&2
  exit 1
fi

for marker in \
  'omnexa_identity.password_credentials' \
  'omnexa_identity.sessions' \
  'omnexa_identity.access_credentials' \
  'omnexa_identity.refresh_credentials' \
  'password_hash' \
  'secret_digest' \
  'credential_version' \
  'tenant_context_hint' \
  'organization_context_hint'; do
  if ! grep -Fq "$marker" "$migration"; then
    echo "ERROR: P02.04 migration missing required authentication/session marker: ${marker}" >&2
    exit 1
  fi
done

if grep -nEi '\b(access_secret|refresh_secret|access_token|refresh_token|role_id|permission_id|policy_id|mfa_secret|passkey|api_key|service_account_id)\b' "$migration"; then
  echo "ERROR: P02.04 migration contains raw secret or future authz/MFA/service-account fields" >&2
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

integration_url="${P02_04_TEST_DATABASE_URL:-}"
P02_04_TEST_DATABASE_URL= go test ./kernel/...
P02_04_TEST_DATABASE_URL= go test -race ./kernel/internal/identity -count=1

if [[ -z "$integration_url" ]]; then
  echo "ERROR: P02_04_TEST_DATABASE_URL is required for canonical P02.04 migration/session evidence" >&2
  exit 1
fi
P02_04_TEST_DATABASE_URL="$integration_url" go test -v ./kernel/internal/identity -run '^TestPostgresAuthenticationSessionLifecycleIntegration$' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/identity"|"$module_prefix/kernel/internal/database"|"$module_prefix/kernel/internal/config"|"$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P02.04 identity package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/identity)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(authorization|tenancy|organization|audit|cache|configuration|developer|jobs|observability|operations|storage)' kernel/internal/identity --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.04 runtime identity source contains future-domain or unrelated kernel coupling" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Role|Permission|Policy|MFA|Passkey|ServiceAccount|APIKey|TenantSetting|Party|Person|Customer|Supplier|Employee)([[:space:]]|$)' kernel/internal/identity --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.04 declares unauthorized authorization/MFA/service-account/settings/business concepts" >&2
  exit 1
fi

if grep -R -nE '(fmt\.(Print|Printf|Println)|log\.|slog\.)' kernel/internal/identity --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.04 runtime identity code must not emit authentication material to ordinary logs" >&2
  exit 1
fi

if grep -R -nE '\.(Reveal\(\))' kernel/internal/identity/postgres_authentication_repository.go kernel/internal/identity/errors.go; then
  echo "ERROR: P02.04 persistence/error paths must never reveal raw session credentials" >&2
  exit 1
fi

for marker in \
  'type ContextReauthorizer interface' \
  'type AccessSecret struct' \
  'type RefreshSecret struct' \
  'func (secret AccessSecret) MarshalJSON' \
  'func (secret RefreshSecret) MarshalJSON' \
  'func (service *AuthenticationService) IssueSession' \
  'func (service *AuthenticationService) ValidateAccess' \
  'func (service *AuthenticationService) RotateRefresh' \
  'service.reauthorize' \
  'credentialVersion'; do
  if ! grep -RFq "$marker" kernel/internal/identity; then
    echo "ERROR: P02.04 authentication/session fail-closed marker missing: ${marker}" >&2
    exit 1
  fi
done

if ! grep -Fq 'passwordIterations      = 600000' kernel/internal/identity/password.go; then
  echo "ERROR: P02.04 governed PBKDF2 work factor marker is missing" >&2
  exit 1
fi

go build ./kernel/...

echo "P02.04 G0 governance/active-package/kernel.identity owner boundary: PASS"
echo "P02.04 G1 format/static/dependency/schema/secret-storage ownership: PASS"
echo "P02.04 G2 unit/race password/session/opaque-secret semantics: PASS"
echo "P02.04 G3 PostgreSQL authentication/session lifecycle contract: PASS"
echo "P02.04 G4 fresh/idempotent/P02.01-upgrade/immutable-ledger migration evidence: PASS"
echo "P02.04 G5 disclosure-safe authentication/current-context reauthorization/no-future-authz scope: PASS"
echo "P02.04 G6 refresh replay/revocation/password-change/account-lifecycle invalidation resilience: PASS"
echo "P02.04 G7 build/package: PASS"
echo "P02.04 G8 pinned identity/database dependency chain: PASS"

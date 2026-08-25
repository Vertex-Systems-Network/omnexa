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
echo "P02.07 Go toolchain: ${actual_go}"

uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
if [[ "$uuid_version" != "v1.6.0" ]]; then
  echo "ERROR: P02.07 requires pinned github.com/google/uuid v1.6.0, got ${uuid_version}" >&2
  exit 1
fi
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.07 requires pinned github.com/jackc/pgx/v5 v5.10.0, got ${pgx_version}" >&2
  exit 1
fi
echo "P02.07 google/uuid version: ${uuid_version}"
echo "P02.07 pgx version: ${pgx_version}"

migration="kernel/migrations/kernel.identity/3_create_strong_authentication.sql"
if [[ ! -f "$migration" ]]; then
  echo "ERROR: missing P02.07 kernel.identity migration: ${migration}" >&2
  exit 1
fi
for marker in \
  'omnexa_identity.mfa_factors' \
  'omnexa_identity.passkey_credentials' \
  'omnexa_identity.authentication_challenges' \
  'omnexa_identity.recovery_code_sets' \
  'omnexa_identity.recovery_codes' \
  'secret_digest' \
  'counter_supported' \
  'identity_sessions_id_principal_unique'; do
  if ! grep -Fq "$marker" "$migration"; then
    echo "ERROR: P02.07 migration missing strong-auth marker: ${marker}" >&2
    exit 1
  fi
done

if grep -nEi '\b(challenge_secret|raw_challenge|raw_recovery_code|private_key|authenticator_secret|api_key|api_secret|service_account_id|saml_assertion|sso_provider_id|tenant_setting_id)\b' "$migration"; then
  echo "ERROR: P02.07 migration contains raw authentication secret or future P02.08+/P24/settings field" >&2
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

integration_url="${P02_07_TEST_DATABASE_URL:-}"
P02_01_TEST_DATABASE_URL= P02_02_TEST_DATABASE_URL= P02_03_TEST_DATABASE_URL= P02_04_TEST_DATABASE_URL= P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= P02_07_TEST_DATABASE_URL= \
  go test ./kernel/...
P02_07_TEST_DATABASE_URL= go test -race ./kernel/internal/identity -count=1

if [[ -z "$integration_url" ]]; then
  echo "ERROR: P02_07_TEST_DATABASE_URL is required for canonical P02.07 strong-authentication evidence" >&2
  exit 1
fi
P02_01_TEST_DATABASE_URL= P02_02_TEST_DATABASE_URL= P02_03_TEST_DATABASE_URL= P02_04_TEST_DATABASE_URL= P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= P02_07_TEST_DATABASE_URL="$integration_url" \
  go test -v ./kernel/internal/identity -run '^TestPostgresStrongAuthenticationIntegration$' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/database"|\
    "$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P02.07 identity package directly imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -f '{{join .Imports "\n"}}' ./kernel/internal/identity)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(authorization|tenancy|organization|audit|cache|configuration|developer|jobs|observability|operations|storage)' kernel/internal/identity --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.07 runtime identity source contains unrelated/future-owner coupling" >&2
  exit 1
fi

# P02.08 service-account/API-credential types and migration are now canonically
# active. Keep the historical P02.07 guard only on concepts that remain future.
if grep -R -nE 'type[[:space:]]+(TenantSetting|SSO|SAML|SCIM|Party|Person|Customer|Supplier|Employee)([[:space:]]|$)' kernel/internal/identity --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.07 declares unauthorized P02.09+/P24/settings/business concepts" >&2
  exit 1
fi

if grep -R -nE '"crypto/(aes|cipher|des|dsa|ecdsa|ed25519|elliptic|hmac|rsa|x509)"' \
  kernel/internal/identity/strong_auth_*.go kernel/internal/identity/postgres_strong_auth_repository.go; then
  echo "ERROR: P02.07 strong-auth runtime contains custom protocol/private-key cryptography instead of injected approved verification" >&2
  exit 1
fi

if grep -R -nE '(fmt\.(Print|Printf|Println)|log\.|slog\.)' kernel/internal/identity --include='*.go' --exclude='*_test.go'; then
  echo "ERROR: P02.07 runtime identity code must not emit RESTRICTED authentication material to ordinary logs" >&2
  exit 1
fi

if grep -nE '\.Reveal\(\)' kernel/internal/identity/postgres_strong_auth_repository.go kernel/internal/identity/strong_auth_errors.go; then
  echo "ERROR: P02.07 persistence/error paths must never reveal raw factor/challenge/recovery material" >&2
  exit 1
fi

for marker in \
  'type PasskeyVerifier interface' \
  'type StrongAuthenticationService struct' \
  'func (service *StrongAuthenticationService) BeginPasskeyEnrollment' \
  'func (service *StrongAuthenticationService) CompletePasskeyEnrollment' \
  'func (service *StrongAuthenticationService) BeginPasskeyAssertion' \
  'func (service *StrongAuthenticationService) CompletePasskeyAssertion' \
  'func (service *StrongAuthenticationService) IssueRecoveryCodes' \
  'func (service *StrongAuthenticationService) UseRecoveryCode' \
  'func (service *StrongAuthenticationService) RemoveFactor' \
  'func (service *StrongAuthenticationService) RequireStepUp' \
  'InvalidateSessionsOnFactorRemoval' \
  'func (secret ChallengeSecret) MarshalJSON' \
  'func (code RecoveryCode) MarshalJSON'; do
  if ! grep -RFq "$marker" kernel/internal/identity; then
    echo "ERROR: P02.07 fail-closed strong-authentication marker missing: ${marker}" >&2
    exit 1
  fi
done

# Direct imports of kernel.authorization are rejected above. This guard separately
# rejects authorization-domain symbols in strong-auth runtime without treating
# explanatory comments containing the word "authorization" as executable authority.
if grep -R -nE '(DecisionAllow|PermissionID|RoleID|ContextService)' kernel/internal/identity/strong_auth_*.go kernel/internal/identity/postgres_strong_auth_repository.go; then
  echo "ERROR: P02.07 strong authentication is attempting to become authorization authority" >&2
  exit 1
fi

go build ./kernel/...

echo "P02.07 G0 governance/active-package/kernel.identity human-user strong-auth owner boundary: PASS"
echo "P02.07 G1 format/static/direct-dependency/no-future-scope/no-custom-protocol-crypto: PASS"
echo "P02.07 G2 unit/race enrollment/assertion/recovery/step-up/redaction semantics: PASS"
echo "P02.07 G3 PostgreSQL strong-authentication lifecycle and P02.04 session integration: PASS"
echo "P02.07 G4 fresh/idempotent/P02.04-upgrade/immutable-ledger migration evidence: PASS"
echo "P02.07 G5 wrong-principal/wrong-session/expiry/replay/recovery/audit negatives: PASS"
echo "P02.07 G6 verifier-failure/counter/session-invalidation fail-closed resilience: PASS"
echo "P02.07 G7 build/package: PASS"
echo "P02.07 G8 pinned identity/database/failure dependency chain: PASS"

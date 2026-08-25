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
uuid_version="$(go list -m -f '{{.Version}}' github.com/google/uuid)"
pgx_version="$(go list -m -f '{{.Version}}' github.com/jackc/pgx/v5)"
if [[ "$uuid_version" != "v1.6.0" || "$pgx_version" != "v5.10.0" ]]; then
  echo "ERROR: P02.10 pinned dependency mismatch: uuid=${uuid_version}, pgx=${pgx_version}" >&2
  exit 1
fi

echo "P02.10 Go toolchain: ${actual_go}"
echo "P02.10 google/uuid version: ${uuid_version}"
echo "P02.10 pgx version: ${pgx_version}"

# P02.10 integrates audit producers and exit proof; it introduces no new durable
# schema. The aggregate integration test replays the complete accepted P02.09
# migration baseline fresh/idempotently as the supported no-op upgrade proof.
if find kernel/migrations -maxdepth 1 -type d -name 'kernel.audit' | grep -q .; then
  echo "ERROR: P02.10 introduced unauthorized audit persistence instead of using P01.11 transport" >&2
  exit 1
fi

for verifier in 01 02 03 04 05 06 07 08 09; do
  script="scripts/verify_p02_${verifier}.sh"
  if [[ ! -f "$script" ]]; then
    echo "ERROR: missing completed P02 regression verifier: ${script}" >&2
    exit 1
  fi
  if ! grep -Fq "bash ${script}" .github/workflows/governance.yml; then
    echo "ERROR: governance workflow does not compose completed P02 verifier: ${script}" >&2
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

integration_url="${P02_10_TEST_DATABASE_URL:-}"
P02_01_TEST_DATABASE_URL= P02_02_TEST_DATABASE_URL= P02_03_TEST_DATABASE_URL= P02_04_TEST_DATABASE_URL= P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= P02_07_TEST_DATABASE_URL= P02_08_TEST_DATABASE_URL= P02_09_TEST_DATABASE_URL= P02_10_TEST_DATABASE_URL= \
  go test ./kernel/...
P02_10_TEST_DATABASE_URL= go test -race ./kernel/internal/audit ./kernel/internal/identity ./kernel/internal/tenancy ./kernel/internal/organization ./kernel/internal/authorization ./kernel/internal/configuration -count=1

go test -v ./kernel/internal/identity -run '^TestAuditAdapter' -count=1
go test -v ./kernel/internal/tenancy -run '^TestAuditedRepository' -count=1

if [[ -z "$integration_url" ]]; then
  echo "ERROR: P02_10_TEST_DATABASE_URL is required for canonical P02 aggregate exit evidence" >&2
  exit 1
fi
P02_01_TEST_DATABASE_URL= P02_02_TEST_DATABASE_URL= P02_03_TEST_DATABASE_URL= P02_04_TEST_DATABASE_URL= P02_05_TEST_DATABASE_URL= P02_06_TEST_DATABASE_URL= P02_07_TEST_DATABASE_URL= P02_08_TEST_DATABASE_URL= P02_09_TEST_DATABASE_URL= P02_10_TEST_DATABASE_URL="$integration_url" \
  go test -v ./kernel/internal/organization -run '^TestPostgresP02AuditAndExitIntegration$' -count=1

for marker in \
  'type AuditAdapter struct' \
  'func (adapter *AuditAdapter) RecordSecurityEvent' \
  'func (adapter *AuditAdapter) RecordServiceAccountEvent'; do
  if ! grep -Fq "$marker" kernel/internal/identity/audit_adapter.go; then
    echo "ERROR: P02.10 identity audit bridge marker missing: ${marker}" >&2
    exit 1
  fi
done
for file in kernel/internal/tenancy/audited_repository.go kernel/internal/organization/audited_repository.go; do
  for marker in 'type AuditedRepository struct' 'audit.RequirementRequired' 'audit.OutcomeSucceeded'; do
    if ! grep -Fq "$marker" "$file"; then
      echo "ERROR: P02.10 required-audit mutation marker missing from ${file}: ${marker}" >&2
      exit 1
    fi
  done
done

if grep -R -nE '(fmt\.(Print|Printf|Println)|log\.|slog\.)' \
  kernel/internal/identity/audit_adapter.go \
  kernel/internal/tenancy/audited_repository.go \
  kernel/internal/organization/audited_repository.go; then
  echo "ERROR: P02.10 audit integration must not emit security metadata to ordinary logs" >&2
  exit 1
fi
if grep -R -nE '\.Reveal\(\)|raw_(secret|token|credential)|password_hash|secret_digest|private_key|recovery_code' \
  kernel/internal/identity/audit_adapter.go \
  kernel/internal/tenancy/audited_repository.go \
  kernel/internal/organization/audited_repository.go; then
  echo "ERROR: P02.10 audit integration references prohibited credential/secret material" >&2
  exit 1
fi
if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(jobs|storage|operations)' \
  kernel/internal/identity/audit_adapter.go \
  kernel/internal/tenancy/audited_repository.go \
  kernel/internal/organization/audited_repository.go; then
  echo "ERROR: P02.10 pulled unrelated/future owner scope into audit integration" >&2
  exit 1
fi
if grep -R -nE 'type[[:space:]]+(Module|Workflow|Event|Agent|Model|Customer|Supplier|Employee|Impersonation)([[:space:]]|$)' \
  kernel/internal/identity/audit_adapter.go \
  kernel/internal/tenancy/audited_repository.go \
  kernel/internal/organization/audited_repository.go; then
  echo "ERROR: P02.10 pulled P03/P04/business/AI/support-impersonation scope forward" >&2
  exit 1
fi

for marker in \
  'TestPostgresP02AuditAndExitIntegration' \
  'cross-tenant organization create unexpectedly succeeded' \
  'required audit failure unexpectedly claimed organization mutation success' \
  'kernel.identity' \
  'kernel.tenancy' \
  'kernel.organization' \
  'kernel.authorization' \
  'kernel.configuration'; do
  if ! grep -Fq "$marker" kernel/internal/organization/p02_exit_integration_test.go; then
    echo "ERROR: P02.10 aggregate exit-proof marker missing: ${marker}" >&2
    exit 1
  fi
done

go build ./kernel/...

echo "P02.10 G0 governance/sole-active-package/P02-exit boundary: PASS"
echo "P02.10 G1 format/static/owner-dependency/no-future-scope/no-secret audit integration: PASS"
echo "P02.10 G2 unit/race identity+tenancy+organization audit semantics: PASS"
echo "P02.10 G3 aggregate PostgreSQL audit producer and completed-P02 contract composition: PASS"
echo "P02.10 G4 complete P02 fresh/idempotent/P02.09 no-op upgrade migration baseline: PASS"
echo "P02.10 G5 same-tenant/cross-tenant/required-audit plus mandatory P02.01-P02.09 security regressions: PASS"
echo "P02.10 G6 audit-failure/session/service-account/revocation resilience via aggregate + retained regressions: PASS"
echo "P02.10 G7 build/package: PASS"
echo "P02.10 G8 pinned dependencies and canonical GitHub-hosted regression composition: PASS"

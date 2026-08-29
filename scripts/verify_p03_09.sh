#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

python scripts/validate_governance.py
python scripts/validate_p03_preparation.py
python scripts/validate_p03_package_specs.py

expected_go="$(tr -d '[:space:]' < .go-version)"
actual_go="$(go env GOVERSION)"
if [[ "$actual_go" != "go${expected_go}" ]]; then
  echo "ERROR: Go toolchain mismatch: got ${actual_go}, want go${expected_go}" >&2
  exit 1
fi

p03_09_files=(
  kernel/internal/modules/migration_ownership_registry.go
  kernel/internal/modules/migration_ownership_registry_test.go
  kernel/internal/database/migration.go
  kernel/migrations/README.md
  docs/roadmap/work-packages/P03.09.md
)
for file in "${p03_09_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.09 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/modules/*.go kernel/internal/database/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P03.09 Go files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P03.09 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

if grep -nE 'database/sql|jackc/pgx|kernel/internal/database|(^|[[:space:]])SQL[[:space:]]+string|(^|[[:space:]])Path[[:space:]]+string|(TenantID|OrganizationID)[[:space:]]+string|(Handler|Callback)[[:space:]]+(func|any|interface)|func[[:space:]].*(Apply|Execute|RunMigration)' kernel/internal/modules/migration_ownership_registry.go; then
  echo "ERROR: P03.09 introduced SQL execution, raw path/scope, callback, or duplicate migration-engine authority" >&2
  exit 1
fi

for marker in \
  'type MigrationChangeClass string' \
  'MigrationCompatible' \
  'MigrationBackfill' \
  'MigrationDestructive' \
  'type MigrationRegistration struct' \
  'type MigrationRecord struct' \
  'type MigrationOwnershipRegistry struct' \
  'func BindMigrationOwnershipRegistry(' \
  'func (r *MigrationOwnershipRegistry) FreshInstallPlan(' \
  'func (r *MigrationOwnershipRegistry) UpgradePlan(' \
  'module.migration.owner_mismatch' \
  'module.migration.target_owner_mismatch' \
  'module.migration.order_conflict' \
  'module.migration.registration_duplicate' \
  'module.migration.strategy_required' \
  'module.migration.introduced_version_order_invalid' \
  'module.migration.upgrade_future_source' \
  'TestMigrationOwnershipRegistryDeterministicFreshInstallAndUpgradePlan' \
  'TestMigrationOwnershipRegistryRejectsIdentityOwnerOrderAndVersionConflicts' \
  'TestMigrationOwnershipRegistryRequiresExplicitBackfillAndDestructiveStrategy' \
  'TestMigrationOwnershipRegistryUpgradePlanningFailsClosed' \
  'TestMigrationOwnershipRegistryMetadataSurfaceHasNoExecutionOrRawScopeFields'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: P03.09 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

# P01 remains the execution authority: immutable checksum ledger, owner-scoped
# advisory lock and transaction boundary must still exist unchanged.
for marker in \
  'schema_migrations' \
  'pg_advisory_lock' \
  'pg_advisory_unlock' \
  'migration.checksum()' \
  'database migration history drift was detected'; do
  if ! grep -Fq "$marker" kernel/internal/database/migration.go; then
    echo "ERROR: retained P01 migration invariant missing: ${marker}" >&2
    exit 1
  fi
done

if ! grep -Fq 'kernel/migrations/<owner>/<version>_<name>.sql' kernel/migrations/README.md; then
  echo "ERROR: retained P01 owner-directory migration convention missing" >&2
  exit 1
fi

# Retained P03.01-P03.08 verifiers execute as preceding canonical-governance
# steps. This focused verifier proves only the new P03.09 ownership/planning boundary.
go vet ./kernel/internal/modules ./kernel/internal/database
go test ./kernel/internal/modules ./kernel/internal/database -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.09 G0 active-package governance boundary: PASS"
echo "P03.09 G1 deterministic module/version/owner + P01 ledger identity binding: PASS"
echo "P03.09 G2 duplicate declaration and owner/version order conflict rejection: PASS"
echo "P03.09 G3 cross-owner target mutation metadata fails closed: PASS"
echo "P03.09 G4 compatible/backfill/destructive classification + strategy/recovery metadata: PASS"
echo "P03.09 G5 deterministic fresh-install and supported-upgrade planning: PASS"
echo "P03.09 G6 retry/double-apply authority retained by P01 checksum ledger + owner advisory lock: PASS"
echo "P03.09 G7 security: no SQL executor/raw path/tenant scope/secret/callback authority: PASS"
echo "P03.09 G8 build/race/package: PASS"
echo "P03.09 G9 P03.10+ and business-feature runtime remain absent: PASS"

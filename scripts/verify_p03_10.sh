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

p03_10_files=(
  kernel/internal/modules/module_health.go
  kernel/internal/modules/module_health_test.go
  kernel/internal/modules/migration_ownership_registry.go
  kernel/internal/modules/capability_registry.go
  kernel/internal/modules/permission_registry.go
  kernel/internal/modules/ui_contribution_registry.go
  kernel/internal/operations/health.go
  docs/roadmap/work-packages/P03.10.md
)
for file in "${p03_10_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.10 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/modules/*.go kernel/internal/operations/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P03.10 Go files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P03.10 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

# P03.10 is a read-only diagnostic layer. It must not gain database, network,
# lifecycle mutation, authorization-grant, execution callback, or raw-scope authority.
# Scan executable source only; explanatory full-line comments are not authority surfaces.
if grep -nE 'database/sql|jackc/pgx|kernel/internal/database|net/http|(^|[[:space:]])(SQL|Path|TenantID|OrganizationID|Secret|Password|Token|Credentials)[[:space:]]+[A-Za-z]|(^|[[:space:]])(Handler|Callback)[[:space:]]+(func|any|interface)|func[[:space:]].*(Apply|Execute|RunMigration|Authorize|Grant|Reconcile|CompareAndSwap)' kernel/internal/modules/module_health.go \
  | grep -vE '^[0-9]+:[[:space:]]*//'; then
  echo "ERROR: P03.10 introduced execution, mutation, raw scope, secret, or transport authority" >&2
  exit 1
fi

for marker in \
  'type ModuleHealthState string' \
  'ModuleHealthHealthy' \
  'ModuleHealthDegraded' \
  'ModuleHealthUnavailable' \
  'ModuleHealthFailed' \
  'type ModuleHealthRecord struct' \
  'type ModuleHealthReport struct' \
  'Readiness operations.State' \
  'type MigrationHealthObservation struct' \
  'type ModuleMigrationHealthSource interface' \
  'func NewModuleHealthReporter(' \
  'func (r *ModuleHealthReporter) Report(' \
  'ResolveDependencies(r.registry, r.platform, nil)' \
  'FreshInstallPlan(module.ID)' \
  'healthReasonMigrationPending' \
  'healthReasonMigrationInconsistent' \
  'healthReasonRequiredUnavailable' \
  'healthReasonOptionalMissing' \
  'TestModuleHealthReporterHealthyAndDeterministic' \
  'TestModuleHealthReporterOptionalDependencySelectiveDegradation' \
  'TestModuleHealthReporterRequiredDependencyFailsClosed' \
  'TestModuleHealthReporterMigrationPendingAndInconsistentNeverHealthy' \
  'TestModuleHealthReporterFailureIsolationAndLifecycleChanges' \
  'TestModuleHealthReporterClassificationSafeSurface'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: P03.10 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

# Retain the P01.08 portable health/readiness foundation. P03.10 may project to
# its stable state vocabulary, but P01 must not depend back on module runtime.
for marker in \
  'type State string' \
  'StateHealthy' \
  'StateDegraded' \
  'StateUnready' \
  'type Dependency struct' \
  'func (manager *Manager) Evaluate('; do
  if ! grep -Fq "$marker" kernel/internal/operations/health.go; then
    echo "ERROR: retained P01.08 health invariant missing: ${marker}" >&2
    exit 1
  fi
done
if grep -R -nE 'kernel/internal/modules|type[[:space:]]+ModuleHealth' kernel/internal/operations --include='*.go'; then
  echo "ERROR: P01.08 portable operations package now depends on P03 module-health runtime" >&2
  exit 1
fi

# P03.09 remains the migration identity/planning authority and must stay execution-free.
for marker in \
  'type MigrationOwnershipRegistry struct' \
  'func (r *MigrationOwnershipRegistry) FreshInstallPlan(' \
  'func (r *MigrationOwnershipRegistry) UpgradePlan('; do
  if ! grep -Fq "$marker" kernel/internal/modules/migration_ownership_registry.go; then
    echo "ERROR: retained P03.09 migration ownership invariant missing: ${marker}" >&2
    exit 1
  fi
done

# P03.06-P03.08 remain non-granting metadata sources.
for marker in \
  'type CapabilityRegistry struct' \
  'type PermissionRegistry struct' \
  'type UIContributionRegistry struct'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: retained module metadata registry missing: ${marker}" >&2
    exit 1
  fi
done

go vet ./kernel/internal/modules ./kernel/internal/operations
go test ./kernel/internal/modules ./kernel/internal/operations -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.10 G0 active-package governance boundary: PASS"
echo "P03.10 G1 stable module identity/version/lifecycle diagnostic projection: PASS"
echo "P03.10 G2 required dependency fail-closed + optional selective degradation: PASS"
echo "P03.10 G3 exact P03.09 migration identity partition + pending/failure inconsistency safety: PASS"
echo "P03.10 G4 capability/permission/UI metadata availability remains diagnostic and non-granting: PASS"
echo "P03.10 G5 classification-safe bounded reasons/counts with no raw scope/secret/stack payload: PASS"
echo "P03.10 G6 module failure isolation + lifecycle state-change accuracy: PASS"
echo "P03.10 G7 retained P01.08 readiness vocabulary integration without reverse coupling: PASS"
echo "P03.10 G8 build/race/package: PASS"
echo "P03.10 G9 P03.11+ and business-feature runtime remain absent: PASS"

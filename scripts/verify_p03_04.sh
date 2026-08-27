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

p03_04_files=(
  kernel/internal/modules/lifecycle.go
  kernel/internal/modules/lifecycle_test.go
  kernel/internal/modules/lifecycle_adversarial_test.go
)
for file in "${p03_04_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.04 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/modules/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P03 module files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P03.04 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

if grep -nE '"(os|os/exec|plugin|net/http|database/sql|path/filepath)"' "${p03_04_files[@]}"; then
  echo "ERROR: P03.04 state machine must not execute packages, filesystem/network operations, or own database access" >&2
  exit 1
fi
if grep -nE '(os\.(ReadDir|ReadFile|Open)|filepath\.(Walk|WalkDir|Glob)|exec\.Command|plugin\.Open|http\.(Get|Post)|http\.Client|sql\.Open)[[:space:]]*[({]' "${p03_04_files[@]}"; then
  echo "ERROR: P03.04 contains an unauthorized execution/discovery primitive" >&2
  exit 1
fi
if grep -nE 'type[[:space:]]+(CapabilityRegistry|PermissionRegistry|UIContributionRegistry|ModuleSettingsRegistry|PackageInstaller)([[:space:]]|$)' "${p03_04_files[@]}"; then
  echo "ERROR: P03.04 pulled P03.05+ or package-install scope forward" >&2
  exit 1
fi

for marker in \
  'type LifecycleState string' \
  'LifecycleInstalled' \
  'LifecycleEnabled' \
  'LifecycleDisabled' \
  'LifecycleSuspended' \
  'LifecycleArchived' \
  'LifecycleDetached' \
  'LifecycleRecoveryRequired' \
  'LifecyclePurged' \
  'type LifecycleStore interface' \
  'CompareAndSwap(' \
  'type LifecycleAuthorizer interface' \
  'type LifecycleAuditor interface' \
  'type LifecycleUpgradeCoordinator interface' \
  'func (m LifecycleManager) Apply(' \
  'func (m LifecycleManager) MarkRecoveryRequired(' \
  'lifecycle.reverse_dependency.active' \
  'lifecycle.reverse_dependency.present' \
  'lifecycle.reverse_dependency.version_mismatch' \
  'lifecycle.version.mismatch' \
  'lifecycle.dependency.version_mismatch' \
  'lifecycle.store.read_failed' \
  'lifecycle.authorization.denied' \
  'lifecycle.audit.failed' \
  'lifecycle.concurrent_conflict' \
  'lifecycle.upgrade.coordinator_unavailable' \
  'lifecycle.failure.state_invalid' \
  'TestLifecycleStateMachineSupportsNonDestructiveDisableAndReenable' \
  'TestLifecycleDependencyPreconditionsAndReverseProtection' \
  'TestLifecyclePurgeRequiresAuthorizationAuditAndDetachedState' \
  'TestLifecycleOperationReplayAndConcurrentConflictAreExplicit' \
  'TestLifecycleFailureAndRecoveryPreserveStableState' \
  'TestLifecycleUpgradeUsesFutureCoordinatorWithoutExecutingMigrations' \
  'TestLifecycleGraphFailureDoesNotMutateExistingState' \
  'TestLifecycleReplayStillRequiresCurrentAuthorization' \
  'TestLifecycleAuthorizationPrecedesUpgradeCoordinator' \
  'TestLifecycleEnableRejectsStaleInstalledModuleVersion' \
  'TestLifecycleDependencyRequiresInstalledVersionBoundToResolverRegistry' \
  'TestLifecycleProviderUpgradeRejectsStaleInstalledReverseDependent' \
  'TestLifecycleReverseDependencyReadFailureFailsClosed' \
  'TestLifecycleRecoveryTargetEnabledProtectsRequiredDependency' \
  'TestLifecycleRecoverToEnabledRechecksRequiredDependency' \
  'TestLifecycleFailureMarkerRejectsImpossibleSourceAction' \
  'TestLifecycleFailedInstallCanEnterAndRecoverFromRecoveryRequired' \
  'TestLifecycleFailedReinstallFromPurgedRecoversToPurged' \
  'TestLifecycleFailureFixturePreservesUnrelatedModuleIntegrity'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: P03.04 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

# Retained P03.01-P03.03 verifiers remain separate mandatory workflow steps.
go vet ./kernel/internal/modules
go test ./kernel/internal/modules -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.04 G0 active-package governance boundary: PASS"
echo "P03.04 G1 deterministic explicit lifecycle state/transition contract: PASS"
echo "P03.04 G2 lifecycle dependency/reverse-dependency/idempotency/recovery tests: PASS"
echo "P03.04 G3 registry/resolver lifecycle integration boundary: PASS"
echo "P03.04 G4 persistence: adapter boundary only; no new schema/migration introduced: N/A"
echo "P03.04 G5 security: authorization + required audit + destructive purge fail closed: PASS"
echo "P03.04 G6 lifecycle/resilience: CAS conflict, retry replay, recovery-required and isolation semantics: PASS"
echo "P03.04 G7 build/package: PASS"
echo "P03.04 G8 supply chain: stdlib-only implementation; module/workspace metadata unchanged: PASS"

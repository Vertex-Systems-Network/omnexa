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

p03_05_files=(
  kernel/internal/configuration/policy_validation.go
  kernel/internal/modules/configuration_binding.go
  kernel/internal/modules/configuration_binding_test.go
  kernel/internal/modules/manifest_versioned.go
  docs/roadmap/work-packages/P03.05.md
)
for file in "${p03_05_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.05 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/configuration/policy_validation.go kernel/internal/modules/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P03.05 Go files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P03.05 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

if grep -nE 'type[[:space:]]+(ModuleSettingsRegistry|FeatureFlagRegistry|CapabilityRegistry|PermissionRegistry|UIContributionRegistry|PackageInstaller)([[:space:]]|$)' kernel/internal/modules/configuration_binding.go; then
  echo "ERROR: P03.05 introduced a duplicate/future registry authority" >&2
  exit 1
fi
if grep -nE '(TenantID|OrganizationID)[[:space:]]+string' kernel/internal/modules/configuration_binding.go; then
  echo "ERROR: P03.05 introduced raw tenant/organization scope input" >&2
  exit 1
fi
if grep -nE 'ClassKillSwitch' kernel/internal/modules/configuration_binding.go; then
  echo "ERROR: P03.05 module feature declarations must not absorb operational kill-switch ownership" >&2
  exit 1
fi

for marker in \
  'type ModuleConfigurationRegistration struct' \
  'type ConfigurationBinding struct' \
  'func BindConfigurationRegistrations(' \
  'configuration.NewRegistry(' \
  'configuration.ValidateSettingPolicy(' \
  'ModuleConfigurationGlobal' \
  'ModuleConfigurationScoped' \
  'module.configuration.definition_missing' \
  'module.configuration.definition_undeclared' \
  'module.configuration.owner_mismatch' \
  'module.configuration.class_mismatch' \
  'module.configuration.scope_invalid' \
  'module.configuration.scope_policy_missing' \
  'module.configuration.scope_policy_invalid' \
  'module.configuration.global_policy_forbidden' \
  'module.configuration.declaration_collision' \
  'module.configuration.unavailable' \
  'RuntimeActive:' \
  'LifecycleDisabled' \
  'LifecycleDetached' \
  'LifecyclePurged' \
  'TestConfigurationBindingUsesValidatedDeclarationsRegistryAndScopePolicies' \
  'TestConfigurationBindingRetainsV1AndV2DeclarationsInValidatedSnapshot' \
  'TestConfigurationBindingRejectsOwnerClassMissingAndUndeclaredDefinitions' \
  'TestConfigurationBindingRejectsInvalidScopeContracts' \
  'TestConfigurationBindingRejectsCrossClassDeclarationCollision' \
  'TestConfigurationBindingLifecycleReadsAreNonDestructiveAndEnabledOnlyIsRuntimeActive' \
  'TestConfigurationBindingDoesNotRequireASecondRegistryForModulesWithoutDeclarations'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules kernel/internal/configuration/policy_validation.go; then
    echo "ERROR: P03.05 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

# P01.10/P02.09 and P03.01-P03.04 verifiers remain mandatory preceding workflow
# steps. This focused verifier proves only the new P03.05 integration boundary.
go vet ./kernel/internal/modules ./kernel/internal/configuration
go test ./kernel/internal/modules -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.05 G0 active-package governance boundary: PASS"
echo "P03.05 G1 manifest declaration snapshot binding without schema evolution: PASS"
echo "P03.05 G2 existing kernel.configuration registry reuse and collision/owner/class validation: PASS"
echo "P03.05 G3 explicit global/scoped registration plus existing P02.09 policy validation: PASS"
echo "P03.05 G4 lifecycle retained-read and enabled-only runtime-active semantics: PASS"
echo "P03.05 G5 tenant/org scope: trusted P02 ScopedService boundary only; no raw scope input: PASS"
echo "P03.05 G6 security: flags grant no permission/authorization authority and kill-switch ownership is not absorbed: PASS"
echo "P03.05 G7 persistence/migration: no new schema or data migration introduced: N/A"
echo "P03.05 G8 build/race/package: PASS"
echo "P03.05 G9 P03.06+ scope remains absent: PASS"

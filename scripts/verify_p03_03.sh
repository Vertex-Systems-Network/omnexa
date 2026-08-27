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

p03_03_files=(
  kernel/internal/modules/semver.go
  kernel/internal/modules/semver_test.go
  kernel/internal/modules/manifest_versioned.go
  kernel/internal/modules/manifest_versioned_test.go
  kernel/internal/modules/resolver.go
  kernel/internal/modules/resolver_test.go
  kernel/internal/modules/resolver_provenance_test.go
  kernel/internal/modules/registry.go
)

for file in "${p03_03_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.03 required source missing: ${file}" >&2
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
  echo "ERROR: P03.03 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

if grep -nE '"(os|os/exec|plugin|net/http|database/sql|path/filepath)"' "${p03_03_files[@]}"; then
  echo "ERROR: P03.03 resolver must remain pure in-memory validation without filesystem/network/database/package execution" >&2
  exit 1
fi
# Match actual qualified execution/discovery APIs, not English substrings such as
# "Global" in test names. Prohibited package imports are independently rejected
# above, so aliases cannot bypass this primitive guard.
if grep -nE '(os\.(ReadDir|ReadFile|Open)|filepath\.(Walk|WalkDir|Glob)|exec\.Command|plugin\.Open|http\.(Get|Post)|http\.Client|sql\.Open)[[:space:]]*[({]' "${p03_03_files[@]}"; then
  echo "ERROR: P03.03 source contains an unauthorized execution or discovery primitive" >&2
  exit 1
fi
if grep -nE 'type[[:space:]]+(LifecycleRuntime|ModuleStateStore|CapabilityRegistry|PermissionRegistry|PackageInstaller)([[:space:]]|$)' "${p03_03_files[@]}"; then
  echo "ERROR: P03.03 pulled lifecycle/authority/package-install scope forward" >&2
  exit 1
fi

for marker in \
  'SchemaVersionV2 = 2' \
  'MaxDependencyConstraintBytes = 256' \
  'MaxDependencyComparators = 16' \
  'type DependencyRequirement struct' \
  'func parseValidatedManifest(' \
  'decoder.DisallowUnknownFields()' \
  'manifest.dependency.self' \
  'manifest.dependency.constraint_invalid' \
  'type PlatformSnapshot struct' \
  'type DependencyObservation struct' \
  'type DependencyResolution struct' \
  'func ResolveDependencies(' \
  'resolver.dependency.version_contract_missing' \
  'resolver.dependency.required_missing' \
  'resolver.dependency.required_incompatible' \
  'resolver.graph.required_cycle' \
  'resolver.dependency.undeclared' \
  'resolver.dependency.private_forbidden' \
  'resolver.dependency.forbidden' \
  'resolver.registry.snapshot_missing' \
  'snapshots map[string]validatedManifestSnapshot' \
  'parseValidatedManifest(payload)' \
  'TestParseValidatedManifestDispatchesV1AndV2' \
  'TestDependencyConstraintGrammarIsBoundedAndExplicit' \
  'TestResolveDependenciesProducesStableRequiredOrder' \
  'TestResolveDependenciesRejectsMissingAndIncompatibleRequiredDependency' \
  'TestResolveDependenciesDegradesOptionalDependencyWithoutGlobalFailure' \
  'TestResolveDependenciesRejectsRequiredCycleButIgnoresOptionalCycleForGlobalOrder' \
  'TestResolveDependenciesAppliesSchemaV1MigrationRules' \
  'TestRegistryBoundSnapshotIsIndependentFromSourcePayloadMutation' \
  'TestRegistrySnapshotCloneCannotMutateBoundEvidence' \
  'TestResolveDependenciesRejectsMismatchedRegistrySnapshot' \
  'TestResolveDependenciesRejectsUndeclaredPrivateAndKernelToBusinessObservations'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: P03.03 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

# Historical P03.01 and P03.02 verification remain separate mandatory steps in
# the canonical governance workflow immediately before this dedicated verifier.
go vet ./kernel/internal/modules
go test ./kernel/internal/modules -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.03 G0 active-package / ADR-0012 governance boundary: PASS"
echo "P03.03 G1 strict v1/v2 parser dispatch, SemVer and bounded constraints: PASS"
echo "P03.03 G2 required/optional/platform/cycle/order/migration/provenance tests: PASS"
echo "P03.03 G3 deterministic registry-bound dependency-resolution contract: PASS"
echo "P03.03 G4 data/migration runtime: N/A — no persistence introduced"
echo "P03.03 G5 security: fail-closed required dependencies, forbidden/private hooks, safe diagnostics, no authority grants: PASS"
echo "P03.03 G6 lifecycle runtime: N/A — resolver output is eligibility/degradation metadata only"
echo "P03.03 G7 build/package: PASS"
echo "P03.03 G8 supply chain: stdlib-only implementation; module/workspace metadata unchanged: PASS"

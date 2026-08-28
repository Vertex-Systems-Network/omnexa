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

p03_08_files=(
  kernel/internal/modules/manifest_versioned.go
  kernel/internal/modules/ui_contribution_registry.go
  kernel/internal/modules/ui_contribution_registry_test.go
  docs/roadmap/work-packages/P03.08.md
)
for file in "${p03_08_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.08 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/modules/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P03.08 Go files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P03.08 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

if grep -nE '(TenantID|OrganizationID)[[:space:]]+string|Handler[[:space:]]+(func|any|interface)|Component[[:space:]]+(func|any|interface)|Payload[[:space:]]+(\[\]byte|any|interface)|HTML[[:space:]]+string|func[[:space:]].*(Render|Invoke|Execute)|database/sql|type[[:space:]]+(MigrationRegistry|ModuleHealthRegistry|PackageTrustRegistry)([[:space:]]|$)' kernel/internal/modules/ui_contribution_registry.go; then
  echo "ERROR: P03.08 introduced executable/raw-scope/private-runtime/future-registry authority" >&2
  exit 1
fi

for marker in \
  'UISlots                 []string' \
  'type UIContributionKind string' \
  'type UIFallbackBehavior string' \
  'type UIContributionRegistration struct' \
  'type UIContributionRecord struct' \
  'type UIContributionRegistry struct' \
  'func BindUIContributionRegistry(' \
  'func (r *UIContributionRegistry) List(' \
  'func (r *UIContributionRegistry) Lookup(' \
  'Available:          state == LifecycleEnabled' \
  'module.ui.registration_duplicate' \
  'module.ui.slot_undeclared' \
  'module.ui.permission_undeclared' \
  'module.ui.feature_flag_undeclared' \
  'module.ui.optional_dependency_undeclared' \
  'version_contract_missing' \
  'dependency_missing' \
  'dependency_incompatible' \
  'dependency_unavailable' \
  'TestUIContributionRegistryRetainsV1AndV2ValidatedUISlots' \
  'TestUIContributionRegistryDeterministicMetadataLookup' \
  'TestUIContributionRegistryRejectsInvalidRegistrationReferences' \
  'TestUIContributionRegistryOwnerLifecycleAvailabilityFailsClosed' \
  'TestUIContributionRegistryOptionalDependencyDegradesOnlyAffectedContribution' \
  'TestUIContributionRegistryOptionalDependencyCompatibilityAndLifecycle' \
  'TestUIContributionRegistrySchemaV1OptionalDependencyFailsSafeAsUnresolved' \
  'TestUIContributionRegistryMetadataSurfaceHasNoExecutionOrAuthorityFields'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: P03.08 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

# Retained P03.01-P03.07 verifiers execute as preceding canonical-governance
# steps. This focused verifier proves only the new P03.08 metadata boundary.
go vet ./kernel/internal/modules
go test ./kernel/internal/modules -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.08 G0 active-package governance boundary: PASS"
echo "P03.08 G1 v1/v2 validated UI-slot declaration retention: PASS"
echo "P03.08 G2 deterministic module/contribution/slot identity and collision validation: PASS"
echo "P03.08 G3 lifecycle availability: enabled-only contribution availability: PASS"
echo "P03.08 G4 optional dependency selective degradation and schema-v1 unresolved fallback: PASS"
echo "P03.08 G5 permission/feature-flag conditions remain descriptive and non-authorizing: PASS"
echo "P03.08 G6 security: no renderer/handler/raw tenant scope/secret/private database authority: PASS"
echo "P03.08 G7 persistence/migration/dependency change: N/A / unchanged: PASS"
echo "P03.08 G8 build/race/package: PASS"
echo "P03.08 G9 P03.09+ and business-feature runtime remain absent: PASS"

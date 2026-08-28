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

p03_06_files=(
  kernel/internal/modules/manifest_versioned.go
  kernel/internal/modules/capability_registry.go
  kernel/internal/modules/capability_registry_test.go
  docs/roadmap/work-packages/P03.06.md
)
for file in "${p03_06_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.06 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/modules/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P03.06 Go files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P03.06 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

if grep -nE '(TenantID|OrganizationID)[[:space:]]+string|Handler[[:space:]]+(func|any|interface)|func[[:space:]].*(Invoke|ExecuteCapability)|database/sql|type[[:space:]]+(PermissionRegistry|UIContributionRegistry|MigrationRegistry|ModuleHealthRegistry|PackageTrustRegistry)([[:space:]]|$)' kernel/internal/modules/capability_registry.go; then
  echo "ERROR: P03.06 introduced forbidden invocation/raw-scope/private-runtime/future-registry authority" >&2
  exit 1
fi

for marker in \
  'CapabilitiesProvided    []string' \
  'CapabilitiesConsumed    []string' \
  'type CapabilityRegistration struct' \
  'type CapabilityQuery struct' \
  'type CapabilityRecord struct' \
  'type CapabilityConsumer struct' \
  'type CapabilityRegistry struct' \
  'func BindCapabilityRegistry(' \
  'func (r *CapabilityRegistry) List(' \
  'func (r *CapabilityRegistry) Consumers()' \
  'func (r *CapabilityRegistry) Lookup(' \
  'func (r *CapabilityRegistry) ResolveConsumer(' \
  'Available:        state == LifecycleEnabled' \
  'module.capability.registration_undeclared' \
  'module.capability.registration_missing' \
  'module.capability.owner_mismatch' \
  'module.capability.declaration_collision' \
  'module.capability.owner_conflict' \
  'module.capability.version_incompatible' \
  'module.capability.provider_unavailable' \
  'module.capability.lifecycle_read_failed' \
  'TestCapabilityRegistryRetainsV1AndV2ValidatedDeclarations' \
  'TestCapabilityRegistryDeterministicProviderConsumerLookup' \
  'TestCapabilityRegistryRejectsUndeclaredOwnerDuplicateAndCollision' \
  'TestCapabilityRegistryRejectsInvalidCapabilityContractAndMetadataRefs' \
  'TestCapabilityRegistryMajorCompatibilityFailsClosed' \
  'TestCapabilityRegistryLifecycleAvailabilityIsEnabledOnlyAndNonDestructive' \
  'TestCapabilityRegistryMetadataSurfaceHasNoInvocationOrRawScopeAuthority'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: P03.06 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

# Retained P03.01-P03.05 verifiers execute as preceding canonical-governance
# steps. This focused verifier proves only the new P03.06 capability boundary.
go vet ./kernel/internal/modules
go test ./kernel/internal/modules -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.06 G0 active-package governance boundary: PASS"
echo "P03.06 G1 v1/v2 validated snapshot capability retention without manifest schema evolution: PASS"
echo "P03.06 G2 stable capability ID/major/provider-owner registration and deterministic lookup: PASS"
echo "P03.06 G3 duplicate/conflicting owner/version and undeclared registration fail closed: PASS"
echo "P03.06 G4 consumer major-version compatibility and missing-provider behavior fail closed: PASS"
echo "P03.06 G5 lifecycle availability: enabled-only active; unavailable historical identity retained: PASS"
echo "P03.06 G6 security: authorization/scope/contract references are descriptive metadata only; no invocation/raw scope/private handler authority: PASS"
echo "P03.06 G7 persistence/migration/dependency change: N/A / unchanged: PASS"
echo "P03.06 G8 build/race/package: PASS"
echo "P03.06 G9 P03.07+ and business-feature runtime remain absent: PASS"

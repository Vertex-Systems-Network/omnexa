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

p03_07_files=(
  kernel/internal/modules/manifest_versioned.go
  kernel/internal/modules/permission_registry.go
  kernel/internal/modules/permission_registry_test.go
  kernel/internal/authorization/permission_catalog.go
  kernel/internal/authorization/permission_catalog_test.go
  kernel/internal/authorization/service.go
  kernel/internal/authorization/service_account.go
  kernel/migrations/kernel.authorization/4_add_module_permission_registration.sql
  docs/roadmap/work-packages/P03.07.md
)
for file in "${p03_07_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.07 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/modules/*.go kernel/internal/authorization/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P03.07 Go files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P03.07 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

if grep -nE 'TenantID[[:space:]]+string|OrganizationID[[:space:]]+string|SuperAdmin|Bypass|type[[:space:]]+(UIContributionRegistry|MigrationRegistry|ModuleHealthRegistry|PackageTrustRegistry)([[:space:]]|$)' kernel/internal/modules/permission_registry.go; then
  echo "ERROR: P03.07 introduced raw-scope/bypass/future-registry authority" >&2
  exit 1
fi

for marker in \
  'Permissions             []string' \
  'type PermissionRegistration struct' \
  'type PermissionRecord struct' \
  'type PermissionRegistry struct' \
  'func BindPermissionRegistry(' \
  'func (r *PermissionRegistry) PermissionAvailable(' \
  'func (r *PermissionRegistry) SynchronizeCatalog(' \
  'Available:      state == LifecycleEnabled' \
  'module.permission.namespace_invalid' \
  'module.permission.declaration_collision' \
  'module.permission.registration_missing' \
  'module.permission.owner_mismatch' \
  'module.permission.capability_undeclared' \
  'type ModulePermissionAvailability interface' \
  'type ModulePermissionDefinition struct' \
  'func NewServiceWithModulePermissionAvailability(' \
  'available, availabilityErr := service.permissionAvailable(ctx, permission)' \
  'func (repository *PostgresRepository) ReconcileModulePermissions(' \
  "source_kind IN ('kernel', 'module')" \
  "source_kind = 'module'" \
  'AND module_id ~' \
  'TestPermissionRegistryRetainsV1AndV2ValidatedDeclarations' \
  'TestPermissionRegistryDeterministicRegistrationAndCapabilityAssociation' \
  'TestPermissionRegistryRejectsReservedNamespaceCollisionAndOwnershipMismatch' \
  'TestPermissionRegistryLifecycleAvailabilityIsEnabledOnlyAndHistorySafe' \
  'TestModulePermissionRequiresLiveAvailabilityAndExistingGrant'; do
  if ! grep -R -Fq "$marker" kernel/internal kernel/migrations; then
    echo "ERROR: P03.07 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

if grep -nE 'DELETE[[:space:]]+FROM[[:space:]]+omnexa_authorization\.permissions|ON[[:space:]]+DELETE[[:space:]]+CASCADE' kernel/migrations/kernel.authorization/4_add_module_permission_registration.sql; then
  echo "ERROR: P03.07 migration may not delete permission identity/history" >&2
  exit 1
fi

# Retained P03.01-P03.06 verifiers execute as preceding canonical-governance
# steps. This focused verifier proves only the new P03.07 permission boundary.
go vet ./kernel/internal/modules ./kernel/internal/authorization
go test ./kernel/internal/modules -count=1
go test ./kernel/internal/authorization -count=1
go test -race ./kernel/internal/modules -count=1
go test -race ./kernel/internal/authorization -count=1
go build ./kernel/...

echo "P03.07 G0 active-package governance boundary: PASS"
echo "P03.07 G1 v1/v2 validated permission declaration retention: PASS"
echo "P03.07 G2 deterministic owner/module registration and namespace/collision validation: PASS"
echo "P03.07 G3 lifecycle availability: enabled-only; identity/history retained: PASS"
echo "P03.07 G4 authorization: availability never grants; existing exact-scope role grant remains required: PASS"
echo "P03.07 G5 security: unknown/unavailable deny; no role-name bypass/raw tenant scope/new auth engine: PASS"
echo "P03.07 G6 capability association remains descriptive/non-invoking: PASS"
echo "P03.07 G7 additive authorization-owned catalog migration preserves permission/role history: PASS"
echo "P03.07 G8 build/race/package: PASS"
echo "P03.07 G9 P03.08+ and business-feature runtime remain absent: PASS"

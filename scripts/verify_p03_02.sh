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

for file in \
  kernel/internal/modules/registry.go \
  kernel/internal/modules/registry_test.go; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.02 required source missing: ${file}" >&2
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
  echo "ERROR: P03.02 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

if grep -nE '"(os|os/exec|plugin|net/http|database/sql|path/filepath)"' kernel/internal/modules/registry.go kernel/internal/modules/registry_test.go; then
  echo "ERROR: P03.02 discovery must not scan filesystem/network resources or execute package code" >&2
  exit 1
fi
if grep -nE '(ReadDir|WalkDir|Glob|exec\.Command|plugin\.Open|http\.(Get|Post|Client)|sql\.Open)' kernel/internal/modules/registry.go kernel/internal/modules/registry_test.go; then
  echo "ERROR: P03.02 source contains an unauthorized discovery/execution primitive" >&2
  exit 1
fi
if grep -nE 'type[[:space:]]+(DependencyGraph|LifecycleRuntime|ModuleStateStore|CapabilityRegistry|PermissionRegistry)([[:space:]]|$)' kernel/internal/modules/registry.go kernel/internal/modules/registry_test.go; then
  echo "ERROR: P03.02 pulled P03.03+ dependency/lifecycle/authority runtime forward" >&2
  exit 1
fi

for marker in \
  'type DiscoverySource struct' \
  'type RegistryRecord struct' \
  'type Registry struct' \
  'func Discover(' \
  'discovery.source.invalid' \
  'discovery.manifest.invalid' \
  'discovery.module.duplicate' \
  'discovery.module.version_conflict' \
  'ParseManifest(payload)' \
  'TestDiscoverEmptySourcesReturnsEmptyRegistry' \
  'TestDiscoverIsDeterministicAcrossEnumerationOrder' \
  'TestDiscoverRejectsDuplicateModuleIdentity' \
  'TestDiscoverRejectsConflictingModuleVersions' \
  'TestDiscoverRejectsMalformedOrUnvalidatedManifest' \
  'TestDiscoverRejectsInvalidExplicitSourceWithoutLeakingRawIdentity' \
  'TestRegistryContainsAvailabilityMetadataOnly'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules/registry.go kernel/internal/modules/registry_test.go; then
    echo "ERROR: P03.02 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

# Retain the completed P03.01 parser/validation contract before exercising the
# new registry/discovery layer.
bash scripts/verify_p03_01.sh

go vet ./kernel/internal/modules
go test ./kernel/internal/modules -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.02 G0 active-package / governance boundary: PASS"
echo "P03.02 G1 deterministic explicit-source registry/discovery: PASS"
echo "P03.02 G2 empty/order/duplicate/conflict/malformed/source-boundary coverage: PASS"
echo "P03.02 G3 validated-manifest registry contract and stable lookup/list identity: PASS"
echo "P03.02 G4 data/migration runtime: N/A — no persistence introduced"
echo "P03.02 G5 security: explicit sources, safe diagnostics, no authority grants, no code/network/filesystem execution: PASS"
echo "P03.02 G6 lifecycle runtime: N/A — discovered metadata only; installed/enabled state remains out of scope"
echo "P03.02 G7 build/package: PASS"
echo "P03.02 G8 supply chain: stdlib-only implementation; module metadata unchanged: PASS"

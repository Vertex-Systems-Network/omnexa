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

p03_11_files=(
  kernel/internal/modules/manifest_versioned.go
  kernel/internal/modules/package_trust_profile.go
  kernel/internal/modules/package_trust_profile_test.go
  kernel/internal/modules/p03_exit_test.go
  docs/roadmap/work-packages/P03.11.md
  docs/governance/P03_EXIT_GATE.md
)
for file in "${p03_11_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.11 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/modules/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P03.11 Go files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P03.11 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

# P03.11 is metadata-only. Trust-root verification, filesystem/network discovery,
# package execution, persistence and authorization decisions remain out of scope.
if grep -nE '"(os|os/exec|plugin|net/http|database/sql|path/filepath|crypto/x509)"' \
  kernel/internal/modules/package_trust_profile.go \
  kernel/internal/modules/package_trust_profile_test.go \
  kernel/internal/modules/p03_exit_test.go; then
  echo "ERROR: P03.11 introduced package execution, trust-root, filesystem, network or database authority" >&2
  exit 1
fi
if grep -nE 'Trusted[[:space:]]+bool|Certified[[:space:]]+bool|SignatureValid[[:space:]]+bool|func[[:space:]].*(VerifySignature|VerifyPublisher|Certify|TrustPackage)' \
  kernel/internal/modules/package_trust_profile.go; then
  echo "ERROR: P03.11 introduced an unauthorized trust/certification decision surface" >&2
  exit 1
fi

for marker in \
  'DataClassification      []string' \
  'Security                SecurityDeclaration' \
  'Publisher               string' \
  'ProvenanceRef           string' \
  'SBOMRef                 string' \
  'cloneSecurityDeclaration' \
  'type PackageTrustAuthority string' \
  'PackageTrustMetadataOnly' \
  'type PublisherIdentityHook struct' \
  'type PackageProvenanceHook struct' \
  'type SBOMIdentityHook struct' \
  'type PackageDeclaredScopeProfile struct' \
  'type PackageTrustProfile struct' \
  'func BuildPackageTrustProfiles(' \
  'func PackageTrustProfileFor(' \
  'sortedSecretReferenceNames' \
  'TestPackageTrustProfileRetainsValidatedV1AndV2Metadata' \
  'TestPackageTrustProfileIsDeterministicAndIndependentFromCallerMutation' \
  'TestPackageTrustProfileOptionalHooksRemainAbsentWithoutClaimingTrust' \
  'TestPackageTrustProfileFailsClosedOnRegistrySnapshotMismatch' \
  'TestPackageTrustProfilePublicSurfaceHasNoTrustDecisionOrSecretValueFields' \
  'TestP03ExitReferenceModuleAggregate' \
  'TestP03ExitUnrelatedModuleIsolationAcrossLifecycleOperations'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: P03.11 required contract marker missing: ${marker}" >&2
    exit 1
  fi
done

# Aggregate EX-01..EX-07 proof must remain visibly mapped to the canonical gate.
for marker in \
  'EX-01-required-dependency-enforcement' \
  'EX-02-optional-dependency-degradation' \
  'EX-03-safe-disable-reenable' \
  'EX-04-upgrade-migration-path' \
  'EX-05-forbidden-dependency-detection' \
  'EX-06-health-state-accuracy' \
  'EX-07-unrelated-module-isolation'; do
  if ! grep -Fq "$marker" kernel/internal/modules/p03_exit_test.go; then
    echo "ERROR: P03 exit aggregate mapping missing: ${marker}" >&2
    exit 1
  fi
done

# Retained runtime boundaries remain authoritative; P03.11 only composes them.
for marker in \
  'func ResolveDependencies(' \
  'func (m LifecycleManager) Apply(' \
  'func (r *MigrationOwnershipRegistry) UpgradePlan(' \
  'func (r *ModuleHealthReporter) Report('; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: retained P03 runtime invariant missing: ${marker}" >&2
    exit 1
  fi
done

# Dedicated tests run explicitly in addition to the full retained package suite.
go test ./kernel/internal/modules -run '^(TestPackageTrustProfile|TestP03Exit)' -count=1
go vet ./kernel/internal/modules ./kernel/internal/operations
go test ./kernel/internal/modules ./kernel/internal/operations -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.11 G0 active-package governance boundary + P03 exit gate mapping: PASS"
echo "P03.11 G1 typed/versioned publisher/provenance/SBOM/scope metadata hooks: PASS"
echo "P03.11 G2 metadata presence remains explicitly non-authoritative; XTRUST-100 not pulled forward: PASS"
echo "P03.11 G3 validated v1/v2 registry snapshots retain immutable trust/profile declarations: PASS"
echo "P03.11 G4 classification-safe secret-name/scope projection with no secret values or package execution: PASS"
echo "P03.11 G5 EX-01..EX-05 dependency/lifecycle/upgrade/migration/forbidden-coupling aggregate proof: PASS"
echo "P03.11 G6 EX-06..EX-07 health/state accuracy and unrelated-module isolation aggregate proof: PASS"
echo "P03.11 G7 retained P03.01-P03.10 plus repository P01/P02 workflow regression chain remains mandatory: PASS"
echo "P03.11 G8 build/race/package + P04/business/AI runtime remain locked: PASS"

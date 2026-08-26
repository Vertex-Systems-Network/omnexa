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

for file in kernel/internal/modules/manifest.go kernel/internal/modules/validate.go kernel/internal/modules/manifest_test.go; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P03.01 required source missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/modules/*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P03.01 files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P03.01 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

if grep -R -nE '"(os|os/exec|plugin|net/http|database/sql)"' kernel/internal/modules; then
  echo "ERROR: P03.01 manifest parsing must remain declarative and side-effect free" >&2
  exit 1
fi
if grep -R -nE 'type[[:space:]]+(Registry|Discovery|DependencyGraph|LifecycleRuntime|ModuleStateStore)([[:space:]]|$)' kernel/internal/modules; then
  echo "ERROR: P03.01 pulled P03.02+ registry/dependency/lifecycle runtime scope forward" >&2
  exit 1
fi
if grep -R -nE '(password|private_key|raw_secret|access_token|client_secret)[[:space:]]*[:=]' kernel/internal/modules; then
  echo "ERROR: P03.01 source contains credential-like value material" >&2
  exit 1
fi

for marker in \
  'MaxManifestBytes' \
  'DisallowUnknownFields' \
  'manifest.parse.too_large' \
  'manifest.schema.unsupported' \
  'manifest.dependency.class_conflict' \
  'manifest.secret.reference_invalid'; do
  if ! grep -R -Fq "$marker" kernel/internal/modules; then
    echo "ERROR: P03.01 required validation marker missing: ${marker}" >&2
    exit 1
  fi
done

go vet ./kernel/internal/modules
go test ./kernel/internal/modules -count=1
go test -race ./kernel/internal/modules -count=1
go build ./kernel/...

echo "P03.01 G0 active-package / governance boundary: PASS"
echo "P03.01 G1 bounded strict parser / deterministic static validation: PASS"
echo "P03.01 G2 positive and negative manifest unit coverage: PASS"
echo "P03.01 G3 schema/dependency/contribution/security declaration contract: PASS"
echo "P03.01 G4 data/migration runtime: N/A — declarations only; no persistence introduced"
echo "P03.01 G5 security: untrusted metadata, no secret values, no authority grants, no code execution: PASS"
echo "P03.01 G6 lifecycle runtime: N/A — lifecycle hooks are declarations only"
echo "P03.01 G7 build/package: PASS"
echo "P03.01 G8 supply chain: stdlib-only implementation; module metadata unchanged: PASS"

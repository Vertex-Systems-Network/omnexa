#!/usr/bin/env bash
set -euo pipefail

GOLANGCI_LINT_VERSION="v2.12.2"
GOVULNCHECK_VERSION="v1.7.0"
EXPECTED_GO="go$(tr -d '[:space:]' < .go-version)"
TOOL_ROOT="${RUNNER_TEMP:-${TMPDIR:-/tmp}}/omnexa-go-quality"

if [[ "$(go env GOVERSION)" != "${EXPECTED_GO}" ]]; then
  echo "ERROR: Go toolchain mismatch: got $(go env GOVERSION), want ${EXPECTED_GO}" >&2
  exit 1
fi

mkdir -p "${TOOL_ROOT}"

GOBIN="${TOOL_ROOT}" GOTOOLCHAIN=local go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}"
GOBIN="${TOOL_ROOT}" GOTOOLCHAIN=local go install "golang.org/x/vuln/cmd/govulncheck@${GOVULNCHECK_VERSION}"

"${TOOL_ROOT}/golangci-lint" version
"${TOOL_ROOT}/govulncheck" -version

"${TOOL_ROOT}/golangci-lint" config verify
"${TOOL_ROOT}/golangci-lint" run ./kernel/...
"${TOOL_ROOT}/govulncheck" ./kernel/...

echo "Omnexa Go code-quality gate: PASS"
echo "golangci-lint: ${GOLANGCI_LINT_VERSION}"
echo "govulncheck: ${GOVULNCHECK_VERSION}"
echo "Go toolchain: ${EXPECTED_GO}"

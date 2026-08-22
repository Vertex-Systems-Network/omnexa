#!/usr/bin/env bash
set -euo pipefail

GOLANGCI_LINT_VERSION="v2.12.2"
GOVULNCHECK_VERSION="v1.7.0"
EXPECTED_GO="go$(tr -d '[:space:]' < .go-version)"
TOOL_ROOT="${RUNNER_TEMP:-${TMPDIR:-/tmp}}/omnexa-go-quality"
FORMAT_DIFF="${TOOL_ROOT}/gofmt.diff"

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

status=0

: > "${FORMAT_DIFF}"
gofmt -d kernel > "${FORMAT_DIFF}"
if [[ -s "${FORMAT_DIFF}" ]]; then
  echo "ERROR: gofmt changes are required:" >&2
  cat "${FORMAT_DIFF}" >&2
  status=1
else
  echo "gofmt: PASS"
fi

if ! "${TOOL_ROOT}/golangci-lint" run ./kernel/...; then
  echo "ERROR: golangci-lint reported findings" >&2
  status=1
else
  echo "golangci-lint: PASS"
fi

if ! "${TOOL_ROOT}/govulncheck" ./kernel/...; then
  echo "ERROR: govulncheck reported reachable vulnerabilities or could not complete" >&2
  status=1
else
  echo "govulncheck: PASS"
fi

if [[ "${status}" -ne 0 ]]; then
  echo "Omnexa Go code-quality gate: FAIL" >&2
  exit "${status}"
fi

echo "Omnexa Go code-quality gate: PASS"
echo "golangci-lint: ${GOLANGCI_LINT_VERSION}"
echo "govulncheck: ${GOVULNCHECK_VERSION}"
echo "Go toolchain: ${EXPECTED_GO}"

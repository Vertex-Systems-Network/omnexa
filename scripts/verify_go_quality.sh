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

# XQ-100 M2 multi-agent enforcement is intentionally executed inside the
# already-required `governance` job. These checks are dependency-free and
# fail before expensive Go tooling when active task/lease/scope state is bad.
python -m py_compile \
  scripts/agent_orchestration_common.py \
  scripts/validate_agent_task.py \
  scripts/validate_agent_leases.py \
  scripts/detect_path_overlap.py \
  scripts/validate_agent_pr_scope.py \
  scripts/validate_agent_base_sha.py \
  scripts/validate_task_dependencies.py
python scripts/validate_agent_task.py
python scripts/validate_agent_leases.py
python scripts/detect_path_overlap.py
python scripts/validate_task_dependencies.py
python scripts/validate_agent_pr_scope.py
python scripts/validate_agent_base_sha.py

mkdir -p "${TOOL_ROOT}"

GOBIN="${TOOL_ROOT}" GOTOOLCHAIN=local go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}"
GOBIN="${TOOL_ROOT}" GOTOOLCHAIN=local go install "golang.org/x/vuln/cmd/govulncheck@${GOVULNCHECK_VERSION}"

"${TOOL_ROOT}/golangci-lint" version
"${TOOL_ROOT}/govulncheck" -version

status=0

if config_output=$("${TOOL_ROOT}/golangci-lint" config verify 2>&1); then
  echo "golangci-lint configuration: PASS"
else
  echo "ERROR: golangci-lint configuration validation failed" >&2
  printf '%s\n' "${config_output}" >&2
  status=1
fi

mapfile -t go_files < <(find kernel -type f -name '*.go' -print | sort)
if [[ ${#go_files[@]} -eq 0 ]]; then
  echo "ERROR: no Go source files found under kernel/" >&2
  exit 1
fi

echo "gofmt input files: ${#go_files[@]}"
: > "${FORMAT_DIFF}"
if gofmt -d "${go_files[@]}" > "${FORMAT_DIFF}"; then
  if [[ -s "${FORMAT_DIFF}" ]]; then
    echo "ERROR: gofmt changes are required:" >&2
    cat "${FORMAT_DIFF}" >&2
    status=1
  else
    echo "gofmt: PASS"
  fi
else
  echo "ERROR: gofmt could not analyze one or more Go files" >&2
  if [[ -s "${FORMAT_DIFF}" ]]; then
    cat "${FORMAT_DIFF}" >&2
  fi
  status=1
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

#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

expected_go="$(tr -d '[:space:]' < .go-version)"
actual_go="$(go env GOVERSION)"
if [[ "$actual_go" != "go${expected_go}" ]]; then
  echo "ERROR: Go toolchain mismatch: got ${actual_go}, want go${expected_go}" >&2
  exit 1
fi

echo "P01.03 Go toolchain: ${actual_go}"

mapfile -t go_files < <(find kernel -type f -name '*.go' -print | sort)
unformatted="$(gofmt -l "${go_files[@]}")"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go vet ./kernel/...
go test ./kernel/...
go test -race ./kernel/internal/failure

go test ./kernel/internal/failure -run 'TestWrapPreservesCauseAndNeverPublishesPrivateCause|TestCodeLookupSurvivesOuterWrapping|TestViolationsAreDeterministicDeduplicatedAndBounded|TestConstructorsRejectUnsafeOrInvalidContractData' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P01.03 failure package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
    net/http|database/sql|log/slog)
      echo "ERROR: P01.03 failure package imports out-of-scope transport/data/logging package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/failure)

if grep -R -nE 'go\.opentelemetry|otel|database/sql|net/http' kernel/internal/failure --include='*.go'; then
  echo "ERROR: P01.03 source contains transport/database/telemetry coupling" >&2
  exit 1
fi

go build ./kernel/...

echo "P01.03 format/static: PASS"
echo "P01.03 unit/race tests: PASS"
echo "P01.03 wrapping/errors.Is-errors.As: PASS"
echo "P01.03 public/private redaction: PASS"
echo "P01.03 validation bounds/order: PASS"
echo "P01.03 dependency boundary: PASS"
echo "P01.03 build/package: PASS"

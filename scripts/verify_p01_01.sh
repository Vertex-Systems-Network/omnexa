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

echo "P01.01 Go toolchain: ${actual_go}"

mapfile -t go_files < <(find kernel -type f -name '*.go' -print | sort)
if [[ ${#go_files[@]} -eq 0 ]]; then
  echo "ERROR: no Go source files found under kernel/" >&2
  exit 1
fi

unformatted="$(gofmt -l "${go_files[@]}")"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: Go module/workspace metadata is not canonical; run go work sync && go mod tidy" >&2
  git status --short -- go.mod go.work go.sum >&2
  exit 1
fi

go vet ./kernel/...
go test ./kernel/...

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel"|"$module_prefix/kernel/"*) ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P01.01 kernel imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/...)

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

commit="${GITHUB_SHA:-$(git rev-parse HEAD)}"
version="0.1.0-dev"
ldflags="-X ${module_prefix}/kernel/internal/buildinfo.Version=${version} -X ${module_prefix}/kernel/internal/buildinfo.Commit=${commit}"

go build -trimpath -buildvcs=false -ldflags "$ldflags" -o "$tmp/omnexa" ./kernel/cmd/omnexa
output="$($tmp/omnexa)"
expected="omnexa-kernel version=${version} commit=${commit}"
if [[ "$output" != "$expected" ]]; then
  echo "ERROR: kernel smoke output mismatch" >&2
  echo "got:  $output" >&2
  echo "want: $expected" >&2
  exit 1
fi

if printf '%s' "$output" | grep -Eq '(/home/|/Users/|[A-Za-z]:\\\\|password|token=|secret=)'; then
  echo "ERROR: build metadata output contains local-path or secret-like material" >&2
  exit 1
fi

echo "P01.01 format/static: PASS"
echo "P01.01 unit tests: PASS"
echo "P01.01 dependency boundary: PASS"
echo "P01.01 build/smoke: PASS"

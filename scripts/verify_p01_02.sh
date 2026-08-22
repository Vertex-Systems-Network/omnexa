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

echo "P01.02 Go toolchain: ${actual_go}"

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
  echo "ERROR: Go module/workspace metadata is not canonical" >&2
  git status --short -- go.mod go.work go.sum >&2
  exit 1
fi

go vet ./kernel/...
go test ./kernel/...
go test -race ./kernel/internal/config

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel"|"$module_prefix/kernel/"*) ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: active P01.02 kernel imports out-of-scope Omnexa package: ${package}" >&2
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

expected="omnexa-kernel version=${version} commit=${commit}"
default_output="$($tmp/omnexa)"
if [[ "$default_output" != "$expected" ]]; then
  echo "ERROR: default configuration smoke output mismatch" >&2
  exit 1
fi

cat > "$tmp/config.json" <<'JSON'
{"environment":"ci"}
JSON
file_output="$(OMNEXA_CONFIG_FILE="$tmp/config.json" "$tmp/omnexa")"
if [[ "$file_output" != "$expected" ]]; then
  echo "ERROR: explicit configuration-file smoke output mismatch" >&2
  exit 1
fi

env_output="$(OMNEXA_CONFIG_FILE="$tmp/config.json" OMNEXA_ENVIRONMENT=test "$tmp/omnexa")"
if [[ "$env_output" != "$expected" ]]; then
  echo "ERROR: environment-over-file precedence smoke output mismatch" >&2
  exit 1
fi

invalid_value="production-secret-like-runtime-value"
if OMNEXA_ENVIRONMENT="$invalid_value" "$tmp/omnexa" >"$tmp/invalid.out" 2>"$tmp/invalid.err"; then
  echo "ERROR: invalid environment configuration unexpectedly succeeded" >&2
  exit 1
fi
if [[ -s "$tmp/invalid.out" ]]; then
  echo "ERROR: invalid configuration produced stdout" >&2
  exit 1
fi
if ! grep -q 'omnexa configuration error' "$tmp/invalid.err"; then
  echo "ERROR: invalid configuration did not produce the safe configuration error boundary" >&2
  exit 1
fi
if grep -Fq "$invalid_value" "$tmp/invalid.err"; then
  echo "ERROR: invalid configuration leaked the raw value" >&2
  exit 1
fi

unknown_value="should-not-be-accepted"
if OMNEXA_ENVIRONMNET="$unknown_value" "$tmp/omnexa" >"$tmp/unknown.out" 2>"$tmp/unknown.err"; then
  echo "ERROR: unknown OMNEXA_ configuration key unexpectedly succeeded" >&2
  exit 1
fi
if grep -Fq "$unknown_value" "$tmp/unknown.err"; then
  echo "ERROR: unknown-key failure leaked the raw value" >&2
  exit 1
fi

echo "P01.02 format/static: PASS"
echo "P01.02 unit/race tests: PASS"
echo "P01.02 precedence/startup smoke: PASS"
echo "P01.02 secret-safe negative tests: PASS"
echo "P01.02 dependency boundary: PASS"
echo "P01.02 build/package: PASS"

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

echo "P01.05 Go toolchain: ${actual_go}"

if [[ -z "${P01_05_TEST_CACHE_ADDRESS:-}" ]]; then
  echo "ERROR: P01_05_TEST_CACHE_ADDRESS is required for P01.05 integration evidence" >&2
  exit 1
fi
if [[ -z "${P01_05_TEST_VALKEY_IMAGE:-}" ]]; then
  echo "ERROR: P01_05_TEST_VALKEY_IMAGE is required for P01.05 lifecycle evidence" >&2
  exit 1
fi

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

valkey_client_version="$(go list -m -f '{{.Version}}' github.com/valkey-io/valkey-go)"
if [[ "$valkey_client_version" != "v1.0.75" ]]; then
  echo "ERROR: valkey-go version = ${valkey_client_version}, want v1.0.75" >&2
  exit 1
fi

echo "P01.05 valkey-go version: ${valkey_client_version}"

go vet ./kernel/...
go test ./kernel/...
go test -race ./kernel/internal/cache -run 'Test(LoadConfiguration|SettingsFromConfig|RenderKey|JSONCodec|UnavailableConnection)' -count=1
go test -v ./kernel/internal/cache -run 'TestUnavailableConnectionIsBoundedAndSafe|TestValkeyFoundationIntegration' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/cache"|"$module_prefix/kernel/internal/config"|"$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P01.05 cache package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/cache)

if grep -R -nE 'database/sql|jackc/pgx|go\.opentelemetry|/modules/|kernel/internal/database|kernel/internal/storage' kernel/internal/cache --include='*.go'; then
  echo "ERROR: P01.05 source contains out-of-scope data/storage/telemetry/module coupling" >&2
  exit 1
fi

mapfile -t valkey_containers < <(docker ps --filter "ancestor=${P01_05_TEST_VALKEY_IMAGE}" --format '{{.ID}}')
if [[ ${#valkey_containers[@]} -ne 1 ]]; then
  echo "ERROR: expected exactly one governed Valkey test container, found ${#valkey_containers[@]}" >&2
  exit 1
fi

valkey_container="${valkey_containers[0]}"
docker restart "$valkey_container" >/dev/null
healthy=false
for _ in $(seq 1 30); do
  status="$(docker inspect --format '{{if .Config.Healthcheck}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$valkey_container")"
  if [[ "$status" == "healthy" || "$status" == "running" ]]; then
    healthy=true
    break
  fi
  sleep 1
done
if [[ "$healthy" != "true" ]]; then
  echo "ERROR: Valkey provider did not recover after governed restart" >&2
  exit 1
fi

P01_05_AFTER_RESTART=1 go test -v ./kernel/internal/cache -run 'TestValkeyReconnectAfterProviderRestart' -count=1

go build ./kernel/...

echo "P01.05 G1 format/static/dependency boundary: PASS"
echo "P01.05 G2 unit/race key/TTL/serialization: PASS"
echo "P01.05 G3 Redis-compatible provider integration: PASS"
echo "P01.05 G5 provider secret redaction and scope guard: PASS"
echo "P01.05 G6 flush/restart/non-authority lifecycle: PASS"
echo "P01.05 G7 build/package: PASS"

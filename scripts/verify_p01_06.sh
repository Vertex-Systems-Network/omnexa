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

echo "P01.06 Go toolchain: ${actual_go}"

if [[ -z "${P01_06_TEST_S3_ENDPOINT:-}" ]]; then
  echo "ERROR: P01_06_TEST_S3_ENDPOINT is required for P01.06 integration evidence" >&2
  exit 1
fi
if [[ -z "${P01_06_TEST_S3_BUCKET:-}" ]]; then
  echo "ERROR: P01_06_TEST_S3_BUCKET is required for P01.06 integration evidence" >&2
  exit 1
fi
if [[ -z "${P01_06_TEST_S3_IMAGE:-}" ]]; then
  echo "ERROR: P01_06_TEST_S3_IMAGE is required for P01.06 lifecycle evidence" >&2
  exit 1
fi
if [[ "$P01_06_TEST_S3_IMAGE" != "adobe/s3mock:5.1.0" ]]; then
  echo "ERROR: P01.06 S3Mock image must be pinned to adobe/s3mock:5.1.0" >&2
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

aws_root_version="$(go list -m -f '{{.Version}}' github.com/aws/aws-sdk-go-v2)"
aws_credentials_version="$(go list -m -f '{{.Version}}' github.com/aws/aws-sdk-go-v2/credentials)"
aws_s3_version="$(go list -m -f '{{.Version}}' github.com/aws/aws-sdk-go-v2/service/s3)"
smithy_version="$(go list -m -f '{{.Version}}' github.com/aws/smithy-go)"

if [[ "$aws_root_version" != "v1.43.7" ]]; then
  echo "ERROR: aws-sdk-go-v2 version = ${aws_root_version}, want v1.43.7" >&2
  exit 1
fi
if [[ "$aws_credentials_version" != "v1.19.37" ]]; then
  echo "ERROR: aws-sdk-go-v2/credentials version = ${aws_credentials_version}, want v1.19.37" >&2
  exit 1
fi
if [[ "$aws_s3_version" != "v1.107.3" ]]; then
  echo "ERROR: aws-sdk-go-v2/service/s3 version = ${aws_s3_version}, want v1.107.3" >&2
  exit 1
fi
if [[ "$smithy_version" != "v1.27.8" ]]; then
  echo "ERROR: smithy-go version = ${smithy_version}, want v1.27.8" >&2
  exit 1
fi

echo "P01.06 aws-sdk-go-v2: ${aws_root_version}"
echo "P01.06 aws-sdk-go-v2/credentials: ${aws_credentials_version}"
echo "P01.06 aws-sdk-go-v2/service/s3: ${aws_s3_version}"
echo "P01.06 smithy-go: ${smithy_version}"
echo "P01.06 S3-compatible provider image: ${P01_06_TEST_S3_IMAGE}"

go vet ./kernel/...
go test ./kernel/...
go test -race ./kernel/internal/storage -run 'Test(LoadConfiguration|SettingsFromConfig|RenderKey|ValidateUpload|IntegrityReader|VerifiedReadCloser|ProviderFailure)' -count=1
go test -v ./kernel/internal/storage -run 'TestUnavailableStorageConnectionIsBoundedAndSafe|TestS3CompatibleStorageFoundationIntegration' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/storage"|"$module_prefix/kernel/internal/config"|"$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P01.06 storage package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/storage)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(cache|database)|go\.opentelemetry' kernel/internal/storage --include='*.go'; then
  echo "ERROR: P01.06 storage source contains out-of-scope module/cache/database/telemetry coupling" >&2
  exit 1
fi

mapfile -t s3_containers < <(docker ps --filter "ancestor=${P01_06_TEST_S3_IMAGE}" --format '{{.ID}}')
if [[ ${#s3_containers[@]} -ne 1 ]]; then
  echo "ERROR: expected exactly one governed S3Mock test container, found ${#s3_containers[@]}" >&2
  exit 1
fi

s3_container="${s3_containers[0]}"
docker restart "$s3_container" >/dev/null
ready=false
for _ in $(seq 1 45); do
  if curl --silent --output /dev/null --connect-timeout 1 --max-time 2 "${P01_06_TEST_S3_ENDPOINT}/" 2>/dev/null; then
    ready=true
    break
  fi
  sleep 1
done
if [[ "$ready" != "true" ]]; then
  echo "ERROR: S3-compatible provider did not recover after governed restart" >&2
  exit 1
fi

P01_06_AFTER_RESTART=1 go test -v ./kernel/internal/storage -run 'TestS3CompatibleReconnectAfterProviderRestart' -count=1

go build ./kernel/...

echo "P01.06 G1 format/static/dependency boundary: PASS"
echo "P01.06 G2 unit/race key/metadata/stream/integrity: PASS"
echo "P01.06 G3 S3-compatible provider contract/integration: PASS"
echo "P01.06 G5 provider secret redaction/path/metadata scope guard: PASS"
echo "P01.06 G6 timeout/cancellation/restart/integrity resilience: PASS"
echo "P01.06 G7 build/package: PASS"

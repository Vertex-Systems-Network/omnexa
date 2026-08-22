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

echo "P01.07 Go toolchain: ${actual_go}"

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
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

otel_version="$(go list -m -f '{{.Version}}' go.opentelemetry.io/otel)"
otel_metric_version="$(go list -m -f '{{.Version}}' go.opentelemetry.io/otel/metric)"
otel_trace_version="$(go list -m -f '{{.Version}}' go.opentelemetry.io/otel/trace)"
otel_sdk_version="$(go list -m -f '{{.Version}}' go.opentelemetry.io/otel/sdk)"
otel_sdk_metric_version="$(go list -m -f '{{.Version}}' go.opentelemetry.io/otel/sdk/metric)"
for item in \
  "otel:${otel_version}" \
  "otel/metric:${otel_metric_version}" \
  "otel/trace:${otel_trace_version}" \
  "otel/sdk:${otel_sdk_version}" \
  "otel/sdk/metric:${otel_sdk_metric_version}"; do
  name="${item%%:*}"
  version="${item#*:}"
  if [[ "$version" != "v1.45.0" ]]; then
    echo "ERROR: ${name} version = ${version}, want v1.45.0" >&2
    exit 1
  fi
done

echo "P01.07 OpenTelemetry Go API: ${otel_version}"
echo "P01.07 OpenTelemetry metric API: ${otel_metric_version}"
echo "P01.07 OpenTelemetry trace API: ${otel_trace_version}"
echo "P01.07 OpenTelemetry SDK: ${otel_sdk_version}"
echo "P01.07 OpenTelemetry metric SDK: ${otel_sdk_metric_version}"

go vet ./kernel/...
go test ./kernel/...
go test -race ./kernel/internal/observability -count=1
go test -v ./kernel/internal/observability -run 'Test(Settings|Strict|Capture|JSON|Correlation|InvalidCorrelation|Provider|DisabledProvider|ExporterFailure|Shutdown)' -count=1

module_prefix="github.com/Vertex-Systems-Network/omnexa"
while IFS= read -r package; do
  case "$package" in
    "$module_prefix/kernel/internal/observability"|"$module_prefix/kernel/internal/buildinfo"|"$module_prefix/kernel/internal/config"|"$module_prefix/kernel/internal/failure") ;;
    "$module_prefix"|"$module_prefix/"*)
      echo "ERROR: P01.07 observability package imports out-of-scope Omnexa package: ${package}" >&2
      exit 1
      ;;
  esac
done < <(go list -deps ./kernel/internal/observability)

if grep -R -nE 'github\.com/Vertex-Systems-Network/omnexa/(modules|platform)|kernel/internal/(cache|database|storage)' kernel/internal/observability --include='*.go'; then
  echo "ERROR: P01.07 observability source contains out-of-scope module/cache/database/storage coupling" >&2
  exit 1
fi

if grep -R -nE 'go\.opentelemetry\.io/otel/exporters/|github\.com/(datadog|newrelic|honeycombio|signalfx)|go\.uber\.org/zap|github\.com/rs/zerolog' kernel/internal/observability --include='*.go'; then
  echo "ERROR: P01.07 must remain vendor-neutral and use the governed slog/OTel SDK boundary" >&2
  exit 1
fi

if grep -R -nE 'type[[:space:]]+(Tenant|User|Customer|Order|Invoice|Product|AuditEvent|Health|Readiness|Job|Scheduler|FeatureFlag|Agent|Planner|Model|Embedding|Vector|Workflow|Event)([[:space:]]|$)|[[:space:]](TenantID|UserID|CustomerID|OrderID|AuditID|AgentID|GoalID|PlanID|TaskID|ExecutionID)[[:space:]]' kernel/internal/observability --include='*.go'; then
  echo "ERROR: P01.07 declares unauthorized tenant/business/audit/health/scheduler/AI concepts" >&2
  exit 1
fi

if grep -R -nE 'otel\.SetTracerProvider|otel\.SetMeterProvider|otel\.SetTextMapPropagator|propagation\.Baggage' kernel/internal/observability --include='*.go'; then
  echo "ERROR: P01.07 must not mutate OpenTelemetry globals or propagate arbitrary baggage" >&2
  exit 1
fi

go build ./kernel/...

echo "P01.07 G1 format/static/dependency/vendor boundary: PASS"
echo "P01.07 G2 unit/race logger/config/context/redaction: PASS"
echo "P01.07 G3 trace/metric SDK lifecycle and capture integration: PASS"
echo "P01.07 G5 secret/classification/error redaction and scope guard: PASS"
echo "P01.07 G6 bounded exporter failure/flush/shutdown and disabled-mode resilience: PASS"
echo "P01.07 G7 build/package: PASS"

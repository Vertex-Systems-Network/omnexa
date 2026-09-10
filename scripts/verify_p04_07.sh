#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

python scripts/validate_governance.py
python scripts/validate_p04_activation.py

expected_go="$(tr -d '[:space:]' < .go-version)"
actual_go="$(go env GOVERSION)"
if [[ "$actual_go" != "go${expected_go}" ]]; then
  echo "ERROR: Go toolchain mismatch: got ${actual_go}, want go${expected_go}" >&2
  exit 1
fi

required_files=(
  kernel/internal/events/schema_registry.go
  kernel/internal/events/schema_registry_test.go
  kernel/internal/events/schema_compatibility.go
  kernel/internal/events/schema_compatibility_test.go
  kernel/internal/events/schema_validation.go
  kernel/internal/events/schema_validation_test.go
  docs/roadmap/work-packages/P04.07.md
  docs/ai/P04.07_IMPLEMENTATION_PLAN.json
  docs/ai/handoffs/P04.07.md
  docs/roadmap/evidence/P04.07_FIRST_SLICE_2026-09-11.md
)
for file in "${required_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "ERROR: P04.07 first-slice required file missing: ${file}" >&2
    exit 1
  fi
done

unformatted="$(gofmt -l kernel/internal/events/schema_*.go)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required for P04.07 schema files:" >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go work sync
go mod tidy
if [[ -n "$(git status --porcelain -- go.mod go.work go.sum)" ]]; then
  echo "ERROR: P04.07 changed Go module/workspace metadata or requires an undeclared dependency" >&2
  git status --short -- go.mod go.work go.sum >&2
  git diff -- go.mod go.work go.sum >&2
  exit 1
fi

# Wave 1 is process-local and provider-neutral. No migration 4, hosted/provider
# registry dependency, broker SDK, remote-reference client, or workflow runtime is authorized.
if compgen -G 'kernel/migrations/kernel.events/4_*' > /dev/null; then
  echo "ERROR: P04.07 Wave 1 must not reserve or add kernel.events migration 4" >&2
  exit 1
fi

if grep -nE '"(net/http|net/url|github.com/segmentio/kafka-go|github.com/IBM/sarama|github.com/nats-io/|github.com/rabbitmq/|github.com/redis/|github.com/confluentinc/|go\.uber\.org/cadence|go\.temporal\.io)' \
  kernel/internal/events/schema_*.go; then
  echo "ERROR: P04.07 introduced unauthorized remote/provider/workflow dependency" >&2
  exit 1
fi

for marker in \
  'maxCanonicalSchemaBytes = 64 * 1024' \
  'maxSchemaDepth          = 8' \
  'maxSchemaFields         = 256' \
  'maxSchemaNodes          = 2048' \
  'type PayloadSchemaIdentity struct' \
  'type SchemaRegistry struct' \
  'func (registry *SchemaRegistry) Register(' \
  'func (registry *SchemaRegistry) Lookup(' \
  'func (registry *SchemaRegistry) Predecessor(' \
  'func CanonicalSchema('; do
  grep -Fq "$marker" kernel/internal/events/schema_registry.go || { echo "ERROR: P04.07 registry marker missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'CompatibilityBackward CompatibilityPolicy = "backward"' \
  'CompatibilityForward  CompatibilityPolicy = "forward"' \
  'CompatibilityFull     CompatibilityPolicy = "full"' \
  'CompatibilityExact    CompatibilityPolicy = "exact"' \
  'func EvaluateSchemaCompatibility(' \
  'func schemaLanguageContained('; do
  grep -Fq "$marker" kernel/internal/events/schema_compatibility.go || { echo "ERROR: P04.07 compatibility marker missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'maxSchemaPayloadBytes             = 64 * 1024' \
  'maxSchemaPayloadDepth             = 16' \
  'maxSchemaPayloadCollectionEntries = 256' \
  'maxSchemaPayloadStringRunes       = 4096' \
  'func ValidateRegisteredPayload(' \
  'func ValidatePayload('; do
  grep -Fq "$marker" kernel/internal/events/schema_validation.go || { echo "ERROR: P04.07 validation marker missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'TestCanonicalSchemaIsOrderIndependentAndDefensive' \
  'TestSchemaRegistryRegistrationIsImmutableAndIdempotent' \
  'TestSchemaRegistryPredecessorChoosesHighestAcceptedLowerVersion' \
  'TestCanonicalSchemaRejectsInvalidStructures' \
  'TestCanonicalSchemaEnforcesDepthAndNodeBounds'; do
  grep -Fq "$marker" kernel/internal/events/schema_registry_test.go || { echo "ERROR: P04.07 registry acceptance test missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'TestEvaluateSchemaCompatibilityExactUsesCanonicalContent' \
  'TestEvaluateSchemaCompatibilityDirectionalClosedObjectRules' \
  'TestEvaluateSchemaCompatibilityIntegerNumberContainment' \
  'TestEvaluateSchemaCompatibilityRecursesObjectsAndArrays' \
  'TestEvaluateSchemaCompatibilityFullRequiresBothDirections' \
  'TestEvaluateSchemaCompatibilityRejectsInvalidInputs'; do
  grep -Fq "$marker" kernel/internal/events/schema_compatibility_test.go || { echo "ERROR: P04.07 compatibility acceptance test missing: ${marker}" >&2; exit 1; }
done

for marker in \
  'TestValidatePayloadAcceptsClosedNestedContract' \
  'TestValidatePayloadRejectsUnknownMissingAndWrongTypes' \
  'TestValidatePayloadIntegerIsStrictSubsetOfNumber' \
  'TestValidatePayloadRejectsMalformedDuplicateNullRootAndTrailingValues' \
  'TestValidatePayloadEnforcesByteStringCollectionAndDepthLimits' \
  'TestValidatePayloadRejectsInvalidSchemaBeforePayloadUse' \
  'TestValidateRegisteredPayloadRequiresExactLocalRegistration' \
  'TestValidatePayloadArrayItemsRemainHomogeneous'; do
  grep -Fq "$marker" kernel/internal/events/schema_validation_test.go || { echo "ERROR: P04.07 validation acceptance test missing: ${marker}" >&2; exit 1; }
done

# Focused first-slice evidence plus retained event-package regression.
go test ./kernel/internal/events -run 'Schema|Payload' -count=1
go vet ./kernel/internal/events
go test ./kernel/internal/events -count=1
go test -race ./kernel/internal/events -run 'Schema|Payload' -count=1
go build ./kernel/...

echo "P04.07 G0 active-package governance + exact bounded first-slice authority: PASS"
echo "P04.07 G1 payload schema identity remains owner/event-type/version bound and separate from envelope version: PASS"
echo "P04.07 G2 immutable same-identity registration + deterministic canonical SHA-256 fingerprint: PASS"
echo "P04.07 G3 highest accepted lower-version predecessor selection remains deterministic: PASS"
echo "P04.07 G4 backward/forward/full/exact compatibility directions + recursive closed-object containment: PASS"
echo "P04.07 G5 integer is a strict subset of number without implicit payload coercion: PASS"
echo "P04.07 G6 exact local registration + closed-object payload validation fail closed before protected mutation: PASS"
echo "P04.07 G7 schema and payload size/depth/collection/string limits remain bounded: PASS"
echo "P04.07 G8 malformed/duplicate/null/trailing/unknown/unregistered payload or schema evidence fails closed: PASS"
echo "P04.07 G9 no migration 4, remote references, hosted/provider registry, broker SDK, workflow runtime, or external I/O: PASS"
echo "P04.07 G10 retained kernel.events vet/regression/race/build boundary: PASS"
echo "P04.07 G11 no P04.08+, business-feature, authorization-semantic, or AI/model/agent product-runtime expansion: PASS"

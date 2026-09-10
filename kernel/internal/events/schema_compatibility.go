package events

import "github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"

const codeSchemaCompatibilityInvalid failure.Code = "events.schema.compatibility_invalid"

// CompatibilityPolicy is the frozen P04.07 Wave 1 compatibility vocabulary.
type CompatibilityPolicy string

const (
	CompatibilityBackward CompatibilityPolicy = "backward"
	CompatibilityForward  CompatibilityPolicy = "forward"
	CompatibilityFull     CompatibilityPolicy = "full"
	CompatibilityExact    CompatibilityPolicy = "exact"
)

// Valid reports whether policy is part of the frozen P04.07 Wave 1 contract.
func (policy CompatibilityPolicy) Valid() bool {
	switch policy {
	case CompatibilityBackward, CompatibilityForward, CompatibilityFull, CompatibilityExact:
		return true
	default:
		return false
	}
}

// CompatibilityReason is a stable bounded machine-readable compatibility result.
type CompatibilityReason string

const (
	CompatibilityReasonCompatible        CompatibilityReason = "compatible"
	CompatibilityReasonExactMismatch     CompatibilityReason = "exact_mismatch"
	CompatibilityReasonBackwardViolation CompatibilityReason = "backward_violation"
	CompatibilityReasonForwardViolation  CompatibilityReason = "forward_violation"
)

// CompatibilityResult is contract evidence only. It grants no publication,
// authorization, protected-handler, or cross-owner authority.
type CompatibilityResult struct {
	Policy     CompatibilityPolicy
	Compatible bool
	Reason     CompatibilityReason
}

// EvaluateSchemaCompatibility deterministically compares one candidate schema
// with its required immediate predecessor under the frozen P04.07 Wave 1 policy.
// Both schemas are validated and normalized before comparison. The caller owns
// predecessor selection; SchemaRegistry.Predecessor provides the frozen
// highest-accepted-lower-version selection rule.
func EvaluateSchemaCompatibility(policy CompatibilityPolicy, candidate, predecessor Schema) (CompatibilityResult, error) {
	if !policy.Valid() {
		return CompatibilityResult{}, schemaFailure(codeSchemaCompatibilityInvalid, failure.CategoryValidation, "event schema compatibility policy is invalid")
	}

	normalizedCandidate, _, candidateFingerprint, err := normalizeSchema(candidate)
	if err != nil {
		return CompatibilityResult{}, err
	}
	normalizedPredecessor, _, predecessorFingerprint, err := normalizeSchema(predecessor)
	if err != nil {
		return CompatibilityResult{}, err
	}

	result := CompatibilityResult{Policy: policy}
	switch policy {
	case CompatibilityExact:
		if candidateFingerprint != predecessorFingerprint {
			result.Reason = CompatibilityReasonExactMismatch
			return result, nil
		}
	case CompatibilityBackward:
		if !schemaLanguageContained(normalizedPredecessor, normalizedCandidate) {
			result.Reason = CompatibilityReasonBackwardViolation
			return result, nil
		}
	case CompatibilityForward:
		if !schemaLanguageContained(normalizedCandidate, normalizedPredecessor) {
			result.Reason = CompatibilityReasonForwardViolation
			return result, nil
		}
	case CompatibilityFull:
		if !schemaLanguageContained(normalizedPredecessor, normalizedCandidate) {
			result.Reason = CompatibilityReasonBackwardViolation
			return result, nil
		}
		if !schemaLanguageContained(normalizedCandidate, normalizedPredecessor) {
			result.Reason = CompatibilityReasonForwardViolation
			return result, nil
		}
	}

	result.Compatible = true
	result.Reason = CompatibilityReasonCompatible
	return result, nil
}

// schemaLanguageContained reports whether every payload accepted by source is
// also accepted by target. Both inputs must already be normalized valid schemas.
func schemaLanguageContained(source, target Schema) bool {
	return fieldsLanguageContained(source.Fields, target.Fields)
}

func fieldsLanguageContained(source, target []Field) bool {
	sourceByName := indexFields(source)
	targetByName := indexFields(target)

	// A source property that target does not define can appear in a valid source
	// payload but is rejected by the target's closed-object contract.
	for name, sourceField := range sourceByName {
		targetField, ok := targetByName[name]
		if !ok || !fieldLanguageContained(sourceField, targetField) {
			return false
		}
	}

	// If target requires a property, every source-valid payload must contain it.
	for name, targetField := range targetByName {
		if !targetField.Required {
			continue
		}
		sourceField, ok := sourceByName[name]
		if !ok || !sourceField.Required {
			return false
		}
	}
	return true
}

func fieldLanguageContained(source, target Field) bool {
	if source.Kind == ValueInteger && target.Kind == ValueNumber {
		return true
	}
	if source.Kind != target.Kind {
		return false
	}

	switch source.Kind {
	case ValueObject:
		return fieldsLanguageContained(source.Fields, target.Fields)
	case ValueArray:
		return source.Items != nil && target.Items != nil && fieldLanguageContained(*source.Items, *target.Items)
	default:
		return true
	}
}

func indexFields(fields []Field) map[string]Field {
	indexed := make(map[string]Field, len(fields))
	for index := range fields {
		indexed[fields[index].Name] = fields[index]
	}
	return indexed
}

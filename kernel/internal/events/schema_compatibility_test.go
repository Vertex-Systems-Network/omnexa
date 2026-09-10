package events

import (
	"reflect"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestEvaluateSchemaCompatibilityExactUsesCanonicalContent(t *testing.T) {
	predecessor := Schema{Fields: []Field{
		{Name: "zeta", Kind: ValueInteger},
		{Name: "alpha", Kind: ValueString, Required: true},
	}}
	candidate := Schema{Fields: []Field{
		{Name: "alpha", Kind: ValueString, Required: true},
		{Name: "zeta", Kind: ValueInteger},
	}}

	result, err := EvaluateSchemaCompatibility(CompatibilityExact, candidate, predecessor)
	if err != nil {
		t.Fatalf("exact compatibility returned error: %v", err)
	}
	if !result.Compatible || result.Reason != CompatibilityReasonCompatible {
		t.Fatalf("expected canonical exact match, got %+v", result)
	}

	candidate.Fields[1].Required = true
	result, err = EvaluateSchemaCompatibility(CompatibilityExact, candidate, predecessor)
	if err != nil {
		t.Fatalf("exact mismatch returned error: %v", err)
	}
	if result.Compatible || result.Reason != CompatibilityReasonExactMismatch {
		t.Fatalf("expected exact mismatch, got %+v", result)
	}
}

func TestEvaluateSchemaCompatibilityDirectionalClosedObjectRules(t *testing.T) {
	base := Schema{Fields: []Field{{Name: "id", Kind: ValueString, Required: true}}}
	withOptional := Schema{Fields: []Field{
		{Name: "id", Kind: ValueString, Required: true},
		{Name: "note", Kind: ValueString},
	}}
	withRequired := Schema{Fields: []Field{
		{Name: "id", Kind: ValueString, Required: true},
		{Name: "note", Kind: ValueString, Required: true},
	}}

	assertCompatibility(t, CompatibilityBackward, withOptional, base, true, CompatibilityReasonCompatible)
	assertCompatibility(t, CompatibilityForward, withOptional, base, false, CompatibilityReasonForwardViolation)
	assertCompatibility(t, CompatibilityFull, withOptional, base, false, CompatibilityReasonForwardViolation)
	assertCompatibility(t, CompatibilityBackward, withRequired, base, false, CompatibilityReasonBackwardViolation)
	assertCompatibility(t, CompatibilityBackward, base, withOptional, false, CompatibilityReasonBackwardViolation)
	assertCompatibility(t, CompatibilityForward, base, withOptional, true, CompatibilityReasonCompatible)
}

func TestEvaluateSchemaCompatibilityIntegerNumberContainment(t *testing.T) {
	integerSchema := Schema{Fields: []Field{{Name: "value", Kind: ValueInteger, Required: true}}}
	numberSchema := Schema{Fields: []Field{{Name: "value", Kind: ValueNumber, Required: true}}}

	assertCompatibility(t, CompatibilityBackward, numberSchema, integerSchema, true, CompatibilityReasonCompatible)
	assertCompatibility(t, CompatibilityForward, numberSchema, integerSchema, false, CompatibilityReasonForwardViolation)
	assertCompatibility(t, CompatibilityBackward, integerSchema, numberSchema, false, CompatibilityReasonBackwardViolation)
	assertCompatibility(t, CompatibilityForward, integerSchema, numberSchema, true, CompatibilityReasonCompatible)
}

func TestEvaluateSchemaCompatibilityRecursesObjectsAndArrays(t *testing.T) {
	predecessor := Schema{Fields: []Field{
		{
			Name:     "profile",
			Kind:     ValueObject,
			Required: true,
			Fields: []Field{
				{Name: "name", Kind: ValueString, Required: true},
			},
		},
		{
			Name:     "scores",
			Kind:     ValueArray,
			Required: true,
			Items:    &Field{Kind: ValueInteger},
		},
	}}
	candidate := Schema{Fields: []Field{
		{
			Name:     "profile",
			Kind:     ValueObject,
			Required: true,
			Fields: []Field{
				{Name: "name", Kind: ValueString, Required: true},
				{Name: "label", Kind: ValueString},
			},
		},
		{
			Name:     "scores",
			Kind:     ValueArray,
			Required: true,
			Items:    &Field{Kind: ValueNumber},
		},
	}}

	assertCompatibility(t, CompatibilityBackward, candidate, predecessor, true, CompatibilityReasonCompatible)
	assertCompatibility(t, CompatibilityForward, candidate, predecessor, false, CompatibilityReasonForwardViolation)
}

func TestEvaluateSchemaCompatibilityFullRequiresBothDirections(t *testing.T) {
	schema := Schema{Fields: []Field{{Name: "id", Kind: ValueString, Required: true}}}
	assertCompatibility(t, CompatibilityFull, schema, schema, true, CompatibilityReasonCompatible)

	candidate := Schema{Fields: []Field{{Name: "id", Kind: ValueNumber, Required: true}}}
	predecessor := Schema{Fields: []Field{{Name: "id", Kind: ValueInteger, Required: true}}}
	assertCompatibility(t, CompatibilityFull, candidate, predecessor, false, CompatibilityReasonForwardViolation)
}

func TestEvaluateSchemaCompatibilityRejectsInvalidInputs(t *testing.T) {
	valid := Schema{Fields: []Field{{Name: "id", Kind: ValueString}}}

	_, err := EvaluateSchemaCompatibility(CompatibilityPolicy("provider-default"), valid, valid)
	if got := failure.CodeOf(err); got != codeSchemaCompatibilityInvalid {
		t.Fatalf("invalid policy code = %q, want %q", got, codeSchemaCompatibilityInvalid)
	}

	invalid := Schema{Fields: []Field{{Name: " id", Kind: ValueString}}}
	_, err = EvaluateSchemaCompatibility(CompatibilityBackward, invalid, valid)
	if got := failure.CodeOf(err); got != codeSchemaInvalid {
		t.Fatalf("invalid candidate code = %q, want %q", got, codeSchemaInvalid)
	}

	_, err = EvaluateSchemaCompatibility(CompatibilityBackward, valid, invalid)
	if got := failure.CodeOf(err); got != codeSchemaInvalid {
		t.Fatalf("invalid predecessor code = %q, want %q", got, codeSchemaInvalid)
	}
}

func TestEvaluateSchemaCompatibilityDoesNotMutateCallerSchemas(t *testing.T) {
	candidate := Schema{Fields: []Field{
		{Name: "zeta", Kind: ValueString},
		{Name: "alpha", Kind: ValueObject, Fields: []Field{
			{Name: "two", Kind: ValueBoolean},
			{Name: "one", Kind: ValueInteger},
		}},
	}}
	predecessor := cloneSchema(candidate)
	beforeCandidate := cloneSchema(candidate)
	beforePredecessor := cloneSchema(predecessor)

	if _, err := EvaluateSchemaCompatibility(CompatibilityExact, candidate, predecessor); err != nil {
		t.Fatalf("compatibility returned error: %v", err)
	}
	if !reflect.DeepEqual(candidate, beforeCandidate) {
		t.Fatal("candidate schema was mutated")
	}
	if !reflect.DeepEqual(predecessor, beforePredecessor) {
		t.Fatal("predecessor schema was mutated")
	}
}

func assertCompatibility(t *testing.T, policy CompatibilityPolicy, candidate, predecessor Schema, wantCompatible bool, wantReason CompatibilityReason) {
	t.Helper()
	result, err := EvaluateSchemaCompatibility(policy, candidate, predecessor)
	if err != nil {
		t.Fatalf("policy %q returned error: %v", policy, err)
	}
	if result.Policy != policy || result.Compatible != wantCompatible || result.Reason != wantReason {
		t.Fatalf("policy %q result = %+v, want compatible=%v reason=%q", policy, result, wantCompatible, wantReason)
	}
}

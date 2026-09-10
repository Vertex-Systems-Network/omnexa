package events

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestValidatePayloadAcceptsClosedNestedContract(t *testing.T) {
	schema := Schema{Fields: []Field{
		{Name: "id", Kind: ValueString, Required: true},
		{Name: "count", Kind: ValueInteger, Required: true},
		{Name: "ratio", Kind: ValueNumber},
		{Name: "active", Kind: ValueBoolean},
		{
			Name: "profile",
			Kind: ValueObject,
			Fields: []Field{
				{Name: "name", Kind: ValueString, Required: true},
			},
		},
		{
			Name:  "scores",
			Kind:  ValueArray,
			Items: &Field{Kind: ValueNumber},
		},
	}}

	payload := []byte(`{"id":"evt-1","count":2,"ratio":2.5,"active":true,"profile":{"name":"Ada"},"scores":[1,2.5,3]}`)
	if err := ValidatePayload(schema, payload); err != nil {
		t.Fatalf("valid payload rejected: %v", err)
	}
}

func TestValidatePayloadRejectsUnknownMissingAndWrongTypes(t *testing.T) {
	schema := Schema{Fields: []Field{
		{Name: "id", Kind: ValueString, Required: true},
		{Name: "active", Kind: ValueBoolean},
	}}

	cases := [][]byte{
		[]byte(`{"id":"ok","extra":true}`),
		[]byte(`{"active":true}`),
		[]byte(`{"id":7}`),
		[]byte(`{"id":"ok","active":"true"}`),
	}
	for _, payload := range cases {
		err := ValidatePayload(schema, payload)
		if got := failure.CodeOf(err); got != codePayloadInvalid {
			t.Fatalf("payload %q code = %q, want %q", payload, got, codePayloadInvalid)
		}
	}
}

func TestValidatePayloadIntegerIsStrictSubsetOfNumber(t *testing.T) {
	integerSchema := Schema{Fields: []Field{{Name: "value", Kind: ValueInteger, Required: true}}}
	numberSchema := Schema{Fields: []Field{{Name: "value", Kind: ValueNumber, Required: true}}}

	if err := ValidatePayload(integerSchema, []byte(`{"value":-12}`)); err != nil {
		t.Fatalf("integer token rejected: %v", err)
	}
	for _, payload := range [][]byte{[]byte(`{"value":1.0}`), []byte(`{"value":1e0}`)} {
		if got := failure.CodeOf(ValidatePayload(integerSchema, payload)); got != codePayloadInvalid {
			t.Fatalf("non-integer token %q code = %q, want %q", payload, got, codePayloadInvalid)
		}
	}
	for _, payload := range [][]byte{[]byte(`{"value":12}`), []byte(`{"value":1.25}`), []byte(`{"value":1e2}`)} {
		if err := ValidatePayload(numberSchema, payload); err != nil {
			t.Fatalf("number payload %q rejected: %v", payload, err)
		}
	}
}

func TestValidatePayloadRejectsMalformedDuplicateNullRootAndTrailingValues(t *testing.T) {
	schema := Schema{Fields: []Field{{Name: "id", Kind: ValueString}}}
	cases := [][]byte{
		[]byte(`{"id":"a","id":"b"}`),
		[]byte(`{"id":null}`),
		[]byte(`["a"]`),
		[]byte(`{"id":"a"} {"id":"b"}`),
		[]byte(`{"id":`),
		{0xff, 0xfe},
		nil,
	}
	for _, payload := range cases {
		if got := failure.CodeOf(ValidatePayload(schema, payload)); got != codePayloadInvalid {
			t.Fatalf("payload %q code = %q, want %q", payload, got, codePayloadInvalid)
		}
	}
}

func TestValidatePayloadEnforcesByteStringCollectionAndDepthLimits(t *testing.T) {
	stringSchema := Schema{Fields: []Field{{Name: "value", Kind: ValueString}}}
	oversizedPayload := []byte(strings.Repeat(" ", maxPayloadBytes+1))
	if got := failure.CodeOf(ValidatePayload(stringSchema, oversizedPayload)); got != codePayloadLimit {
		t.Fatalf("oversized bytes code = %q, want %q", got, codePayloadLimit)
	}

	longString := []byte(`{"value":"` + strings.Repeat("x", maxPayloadStringRunes+1) + `"}`)
	if got := failure.CodeOf(ValidatePayload(stringSchema, longString)); got != codePayloadLimit {
		t.Fatalf("long string code = %q, want %q", got, codePayloadLimit)
	}

	arraySchema := Schema{Fields: []Field{{Name: "value", Kind: ValueArray, Items: &Field{Kind: ValueInteger}}}}
	items := make([]string, maxPayloadCollectionEntries+1)
	for index := range items {
		items[index] = "1"
	}
	largeArray := []byte(`{"value":[` + strings.Join(items, ",") + `]}`)
	if got := failure.CodeOf(ValidatePayload(arraySchema, largeArray)); got != codePayloadLimit {
		t.Fatalf("large array code = %q, want %q", got, codePayloadLimit)
	}

	var object strings.Builder
	object.WriteByte('{')
	for index := 0; index <= maxPayloadCollectionEntries; index++ {
		if index > 0 {
			object.WriteByte(',')
		}
		fmt.Fprintf(&object, `"k%d":true`, index)
	}
	object.WriteByte('}')
	if got := failure.CodeOf(ValidatePayload(Schema{}, []byte(object.String()))); got != codePayloadLimit {
		t.Fatalf("large object code = %q, want %q", got, codePayloadLimit)
	}

	tooDeep := []byte(`{"value":` + strings.Repeat("[", maxPayloadDepth) + `0` + strings.Repeat("]", maxPayloadDepth) + `}`)
	if got := failure.CodeOf(ValidatePayload(arraySchema, tooDeep)); got != codePayloadLimit {
		t.Fatalf("deep payload code = %q, want %q", got, codePayloadLimit)
	}
}

func TestValidatePayloadRejectsInvalidSchemaBeforePayloadUse(t *testing.T) {
	invalid := Schema{Fields: []Field{{Name: " value", Kind: ValueString}}}
	err := ValidatePayload(invalid, []byte(`{"value":"secret"}`))
	if got := failure.CodeOf(err); got != codeSchemaInvalid {
		t.Fatalf("invalid schema code = %q, want %q", got, codeSchemaInvalid)
	}
}

func TestValidateRegisteredPayloadRequiresExactLocalRegistration(t *testing.T) {
	registry := NewSchemaRegistry()
	identity := PayloadSchemaIdentity{
		Owner:     Producer("urn:omnexa:module:kernel.events"),
		EventType: EventType("kernel.events.schema_changed.v1"),
		Version:   1,
	}
	schema := Schema{Fields: []Field{{Name: "id", Kind: ValueString, Required: true}}}
	if _, err := registry.Register(identity, schema); err != nil {
		t.Fatalf("register: %v", err)
	}

	if err := ValidateRegisteredPayload(registry, identity, []byte(`{"id":"evt-1"}`)); err != nil {
		t.Fatalf("registered payload rejected: %v", err)
	}

	unknown := identity
	unknown.Version = 2
	if got := failure.CodeOf(ValidateRegisteredPayload(registry, unknown, []byte(`{"id":"evt-1"}`))); got != codeSchemaNotFound {
		t.Fatalf("unknown schema code = %q, want %q", got, codeSchemaNotFound)
	}

	if got := failure.CodeOf(ValidateRegisteredPayload(nil, identity, []byte(`{"id":"evt-1"}`))); got != codeSchemaInvalid {
		t.Fatalf("nil registry code = %q, want %q", got, codeSchemaInvalid)
	}
}

func TestValidatePayloadArrayItemsRemainHomogeneous(t *testing.T) {
	schema := Schema{Fields: []Field{{
		Name:     "items",
		Kind:     ValueArray,
		Required: true,
		Items: &Field{
			Kind: ValueObject,
			Fields: []Field{
				{Name: "enabled", Kind: ValueBoolean, Required: true},
			},
		},
	}}}

	if err := ValidatePayload(schema, []byte(`{"items":[{"enabled":true},{"enabled":false}]}`)); err != nil {
		t.Fatalf("homogeneous object array rejected: %v", err)
	}
	if got := failure.CodeOf(ValidatePayload(schema, []byte(`{"items":[{"enabled":true},{"enabled":"false"}]}`))); got != codePayloadInvalid {
		t.Fatalf("heterogeneous item code = %q, want %q", got, codePayloadInvalid)
	}
}

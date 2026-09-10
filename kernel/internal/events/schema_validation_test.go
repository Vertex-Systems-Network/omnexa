package events

import (
	"fmt"
	"strings"
	"testing"
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
		assertSchemaFailureCode(t, ValidatePayload(schema, payload), codePayloadInvalid)
	}
}

func TestValidatePayloadIntegerIsStrictSubsetOfNumber(t *testing.T) {
	integerSchema := Schema{Fields: []Field{{Name: "value", Kind: ValueInteger, Required: true}}}
	numberSchema := Schema{Fields: []Field{{Name: "value", Kind: ValueNumber, Required: true}}}

	if err := ValidatePayload(integerSchema, []byte(`{"value":-12}`)); err != nil {
		t.Fatalf("integer token rejected: %v", err)
	}
	for _, payload := range [][]byte{[]byte(`{"value":1.0}`), []byte(`{"value":1e0}`)} {
		assertSchemaFailureCode(t, ValidatePayload(integerSchema, payload), codePayloadInvalid)
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
		assertSchemaFailureCode(t, ValidatePayload(schema, payload), codePayloadInvalid)
	}
}

func TestValidatePayloadEnforcesByteStringCollectionAndDepthLimits(t *testing.T) {
	stringSchema := Schema{Fields: []Field{{Name: "value", Kind: ValueString}}}
	oversizedPayload := []byte(strings.Repeat(" ", maxPayloadBytes+1))
	assertSchemaFailureCode(t, ValidatePayload(stringSchema, oversizedPayload), codePayloadLimit)

	longString := []byte(`{"value":"` + strings.Repeat("x", maxPayloadStringRunes+1) + `"}`)
	assertSchemaFailureCode(t, ValidatePayload(stringSchema, longString), codePayloadLimit)

	arraySchema := Schema{Fields: []Field{{Name: "value", Kind: ValueArray, Items: &Field{Kind: ValueInteger}}}}
	items := make([]string, maxPayloadCollectionEntries+1)
	for index := range items {
		items[index] = "1"
	}
	largeArray := []byte(`{"value":[` + strings.Join(items, ",") + `]}`)
	assertSchemaFailureCode(t, ValidatePayload(arraySchema, largeArray), codePayloadLimit)

	var object strings.Builder
	object.WriteByte('{')
	for index := 0; index <= maxPayloadCollectionEntries; index++ {
		if index > 0 {
			object.WriteByte(',')
		}
		fmt.Fprintf(&object, `"k%d":true`, index)
	}
	object.WriteByte('}')
	assertSchemaFailureCode(t, ValidatePayload(Schema{}, []byte(object.String())), codePayloadLimit)

	tooDeep := []byte(`{"value":` + strings.Repeat("[", maxPayloadDepth) + `0` + strings.Repeat("]", maxPayloadDepth) + `}`)
	assertSchemaFailureCode(t, ValidatePayload(arraySchema, tooDeep), codePayloadLimit)
}

func TestValidatePayloadRejectsInvalidSchemaBeforePayloadUse(t *testing.T) {
	invalid := Schema{Fields: []Field{{Name: " value", Kind: ValueString}}}
	assertSchemaFailureCode(t, ValidatePayload(invalid, []byte(`{"value":"secret"}`)), codeSchemaInvalid)
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
	assertSchemaFailureCode(t, ValidateRegisteredPayload(registry, unknown, []byte(`{"id":"evt-1"}`)), codeSchemaNotFound)
	assertSchemaFailureCode(t, ValidateRegisteredPayload(nil, identity, []byte(`{"id":"evt-1"}`)), codeSchemaInvalid)
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
	assertSchemaFailureCode(t, ValidatePayload(schema, []byte(`{"items":[{"enabled":true},{"enabled":"false"}]}`)), codePayloadInvalid)
}

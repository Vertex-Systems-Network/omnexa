package events

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestCanonicalSchemaIsOrderIndependentAndDefensive(t *testing.T) {
	t.Parallel()
	left := Schema{Fields: []Field{
		{Name: "zeta", Kind: ValueString},
		{Name: "account", Kind: ValueObject, Required: true, Fields: []Field{
			{Name: "name", Kind: ValueString, Required: true},
			{Name: "age", Kind: ValueInteger},
		}},
	}}
	right := Schema{Fields: []Field{
		{Name: "account", Kind: ValueObject, Required: true, Fields: []Field{
			{Name: "age", Kind: ValueInteger},
			{Name: "name", Kind: ValueString, Required: true},
		}},
		{Name: "zeta", Kind: ValueString},
	}}

	leftBytes, leftFingerprint, err := CanonicalSchema(left)
	if err != nil {
		t.Fatalf("canonicalize left: %v", err)
	}
	rightBytes, rightFingerprint, err := CanonicalSchema(right)
	if err != nil {
		t.Fatalf("canonicalize right: %v", err)
	}
	if !bytes.Equal(leftBytes, rightBytes) {
		t.Fatalf("canonical bytes differ:\nleft=%s\nright=%s", leftBytes, rightBytes)
	}
	if leftFingerprint != rightFingerprint || len(leftFingerprint) != 64 {
		t.Fatalf("unexpected fingerprints: %q %q", leftFingerprint, rightFingerprint)
	}
	if left.Fields[0].Name != "zeta" || left.Fields[1].Fields[0].Name != "name" {
		t.Fatal("canonicalization mutated caller-owned schema ordering")
	}
}

func TestSchemaRegistryRegistrationIsImmutableAndIdempotent(t *testing.T) {
	t.Parallel()
	registry := NewSchemaRegistry()
	identity := testSchemaIdentity(2)
	schema := Schema{Fields: []Field{{Name: "id", Kind: ValueString, Required: true}}}

	first, err := registry.Register(identity, schema)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	second, err := registry.Register(identity, Schema{Fields: []Field{{Name: "id", Kind: ValueString, Required: true}}})
	if err != nil {
		t.Fatalf("idempotent register: %v", err)
	}
	if first.Fingerprint != second.Fingerprint || !bytes.Equal(first.Canonical, second.Canonical) {
		t.Fatal("idempotent registration returned different immutable content")
	}

	first.Canonical[0] = 'X'
	first.Schema.Fields[0].Name = "mutated"
	lookedUp, err := registry.Lookup(identity)
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if lookedUp.Schema.Fields[0].Name != "id" || len(lookedUp.Canonical) == 0 || lookedUp.Canonical[0] == 'X' {
		t.Fatal("registry leaked caller-owned mutable state")
	}

	_, err = registry.Register(identity, Schema{Fields: []Field{{Name: "id", Kind: ValueInteger, Required: true}}})
	assertSchemaFailureCode(t, err, codeSchemaConflict)
}

func TestSchemaRegistryPredecessorChoosesHighestAcceptedLowerVersion(t *testing.T) {
	t.Parallel()
	registry := NewSchemaRegistry()
	for _, version := range []uint{1, 2, 4} {
		_, err := registry.Register(testSchemaIdentity(version), Schema{Fields: []Field{{Name: "version", Kind: ValueInteger}}})
		if err != nil {
			t.Fatalf("register v%d: %v", version, err)
		}
	}
	predecessor, found, err := registry.Predecessor(testSchemaIdentity(4))
	if err != nil {
		t.Fatalf("predecessor v4: %v", err)
	}
	if !found || predecessor.Identity.Version != 2 {
		t.Fatalf("expected v2 predecessor, got found=%v version=%d", found, predecessor.Identity.Version)
	}
	_, found, err = registry.Predecessor(testSchemaIdentity(1))
	if err != nil || found {
		t.Fatalf("v1 must have no predecessor: found=%v err=%v", found, err)
	}
}

func TestSchemaRegistryRejectsInvalidIdentityAndMissingLookup(t *testing.T) {
	t.Parallel()
	registry := NewSchemaRegistry()
	_, err := registry.Register(PayloadSchemaIdentity{}, Schema{})
	assertSchemaFailureCode(t, err, codeSchemaInvalid)

	_, err = registry.Lookup(testSchemaIdentity(1))
	assertSchemaFailureCode(t, err, codeSchemaNotFound)

	var nilRegistry *SchemaRegistry
	_, err = nilRegistry.Lookup(testSchemaIdentity(1))
	assertSchemaFailureCode(t, err, codeSchemaInvalid)
}

func TestCanonicalSchemaRejectsInvalidStructures(t *testing.T) {
	t.Parallel()
	tooLongName := strings.Repeat("a", maxSchemaFieldNameRunes+1)
	cases := map[string]Schema{
		"duplicate field": {Fields: []Field{{Name: "id", Kind: ValueString}, {Name: "id", Kind: ValueString}}},
		"trimmed field":   {Fields: []Field{{Name: " id", Kind: ValueString}}},
		"oversize name":   {Fields: []Field{{Name: tooLongName, Kind: ValueString}}},
		"unknown kind":    {Fields: []Field{{Name: "id", Kind: ValueKind("uuid")}}},
		"primitive child": {Fields: []Field{{Name: "id", Kind: ValueString, Fields: []Field{{Name: "nested", Kind: ValueString}}}}},
		"object item":     {Fields: []Field{{Name: "obj", Kind: ValueObject, Items: &Field{Kind: ValueString}}}},
		"array no item":   {Fields: []Field{{Name: "items", Kind: ValueArray}}},
		"named array item": {Fields: []Field{{Name: "items", Kind: ValueArray, Items: &Field{
			Name: "forbidden", Kind: ValueString,
		}}}},
	}
	for name, schema := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, _, err := CanonicalSchema(schema)
			assertSchemaFailureCode(t, err, codeSchemaInvalid)
		})
	}
}

func TestCanonicalSchemaSupportsNestedHomogeneousArrays(t *testing.T) {
	t.Parallel()
	arrayOfObjects := Field{Kind: ValueObject, Fields: []Field{
		{Name: "enabled", Kind: ValueBoolean, Required: true},
		{Name: "score", Kind: ValueNumber},
	}}
	schema := Schema{Fields: []Field{{Name: "items", Kind: ValueArray, Required: true, Items: &arrayOfObjects}}}
	canonical, fingerprint, err := CanonicalSchema(schema)
	if err != nil {
		t.Fatalf("canonicalize array schema: %v", err)
	}
	if len(canonical) == 0 || len(fingerprint) != 64 {
		t.Fatalf("missing canonical result: %q %q", canonical, fingerprint)
	}
}

func TestCanonicalSchemaEnforcesDepthAndNodeBounds(t *testing.T) {
	t.Parallel()
	item := &Field{Kind: ValueString}
	for i := 0; i < maxSchemaDepth+1; i++ {
		item = &Field{Kind: ValueArray, Items: item}
	}
	_, _, err := CanonicalSchema(Schema{Fields: []Field{{Name: "deep", Kind: ValueArray, Items: item}}})
	assertSchemaFailureCode(t, err, codeSchemaInvalid)

	fields := make([]Field, maxSchemaFields+1)
	for index := range fields {
		fields[index] = Field{Name: strings.Repeat("a", index/26) + string(rune('a'+index%26)), Kind: ValueString}
	}
	_, _, err = CanonicalSchema(Schema{Fields: fields})
	assertSchemaFailureCode(t, err, codeSchemaInvalid)
}

func testSchemaIdentity(version uint) PayloadSchemaIdentity {
	return PayloadSchemaIdentity{
		Owner:     Producer("urn:omnexa:module:kernel.events"),
		EventType: EventType("kernel.events.schema_changed.v1"),
		Version:   version,
	}
}

func assertSchemaFailureCode(t *testing.T, err error, expected failure.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected failure code %s, got nil", expected)
	}
	code, ok := failure.CodeOf(err)
	if !ok || code != expected {
		t.Fatalf("expected failure code %s, got %v (%v)", expected, code, err)
	}
}

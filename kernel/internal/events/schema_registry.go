package events

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	maxCanonicalSchemaBytes = 64 * 1024
	maxSchemaDepth          = 8
	maxSchemaFields         = 256
	maxSchemaNodes          = 2048
	maxSchemaFieldNameRunes = 128
)

const (
	codeSchemaInvalid  failure.Code = "events.schema.invalid"
	codeSchemaConflict failure.Code = "events.schema.conflict"
	codeSchemaNotFound failure.Code = "events.schema.not_found"
)

// ValueKind is the provider-neutral P04.07 payload value vocabulary.
type ValueKind string

const (
	ValueString  ValueKind = "string"
	ValueInteger ValueKind = "integer"
	ValueNumber  ValueKind = "number"
	ValueBoolean ValueKind = "boolean"
	ValueObject  ValueKind = "object"
	ValueArray   ValueKind = "array"
)

// Valid reports whether kind is part of the frozen P04.07 Wave 1 vocabulary.
func (kind ValueKind) Valid() bool {
	switch kind {
	case ValueString, ValueInteger, ValueNumber, ValueBoolean, ValueObject, ValueArray:
		return true
	default:
		return false
	}
}

// Schema is one closed object payload contract. Unknown payload properties are
// rejected by the validation layer; Fields therefore describe the complete root.
type Schema struct {
	Fields []Field `json:"fields"`
}

// Field describes one object property. Object values use Fields. Array values
// use Items as a nameless homogeneous item descriptor. Primitive values use neither.
type Field struct {
	Name     string    `json:"name,omitempty"`
	Kind     ValueKind `json:"kind"`
	Required bool      `json:"required,omitempty"`
	Fields   []Field   `json:"fields,omitempty"`
	Items    *Field    `json:"items,omitempty"`
}

// PayloadSchemaIdentity is the immutable registry key for one payload schema version.
// Owner and EventType are transport-neutral event identities; Version is independent
// from the envelope/event contract version and must be positive.
type PayloadSchemaIdentity struct {
	Owner     Producer
	EventType EventType
	Version   uint
}

// RegisteredSchema is an immutable defensive snapshot returned by SchemaRegistry.
type RegisteredSchema struct {
	Identity    PayloadSchemaIdentity
	Schema      Schema
	Fingerprint string
	Canonical   []byte
}

// SchemaRegistry is the P04.07 Wave 1 process-local immutable schema registry.
// It makes no durability or external-provider claim.
type SchemaRegistry struct {
	mu      sync.RWMutex
	entries map[PayloadSchemaIdentity]RegisteredSchema
}

// NewSchemaRegistry creates an empty process-local registry.
func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{entries: make(map[PayloadSchemaIdentity]RegisteredSchema)}
}

// Register validates and canonicalizes a schema before immutable registration.
// Registering identical canonical content at the same identity is idempotent;
// registering different content at that identity fails closed.
func (registry *SchemaRegistry) Register(identity PayloadSchemaIdentity, schema Schema) (RegisteredSchema, error) {
	if registry == nil || !identity.valid() {
		return RegisteredSchema{}, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema identity is invalid")
	}
	normalized, canonical, fingerprint, err := normalizeSchema(schema)
	if err != nil {
		return RegisteredSchema{}, err
	}
	candidate := RegisteredSchema{
		Identity:    identity,
		Schema:      normalized,
		Fingerprint: fingerprint,
		Canonical:   canonical,
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()
	if existing, ok := registry.entries[identity]; ok {
		if existing.Fingerprint != fingerprint || !equalBytes(existing.Canonical, canonical) {
			return RegisteredSchema{}, schemaFailure(codeSchemaConflict, failure.CategoryConflict, "event schema identity already has different content")
		}
		return cloneRegisteredSchema(existing), nil
	}
	registry.entries[identity] = cloneRegisteredSchema(candidate)
	return cloneRegisteredSchema(candidate), nil
}

// Lookup returns a defensive immutable snapshot for the exact identity.
func (registry *SchemaRegistry) Lookup(identity PayloadSchemaIdentity) (RegisteredSchema, error) {
	if registry == nil || !identity.valid() {
		return RegisteredSchema{}, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema identity is invalid")
	}
	registry.mu.RLock()
	entry, ok := registry.entries[identity]
	registry.mu.RUnlock()
	if !ok {
		return RegisteredSchema{}, schemaFailure(codeSchemaNotFound, failure.CategoryNotFound, "event schema is not registered")
	}
	return cloneRegisteredSchema(entry), nil
}

// Predecessor returns the highest accepted version lower than identity.Version
// for the same owner and event type. Version 1, or a key with no predecessor,
// returns found=false without manufacturing a compatibility requirement.
func (registry *SchemaRegistry) Predecessor(identity PayloadSchemaIdentity) (entry RegisteredSchema, found bool, err error) {
	if registry == nil || !identity.valid() {
		return RegisteredSchema{}, false, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema identity is invalid")
	}
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	var selected RegisteredSchema
	var selectedVersion uint
	for key, candidate := range registry.entries {
		if key.Owner != identity.Owner || key.EventType != identity.EventType || key.Version >= identity.Version {
			continue
		}
		if !found || key.Version > selectedVersion {
			selected = candidate
			selectedVersion = key.Version
			found = true
		}
	}
	if !found {
		return RegisteredSchema{}, false, nil
	}
	return cloneRegisteredSchema(selected), true, nil
}

// CanonicalSchema validates a schema and returns deterministic canonical bytes
// plus its lowercase SHA-256 fingerprint without registering it.
func CanonicalSchema(schema Schema) (canonical []byte, fingerprint string, err error) {
	_, canonical, fingerprint, err = normalizeSchema(schema)
	return canonical, fingerprint, err
}

func (identity PayloadSchemaIdentity) valid() bool {
	return identity.Owner.Valid() && identity.EventType.Valid() && identity.Version > 0
}

func normalizeSchema(schema Schema) (Schema, []byte, string, error) {
	nodes := 1
	normalizedFields, err := normalizeFields(schema.Fields, 1, &nodes)
	if err != nil {
		return Schema{}, nil, "", err
	}
	normalized := Schema{Fields: normalizedFields}
	canonical, marshalErr := json.Marshal(normalized)
	if marshalErr != nil || len(canonical) > maxCanonicalSchemaBytes {
		return Schema{}, nil, "", schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
	}
	digest := sha256.Sum256(canonical)
	return normalized, canonical, hex.EncodeToString(digest[:]), nil
}

func normalizeFields(fields []Field, depth int, nodes *int) ([]Field, error) {
	if depth > maxSchemaDepth || len(fields) > maxSchemaFields {
		return nil, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
	}
	result := make([]Field, len(fields))
	seen := make(map[string]struct{}, len(fields))
	for index := range fields {
		field := fields[index]
		if !validSchemaFieldName(field.Name) {
			return nil, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
		}
		if _, exists := seen[field.Name]; exists {
			return nil, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
		}
		seen[field.Name] = struct{}{}
		normalized, err := normalizeField(field, depth, nodes, true)
		if err != nil {
			return nil, err
		}
		result[index] = normalized
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func normalizeField(field Field, depth int, nodes *int, named bool) (Field, error) {
	*nodes++
	if *nodes > maxSchemaNodes || depth > maxSchemaDepth || !field.Kind.Valid() {
		return Field{}, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
	}
	if named {
		if !validSchemaFieldName(field.Name) {
			return Field{}, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
		}
	} else if field.Name != "" || field.Required {
		return Field{}, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
	}

	normalized := Field{Name: field.Name, Kind: field.Kind, Required: field.Required}
	switch field.Kind {
	case ValueObject:
		if field.Items != nil {
			return Field{}, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
		}
		children, err := normalizeFields(field.Fields, depth+1, nodes)
		if err != nil {
			return Field{}, err
		}
		normalized.Fields = children
	case ValueArray:
		if len(field.Fields) != 0 || field.Items == nil {
			return Field{}, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
		}
		item, err := normalizeField(*field.Items, depth+1, nodes, false)
		if err != nil {
			return Field{}, err
		}
		normalized.Items = &item
	default:
		if len(field.Fields) != 0 || field.Items != nil {
			return Field{}, schemaFailure(codeSchemaInvalid, failure.CategoryValidation, "event schema is invalid")
		}
	}
	return normalized, nil
}

func validSchemaFieldName(name string) bool {
	if name == "" || strings.TrimSpace(name) != name || !utf8.ValidString(name) || utf8.RuneCountInString(name) > maxSchemaFieldNameRunes {
		return false
	}
	return !strings.ContainsAny(name, "\x00\r\n")
}

func cloneRegisteredSchema(entry RegisteredSchema) RegisteredSchema {
	return RegisteredSchema{
		Identity:    entry.Identity,
		Schema:      cloneSchema(entry.Schema),
		Fingerprint: entry.Fingerprint,
		Canonical:   append([]byte(nil), entry.Canonical...),
	}
}

func cloneSchema(schema Schema) Schema {
	return Schema{Fields: cloneFields(schema.Fields)}
}

func cloneFields(fields []Field) []Field {
	if len(fields) == 0 {
		return nil
	}
	cloned := make([]Field, len(fields))
	for index := range fields {
		cloned[index] = fields[index]
		cloned[index].Fields = cloneFields(fields[index].Fields)
		if fields[index].Items != nil {
			item := cloneFields([]Field{*fields[index].Items})[0]
			cloned[index].Items = &item
		}
	}
	return cloned
}

func equalBytes(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func schemaFailure(code failure.Code, category failure.Category, title string) error {
	value, err := failure.New(code, category, title, failure.WithRetryable(false))
	if err != nil {
		return errors.New("event schema failure could not be classified safely")
	}
	return value
}

package events

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	maxPayloadBytes             = 64 * 1024
	maxPayloadDepth             = 16
	maxPayloadCollectionEntries = 256
	maxPayloadStringRunes       = 4096
)

const (
	codePayloadInvalid failure.Code = "events.schema.payload_invalid"
	codePayloadLimit   failure.Code = "events.schema.payload_limit"
)

type decodedPayloadValue struct {
	kind      ValueKind
	object    map[string]decodedPayloadValue
	array     []decodedPayloadValue
	isInteger bool
}

// ValidateRegisteredPayload resolves the exact immutable schema identity from
// local accepted registry state and validates payload before a protected
// handler may use it. Validation success grants no authorization or mutation.
func ValidateRegisteredPayload(registry *SchemaRegistry, identity PayloadSchemaIdentity, payload []byte) error {
	registered, err := registry.Lookup(identity)
	if err != nil {
		return err
	}
	return ValidatePayload(registered.Schema, payload)
}

// ValidatePayload validates raw JSON bytes against one provider-neutral closed
// object schema under the frozen P04.07 Wave 1 limits. It is synchronous and
// performs no network, filesystem, plugin, shell, callback, or registry I/O.
func ValidatePayload(schema Schema, payload []byte) error {
	normalized, _, _, err := normalizeSchema(schema)
	if err != nil {
		return err
	}
	if len(payload) == 0 || !utf8.Valid(payload) {
		return payloadInvalidFailure()
	}
	if len(payload) > maxPayloadBytes {
		return payloadLimitFailure()
	}

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()
	value, err := decodePayloadValue(decoder, 0)
	if err != nil {
		return err
	}
	if value.kind != ValueObject {
		return payloadInvalidFailure()
	}
	if err := requirePayloadEOF(decoder); err != nil {
		return err
	}
	if !validateDecodedObject(normalized.Fields, value.object) {
		return payloadInvalidFailure()
	}
	return nil
}

func decodePayloadValue(decoder *json.Decoder, containerDepth int) (decodedPayloadValue, error) {
	token, err := decoder.Token()
	if err != nil {
		return decodedPayloadValue{}, payloadInvalidFailure()
	}

	switch value := token.(type) {
	case json.Delim:
		switch value {
		case '{':
			return decodePayloadObject(decoder, containerDepth)
		case '[':
			return decodePayloadArray(decoder, containerDepth)
		default:
			return decodedPayloadValue{}, payloadInvalidFailure()
		}
	case string:
		if utf8.RuneCountInString(value) > maxPayloadStringRunes {
			return decodedPayloadValue{}, payloadLimitFailure()
		}
		return decodedPayloadValue{kind: ValueString}, nil
	case json.Number:
		return decodedPayloadValue{
			kind:      ValueNumber,
			isInteger: isJSONIntegerToken(value.String()),
		}, nil
	case bool:
		return decodedPayloadValue{kind: ValueBoolean}, nil
	case nil:
		return decodedPayloadValue{}, payloadInvalidFailure()
	default:
		return decodedPayloadValue{}, payloadInvalidFailure()
	}
}

func decodePayloadObject(decoder *json.Decoder, containerDepth int) (decodedPayloadValue, error) {
	if containerDepth >= maxPayloadDepth {
		return decodedPayloadValue{}, payloadLimitFailure()
	}
	nextDepth := containerDepth + 1
	values := make(map[string]decodedPayloadValue)
	entries := 0
	for decoder.More() {
		entries++
		if entries > maxPayloadCollectionEntries {
			return decodedPayloadValue{}, payloadLimitFailure()
		}
		keyToken, err := decoder.Token()
		if err != nil {
			return decodedPayloadValue{}, payloadInvalidFailure()
		}
		key, ok := keyToken.(string)
		if !ok {
			return decodedPayloadValue{}, payloadInvalidFailure()
		}
		if utf8.RuneCountInString(key) > maxPayloadStringRunes {
			return decodedPayloadValue{}, payloadLimitFailure()
		}
		if _, exists := values[key]; exists {
			return decodedPayloadValue{}, payloadInvalidFailure()
		}
		decoded, err := decodePayloadValue(decoder, nextDepth)
		if err != nil {
			return decodedPayloadValue{}, err
		}
		values[key] = decoded
	}
	closing, err := decoder.Token()
	if err != nil || closing != json.Delim('}') {
		return decodedPayloadValue{}, payloadInvalidFailure()
	}
	return decodedPayloadValue{kind: ValueObject, object: values}, nil
}

func decodePayloadArray(decoder *json.Decoder, containerDepth int) (decodedPayloadValue, error) {
	if containerDepth >= maxPayloadDepth {
		return decodedPayloadValue{}, payloadLimitFailure()
	}
	nextDepth := containerDepth + 1
	values := make([]decodedPayloadValue, 0)
	for decoder.More() {
		if len(values) >= maxPayloadCollectionEntries {
			return decodedPayloadValue{}, payloadLimitFailure()
		}
		decoded, err := decodePayloadValue(decoder, nextDepth)
		if err != nil {
			return decodedPayloadValue{}, err
		}
		values = append(values, decoded)
	}
	closing, err := decoder.Token()
	if err != nil || closing != json.Delim(']') {
		return decodedPayloadValue{}, payloadInvalidFailure()
	}
	return decodedPayloadValue{kind: ValueArray, array: values}, nil
}

func requirePayloadEOF(decoder *json.Decoder) error {
	if _, err := decoder.Token(); err == io.EOF {
		return nil
	}
	return payloadInvalidFailure()
}

func validateDecodedObject(fields []Field, values map[string]decodedPayloadValue) bool {
	fieldsByName := indexValidationFields(fields)
	for name, value := range values {
		field, ok := fieldsByName[name]
		if !ok || !validateDecodedValue(field, value) {
			return false
		}
	}
	for index := range fields {
		if !fields[index].Required {
			continue
		}
		if _, ok := values[fields[index].Name]; !ok {
			return false
		}
	}
	return true
}

func validateDecodedValue(field Field, value decodedPayloadValue) bool {
	switch field.Kind {
	case ValueString:
		return value.kind == ValueString
	case ValueInteger:
		return value.kind == ValueNumber && value.isInteger
	case ValueNumber:
		return value.kind == ValueNumber
	case ValueBoolean:
		return value.kind == ValueBoolean
	case ValueObject:
		return value.kind == ValueObject && validateDecodedObject(field.Fields, value.object)
	case ValueArray:
		if value.kind != ValueArray || field.Items == nil {
			return false
		}
		for index := range value.array {
			if !validateDecodedValue(*field.Items, value.array[index]) {
				return false
			}
		return true
	default:
		return false
	}
}

func indexValidationFields(fields []Field) map[string]Field {
	indexed := make(map[string]Field, len(fields))
	for index := range fields {
		indexed[fields[index].Name] = fields[index]
	}
	return indexed
}

func isJSONIntegerToken(value string) bool {
	return value != "" && !strings.ContainsAny(value, ".eE")
}

func payloadInvalidFailure() error {
	return schemaFailure(codePayloadInvalid, failure.CategoryValidation, "event payload does not conform to the accepted schema")
}

func payloadLimitFailure() error {
	return schemaFailure(codePayloadLimit, failure.CategoryValidation, "event payload exceeds validation limits")
}

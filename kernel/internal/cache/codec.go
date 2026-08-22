package cache

import (
	"encoding/json"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

// Codec is the typed serialization boundary for cache payloads. Domain types
// remain owned by their callers; P01.05 defines only transport serialization.
type Codec[T any] interface {
	Encode(value T) ([]byte, error)
	Decode(data []byte) (T, error)
}

// JSONCodec provides deterministic type-directed JSON encoding/decoding without
// embedding provider or business-domain semantics in the cache substrate.
type JSONCodec[T any] struct{}

func (JSONCodec[T]) Encode(value T) ([]byte, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, safeWrappedFailure(
			err,
			codeSerializationFailed,
			failure.CategoryValidation,
			"cache value could not be serialized",
		)
	}
	return encoded, nil
}

func (JSONCodec[T]) Decode(data []byte) (T, error) {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return value, safeWrappedFailure(
			err,
			codeSerializationFailed,
			failure.CategoryInvariant,
			"cache value could not be deserialized",
		)
	}
	return value, nil
}

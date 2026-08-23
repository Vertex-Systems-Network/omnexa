// Package configuration provides the P01.10 governed runtime feature-flag and configuration registry.
package configuration

import (
	"sort"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

// Registry is an immutable deterministic definition registry.
type Registry struct {
	definitions map[Key]Definition
	order       []Key
}

// NewRegistry validates definitions, rejects duplicates, and freezes them in key order.
func NewRegistry(definitions ...Definition) (*Registry, error) {
	if len(definitions) == 0 {
		return nil, classifiedFailure(codeRegistryInvalid, failure.CategoryInvariant, "runtime configuration registry is empty")
	}

	items := make(map[Key]Definition, len(definitions))
	order := make([]Key, 0, len(definitions))
	for _, definition := range definitions {
		if err := validateDefinition(definition); err != nil {
			return nil, err
		}
		if _, exists := items[definition.Key]; exists {
			return nil, classifiedFailure(codeDefinitionDuplicate, failure.CategoryConflict, "runtime configuration definition is duplicated")
		}
		items[definition.Key] = definition
		order = append(order, definition.Key)
	}

	sort.Slice(order, func(left, right int) bool {
		return order[left] < order[right]
	})
	return &Registry{definitions: items, order: order}, nil
}

// Definition returns one immutable definition copy.
func (registry *Registry) Definition(key Key) (Definition, bool) {
	if registry == nil {
		return Definition{}, false
	}
	definition, ok := registry.definitions[key]
	return definition, ok
}

// Definitions returns definitions in deterministic key order.
func (registry *Registry) Definitions() []Definition {
	if registry == nil {
		return []Definition{}
	}
	result := make([]Definition, 0, len(registry.order))
	for _, key := range registry.order {
		result = append(result, registry.definitions[key])
	}
	return result
}

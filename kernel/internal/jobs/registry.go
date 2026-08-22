package jobs

import (
	"context"
	"sort"
	"sync"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

// Handler executes one job attempt. Scheduler identity never grants authority;
// handlers remain responsible for their owning capability's authorization and
// duplicate-safe side-effect contract.
type Handler func(context.Context, Invocation) error

// Registry owns the deterministic kernel-local job type boundary.
type Registry struct {
	mu       sync.RWMutex
	handlers map[Type]Handler
	frozen   bool
}

// NewRegistry returns an empty mutable registry. An Executor freezes it before
// worker execution begins.
func NewRegistry() *Registry {
	return &Registry{handlers: make(map[Type]Handler)}
}

// Register adds one unique job type before registry freeze.
func (registry *Registry) Register(jobType Type, handler Handler) error {
	if registry == nil {
		return classifiedFailure(codeRegistryInvalid, failure.CategoryInvariant, "job registry is invalid", false)
	}
	if !jobType.valid() || handler == nil {
		return classifiedFailure(codeTypeInvalid, failure.CategoryValidation, "job registration is invalid", false)
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()
	if registry.frozen {
		return classifiedFailure(codeRegistryFrozen, failure.CategoryConflict, "job registry is frozen", false)
	}
	if _, exists := registry.handlers[jobType]; exists {
		return classifiedFailure(codeTypeDuplicate, failure.CategoryConflict, "job type is already registered", false)
	}
	registry.handlers[jobType] = handler
	return nil
}

// Types returns registered job types in deterministic lexical order.
func (registry *Registry) Types() []Type {
	if registry == nil {
		return nil
	}
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	result := make([]Type, 0, len(registry.handlers))
	for jobType := range registry.handlers {
		result = append(result, jobType)
	}
	sort.Slice(result, func(left, right int) bool {
		return result[left] < result[right]
	})
	return result
}

func (registry *Registry) freeze() error {
	if registry == nil {
		return classifiedFailure(codeRegistryInvalid, failure.CategoryInvariant, "job registry is invalid", false)
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.frozen = true
	return nil
}

func (registry *Registry) lookup(jobType Type) (Handler, error) {
	if registry == nil {
		return nil, classifiedFailure(codeRegistryInvalid, failure.CategoryInvariant, "job registry is invalid", false)
	}
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	handler, exists := registry.handlers[jobType]
	if !exists {
		return nil, classifiedFailure(codeTypeUnknown, failure.CategoryNotFound, "job type is not registered", false)
	}
	return handler, nil
}

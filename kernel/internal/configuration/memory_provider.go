package configuration

import (
	"context"
	"sync"
)

type providerAddress struct {
	key   Key
	scope string
}

// MemoryProvider is a deterministic exact-context provider intended for tests and
// local composition. It is non-authoritative and performs no scope or authorization inference.
type MemoryProvider struct {
	mu       sync.RWMutex
	values   map[providerAddress]ProviderResult
	revision uint64
	failure  error
}

// NewMemoryProvider returns an empty deterministic provider.
func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{values: make(map[providerAddress]ProviderResult)}
}

// Set stores one exact override and returns its monotonically increasing revision.
func (provider *MemoryProvider) Set(key Key, scope EvaluationContext, value Value) (uint64, error) {
	if provider == nil || !keyPattern.MatchString(string(key)) || !value.valid() {
		return 0, safeFailure(codeEvaluatorInvalid, "runtime configuration test provider input is invalid")
	}
	if err := scope.validate(); err != nil {
		return 0, err
	}

	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.revision++
	result := ProviderResult{Value: value, Revision: provider.revision}
	provider.values[providerAddress{key: key, scope: scope.cacheKey()}] = result
	return result.Revision, nil
}

// Delete removes one exact override and advances the provider revision.
func (provider *MemoryProvider) Delete(key Key, scope EvaluationContext) (uint64, error) {
	if provider == nil || !keyPattern.MatchString(string(key)) {
		return 0, safeFailure(codeEvaluatorInvalid, "runtime configuration test provider input is invalid")
	}
	if err := scope.validate(); err != nil {
		return 0, err
	}

	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.revision++
	delete(provider.values, providerAddress{key: key, scope: scope.cacheKey()})
	return provider.revision, nil
}

// SetFailure makes every Resolve return the supplied private provider failure.
// Passing nil clears failure mode. The error text is never copied into Evaluation.
func (provider *MemoryProvider) SetFailure(err error) {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.failure = err
}

// Resolve implements Provider with exact-context matching only.
func (provider *MemoryProvider) Resolve(ctx context.Context, key Key, scope EvaluationContext) (ProviderResult, error) {
	if provider == nil {
		return ProviderResult{}, ErrProviderValueNotFound
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return ProviderResult{}, err
	}

	provider.mu.RLock()
	defer provider.mu.RUnlock()
	if provider.failure != nil {
		return ProviderResult{}, provider.failure
	}
	result, ok := provider.values[providerAddress{key: key, scope: scope.cacheKey()}]
	if !ok {
		return ProviderResult{}, ErrProviderValueNotFound
	}
	return result, nil
}

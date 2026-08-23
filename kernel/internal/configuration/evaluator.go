package configuration

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	defaultCacheTTL               = 30 * time.Second
	defaultRefreshTimeout         = time.Second
	defaultMaxCacheEntries        = 256
	defaultMaxConcurrentRefreshes = 8
	maxCacheTTL                   = 24 * time.Hour
	maxRefreshTimeout             = 30 * time.Second
	maxCacheEntries               = 4096
	maxConcurrentRefreshes        = 64
)

// EvaluatorOptions bounds runtime cache and refresh behavior. Zero values select safe defaults.
type EvaluatorOptions struct {
	CacheTTL               time.Duration
	RefreshTimeout         time.Duration
	MaxCacheEntries        int
	MaxConcurrentRefreshes int
	Now                    func() time.Time
}

type cacheAddress struct {
	key   Key
	scope string
}

type cacheEntry struct {
	value             Value
	source            Source
	reason            Reason
	definitionVersion uint64
	providerRevision  uint64
	evaluatedAt       time.Time
	expiresAt         time.Time
	degraded          bool
	sequence          uint64
}

// Evaluator resolves registered definitions through an optional runtime provider.
// Its cache is bounded and non-authoritative; provider failure uses the definition fallback.
type Evaluator struct {
	registry       *Registry
	provider       Provider
	cacheTTL       time.Duration
	refreshTimeout time.Duration
	maxEntries     int
	now            func() time.Time
	refreshSlots   chan struct{}

	mu       sync.Mutex
	cache    map[cacheAddress]cacheEntry
	sequence uint64
}

// NewEvaluator constructs a bounded runtime evaluator. A nil provider means definitions
// resolve to their defaults until a governed provider is supplied by a later composition root.
func NewEvaluator(registry *Registry, provider Provider, options EvaluatorOptions) (*Evaluator, error) {
	if registry == nil || len(registry.definitions) == 0 {
		return nil, classifiedFailure(codeEvaluatorInvalid, failure.CategoryInvariant, "runtime configuration evaluator is invalid")
	}

	normalized, err := normalizeEvaluatorOptions(options)
	if err != nil {
		return nil, err
	}
	return &Evaluator{
		registry:       registry,
		provider:       provider,
		cacheTTL:       normalized.CacheTTL,
		refreshTimeout: normalized.RefreshTimeout,
		maxEntries:     normalized.MaxCacheEntries,
		now:            normalized.Now,
		refreshSlots:   make(chan struct{}, normalized.MaxConcurrentRefreshes),
		cache:          make(map[cacheAddress]cacheEntry, normalized.MaxCacheEntries),
	}, nil
}

func normalizeEvaluatorOptions(options EvaluatorOptions) (EvaluatorOptions, error) {
	if options.CacheTTL == 0 {
		options.CacheTTL = defaultCacheTTL
	}
	if options.RefreshTimeout == 0 {
		options.RefreshTimeout = defaultRefreshTimeout
	}
	if options.MaxCacheEntries == 0 {
		options.MaxCacheEntries = defaultMaxCacheEntries
	}
	if options.MaxConcurrentRefreshes == 0 {
		options.MaxConcurrentRefreshes = defaultMaxConcurrentRefreshes
	}
	if options.Now == nil {
		options.Now = time.Now
	}

	if options.CacheTTL < time.Millisecond || options.CacheTTL > maxCacheTTL ||
		options.RefreshTimeout < time.Millisecond || options.RefreshTimeout > maxRefreshTimeout ||
		options.MaxCacheEntries < 1 || options.MaxCacheEntries > maxCacheEntries ||
		options.MaxConcurrentRefreshes < 1 || options.MaxConcurrentRefreshes > maxConcurrentRefreshes {
		return EvaluatorOptions{}, classifiedFailure(codeEvaluatorInvalid, failure.CategoryValidation, "runtime configuration evaluator options are invalid")
	}
	return options, nil
}

// Evaluate resolves a value from cache or performs a bounded provider refresh.
func (evaluator *Evaluator) Evaluate(ctx context.Context, key Key, scope EvaluationContext) (Evaluation, error) {
	definition, err := evaluator.definitionAndScope(key, scope)
	if err != nil {
		return Evaluation{}, err
	}
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return Evaluation{}, evaluationContextFailure(err)
	}

	address := cacheAddress{key: key, scope: scope.cacheKey()}
	now := evaluator.now().UTC()
	if cached, ok := evaluator.cached(address, now); ok {
		return evaluationFromEntry(key, cached, true, nil), nil
	}
	return evaluator.refresh(ctx, definition, scope, address, now)
}

// Refresh bypasses cache freshness and resolves the exact definition/context pair again.
func (evaluator *Evaluator) Refresh(ctx context.Context, key Key, scope EvaluationContext) (Evaluation, error) {
	definition, err := evaluator.definitionAndScope(key, scope)
	if err != nil {
		return Evaluation{}, err
	}
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return Evaluation{}, evaluationContextFailure(err)
	}
	return evaluator.refresh(ctx, definition, scope, cacheAddress{key: key, scope: scope.cacheKey()}, evaluator.now().UTC())
}

// Invalidate removes one exact cached context and returns value-free change metadata.
func (evaluator *Evaluator) Invalidate(key Key, scope EvaluationContext) (Change, bool, error) {
	definition, err := evaluator.definitionAndScope(key, scope)
	if err != nil {
		return Change{}, false, err
	}
	address := cacheAddress{key: key, scope: scope.cacheKey()}
	now := evaluator.now().UTC()

	evaluator.mu.Lock()
	defer evaluator.mu.Unlock()
	entry, exists := evaluator.cache[address]
	if !exists {
		return Change{}, false, nil
	}
	delete(evaluator.cache, address)
	return Change{
		Key:               key,
		Action:            ChangeInvalidated,
		Source:            entry.source,
		DefinitionVersion: definition.Version,
		ProviderRevision:  entry.providerRevision,
		At:                now,
	}, true, nil
}

// InvalidateKey removes every bounded cached scope for one definition in deterministic scope order.
func (evaluator *Evaluator) InvalidateKey(key Key) ([]Change, error) {
	definition, ok := evaluator.registry.Definition(key)
	if !ok {
		return nil, classifiedFailure(codeDefinitionUnknown, failure.CategoryNotFound, "runtime configuration definition is unknown")
	}
	now := evaluator.now().UTC()

	evaluator.mu.Lock()
	defer evaluator.mu.Unlock()
	addresses := make([]cacheAddress, 0)
	for address := range evaluator.cache {
		if address.key == key {
			addresses = append(addresses, address)
		}
	}
	sort.Slice(addresses, func(left, right int) bool {
		return addresses[left].scope < addresses[right].scope
	})
	changes := make([]Change, 0, len(addresses))
	for _, address := range addresses {
		entry := evaluator.cache[address]
		delete(evaluator.cache, address)
		changes = append(changes, Change{
			Key:               key,
			Action:            ChangeInvalidated,
			Source:            entry.source,
			DefinitionVersion: definition.Version,
			ProviderRevision:  entry.providerRevision,
			At:                now,
		})
	}
	return changes, nil
}

// CacheSize reports bounded cache occupancy for diagnostics/tests without exposing values or scope IDs.
func (evaluator *Evaluator) CacheSize() int {
	if evaluator == nil {
		return 0
	}
	evaluator.mu.Lock()
	defer evaluator.mu.Unlock()
	return len(evaluator.cache)
}

func (evaluator *Evaluator) definitionAndScope(key Key, scope EvaluationContext) (Definition, error) {
	if evaluator == nil || evaluator.registry == nil {
		return Definition{}, classifiedFailure(codeEvaluatorInvalid, failure.CategoryInvariant, "runtime configuration evaluator is invalid")
	}
	definition, ok := evaluator.registry.Definition(key)
	if !ok {
		return Definition{}, classifiedFailure(codeDefinitionUnknown, failure.CategoryNotFound, "runtime configuration definition is unknown")
	}
	if err := scope.validate(); err != nil {
		return Definition{}, err
	}
	return definition, nil
}

func (evaluator *Evaluator) cached(address cacheAddress, now time.Time) (cacheEntry, bool) {
	evaluator.mu.Lock()
	defer evaluator.mu.Unlock()
	entry, exists := evaluator.cache[address]
	return entry, exists && now.Before(entry.expiresAt)
}

func (evaluator *Evaluator) refresh(ctx context.Context, definition Definition, scope EvaluationContext, address cacheAddress, now time.Time) (Evaluation, error) {
	value := definition.Default
	source := SourceDefault
	reason := ReasonNone
	providerRevision := uint64(0)
	degraded := false

	if evaluator.provider != nil {
		result, providerReason, notFound, err := evaluator.resolveBounded(ctx, definition.Key, scope)
		if err != nil {
			return Evaluation{}, err
		}
		switch {
		case notFound:
			// Default is the deterministic no-override behavior.
		case providerReason != ReasonNone:
			value = definition.Fallback
			source = SourceFallback
			reason = providerReason
			degraded = true
		case result.Revision == 0 || !result.Value.valid() || result.Value.Kind() != definition.Kind:
			value = definition.Fallback
			source = SourceFallback
			reason = ReasonInvalidResult
			degraded = true
		default:
			value = result.Value
			source = SourceProvider
			providerRevision = result.Revision
		}
	}

	entry := cacheEntry{
		value:             value,
		source:            source,
		reason:            reason,
		definitionVersion: definition.Version,
		providerRevision:  providerRevision,
		evaluatedAt:       now,
		expiresAt:         now.Add(evaluator.cacheTTL),
		degraded:          degraded,
	}
	change := evaluator.store(address, entry)
	return evaluationFromEntry(definition.Key, entry, false, change), nil
}

type providerOutcome struct {
	result   ProviderResult
	err      error
	panicked bool
}

func (evaluator *Evaluator) resolveBounded(parent context.Context, key Key, scope EvaluationContext) (ProviderResult, Reason, bool, error) {
	bounded, cancel := context.WithTimeout(parent, evaluator.refreshTimeout)
	defer cancel()

	select {
	case evaluator.refreshSlots <- struct{}{}:
	case <-bounded.Done():
		if err := parent.Err(); err != nil {
			return ProviderResult{}, ReasonNone, false, evaluationContextFailure(err)
		}
		return ProviderResult{}, ReasonTimeout, false, nil
	}

	completed := make(chan providerOutcome, 1)
	go func() {
		defer func() { <-evaluator.refreshSlots }()
		completed <- runProviderSafely(bounded, evaluator.provider, key, scope)
	}()

	select {
	case outcome := <-completed:
		if err := parent.Err(); err != nil {
			return ProviderResult{}, ReasonNone, false, evaluationContextFailure(err)
		}
		if outcome.panicked {
			return ProviderResult{}, ReasonPanic, false, nil
		}
		if outcome.err == nil {
			return outcome.result, ReasonNone, false, nil
		}
		if errors.Is(outcome.err, ErrProviderValueNotFound) {
			return ProviderResult{}, ReasonNone, true, nil
		}
		if errors.Is(outcome.err, context.DeadlineExceeded) || errors.Is(bounded.Err(), context.DeadlineExceeded) {
			return ProviderResult{}, ReasonTimeout, false, nil
		}
		return ProviderResult{}, ReasonUnavailable, false, nil
	case <-bounded.Done():
		if err := parent.Err(); err != nil {
			return ProviderResult{}, ReasonNone, false, evaluationContextFailure(err)
		}
		return ProviderResult{}, ReasonTimeout, false, nil
	}
}

func runProviderSafely(ctx context.Context, provider Provider, key Key, scope EvaluationContext) (outcome providerOutcome) {
	defer func() {
		if recover() != nil {
			outcome = providerOutcome{panicked: true}
		}
	}()
	outcome.result, outcome.err = provider.Resolve(ctx, key, scope)
	return outcome
}

func (evaluator *Evaluator) store(address cacheAddress, entry cacheEntry) *Change {
	evaluator.mu.Lock()
	defer evaluator.mu.Unlock()

	previous, existed := evaluator.cache[address]
	changed := !existed || !cacheEntriesEquivalent(previous, entry)
	if !existed && len(evaluator.cache) >= evaluator.maxEntries {
		evaluator.evictOldestLocked()
	}
	evaluator.sequence++
	entry.sequence = evaluator.sequence
	evaluator.cache[address] = entry
	if !changed {
		return nil
	}
	return &Change{
		Key:               address.key,
		Action:            ChangeResolved,
		Source:            entry.source,
		DefinitionVersion: entry.definitionVersion,
		ProviderRevision:  entry.providerRevision,
		At:                entry.evaluatedAt,
	}
}

func (evaluator *Evaluator) evictOldestLocked() {
	var selected cacheAddress
	var selectedEntry cacheEntry
	found := false
	for address, entry := range evaluator.cache {
		if !found || entry.sequence < selectedEntry.sequence {
			selected = address
			selectedEntry = entry
			found = true
		}
	}
	if found {
		delete(evaluator.cache, selected)
	}
}

func cacheEntriesEquivalent(left, right cacheEntry) bool {
	return left.value.equal(right.value) &&
		left.source == right.source &&
		left.reason == right.reason &&
		left.definitionVersion == right.definitionVersion &&
		left.providerRevision == right.providerRevision &&
		left.degraded == right.degraded
}

func evaluationFromEntry(key Key, entry cacheEntry, cacheHit bool, change *Change) Evaluation {
	return Evaluation{
		Key:               key,
		Value:             entry.value,
		Source:            entry.source,
		Reason:            entry.reason,
		DefinitionVersion: entry.definitionVersion,
		ProviderRevision:  entry.providerRevision,
		EvaluatedAt:       entry.evaluatedAt,
		ExpiresAt:         entry.expiresAt,
		CacheHit:          cacheHit,
		Degraded:          entry.degraded,
		Change:            change,
	}
}

func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

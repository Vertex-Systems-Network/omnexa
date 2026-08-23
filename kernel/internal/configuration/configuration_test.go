package configuration

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	testTenantID       = "018f47a6-7b7e-7c14-a847-0af56b2a44fe"
	testOrganizationID = "018f47b1-f3a8-79da-8d9d-4a64d65c2fd5"
	testUserID         = "018f47c2-1111-7abc-8def-1234567890ab"
)

func TestRegistryIsTypedDeterministicAndRejectsUnsafeDefinitions(t *testing.T) {
	t.Parallel()

	kill := killSwitchDefinition()
	feature := featureFlagDefinition()
	registry, err := NewRegistry(kill, feature)
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}
	definitions := registry.Definitions()
	if len(definitions) != 2 || definitions[0].Key != feature.Key || definitions[1].Key != kill.Key {
		t.Fatalf("Definitions() = %#v, want deterministic key order", definitions)
	}

	if _, err := NewRegistry(feature, feature); !failure.IsCode(err, codeDefinitionDuplicate) {
		t.Fatalf("duplicate error = %v, want %s", err, codeDefinitionDuplicate)
	}

	invalidKey := feature
	invalidKey.Key = "not_snake_dot_path"
	if _, err := NewRegistry(invalidKey); !failure.IsCode(err, codeDefinitionInvalid) {
		t.Fatalf("invalid key error = %v, want %s", err, codeDefinitionInvalid)
	}

	invalidOwner := feature
	invalidOwner.Owner = "Kernel Configuration"
	if _, err := NewRegistry(invalidOwner); !failure.IsCode(err, codeDefinitionInvalid) {
		t.Fatalf("invalid owner error = %v, want %s", err, codeDefinitionInvalid)
	}

	nonBooleanFeature := feature
	nonBooleanFeature.Kind = KindString
	nonBooleanFeature.Default = StringValue("off")
	nonBooleanFeature.Fallback = StringValue("off")
	if _, err := NewRegistry(nonBooleanFeature); !failure.IsCode(err, codeDefinitionInvalid) {
		t.Fatalf("non-boolean feature error = %v, want %s", err, codeDefinitionInvalid)
	}

	unsafeKill := kill
	unsafeKill.Fallback = BoolValue(false)
	if _, err := NewRegistry(unsafeKill); !failure.IsCode(err, codeDefinitionInvalid) {
		t.Fatalf("unsafe kill-switch fallback error = %v, want %s", err, codeDefinitionInvalid)
	}
}

func TestEvaluationUsesProviderDefaultAndFallbackDeterministically(t *testing.T) {
	t.Parallel()

	definition := featureFlagDefinition()
	registry := mustRegistry(t, definition)
	provider := NewMemoryProvider()
	evaluator := mustEvaluator(t, registry, provider, EvaluatorOptions{})
	scope := EvaluationContext{TenantID: testTenantID}

	initial, err := evaluator.Evaluate(context.Background(), definition.Key, scope)
	if err != nil {
		t.Fatalf("Evaluate(default) error = %v", err)
	}
	assertBoolEvaluation(t, initial, false, SourceDefault, false)

	revision, err := provider.Set(definition.Key, scope, BoolValue(true))
	if err != nil {
		t.Fatalf("provider.Set() error = %v", err)
	}
	provided, err := evaluator.Refresh(context.Background(), definition.Key, scope)
	if err != nil {
		t.Fatalf("Refresh(provider) error = %v", err)
	}
	assertBoolEvaluation(t, provided, true, SourceProvider, false)
	if provided.ProviderRevision != revision || provided.Change == nil || provided.Change.ProviderRevision != revision {
		t.Fatalf("provider revision/change = %#v, want revision %d", provided, revision)
	}

	provider.SetFailure(errors.New("private-provider-token=do-not-expose"))
	fallback, err := evaluator.Refresh(context.Background(), definition.Key, scope)
	if err != nil {
		t.Fatalf("Refresh(fallback) error = %v", err)
	}
	assertBoolEvaluation(t, fallback, false, SourceFallback, true)
	if fallback.Reason != ReasonUnavailable || fallback.ProviderRevision != 0 {
		t.Fatalf("fallback metadata = %#v", fallback)
	}
}

func TestKillSwitchProviderFailureFailsClosed(t *testing.T) {
	t.Parallel()

	definition := killSwitchDefinition()
	registry := mustRegistry(t, definition)
	provider := NewMemoryProvider()
	provider.SetFailure(errors.New("provider outage"))
	evaluator := mustEvaluator(t, registry, provider, EvaluatorOptions{})

	evaluation, err := evaluator.Evaluate(context.Background(), definition.Key, EvaluationContext{})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	assertBoolEvaluation(t, evaluation, true, SourceFallback, true)
	if evaluation.Reason != ReasonUnavailable {
		t.Fatalf("Reason = %q, want %q", evaluation.Reason, ReasonUnavailable)
	}
}

func TestEvaluationContextAcceptsOnlyCanonicalUUIDv7Metadata(t *testing.T) {
	t.Parallel()

	definition := featureFlagDefinition()
	evaluator := mustEvaluator(t, mustRegistry(t, definition), nil, EvaluatorOptions{})
	valid := EvaluationContext{TenantID: testTenantID, OrganizationID: testOrganizationID, UserID: testUserID}
	if _, err := evaluator.Evaluate(context.Background(), definition.Key, valid); err != nil {
		t.Fatalf("valid scope error = %v", err)
	}

	invalid := valid
	invalid.UserID = "550e8400-e29b-41d4-a716-446655440000"
	if _, err := evaluator.Evaluate(context.Background(), definition.Key, invalid); !failure.IsCode(err, codeContextInvalid) {
		t.Fatalf("invalid UUID error = %v, want %s", err, codeContextInvalid)
	}

	invalid.UserID = "018F47C2-1111-7ABC-8DEF-1234567890AB"
	if _, err := evaluator.Evaluate(context.Background(), definition.Key, invalid); !failure.IsCode(err, codeContextInvalid) {
		t.Fatalf("non-canonical UUID error = %v, want %s", err, codeContextInvalid)
	}
}

func TestCacheRefreshInvalidationAndChangeMetadataAreBounded(t *testing.T) {
	t.Parallel()

	definition := featureFlagDefinition()
	registry := mustRegistry(t, definition)
	provider := NewMemoryProvider()
	now := time.Date(2026, time.August, 23, 0, 0, 0, 0, time.UTC)
	evaluator := mustEvaluator(t, registry, provider, EvaluatorOptions{
		CacheTTL:        time.Minute,
		MaxCacheEntries: 2,
		Now:             func() time.Time { return now },
	})
	scope := EvaluationContext{TenantID: testTenantID}

	revisionOne, err := provider.Set(definition.Key, scope, BoolValue(false))
	if err != nil {
		t.Fatalf("provider.Set(false) error = %v", err)
	}
	first, err := evaluator.Evaluate(context.Background(), definition.Key, scope)
	if err != nil {
		t.Fatalf("Evaluate(first) error = %v", err)
	}
	if first.Change == nil || first.Change.Action != ChangeResolved || first.ProviderRevision != revisionOne {
		t.Fatalf("first change = %#v", first)
	}

	revisionTwo, err := provider.Set(definition.Key, scope, BoolValue(true))
	if err != nil {
		t.Fatalf("provider.Set(true) error = %v", err)
	}
	cached, err := evaluator.Evaluate(context.Background(), definition.Key, scope)
	if err != nil {
		t.Fatalf("Evaluate(cached) error = %v", err)
	}
	assertBoolEvaluation(t, cached, false, SourceProvider, false)
	if !cached.CacheHit || cached.ProviderRevision != revisionOne || cached.Change != nil {
		t.Fatalf("cached metadata = %#v", cached)
	}

	now = now.Add(time.Second)
	refreshed, err := evaluator.Refresh(context.Background(), definition.Key, scope)
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	assertBoolEvaluation(t, refreshed, true, SourceProvider, false)
	if refreshed.ProviderRevision != revisionTwo || refreshed.Change == nil || refreshed.Change.ProviderRevision != revisionTwo {
		t.Fatalf("refreshed metadata = %#v", refreshed)
	}

	change, err := evaluator.Invalidate(definition.Key, scope)
	if err != nil {
		t.Fatalf("Invalidate() error = %v", err)
	}
	if change == nil || change.Action != ChangeInvalidated || change.ProviderRevision != revisionTwo {
		t.Fatalf("invalidation change = %#v", change)
	}
	if evaluator.CacheSize() != 0 {
		t.Fatalf("CacheSize() = %d, want 0", evaluator.CacheSize())
	}
}

func TestCacheCapacityRemainsBoundedAcrossScopedInputs(t *testing.T) {
	t.Parallel()

	definition := featureFlagDefinition()
	evaluator := mustEvaluator(t, mustRegistry(t, definition), nil, EvaluatorOptions{MaxCacheEntries: 2})
	scopes := []EvaluationContext{
		{},
		{TenantID: testTenantID},
		{TenantID: testOrganizationID},
	}
	for _, scope := range scopes {
		if _, err := evaluator.Evaluate(context.Background(), definition.Key, scope); err != nil {
			t.Fatalf("Evaluate(%+v) error = %v", scope, err)
		}
	}
	if evaluator.CacheSize() != 2 {
		t.Fatalf("CacheSize() = %d, want bounded size 2", evaluator.CacheSize())
	}
}

func TestRefreshTimeoutBoundsProviderThatIgnoresContext(t *testing.T) {
	t.Parallel()

	definition := killSwitchDefinition()
	provider := newBlockingProvider(BoolValue(false))
	evaluator := mustEvaluator(t, mustRegistry(t, definition), provider, EvaluatorOptions{
		RefreshTimeout:         20 * time.Millisecond,
		MaxConcurrentRefreshes: 1,
	})

	first, err := evaluator.Evaluate(context.Background(), definition.Key, EvaluationContext{})
	if err != nil {
		t.Fatalf("first Evaluate() error = %v", err)
	}
	assertBoolEvaluation(t, first, true, SourceFallback, true)
	if first.Reason != ReasonTimeout {
		t.Fatalf("first reason = %q, want %q", first.Reason, ReasonTimeout)
	}

	second, err := evaluator.Refresh(context.Background(), definition.Key, EvaluationContext{})
	if err != nil {
		t.Fatalf("second Refresh() error = %v", err)
	}
	assertBoolEvaluation(t, second, true, SourceFallback, true)
	if second.Reason != ReasonTimeout || provider.calls.Load() != 1 {
		t.Fatalf("second reason/calls = %q/%d, want timeout/1", second.Reason, provider.calls.Load())
	}

	close(provider.release)
	select {
	case <-provider.done:
	case <-time.After(time.Second):
		t.Fatal("blocked provider did not exit after release")
	}
}

func TestCallerCancellationIsPropagatedInsteadOfConvertedToFallback(t *testing.T) {
	t.Parallel()

	definition := killSwitchDefinition()
	evaluator := mustEvaluator(t, mustRegistry(t, definition), NewMemoryProvider(), EvaluatorOptions{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := evaluator.Evaluate(ctx, definition.Key, EvaluationContext{})
	if !failure.IsCode(err, codeEvaluationCanceled) || !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation error = %v, want structured cancellation wrapping context.Canceled", err)
	}
}

func TestProviderPanicAndInvalidResultUseSafeFallback(t *testing.T) {
	t.Parallel()

	definition := killSwitchDefinition()
	registry := mustRegistry(t, definition)

	panicEvaluation, err := mustEvaluator(t, registry, panicProvider{}, EvaluatorOptions{}).Evaluate(context.Background(), definition.Key, EvaluationContext{})
	if err != nil {
		t.Fatalf("panic provider Evaluate() error = %v", err)
	}
	assertBoolEvaluation(t, panicEvaluation, true, SourceFallback, true)
	if panicEvaluation.Reason != ReasonPanic {
		t.Fatalf("panic reason = %q, want %q", panicEvaluation.Reason, ReasonPanic)
	}

	invalidEvaluation, err := mustEvaluator(t, registry, invalidProvider{}, EvaluatorOptions{}).Evaluate(context.Background(), definition.Key, EvaluationContext{})
	if err != nil {
		t.Fatalf("invalid provider Evaluate() error = %v", err)
	}
	assertBoolEvaluation(t, invalidEvaluation, true, SourceFallback, true)
	if invalidEvaluation.Reason != ReasonInvalidResult {
		t.Fatalf("invalid reason = %q, want %q", invalidEvaluation.Reason, ReasonInvalidResult)
	}
}

func TestConcurrentEvaluationIsRaceSafe(t *testing.T) {
	definition := featureFlagDefinition()
	provider := NewMemoryProvider()
	if _, err := provider.Set(definition.Key, EvaluationContext{}, BoolValue(true)); err != nil {
		t.Fatalf("provider.Set() error = %v", err)
	}
	evaluator := mustEvaluator(t, mustRegistry(t, definition), provider, EvaluatorOptions{MaxCacheEntries: 4})

	var wait sync.WaitGroup
	for range 32 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for range 20 {
				evaluation, err := evaluator.Evaluate(context.Background(), definition.Key, EvaluationContext{})
				if err != nil {
					t.Errorf("Evaluate() error = %v", err)
					return
				}
				value, ok := evaluation.Value.Bool()
				if !ok || !value {
					t.Errorf("Evaluate() bool = %v/%v, want true/true", value, ok)
					return
				}
			}
		}()
	}
	wait.Wait()
	if evaluator.CacheSize() > 4 {
		t.Fatalf("CacheSize() = %d, want <= 4", evaluator.CacheSize())
	}
}

func featureFlagDefinition() Definition {
	return Definition{
		Key:         "platform.feature.preview",
		Description: "Synthetic platform feature gate used by the P01.10 contract tests.",
		Owner:       "kernel.configuration",
		Kind:        KindBool,
		Class:       ClassFeatureFlag,
		Version:     1,
		Default:     BoolValue(false),
		Fallback:    BoolValue(false),
	}
}

func killSwitchDefinition() Definition {
	return Definition{
		Key:         "platform.runtime.emergency_stop",
		Description: "Synthetic disable-only operational control used by the P01.10 contract tests.",
		Owner:       "kernel.configuration",
		Kind:        KindBool,
		Class:       ClassKillSwitch,
		Version:     1,
		Default:     BoolValue(false),
		Fallback:    BoolValue(true),
	}
}

func mustRegistry(t *testing.T, definitions ...Definition) *Registry {
	t.Helper()
	registry, err := NewRegistry(definitions...)
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}
	return registry
}

func mustEvaluator(t *testing.T, registry *Registry, provider Provider, options EvaluatorOptions) *Evaluator {
	t.Helper()
	evaluator, err := NewEvaluator(registry, provider, options)
	if err != nil {
		t.Fatalf("NewEvaluator() error = %v", err)
	}
	return evaluator
}

func assertBoolEvaluation(t *testing.T, evaluation Evaluation, expected bool, source Source, degraded bool) {
	t.Helper()
	value, ok := evaluation.Value.Bool()
	if !ok || value != expected || evaluation.Source != source || evaluation.Degraded != degraded {
		t.Fatalf("evaluation = %#v, bool=%v/%v; want bool=%v source=%s degraded=%v", evaluation, value, ok, expected, source, degraded)
	}
}

type panicProvider struct{}

func (panicProvider) Resolve(context.Context, Key, EvaluationContext) (ProviderResult, error) {
	panic("private provider panic")
}

type invalidProvider struct{}

func (invalidProvider) Resolve(context.Context, Key, EvaluationContext) (ProviderResult, error) {
	return ProviderResult{Value: StringValue("wrong-kind"), Revision: 1}, nil
}

type blockingProvider struct {
	value   Value
	release chan struct{}
	done    chan struct{}
	calls   atomic.Int64
	once    sync.Once
}

func newBlockingProvider(value Value) *blockingProvider {
	return &blockingProvider{value: value, release: make(chan struct{}), done: make(chan struct{})}
}

func (provider *blockingProvider) Resolve(context.Context, Key, EvaluationContext) (ProviderResult, error) {
	provider.calls.Add(1)
	<-provider.release
	provider.once.Do(func() { close(provider.done) })
	return ProviderResult{Value: provider.value, Revision: 1}, nil
}

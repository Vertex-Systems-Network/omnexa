package modules

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/configuration"
)

const (
	p0305SettingKey = configuration.Key("omnexa.inventory.reorder_threshold")
	p0305FlagKey    = configuration.Key("omnexa.inventory.batch_reservation")
)

func TestConfigurationBindingUsesValidatedDeclarationsRegistryAndScopePolicies(t *testing.T) {
	t.Parallel()

	registry := p0305RegistryV1(t, []string{string(p0305SettingKey)}, []string{string(p0305FlagKey)})
	binding, err := BindConfigurationRegistrations(registry, NewMemoryLifecycleStore(), []ModuleConfigurationRegistration{
		p0305GlobalRegistration(p0305FlagDefinition()),
		p0305ScopedRegistration(p0305SettingDefinition(), p0305SettingPolicy()),
	})
	if err != nil {
		t.Fatalf("BindConfigurationRegistrations() error = %v", err)
	}

	configRegistry, ok := binding.Registry()
	if !ok {
		t.Fatal("Registry() missing governed configuration registry")
	}
	definitions := configRegistry.Definitions()
	if len(definitions) != 2 || definitions[0].Key != p0305FlagKey || definitions[1].Key != p0305SettingKey {
		t.Fatalf("Definitions() = %#v, want deterministic key order", definitions)
	}
	policies := binding.ScopedPolicies()
	if len(policies) != 1 || policies[0].Key != p0305SettingKey || !policies[0].AllowOrganizationOverride {
		t.Fatalf("ScopedPolicies() = %#v, want one validated setting policy", policies)
	}
}

func TestConfigurationBindingRetainsV1AndV2DeclarationsInValidatedSnapshot(t *testing.T) {
	t.Parallel()

	v1 := p0305RegistryV1(t, []string{string(p0305SettingKey)}, []string{string(p0305FlagKey)})
	v1Snapshot, ok := v1.manifestSnapshot("omnexa.inventory")
	if !ok || len(v1Snapshot.Settings) != 1 || v1Snapshot.Settings[0] != string(p0305SettingKey) || len(v1Snapshot.FeatureFlags) != 1 || v1Snapshot.FeatureFlags[0] != string(p0305FlagKey) {
		t.Fatalf("unexpected v1 configuration snapshot: %#v", v1Snapshot)
	}

	manifest := validManifestV2()
	manifest.Settings = []string{string(p0305SettingKey)}
	manifest.FeatureFlags = []string{string(p0305FlagKey)}
	registry, err := Discover([]DiscoverySource{{
		ID:        "source.inventory",
		Version:   "1.0.0",
		Manifests: [][]byte{manifestV2Bytes(t, manifest)},
	}})
	if err != nil {
		t.Fatalf("Discover(v2) error = %v", err)
	}
	v2Snapshot, ok := registry.manifestSnapshot("omnexa.inventory")
	if !ok || len(v2Snapshot.Settings) != 1 || v2Snapshot.Settings[0] != string(p0305SettingKey) || len(v2Snapshot.FeatureFlags) != 1 || v2Snapshot.FeatureFlags[0] != string(p0305FlagKey) {
		t.Fatalf("unexpected v2 configuration snapshot: %#v", v2Snapshot)
	}
}

func TestConfigurationBindingRejectsOwnerClassMissingAndUndeclaredDefinitions(t *testing.T) {
	t.Parallel()

	registry := p0305RegistryV1(t, []string{string(p0305SettingKey)}, []string{string(p0305FlagKey)})
	store := NewMemoryLifecycleStore()

	wrongOwner := p0305SettingDefinition()
	wrongOwner.Owner = "omnexa.catalog"
	assertP0305BindingCode(t, p0305BindError(registry, store, []ModuleConfigurationRegistration{
		p0305ScopedRegistration(wrongOwner, p0305SettingPolicy()),
		p0305GlobalRegistration(p0305FlagDefinition()),
	}), "module.configuration.owner_mismatch")

	wrongClass := p0305FlagDefinition()
	wrongClass.Class = configuration.ClassKillSwitch
	assertP0305BindingCode(t, p0305BindError(registry, store, []ModuleConfigurationRegistration{
		p0305ScopedRegistration(p0305SettingDefinition(), p0305SettingPolicy()),
		p0305GlobalRegistration(wrongClass),
	}), "module.configuration.class_mismatch")

	assertP0305BindingCode(t, p0305BindError(registry, store, []ModuleConfigurationRegistration{
		p0305ScopedRegistration(p0305SettingDefinition(), p0305SettingPolicy()),
	}), "module.configuration.definition_missing")

	extra := p0305SettingDefinition()
	extra.Key = "omnexa.inventory.undeclared_setting"
	assertP0305BindingCode(t, p0305BindError(registry, store, []ModuleConfigurationRegistration{
		p0305ScopedRegistration(p0305SettingDefinition(), p0305SettingPolicy()),
		p0305GlobalRegistration(p0305FlagDefinition()),
		p0305GlobalRegistration(extra),
	}), "module.configuration.definition_undeclared")
}

func TestConfigurationBindingRejectsInvalidScopeContracts(t *testing.T) {
	t.Parallel()

	registry := p0305RegistryV1(t, []string{string(p0305SettingKey)}, []string{string(p0305FlagKey)})
	store := NewMemoryLifecycleStore()

	assertP0305BindingCode(t, p0305BindError(registry, store, []ModuleConfigurationRegistration{
		{Definition: p0305SettingDefinition(), Scope: ModuleConfigurationScope("unknown")},
		p0305GlobalRegistration(p0305FlagDefinition()),
	}), "module.configuration.scope_invalid")

	assertP0305BindingCode(t, p0305BindError(registry, store, []ModuleConfigurationRegistration{
		{Definition: p0305SettingDefinition(), Scope: ModuleConfigurationScoped},
		p0305GlobalRegistration(p0305FlagDefinition()),
	}), "module.configuration.scope_policy_missing")

	policy := p0305SettingPolicy()
	assertP0305BindingCode(t, p0305BindError(registry, store, []ModuleConfigurationRegistration{
		{Definition: p0305SettingDefinition(), Scope: ModuleConfigurationGlobal, Policy: &policy},
		p0305GlobalRegistration(p0305FlagDefinition()),
	}), "module.configuration.global_policy_forbidden")

	keyMismatch := p0305SettingPolicy()
	keyMismatch.Key = p0305FlagKey
	assertP0305BindingCode(t, p0305BindError(registry, store, []ModuleConfigurationRegistration{
		p0305ScopedRegistration(p0305SettingDefinition(), keyMismatch),
		p0305GlobalRegistration(p0305FlagDefinition()),
	}), "module.configuration.scope_policy_key_mismatch")

	invalidPolicy := p0305SettingPolicy()
	invalidPolicy.Classification = configuration.DataInternal
	invalidPolicy.ProtectedRead = false
	assertP0305BindingCode(t, p0305BindError(registry, store, []ModuleConfigurationRegistration{
		p0305ScopedRegistration(p0305SettingDefinition(), invalidPolicy),
		p0305GlobalRegistration(p0305FlagDefinition()),
	}), "module.configuration.scope_policy_invalid")
}

func TestConfigurationBindingRejectsCrossClassDeclarationCollision(t *testing.T) {
	t.Parallel()

	registry := p0305RegistryV1(t, []string{string(p0305SettingKey)}, []string{string(p0305SettingKey)})
	assertP0305BindingCode(t, p0305BindError(registry, NewMemoryLifecycleStore(), nil), "module.configuration.declaration_collision")
}

func TestConfigurationBindingLifecycleReadsAreNonDestructiveAndEnabledOnlyIsRuntimeActive(t *testing.T) {
	registry := p0305RegistryV1(t, []string{string(p0305SettingKey)}, []string{string(p0305FlagKey)})
	store := NewMemoryLifecycleStore()
	binding, err := BindConfigurationRegistrations(registry, store, []ModuleConfigurationRegistration{
		p0305ScopedRegistration(p0305SettingDefinition(), p0305SettingPolicy()),
		p0305GlobalRegistration(p0305FlagDefinition()),
	})
	if err != nil {
		t.Fatalf("BindConfigurationRegistrations() error = %v", err)
	}

	assertP0305ResolveCode(t, binding, p0305SettingKey, "module.configuration.unavailable")

	ctx := context.Background()
	if err := store.CompareAndSwap(ctx, "omnexa.inventory", 0, LifecycleRecord{
		ModuleID: "omnexa.inventory", Version: "1.2.3", State: LifecycleInstalled, Revision: 1,
	}); err != nil {
		t.Fatalf("seed installed state: %v", err)
	}
	installed, err := binding.Resolve(ctx, p0305SettingKey)
	if err != nil || installed.RuntimeActive || installed.LifecycleState != LifecycleInstalled || installed.Scope != ModuleConfigurationScoped {
		t.Fatalf("installed Resolve() = %#v, err=%v", installed, err)
	}

	if err := store.CompareAndSwap(ctx, "omnexa.inventory", 1, LifecycleRecord{
		ModuleID: "omnexa.inventory", Version: "1.2.3", State: LifecycleEnabled, Revision: 2,
	}); err != nil {
		t.Fatalf("seed enabled state: %v", err)
	}
	enabled, err := binding.Resolve(ctx, p0305FlagKey)
	if err != nil || !enabled.RuntimeActive || enabled.LifecycleState != LifecycleEnabled || enabled.Scope != ModuleConfigurationGlobal {
		t.Fatalf("enabled Resolve() = %#v, err=%v", enabled, err)
	}

	if err := store.CompareAndSwap(ctx, "omnexa.inventory", 2, LifecycleRecord{
		ModuleID: "omnexa.inventory", Version: "1.2.3", State: LifecycleDisabled, Revision: 3,
	}); err != nil {
		t.Fatalf("seed disabled state: %v", err)
	}
	disabled, err := binding.Resolve(ctx, p0305SettingKey)
	if err != nil || disabled.RuntimeActive || disabled.LifecycleState != LifecycleDisabled {
		t.Fatalf("disabled Resolve() = %#v, err=%v", disabled, err)
	}
	if disabled.Definition.Key != p0305SettingKey {
		t.Fatalf("disabled definition was not retained: %#v", disabled.Definition)
	}

	if err := store.CompareAndSwap(ctx, "omnexa.inventory", 3, LifecycleRecord{
		ModuleID: "omnexa.inventory", Version: "1.2.3", State: LifecycleDetached, Revision: 4,
	}); err != nil {
		t.Fatalf("seed detached state: %v", err)
	}
	detached, err := binding.Resolve(ctx, p0305SettingKey)
	if err != nil || detached.RuntimeActive || detached.LifecycleState != LifecycleDetached {
		t.Fatalf("detached Resolve() = %#v, err=%v", detached, err)
	}

	if err := store.CompareAndSwap(ctx, "omnexa.inventory", 4, LifecycleRecord{
		ModuleID: "omnexa.inventory", State: LifecyclePurged, Revision: 5,
	}); err != nil {
		t.Fatalf("seed purged state: %v", err)
	}
	assertP0305ResolveCode(t, binding, p0305SettingKey, "module.configuration.unavailable")
}

func TestConfigurationBindingDoesNotRequireASecondRegistryForModulesWithoutDeclarations(t *testing.T) {
	t.Parallel()

	registry := p0305RegistryV1(t, []string{}, []string{})
	binding, err := BindConfigurationRegistrations(registry, NewMemoryLifecycleStore(), nil)
	if err != nil {
		t.Fatalf("BindConfigurationRegistrations(empty) error = %v", err)
	}
	if _, ok := binding.Registry(); ok {
		t.Fatal("empty module declaration set unexpectedly created a configuration registry")
	}
	if policies := binding.ScopedPolicies(); len(policies) != 0 {
		t.Fatalf("ScopedPolicies(empty) = %#v, want empty", policies)
	}
}

func p0305RegistryV1(t *testing.T, settings, flags []string) Registry {
	t.Helper()
	manifest := validManifest()
	manifest.Settings = append([]string(nil), settings...)
	manifest.FeatureFlags = append([]string(nil), flags...)
	payload, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	registry, err := Discover([]DiscoverySource{{
		ID:        "source.inventory",
		Version:   "1.0.0",
		Manifests: [][]byte{payload},
	}})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	return registry
}

func p0305GlobalRegistration(definition configuration.Definition) ModuleConfigurationRegistration {
	return ModuleConfigurationRegistration{Definition: definition, Scope: ModuleConfigurationGlobal}
}

func p0305ScopedRegistration(definition configuration.Definition, policy configuration.SettingPolicy) ModuleConfigurationRegistration {
	return ModuleConfigurationRegistration{Definition: definition, Scope: ModuleConfigurationScoped, Policy: &policy}
}

func p0305BindError(registry Registry, store LifecycleStore, registrations []ModuleConfigurationRegistration) error {
	_, err := BindConfigurationRegistrations(registry, store, registrations)
	return err
}

func p0305SettingDefinition() configuration.Definition {
	return configuration.Definition{
		Key:         p0305SettingKey,
		Description: "Inventory reorder threshold used by P03.05 module configuration tests.",
		Owner:       "omnexa.inventory",
		Kind:        configuration.KindInt,
		Class:       configuration.ClassRuntimeConfig,
		Version:     1,
		Default:     configuration.IntValue(10),
		Fallback:    configuration.IntValue(10),
	}
}

func p0305FlagDefinition() configuration.Definition {
	return configuration.Definition{
		Key:         p0305FlagKey,
		Description: "Inventory batch reservation feature flag used by P03.05 tests.",
		Owner:       "omnexa.inventory",
		Kind:        configuration.KindBool,
		Class:       configuration.ClassFeatureFlag,
		Version:     1,
		Default:     configuration.BoolValue(false),
		Fallback:    configuration.BoolValue(false),
	}
}

func p0305SettingPolicy() configuration.SettingPolicy {
	return configuration.SettingPolicy{
		Key:                       p0305SettingKey,
		Classification:            configuration.DataInternal,
		AllowOrganizationOverride: true,
		ProtectedRead:             true,
		SecuritySignificant:       false,
	}
}

func assertP0305BindingCode(t *testing.T, err error, code string) {
	t.Helper()
	var bindingErr *ModuleConfigurationError
	if !errors.As(err, &bindingErr) || bindingErr.Diagnostic.Code != code {
		t.Fatalf("binding error = %#v, want code %s", err, code)
	}
}

func assertP0305ResolveCode(t *testing.T, binding *ConfigurationBinding, key configuration.Key, code string) {
	t.Helper()
	_, err := binding.Resolve(context.Background(), key)
	assertP0305BindingCode(t, err, code)
}

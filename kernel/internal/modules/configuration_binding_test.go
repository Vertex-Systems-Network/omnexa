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

func TestConfigurationBindingUsesValidatedManifestDeclarationsAndExistingRegistry(t *testing.T) {
	t.Parallel()

	registry := p0305RegistryV1(t, []string{string(p0305SettingKey)}, []string{string(p0305FlagKey)})
	binding, err := BindConfigurationDefinitions(registry, NewMemoryLifecycleStore(), []configuration.Definition{
		p0305FlagDefinition(),
		p0305SettingDefinition(),
	})
	if err != nil {
		t.Fatalf("BindConfigurationDefinitions() error = %v", err)
	}

	configRegistry, ok := binding.Registry()
	if !ok {
		t.Fatal("Registry() missing governed configuration registry")
	}
	definitions := configRegistry.Definitions()
	if len(definitions) != 2 || definitions[0].Key != p0305FlagKey || definitions[1].Key != p0305SettingKey {
		t.Fatalf("Definitions() = %#v, want deterministic key order", definitions)
	}
	if definitions[0].Owner != "omnexa.inventory" || definitions[1].Owner != "omnexa.inventory" {
		t.Fatalf("definition owners = %q/%q, want exact module owner", definitions[0].Owner, definitions[1].Owner)
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
	assertP0305BindingCode(t, BindConfigurationDefinitions(registry, store, []configuration.Definition{wrongOwner, p0305FlagDefinition()}), "module.configuration.owner_mismatch")

	wrongClass := p0305FlagDefinition()
	wrongClass.Class = configuration.ClassRuntimeConfig
	assertP0305BindingCode(t, BindConfigurationDefinitions(registry, store, []configuration.Definition{p0305SettingDefinition(), wrongClass}), "module.configuration.class_mismatch")

	assertP0305BindingCode(t, BindConfigurationDefinitions(registry, store, []configuration.Definition{p0305SettingDefinition()}), "module.configuration.definition_missing")

	extra := p0305SettingDefinition()
	extra.Key = "omnexa.inventory.undeclared_setting"
	assertP0305BindingCode(t, BindConfigurationDefinitions(registry, store, []configuration.Definition{p0305SettingDefinition(), p0305FlagDefinition(), extra}), "module.configuration.definition_undeclared")
}

func TestConfigurationBindingRejectsCrossClassDeclarationCollision(t *testing.T) {
	t.Parallel()

	registry := p0305RegistryV1(t, []string{string(p0305SettingKey)}, []string{string(p0305SettingKey)})
	assertP0305BindingCode(t, BindConfigurationDefinitions(registry, NewMemoryLifecycleStore(), nil), "module.configuration.declaration_collision")
}

func TestConfigurationBindingLifecycleReadsAreNonDestructiveAndEnabledOnlyIsRuntimeActive(t *testing.T) {
	registry := p0305RegistryV1(t, []string{string(p0305SettingKey)}, []string{string(p0305FlagKey)})
	store := NewMemoryLifecycleStore()
	binding, err := BindConfigurationDefinitions(registry, store, []configuration.Definition{p0305SettingDefinition(), p0305FlagDefinition()})
	if err != nil {
		t.Fatalf("BindConfigurationDefinitions() error = %v", err)
	}

	assertP0305ResolveCode(t, binding, p0305SettingKey, "module.configuration.unavailable")

	ctx := context.Background()
	if err := store.CompareAndSwap(ctx, "omnexa.inventory", 0, LifecycleRecord{
		ModuleID: "omnexa.inventory", Version: "1.2.3", State: LifecycleInstalled, Revision: 1,
	}); err != nil {
		t.Fatalf("seed installed state: %v", err)
	}
	installed, err := binding.Resolve(ctx, p0305SettingKey)
	if err != nil || installed.RuntimeActive || installed.LifecycleState != LifecycleInstalled {
		t.Fatalf("installed Resolve() = %#v, err=%v", installed, err)
	}

	if err := store.CompareAndSwap(ctx, "omnexa.inventory", 1, LifecycleRecord{
		ModuleID: "omnexa.inventory", Version: "1.2.3", State: LifecycleEnabled, Revision: 2,
	}); err != nil {
		t.Fatalf("seed enabled state: %v", err)
	}
	enabled, err := binding.Resolve(ctx, p0305FlagKey)
	if err != nil || !enabled.RuntimeActive || enabled.LifecycleState != LifecycleEnabled {
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
	binding, err := BindConfigurationDefinitions(registry, NewMemoryLifecycleStore(), nil)
	if err != nil {
		t.Fatalf("BindConfigurationDefinitions(empty) error = %v", err)
	}
	if _, ok := binding.Registry(); ok {
		t.Fatal("empty module declaration set unexpectedly created a configuration registry")
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

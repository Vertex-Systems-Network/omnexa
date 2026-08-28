package modules

import (
	"context"
	"sort"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/configuration"
)

// ModuleConfigurationKind identifies which manifest declaration class owns a
// configuration key. It is descriptive metadata only and grants no authority.
type ModuleConfigurationKind string

const (
	ModuleConfigurationSetting     ModuleConfigurationKind = "setting"
	ModuleConfigurationFeatureFlag ModuleConfigurationKind = "feature_flag"
)

// ModuleConfigurationDiagnostic is a stable, value-free P03.05 failure record.
// It never includes configuration values, tenant identifiers, raw manifests, or
// authorization details.
type ModuleConfigurationDiagnostic struct {
	Code     string `json:"code"`
	ModuleID string `json:"module_id,omitempty"`
	Key      string `json:"key,omitempty"`
}

// ModuleConfigurationError is the fail-closed P03.05 binding error.
type ModuleConfigurationError struct {
	Diagnostic ModuleConfigurationDiagnostic
}

func (e *ModuleConfigurationError) Error() string { return "module configuration binding failed" }

type moduleConfigurationDeclaration struct {
	moduleID string
	kind     ModuleConfigurationKind
}

// BoundConfiguration is one lifecycle-aware projection of an existing governed
// configuration definition. RuntimeActive is true only while the owning module
// is enabled; retained lifecycle states remain readable for safe history and
// re-enable without becoming runtime authority.
type BoundConfiguration struct {
	ModuleID       string
	Kind           ModuleConfigurationKind
	Definition     configuration.Definition
	LifecycleState LifecycleState
	RuntimeActive  bool
}

// ConfigurationBinding binds the validated discovery snapshot to the existing
// kernel.configuration registry. It does not store setting values, tenant scope,
// permissions, or lifecycle state itself.
type ConfigurationBinding struct {
	registry      *configuration.Registry
	declarations  map[configuration.Key]moduleConfigurationDeclaration
	lifecycleStore LifecycleStore
}

// BindConfigurationDefinitions validates exact manifest declaration-to-definition
// ownership and builds the existing governed configuration registry. Manifest
// schema is not reinterpreted or reparsed here: declarations come only from the
// validated snapshot retained by Discover.
func BindConfigurationDefinitions(
	moduleRegistry Registry,
	lifecycleStore LifecycleStore,
	definitions []configuration.Definition,
) (*ConfigurationBinding, error) {
	if lifecycleStore == nil {
		return nil, configurationBindingError("module.configuration.lifecycle_store_required", "", "")
	}

	declarations, err := collectConfigurationDeclarations(moduleRegistry)
	if err != nil {
		return nil, err
	}

	byKey := make(map[configuration.Key]configuration.Definition, len(definitions))
	for _, definition := range definitions {
		if _, exists := byKey[definition.Key]; exists {
			return nil, configurationBindingError("module.configuration.definition_duplicate", "", string(definition.Key))
		}
		declaration, ok := declarations[definition.Key]
		if !ok {
			return nil, configurationBindingError("module.configuration.definition_undeclared", "", string(definition.Key))
		}
		if definition.Owner != declaration.moduleID {
			return nil, configurationBindingError("module.configuration.owner_mismatch", declaration.moduleID, string(definition.Key))
		}
		if !configurationClassMatches(declaration.kind, definition.Class) {
			return nil, configurationBindingError("module.configuration.class_mismatch", declaration.moduleID, string(definition.Key))
		}
		if _, validationErr := configuration.NewRegistry(definition); validationErr != nil {
			return nil, configurationBindingError("module.configuration.definition_invalid", declaration.moduleID, string(definition.Key))
		}
		byKey[definition.Key] = definition
	}

	for key, declaration := range declarations {
		if _, ok := byKey[key]; !ok {
			return nil, configurationBindingError("module.configuration.definition_missing", declaration.moduleID, string(key))
		}
	}

	binding := &ConfigurationBinding{
		declarations:  copyConfigurationDeclarations(declarations),
		lifecycleStore: lifecycleStore,
	}
	if len(byKey) == 0 {
		return binding, nil
	}

	ordered := make([]configuration.Definition, 0, len(byKey))
	for _, definition := range byKey {
		ordered = append(ordered, definition)
	}
	sort.Slice(ordered, func(left, right int) bool { return ordered[left].Key < ordered[right].Key })

	registry, registryErr := configuration.NewRegistry(ordered...)
	if registryErr != nil {
		return nil, configurationBindingError("module.configuration.registry_invalid", "", "")
	}
	binding.registry = registry
	return binding, nil
}

// Registry exposes the existing governed configuration registry for integration
// with kernel.configuration services. Empty declaration sets intentionally have
// no configuration registry rather than weakening configuration.NewRegistry's
// non-empty invariant.
func (binding *ConfigurationBinding) Registry() (*configuration.Registry, bool) {
	if binding == nil || binding.registry == nil {
		return nil, false
	}
	return binding.registry, true
}

// Resolve returns a typed definition only when the owning module has reached a
// lifecycle state where retained configuration is meaningful. It never accepts
// client-supplied tenant/org identifiers; scoped values remain owned by the
// trusted kernel.configuration.ScopedService boundary.
func (binding *ConfigurationBinding) Resolve(ctx context.Context, key configuration.Key) (BoundConfiguration, error) {
	if binding == nil || binding.lifecycleStore == nil {
		return BoundConfiguration{}, configurationBindingError("module.configuration.binding_invalid", "", string(key))
	}
	declaration, ok := binding.declarations[key]
	if !ok || binding.registry == nil {
		return BoundConfiguration{}, configurationBindingError("module.configuration.definition_undeclared", "", string(key))
	}
	definition, ok := binding.registry.Definition(key)
	if !ok {
		return BoundConfiguration{}, configurationBindingError("module.configuration.registry_inconsistent", declaration.moduleID, string(key))
	}

	record, found, err := binding.lifecycleStore.Load(ctx, declaration.moduleID)
	if err != nil {
		return BoundConfiguration{}, configurationBindingError("module.configuration.lifecycle_read_failed", declaration.moduleID, string(key))
	}
	state := LifecycleAvailable
	if found {
		if record.ModuleID != declaration.moduleID {
			return BoundConfiguration{}, configurationBindingError("module.configuration.lifecycle_identity_mismatch", declaration.moduleID, string(key))
		}
		state = record.State
	}
	if !configurationReadableLifecycle(state) {
		return BoundConfiguration{}, configurationBindingError("module.configuration.unavailable", declaration.moduleID, string(key))
	}

	return BoundConfiguration{
		ModuleID:       declaration.moduleID,
		Kind:           declaration.kind,
		Definition:     definition,
		LifecycleState: state,
		RuntimeActive:  state == LifecycleEnabled,
	}, nil
}

func collectConfigurationDeclarations(registry Registry) (map[configuration.Key]moduleConfigurationDeclaration, error) {
	declarations := make(map[configuration.Key]moduleConfigurationDeclaration)
	for _, record := range registry.List() {
		snapshot, ok := registry.manifestSnapshot(record.ID)
		if !ok || snapshot.ID != record.ID || snapshot.Owner != record.Owner || snapshot.Version != record.Version {
			return nil, configurationBindingError("module.configuration.snapshot_invalid", record.ID, "")
		}
		for _, raw := range snapshot.Settings {
			if err := addConfigurationDeclaration(declarations, configuration.Key(raw), record.ID, ModuleConfigurationSetting); err != nil {
				return nil, err
			}
		}
		for _, raw := range snapshot.FeatureFlags {
			if err := addConfigurationDeclaration(declarations, configuration.Key(raw), record.ID, ModuleConfigurationFeatureFlag); err != nil {
				return nil, err
			}
		}
	}
	return declarations, nil
}

func addConfigurationDeclaration(
	declarations map[configuration.Key]moduleConfigurationDeclaration,
	key configuration.Key,
	moduleID string,
	kind ModuleConfigurationKind,
) error {
	if existing, ok := declarations[key]; ok {
		return configurationBindingError("module.configuration.declaration_collision", existing.moduleID, string(key))
	}
	declarations[key] = moduleConfigurationDeclaration{moduleID: moduleID, kind: kind}
	return nil
}

func copyConfigurationDeclarations(values map[configuration.Key]moduleConfigurationDeclaration) map[configuration.Key]moduleConfigurationDeclaration {
	result := make(map[configuration.Key]moduleConfigurationDeclaration, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}

func configurationClassMatches(kind ModuleConfigurationKind, class configuration.Class) bool {
	switch kind {
	case ModuleConfigurationSetting:
		return class == configuration.ClassRuntimeConfig
	case ModuleConfigurationFeatureFlag:
		return class == configuration.ClassFeatureFlag
	default:
		return false
	}
}

func configurationReadableLifecycle(state LifecycleState) bool {
	switch state {
	case LifecycleInstalled,
		LifecycleEnabled,
		LifecycleDisabled,
		LifecycleSuspended,
		LifecycleArchived,
		LifecycleDetached,
		LifecycleRecoveryRequired:
		return true
	case LifecycleAvailable, LifecyclePurged:
		return false
	default:
		return false
	}
}

func configurationBindingError(code, moduleID, key string) error {
	return &ModuleConfigurationError{Diagnostic: ModuleConfigurationDiagnostic{
		Code:     code,
		ModuleID: moduleID,
		Key:      key,
	}}
}

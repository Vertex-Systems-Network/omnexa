package modules

import (
	"context"
	"sort"
	"strconv"
	"strings"
)

const maxCapabilityMajor uint32 = 9999

// CapabilityRegistration binds one manifest-declared provided capability to
// descriptive public contract metadata. These references do not grant
// permission, tenant scope, handler access, or invocation authority.
type CapabilityRegistration struct {
	ModuleID         string `json:"module_id"`
	Owner            string `json:"owner"`
	Declaration      string `json:"declaration"`
	AuthorizationRef string `json:"authorization_ref"`
	ScopeRef         string `json:"scope_ref"`
	ContractRef      string `json:"contract_ref"`
}

// CapabilityQuery requires exact stable ID, major version, and owner identity.
// Lookup remains metadata-only and never invokes a capability.
type CapabilityQuery struct {
	ID    string
	Major uint32
	Owner string
}

// CapabilityRecord is a lifecycle-aware projection of one static provider
// declaration. Available is true only while the owning module is enabled.
type CapabilityRecord struct {
	ID               string         `json:"id"`
	Major            uint32         `json:"major"`
	Declaration      string         `json:"declaration"`
	ProviderModuleID string         `json:"provider_module_id"`
	ProviderOwner    string         `json:"provider_owner"`
	AuthorizationRef string         `json:"authorization_ref"`
	ScopeRef         string         `json:"scope_ref"`
	ContractRef      string         `json:"contract_ref"`
	LifecycleState   LifecycleState `json:"lifecycle_state"`
	Available        bool           `json:"available"`
}

// CapabilityConsumer is static validated declaration metadata. It expresses a
// dependency on a capability contract but grants no right to invoke it.
type CapabilityConsumer struct {
	ID               string `json:"id"`
	Major            uint32 `json:"major"`
	Declaration      string `json:"declaration"`
	ConsumerModuleID string `json:"consumer_module_id"`
	ConsumerOwner    string `json:"consumer_owner"`
}

// CapabilityDiagnostic is stable and value-free. It intentionally excludes raw
// manifests, private implementation details, tenant identifiers, secrets,
// authorization results, and database details.
type CapabilityDiagnostic struct {
	Code         string `json:"code"`
	ModuleID     string `json:"module_id,omitempty"`
	CapabilityID string `json:"capability_id,omitempty"`
	Major        uint32 `json:"major,omitempty"`
}

// CapabilityError is a fail-closed P03.06 registry error.
type CapabilityError struct {
	Diagnostic CapabilityDiagnostic
}

func (e *CapabilityError) Error() string { return "module capability registry failed" }

type parsedCapability struct {
	ID          string
	Major       uint32
	Declaration string
}

type capabilityProviderDeclaration struct {
	identity parsedCapability
	moduleID string
	owner    string
}

type boundCapabilityProvider struct {
	declaration  capabilityProviderDeclaration
	registration CapabilityRegistration
}

// CapabilityRegistry binds validated manifest declarations to explicit public
// contract metadata and accepted lifecycle state. It stores no handlers,
// permissions, tenant scope values, configuration values, or business data.
type CapabilityRegistry struct {
	providers        []boundCapabilityProvider
	byKey            map[string]boundCapabilityProvider
	consumers        []CapabilityConsumer
	consumerByModule map[string]map[string]CapabilityConsumer
	majorsByID       map[string][]uint32
	lifecycleStore   LifecycleStore
}

// BindCapabilityRegistry builds one deterministic P03.06 capability metadata
// registry from the validated snapshots already retained by module discovery.
// Raw manifests are never reparsed at this boundary.
func BindCapabilityRegistry(
	moduleRegistry Registry,
	lifecycleStore LifecycleStore,
	registrations []CapabilityRegistration,
) (*CapabilityRegistry, error) {
	if lifecycleStore == nil {
		return nil, capabilityError("module.capability.lifecycle_store_required", "", "", 0)
	}

	declarations, consumers, err := collectCapabilityDeclarations(moduleRegistry)
	if err != nil {
		return nil, err
	}

	registrationsByKey := make(map[string]CapabilityRegistration, len(registrations))
	for _, registration := range registrations {
		identity, ok := parseCapabilityDeclaration(registration.Declaration)
		if !ok {
			return nil, capabilityError("module.capability.registration_invalid", registration.ModuleID, "", 0)
		}
		key := capabilityKey(identity.ID, identity.Major)
		if _, exists := registrationsByKey[key]; exists {
			return nil, capabilityError("module.capability.registration_duplicate", registration.ModuleID, identity.ID, identity.Major)
		}
		declaration, declared := declarations[key]
		if !declared {
			return nil, capabilityError("module.capability.registration_undeclared", registration.ModuleID, identity.ID, identity.Major)
		}
		if registration.ModuleID != declaration.moduleID {
			return nil, capabilityError("module.capability.module_mismatch", declaration.moduleID, identity.ID, identity.Major)
		}
		if registration.Owner != declaration.owner {
			return nil, capabilityError("module.capability.owner_mismatch", declaration.moduleID, identity.ID, identity.Major)
		}
		if !validCapabilityMetadataRef(registration.AuthorizationRef) ||
			!validCapabilityMetadataRef(registration.ScopeRef) ||
			!validCapabilityMetadataRef(registration.ContractRef) {
			return nil, capabilityError("module.capability.metadata_ref_invalid", declaration.moduleID, identity.ID, identity.Major)
		}
		registrationsByKey[key] = registration
	}

	keys := make([]string, 0, len(declarations))
	for key := range declarations {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	providers := make([]boundCapabilityProvider, 0, len(keys))
	byKey := make(map[string]boundCapabilityProvider, len(keys))
	majorsByID := make(map[string][]uint32)
	for _, key := range keys {
		declaration := declarations[key]
		registration, ok := registrationsByKey[key]
		if !ok {
			return nil, capabilityError(
				"module.capability.registration_missing",
				declaration.moduleID,
				declaration.identity.ID,
				declaration.identity.Major,
			)
		}
		provider := boundCapabilityProvider{declaration: declaration, registration: registration}
		providers = append(providers, provider)
		byKey[key] = provider
		majorsByID[declaration.identity.ID] = append(majorsByID[declaration.identity.ID], declaration.identity.Major)
	}
	for id := range majorsByID {
		sort.Slice(majorsByID[id], func(i, j int) bool { return majorsByID[id][i] < majorsByID[id][j] })
	}

	sort.Slice(consumers, func(i, j int) bool {
		if consumers[i].ID != consumers[j].ID {
			return consumers[i].ID < consumers[j].ID
		}
		if consumers[i].Major != consumers[j].Major {
			return consumers[i].Major < consumers[j].Major
		}
		return consumers[i].ConsumerModuleID < consumers[j].ConsumerModuleID
	})
	consumerByModule := make(map[string]map[string]CapabilityConsumer)
	for _, consumer := range consumers {
		moduleConsumers := consumerByModule[consumer.ConsumerModuleID]
		if moduleConsumers == nil {
			moduleConsumers = make(map[string]CapabilityConsumer)
			consumerByModule[consumer.ConsumerModuleID] = moduleConsumers
		}
		moduleConsumers[capabilityKey(consumer.ID, consumer.Major)] = consumer
	}

	return &CapabilityRegistry{
		providers:        providers,
		byKey:            byKey,
		consumers:        append([]CapabilityConsumer(nil), consumers...),
		consumerByModule: consumerByModule,
		majorsByID:       majorsByID,
		lifecycleStore:   lifecycleStore,
	}, nil
}

// List returns all declared provider identities in deterministic order with
// current lifecycle-derived availability. Unavailable providers remain present
// for history/reference but are never marked active.
func (r *CapabilityRegistry) List(ctx context.Context) ([]CapabilityRecord, error) {
	if r == nil || r.lifecycleStore == nil {
		return nil, capabilityError("module.capability.registry_invalid", "", "", 0)
	}
	result := make([]CapabilityRecord, 0, len(r.providers))
	for _, provider := range r.providers {
		record, err := r.resolveProvider(ctx, provider)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

// Consumers returns a deterministic immutable-by-copy view of validated
// capability-consumer declarations.
func (r *CapabilityRegistry) Consumers() []CapabilityConsumer {
	if r == nil || len(r.consumers) == 0 {
		return []CapabilityConsumer{}
	}
	return append([]CapabilityConsumer(nil), r.consumers...)
}

// Lookup resolves exact stable capability ID, major version, and owner identity.
// The returned record may be unavailable; callers must not interpret metadata
// lookup as permission or invocation authority.
func (r *CapabilityRegistry) Lookup(ctx context.Context, query CapabilityQuery) (CapabilityRecord, error) {
	if r == nil || r.lifecycleStore == nil {
		return CapabilityRecord{}, capabilityError("module.capability.registry_invalid", "", query.ID, query.Major)
	}
	if !validCapabilityIdentity(query.ID, query.Major) || !validBoundedPattern(query.Owner, ownerPattern, maxIdentifierLength) {
		return CapabilityRecord{}, capabilityError("module.capability.query_invalid", "", "", 0)
	}
	provider, ok := r.byKey[capabilityKey(query.ID, query.Major)]
	if !ok {
		return CapabilityRecord{}, capabilityError("module.capability.provider_missing", "", query.ID, query.Major)
	}
	if provider.declaration.owner != query.Owner {
		return CapabilityRecord{}, capabilityError(
			"module.capability.owner_mismatch",
			provider.declaration.moduleID,
			query.ID,
			query.Major,
		)
	}
	return r.resolveProvider(ctx, provider)
}

// ResolveConsumer verifies that a module actually declared consumption of the
// exact capability major and returns only an enabled provider. Missing,
// incompatible, or lifecycle-unavailable providers fail closed.
func (r *CapabilityRegistry) ResolveConsumer(
	ctx context.Context,
	consumerModuleID string,
	declaration string,
) (CapabilityRecord, error) {
	if r == nil || r.lifecycleStore == nil {
		return CapabilityRecord{}, capabilityError("module.capability.registry_invalid", consumerModuleID, "", 0)
	}
	if !validBoundedPattern(consumerModuleID, moduleIDPattern, maxIdentifierLength) {
		return CapabilityRecord{}, capabilityError("module.capability.consumer_invalid", "", "", 0)
	}
	identity, ok := parseCapabilityDeclaration(declaration)
	if !ok {
		return CapabilityRecord{}, capabilityError("module.capability.consumer_invalid", consumerModuleID, "", 0)
	}
	moduleConsumers := r.consumerByModule[consumerModuleID]
	if moduleConsumers == nil {
		return CapabilityRecord{}, capabilityError(
			"module.capability.consumer_undeclared",
			consumerModuleID,
			identity.ID,
			identity.Major,
		)
	}
	if _, declared := moduleConsumers[capabilityKey(identity.ID, identity.Major)]; !declared {
		return CapabilityRecord{}, capabilityError(
			"module.capability.consumer_undeclared",
			consumerModuleID,
			identity.ID,
			identity.Major,
		)
	}

	provider, ok := r.byKey[capabilityKey(identity.ID, identity.Major)]
	if !ok {
		if len(r.majorsByID[identity.ID]) > 0 {
			return CapabilityRecord{}, capabilityError(
				"module.capability.version_incompatible",
				consumerModuleID,
				identity.ID,
				identity.Major,
			)
		}
		return CapabilityRecord{}, capabilityError(
			"module.capability.provider_missing",
			consumerModuleID,
			identity.ID,
			identity.Major,
		)
	}

	record, err := r.resolveProvider(ctx, provider)
	if err != nil {
		return CapabilityRecord{}, err
	}
	if !record.Available {
		return CapabilityRecord{}, capabilityError(
			"module.capability.provider_unavailable",
			record.ProviderModuleID,
			record.ID,
			record.Major,
		)
	}
	return record, nil
}

func (r *CapabilityRegistry) resolveProvider(ctx context.Context, provider boundCapabilityProvider) (CapabilityRecord, error) {
	lifecycle, found, err := r.lifecycleStore.Load(ctx, provider.declaration.moduleID)
	if err != nil {
		return CapabilityRecord{}, capabilityError(
			"module.capability.lifecycle_read_failed",
			provider.declaration.moduleID,
			provider.declaration.identity.ID,
			provider.declaration.identity.Major,
		)
	}

	state := LifecycleAvailable
	if found {
		if lifecycle.ModuleID != provider.declaration.moduleID {
			return CapabilityRecord{}, capabilityError(
				"module.capability.lifecycle_identity_mismatch",
				provider.declaration.moduleID,
				provider.declaration.identity.ID,
				provider.declaration.identity.Major,
			)
		}
		if !knownCapabilityLifecycleState(lifecycle.State) {
			return CapabilityRecord{}, capabilityError(
				"module.capability.lifecycle_state_invalid",
				provider.declaration.moduleID,
				provider.declaration.identity.ID,
				provider.declaration.identity.Major,
			)
		}
		state = lifecycle.State
	}

	return CapabilityRecord{
		ID:               provider.declaration.identity.ID,
		Major:            provider.declaration.identity.Major,
		Declaration:      provider.declaration.identity.Declaration,
		ProviderModuleID: provider.declaration.moduleID,
		ProviderOwner:    provider.declaration.owner,
		AuthorizationRef: provider.registration.AuthorizationRef,
		ScopeRef:         provider.registration.ScopeRef,
		ContractRef:      provider.registration.ContractRef,
		LifecycleState:   state,
		Available:        state == LifecycleEnabled,
	}, nil
}

func collectCapabilityDeclarations(
	registry Registry,
) (map[string]capabilityProviderDeclaration, []CapabilityConsumer, error) {
	providers := make(map[string]capabilityProviderDeclaration)
	consumers := make([]CapabilityConsumer, 0)
	ownerByID := make(map[string]string)

	for _, record := range registry.List() {
		snapshot, ok := registry.manifestSnapshot(record.ID)
		if !ok || snapshot.ID != record.ID || snapshot.Owner != record.Owner || snapshot.Version != record.Version {
			return nil, nil, capabilityError("module.capability.snapshot_invalid", record.ID, "", 0)
		}

		for _, raw := range snapshot.CapabilitiesProvided {
			identity, valid := parseCapabilityDeclaration(raw)
			if !valid {
				return nil, nil, capabilityError("module.capability.declaration_invalid", record.ID, "", 0)
			}
			if owner, exists := ownerByID[identity.ID]; exists && owner != record.Owner {
				return nil, nil, capabilityError("module.capability.owner_conflict", record.ID, identity.ID, identity.Major)
			}
			ownerByID[identity.ID] = record.Owner

			key := capabilityKey(identity.ID, identity.Major)
			if existing, exists := providers[key]; exists {
				return nil, nil, capabilityError(
					"module.capability.declaration_collision",
					existing.moduleID,
					identity.ID,
					identity.Major,
				)
			}
			providers[key] = capabilityProviderDeclaration{
				identity: identity,
				moduleID: record.ID,
				owner:    record.Owner,
			}
		}

		for _, raw := range snapshot.CapabilitiesConsumed {
			identity, valid := parseCapabilityDeclaration(raw)
			if !valid {
				return nil, nil, capabilityError("module.capability.consumer_declaration_invalid", record.ID, "", 0)
			}
			consumers = append(consumers, CapabilityConsumer{
				ID:               identity.ID,
				Major:            identity.Major,
				Declaration:      identity.Declaration,
				ConsumerModuleID: record.ID,
				ConsumerOwner:    record.Owner,
			})
		}
	}
	return providers, consumers, nil
}

func parseCapabilityDeclaration(value string) (parsedCapability, bool) {
	if !validBoundedPattern(value, identifierPattern, maxIdentifierLength) {
		return parsedCapability{}, false
	}
	index := strings.LastIndex(value, ".v")
	if index <= 0 || index+2 >= len(value) {
		return parsedCapability{}, false
	}
	id := value[:index]
	majorText := value[index+2:]
	if !validBoundedPattern(id, identifierPattern, maxIdentifierLength) ||
		len(majorText) > 4 ||
		majorText[0] == '0' {
		return parsedCapability{}, false
	}
	major64, err := strconv.ParseUint(majorText, 10, 32)
	if err != nil || major64 == 0 || major64 > uint64(maxCapabilityMajor) {
		return parsedCapability{}, false
	}
	return parsedCapability{ID: id, Major: uint32(major64), Declaration: value}, true
}

func validCapabilityIdentity(id string, major uint32) bool {
	return validBoundedPattern(id, identifierPattern, maxIdentifierLength) && major > 0 && major <= maxCapabilityMajor
}

func validCapabilityMetadataRef(value string) bool {
	return validBoundedPattern(value, identifierPattern, maxReferenceLength)
}

func capabilityKey(id string, major uint32) string {
	return id + "\x00" + strconv.FormatUint(uint64(major), 10)
}

func knownCapabilityLifecycleState(state LifecycleState) bool {
	switch state {
	case LifecycleAvailable,
		LifecycleInstalled,
		LifecycleEnabled,
		LifecycleDisabled,
		LifecycleSuspended,
		LifecycleArchived,
		LifecycleDetached,
		LifecycleRecoveryRequired,
		LifecyclePurged:
		return true
	default:
		return false
	}
}

func capabilityError(code, moduleID, capabilityID string, major uint32) error {
	return &CapabilityError{Diagnostic: CapabilityDiagnostic{
		Code:         code,
		ModuleID:     moduleID,
		CapabilityID: capabilityID,
		Major:        major,
	}}
}

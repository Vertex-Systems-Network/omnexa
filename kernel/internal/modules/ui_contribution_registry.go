package modules

import (
	"context"
	"sort"
	"strings"
)

const maxUIContributionContractVersion = 9999

// UIContributionKind is declarative presentation metadata only. It does not
// identify a renderer, component, handler, executable callback, or frontend
// authorization path.
type UIContributionKind string

const (
	UIContributionNavigation   UIContributionKind = "navigation"
	UIContributionPage         UIContributionKind = "page"
	UIContributionWidget       UIContributionKind = "widget"
	UIContributionSettings     UIContributionKind = "settings"
	UIContributionBuilderBlock UIContributionKind = "builder_block"
)

// UIFallbackBehavior describes consumer-facing degradation metadata. It never
// carries executable fallback behavior.
type UIFallbackBehavior string

const (
	UIFallbackHidden   UIFallbackBehavior = "hidden"
	UIFallbackDegraded UIFallbackBehavior = "degraded"
)

// UIContributionRegistration binds one explicit contribution identity to an
// owning discovered module and one validated manifest-declared slot. Permission,
// feature flag and optional dependency values are descriptive references only.
type UIContributionRegistration struct {
	ModuleID           string               `json:"module_id"`
	Owner              string               `json:"owner"`
	ID                 string               `json:"id"`
	Slot               string               `json:"slot"`
	Kind               UIContributionKind   `json:"kind"`
	ContractVersion    uint32               `json:"contract_version"`
	Permission         string               `json:"permission,omitempty"`
	FeatureFlag        string               `json:"feature_flag,omitempty"`
	OptionalDependency string               `json:"optional_dependency,omitempty"`
	Fallback           UIFallbackBehavior   `json:"fallback"`
}

// UIContributionRecord is a lifecycle/dependency-aware metadata projection.
// Available means only that the owning module is lifecycle-enabled. Degraded
// reports an optional-dependency condition and never changes backend authority.
type UIContributionRecord struct {
	ModuleID           string               `json:"module_id"`
	Owner              string               `json:"owner"`
	ID                 string               `json:"id"`
	Slot               string               `json:"slot"`
	Kind               UIContributionKind   `json:"kind"`
	ContractVersion    uint32               `json:"contract_version"`
	Permission         string               `json:"permission,omitempty"`
	FeatureFlag        string               `json:"feature_flag,omitempty"`
	OptionalDependency string               `json:"optional_dependency,omitempty"`
	Fallback           UIFallbackBehavior   `json:"fallback"`
	LifecycleState     LifecycleState       `json:"lifecycle_state"`
	Available          bool                 `json:"available"`
	Degraded           bool                 `json:"degraded"`
	DegradationReason  string               `json:"degradation_reason,omitempty"`
}

// UIContributionDiagnostic is stable and value-safe. It intentionally excludes
// raw manifests, tenant/organization identifiers, authorization decisions,
// configuration values, secrets and rendering implementation details.
type UIContributionDiagnostic struct {
	Code           string `json:"code"`
	ModuleID       string `json:"module_id,omitempty"`
	ContributionID string `json:"contribution_id,omitempty"`
}

// UIContributionError is a fail-closed P03.08 registry error.
type UIContributionError struct {
	Diagnostic UIContributionDiagnostic
}

func (e *UIContributionError) Error() string { return "module UI contribution registry failed" }

type boundUIContribution struct {
	registration UIContributionRegistration
	snapshot     validatedManifestSnapshot
}

// UIContributionRegistry is a deterministic declarative metadata registry. It
// stores no executable UI objects, tenant scope, permission grants, feature-flag
// values, secrets, database handles, or private module implementation objects.
type UIContributionRegistry struct {
	registry       Registry
	lifecycleStore LifecycleStore
	contributions  []boundUIContribution
	byKey          map[string]boundUIContribution
}

// BindUIContributionRegistry validates registrations against the exact validated
// discovery snapshots that established module identity. Raw manifests are never
// reparsed at this boundary.
func BindUIContributionRegistry(
	registry Registry,
	lifecycleStore LifecycleStore,
	registrations []UIContributionRegistration,
) (*UIContributionRegistry, error) {
	if lifecycleStore == nil {
		return nil, uiContributionError("module.ui.lifecycle_store_required", "", "")
	}

	contributions := make([]boundUIContribution, 0, len(registrations))
	byKey := make(map[string]boundUIContribution, len(registrations))
	for _, registration := range registrations {
		record, ok := registry.Lookup(registration.ModuleID)
		if !ok {
			return nil, uiContributionError("module.ui.module_undiscovered", registration.ModuleID, safeUIContributionID(registration.ID))
		}
		snapshot, ok := registry.manifestSnapshot(record.ID)
		if !ok || snapshot.ID != record.ID || snapshot.Owner != record.Owner || snapshot.Version != record.Version {
			return nil, uiContributionError("module.ui.snapshot_invalid", record.ID, safeUIContributionID(registration.ID))
		}
		if registration.Owner != record.Owner {
			return nil, uiContributionError("module.ui.owner_mismatch", record.ID, safeUIContributionID(registration.ID))
		}
		if !validUIContributionIdentity(registration.ModuleID, registration.ID) {
			return nil, uiContributionError("module.ui.registration_invalid", record.ID, "")
		}
		if !validUIContributionKind(registration.Kind) ||
			registration.ContractVersion == 0 ||
			registration.ContractVersion > maxUIContributionContractVersion ||
			!validUIFallback(registration.Fallback) {
			return nil, uiContributionError("module.ui.registration_invalid", record.ID, registration.ID)
		}
		if !snapshotDeclaresUIContributionSlot(snapshot, registration.Slot) {
			return nil, uiContributionError("module.ui.slot_undeclared", record.ID, registration.ID)
		}
		if registration.Permission != "" && !snapshotDeclaresString(snapshot.Permissions, registration.Permission) {
			return nil, uiContributionError("module.ui.permission_undeclared", record.ID, registration.ID)
		}
		if registration.FeatureFlag != "" && !snapshotDeclaresString(snapshot.FeatureFlags, registration.FeatureFlag) {
			return nil, uiContributionError("module.ui.feature_flag_undeclared", record.ID, registration.ID)
		}
		if registration.OptionalDependency != "" && !snapshotDeclaresOptionalDependency(snapshot, registration.OptionalDependency) {
			return nil, uiContributionError("module.ui.optional_dependency_undeclared", record.ID, registration.ID)
		}

		key := uiContributionKey(registration.ModuleID, registration.ID)
		if _, exists := byKey[key]; exists {
			return nil, uiContributionError("module.ui.registration_duplicate", record.ID, registration.ID)
		}
		bound := boundUIContribution{registration: registration, snapshot: snapshot.clone()}
		contributions = append(contributions, bound)
		byKey[key] = bound
	}

	sort.Slice(contributions, func(i, j int) bool {
		left := contributions[i].registration
		right := contributions[j].registration
		if left.ModuleID != right.ModuleID {
			return left.ModuleID < right.ModuleID
		}
		return left.ID < right.ID
	})

	return &UIContributionRegistry{
		registry:       registry,
		lifecycleStore: lifecycleStore,
		contributions:  contributions,
		byKey:          byKey,
	}, nil
}

// List returns deterministic lifecycle/dependency-aware metadata. Permission and
// feature-flag requirements remain references only and are never evaluated here.
func (r *UIContributionRegistry) List(ctx context.Context) ([]UIContributionRecord, error) {
	if r == nil || r.lifecycleStore == nil {
		return nil, uiContributionError("module.ui.registry_invalid", "", "")
	}
	result := make([]UIContributionRecord, 0, len(r.contributions))
	for _, contribution := range r.contributions {
		record, err := r.resolve(ctx, contribution)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

// Lookup returns exactly one module/contribution metadata record. It does not
// render, execute, invoke, authorize or read tenant-scoped configuration.
func (r *UIContributionRegistry) Lookup(
	ctx context.Context,
	moduleID string,
	contributionID string,
) (UIContributionRecord, error) {
	if r == nil || r.lifecycleStore == nil ||
		!validBoundedPattern(moduleID, moduleIDPattern, maxIdentifierLength) ||
		!validBoundedPattern(contributionID, identifierPattern, maxIdentifierLength) {
		return UIContributionRecord{}, uiContributionError("module.ui.query_invalid", "", "")
	}
	contribution, ok := r.byKey[uiContributionKey(moduleID, contributionID)]
	if !ok {
		return UIContributionRecord{}, uiContributionError("module.ui.unknown", moduleID, contributionID)
	}
	return r.resolve(ctx, contribution)
}

func (r *UIContributionRegistry) resolve(
	ctx context.Context,
	contribution boundUIContribution,
) (UIContributionRecord, error) {
	registration := contribution.registration
	state, err := r.moduleLifecycleState(ctx, registration.ModuleID, registration.ID)
	if err != nil {
		return UIContributionRecord{}, err
	}

	degraded := false
	reason := ""
	if registration.OptionalDependency != "" {
		degraded, reason, err = r.optionalDependencyDegradation(
			ctx,
			contribution.snapshot,
			registration.OptionalDependency,
			registration.ID,
		)
		if err != nil {
			return UIContributionRecord{}, err
		}
	}

	return UIContributionRecord{
		ModuleID:           registration.ModuleID,
		Owner:              registration.Owner,
		ID:                 registration.ID,
		Slot:               registration.Slot,
		Kind:               registration.Kind,
		ContractVersion:    registration.ContractVersion,
		Permission:         registration.Permission,
		FeatureFlag:        registration.FeatureFlag,
		OptionalDependency: registration.OptionalDependency,
		Fallback:           registration.Fallback,
		LifecycleState:     state,
		Available:          state == LifecycleEnabled,
		Degraded:           degraded,
		DegradationReason:  reason,
	}, nil
}

func (r *UIContributionRegistry) moduleLifecycleState(
	ctx context.Context,
	moduleID string,
	contributionID string,
) (LifecycleState, error) {
	lifecycle, found, err := r.lifecycleStore.Load(ctx, moduleID)
	if err != nil {
		return "", uiContributionError("module.ui.lifecycle_read_failed", moduleID, contributionID)
	}
	state := LifecycleAvailable
	if found {
		if lifecycle.ModuleID != moduleID {
			return "", uiContributionError("module.ui.lifecycle_identity_mismatch", moduleID, contributionID)
		}
		if !knownCapabilityLifecycleState(lifecycle.State) {
			return "", uiContributionError("module.ui.lifecycle_state_invalid", moduleID, contributionID)
		}
		state = lifecycle.State
	}
	return state, nil
}

func (r *UIContributionRegistry) optionalDependencyDegradation(
	ctx context.Context,
	snapshot validatedManifestSnapshot,
	dependencyID string,
	contributionID string,
) (bool, string, error) {
	requirement, ok := snapshotOptionalDependency(snapshot, dependencyID)
	if !ok {
		return false, "", uiContributionError("module.ui.optional_dependency_undeclared", snapshot.ID, contributionID)
	}
	if snapshot.SchemaVersion == SchemaVersion || requirement.Constraint == "" {
		return true, "version_contract_missing", nil
	}

	provider, exists := r.registry.Lookup(dependencyID)
	if !exists {
		return true, "dependency_missing", nil
	}
	compatible, valid := matchesDependencyConstraint(provider.Version, requirement.Constraint)
	if !valid || !compatible {
		return true, "dependency_incompatible", nil
	}
	state, err := r.moduleLifecycleState(ctx, dependencyID, contributionID)
	if err != nil {
		return false, "", err
	}
	if state != LifecycleEnabled {
		return true, "dependency_unavailable", nil
	}
	return false, "", nil
}

func validUIContributionIdentity(moduleID, contributionID string) bool {
	if !validBoundedPattern(moduleID, moduleIDPattern, maxIdentifierLength) ||
		!validBoundedPattern(contributionID, identifierPattern, maxIdentifierLength) {
		return false
	}
	moduleNamespace := strings.TrimPrefix(moduleID, "omnexa.")
	moduleNamespace = strings.Split(moduleNamespace, ".")[0]
	moduleNamespace = strings.Split(moduleNamespace, "_")[0]
	moduleNamespace = strings.Split(moduleNamespace, "-")[0]
	parts := strings.Split(contributionID, ".")
	return moduleNamespace != "" && len(parts) >= 2 && parts[0] == moduleNamespace
}

func validUIContributionKind(kind UIContributionKind) bool {
	switch kind {
	case UIContributionNavigation,
		UIContributionPage,
		UIContributionWidget,
		UIContributionSettings,
		UIContributionBuilderBlock:
		return true
	default:
		return false
	}
}

func validUIFallback(fallback UIFallbackBehavior) bool {
	return fallback == UIFallbackHidden || fallback == UIFallbackDegraded
}

func snapshotDeclaresUIContributionSlot(snapshot validatedManifestSnapshot, slot string) bool {
	return validBoundedPattern(slot, identifierPattern, maxIdentifierLength) && snapshotDeclaresString(snapshot.UISlots, slot)
}

func snapshotDeclaresString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func snapshotDeclaresOptionalDependency(snapshot validatedManifestSnapshot, dependencyID string) bool {
	_, ok := snapshotOptionalDependency(snapshot, dependencyID)
	return ok
}

func snapshotOptionalDependency(snapshot validatedManifestSnapshot, dependencyID string) (DependencyRequirement, bool) {
	for _, requirement := range snapshot.OptionalDependencies {
		if requirement.ID == dependencyID {
			return requirement, true
		}
	}
	return DependencyRequirement{}, false
}

func uiContributionKey(moduleID, contributionID string) string {
	return moduleID + "\x00" + contributionID
}

func safeUIContributionID(value string) string {
	if validBoundedPattern(value, identifierPattern, maxIdentifierLength) {
		return value
	}
	return ""
}

func uiContributionError(code, moduleID, contributionID string) error {
	return &UIContributionError{Diagnostic: UIContributionDiagnostic{
		Code:           code,
		ModuleID:       moduleID,
		ContributionID: contributionID,
	}}
}

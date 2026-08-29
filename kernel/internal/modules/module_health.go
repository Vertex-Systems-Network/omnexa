package modules

import (
	"context"
	"errors"
	"sort"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/operations"
)

// ModuleHealthState is the P03.10 module-level diagnostic state. It is evidence
// only and grants no authorization, lifecycle, migration, or invocation authority.
type ModuleHealthState string

const (
	ModuleHealthHealthy     ModuleHealthState = "healthy"
	ModuleHealthDegraded    ModuleHealthState = "degraded"
	ModuleHealthUnavailable ModuleHealthState = "unavailable"
	ModuleHealthFailed      ModuleHealthState = "failed"
)

const (
	healthReasonLifecycleReadFailed        = "module.health.lifecycle.read_failed"
	healthReasonLifecycleIdentityMismatch  = "module.health.lifecycle.identity_mismatch"
	healthReasonLifecycleStateInvalid      = "module.health.lifecycle.state_invalid"
	healthReasonLifecycleNotEnabled        = "module.health.lifecycle.not_enabled"
	healthReasonLifecycleRecoveryRequired  = "module.health.lifecycle.recovery_required"
	healthReasonDependencyResolutionFailed = "module.health.dependency.resolution_failed"
	healthReasonRequiredUnavailable        = "module.health.dependency.required_unavailable"
	healthReasonOptionalContractMissing    = "module.health.dependency.optional_version_contract_missing"
	healthReasonOptionalMissing            = "module.health.dependency.optional_missing"
	healthReasonOptionalIncompatible       = "module.health.dependency.optional_incompatible"
	healthReasonOptionalUnavailable        = "module.health.dependency.optional_unavailable"
	healthReasonMigrationRegistryInvalid   = "module.health.migration.registry_invalid"
	healthReasonMigrationObserverMissing   = "module.health.migration.observer_missing"
	healthReasonMigrationObservationFailed = "module.health.migration.observation_failed"
	healthReasonMigrationInconsistent      = "module.health.migration.inconsistent"
	healthReasonMigrationPending           = "module.health.migration.pending"
	healthReasonMigrationFailed            = "module.health.migration.failed"
)

// ModuleSurfaceSummary is a bounded availability count for declarative module
// metadata. It contains no handlers, grants, tenant scope, or configuration values.
type ModuleSurfaceSummary struct {
	Declared  int `json:"declared"`
	Available int `json:"available"`
	Degraded  int `json:"degraded,omitempty"`
}

// ModuleDependencySummary reports dependency counts without exposing manifest
// constraints or topology beyond the dependency identities already represented by
// bounded diagnostic reason codes.
type ModuleDependencySummary struct {
	Required         int `json:"required"`
	Optional         int `json:"optional"`
	OptionalDegraded int `json:"optional_degraded"`
}

// ModuleMigrationSummary reports only bounded migration counts and state. Exact
// migration identities remain internal validation evidence and are not projected.
type ModuleMigrationSummary struct {
	State   ModuleHealthState `json:"state"`
	Total   int               `json:"total"`
	Applied int               `json:"applied"`
	Pending int               `json:"pending"`
	Failed  int               `json:"failed"`
}

// ModuleHealthRecord is the classification-safe health projection for one module.
// Reasons are stable codes only; raw errors, paths, SQL, credentials, tenant/org
// authority, stack traces and provider payloads are intentionally absent.
type ModuleHealthRecord struct {
	ModuleID        string                  `json:"module_id"`
	Version         string                  `json:"version"`
	LifecycleState  LifecycleState          `json:"lifecycle_state"`
	State           ModuleHealthState       `json:"state"`
	Reasons         []string                `json:"reasons"`
	Dependencies    ModuleDependencySummary `json:"dependencies"`
	Migrations      ModuleMigrationSummary  `json:"migrations"`
	Capabilities    ModuleSurfaceSummary    `json:"capabilities"`
	Permissions     ModuleSurfaceSummary    `json:"permissions"`
	UIContributions ModuleSurfaceSummary    `json:"ui_contributions"`
}

// ModuleHealthReport reuses the retained P01.08 readiness vocabulary while
// preserving richer P03.10 per-module diagnostic states.
type ModuleHealthReport struct {
	Readiness operations.State     `json:"readiness"`
	Modules   []ModuleHealthRecord `json:"modules"`
}

// MigrationHealthObservation is caller-supplied runtime evidence partitioned by
// exact P03.09 migration record identity. P03.10 never queries or executes SQL.
type MigrationHealthObservation struct {
	AppliedIDs []string
	PendingIDs []string
	FailedIDs  []string
}

// ModuleMigrationHealthSource is the narrow runtime evidence seam for migration
// state. The caller receives an immutable copy of P03.09-owned expected records;
// no database handle, SQL text, filesystem path, or execution callback is exposed.
type ModuleMigrationHealthSource interface {
	ObserveModuleMigrations(context.Context, string, []MigrationRecord) (MigrationHealthObservation, error)
}

// ModuleHealthDiagnostic is a stable configuration failure code.
type ModuleHealthDiagnostic struct {
	Code string `json:"code"`
}

// ModuleHealthError is intentionally value-free and classification-safe.
type ModuleHealthError struct {
	Diagnostic ModuleHealthDiagnostic
}

func (e *ModuleHealthError) Error() string { return "module health reporter configuration failed" }

// ModuleHealthReporter aggregates accepted P03 metadata and caller-owned runtime
// evidence. It is read-only: no lifecycle mutation, permission grant, migration
// execution, UI rendering, network access, or persistence is performed.
type ModuleHealthReporter struct {
	registry           Registry
	platform           PlatformSnapshot
	lifecycleStore     LifecycleStore
	migrationRegistry  *MigrationOwnershipRegistry
	migrationSource    ModuleMigrationHealthSource
	capabilityRegistry *CapabilityRegistry
	permissionRegistry *PermissionRegistry
	uiRegistry         *UIContributionRegistry
}

// NewModuleHealthReporter binds the reporter to already-validated P03 registries
// and the retained lifecycle store. Migration observations remain optional at
// construction time; modules with expected migrations fail closed as unavailable
// until a source supplies exact identity-partitioned evidence.
func NewModuleHealthReporter(
	registry Registry,
	platform PlatformSnapshot,
	lifecycleStore LifecycleStore,
	migrationRegistry *MigrationOwnershipRegistry,
	migrationSource ModuleMigrationHealthSource,
	capabilityRegistry *CapabilityRegistry,
	permissionRegistry *PermissionRegistry,
	uiRegistry *UIContributionRegistry,
) (*ModuleHealthReporter, error) {
	if lifecycleStore == nil || migrationRegistry == nil || capabilityRegistry == nil || permissionRegistry == nil || uiRegistry == nil {
		return nil, moduleHealthError("module.health.configuration.required_input_missing")
	}
	if _, ok := parseStrictSemVer(platform.Version); !ok {
		return nil, moduleHealthError("module.health.configuration.platform_invalid")
	}
	if _, diagnostics := validatePlatformSnapshot(platform); len(diagnostics) != 0 {
		return nil, moduleHealthError("module.health.configuration.platform_invalid")
	}
	for _, module := range registry.List() {
		current, ok := migrationRegistry.moduleCurrentVersions[module.ID]
		if !ok || current != module.Version {
			return nil, moduleHealthError("module.health.configuration.migration_registry_mismatch")
		}
	}
	return &ModuleHealthReporter{
		registry:           registry,
		platform:           platform,
		lifecycleStore:     lifecycleStore,
		migrationRegistry:  migrationRegistry,
		migrationSource:    migrationSource,
		capabilityRegistry: capabilityRegistry,
		permissionRegistry: permissionRegistry,
		uiRegistry:         uiRegistry,
	}, nil
}

// Report returns deterministic, failure-isolated diagnostics. A failing module
// may make its required dependents unavailable, but unrelated modules continue to
// be evaluated and reported from their own authoritative inputs.
func (r *ModuleHealthReporter) Report(ctx context.Context) ModuleHealthReport {
	if ctx == nil {
		ctx = context.Background()
	}
	if r == nil {
		return ModuleHealthReport{Readiness: operations.StateUnready, Modules: []ModuleHealthRecord{}}
	}

	resolverReasons, globalResolverReasons := r.dependencyResolutionReasons()
	modules := make([]ModuleHealthRecord, 0, len(r.registry.List()))
	for _, module := range r.registry.List() {
		record := r.reportModule(ctx, module, resolverReasons[module.ID], globalResolverReasons)
		modules = append(modules, record)
	}
	return ModuleHealthReport{Readiness: moduleReadiness(modules), Modules: modules}
}

func (r *ModuleHealthReporter) reportModule(
	ctx context.Context,
	module RegistryRecord,
	resolverReasons []string,
	globalResolverReasons []string,
) ModuleHealthRecord {
	record := ModuleHealthRecord{
		ModuleID: module.ID,
		Version:  module.Version,
		State:    ModuleHealthHealthy,
		Reasons:  []string{},
		Migrations: ModuleMigrationSummary{
			State: ModuleHealthHealthy,
		},
	}

	state, lifecycleHealth, lifecycleReason := r.moduleLifecycle(ctx, module.ID)
	record.LifecycleState = state
	record.State = worsenModuleHealth(record.State, lifecycleHealth)
	appendHealthReason(&record.Reasons, lifecycleReason)

	for _, reason := range globalResolverReasons {
		record.State = worsenModuleHealth(record.State, ModuleHealthUnavailable)
		appendHealthReason(&record.Reasons, reason)
	}
	for _, reason := range resolverReasons {
		record.State = worsenModuleHealth(record.State, ModuleHealthUnavailable)
		appendHealthReason(&record.Reasons, reason)
	}

	optionalIssues := make(map[string]struct{})
	snapshot, snapshotOK := r.registry.manifestSnapshot(module.ID)
	if !snapshotOK || snapshot.ID != module.ID || snapshot.Version != module.Version || snapshot.Owner != module.Owner {
		record.State = worsenModuleHealth(record.State, ModuleHealthFailed)
		appendHealthReason(&record.Reasons, "module.health.registry.snapshot_invalid")
	} else {
		record.Dependencies.Required = len(snapshot.RequiredDependencies)
		record.Dependencies.Optional = len(snapshot.OptionalDependencies)
		for _, requirement := range snapshot.RequiredDependencies {
			if r.requiredDependencyUnavailable(ctx, requirement) {
				record.State = worsenModuleHealth(record.State, ModuleHealthUnavailable)
				appendHealthReason(&record.Reasons, healthReasonRequiredUnavailable)
			}
		}
		for _, requirement := range snapshot.OptionalDependencies {
			degraded, reason := r.optionalDependencyDegraded(ctx, snapshot, requirement)
			if !degraded {
				continue
			}
			record.Dependencies.OptionalDegraded++
			optionalIssues[requirement.ID] = struct{}{}
			record.State = worsenModuleHealth(record.State, ModuleHealthDegraded)
			appendHealthReason(&record.Reasons, reason)
		}
	}

	migrationSummary, migrationHealth, migrationReason := r.migrationHealth(ctx, module)
	record.Migrations = migrationSummary
	record.State = worsenModuleHealth(record.State, migrationHealth)
	appendHealthReason(&record.Reasons, migrationReason)

	record.Capabilities = r.capabilitySummary(module.ID, state)
	record.Permissions = r.permissionSummary(module.ID, state)
	record.UIContributions = r.uiSummary(module.ID, state, optionalIssues)

	sort.Strings(record.Reasons)
	return record
}

func (r *ModuleHealthReporter) dependencyResolutionReasons() (map[string][]string, []string) {
	perModule := make(map[string][]string)
	_, err := ResolveDependencies(r.registry, r.platform, nil)
	if err == nil {
		return perModule, nil
	}
	var resolutionErr *ResolutionErrors
	if !errors.As(err, &resolutionErr) {
		return perModule, []string{healthReasonDependencyResolutionFailed}
	}
	global := make([]string, 0)
	for _, diagnostic := range resolutionErr.Diagnostics() {
		if diagnostic.ModuleID == "" {
			appendHealthReason(&global, diagnostic.Code)
			continue
		}
		reasons := perModule[diagnostic.ModuleID]
		appendHealthReason(&reasons, diagnostic.Code)
		perModule[diagnostic.ModuleID] = reasons
	}
	for moduleID := range perModule {
		sort.Strings(perModule[moduleID])
	}
	sort.Strings(global)
	return perModule, global
}

func (r *ModuleHealthReporter) moduleLifecycle(ctx context.Context, moduleID string) (LifecycleState, ModuleHealthState, string) {
	lifecycle, found, err := r.lifecycleStore.Load(ctx, moduleID)
	if err != nil {
		return LifecycleAvailable, ModuleHealthFailed, healthReasonLifecycleReadFailed
	}
	if !found {
		return LifecycleAvailable, ModuleHealthUnavailable, healthReasonLifecycleNotEnabled
	}
	if lifecycle.ModuleID != moduleID {
		return LifecycleAvailable, ModuleHealthFailed, healthReasonLifecycleIdentityMismatch
	}
	if !knownCapabilityLifecycleState(lifecycle.State) {
		return lifecycle.State, ModuleHealthFailed, healthReasonLifecycleStateInvalid
	}
	if lifecycle.State == LifecycleRecoveryRequired {
		return lifecycle.State, ModuleHealthFailed, healthReasonLifecycleRecoveryRequired
	}
	if lifecycle.State != LifecycleEnabled {
		return lifecycle.State, ModuleHealthUnavailable, healthReasonLifecycleNotEnabled
	}
	return lifecycle.State, ModuleHealthHealthy, ""
}

func (r *ModuleHealthReporter) requiredDependencyUnavailable(ctx context.Context, requirement DependencyRequirement) bool {
	provider, exists := r.registry.Lookup(requirement.ID)
	if !exists {
		return true
	}
	if requirement.Constraint == "" {
		return true
	}
	compatible, valid := matchesDependencyConstraint(provider.Version, requirement.Constraint)
	if !valid || !compatible {
		return true
	}
	state, health, _ := r.moduleLifecycle(ctx, requirement.ID)
	return health == ModuleHealthFailed || state != LifecycleEnabled
}

func (r *ModuleHealthReporter) optionalDependencyDegraded(
	ctx context.Context,
	snapshot validatedManifestSnapshot,
	requirement DependencyRequirement,
) (bool, string) {
	if snapshot.SchemaVersion == SchemaVersion || requirement.Constraint == "" {
		return true, healthReasonOptionalContractMissing
	}
	provider, exists := r.registry.Lookup(requirement.ID)
	if !exists {
		return true, healthReasonOptionalMissing
	}
	compatible, valid := matchesDependencyConstraint(provider.Version, requirement.Constraint)
	if !valid || !compatible {
		return true, healthReasonOptionalIncompatible
	}
	state, health, _ := r.moduleLifecycle(ctx, requirement.ID)
	if health == ModuleHealthFailed || state != LifecycleEnabled {
		return true, healthReasonOptionalUnavailable
	}
	return false, ""
}

func (r *ModuleHealthReporter) migrationHealth(
	ctx context.Context,
	module RegistryRecord,
) (ModuleMigrationSummary, ModuleHealthState, string) {
	summary := ModuleMigrationSummary{State: ModuleHealthHealthy}
	expected, err := r.migrationRegistry.FreshInstallPlan(module.ID)
	if err != nil || r.migrationRegistry.moduleCurrentVersions[module.ID] != module.Version {
		summary.State = ModuleHealthFailed
		return summary, ModuleHealthFailed, healthReasonMigrationRegistryInvalid
	}
	summary.Total = len(expected)
	if len(expected) == 0 {
		return summary, ModuleHealthHealthy, ""
	}
	if r.migrationSource == nil {
		summary.State = ModuleHealthUnavailable
		summary.Pending = len(expected)
		return summary, ModuleHealthUnavailable, healthReasonMigrationObserverMissing
	}

	observation, observeErr := r.migrationSource.ObserveModuleMigrations(ctx, module.ID, append([]MigrationRecord(nil), expected...))
	if observeErr != nil {
		summary.State = ModuleHealthFailed
		return summary, ModuleHealthFailed, healthReasonMigrationObservationFailed
	}
	applied, pending, failed, valid := validateMigrationObservation(expected, observation)
	summary.Applied = applied
	summary.Pending = pending
	summary.Failed = failed
	if !valid {
		summary.State = ModuleHealthFailed
		return summary, ModuleHealthFailed, healthReasonMigrationInconsistent
	}
	if failed > 0 {
		summary.State = ModuleHealthFailed
		return summary, ModuleHealthFailed, healthReasonMigrationFailed
	}
	if pending > 0 {
		summary.State = ModuleHealthUnavailable
		return summary, ModuleHealthUnavailable, healthReasonMigrationPending
	}
	return summary, ModuleHealthHealthy, ""
}

func validateMigrationObservation(expected []MigrationRecord, observation MigrationHealthObservation) (int, int, int, bool) {
	expectedIDs := make(map[string]struct{}, len(expected))
	for _, record := range expected {
		expectedIDs[record.ID] = struct{}{}
	}
	seen := make(map[string]struct{}, len(expected))
	validIDs := func(ids []string) bool {
		for _, id := range ids {
			if _, ok := expectedIDs[id]; !ok {
				return false
			}
			if _, duplicate := seen[id]; duplicate {
				return false
			}
			seen[id] = struct{}{}
		}
		return true
	}
	if !validIDs(observation.AppliedIDs) || !validIDs(observation.PendingIDs) || !validIDs(observation.FailedIDs) {
		return len(observation.AppliedIDs), len(observation.PendingIDs), len(observation.FailedIDs), false
	}
	return len(observation.AppliedIDs), len(observation.PendingIDs), len(observation.FailedIDs), len(seen) == len(expectedIDs)
}

func (r *ModuleHealthReporter) capabilitySummary(moduleID string, lifecycle LifecycleState) ModuleSurfaceSummary {
	summary := ModuleSurfaceSummary{}
	for _, provider := range r.capabilityRegistry.providers {
		if provider.declaration.moduleID != moduleID {
			continue
		}
		summary.Declared++
		if lifecycle == LifecycleEnabled {
			summary.Available++
		}
	}
	return summary
}

func (r *ModuleHealthReporter) permissionSummary(moduleID string, lifecycle LifecycleState) ModuleSurfaceSummary {
	summary := ModuleSurfaceSummary{}
	for _, declaration := range r.permissionRegistry.declarations {
		if declaration.moduleID != moduleID {
			continue
		}
		summary.Declared++
		if lifecycle == LifecycleEnabled {
			summary.Available++
		}
	}
	return summary
}

func (r *ModuleHealthReporter) uiSummary(
	moduleID string,
	lifecycle LifecycleState,
	optionalIssues map[string]struct{},
) ModuleSurfaceSummary {
	summary := ModuleSurfaceSummary{}
	for _, contribution := range r.uiRegistry.contributions {
		registration := contribution.registration
		if registration.ModuleID != moduleID {
			continue
		}
		summary.Declared++
		if lifecycle == LifecycleEnabled {
			summary.Available++
		}
		if registration.OptionalDependency != "" {
			if _, degraded := optionalIssues[registration.OptionalDependency]; degraded {
				summary.Degraded++
			}
		}
	}
	return summary
}

func moduleReadiness(records []ModuleHealthRecord) operations.State {
	state := operations.StateHealthy
	for _, record := range records {
		switch record.State {
		case ModuleHealthFailed, ModuleHealthUnavailable:
			return operations.StateUnready
		case ModuleHealthDegraded:
			state = operations.StateDegraded
		}
	}
	return state
}

func worsenModuleHealth(current, candidate ModuleHealthState) ModuleHealthState {
	if moduleHealthRank(candidate) > moduleHealthRank(current) {
		return candidate
	}
	return current
}

func moduleHealthRank(state ModuleHealthState) int {
	switch state {
	case ModuleHealthHealthy:
		return 0
	case ModuleHealthDegraded:
		return 1
	case ModuleHealthUnavailable:
		return 2
	case ModuleHealthFailed:
		return 3
	default:
		return 3
	}
}

func appendHealthReason(reasons *[]string, reason string) {
	if reason == "" {
		return
	}
	for _, existing := range *reasons {
		if existing == reason {
			return
		}
	}
	*reasons = append(*reasons, reason)
}

func moduleHealthError(code string) error {
	return &ModuleHealthError{Diagnostic: ModuleHealthDiagnostic{Code: code}}
}

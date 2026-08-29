package modules

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/operations"
)

type moduleHealthFixture struct {
	registry     Registry
	store        *MemoryLifecycleStore
	migrations   *MigrationOwnershipRegistry
	capabilities *CapabilityRegistry
	permissions  *PermissionRegistry
	ui           *UIContributionRegistry
}

type fixtureMigrationHealthSource struct {
	modes map[string]string
}

func (s fixtureMigrationHealthSource) ObserveModuleMigrations(
	_ context.Context,
	moduleID string,
	expected []MigrationRecord,
) (MigrationHealthObservation, error) {
	ids := make([]string, 0, len(expected))
	for _, record := range expected {
		ids = append(ids, record.ID)
	}
	if len(ids) == 0 {
		return MigrationHealthObservation{}, nil
	}
	switch s.modes[moduleID] {
	case "pending":
		return MigrationHealthObservation{AppliedIDs: append([]string(nil), ids[1:]...), PendingIDs: []string{ids[0]}}, nil
	case "failed":
		return MigrationHealthObservation{AppliedIDs: append([]string(nil), ids[1:]...), FailedIDs: []string{ids[0]}}, nil
	case "inconsistent":
		return MigrationHealthObservation{AppliedIDs: append([]string(nil), ids...), PendingIDs: []string{"unknown.migration.identity"}}, nil
	case "error":
		return MigrationHealthObservation{}, errors.New("synthetic migration observer failure")
	default:
		return MigrationHealthObservation{AppliedIDs: ids}, nil
	}
}

type selectiveLifecycleStore struct {
	delegate LifecycleStore
	failID   string
}

func (s selectiveLifecycleStore) Load(ctx context.Context, moduleID string) (LifecycleRecord, bool, error) {
	if moduleID == s.failID {
		return LifecycleRecord{}, false, errors.New("synthetic lifecycle read failure")
	}
	return s.delegate.Load(ctx, moduleID)
}

func (s selectiveLifecycleStore) CompareAndSwap(
	ctx context.Context,
	moduleID string,
	expectedRevision uint64,
	next LifecycleRecord,
) error {
	return s.delegate.CompareAndSwap(ctx, moduleID, expectedRevision, next)
}

func TestModuleHealthReporterHealthyAndDeterministic(t *testing.T) {
	fixture := newModuleHealthFixture(t, true, true)
	reporter := fixture.reporter(t, fixture.store, fixtureMigrationHealthSource{})

	first := reporter.Report(context.Background())
	second := reporter.Report(context.Background())
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("health report is not deterministic:\nfirst=%#v\nsecond=%#v", first, second)
	}
	if first.Readiness != operations.StateHealthy {
		t.Fatalf("readiness = %q, want healthy", first.Readiness)
	}
	ids := make([]string, 0, len(first.Modules))
	for _, record := range first.Modules {
		ids = append(ids, record.ModuleID)
	}
	if !sort.StringsAreSorted(ids) {
		t.Fatalf("module order is not deterministic: %v", ids)
	}

	inventory := healthRecordByID(t, first, "omnexa.inventory")
	if inventory.State != ModuleHealthHealthy || inventory.LifecycleState != LifecycleEnabled {
		t.Fatalf("unexpected inventory health: %#v", inventory)
	}
	if inventory.Dependencies.Required != 1 || inventory.Dependencies.Optional != 1 || inventory.Dependencies.OptionalDegraded != 0 {
		t.Fatalf("unexpected dependency summary: %#v", inventory.Dependencies)
	}
	if inventory.Migrations.Total != 1 || inventory.Migrations.Applied != 1 || inventory.Migrations.Pending != 0 || inventory.Migrations.Failed != 0 || inventory.Migrations.State != ModuleHealthHealthy {
		t.Fatalf("unexpected migration summary: %#v", inventory.Migrations)
	}
	if inventory.Capabilities.Declared != 1 || inventory.Capabilities.Available != 1 ||
		inventory.Permissions.Declared != 1 || inventory.Permissions.Available != 1 ||
		inventory.UIContributions.Declared != 1 || inventory.UIContributions.Available != 1 {
		t.Fatalf("unexpected surface summaries: capabilities=%#v permissions=%#v ui=%#v", inventory.Capabilities, inventory.Permissions, inventory.UIContributions)
	}
	if len(inventory.Reasons) != 0 {
		t.Fatalf("healthy inventory has reasons: %v", inventory.Reasons)
	}
}

func TestModuleHealthReporterOptionalDependencySelectiveDegradation(t *testing.T) {
	fixture := newModuleHealthFixture(t, true, false)
	reporter := fixture.reporter(t, fixture.store, fixtureMigrationHealthSource{})
	report := reporter.Report(context.Background())

	inventory := healthRecordByID(t, report, "omnexa.inventory")
	catalog := healthRecordByID(t, report, "omnexa.catalog")
	if inventory.State != ModuleHealthDegraded || inventory.Dependencies.OptionalDegraded != 1 {
		t.Fatalf("optional absence did not selectively degrade inventory: %#v", inventory)
	}
	if inventory.UIContributions.Degraded != 1 {
		t.Fatalf("UI optional dependency degradation not summarized: %#v", inventory.UIContributions)
	}
	if !containsHealthReason(inventory.Reasons, healthReasonOptionalMissing) {
		t.Fatalf("missing bounded optional reason: %v", inventory.Reasons)
	}
	if catalog.State != ModuleHealthHealthy {
		t.Fatalf("unrelated catalog module was degraded: %#v", catalog)
	}
	if report.Readiness != operations.StateDegraded {
		t.Fatalf("readiness = %q, want degraded", report.Readiness)
	}
}

func TestModuleHealthReporterRequiredDependencyFailsClosed(t *testing.T) {
	fixture := newModuleHealthFixture(t, false, true)
	reporter := fixture.reporter(t, fixture.store, fixtureMigrationHealthSource{})
	report := reporter.Report(context.Background())

	inventory := healthRecordByID(t, report, "omnexa.inventory")
	analytics := healthRecordByID(t, report, "omnexa.analytics")
	if inventory.State != ModuleHealthUnavailable {
		t.Fatalf("required missing dependency did not fail closed: %#v", inventory)
	}
	if !containsHealthReason(inventory.Reasons, "resolver.dependency.required_missing") {
		t.Fatalf("retained resolver reason missing: %v", inventory.Reasons)
	}
	if analytics.State != ModuleHealthHealthy {
		t.Fatalf("unrelated analytics module was affected: %#v", analytics)
	}
	if report.Readiness != operations.StateUnready {
		t.Fatalf("readiness = %q, want unready", report.Readiness)
	}
}

func TestModuleHealthReporterMigrationPendingAndInconsistentNeverHealthy(t *testing.T) {
	fixture := newModuleHealthFixture(t, true, true)

	pending := fixture.reporter(t, fixture.store, fixtureMigrationHealthSource{modes: map[string]string{"omnexa.inventory": "pending"}}).Report(context.Background())
	pendingInventory := healthRecordByID(t, pending, "omnexa.inventory")
	if pendingInventory.State != ModuleHealthUnavailable || pendingInventory.Migrations.State != ModuleHealthUnavailable || pendingInventory.Migrations.Pending != 1 {
		t.Fatalf("pending migration was reported too healthy: %#v", pendingInventory)
	}
	if !containsHealthReason(pendingInventory.Reasons, healthReasonMigrationPending) {
		t.Fatalf("pending migration reason missing: %v", pendingInventory.Reasons)
	}

	inconsistent := fixture.reporter(t, fixture.store, fixtureMigrationHealthSource{modes: map[string]string{"omnexa.inventory": "inconsistent"}}).Report(context.Background())
	inconsistentInventory := healthRecordByID(t, inconsistent, "omnexa.inventory")
	if inconsistentInventory.State != ModuleHealthFailed || inconsistentInventory.Migrations.State != ModuleHealthFailed {
		t.Fatalf("inconsistent migration observation was reported too healthy: %#v", inconsistentInventory)
	}
	if !containsHealthReason(inconsistentInventory.Reasons, healthReasonMigrationInconsistent) {
		t.Fatalf("migration inconsistency reason missing: %v", inconsistentInventory.Reasons)
	}
}

func TestModuleHealthReporterFailureIsolationAndLifecycleChanges(t *testing.T) {
	fixture := newModuleHealthFixture(t, true, true)
	failingStore := selectiveLifecycleStore{delegate: fixture.store, failID: "omnexa.inventory"}
	report := fixture.reporter(t, failingStore, fixtureMigrationHealthSource{}).Report(context.Background())
	if healthRecordByID(t, report, "omnexa.inventory").State != ModuleHealthFailed {
		t.Fatalf("lifecycle read failure did not fail module health: %#v", report)
	}
	if healthRecordByID(t, report, "omnexa.catalog").State != ModuleHealthHealthy || healthRecordByID(t, report, "omnexa.analytics").State != ModuleHealthHealthy {
		t.Fatalf("one module failure contaminated unrelated module health: %#v", report)
	}

	current, found, err := fixture.store.Load(context.Background(), "omnexa.inventory")
	if err != nil || !found {
		t.Fatalf("load inventory lifecycle: found=%v err=%v", found, err)
	}
	next := current
	next.State = LifecycleDisabled
	next.Revision = current.Revision + 1
	if casErr := fixture.store.CompareAndSwap(context.Background(), current.ModuleID, current.Revision, next); casErr != nil {
		t.Fatalf("disable inventory fixture: %v", casErr)
	}
	changed := fixture.reporter(t, fixture.store, fixtureMigrationHealthSource{}).Report(context.Background())
	inventory := healthRecordByID(t, changed, "omnexa.inventory")
	if inventory.State != ModuleHealthUnavailable || inventory.LifecycleState != LifecycleDisabled {
		t.Fatalf("lifecycle state change not reflected: %#v", inventory)
	}
	if inventory.Capabilities.Available != 0 || inventory.Permissions.Available != 0 || inventory.UIContributions.Available != 0 {
		t.Fatalf("disabled module surfaces remained available: %#v", inventory)
	}
}

func TestModuleHealthReporterClassificationSafeSurface(t *testing.T) {
	for _, value := range []any{
		ModuleHealthRecord{},
		ModuleHealthReport{},
		ModuleSurfaceSummary{},
		ModuleDependencySummary{},
		ModuleMigrationSummary{},
		MigrationHealthObservation{},
		ModuleHealthDiagnostic{},
	} {
		typeOf := reflect.TypeOf(value)
		for index := 0; index < typeOf.NumField(); index++ {
			field := typeOf.Field(index)
			lower := strings.ToLower(field.Name)
			for _, forbidden := range []string{"secret", "credential", "password", "token", "tenant", "organization", "stack", "path", "sql", "handler", "callback", "rawerror"} {
				if strings.Contains(lower, forbidden) {
					t.Fatalf("classification-unsafe field %s.%s", typeOf.Name(), field.Name)
				}
			}
			if field.Type.Kind() == reflect.Func {
				t.Fatalf("executable field exposed in %s.%s", typeOf.Name(), field.Name)
			}
		}
	}
}

func TestModuleHealthReporterRejectsInvalidConfiguration(t *testing.T) {
	fixture := newModuleHealthFixture(t, true, true)
	_, err := NewModuleHealthReporter(
		fixture.registry,
		PlatformSnapshot{Version: "1.0.0"},
		nil,
		fixture.migrations,
		fixtureMigrationHealthSource{},
		fixture.capabilities,
		fixture.permissions,
		fixture.ui,
	)
	var healthErr *ModuleHealthError
	if !errors.As(err, &healthErr) || healthErr.Diagnostic.Code != "module.health.configuration.required_input_missing" {
		t.Fatalf("unexpected configuration error: %#v", err)
	}
}

func newModuleHealthFixture(t *testing.T, includeCatalog, includeAnalytics bool) moduleHealthFixture {
	t.Helper()

	inventory := healthManifestFixture("omnexa.inventory", "Inventory", "inventory.team", "1.2.3")
	inventory.Dependencies = []DependencyRequirement{{ID: "omnexa.catalog", Constraint: ">=1.0.0 <2.0.0"}}
	inventory.OptionalDependencies = []DependencyRequirement{{ID: "omnexa.analytics", Constraint: ">=1.0.0 <2.0.0"}}
	inventory.CapabilitiesProvided = []string{"inventory.reserve-stock.v1"}
	inventory.Permissions = []string{"inventory.stock.read"}
	inventory.UISlots = []string{"inventory.dashboard"}

	manifests := []manifestV2{inventory}
	if includeCatalog {
		manifests = append(manifests, healthManifestFixture("omnexa.catalog", "Catalog", "catalog.team", "1.1.0"))
	}
	if includeAnalytics {
		manifests = append(manifests, healthManifestFixture("omnexa.analytics", "Analytics", "analytics.team", "1.0.0"))
	}

	sources := make([]DiscoverySource, 0, len(manifests))
	for _, manifest := range manifests {
		sources = append(sources, DiscoverySource{
			ID:        strings.TrimPrefix(manifest.ID, "omnexa.") + ".fixture",
			Version:   "1.0.0",
			Manifests: [][]byte{manifestV2Bytes(t, manifest)},
		})
	}
	registry, err := Discover(sources)
	if err != nil {
		t.Fatalf("discover module-health fixture: %v", err)
	}

	store := NewMemoryLifecycleStore()
	for _, module := range registry.List() {
		record := LifecycleRecord{ModuleID: module.ID, Version: module.Version, State: LifecycleEnabled, Revision: 1}
		if casErr := store.CompareAndSwap(context.Background(), module.ID, 0, record); casErr != nil {
			t.Fatalf("seed lifecycle %s: %v", module.ID, casErr)
		}
	}

	capabilities, err := BindCapabilityRegistry(registry, store, []CapabilityRegistration{{
		ModuleID:         "omnexa.inventory",
		Owner:            "inventory.team",
		Declaration:      "inventory.reserve-stock.v1",
		AuthorizationRef: "inventory.authorization",
		ScopeRef:         "inventory.scope",
		ContractRef:      "inventory.contract",
	}})
	if err != nil {
		t.Fatalf("bind capability registry: %v", err)
	}
	permissions, err := BindPermissionRegistry(registry, store, []PermissionRegistration{{
		ModuleID:   "omnexa.inventory",
		Owner:      "inventory.team",
		Permission: "inventory.stock.read",
		Capability: "inventory.reserve-stock.v1",
	}})
	if err != nil {
		t.Fatalf("bind permission registry: %v", err)
	}
	ui, err := BindUIContributionRegistry(registry, store, []UIContributionRegistration{{
		ModuleID:           "omnexa.inventory",
		Owner:              "inventory.team",
		ID:                 "inventory.dashboard.widget",
		Slot:               "inventory.dashboard",
		Kind:               UIContributionWidget,
		ContractVersion:    1,
		Permission:         "inventory.stock.read",
		OptionalDependency: "omnexa.analytics",
		Fallback:           UIFallbackDegraded,
	}})
	if err != nil {
		t.Fatalf("bind UI contribution registry: %v", err)
	}

	migrations, err := BindMigrationOwnershipRegistry(registry, []MigrationRegistration{{
		ModuleID:            "omnexa.inventory",
		ModuleVersion:       "1.2.3",
		Owner:               "inventory.team",
		Declaration:         "inventory.initial",
		Version:             1,
		Name:                "initial",
		IntroducedInVersion: "1.2.3",
		TargetOwner:         "inventory.team",
		ChangeClass:         MigrationCompatible,
	}})
	if err != nil {
		t.Fatalf("bind migration ownership registry: %v", err)
	}

	return moduleHealthFixture{
		registry:     registry,
		store:        store,
		migrations:   migrations,
		capabilities: capabilities,
		permissions:  permissions,
		ui:           ui,
	}
}

func (f moduleHealthFixture) reporter(t *testing.T, lifecycle LifecycleStore, source ModuleMigrationHealthSource) *ModuleHealthReporter {
	t.Helper()
	reporter, err := NewModuleHealthReporter(
		f.registry,
		PlatformSnapshot{Version: "1.0.0", Capabilities: []string{}},
		lifecycle,
		f.migrations,
		source,
		f.capabilities,
		f.permissions,
		f.ui,
	)
	if err != nil {
		t.Fatalf("new module health reporter: %v", err)
	}
	return reporter
}

func healthManifestFixture(id, name, owner, version string) manifestV2 {
	manifest := validManifestV2()
	manifest.ID = id
	manifest.Name = name
	manifest.Owner = owner
	manifest.Version = version
	manifest.Dependencies = []DependencyRequirement{}
	manifest.OptionalDependencies = []DependencyRequirement{}
	manifest.PlatformDependencies = []string{}
	manifest.CapabilitiesProvided = []string{}
	manifest.CapabilitiesConsumed = []string{}
	manifest.Permissions = []string{}
	manifest.EventsPublished = []string{}
	manifest.EventsConsumed = []string{}
	manifest.WorkflowTriggers = []string{}
	manifest.WorkflowActions = []string{}
	manifest.UISlots = []string{}
	manifest.Settings = []string{}
	manifest.FeatureFlags = []string{}
	manifest.DataClassification = []string{"INTERNAL"}
	manifest.Migrations = []string{}
	manifest.LifecycleHooks = []string{}
	manifest.HealthChecks = []string{}
	manifest.Security.SecretReferences = []SecretReference{}
	manifest.Security.NetworkDestinations = []string{}
	manifest.Security.ExposedEndpoints = []string{}
	manifest.Security.FileTypes = []string{}
	manifest.Security.PrivilegedOperations = []string{}
	manifest.Security.AICapabilities = []string{}
	return manifest
}

func healthRecordByID(t *testing.T, report ModuleHealthReport, moduleID string) ModuleHealthRecord {
	t.Helper()
	for _, record := range report.Modules {
		if record.ModuleID == moduleID {
			return record
		}
	}
	t.Fatalf("module %q missing from report: %#v", moduleID, report)
	return ModuleHealthRecord{}
}

func containsHealthReason(reasons []string, target string) bool {
	for _, reason := range reasons {
		if reason == target {
			return true
		}
	}
	return false
}

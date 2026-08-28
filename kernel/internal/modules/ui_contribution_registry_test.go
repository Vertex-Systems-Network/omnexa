package modules

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestUIContributionRegistryRetainsV1AndV2ValidatedUISlots(t *testing.T) {
	v1 := validManifest()
	v1.UISlots = []string{"inventory.dashboard", "inventory.settings"}
	v1Snapshot, err := parseValidatedManifest(manifestBytes(t, v1))
	if err != nil {
		t.Fatalf("parse v1 snapshot: %v", err)
	}
	if !reflect.DeepEqual(v1Snapshot.UISlots, v1.UISlots) {
		t.Fatalf("v1 UI slots missing from validated snapshot: %#v", v1Snapshot.UISlots)
	}

	v2 := validManifestV2()
	v2.UISlots = []string{"inventory.dashboard", "inventory.builder"}
	v2Snapshot, err := parseValidatedManifest(manifestV2Bytes(t, v2))
	if err != nil {
		t.Fatalf("parse v2 snapshot: %v", err)
	}
	if !reflect.DeepEqual(v2Snapshot.UISlots, v2.UISlots) {
		t.Fatalf("v2 UI slots missing from validated snapshot: %#v", v2Snapshot.UISlots)
	}

	clone := v2Snapshot.clone()
	clone.UISlots[0] = "mutated.slot"
	if v2Snapshot.UISlots[0] != "inventory.dashboard" {
		t.Fatal("validated snapshot clone leaked UI slot slice mutation")
	}
}

func TestUIContributionRegistryDeterministicMetadataLookup(t *testing.T) {
	ctx := context.Background()
	registry := uiRegistryWithInventoryAndAnalytics(t, "1.2.3")
	store := NewMemoryLifecycleStore()
	seedUILifecycle(t, store, "omnexa.inventory", LifecycleEnabled)
	seedUILifecycle(t, store, "omnexa.analytics", LifecycleEnabled)

	uiRegistry, err := BindUIContributionRegistry(registry, store, []UIContributionRegistration{
		validUIContribution("inventory.stock.widget", UIContributionWidget, "omnexa.analytics"),
		validUIContribution("inventory.stock.page", UIContributionPage, ""),
	})
	if err != nil {
		t.Fatalf("BindUIContributionRegistry() error = %v", err)
	}

	records, err := uiRegistry.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(records) != 2 || records[0].ID != "inventory.stock.page" || records[1].ID != "inventory.stock.widget" {
		t.Fatalf("unexpected deterministic records: %#v", records)
	}
	widget := records[1]
	if widget.ModuleID != "omnexa.inventory" || widget.Owner != "inventory.team" ||
		widget.Slot != "inventory.dashboard" || widget.Kind != UIContributionWidget ||
		widget.ContractVersion != 1 || widget.Permission != "inventory.stock.read" ||
		widget.FeatureFlag != "inventory.batch-reservation" || widget.OptionalDependency != "omnexa.analytics" ||
		widget.Fallback != UIFallbackDegraded || !widget.Available || widget.Degraded {
		t.Fatalf("unexpected widget record: %#v", widget)
	}

	lookup, err := uiRegistry.Lookup(ctx, "omnexa.inventory", "inventory.stock.widget")
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}
	if !reflect.DeepEqual(lookup, widget) {
		t.Fatalf("lookup mismatch: list=%#v lookup=%#v", widget, lookup)
	}
}

func TestUIContributionRegistryRejectsInvalidRegistrationReferences(t *testing.T) {
	registry := uiRegistryWithInventoryOnly(t, true)
	store := NewMemoryLifecycleStore()

	tests := []struct {
		name   string
		mutate func(*UIContributionRegistration)
		code   string
	}{
		{name: "owner mismatch", mutate: func(value *UIContributionRegistration) { value.Owner = "other.team" }, code: "module.ui.owner_mismatch"},
		{name: "foreign identity", mutate: func(value *UIContributionRegistration) { value.ID = "analytics.stock.widget" }, code: "module.ui.registration_invalid"},
		{name: "unknown slot", mutate: func(value *UIContributionRegistration) { value.Slot = "inventory.unknown" }, code: "module.ui.slot_undeclared"},
		{name: "unknown kind", mutate: func(value *UIContributionRegistration) { value.Kind = "script" }, code: "module.ui.registration_invalid"},
		{name: "zero contract", mutate: func(value *UIContributionRegistration) { value.ContractVersion = 0 }, code: "module.ui.registration_invalid"},
		{name: "unknown permission", mutate: func(value *UIContributionRegistration) { value.Permission = "inventory.secret.execute" }, code: "module.ui.permission_undeclared"},
		{name: "unknown feature flag", mutate: func(value *UIContributionRegistration) { value.FeatureFlag = "inventory.hidden-mode" }, code: "module.ui.feature_flag_undeclared"},
		{name: "unknown optional dependency", mutate: func(value *UIContributionRegistration) { value.OptionalDependency = "omnexa.unknown" }, code: "module.ui.optional_dependency_undeclared"},
		{name: "invalid fallback", mutate: func(value *UIContributionRegistration) { value.Fallback = "execute" }, code: "module.ui.registration_invalid"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration := validUIContribution("inventory.stock.widget", UIContributionWidget, "omnexa.analytics")
			test.mutate(&registration)
			_, err := BindUIContributionRegistry(registry, store, []UIContributionRegistration{registration})
			assertUIContributionErrorCode(t, err, test.code)
		})
	}

	duplicate := validUIContribution("inventory.stock.widget", UIContributionWidget, "omnexa.analytics")
	_, err := BindUIContributionRegistry(registry, store, []UIContributionRegistration{duplicate, duplicate})
	assertUIContributionErrorCode(t, err, "module.ui.registration_duplicate")
}

func TestUIContributionRegistryOwnerLifecycleAvailabilityFailsClosed(t *testing.T) {
	ctx := context.Background()
	states := []LifecycleState{
		LifecycleAvailable,
		LifecycleInstalled,
		LifecycleDisabled,
		LifecycleSuspended,
		LifecycleArchived,
		LifecycleDetached,
		LifecycleRecoveryRequired,
		LifecyclePurged,
	}

	for _, state := range states {
		t.Run(string(state), func(t *testing.T) {
			registry := uiRegistryWithInventoryOnly(t, true)
			store := NewMemoryLifecycleStore()
			if state != LifecycleAvailable {
				seedUILifecycle(t, store, "omnexa.inventory", state)
			}
			uiRegistry, err := BindUIContributionRegistry(registry, store, []UIContributionRegistration{
				validUIContribution("inventory.stock.page", UIContributionPage, ""),
			})
			if err != nil {
				t.Fatalf("BindUIContributionRegistry() error = %v", err)
			}
			record, err := uiRegistry.Lookup(ctx, "omnexa.inventory", "inventory.stock.page")
			if err != nil {
				t.Fatalf("Lookup() error = %v", err)
			}
			if record.Available {
				t.Fatalf("lifecycle %s must not report UI contribution available: %#v", state, record)
			}
		})
	}

	registry := uiRegistryWithInventoryOnly(t, true)
	store := NewMemoryLifecycleStore()
	seedUILifecycle(t, store, "omnexa.inventory", LifecycleEnabled)
	uiRegistry, err := BindUIContributionRegistry(registry, store, []UIContributionRegistration{
		validUIContribution("inventory.stock.page", UIContributionPage, ""),
	})
	if err != nil {
		t.Fatalf("BindUIContributionRegistry() error = %v", err)
	}
	record, err := uiRegistry.Lookup(ctx, "omnexa.inventory", "inventory.stock.page")
	if err != nil || !record.Available {
		t.Fatalf("enabled contribution availability = %#v, %v", record, err)
	}
}

func TestUIContributionRegistryOptionalDependencyDegradesOnlyAffectedContribution(t *testing.T) {
	ctx := context.Background()
	registry := uiRegistryWithInventoryOnly(t, true)
	store := NewMemoryLifecycleStore()
	seedUILifecycle(t, store, "omnexa.inventory", LifecycleEnabled)
	uiRegistry, err := BindUIContributionRegistry(registry, store, []UIContributionRegistration{
		validUIContribution("inventory.stock.widget", UIContributionWidget, "omnexa.analytics"),
		validUIContribution("inventory.stock.page", UIContributionPage, ""),
	})
	if err != nil {
		t.Fatalf("BindUIContributionRegistry() error = %v", err)
	}
	records, err := uiRegistry.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if records[0].ID != "inventory.stock.page" || records[0].Degraded || !records[0].Available {
		t.Fatalf("unrelated contribution was corrupted: %#v", records[0])
	}
	if records[1].ID != "inventory.stock.widget" || !records[1].Degraded || records[1].DegradationReason != "dependency_missing" || !records[1].Available {
		t.Fatalf("missing optional dependency was not selectively degraded: %#v", records[1])
	}
}

func TestUIContributionRegistryOptionalDependencyCompatibilityAndLifecycle(t *testing.T) {
	ctx := context.Background()

	t.Run("incompatible version", func(t *testing.T) {
		registry := uiRegistryWithInventoryAndAnalytics(t, "2.0.0")
		store := NewMemoryLifecycleStore()
		seedUILifecycle(t, store, "omnexa.inventory", LifecycleEnabled)
		seedUILifecycle(t, store, "omnexa.analytics", LifecycleEnabled)
		record := lookupUIContribution(t, ctx, registry, store, true)
		if !record.Degraded || record.DegradationReason != "dependency_incompatible" {
			t.Fatalf("unexpected incompatible dependency record: %#v", record)
		}
	})

	t.Run("dependency disabled", func(t *testing.T) {
		registry := uiRegistryWithInventoryAndAnalytics(t, "1.2.3")
		store := NewMemoryLifecycleStore()
		seedUILifecycle(t, store, "omnexa.inventory", LifecycleEnabled)
		seedUILifecycle(t, store, "omnexa.analytics", LifecycleDisabled)
		record := lookupUIContribution(t, ctx, registry, store, true)
		if !record.Degraded || record.DegradationReason != "dependency_unavailable" {
			t.Fatalf("unexpected unavailable dependency record: %#v", record)
		}
	})

	t.Run("dependency enabled", func(t *testing.T) {
		registry := uiRegistryWithInventoryAndAnalytics(t, "1.2.3")
		store := NewMemoryLifecycleStore()
		seedUILifecycle(t, store, "omnexa.inventory", LifecycleEnabled)
		seedUILifecycle(t, store, "omnexa.analytics", LifecycleEnabled)
		record := lookupUIContribution(t, ctx, registry, store, true)
		if record.Degraded || !record.Available {
			t.Fatalf("enabled compatible dependency unexpectedly degraded: %#v", record)
		}
	})
}

func TestUIContributionRegistrySchemaV1OptionalDependencyFailsSafeAsUnresolved(t *testing.T) {
	ctx := context.Background()
	registry := uiRegistryWithInventoryOnly(t, false)
	store := NewMemoryLifecycleStore()
	seedUILifecycle(t, store, "omnexa.inventory", LifecycleEnabled)
	record := lookupUIContribution(t, ctx, registry, store, true)
	if !record.Degraded || record.DegradationReason != "version_contract_missing" || !record.Available {
		t.Fatalf("schema-v1 optional dependency must degrade safely: %#v", record)
	}
}

func TestUIContributionRegistryMetadataSurfaceHasNoExecutionOrAuthorityFields(t *testing.T) {
	for _, value := range []any{UIContributionRegistration{}, UIContributionRecord{}} {
		typ := reflect.TypeOf(value)
		for _, forbidden := range []string{
			"Handler", "Component", "Render", "Invoke", "Execute", "TenantID", "OrganizationID",
			"DatabaseTable", "Secret", "Payload", "HTML", "AuthorizationDecision", "FeatureFlagValue",
		} {
			if _, ok := typ.FieldByName(forbidden); ok {
				t.Fatalf("%s must not expose forbidden authority/execution field %s", typ.Name(), forbidden)
			}
		}
	}
}

func validUIContribution(id string, kind UIContributionKind, optionalDependency string) UIContributionRegistration {
	return UIContributionRegistration{
		ModuleID:           "omnexa.inventory",
		Owner:              "inventory.team",
		ID:                 id,
		Slot:               "inventory.dashboard",
		Kind:               kind,
		ContractVersion:    1,
		Permission:         "inventory.stock.read",
		FeatureFlag:        "inventory.batch-reservation",
		OptionalDependency: optionalDependency,
		Fallback:           UIFallbackDegraded,
	}
}

func uiRegistryWithInventoryOnly(t *testing.T, schemaV2 bool) Registry {
	t.Helper()
	var payload []byte
	if schemaV2 {
		manifest := validManifestV2()
		manifest.Dependencies = nil
		manifest.UISlots = []string{"inventory.dashboard", "inventory.settings"}
		payload = manifestV2Bytes(t, manifest)
	} else {
		manifest := validManifest()
		manifest.Dependencies = nil
		manifest.UISlots = []string{"inventory.dashboard", "inventory.settings"}
		payload = manifestBytes(t, manifest)
	}
	registry, err := Discover([]DiscoverySource{{ID: "repo.inventory", Version: "1.0.0", Manifests: [][]byte{payload}}})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	return registry
}

func uiRegistryWithInventoryAndAnalytics(t *testing.T, analyticsVersion string) Registry {
	t.Helper()
	inventory := validManifestV2()
	inventory.Dependencies = nil
	inventory.UISlots = []string{"inventory.dashboard", "inventory.settings"}

	analytics := validManifestV2()
	analytics.ID = "omnexa.analytics"
	analytics.Name = "Analytics"
	analytics.Version = analyticsVersion
	analytics.Owner = "analytics.team"
	analytics.Dependencies = nil
	analytics.OptionalDependencies = nil
	analytics.UISlots = []string{"analytics.dashboard"}
	analytics.Permissions = []string{"analytics.report.read"}
	analytics.FeatureFlags = []string{"analytics.summary"}

	registry, err := Discover([]DiscoverySource{
		{ID: "repo.analytics", Version: "1.0.0", Manifests: [][]byte{manifestV2Bytes(t, analytics)}},
		{ID: "repo.inventory", Version: "1.0.0", Manifests: [][]byte{manifestV2Bytes(t, inventory)}},
	})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	return registry
}

func seedUILifecycle(t *testing.T, store *MemoryLifecycleStore, moduleID string, state LifecycleState) {
	t.Helper()
	if err := store.CompareAndSwap(context.Background(), moduleID, 0, LifecycleRecord{
		ModuleID: moduleID,
		Version:  "1.2.3",
		State:    state,
		Revision: 1,
	}); err != nil {
		t.Fatalf("seed lifecycle %s: %v", moduleID, err)
	}
}

func lookupUIContribution(
	t *testing.T,
	ctx context.Context,
	registry Registry,
	store *MemoryLifecycleStore,
	withOptionalDependency bool,
) UIContributionRecord {
	t.Helper()
	optionalDependency := ""
	if withOptionalDependency {
		optionalDependency = "omnexa.analytics"
	}
	uiRegistry, err := BindUIContributionRegistry(registry, store, []UIContributionRegistration{
		validUIContribution("inventory.stock.widget", UIContributionWidget, optionalDependency),
	})
	if err != nil {
		t.Fatalf("BindUIContributionRegistry() error = %v", err)
	}
	record, err := uiRegistry.Lookup(ctx, "omnexa.inventory", "inventory.stock.widget")
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}
	return record
}

func assertUIContributionErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	var uiErr *UIContributionError
	if !errors.As(err, &uiErr) {
		t.Fatalf("expected UIContributionError %s, got %T: %v", code, err, err)
	}
	if uiErr.Diagnostic.Code != code {
		t.Fatalf("unexpected UI contribution error: got %#v want code %s", uiErr.Diagnostic, code)
	}
}

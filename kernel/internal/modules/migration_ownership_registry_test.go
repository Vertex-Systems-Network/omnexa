package modules

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestMigrationOwnershipRegistryDeterministicFreshInstallAndUpgradePlan(t *testing.T) {
	registry := migrationOwnershipTestModuleRegistry(t, "1.2.0")
	registrations := migrationOwnershipTestRegistrations()

	bound, err := BindMigrationOwnershipRegistry(registry, []MigrationRegistration{
		registrations[2],
		registrations[0],
		registrations[1],
	})
	if err != nil {
		t.Fatalf("BindMigrationOwnershipRegistry() error = %v", err)
	}

	listed := bound.List()
	assertMigrationVersions(t, listed, []int64{1, 2, 3})
	if listed[1].ID != "omnexa.inventory@1.2.0:inventory.team:2:add_location" {
		t.Fatalf("unexpected deterministic migration ID: %q", listed[1].ID)
	}

	fresh, err := bound.FreshInstallPlan("omnexa.inventory")
	if err != nil {
		t.Fatalf("FreshInstallPlan() error = %v", err)
	}
	assertMigrationVersions(t, fresh, []int64{1, 2, 3})

	upgrade, err := bound.UpgradePlan("omnexa.inventory", "1.0.0")
	if err != nil {
		t.Fatalf("UpgradePlan() error = %v", err)
	}
	assertMigrationVersions(t, upgrade, []int64{2, 3})

	current, err := bound.UpgradePlan("omnexa.inventory", "1.2.0")
	if err != nil {
		t.Fatalf("UpgradePlan(current) error = %v", err)
	}
	if len(current) != 0 {
		t.Fatalf("current-version upgrade plan = %#v, want empty", current)
	}

	lookup, err := bound.Lookup("omnexa.inventory", "inventory.0002.add_location")
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}
	if lookup.Version != 2 || lookup.ChangeClass != MigrationBackfill {
		t.Fatalf("unexpected lookup record: %#v", lookup)
	}

	fresh[0].Version = 999
	freshAgain, err := bound.FreshInstallPlan("omnexa.inventory")
	if err != nil {
		t.Fatal(err)
	}
	if freshAgain[0].Version != 1 {
		t.Fatal("FreshInstallPlan returned mutable internal storage")
	}

	second, err := BindMigrationOwnershipRegistry(registry, registrations)
	if err != nil {
		t.Fatalf("second bind error = %v", err)
	}
	if !reflect.DeepEqual(bound.List(), second.List()) {
		t.Fatalf("registry output changed with input order:\nfirst=%#v\nsecond=%#v", bound.List(), second.List())
	}
}

func TestMigrationOwnershipRegistryRejectsIdentityOwnerOrderAndVersionConflicts(t *testing.T) {
	registry := migrationOwnershipTestModuleRegistry(t, "1.2.0")
	base := migrationOwnershipTestRegistrations()

	tests := []struct {
		name string
		code string
		edit func([]MigrationRegistration) []MigrationRegistration
	}{
		{
			name: "missing module",
			code: "module.migration.module_missing",
			edit: func(values []MigrationRegistration) []MigrationRegistration {
				values[0].ModuleID = "omnexa.missing"
				return values
			},
		},
		{
			name: "module version mismatch",
			code: "module.migration.module_version_mismatch",
			edit: func(values []MigrationRegistration) []MigrationRegistration {
				values[0].ModuleVersion = "1.1.0"
				return values
			},
		},
		{
			name: "owner mismatch",
			code: "module.migration.owner_mismatch",
			edit: func(values []MigrationRegistration) []MigrationRegistration {
				values[0].Owner = "other.team"
				return values
			},
		},
		{
			name: "cross owner target",
			code: "module.migration.target_owner_mismatch",
			edit: func(values []MigrationRegistration) []MigrationRegistration {
				values[0].TargetOwner = "other.team"
				return values
			},
		},
		{
			name: "future introduced version",
			code: "module.migration.introduced_version_future",
			edit: func(values []MigrationRegistration) []MigrationRegistration {
				values[0].IntroducedInVersion = "1.3.0"
				return values
			},
		},
		{
			name: "duplicate declaration",
			code: "module.migration.registration_duplicate",
			edit: func(values []MigrationRegistration) []MigrationRegistration {
				duplicate := values[0]
				duplicate.Version = 4
				duplicate.Name = "duplicate"
				return append(values, duplicate)
			},
		},
		{
			name: "owner ledger order conflict",
			code: "module.migration.order_conflict",
			edit: func(values []MigrationRegistration) []MigrationRegistration {
				conflict := values[0]
				conflict.Declaration = "inventory.0004.conflict"
				conflict.Name = "conflict"
				conflict.Version = values[1].Version
				return append(values, conflict)
			},
		},
		{
			name: "introduced version regresses",
			code: "module.migration.introduced_version_order_invalid",
			edit: func(values []MigrationRegistration) []MigrationRegistration {
				values[1].IntroducedInVersion = "0.9.0"
				return values
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values := append([]MigrationRegistration(nil), base...)
			_, err := BindMigrationOwnershipRegistry(registry, test.edit(values))
			assertMigrationOwnershipCode(t, err, test.code)
		})
	}
}

func TestMigrationOwnershipRegistryRequiresExplicitBackfillAndDestructiveStrategy(t *testing.T) {
	registry := migrationOwnershipTestModuleRegistry(t, "1.2.0")

	for _, class := range []MigrationChangeClass{MigrationBackfill, MigrationDestructive} {
		registration := migrationOwnershipTestRegistrations()[0]
		registration.ChangeClass = class
		registration.StrategyRef = ""
		registration.RecoveryRef = ""

		_, err := BindMigrationOwnershipRegistry(registry, []MigrationRegistration{registration})
		assertMigrationOwnershipCode(t, err, "module.migration.strategy_required")
	}

	invalid := migrationOwnershipTestRegistrations()[1]
	invalid.StrategyRef = "contains spaces"
	_, err := BindMigrationOwnershipRegistry(registry, []MigrationRegistration{invalid})
	assertMigrationOwnershipCode(t, err, "module.migration.strategy_ref_invalid")

	unknown := migrationOwnershipTestRegistrations()[0]
	unknown.ChangeClass = MigrationChangeClass("manual")
	_, err = BindMigrationOwnershipRegistry(registry, []MigrationRegistration{unknown})
	assertMigrationOwnershipCode(t, err, "module.migration.change_class_invalid")
}

func TestMigrationOwnershipRegistryUpgradePlanningFailsClosed(t *testing.T) {
	registry := migrationOwnershipTestModuleRegistry(t, "1.2.0")
	bound, err := BindMigrationOwnershipRegistry(registry, migrationOwnershipTestRegistrations())
	if err != nil {
		t.Fatal(err)
	}

	_, err = bound.UpgradePlan("omnexa.inventory", "not-semver")
	assertMigrationOwnershipCode(t, err, "module.migration.upgrade_version_invalid")

	_, err = bound.UpgradePlan("omnexa.inventory", "2.0.0")
	assertMigrationOwnershipCode(t, err, "module.migration.upgrade_future_source")

	_, err = bound.UpgradePlan("omnexa.missing", "1.0.0")
	assertMigrationOwnershipCode(t, err, "module.migration.module_missing")
}

func TestMigrationOwnershipRegistryMetadataSurfaceHasNoExecutionOrRawScopeFields(t *testing.T) {
	forbidden := map[string]struct{}{
		"SQL":            {},
		"Path":           {},
		"Handler":        {},
		"Callback":       {},
		"TenantID":       {},
		"OrganizationID": {},
		"Secret":         {},
	}

	for _, value := range []any{MigrationRegistration{}, MigrationRecord{}} {
		typeOf := reflect.TypeOf(value)
		for index := 0; index < typeOf.NumField(); index++ {
			field := typeOf.Field(index)
			if _, blocked := forbidden[field.Name]; blocked {
				t.Fatalf("%s exposes forbidden field %s", typeOf.Name(), field.Name)
			}
			if field.Type.Kind() == reflect.Func {
				t.Fatalf("%s exposes executable field %s", typeOf.Name(), field.Name)
			}
		}
	}
}

func migrationOwnershipTestModuleRegistry(t *testing.T, version string) Registry {
	t.Helper()
	manifest := validManifest()
	manifest.ID = "omnexa.inventory"
	manifest.Owner = "inventory.team"
	manifest.Version = version

	payload, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	registry, err := Discover([]DiscoverySource{{
		ID:        "test.source",
		Version:   "1.0.0",
		Manifests: [][]byte{payload},
	}})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	return registry
}

func migrationOwnershipTestRegistrations() []MigrationRegistration {
	return []MigrationRegistration{
		{
			ModuleID:            "omnexa.inventory",
			ModuleVersion:       "1.2.0",
			Owner:               "inventory.team",
			Declaration:         "inventory.0001.initial",
			Version:             1,
			Name:                "initial",
			IntroducedInVersion: "1.0.0",
			TargetOwner:         "inventory.team",
			ChangeClass:         MigrationCompatible,
		},
		{
			ModuleID:            "omnexa.inventory",
			ModuleVersion:       "1.2.0",
			Owner:               "inventory.team",
			Declaration:         "inventory.0002.add_location",
			Version:             2,
			Name:                "add_location",
			IntroducedInVersion: "1.1.0",
			TargetOwner:         "inventory.team",
			ChangeClass:         MigrationBackfill,
			StrategyRef:         "inventory.backfill.location",
			RecoveryRef:         "inventory.recovery.location",
		},
		{
			ModuleID:            "omnexa.inventory",
			ModuleVersion:       "1.2.0",
			Owner:               "inventory.team",
			Declaration:         "inventory.0003.retire_legacy",
			Version:             3,
			Name:                "retire_legacy",
			IntroducedInVersion: "1.2.0",
			TargetOwner:         "inventory.team",
			ChangeClass:         MigrationDestructive,
			StrategyRef:         "inventory.destructive.retire_legacy",
			RecoveryRef:         "inventory.recovery.forward_fix",
		},
	}
}

func assertMigrationOwnershipCode(t *testing.T, err error, code string) {
	t.Helper()
	var ownershipErr *MigrationOwnershipError
	if !errors.As(err, &ownershipErr) {
		t.Fatalf("expected MigrationOwnershipError %q, got %T: %v", code, err, err)
	}
	if ownershipErr.Diagnostic.Code != code {
		t.Fatalf("error code = %q, want %q", ownershipErr.Diagnostic.Code, code)
	}
}

func assertMigrationVersions(t *testing.T, records []MigrationRecord, want []int64) {
	t.Helper()
	if len(records) != len(want) {
		t.Fatalf("migration count = %d, want %d: %#v", len(records), len(want), records)
	}
	for index, version := range want {
		if records[index].Version != version {
			t.Fatalf("record[%d].Version = %d, want %d", index, records[index].Version, version)
		}
	}
}

package modules

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/authorization"
)

func TestPermissionRegistryRetainsV1AndV2ValidatedDeclarations(t *testing.T) {
	v1 := validManifest()
	v1.Permissions = []string{"inventory.stock.read", "inventory.stock.manage"}
	v1Snapshot, err := parseValidatedManifest(manifestBytes(t, v1))
	if err != nil {
		t.Fatalf("parse v1 snapshot: %v", err)
	}
	if !reflect.DeepEqual(v1Snapshot.Permissions, v1.Permissions) {
		t.Fatalf("v1 permission declarations missing: %#v", v1Snapshot.Permissions)
	}

	v2 := validManifestV2()
	v2.Permissions = []string{"inventory.stock.read", "inventory.stock.manage"}
	v2Snapshot, err := parseValidatedManifest(manifestV2Bytes(t, v2))
	if err != nil {
		t.Fatalf("parse v2 snapshot: %v", err)
	}
	if !reflect.DeepEqual(v2Snapshot.Permissions, v2.Permissions) {
		t.Fatalf("v2 permission declarations missing: %#v", v2Snapshot.Permissions)
	}

	clone := v2Snapshot.clone()
	clone.Permissions[0] = "mutated.permission.value"
	if v2Snapshot.Permissions[0] != "inventory.stock.read" {
		t.Fatal("validated snapshot clone leaked permission slice mutation")
	}
}

func TestPermissionRegistryDeterministicRegistrationAndCapabilityAssociation(t *testing.T) {
	ctx := context.Background()
	registry := permissionModuleRegistry(t, "inventory.stock.read", "inventory.stock.manage")
	store := enabledCapabilityLifecycleStore(t, "omnexa.inventory")

	permissions, err := BindPermissionRegistry(registry, store, []PermissionRegistration{
		{
			ModuleID: "omnexa.inventory", Owner: "inventory.team",
			Permission: "inventory.stock.manage", Capability: "inventory.reserve-stock.v1",
		},
		{
			ModuleID: "omnexa.inventory", Owner: "inventory.team",
			Permission: "inventory.stock.read",
		},
	})
	if err != nil {
		t.Fatalf("BindPermissionRegistry() error = %v", err)
	}

	records, err := permissions.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(records) != 2 || records[0].Permission != "inventory.stock.manage" || records[1].Permission != "inventory.stock.read" {
		t.Fatalf("unexpected deterministic permission ordering: %#v", records)
	}
	for _, record := range records {
		if record.ModuleID != "omnexa.inventory" || record.Owner != "inventory.team" || !record.Available || record.LifecycleState != LifecycleEnabled {
			t.Fatalf("unexpected permission record: %#v", record)
		}
	}
	if records[0].Capability != "inventory.reserve-stock.v1" {
		t.Fatalf("capability association not retained: %#v", records[0])
	}

	available, err := permissions.PermissionAvailable(ctx, authorization.PermissionID("inventory.stock.read"))
	if err != nil || !available {
		t.Fatalf("PermissionAvailable() = %v, %v; want true, nil", available, err)
	}
	unknown, err := permissions.PermissionAvailable(ctx, authorization.PermissionID("inventory.stock.delete"))
	if err != nil || unknown {
		t.Fatalf("unknown PermissionAvailable() = %v, %v; want false, nil", unknown, err)
	}
}

func TestPermissionRegistryRejectsReservedNamespaceCollisionAndOwnershipMismatch(t *testing.T) {
	t.Run("reserved namespace", func(t *testing.T) {
		manifest := validManifest()
		manifest.Permissions = []string{"authorization.role.read"}
		registry, err := Discover([]DiscoverySource{{ID: "repo.inventory", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, manifest)}}})
		if err != nil {
			t.Fatalf("Discover() error = %v", err)
		}
		_, err = BindPermissionRegistry(registry, NewMemoryLifecycleStore(), []PermissionRegistration{{
			ModuleID: manifest.ID, Owner: manifest.Owner, Permission: manifest.Permissions[0],
		}})
		assertPermissionErrorCode(t, err, "module.permission.namespace_invalid")
	})

	t.Run("cross module declaration collision", func(t *testing.T) {
		first := validManifest()
		first.Permissions = []string{"inventory.stock.read"}
		second := validManifest()
		second.ID = "omnexa.inventory_reports"
		second.Name = "Inventory Reports"
		second.Owner = "reports.team"
		second.Dependencies = nil
		second.OptionalDependencies = nil
		second.Permissions = []string{"inventory.stock.read"}
		registry, err := Discover([]DiscoverySource{
			{ID: "repo.inventory", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, first)}},
			{ID: "repo.reports", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, second)}},
		})
		if err != nil {
			t.Fatalf("Discover() error = %v", err)
		}
		_, err = BindPermissionRegistry(registry, NewMemoryLifecycleStore(), nil)
		assertPermissionErrorCode(t, err, "module.permission.declaration_collision")
	})

	t.Run("owner mismatch", func(t *testing.T) {
		registry := permissionModuleRegistry(t, "inventory.stock.read")
		_, err := BindPermissionRegistry(registry, NewMemoryLifecycleStore(), []PermissionRegistration{{
			ModuleID: "omnexa.inventory", Owner: "other.team", Permission: "inventory.stock.read",
		}})
		assertPermissionErrorCode(t, err, "module.permission.owner_mismatch")
	})

	t.Run("duplicate explicit registration", func(t *testing.T) {
		registry := permissionModuleRegistry(t, "inventory.stock.read")
		registration := PermissionRegistration{
			ModuleID: "omnexa.inventory", Owner: "inventory.team", Permission: "inventory.stock.read",
		}
		_, err := BindPermissionRegistry(registry, NewMemoryLifecycleStore(), []PermissionRegistration{registration, registration})
		assertPermissionErrorCode(t, err, "module.permission.registration_duplicate")
	})
}

func TestPermissionRegistryCapabilityAssociationMustBeOwnedDeclaration(t *testing.T) {
	registry := permissionModuleRegistry(t, "inventory.stock.manage")
	_, err := BindPermissionRegistry(registry, NewMemoryLifecycleStore(), []PermissionRegistration{{
		ModuleID: "omnexa.inventory", Owner: "inventory.team", Permission: "inventory.stock.manage",
		Capability: "catalog.product.read.v1",
	}})
	assertPermissionErrorCode(t, err, "module.permission.capability_undeclared")
}

func TestPermissionRegistryLifecycleAvailabilityIsEnabledOnlyAndHistorySafe(t *testing.T) {
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
			registry := permissionModuleRegistry(t, "inventory.stock.read")
			store := NewMemoryLifecycleStore()
			if state != LifecycleAvailable {
				if err := store.CompareAndSwap(ctx, "omnexa.inventory", 0, LifecycleRecord{
					ModuleID: "omnexa.inventory", Version: "1.2.3", State: state, Revision: 1,
				}); err != nil {
					t.Fatalf("seed lifecycle: %v", err)
				}
			}
			permissions, err := BindPermissionRegistry(registry, store, []PermissionRegistration{{
				ModuleID: "omnexa.inventory", Owner: "inventory.team", Permission: "inventory.stock.read",
			}})
			if err != nil {
				t.Fatalf("BindPermissionRegistry() error = %v", err)
			}
			record, err := permissions.Lookup(ctx, authorization.PermissionID("inventory.stock.read"))
			if err != nil {
				t.Fatalf("Lookup() error = %v", err)
			}
			if record.Available {
				t.Fatalf("lifecycle %s must not authorize availability: %#v", state, record)
			}
			if record.Permission != "inventory.stock.read" || record.ModuleID != "omnexa.inventory" || record.Owner != "inventory.team" {
				t.Fatalf("historical permission identity lost: %#v", record)
			}
		})
	}
}

func TestPermissionRegistrySynchronizeCatalogIsMetadataOnly(t *testing.T) {
	ctx := context.Background()
	registry := permissionModuleRegistry(t, "inventory.stock.read")
	store := enabledCapabilityLifecycleStore(t, "omnexa.inventory")
	permissions, err := BindPermissionRegistry(registry, store, []PermissionRegistration{{
		ModuleID: "omnexa.inventory", Owner: "inventory.team", Permission: "inventory.stock.read",
	}})
	if err != nil {
		t.Fatalf("BindPermissionRegistry() error = %v", err)
	}
	catalog := &recordingPermissionCatalog{}
	observed := time.Unix(1_700_001_000, 0).UTC()
	if err := permissions.SynchronizeCatalog(ctx, catalog, observed); err != nil {
		t.Fatalf("SynchronizeCatalog() error = %v", err)
	}
	if !catalog.observedAt.Equal(observed) || len(catalog.definitions) != 1 {
		t.Fatalf("unexpected catalog sync: %#v at %v", catalog.definitions, catalog.observedAt)
	}
	definition := catalog.definitions[0]
	if definition.Permission != "inventory.stock.read" || definition.ModuleID != "omnexa.inventory" || definition.Owner != "inventory.team" || !definition.Available {
		t.Fatalf("unexpected catalog definition: %#v", definition)
	}
}

func permissionModuleRegistry(t *testing.T, permissions ...string) Registry {
	t.Helper()
	manifest := validManifest()
	manifest.Dependencies = nil
	manifest.OptionalDependencies = nil
	manifest.Permissions = append([]string(nil), permissions...)
	registry, err := Discover([]DiscoverySource{{
		ID: "repo.inventory", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, manifest)},
	}})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	return registry
}

func assertPermissionErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	var permissionErr *PermissionError
	if !errors.As(err, &permissionErr) {
		t.Fatalf("expected PermissionError %s, got %T: %v", code, err, err)
	}
	if permissionErr.Diagnostic.Code != code {
		t.Fatalf("unexpected permission error: got %#v want code %s", permissionErr.Diagnostic, code)
	}
}

type recordingPermissionCatalog struct {
	definitions []authorization.ModulePermissionDefinition
	observedAt  time.Time
}

func (catalog *recordingPermissionCatalog) ReconcileModulePermissions(
	_ context.Context,
	definitions []authorization.ModulePermissionDefinition,
	observedAt time.Time,
) error {
	catalog.definitions = append([]authorization.ModulePermissionDefinition(nil), definitions...)
	catalog.observedAt = observedAt
	return nil
}

package modules

import "testing"

func TestRegistryManifestSnapshotReturnsIndependentClone(t *testing.T) {
	catalog := simpleManifestV2("omnexa.catalog", "2.0.0")
	inventory := simpleManifestV2("omnexa.inventory", "1.0.0")
	inventory.Dependencies = []DependencyRequirement{{ID: catalog.ID, Constraint: "=2.0.0"}}
	registry := discoverV2Registry(t, inventory, catalog)

	snapshot, ok := registry.manifestSnapshot(inventory.ID)
	if !ok || len(snapshot.RequiredDependencies) != 1 {
		t.Fatalf("expected bound inventory snapshot, got %#v", snapshot)
	}
	snapshot.RequiredDependencies[0].ID = "omnexa.mutated"
	snapshot.RequiredDependencies[0].Constraint = "=999.0.0"
	snapshot.PlatformDependencies = append(snapshot.PlatformDependencies, "kernel.mutated")

	fresh, ok := registry.manifestSnapshot(inventory.ID)
	if !ok {
		t.Fatal("expected fresh registry-bound snapshot")
	}
	if len(fresh.RequiredDependencies) != 1 || fresh.RequiredDependencies[0].ID != catalog.ID || fresh.RequiredDependencies[0].Constraint != "=2.0.0" {
		t.Fatalf("caller mutation changed registry evidence: %#v", fresh)
	}
	if len(fresh.PlatformDependencies) != 0 {
		t.Fatalf("caller mutation changed registry platform evidence: %#v", fresh.PlatformDependencies)
	}
}

func TestResolveDependenciesFailsClosedWhenRegistrySnapshotMismatchesRecord(t *testing.T) {
	inventory := simpleManifestV2("omnexa.inventory", "1.0.0")
	registry := discoverV2Registry(t, inventory)

	snapshot := registry.snapshots[inventory.ID]
	snapshot.Version = "9.9.9"
	registry.snapshots[inventory.ID] = snapshot

	_, err := ResolveDependencies(registry, PlatformSnapshot{Version: "1.0.0"}, nil)
	assertResolutionCode(t, resolutionDiagnostics(t, err), "resolver.registry.snapshot_mismatch")
}

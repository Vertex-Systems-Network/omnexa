package modules

import (
	"errors"
	"reflect"
	"testing"
)

func simpleManifestV2(id, version string) manifestV2 {
	manifest := validManifestV2()
	manifest.ID = id
	manifest.Name = "Test Module"
	manifest.Version = version
	manifest.Owner = "test.team"
	manifest.RequiredPlatformVersion = ">=1.0.0"
	manifest.Dependencies = nil
	manifest.OptionalDependencies = nil
	manifest.PlatformDependencies = nil
	return manifest
}

func discoverV2Registry(t *testing.T, manifests ...manifestV2) Registry {
	t.Helper()
	payloads := make([][]byte, 0, len(manifests))
	for _, manifest := range manifests {
		payloads = append(payloads, manifestV2Bytes(t, manifest))
	}
	registry, err := Discover([]DiscoverySource{{
		ID:        "repo.test",
		Version:   "1.0.0",
		Manifests: payloads,
	}})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	return registry
}

func resolutionDiagnostics(t *testing.T, err error) []ResolutionDiagnostic {
	t.Helper()
	var resolutionErr *ResolutionErrors
	if !errors.As(err, &resolutionErr) {
		t.Fatalf("expected ResolutionErrors, got %T: %v", err, err)
	}
	return resolutionErr.Diagnostics()
}

func assertResolutionCode(t *testing.T, diagnostics []ResolutionDiagnostic, code string) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			return
		}
	}
	t.Fatalf("missing resolution code %s in %#v", code, diagnostics)
}

func TestResolveDependenciesProducesStableRequiredOrder(t *testing.T) {
	catalog := simpleManifestV2("omnexa.catalog", "2.1.0")
	inventory := simpleManifestV2("omnexa.inventory", "1.5.0")
	inventory.Dependencies = []DependencyRequirement{{ID: catalog.ID, Constraint: ">=2.0.0 <3.0.0"}}
	commerce := simpleManifestV2("omnexa.commerce", "1.0.0")
	commerce.Dependencies = []DependencyRequirement{{ID: inventory.ID, Constraint: ">=1.0.0 <2.0.0"}}

	forward := discoverV2Registry(t, commerce, catalog, inventory)
	reverse := discoverV2Registry(t, inventory, catalog, commerce)
	platform := PlatformSnapshot{Version: "1.5.0"}

	first, err := ResolveDependencies(forward, platform, nil)
	if err != nil {
		t.Fatalf("ResolveDependencies(forward) error = %v", err)
	}
	second, err := ResolveDependencies(reverse, platform, nil)
	if err != nil {
		t.Fatalf("ResolveDependencies(reverse) error = %v", err)
	}
	want := []string{"omnexa.catalog", "omnexa.inventory", "omnexa.commerce"}
	if !reflect.DeepEqual(first.Order, want) || !reflect.DeepEqual(second.Order, want) {
		t.Fatalf("unexpected deterministic order: first=%#v second=%#v want=%#v", first.Order, second.Order, want)
	}
}

func TestResolveDependenciesRejectsMissingAndIncompatibleRequiredDependency(t *testing.T) {
	inventory := simpleManifestV2("omnexa.inventory", "1.0.0")
	inventory.Dependencies = []DependencyRequirement{{ID: "omnexa.catalog", Constraint: ">=2.0.0 <3.0.0"}}

	_, err := ResolveDependencies(discoverV2Registry(t, inventory), PlatformSnapshot{Version: "1.0.0"}, nil)
	assertResolutionCode(t, resolutionDiagnostics(t, err), "resolver.dependency.required_missing")

	catalog := simpleManifestV2("omnexa.catalog", "3.0.0")
	_, err = ResolveDependencies(discoverV2Registry(t, inventory, catalog), PlatformSnapshot{Version: "1.0.0"}, nil)
	assertResolutionCode(t, resolutionDiagnostics(t, err), "resolver.dependency.required_incompatible")
}

func TestResolveDependenciesDegradesOptionalDependencyWithoutGlobalFailure(t *testing.T) {
	inventory := simpleManifestV2("omnexa.inventory", "1.0.0")
	inventory.OptionalDependencies = []DependencyRequirement{{ID: "omnexa.analytics", Constraint: ">=1.0.0 <2.0.0"}}

	missing, err := ResolveDependencies(discoverV2Registry(t, inventory), PlatformSnapshot{Version: "1.0.0"}, nil)
	if err != nil {
		t.Fatalf("optional absence must not globally fail: %v", err)
	}
	if len(missing.Optional) != 1 || missing.Optional[0].State != optionalMissing || missing.Optional[0].Reason != "dependency_missing" {
		t.Fatalf("unexpected missing optional metadata: %#v", missing.Optional)
	}

	analytics := simpleManifestV2("omnexa.analytics", "2.0.0")
	incompatible, err := ResolveDependencies(discoverV2Registry(t, inventory, analytics), PlatformSnapshot{Version: "1.0.0"}, nil)
	if err != nil {
		t.Fatalf("optional incompatibility must not globally fail: %v", err)
	}
	if len(incompatible.Optional) != 1 || incompatible.Optional[0].State != optionalIncompatible || incompatible.Optional[0].Reason != "version_incompatible" {
		t.Fatalf("unexpected incompatible optional metadata: %#v", incompatible.Optional)
	}
}

func TestResolveDependenciesRejectsRequiredCycleButIgnoresOptionalCycleForGlobalOrder(t *testing.T) {
	first := simpleManifestV2("omnexa.alpha", "1.0.0")
	second := simpleManifestV2("omnexa.beta", "1.0.0")
	first.Dependencies = []DependencyRequirement{{ID: second.ID, Constraint: "=1.0.0"}}
	second.Dependencies = []DependencyRequirement{{ID: first.ID, Constraint: "=1.0.0"}}

	_, err := ResolveDependencies(discoverV2Registry(t, first, second), PlatformSnapshot{Version: "1.0.0"}, nil)
	assertResolutionCode(t, resolutionDiagnostics(t, err), "resolver.graph.required_cycle")

	first.Dependencies = nil
	second.Dependencies = nil
	first.OptionalDependencies = []DependencyRequirement{{ID: second.ID, Constraint: "=1.0.0"}}
	second.OptionalDependencies = []DependencyRequirement{{ID: first.ID, Constraint: "=1.0.0"}}
	resolved, err := ResolveDependencies(discoverV2Registry(t, second, first), PlatformSnapshot{Version: "1.0.0"}, nil)
	if err != nil {
		t.Fatalf("optional cycle must not globally fail required graph: %v", err)
	}
	want := []string{"omnexa.alpha", "omnexa.beta"}
	if !reflect.DeepEqual(resolved.Order, want) {
		t.Fatalf("unexpected optional-cycle order: got %#v want %#v", resolved.Order, want)
	}
}

func TestResolveDependenciesAppliesSchemaV1MigrationRules(t *testing.T) {
	inventory := validManifest()
	inventory.Dependencies = []string{"omnexa.catalog"}
	inventory.OptionalDependencies = nil
	inventory.PlatformDependencies = nil
	catalog := validManifest()
	catalog.ID = "omnexa.catalog"
	catalog.Name = "Catalog"
	catalog.Owner = "catalog.team"
	catalog.Dependencies = nil
	catalog.OptionalDependencies = nil
	catalog.PlatformDependencies = nil

	registry, err := Discover([]DiscoverySource{{
		ID:        "repo.v1",
		Version:   "1.0.0",
		Manifests: [][]byte{manifestBytes(t, inventory), manifestBytes(t, catalog)},
	}})
	if err != nil {
		t.Fatalf("Discover(v1) error = %v", err)
	}
	_, err = ResolveDependencies(registry, PlatformSnapshot{Version: "1.0.0"}, nil)
	assertResolutionCode(t, resolutionDiagnostics(t, err), "resolver.dependency.version_contract_missing")

	inventory.Dependencies = nil
	inventory.OptionalDependencies = []string{"omnexa.catalog"}
	registry, err = Discover([]DiscoverySource{{
		ID:        "repo.v1",
		Version:   "1.0.0",
		Manifests: [][]byte{manifestBytes(t, inventory), manifestBytes(t, catalog)},
	}})
	if err != nil {
		t.Fatalf("Discover(v1 optional) error = %v", err)
	}
	resolved, err := ResolveDependencies(registry, PlatformSnapshot{Version: "1.0.0"}, nil)
	if err != nil {
		t.Fatalf("v1 optional dependency must degrade without global failure: %v", err)
	}
	if len(resolved.Optional) != 1 || resolved.Optional[0].State != optionalUnresolved || resolved.Optional[0].Reason != "version_contract_missing" {
		t.Fatalf("unexpected v1 optional metadata: %#v", resolved.Optional)
	}
}

func TestRegistryBoundSnapshotIsIndependentFromSourcePayloadMutation(t *testing.T) {
	catalog := simpleManifestV2("omnexa.catalog", "2.0.0")
	inventory := simpleManifestV2("omnexa.inventory", "1.0.0")
	inventory.Dependencies = []DependencyRequirement{{ID: catalog.ID, Constraint: "=2.0.0"}}
	inventoryPayload := manifestV2Bytes(t, inventory)
	catalogPayload := manifestV2Bytes(t, catalog)

	registry, err := Discover([]DiscoverySource{{
		ID:        "repo.bound",
		Version:   "1.0.0",
		Manifests: [][]byte{inventoryPayload, catalogPayload},
	}})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	for index := range inventoryPayload {
		inventoryPayload[index] = 'x'
	}
	resolved, err := ResolveDependencies(registry, PlatformSnapshot{Version: "1.0.0"}, nil)
	if err != nil {
		t.Fatalf("resolver drifted after source payload mutation: %v", err)
	}
	want := []string{"omnexa.catalog", "omnexa.inventory"}
	if !reflect.DeepEqual(resolved.Order, want) {
		t.Fatalf("unexpected bound-snapshot order: %#v", resolved.Order)
	}
}

func TestResolveDependenciesFailsClosedWhenRegistrySnapshotIsMissing(t *testing.T) {
	record := RegistryRecord{ID: "omnexa.inventory", Version: "1.0.0", Owner: "test.team", SourceID: "repo.test", SourceVersion: "1.0.0"}
	registry := Registry{
		records: []RegistryRecord{record},
		byID:    map[string]RegistryRecord{record.ID: record},
	}
	_, err := ResolveDependencies(registry, PlatformSnapshot{Version: "1.0.0"}, nil)
	assertResolutionCode(t, resolutionDiagnostics(t, err), "resolver.registry.snapshot_missing")
}

func TestResolveDependenciesEnforcesPlatformVersionAndCapabilities(t *testing.T) {
	inventory := simpleManifestV2("omnexa.inventory", "1.0.0")
	inventory.RequiredPlatformVersion = ">=2.0.0"
	inventory.PlatformDependencies = []string{"kernel.identity"}
	_, err := ResolveDependencies(discoverV2Registry(t, inventory), PlatformSnapshot{Version: "1.5.0"}, nil)
	diagnostics := resolutionDiagnostics(t, err)
	assertResolutionCode(t, diagnostics, "resolver.platform.version_incompatible")
	assertResolutionCode(t, diagnostics, "resolver.platform.dependency_missing")
}

func TestResolveDependenciesRejectsUndeclaredPrivateAndKernelToBusinessObservations(t *testing.T) {
	alpha := simpleManifestV2("omnexa.alpha", "1.0.0")
	beta := simpleManifestV2("omnexa.beta", "1.0.0")
	alpha.OptionalDependencies = []DependencyRequirement{{ID: beta.ID, Constraint: "=1.0.0"}}
	registry := discoverV2Registry(t, alpha, beta)

	_, err := ResolveDependencies(registry, PlatformSnapshot{Version: "1.0.0"}, []DependencyObservation{
		{ConsumerID: alpha.ID, ProviderID: "omnexa.gamma"},
		{ConsumerID: alpha.ID, ProviderID: beta.ID, Private: true},
		{ConsumerID: "kernel.modules", ProviderID: alpha.ID},
	})
	diagnostics := resolutionDiagnostics(t, err)
	assertResolutionCode(t, diagnostics, "resolver.dependency.undeclared")
	assertResolutionCode(t, diagnostics, "resolver.dependency.private_forbidden")
	assertResolutionCode(t, diagnostics, "resolver.dependency.forbidden")
}

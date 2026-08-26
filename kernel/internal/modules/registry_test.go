package modules

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestDiscoverEmptySourcesReturnsEmptyRegistry(t *testing.T) {
	registry, err := Discover(nil)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if got := registry.List(); len(got) != 0 {
		t.Fatalf("expected empty registry, got %#v", got)
	}
	if _, ok := registry.Lookup("omnexa.inventory"); ok {
		t.Fatal("empty registry unexpectedly resolved a module")
	}
}

func TestDiscoverIsDeterministicAcrossEnumerationOrder(t *testing.T) {
	inventory := validManifest()
	catalog := validManifest()
	catalog.ID = "omnexa.catalog"
	catalog.Name = "Catalog"
	catalog.Version = "2.0.0"
	catalog.Owner = "catalog.team"

	forward := []DiscoverySource{
		{ID: "repo.alpha", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, inventory)}},
		{ID: "repo.beta", Version: "1.1.0", Manifests: [][]byte{manifestBytes(t, catalog)}},
	}
	reverse := []DiscoverySource{
		{ID: "repo.beta", Version: "1.1.0", Manifests: [][]byte{manifestBytes(t, catalog)}},
		{ID: "repo.alpha", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, inventory)}},
	}

	first, err := Discover(forward)
	if err != nil {
		t.Fatalf("Discover(forward) error = %v", err)
	}
	second, err := Discover(reverse)
	if err != nil {
		t.Fatalf("Discover(reverse) error = %v", err)
	}
	if !reflect.DeepEqual(first.List(), second.List()) {
		t.Fatalf("registry ordering changed with enumeration order:\nfirst=%#v\nsecond=%#v", first.List(), second.List())
	}
	if got := first.List(); len(got) != 2 || got[0].ID != "omnexa.catalog" || got[1].ID != "omnexa.inventory" {
		t.Fatalf("unexpected stable ordering: %#v", got)
	}
}

func TestDiscoverRejectsDuplicateModuleIdentity(t *testing.T) {
	manifest := validManifest()
	_, err := Discover([]DiscoverySource{
		{ID: "repo.alpha", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, manifest)}},
		{ID: "repo.beta", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, manifest)}},
	})
	diagnostics := discoveryDiagnostics(t, err)
	assertDiscoveryCode(t, diagnostics, "discovery.module.duplicate")
}

func TestDiscoverRejectsConflictingModuleVersions(t *testing.T) {
	first := validManifest()
	second := validManifest()
	second.Version = "1.2.4"

	_, err := Discover([]DiscoverySource{
		{ID: "repo.alpha", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, first)}},
		{ID: "repo.beta", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, second)}},
	})
	diagnostics := discoveryDiagnostics(t, err)
	assertDiscoveryCode(t, diagnostics, "discovery.module.version_conflict")
}

func TestDiscoverRejectsMalformedOrUnvalidatedManifest(t *testing.T) {
	manifest := validManifest()
	manifest.ID = "INVALID MODULE ID"

	_, err := Discover([]DiscoverySource{{
		ID:        "repo.alpha",
		Version:   "1.0.0",
		Manifests: [][]byte{manifestBytes(t, manifest), []byte(`{"not":"a manifest"}`)},
	}})
	diagnostics := discoveryDiagnostics(t, err)
	assertDiscoveryCode(t, diagnostics, "discovery.manifest.invalid")
	for _, diagnostic := range diagnostics {
		if diagnostic.SourceID != "repo.alpha" {
			t.Fatalf("unexpected diagnostic source identity: %#v", diagnostic)
		}
		if diagnostic.ModuleID != "" {
			t.Fatalf("malformed manifest diagnostic leaked an unvalidated module identity: %#v", diagnostic)
		}
	}
}

func TestDiscoverRejectsInvalidExplicitSourceWithoutLeakingRawIdentity(t *testing.T) {
	_, err := Discover([]DiscoverySource{{
		ID:        "C:\\secret\\modules",
		Version:   "latest",
		Manifests: [][]byte{manifestBytes(t, validManifest())},
	}})
	diagnostics := discoveryDiagnostics(t, err)
	assertDiscoveryCode(t, diagnostics, "discovery.source.invalid")
	if len(diagnostics) != 1 || diagnostics[0].SourceID != "" {
		t.Fatalf("invalid source identity must be redacted: %#v", diagnostics)
	}
}

func TestRegistryContainsAvailabilityMetadataOnly(t *testing.T) {
	registry, err := Discover([]DiscoverySource{{
		ID:        "repo.alpha",
		Version:   "1.0.0",
		Manifests: [][]byte{manifestBytes(t, validManifest())},
	}})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	record, ok := registry.Lookup("omnexa.inventory")
	if !ok {
		t.Fatal("expected discovered module")
	}
	if record.Owner != "inventory.team" || record.SourceID != "repo.alpha" || record.SourceVersion != "1.0.0" {
		t.Fatalf("unexpected registry record: %#v", record)
	}

	typeOfRecord := reflect.TypeOf(record)
	for _, forbidden := range []string{"Installed", "Enabled", "Permissions", "Capabilities"} {
		if _, ok := typeOfRecord.FieldByName(forbidden); ok {
			t.Fatalf("registry record must not contain lifecycle/authority field %s", forbidden)
		}
	}
}

func manifestBytes(t *testing.T, manifest Manifest) []byte {
	t.Helper()
	payload, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	return payload
}

func discoveryDiagnostics(t *testing.T, err error) []DiscoveryDiagnostic {
	t.Helper()
	var discoveryErr *DiscoveryErrors
	if !errors.As(err, &discoveryErr) {
		t.Fatalf("expected DiscoveryErrors, got %T: %v", err, err)
	}
	return discoveryErr.Diagnostics()
}

func assertDiscoveryCode(t *testing.T, diagnostics []DiscoveryDiagnostic, code string) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			return
		}
	}
	t.Fatalf("missing discovery code %s in %#v", code, diagnostics)
}

package modules

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestPackageTrustProfileRetainsValidatedV1AndV2Metadata(t *testing.T) {
	legacy := validManifest()
	legacy.ID = "omnexa.legacy"
	legacy.Name = "Legacy"
	legacy.Owner = "legacy.team"
	legacy.Dependencies = nil
	legacy.OptionalDependencies = nil
	legacy.PlatformDependencies = nil
	legacy.Publisher = "legacy-publisher"
	legacy.ProvenanceRef = "oci://registry.example.test/omnexa/legacy@sha256:abcdef"
	legacy.SBOMRef = "sbom://legacy/1.2.3"
	legacy.DataClassification = []string{"INTERNAL", "PUBLIC"}
	legacy.Security.SecretReferences = []SecretReference{{Name: "legacy-api", Reference: "secret://vault/legacy/api"}}
	legacy.Security.NetworkDestinations = []string{"https://legacy.example.test"}

	current := validManifestV2()
	current.ID = "omnexa.current"
	current.Name = "Current"
	current.Owner = "current.team"
	current.Dependencies = nil
	current.OptionalDependencies = nil
	current.PlatformDependencies = nil
	current.Publisher = "current-publisher"
	current.ProvenanceRef = "oci://registry.example.test/omnexa/current@sha256:012345"
	current.SBOMRef = "sbom://current/1.2.3"
	current.DataClassification = []string{"CONFIDENTIAL", "INTERNAL"}
	current.Security.SecretReferences = []SecretReference{{Name: "current-api", Reference: "secret://vault/current/api"}}
	current.Security.NetworkDestinations = []string{"https://current.example.test"}

	legacyPayload, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	currentPayload, err := json.Marshal(current)
	if err != nil {
		t.Fatal(err)
	}
	registry, err := Discover([]DiscoverySource{{
		ID:        "trust.reference",
		Version:   "1.0.0",
		Manifests: [][]byte{legacyPayload, currentPayload},
	}})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	profiles, err := BuildPackageTrustProfiles(registry)
	if err != nil {
		t.Fatalf("BuildPackageTrustProfiles() error = %v", err)
	}
	if len(profiles) != 2 || profiles[0].ModuleID != current.ID || profiles[1].ModuleID != legacy.ID {
		t.Fatalf("unexpected deterministic profile order: %#v", profiles)
	}

	for _, profile := range profiles {
		if profile.ContractVersion != PackageTrustContractVersion || profile.Authority != PackageTrustMetadataOnly {
			t.Fatalf("profile claimed unexpected authority/version: %#v", profile)
		}
		if profile.Publisher == nil || profile.Publisher.ContractVersion != PackageTrustContractVersion ||
			profile.Provenance == nil || profile.Provenance.ContractVersion != PackageTrustContractVersion ||
			profile.SBOM == nil || profile.SBOM.ContractVersion != PackageTrustContractVersion ||
			profile.Scope.ContractVersion != PackageTrustContractVersion {
			t.Fatalf("typed trust hook contract missing: %#v", profile)
		}
		if len(profile.Scope.SecretReferences) != 1 || !strings.HasSuffix(profile.Scope.SecretReferences[0], "-api") {
			t.Fatalf("secret declaration name not retained safely: %#v", profile.Scope.SecretReferences)
		}
		encoded, marshalErr := json.Marshal(profile)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		if strings.Contains(string(encoded), "secret://") || strings.Contains(string(encoded), "/vault/") {
			t.Fatalf("profile leaked secret locator: %s", encoded)
		}
	}
}

func TestPackageTrustProfileIsDeterministicAndIndependentFromCallerMutation(t *testing.T) {
	manifest := simpleManifestV2("omnexa.inventory", "1.2.3")
	manifest.Publisher = "vertex-systems-network"
	manifest.ProvenanceRef = "oci://registry.example.test/omnexa/inventory@sha256:012345"
	manifest.SBOMRef = "sbom://inventory/1.2.3"
	manifest.DataClassification = []string{"INTERNAL", "PUBLIC"}
	manifest.CapabilitiesProvided = []string{"inventory.z", "inventory.a"}
	manifest.Security.SecretReferences = []SecretReference{
		{Name: "z-secret", Reference: "secret://vault/inventory/z"},
		{Name: "a-secret", Reference: "secret://vault/inventory/a"},
	}
	manifest.Security.NetworkDestinations = []string{"https://z.example.test", "https://a.example.test"}
	registry := discoverV2Registry(t, manifest)

	first, err := BuildPackageTrustProfiles(registry)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildPackageTrustProfiles(registry)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("profile output is not deterministic:\nfirst=%#v\nsecond=%#v", first, second)
	}
	if len(first) != 1 || !reflect.DeepEqual(first[0].Scope.CapabilitiesProvided, []string{"inventory.a", "inventory.z"}) ||
		!reflect.DeepEqual(first[0].Scope.SecretReferences, []string{"a-secret", "z-secret"}) {
		t.Fatalf("profile declarations are not canonically ordered: %#v", first)
	}

	first[0].Publisher.Identity = "mutated"
	first[0].Scope.CapabilitiesProvided[0] = "mutated.capability"
	first[0].Scope.SecretReferences[0] = "mutated-secret"
	fresh, found, err := PackageTrustProfileFor(registry, manifest.ID)
	if err != nil || !found {
		t.Fatalf("PackageTrustProfileFor() found=%v err=%v", found, err)
	}
	if fresh.Publisher == nil || fresh.Publisher.Identity != manifest.Publisher ||
		fresh.Scope.CapabilitiesProvided[0] != "inventory.a" || fresh.Scope.SecretReferences[0] != "a-secret" {
		t.Fatalf("caller mutation changed registry-bound evidence: %#v", fresh)
	}
}

func TestPackageTrustProfileOptionalHooksRemainAbsentWithoutClaimingTrust(t *testing.T) {
	manifest := simpleManifestV2("omnexa.inventory", "1.0.0")
	manifest.Publisher = ""
	manifest.ProvenanceRef = ""
	manifest.SBOMRef = ""
	registry := discoverV2Registry(t, manifest)

	profile, found, err := PackageTrustProfileFor(registry, manifest.ID)
	if err != nil || !found {
		t.Fatalf("PackageTrustProfileFor() found=%v err=%v", found, err)
	}
	if profile.Publisher != nil || profile.Provenance != nil || profile.SBOM != nil {
		t.Fatalf("absent optional hooks became present: %#v", profile)
	}
	if profile.Authority != PackageTrustMetadataOnly {
		t.Fatalf("optional metadata absence changed authority semantics: %#v", profile)
	}
}

func TestPackageTrustProfileFailsClosedOnRegistrySnapshotMismatch(t *testing.T) {
	manifest := simpleManifestV2("omnexa.inventory", "1.0.0")
	registry := discoverV2Registry(t, manifest)
	snapshot := registry.snapshots[manifest.ID]
	snapshot.Owner = "mutated.team"
	registry.snapshots[manifest.ID] = snapshot

	_, err := BuildPackageTrustProfiles(registry)
	var trustErr *PackageTrustError
	if !errors.As(err, &trustErr) || trustErr.Diagnostic.Code != "package.trust.registry.snapshot_mismatch" || trustErr.Diagnostic.ModuleID != manifest.ID {
		t.Fatalf("unexpected trust profile error: %#v", err)
	}
}

func TestPackageTrustProfilePublicSurfaceHasNoTrustDecisionOrSecretValueFields(t *testing.T) {
	for _, value := range []any{
		PublisherIdentityHook{},
		PackageProvenanceHook{},
		SBOMIdentityHook{},
		PackageDeclaredScopeProfile{},
		PackageTrustProfile{},
		PackageTrustDiagnostic{},
	} {
		typeOf := reflect.TypeOf(value)
		for index := 0; index < typeOf.NumField(); index++ {
			field := typeOf.Field(index)
			lower := strings.ToLower(field.Name)
			for _, forbidden := range []string{
				"trusted", "certified", "signaturevalid", "verificationresult", "secretvalue",
				"password", "credential", "token", "rawmanifest", "handler", "callback",
			} {
				if strings.Contains(lower, forbidden) {
					t.Fatalf("authoritative or unsafe field %s.%s", typeOf.Name(), field.Name)
				}
			}
			if field.Type.Kind() == reflect.Func {
				t.Fatalf("executable field exposed in %s.%s", typeOf.Name(), field.Name)
			}
		}
	}
}

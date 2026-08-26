package modules

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func validManifest() Manifest {
	return Manifest{
		SchemaVersion:           SchemaVersion,
		ID:                      "omnexa.inventory",
		Name:                    "Inventory",
		Version:                 "1.2.3",
		ContractVersion:         1,
		Status:                  "stable",
		Owner:                   "inventory.team",
		Runtime:                 "go",
		RequiredPlatformVersion: ">=1.0.0",
		Dependencies:            []string{"omnexa.catalog"},
		OptionalDependencies:    []string{"omnexa.analytics"},
		PlatformDependencies:    []string{"kernel.identity", "kernel.tenancy"},
		CapabilitiesProvided:    []string{"inventory.reserve-stock.v1"},
		CapabilitiesConsumed:    []string{"catalog.product.read.v1"},
		Permissions:             []string{"inventory.stock.read"},
		EventsPublished:         []string{"inventory.stock-reserved.v1"},
		EventsConsumed:          []string{"catalog.product-archived.v1"},
		WorkflowTriggers:        []string{"inventory.low-stock"},
		WorkflowActions:         []string{"inventory.reserve-stock"},
		UISlots:                 []string{"inventory.dashboard"},
		Settings:                []string{"inventory.reorder-threshold"},
		FeatureFlags:            []string{"inventory.batch-reservation"},
		DataClassification:      []string{"INTERNAL"},
		Migrations:              []string{"inventory.0001.initial"},
		LifecycleHooks:          []string{"install", "enable", "disable", "health_check"},
		HealthChecks:            []string{"inventory.database"},
		Security: SecurityDeclaration{
			SecretReferences:     []SecretReference{{Name: "erp-api", Reference: "secret://vault/inventory/erp-api"}},
			NetworkDestinations:  []string{"https://inventory.example.test"},
			ExposedEndpoints:     []string{"/api/v1/inventory"},
			FileTypes:            []string{"text/csv"},
			PrivilegedOperations: []string{"inventory.stock-adjustment"},
			AICapabilities:       []string{},
		},
		Publisher:     "vertex-systems-network",
		ProvenanceRef: "oci://registry.example.test/omnexa/inventory@sha256:0123456789abcdef",
		SBOMRef:       "sbom://inventory/1.2.3",
	}
}

func TestParseManifestAcceptsReferenceManifest(t *testing.T) {
	payload, err := json.Marshal(validManifest())
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}

	manifest, err := ParseManifest(payload)
	if err != nil {
		t.Fatalf("ParseManifest() error = %v", err)
	}
	if manifest.ID != "omnexa.inventory" || manifest.SchemaVersion != SchemaVersion {
		t.Fatalf("unexpected manifest identity: %#v", manifest)
	}
}

func TestParseManifestRejectsUnknownField(t *testing.T) {
	payload := manifestJSON(t)
	payload = strings.TrimSuffix(payload, "}") + `,"authorization_grant":true}`

	_, err := ParseManifest([]byte(payload))
	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected ParseError, got %T: %v", err, err)
	}
	if parseErr.Code != "manifest.parse.unknown_field" || parseErr.Path != "authorization_grant" {
		t.Fatalf("unexpected parse error: %#v", parseErr)
	}
}

func TestParseManifestRejectsMissingBaselineField(t *testing.T) {
	var raw map[string]any
	if err := json.Unmarshal([]byte(manifestJSON(t)), &raw); err != nil {
		t.Fatal(err)
	}
	delete(raw, "permissions")
	payload, _ := json.Marshal(raw)

	_, err := ParseManifest(payload)
	var parseErr *ParseError
	if !errors.As(err, &parseErr) || parseErr.Code != "manifest.parse.missing_field" || parseErr.Path != "permissions" {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestParseManifestRejectsOversizedInput(t *testing.T) {
	payload := make([]byte, MaxManifestBytes+1)
	for index := range payload {
		payload[index] = 'x'
	}
	_, err := ParseManifest(payload)
	var parseErr *ParseError
	if !errors.As(err, &parseErr) || parseErr.Code != "manifest.parse.too_large" {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestValidateManifestRejectsUnsupportedSchemaAndConflictingDependencyClass(t *testing.T) {
	manifest := validManifest()
	manifest.SchemaVersion = 2
	manifest.OptionalDependencies = append(manifest.OptionalDependencies, "omnexa.catalog")

	issues := validationIssues(t, ValidateManifest(manifest))
	assertIssue(t, issues, "schema_version", "manifest.schema.unsupported")
	assertIssue(t, issues, "optional_dependencies[1]", "manifest.dependency.class_conflict")
}

func TestValidateManifestRejectsDuplicateDeclarationAndInvalidLifecycleHook(t *testing.T) {
	manifest := validManifest()
	manifest.Permissions = []string{"inventory.stock.read", "inventory.stock.read"}
	manifest.LifecycleHooks = append(manifest.LifecycleHooks, "execute-arbitrary-code")

	issues := validationIssues(t, ValidateManifest(manifest))
	assertIssue(t, issues, "permissions[1]", "manifest.declaration.duplicate")
	assertIssue(t, issues, "lifecycle_hooks[4]", "manifest.declaration.invalid")
}

func TestValidateManifestRejectsRawSecretAndInsecureNetworkDestination(t *testing.T) {
	manifest := validManifest()
	manifest.Security.SecretReferences[0].Reference = "plain-text-api-key"
	manifest.Security.NetworkDestinations = []string{"http://inventory.example.test"}

	issues := validationIssues(t, ValidateManifest(manifest))
	assertIssue(t, issues, "security.secret_references[0].reference", "manifest.secret.reference_invalid")
	assertIssue(t, issues, "security.network_destinations[0]", "manifest.network_destination.invalid")
}

func TestValidationIssuesAreDeterministicAndSafe(t *testing.T) {
	manifest := validManifest()
	manifest.ID = "INJECTED SECRET = super-secret"
	manifest.Owner = ""
	manifest.Runtime = "javascript"

	err := ValidateManifest(manifest)
	issues := validationIssues(t, err)
	for index := 1; index < len(issues); index++ {
		previous := issues[index-1]
		current := issues[index]
		if previous.Path > current.Path || (previous.Path == current.Path && previous.Code > current.Code) {
			t.Fatalf("issues are not sorted: %#v", issues)
		}
	}
	if strings.Contains(err.Error(), "super-secret") {
		t.Fatal("validation error leaked untrusted manifest value")
	}
}

func TestParseManifestRejectsTrailingJSON(t *testing.T) {
	payload := manifestJSON(t) + ` {}`
	_, err := ParseManifest([]byte(payload))
	var parseErr *ParseError
	if !errors.As(err, &parseErr) || parseErr.Code != "manifest.parse.invalid_json" {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func manifestJSON(t *testing.T) string {
	t.Helper()
	payload, err := json.Marshal(validManifest())
	if err != nil {
		t.Fatal(err)
	}
	return string(payload)
}

func validationIssues(t *testing.T, err error) []ValidationIssue {
	t.Helper()
	var validationErr *ValidationErrors
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationErrors, got %T: %v", err, err)
	}
	return validationErr.Issues()
}

func assertIssue(t *testing.T, issues []ValidationIssue, path, code string) {
	t.Helper()
	for _, issue := range issues {
		if issue.Path == path && issue.Code == code {
			return
		}
	}
	t.Fatalf("missing issue %s at %s in %#v", code, path, issues)
}

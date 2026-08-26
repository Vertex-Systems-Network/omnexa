package modules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
)

const (
	// SchemaVersion is the only manifest schema version accepted by P03.01.
	SchemaVersion = 1
	// MaxManifestBytes bounds untrusted manifest input before decoding.
	MaxManifestBytes = 256 * 1024
	// MaxDeclarationItems bounds every declaration list in schema v1.
	MaxDeclarationItems = 256
)

// Manifest is the canonical P03.01 machine-readable module manifest.
// It is declarative metadata only; parsing or validation never executes module code.
type Manifest struct {
	SchemaVersion           int                 `json:"schema_version"`
	ID                      string              `json:"id"`
	Name                    string              `json:"name"`
	Version                 string              `json:"version"`
	ContractVersion         int                 `json:"contract_version"`
	Status                  string              `json:"status"`
	Owner                   string              `json:"owner"`
	Runtime                 string              `json:"runtime"`
	RequiredPlatformVersion string              `json:"required_platform_version"`
	Dependencies            []string            `json:"dependencies"`
	OptionalDependencies    []string            `json:"optional_dependencies"`
	PlatformDependencies    []string            `json:"platform_dependencies"`
	CapabilitiesProvided    []string            `json:"capabilities_provided"`
	CapabilitiesConsumed    []string            `json:"capabilities_consumed"`
	Permissions             []string            `json:"permissions"`
	EventsPublished         []string            `json:"events_published"`
	EventsConsumed          []string            `json:"events_consumed"`
	WorkflowTriggers        []string            `json:"workflow_triggers"`
	WorkflowActions         []string            `json:"workflow_actions"`
	UISlots                 []string            `json:"ui_slots"`
	Settings                []string            `json:"settings"`
	FeatureFlags            []string            `json:"feature_flags"`
	DataClassification      []string            `json:"data_classification"`
	Migrations              []string            `json:"migrations"`
	LifecycleHooks          []string            `json:"lifecycle_hooks"`
	HealthChecks            []string            `json:"health_checks"`
	Security                SecurityDeclaration `json:"security"`
	Publisher               string              `json:"publisher,omitempty"`
	ProvenanceRef           string              `json:"provenance_ref,omitempty"`
	SBOMRef                 string              `json:"sbom_ref,omitempty"`
}

// SecurityDeclaration records declared security scope. These declarations do not
// create permissions or authorization grants.
type SecurityDeclaration struct {
	SecretReferences     []SecretReference `json:"secret_references"`
	NetworkDestinations  []string          `json:"network_destinations"`
	ExposedEndpoints     []string          `json:"exposed_endpoints"`
	FileTypes            []string          `json:"file_types"`
	PrivilegedOperations []string          `json:"privileged_operations"`
	AICapabilities       []string          `json:"ai_capabilities"`
}

// SecretReference names a symbolic secret reference. Secret values never belong
// in module manifests.
type SecretReference struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

// ParseError is a stable, safe parser failure. Error intentionally excludes raw
// untrusted manifest content and decoder text.
type ParseError struct {
	Code string
	Path string
}

func (e *ParseError) Error() string {
	if e == nil {
		return "module manifest parse failed"
	}
	return fmt.Sprintf("module manifest parse failed (%s)", e.Code)
}

var requiredFields = []string{
	"schema_version",
	"id",
	"name",
	"version",
	"contract_version",
	"status",
	"owner",
	"runtime",
	"required_platform_version",
	"dependencies",
	"optional_dependencies",
	"platform_dependencies",
	"capabilities_provided",
	"capabilities_consumed",
	"permissions",
	"events_published",
	"events_consumed",
	"workflow_triggers",
	"workflow_actions",
	"ui_slots",
	"settings",
	"feature_flags",
	"data_classification",
	"migrations",
	"lifecycle_hooks",
	"health_checks",
	"security",
}

var requiredSecurityFields = []string{
	"secret_references",
	"network_destinations",
	"exposed_endpoints",
	"file_types",
	"privileged_operations",
	"ai_capabilities",
}

var allowedFields = func() map[string]struct{} {
	fields := append([]string{}, requiredFields...)
	fields = append(fields, "publisher", "provenance_ref", "sbom_ref")
	result := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		result[field] = struct{}{}
	}
	return result
}()

// ParseManifest decodes one bounded JSON manifest, rejects unknown schema-v1
// fields, verifies presence of baseline fields, and runs deterministic semantic
// validation. It never loads packages, hooks, files, environment variables, or
// network resources.
func ParseManifest(data []byte) (Manifest, error) {
	if len(data) == 0 {
		return Manifest{}, &ParseError{Code: "manifest.parse.empty"}
	}
	if len(data) > MaxManifestBytes {
		return Manifest{}, &ParseError{Code: "manifest.parse.too_large"}
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil || raw == nil {
		return Manifest{}, &ParseError{Code: "manifest.parse.invalid_json"}
	}

	unknown := make([]string, 0)
	for field := range raw {
		if _, ok := allowedFields[field]; !ok {
			unknown = append(unknown, field)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return Manifest{}, &ParseError{Code: "manifest.parse.unknown_field", Path: unknown[0]}
	}

	missing := make([]string, 0)
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return Manifest{}, &ParseError{Code: "manifest.parse.missing_field", Path: missing[0]}
	}

	var rawSecurity map[string]json.RawMessage
	if err := json.Unmarshal(raw["security"], &rawSecurity); err != nil || rawSecurity == nil {
		return Manifest{}, &ParseError{Code: "manifest.parse.invalid_security"}
	}
	missingSecurity := make([]string, 0)
	for _, field := range requiredSecurityFields {
		if _, ok := rawSecurity[field]; !ok {
			missingSecurity = append(missingSecurity, field)
		}
	}
	if len(missingSecurity) > 0 {
		sort.Strings(missingSecurity)
		return Manifest{}, &ParseError{Code: "manifest.parse.missing_security_field", Path: "security." + missingSecurity[0]}
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, &ParseError{Code: "manifest.parse.invalid_json"}
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return Manifest{}, &ParseError{Code: "manifest.parse.trailing_content"}
	}

	if err := ValidateManifest(manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

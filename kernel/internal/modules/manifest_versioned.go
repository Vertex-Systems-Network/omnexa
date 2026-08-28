package modules

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"sort"
)

const SchemaVersionV2 = 2

// DependencyRequirement is the schema-v2 module dependency declaration.
// Constraint uses the bounded comparator grammar defined by ADR-0012.
type DependencyRequirement struct {
	ID         string `json:"id"`
	Constraint string `json:"constraint"`
}

type manifestV2 struct {
	SchemaVersion           int                     `json:"schema_version"`
	ID                      string                  `json:"id"`
	Name                    string                  `json:"name"`
	Version                 string                  `json:"version"`
	ContractVersion         int                     `json:"contract_version"`
	Status                  string                  `json:"status"`
	Owner                   string                  `json:"owner"`
	Runtime                 string                  `json:"runtime"`
	RequiredPlatformVersion string                  `json:"required_platform_version"`
	Dependencies            []DependencyRequirement `json:"dependencies"`
	OptionalDependencies    []DependencyRequirement `json:"optional_dependencies"`
	PlatformDependencies    []string                `json:"platform_dependencies"`
	CapabilitiesProvided    []string                `json:"capabilities_provided"`
	CapabilitiesConsumed    []string                `json:"capabilities_consumed"`
	Permissions             []string                `json:"permissions"`
	EventsPublished         []string                `json:"events_published"`
	EventsConsumed          []string                `json:"events_consumed"`
	WorkflowTriggers        []string                `json:"workflow_triggers"`
	WorkflowActions         []string                `json:"workflow_actions"`
	UISlots                 []string                `json:"ui_slots"`
	Settings                []string                `json:"settings"`
	FeatureFlags            []string                `json:"feature_flags"`
	DataClassification      []string                `json:"data_classification"`
	Migrations              []string                `json:"migrations"`
	LifecycleHooks          []string                `json:"lifecycle_hooks"`
	HealthChecks            []string                `json:"health_checks"`
	Security                SecurityDeclaration     `json:"security"`
	Publisher               string                  `json:"publisher,omitempty"`
	ProvenanceRef           string                  `json:"provenance_ref,omitempty"`
	SBOMRef                 string                  `json:"sbom_ref,omitempty"`
}

// validatedManifestSnapshot is the package-private immutable-by-convention
// normalized representation bound atomically to one RegistryRecord.
type validatedManifestSnapshot struct {
	SchemaVersion           int
	ID                      string
	Version                 string
	Owner                   string
	RequiredPlatformVersion string
	RequiredDependencies    []DependencyRequirement
	OptionalDependencies    []DependencyRequirement
	PlatformDependencies    []string
	CapabilitiesProvided    []string
	CapabilitiesConsumed    []string
	Permissions             []string
	UISlots                 []string
	Settings                []string
	FeatureFlags            []string
}

func (s validatedManifestSnapshot) clone() validatedManifestSnapshot {
	return validatedManifestSnapshot{
		SchemaVersion:           s.SchemaVersion,
		ID:                      s.ID,
		Version:                 s.Version,
		Owner:                   s.Owner,
		RequiredPlatformVersion: s.RequiredPlatformVersion,
		RequiredDependencies:    copyDependencyRequirements(s.RequiredDependencies),
		OptionalDependencies:    copyDependencyRequirements(s.OptionalDependencies),
		PlatformDependencies:    append([]string(nil), s.PlatformDependencies...),
		CapabilitiesProvided:    append([]string(nil), s.CapabilitiesProvided...),
		CapabilitiesConsumed:    append([]string(nil), s.CapabilitiesConsumed...),
		Permissions:             append([]string(nil), s.Permissions...),
		UISlots:                 append([]string(nil), s.UISlots...),
		Settings:                append([]string(nil), s.Settings...),
		FeatureFlags:            append([]string(nil), s.FeatureFlags...),
	}
}

func (s validatedManifestSnapshot) declaresModuleDependency(id string) bool {
	for _, requirement := range s.RequiredDependencies {
		if requirement.ID == id {
			return true
		}
	}
	for _, requirement := range s.OptionalDependencies {
		if requirement.ID == id {
			return true
		}
	}
	return false
}

func copyDependencyRequirements(values []DependencyRequirement) []DependencyRequirement {
	result := make([]DependencyRequirement, len(values))
	copy(result, values)
	return result
}

// parseValidatedManifest is the P03.03 version-dispatch boundary. It preserves
// the accepted P03.01 schema-v1 parser while adding a separate strict schema-v2
// decoder. It never falls back between schema versions after dispatch.
func parseValidatedManifest(data []byte) (validatedManifestSnapshot, error) {
	if len(data) == 0 {
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.empty"}
	}
	if len(data) > MaxManifestBytes {
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.too_large"}
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil || raw == nil {
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.invalid_json"}
	}
	versionRaw, ok := raw["schema_version"]
	if !ok {
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.missing_field", Path: "schema_version"}
	}
	var version int
	if err := json.Unmarshal(versionRaw, &version); err != nil {
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.schema.invalid", Path: "schema_version"}
	}

	switch version {
	case SchemaVersion:
		manifest, err := ParseManifest(data)
		if err != nil {
			return validatedManifestSnapshot{}, err
		}
		return snapshotFromV1(manifest), nil
	case SchemaVersionV2:
		return parseManifestV2(data, raw)
	default:
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.schema.unsupported", Path: "schema_version"}
	}
}

func snapshotFromV1(manifest Manifest) validatedManifestSnapshot {
	required := make([]DependencyRequirement, 0, len(manifest.Dependencies))
	for _, id := range manifest.Dependencies {
		required = append(required, DependencyRequirement{ID: id})
	}
	optional := make([]DependencyRequirement, 0, len(manifest.OptionalDependencies))
	for _, id := range manifest.OptionalDependencies {
		optional = append(optional, DependencyRequirement{ID: id})
	}
	return validatedManifestSnapshot{
		SchemaVersion:           SchemaVersion,
		ID:                      manifest.ID,
		Version:                 manifest.Version,
		Owner:                   manifest.Owner,
		RequiredPlatformVersion: manifest.RequiredPlatformVersion,
		RequiredDependencies:    required,
		OptionalDependencies:    optional,
		PlatformDependencies:    append([]string(nil), manifest.PlatformDependencies...),
		CapabilitiesProvided:    append([]string(nil), manifest.CapabilitiesProvided...),
		CapabilitiesConsumed:    append([]string(nil), manifest.CapabilitiesConsumed...),
		Permissions:             append([]string(nil), manifest.Permissions...),
		UISlots:                 append([]string(nil), manifest.UISlots...),
		Settings:                append([]string(nil), manifest.Settings...),
		FeatureFlags:            append([]string(nil), manifest.FeatureFlags...),
	}
}

func parseManifestV2(data []byte, raw map[string]json.RawMessage) (validatedManifestSnapshot, error) {
	unknown := make([]string, 0)
	for field := range raw {
		if _, ok := allowedFields[field]; !ok {
			unknown = append(unknown, field)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.unknown_field", Path: unknown[0]}
	}

	missing := make([]string, 0)
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.missing_field", Path: missing[0]}
	}

	var rawSecurity map[string]json.RawMessage
	if err := json.Unmarshal(raw["security"], &rawSecurity); err != nil || rawSecurity == nil {
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.invalid_security"}
	}
	missingSecurity := make([]string, 0)
	for _, field := range requiredSecurityFields {
		if _, ok := rawSecurity[field]; !ok {
			missingSecurity = append(missingSecurity, field)
		}
	}
	if len(missingSecurity) > 0 {
		sort.Strings(missingSecurity)
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.missing_security_field", Path: "security." + missingSecurity[0]}
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var manifest manifestV2
	if err := decoder.Decode(&manifest); err != nil {
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.invalid_json"}
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return validatedManifestSnapshot{}, &ParseError{Code: "manifest.parse.trailing_content"}
	}
	if err := validateManifestV2(manifest); err != nil {
		return validatedManifestSnapshot{}, err
	}

	return validatedManifestSnapshot{
		SchemaVersion:           SchemaVersionV2,
		ID:                      manifest.ID,
		Version:                 manifest.Version,
		Owner:                   manifest.Owner,
		RequiredPlatformVersion: manifest.RequiredPlatformVersion,
		RequiredDependencies:    copyDependencyRequirements(manifest.Dependencies),
		OptionalDependencies:    copyDependencyRequirements(manifest.OptionalDependencies),
		PlatformDependencies:    append([]string(nil), manifest.PlatformDependencies...),
		CapabilitiesProvided:    append([]string(nil), manifest.CapabilitiesProvided...),
		CapabilitiesConsumed:    append([]string(nil), manifest.CapabilitiesConsumed...),
		Permissions:             append([]string(nil), manifest.Permissions...),
		UISlots:                 append([]string(nil), manifest.UISlots...),
		Settings:                append([]string(nil), manifest.Settings...),
		FeatureFlags:            append([]string(nil), manifest.FeatureFlags...),
	}, nil
}

func validateManifestV2(manifest manifestV2) error {
	collector := &validationCollector{}
	if manifest.SchemaVersion != SchemaVersionV2 {
		collector.add("schema_version", "manifest.schema.unsupported")
	}

	common := Manifest{
		SchemaVersion:           SchemaVersion,
		ID:                      manifest.ID,
		Name:                    manifest.Name,
		Version:                 manifest.Version,
		ContractVersion:         manifest.ContractVersion,
		Status:                  manifest.Status,
		Owner:                   manifest.Owner,
		Runtime:                 manifest.Runtime,
		RequiredPlatformVersion: manifest.RequiredPlatformVersion,
		Dependencies:            nil,
		OptionalDependencies:    nil,
		PlatformDependencies:    append([]string(nil), manifest.PlatformDependencies...),
		CapabilitiesProvided:    append([]string(nil), manifest.CapabilitiesProvided...),
		CapabilitiesConsumed:    append([]string(nil), manifest.CapabilitiesConsumed...),
		Permissions:             append([]string(nil), manifest.Permissions...),
		EventsPublished:         append([]string(nil), manifest.EventsPublished...),
		EventsConsumed:          append([]string(nil), manifest.EventsConsumed...),
		WorkflowTriggers:        append([]string(nil), manifest.WorkflowTriggers...),
		WorkflowActions:         append([]string(nil), manifest.WorkflowActions...),
		UISlots:                 append([]string(nil), manifest.UISlots...),
		Settings:                append([]string(nil), manifest.Settings...),
		FeatureFlags:            append([]string(nil), manifest.FeatureFlags...),
		DataClassification:      append([]string(nil), manifest.DataClassification...),
		Migrations:              append([]string(nil), manifest.Migrations...),
		LifecycleHooks:          append([]string(nil), manifest.LifecycleHooks...),
		HealthChecks:            append([]string(nil), manifest.HealthChecks...),
		Security:                manifest.Security,
		Publisher:               manifest.Publisher,
		ProvenanceRef:           manifest.ProvenanceRef,
		SBOMRef:                 manifest.SBOMRef,
	}
	if err := ValidateManifest(common); err != nil {
		var validationErr *ValidationErrors
		if errors.As(err, &validationErr) {
			collector.issues = append(collector.issues, validationErr.Issues()...)
		} else {
			collector.add("manifest", "manifest.validation.invalid")
		}
	}
	if semverPattern.MatchString(manifest.Version) && len(manifest.Version) <= maxIdentifierLength {
		if _, ok := parseStrictSemVer(manifest.Version); !ok {
			collector.add("version", "manifest.version.invalid")
		}
	}

	validateV2DependencyClasses(collector, manifest.ID, manifest.Dependencies, manifest.OptionalDependencies)
	return collector.err()
}

func validateV2DependencyClasses(c *validationCollector, moduleID string, required, optional []DependencyRequirement) {
	classes := []struct {
		path   string
		values []DependencyRequirement
	}{
		{path: "dependencies", values: required},
		{path: "optional_dependencies", values: optional},
	}

	seenAcross := make(map[string]string)
	for _, class := range classes {
		if len(class.values) > MaxDeclarationItems {
			c.add(class.path, "manifest.list.too_many_items")
			continue
		}
		seenWithin := make(map[string]struct{}, len(class.values))
		for index, requirement := range class.values {
			base := class.path + "[" + decimalIndex(index) + "]"
			if !validBoundedPattern(requirement.ID, moduleIDPattern, maxIdentifierLength) {
				c.add(base+".id", "manifest.dependency.invalid")
			}
			if requirement.ID == moduleID {
				c.add(base+".id", "manifest.dependency.self")
			}
			switch _, failure := parseDependencyConstraint(requirement.Constraint); failure {
			case constraintTooLong:
				c.add(base+".constraint", "manifest.dependency.constraint_too_long")
			case constraintTooManyComparators:
				c.add(base+".constraint", "manifest.dependency.constraint_too_many_comparators")
			case constraintInvalid:
				c.add(base+".constraint", "manifest.dependency.constraint_invalid")
			}
			if _, exists := seenWithin[requirement.ID]; exists {
				c.add(base+".id", "manifest.dependency.duplicate")
			} else {
				seenWithin[requirement.ID] = struct{}{}
			}
			if previous, exists := seenAcross[requirement.ID]; exists && previous != class.path {
				c.add(base+".id", "manifest.dependency.class_conflict")
			} else if !exists {
				seenAcross[requirement.ID] = class.path
			}
		}
	}
}

func decimalIndex(value int) string {
	if value == 0 {
		return "0"
	}
	var digits [20]byte
	position := len(digits)
	for value > 0 {
		position--
		digits[position] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[position:])
}

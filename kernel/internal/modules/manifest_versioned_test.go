package modules

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func validManifestV2() manifestV2 {
	base := validManifest()
	return manifestV2{
		SchemaVersion:           SchemaVersionV2,
		ID:                      base.ID,
		Name:                    base.Name,
		Version:                 base.Version,
		ContractVersion:         base.ContractVersion,
		Status:                  base.Status,
		Owner:                   base.Owner,
		Runtime:                 base.Runtime,
		RequiredPlatformVersion: base.RequiredPlatformVersion,
		Dependencies: []DependencyRequirement{{
			ID:         "omnexa.catalog",
			Constraint: ">=2.0.0 <3.0.0",
		}},
		OptionalDependencies: []DependencyRequirement{{
			ID:         "omnexa.analytics",
			Constraint: ">=1.0.0 <2.0.0",
		}},
		PlatformDependencies: append([]string(nil), base.PlatformDependencies...),
		CapabilitiesProvided: append([]string(nil), base.CapabilitiesProvided...),
		CapabilitiesConsumed: append([]string(nil), base.CapabilitiesConsumed...),
		Permissions:          append([]string(nil), base.Permissions...),
		EventsPublished:      append([]string(nil), base.EventsPublished...),
		EventsConsumed:       append([]string(nil), base.EventsConsumed...),
		WorkflowTriggers:     append([]string(nil), base.WorkflowTriggers...),
		WorkflowActions:      append([]string(nil), base.WorkflowActions...),
		UISlots:              append([]string(nil), base.UISlots...),
		Settings:             append([]string(nil), base.Settings...),
		FeatureFlags:         append([]string(nil), base.FeatureFlags...),
		DataClassification:   append([]string(nil), base.DataClassification...),
		Migrations:           append([]string(nil), base.Migrations...),
		LifecycleHooks:       append([]string(nil), base.LifecycleHooks...),
		HealthChecks:         append([]string(nil), base.HealthChecks...),
		Security:             base.Security,
		Publisher:            base.Publisher,
		ProvenanceRef:        base.ProvenanceRef,
		SBOMRef:              base.SBOMRef,
	}
}

func manifestV2Bytes(t *testing.T, manifest manifestV2) []byte {
	t.Helper()
	payload, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal schema-v2 manifest: %v", err)
	}
	return payload
}

func TestParseValidatedManifestDispatchesV1AndV2(t *testing.T) {
	v1Payload, err := json.Marshal(validManifest())
	if err != nil {
		t.Fatal(err)
	}
	v1, err := parseValidatedManifest(v1Payload)
	if err != nil {
		t.Fatalf("parse v1: %v", err)
	}
	if v1.SchemaVersion != SchemaVersion || v1.RequiredDependencies[0].Constraint != "" {
		t.Fatalf("unexpected v1 snapshot: %#v", v1)
	}

	v2 := validManifestV2()
	v2Snapshot, err := parseValidatedManifest(manifestV2Bytes(t, v2))
	if err != nil {
		t.Fatalf("parse v2: %v", err)
	}
	if v2Snapshot.SchemaVersion != SchemaVersionV2 || v2Snapshot.RequiredDependencies[0].Constraint != ">=2.0.0 <3.0.0" {
		t.Fatalf("unexpected v2 snapshot: %#v", v2Snapshot)
	}
}

func TestParseValidatedManifestRejectsUnsupportedSchemaWithoutFallback(t *testing.T) {
	manifest := validManifestV2()
	manifest.SchemaVersion = 3
	_, err := parseValidatedManifest(manifestV2Bytes(t, manifest))
	var parseErr *ParseError
	if !errors.As(err, &parseErr) || parseErr.Code != "manifest.schema.unsupported" || parseErr.Path != "schema_version" {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestParseValidatedManifestRejectsUnknownDependencyRecordField(t *testing.T) {
	var raw map[string]any
	if err := json.Unmarshal(manifestV2Bytes(t, validManifestV2()), &raw); err != nil {
		t.Fatal(err)
	}
	raw["dependencies"] = []any{map[string]any{
		"id":         "omnexa.catalog",
		"constraint": ">=2.0.0 <3.0.0",
		"authority":  true,
	}}
	payload, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	_, err = parseValidatedManifest(payload)
	var parseErr *ParseError
	if !errors.As(err, &parseErr) || parseErr.Code != "manifest.parse.invalid_json" {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestManifestV2RejectsSelfDependencyAndInvalidConstraint(t *testing.T) {
	manifest := validManifestV2()
	manifest.Dependencies = []DependencyRequirement{
		{ID: manifest.ID, Constraint: ">=1.0.0"},
		{ID: "omnexa.catalog", Constraint: ">=1.0.0  <2.0.0"},
	}
	issues := validationIssues(t, validateManifestV2(manifest))
	assertIssue(t, issues, "dependencies[0].id", "manifest.dependency.self")
	assertIssue(t, issues, "dependencies[1].constraint", "manifest.dependency.constraint_invalid")
}

func TestManifestV2RejectsDuplicateAndCrossClassDependency(t *testing.T) {
	manifest := validManifestV2()
	manifest.Dependencies = []DependencyRequirement{
		{ID: "omnexa.catalog", Constraint: ">=2.0.0"},
		{ID: "omnexa.catalog", Constraint: "<3.0.0"},
	}
	manifest.OptionalDependencies = []DependencyRequirement{{ID: "omnexa.catalog", Constraint: ">=2.0.0"}}
	issues := validationIssues(t, validateManifestV2(manifest))
	assertIssue(t, issues, "dependencies[1].id", "manifest.dependency.duplicate")
	assertIssue(t, issues, "optional_dependencies[0].id", "manifest.dependency.class_conflict")
}

func TestManifestV2UsesStrictSemVerForModuleVersion(t *testing.T) {
	manifest := validManifestV2()
	manifest.Version = "1.0.0-01"
	issues := validationIssues(t, validateManifestV2(manifest))
	assertIssue(t, issues, "version", "manifest.version.invalid")
}

func TestManifestV2ConstraintLengthBoundDoesNotLeakContent(t *testing.T) {
	manifest := validManifestV2()
	manifest.Dependencies[0].Constraint = ">=" + strings.Repeat("1", MaxDependencyConstraintBytes)
	err := validateManifestV2(manifest)
	issues := validationIssues(t, err)
	assertIssue(t, issues, "dependencies[0].constraint", "manifest.dependency.constraint_too_long")
	if strings.Contains(err.Error(), manifest.Dependencies[0].Constraint) {
		t.Fatal("validation error leaked raw dependency constraint")
	}
}

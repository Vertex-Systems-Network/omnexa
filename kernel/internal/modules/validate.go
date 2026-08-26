package modules

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

const (
	maxIdentifierLength = 128
	maxNameLength       = 160
	maxReferenceLength  = 512
)

var (
	moduleIDPattern        = regexp.MustCompile(`^omnexa\.[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$`)
	platformIDPattern      = regexp.MustCompile(`^kernel\.[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$`)
	identifierPattern      = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[._:-][a-z0-9][a-z0-9_-]*)*$`)
	ownerPattern           = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$`)
	semverPattern          = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
	platformVersionPattern = regexp.MustCompile(`^(?:>=|>|<=|<|=)(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?$`)
	mimeTypePattern        = regexp.MustCompile(`^[a-z0-9][a-z0-9!#$&^_.+-]*/[a-z0-9][a-z0-9!#$&^_.+-]*$`)
)

var allowedStatuses = map[string]struct{}{
	"stable":       {},
	"beta":         {},
	"experimental": {},
}

var allowedRuntimes = map[string]struct{}{
	"go":         {},
	"typescript": {},
	"rust":       {},
	"python":     {},
}

var allowedDataClassifications = map[string]struct{}{
	"PUBLIC":       {},
	"INTERNAL":     {},
	"CONFIDENTIAL": {},
	"RESTRICTED":   {},
}

var allowedLifecycleHooks = map[string]struct{}{
	"pre_install":  {},
	"install":      {},
	"post_install": {},
	"pre_enable":   {},
	"enable":       {},
	"disable":      {},
	"upgrade":      {},
	"rollback":     {},
	"archive":      {},
	"export":       {},
	"detach":       {},
	"purge":        {},
	"health_check": {},
}

// ValidationIssue is one stable, non-secret validation result.
type ValidationIssue struct {
	Code string `json:"code"`
	Path string `json:"path"`
}

// ValidationErrors is sorted deterministically by path then code.
type ValidationErrors struct {
	issues []ValidationIssue
}

func (e *ValidationErrors) Error() string {
	return "module manifest validation failed"
}

// Issues returns a copy so callers cannot mutate the validation result.
func (e *ValidationErrors) Issues() []ValidationIssue {
	if e == nil {
		return nil
	}
	result := make([]ValidationIssue, len(e.issues))
	copy(result, e.issues)
	return result
}

type validationCollector struct {
	issues []ValidationIssue
}

func (c *validationCollector) add(path, code string) {
	c.issues = append(c.issues, ValidationIssue{Code: code, Path: path})
}

func (c *validationCollector) err() error {
	if len(c.issues) == 0 {
		return nil
	}
	sort.SliceStable(c.issues, func(i, j int) bool {
		if c.issues[i].Path == c.issues[j].Path {
			return c.issues[i].Code < c.issues[j].Code
		}
		return c.issues[i].Path < c.issues[j].Path
	})
	return &ValidationErrors{issues: append([]ValidationIssue(nil), c.issues...)}
}

// ValidateManifest validates only declarations. It does not resolve dependencies,
// register capabilities, grant permissions, execute lifecycle hooks, or inspect
// external state.
func ValidateManifest(manifest Manifest) error {
	collector := &validationCollector{}

	if manifest.SchemaVersion != SchemaVersion {
		collector.add("schema_version", "manifest.schema.unsupported")
	}
	if !validBoundedPattern(manifest.ID, moduleIDPattern, maxIdentifierLength) {
		collector.add("id", "manifest.id.invalid")
	}
	if strings.TrimSpace(manifest.Name) == "" || len(manifest.Name) > maxNameLength {
		collector.add("name", "manifest.name.invalid")
	}
	if !semverPattern.MatchString(manifest.Version) || len(manifest.Version) > maxIdentifierLength {
		collector.add("version", "manifest.version.invalid")
	}
	if manifest.ContractVersion < 1 {
		collector.add("contract_version", "manifest.contract_version.invalid")
	}
	if _, ok := allowedStatuses[manifest.Status]; !ok {
		collector.add("status", "manifest.status.invalid")
	}
	if !validBoundedPattern(manifest.Owner, ownerPattern, maxIdentifierLength) {
		collector.add("owner", "manifest.owner.invalid")
	}
	if _, ok := allowedRuntimes[manifest.Runtime]; !ok {
		collector.add("runtime", "manifest.runtime.invalid")
	}
	if !platformVersionPattern.MatchString(manifest.RequiredPlatformVersion) || len(manifest.RequiredPlatformVersion) > maxIdentifierLength {
		collector.add("required_platform_version", "manifest.platform_version.invalid")
	}

	validateDependencyClasses(collector, manifest)
	validateIdentifierList(collector, "capabilities_provided", manifest.CapabilitiesProvided)
	validateIdentifierList(collector, "capabilities_consumed", manifest.CapabilitiesConsumed)
	validateIdentifierList(collector, "permissions", manifest.Permissions)
	validateIdentifierList(collector, "events_published", manifest.EventsPublished)
	validateIdentifierList(collector, "events_consumed", manifest.EventsConsumed)
	validateIdentifierList(collector, "workflow_triggers", manifest.WorkflowTriggers)
	validateIdentifierList(collector, "workflow_actions", manifest.WorkflowActions)
	validateIdentifierList(collector, "ui_slots", manifest.UISlots)
	validateIdentifierList(collector, "settings", manifest.Settings)
	validateIdentifierList(collector, "feature_flags", manifest.FeatureFlags)
	validateIdentifierList(collector, "migrations", manifest.Migrations)
	validateIdentifierList(collector, "health_checks", manifest.HealthChecks)

	validateEnumList(collector, "data_classification", manifest.DataClassification, allowedDataClassifications)
	validateEnumList(collector, "lifecycle_hooks", manifest.LifecycleHooks, allowedLifecycleHooks)
	validateSecurity(collector, manifest.Security)

	for path, value := range map[string]string{
		"publisher":      manifest.Publisher,
		"provenance_ref": manifest.ProvenanceRef,
		"sbom_ref":       manifest.SBOMRef,
	} {
		if len(value) > maxReferenceLength {
			collector.add(path, "manifest.reference.too_long")
		}
	}

	return collector.err()
}

func validateDependencyClasses(c *validationCollector, manifest Manifest) {
	classes := []struct {
		path    string
		values  []string
		pattern *regexp.Regexp
	}{
		{path: "dependencies", values: manifest.Dependencies, pattern: moduleIDPattern},
		{path: "optional_dependencies", values: manifest.OptionalDependencies, pattern: moduleIDPattern},
		{path: "platform_dependencies", values: manifest.PlatformDependencies, pattern: platformIDPattern},
	}

	seenAcross := make(map[string]string)
	for _, class := range classes {
		if len(class.values) > MaxDeclarationItems {
			c.add(class.path, "manifest.list.too_many_items")
			continue
		}
		seenWithin := make(map[string]struct{}, len(class.values))
		for index, value := range class.values {
			path := fmt.Sprintf("%s[%d]", class.path, index)
			if !validBoundedPattern(value, class.pattern, maxIdentifierLength) {
				c.add(path, "manifest.dependency.invalid")
			}
			if _, exists := seenWithin[value]; exists {
				c.add(path, "manifest.dependency.duplicate")
			} else {
				seenWithin[value] = struct{}{}
			}
			if previous, exists := seenAcross[value]; exists && previous != class.path {
				c.add(path, "manifest.dependency.class_conflict")
			} else if !exists {
				seenAcross[value] = class.path
			}
		}
	}
}

func validateIdentifierList(c *validationCollector, path string, values []string) {
	if len(values) > MaxDeclarationItems {
		c.add(path, "manifest.list.too_many_items")
		return
	}
	seen := make(map[string]struct{}, len(values))
	for index, value := range values {
		itemPath := fmt.Sprintf("%s[%d]", path, index)
		if !validBoundedPattern(value, identifierPattern, maxIdentifierLength) {
			c.add(itemPath, "manifest.declaration.invalid")
		}
		if _, exists := seen[value]; exists {
			c.add(itemPath, "manifest.declaration.duplicate")
		} else {
			seen[value] = struct{}{}
		}
	}
}

func validateEnumList(c *validationCollector, path string, values []string, allowed map[string]struct{}) {
	if len(values) > MaxDeclarationItems {
		c.add(path, "manifest.list.too_many_items")
		return
	}
	seen := make(map[string]struct{}, len(values))
	for index, value := range values {
		itemPath := fmt.Sprintf("%s[%d]", path, index)
		if _, ok := allowed[value]; !ok {
			c.add(itemPath, "manifest.declaration.invalid")
		}
		if _, exists := seen[value]; exists {
			c.add(itemPath, "manifest.declaration.duplicate")
		} else {
			seen[value] = struct{}{}
		}
	}
}

func validateSecurity(c *validationCollector, security SecurityDeclaration) {
	if len(security.SecretReferences) > MaxDeclarationItems {
		c.add("security.secret_references", "manifest.list.too_many_items")
	} else {
		seen := make(map[string]struct{}, len(security.SecretReferences))
		for index, secret := range security.SecretReferences {
			base := fmt.Sprintf("security.secret_references[%d]", index)
			if !validBoundedPattern(secret.Name, identifierPattern, maxIdentifierLength) {
				c.add(base+".name", "manifest.secret.name_invalid")
			}
			if !validSecretReference(secret.Reference) {
				c.add(base+".reference", "manifest.secret.reference_invalid")
			}
			if _, exists := seen[secret.Name]; exists {
				c.add(base+".name", "manifest.secret.duplicate")
			} else {
				seen[secret.Name] = struct{}{}
			}
		}
	}

	validateBoundedStrings(c, "security.network_destinations", security.NetworkDestinations, validNetworkDestination, "manifest.network_destination.invalid")
	validateBoundedStrings(c, "security.exposed_endpoints", security.ExposedEndpoints, validEndpoint, "manifest.endpoint.invalid")
	validateBoundedStrings(c, "security.file_types", security.FileTypes, func(value string) bool {
		return len(value) <= maxIdentifierLength && mimeTypePattern.MatchString(value)
	}, "manifest.file_type.invalid")
	validateSecurityIdentifierList(c, "security.privileged_operations", security.PrivilegedOperations)
	validateSecurityIdentifierList(c, "security.ai_capabilities", security.AICapabilities)
}

func validateSecurityIdentifierList(c *validationCollector, path string, values []string) {
	if len(values) > MaxDeclarationItems {
		c.add(path, "manifest.list.too_many_items")
		return
	}
	seen := make(map[string]struct{}, len(values))
	for index, value := range values {
		itemPath := fmt.Sprintf("%s[%d]", path, index)
		if !validBoundedPattern(value, identifierPattern, maxIdentifierLength) {
			c.add(itemPath, "manifest.security_declaration.invalid")
		}
		if _, exists := seen[value]; exists {
			c.add(itemPath, "manifest.declaration.duplicate")
		} else {
			seen[value] = struct{}{}
		}
	}
}

func validateBoundedStrings(c *validationCollector, path string, values []string, validator func(string) bool, invalidCode string) {
	if len(values) > MaxDeclarationItems {
		c.add(path, "manifest.list.too_many_items")
		return
	}
	seen := make(map[string]struct{}, len(values))
	for index, value := range values {
		itemPath := fmt.Sprintf("%s[%d]", path, index)
		if !validator(value) {
			c.add(itemPath, invalidCode)
		}
		if _, exists := seen[value]; exists {
			c.add(itemPath, "manifest.declaration.duplicate")
		} else {
			seen[value] = struct{}{}
		}
	}
}

func validSecretReference(value string) bool {
	if len(value) == 0 || len(value) > maxReferenceLength {
		return false
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "secret" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	return parsed.Path != "" && parsed.Path != "/"
}

func validNetworkDestination(value string) bool {
	if len(value) == 0 || len(value) > maxReferenceLength {
		return false
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	return parsed.Path == "" || parsed.Path == "/"
}

func validEndpoint(value string) bool {
	if len(value) == 0 || len(value) > maxReferenceLength || !strings.HasPrefix(value, "/") || strings.ContainsAny(value, "?#") {
		return false
	}
	return !strings.Contains(value, "//")
}

func validBoundedPattern(value string, pattern *regexp.Regexp, maxLength int) bool {
	return len(value) > 0 && len(value) <= maxLength && pattern.MatchString(value)
}

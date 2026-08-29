package modules

import (
	"regexp"
	"sort"
	"strconv"
)

var migrationOwnershipNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

// MigrationChangeClass makes compatibility, backfill and destructive intent
// explicit before a migration reaches the retained P01 execution boundary.
type MigrationChangeClass string

const (
	MigrationCompatible MigrationChangeClass = "compatible"
	MigrationBackfill   MigrationChangeClass = "backfill"
	MigrationDestructive MigrationChangeClass = "destructive"
)

// MigrationRegistration is declarative P03.09 metadata only. It contains no
// SQL, filesystem path, callback, tenant scope or execution authority.
type MigrationRegistration struct {
	ModuleID             string               `json:"module_id"`
	ModuleVersion        string               `json:"module_version"`
	Owner                string               `json:"owner"`
	Declaration          string               `json:"declaration"`
	Version              int64                `json:"version"`
	Name                 string               `json:"name"`
	IntroducedInVersion  string               `json:"introduced_in_version"`
	TargetOwner          string               `json:"target_owner"`
	ChangeClass          MigrationChangeClass `json:"change_class"`
	StrategyRef          string               `json:"strategy_ref,omitempty"`
	RecoveryRef          string               `json:"recovery_ref,omitempty"`
}

// MigrationRecord is the immutable-by-copy normalized identity used to plan
// fresh installs and supported upgrades. Owner+Version maps directly to the P01
// schema_migrations ledger key; P03.09 never applies the migration itself.
type MigrationRecord struct {
	ID                   string               `json:"id"`
	ModuleID             string               `json:"module_id"`
	ModuleVersion        string               `json:"module_version"`
	Owner                string               `json:"owner"`
	Declaration          string               `json:"declaration"`
	Version              int64                `json:"version"`
	Name                 string               `json:"name"`
	IntroducedInVersion  string               `json:"introduced_in_version"`
	TargetOwner          string               `json:"target_owner"`
	ChangeClass          MigrationChangeClass `json:"change_class"`
	StrategyRef          string               `json:"strategy_ref,omitempty"`
	RecoveryRef          string               `json:"recovery_ref,omitempty"`
}

// MigrationOwnershipDiagnostic is stable and value-bounded. It intentionally
// excludes raw SQL, filesystem paths, secrets, tenant IDs and database details.
type MigrationOwnershipDiagnostic struct {
	Code        string `json:"code"`
	ModuleID    string `json:"module_id,omitempty"`
	Declaration string `json:"declaration,omitempty"`
	Version     int64  `json:"version,omitempty"`
}

// MigrationOwnershipError is the fail-closed P03.09 registry/planning error.
type MigrationOwnershipError struct {
	Diagnostic MigrationOwnershipDiagnostic
}

func (e *MigrationOwnershipError) Error() string { return "module migration ownership registry failed" }

// MigrationOwnershipRegistry binds migrations to accepted module identity and
// authoritative owner metadata. It is a planner only; P01 remains the sole SQL
// runner, immutable ledger and owner advisory-lock authority.
type MigrationOwnershipRegistry struct {
	records               []MigrationRecord
	byID                  map[string]MigrationRecord
	byModule              map[string][]MigrationRecord
	moduleCurrentVersions map[string]string
}

// BindMigrationOwnershipRegistry creates a deterministic registry from explicit
// reviewed metadata and the already-validated P03 module registry.
func BindMigrationOwnershipRegistry(
	moduleRegistry Registry,
	registrations []MigrationRegistration,
) (*MigrationOwnershipRegistry, error) {
	moduleCurrentVersions := make(map[string]string)
	for _, module := range moduleRegistry.List() {
		if !validBoundedPattern(module.ID, moduleIDPattern, maxIdentifierLength) ||
			!validBoundedPattern(module.Owner, ownerPattern, maxIdentifierLength) {
			return nil, migrationOwnershipError("module.migration.registry_module_invalid", module.ID, "", 0)
		}
		if _, ok := parseStrictSemVer(module.Version); !ok {
			return nil, migrationOwnershipError("module.migration.registry_module_invalid", module.ID, "", 0)
		}
		moduleCurrentVersions[module.ID] = module.Version
	}

	records := make([]MigrationRecord, 0, len(registrations))
	byID := make(map[string]MigrationRecord, len(registrations))
	byDeclaration := make(map[string]struct{}, len(registrations))
	byOwnerVersion := make(map[string]MigrationRecord, len(registrations))

	for _, registration := range registrations {
		module, ok := moduleRegistry.Lookup(registration.ModuleID)
		if !ok {
			return nil, migrationOwnershipError("module.migration.module_missing", registration.ModuleID, registration.Declaration, registration.Version)
		}
		if registration.ModuleVersion != module.Version {
			return nil, migrationOwnershipError("module.migration.module_version_mismatch", module.ID, registration.Declaration, registration.Version)
		}
		if registration.Owner != module.Owner {
			return nil, migrationOwnershipError("module.migration.owner_mismatch", module.ID, registration.Declaration, registration.Version)
		}
		if registration.TargetOwner != module.Owner {
			return nil, migrationOwnershipError("module.migration.target_owner_mismatch", module.ID, registration.Declaration, registration.Version)
		}
		if !validBoundedPattern(registration.ModuleID, moduleIDPattern, maxIdentifierLength) ||
			!validBoundedPattern(registration.Owner, ownerPattern, maxIdentifierLength) ||
			!validBoundedPattern(registration.TargetOwner, ownerPattern, maxIdentifierLength) ||
			!validBoundedPattern(registration.Declaration, identifierPattern, maxIdentifierLength) ||
			registration.Version <= 0 ||
			!migrationOwnershipNamePattern.MatchString(registration.Name) {
			return nil, migrationOwnershipError("module.migration.registration_invalid", module.ID, registration.Declaration, registration.Version)
		}

		currentVersion, currentOK := parseStrictSemVer(module.Version)
		introducedVersion, introducedOK := parseStrictSemVer(registration.IntroducedInVersion)
		if !currentOK || !introducedOK {
			return nil, migrationOwnershipError("module.migration.introduced_version_invalid", module.ID, registration.Declaration, registration.Version)
		}
		if compareStrictSemVer(introducedVersion, currentVersion) > 0 {
			return nil, migrationOwnershipError("module.migration.introduced_version_future", module.ID, registration.Declaration, registration.Version)
		}
		if !knownMigrationChangeClass(registration.ChangeClass) {
			return nil, migrationOwnershipError("module.migration.change_class_invalid", module.ID, registration.Declaration, registration.Version)
		}
		if err := validateMigrationStrategy(registration); err != nil {
			return nil, err
		}

		declarationKey := registration.ModuleID + "\x00" + registration.Declaration
		if _, exists := byDeclaration[declarationKey]; exists {
			return nil, migrationOwnershipError("module.migration.registration_duplicate", module.ID, registration.Declaration, registration.Version)
		}
		byDeclaration[declarationKey] = struct{}{}

		ownerVersionKey := registration.Owner + "\x00" + strconv.FormatInt(registration.Version, 10)
		if existing, exists := byOwnerVersion[ownerVersionKey]; exists {
			return nil, migrationOwnershipError("module.migration.order_conflict", existing.ModuleID, existing.Declaration, existing.Version)
		}

		record := MigrationRecord{
			ID:                  migrationOwnershipID(registration),
			ModuleID:            registration.ModuleID,
			ModuleVersion:       registration.ModuleVersion,
			Owner:               registration.Owner,
			Declaration:         registration.Declaration,
			Version:             registration.Version,
			Name:                registration.Name,
			IntroducedInVersion: registration.IntroducedInVersion,
			TargetOwner:         registration.TargetOwner,
			ChangeClass:         registration.ChangeClass,
			StrategyRef:         registration.StrategyRef,
			RecoveryRef:         registration.RecoveryRef,
		}
		if _, exists := byID[record.ID]; exists {
			return nil, migrationOwnershipError("module.migration.identity_duplicate", module.ID, registration.Declaration, registration.Version)
		}
		byID[record.ID] = record
		byOwnerVersion[ownerVersionKey] = record
		records = append(records, record)
	}

	sortMigrationRecords(records)
	byModule := make(map[string][]MigrationRecord)
	for _, record := range records {
		byModule[record.ModuleID] = append(byModule[record.ModuleID], record)
	}
	for moduleID, moduleRecords := range byModule {
		if err := validateMigrationVersionProgression(moduleID, moduleRecords); err != nil {
			return nil, err
		}
	}

	return &MigrationOwnershipRegistry{
		records:               append([]MigrationRecord(nil), records...),
		byID:                  byID,
		byModule:              cloneMigrationRecordMap(byModule),
		moduleCurrentVersions: moduleCurrentVersions,
	}, nil
}

// List returns the complete deterministic owner/version ordered metadata set.
func (r *MigrationOwnershipRegistry) List() []MigrationRecord {
	if r == nil || len(r.records) == 0 {
		return []MigrationRecord{}
	}
	return append([]MigrationRecord(nil), r.records...)
}

// Lookup resolves one exact module-scoped declaration without granting apply authority.
func (r *MigrationOwnershipRegistry) Lookup(moduleID, declaration string) (MigrationRecord, error) {
	if r == nil || !validBoundedPattern(moduleID, moduleIDPattern, maxIdentifierLength) ||
		!validBoundedPattern(declaration, identifierPattern, maxIdentifierLength) {
		return MigrationRecord{}, migrationOwnershipError("module.migration.query_invalid", moduleID, declaration, 0)
	}
	for _, record := range r.byModule[moduleID] {
		if record.Declaration == declaration {
			return record, nil
		}
	}
	return MigrationRecord{}, migrationOwnershipError("module.migration.record_missing", moduleID, declaration, 0)
}

// FreshInstallPlan returns every migration owned by one module in immutable P01
// ledger order. It plans metadata only; callers still execute through P01.
func (r *MigrationOwnershipRegistry) FreshInstallPlan(moduleID string) ([]MigrationRecord, error) {
	if err := r.validateModuleQuery(moduleID); err != nil {
		return nil, err
	}
	return append([]MigrationRecord(nil), r.byModule[moduleID]...), nil
}

// UpgradePlan returns migrations introduced after a supported source module
// version through the currently discovered module version. Future/downgrade
// targets fail closed; execution remains the P01 migrator's responsibility.
func (r *MigrationOwnershipRegistry) UpgradePlan(moduleID, fromModuleVersion string) ([]MigrationRecord, error) {
	if err := r.validateModuleQuery(moduleID); err != nil {
		return nil, err
	}
	currentText := r.moduleCurrentVersions[moduleID]
	current, currentOK := parseStrictSemVer(currentText)
	from, fromOK := parseStrictSemVer(fromModuleVersion)
	if !currentOK || !fromOK {
		return nil, migrationOwnershipError("module.migration.upgrade_version_invalid", moduleID, "", 0)
	}
	if compareStrictSemVer(from, current) > 0 {
		return nil, migrationOwnershipError("module.migration.upgrade_future_source", moduleID, "", 0)
	}
	if compareStrictSemVer(from, current) == 0 {
		return []MigrationRecord{}, nil
	}

	plan := make([]MigrationRecord, 0)
	for _, record := range r.byModule[moduleID] {
		introduced, ok := parseStrictSemVer(record.IntroducedInVersion)
		if !ok {
			return nil, migrationOwnershipError("module.migration.registry_invalid", moduleID, record.Declaration, record.Version)
		}
		if compareStrictSemVer(introduced, from) > 0 && compareStrictSemVer(introduced, current) <= 0 {
			plan = append(plan, record)
		}
	}
	return plan, nil
}

func (r *MigrationOwnershipRegistry) validateModuleQuery(moduleID string) error {
	if r == nil || !validBoundedPattern(moduleID, moduleIDPattern, maxIdentifierLength) {
		return migrationOwnershipError("module.migration.query_invalid", moduleID, "", 0)
	}
	if _, ok := r.moduleCurrentVersions[moduleID]; !ok {
		return migrationOwnershipError("module.migration.module_missing", moduleID, "", 0)
	}
	return nil
}

func validateMigrationStrategy(registration MigrationRegistration) error {
	strategyValid := registration.StrategyRef == "" || validBoundedPattern(registration.StrategyRef, identifierPattern, maxReferenceLength)
	recoveryValid := registration.RecoveryRef == "" || validBoundedPattern(registration.RecoveryRef, identifierPattern, maxReferenceLength)
	if !strategyValid || !recoveryValid {
		return migrationOwnershipError("module.migration.strategy_ref_invalid", registration.ModuleID, registration.Declaration, registration.Version)
	}
	if registration.ChangeClass == MigrationBackfill || registration.ChangeClass == MigrationDestructive {
		if registration.StrategyRef == "" || registration.RecoveryRef == "" {
			return migrationOwnershipError("module.migration.strategy_required", registration.ModuleID, registration.Declaration, registration.Version)
		}
	}
	return nil
}

func validateMigrationVersionProgression(moduleID string, records []MigrationRecord) error {
	var previous strictSemVer
	hasPrevious := false
	for _, record := range records {
		introduced, ok := parseStrictSemVer(record.IntroducedInVersion)
		if !ok {
			return migrationOwnershipError("module.migration.introduced_version_invalid", moduleID, record.Declaration, record.Version)
		}
		if hasPrevious && compareStrictSemVer(introduced, previous) < 0 {
			return migrationOwnershipError("module.migration.introduced_version_order_invalid", moduleID, record.Declaration, record.Version)
		}
		previous = introduced
		hasPrevious = true
	}
	return nil
}

func sortMigrationRecords(records []MigrationRecord) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].Owner != records[j].Owner {
			return records[i].Owner < records[j].Owner
		}
		if records[i].Version != records[j].Version {
			return records[i].Version < records[j].Version
		}
		if records[i].ModuleID != records[j].ModuleID {
			return records[i].ModuleID < records[j].ModuleID
		}
		return records[i].Declaration < records[j].Declaration
	})
}

func cloneMigrationRecordMap(source map[string][]MigrationRecord) map[string][]MigrationRecord {
	result := make(map[string][]MigrationRecord, len(source))
	for moduleID, records := range source {
		result[moduleID] = append([]MigrationRecord(nil), records...)
	}
	return result
}

func knownMigrationChangeClass(class MigrationChangeClass) bool {
	switch class {
	case MigrationCompatible, MigrationBackfill, MigrationDestructive:
		return true
	default:
		return false
	}
}

func migrationOwnershipID(registration MigrationRegistration) string {
	return registration.ModuleID + "@" + registration.ModuleVersion + ":" + registration.Owner + ":" +
		strconv.FormatInt(registration.Version, 10) + ":" + registration.Name
}

func migrationOwnershipError(code, moduleID, declaration string, version int64) error {
	return &MigrationOwnershipError{Diagnostic: MigrationOwnershipDiagnostic{
		Code:        code,
		ModuleID:    moduleID,
		Declaration: declaration,
		Version:     version,
	}}
}

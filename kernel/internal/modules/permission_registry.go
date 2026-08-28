package modules

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/authorization"
)

var reservedPermissionNamespaces = map[string]struct{}{
	"authorization": {},
	"configuration": {},
	"kernel":        {},
	"platform":      {},
}

// PermissionRegistration binds one validated manifest permission declaration to
// its module/owner identity and optional descriptive capability association. It
// does not create a role, assignment, tenant scope, or authorization grant.
type PermissionRegistration struct {
	ModuleID   string `json:"module_id"`
	Owner      string `json:"owner"`
	Permission string `json:"permission"`
	Capability string `json:"capability,omitempty"`
}

// PermissionRecord is the deterministic lifecycle-aware metadata projection for
// one module permission. Available is true only while the owning module is
// enabled; availability is a precondition, never a grant.
type PermissionRecord struct {
	Permission     authorization.PermissionID `json:"permission"`
	ModuleID       string                     `json:"module_id"`
	Owner          string                     `json:"owner"`
	Capability     string                     `json:"capability,omitempty"`
	LifecycleState LifecycleState             `json:"lifecycle_state"`
	Available      bool                       `json:"available"`
}

// PermissionDiagnostic is stable and value-safe. It excludes raw manifests,
// role/principal state, tenant identifiers, secrets and authorization results.
type PermissionDiagnostic struct {
	Code       string `json:"code"`
	ModuleID   string `json:"module_id,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// PermissionError is a fail-closed P03.07 registration/availability failure.
type PermissionError struct {
	Diagnostic PermissionDiagnostic
}

func (e *PermissionError) Error() string { return "module permission registration failed" }

type permissionDeclaration struct {
	permission authorization.PermissionID
	moduleID   string
	owner      string
	capability string
}

// PermissionRegistry consumes only validated discovery snapshots plus accepted
// lifecycle state. It contains no roles, assignments, authorization decisions or
// raw tenant/organization identifiers.
type PermissionRegistry struct {
	declarations   []permissionDeclaration
	byPermission   map[authorization.PermissionID]permissionDeclaration
	lifecycleStore LifecycleStore
}

// BindPermissionRegistry constructs a deterministic P03.07 registry. Raw
// manifests are never reparsed and every declared permission requires exactly
// one explicit owner/module registration.
func BindPermissionRegistry(
	registry Registry,
	lifecycleStore LifecycleStore,
	registrations []PermissionRegistration,
) (*PermissionRegistry, error) {
	if lifecycleStore == nil {
		return nil, permissionError("module.permission.lifecycle_store_required", "", "")
	}

	declared, err := collectPermissionDeclarations(registry)
	if err != nil {
		return nil, err
	}

	registrationsByPermission := make(map[authorization.PermissionID]PermissionRegistration, len(registrations))
	for _, registration := range registrations {
		permission := authorization.PermissionID(registration.Permission)
		if !permission.Valid() {
			return nil, permissionError("module.permission.registration_invalid", registration.ModuleID, "")
		}
		if _, exists := registrationsByPermission[permission]; exists {
			return nil, permissionError("module.permission.registration_duplicate", registration.ModuleID, registration.Permission)
		}
		declaration, exists := declared[permission]
		if !exists {
			return nil, permissionError("module.permission.registration_undeclared", registration.ModuleID, registration.Permission)
		}
		if registration.ModuleID != declaration.moduleID {
			return nil, permissionError("module.permission.module_mismatch", declaration.moduleID, registration.Permission)
		}
		if registration.Owner != declaration.owner {
			return nil, permissionError("module.permission.owner_mismatch", declaration.moduleID, registration.Permission)
		}
		if registration.Capability != "" && !moduleDeclaresCapability(registry, declaration.moduleID, registration.Capability) {
			return nil, permissionError("module.permission.capability_undeclared", declaration.moduleID, registration.Permission)
		}
		registration.Capability = strings.TrimSpace(registration.Capability)
		registrationsByPermission[permission] = registration
	}

	permissions := make([]authorization.PermissionID, 0, len(declared))
	for permission := range declared {
		permissions = append(permissions, permission)
	}
	sort.Slice(permissions, func(i, j int) bool { return permissions[i] < permissions[j] })

	declarations := make([]permissionDeclaration, 0, len(permissions))
	byPermission := make(map[authorization.PermissionID]permissionDeclaration, len(permissions))
	for _, permission := range permissions {
		declaration := declared[permission]
		registration, ok := registrationsByPermission[permission]
		if !ok {
			return nil, permissionError("module.permission.registration_missing", declaration.moduleID, string(permission))
		}
		declaration.capability = registration.Capability
		declarations = append(declarations, declaration)
		byPermission[permission] = declaration
	}

	return &PermissionRegistry{
		declarations:   declarations,
		byPermission:   byPermission,
		lifecycleStore: lifecycleStore,
	}, nil
}

// List returns all declared module permissions in deterministic order. Identity
// remains visible while unavailable so policy/history references are not erased.
func (r *PermissionRegistry) List(ctx context.Context) ([]PermissionRecord, error) {
	if r == nil || r.lifecycleStore == nil {
		return nil, permissionError("module.permission.registry_invalid", "", "")
	}
	result := make([]PermissionRecord, 0, len(r.declarations))
	for _, declaration := range r.declarations {
		record, err := r.resolve(ctx, declaration)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

// Lookup returns one lifecycle-aware metadata record. Lookup itself is not an
// authorization decision and never creates invocation authority.
func (r *PermissionRegistry) Lookup(ctx context.Context, permission authorization.PermissionID) (PermissionRecord, error) {
	if r == nil || r.lifecycleStore == nil || !permission.Valid() {
		return PermissionRecord{}, permissionError("module.permission.query_invalid", "", string(permission))
	}
	declaration, ok := r.byPermission[permission]
	if !ok {
		return PermissionRecord{}, permissionError("module.permission.unknown", "", string(permission))
	}
	return r.resolve(ctx, declaration)
}

// PermissionAvailable implements authorization.ModulePermissionAvailability.
// Unknown or non-enabled module permissions return false, preserving deny-by-
// default behavior in kernel.authorization.
func (r *PermissionRegistry) PermissionAvailable(
	ctx context.Context,
	permission authorization.PermissionID,
) (bool, error) {
	if r == nil || r.lifecycleStore == nil || !permission.Valid() {
		return false, nil
	}
	declaration, ok := r.byPermission[permission]
	if !ok {
		return false, nil
	}
	record, err := r.resolve(ctx, declaration)
	if err != nil {
		return false, err
	}
	return record.Available, nil
}

// SynchronizeCatalog writes only authorization-owned permission reference
// metadata. It never writes roles/assignments and never deletes permission IDs.
func (r *PermissionRegistry) SynchronizeCatalog(
	ctx context.Context,
	catalog authorization.ModulePermissionCatalog,
	observedAt time.Time,
) error {
	if r == nil || catalog == nil || observedAt.IsZero() {
		return permissionError("module.permission.catalog_invalid", "", "")
	}
	records, err := r.List(ctx)
	if err != nil {
		return err
	}
	definitions := make([]authorization.ModulePermissionDefinition, 0, len(records))
	for _, record := range records {
		definitions = append(definitions, authorization.ModulePermissionDefinition{
			Permission: record.Permission,
			ModuleID:   record.ModuleID,
			Owner:      record.Owner,
			Capability: record.Capability,
			Available:  record.Available,
		})
	}
	return catalog.ReconcileModulePermissions(ctx, definitions, observedAt)
}

func (r *PermissionRegistry) resolve(
	ctx context.Context,
	declaration permissionDeclaration,
) (PermissionRecord, error) {
	lifecycle, found, err := r.lifecycleStore.Load(ctx, declaration.moduleID)
	if err != nil {
		return PermissionRecord{}, permissionError("module.permission.lifecycle_read_failed", declaration.moduleID, string(declaration.permission))
	}
	state := LifecycleAvailable
	if found {
		if lifecycle.ModuleID != declaration.moduleID {
			return PermissionRecord{}, permissionError("module.permission.lifecycle_identity_mismatch", declaration.moduleID, string(declaration.permission))
		}
		if !knownCapabilityLifecycleState(lifecycle.State) {
			return PermissionRecord{}, permissionError("module.permission.lifecycle_state_invalid", declaration.moduleID, string(declaration.permission))
		}
		state = lifecycle.State
	}
	return PermissionRecord{
		Permission:     declaration.permission,
		ModuleID:       declaration.moduleID,
		Owner:          declaration.owner,
		Capability:     declaration.capability,
		LifecycleState: state,
		Available:      state == LifecycleEnabled,
	}, nil
}

func collectPermissionDeclarations(registry Registry) (map[authorization.PermissionID]permissionDeclaration, error) {
	result := make(map[authorization.PermissionID]permissionDeclaration)
	for _, record := range registry.List() {
		snapshot, ok := registry.manifestSnapshot(record.ID)
		if !ok || snapshot.ID != record.ID || snapshot.Owner != record.Owner || snapshot.Version != record.Version {
			return nil, permissionError("module.permission.snapshot_invalid", record.ID, "")
		}
		for _, raw := range snapshot.Permissions {
			permission := authorization.PermissionID(raw)
			if !permission.Valid() {
				return nil, permissionError("module.permission.declaration_invalid", record.ID, "")
			}
			if !permissionNamespaceAllowed(record.ID, raw) {
				return nil, permissionError("module.permission.namespace_invalid", record.ID, raw)
			}
			if existing, exists := result[permission]; exists {
				return nil, permissionError("module.permission.declaration_collision", existing.moduleID, raw)
			}
			result[permission] = permissionDeclaration{
				permission: permission,
				moduleID:   record.ID,
				owner:      record.Owner,
			}
		}
	}
	return result, nil
}

func permissionNamespaceAllowed(moduleID, permission string) bool {
	parts := strings.Split(permission, ".")
	if len(parts) < 3 {
		return false
	}
	if _, reserved := reservedPermissionNamespaces[parts[0]]; reserved {
		return false
	}
	moduleNamespace := strings.TrimPrefix(moduleID, "omnexa.")
	moduleNamespace = strings.Split(moduleNamespace, ".")[0]
	moduleNamespace = strings.Split(moduleNamespace, "_")[0]
	moduleNamespace = strings.Split(moduleNamespace, "-")[0]
	return moduleNamespace != "" && parts[0] == moduleNamespace
}

func moduleDeclaresCapability(registry Registry, moduleID, capability string) bool {
	identity, valid := parseCapabilityDeclaration(capability)
	if !valid {
		return false
	}
	snapshot, ok := registry.manifestSnapshot(moduleID)
	if !ok {
		return false
	}
	for _, declared := range snapshot.CapabilitiesProvided {
		parsed, parsedOK := parseCapabilityDeclaration(declared)
		if parsedOK && parsed.ID == identity.ID && parsed.Major == identity.Major && parsed.Declaration == identity.Declaration {
			return true
		}
	}
	return false
}

func permissionError(code, moduleID, permission string) error {
	return &PermissionError{Diagnostic: PermissionDiagnostic{
		Code: code, ModuleID: moduleID, Permission: permission,
	}}
}

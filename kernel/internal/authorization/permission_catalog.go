package authorization

import (
	"context"
	"regexp"
	"sort"
	"time"
)

const (
	permissionSourceKernel = "kernel"
	permissionSourceModule = "module"
)

var (
	modulePermissionModuleIDPattern = regexp.MustCompile(`^omnexa\.[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$`)
	modulePermissionOwnerPattern    = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$`)
)

// ModulePermissionDefinition is authorization-owned reference metadata for one
// module-declared permission. It is not a role/principal grant and contains no
// tenant or organization scope.
type ModulePermissionDefinition struct {
	Permission PermissionID
	ModuleID   string
	Owner      string
	Capability string
	Available  bool
}

// ModulePermissionAvailability is the live P03.07 precondition consulted before
// evaluating a module permission through the existing deny-by-default RBAC
// repository. It grants no authority by itself.
type ModulePermissionAvailability interface {
	PermissionAvailable(context.Context, PermissionID) (bool, error)
}

// ModulePermissionCatalog is the non-destructive authorization-owned catalog
// synchronization boundary used by the module runtime. Implementations must
// preserve permission identity/history rather than delete referenced rows.
type ModulePermissionCatalog interface {
	ReconcileModulePermissions(context.Context, []ModulePermissionDefinition, time.Time) error
}

// ReconcileModulePermissions synchronizes module permission reference metadata
// without mutating roles or assignments. Missing declarations are marked
// unavailable rather than deleted so existing policy/audit references survive.
func (repository *PostgresRepository) ReconcileModulePermissions(
	ctx context.Context,
	definitions []ModulePermissionDefinition,
	observedAt time.Time,
) error {
	if repository == nil || repository.pool == nil || observedAt.IsZero() {
		return repositoryInvalidFailure()
	}

	normalized := append([]ModulePermissionDefinition(nil), definitions...)
	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].Permission < normalized[j].Permission
	})
	seen := make(map[PermissionID]struct{}, len(normalized))
	for _, definition := range normalized {
		if !validModulePermissionDefinition(definition) || kernelPermission(definition.Permission) {
			return invalidPermissionFailure()
		}
		if _, exists := seen[definition.Permission]; exists {
			return invalidPermissionFailure()
		}
		seen[definition.Permission] = struct{}{}
	}

	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return repositoryFailure(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	instant := observedAt.UTC()
	activeIDs := make([]string, 0, len(normalized))
	for _, definition := range normalized {
		activeIDs = append(activeIDs, string(definition.Permission))
		command, execErr := tx.Exec(
			ctx,
			`INSERT INTO omnexa_authorization.permissions
			 (permission_id, created_at, source_kind, module_id, owner_name, capability_ref, available, updated_at)
			 VALUES ($1, $2, 'module', $3, $4, NULLIF($5, ''), $6, $2)
			 ON CONFLICT (permission_id) DO UPDATE
			 SET capability_ref = EXCLUDED.capability_ref,
			     available = EXCLUDED.available,
			     updated_at = EXCLUDED.updated_at
			 WHERE omnexa_authorization.permissions.source_kind = 'module'
			   AND omnexa_authorization.permissions.module_id = EXCLUDED.module_id
			   AND omnexa_authorization.permissions.owner_name = EXCLUDED.owner_name`,
			string(definition.Permission), instant, definition.ModuleID, definition.Owner,
			definition.Capability, definition.Available,
		)
		if execErr != nil {
			return permissionCatalogPersistenceFailure(execErr)
		}
		if command.RowsAffected() != 1 {
			return permissionCatalogConflictFailure()
		}
	}

	if _, err = tx.Exec(
		ctx,
		`UPDATE omnexa_authorization.permissions
		 SET available = FALSE, updated_at = $2
		 WHERE source_kind = 'module'
		   AND NOT (permission_id = ANY($1::text[]))`,
		activeIDs, instant,
	); err != nil {
		return permissionCatalogPersistenceFailure(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return repositoryFailure(err)
	}
	return nil
}

func validModulePermissionDefinition(definition ModulePermissionDefinition) bool {
	if !definition.Permission.Valid() ||
		!modulePermissionModuleIDPattern.MatchString(definition.ModuleID) ||
		!modulePermissionOwnerPattern.MatchString(definition.Owner) {
		return false
	}
	return definition.Capability == "" || validReference(definition.Capability, 128)
}

func kernelPermission(permission PermissionID) bool {
	if BuiltinPermission(permission) {
		return true
	}
	switch permission {
	case PermissionID("configuration.setting.read"), PermissionID("configuration.setting.manage"):
		return true
	default:
		return false
	}
}

func (service *Service) permissionAvailable(ctx context.Context, permission PermissionID) (bool, error) {
	if kernelPermission(permission) {
		return true, nil
	}
	if service == nil || service.modulePermissions == nil {
		return false, nil
	}
	return service.modulePermissions.PermissionAvailable(ctx, permission)
}

func permissionCatalogPersistenceFailure(cause error) error {
	return repositoryFailure(cause)
}

func permissionCatalogConflictFailure() error {
	return invalidPermissionFailure()
}

package authorization

import (
	"context"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository persists only kernel.authorization-owned direct RBAC state.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates the P02.05 owner-bounded PostgreSQL repository.
func NewPostgresRepository(pool *pgxpool.Pool) (*PostgresRepository, error) {
	if pool == nil {
		return nil, repositoryInvalidFailure()
	}
	return &PostgresRepository{pool: pool}, nil
}

func (repository *PostgresRepository) permissionsExist(ctx context.Context, permissions []PermissionID) (bool, error) {
	if repository == nil || repository.pool == nil {
		return false, repositoryInvalidFailure()
	}
	if len(permissions) == 0 || len(permissions) > maxRolePermissions {
		return false, invalidPermissionFailure()
	}
	values := make([]string, 0, len(permissions))
	seen := make(map[PermissionID]struct{}, len(permissions))
	for _, permission := range permissions {
		if !permission.Valid() {
			return false, invalidPermissionFailure()
		}
		if _, exists := seen[permission]; exists {
			continue
		}
		seen[permission] = struct{}{}
		values = append(values, string(permission))
	}
	var count int
	if err := repository.pool.QueryRow(
		ctx,
		`SELECT count(*) FROM omnexa_authorization.permissions WHERE permission_id = ANY($1)`,
		values,
	).Scan(&count); err != nil {
		return false, repositoryFailure(err)
	}
	return count == len(values), nil
}

func (repository *PostgresRepository) createRole(ctx context.Context, role Role) error {
	if repository == nil || repository.pool == nil {
		return repositoryInvalidFailure()
	}
	if err := role.validate(); err != nil {
		return err
	}
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return repositoryFailure(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err = tx.Exec(
		ctx,
		`INSERT INTO omnexa_authorization.roles
		 (id, tenant_id, organization_id, name, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		string(role.ID()), string(role.Scope().TenantID()), organizationArgument(role.Scope()),
		role.Name(), role.CreatedAt(), role.UpdatedAt(),
	); err != nil {
		return rolePersistenceFailure(err)
	}
	for _, permission := range role.Permissions() {
		if _, err = tx.Exec(
			ctx,
			`INSERT INTO omnexa_authorization.role_permissions (role_id, permission_id)
			 VALUES ($1, $2)`,
			string(role.ID()), string(permission),
		); err != nil {
			return rolePersistenceFailure(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return repositoryFailure(err)
	}
	return nil
}

func (repository *PostgresRepository) getRole(ctx context.Context, scope Scope, id RoleID) (Role, error) {
	if repository == nil || repository.pool == nil {
		return Role{}, repositoryInvalidFailure()
	}
	if !scope.Valid() || !id.Valid() {
		return Role{}, invalidRoleFailure()
	}

	var storedID string
	var tenantID string
	var organizationID *string
	var name string
	var createdAt time.Time
	var updatedAt time.Time
	queryErr := repository.pool.QueryRow(
		ctx,
		`SELECT id, tenant_id, organization_id, name, created_at, updated_at
		 FROM omnexa_authorization.roles
		 WHERE id = $1
		   AND tenant_id = $2
		   AND organization_id IS NOT DISTINCT FROM $3::uuid`,
		string(id), string(scope.TenantID()), organizationArgument(scope),
	).Scan(&storedID, &tenantID, &organizationID, &name, &createdAt, &updatedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return Role{}, roleNotFoundFailure()
	}
	if queryErr != nil {
		return Role{}, repositoryFailure(queryErr)
	}
	storedScope, scopeErr := scopeFromStored(tenantID, organizationID)
	if scopeErr != nil {
		return Role{}, invalidStoredRoleFailure()
	}
	permissions, permissionsErr := repository.rolePermissions(ctx, RoleID(storedID))
	if permissionsErr != nil {
		return Role{}, permissionsErr
	}
	return rehydrateRole(RoleID(storedID), storedScope, name, permissions, createdAt, updatedAt)
}

func (repository *PostgresRepository) replaceRolePermissions(
	ctx context.Context,
	scope Scope,
	id RoleID,
	permissions []PermissionID,
	changedAt time.Time,
) (Role, error) {
	if repository == nil || repository.pool == nil {
		return Role{}, repositoryInvalidFailure()
	}
	if !scope.Valid() || !id.Valid() || changedAt.IsZero() {
		return Role{}, invalidRoleFailure()
	}
	normalized, err := normalizePermissions(permissions)
	if err != nil {
		return Role{}, err
	}
	instant := changedAt.UTC()
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return Role{}, repositoryFailure(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	command, err := tx.Exec(
		ctx,
		`UPDATE omnexa_authorization.roles
		 SET updated_at = $4
		 WHERE id = $1
		   AND tenant_id = $2
		   AND organization_id IS NOT DISTINCT FROM $3::uuid
		   AND updated_at <= $4`,
		string(id), string(scope.TenantID()), organizationArgument(scope), instant,
	)
	if err != nil {
		return Role{}, rolePersistenceFailure(err)
	}
	if command.RowsAffected() != 1 {
		return Role{}, roleConflictFailure()
	}
	if _, err = tx.Exec(ctx, `DELETE FROM omnexa_authorization.role_permissions WHERE role_id = $1`, string(id)); err != nil {
		return Role{}, rolePersistenceFailure(err)
	}
	for _, permission := range normalized {
		if _, err = tx.Exec(
			ctx,
			`INSERT INTO omnexa_authorization.role_permissions (role_id, permission_id) VALUES ($1, $2)`,
			string(id), string(permission),
		); err != nil {
			return Role{}, rolePersistenceFailure(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return Role{}, repositoryFailure(err)
	}
	return repository.getRole(ctx, scope, id)
}

func (repository *PostgresRepository) createAssignment(ctx context.Context, assignment Assignment) error {
	if repository == nil || repository.pool == nil {
		return repositoryInvalidFailure()
	}
	if err := assignment.validate(); err != nil {
		return err
	}
	role, err := repository.getRole(ctx, assignment.Scope(), assignment.RoleID())
	if err != nil {
		return err
	}
	if !role.Scope().Equal(assignment.Scope()) {
		return scopeDeniedFailure()
	}
	_, err = repository.pool.Exec(
		ctx,
		`INSERT INTO omnexa_authorization.role_assignments
		 (id, role_id, principal_id, assignment_state, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		string(assignment.ID()), string(assignment.RoleID()), string(assignment.PrincipalID()),
		string(assignment.State()), assignment.CreatedAt(), assignment.UpdatedAt(),
	)
	if err != nil {
		return assignmentPersistenceFailure(err)
	}
	return nil
}

func (repository *PostgresRepository) getAssignment(ctx context.Context, scope Scope, id AssignmentID) (Assignment, error) {
	if repository == nil || repository.pool == nil {
		return Assignment{}, repositoryInvalidFailure()
	}
	if !scope.Valid() || !id.Valid() {
		return Assignment{}, invalidAssignmentFailure()
	}

	var storedID string
	var roleID string
	var principalID string
	var state string
	var createdAt time.Time
	var updatedAt time.Time
	var tenantID string
	var organizationID *string
	queryErr := repository.pool.QueryRow(
		ctx,
		`SELECT a.id, a.role_id, a.principal_id, a.assignment_state, a.created_at, a.updated_at,
		        r.tenant_id, r.organization_id
		 FROM omnexa_authorization.role_assignments AS a
		 JOIN omnexa_authorization.roles AS r ON r.id = a.role_id
		 JOIN omnexa_identity.principals AS p ON p.id = a.principal_id
		 WHERE a.id = $1
		   AND p.principal_type = 'human_user'
		   AND r.tenant_id = $2
		   AND r.organization_id IS NOT DISTINCT FROM $3::uuid`,
		string(id), string(scope.TenantID()), organizationArgument(scope),
	).Scan(&storedID, &roleID, &principalID, &state, &createdAt, &updatedAt, &tenantID, &organizationID)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return Assignment{}, assignmentNotFoundFailure()
	}
	if queryErr != nil {
		return Assignment{}, repositoryFailure(queryErr)
	}
	storedScope, scopeErr := scopeFromStored(tenantID, organizationID)
	if scopeErr != nil {
		return Assignment{}, invalidStoredAssignmentFailure()
	}
	assignment, assignmentErr := newAssignmentAt(
		AssignmentID(storedID), RoleID(roleID), identity.UserID(principalID), storedScope,
		AssignmentState(state), createdAt, updatedAt,
	)
	if assignmentErr != nil {
		return Assignment{}, invalidStoredAssignmentFailure()
	}
	return assignment, nil
}

func (repository *PostgresRepository) revokeAssignment(
	ctx context.Context,
	scope Scope,
	id AssignmentID,
	changedAt time.Time,
) (Assignment, error) {
	if repository == nil || repository.pool == nil {
		return Assignment{}, repositoryInvalidFailure()
	}
	if !scope.Valid() || !id.Valid() || changedAt.IsZero() {
		return Assignment{}, invalidAssignmentFailure()
	}
	instant := changedAt.UTC()

	var storedID string
	var roleID string
	var principalID string
	var state string
	var createdAt time.Time
	var updatedAt time.Time
	queryErr := repository.pool.QueryRow(
		ctx,
		`UPDATE omnexa_authorization.role_assignments AS a
		 SET assignment_state = 'revoked', updated_at = $4
		 FROM omnexa_authorization.roles AS r, omnexa_identity.principals AS p
		 WHERE a.id = $1
		   AND a.role_id = r.id
		   AND p.id = a.principal_id
		   AND p.principal_type = 'human_user'
		   AND r.tenant_id = $2
		   AND r.organization_id IS NOT DISTINCT FROM $3::uuid
		   AND a.assignment_state = 'active'
		   AND a.updated_at <= $4
		 RETURNING a.id, a.role_id, a.principal_id, a.assignment_state, a.created_at, a.updated_at`,
		string(id), string(scope.TenantID()), organizationArgument(scope), instant,
	).Scan(&storedID, &roleID, &principalID, &state, &createdAt, &updatedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return Assignment{}, assignmentConflictFailure()
	}
	if queryErr != nil {
		return Assignment{}, assignmentPersistenceFailure(queryErr)
	}
	assignment, assignmentErr := newAssignmentAt(
		AssignmentID(storedID), RoleID(roleID), identity.UserID(principalID), scope,
		AssignmentState(state), createdAt, updatedAt,
	)
	if assignmentErr != nil {
		return Assignment{}, invalidStoredAssignmentFailure()
	}
	return assignment, nil
}

func (repository *PostgresRepository) hasPermission(ctx context.Context, subject Subject, permission PermissionID) (bool, error) {
	if repository == nil || repository.pool == nil {
		return false, repositoryInvalidFailure()
	}
	if !subject.Valid() || !permission.Valid() {
		return false, invalidSubjectFailure()
	}
	var allowed bool
	queryErr := repository.pool.QueryRow(
		ctx,
		`SELECT EXISTS (
		   SELECT 1
		   FROM omnexa_authorization.role_assignments AS a
		   JOIN omnexa_authorization.roles AS r ON r.id = a.role_id
		   JOIN omnexa_authorization.role_permissions AS rp ON rp.role_id = r.id
		   JOIN omnexa_identity.principals AS p ON p.id = a.principal_id
		   JOIN omnexa_tenancy.tenants AS t ON t.id = r.tenant_id
		   WHERE a.principal_id = $1
		     AND a.assignment_state = 'active'
		     AND p.principal_type = 'human_user'
		     AND p.lifecycle_state = 'active'
		     AND t.lifecycle_state = 'active'
		     AND r.tenant_id = $2
		     AND r.organization_id IS NOT DISTINCT FROM $3::uuid
		     AND rp.permission_id = $4
		 )`,
		string(subject.PrincipalID()), string(subject.Scope().TenantID()), organizationArgument(subject.Scope()), string(permission),
	).Scan(&allowed)
	if queryErr != nil {
		return false, repositoryFailure(queryErr)
	}
	return allowed, nil
}

func (repository *PostgresRepository) rolePermissions(ctx context.Context, roleID RoleID) ([]PermissionID, error) {
	rows, err := repository.pool.Query(
		ctx,
		`SELECT permission_id FROM omnexa_authorization.role_permissions WHERE role_id = $1 ORDER BY permission_id`,
		string(roleID),
	)
	if err != nil {
		return nil, repositoryFailure(err)
	}
	defer rows.Close()
	permissions := make([]PermissionID, 0, 8)
	for rows.Next() {
		var permission string
		if scanErr := rows.Scan(&permission); scanErr != nil {
			return nil, repositoryFailure(scanErr)
		}
		permissions = append(permissions, PermissionID(permission))
	}
	if rows.Err() != nil {
		return nil, repositoryFailure(rows.Err())
	}
	return permissions, nil
}

func organizationArgument(scope Scope) any {
	if scope.Kind() != ScopeOrganization {
		return nil
	}
	return string(scope.OrganizationID())
}

func scopeFromStored(tenantID string, organizationID *string) (Scope, error) {
	scope := Scope{tenantID: tenancy.TenantID(tenantID)}
	if organizationID == nil {
		scope.kind = ScopeTenant
	} else {
		scope.kind = ScopeOrganization
		scope.organizationID = organization.NodeID(*organizationID)
	}
	if !scope.Valid() {
		return Scope{}, scopeDeniedFailure()
	}
	return scope, nil
}

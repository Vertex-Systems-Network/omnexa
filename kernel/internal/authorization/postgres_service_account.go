package authorization

import (
	"context"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/jackc/pgx/v5"
)

func (repository *PostgresRepository) createServiceAccountAssignment(ctx context.Context, assignment ServiceAccountAssignment) error {
	if repository == nil || repository.pool == nil || assignment.validate() != nil {
		return invalidAssignmentFailure()
	}
	role, err := repository.getRole(ctx, assignment.Scope(), assignment.RoleID())
	if err != nil {
		return err
	}
	if !role.Scope().Equal(assignment.Scope()) {
		return scopeDeniedFailure()
	}
	command, err := repository.pool.Exec(
		ctx,
		`INSERT INTO omnexa_authorization.role_assignments
		 (id, role_id, principal_id, assignment_state, created_at, updated_at)
		 SELECT $1, $2, p.id, $4, $5, $6
		 FROM omnexa_identity.principals p
		 JOIN omnexa_identity.service_accounts sa ON sa.principal_id = p.id
		 WHERE p.id = $3
		   AND p.principal_type = 'service_account'
		   AND p.lifecycle_state = 'active'
		   AND sa.tenant_id = $7
		   AND sa.organization_id IS NOT DISTINCT FROM $8::uuid`,
		string(assignment.ID()), string(assignment.RoleID()), string(assignment.PrincipalID()),
		string(assignment.State()), assignment.CreatedAt(), assignment.UpdatedAt(),
		string(assignment.Scope().TenantID()), organizationArgument(assignment.Scope()),
	)
	if err != nil {
		return assignmentPersistenceFailure(err)
	}
	if command.RowsAffected() != 1 {
		return scopeDeniedFailure()
	}
	return nil
}

func (repository *PostgresRepository) getServiceAccountAssignment(
	ctx context.Context,
	scope Scope,
	id AssignmentID,
) (ServiceAccountAssignment, error) {
	if repository == nil || repository.pool == nil || !scope.Valid() || !id.Valid() {
		return ServiceAccountAssignment{}, invalidAssignmentFailure()
	}
	var storedID, roleID, principalID, state, tenantID string
	var organizationID *string
	var createdAt, updatedAt time.Time
	err := repository.pool.QueryRow(
		ctx,
		`SELECT a.id, a.role_id, a.principal_id, a.assignment_state, a.created_at, a.updated_at,
		        r.tenant_id, r.organization_id
		 FROM omnexa_authorization.role_assignments a
		 JOIN omnexa_authorization.roles r ON r.id = a.role_id
		 JOIN omnexa_identity.principals p ON p.id = a.principal_id
		 JOIN omnexa_identity.service_accounts sa ON sa.principal_id = p.id
		 WHERE a.id = $1
		   AND p.principal_type = 'service_account'
		   AND r.tenant_id = $2
		   AND r.organization_id IS NOT DISTINCT FROM $3::uuid
		   AND sa.tenant_id = r.tenant_id
		   AND sa.organization_id IS NOT DISTINCT FROM r.organization_id`,
		string(id), string(scope.TenantID()), organizationArgument(scope),
	).Scan(&storedID, &roleID, &principalID, &state, &createdAt, &updatedAt, &tenantID, &organizationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ServiceAccountAssignment{}, assignmentNotFoundFailure()
	}
	if err != nil {
		return ServiceAccountAssignment{}, repositoryFailure(err)
	}
	storedScope, err := scopeFromStored(tenantID, organizationID)
	if err != nil {
		return ServiceAccountAssignment{}, invalidStoredAssignmentFailure()
	}
	assignment, err := newServiceAccountAssignmentAt(
		AssignmentID(storedID), RoleID(roleID), identity.ServiceAccountID(principalID), storedScope,
		AssignmentState(state), createdAt, updatedAt,
	)
	if err != nil {
		return ServiceAccountAssignment{}, invalidStoredAssignmentFailure()
	}
	return assignment, nil
}

func (repository *PostgresRepository) revokeServiceAccountAssignment(
	ctx context.Context,
	scope Scope,
	id AssignmentID,
	changedAt time.Time,
) (ServiceAccountAssignment, error) {
	if repository == nil || repository.pool == nil || !scope.Valid() || !id.Valid() || changedAt.IsZero() {
		return ServiceAccountAssignment{}, invalidAssignmentFailure()
	}
	var storedID, roleID, principalID, state string
	var createdAt, updatedAt time.Time
	err := repository.pool.QueryRow(
		ctx,
		`UPDATE omnexa_authorization.role_assignments a
		 SET assignment_state = 'revoked', updated_at = $4
		 FROM omnexa_authorization.roles r, omnexa_identity.principals p, omnexa_identity.service_accounts sa
		 WHERE a.id = $1
		   AND a.role_id = r.id
		   AND p.id = a.principal_id
		   AND p.principal_type = 'service_account'
		   AND sa.principal_id = p.id
		   AND sa.tenant_id = r.tenant_id
		   AND sa.organization_id IS NOT DISTINCT FROM r.organization_id
		   AND r.tenant_id = $2
		   AND r.organization_id IS NOT DISTINCT FROM $3::uuid
		   AND a.assignment_state = 'active'
		   AND a.updated_at <= $4
		 RETURNING a.id, a.role_id, a.principal_id, a.assignment_state, a.created_at, a.updated_at`,
		string(id), string(scope.TenantID()), organizationArgument(scope), changedAt.UTC(),
	).Scan(&storedID, &roleID, &principalID, &state, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ServiceAccountAssignment{}, assignmentConflictFailure()
	}
	if err != nil {
		return ServiceAccountAssignment{}, assignmentPersistenceFailure(err)
	}
	return newServiceAccountAssignmentAt(
		AssignmentID(storedID), RoleID(roleID), identity.ServiceAccountID(principalID), scope,
		AssignmentState(state), createdAt, updatedAt,
	)
}

func (repository *PostgresRepository) hasServiceAccountPermission(
	ctx context.Context,
	subject ServiceAccountSubject,
	permission PermissionID,
) (bool, error) {
	if repository == nil || repository.pool == nil || !subject.Valid() || !permission.Valid() {
		return false, invalidSubjectFailure()
	}
	var allowed bool
	err := repository.pool.QueryRow(
		ctx,
		`SELECT EXISTS (
		   SELECT 1
		   FROM omnexa_authorization.role_assignments a
		   JOIN omnexa_authorization.roles r ON r.id = a.role_id
		   JOIN omnexa_authorization.role_permissions rp ON rp.role_id = r.id
		   JOIN omnexa_identity.principals p ON p.id = a.principal_id
		   JOIN omnexa_identity.service_accounts sa ON sa.principal_id = p.id
		   JOIN omnexa_tenancy.tenants t ON t.id = r.tenant_id
		   WHERE a.principal_id = $1
		     AND a.assignment_state = 'active'
		     AND p.principal_type = 'service_account'
		     AND p.lifecycle_state = 'active'
		     AND t.lifecycle_state = 'active'
		     AND sa.tenant_id = r.tenant_id
		     AND sa.organization_id IS NOT DISTINCT FROM r.organization_id
		     AND r.tenant_id = $2
		     AND r.organization_id IS NOT DISTINCT FROM $3::uuid
		     AND rp.permission_id = $4
		 )`,
		string(subject.PrincipalID()), string(subject.Scope().TenantID()), organizationArgument(subject.Scope()), string(permission),
	).Scan(&allowed)
	if err != nil {
		return false, repositoryFailure(err)
	}
	return allowed, nil
}

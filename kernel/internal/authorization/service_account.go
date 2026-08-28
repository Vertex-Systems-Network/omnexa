package authorization

import (
	"context"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/google/uuid"
)

// ServiceAccountSubject is a distinct non-human direct-RBAC subject. It is built
// only from a current active ServiceAccount or an authenticated current API
// credential and never reuses identity.UserID.
type ServiceAccountSubject struct {
	principalID identity.ServiceAccountID
	scope       Scope
}

// ServiceAccountSubjectFromAccount constructs a target suitable for a privileged
// human role assignment. The service principal must already be active.
func ServiceAccountSubjectFromAccount(account identity.ServiceAccount) (ServiceAccountSubject, error) {
	if !account.ID().Valid() || account.PrincipalType() != identity.PrincipalTypeServiceAccount || account.State() != identity.LifecycleActive {
		return ServiceAccountSubject{}, invalidSubjectFailure()
	}
	scope, err := scopeFromServiceAccountBinding(account.Binding())
	if err != nil {
		return ServiceAccountSubject{}, err
	}
	return ServiceAccountSubject{principalID: account.ID(), scope: scope}, nil
}

// ServiceAccountSubjectFromAuthentication constructs the runtime RBAC subject from
// P02.08 credential proof. Credential possession alone still grants no permission.
func ServiceAccountSubjectFromAuthentication(authenticated identity.AuthenticatedServiceAccount) (ServiceAccountSubject, error) {
	if !authenticated.Valid() {
		return ServiceAccountSubject{}, invalidSubjectFailure()
	}
	return ServiceAccountSubjectFromAccount(authenticated.ServiceAccount())
}

func (subject ServiceAccountSubject) Valid() bool {
	return subject.principalID.Valid() && subject.scope.Valid()
}

func (subject ServiceAccountSubject) PrincipalID() identity.ServiceAccountID {
	return subject.principalID
}
func (subject ServiceAccountSubject) Scope() Scope { return subject.scope }

func scopeFromServiceAccountBinding(binding identity.ServiceAccountBinding) (Scope, error) {
	if !binding.Valid() {
		return Scope{}, invalidSubjectFailure()
	}
	scope := Scope{tenantID: tenancy.TenantID(binding.TenantID())}
	if binding.OrganizationID() == "" {
		scope.kind = ScopeTenant
	} else {
		scope.kind = ScopeOrganization
		scope.organizationID = organization.NodeID(binding.OrganizationID())
	}
	if !scope.Valid() {
		return Scope{}, invalidSubjectFailure()
	}
	return scope, nil
}

// ServiceAccountAssignment is a direct role grant to one non-human service
// principal. It shares role/permission semantics with P02.05 but is not a human Assignment.
type ServiceAccountAssignment struct {
	id          AssignmentID
	roleID      RoleID
	principalID identity.ServiceAccountID
	scope       Scope
	state       AssignmentState
	createdAt   time.Time
	updatedAt   time.Time
}

func newServiceAccountAssignment(role Role, target ServiceAccountSubject, createdAt time.Time) (ServiceAccountAssignment, error) {
	identifier, err := uuid.NewV7()
	if err != nil {
		return ServiceAccountAssignment{}, identifierFailure(err)
	}
	return newServiceAccountAssignmentAt(
		AssignmentID(identifier.String()), role.ID(), target.PrincipalID(), role.Scope(), AssignmentActive, createdAt, createdAt,
	)
}

func newServiceAccountAssignmentAt(
	id AssignmentID,
	roleID RoleID,
	principalID identity.ServiceAccountID,
	scope Scope,
	state AssignmentState,
	createdAt time.Time,
	updatedAt time.Time,
) (ServiceAccountAssignment, error) {
	assignment := ServiceAccountAssignment{
		id: id, roleID: roleID, principalID: principalID, scope: scope, state: state,
		createdAt: createdAt.UTC(), updatedAt: updatedAt.UTC(),
	}
	if assignment.validate() != nil {
		return ServiceAccountAssignment{}, invalidAssignmentFailure()
	}
	return assignment, nil
}

func (assignment ServiceAccountAssignment) validate() error {
	if !assignment.id.Valid() || !assignment.roleID.Valid() || !assignment.principalID.Valid() || !assignment.scope.Valid() {
		return invalidAssignmentFailure()
	}
	if assignment.state != AssignmentActive && assignment.state != AssignmentRevoked {
		return invalidAssignmentFailure()
	}
	if assignment.createdAt.IsZero() || assignment.updatedAt.IsZero() || assignment.updatedAt.Before(assignment.createdAt) {
		return invalidAssignmentFailure()
	}
	return nil
}

func (assignment ServiceAccountAssignment) ID() AssignmentID { return assignment.id }
func (assignment ServiceAccountAssignment) RoleID() RoleID   { return assignment.roleID }
func (assignment ServiceAccountAssignment) PrincipalID() identity.ServiceAccountID {
	return assignment.principalID
}
func (assignment ServiceAccountAssignment) Scope() Scope           { return assignment.scope }
func (assignment ServiceAccountAssignment) State() AssignmentState { return assignment.state }
func (assignment ServiceAccountAssignment) CreatedAt() time.Time   { return assignment.createdAt }
func (assignment ServiceAccountAssignment) UpdatedAt() time.Time   { return assignment.updatedAt }

// CheckServiceAccount evaluates one direct permission using the same accepted
// P02.05 roles/permissions/assignment owner. Module availability is a fail-closed
// precondition and never substitutes for the existing exact-scope role grant.
func (service *Service) CheckServiceAccount(ctx context.Context, subject ServiceAccountSubject, permission PermissionID) (Decision, error) {
	if service == nil || service.repository == nil || service.audit == nil {
		return DecisionDeny, serviceInvalidFailure()
	}
	if !subject.Valid() {
		return DecisionDeny, invalidSubjectFailure()
	}
	if !permission.Valid() {
		return DecisionDeny, invalidPermissionFailure()
	}
	available, availabilityErr := service.permissionAvailable(ctx, permission)
	if availabilityErr != nil {
		return DecisionDeny, availabilityErr
	}
	if !available {
		return DecisionDeny, nil
	}
	allowed, err := service.repository.hasServiceAccountPermission(ctx, subject, permission)
	if err != nil {
		return DecisionDeny, err
	}
	if !allowed {
		return DecisionDeny, nil
	}
	return DecisionAllow, nil
}

func (service *Service) RequireServiceAccount(ctx context.Context, subject ServiceAccountSubject, permission PermissionID) error {
	decision, err := service.CheckServiceAccount(ctx, subject, permission)
	if err != nil {
		return err
	}
	if !decision.Allowed() {
		return deniedFailure()
	}
	return nil
}

// AssignRoleToServiceAccount creates one direct grant after the same human actor
// authority, anti-escalation, exact-scope, and audit controls used by P02.05.
func (service *Service) AssignRoleToServiceAccount(
	ctx context.Context,
	actor Subject,
	target ServiceAccountSubject,
	roleID RoleID,
	metadata MutationMetadata,
) (ServiceAccountAssignment, error) {
	if validationErr := service.validateMutation(actor, metadata); validationErr != nil {
		return ServiceAccountAssignment{}, validationErr
	}
	if !target.Valid() || !roleID.Valid() {
		return ServiceAccountAssignment{}, invalidAssignmentFailure()
	}
	if !actor.Scope().Equal(target.Scope()) {
		scopeErr := scopeDeniedFailure()
		return ServiceAccountAssignment{}, service.deniedServiceAssignment(ctx, actor, roleID, metadata, scopeErr)
	}
	if authorityErr := service.requireMutationAuthority(ctx, actor, PermissionAssignmentManage); authorityErr != nil {
		return ServiceAccountAssignment{}, service.deniedServiceAssignment(ctx, actor, roleID, metadata, authorityErr)
	}
	role, roleErr := service.repository.getRole(ctx, actor.Scope(), roleID)
	if roleErr != nil {
		return ServiceAccountAssignment{}, roleErr
	}
	if !role.Scope().Equal(target.Scope()) {
		return ServiceAccountAssignment{}, service.deniedServiceAssignment(ctx, actor, roleID, metadata, scopeDeniedFailure())
	}
	if grantErr := service.ensureGrantable(ctx, actor, role.Permissions()); grantErr != nil {
		return ServiceAccountAssignment{}, service.deniedServiceAssignment(ctx, actor, roleID, metadata, grantErr)
	}
	assignment, assignmentErr := newServiceAccountAssignment(role, target, service.now())
	if assignmentErr != nil {
		return ServiceAccountAssignment{}, assignmentErr
	}
	if createErr := service.repository.createServiceAccountAssignment(ctx, assignment); createErr != nil {
		return ServiceAccountAssignment{}, createErr
	}
	if auditErr := service.auditMutation(
		ctx, actor, "authorization.service_account_assignment.create", "authorization.service_account_assignment",
		string(assignment.ID()), audit.OutcomeSucceeded, metadata,
	); auditErr != nil {
		return ServiceAccountAssignment{}, auditErr
	}
	return assignment, nil
}

// RevokeServiceAccountAssignment terminally revokes one service-principal grant
// using the same privileged human assignment-management boundary.
func (service *Service) RevokeServiceAccountAssignment(
	ctx context.Context,
	actor Subject,
	assignmentID AssignmentID,
	metadata MutationMetadata,
) (ServiceAccountAssignment, error) {
	if validationErr := service.validateMutation(actor, metadata); validationErr != nil {
		return ServiceAccountAssignment{}, validationErr
	}
	if !assignmentID.Valid() {
		return ServiceAccountAssignment{}, invalidAssignmentFailure()
	}
	if authorityErr := service.requireMutationAuthority(ctx, actor, PermissionAssignmentManage); authorityErr != nil {
		return ServiceAccountAssignment{}, service.deniedServiceAssignmentByID(ctx, actor, assignmentID, metadata, authorityErr)
	}
	assignment, err := service.repository.getServiceAccountAssignment(ctx, actor.Scope(), assignmentID)
	if err != nil {
		return ServiceAccountAssignment{}, err
	}
	if !assignment.Scope().Equal(actor.Scope()) {
		return ServiceAccountAssignment{}, service.deniedServiceAssignmentByID(ctx, actor, assignmentID, metadata, scopeDeniedFailure())
	}
	revoked, err := service.repository.revokeServiceAccountAssignment(ctx, actor.Scope(), assignmentID, service.now())
	if err != nil {
		return ServiceAccountAssignment{}, err
	}
	if auditErr := service.auditMutation(
		ctx, actor, "authorization.service_account_assignment.revoke", "authorization.service_account_assignment",
		string(assignmentID), audit.OutcomeSucceeded, metadata,
	); auditErr != nil {
		return ServiceAccountAssignment{}, auditErr
	}
	return revoked, nil
}

func (service *Service) deniedServiceAssignment(
	ctx context.Context,
	actor Subject,
	roleID RoleID,
	metadata MutationMetadata,
	denial error,
) error {
	return service.deniedMutation(
		ctx, actor, "authorization.service_account_assignment.create", "authorization.role", string(roleID), metadata, denial,
	)
}

func (service *Service) deniedServiceAssignmentByID(
	ctx context.Context,
	actor Subject,
	assignmentID AssignmentID,
	metadata MutationMetadata,
	denial error,
) error {
	return service.deniedMutation(
		ctx, actor, "authorization.service_account_assignment.revoke", "authorization.service_account_assignment",
		string(assignmentID), metadata, denial,
	)
}

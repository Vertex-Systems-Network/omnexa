package authorization

import (
	"context"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
)

const actorKindUser = "user"

// Service is the P02.05 direct-RBAC capability boundary. It evaluates only
// exact-scope direct role grants and owns all protected RBAC mutations.
type Service struct {
	repository        repository
	audit             *audit.Writer
	now               func() time.Time
	modulePermissions ModulePermissionAvailability
}

// NewService creates an authorization service with required protected audit
// transport. Kernel permissions keep their accepted P02 behavior. Module
// permissions fail closed unless a P03.07 live availability provider is bound.
func NewService(repository repository, auditWriter *audit.Writer) (*Service, error) {
	return newServiceWithClock(repository, auditWriter, func() time.Time { return time.Now().UTC() })
}

// NewServiceWithModulePermissionAvailability binds the live P03.07 module
// permission availability precondition while preserving this service as the
// sole deny-by-default role/policy enforcement boundary.
func NewServiceWithModulePermissionAvailability(
	repository repository,
	auditWriter *audit.Writer,
	availability ModulePermissionAvailability,
) (*Service, error) {
	if availability == nil {
		return nil, serviceInvalidFailure()
	}
	return newServiceWithClockAndModulePermissions(
		repository,
		auditWriter,
		func() time.Time { return time.Now().UTC() },
		availability,
	)
}

func newServiceWithClock(repository repository, auditWriter *audit.Writer, now func() time.Time) (*Service, error) {
	return newServiceWithClockAndModulePermissions(repository, auditWriter, now, nil)
}

func newServiceWithClockAndModulePermissions(
	repository repository,
	auditWriter *audit.Writer,
	now func() time.Time,
	availability ModulePermissionAvailability,
) (*Service, error) {
	if repository == nil || auditWriter == nil || now == nil {
		return nil, serviceInvalidFailure()
	}
	return &Service{
		repository:        repository,
		audit:             auditWriter,
		now:               now,
		modulePermissions: availability,
	}, nil
}

// Check evaluates one direct permission at the subject's exact trusted scope.
// No role-name shortcut, tenant-to-organization inheritance, relationship rule,
// contextual condition or internal-caller bypass exists. Module permission
// availability is only a fail-closed precondition; an allow still requires the
// existing repository-backed exact-scope role grant.
func (service *Service) Check(ctx context.Context, subject Subject, permission PermissionID) (Decision, error) {
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
	allowed, permissionErr := service.repository.hasPermission(ctx, subject, permission)
	if permissionErr != nil {
		return DecisionDeny, permissionErr
	}
	if !allowed {
		return DecisionDeny, nil
	}
	return DecisionAllow, nil
}

// Require returns a disclosure-safe authorization failure unless Check allows
// the exact permission at the exact subject scope.
func (service *Service) Require(ctx context.Context, subject Subject, permission PermissionID) error {
	decision, checkErr := service.Check(ctx, subject, permission)
	if checkErr != nil {
		return checkErr
	}
	if !decision.Allowed() {
		return deniedFailure()
	}
	return nil
}

// CreateRole creates one exact-scope Role after server-side privilege and
// anti-escalation checks. Role names never alter these checks.
func (service *Service) CreateRole(
	ctx context.Context,
	actor Subject,
	name string,
	permissions []PermissionID,
	metadata MutationMetadata,
) (Role, error) {
	if validationErr := service.validateMutation(actor, metadata); validationErr != nil {
		return Role{}, validationErr
	}
	role, roleErr := newRole(actor.Scope(), name, permissions, service.now())
	if roleErr != nil {
		return Role{}, roleErr
	}
	if authorityErr := service.requireMutationAuthority(ctx, actor, PermissionRoleManage); authorityErr != nil {
		return Role{}, service.deniedMutation(ctx, actor, "authorization.role.create", "authorization.role", string(role.ID()), metadata, authorityErr)
	}
	exists, existenceErr := service.repository.permissionsExist(ctx, role.Permissions())
	if existenceErr != nil {
		return Role{}, existenceErr
	}
	if !exists {
		return Role{}, invalidPermissionFailure()
	}
	if grantErr := service.ensureGrantable(ctx, actor, role.Permissions()); grantErr != nil {
		return Role{}, service.deniedMutation(ctx, actor, "authorization.role.create", "authorization.role", string(role.ID()), metadata, grantErr)
	}
	if createErr := service.repository.createRole(ctx, role); createErr != nil {
		return Role{}, createErr
	}
	if auditErr := service.auditMutation(ctx, actor, "authorization.role.create", "authorization.role", string(role.ID()), audit.OutcomeSucceeded, metadata); auditErr != nil {
		return Role{}, auditErr
	}
	return role, nil
}

// ReplaceRolePermissions atomically replaces one role's direct composition after
// verifying the actor already possesses every permission they attempt to grant.
func (service *Service) ReplaceRolePermissions(
	ctx context.Context,
	actor Subject,
	roleID RoleID,
	permissions []PermissionID,
	metadata MutationMetadata,
) (Role, error) {
	if validationErr := service.validateMutation(actor, metadata); validationErr != nil {
		return Role{}, validationErr
	}
	if !roleID.Valid() {
		return Role{}, invalidRoleFailure()
	}
	if authorityErr := service.requireMutationAuthority(ctx, actor, PermissionRoleManage); authorityErr != nil {
		return Role{}, service.deniedMutation(ctx, actor, "authorization.role.permissions.replace", "authorization.role", string(roleID), metadata, authorityErr)
	}
	role, roleErr := service.repository.getRole(ctx, actor.Scope(), roleID)
	if roleErr != nil {
		return Role{}, roleErr
	}
	normalized, normalizeErr := normalizePermissions(permissions)
	if normalizeErr != nil {
		return Role{}, normalizeErr
	}
	exists, existenceErr := service.repository.permissionsExist(ctx, normalized)
	if existenceErr != nil {
		return Role{}, existenceErr
	}
	if !exists {
		return Role{}, invalidPermissionFailure()
	}
	if grantErr := service.ensureGrantable(ctx, actor, normalized); grantErr != nil {
		return Role{}, service.deniedMutation(ctx, actor, "authorization.role.permissions.replace", "authorization.role", string(role.ID()), metadata, grantErr)
	}
	updated, replaceErr := service.repository.replaceRolePermissions(ctx, actor.Scope(), role.ID(), normalized, service.now())
	if replaceErr != nil {
		return Role{}, replaceErr
	}
	if auditErr := service.auditMutation(ctx, actor, "authorization.role.permissions.replace", "authorization.role", string(role.ID()), audit.OutcomeSucceeded, metadata); auditErr != nil {
		return Role{}, auditErr
	}
	return updated, nil
}

// AssignRole creates one direct role grant. Actor, target and role must use the
// same exact trusted scope, and the actor may not grant permissions they lack.
func (service *Service) AssignRole(
	ctx context.Context,
	actor Subject,
	target Subject,
	roleID RoleID,
	metadata MutationMetadata,
) (Assignment, error) {
	if validationErr := service.validateMutation(actor, metadata); validationErr != nil {
		return Assignment{}, validationErr
	}
	if !target.Valid() || !roleID.Valid() {
		return Assignment{}, invalidAssignmentFailure()
	}
	if !actor.Scope().Equal(target.Scope()) {
		scopeErr := scopeDeniedFailure()
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.create", "authorization.role", string(roleID), metadata, scopeErr)
	}
	if authorityErr := service.requireMutationAuthority(ctx, actor, PermissionAssignmentManage); authorityErr != nil {
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.create", "authorization.role", string(roleID), metadata, authorityErr)
	}
	role, roleErr := service.repository.getRole(ctx, actor.Scope(), roleID)
	if roleErr != nil {
		return Assignment{}, roleErr
	}
	if !role.Scope().Equal(target.Scope()) {
		scopeErr := scopeDeniedFailure()
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.create", "authorization.role", string(role.ID()), metadata, scopeErr)
	}
	if grantErr := service.ensureGrantable(ctx, actor, role.Permissions()); grantErr != nil {
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.create", "authorization.role", string(role.ID()), metadata, grantErr)
	}
	assignment, assignmentErr := newAssignment(role, target, service.now())
	if assignmentErr != nil {
		return Assignment{}, assignmentErr
	}
	if createErr := service.repository.createAssignment(ctx, assignment); createErr != nil {
		return Assignment{}, createErr
	}
	if auditErr := service.auditMutation(ctx, actor, "authorization.assignment.create", "authorization.assignment", string(assignment.ID()), audit.OutcomeSucceeded, metadata); auditErr != nil {
		return Assignment{}, auditErr
	}
	return assignment, nil
}

// RevokeAssignment terminally revokes one direct grant inside the actor's exact
// scope after server-side assignment-management authorization.
func (service *Service) RevokeAssignment(
	ctx context.Context,
	actor Subject,
	assignmentID AssignmentID,
	metadata MutationMetadata,
) (Assignment, error) {
	if validationErr := service.validateMutation(actor, metadata); validationErr != nil {
		return Assignment{}, validationErr
	}
	if !assignmentID.Valid() {
		return Assignment{}, invalidAssignmentFailure()
	}
	if authorityErr := service.requireMutationAuthority(ctx, actor, PermissionAssignmentManage); authorityErr != nil {
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.revoke", "authorization.assignment", string(assignmentID), metadata, authorityErr)
	}
	assignment, assignmentErr := service.repository.getAssignment(ctx, actor.Scope(), assignmentID)
	if assignmentErr != nil {
		return Assignment{}, assignmentErr
	}
	if !assignment.Scope().Equal(actor.Scope()) {
		scopeErr := scopeDeniedFailure()
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.revoke", "authorization.assignment", string(assignmentID), metadata, scopeErr)
	}
	revoked, revokeErr := service.repository.revokeAssignment(ctx, actor.Scope(), assignment.ID(), service.now())
	if revokeErr != nil {
		return Assignment{}, revokeErr
	}
	if auditErr := service.auditMutation(ctx, actor, "authorization.assignment.revoke", "authorization.assignment", string(assignment.ID()), audit.OutcomeSucceeded, metadata); auditErr != nil {
		return Assignment{}, auditErr
	}
	return revoked, nil
}

func (service *Service) validateMutation(actor Subject, metadata MutationMetadata) error {
	if service == nil || service.repository == nil || service.audit == nil || service.now == nil {
		return serviceInvalidFailure()
	}
	if !actor.Valid() {
		return invalidSubjectFailure()
	}
	return metadata.validate()
}

func (service *Service) requireMutationAuthority(ctx context.Context, actor Subject, permission PermissionID) error {
	return service.Require(ctx, actor, permission)
}

func (service *Service) ensureGrantable(ctx context.Context, actor Subject, permissions []PermissionID) error {
	for _, permission := range permissions {
		if requireErr := service.Require(ctx, actor, permission); requireErr != nil {
			return requireErr
		}
	}
	return nil
}

func (service *Service) deniedMutation(
	ctx context.Context,
	actor Subject,
	action string,
	targetKind string,
	targetReference string,
	metadata MutationMetadata,
	denial error,
) error {
	if auditErr := service.auditMutation(ctx, actor, action, targetKind, targetReference, audit.OutcomeDenied, metadata); auditErr != nil {
		return auditErr
	}
	return denial
}

func (service *Service) auditMutation(
	ctx context.Context,
	actor Subject,
	action string,
	targetKind string,
	targetReference string,
	outcome audit.Outcome,
	metadata MutationMetadata,
) error {
	scope := audit.Scope{TenantID: string(actor.Scope().TenantID())}
	if actor.Scope().Kind() == ScopeOrganization {
		scope.OrganizationID = string(actor.Scope().OrganizationID())
	}
	_, writeErr := service.audit.Write(ctx, audit.RequirementRequired, audit.RecordInput{
		Classification: audit.ClassificationInternal,
		Actor: audit.Actor{
			Kind:      actorKindUser,
			Reference: string(actor.PrincipalID()),
		},
		Action: action,
		Target: audit.Target{
			Kind:      targetKind,
			Reference: targetReference,
		},
		Scope:         scope,
		Outcome:       outcome,
		CorrelationID: metadata.CorrelationID,
		Reason:        metadata.Reason,
		Privileged:    true,
	})
	return writeErr
}

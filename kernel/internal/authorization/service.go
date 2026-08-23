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
	repository repository
	audit      *audit.Writer
	now        func() time.Time
}

// NewService creates an authorization service with required protected audit
// transport. A nil repository or audit writer fails closed.
func NewService(repository repository, auditWriter *audit.Writer) (*Service, error) {
	return newServiceWithClock(repository, auditWriter, func() time.Time { return time.Now().UTC() })
}

func newServiceWithClock(repository repository, auditWriter *audit.Writer, now func() time.Time) (*Service, error) {
	if repository == nil || auditWriter == nil || now == nil {
		return nil, serviceInvalidFailure()
	}
	return &Service{repository: repository, audit: auditWriter, now: now}, nil
}

// Check evaluates one direct permission at the subject's exact trusted scope.
// No role-name shortcut, tenant-to-organization inheritance, relationship rule,
// contextual condition or internal-caller bypass exists in P02.05.
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
	allowed, err := service.repository.hasPermission(ctx, subject, permission)
	if err != nil {
		return DecisionDeny, err
	}
	if !allowed {
		return DecisionDeny, nil
	}
	return DecisionAllow, nil
}

// Require returns a disclosure-safe authorization failure unless Check allows
// the exact permission at the exact subject scope.
func (service *Service) Require(ctx context.Context, subject Subject, permission PermissionID) error {
	decision, err := service.Check(ctx, subject, permission)
	if err != nil {
		return err
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
	if err := service.validateMutation(actor, metadata); err != nil {
		return Role{}, err
	}
	role, err := newRole(actor.Scope(), name, permissions, service.now())
	if err != nil {
		return Role{}, err
	}
	if err := service.requireMutationAuthority(ctx, actor, PermissionRoleManage); err != nil {
		return Role{}, service.deniedMutation(ctx, actor, "authorization.role.create", "authorization.role", string(role.ID()), metadata, err)
	}
	exists, err := service.repository.permissionsExist(ctx, role.Permissions())
	if err != nil {
		return Role{}, err
	}
	if !exists {
		return Role{}, invalidPermissionFailure()
	}
	if err := service.ensureGrantable(ctx, actor, role.Permissions()); err != nil {
		return Role{}, service.deniedMutation(ctx, actor, "authorization.role.create", "authorization.role", string(role.ID()), metadata, err)
	}
	if err := service.repository.createRole(ctx, role); err != nil {
		return Role{}, err
	}
	if err := service.auditMutation(ctx, actor, "authorization.role.create", "authorization.role", string(role.ID()), audit.OutcomeSucceeded, metadata); err != nil {
		return Role{}, err
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
	if err := service.validateMutation(actor, metadata); err != nil {
		return Role{}, err
	}
	if !roleID.Valid() {
		return Role{}, invalidRoleFailure()
	}
	if err := service.requireMutationAuthority(ctx, actor, PermissionRoleManage); err != nil {
		return Role{}, service.deniedMutation(ctx, actor, "authorization.role.permissions.replace", "authorization.role", string(roleID), metadata, err)
	}
	role, err := service.repository.getRole(ctx, actor.Scope(), roleID)
	if err != nil {
		return Role{}, err
	}
	normalized, err := normalizePermissions(permissions)
	if err != nil {
		return Role{}, err
	}
	exists, err := service.repository.permissionsExist(ctx, normalized)
	if err != nil {
		return Role{}, err
	}
	if !exists {
		return Role{}, invalidPermissionFailure()
	}
	if err := service.ensureGrantable(ctx, actor, normalized); err != nil {
		return Role{}, service.deniedMutation(ctx, actor, "authorization.role.permissions.replace", "authorization.role", string(role.ID()), metadata, err)
	}
	updated, err := service.repository.replaceRolePermissions(ctx, actor.Scope(), role.ID(), normalized, service.now())
	if err != nil {
		return Role{}, err
	}
	if err := service.auditMutation(ctx, actor, "authorization.role.permissions.replace", "authorization.role", string(role.ID()), audit.OutcomeSucceeded, metadata); err != nil {
		return Role{}, err
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
	if err := service.validateMutation(actor, metadata); err != nil {
		return Assignment{}, err
	}
	if !target.Valid() || !roleID.Valid() {
		return Assignment{}, invalidAssignmentFailure()
	}
	if !actor.Scope().Equal(target.Scope()) {
		err := scopeDeniedFailure()
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.create", "authorization.role", string(roleID), metadata, err)
	}
	if err := service.requireMutationAuthority(ctx, actor, PermissionAssignmentManage); err != nil {
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.create", "authorization.role", string(roleID), metadata, err)
	}
	role, err := service.repository.getRole(ctx, actor.Scope(), roleID)
	if err != nil {
		return Assignment{}, err
	}
	if !role.Scope().Equal(target.Scope()) {
		err := scopeDeniedFailure()
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.create", "authorization.role", string(role.ID()), metadata, err)
	}
	if err := service.ensureGrantable(ctx, actor, role.Permissions()); err != nil {
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.create", "authorization.role", string(role.ID()), metadata, err)
	}
	assignment, err := newAssignment(role, target, service.now())
	if err != nil {
		return Assignment{}, err
	}
	if err := service.repository.createAssignment(ctx, assignment); err != nil {
		return Assignment{}, err
	}
	if err := service.auditMutation(ctx, actor, "authorization.assignment.create", "authorization.assignment", string(assignment.ID()), audit.OutcomeSucceeded, metadata); err != nil {
		return Assignment{}, err
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
	if err := service.validateMutation(actor, metadata); err != nil {
		return Assignment{}, err
	}
	if !assignmentID.Valid() {
		return Assignment{}, invalidAssignmentFailure()
	}
	if err := service.requireMutationAuthority(ctx, actor, PermissionAssignmentManage); err != nil {
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.revoke", "authorization.assignment", string(assignmentID), metadata, err)
	}
	assignment, err := service.repository.getAssignment(ctx, actor.Scope(), assignmentID)
	if err != nil {
		return Assignment{}, err
	}
	if !assignment.Scope().Equal(actor.Scope()) {
		err := scopeDeniedFailure()
		return Assignment{}, service.deniedMutation(ctx, actor, "authorization.assignment.revoke", "authorization.assignment", string(assignmentID), metadata, err)
	}
	revoked, err := service.repository.revokeAssignment(ctx, actor.Scope(), assignment.ID(), service.now())
	if err != nil {
		return Assignment{}, err
	}
	if err := service.auditMutation(ctx, actor, "authorization.assignment.revoke", "authorization.assignment", string(assignment.ID()), audit.OutcomeSucceeded, metadata); err != nil {
		return Assignment{}, err
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
		if err := service.Require(ctx, actor, permission); err != nil {
			return err
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
	_, err := service.audit.Write(ctx, audit.RequirementRequired, audit.RecordInput{
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
	return err
}

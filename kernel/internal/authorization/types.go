// Package authorization implements the P02.05 deny-by-default direct RBAC foundation.
package authorization

import (
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/google/uuid"
)

const (
	maxRoleNameRunes    = 128
	maxMutationRunes    = 512
	maxCorrelationRunes = 128
	maxRolePermissions  = 128
)

var permissionPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:\.[a-z][a-z0-9_]*){2,7}$`)

// PermissionID is a stable capability-oriented authorization identifier.
type PermissionID string

const (
	PermissionRoleRead         PermissionID = "authorization.role.read"
	PermissionRoleManage       PermissionID = "authorization.role.manage"
	PermissionAssignmentRead   PermissionID = "authorization.assignment.read"
	PermissionAssignmentManage PermissionID = "authorization.assignment.manage"
)

// Valid reports whether permission is a canonical stable permission identifier.
func (permission PermissionID) Valid() bool {
	return permissionPattern.MatchString(string(permission))
}

// BuiltinPermission reports whether permission is part of the P02.05 kernel
// reference vocabulary. P03 owns future module permission registration.
func BuiltinPermission(permission PermissionID) bool {
	switch permission {
	case PermissionRoleRead, PermissionRoleManage, PermissionAssignmentRead, PermissionAssignmentManage:
		return true
	default:
		return false
	}
}

// ScopeKind is the exact direct-RBAC scope vocabulary. P02.05 deliberately has
// no inherited, object, relationship or contextual-policy scope semantics.
type ScopeKind string

const (
	ScopeTenant       ScopeKind = "tenant"
	ScopeOrganization ScopeKind = "organization"
)

// Scope is an exact trusted RBAC scope derived from current P02.02/P02.03
// contexts. Private fields prevent raw client identifiers from constructing it.
type Scope struct {
	kind           ScopeKind
	tenantID       tenancy.TenantID
	organizationID organization.NodeID
}

// Subject is one human principal paired with an exact trusted RBAC scope.
// It is execution-scoped and must not be persisted as a long-lived authority cache.
type Subject struct {
	principalID identity.UserID
	scope       Scope
}

// SubjectFromTenantContext creates an exact tenant-scoped subject from a trusted
// active P02.02 context.
func SubjectFromTenantContext(trusted tenancy.TrustedContext) (Subject, error) {
	if !trusted.Valid() {
		return Subject{}, invalidSubjectFailure()
	}
	return Subject{
		principalID: trusted.PrincipalID(),
		scope: Scope{
			kind:     ScopeTenant,
			tenantID: trusted.TenantID(),
		},
	}, nil
}

// SubjectFromOrganizationContext creates an exact organization-scoped subject
// from a trusted active P02.03 relationship context.
func SubjectFromOrganizationContext(scoped organization.ScopedContext) (Subject, error) {
	if !scoped.Valid() {
		return Subject{}, invalidSubjectFailure()
	}
	return Subject{
		principalID: scoped.PrincipalID(),
		scope: Scope{
			kind:           ScopeOrganization,
			tenantID:       scoped.TenantID(),
			organizationID: scoped.NodeID(),
		},
	}, nil
}

// PrincipalID returns the current human principal identifier.
func (subject Subject) PrincipalID() identity.UserID { return subject.principalID }

// Scope returns the exact trusted RBAC scope.
func (subject Subject) Scope() Scope { return subject.scope }

// Valid reports whether subject contains a canonical principal and exact scope.
func (subject Subject) Valid() bool {
	return subject.principalID.Valid() && subject.scope.Valid()
}

// Kind returns tenant or organization.
func (scope Scope) Kind() ScopeKind { return scope.kind }

// TenantID returns the trusted tenant isolation boundary.
func (scope Scope) TenantID() tenancy.TenantID { return scope.tenantID }

// OrganizationID returns the exact organization scope, or the zero value for a
// tenant-scoped role. Tenant scope never implicitly authorizes child organizations.
func (scope Scope) OrganizationID() organization.NodeID { return scope.organizationID }

// Valid reports whether scope is a canonical exact P02.05 RBAC scope.
func (scope Scope) Valid() bool {
	if !scope.tenantID.Valid() {
		return false
	}
	switch scope.kind {
	case ScopeTenant:
		return scope.organizationID == ""
	case ScopeOrganization:
		return scope.organizationID.Valid()
	default:
		return false
	}
}

// Equal reports exact scope equality. No parent/child inheritance occurs in P02.05.
func (scope Scope) Equal(other Scope) bool {
	return scope.Valid() && other.Valid() &&
		scope.kind == other.kind &&
		scope.tenantID == other.tenantID &&
		scope.organizationID == other.organizationID
}

// RoleID is the stable UUIDv7 identity of one permission composition.
type RoleID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id RoleID) Valid() bool {
	parsed, err := uuid.Parse(string(id))
	return err == nil && parsed.Version() == 7
}

// Role is an exact-scope permission composition. Name is descriptive only and
// never grants authority or implies a platform bypass.
type Role struct {
	id          RoleID
	scope       Scope
	name        string
	permissions []PermissionID
	createdAt   time.Time
	updatedAt   time.Time
}

func newRole(scope Scope, name string, permissions []PermissionID, createdAt time.Time) (Role, error) {
	identifier, err := uuid.NewV7()
	if err != nil {
		return Role{}, identifierFailure(err)
	}
	return newRoleAt(RoleID(identifier.String()), scope, name, permissions, createdAt)
}

func newRoleAt(id RoleID, scope Scope, name string, permissions []PermissionID, createdAt time.Time) (Role, error) {
	normalizedName, nameErr := normalizeRoleName(name)
	normalizedPermissions, permissionsErr := normalizePermissions(permissions)
	if nameErr != nil || permissionsErr != nil || !id.Valid() || !scope.Valid() || createdAt.IsZero() {
		return Role{}, invalidRoleFailure()
	}
	instant := createdAt.UTC()
	return Role{
		id:          id,
		scope:       scope,
		name:        normalizedName,
		permissions: normalizedPermissions,
		createdAt:   instant,
		updatedAt:   instant,
	}, nil
}

func rehydrateRole(id RoleID, scope Scope, name string, permissions []PermissionID, createdAt, updatedAt time.Time) (Role, error) {
	normalizedName, nameErr := normalizeRoleName(name)
	normalizedPermissions, permissionsErr := normalizePermissions(permissions)
	role := Role{
		id:          id,
		scope:       scope,
		name:        normalizedName,
		permissions: normalizedPermissions,
		createdAt:   createdAt.UTC(),
		updatedAt:   updatedAt.UTC(),
	}
	if nameErr != nil || permissionsErr != nil || role.validate() != nil {
		return Role{}, invalidStoredRoleFailure()
	}
	return role, nil
}

func (role Role) validate() error {
	if !role.id.Valid() || !role.scope.Valid() || role.createdAt.IsZero() || role.updatedAt.IsZero() || role.updatedAt.Before(role.createdAt) {
		return invalidRoleFailure()
	}
	if _, err := normalizeRoleName(role.name); err != nil {
		return invalidRoleFailure()
	}
	if _, err := normalizePermissions(role.permissions); err != nil {
		return invalidRoleFailure()
	}
	return nil
}

// ID returns the stable UUIDv7 Role identifier.
func (role Role) ID() RoleID { return role.id }

// Scope returns the exact trusted scope owned by the role.
func (role Role) Scope() Scope { return role.scope }

// Name returns the descriptive role name. It conveys no authority by itself.
func (role Role) Name() string { return role.name }

// Permissions returns a sorted defensive copy of the role composition.
func (role Role) Permissions() []PermissionID {
	return append([]PermissionID(nil), role.permissions...)
}

// HasPermission reports direct membership in this role composition only.
func (role Role) HasPermission(permission PermissionID) bool {
	index := sort.Search(len(role.permissions), func(index int) bool { return role.permissions[index] >= permission })
	return index < len(role.permissions) && role.permissions[index] == permission
}

// CreatedAt returns the canonical UTC creation instant.
func (role Role) CreatedAt() time.Time { return role.createdAt }

// UpdatedAt returns the canonical UTC last composition-change instant.
func (role Role) UpdatedAt() time.Time { return role.updatedAt }

// AssignmentID is the stable UUIDv7 identity of one direct role assignment.
type AssignmentID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id AssignmentID) Valid() bool {
	parsed, err := uuid.Parse(string(id))
	return err == nil && parsed.Version() == 7
}

// AssignmentState is the direct role-assignment lifecycle vocabulary.
type AssignmentState string

const (
	AssignmentActive  AssignmentState = "active"
	AssignmentRevoked AssignmentState = "revoked"
)

// Assignment binds one role to one human principal at the role's exact scope.
type Assignment struct {
	id          AssignmentID
	roleID      RoleID
	principalID identity.UserID
	scope       Scope
	state       AssignmentState
	createdAt   time.Time
	updatedAt   time.Time
}

func newAssignment(role Role, target Subject, createdAt time.Time) (Assignment, error) {
	identifier, err := uuid.NewV7()
	if err != nil {
		return Assignment{}, identifierFailure(err)
	}
	return newAssignmentAt(AssignmentID(identifier.String()), role.ID(), target.PrincipalID(), role.Scope(), AssignmentActive, createdAt, createdAt)
}

func newAssignmentAt(
	id AssignmentID,
	roleID RoleID,
	principalID identity.UserID,
	scope Scope,
	state AssignmentState,
	createdAt time.Time,
	updatedAt time.Time,
) (Assignment, error) {
	assignment := Assignment{
		id:          id,
		roleID:      roleID,
		principalID: principalID,
		scope:       scope,
		state:       state,
		createdAt:   createdAt.UTC(),
		updatedAt:   updatedAt.UTC(),
	}
	if assignment.validate() != nil {
		return Assignment{}, invalidAssignmentFailure()
	}
	return assignment, nil
}

func (assignment Assignment) validate() error {
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

// ID returns the stable assignment UUIDv7 identifier.
func (assignment Assignment) ID() AssignmentID { return assignment.id }

// RoleID returns the referenced role identifier.
func (assignment Assignment) RoleID() RoleID { return assignment.roleID }

// PrincipalID returns the human principal receiving the direct role grant.
func (assignment Assignment) PrincipalID() identity.UserID { return assignment.principalID }

// Scope returns the role's exact direct-RBAC scope.
func (assignment Assignment) Scope() Scope { return assignment.scope }

// State returns active or revoked.
func (assignment Assignment) State() AssignmentState { return assignment.state }

// CreatedAt returns the canonical UTC creation instant.
func (assignment Assignment) CreatedAt() time.Time { return assignment.createdAt }

// UpdatedAt returns the canonical UTC last lifecycle-change instant.
func (assignment Assignment) UpdatedAt() time.Time { return assignment.updatedAt }

// Decision is the deliberately small direct RBAC result vocabulary.
type Decision string

const (
	DecisionAllow Decision = "allow"
	DecisionDeny  Decision = "deny"
)

// Allowed reports whether the decision grants the exact requested permission.
func (decision Decision) Allowed() bool { return decision == DecisionAllow }

// MutationMetadata is bounded, classification-safe metadata required for
// privileged RBAC changes and the protected audit boundary.
type MutationMetadata struct {
	CorrelationID string
	Reason        string
}

func (metadata MutationMetadata) validate() error {
	if !validReference(metadata.CorrelationID, maxCorrelationRunes) || !validText(metadata.Reason, maxMutationRunes) || strings.TrimSpace(metadata.Reason) == "" {
		return invalidMutationMetadataFailure()
	}
	return nil
}

func normalizeRoleName(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || utf8.RuneCountInString(trimmed) > maxRoleNameRunes || !validText(trimmed, maxRoleNameRunes) {
		return "", invalidRoleFailure()
	}
	return trimmed, nil
}

func normalizePermissions(input []PermissionID) ([]PermissionID, error) {
	if len(input) == 0 || len(input) > maxRolePermissions {
		return nil, invalidRoleFailure()
	}
	seen := make(map[PermissionID]struct{}, len(input))
	result := make([]PermissionID, 0, len(input))
	for _, permission := range input {
		if !permission.Valid() {
			return nil, invalidPermissionFailure()
		}
		if _, exists := seen[permission]; exists {
			continue
		}
		seen[permission] = struct{}{}
		result = append(result, permission)
	}
	if len(result) == 0 {
		return nil, invalidRoleFailure()
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result, nil
}

func validReference(value string, maxRunes int) bool {
	if strings.TrimSpace(value) == "" || utf8.RuneCountInString(value) > maxRunes {
		return false
	}
	for _, char := range value {
		if unicode.IsSpace(char) || unicode.IsControl(char) {
			return false
		}
	}
	return true
}

func validText(value string, maxRunes int) bool {
	if utf8.RuneCountInString(value) > maxRunes {
		return false
	}
	for _, char := range value {
		if unicode.IsControl(char) && char != '\t' {
			return false
		}
	}
	return true
}

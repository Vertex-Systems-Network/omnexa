package authorization

import (
	"context"
	"time"
)

// repository is deliberately owner-private. External kernel/domain code invokes
// the Service rather than bypassing server-side privileged mutation checks.
type repository interface {
	permissionsExist(context.Context, []PermissionID) (bool, error)
	createRole(context.Context, Role) error
	getRole(context.Context, Scope, RoleID) (Role, error)
	replaceRolePermissions(context.Context, Scope, RoleID, []PermissionID, time.Time) (Role, error)
	createAssignment(context.Context, Assignment) error
	getAssignment(context.Context, Scope, AssignmentID) (Assignment, error)
	revokeAssignment(context.Context, Scope, AssignmentID, time.Time) (Assignment, error)
	hasPermission(context.Context, Subject, PermissionID) (bool, error)
	createServiceAccountAssignment(context.Context, ServiceAccountAssignment) error
	getServiceAccountAssignment(context.Context, Scope, AssignmentID) (ServiceAccountAssignment, error)
	revokeServiceAccountAssignment(context.Context, Scope, AssignmentID, time.Time) (ServiceAccountAssignment, error)
	hasServiceAccountPermission(context.Context, ServiceAccountSubject, PermissionID) (bool, error)
}

package organization

import (
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

// ScopedContext is a relationship-derived organization/sub-scope context primitive
// for later policy evaluation. It is not an authorization decision and grants no
// role, permission or protected action by itself.
type ScopedContext struct {
	tenantID     tenancy.TenantID
	principalID  identity.UserID
	membershipID MembershipID
	scopeID      NodeID
}

// Scope is a tenant-contained organization query/persistence token. It is only
// produced from a valid active scoped relationship context.
type Scope struct {
	tenantID tenancy.TenantID
	nodeID   NodeID
}

func newScopedContext(trusted tenancy.TrustedContext, membership Membership) (ScopedContext, error) {
	if !trusted.Valid() {
		return ScopedContext{}, contextUntrustedFailure()
	}
	if membershipValidationErr := membership.validate(); membershipValidationErr != nil {
		return ScopedContext{}, contextUntrustedFailure()
	}
	if membership.state != MembershipStateActive ||
		membership.tenantID != trusted.TenantID() ||
		membership.principalID != trusted.PrincipalID() {
		return ScopedContext{}, contextUntrustedFailure()
	}
	return ScopedContext{
		tenantID:     membership.tenantID,
		principalID:  membership.principalID,
		membershipID: membership.id,
		scopeID:      membership.scopeID,
	}, nil
}

// Valid reports whether context contains a complete trusted tenant + scoped relationship.
func (scoped ScopedContext) Valid() bool {
	return scoped.tenantID.Valid() &&
		scoped.principalID.Valid() &&
		scoped.membershipID.Valid() &&
		scoped.scopeID.Valid()
}

// TenantID returns the enclosing trusted Tenant.
func (scoped ScopedContext) TenantID() tenancy.TenantID { return scoped.tenantID }

// PrincipalID returns the related P02.01 human User.
func (scoped ScopedContext) PrincipalID() identity.UserID { return scoped.principalID }

// MembershipID returns the active scoped relationship used to derive this context.
func (scoped ScopedContext) MembershipID() MembershipID { return scoped.membershipID }

// NodeID returns the organization/sub-scope node selected by this context.
func (scoped ScopedContext) NodeID() NodeID { return scoped.scopeID }

// ScopeFor returns a scope token only when target exactly matches this relationship.
// The token remains context for later policy evaluation; it is not permission.
func (scoped ScopedContext) ScopeFor(target NodeID) (Scope, error) {
	if !scoped.Valid() || !target.Valid() || target != scoped.scopeID {
		return Scope{}, scopeDeniedFailure()
	}
	return Scope{tenantID: scoped.tenantID, nodeID: scoped.scopeID}, nil
}

// Valid reports whether scope contains canonical Tenant and node identifiers.
func (scope Scope) Valid() bool { return scope.tenantID.Valid() && scope.nodeID.Valid() }

// TenantID returns the enclosing Tenant identifier.
func (scope Scope) TenantID() tenancy.TenantID { return scope.tenantID }

// NodeID returns the selected organization hierarchy identifier.
func (scope Scope) NodeID() NodeID { return scope.nodeID }

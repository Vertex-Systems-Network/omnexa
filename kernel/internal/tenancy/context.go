package tenancy

import "github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"

// TrustedContext is an execution-scoped tenant context derived from an active
// authoritative Tenant/User membership. Its fields are private so callers cannot
// construct tenant authority by copying a client-supplied tenant_id into a struct.
type TrustedContext struct {
	tenantID     TenantID
	principalID  identity.UserID
	membershipID MembershipID
}

// Scope is a tenant-safe persistence/query token. It can only be produced from a
// valid TrustedContext after an explicit target-tenant equality check.
type Scope struct {
	tenantID TenantID
}

func newTrustedContext(tenant Tenant, membership Membership) (TrustedContext, error) {
	if tenantValidationErr := tenant.validate(); tenantValidationErr != nil {
		return TrustedContext{}, contextUntrustedFailure()
	}
	if membershipValidationErr := membership.validate(); membershipValidationErr != nil {
		return TrustedContext{}, contextUntrustedFailure()
	}
	if tenant.state != TenantStateActive || membership.state != MembershipStateActive || membership.tenantID != tenant.id {
		return TrustedContext{}, contextUntrustedFailure()
	}
	return TrustedContext{
		tenantID:     tenant.id,
		principalID:  membership.principalID,
		membershipID: membership.id,
	}, nil
}

// Valid reports whether context contains a complete trusted tenant relationship.
// A valid value should be resolved per execution and must not be persisted or
// reused as a long-lived authorization cache.
func (trusted TrustedContext) Valid() bool {
	return trusted.tenantID.Valid() && trusted.principalID.Valid() && trusted.membershipID.Valid()
}

// TenantID returns the authoritative tenant selected by the resolved relationship.
func (trusted TrustedContext) TenantID() TenantID { return trusted.tenantID }

// PrincipalID returns the trusted P02.01 human User identifier used to resolve context.
func (trusted TrustedContext) PrincipalID() identity.UserID { return trusted.principalID }

// MembershipID returns the authoritative tenancy relationship used to resolve context.
func (trusted TrustedContext) MembershipID() MembershipID { return trusted.membershipID }

// ScopeFor returns a tenant-safe scope only when target exactly matches the
// resolved trusted tenant. A changed/forged target tenant fails closed.
func (trusted TrustedContext) ScopeFor(target TenantID) (Scope, error) {
	if !trusted.Valid() || !target.Valid() || target != trusted.tenantID {
		return Scope{}, crossTenantDeniedFailure()
	}
	return Scope{tenantID: trusted.tenantID}, nil
}

// Valid reports whether scope contains a canonical tenant identifier.
func (scope Scope) Valid() bool { return scope.tenantID.Valid() }

// TenantID returns the tenant identifier carried by this trusted scope token.
func (scope Scope) TenantID() TenantID { return scope.tenantID }

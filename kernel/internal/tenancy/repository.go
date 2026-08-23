package tenancy

import (
	"context"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
)

// Repository is the owner-bounded persistence and trusted-context contract for
// kernel.tenancy. Every relationship operation carries an explicit Tenant scope;
// no global-tenant fallback is available.
type Repository interface {
	CreateTenant(context.Context, Tenant) error
	GetTenant(context.Context, TenantID) (Tenant, error)
	TransitionTenant(context.Context, TenantID, TenantState, TenantState, time.Time) (Tenant, error)
	CreateMembership(context.Context, Membership) error
	RevokeMembership(context.Context, TenantID, MembershipID, time.Time) (Membership, error)
	ResolveContext(context.Context, identity.UserID, TenantID) (TrustedContext, error)
}

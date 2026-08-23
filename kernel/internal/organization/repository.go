package organization

import (
	"context"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

// TenantContextResolver is the narrow P02.02 capability consumed by
// kernel.organization. It proves an active Tenant/User relationship without
// granting organization authority.
type TenantContextResolver interface {
	ResolveContext(context.Context, identity.UserID, tenancy.TenantID) (tenancy.TrustedContext, error)
}

// Repository is the kernel.organization owner-bounded persistence and hierarchy
// contract. Every mutation/read is enclosed by a P02.02 trusted Tenant scope.
type Repository interface {
	CreateNode(context.Context, tenancy.Scope, Node) error
	GetNode(context.Context, tenancy.Scope, NodeID) (Node, error)
	MoveNode(context.Context, tenancy.Scope, NodeID, NodeID, time.Time) (Node, error)
	Ancestors(context.Context, tenancy.Scope, NodeID) ([]Node, error)
	CreateMembership(context.Context, tenancy.Scope, Membership) error
	RevokeMembership(context.Context, tenancy.Scope, MembershipID, time.Time) (Membership, error)
	ResolveContext(context.Context, tenancy.TrustedContext, NodeID) (ScopedContext, error)
}

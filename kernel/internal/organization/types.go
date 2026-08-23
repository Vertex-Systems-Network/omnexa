// Package organization implements the P02.03 tenant-contained organization hierarchy
// and scoped membership relationship foundation.
package organization

import (
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/google/uuid"
)

// NodeID is the stable UUIDv7 identifier of one kernel.organization hierarchy node.
type NodeID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id NodeID) Valid() bool {
	parsed, err := uuid.Parse(string(id))
	return err == nil && parsed.Version() == 7
}

// NodeKind is the canonical organization hierarchy vocabulary authorized by P02.03.
type NodeKind string

const (
	NodeKindOrganization NodeKind = "organization"
	NodeKindLegalEntity  NodeKind = "legal_entity"
	NodeKindBusinessUnit NodeKind = "business_unit"
	NodeKindBranch       NodeKind = "branch"
	NodeKindTeam         NodeKind = "team"
	NodeKindLocation     NodeKind = "location"
)

// Valid reports whether kind is a canonical P02.03 hierarchy type.
func (kind NodeKind) Valid() bool {
	switch kind {
	case NodeKindOrganization, NodeKindLegalEntity, NodeKindBusinessUnit, NodeKindBranch, NodeKindTeam, NodeKindLocation:
		return true
	default:
		return false
	}
}

// parentKindAllowed is deliberately strict. Organization is the tenant-contained
// root; non-root kinds may only descend through structurally compatible scopes.
// Team and Location are terminal P02.03 hierarchy scopes.
func parentKindAllowed(child, parent NodeKind) bool {
	switch child {
	case NodeKindLegalEntity:
		return parent == NodeKindOrganization
	case NodeKindBusinessUnit:
		return parent == NodeKindOrganization || parent == NodeKindLegalEntity
	case NodeKindBranch:
		return parent == NodeKindOrganization || parent == NodeKindLegalEntity || parent == NodeKindBusinessUnit
	case NodeKindTeam, NodeKindLocation:
		return parent == NodeKindOrganization || parent == NodeKindLegalEntity || parent == NodeKindBusinessUnit || parent == NodeKindBranch
	default:
		return false
	}
}

// Node is a kernel.organization-owned hierarchy record. It is an access/operating
// scope primitive and is intentionally not a business Party Organization record.
type Node struct {
	id        NodeID
	tenantID  tenancy.TenantID
	kind      NodeKind
	parentID  *NodeID
	createdAt time.Time
	updatedAt time.Time
}

// NewOrganization creates a root Organization inside one Tenant.
func NewOrganization(tenantID tenancy.TenantID) (Node, error) {
	identifier, identifierErr := uuid.NewV7()
	if identifierErr != nil {
		return Node{}, identifierFailure(identifierErr)
	}
	return newNodeAt(NodeID(identifier.String()), tenantID, NodeKindOrganization, nil, time.Now().UTC())
}

// NewChild creates one non-root hierarchy scope. Parent existence, parent type,
// tenant equality and cycle safety are validated by the owner repository.
func NewChild(tenantID tenancy.TenantID, kind NodeKind, parentID NodeID) (Node, error) {
	if kind == NodeKindOrganization {
		return Node{}, invalidNodeFailure()
	}
	identifier, identifierErr := uuid.NewV7()
	if identifierErr != nil {
		return Node{}, identifierFailure(identifierErr)
	}
	return newNodeAt(NodeID(identifier.String()), tenantID, kind, &parentID, time.Now().UTC())
}

func newNodeAt(id NodeID, tenantID tenancy.TenantID, kind NodeKind, parentID *NodeID, createdAt time.Time) (Node, error) {
	node := Node{
		id:        id,
		tenantID:  tenantID,
		kind:      kind,
		parentID:  cloneNodeID(parentID),
		createdAt: createdAt.UTC(),
		updatedAt: createdAt.UTC(),
	}
	if validationErr := node.validate(); validationErr != nil {
		return Node{}, validationErr
	}
	return node, nil
}

func rehydrateNode(
	id NodeID,
	tenantID tenancy.TenantID,
	kind NodeKind,
	parentID *NodeID,
	createdAt time.Time,
	updatedAt time.Time,
) (Node, error) {
	node := Node{
		id:        id,
		tenantID:  tenantID,
		kind:      kind,
		parentID:  cloneNodeID(parentID),
		createdAt: createdAt.UTC(),
		updatedAt: updatedAt.UTC(),
	}
	if validationErr := node.validate(); validationErr != nil {
		return Node{}, invalidStoredNodeFailure()
	}
	return node, nil
}

func cloneNodeID(id *NodeID) *NodeID {
	if id == nil {
		return nil
	}
	value := *id
	return &value
}

func (node Node) validate() error {
	if !node.id.Valid() || !node.tenantID.Valid() || !node.kind.Valid() {
		return invalidNodeFailure()
	}
	if node.kind == NodeKindOrganization && node.parentID != nil {
		return invalidNodeFailure()
	}
	if node.kind != NodeKindOrganization && (node.parentID == nil || !node.parentID.Valid()) {
		return invalidNodeFailure()
	}
	if node.parentID != nil && *node.parentID == node.id {
		return hierarchyCycleFailure()
	}
	if node.createdAt.IsZero() || node.updatedAt.IsZero() || node.updatedAt.Before(node.createdAt) {
		return invalidNodeFailure()
	}
	return nil
}

// ID returns the stable hierarchy node UUIDv7 identifier.
func (node Node) ID() NodeID { return node.id }

// TenantID returns the enclosing Tenant isolation boundary.
func (node Node) TenantID() tenancy.TenantID { return node.tenantID }

// Kind returns the canonical hierarchy type.
func (node Node) Kind() NodeKind { return node.kind }

// ParentID returns a copy of the parent identifier, or nil for root Organization.
func (node Node) ParentID() *NodeID { return cloneNodeID(node.parentID) }

// CreatedAt returns the canonical UTC creation instant.
func (node Node) CreatedAt() time.Time { return node.createdAt }

// UpdatedAt returns the canonical UTC last hierarchy-change instant.
func (node Node) UpdatedAt() time.Time { return node.updatedAt }

// MembershipID is the stable UUIDv7 identifier of one scoped User relationship.
type MembershipID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id MembershipID) Valid() bool {
	parsed, err := uuid.Parse(string(id))
	return err == nil && parsed.Version() == 7
}

// MembershipState is the deliberately minimal P02.03 scoped relationship state.
type MembershipState string

const (
	MembershipStateActive  MembershipState = "active"
	MembershipStateRevoked MembershipState = "revoked"
)

// Valid reports whether state is a canonical P02.03 relationship state.
func (state MembershipState) Valid() bool {
	return state == MembershipStateActive || state == MembershipStateRevoked
}

// Membership relates a P02.01 human User to one organization hierarchy scope.
// It contains no role, permission, policy, session or employment semantics.
type Membership struct {
	id          MembershipID
	tenantID    tenancy.TenantID
	scopeID     NodeID
	principalID identity.UserID
	state       MembershipState
	createdAt   time.Time
	updatedAt   time.Time
}

// NewMembership creates one active scoped relationship. The repository verifies
// that both the scope and the target User relationship belong to the same Tenant.
func NewMembership(tenantID tenancy.TenantID, scopeID NodeID, principalID identity.UserID) (Membership, error) {
	identifier, identifierErr := uuid.NewV7()
	if identifierErr != nil {
		return Membership{}, identifierFailure(identifierErr)
	}
	return newMembershipAt(
		MembershipID(identifier.String()),
		tenantID,
		scopeID,
		principalID,
		time.Now().UTC(),
	)
}

func newMembershipAt(
	id MembershipID,
	tenantID tenancy.TenantID,
	scopeID NodeID,
	principalID identity.UserID,
	createdAt time.Time,
) (Membership, error) {
	membership := Membership{
		id:          id,
		tenantID:    tenantID,
		scopeID:     scopeID,
		principalID: principalID,
		state:       MembershipStateActive,
		createdAt:   createdAt.UTC(),
		updatedAt:   createdAt.UTC(),
	}
	if validationErr := membership.validate(); validationErr != nil {
		return Membership{}, validationErr
	}
	return membership, nil
}

func rehydrateMembership(
	id MembershipID,
	tenantID tenancy.TenantID,
	scopeID NodeID,
	principalID identity.UserID,
	state MembershipState,
	createdAt time.Time,
	updatedAt time.Time,
) (Membership, error) {
	membership := Membership{
		id:          id,
		tenantID:    tenantID,
		scopeID:     scopeID,
		principalID: principalID,
		state:       state,
		createdAt:   createdAt.UTC(),
		updatedAt:   updatedAt.UTC(),
	}
	if validationErr := membership.validate(); validationErr != nil {
		return Membership{}, invalidStoredMembershipFailure()
	}
	return membership, nil
}

func (membership Membership) validate() error {
	if !membership.id.Valid() || !membership.tenantID.Valid() || !membership.scopeID.Valid() || !membership.principalID.Valid() || !membership.state.Valid() {
		return invalidMembershipFailure()
	}
	if membership.createdAt.IsZero() || membership.updatedAt.IsZero() || membership.updatedAt.Before(membership.createdAt) {
		return invalidMembershipFailure()
	}
	return nil
}

// ID returns the stable scoped relationship identifier.
func (membership Membership) ID() MembershipID { return membership.id }

// TenantID returns the enclosing Tenant isolation boundary.
func (membership Membership) TenantID() tenancy.TenantID { return membership.tenantID }

// ScopeID returns the hierarchy node this relationship references.
func (membership Membership) ScopeID() NodeID { return membership.scopeID }

// PrincipalID returns the P02.01 human User referenced by this relationship.
func (membership Membership) PrincipalID() identity.UserID { return membership.principalID }

// State returns the scoped relationship lifecycle state.
func (membership Membership) State() MembershipState { return membership.state }

// CreatedAt returns the canonical UTC creation instant.
func (membership Membership) CreatedAt() time.Time { return membership.createdAt }

// UpdatedAt returns the canonical UTC last-change instant.
func (membership Membership) UpdatedAt() time.Time { return membership.updatedAt }

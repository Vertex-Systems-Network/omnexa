// Package tenancy implements the P02.02 Tenant lifecycle and trusted tenant-context foundation.
package tenancy

import (
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/google/uuid"
)

// TenantID is the stable UUIDv7 identifier of an Omnexa tenant isolation boundary.
type TenantID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id TenantID) Valid() bool {
	parsed, err := uuid.Parse(string(id))
	return err == nil && parsed.Version() == 7
}

// TenantState is the transport-neutral tenant lifecycle vocabulary.
type TenantState string

const (
	TenantStateProvisioned TenantState = "provisioned"
	TenantStateActive      TenantState = "active"
	TenantStateSuspended   TenantState = "suspended"
	TenantStateDisabled    TenantState = "disabled"
)

// Valid reports whether state is a canonical P02.02 tenant lifecycle state.
func (state TenantState) Valid() bool {
	switch state {
	case TenantStateProvisioned, TenantStateActive, TenantStateSuspended, TenantStateDisabled:
		return true
	default:
		return false
	}
}

// Tenant is the immutable authoritative tenant isolation record.
type Tenant struct {
	id        TenantID
	state     TenantState
	createdAt time.Time
	updatedAt time.Time
}

// NewTenant creates one provisioned Tenant with a UUIDv7 identifier and UTC timestamps.
func NewTenant() (Tenant, error) {
	identifier, identifierErr := uuid.NewV7()
	if identifierErr != nil {
		return Tenant{}, identifierFailure(identifierErr)
	}
	return newTenantAt(TenantID(identifier.String()), time.Now().UTC())
}

func newTenantAt(id TenantID, createdAt time.Time) (Tenant, error) {
	if !id.Valid() || createdAt.IsZero() {
		return Tenant{}, invalidTenantFailure()
	}
	instant := createdAt.UTC()
	return Tenant{
		id:        id,
		state:     TenantStateProvisioned,
		createdAt: instant,
		updatedAt: instant,
	}, nil
}

func rehydrateTenant(id TenantID, state TenantState, createdAt, updatedAt time.Time) (Tenant, error) {
	tenant := Tenant{
		id:        id,
		state:     state,
		createdAt: createdAt.UTC(),
		updatedAt: updatedAt.UTC(),
	}
	if validationErr := tenant.validate(); validationErr != nil {
		return Tenant{}, invalidStoredTenantFailure()
	}
	return tenant, nil
}

// ID returns the stable Tenant UUIDv7 identifier.
func (tenant Tenant) ID() TenantID { return tenant.id }

// State returns the current non-authorization Tenant lifecycle state.
func (tenant Tenant) State() TenantState { return tenant.state }

// CreatedAt returns the canonical UTC creation instant.
func (tenant Tenant) CreatedAt() time.Time { return tenant.createdAt }

// UpdatedAt returns the canonical UTC last lifecycle-change instant.
func (tenant Tenant) UpdatedAt() time.Time { return tenant.updatedAt }

// Transition returns a new Tenant after validating one explicit lifecycle change.
// Same-state transitions and resurrection from disabled fail closed.
func (tenant Tenant) Transition(next TenantState, changedAt time.Time) (Tenant, error) {
	if validationErr := tenant.validate(); validationErr != nil || !next.Valid() || !tenantTransitionAllowed(tenant.state, next) {
		return Tenant{}, tenantTransitionFailure()
	}
	instant := changedAt.UTC()
	if changedAt.IsZero() || instant.Before(tenant.updatedAt) {
		return Tenant{}, tenantTransitionFailure()
	}
	updated := tenant
	updated.state = next
	updated.updatedAt = instant
	return updated, nil
}

func (tenant Tenant) validate() error {
	if !tenant.id.Valid() || !tenant.state.Valid() {
		return invalidTenantFailure()
	}
	if tenant.createdAt.IsZero() || tenant.updatedAt.IsZero() || tenant.updatedAt.Before(tenant.createdAt) {
		return invalidTenantFailure()
	}
	return nil
}

func tenantTransitionAllowed(from, to TenantState) bool {
	switch from {
	case TenantStateProvisioned:
		return to == TenantStateActive || to == TenantStateDisabled
	case TenantStateActive:
		return to == TenantStateSuspended || to == TenantStateDisabled
	case TenantStateSuspended:
		return to == TenantStateActive || to == TenantStateDisabled
	case TenantStateDisabled:
		return false
	default:
		return false
	}
}

// MembershipID is the stable UUIDv7 identifier of one tenant/User relationship.
type MembershipID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id MembershipID) Valid() bool {
	parsed, err := uuid.Parse(string(id))
	return err == nil && parsed.Version() == 7
}

// MembershipState is the deliberately minimal P02.02 tenant relationship state.
type MembershipState string

const (
	MembershipStateActive  MembershipState = "active"
	MembershipStateRevoked MembershipState = "revoked"
)

// Valid reports whether state is a canonical P02.02 membership state.
func (state MembershipState) Valid() bool {
	return state == MembershipStateActive || state == MembershipStateRevoked
}

// Membership is the kernel.tenancy-owned relationship between a P02.01 human User
// and a Tenant. It carries no role, permission, organization or session semantics.
type Membership struct {
	id          MembershipID
	tenantID    TenantID
	principalID identity.UserID
	state       MembershipState
	createdAt   time.Time
	updatedAt   time.Time
}

// NewMembership creates one active Tenant/User relationship. The User identifier
// remains owned by kernel.identity; this record owns only the tenancy relationship.
func NewMembership(tenantID TenantID, principalID identity.UserID) (Membership, error) {
	identifier, identifierErr := uuid.NewV7()
	if identifierErr != nil {
		return Membership{}, identifierFailure(identifierErr)
	}
	return newMembershipAt(MembershipID(identifier.String()), tenantID, principalID, time.Now().UTC())
}

func newMembershipAt(id MembershipID, tenantID TenantID, principalID identity.UserID, createdAt time.Time) (Membership, error) {
	if !id.Valid() || !tenantID.Valid() || !principalID.Valid() || createdAt.IsZero() {
		return Membership{}, invalidMembershipFailure()
	}
	instant := createdAt.UTC()
	return Membership{
		id:          id,
		tenantID:    tenantID,
		principalID: principalID,
		state:       MembershipStateActive,
		createdAt:   instant,
		updatedAt:   instant,
	}, nil
}

func rehydrateMembership(
	id MembershipID,
	tenantID TenantID,
	principalID identity.UserID,
	state MembershipState,
	createdAt time.Time,
	updatedAt time.Time,
) (Membership, error) {
	membership := Membership{
		id:          id,
		tenantID:    tenantID,
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

// ID returns the stable relationship UUIDv7 identifier.
func (membership Membership) ID() MembershipID { return membership.id }

// TenantID returns the authoritative Tenant referenced by this relationship.
func (membership Membership) TenantID() TenantID { return membership.tenantID }

// PrincipalID returns the P02.01 human User identifier referenced by this relationship.
func (membership Membership) PrincipalID() identity.UserID { return membership.principalID }

// State returns the minimal relationship lifecycle state.
func (membership Membership) State() MembershipState { return membership.state }

// CreatedAt returns the canonical UTC relationship creation instant.
func (membership Membership) CreatedAt() time.Time { return membership.createdAt }

// UpdatedAt returns the canonical UTC relationship last-change instant.
func (membership Membership) UpdatedAt() time.Time { return membership.updatedAt }

// Revoke returns a terminal revoked relationship. A later explicit re-binding may
// create a new relationship identifier; revoked records are never silently revived.
func (membership Membership) Revoke(changedAt time.Time) (Membership, error) {
	if validationErr := membership.validate(); validationErr != nil || membership.state != MembershipStateActive {
		return Membership{}, membershipTransitionFailure()
	}
	instant := changedAt.UTC()
	if changedAt.IsZero() || instant.Before(membership.updatedAt) {
		return Membership{}, membershipTransitionFailure()
	}
	updated := membership
	updated.state = MembershipStateRevoked
	updated.updatedAt = instant
	return updated, nil
}

func (membership Membership) validate() error {
	if !membership.id.Valid() || !membership.tenantID.Valid() || !membership.principalID.Valid() || !membership.state.Valid() {
		return invalidMembershipFailure()
	}
	if membership.createdAt.IsZero() || membership.updatedAt.IsZero() || membership.updatedAt.Before(membership.createdAt) {
		return invalidMembershipFailure()
	}
	return nil
}

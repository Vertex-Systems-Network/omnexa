package tenancy

import (
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/google/uuid"
)

const (
	fixedTenantAID    TenantID        = "01890f3e-7b9a-7cc0-98c4-dc0c0c073991"
	fixedTenantBID    TenantID        = "01890f3e-7b9a-7cc0-98c4-dc0c0c073992"
	fixedMembershipID MembershipID    = "01890f3e-7b9a-7cc0-98c4-dc0c0c073993"
	fixedUserID       identity.UserID = "01890f3e-7b9a-7cc0-98c4-dc0c0c07398f"
)

func TestNewTenantUsesUUIDv7UTCAndProvisionedState(t *testing.T) {
	tenant, tenantErr := NewTenant()
	if tenantErr != nil {
		t.Fatalf("NewTenant() error = %v", tenantErr)
	}
	parsed, parseErr := uuid.Parse(string(tenant.ID()))
	if parseErr != nil || parsed.Version() != 7 {
		t.Fatalf("NewTenant() ID = %q, want UUIDv7", tenant.ID())
	}
	if tenant.State() != TenantStateProvisioned {
		t.Fatalf("NewTenant() state = %q, want %q", tenant.State(), TenantStateProvisioned)
	}
	if tenant.CreatedAt().Location() != time.UTC || tenant.UpdatedAt().Location() != time.UTC {
		t.Fatalf("NewTenant() timestamps must be UTC")
	}
}

func TestTenantLifecycleTransitionsAreDeterministicAndFailClosed(t *testing.T) {
	createdAt := time.Date(2026, time.August, 23, 11, 0, 0, 0, time.UTC)
	tenant, tenantErr := newTenantAt(fixedTenantAID, createdAt)
	if tenantErr != nil {
		t.Fatalf("newTenantAt() error = %v", tenantErr)
	}

	active, activeErr := tenant.Transition(TenantStateActive, createdAt.Add(time.Minute))
	if activeErr != nil {
		t.Fatalf("provisioned -> active error = %v", activeErr)
	}
	suspended, suspendedErr := active.Transition(TenantStateSuspended, createdAt.Add(2*time.Minute))
	if suspendedErr != nil {
		t.Fatalf("active -> suspended error = %v", suspendedErr)
	}
	reactivated, reactivateErr := suspended.Transition(TenantStateActive, createdAt.Add(3*time.Minute))
	if reactivateErr != nil {
		t.Fatalf("suspended -> active error = %v", reactivateErr)
	}
	disabled, disableErr := reactivated.Transition(TenantStateDisabled, createdAt.Add(4*time.Minute))
	if disableErr != nil {
		t.Fatalf("active -> disabled error = %v", disableErr)
	}

	if _, transitionErr := disabled.Transition(TenantStateActive, createdAt.Add(5*time.Minute)); !failure.IsCode(transitionErr, codeTenantTransitionInvalid) {
		t.Fatalf("disabled -> active error = %v, want %s", transitionErr, codeTenantTransitionInvalid)
	}
	if _, transitionErr := active.Transition(TenantStateActive, createdAt.Add(2*time.Minute)); !failure.IsCode(transitionErr, codeTenantTransitionInvalid) {
		t.Fatalf("same-state transition error = %v, want %s", transitionErr, codeTenantTransitionInvalid)
	}
	if _, transitionErr := active.Transition(TenantStateSuspended, createdAt); !failure.IsCode(transitionErr, codeTenantTransitionInvalid) {
		t.Fatalf("backdated transition error = %v, want %s", transitionErr, codeTenantTransitionInvalid)
	}
}

func TestMembershipIsMinimalAndRevocationIsTerminal(t *testing.T) {
	createdAt := time.Date(2026, time.August, 23, 11, 0, 0, 0, time.UTC)
	membership, membershipErr := newMembershipAt(fixedMembershipID, fixedTenantAID, fixedUserID, createdAt)
	if membershipErr != nil {
		t.Fatalf("newMembershipAt() error = %v", membershipErr)
	}
	if membership.TenantID() != fixedTenantAID || membership.PrincipalID() != fixedUserID || membership.State() != MembershipStateActive {
		t.Fatalf("membership identity/state mismatch")
	}

	revoked, revokeErr := membership.Revoke(createdAt.Add(time.Minute))
	if revokeErr != nil {
		t.Fatalf("membership.Revoke() error = %v", revokeErr)
	}
	if revoked.State() != MembershipStateRevoked {
		t.Fatalf("revoked membership state = %q", revoked.State())
	}
	if _, secondRevokeErr := revoked.Revoke(createdAt.Add(2 * time.Minute)); !failure.IsCode(secondRevokeErr, codeMembershipTransition) {
		t.Fatalf("second revoke error = %v, want %s", secondRevokeErr, codeMembershipTransition)
	}
}

func TestTrustedContextCreatesScopeOnlyForResolvedTenant(t *testing.T) {
	createdAt := time.Date(2026, time.August, 23, 11, 0, 0, 0, time.UTC)
	tenant, tenantErr := newTenantAt(fixedTenantAID, createdAt)
	if tenantErr != nil {
		t.Fatalf("newTenantAt() error = %v", tenantErr)
	}
	activeTenant, transitionErr := tenant.Transition(TenantStateActive, createdAt.Add(time.Minute))
	if transitionErr != nil {
		t.Fatalf("Tenant.Transition() error = %v", transitionErr)
	}
	membership, membershipErr := newMembershipAt(fixedMembershipID, fixedTenantAID, fixedUserID, createdAt)
	if membershipErr != nil {
		t.Fatalf("newMembershipAt() error = %v", membershipErr)
	}

	trusted, contextErr := newTrustedContext(activeTenant, membership)
	if contextErr != nil {
		t.Fatalf("newTrustedContext() error = %v", contextErr)
	}
	if !trusted.Valid() || trusted.TenantID() != fixedTenantAID || trusted.PrincipalID() != fixedUserID {
		t.Fatalf("trusted context mismatch")
	}
	scope, scopeErr := trusted.ScopeFor(fixedTenantAID)
	if scopeErr != nil || !scope.Valid() || scope.TenantID() != fixedTenantAID {
		t.Fatalf("same-tenant ScopeFor() = (%v, %v)", scope, scopeErr)
	}
	if _, forgedErr := trusted.ScopeFor(fixedTenantBID); !failure.IsCode(forgedErr, codeCrossTenantDenied) {
		t.Fatalf("cross-tenant ScopeFor() error = %v, want %s", forgedErr, codeCrossTenantDenied)
	}
	if _, zeroErr := (TrustedContext{}).ScopeFor(fixedTenantAID); !failure.IsCode(zeroErr, codeCrossTenantDenied) {
		t.Fatalf("zero TrustedContext ScopeFor() error = %v, want %s", zeroErr, codeCrossTenantDenied)
	}
}

func TestTrustedContextRejectsInactiveTenantAndRevokedRelationship(t *testing.T) {
	createdAt := time.Date(2026, time.August, 23, 11, 0, 0, 0, time.UTC)
	tenant, tenantErr := newTenantAt(fixedTenantAID, createdAt)
	if tenantErr != nil {
		t.Fatalf("newTenantAt() error = %v", tenantErr)
	}
	membership, membershipErr := newMembershipAt(fixedMembershipID, fixedTenantAID, fixedUserID, createdAt)
	if membershipErr != nil {
		t.Fatalf("newMembershipAt() error = %v", membershipErr)
	}
	if _, contextErr := newTrustedContext(tenant, membership); !failure.IsCode(contextErr, codeContextUntrusted) {
		t.Fatalf("provisioned tenant context error = %v, want %s", contextErr, codeContextUntrusted)
	}

	activeTenant, transitionErr := tenant.Transition(TenantStateActive, createdAt.Add(time.Minute))
	if transitionErr != nil {
		t.Fatalf("Tenant.Transition() error = %v", transitionErr)
	}
	revoked, revokeErr := membership.Revoke(createdAt.Add(time.Minute))
	if revokeErr != nil {
		t.Fatalf("Membership.Revoke() error = %v", revokeErr)
	}
	if _, contextErr := newTrustedContext(activeTenant, revoked); !failure.IsCode(contextErr, codeContextUntrusted) {
		t.Fatalf("revoked membership context error = %v, want %s", contextErr, codeContextUntrusted)
	}
}

func TestMembershipRejectsInvalidTenantOrIdentity(t *testing.T) {
	createdAt := time.Date(2026, time.August, 23, 11, 0, 0, 0, time.UTC)
	cases := []struct {
		name        string
		tenantID    TenantID
		principalID identity.UserID
	}{
		{name: "missing tenant", tenantID: "", principalID: fixedUserID},
		{name: "non-v7 tenant", tenantID: "00000000-0000-4000-8000-000000000000", principalID: fixedUserID},
		{name: "missing principal", tenantID: fixedTenantAID, principalID: ""},
		{name: "non-v7 principal", tenantID: fixedTenantAID, principalID: "00000000-0000-4000-8000-000000000000"},
	}
	for _, testCase := range cases {
		_, membershipErr := newMembershipAt(fixedMembershipID, testCase.tenantID, testCase.principalID, createdAt)
		if !failure.IsCode(membershipErr, codeMembershipInvalid) {
			t.Fatalf("%s: newMembershipAt() error = %v, want %s", testCase.name, membershipErr, codeMembershipInvalid)
		}
	}
}

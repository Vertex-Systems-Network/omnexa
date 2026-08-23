package organization

import (
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/google/uuid"
)

func testTenantID(t *testing.T) tenancy.TenantID {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("uuid.NewV7() error = %v", err)
	}
	return tenancy.TenantID(value.String())
}

func testNodeID(t *testing.T) NodeID {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("uuid.NewV7() error = %v", err)
	}
	return NodeID(value.String())
}

func testUserID(t *testing.T) identity.UserID {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("uuid.NewV7() error = %v", err)
	}
	return identity.UserID(value.String())
}

func TestNodeKindsAreBoundedToP0203Vocabulary(t *testing.T) {
	valid := []NodeKind{
		NodeKindOrganization,
		NodeKindLegalEntity,
		NodeKindBusinessUnit,
		NodeKindBranch,
		NodeKindTeam,
		NodeKindLocation,
	}
	for _, kind := range valid {
		if !kind.Valid() {
			t.Fatalf("NodeKind(%q).Valid() = false", kind)
		}
	}
	for _, kind := range []NodeKind{"company", "workspace", "party_organization", "warehouse", "department", ""} {
		if kind.Valid() {
			t.Fatalf("unauthorized NodeKind(%q).Valid() = true", kind)
		}
	}
}

func TestHierarchyShapeAndTransitionsFailClosed(t *testing.T) {
	tenantID := testTenantID(t)
	createdAt := time.Date(2026, time.August, 23, 13, 0, 0, 0, time.UTC)
	rootID := testNodeID(t)

	root, rootErr := newNodeAt(rootID, tenantID, NodeKindOrganization, nil, createdAt)
	if rootErr != nil {
		t.Fatalf("newNodeAt(root) error = %v", rootErr)
	}
	if root.ParentID() != nil {
		t.Fatalf("root ParentID() = %v, want nil", root.ParentID())
	}

	parentID := testNodeID(t)
	if _, err := newNodeAt(testNodeID(t), tenantID, NodeKindOrganization, &parentID, createdAt); !failure.IsCode(err, codeNodeInvalid) {
		t.Fatalf("Organization with parent error = %v, want %s", err, codeNodeInvalid)
	}
	if _, err := newNodeAt(testNodeID(t), tenantID, NodeKindBranch, nil, createdAt); !failure.IsCode(err, codeNodeInvalid) {
		t.Fatalf("Branch without parent error = %v, want %s", err, codeNodeInvalid)
	}
	selfID := testNodeID(t)
	if _, err := newNodeAt(selfID, tenantID, NodeKindTeam, &selfID, createdAt); !failure.IsCode(err, codeHierarchyCycle) {
		t.Fatalf("self-parent error = %v, want %s", err, codeHierarchyCycle)
	}

	allowed := []struct {
		child  NodeKind
		parent NodeKind
	}{
		{NodeKindLegalEntity, NodeKindOrganization},
		{NodeKindBusinessUnit, NodeKindOrganization},
		{NodeKindBusinessUnit, NodeKindLegalEntity},
		{NodeKindBranch, NodeKindOrganization},
		{NodeKindBranch, NodeKindLegalEntity},
		{NodeKindBranch, NodeKindBusinessUnit},
		{NodeKindTeam, NodeKindBranch},
		{NodeKindLocation, NodeKindBranch},
	}
	for _, pair := range allowed {
		if !parentKindAllowed(pair.child, pair.parent) {
			t.Fatalf("parentKindAllowed(%q, %q) = false", pair.child, pair.parent)
		}
	}
	for _, pair := range []struct {
		child  NodeKind
		parent NodeKind
	}{
		{NodeKindOrganization, NodeKindOrganization},
		{NodeKindLegalEntity, NodeKindBusinessUnit},
		{NodeKindBusinessUnit, NodeKindBranch},
		{NodeKindBranch, NodeKindTeam},
		{NodeKindTeam, NodeKindLocation},
		{NodeKindLocation, NodeKindTeam},
	} {
		if parentKindAllowed(pair.child, pair.parent) {
			t.Fatalf("parentKindAllowed(%q, %q) = true for invalid transition", pair.child, pair.parent)
		}
	}
}

func TestMembershipCarriesRelationshipOnly(t *testing.T) {
	createdAt := time.Date(2026, time.August, 23, 13, 10, 0, 0, time.UTC)
	membership, err := newMembershipAt(
		MembershipID(testNodeID(t)),
		testTenantID(t),
		testNodeID(t),
		testUserID(t),
		createdAt,
	)
	if err != nil {
		t.Fatalf("newMembershipAt() error = %v", err)
	}
	if membership.State() != MembershipStateActive {
		t.Fatalf("membership.State() = %q, want active", membership.State())
	}
	if membership.CreatedAt() != createdAt || membership.UpdatedAt() != createdAt {
		t.Fatalf("membership timestamps changed")
	}
}

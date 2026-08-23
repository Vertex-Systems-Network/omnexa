package authorization

import (
	"context"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

func TestRoleCompositionIsDeterministicAndRoleNameIsNotAuthority(t *testing.T) {
	scope := testTenantScope()
	role, err := newRoleAt(
		RoleID("01890f3e-7b9a-7cc0-98c4-dc0c0c073980"),
		scope,
		"superuser",
		[]PermissionID{PermissionRoleRead, PermissionAssignmentRead, PermissionRoleRead},
		time.Unix(1_700_000_000, 0).UTC(),
	)
	if err != nil {
		t.Fatalf("newRoleAt() error = %v", err)
	}
	want := []PermissionID{PermissionAssignmentRead, PermissionRoleRead}
	got := role.Permissions()
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("Permissions() = %v, want %v", got, want)
	}
	if role.HasPermission(PermissionAssignmentManage) {
		t.Fatal("role name superuser unexpectedly created bypass authority")
	}
}

func TestServiceDeniesByDefaultAndRejectsExactScopeSubstitution(t *testing.T) {
	repository := newFakeRepository()
	writer, _ := testAuditWriter(t)
	service, err := NewService(repository, writer)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	subjectA := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073981", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	subjectB := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073982", "01890f3e-7b9a-7cc0-98c4-dc0c0c073992")
	decision, err := service.Check(context.Background(), subjectA, PermissionRoleRead)
	if err != nil || decision != DecisionDeny {
		t.Fatalf("unassigned Check() = %q, %v; want deny, nil", decision, err)
	}
	if err := service.Require(context.Background(), subjectA, PermissionRoleRead); !failure.IsCode(err, codeDenied) {
		t.Fatalf("Require() error = %v, want %s", err, codeDenied)
	}
	if subjectA.Scope().Equal(subjectB.Scope()) {
		t.Fatal("different tenant scopes unexpectedly compare equal")
	}
	organizationScope := Scope{
		kind:           ScopeOrganization,
		tenantID:       subjectA.Scope().TenantID(),
		organizationID: organization.NodeID("01890f3e-7b9a-7cc0-98c4-dc0c0c0739a1"),
	}
	if subjectA.Scope().Equal(organizationScope) {
		t.Fatal("tenant scope unexpectedly inherits into organization scope")
	}
}

func TestPrivilegedMutationsRequireCapabilityPreventEscalationAndAudit(t *testing.T) {
	repository := newFakeRepository()
	writer, sink := testAuditWriter(t)
	fixed := time.Unix(1_700_000_100, 0).UTC()
	service, err := newServiceWithClock(repository, writer, func() time.Time { return fixed })
	if err != nil {
		t.Fatalf("newServiceWithClock() error = %v", err)
	}
	actor := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073983", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	seedGrant(t, repository, actor, []PermissionID{PermissionRoleManage, PermissionRoleRead})
	metadata := MutationMetadata{CorrelationID: "p02-05-unit", Reason: "exercise privileged role mutation"}

	role, err := service.CreateRole(context.Background(), actor, "admin", []PermissionID{PermissionRoleRead}, metadata)
	if err != nil {
		t.Fatalf("CreateRole(allowed) error = %v", err)
	}
	if role.Name() != "admin" || role.HasPermission(PermissionAssignmentManage) {
		t.Fatal("admin role name changed permission composition")
	}

	_, err = service.CreateRole(context.Background(), actor, "escalated", []PermissionID{PermissionAssignmentManage}, metadata)
	if !failure.IsCode(err, codeDenied) {
		t.Fatalf("CreateRole(escalation) error = %v, want %s", err, codeDenied)
	}
	if sink.Len() != 2 {
		t.Fatalf("audit record count = %d, want 2", sink.Len())
	}
	records := sink.Snapshot()
	if records[0].Outcome() != audit.OutcomeSucceeded || records[1].Outcome() != audit.OutcomeDenied {
		t.Fatalf("audit outcomes = %q, %q", records[0].Outcome(), records[1].Outcome())
	}
	for _, record := range records {
		if !record.Privileged() || record.Scope().TenantID != string(actor.Scope().TenantID()) || len(record.Fields()) != 0 {
			t.Fatalf("unsafe or incomplete audit record: action=%s scope=%+v fields=%v", record.Action(), record.Scope(), record.Fields())
		}
	}
}

func TestAssignmentRequiresExactActorTargetScopeAndRevocationRemovesAuthority(t *testing.T) {
	repository := newFakeRepository()
	writer, sink := testAuditWriter(t)
	clock := time.Unix(1_700_000_200, 0).UTC()
	service, err := newServiceWithClock(repository, writer, func() time.Time {
		clock = clock.Add(time.Second)
		return clock
	})
	if err != nil {
		t.Fatalf("newServiceWithClock() error = %v", err)
	}
	actor := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073984", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	target := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073985", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	wrongTenant := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073986", "01890f3e-7b9a-7cc0-98c4-dc0c0c073992")
	seedGrant(t, repository, actor, []PermissionID{PermissionAssignmentManage, PermissionRoleRead})
	role, roleErr := newRole(actor.Scope(), "reader", []PermissionID{PermissionRoleRead}, clock)
	if roleErr != nil {
		t.Fatalf("newRole() error = %v", roleErr)
	}
	if err := repository.createRole(context.Background(), role); err != nil {
		t.Fatalf("seed role error = %v", err)
	}
	metadata := MutationMetadata{CorrelationID: "p02-05-assignment", Reason: "exercise direct assignment lifecycle"}

	if _, err := service.AssignRole(context.Background(), actor, wrongTenant, role.ID(), metadata); !failure.IsCode(err, codeScopeDenied) {
		t.Fatalf("cross-tenant AssignRole() error = %v, want %s", err, codeScopeDenied)
	}
	assignment, err := service.AssignRole(context.Background(), actor, target, role.ID(), metadata)
	if err != nil {
		t.Fatalf("AssignRole() error = %v", err)
	}
	decision, err := service.Check(context.Background(), target, PermissionRoleRead)
	if err != nil || decision != DecisionAllow {
		t.Fatalf("assigned Check() = %q, %v; want allow, nil", decision, err)
	}
	if _, err := service.RevokeAssignment(context.Background(), actor, assignment.ID(), metadata); err != nil {
		t.Fatalf("RevokeAssignment() error = %v", err)
	}
	decision, err = service.Check(context.Background(), target, PermissionRoleRead)
	if err != nil || decision != DecisionDeny {
		t.Fatalf("revoked Check() = %q, %v; want deny, nil", decision, err)
	}
	if sink.Len() != 3 {
		t.Fatalf("audit record count = %d, want 3", sink.Len())
	}
}

func testTenantScope() Scope {
	return Scope{kind: ScopeTenant, tenantID: tenancy.TenantID("01890f3e-7b9a-7cc0-98c4-dc0c0c073991")}
}

func testTenantSubject(principalID, tenantID string) Subject {
	return Subject{
		principalID: identity.UserID(principalID),
		scope:       Scope{kind: ScopeTenant, tenantID: tenancy.TenantID(tenantID)},
	}
}

func testAuditWriter(t *testing.T) (*audit.Writer, *audit.MemorySink) {
	t.Helper()
	sink, err := audit.NewMemorySink(64)
	if err != nil {
		t.Fatalf("audit.NewMemorySink() error = %v", err)
	}
	writer, err := audit.NewWriter(sink, nil)
	if err != nil {
		t.Fatalf("audit.NewWriter() error = %v", err)
	}
	return writer, sink
}

func seedGrant(t *testing.T, repository *fakeRepository, subject Subject, permissions []PermissionID) {
	t.Helper()
	role, err := newRole(subject.Scope(), "fixture authority", permissions, time.Unix(1_700_000_000, 0).UTC())
	if err != nil {
		t.Fatalf("seed newRole() error = %v", err)
	}
	if err := repository.createRole(context.Background(), role); err != nil {
		t.Fatalf("seed createRole() error = %v", err)
	}
	assignment, err := newAssignment(role, subject, time.Unix(1_700_000_001, 0).UTC())
	if err != nil {
		t.Fatalf("seed newAssignment() error = %v", err)
	}
	if err := repository.createAssignment(context.Background(), assignment); err != nil {
		t.Fatalf("seed createAssignment() error = %v", err)
	}
}

type fakeRepository struct {
	roles       map[RoleID]Role
	assignments map[AssignmentID]Assignment
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{roles: make(map[RoleID]Role), assignments: make(map[AssignmentID]Assignment)}
}

func (repository *fakeRepository) permissionsExist(_ context.Context, permissions []PermissionID) (bool, error) {
	for _, permission := range permissions {
		if !BuiltinPermission(permission) {
			return false, nil
		}
	}
	return true, nil
}

func (repository *fakeRepository) createRole(_ context.Context, role Role) error {
	if _, exists := repository.roles[role.ID()]; exists {
		return roleConflictFailure()
	}
	repository.roles[role.ID()] = role
	return nil
}

func (repository *fakeRepository) getRole(_ context.Context, scope Scope, id RoleID) (Role, error) {
	role, exists := repository.roles[id]
	if !exists || !role.Scope().Equal(scope) {
		return Role{}, roleNotFoundFailure()
	}
	return role, nil
}

func (repository *fakeRepository) replaceRolePermissions(_ context.Context, scope Scope, id RoleID, permissions []PermissionID, changedAt time.Time) (Role, error) {
	role, err := repository.getRole(context.Background(), scope, id)
	if err != nil {
		return Role{}, err
	}
	updated, err := rehydrateRole(role.ID(), role.Scope(), role.Name(), permissions, role.CreatedAt(), changedAt)
	if err != nil {
		return Role{}, err
	}
	repository.roles[id] = updated
	return updated, nil
}

func (repository *fakeRepository) createAssignment(_ context.Context, assignment Assignment) error {
	for _, current := range repository.assignments {
		if current.RoleID() == assignment.RoleID() && current.PrincipalID() == assignment.PrincipalID() && current.State() == AssignmentActive {
			return assignmentConflictFailure()
		}
	}
	repository.assignments[assignment.ID()] = assignment
	return nil
}

func (repository *fakeRepository) getAssignment(_ context.Context, scope Scope, id AssignmentID) (Assignment, error) {
	assignment, exists := repository.assignments[id]
	if !exists || !assignment.Scope().Equal(scope) {
		return Assignment{}, assignmentNotFoundFailure()
	}
	return assignment, nil
}

func (repository *fakeRepository) revokeAssignment(_ context.Context, scope Scope, id AssignmentID, changedAt time.Time) (Assignment, error) {
	assignment, err := repository.getAssignment(context.Background(), scope, id)
	if err != nil {
		return Assignment{}, err
	}
	if assignment.State() != AssignmentActive {
		return Assignment{}, assignmentConflictFailure()
	}
	revoked, err := newAssignmentAt(
		assignment.ID(), assignment.RoleID(), assignment.PrincipalID(), assignment.Scope(), AssignmentRevoked,
		assignment.CreatedAt(), changedAt,
	)
	if err != nil {
		return Assignment{}, err
	}
	repository.assignments[id] = revoked
	return revoked, nil
}

func (repository *fakeRepository) hasPermission(_ context.Context, subject Subject, permission PermissionID) (bool, error) {
	for _, assignment := range repository.assignments {
		if assignment.State() != AssignmentActive || assignment.PrincipalID() != subject.PrincipalID() || !assignment.Scope().Equal(subject.Scope()) {
			continue
		}
		role, exists := repository.roles[assignment.RoleID()]
		if exists && role.Scope().Equal(subject.Scope()) && role.HasPermission(permission) {
			return true, nil
		}
	}
	return false, nil
}

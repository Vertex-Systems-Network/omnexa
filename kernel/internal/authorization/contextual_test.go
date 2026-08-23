package authorization

import (
	"context"
	"errors"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

func TestContextServiceAllowsOnlyAfterRBACRelationshipAndConstraint(t *testing.T) {
	repository := newFakeRepository()
	writer, sink := testAuditWriter(t)
	rbac, serviceErr := NewService(repository, writer)
	if serviceErr != nil {
		t.Fatalf("NewService() error = %v", serviceErr)
	}
	subject := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073981", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	seedGrant(t, repository, subject, []PermissionID{PermissionRoleRead})
	object := mustObjectReference(t, "project.task", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a1")
	evidence, evidenceErr := NewTenantRelationshipEvidence(subject.PrincipalID(), subject.Scope().TenantID(), object)
	if evidenceErr != nil {
		t.Fatalf("NewTenantRelationshipEvidence() error = %v", evidenceErr)
	}
	resolution, resolutionErr := RelationshipFound(evidence)
	if resolutionErr != nil {
		t.Fatalf("RelationshipFound() error = %v", resolutionErr)
	}
	resolver := &fixedRelationshipResolver{resolution: resolution}
	evaluator := &fixedConstraintEvaluator{allow: true}
	service := mustContextService(t, rbac, resolver, evaluator)
	boundary := mustCapabilityBoundary(t)
	request := mustContextRequest(t, subject, object, boundary, AccessRead, CallerInteractive, "p02-06-allow")

	decision, decisionErr := service.Check(context.Background(), request)
	if decisionErr != nil || decision != DecisionAllow {
		t.Fatalf("Check() = %q, %v; want allow, nil", decision, decisionErr)
	}
	if resolver.calls != 1 || evaluator.calls != 1 {
		t.Fatalf("resolver/evaluator calls = %d/%d, want 1/1", resolver.calls, evaluator.calls)
	}
	if sink.Len() != 0 {
		t.Fatalf("ordinary successful read emitted %d audit records, want 0", sink.Len())
	}
}

func TestContextServiceMissingPermissionDeniesBeforeRelationshipOrConstraint(t *testing.T) {
	repository := newFakeRepository()
	writer, sink := testAuditWriter(t)
	rbac, serviceErr := NewService(repository, writer)
	if serviceErr != nil {
		t.Fatalf("NewService() error = %v", serviceErr)
	}
	subject := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073982", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	object := mustObjectReference(t, "project.task", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a2")
	resolver := &fixedRelationshipResolver{}
	evaluator := &fixedConstraintEvaluator{allow: true}
	service := mustContextService(t, rbac, resolver, evaluator)
	request := mustContextRequest(t, subject, object, mustCapabilityBoundary(t), AccessRead, CallerInternal, "p02-06-missing-permission")

	decision, decisionErr := service.Check(context.Background(), request)
	if decisionErr != nil || decision != DecisionDeny {
		t.Fatalf("Check() = %q, %v; want deny, nil", decision, decisionErr)
	}
	if resolver.calls != 0 || evaluator.calls != 0 {
		t.Fatalf("missing RBAC authority reached resolver/evaluator: %d/%d", resolver.calls, evaluator.calls)
	}
	assertSafeContextAudit(t, sink, subject, object, audit.OutcomeDenied, false)
}

func TestContextServiceRejectsWrongTenantOrganizationAndObjectRelationship(t *testing.T) {
	boundary := mustCapabilityBoundary(t)
	object := mustObjectReference(t, "project.task", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a3")
	otherObject := mustObjectReference(t, "project.task", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a4")
	principal := identity.UserID("01890f3e-7b9a-7cc0-98c4-dc0c0c073983")
	tenantA := tenancy.TenantID("01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	tenantB := tenancy.TenantID("01890f3e-7b9a-7cc0-98c4-dc0c0c073992")
	orgA := organization.NodeID("01890f3e-7b9a-7cc0-98c4-dc0c0c0739b1")
	orgB := organization.NodeID("01890f3e-7b9a-7cc0-98c4-dc0c0c0739b2")

	tests := []struct {
		name     string
		subject  Subject
		evidence RelationshipEvidence
	}{
		{
			name:     "wrong tenant",
			subject:  testTenantSubject(string(principal), string(tenantA)),
			evidence: mustTenantRelationshipEvidence(t, principal, tenantB, object),
		},
		{
			name:     "wrong organization",
			subject:  testOrganizationSubject(principal, tenantA, orgA),
			evidence: mustOrganizationRelationshipEvidence(t, principal, tenantA, orgB, object),
		},
		{
			name:     "wrong object",
			subject:  testOrganizationSubject(principal, tenantA, orgA),
			evidence: mustOrganizationRelationshipEvidence(t, principal, tenantA, orgA, otherObject),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := newFakeRepository()
			writer, sink := testAuditWriter(t)
			rbac, rbacErr := NewService(repository, writer)
			if rbacErr != nil {
				t.Fatalf("NewService() error = %v", rbacErr)
			}
			seedGrant(t, repository, test.subject, []PermissionID{PermissionRoleRead})
			resolution, relationshipErr := RelationshipFound(test.evidence)
			if relationshipErr != nil {
				t.Fatalf("RelationshipFound() error = %v", relationshipErr)
			}
			resolver := &fixedRelationshipResolver{resolution: resolution}
			evaluator := &fixedConstraintEvaluator{allow: true}
			service := mustContextService(t, rbac, resolver, evaluator)
			request := mustContextRequest(t, test.subject, object, boundary, AccessRead, CallerInteractive, "p02-06-wrong-scope")

			decision, decisionErr := service.Check(context.Background(), request)
			if decisionErr != nil || decision != DecisionDeny {
				t.Fatalf("Check() = %q, %v; want deny, nil", decision, decisionErr)
			}
			if resolver.calls != 1 || evaluator.calls != 0 {
				t.Fatalf("wrong relationship reached constraint evaluator: resolver/evaluator=%d/%d", resolver.calls, evaluator.calls)
			}
			assertSafeContextAudit(t, sink, test.subject, object, audit.OutcomeDenied, false)
		})
	}
}

func TestContextServiceNoRelationshipAndConstraintCanOnlyNarrow(t *testing.T) {
	subject := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073984", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	object := mustObjectReference(t, "project.task", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a5")
	boundary := mustCapabilityBoundary(t)

	t.Run("no relationship denies before constraints", func(t *testing.T) {
		repository := newFakeRepository()
		writer, _ := testAuditWriter(t)
		rbac, _ := NewService(repository, writer)
		seedGrant(t, repository, subject, []PermissionID{PermissionRoleRead})
		resolver := &fixedRelationshipResolver{resolution: NoRelationship()}
		evaluator := &fixedConstraintEvaluator{allow: true}
		service := mustContextService(t, rbac, resolver, evaluator)
		request := mustContextRequest(t, subject, object, boundary, AccessRead, CallerInteractive, "p02-06-no-relationship")
		decision, checkErr := service.Check(context.Background(), request)
		if checkErr != nil || decision != DecisionDeny {
			t.Fatalf("Check() = %q, %v; want deny", decision, checkErr)
		}
		if evaluator.calls != 0 {
			t.Fatalf("constraint evaluator called %d times without relationship", evaluator.calls)
		}
	})

	t.Run("constraint false narrows otherwise valid authority", func(t *testing.T) {
		repository := newFakeRepository()
		writer, _ := testAuditWriter(t)
		rbac, _ := NewService(repository, writer)
		seedGrant(t, repository, subject, []PermissionID{PermissionRoleRead})
		evidence := mustTenantRelationshipEvidence(t, subject.PrincipalID(), subject.Scope().TenantID(), object)
		resolution, _ := RelationshipFound(evidence)
		resolver := &fixedRelationshipResolver{resolution: resolution}
		evaluator := &fixedConstraintEvaluator{allow: false}
		service := mustContextService(t, rbac, resolver, evaluator)
		request := mustContextRequest(t, subject, object, boundary, AccessRead, CallerInteractive, "p02-06-constraint-deny")
		decision, checkErr := service.Check(context.Background(), request)
		if checkErr != nil || decision != DecisionDeny {
			t.Fatalf("Check() = %q, %v; want deny", decision, checkErr)
		}
		if evaluator.calls != 1 {
			t.Fatalf("constraint evaluator calls = %d, want 1", evaluator.calls)
		}
	})
}

func TestContextServiceInternalAndBackgroundOriginsNeverBypass(t *testing.T) {
	for _, origin := range []CallerOrigin{CallerInternal, CallerBackground} {
		t.Run(string(origin), func(t *testing.T) {
			repository := newFakeRepository()
			writer, _ := testAuditWriter(t)
			rbac, _ := NewService(repository, writer)
			subject := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073985", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
			object := mustObjectReference(t, "project.task", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a6")
			resolver := &fixedRelationshipResolver{}
			evaluator := &fixedConstraintEvaluator{allow: true}
			service := mustContextService(t, rbac, resolver, evaluator)
			request := mustContextRequest(t, subject, object, mustCapabilityBoundary(t), AccessRead, origin, "p02-06-origin")
			decision, checkErr := service.Check(context.Background(), request)
			if checkErr != nil || decision != DecisionDeny {
				t.Fatalf("Check(%s) = %q, %v; want deny", origin, decision, checkErr)
			}
			if resolver.calls != 0 || evaluator.calls != 0 {
				t.Fatalf("%s origin bypassed RBAC to resolver/evaluator", origin)
			}
		})
	}
}

func TestContextServiceSensitiveFieldAndExportRequireDistinctPermissionAndAudit(t *testing.T) {
	repository := newFakeRepository()
	writer, sink := testAuditWriter(t)
	rbac, _ := NewService(repository, writer)
	subject := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073986", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	seedGrant(t, repository, subject, []PermissionID{PermissionRoleRead})
	object := mustObjectReference(t, "classified.record", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a7")
	evidence := mustTenantRelationshipEvidence(t, subject.PrincipalID(), subject.Scope().TenantID(), object)
	resolution, _ := RelationshipFound(evidence)
	resolver := &fixedRelationshipResolver{resolution: resolution}
	evaluator := &fixedConstraintEvaluator{allow: true}
	service := mustContextService(t, rbac, resolver, evaluator)
	boundary := mustCapabilityBoundary(t)

	sensitiveRequest := mustContextRequest(t, subject, object, boundary, AccessSensitiveField, CallerInteractive, "p02-06-sensitive")
	decision, checkErr := service.Check(context.Background(), sensitiveRequest)
	if checkErr != nil || decision != DecisionDeny {
		t.Fatalf("read-only sensitive Check() = %q, %v; want deny", decision, checkErr)
	}
	seedGrant(t, repository, subject, []PermissionID{PermissionAssignmentRead})
	decision, checkErr = service.Check(context.Background(), sensitiveRequest)
	if checkErr != nil || decision != DecisionAllow {
		t.Fatalf("sensitive permission Check() = %q, %v; want allow", decision, checkErr)
	}

	exportRequest := mustContextRequest(t, subject, object, boundary, AccessExport, CallerInteractive, "p02-06-export")
	decision, checkErr = service.Check(context.Background(), exportRequest)
	if checkErr != nil || decision != DecisionDeny {
		t.Fatalf("missing export permission Check() = %q, %v; want deny", decision, checkErr)
	}
	if sink.Len() != 3 {
		t.Fatalf("privileged decision audit count = %d, want 3", sink.Len())
	}
	records := sink.Snapshot()
	if records[0].Outcome() != audit.OutcomeDenied || !records[0].Privileged() ||
		records[1].Outcome() != audit.OutcomeSucceeded || !records[1].Privileged() ||
		records[2].Outcome() != audit.OutcomeDenied || !records[2].Privileged() {
		t.Fatalf("unexpected privileged audit outcomes")
	}
	for _, record := range records {
		if len(record.Fields()) != 0 || record.Target().Kind != contextAuditTargetKind || record.Target().Reference != string(object.ID()) {
			t.Fatalf("unsafe contextual audit target/fields: target=%+v fields=%v", record.Target(), record.Fields())
		}
	}
}

func TestContextServiceResolverAndConstraintFailuresFailClosed(t *testing.T) {
	subject := testTenantSubject("01890f3e-7b9a-7cc0-98c4-dc0c0c073987", "01890f3e-7b9a-7cc0-98c4-dc0c0c073991")
	object := mustObjectReference(t, "project.task", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a8")
	boundary := mustCapabilityBoundary(t)

	t.Run("resolver failure", func(t *testing.T) {
		repository := newFakeRepository()
		writer, _ := testAuditWriter(t)
		rbac, _ := NewService(repository, writer)
		seedGrant(t, repository, subject, []PermissionID{PermissionRoleRead})
		resolver := &fixedRelationshipResolver{err: errors.New("synthetic resolver failure")}
		service := mustContextService(t, rbac, resolver, &fixedConstraintEvaluator{allow: true})
		request := mustContextRequest(t, subject, object, boundary, AccessRead, CallerInteractive, "p02-06-resolver-failure")
		decision, checkErr := service.Check(context.Background(), request)
		if decision != DecisionDeny || !failure.IsCode(checkErr, codeRelationshipResolutionFailed) {
			t.Fatalf("Check() = %q, %v; want deny/%s", decision, checkErr, codeRelationshipResolutionFailed)
		}
	})

	t.Run("constraint failure", func(t *testing.T) {
		repository := newFakeRepository()
		writer, _ := testAuditWriter(t)
		rbac, _ := NewService(repository, writer)
		seedGrant(t, repository, subject, []PermissionID{PermissionRoleRead})
		evidence := mustTenantRelationshipEvidence(t, subject.PrincipalID(), subject.Scope().TenantID(), object)
		resolution, _ := RelationshipFound(evidence)
		resolver := &fixedRelationshipResolver{resolution: resolution}
		evaluator := &fixedConstraintEvaluator{err: errors.New("synthetic constraint failure")}
		service := mustContextService(t, rbac, resolver, evaluator)
		request := mustContextRequest(t, subject, object, boundary, AccessRead, CallerInteractive, "p02-06-constraint-failure")
		decision, checkErr := service.Check(context.Background(), request)
		if decision != DecisionDeny || !failure.IsCode(checkErr, codeContextConstraintFailed) {
			t.Fatalf("Check() = %q, %v; want deny/%s", decision, checkErr, codeContextConstraintFailed)
		}
	})
}

func mustObjectReference(t *testing.T, kind, id string) ObjectReference {
	t.Helper()
	object, objectErr := NewObjectReference(ObjectKind(kind), ObjectID(id))
	if objectErr != nil {
		t.Fatalf("NewObjectReference() error = %v", objectErr)
	}
	return object
}

func mustCapabilityBoundary(t *testing.T) CapabilityBoundary {
	t.Helper()
	boundary, boundaryErr := NewCapabilityBoundary(PermissionRoleRead, PermissionAssignmentRead, PermissionAssignmentManage)
	if boundaryErr != nil {
		t.Fatalf("NewCapabilityBoundary() error = %v", boundaryErr)
	}
	return boundary
}

func mustContextRequest(
	t *testing.T,
	subject Subject,
	object ObjectReference,
	boundary CapabilityBoundary,
	access AccessKind,
	origin CallerOrigin,
	correlationID string,
) ContextRequest {
	t.Helper()
	request, requestErr := NewContextRequest(subject, object, boundary, access, origin, ContextMetadata{CorrelationID: correlationID})
	if requestErr != nil {
		t.Fatalf("NewContextRequest() error = %v", requestErr)
	}
	return request
}

func testOrganizationSubject(principal identity.UserID, tenant tenancy.TenantID, organizationID organization.NodeID) Subject {
	return Subject{principalID: principal, scope: Scope{kind: ScopeOrganization, tenantID: tenant, organizationID: organizationID}}
}

func mustTenantRelationshipEvidence(t *testing.T, principal identity.UserID, tenant tenancy.TenantID, object ObjectReference) RelationshipEvidence {
	t.Helper()
	evidence, evidenceErr := NewTenantRelationshipEvidence(principal, tenant, object)
	if evidenceErr != nil {
		t.Fatalf("NewTenantRelationshipEvidence() error = %v", evidenceErr)
	}
	return evidence
}

func mustOrganizationRelationshipEvidence(
	t *testing.T,
	principal identity.UserID,
	tenant tenancy.TenantID,
	organizationID organization.NodeID,
	object ObjectReference,
) RelationshipEvidence {
	t.Helper()
	evidence, evidenceErr := NewOrganizationRelationshipEvidence(principal, tenant, organizationID, object)
	if evidenceErr != nil {
		t.Fatalf("NewOrganizationRelationshipEvidence() error = %v", evidenceErr)
	}
	return evidence
}

func mustContextService(t *testing.T, rbac *Service, resolver RelationshipResolver, evaluator ContextConstraintEvaluator) *ContextService {
	t.Helper()
	service, serviceErr := NewContextService(rbac, resolver, evaluator)
	if serviceErr != nil {
		t.Fatalf("NewContextService() error = %v", serviceErr)
	}
	return service
}

func assertSafeContextAudit(
	t *testing.T,
	sink *audit.MemorySink,
	subject Subject,
	object ObjectReference,
	outcome audit.Outcome,
	privileged bool,
) {
	t.Helper()
	if sink.Len() != 1 {
		t.Fatalf("context audit count = %d, want 1", sink.Len())
	}
	record := sink.Snapshot()[0]
	if record.Action() != contextAuditAction || record.Outcome() != outcome || record.Privileged() != privileged {
		t.Fatalf("context audit action/outcome/privileged = %q/%q/%v", record.Action(), record.Outcome(), record.Privileged())
	}
	if record.Target().Kind != contextAuditTargetKind || record.Target().Reference != string(object.ID()) || len(record.Fields()) != 0 {
		t.Fatalf("unsafe context audit target/fields: target=%+v fields=%v", record.Target(), record.Fields())
	}
	if record.Scope().TenantID != string(subject.Scope().TenantID()) {
		t.Fatalf("audit tenant = %q, want %q", record.Scope().TenantID, subject.Scope().TenantID())
	}
	wantOrganization := ""
	if subject.Scope().Kind() == ScopeOrganization {
		wantOrganization = string(subject.Scope().OrganizationID())
	}
	if record.Scope().OrganizationID != wantOrganization {
		t.Fatalf("audit organization = %q, want %q", record.Scope().OrganizationID, wantOrganization)
	}
}

type fixedRelationshipResolver struct {
	resolution RelationshipResolution
	err        error
	calls      int
}

func (resolver *fixedRelationshipResolver) Resolve(_ context.Context, query RelationshipQuery) (RelationshipResolution, error) {
	resolver.calls++
	if !query.valid() {
		return RelationshipResolution{}, invalidContextFailure()
	}
	if resolver.err != nil {
		return RelationshipResolution{}, resolver.err
	}
	return resolver.resolution, nil
}

type fixedConstraintEvaluator struct {
	allow bool
	err   error
	calls int
}

func (evaluator *fixedConstraintEvaluator) Evaluate(_ context.Context, facts ContextFacts) (bool, error) {
	evaluator.calls++
	if !facts.valid() {
		return false, invalidContextFailure()
	}
	if evaluator.err != nil {
		return false, evaluator.err
	}
	return evaluator.allow, nil
}

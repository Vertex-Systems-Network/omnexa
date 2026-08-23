package authorization

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresContextAwareAuthorizationIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_06_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_06_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	pool, poolErr := pgxpool.New(ctx, databaseURL)
	if poolErr != nil {
		t.Fatalf("pgxpool.New() error = %v", poolErr)
	}
	defer pool.Close()
	if pingErr := pool.Ping(ctx); pingErr != nil {
		t.Fatalf("pool.Ping() error = %v", pingErr)
	}

	foundation, foundationErr := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if foundationErr != nil {
		t.Fatalf("NewMigrator(kernel.foundation) error = %v", foundationErr)
	}
	if runErr := foundation.Run(ctx); runErr != nil {
		t.Fatalf("foundation migration error = %v", runErr)
	}
	resetP0205Database(t, ctx, pool)
	defer resetP0205Database(t, context.Background(), pool)

	runOwnerMigrations(t, ctx, pool, "kernel.identity", []database.Migration{
		{Version: 1, Name: "create_identity_foundation", SQL: readMigration(t, migrationIdentityFoundation)},
		{Version: 2, Name: "create_authentication_sessions", SQL: readMigration(t, migrationIdentityAuthentication)},
	})
	runOwnerMigrations(t, ctx, pool, "kernel.tenancy", []database.Migration{{
		Version: 1, Name: "create_tenancy_foundation", SQL: readMigration(t, migrationTenancyFoundation),
	}})
	runOwnerMigrations(t, ctx, pool, "kernel.organization", []database.Migration{{
		Version: 1, Name: "create_organization_foundation", SQL: readMigration(t, migrationOrganizationFoundation),
	}})
	runOwnerMigrations(t, ctx, pool, "kernel.authorization", []database.Migration{{
		Version: 1, Name: "create_rbac_foundation", SQL: readMigration(t, migrationAuthorizationFoundation),
	}})

	identityRepository, identityRepositoryErr := identity.NewPostgresRepository(pool)
	if identityRepositoryErr != nil {
		t.Fatalf("identity.NewPostgresRepository() error = %v", identityRepositoryErr)
	}
	user, userErr := identity.NewUser("p02-06-context@example.com")
	if userErr != nil {
		t.Fatalf("identity.NewUser() error = %v", userErr)
	}
	if createErr := identityRepository.Create(ctx, user); createErr != nil {
		t.Fatalf("identity Create() error = %v", createErr)
	}
	user, userErr = identityRepository.Transition(ctx, user.ID(), identity.LifecycleProvisioned, identity.LifecycleActive, time.Now().UTC().Add(time.Second))
	if userErr != nil {
		t.Fatalf("identity activate error = %v", userErr)
	}

	tenancyRepository, tenancyRepositoryErr := tenancy.NewPostgresRepository(pool)
	if tenancyRepositoryErr != nil {
		t.Fatalf("tenancy.NewPostgresRepository() error = %v", tenancyRepositoryErr)
	}
	tenantA, tenantAErr := tenancy.NewTenant()
	if tenantAErr != nil {
		t.Fatalf("tenancy.NewTenant(A) error = %v", tenantAErr)
	}
	tenantB, tenantBErr := tenancy.NewTenant()
	if tenantBErr != nil {
		t.Fatalf("tenancy.NewTenant(B) error = %v", tenantBErr)
	}
	for _, tenant := range []tenancy.Tenant{tenantA, tenantB} {
		if createErr := tenancyRepository.CreateTenant(ctx, tenant); createErr != nil {
			t.Fatalf("CreateTenant(%s) error = %v", tenant.ID(), createErr)
		}
		if _, transitionErr := tenancyRepository.TransitionTenant(
			ctx, tenant.ID(), tenancy.TenantStateProvisioned, tenancy.TenantStateActive, time.Now().UTC().Add(2*time.Second),
		); transitionErr != nil {
			t.Fatalf("activate tenant %s error = %v", tenant.ID(), transitionErr)
		}
	}
	membership, membershipErr := tenancy.NewMembership(tenantA.ID(), user.ID())
	if membershipErr != nil {
		t.Fatalf("tenancy.NewMembership() error = %v", membershipErr)
	}
	if createMembershipErr := tenancyRepository.CreateMembership(ctx, membership); createMembershipErr != nil {
		t.Fatalf("tenancy CreateMembership() error = %v", createMembershipErr)
	}
	trusted := mustTenantContext(t, ctx, tenancyRepository, user.ID(), tenantA.ID())
	tenantSubject, subjectErr := SubjectFromTenantContext(trusted)
	if subjectErr != nil {
		t.Fatalf("SubjectFromTenantContext() error = %v", subjectErr)
	}

	organizationRepository, organizationRepositoryErr := organization.NewPostgresRepository(pool, tenancyRepository)
	if organizationRepositoryErr != nil {
		t.Fatalf("organization.NewPostgresRepository() error = %v", organizationRepositoryErr)
	}
	tenantScope, tenantScopeErr := trusted.ScopeFor(tenantA.ID())
	if tenantScopeErr != nil {
		t.Fatalf("trusted.ScopeFor() error = %v", tenantScopeErr)
	}
	orgA, orgAErr := organization.NewOrganization(tenantA.ID())
	if orgAErr != nil {
		t.Fatalf("organization.NewOrganization() error = %v", orgAErr)
	}
	if createNodeErr := organizationRepository.CreateNode(ctx, tenantScope, orgA); createNodeErr != nil {
		t.Fatalf("CreateNode() error = %v", createNodeErr)
	}
	orgMembership, orgMembershipErr := organization.NewMembership(orgA.TenantID(), orgA.ID(), user.ID())
	if orgMembershipErr != nil {
		t.Fatalf("organization.NewMembership() error = %v", orgMembershipErr)
	}
	if createOrgMembershipErr := organizationRepository.CreateMembership(ctx, tenantScope, orgMembership); createOrgMembershipErr != nil {
		t.Fatalf("organization CreateMembership() error = %v", createOrgMembershipErr)
	}
	orgContext, orgContextErr := organizationRepository.ResolveContext(ctx, trusted, orgA.ID())
	if orgContextErr != nil {
		t.Fatalf("organization ResolveContext() error = %v", orgContextErr)
	}
	orgSubject, orgSubjectErr := SubjectFromOrganizationContext(orgContext)
	if orgSubjectErr != nil {
		t.Fatalf("SubjectFromOrganizationContext() error = %v", orgSubjectErr)
	}

	repository, repositoryErr := NewPostgresRepository(pool)
	if repositoryErr != nil {
		t.Fatalf("NewPostgresRepository() error = %v", repositoryErr)
	}
	sink, sinkErr := audit.NewMemorySink(128)
	if sinkErr != nil {
		t.Fatalf("audit.NewMemorySink() error = %v", sinkErr)
	}
	writer, writerErr := audit.NewWriter(sink, nil)
	if writerErr != nil {
		t.Fatalf("audit.NewWriter() error = %v", writerErr)
	}
	rbac, rbacErr := NewService(repository, writer)
	if rbacErr != nil {
		t.Fatalf("NewService() error = %v", rbacErr)
	}
	seedPostgresGrant(t, ctx, repository, tenantSubject, "p02-06 tenant reader", []PermissionID{PermissionRoleRead})
	seedPostgresGrant(t, ctx, repository, orgSubject, "p02-06 organization reader", []PermissionID{PermissionRoleRead})

	object := mustObjectReference(t, "project.task", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739d1")
	otherObject := mustObjectReference(t, "project.task", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739d2")
	boundary := mustCapabilityBoundary(t)
	evaluator := &fixedConstraintEvaluator{allow: true}

	tenantEvidence := mustTenantRelationshipEvidence(t, tenantSubject.PrincipalID(), tenantSubject.Scope().TenantID(), object)
	tenantResolution, _ := RelationshipFound(tenantEvidence)
	tenantResolver := &fixedRelationshipResolver{resolution: tenantResolution}
	tenantPolicy := mustContextService(t, rbac, tenantResolver, evaluator)
	readRequest := mustContextRequest(t, tenantSubject, object, boundary, AccessRead, CallerInteractive, "p02-06-integration-read")
	decision, checkErr := tenantPolicy.Check(ctx, readRequest)
	if checkErr != nil || decision != DecisionAllow {
		t.Fatalf("same-tenant contextual Check() = %q, %v; want allow", decision, checkErr)
	}

	wrongTenantEvidence := mustTenantRelationshipEvidence(t, tenantSubject.PrincipalID(), tenantB.ID(), object)
	wrongTenantResolution, _ := RelationshipFound(wrongTenantEvidence)
	wrongTenantPolicy := mustContextService(t, rbac, &fixedRelationshipResolver{resolution: wrongTenantResolution}, &fixedConstraintEvaluator{allow: true})
	decision, checkErr = wrongTenantPolicy.Check(ctx, readRequest)
	if checkErr != nil || decision != DecisionDeny {
		t.Fatalf("wrong-tenant relationship Check() = %q, %v; want deny", decision, checkErr)
	}

	wrongObjectEvidence := mustTenantRelationshipEvidence(t, tenantSubject.PrincipalID(), tenantA.ID(), otherObject)
	wrongObjectResolution, _ := RelationshipFound(wrongObjectEvidence)
	wrongObjectPolicy := mustContextService(t, rbac, &fixedRelationshipResolver{resolution: wrongObjectResolution}, &fixedConstraintEvaluator{allow: true})
	decision, checkErr = wrongObjectPolicy.Check(ctx, readRequest)
	if checkErr != nil || decision != DecisionDeny {
		t.Fatalf("wrong-object relationship Check() = %q, %v; want deny", decision, checkErr)
	}

	orgEvidence := mustOrganizationRelationshipEvidence(t, orgSubject.PrincipalID(), tenantA.ID(), orgA.ID(), object)
	orgResolution, _ := RelationshipFound(orgEvidence)
	orgPolicy := mustContextService(t, rbac, &fixedRelationshipResolver{resolution: orgResolution}, &fixedConstraintEvaluator{allow: true})
	orgRequest := mustContextRequest(t, orgSubject, object, boundary, AccessRead, CallerBackground, "p02-06-integration-org")
	decision, checkErr = orgPolicy.Check(ctx, orgRequest)
	if checkErr != nil || decision != DecisionAllow {
		t.Fatalf("same-org contextual Check(background) = %q, %v; want allow", decision, checkErr)
	}

	sensitiveRequest := mustContextRequest(t, tenantSubject, object, boundary, AccessSensitiveField, CallerInternal, "p02-06-integration-sensitive")
	decision, checkErr = tenantPolicy.Check(ctx, sensitiveRequest)
	if checkErr != nil || decision != DecisionDeny {
		t.Fatalf("ordinary read implied sensitive-field authority: %q, %v", decision, checkErr)
	}
	seedPostgresGrant(t, ctx, repository, tenantSubject, "p02-06 sensitive reader", []PermissionID{PermissionAssignmentRead})
	decision, checkErr = tenantPolicy.Check(ctx, sensitiveRequest)
	if checkErr != nil || decision != DecisionAllow {
		t.Fatalf("explicit sensitive permission Check() = %q, %v; want allow", decision, checkErr)
	}

	var authorizationMigrationCount int
	if migrationErr := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.authorization'",
	).Scan(&authorizationMigrationCount); migrationErr != nil {
		t.Fatalf("authorization migration ledger query error = %v", migrationErr)
	}
	if authorizationMigrationCount != 1 {
		t.Fatalf("P02.06 introduced unexpected authorization migration count = %d, want 1", authorizationMigrationCount)
	}
	if sink.Len() < 4 {
		t.Fatalf("contextual audit count = %d, want at least 4 denial/privileged facts", sink.Len())
	}
	for _, record := range sink.Snapshot() {
		if record.Action() != contextAuditAction || len(record.Fields()) != 0 || record.Target().Kind != contextAuditTargetKind {
			t.Fatalf("unsafe contextual audit record: action=%q target=%+v fields=%v", record.Action(), record.Target(), record.Fields())
		}
	}
}

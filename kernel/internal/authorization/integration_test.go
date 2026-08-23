package authorization

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5/pgxpool"
)

type migrationFixture int

const (
	migrationIdentityFoundation migrationFixture = iota + 1
	migrationIdentityAuthentication
	migrationTenancyFoundation
	migrationOrganizationFoundation
	migrationAuthorizationFoundation
)

func TestPostgresRBACFoundationIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_05_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_05_TEST_DATABASE_URL is not set")
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

	identityV1 := readMigration(t, migrationIdentityFoundation)
	identityV2 := readMigration(t, migrationIdentityAuthentication)
	runOwnerMigrations(t, ctx, pool, "kernel.identity", []database.Migration{
		{Version: 1, Name: "create_identity_foundation", SQL: identityV1},
		{Version: 2, Name: "create_authentication_sessions", SQL: identityV2},
	})
	tenancySQL := readMigration(t, migrationTenancyFoundation)
	runOwnerMigrations(t, ctx, pool, "kernel.tenancy", []database.Migration{{Version: 1, Name: "create_tenancy_foundation", SQL: tenancySQL}})
	organizationSQL := readMigration(t, migrationOrganizationFoundation)
	runOwnerMigrations(t, ctx, pool, "kernel.organization", []database.Migration{{Version: 1, Name: "create_organization_foundation", SQL: organizationSQL}})

	identityRepository, identityRepositoryErr := identity.NewPostgresRepository(pool)
	if identityRepositoryErr != nil {
		t.Fatalf("identity.NewPostgresRepository() error = %v", identityRepositoryErr)
	}
	users := make([]identity.User, 0, 5)
	for _, email := range []string{
		"rbac-admin-a@example.com",
		"rbac-member-a@example.com",
		"rbac-manager-a@example.com",
		"rbac-admin-b@example.com",
		"rbac-member-b@example.com",
	} {
		user, userErr := identity.NewUser(email)
		if userErr != nil {
			t.Fatalf("identity.NewUser(%s) error = %v", email, userErr)
		}
		if createErr := identityRepository.Create(ctx, user); createErr != nil {
			t.Fatalf("identity Create(%s) error = %v", email, createErr)
		}
		active, activeErr := identityRepository.Transition(
			ctx, user.ID(), identity.LifecycleProvisioned, identity.LifecycleActive, time.Now().UTC().Add(time.Second),
		)
		if activeErr != nil {
			t.Fatalf("identity activate(%s) error = %v", email, activeErr)
		}
		users = append(users, active)
	}
	adminA, memberA, managerA, adminB, memberB := users[0], users[1], users[2], users[3], users[4]

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
	for _, pair := range []struct {
		tenant tenancy.Tenant
		user   identity.User
	}{
		{tenantA, adminA}, {tenantA, memberA}, {tenantA, managerA}, {tenantB, adminB}, {tenantB, memberB},
	} {
		membership, membershipErr := tenancy.NewMembership(pair.tenant.ID(), pair.user.ID())
		if membershipErr != nil {
			t.Fatalf("tenancy.NewMembership() error = %v", membershipErr)
		}
		if createErr := tenancyRepository.CreateMembership(ctx, membership); createErr != nil {
			t.Fatalf("tenancy CreateMembership() error = %v", createErr)
		}
	}

	trustedAdminA := mustTenantContext(t, ctx, tenancyRepository, adminA.ID(), tenantA.ID())
	trustedMemberA := mustTenantContext(t, ctx, tenancyRepository, memberA.ID(), tenantA.ID())
	trustedManagerA := mustTenantContext(t, ctx, tenancyRepository, managerA.ID(), tenantA.ID())
	trustedAdminB := mustTenantContext(t, ctx, tenancyRepository, adminB.ID(), tenantB.ID())
	trustedMemberB := mustTenantContext(t, ctx, tenancyRepository, memberB.ID(), tenantB.ID())

	organizationRepository, organizationRepositoryErr := organization.NewPostgresRepository(pool, tenancyRepository)
	if organizationRepositoryErr != nil {
		t.Fatalf("organization.NewPostgresRepository() error = %v", organizationRepositoryErr)
	}
	tenantScopeA, tenantScopeAErr := trustedAdminA.ScopeFor(tenantA.ID())
	if tenantScopeAErr != nil {
		t.Fatalf("trustedAdminA.ScopeFor() error = %v", tenantScopeAErr)
	}
	tenantScopeB, tenantScopeBErr := trustedAdminB.ScopeFor(tenantB.ID())
	if tenantScopeBErr != nil {
		t.Fatalf("trustedAdminB.ScopeFor() error = %v", tenantScopeBErr)
	}
	orgA, orgAErr := organization.NewOrganization(tenantA.ID())
	if orgAErr != nil {
		t.Fatalf("organization.NewOrganization(A) error = %v", orgAErr)
	}
	orgB, orgBErr := organization.NewOrganization(tenantB.ID())
	if orgBErr != nil {
		t.Fatalf("organization.NewOrganization(B) error = %v", orgBErr)
	}
	if createErr := organizationRepository.CreateNode(ctx, tenantScopeA, orgA); createErr != nil {
		t.Fatalf("CreateNode(orgA) error = %v", createErr)
	}
	if createErr := organizationRepository.CreateNode(ctx, tenantScopeB, orgB); createErr != nil {
		t.Fatalf("CreateNode(orgB) error = %v", createErr)
	}
	for _, pair := range []struct {
		trusted tenancy.TrustedContext
		scope   tenancy.Scope
		org     organization.Node
		user    identity.User
	}{
		{trustedAdminA, tenantScopeA, orgA, adminA},
		{trustedMemberA, tenantScopeA, orgA, memberA},
		{trustedAdminB, tenantScopeB, orgB, adminB},
		{trustedMemberB, tenantScopeB, orgB, memberB},
	} {
		membership, membershipErr := organization.NewMembership(pair.org.TenantID(), pair.org.ID(), pair.user.ID())
		if membershipErr != nil {
			t.Fatalf("organization.NewMembership() error = %v", membershipErr)
		}
		if createErr := organizationRepository.CreateMembership(ctx, pair.scope, membership); createErr != nil {
			t.Fatalf("organization CreateMembership() error = %v", createErr)
		}
	}

	authorizationSQL := readMigration(t, migrationAuthorizationFoundation)
	authorizationMigrations := []database.Migration{{Version: 1, Name: "create_rbac_foundation", SQL: authorizationSQL}}
	runOwnerMigrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	runOwnerMigrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	var ledgerCount int
	if ledgerErr := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.authorization' AND version = 1",
	).Scan(&ledgerCount); ledgerErr != nil {
		t.Fatalf("authorization ledger query error = %v", ledgerErr)
	}
	if ledgerCount != 1 {
		t.Fatalf("authorization ledger count = %d, want 1", ledgerCount)
	}
	if _, resolveErr := tenancyRepository.ResolveContext(ctx, adminA.ID(), tenantA.ID()); resolveErr != nil {
		t.Fatalf("P02.05 migration changed P02.02 trusted tenancy state: %v", resolveErr)
	}
	if _, resolveErr := organizationRepository.ResolveContext(ctx, trustedAdminA, orgA.ID()); resolveErr != nil {
		t.Fatalf("P02.05 migration changed P02.03 organization context: %v", resolveErr)
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
	service, serviceErr := NewService(repository, writer)
	if serviceErr != nil {
		t.Fatalf("NewService() error = %v", serviceErr)
	}

	adminATenant, adminATenantErr := SubjectFromTenantContext(trustedAdminA)
	if adminATenantErr != nil {
		t.Fatalf("SubjectFromTenantContext(adminA) error = %v", adminATenantErr)
	}
	memberATenant, memberATenantErr := SubjectFromTenantContext(trustedMemberA)
	if memberATenantErr != nil {
		t.Fatalf("SubjectFromTenantContext(memberA) error = %v", memberATenantErr)
	}
	managerATenant, managerATenantErr := SubjectFromTenantContext(trustedManagerA)
	if managerATenantErr != nil {
		t.Fatalf("SubjectFromTenantContext(managerA) error = %v", managerATenantErr)
	}
	memberBTenant, memberBTenantErr := SubjectFromTenantContext(trustedMemberB)
	if memberBTenantErr != nil {
		t.Fatalf("SubjectFromTenantContext(memberB) error = %v", memberBTenantErr)
	}

	adminAOrgContext, adminAOrgContextErr := organizationRepository.ResolveContext(ctx, trustedAdminA, orgA.ID())
	if adminAOrgContextErr != nil {
		t.Fatalf("ResolveContext(adminA/orgA) error = %v", adminAOrgContextErr)
	}
	memberAOrgContext, memberAOrgContextErr := organizationRepository.ResolveContext(ctx, trustedMemberA, orgA.ID())
	if memberAOrgContextErr != nil {
		t.Fatalf("ResolveContext(memberA/orgA) error = %v", memberAOrgContextErr)
	}
	adminAOrg, adminAOrgErr := SubjectFromOrganizationContext(adminAOrgContext)
	if adminAOrgErr != nil {
		t.Fatalf("SubjectFromOrganizationContext(adminA) error = %v", adminAOrgErr)
	}
	memberAOrg, memberAOrgErr := SubjectFromOrganizationContext(memberAOrgContext)
	if memberAOrgErr != nil {
		t.Fatalf("SubjectFromOrganizationContext(memberA) error = %v", memberAOrgErr)
	}

	allPermissions := []PermissionID{
		PermissionRoleRead, PermissionRoleManage, PermissionAssignmentRead, PermissionAssignmentManage,
	}
	seedPostgresGrant(t, ctx, repository, adminATenant, "fixture tenant A authority", allPermissions)
	seedPostgresGrant(t, ctx, repository, managerATenant, "fixture bounded role manager", []PermissionID{PermissionRoleManage})
	seedPostgresGrant(t, ctx, repository, adminAOrg, "fixture organization A authority", allPermissions)

	metadata := MutationMetadata{CorrelationID: "p02-05-integration", Reason: "synthetic P02.05 authorization evidence"}
	decision, decisionErr := service.Check(ctx, memberATenant, PermissionRoleRead)
	if decisionErr != nil || decision != DecisionDeny {
		t.Fatalf("unassigned tenant member Check() = %q, %v; want deny", decision, decisionErr)
	}
	decision, decisionErr = service.Check(ctx, adminATenant, PermissionRoleManage)
	if decisionErr != nil || decision != DecisionAllow {
		t.Fatalf("bootstrap tenant admin Check() = %q, %v; want allow", decision, decisionErr)
	}

	readerRole, readerRoleErr := service.CreateRole(ctx, adminATenant, "superuser", []PermissionID{PermissionRoleRead}, metadata)
	if readerRoleErr != nil {
		t.Fatalf("CreateRole(superuser name) error = %v", readerRoleErr)
	}
	if readerRole.Name() != "superuser" || readerRole.HasPermission(PermissionAssignmentManage) {
		t.Fatal("superuser role name created implicit bypass authority")
	}
	assignment, assignmentErr := service.AssignRole(ctx, adminATenant, memberATenant, readerRole.ID(), metadata)
	if assignmentErr != nil {
		t.Fatalf("AssignRole(memberA) error = %v", assignmentErr)
	}
	decision, decisionErr = service.Check(ctx, memberATenant, PermissionRoleRead)
	if decisionErr != nil || decision != DecisionAllow {
		t.Fatalf("assigned member Check(role.read) = %q, %v; want allow", decision, decisionErr)
	}
	decision, decisionErr = service.Check(ctx, memberATenant, PermissionAssignmentManage)
	if decisionErr != nil || decision != DecisionDeny {
		t.Fatalf("superuser-name member Check(assignment.manage) = %q, %v; want deny", decision, decisionErr)
	}

	if _, escalationErr := service.CreateRole(ctx, managerATenant, "escalation-attempt", []PermissionID{PermissionAssignmentManage}, metadata); !failure.IsCode(escalationErr, codeDenied) {
		t.Fatalf("bounded manager escalation error = %v, want %s", escalationErr, codeDenied)
	}
	if _, unprivilegedErr := service.AssignRole(ctx, memberATenant, memberATenant, readerRole.ID(), metadata); !failure.IsCode(unprivilegedErr, codeDenied) {
		t.Fatalf("unprivileged assignment error = %v, want %s", unprivilegedErr, codeDenied)
	}
	if _, crossTenantErr := service.AssignRole(ctx, adminATenant, memberBTenant, readerRole.ID(), metadata); !failure.IsCode(crossTenantErr, codeScopeDenied) {
		t.Fatalf("cross-tenant scope substitution error = %v, want %s", crossTenantErr, codeScopeDenied)
	}

	decision, decisionErr = service.Check(ctx, memberAOrg, PermissionRoleRead)
	if decisionErr != nil || decision != DecisionDeny {
		t.Fatalf("tenant grant inherited into organization: decision=%q err=%v", decision, decisionErr)
	}
	orgReader, orgReaderErr := service.CreateRole(ctx, adminAOrg, "organization reader", []PermissionID{PermissionRoleRead}, metadata)
	if orgReaderErr != nil {
		t.Fatalf("CreateRole(org reader) error = %v", orgReaderErr)
	}
	if _, orgAssignErr := service.AssignRole(ctx, adminAOrg, memberAOrg, orgReader.ID(), metadata); orgAssignErr != nil {
		t.Fatalf("AssignRole(memberA org) error = %v", orgAssignErr)
	}
	decision, decisionErr = service.Check(ctx, memberAOrg, PermissionRoleRead)
	if decisionErr != nil || decision != DecisionAllow {
		t.Fatalf("organization direct grant Check() = %q, %v; want allow", decision, decisionErr)
	}

	updatedRole, updateRoleErr := service.ReplaceRolePermissions(
		ctx, adminATenant, readerRole.ID(), []PermissionID{PermissionAssignmentRead, PermissionRoleRead}, metadata,
	)
	if updateRoleErr != nil {
		t.Fatalf("ReplaceRolePermissions() error = %v", updateRoleErr)
	}
	wantPermissions := []PermissionID{PermissionAssignmentRead, PermissionRoleRead}
	gotPermissions := updatedRole.Permissions()
	if len(gotPermissions) != len(wantPermissions) || gotPermissions[0] != wantPermissions[0] || gotPermissions[1] != wantPermissions[1] {
		t.Fatalf("updated permissions = %v, want %v", gotPermissions, wantPermissions)
	}
	decision, decisionErr = service.Check(ctx, memberATenant, PermissionAssignmentRead)
	if decisionErr != nil || decision != DecisionAllow {
		t.Fatalf("updated role did not affect direct grant: decision=%q err=%v", decision, decisionErr)
	}

	if _, revokeErr := service.RevokeAssignment(ctx, adminATenant, assignment.ID(), metadata); revokeErr != nil {
		t.Fatalf("RevokeAssignment() error = %v", revokeErr)
	}
	decision, decisionErr = service.Check(ctx, memberATenant, PermissionRoleRead)
	if decisionErr != nil || decision != DecisionDeny {
		t.Fatalf("revoked role still allows: decision=%q err=%v", decision, decisionErr)
	}
	if _, secondRevokeErr := service.RevokeAssignment(ctx, adminATenant, assignment.ID(), metadata); !failure.IsCode(secondRevokeErr, codeAssignmentConflict) {
		t.Fatalf("second revoke error = %v, want %s", secondRevokeErr, codeAssignmentConflict)
	}

	if sink.Len() < 9 {
		t.Fatalf("audit record count = %d, want at least 9 privileged mutation facts", sink.Len())
	}
	for _, record := range sink.Snapshot() {
		if !record.Privileged() || !strings.HasPrefix(record.Action(), "authorization.") || record.Scope().TenantID == "" || len(record.Fields()) != 0 {
			t.Fatalf("unsafe authorization audit record: action=%q scope=%+v fields=%v", record.Action(), record.Scope(), record.Fields())
		}
		projection := fmt.Sprintf("%v %v %v %v", record.Actor(), record.Target(), record.Scope(), record.Fields())
		if strings.Contains(projection, "superuser") || strings.Contains(projection, "escalation-attempt") {
			t.Fatalf("role display name leaked into audit payload: %s", projection)
		}
	}

	var permissionCount int
	if permissionQueryErr := pool.QueryRow(ctx, "SELECT count(*) FROM omnexa_authorization.permissions").Scan(&permissionCount); permissionQueryErr != nil {
		t.Fatalf("permission count query error = %v", permissionQueryErr)
	}
	if permissionCount != 4 {
		t.Fatalf("P02.05 permission reference count = %d, want 4", permissionCount)
	}
	var futureColumnCount int
	if columnQueryErr := pool.QueryRow(
		ctx,
		`SELECT count(*)
		 FROM information_schema.columns
		 WHERE table_schema = 'omnexa_authorization'
		   AND column_name ~ '(policy|relationship|object|module|mfa|passkey|service_account|api_key)'`,
	).Scan(&futureColumnCount); columnQueryErr != nil {
		t.Fatalf("future-column query error = %v", columnQueryErr)
	}
	if futureColumnCount != 0 {
		t.Fatalf("authorization schema contains %d future-scope columns", futureColumnCount)
	}

	driftMigrator, driftMigratorErr := database.NewMigrator(pool, "kernel.authorization", []database.Migration{{
		Version: 1,
		Name:    "create_rbac_foundation",
		SQL:     authorizationSQL + "\n-- synthetic forbidden rewrite",
	}}, 5*time.Second)
	if driftMigratorErr != nil {
		t.Fatalf("NewMigrator(drift) error = %v", driftMigratorErr)
	}
	if driftErr := driftMigrator.Run(ctx); driftErr == nil {
		t.Fatal("rewritten applied P02.05 migration unexpectedly passed ledger drift check")
	}
}

func mustTenantContext(
	t *testing.T,
	ctx context.Context,
	repository tenancy.Repository,
	principalID identity.UserID,
	tenantID tenancy.TenantID,
) tenancy.TrustedContext {
	t.Helper()
	trusted, resolveErr := repository.ResolveContext(ctx, principalID, tenantID)
	if resolveErr != nil {
		t.Fatalf("ResolveContext(%s/%s) error = %v", principalID, tenantID, resolveErr)
	}
	return trusted
}

func seedPostgresGrant(
	t *testing.T,
	ctx context.Context,
	repository *PostgresRepository,
	subject Subject,
	name string,
	permissions []PermissionID,
) {
	t.Helper()
	role, roleErr := newRole(subject.Scope(), name, permissions, time.Now().UTC())
	if roleErr != nil {
		t.Fatalf("fixture newRole(%s) error = %v", name, roleErr)
	}
	if createRoleErr := repository.createRole(ctx, role); createRoleErr != nil {
		t.Fatalf("fixture createRole(%s) error = %v", name, createRoleErr)
	}
	assignment, assignmentErr := newAssignment(role, subject, time.Now().UTC().Add(time.Millisecond))
	if assignmentErr != nil {
		t.Fatalf("fixture newAssignment(%s) error = %v", name, assignmentErr)
	}
	if createAssignmentErr := repository.createAssignment(ctx, assignment); createAssignmentErr != nil {
		t.Fatalf("fixture createAssignment(%s) error = %v", name, createAssignmentErr)
	}
}

func readMigration(t *testing.T, fixture migrationFixture) string {
	t.Helper()
	var content []byte
	var readErr error
	switch fixture {
	case migrationIdentityFoundation:
		content, readErr = os.ReadFile("../../migrations/kernel.identity/1_create_identity_foundation.sql")
	case migrationIdentityAuthentication:
		content, readErr = os.ReadFile("../../migrations/kernel.identity/2_create_authentication_sessions.sql")
	case migrationTenancyFoundation:
		content, readErr = os.ReadFile("../../migrations/kernel.tenancy/1_create_tenancy_foundation.sql")
	case migrationOrganizationFoundation:
		content, readErr = os.ReadFile("../../migrations/kernel.organization/1_create_organization_foundation.sql")
	case migrationAuthorizationFoundation:
		content, readErr = os.ReadFile("../../migrations/kernel.authorization/1_create_rbac_foundation.sql")
	default:
		t.Fatalf("unknown migration fixture %d", fixture)
	}
	if readErr != nil {
		t.Fatalf("read migration fixture %d error = %v", fixture, readErr)
	}
	return string(content)
}

func runOwnerMigrations(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	owner string,
	migrations []database.Migration,
) {
	t.Helper()
	migrator, migratorErr := database.NewMigrator(pool, owner, migrations, 5*time.Second)
	if migratorErr != nil {
		t.Fatalf("NewMigrator(%s) error = %v", owner, migratorErr)
	}
	if runErr := migrator.Run(ctx); runErr != nil {
		t.Fatalf("%s migrator Run() error = %v", owner, runErr)
	}
}

func resetP0205Database(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	if ctx == nil {
		ctx = context.Background()
	}
	for _, statement := range []string{
		"DROP SCHEMA IF EXISTS omnexa_authorization CASCADE",
		"DROP SCHEMA IF EXISTS omnexa_organization CASCADE",
		"DROP SCHEMA IF EXISTS omnexa_tenancy CASCADE",
		"DROP SCHEMA IF EXISTS omnexa_identity CASCADE",
		"DELETE FROM omnexa_kernel.schema_migrations WHERE owner IN ('kernel.authorization','kernel.organization','kernel.tenancy','kernel.identity')",
	} {
		if _, execErr := pool.Exec(ctx, statement); execErr != nil {
			t.Fatalf("reset statement %q error = %v", statement, execErr)
		}
	}
}

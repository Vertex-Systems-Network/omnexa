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

func TestPostgresRBACFoundationIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_05_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_05_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("pool.Ping() error = %v", err)
	}

	foundation, err := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(kernel.foundation) error = %v", err)
	}
	if err := foundation.Run(ctx); err != nil {
		t.Fatalf("foundation migration error = %v", err)
	}
	resetP0205Database(t, ctx, pool)
	defer resetP0205Database(t, context.Background(), pool)

	identityV1 := readMigration(t, "../../migrations/kernel.identity/1_create_identity_foundation.sql")
	identityV2 := readMigration(t, "../../migrations/kernel.identity/2_create_authentication_sessions.sql")
	runOwnerMigrations(t, ctx, pool, "kernel.identity", []database.Migration{
		{Version: 1, Name: "create_identity_foundation", SQL: identityV1},
		{Version: 2, Name: "create_authentication_sessions", SQL: identityV2},
	})
	tenancySQL := readMigration(t, "../../migrations/kernel.tenancy/1_create_tenancy_foundation.sql")
	runOwnerMigrations(t, ctx, pool, "kernel.tenancy", []database.Migration{{Version: 1, Name: "create_tenancy_foundation", SQL: tenancySQL}})
	organizationSQL := readMigration(t, "../../migrations/kernel.organization/1_create_organization_foundation.sql")
	runOwnerMigrations(t, ctx, pool, "kernel.organization", []database.Migration{{Version: 1, Name: "create_organization_foundation", SQL: organizationSQL}})

	identityRepository, err := identity.NewPostgresRepository(pool)
	if err != nil {
		t.Fatalf("identity.NewPostgresRepository() error = %v", err)
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

	tenancyRepository, err := tenancy.NewPostgresRepository(pool)
	if err != nil {
		t.Fatalf("tenancy.NewPostgresRepository() error = %v", err)
	}
	tenantA, err := tenancy.NewTenant()
	if err != nil {
		t.Fatalf("tenancy.NewTenant(A) error = %v", err)
	}
	tenantB, err := tenancy.NewTenant()
	if err != nil {
		t.Fatalf("tenancy.NewTenant(B) error = %v", err)
	}
	for _, tenant := range []tenancy.Tenant{tenantA, tenantB} {
		if err := tenancyRepository.CreateTenant(ctx, tenant); err != nil {
			t.Fatalf("CreateTenant(%s) error = %v", tenant.ID(), err)
		}
		if _, err := tenancyRepository.TransitionTenant(
			ctx, tenant.ID(), tenancy.TenantStateProvisioned, tenancy.TenantStateActive, time.Now().UTC().Add(2*time.Second),
		); err != nil {
			t.Fatalf("activate tenant %s error = %v", tenant.ID(), err)
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
		if err := tenancyRepository.CreateMembership(ctx, membership); err != nil {
			t.Fatalf("tenancy CreateMembership() error = %v", err)
		}
	}

	trustedAdminA := mustTenantContext(t, ctx, tenancyRepository, adminA.ID(), tenantA.ID())
	trustedMemberA := mustTenantContext(t, ctx, tenancyRepository, memberA.ID(), tenantA.ID())
	trustedManagerA := mustTenantContext(t, ctx, tenancyRepository, managerA.ID(), tenantA.ID())
	trustedAdminB := mustTenantContext(t, ctx, tenancyRepository, adminB.ID(), tenantB.ID())
	trustedMemberB := mustTenantContext(t, ctx, tenancyRepository, memberB.ID(), tenantB.ID())

	organizationRepository, err := organization.NewPostgresRepository(pool, tenancyRepository)
	if err != nil {
		t.Fatalf("organization.NewPostgresRepository() error = %v", err)
	}
	tenantScopeA, err := trustedAdminA.ScopeFor(tenantA.ID())
	if err != nil {
		t.Fatalf("trustedAdminA.ScopeFor() error = %v", err)
	}
	tenantScopeB, err := trustedAdminB.ScopeFor(tenantB.ID())
	if err != nil {
		t.Fatalf("trustedAdminB.ScopeFor() error = %v", err)
	}
	orgA, err := organization.NewOrganization(tenantA.ID())
	if err != nil {
		t.Fatalf("organization.NewOrganization(A) error = %v", err)
	}
	orgB, err := organization.NewOrganization(tenantB.ID())
	if err != nil {
		t.Fatalf("organization.NewOrganization(B) error = %v", err)
	}
	if err := organizationRepository.CreateNode(ctx, tenantScopeA, orgA); err != nil {
		t.Fatalf("CreateNode(orgA) error = %v", err)
	}
	if err := organizationRepository.CreateNode(ctx, tenantScopeB, orgB); err != nil {
		t.Fatalf("CreateNode(orgB) error = %v", err)
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
		if err := organizationRepository.CreateMembership(ctx, pair.scope, membership); err != nil {
			t.Fatalf("organization CreateMembership() error = %v", err)
		}
	}

	authorizationSQL := readMigration(t, "../../migrations/kernel.authorization/1_create_rbac_foundation.sql")
	authorizationMigrations := []database.Migration{{Version: 1, Name: "create_rbac_foundation", SQL: authorizationSQL}}
	runOwnerMigrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	runOwnerMigrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	var ledgerCount int
	if err := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.authorization' AND version = 1",
	).Scan(&ledgerCount); err != nil {
		t.Fatalf("authorization ledger query error = %v", err)
	}
	if ledgerCount != 1 {
		t.Fatalf("authorization ledger count = %d, want 1", ledgerCount)
	}
	if _, err := tenancyRepository.ResolveContext(ctx, adminA.ID(), tenantA.ID()); err != nil {
		t.Fatalf("P02.05 migration changed P02.02 trusted tenancy state: %v", err)
	}
	if _, err := organizationRepository.ResolveContext(ctx, trustedAdminA, orgA.ID()); err != nil {
		t.Fatalf("P02.05 migration changed P02.03 organization context: %v", err)
	}

	repository, err := NewPostgresRepository(pool)
	if err != nil {
		t.Fatalf("NewPostgresRepository() error = %v", err)
	}
	sink, err := audit.NewMemorySink(128)
	if err != nil {
		t.Fatalf("audit.NewMemorySink() error = %v", err)
	}
	writer, err := audit.NewWriter(sink, nil)
	if err != nil {
		t.Fatalf("audit.NewWriter() error = %v", err)
	}
	service, err := NewService(repository, writer)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	adminATenant, err := SubjectFromTenantContext(trustedAdminA)
	if err != nil {
		t.Fatalf("SubjectFromTenantContext(adminA) error = %v", err)
	}
	memberATenant, err := SubjectFromTenantContext(trustedMemberA)
	if err != nil {
		t.Fatalf("SubjectFromTenantContext(memberA) error = %v", err)
	}
	managerATenant, err := SubjectFromTenantContext(trustedManagerA)
	if err != nil {
		t.Fatalf("SubjectFromTenantContext(managerA) error = %v", err)
	}
	memberBTenant, err := SubjectFromTenantContext(trustedMemberB)
	if err != nil {
		t.Fatalf("SubjectFromTenantContext(memberB) error = %v", err)
	}

	adminAOrgContext, err := organizationRepository.ResolveContext(ctx, trustedAdminA, orgA.ID())
	if err != nil {
		t.Fatalf("ResolveContext(adminA/orgA) error = %v", err)
	}
	memberAOrgContext, err := organizationRepository.ResolveContext(ctx, trustedMemberA, orgA.ID())
	if err != nil {
		t.Fatalf("ResolveContext(memberA/orgA) error = %v", err)
	}
	adminAOrg, err := SubjectFromOrganizationContext(adminAOrgContext)
	if err != nil {
		t.Fatalf("SubjectFromOrganizationContext(adminA) error = %v", err)
	}
	memberAOrg, err := SubjectFromOrganizationContext(memberAOrgContext)
	if err != nil {
		t.Fatalf("SubjectFromOrganizationContext(memberA) error = %v", err)
	}

	allPermissions := []PermissionID{
		PermissionRoleRead, PermissionRoleManage, PermissionAssignmentRead, PermissionAssignmentManage,
	}
	seedPostgresGrant(t, ctx, repository, adminATenant, "fixture tenant A authority", allPermissions)
	seedPostgresGrant(t, ctx, repository, managerATenant, "fixture bounded role manager", []PermissionID{PermissionRoleManage})
	seedPostgresGrant(t, ctx, repository, adminAOrg, "fixture organization A authority", allPermissions)

	metadata := MutationMetadata{CorrelationID: "p02-05-integration", Reason: "synthetic P02.05 authorization evidence"}
	decision, err := service.Check(ctx, memberATenant, PermissionRoleRead)
	if err != nil || decision != DecisionDeny {
		t.Fatalf("unassigned tenant member Check() = %q, %v; want deny", decision, err)
	}
	decision, err = service.Check(ctx, adminATenant, PermissionRoleManage)
	if err != nil || decision != DecisionAllow {
		t.Fatalf("bootstrap tenant admin Check() = %q, %v; want allow", decision, err)
	}

	readerRole, err := service.CreateRole(ctx, adminATenant, "superuser", []PermissionID{PermissionRoleRead}, metadata)
	if err != nil {
		t.Fatalf("CreateRole(superuser name) error = %v", err)
	}
	if readerRole.Name() != "superuser" || readerRole.HasPermission(PermissionAssignmentManage) {
		t.Fatal("superuser role name created implicit bypass authority")
	}
	assignment, err := service.AssignRole(ctx, adminATenant, memberATenant, readerRole.ID(), metadata)
	if err != nil {
		t.Fatalf("AssignRole(memberA) error = %v", err)
	}
	decision, err = service.Check(ctx, memberATenant, PermissionRoleRead)
	if err != nil || decision != DecisionAllow {
		t.Fatalf("assigned member Check(role.read) = %q, %v; want allow", decision, err)
	}
	decision, err = service.Check(ctx, memberATenant, PermissionAssignmentManage)
	if err != nil || decision != DecisionDeny {
		t.Fatalf("superuser-name member Check(assignment.manage) = %q, %v; want deny", decision, err)
	}

	if _, err := service.CreateRole(ctx, managerATenant, "escalation-attempt", []PermissionID{PermissionAssignmentManage}, metadata); !failure.IsCode(err, codeDenied) {
		t.Fatalf("bounded manager escalation error = %v, want %s", err, codeDenied)
	}
	if _, err := service.AssignRole(ctx, memberATenant, memberATenant, readerRole.ID(), metadata); !failure.IsCode(err, codeDenied) {
		t.Fatalf("unprivileged assignment error = %v, want %s", err, codeDenied)
	}
	if _, err := service.AssignRole(ctx, adminATenant, memberBTenant, readerRole.ID(), metadata); !failure.IsCode(err, codeScopeDenied) {
		t.Fatalf("cross-tenant scope substitution error = %v, want %s", err, codeScopeDenied)
	}

	decision, err = service.Check(ctx, memberAOrg, PermissionRoleRead)
	if err != nil || decision != DecisionDeny {
		t.Fatalf("tenant grant inherited into organization: decision=%q err=%v", decision, err)
	}
	orgReader, err := service.CreateRole(ctx, adminAOrg, "organization reader", []PermissionID{PermissionRoleRead}, metadata)
	if err != nil {
		t.Fatalf("CreateRole(org reader) error = %v", err)
	}
	if _, err := service.AssignRole(ctx, adminAOrg, memberAOrg, orgReader.ID(), metadata); err != nil {
		t.Fatalf("AssignRole(memberA org) error = %v", err)
	}
	decision, err = service.Check(ctx, memberAOrg, PermissionRoleRead)
	if err != nil || decision != DecisionAllow {
		t.Fatalf("organization direct grant Check() = %q, %v; want allow", decision, err)
	}

	updatedRole, err := service.ReplaceRolePermissions(
		ctx, adminATenant, readerRole.ID(), []PermissionID{PermissionAssignmentRead, PermissionRoleRead}, metadata,
	)
	if err != nil {
		t.Fatalf("ReplaceRolePermissions() error = %v", err)
	}
	wantPermissions := []PermissionID{PermissionAssignmentRead, PermissionRoleRead}
	gotPermissions := updatedRole.Permissions()
	if len(gotPermissions) != len(wantPermissions) || gotPermissions[0] != wantPermissions[0] || gotPermissions[1] != wantPermissions[1] {
		t.Fatalf("updated permissions = %v, want %v", gotPermissions, wantPermissions)
	}
	decision, err = service.Check(ctx, memberATenant, PermissionAssignmentRead)
	if err != nil || decision != DecisionAllow {
		t.Fatalf("updated role did not affect direct grant: decision=%q err=%v", decision, err)
	}

	if _, err := service.RevokeAssignment(ctx, adminATenant, assignment.ID(), metadata); err != nil {
		t.Fatalf("RevokeAssignment() error = %v", err)
	}
	decision, err = service.Check(ctx, memberATenant, PermissionRoleRead)
	if err != nil || decision != DecisionDeny {
		t.Fatalf("revoked role still allows: decision=%q err=%v", decision, err)
	}
	if _, err := service.RevokeAssignment(ctx, adminATenant, assignment.ID(), metadata); !failure.IsCode(err, codeAssignmentConflict) {
		t.Fatalf("second revoke error = %v, want %s", err, codeAssignmentConflict)
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
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM omnexa_authorization.permissions").Scan(&permissionCount); err != nil {
		t.Fatalf("permission count query error = %v", err)
	}
	if permissionCount != 4 {
		t.Fatalf("P02.05 permission reference count = %d, want 4", permissionCount)
	}
	var futureColumnCount int
	if err := pool.QueryRow(
		ctx,
		`SELECT count(*)
		 FROM information_schema.columns
		 WHERE table_schema = 'omnexa_authorization'
		   AND column_name ~ '(policy|relationship|object|module|mfa|passkey|service_account|api_key)'`,
	).Scan(&futureColumnCount); err != nil {
		t.Fatalf("future-column query error = %v", err)
	}
	if futureColumnCount != 0 {
		t.Fatalf("authorization schema contains %d future-scope columns", futureColumnCount)
	}

	driftMigrator, err := database.NewMigrator(pool, "kernel.authorization", []database.Migration{{
		Version: 1,
		Name:    "create_rbac_foundation",
		SQL:     authorizationSQL + "\n-- synthetic forbidden rewrite",
	}}, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(drift) error = %v", err)
	}
	if err := driftMigrator.Run(ctx); err == nil {
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
	trusted, err := repository.ResolveContext(ctx, principalID, tenantID)
	if err != nil {
		t.Fatalf("ResolveContext(%s/%s) error = %v", principalID, tenantID, err)
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
	role, err := newRole(subject.Scope(), name, permissions, time.Now().UTC())
	if err != nil {
		t.Fatalf("fixture newRole(%s) error = %v", name, err)
	}
	if err := repository.createRole(ctx, role); err != nil {
		t.Fatalf("fixture createRole(%s) error = %v", name, err)
	}
	assignment, err := newAssignment(role, subject, time.Now().UTC().Add(time.Millisecond))
	if err != nil {
		t.Fatalf("fixture newAssignment(%s) error = %v", name, err)
	}
	if err := repository.createAssignment(ctx, assignment); err != nil {
		t.Fatalf("fixture createAssignment(%s) error = %v", name, err)
	}
}

func readMigration(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read migration %s error = %v", path, err)
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
	migrator, err := database.NewMigrator(pool, owner, migrations, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(%s) error = %v", owner, err)
	}
	if err := migrator.Run(ctx); err != nil {
		t.Fatalf("%s migrator Run() error = %v", owner, err)
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
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("reset statement %q error = %v", statement, err)
		}
	}
}

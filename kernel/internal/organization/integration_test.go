package organization

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresOrganizationFoundationIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_03_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_03_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	pool, poolErr := pgxpool.New(ctx, databaseURL)
	if poolErr != nil {
		t.Fatalf("pgxpool.New() error = %v", poolErr)
	}
	defer pool.Close()
	if pingErr := pool.Ping(ctx); pingErr != nil {
		t.Fatalf("pool.Ping() error = %v", pingErr)
	}

	foundationMigrator, foundationMigratorErr := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if foundationMigratorErr != nil {
		t.Fatalf("NewMigrator(kernel.foundation) error = %v", foundationMigratorErr)
	}
	if runErr := foundationMigrator.Run(ctx); runErr != nil {
		t.Fatalf("foundation migrator Run() error = %v", runErr)
	}

	resetP0203Database(t, ctx, pool)
	defer resetP0203Database(t, context.Background(), pool)

	identityMigrationBytes, identityMigrationReadErr := os.ReadFile("../../migrations/kernel.identity/1_create_identity_foundation.sql")
	if identityMigrationReadErr != nil {
		t.Fatalf("read identity migration error = %v", identityMigrationReadErr)
	}
	runMigrationSQL(t, ctx, pool, "kernel.identity", 1, "create_identity_foundation", string(identityMigrationBytes))

	tenancyMigrationBytes, tenancyMigrationReadErr := os.ReadFile("../../migrations/kernel.tenancy/1_create_tenancy_foundation.sql")
	if tenancyMigrationReadErr != nil {
		t.Fatalf("read tenancy migration error = %v", tenancyMigrationReadErr)
	}
	runMigrationSQL(t, ctx, pool, "kernel.tenancy", 1, "create_tenancy_foundation", string(tenancyMigrationBytes))

	identityRepository, identityRepositoryErr := identity.NewPostgresRepository(pool)
	if identityRepositoryErr != nil {
		t.Fatalf("identity.NewPostgresRepository() error = %v", identityRepositoryErr)
	}
	userA, userAErr := identity.NewUser("org-a-owner@example.com")
	if userAErr != nil {
		t.Fatalf("identity.NewUser(A) error = %v", userAErr)
	}
	userB, userBErr := identity.NewUser("org-b-owner@example.com")
	if userBErr != nil {
		t.Fatalf("identity.NewUser(B) error = %v", userBErr)
	}
	userA2, userA2Err := identity.NewUser("org-a-member@example.com")
	if userA2Err != nil {
		t.Fatalf("identity.NewUser(A2) error = %v", userA2Err)
	}
	for _, user := range []identity.User{userA, userB, userA2} {
		if createErr := identityRepository.Create(ctx, user); createErr != nil {
			t.Fatalf("identityRepository.Create(%s) error = %v", user.ID(), createErr)
		}
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
	}
	activeAt := time.Now().UTC().Add(time.Second)
	if _, err := tenancyRepository.TransitionTenant(ctx, tenantA.ID(), tenancy.TenantStateProvisioned, tenancy.TenantStateActive, activeAt); err != nil {
		t.Fatalf("activate tenant A error = %v", err)
	}
	if _, err := tenancyRepository.TransitionTenant(ctx, tenantB.ID(), tenancy.TenantStateProvisioned, tenancy.TenantStateActive, activeAt); err != nil {
		t.Fatalf("activate tenant B error = %v", err)
	}

	tenantMemberships := []struct {
		tenant tenancy.Tenant
		user   identity.User
	}{
		{tenantA, userA},
		{tenantA, userA2},
		{tenantB, userB},
	}
	for _, pair := range tenantMemberships {
		membership, err := tenancy.NewMembership(pair.tenant.ID(), pair.user.ID())
		if err != nil {
			t.Fatalf("tenancy.NewMembership() error = %v", err)
		}
		if createErr := tenancyRepository.CreateMembership(ctx, membership); createErr != nil {
			t.Fatalf("tenancy CreateMembership() error = %v", createErr)
		}
	}

	trustedA, trustedAErr := tenancyRepository.ResolveContext(ctx, userA.ID(), tenantA.ID())
	if trustedAErr != nil {
		t.Fatalf("ResolveContext(A) error = %v", trustedAErr)
	}
	scopeA, scopeAErr := trustedA.ScopeFor(tenantA.ID())
	if scopeAErr != nil {
		t.Fatalf("ScopeFor(A) error = %v", scopeAErr)
	}
	trustedB, trustedBErr := tenancyRepository.ResolveContext(ctx, userB.ID(), tenantB.ID())
	if trustedBErr != nil {
		t.Fatalf("ResolveContext(B) error = %v", trustedBErr)
	}
	scopeB, scopeBErr := trustedB.ScopeFor(tenantB.ID())
	if scopeBErr != nil {
		t.Fatalf("ScopeFor(B) error = %v", scopeBErr)
	}

	organizationMigrationBytes, organizationMigrationReadErr := os.ReadFile("../../migrations/kernel.organization/1_create_organization_foundation.sql")
	if organizationMigrationReadErr != nil {
		t.Fatalf("read organization migration error = %v", organizationMigrationReadErr)
	}
	organizationMigrationSQL := string(organizationMigrationBytes)
	organizationMigrations := []database.Migration{{
		Version: 1,
		Name:    "create_organization_foundation",
		SQL:     organizationMigrationSQL,
	}}
	organizationMigrator, organizationMigratorErr := database.NewMigrator(pool, "kernel.organization", organizationMigrations, 5*time.Second)
	if organizationMigratorErr != nil {
		t.Fatalf("NewMigrator(kernel.organization) error = %v", organizationMigratorErr)
	}
	if runErr := organizationMigrator.Run(ctx); runErr != nil {
		t.Fatalf("fresh organization migration error = %v", runErr)
	}
	if runErr := organizationMigrator.Run(ctx); runErr != nil {
		t.Fatalf("idempotent organization migration error = %v", runErr)
	}

	var ledgerCount int
	if queryErr := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.organization' AND version = 1",
	).Scan(&ledgerCount); queryErr != nil {
		t.Fatalf("organization ledger query error = %v", queryErr)
	}
	if ledgerCount != 1 {
		t.Fatalf("organization migration ledger count = %d, want 1", ledgerCount)
	}
	if _, err := tenancyRepository.ResolveContext(ctx, userA.ID(), tenantA.ID()); err != nil {
		t.Fatalf("P02.03 upgrade changed P02.02 tenancy state: %v", err)
	}

	repository, repositoryErr := NewPostgresRepository(pool, tenancyRepository)
	if repositoryErr != nil {
		t.Fatalf("NewPostgresRepository() error = %v", repositoryErr)
	}

	orgA, orgAErr := NewOrganization(tenantA.ID())
	if orgAErr != nil {
		t.Fatalf("NewOrganization(A) error = %v", orgAErr)
	}
	orgB, orgBErr := NewOrganization(tenantB.ID())
	if orgBErr != nil {
		t.Fatalf("NewOrganization(B) error = %v", orgBErr)
	}
	if err := repository.CreateNode(ctx, scopeA, orgA); err != nil {
		t.Fatalf("CreateNode(orgA) error = %v", err)
	}
	if err := repository.CreateNode(ctx, scopeB, orgB); err != nil {
		t.Fatalf("CreateNode(orgB) error = %v", err)
	}

	legalA := mustChild(t, tenantA.ID(), NodeKindLegalEntity, orgA.ID())
	businessUnitA := mustChild(t, tenantA.ID(), NodeKindBusinessUnit, legalA.ID())
	branchA := mustChild(t, tenantA.ID(), NodeKindBranch, businessUnitA.ID())
	teamA := mustChild(t, tenantA.ID(), NodeKindTeam, branchA.ID())
	locationA := mustChild(t, tenantA.ID(), NodeKindLocation, branchA.ID())
	for _, node := range []Node{legalA, businessUnitA, branchA, teamA, locationA} {
		if createErr := repository.CreateNode(ctx, scopeA, node); createErr != nil {
			t.Fatalf("CreateNode(%s/%s) error = %v", node.Kind(), node.ID(), createErr)
		}
	}

	ancestors, ancestorsErr := repository.Ancestors(ctx, scopeA, teamA.ID())
	if ancestorsErr != nil {
		t.Fatalf("Ancestors(teamA) error = %v", ancestorsErr)
	}
	wantKinds := []NodeKind{NodeKindBranch, NodeKindBusinessUnit, NodeKindLegalEntity, NodeKindOrganization}
	if len(ancestors) != len(wantKinds) {
		t.Fatalf("Ancestors(teamA) count = %d, want %d", len(ancestors), len(wantKinds))
	}
	for index, wantKind := range wantKinds {
		if ancestors[index].Kind() != wantKind {
			t.Fatalf("Ancestors(teamA)[%d].Kind() = %q, want %q", index, ancestors[index].Kind(), wantKind)
		}
	}

	crossTenantChild := mustChild(t, tenantA.ID(), NodeKindBranch, orgB.ID())
	if createErr := repository.CreateNode(ctx, scopeA, crossTenantChild); !failure.IsCode(createErr, codeHierarchyParentInvalid) {
		t.Fatalf("cross-tenant parent error = %v, want %s", createErr, codeHierarchyParentInvalid)
	}
	invalidTransition := mustChild(t, tenantA.ID(), NodeKindLegalEntity, businessUnitA.ID())
	if createErr := repository.CreateNode(ctx, scopeA, invalidTransition); !failure.IsCode(createErr, codeHierarchyTransitionInvalid) {
		t.Fatalf("invalid hierarchy transition error = %v, want %s", createErr, codeHierarchyTransitionInvalid)
	}
	if _, moveErr := repository.MoveNode(ctx, scopeA, teamA.ID(), teamA.ID(), time.Now().UTC().Add(2*time.Second)); !failure.IsCode(moveErr, codeHierarchyCycle) {
		t.Fatalf("self-cycle move error = %v, want %s", moveErr, codeHierarchyCycle)
	}
	if _, crossTenantMoveErr := repository.MoveNode(ctx, scopeA, locationA.ID(), orgB.ID(), time.Now().UTC().Add(3*time.Second)); !failure.IsCode(crossTenantMoveErr, codeHierarchyParentInvalid) {
		t.Fatalf("cross-tenant move error = %v, want %s", crossTenantMoveErr, codeHierarchyParentInvalid)
	}
	movedLocation, moveLocationErr := repository.MoveNode(ctx, scopeA, locationA.ID(), businessUnitA.ID(), time.Now().UTC().Add(4*time.Second))
	if moveLocationErr != nil {
		t.Fatalf("valid location move error = %v", moveLocationErr)
	}
	if movedLocation.ParentID() == nil || *movedLocation.ParentID() != businessUnitA.ID() {
		t.Fatalf("moved location parent mismatch")
	}

	scopedMembership, scopedMembershipErr := NewMembership(tenantA.ID(), branchA.ID(), userA.ID())
	if scopedMembershipErr != nil {
		t.Fatalf("NewMembership(A branch) error = %v", scopedMembershipErr)
	}
	if createErr := repository.CreateMembership(ctx, scopeA, scopedMembership); createErr != nil {
		t.Fatalf("CreateMembership(A branch) error = %v", createErr)
	}
	scopedContext, scopedContextErr := repository.ResolveContext(ctx, trustedA, branchA.ID())
	if scopedContextErr != nil {
		t.Fatalf("ResolveContext(A branch) error = %v", scopedContextErr)
	}
	if !scopedContext.Valid() || scopedContext.TenantID() != tenantA.ID() || scopedContext.NodeID() != branchA.ID() {
		t.Fatalf("scoped context mismatch")
	}
	if _, err := scopedContext.ScopeFor(branchA.ID()); err != nil {
		t.Fatalf("ScopedContext.ScopeFor(branchA) error = %v", err)
	}
	if _, err := scopedContext.ScopeFor(teamA.ID()); !failure.IsCode(err, codeScopeDenied) {
		t.Fatalf("ScopedContext.ScopeFor(other scope) error = %v, want %s", err, codeScopeDenied)
	}

	crossTenantMembership, crossTenantMembershipErr := NewMembership(tenantB.ID(), orgB.ID(), userB.ID())
	if crossTenantMembershipErr != nil {
		t.Fatalf("NewMembership(B) error = %v", crossTenantMembershipErr)
	}
	if createErr := repository.CreateMembership(ctx, scopeA, crossTenantMembership); !failure.IsCode(createErr, codeScopeDenied) {
		t.Fatalf("cross-tenant membership error = %v, want %s", createErr, codeScopeDenied)
	}
	wrongTenantPrincipalMembership, wrongTenantPrincipalMembershipErr := NewMembership(tenantA.ID(), branchA.ID(), userB.ID())
	if wrongTenantPrincipalMembershipErr != nil {
		t.Fatalf("NewMembership(A/userB) error = %v", wrongTenantPrincipalMembershipErr)
	}
	if createErr := repository.CreateMembership(ctx, scopeA, wrongTenantPrincipalMembership); !failure.IsCode(createErr, codeScopeDenied) {
		t.Fatalf("wrong-tenant principal membership error = %v, want %s", createErr, codeScopeDenied)
	}

	if _, err := repository.GetNode(ctx, scopeB, branchA.ID()); !failure.IsCode(err, codeNodeNotFound) {
		t.Fatalf("cross-tenant GetNode error = %v, want %s", err, codeNodeNotFound)
	}
	if _, err := repository.RevokeMembership(ctx, scopeB, scopedMembership.ID(), time.Now().UTC().Add(5*time.Second)); !failure.IsCode(err, codeMembershipConflict) {
		t.Fatalf("wrong-tenant revoke error = %v, want %s", err, codeMembershipConflict)
	}
	if _, err := repository.ResolveContext(ctx, trustedA, branchA.ID()); err != nil {
		t.Fatalf("wrong-tenant revoke changed relationship: %v", err)
	}
	if _, err := repository.RevokeMembership(ctx, scopeA, scopedMembership.ID(), time.Now().UTC().Add(6*time.Second)); err != nil {
		t.Fatalf("RevokeMembership(A) error = %v", err)
	}
	if _, err := repository.ResolveContext(ctx, trustedA, branchA.ID()); !failure.IsCode(err, codeContextUntrusted) {
		t.Fatalf("revoked scoped context error = %v, want %s", err, codeContextUntrusted)
	}

	replacementMembership, replacementMembershipErr := NewMembership(tenantA.ID(), branchA.ID(), userA.ID())
	if replacementMembershipErr != nil {
		t.Fatalf("NewMembership(replacement) error = %v", replacementMembershipErr)
	}
	if createErr := repository.CreateMembership(ctx, scopeA, replacementMembership); createErr != nil {
		t.Fatalf("CreateMembership(replacement) error = %v", createErr)
	}
	if _, err := repository.ResolveContext(ctx, trustedA, branchA.ID()); err != nil {
		t.Fatalf("replacement scoped relationship did not restore context: %v", err)
	}
	if _, err := tenancyRepository.TransitionTenant(
		ctx,
		tenantA.ID(),
		tenancy.TenantStateActive,
		tenancy.TenantStateSuspended,
		time.Now().UTC().Add(7*time.Second),
	); err != nil {
		t.Fatalf("suspend tenant A error = %v", err)
	}
	if _, err := repository.ResolveContext(ctx, trustedA, branchA.ID()); !failure.IsCode(err, codeContextUntrusted) {
		t.Fatalf("stale trusted tenant context error = %v, want %s", err, codeContextUntrusted)
	}

	var forbiddenColumns int
	if queryErr := pool.QueryRow(
		ctx,
		`SELECT count(*)
		 FROM information_schema.columns
		 WHERE table_schema = 'omnexa_organization'
		   AND column_name IN (
		       'company_id', 'workspace_id', 'party_id', 'person_id', 'employee_id',
		       'role_id', 'permission_id', 'policy_id', 'session_id', 'password_hash'
		   )`,
	).Scan(&forbiddenColumns); queryErr != nil {
		t.Fatalf("forbidden-column query error = %v", queryErr)
	}
	if forbiddenColumns != 0 {
		t.Fatalf("organization schema contains %d unauthorized business/auth/authz columns", forbiddenColumns)
	}

	drifted := []database.Migration{{
		Version: 1,
		Name:    "create_organization_foundation",
		SQL:     organizationMigrationSQL + "\n-- rewritten migration must fail closed\n",
	}}
	driftMigrator, driftMigratorErr := database.NewMigrator(pool, "kernel.organization", drifted, 5*time.Second)
	if driftMigratorErr != nil {
		t.Fatalf("NewMigrator(drifted organization) error = %v", driftMigratorErr)
	}
	if driftErr := driftMigrator.Run(ctx); driftErr == nil {
		t.Fatalf("rewritten applied organization migration unexpectedly succeeded")
	}
}

func mustChild(t *testing.T, tenantID tenancy.TenantID, kind NodeKind, parentID NodeID) Node {
	t.Helper()
	node, err := NewChild(tenantID, kind, parentID)
	if err != nil {
		t.Fatalf("NewChild(%s) error = %v", kind, err)
	}
	return node
}

func runMigrationSQL(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	owner string,
	version int64,
	name string,
	sql string,
) {
	t.Helper()
	migrations := []database.Migration{{
		Version: version,
		Name:    name,
		SQL:     sql,
	}}
	migrator, err := database.NewMigrator(pool, owner, migrations, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(%s) error = %v", owner, err)
	}
	if err := migrator.Run(ctx); err != nil {
		t.Fatalf("%s migration error = %v", owner, err)
	}
}

func resetP0203Database(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	cleanupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := pool.Exec(cleanupCtx, "DROP SCHEMA IF EXISTS omnexa_organization CASCADE"); err != nil {
		t.Fatalf("drop organization schema error = %v", err)
	}
	if _, err := pool.Exec(cleanupCtx, "DROP SCHEMA IF EXISTS omnexa_tenancy CASCADE"); err != nil {
		t.Fatalf("drop tenancy schema error = %v", err)
	}
	if _, err := pool.Exec(cleanupCtx, "DROP SCHEMA IF EXISTS omnexa_identity CASCADE"); err != nil {
		t.Fatalf("drop identity schema error = %v", err)
	}
	if _, err := pool.Exec(
		cleanupCtx,
		"DELETE FROM omnexa_kernel.schema_migrations WHERE owner IN ('kernel.organization', 'kernel.tenancy', 'kernel.identity')",
	); err != nil {
		t.Fatalf("reset P02 migration ledger error = %v", err)
	}
}

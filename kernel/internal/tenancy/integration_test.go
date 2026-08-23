package tenancy

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresTenancyFoundationIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_02_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_02_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
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

	if _, dropTenancyErr := pool.Exec(ctx, "DROP SCHEMA IF EXISTS omnexa_tenancy CASCADE"); dropTenancyErr != nil {
		t.Fatalf("drop tenancy schema error = %v", dropTenancyErr)
	}
	if _, dropIdentityErr := pool.Exec(ctx, "DROP SCHEMA IF EXISTS omnexa_identity CASCADE"); dropIdentityErr != nil {
		t.Fatalf("drop identity schema error = %v", dropIdentityErr)
	}
	if _, resetLedgerErr := pool.Exec(
		ctx,
		"DELETE FROM omnexa_kernel.schema_migrations WHERE owner IN ('kernel.tenancy', 'kernel.identity')",
	); resetLedgerErr != nil {
		t.Fatalf("reset P02 migration ledger error = %v", resetLedgerErr)
	}
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_, _ = pool.Exec(cleanupCtx, "DROP SCHEMA IF EXISTS omnexa_tenancy CASCADE")
		_, _ = pool.Exec(cleanupCtx, "DROP SCHEMA IF EXISTS omnexa_identity CASCADE")
		_, _ = pool.Exec(
			cleanupCtx,
			"DELETE FROM omnexa_kernel.schema_migrations WHERE owner IN ('kernel.tenancy', 'kernel.identity')",
		)
	}()

	identityMigrationSQL, identityMigrationReadErr := os.ReadFile("../../migrations/kernel.identity/1_create_identity_foundation.sql")
	if identityMigrationReadErr != nil {
		t.Fatalf("read identity migration error = %v", identityMigrationReadErr)
	}
	identityMigrations := []database.Migration{{
		Version: 1,
		Name:    "create_identity_foundation",
		SQL:     string(identityMigrationSQL),
	}}
	identityMigrator, identityMigratorErr := database.NewMigrator(pool, "kernel.identity", identityMigrations, 5*time.Second)
	if identityMigratorErr != nil {
		t.Fatalf("NewMigrator(kernel.identity) error = %v", identityMigratorErr)
	}
	if runErr := identityMigrator.Run(ctx); runErr != nil {
		t.Fatalf("identity prerequisite migration error = %v", runErr)
	}

	identityRepository, identityRepositoryErr := identity.NewPostgresRepository(pool)
	if identityRepositoryErr != nil {
		t.Fatalf("identity.NewPostgresRepository() error = %v", identityRepositoryErr)
	}
	user, userErr := identity.NewUser("tenant-owner@example.com")
	if userErr != nil {
		t.Fatalf("identity.NewUser() error = %v", userErr)
	}
	if createUserErr := identityRepository.Create(ctx, user); createUserErr != nil {
		t.Fatalf("identity repository Create() error = %v", createUserErr)
	}

	tenancyMigrationSQL, tenancyMigrationReadErr := os.ReadFile("../../migrations/kernel.tenancy/1_create_tenancy_foundation.sql")
	if tenancyMigrationReadErr != nil {
		t.Fatalf("read tenancy migration error = %v", tenancyMigrationReadErr)
	}
	tenancyMigrations := []database.Migration{{
		Version: 1,
		Name:    "create_tenancy_foundation",
		SQL:     string(tenancyMigrationSQL),
	}}
	tenancyMigrator, tenancyMigratorErr := database.NewMigrator(pool, "kernel.tenancy", tenancyMigrations, 5*time.Second)
	if tenancyMigratorErr != nil {
		t.Fatalf("NewMigrator(kernel.tenancy) error = %v", tenancyMigratorErr)
	}
	if runErr := tenancyMigrator.Run(ctx); runErr != nil {
		t.Fatalf("fresh tenancy migration error = %v", runErr)
	}
	if runErr := tenancyMigrator.Run(ctx); runErr != nil {
		t.Fatalf("idempotent tenancy migration error = %v", runErr)
	}

	var tenancyLedgerCount int
	if queryErr := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.tenancy' AND version = 1",
	).Scan(&tenancyLedgerCount); queryErr != nil {
		t.Fatalf("tenancy ledger query error = %v", queryErr)
	}
	if tenancyLedgerCount != 1 {
		t.Fatalf("tenancy migration ledger count = %d, want 1", tenancyLedgerCount)
	}

	storedUserAfterUpgrade, storedUserErr := identityRepository.Get(ctx, user.ID())
	if storedUserErr != nil {
		t.Fatalf("identity user missing after tenancy upgrade: %v", storedUserErr)
	}
	if storedUserAfterUpgrade.ID() != user.ID() || storedUserAfterUpgrade.PrimaryEmail() != user.PrimaryEmail() {
		t.Fatalf("P02.02 migration changed P02.01 identity state")
	}

	repository, repositoryErr := NewPostgresRepository(pool)
	if repositoryErr != nil {
		t.Fatalf("NewPostgresRepository() error = %v", repositoryErr)
	}
	createdAt := time.Date(2026, time.August, 23, 11, 30, 0, 0, time.UTC)
	tenantA, tenantAErr := newTenantAt(fixedTenantAID, createdAt)
	if tenantAErr != nil {
		t.Fatalf("newTenantAt(A) error = %v", tenantAErr)
	}
	tenantB, tenantBErr := newTenantAt(fixedTenantBID, createdAt)
	if tenantBErr != nil {
		t.Fatalf("newTenantAt(B) error = %v", tenantBErr)
	}
	if createTenantAErr := repository.CreateTenant(ctx, tenantA); createTenantAErr != nil {
		t.Fatalf("CreateTenant(A) error = %v", createTenantAErr)
	}
	if createTenantBErr := repository.CreateTenant(ctx, tenantB); createTenantBErr != nil {
		t.Fatalf("CreateTenant(B) error = %v", createTenantBErr)
	}

	activeAt := createdAt.Add(time.Minute)
	activeTenantA, activateAErr := repository.TransitionTenant(ctx, tenantA.ID(), TenantStateProvisioned, TenantStateActive, activeAt)
	if activateAErr != nil {
		t.Fatalf("TransitionTenant(A active) error = %v", activateAErr)
	}
	if _, activateBErr := repository.TransitionTenant(ctx, tenantB.ID(), TenantStateProvisioned, TenantStateActive, activeAt); activateBErr != nil {
		t.Fatalf("TransitionTenant(B active) error = %v", activateBErr)
	}
	if activeTenantA.State() != TenantStateActive {
		t.Fatalf("Tenant A state = %q, want active", activeTenantA.State())
	}

	membership, membershipErr := newMembershipAt(fixedMembershipID, tenantA.ID(), user.ID(), activeAt)
	if membershipErr != nil {
		t.Fatalf("newMembershipAt() error = %v", membershipErr)
	}
	if createMembershipErr := repository.CreateMembership(ctx, membership); createMembershipErr != nil {
		t.Fatalf("CreateMembership() error = %v", createMembershipErr)
	}

	trusted, resolveErr := repository.ResolveContext(ctx, user.ID(), tenantA.ID())
	if resolveErr != nil {
		t.Fatalf("ResolveContext(same tenant) error = %v", resolveErr)
	}
	if !trusted.Valid() || trusted.TenantID() != tenantA.ID() || trusted.PrincipalID() != user.ID() {
		t.Fatalf("resolved trusted context mismatch")
	}
	if _, sameScopeErr := trusted.ScopeFor(tenantA.ID()); sameScopeErr != nil {
		t.Fatalf("ScopeFor(same tenant) error = %v", sameScopeErr)
	}
	if _, crossScopeErr := trusted.ScopeFor(tenantB.ID()); !failure.IsCode(crossScopeErr, codeCrossTenantDenied) {
		t.Fatalf("ScopeFor(cross tenant) error = %v, want %s", crossScopeErr, codeCrossTenantDenied)
	}
	if _, forgedTenantErr := repository.ResolveContext(ctx, user.ID(), tenantB.ID()); !failure.IsCode(forgedTenantErr, codeContextUntrusted) {
		t.Fatalf("forged tenant selector error = %v, want %s", forgedTenantErr, codeContextUntrusted)
	}
	if _, missingScopeErr := repository.ResolveContext(ctx, user.ID(), TenantID("")); !failure.IsCode(missingScopeErr, codeContextUntrusted) {
		t.Fatalf("missing tenant selector error = %v, want %s", missingScopeErr, codeContextUntrusted)
	}

	duplicateMembership, duplicateMembershipErr := NewMembership(tenantA.ID(), user.ID())
	if duplicateMembershipErr != nil {
		t.Fatalf("NewMembership(duplicate) error = %v", duplicateMembershipErr)
	}
	if createDuplicateErr := repository.CreateMembership(ctx, duplicateMembership); !failure.IsCode(createDuplicateErr, codeMembershipConflict) {
		t.Fatalf("duplicate active membership error = %v, want %s", createDuplicateErr, codeMembershipConflict)
	}

	if _, wrongTenantRevokeErr := repository.RevokeMembership(
		ctx,
		tenantB.ID(),
		membership.ID(),
		activeAt.Add(time.Minute),
	); !failure.IsCode(wrongTenantRevokeErr, codeMembershipConflict) {
		t.Fatalf("wrong-tenant revoke error = %v, want %s", wrongTenantRevokeErr, codeMembershipConflict)
	}
	if _, stillTrustedErr := repository.ResolveContext(ctx, user.ID(), tenantA.ID()); stillTrustedErr != nil {
		t.Fatalf("wrong-tenant revoke changed membership: %v", stillTrustedErr)
	}

	revoked, revokeErr := repository.RevokeMembership(ctx, tenantA.ID(), membership.ID(), activeAt.Add(2*time.Minute))
	if revokeErr != nil {
		t.Fatalf("RevokeMembership() error = %v", revokeErr)
	}
	if revoked.State() != MembershipStateRevoked {
		t.Fatalf("revoked membership state = %q", revoked.State())
	}
	if _, revokedContextErr := repository.ResolveContext(ctx, user.ID(), tenantA.ID()); !failure.IsCode(revokedContextErr, codeContextUntrusted) {
		t.Fatalf("revoked membership context error = %v, want %s", revokedContextErr, codeContextUntrusted)
	}

	replacementMembership, replacementErr := NewMembership(tenantA.ID(), user.ID())
	if replacementErr != nil {
		t.Fatalf("NewMembership(replacement) error = %v", replacementErr)
	}
	if createReplacementErr := repository.CreateMembership(ctx, replacementMembership); createReplacementErr != nil {
		t.Fatalf("CreateMembership(replacement) error = %v", createReplacementErr)
	}
	if _, replacementContextErr := repository.ResolveContext(ctx, user.ID(), tenantA.ID()); replacementContextErr != nil {
		t.Fatalf("replacement relationship did not restore tenant context: %v", replacementContextErr)
	}

	suspendedAt := activeAt.Add(3 * time.Minute)
	if _, suspendErr := repository.TransitionTenant(ctx, tenantA.ID(), TenantStateActive, TenantStateSuspended, suspendedAt); suspendErr != nil {
		t.Fatalf("TransitionTenant(A suspended) error = %v", suspendErr)
	}
	if _, suspendedContextErr := repository.ResolveContext(ctx, user.ID(), tenantA.ID()); !failure.IsCode(suspendedContextErr, codeContextUntrusted) {
		t.Fatalf("suspended tenant context error = %v, want %s", suspendedContextErr, codeContextUntrusted)
	}
	if _, staleTransitionErr := repository.TransitionTenant(
		ctx,
		tenantA.ID(),
		TenantStateActive,
		TenantStateDisabled,
		suspendedAt.Add(time.Minute),
	); !failure.IsCode(staleTransitionErr, codeTenantConflict) {
		t.Fatalf("stale tenant transition error = %v, want %s", staleTransitionErr, codeTenantConflict)
	}

	var forbiddenColumns int
	if queryErr := pool.QueryRow(
		ctx,
		`SELECT count(*)
		 FROM information_schema.columns
		 WHERE table_schema = 'omnexa_tenancy'
		   AND column_name IN (
		       'organization_id', 'role_id', 'permission_id', 'policy_id',
		       'password', 'password_hash', 'session_id', 'mfa_secret',
		       'passkey', 'api_key', 'setting_key', 'global_tenant_id'
		   )`,
	).Scan(&forbiddenColumns); queryErr != nil {
		t.Fatalf("forbidden-column query error = %v", queryErr)
	}
	if forbiddenColumns != 0 {
		t.Fatalf("tenancy schema contains %d premature organization/auth/authz/settings columns", forbiddenColumns)
	}

	var identityPrincipalCount int
	if queryErr := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_identity.principals WHERE id = $1 AND principal_type = 'human_user'",
		string(user.ID()),
	).Scan(&identityPrincipalCount); queryErr != nil {
		t.Fatalf("identity ownership query error = %v", queryErr)
	}
	if identityPrincipalCount != 1 {
		t.Fatalf("tenancy operations changed identity owner state")
	}

	drifted := []database.Migration{{
		Version: 1,
		Name:    "create_tenancy_foundation",
		SQL:     string(tenancyMigrationSQL) + "\n-- rewritten migration must fail closed\n",
	}}
	driftMigrator, driftMigratorErr := database.NewMigrator(pool, "kernel.tenancy", drifted, 5*time.Second)
	if driftMigratorErr != nil {
		t.Fatalf("NewMigrator(drifted tenancy) error = %v", driftMigratorErr)
	}
	if driftErr := driftMigrator.Run(ctx); driftErr == nil {
		t.Fatalf("rewritten applied tenancy migration unexpectedly succeeded")
	}
}

package organization

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresP02AuditAndExitIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_10_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_10_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		t.Fatalf("pool.Ping() error = %v", err)
	}
	foundation, err := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if err != nil {
		t.Fatalf("database.NewMigrator(kernel.foundation) error = %v", err)
	}
	if err = foundation.Run(ctx); err != nil {
		t.Fatalf("foundation.Run() error = %v", err)
	}
	resetP0203Database(t, ctx, pool)
	defer resetP0203Database(t, context.Background(), pool)

	identityMigrations := []database.Migration{
		p0210Migration(t, 1, "create_identity_foundation", "identity-1"),
		p0210Migration(t, 2, "create_authentication_sessions", "identity-2"),
		p0210Migration(t, 3, "create_strong_authentication", "identity-3"),
		p0210Migration(t, 4, "create_service_accounts_api_credentials", "identity-4"),
	}
	tenancyMigrations := []database.Migration{p0210Migration(t, 1, "create_tenancy_foundation", "tenancy-1")}
	organizationMigrations := []database.Migration{p0210Migration(t, 1, "create_organization_foundation", "organization-1")}
	authorizationMigrations := []database.Migration{
		p0210Migration(t, 1, "create_rbac_foundation", "authorization-1"),
		p0210Migration(t, 2, "allow_service_account_assignments", "authorization-2"),
		p0210Migration(t, 3, "add_configuration_permissions", "authorization-3"),
	}
	configurationMigrations := []database.Migration{p0210Migration(t, 1, "create_scoped_settings", "configuration-1")}

	p0210RunMigrations(t, ctx, pool, "kernel.identity", identityMigrations)
	p0210RunMigrations(t, ctx, pool, "kernel.tenancy", tenancyMigrations)
	p0210RunMigrations(t, ctx, pool, "kernel.organization", organizationMigrations)
	p0210RunMigrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	p0210RunMigrations(t, ctx, pool, "kernel.configuration", configurationMigrations)
	// P02.10 adds no persistence; rerunning the complete accepted P02.09 baseline is
	// the supported no-op upgrade/idempotency proof for the P02 exit package.
	p0210RunMigrations(t, ctx, pool, "kernel.identity", identityMigrations)
	p0210RunMigrations(t, ctx, pool, "kernel.tenancy", tenancyMigrations)
	p0210RunMigrations(t, ctx, pool, "kernel.organization", organizationMigrations)
	p0210RunMigrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	p0210RunMigrations(t, ctx, pool, "kernel.configuration", configurationMigrations)
	p0210AssertLedger(t, ctx, pool, "kernel.identity", 4)
	p0210AssertLedger(t, ctx, pool, "kernel.tenancy", 1)
	p0210AssertLedger(t, ctx, pool, "kernel.organization", 1)
	p0210AssertLedger(t, ctx, pool, "kernel.authorization", 3)
	p0210AssertLedger(t, ctx, pool, "kernel.configuration", 1)

	identityRepository, err := identity.NewPostgresRepository(pool)
	if err != nil {
		t.Fatalf("identity.NewPostgresRepository() error = %v", err)
	}
	user, err := identity.NewUser("p02-exit-owner@example.com")
	if err != nil {
		t.Fatalf("identity.NewUser() error = %v", err)
	}
	if err = identityRepository.Create(ctx, user); err != nil {
		t.Fatalf("identityRepository.Create() error = %v", err)
	}

	sink, err := audit.NewMemorySink(64)
	if err != nil {
		t.Fatalf("audit.NewMemorySink() error = %v", err)
	}
	writer, err := audit.NewWriter(sink, nil)
	if err != nil {
		t.Fatalf("audit.NewWriter() error = %v", err)
	}
	identityAudit, err := identity.NewAuditAdapter(writer)
	if err != nil {
		t.Fatalf("identity.NewAuditAdapter() error = %v", err)
	}
	identityAudit.RecordSecurityEvent(identity.SecurityAuditEvent{
		Action:      identity.SecurityAuditAuthenticationOK,
		PrincipalID: user.ID(),
		Succeeded:   true,
		OccurredAt:  time.Date(2026, time.August, 26, 2, 15, 0, 0, time.UTC),
	})

	tenancyBase, err := tenancy.NewPostgresRepository(pool)
	if err != nil {
		t.Fatalf("tenancy.NewPostgresRepository() error = %v", err)
	}
	tenancyRepository, err := tenancy.NewAuditedRepository(tenancyBase, writer)
	if err != nil {
		t.Fatalf("tenancy.NewAuditedRepository() error = %v", err)
	}
	tenantA, err := tenancy.NewTenant()
	if err != nil {
		t.Fatalf("tenancy.NewTenant(A) error = %v", err)
	}
	if err = tenancyRepository.CreateTenant(ctx, tenantA); err != nil {
		t.Fatalf("CreateTenant(A) error = %v", err)
	}
	activeAt := time.Now().UTC().Add(time.Second)
	if _, err = tenancyRepository.TransitionTenant(ctx, tenantA.ID(), tenancy.TenantStateProvisioned, tenancy.TenantStateActive, activeAt); err != nil {
		t.Fatalf("TransitionTenant(A) error = %v", err)
	}
	membershipA, err := tenancy.NewMembership(tenantA.ID(), user.ID())
	if err != nil {
		t.Fatalf("tenancy.NewMembership(A) error = %v", err)
	}
	if err = tenancyRepository.CreateMembership(ctx, membershipA); err != nil {
		t.Fatalf("CreateMembership(A) error = %v", err)
	}
	trustedA, err := tenancyRepository.ResolveContext(ctx, user.ID(), tenantA.ID())
	if err != nil {
		t.Fatalf("ResolveContext(A) error = %v", err)
	}
	scopeA, err := trustedA.ScopeFor(tenantA.ID())
	if err != nil {
		t.Fatalf("ScopeFor(A) error = %v", err)
	}

	organizationBase, err := NewPostgresRepository(pool, tenancyRepository)
	if err != nil {
		t.Fatalf("organization.NewPostgresRepository() error = %v", err)
	}
	organizationRepository, err := NewAuditedRepository(organizationBase, writer)
	if err != nil {
		t.Fatalf("organization.NewAuditedRepository() error = %v", err)
	}
	rootA, err := NewOrganization(tenantA.ID())
	if err != nil {
		t.Fatalf("NewOrganization(A) error = %v", err)
	}
	if err = organizationRepository.CreateNode(ctx, scopeA, rootA); err != nil {
		t.Fatalf("CreateNode(A) error = %v", err)
	}
	organizationMembership, err := NewMembership(tenantA.ID(), rootA.ID(), user.ID())
	if err != nil {
		t.Fatalf("organization.NewMembership() error = %v", err)
	}
	if err = organizationRepository.CreateMembership(ctx, scopeA, organizationMembership); err != nil {
		t.Fatalf("organization CreateMembership() error = %v", err)
	}

	tenantB, err := tenancy.NewTenant()
	if err != nil {
		t.Fatalf("tenancy.NewTenant(B) error = %v", err)
	}
	if err = tenancyRepository.CreateTenant(ctx, tenantB); err != nil {
		t.Fatalf("CreateTenant(B) error = %v", err)
	}
	crossTenantRoot, err := NewOrganization(tenantB.ID())
	if err != nil {
		t.Fatalf("NewOrganization(B) error = %v", err)
	}
	if err = organizationRepository.CreateNode(ctx, scopeA, crossTenantRoot); err == nil {
		t.Fatal("cross-tenant organization create unexpectedly succeeded")
	}

	failingWriter, err := audit.NewWriter(p0210FailingSink{}, nil)
	if err != nil {
		t.Fatalf("audit.NewWriter(failing) error = %v", err)
	}
	failingOrganization, err := NewAuditedRepository(organizationBase, failingWriter)
	if err != nil {
		t.Fatalf("NewAuditedRepository(failing) error = %v", err)
	}
	branchA, err := NewChild(tenantA.ID(), NodeKindBranch, rootA.ID())
	if err != nil {
		t.Fatalf("NewChild(branchA) error = %v", err)
	}
	if err = failingOrganization.CreateNode(ctx, scopeA, branchA); err == nil {
		t.Fatal("required audit failure unexpectedly claimed organization mutation success")
	}
	if _, getErr := organizationBase.GetNode(ctx, scopeA, branchA.ID()); getErr != nil {
		t.Fatalf("underlying mutation was not persisted before required audit failure: %v", getErr)
	}
	if health := failingWriter.Health(); health.Submitted != 1 || health.Failed != 1 {
		t.Fatalf("required-audit writer health = %#v", health)
	}

	actions := make(map[string]bool)
	for _, record := range sink.Snapshot() {
		actions[record.Action()] = true
		if record.Classification() != audit.ClassificationConfidential || len(record.Fields()) != 0 {
			t.Fatalf("classification-safe audit invariant failed for %q: %q/%#v", record.Action(), record.Classification(), record.Fields())
		}
		if err = record.Verify(); err != nil {
			t.Fatalf("record.Verify(%q) error = %v", record.Action(), err)
		}
	}
	for _, action := range []string{
		"identity.authentication.succeeded",
		"tenancy.tenant.create",
		"tenancy.tenant.transition",
		"tenancy.membership.create",
		"organization.node.create",
		"organization.membership.create",
	} {
		if !actions[action] {
			t.Fatalf("missing P02.10 audit action %q; got %#v", action, actions)
		}
	}
}

type p0210FailingSink struct{}

func (p0210FailingSink) Append(context.Context, audit.Record) error {
	return errors.New("synthetic required audit failure")
}

func p0210RunMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool, owner string, migrations []database.Migration) {
	t.Helper()
	migrator, err := database.NewMigrator(pool, owner, migrations, 5*time.Second)
	if err != nil {
		t.Fatalf("database.NewMigrator(%s) error = %v", owner, err)
	}
	if err = migrator.Run(ctx); err != nil {
		t.Fatalf("migrator.Run(%s) error = %v", owner, err)
	}
}

func p0210AssertLedger(t *testing.T, ctx context.Context, pool *pgxpool.Pool, owner string, want int) {
	t.Helper()
	var got int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = $1", owner).Scan(&got); err != nil {
		t.Fatalf("ledger query(%s) error = %v", owner, err)
	}
	if got != want {
		t.Fatalf("ledger count(%s) = %d, want %d", owner, got, want)
	}
}

func p0210Migration(t *testing.T, version int64, name, fixture string) database.Migration {
	t.Helper()
	var (
		contents []byte
		err      error
	)
	switch fixture {
	case "identity-1":
		contents, err = os.ReadFile("../../migrations/kernel.identity/1_create_identity_foundation.sql")
	case "identity-2":
		contents, err = os.ReadFile("../../migrations/kernel.identity/2_create_authentication_sessions.sql")
	case "identity-3":
		contents, err = os.ReadFile("../../migrations/kernel.identity/3_create_strong_authentication.sql")
	case "identity-4":
		contents, err = os.ReadFile("../../migrations/kernel.identity/4_create_service_accounts_api_credentials.sql")
	case "tenancy-1":
		contents, err = os.ReadFile("../../migrations/kernel.tenancy/1_create_tenancy_foundation.sql")
	case "organization-1":
		contents, err = os.ReadFile("../../migrations/kernel.organization/1_create_organization_foundation.sql")
	case "authorization-1":
		contents, err = os.ReadFile("../../migrations/kernel.authorization/1_create_rbac_foundation.sql")
	case "authorization-2":
		contents, err = os.ReadFile("../../migrations/kernel.authorization/2_allow_service_account_assignments.sql")
	case "authorization-3":
		contents, err = os.ReadFile("../../migrations/kernel.authorization/3_add_configuration_permissions.sql")
	case "configuration-1":
		contents, err = os.ReadFile("../../migrations/kernel.configuration/1_create_scoped_settings.sql")
	default:
		t.Fatalf("unsupported P02.10 migration fixture %q", fixture)
	}
	if err != nil {
		t.Fatalf("read P02.10 migration fixture %q error = %v", fixture, err)
	}
	return database.Migration{Version: version, Name: name, SQL: string(contents)}
}

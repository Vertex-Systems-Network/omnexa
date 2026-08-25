package configuration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/authorization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresTenantScopedSettingsIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_09_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_09_TEST_DATABASE_URL is not set")
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
		t.Fatalf("foundation migrator create error = %v", err)
	}
	if err = foundation.Run(ctx); err != nil {
		t.Fatalf("foundation migrator run error = %v", err)
	}
	defer resetP0209Database(t, context.Background(), pool)

	identityMigrations := []database.Migration{
		readP0209Migration(t, 1, "create_identity_foundation", "../../migrations/kernel.identity/1_create_identity_foundation.sql"),
		readP0209Migration(t, 2, "create_authentication_sessions", "../../migrations/kernel.identity/2_create_authentication_sessions.sql"),
		readP0209Migration(t, 3, "create_strong_authentication", "../../migrations/kernel.identity/3_create_strong_authentication.sql"),
		readP0209Migration(t, 4, "create_service_accounts_api_credentials", "../../migrations/kernel.identity/4_create_service_accounts_api_credentials.sql"),
	}
	tenancyMigrations := []database.Migration{
		readP0209Migration(t, 1, "create_tenancy_foundation", "../../migrations/kernel.tenancy/1_create_tenancy_foundation.sql"),
	}
	organizationMigrations := []database.Migration{
		readP0209Migration(t, 1, "create_organization_foundation", "../../migrations/kernel.organization/1_create_organization_foundation.sql"),
	}
	authorizationMigrations := []database.Migration{
		readP0209Migration(t, 1, "create_rbac_foundation", "../../migrations/kernel.authorization/1_create_rbac_foundation.sql"),
		readP0209Migration(t, 2, "allow_service_account_assignments", "../../migrations/kernel.authorization/2_allow_service_account_assignments.sql"),
		readP0209Migration(t, 3, "add_configuration_permissions", "../../migrations/kernel.authorization/3_add_configuration_permissions.sql"),
	}
	configurationMigrations := []database.Migration{
		readP0209Migration(t, 1, "create_scoped_settings", "../../migrations/kernel.configuration/1_create_scoped_settings.sql"),
	}

	resetP0209Database(t, ctx, pool)
	runP0209Migrations(t, ctx, pool, "kernel.identity", identityMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.tenancy", tenancyMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.organization", organizationMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.configuration", configurationMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.configuration", configurationMigrations)
	assertP0209LedgerCount(t, ctx, pool, "kernel.authorization", 3)
	assertP0209LedgerCount(t, ctx, pool, "kernel.configuration", 1)

	resetP0209Database(t, ctx, pool)
	runP0209Migrations(t, ctx, pool, "kernel.identity", identityMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.tenancy", tenancyMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.organization", organizationMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.authorization", authorizationMigrations[:2])

	fixed := time.Date(2026, time.August, 25, 21, 0, 0, 0, time.UTC)
	ids := newP0209FixtureIDs()
	seedP0209BaseState(t, ctx, pool, fixed, ids)

	tenancyRepository, err := tenancy.NewPostgresRepository(pool)
	if err != nil {
		t.Fatalf("tenancy.NewPostgresRepository() error = %v", err)
	}
	trustedA, err := tenancyRepository.ResolveContext(ctx, identity.UserID(ids.actor), tenancy.TenantID(ids.tenantA))
	if err != nil {
		t.Fatalf("ResolveContext(tenant A) error = %v", err)
	}
	actorA, err := authorization.SubjectFromTenantContext(trustedA)
	if err != nil {
		t.Fatalf("SubjectFromTenantContext(A) error = %v", err)
	}
	sink, err := audit.NewMemorySink(128)
	if err != nil {
		t.Fatalf("audit.NewMemorySink() error = %v", err)
	}
	writer, err := audit.NewWriter(sink, nil)
	if err != nil {
		t.Fatalf("audit.NewWriter() error = %v", err)
	}
	rbacRepository, err := authorization.NewPostgresRepository(pool)
	if err != nil {
		t.Fatalf("authorization.NewPostgresRepository() error = %v", err)
	}
	rbac, err := authorization.NewService(rbacRepository, writer)
	if err != nil {
		t.Fatalf("authorization.NewService() error = %v", err)
	}
	if decision, checkErr := rbac.Check(ctx, actorA, authorization.PermissionRoleRead); checkErr != nil || decision != authorization.DecisionAllow {
		t.Fatalf("accepted P02.08 baseline human RBAC = %q/%v", decision, checkErr)
	}

	runP0209Migrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	runP0209Migrations(t, ctx, pool, "kernel.configuration", configurationMigrations)
	if decision, checkErr := rbac.Check(ctx, actorA, authorization.PermissionRoleRead); checkErr != nil || decision != authorization.DecisionAllow {
		t.Fatalf("P02.09 upgrade changed accepted human RBAC = %q/%v", decision, checkErr)
	}
	seedP0209ConfigurationAuthority(t, ctx, pool, fixed, ids)

	var permissionCount int
	if err = pool.QueryRow(ctx, `SELECT count(*) FROM omnexa_authorization.permissions WHERE permission_id IN ('configuration.setting.read','configuration.setting.manage')`).Scan(&permissionCount); err != nil || permissionCount != 2 {
		t.Fatalf("configuration permission seed count/error = %d/%v", permissionCount, err)
	}
	var forbiddenColumns int
	if err = pool.QueryRow(ctx, `SELECT count(*) FROM information_schema.columns WHERE table_schema='omnexa_configuration' AND table_name='setting_overrides' AND column_name ~ '(secret|password|credential|token|key_material)'`).Scan(&forbiddenColumns); err != nil || forbiddenColumns != 0 {
		t.Fatalf("forbidden secret columns count/error = %d/%v", forbiddenColumns, err)
	}

	organizationRepository, err := organization.NewPostgresRepository(pool, tenancyRepository)
	if err != nil {
		t.Fatalf("organization.NewPostgresRepository() error = %v", err)
	}
	trustedB, err := tenancyRepository.ResolveContext(ctx, identity.UserID(ids.actor), tenancy.TenantID(ids.tenantB))
	if err != nil {
		t.Fatalf("ResolveContext(tenant B) error = %v", err)
	}
	orgContextA, err := organizationRepository.ResolveContext(ctx, trustedA, organization.NodeID(ids.orgA))
	if err != nil {
		t.Fatalf("organization ResolveContext(A) error = %v", err)
	}
	orgContextB, err := organizationRepository.ResolveContext(ctx, trustedA, organization.NodeID(ids.orgB))
	if err != nil {
		t.Fatalf("organization ResolveContext(B) error = %v", err)
	}
	tenantScopeA, err := ScopeFromTenantContext(trustedA)
	if err != nil {
		t.Fatalf("ScopeFromTenantContext(A) error = %v", err)
	}
	tenantScopeB, err := ScopeFromTenantContext(trustedB)
	if err != nil {
		t.Fatalf("ScopeFromTenantContext(B) error = %v", err)
	}
	orgScopeA, err := ScopeFromOrganizationContext(orgContextA)
	if err != nil {
		t.Fatalf("ScopeFromOrganizationContext(A) error = %v", err)
	}
	orgScopeB, err := ScopeFromOrganizationContext(orgContextB)
	if err != nil {
		t.Fatalf("ScopeFromOrganizationContext(B) error = %v", err)
	}

	registry := mustRegistry(t, p0209PublicDefinition(), p0209ProtectedDefinition())
	repository, err := NewPostgresScopedRepository(pool)
	if err != nil {
		t.Fatalf("NewPostgresScopedRepository() error = %v", err)
	}
	service, err := NewScopedService(registry, repository, rbac, writer, []SettingPolicy{
		{Key: p0209PublicDefinition().Key, Classification: DataPublic, AllowOrganizationOverride: true},
		{Key: p0209ProtectedDefinition().Key, Classification: DataConfidential, AllowOrganizationOverride: true, ProtectedRead: true, SecuritySignificant: true},
	}, EvaluatorOptions{CacheTTL: time.Minute})
	if err != nil {
		t.Fatalf("NewScopedService() error = %v", err)
	}
	metadata := SettingMutationMetadata{CorrelationID: "p02-09-integration", Reason: "exercise trusted scoped setting mutation"}

	if _, err = service.Upsert(ctx, tenantScopeA, p0209PublicDefinition().Key, StringValue("tenant-a"), metadata); err != nil {
		t.Fatalf("tenant public Upsert() error = %v", err)
	}
	orgBFallback, err := service.Resolve(ctx, orgScopeB, p0209PublicDefinition().Key)
	if err != nil {
		t.Fatalf("org B tenant fallback Resolve() error = %v", err)
	}
	assertP0209StringEvaluation(t, orgBFallback, "tenant-a")

	if _, err = service.Upsert(ctx, orgScopeA, p0209PublicDefinition().Key, StringValue("org-a"), metadata); err != nil {
		t.Fatalf("org A public Upsert() error = %v", err)
	}
	orgAValue, err := service.Resolve(ctx, orgScopeA, p0209PublicDefinition().Key)
	if err != nil {
		t.Fatalf("org A Resolve() error = %v", err)
	}
	assertP0209StringEvaluation(t, orgAValue, "org-a")
	tenantAValue, err := service.Resolve(ctx, tenantScopeA, p0209PublicDefinition().Key)
	if err != nil {
		t.Fatalf("tenant A Resolve() error = %v", err)
	}
	assertP0209StringEvaluation(t, tenantAValue, "tenant-a")

	tenantBPublic, err := service.Resolve(ctx, tenantScopeB, p0209PublicDefinition().Key)
	if err != nil {
		t.Fatalf("tenant B public Resolve() error = %v", err)
	}
	assertP0209StringEvaluation(t, tenantBPublic, "en")
	if _, err = service.Upsert(ctx, tenantScopeB, p0209PublicDefinition().Key, StringValue("tenant-b"), metadata); err == nil {
		t.Fatal("wrong-tenant public mutation without scoped manage authority unexpectedly succeeded")
	}
	if _, err = service.Upsert(ctx, orgScopeB, p0209PublicDefinition().Key, StringValue("org-b"), metadata); err == nil {
		t.Fatal("wrong-organization public mutation without scoped manage authority unexpectedly succeeded")
	}

	protectedMarker := "confidential-p02-09-marker"
	if _, err = service.Upsert(ctx, tenantScopeA, p0209ProtectedDefinition().Key, StringValue(protectedMarker), metadata); err != nil {
		t.Fatalf("protected tenant Upsert() error = %v", err)
	}
	protectedA, err := service.Resolve(ctx, tenantScopeA, p0209ProtectedDefinition().Key)
	if err != nil {
		t.Fatalf("protected tenant A Resolve() error = %v", err)
	}
	assertP0209StringEvaluation(t, protectedA, protectedMarker)
	if protectedA.Classification != DataConfidential {
		t.Fatalf("protected classification = %q", protectedA.Classification)
	}
	if _, err = service.Resolve(ctx, tenantScopeB, p0209ProtectedDefinition().Key); err == nil {
		t.Fatal("wrong-tenant protected read unexpectedly succeeded")
	}
	if _, err = service.Resolve(ctx, orgScopeB, p0209ProtectedDefinition().Key); err == nil {
		t.Fatal("wrong-organization protected read unexpectedly succeeded")
	}

	records := sink.Snapshot()
	if len(records) != 3 {
		t.Fatalf("configuration audit records = %d, want 3", len(records))
	}
	for _, record := range records {
		if record.Action() != "configuration.setting.upsert" || record.Target().Kind != "configuration.setting" || len(record.Fields()) != 0 || !record.Privileged() {
			t.Fatalf("unsafe configuration audit record: action=%s target=%+v fields=%v", record.Action(), record.Target(), record.Fields())
		}
		if record.Target().Reference == protectedMarker {
			t.Fatal("setting value leaked into audit target")
		}
	}

	var tenantAStored string
	if err = pool.QueryRow(ctx, `SELECT value_text FROM omnexa_configuration.setting_overrides WHERE tenant_id=$1 AND organization_id IS NULL AND setting_key=$2`, ids.tenantA, string(p0209PublicDefinition().Key)).Scan(&tenantAStored); err != nil || tenantAStored != "tenant-a" {
		t.Fatalf("tenant A stored value/error = %q/%v", tenantAStored, err)
	}
	var tenantBCount int
	if err = pool.QueryRow(ctx, `SELECT count(*) FROM omnexa_configuration.setting_overrides WHERE tenant_id=$1`, ids.tenantB).Scan(&tenantBCount); err != nil || tenantBCount != 0 {
		t.Fatalf("wrong-tenant mutation persisted rows = %d/%v", tenantBCount, err)
	}
}

type p0209FixtureIDs struct {
	actor             string
	tenantA           string
	tenantB           string
	membershipA       string
	membershipB       string
	orgA              string
	orgB              string
	orgMembershipA    string
	orgMembershipB    string
	tenantRoleA       string
	tenantAssignmentA string
	orgRoleA          string
	orgAssignmentA    string
}

func newP0209FixtureIDs() p0209FixtureIDs {
	return p0209FixtureIDs{
		actor:             "01890f3e-7b9a-7cc0-98c4-dc0c0c074901",
		tenantA:           "01890f3e-7b9a-7cc0-98c4-dc0c0c074911",
		tenantB:           "01890f3e-7b9a-7cc0-98c4-dc0c0c074912",
		membershipA:       "01890f3e-7b9a-7cc0-98c4-dc0c0c074921",
		membershipB:       "01890f3e-7b9a-7cc0-98c4-dc0c0c074922",
		orgA:              "01890f3e-7b9a-7cc0-98c4-dc0c0c074931",
		orgB:              "01890f3e-7b9a-7cc0-98c4-dc0c0c074932",
		orgMembershipA:    "01890f3e-7b9a-7cc0-98c4-dc0c0c074941",
		orgMembershipB:    "01890f3e-7b9a-7cc0-98c4-dc0c0c074942",
		tenantRoleA:       "01890f3e-7b9a-7cc0-98c4-dc0c0c074951",
		tenantAssignmentA: "01890f3e-7b9a-7cc0-98c4-dc0c0c074952",
		orgRoleA:          "01890f3e-7b9a-7cc0-98c4-dc0c0c074961",
		orgAssignmentA:    "01890f3e-7b9a-7cc0-98c4-dc0c0c074962",
	}
}

func p0209PublicDefinition() Definition {
	return Definition{
		Key: "platform.ui.locale", Description: "Synthetic tenant/org locale setting.", Owner: "kernel.configuration",
		Kind: KindString, Class: ClassRuntimeConfig, Version: 1,
		Default: StringValue("en"), Fallback: StringValue("en"),
	}
}

func p0209ProtectedDefinition() Definition {
	return Definition{
		Key: "platform.security.login_notice", Description: "Synthetic protected security setting.", Owner: "kernel.configuration",
		Kind: KindString, Class: ClassRuntimeConfig, Version: 1,
		Default: StringValue("baseline"), Fallback: StringValue("baseline"),
	}
}

func assertP0209StringEvaluation(t *testing.T, evaluation ScopedEvaluation, want string) {
	t.Helper()
	got, ok := evaluation.Evaluation.Value.String()
	if !ok || got != want {
		t.Fatalf("scoped evaluation = %q/%v, want %q/true", got, ok, want)
	}
}

func readP0209Migration(t *testing.T, version int64, name, path string) database.Migration {
	t.Helper()
	var (
		contents []byte
		err      error
	)
	switch path {
	case "../../migrations/kernel.identity/1_create_identity_foundation.sql":
		contents, err = os.ReadFile("../../migrations/kernel.identity/1_create_identity_foundation.sql")
	case "../../migrations/kernel.identity/2_create_authentication_sessions.sql":
		contents, err = os.ReadFile("../../migrations/kernel.identity/2_create_authentication_sessions.sql")
	case "../../migrations/kernel.identity/3_create_strong_authentication.sql":
		contents, err = os.ReadFile("../../migrations/kernel.identity/3_create_strong_authentication.sql")
	case "../../migrations/kernel.identity/4_create_service_accounts_api_credentials.sql":
		contents, err = os.ReadFile("../../migrations/kernel.identity/4_create_service_accounts_api_credentials.sql")
	case "../../migrations/kernel.tenancy/1_create_tenancy_foundation.sql":
		contents, err = os.ReadFile("../../migrations/kernel.tenancy/1_create_tenancy_foundation.sql")
	case "../../migrations/kernel.organization/1_create_organization_foundation.sql":
		contents, err = os.ReadFile("../../migrations/kernel.organization/1_create_organization_foundation.sql")
	case "../../migrations/kernel.authorization/1_create_rbac_foundation.sql":
		contents, err = os.ReadFile("../../migrations/kernel.authorization/1_create_rbac_foundation.sql")
	case "../../migrations/kernel.authorization/2_allow_service_account_assignments.sql":
		contents, err = os.ReadFile("../../migrations/kernel.authorization/2_allow_service_account_assignments.sql")
	case "../../migrations/kernel.authorization/3_add_configuration_permissions.sql":
		contents, err = os.ReadFile("../../migrations/kernel.authorization/3_add_configuration_permissions.sql")
	case "../../migrations/kernel.configuration/1_create_scoped_settings.sql":
		contents, err = os.ReadFile("../../migrations/kernel.configuration/1_create_scoped_settings.sql")
	default:
		t.Fatalf("unsupported migration fixture %q", path)
	}
	if err != nil {
		t.Fatalf("read migration %s error = %v", path, err)
	}
	return database.Migration{Version: version, Name: name, SQL: string(contents)}
}

func runP0209Migrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool, owner string, migrations []database.Migration) {
	t.Helper()
	migrator, err := database.NewMigrator(pool, owner, migrations, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(%s) error = %v", owner, err)
	}
	if err = migrator.Run(ctx); err != nil {
		t.Fatalf("migrator Run(%s) error = %v", owner, err)
	}
}

func assertP0209LedgerCount(t *testing.T, ctx context.Context, pool *pgxpool.Pool, owner string, want int) {
	t.Helper()
	var got int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner=$1", owner).Scan(&got); err != nil || got != want {
		t.Fatalf("migration ledger %s count/error = %d/%v, want %d", owner, got, err, want)
	}
}

func seedP0209BaseState(t *testing.T, ctx context.Context, pool *pgxpool.Pool, at time.Time, ids p0209FixtureIDs) {
	t.Helper()
	actor := ids.actor
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO omnexa_identity.principals (id, principal_type, lifecycle_state, created_at, updated_at) VALUES ($1,'human_user','active',$2,$2)`, []any{actor, at}},
		{`INSERT INTO omnexa_identity.users (principal_id, primary_email) VALUES ($1,'p0209-admin@example.com')`, []any{actor}},
		{`INSERT INTO omnexa_tenancy.tenants (id, lifecycle_state, created_at, updated_at) VALUES ($1,'active',$2,$2),($3,'active',$2,$2)`, []any{ids.tenantA, at, ids.tenantB}},
		{`INSERT INTO omnexa_tenancy.tenant_memberships (id, tenant_id, principal_id, relationship_state, created_at, updated_at) VALUES ($1,$2,$3,'active',$4,$4),($5,$6,$3,'active',$4,$4)`, []any{ids.membershipA, ids.tenantA, actor, at, ids.membershipB, ids.tenantB}},
		{`INSERT INTO omnexa_organization.nodes (id, tenant_id, node_kind, parent_id, created_at, updated_at) VALUES ($1,$2,'organization',NULL,$4,$4),($3,$2,'organization',NULL,$4,$4)`, []any{ids.orgA, ids.tenantA, ids.orgB, at}},
		{`INSERT INTO omnexa_organization.scoped_memberships (id, tenant_id, scope_id, principal_id, relationship_state, created_at, updated_at) VALUES ($1,$2,$3,$4,'active',$5,$5),($6,$2,$7,$4,'active',$5,$5)`, []any{ids.orgMembershipA, ids.tenantA, ids.orgA, actor, at, ids.orgMembershipB, ids.orgB}},
		{`INSERT INTO omnexa_authorization.roles (id, tenant_id, organization_id, name, created_at, updated_at) VALUES ($1,$2,NULL,'p0209 baseline',$3,$3)`, []any{ids.tenantRoleA, ids.tenantA, at}},
		{`INSERT INTO omnexa_authorization.role_permissions (role_id, permission_id) VALUES ($1,'authorization.role.read')`, []any{ids.tenantRoleA}},
		{`INSERT INTO omnexa_authorization.role_assignments (id, role_id, principal_id, assignment_state, created_at, updated_at) VALUES ($1,$2,$3,'active',$4,$4)`, []any{ids.tenantAssignmentA, ids.tenantRoleA, actor, at}},
	}
	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement.query, statement.args...); err != nil {
			t.Fatalf("seed P02.09 base state error = %v", err)
		}
	}
}

func seedP0209ConfigurationAuthority(t *testing.T, ctx context.Context, pool *pgxpool.Pool, at time.Time, ids p0209FixtureIDs) {
	t.Helper()
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO omnexa_authorization.role_permissions (role_id, permission_id) VALUES ($1,'configuration.setting.read'),($1,'configuration.setting.manage')`, []any{ids.tenantRoleA}},
		{`INSERT INTO omnexa_authorization.roles (id, tenant_id, organization_id, name, created_at, updated_at) VALUES ($1,$2,$3,'p0209 org config admin',$4,$4)`, []any{ids.orgRoleA, ids.tenantA, ids.orgA, at}},
		{`INSERT INTO omnexa_authorization.role_permissions (role_id, permission_id) VALUES ($1,'configuration.setting.read'),($1,'configuration.setting.manage')`, []any{ids.orgRoleA}},
		{`INSERT INTO omnexa_authorization.role_assignments (id, role_id, principal_id, assignment_state, created_at, updated_at) VALUES ($1,$2,$3,'active',$4,$4)`, []any{ids.orgAssignmentA, ids.orgRoleA, ids.actor, at}},
	}
	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement.query, statement.args...); err != nil {
			t.Fatalf("seed P02.09 configuration authority error = %v", err)
		}
	}
}

func resetP0209Database(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	for _, schema := range []string{"omnexa_configuration", "omnexa_authorization", "omnexa_organization", "omnexa_tenancy", "omnexa_identity"} {
		if _, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS "+schema+" CASCADE"); err != nil {
			t.Fatalf("drop schema %s error = %v", schema, err)
		}
	}
	if _, err := pool.Exec(ctx, `DELETE FROM omnexa_kernel.schema_migrations WHERE owner IN ('kernel.configuration','kernel.authorization','kernel.organization','kernel.tenancy','kernel.identity')`); err != nil {
		t.Fatalf("reset migration ledger error = %v", err)
	}
}

package authorization

import (
	"context"
	"crypto/sha256"
	"os"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/identity"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresServiceAccountCredentialIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_08_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_08_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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
	if err != nil || foundation.Run(ctx) != nil {
		t.Fatalf("foundation migrator error = %v", err)
	}
	defer resetP0208Database(t, context.Background(), pool)

	identityMigrations := []database.Migration{
		readP0208Migration(t, 1, "create_identity_foundation", "../../migrations/kernel.identity/1_create_identity_foundation.sql"),
		readP0208Migration(t, 2, "create_authentication_sessions", "../../migrations/kernel.identity/2_create_authentication_sessions.sql"),
		readP0208Migration(t, 3, "create_strong_authentication", "../../migrations/kernel.identity/3_create_strong_authentication.sql"),
		readP0208Migration(t, 4, "create_service_accounts_api_credentials", "../../migrations/kernel.identity/4_create_service_accounts_api_credentials.sql"),
	}
	tenancyMigrations := []database.Migration{
		readP0208Migration(t, 1, "create_tenancy_foundation", "../../migrations/kernel.tenancy/1_create_tenancy_foundation.sql"),
	}
	organizationMigrations := []database.Migration{
		readP0208Migration(t, 1, "create_organization_foundation", "../../migrations/kernel.organization/1_create_organization_foundation.sql"),
	}
	authorizationMigrations := []database.Migration{
		readP0208Migration(t, 1, "create_rbac_foundation", "../../migrations/kernel.authorization/1_create_rbac_foundation.sql"),
		readP0208Migration(t, 2, "allow_service_account_assignments", "../../migrations/kernel.authorization/2_allow_service_account_assignments.sql"),
	}

	// Fresh/idempotent evidence for both changed schema owners.
	resetP0208Database(t, ctx, pool)
	runP0208Migrations(t, ctx, pool, "kernel.identity", identityMigrations)
	runP0208Migrations(t, ctx, pool, "kernel.tenancy", tenancyMigrations)
	runP0208Migrations(t, ctx, pool, "kernel.organization", organizationMigrations)
	runP0208Migrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	runP0208Migrations(t, ctx, pool, "kernel.identity", identityMigrations)
	runP0208Migrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	assertP0208LedgerCount(t, ctx, pool, "kernel.identity", 4)
	assertP0208LedgerCount(t, ctx, pool, "kernel.authorization", 2)

	// Supported upgrade evidence starts from the exact accepted P02.07 identity
	// and P02.05 authorization baselines, preserves an existing human grant, then
	// adds only P02.08 migrations.
	resetP0208Database(t, ctx, pool)
	runP0208Migrations(t, ctx, pool, "kernel.identity", identityMigrations[:3])
	runP0208Migrations(t, ctx, pool, "kernel.tenancy", tenancyMigrations)
	runP0208Migrations(t, ctx, pool, "kernel.organization", organizationMigrations)
	runP0208Migrations(t, ctx, pool, "kernel.authorization", authorizationMigrations[:1])

	fixed := time.Date(2026, time.August, 25, 13, 0, 0, 0, time.UTC)
	tenantID := "01890f3e-7b9a-7cc0-98c4-dc0c0c073991"
	actorID := "01890f3e-7b9a-7cc0-98c4-dc0c0c073981"
	membershipID := "01890f3e-7b9a-7cc0-98c4-dc0c0c073982"
	adminRoleID := "01890f3e-7b9a-7cc0-98c4-dc0c0c073983"
	adminAssignmentID := "01890f3e-7b9a-7cc0-98c4-dc0c0c073984"
	seedP0208HumanAuthority(t, ctx, pool, fixed, tenantID, actorID, membershipID, adminRoleID, adminAssignmentID)

	tenancyRepository, _ := tenancy.NewPostgresRepository(pool)
	trustedActor, err := tenancyRepository.ResolveContext(ctx, identity.UserID(actorID), tenancy.TenantID(tenantID))
	if err != nil {
		t.Fatalf("ResolveContext(actor) error = %v", err)
	}
	actor, err := SubjectFromTenantContext(trustedActor)
	if err != nil {
		t.Fatalf("SubjectFromTenantContext() error = %v", err)
	}
	sink, _ := audit.NewMemorySink(128)
	writer, _ := audit.NewWriter(sink, nil)
	rbacRepository, _ := NewPostgresRepository(pool)
	rbac, err := NewService(rbacRepository, writer)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	if decision, checkErr := rbac.Check(ctx, actor, PermissionRoleRead); checkErr != nil || decision != DecisionAllow {
		t.Fatalf("baseline human RBAC = %q/%v", decision, checkErr)
	}

	runP0208Migrations(t, ctx, pool, "kernel.identity", identityMigrations)
	runP0208Migrations(t, ctx, pool, "kernel.authorization", authorizationMigrations)
	if decision, checkErr := rbac.Check(ctx, actor, PermissionRoleRead); checkErr != nil || decision != DecisionAllow {
		t.Fatalf("P02.08 upgrade changed accepted human RBAC: %q/%v", decision, checkErr)
	}

	serviceRepository, err := identity.NewPostgresServiceAccountRepository(pool)
	if err != nil {
		t.Fatalf("NewPostgresServiceAccountRepository() error = %v", err)
	}
	serviceIdentity, err := identity.NewServiceAccountService(serviceRepository, nil)
	if err != nil {
		t.Fatalf("NewServiceAccountService() error = %v", err)
	}
	binding, _ := identity.NewServiceAccountBinding(tenantID, "")
	account, err := serviceIdentity.Create(ctx, "p02-08-ci-agent", binding, fixed.Add(time.Minute))
	if err != nil {
		t.Fatalf("service account Create() error = %v", err)
	}
	account, err = serviceIdentity.Transition(ctx, account.ID(), identity.LifecycleProvisioned, identity.LifecycleActive, fixed.Add(2*time.Minute))
	if err != nil {
		t.Fatalf("service account activate error = %v", err)
	}
	issuedAt := fixed.Add(3 * time.Minute)
	issued, err := serviceIdentity.IssueCredential(ctx, account.ID(), issuedAt.Add(time.Hour), issuedAt)
	if err != nil {
		t.Fatalf("IssueCredential() error = %v", err)
	}
	raw := issued.Secret().Reveal()
	digest := sha256.Sum256([]byte(raw))
	var digestCount int
	if err = pool.QueryRow(ctx, "SELECT count(*) FROM omnexa_identity.api_credentials WHERE secret_digest = $1", digest[:]).Scan(&digestCount); err != nil || digestCount != 1 {
		t.Fatalf("credential digest persistence count/error = %d/%v", digestCount, err)
	}
	var forbiddenColumns int
	if err = pool.QueryRow(
		ctx,
		`SELECT count(*) FROM information_schema.columns
		 WHERE table_schema = 'omnexa_identity'
		   AND table_name = 'api_credentials'
		   AND column_name IN ('secret','raw_secret','api_key','api_secret','credential_value')`,
	).Scan(&forbiddenColumns); err != nil || forbiddenColumns != 0 {
		t.Fatalf("raw-secret persistence columns count/error = %d/%v", forbiddenColumns, err)
	}

	parsed, _ := identity.ParseAPICredentialSecret(raw)
	authenticated, err := serviceIdentity.VerifyCredential(ctx, parsed, binding, issuedAt.Add(time.Second))
	if err != nil {
		t.Fatalf("VerifyCredential() error = %v", err)
	}
	wrongTenant, _ := identity.NewServiceAccountBinding("01890f3e-7b9a-7cc0-98c4-dc0c0c073992", "")
	if _, err = serviceIdentity.VerifyCredential(ctx, parsed, wrongTenant, issuedAt.Add(2*time.Second)); err == nil {
		t.Fatal("wrong-tenant credential verification unexpectedly succeeded")
	}
	wrongOrganization, _ := identity.NewServiceAccountBinding(tenantID, "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a1")
	if _, err = serviceIdentity.VerifyCredential(ctx, parsed, wrongOrganization, issuedAt.Add(2*time.Second)); err == nil {
		t.Fatal("wrong-organization credential verification unexpectedly succeeded")
	}

	target, err := ServiceAccountSubjectFromAccount(account)
	if err != nil {
		t.Fatalf("ServiceAccountSubjectFromAccount() error = %v", err)
	}
	metadata := MutationMetadata{CorrelationID: "p02-08-integration", Reason: "assign bounded service-account test permission"}
	role, err := rbac.CreateRole(ctx, actor, "service reader", []PermissionID{PermissionRoleRead}, metadata)
	if err != nil {
		t.Fatalf("CreateRole(service reader) error = %v", err)
	}
	assignment, err := rbac.AssignRoleToServiceAccount(ctx, actor, target, role.ID(), metadata)
	if err != nil {
		t.Fatalf("AssignRoleToServiceAccount() error = %v", err)
	}
	runtimeSubject, err := ServiceAccountSubjectFromAuthentication(authenticated)
	if err != nil {
		t.Fatalf("ServiceAccountSubjectFromAuthentication() error = %v", err)
	}
	if decision, checkErr := rbac.CheckServiceAccount(ctx, runtimeSubject, PermissionRoleRead); checkErr != nil || decision != DecisionAllow {
		t.Fatalf("service-account allowed RBAC = %q/%v", decision, checkErr)
	}
	if decision, checkErr := rbac.CheckServiceAccount(ctx, runtimeSubject, PermissionAssignmentManage); checkErr != nil || decision != DecisionDeny {
		t.Fatalf("service-account wrong-scope permission = %q/%v", decision, checkErr)
	}
	if _, err = rbac.RevokeServiceAccountAssignment(ctx, actor, assignment.ID(), metadata); err != nil {
		t.Fatalf("RevokeServiceAccountAssignment() error = %v", err)
	}
	if decision, checkErr := rbac.CheckServiceAccount(ctx, runtimeSubject, PermissionRoleRead); checkErr != nil || decision != DecisionDeny {
		t.Fatalf("revoked service assignment = %q/%v", decision, checkErr)
	}

	rotatedAt := issuedAt.Add(5 * time.Minute)
	rotated, err := serviceIdentity.RotateCredential(ctx, authenticated, rotatedAt.Add(time.Hour), rotatedAt)
	if err != nil {
		t.Fatalf("RotateCredential() error = %v", err)
	}
	if _, err = serviceIdentity.VerifyCredential(ctx, parsed, binding, rotatedAt.Add(time.Second)); err == nil {
		t.Fatal("superseded credential remained valid")
	}
	rotatedParsed, _ := identity.ParseAPICredentialSecret(rotated.Secret().Reveal())
	rotatedAuth, err := serviceIdentity.VerifyCredential(ctx, rotatedParsed, binding, rotatedAt.Add(time.Second))
	if err != nil {
		t.Fatalf("rotated credential verify error = %v", err)
	}
	revokedAt := rotatedAt.Add(2 * time.Second)
	if _, err = serviceIdentity.RevokeCredential(ctx, rotatedAuth, revokedAt); err != nil {
		t.Fatalf("RevokeCredential() error = %v", err)
	}
	if _, err = serviceIdentity.VerifyCredential(ctx, rotatedParsed, binding, revokedAt.Add(time.Second)); err == nil {
		t.Fatal("revoked credential remained valid")
	}

	expiring, err := serviceIdentity.IssueCredential(ctx, account.ID(), revokedAt.Add(2*time.Second), revokedAt.Add(time.Second))
	if err != nil {
		t.Fatalf("IssueCredential(expiring) error = %v", err)
	}
	expiringParsed, _ := identity.ParseAPICredentialSecret(expiring.Secret().Reveal())
	if _, err = serviceIdentity.VerifyCredential(ctx, expiringParsed, binding, expiring.Credential().ExpiresAt()); err == nil {
		t.Fatal("expired credential remained valid")
	}
}

func readP0208Migration(t *testing.T, version int64, name, path string) database.Migration {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read migration %s error = %v", path, err)
	}
	return database.Migration{Version: version, Name: name, SQL: string(contents)}
}

func runP0208Migrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool, owner string, migrations []database.Migration) {
	t.Helper()
	migrator, err := database.NewMigrator(pool, owner, migrations, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(%s) error = %v", owner, err)
	}
	if err = migrator.Run(ctx); err != nil {
		t.Fatalf("migrator Run(%s) error = %v", owner, err)
	}
}

func assertP0208LedgerCount(t *testing.T, ctx context.Context, pool *pgxpool.Pool, owner string, want int) {
	t.Helper()
	var got int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = $1", owner).Scan(&got); err != nil || got != want {
		t.Fatalf("migration ledger %s count/error = %d/%v, want %d", owner, got, err, want)
	}
}

func seedP0208HumanAuthority(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	at time.Time,
	tenantID, actorID, membershipID, roleID, assignmentID string,
) {
	t.Helper()
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO omnexa_identity.principals (id, principal_type, lifecycle_state, created_at, updated_at) VALUES ($1,'human_user','active',$2,$2)`, []any{actorID, at}},
		{`INSERT INTO omnexa_identity.users (principal_id, primary_email) VALUES ($1,'p0208-admin@example.com')`, []any{actorID}},
		{`INSERT INTO omnexa_tenancy.tenants (id, lifecycle_state, created_at, updated_at) VALUES ($1,'active',$2,$2)`, []any{tenantID, at}},
		{`INSERT INTO omnexa_tenancy.tenant_memberships (id, tenant_id, principal_id, relationship_state, created_at, updated_at) VALUES ($1,$2,$3,'active',$4,$4)`, []any{membershipID, tenantID, actorID, at}},
		{`INSERT INTO omnexa_authorization.roles (id, tenant_id, organization_id, name, created_at, updated_at) VALUES ($1,$2,NULL,'p0208 bootstrap admin',$3,$3)`, []any{roleID, tenantID, at}},
	}
	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement.query, statement.args...); err != nil {
			t.Fatalf("seed P02.08 prerequisite error = %v", err)
		}
	}
	for _, permission := range []PermissionID{PermissionRoleRead, PermissionRoleManage, PermissionAssignmentRead, PermissionAssignmentManage} {
		if _, err := pool.Exec(ctx, `INSERT INTO omnexa_authorization.role_permissions (role_id, permission_id) VALUES ($1,$2)`, roleID, string(permission)); err != nil {
			t.Fatalf("seed permission %s error = %v", permission, err)
		}
	}
	if _, err := pool.Exec(
		ctx,
		`INSERT INTO omnexa_authorization.role_assignments (id, role_id, principal_id, assignment_state, created_at, updated_at)
		 VALUES ($1,$2,$3,'active',$4,$4)`,
		assignmentID, roleID, actorID, at,
	); err != nil {
		t.Fatalf("seed admin assignment error = %v", err)
	}
}

func resetP0208Database(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	for _, schema := range []string{"omnexa_authorization", "omnexa_organization", "omnexa_tenancy", "omnexa_identity"} {
		if _, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS "+schema+" CASCADE"); err != nil {
			t.Fatalf("drop schema %s error = %v", schema, err)
		}
	}
	if _, err := pool.Exec(ctx, `DELETE FROM omnexa_kernel.schema_migrations WHERE owner IN ('kernel.identity','kernel.tenancy','kernel.organization','kernel.authorization')`); err != nil {
		t.Fatalf("reset migration ledger error = %v", err)
	}
}

var _ = failure.IsCode

package identity

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresIdentityFoundationIntegration(t *testing.T) {
	databaseURL := os.Getenv("P02_01_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("P02_01_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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

	if _, dropErr := pool.Exec(ctx, "DROP SCHEMA IF EXISTS omnexa_identity CASCADE"); dropErr != nil {
		t.Fatalf("drop identity schema error = %v", dropErr)
	}
	if _, resetErr := pool.Exec(ctx, "DELETE FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.identity'"); resetErr != nil {
		t.Fatalf("reset identity migration ledger error = %v", resetErr)
	}
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_, _ = pool.Exec(cleanupCtx, "DROP SCHEMA IF EXISTS omnexa_identity CASCADE")
		_, _ = pool.Exec(cleanupCtx, "DELETE FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.identity'")
	}()

	migrationSQL, migrationReadErr := os.ReadFile("../../migrations/kernel.identity/1_create_identity_foundation.sql")
	if migrationReadErr != nil {
		t.Fatalf("read identity migration error = %v", migrationReadErr)
	}
	migrations := []database.Migration{{
		Version: 1,
		Name:    "create_identity_foundation",
		SQL:     string(migrationSQL),
	}}
	migrator, migratorErr := database.NewMigrator(pool, "kernel.identity", migrations, 5*time.Second)
	if migratorErr != nil {
		t.Fatalf("NewMigrator(kernel.identity) error = %v", migratorErr)
	}
	if runErr := migrator.Run(ctx); runErr != nil {
		t.Fatalf("fresh identity migration error = %v", runErr)
	}
	if runErr := migrator.Run(ctx); runErr != nil {
		t.Fatalf("idempotent identity migration error = %v", runErr)
	}

	var ledgerCount int
	if queryErr := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.identity' AND version = 1",
	).Scan(&ledgerCount); queryErr != nil {
		t.Fatalf("identity ledger query error = %v", queryErr)
	}
	if ledgerCount != 1 {
		t.Fatalf("identity migration ledger count = %d, want 1", ledgerCount)
	}

	repository, repositoryErr := NewPostgresRepository(pool)
	if repositoryErr != nil {
		t.Fatalf("NewPostgresRepository() error = %v", repositoryErr)
	}
	createdAt := time.Date(2026, time.August, 23, 10, 0, 0, 0, time.UTC)
	user, userErr := newUserAt(fixedUserID, "owner@example.com", createdAt)
	if userErr != nil {
		t.Fatalf("newUserAt() error = %v", userErr)
	}
	if createErr := repository.Create(ctx, user); createErr != nil {
		t.Fatalf("repository.Create() error = %v", createErr)
	}

	stored, getErr := repository.Get(ctx, user.ID())
	if getErr != nil {
		t.Fatalf("repository.Get() error = %v", getErr)
	}
	if stored.ID() != user.ID() || stored.PrimaryEmail() != "owner@example.com" || stored.State() != LifecycleProvisioned {
		t.Fatalf("stored User mismatch: id=%q state=%q", stored.ID(), stored.State())
	}

	active, transitionErr := repository.Transition(
		ctx,
		user.ID(),
		LifecycleProvisioned,
		LifecycleActive,
		createdAt.Add(time.Minute),
	)
	if transitionErr != nil {
		t.Fatalf("repository.Transition() error = %v", transitionErr)
	}
	if active.State() != LifecycleActive {
		t.Fatalf("transition state = %q, want %q", active.State(), LifecycleActive)
	}
	if _, staleErr := repository.Transition(
		ctx,
		user.ID(),
		LifecycleProvisioned,
		LifecycleActive,
		createdAt.Add(2*time.Minute),
	); !failure.IsCode(staleErr, codeUserConflict) {
		t.Fatalf("stale transition error = %v, want %s", staleErr, codeUserConflict)
	}
	if duplicateErr := repository.Create(ctx, user); !failure.IsCode(duplicateErr, codeUserConflict) {
		t.Fatalf("duplicate create error = %v, want %s", duplicateErr, codeUserConflict)
	}

	missingID := UserID("01890f3e-7b9a-7cc0-98c4-dc0c0c073990")
	if _, missingErr := repository.Get(ctx, missingID); !failure.IsCode(missingErr, codeUserNotFound) {
		t.Fatalf("missing Get() error = %v, want %s", missingErr, codeUserNotFound)
	}

	var forbiddenColumns int
	if queryErr := pool.QueryRow(
		ctx,
		`SELECT count(*)
		 FROM information_schema.columns
		 WHERE table_schema = 'omnexa_identity'
		   AND column_name IN (
		       'tenant_id', 'organization_id', 'password', 'password_hash',
		       'session_id', 'role_id', 'permission_id', 'mfa_secret', 'api_key'
		   )`,
	).Scan(&forbiddenColumns); queryErr != nil {
		t.Fatalf("forbidden-column query error = %v", queryErr)
	}
	if forbiddenColumns != 0 {
		t.Fatalf("identity schema contains %d premature authority/credential columns", forbiddenColumns)
	}

	var nonHumanRows int
	if queryErr := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_identity.principals WHERE principal_type <> 'human_user'",
	).Scan(&nonHumanRows); queryErr != nil {
		t.Fatalf("non-human principal query error = %v", queryErr)
	}
	if nonHumanRows != 0 {
		t.Fatalf("P02.01 persisted %d unauthorized non-human principals", nonHumanRows)
	}

	drifted := []database.Migration{{
		Version: 1,
		Name:    "create_identity_foundation",
		SQL:     string(migrationSQL) + "\n-- rewritten migration must fail closed\n",
	}}
	driftMigrator, driftMigratorErr := database.NewMigrator(pool, "kernel.identity", drifted, 5*time.Second)
	if driftMigratorErr != nil {
		t.Fatalf("NewMigrator(drifted) error = %v", driftMigratorErr)
	}
	if driftErr := driftMigrator.Run(ctx); driftErr == nil {
		t.Fatalf("rewritten applied identity migration unexpectedly succeeded")
	}
}

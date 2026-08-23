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

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("pool.Ping() error = %v", err)
	}

	foundationMigrator, err := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(kernel.foundation) error = %v", err)
	}
	if err := foundationMigrator.Run(ctx); err != nil {
		t.Fatalf("foundation migrator Run() error = %v", err)
	}

	if _, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS omnexa_identity CASCADE"); err != nil {
		t.Fatalf("drop identity schema error = %v", err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.identity'"); err != nil {
		t.Fatalf("reset identity migration ledger error = %v", err)
	}
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_, _ = pool.Exec(cleanupCtx, "DROP SCHEMA IF EXISTS omnexa_identity CASCADE")
		_, _ = pool.Exec(cleanupCtx, "DELETE FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.identity'")
	}()

	migrationSQL, err := os.ReadFile("../../migrations/kernel.identity/1_create_identity_foundation.sql")
	if err != nil {
		t.Fatalf("read identity migration error = %v", err)
	}
	migrations := []database.Migration{{
		Version: 1,
		Name:    "create_identity_foundation",
		SQL:     string(migrationSQL),
	}}
	migrator, err := database.NewMigrator(pool, "kernel.identity", migrations, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(kernel.identity) error = %v", err)
	}
	if err := migrator.Run(ctx); err != nil {
		t.Fatalf("fresh identity migration error = %v", err)
	}
	if err := migrator.Run(ctx); err != nil {
		t.Fatalf("idempotent identity migration error = %v", err)
	}

	var ledgerCount int
	if err := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_kernel.schema_migrations WHERE owner = 'kernel.identity' AND version = 1",
	).Scan(&ledgerCount); err != nil {
		t.Fatalf("identity ledger query error = %v", err)
	}
	if ledgerCount != 1 {
		t.Fatalf("identity migration ledger count = %d, want 1", ledgerCount)
	}

	repository, err := NewPostgresRepository(pool)
	if err != nil {
		t.Fatalf("NewPostgresRepository() error = %v", err)
	}
	createdAt := time.Date(2026, time.August, 23, 10, 0, 0, 0, time.UTC)
	user, err := newUserAt(fixedUserID, "owner@example.com", createdAt)
	if err != nil {
		t.Fatalf("newUserAt() error = %v", err)
	}
	if err := repository.Create(ctx, user); err != nil {
		t.Fatalf("repository.Create() error = %v", err)
	}

	stored, err := repository.Get(ctx, user.ID())
	if err != nil {
		t.Fatalf("repository.Get() error = %v", err)
	}
	if stored.ID() != user.ID() || stored.PrimaryEmail() != "owner@example.com" || stored.State() != LifecycleProvisioned {
		t.Fatalf("stored User mismatch: id=%q state=%q", stored.ID(), stored.State())
	}

	active, err := repository.Transition(
		ctx,
		user.ID(),
		LifecycleProvisioned,
		LifecycleActive,
		createdAt.Add(time.Minute),
	)
	if err != nil {
		t.Fatalf("repository.Transition() error = %v", err)
	}
	if active.State() != LifecycleActive {
		t.Fatalf("transition state = %q, want %q", active.State(), LifecycleActive)
	}
	if _, err := repository.Transition(
		ctx,
		user.ID(),
		LifecycleProvisioned,
		LifecycleActive,
		createdAt.Add(2*time.Minute),
	); !failure.IsCode(err, codeUserConflict) {
		t.Fatalf("stale transition error = %v, want %s", err, codeUserConflict)
	}
	if err := repository.Create(ctx, user); !failure.IsCode(err, codeUserConflict) {
		t.Fatalf("duplicate create error = %v, want %s", err, codeUserConflict)
	}

	missingID := UserID("01890f3e-7b9a-7cc0-98c4-dc0c0c073990")
	if _, err := repository.Get(ctx, missingID); !failure.IsCode(err, codeUserNotFound) {
		t.Fatalf("missing Get() error = %v, want %s", err, codeUserNotFound)
	}

	var forbiddenColumns int
	if err := pool.QueryRow(
		ctx,
		`SELECT count(*)
		 FROM information_schema.columns
		 WHERE table_schema = 'omnexa_identity'
		   AND column_name IN (
		       'tenant_id', 'organization_id', 'password', 'password_hash',
		       'session_id', 'role_id', 'permission_id', 'mfa_secret', 'api_key'
		   )`,
	).Scan(&forbiddenColumns); err != nil {
		t.Fatalf("forbidden-column query error = %v", err)
	}
	if forbiddenColumns != 0 {
		t.Fatalf("identity schema contains %d premature authority/credential columns", forbiddenColumns)
	}

	var nonHumanRows int
	if err := pool.QueryRow(
		ctx,
		"SELECT count(*) FROM omnexa_identity.principals WHERE principal_type <> 'human_user'",
	).Scan(&nonHumanRows); err != nil {
		t.Fatalf("non-human principal query error = %v", err)
	}
	if nonHumanRows != 0 {
		t.Fatalf("P02.01 persisted %d unauthorized non-human principals", nonHumanRows)
	}

	drifted := []database.Migration{{
		Version: 1,
		Name:    "create_identity_foundation",
		SQL:     string(migrationSQL) + "\n-- rewritten migration must fail closed\n",
	}}
	driftMigrator, err := database.NewMigrator(pool, "kernel.identity", drifted, 5*time.Second)
	if err != nil {
		t.Fatalf("NewMigrator(drifted) error = %v", err)
	}
	if err := driftMigrator.Run(ctx); err == nil {
		t.Fatalf("rewritten applied identity migration unexpectedly succeeded")
	}
}

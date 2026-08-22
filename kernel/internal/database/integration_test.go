package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const integrationDatabaseEnv = "P01_04_TEST_DATABASE_URL"

func TestUnavailableConnectionIsBoundedAndSafe(t *testing.T) {
	const secret = "bounded-unavailable-secret"
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			databaseURLKey:            "postgres://omnexa:" + secret + "@127.0.0.1:1/none?sslmode=disable",
			databaseConnectTimeoutKey: "200ms",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}

	started := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	pool, err := NewPool(ctx, resolved)
	if pool != nil {
		pool.Close()
	}
	if err == nil {
		t.Fatal("NewPool() error = nil, want unavailable connection failure")
	}
	assertFailureCode(t, err, codeConnectionUnavailable)
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("unavailable connection was not bounded: %v", elapsed)
	}
	if strings.Contains(err.Error(), secret) || strings.Contains(err.Error(), "postgres://") {
		t.Fatalf("connection failure leaked restricted configuration: %q", err)
	}
}

func TestPostgreSQLFoundationIntegration(t *testing.T) {
	if os.Getenv(integrationDatabaseEnv) == "" {
		t.Skip(integrationDatabaseEnv + " is not configured")
	}

	t.Run("server version and ping", func(t *testing.T) {
		pool := openIntegrationPool(t, 4)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		var version string
		if err := pool.QueryRow(ctx, "SHOW server_version").Scan(&version); err != nil {
			t.Fatalf("SHOW server_version error = %v", err)
		}
		t.Logf("PostgreSQL server_version=%s", version)
		if !strings.HasPrefix(version, "18.6") {
			t.Fatalf("PostgreSQL server_version = %q, want 18.6 baseline", version)
		}
	})

	t.Run("pool exhaustion is bounded", func(t *testing.T) {
		pool := openIntegrationPool(t, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		first, err := pool.Acquire(ctx)
		if err != nil {
			t.Fatalf("first Acquire() error = %v", err)
		}
		defer first.Release()

		waitCtx, waitCancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer waitCancel()
		started := time.Now()
		_, err = pool.Acquire(waitCtx)
		if err == nil {
			t.Fatal("second Acquire() error = nil, want bounded exhaustion failure")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("second Acquire() error = %v, want context deadline", err)
		}
		if elapsed := time.Since(started); elapsed > time.Second {
			t.Fatalf("pool exhaustion wait was not bounded: %v", elapsed)
		}
	})

	t.Run("transaction commit and rollback", func(t *testing.T) {
		pool := openIntegrationPool(t, 4)
		resetMigrationSchema(t, pool)
		bootstrap, err := NewMigrator(pool, "kernel.data", nil, time.Second)
		if err != nil {
			t.Fatalf("NewMigrator() error = %v", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := bootstrap.Run(ctx); err != nil {
			t.Fatalf("bootstrap Run() error = %v", err)
		}
		if _, err := pool.Exec(ctx, "CREATE TABLE omnexa_kernel.synthetic_transaction (id bigint PRIMARY KEY)"); err != nil {
			t.Fatalf("create synthetic transaction table error = %v", err)
		}

		if err := InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
			_, err := tx.Exec(ctx, "INSERT INTO omnexa_kernel.synthetic_transaction (id) VALUES (1)")
			return err
		}); err != nil {
			t.Fatalf("commit transaction error = %v", err)
		}

		sentinel := errors.New("synthetic callback rejection")
		err = InTransaction(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, "INSERT INTO omnexa_kernel.synthetic_transaction (id) VALUES (2)"); err != nil {
				return err
			}
			return sentinel
		})
		if !errors.Is(err, sentinel) {
			t.Fatalf("rollback transaction error = %v, want callback sentinel", err)
		}

		var count int
		if err := pool.QueryRow(ctx, "SELECT count(*) FROM omnexa_kernel.synthetic_transaction").Scan(&count); err != nil {
			t.Fatalf("transaction row count error = %v", err)
		}
		if count != 1 {
			t.Fatalf("transaction row count = %d, want 1 committed row", count)
		}
	})

	t.Run("fresh and idempotent migrations", func(t *testing.T) {
		pool := openIntegrationPool(t, 4)
		resetMigrationSchema(t, pool)
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		migrator := mustMigrator(t, pool, syntheticMigrations(), time.Second)
		if err := migrator.Run(ctx); err != nil {
			t.Fatalf("fresh Run() error = %v", err)
		}
		assertLedgerVersions(t, pool, []int64{1, 2})
		assertFixtureLabelColumn(t, pool, true)

		if err := migrator.Run(ctx); err != nil {
			t.Fatalf("idempotent Run() error = %v", err)
		}
		assertLedgerVersions(t, pool, []int64{1, 2})
		assertNoPulledForwardSchemas(t, pool)
	})

	t.Run("deterministic upgrade from prior synthetic version", func(t *testing.T) {
		pool := openIntegrationPool(t, 4)
		resetMigrationSchema(t, pool)
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		prior := mustMigrator(t, pool, syntheticMigrations()[:1], time.Second)
		if err := prior.Run(ctx); err != nil {
			t.Fatalf("prior-version Run() error = %v", err)
		}
		assertLedgerVersions(t, pool, []int64{1})
		assertFixtureLabelColumn(t, pool, false)

		current := mustMigrator(t, pool, syntheticMigrations(), time.Second)
		if err := current.Run(ctx); err != nil {
			t.Fatalf("upgrade Run() error = %v", err)
		}
		assertLedgerVersions(t, pool, []int64{1, 2})
		assertFixtureLabelColumn(t, pool, true)
	})

	t.Run("failed migration rolls back and does not advance ledger", func(t *testing.T) {
		pool := openIntegrationPool(t, 4)
		resetMigrationSchema(t, pool)
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		prior := mustMigrator(t, pool, syntheticMigrations()[:1], time.Second)
		if err := prior.Run(ctx); err != nil {
			t.Fatalf("prior-version Run() error = %v", err)
		}
		failing := []Migration{
			syntheticMigrations()[0],
			{
				Version: 2,
				Name:    "fail_after_change",
				SQL: `DO $$
BEGIN
    ALTER TABLE omnexa_kernel.synthetic_fixture ADD COLUMN doomed boolean NOT NULL DEFAULT false;
    RAISE EXCEPTION 'synthetic migration failure';
END
$$;`,
			},
		}
		err := mustMigrator(t, pool, failing, time.Second).Run(ctx)
		if err == nil {
			t.Fatal("failing Run() error = nil")
		}
		assertFailureCode(t, err, codeMigrationApply)
		if strings.Contains(err.Error(), "synthetic migration failure") {
			t.Fatalf("migration failure published provider SQL detail: %q", err)
		}
		assertLedgerVersions(t, pool, []int64{1})
		assertColumnExists(t, pool, "doomed", false)
	})

	t.Run("ledger detects rewritten migration", func(t *testing.T) {
		pool := openIntegrationPool(t, 4)
		resetMigrationSchema(t, pool)
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		original := mustMigrator(t, pool, syntheticMigrations()[:1], time.Second)
		if err := original.Run(ctx); err != nil {
			t.Fatalf("original Run() error = %v", err)
		}
		rewritten := []Migration{{Version: 1, Name: "create_fixture", SQL: "CREATE TABLE omnexa_kernel.synthetic_fixture (id text PRIMARY KEY);"}}
		err := mustMigrator(t, pool, rewritten, time.Second).Run(ctx)
		if err == nil {
			t.Fatal("rewritten Run() error = nil, want drift failure")
		}
		assertFailureCode(t, err, codeMigrationDrift)
		assertLedgerVersions(t, pool, []int64{1})
	})

	t.Run("advisory lock contention is bounded", func(t *testing.T) {
		pool := openIntegrationPool(t, 4)
		resetMigrationSchema(t, pool)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		holder, err := pool.Acquire(ctx)
		if err != nil {
			t.Fatalf("Acquire() lock holder error = %v", err)
		}
		defer holder.Release()
		lockKey := advisoryLockKey("kernel.data")
		if _, err := holder.Exec(ctx, "SELECT pg_advisory_lock($1)", lockKey); err != nil {
			t.Fatalf("manual advisory lock error = %v", err)
		}
		defer func() {
			var unlocked bool
			_ = holder.QueryRow(context.Background(), "SELECT pg_advisory_unlock($1)", lockKey).Scan(&unlocked)
		}()

		migrator := mustMigrator(t, pool, nil, 150*time.Millisecond)
		started := time.Now()
		err = migrator.Run(ctx)
		if err == nil {
			t.Fatal("contended Run() error = nil")
		}
		assertFailureCode(t, err, codeMigrationLock)
		if elapsed := time.Since(started); elapsed > time.Second {
			t.Fatalf("migration lock contention was not bounded: %v", elapsed)
		}
	})
}

func openIntegrationPool(t *testing.T, maxConnections int) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv(integrationDatabaseEnv)
	if dsn == "" {
		t.Skip(integrationDatabaseEnv + " is not configured")
	}
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			databaseURLKey:                   dsn,
			databaseConnectTimeoutKey:        "2s",
			databaseMaxConnectionsKey:        strconv.Itoa(maxConnections),
			databaseMinConnectionsKey:        "0",
			databaseMaxConnectionLifetimeKey: "10m",
			databaseMaxConnectionIdleTimeKey: "2m",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	pool, err := NewPool(ctx, resolved)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func resetMigrationSchema(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS omnexa_kernel CASCADE"); err != nil {
		t.Fatalf("reset migration schema error = %v", err)
	}
}

func syntheticMigrations() []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "create_fixture",
			SQL:     "CREATE TABLE omnexa_kernel.synthetic_fixture (id bigint PRIMARY KEY);",
		},
		{
			Version: 2,
			Name:    "add_fixture_label",
			SQL:     "ALTER TABLE omnexa_kernel.synthetic_fixture ADD COLUMN label text NOT NULL DEFAULT 'synthetic';",
		},
	}
}

func mustMigrator(t *testing.T, pool *pgxpool.Pool, migrations []Migration, timeout time.Duration) *Migrator {
	t.Helper()
	migrator, err := NewMigrator(pool, "kernel.data", migrations, timeout)
	if err != nil {
		t.Fatalf("NewMigrator() error = %v", err)
	}
	return migrator
}

func assertLedgerVersions(t *testing.T, pool *pgxpool.Pool, want []int64) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := pool.Query(ctx, "SELECT version FROM omnexa_kernel.schema_migrations WHERE owner = $1 ORDER BY version", "kernel.data")
	if err != nil {
		t.Fatalf("ledger query error = %v", err)
	}
	defer rows.Close()
	var got []int64
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			t.Fatalf("ledger Scan() error = %v", err)
		}
		got = append(got, version)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("ledger rows error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ledger versions = %v, want %v", got, want)
	}
}

func assertFixtureLabelColumn(t *testing.T, pool *pgxpool.Pool, want bool) {
	t.Helper()
	assertColumnExists(t, pool, "label", want)
}

func assertColumnExists(t *testing.T, pool *pgxpool.Pool, column string, want bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var exists bool
	err := pool.QueryRow(ctx, `
SELECT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'omnexa_kernel'
      AND table_name = 'synthetic_fixture'
      AND column_name = $1
)`, column).Scan(&exists)
	if err != nil {
		t.Fatalf("column existence query error = %v", err)
	}
	if exists != want {
		t.Fatalf("column %q exists = %v, want %v", column, exists, want)
	}
}

func assertNoPulledForwardSchemas(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := pool.Query(ctx, `
SELECT nspname
FROM pg_namespace
WHERE nspname NOT LIKE 'pg_%'
  AND nspname <> 'information_schema'
ORDER BY nspname`)
	if err != nil {
		t.Fatalf("schema scope query error = %v", err)
	}
	defer rows.Close()
	allowed := map[string]bool{"public": true, migrationSchema: true}
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			t.Fatalf("schema Scan() error = %v", err)
		}
		if !allowed[schema] {
			t.Fatalf("P01.04 created or encountered out-of-scope schema %q", schema)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("schema rows error = %v", err)
	}

	var owner string
	var table string
	err = pool.QueryRow(ctx, `
SELECT schemaname, tablename
FROM pg_tables
WHERE schemaname = 'omnexa_kernel'
  AND tablename = 'schema_migrations'`).Scan(&owner, &table)
	if err != nil {
		t.Fatalf("foundation ledger table missing: %v", err)
	}
	if owner != migrationSchema || table != migrationTable {
		t.Fatalf("unexpected foundation ledger identity %s.%s", owner, table)
	}

	for _, forbidden := range []string{"tenant", "organization", "module", "business", "outbox", "inbox", "cache", "storage", "telemetry"} {
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM pg_tables WHERE lower(schemaname) LIKE '%%%s%%' OR lower(tablename) LIKE '%%%s%%'`, forbidden, forbidden)
		if err := pool.QueryRow(ctx, query).Scan(&count); err != nil {
			t.Fatalf("scope guard query for %q error = %v", forbidden, err)
		}
		if count != 0 {
			t.Fatalf("P01.04 scope guard found forbidden %q table/schema marker", forbidden)
		}
	}
}

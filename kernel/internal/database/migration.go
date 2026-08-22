package database

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	migrationSchema = "omnexa_kernel"
	migrationTable  = "schema_migrations"
)

var (
	migrationOwnerPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:\.[a-z][a-z0-9]*)*$`)
	migrationNamePattern  = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)
)

const ensureMigrationLedgerSQL = `
CREATE SCHEMA IF NOT EXISTS omnexa_kernel;
CREATE TABLE IF NOT EXISTS omnexa_kernel.schema_migrations (
    owner text NOT NULL,
    version bigint NOT NULL CHECK (version > 0),
    name text NOT NULL,
    checksum char(64) NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (owner, version)
);`

// Migration is one immutable, owner-scoped schema transition. SQL must be
// reviewable source code; the runner never discovers or executes hidden SQL.
type Migration struct {
	Version int64
	Name    string
	SQL     string
}

func (migration Migration) checksum() string {
	digest := sha256.Sum256([]byte(migration.SQL))
	return hex.EncodeToString(digest[:])
}

// Migrator coordinates and applies one authoritative owner's ordered migration
// set. The P01.04 foundation ledger itself lives in omnexa_kernel only.
type Migrator struct {
	pool        *pgxpool.Pool
	owner       string
	migrations  []Migration
	lockTimeout time.Duration
}

// NewMigrator validates an immutable migration sequence. Versions must start
// above zero and be strictly increasing; previously applied entries may never be
// removed or rewritten without triggering drift detection.
func NewMigrator(pool *pgxpool.Pool, owner string, migrations []Migration, lockTimeout time.Duration) (*Migrator, error) {
	if pool == nil || !migrationOwnerPattern.MatchString(owner) || lockTimeout <= 0 {
		return nil, safeFailure(
			codeMigrationInvalid,
			failure.CategoryValidation,
			"database migration configuration is invalid",
		)
	}

	copyOfMigrations := append([]Migration(nil), migrations...)
	var previous int64
	for index, migration := range copyOfMigrations {
		if migration.Version <= 0 || migration.Version <= previous || !migrationNamePattern.MatchString(migration.Name) || strings.TrimSpace(migration.SQL) == "" {
			return nil, safeFailure(
				codeMigrationInvalid,
				failure.CategoryValidation,
				"database migration set is invalid",
			)
		}
		if index > 0 && migration.Version == previous {
			return nil, safeFailure(
				codeMigrationInvalid,
				failure.CategoryValidation,
				"database migration set is invalid",
			)
		}
		previous = migration.Version
	}

	return &Migrator{
		pool:        pool,
		owner:       owner,
		migrations:  copyOfMigrations,
		lockTimeout: lockTimeout,
	}, nil
}

// Run acquires one dedicated PostgreSQL session, obtains an owner-specific
// advisory lock, reconciles the immutable ledger, and applies pending versions.
func (migrator *Migrator) Run(ctx context.Context) error {
	if migrator == nil || migrator.pool == nil {
		return safeFailure(
			codeMigrationInvalid,
			failure.CategoryValidation,
			"database migration configuration is invalid",
		)
	}

	lockCtx, cancel := context.WithTimeout(ctx, migrator.lockTimeout)
	connection, err := migrator.pool.Acquire(lockCtx)
	if err != nil {
		cancel()
		return classifyMigrationLockFailure(err)
	}

	release := true
	defer func() {
		if release {
			connection.Release()
		}
	}()

	lockKey := advisoryLockKey(migrator.owner)
	if _, err := connection.Exec(lockCtx, "SELECT pg_advisory_lock($1)", lockKey); err != nil {
		cancel()
		return classifyMigrationLockFailure(err)
	}
	cancel()

	runErr := migrator.runLocked(ctx, connection)

	unlockCtx, unlockCancel := context.WithTimeout(context.WithoutCancel(ctx), migrator.lockTimeout)
	var unlocked bool
	unlockErr := connection.QueryRow(unlockCtx, "SELECT pg_advisory_unlock($1)", lockKey).Scan(&unlocked)
	unlockCancel()
	if unlockErr != nil || !unlocked {
		raw := connection.Hijack()
		release = false
		closeCtx, closeCancel := context.WithTimeout(context.Background(), migrator.lockTimeout)
		_ = raw.Close(closeCtx)
		closeCancel()
		if runErr != nil {
			return runErr
		}
		if unlockErr == nil {
			unlockErr = errors.New("database migration advisory lock was not owned by the session")
		}
		return safeWrappedFailure(
			unlockErr,
			codeMigrationLock,
			failure.CategoryInternal,
			"database migration lock could not be released safely",
		)
	}

	return runErr
}

func (migrator *Migrator) runLocked(ctx context.Context, connection *pgxpool.Conn) error {
	if _, err := connection.Exec(ctx, ensureMigrationLedgerSQL); err != nil {
		return safeWrappedFailure(
			err,
			codeMigrationLedger,
			failure.CategoryDependency,
			"database migration ledger could not be prepared",
		)
	}

	applied, err := migrator.readLedger(ctx, connection)
	if err != nil {
		return err
	}

	expected := make(map[int64]Migration, len(migrator.migrations))
	for _, migration := range migrator.migrations {
		expected[migration.Version] = migration
	}
	for version, record := range applied {
		migration, exists := expected[version]
		if !exists || record.name != migration.Name || record.checksum != migration.checksum() {
			return safeFailure(
				codeMigrationDrift,
				failure.CategoryInvariant,
				"database migration history drift was detected",
			)
		}
	}

	for _, migration := range migrator.migrations {
		if _, exists := applied[migration.Version]; exists {
			continue
		}
		if err := migrator.apply(ctx, connection, migration); err != nil {
			return err
		}
	}
	return nil
}

type ledgerRecord struct {
	name     string
	checksum string
}

func (migrator *Migrator) readLedger(ctx context.Context, connection *pgxpool.Conn) (map[int64]ledgerRecord, error) {
	rows, err := connection.Query(
		ctx,
		"SELECT version, name, checksum FROM omnexa_kernel.schema_migrations WHERE owner = $1 ORDER BY version",
		migrator.owner,
	)
	if err != nil {
		return nil, safeWrappedFailure(
			err,
			codeMigrationLedger,
			failure.CategoryDependency,
			"database migration ledger could not be read",
		)
	}
	defer rows.Close()

	result := make(map[int64]ledgerRecord)
	for rows.Next() {
		var version int64
		var name string
		var checksum string
		if err := rows.Scan(&version, &name, &checksum); err != nil {
			return nil, safeWrappedFailure(
				err,
				codeMigrationLedger,
				failure.CategoryDependency,
				"database migration ledger could not be read",
			)
		}
		result[version] = ledgerRecord{name: name, checksum: checksum}
	}
	if err := rows.Err(); err != nil {
		return nil, safeWrappedFailure(
			err,
			codeMigrationLedger,
			failure.CategoryDependency,
			"database migration ledger could not be read",
		)
	}
	return result, nil
}

func (migrator *Migrator) apply(ctx context.Context, connection *pgxpool.Conn, migration Migration) error {
	tx, err := connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return safeWrappedFailure(
			err,
			codeMigrationApply,
			failure.CategoryDependency,
			"database migration could not begin",
		)
	}

	finished := false
	defer func() {
		if finished {
			return
		}
		rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), transactionRollbackTimeout)
		defer cancel()
		_ = tx.Rollback(rollbackCtx)
	}()

	if _, err := tx.Exec(ctx, migration.SQL); err != nil {
		finished = true
		return migrationRollbackFailure(ctx, tx, err)
	}
	if _, err := tx.Exec(
		ctx,
		"INSERT INTO omnexa_kernel.schema_migrations (owner, version, name, checksum) VALUES ($1, $2, $3, $4)",
		migrator.owner,
		migration.Version,
		migration.Name,
		migration.checksum(),
	); err != nil {
		finished = true
		return migrationRollbackFailure(ctx, tx, err)
	}
	if err := tx.Commit(ctx); err != nil {
		finished = true
		return safeWrappedFailure(
			err,
			codeMigrationApply,
			failure.CategoryDependency,
			"database migration could not be committed",
		)
	}
	finished = true
	return nil
}

func migrationRollbackFailure(ctx context.Context, tx pgx.Tx, cause error) error {
	rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), transactionRollbackTimeout)
	rollbackErr := tx.Rollback(rollbackCtx)
	cancel()
	if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
		cause = errors.Join(cause, rollbackErr)
	}
	return safeWrappedFailure(
		cause,
		codeMigrationApply,
		failure.CategoryDependency,
		"database migration could not be applied",
	)
}

func classifyMigrationLockFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) {
		return safeWrappedFailure(
			cause,
			codeMigrationLock,
			failure.CategoryTimeout,
			"database migration lock timed out",
			failure.WithRetryable(true),
		)
	}
	return safeWrappedFailure(
		cause,
		codeMigrationLock,
		failure.CategoryUnavailable,
		"database migration lock is unavailable",
		failure.WithRetryable(true),
	)
}

func advisoryLockKey(owner string) int64 {
	digest := sha256.Sum256([]byte("omnexa:migrations:" + owner))
	return int64(binary.BigEndian.Uint64(digest[:8]))
}

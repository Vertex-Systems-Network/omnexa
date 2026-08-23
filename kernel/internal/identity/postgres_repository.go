package identity

import (
	"context"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository persists only kernel.identity-owned principal/User state.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates the P02.01 owner-bounded PostgreSQL repository.
func NewPostgresRepository(pool *pgxpool.Pool) (*PostgresRepository, error) {
	if pool == nil {
		return nil, repositoryInvalidFailure()
	}
	return &PostgresRepository{pool: pool}, nil
}

// Create atomically persists the human principal and User-specific PII attribute.
func (repository *PostgresRepository) Create(ctx context.Context, user User) error {
	if repository == nil || repository.pool == nil {
		return repositoryInvalidFailure()
	}
	if err := user.validate(); err != nil {
		return err
	}

	return database.InTransaction(ctx, repository.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO omnexa_identity.principals
			 (id, principal_type, lifecycle_state, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5)`,
			string(user.id), string(user.principal), string(user.state), user.createdAt, user.updatedAt,
		); err != nil {
			return repositoryFailure(err)
		}
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO omnexa_identity.users (principal_id, primary_email)
			 VALUES ($1, $2)`,
			string(user.id), user.primaryEmail,
		); err != nil {
			return repositoryFailure(err)
		}
		return nil
	})
}

// Get retrieves one human User by UUIDv7 identifier. No tenant or authorization
// decision is made by this repository.
func (repository *PostgresRepository) Get(ctx context.Context, id UserID) (User, error) {
	if repository == nil || repository.pool == nil {
		return User{}, repositoryInvalidFailure()
	}
	if !id.Valid() {
		return User{}, invalidUserFailure()
	}

	var storedID string
	var principal string
	var state string
	var primaryEmail string
	var createdAt time.Time
	var updatedAt time.Time
	err := repository.pool.QueryRow(
		ctx,
		`SELECT p.id, p.principal_type, p.lifecycle_state, u.primary_email, p.created_at, p.updated_at
		 FROM omnexa_identity.principals AS p
		 JOIN omnexa_identity.users AS u ON u.principal_id = p.id
		 WHERE p.id = $1 AND p.principal_type = 'human_user'`,
		string(id),
	).Scan(&storedID, &principal, &state, &primaryEmail, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, userNotFoundFailure()
	}
	if err != nil {
		return User{}, repositoryFailure(err)
	}
	return rehydrateUser(UserID(storedID), PrincipalType(principal), LifecycleState(state), primaryEmail, createdAt, updatedAt)
}

// Transition atomically applies one validated lifecycle change using the caller's
// expected current state. Concurrent/stale transitions fail closed as conflicts.
func (repository *PostgresRepository) Transition(
	ctx context.Context,
	id UserID,
	from LifecycleState,
	to LifecycleState,
	changedAt time.Time,
) (User, error) {
	if repository == nil || repository.pool == nil {
		return User{}, repositoryInvalidFailure()
	}
	if !id.Valid() || !from.Valid() || !to.Valid() || !transitionAllowed(from, to) || changedAt.IsZero() {
		return User{}, transitionFailure()
	}
	instant := changedAt.UTC()

	var storedID string
	var principal string
	var state string
	var primaryEmail string
	var createdAt time.Time
	var updatedAt time.Time
	err := repository.pool.QueryRow(
		ctx,
		`UPDATE omnexa_identity.principals AS p
		 SET lifecycle_state = $3, updated_at = $4
		 FROM omnexa_identity.users AS u
		 WHERE p.id = $1
		   AND p.id = u.principal_id
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = $2
		   AND p.updated_at <= $4
		 RETURNING p.id, p.principal_type, p.lifecycle_state, u.primary_email, p.created_at, p.updated_at`,
		string(id), string(from), string(to), instant,
	).Scan(&storedID, &principal, &state, &primaryEmail, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		existing, getErr := repository.Get(ctx, id)
		if getErr != nil {
			return User{}, getErr
		}
		if existing.state != from || instant.Before(existing.updatedAt) {
			return User{}, userConflictFailure()
		}
		return User{}, invalidStoredUserFailure()
	}
	if err != nil {
		return User{}, repositoryFailure(err)
	}
	return rehydrateUser(UserID(storedID), PrincipalType(principal), LifecycleState(state), primaryEmail, createdAt, updatedAt)
}

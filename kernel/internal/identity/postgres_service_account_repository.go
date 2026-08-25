package identity

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresServiceAccountRepository persists only P02.08 kernel.identity-owned
// service-principal and credential state. Raw credential values never cross this API.
type PostgresServiceAccountRepository struct{ pool *pgxpool.Pool }

func NewPostgresServiceAccountRepository(pool *pgxpool.Pool) (*PostgresServiceAccountRepository, error) {
	if pool == nil {
		return nil, repositoryInvalidFailure()
	}
	return &PostgresServiceAccountRepository{pool: pool}, nil
}

func (repository *PostgresServiceAccountRepository) CreateServiceAccount(ctx context.Context, account ServiceAccount) error {
	if repository == nil || repository.pool == nil || account.validate() != nil {
		return serviceAccountInvalidFailure()
	}
	return database.InTransaction(ctx, repository.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO omnexa_identity.principals
			 (id, principal_type, lifecycle_state, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5)`,
			string(account.id), string(account.principal), string(account.state), account.createdAt, account.updatedAt,
		); err != nil {
			return serviceAccountRepositoryFailure(err)
		}
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO omnexa_identity.service_accounts
			 (principal_id, name, tenant_id, organization_id)
			 VALUES ($1, $2, $3, NULLIF($4, '')::uuid)`,
			string(account.id), account.name, account.binding.tenantID, account.binding.organizationID,
		); err != nil {
			return serviceAccountRepositoryFailure(err)
		}
		return nil
	})
}

func (repository *PostgresServiceAccountRepository) GetServiceAccount(ctx context.Context, id ServiceAccountID) (ServiceAccount, error) {
	if repository == nil || repository.pool == nil || !id.Valid() {
		return ServiceAccount{}, serviceAccountInvalidFailure()
	}
	var storedID, principal, state, name, tenantID, organizationID string
	var createdAt, updatedAt time.Time
	err := repository.pool.QueryRow(
		ctx,
		`SELECT p.id, p.principal_type, p.lifecycle_state, sa.name,
		        sa.tenant_id::text, COALESCE(sa.organization_id::text, ''),
		        p.created_at, p.updated_at
		 FROM omnexa_identity.principals p
		 JOIN omnexa_identity.service_accounts sa ON sa.principal_id = p.id
		 WHERE p.id = $1 AND p.principal_type = 'service_account'`,
		string(id),
	).Scan(&storedID, &principal, &state, &name, &tenantID, &organizationID, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ServiceAccount{}, serviceAccountNotFoundFailure()
	}
	if err != nil {
		return ServiceAccount{}, serviceAccountRepositoryFailure(err)
	}
	binding, err := NewServiceAccountBinding(tenantID, organizationID)
	if err != nil {
		return ServiceAccount{}, serviceAccountStoredInvalidFailure()
	}
	return rehydrateServiceAccount(ServiceAccountID(storedID), PrincipalType(principal), LifecycleState(state), name, binding, createdAt, updatedAt)
}

func (repository *PostgresServiceAccountRepository) TransitionServiceAccount(
	ctx context.Context,
	id ServiceAccountID,
	from LifecycleState,
	to LifecycleState,
	changedAt time.Time,
) (ServiceAccount, error) {
	if repository == nil || repository.pool == nil || !id.Valid() || !from.Valid() || !to.Valid() || !transitionAllowed(from, to) || changedAt.IsZero() {
		return ServiceAccount{}, serviceAccountTransitionFailure()
	}
	instant := changedAt.UTC()
	var storedID, principal, state, name, tenantID, organizationID string
	var createdAt, updatedAt time.Time
	err := repository.pool.QueryRow(
		ctx,
		`UPDATE omnexa_identity.principals p
		 SET lifecycle_state = $3, updated_at = $4
		 FROM omnexa_identity.service_accounts sa
		 WHERE p.id = $1
		   AND p.id = sa.principal_id
		   AND p.principal_type = 'service_account'
		   AND p.lifecycle_state = $2
		   AND p.updated_at <= $4
		 RETURNING p.id, p.principal_type, p.lifecycle_state, sa.name,
		           sa.tenant_id::text, COALESCE(sa.organization_id::text, ''),
		           p.created_at, p.updated_at`,
		string(id), string(from), string(to), instant,
	).Scan(&storedID, &principal, &state, &name, &tenantID, &organizationID, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ServiceAccount{}, serviceAccountConflictFailure()
	}
	if err != nil {
		return ServiceAccount{}, serviceAccountRepositoryFailure(err)
	}
	binding, err := NewServiceAccountBinding(tenantID, organizationID)
	if err != nil {
		return ServiceAccount{}, serviceAccountStoredInvalidFailure()
	}
	return rehydrateServiceAccount(ServiceAccountID(storedID), PrincipalType(principal), LifecycleState(state), name, binding, createdAt, updatedAt)
}

func (repository *PostgresServiceAccountRepository) CreateAPICredential(
	ctx context.Context,
	credential APICredential,
	digest [sha256.Size]byte,
) error {
	if repository == nil || repository.pool == nil || credential.validate() != nil {
		return apiCredentialInvalidFailure()
	}
	command, err := repository.pool.Exec(
		ctx,
		`INSERT INTO omnexa_identity.api_credentials
		 (id, service_account_id, secret_digest, created_at, expires_at)
		 SELECT $1, sa.principal_id, $3, $4, $5
		 FROM omnexa_identity.service_accounts sa
		 JOIN omnexa_identity.principals p ON p.id = sa.principal_id
		 WHERE sa.principal_id = $2
		   AND p.principal_type = 'service_account'
		   AND p.lifecycle_state = 'active'`,
		string(credential.id), string(credential.serviceAccountID), digest[:], credential.createdAt, credential.expiresAt,
	)
	if err != nil {
		return serviceAccountRepositoryFailure(err)
	}
	if command.RowsAffected() != 1 {
		return apiCredentialAuthenticationFailure()
	}
	return nil
}

func (repository *PostgresServiceAccountRepository) AuthenticateAPICredential(
	ctx context.Context,
	credentialID APICredentialID,
	digest [sha256.Size]byte,
	at time.Time,
) (ServiceAccount, APICredential, error) {
	if repository == nil || repository.pool == nil || !credentialID.Valid() || at.IsZero() {
		return ServiceAccount{}, APICredential{}, apiCredentialAuthenticationFailure()
	}
	tx, err := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ServiceAccount{}, APICredential{}, serviceAccountRepositoryFailure(err)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	account, credential, err := scanServiceCredential(ctx, tx, credentialID, digest, at.UTC())
	if err != nil {
		return ServiceAccount{}, APICredential{}, err
	}
	if _, err = tx.Exec(
		ctx,
		`UPDATE omnexa_identity.api_credentials
		 SET last_used_at = CASE
		     WHEN last_used_at IS NULL OR last_used_at < $2 THEN $2
		     ELSE last_used_at
		 END
		 WHERE id = $1`,
		string(credentialID), at.UTC(),
	); err != nil {
		return ServiceAccount{}, APICredential{}, serviceAccountRepositoryFailure(err)
	}
	credential.lastUsedAt = at.UTC()
	if err = tx.Commit(ctx); err != nil {
		return ServiceAccount{}, APICredential{}, serviceAccountRepositoryFailure(err)
	}
	return account, credential, nil
}

func scanServiceCredential(
	ctx context.Context,
	tx pgx.Tx,
	credentialID APICredentialID,
	digest [sha256.Size]byte,
	at time.Time,
) (ServiceAccount, APICredential, error) {
	var accountID, principal, state, name, tenantID, organizationID string
	var accountCreatedAt, accountUpdatedAt time.Time
	var storedCredentialID string
	var credentialCreatedAt, expiresAt time.Time
	var lastUsedAt, supersededAt, revokedAt *time.Time
	err := tx.QueryRow(
		ctx,
		`SELECT p.id, p.principal_type, p.lifecycle_state, sa.name,
		        sa.tenant_id::text, COALESCE(sa.organization_id::text, ''),
		        p.created_at, p.updated_at,
		        c.id, c.created_at, c.expires_at, c.last_used_at, c.superseded_at, c.revoked_at
		 FROM omnexa_identity.api_credentials c
		 JOIN omnexa_identity.service_accounts sa ON sa.principal_id = c.service_account_id
		 JOIN omnexa_identity.principals p ON p.id = sa.principal_id
		 WHERE c.id = $1
		   AND c.secret_digest = $2
		   AND c.revoked_at IS NULL
		   AND c.superseded_at IS NULL
		   AND c.expires_at > $3
		   AND p.principal_type = 'service_account'
		   AND p.lifecycle_state = 'active'
		 FOR UPDATE OF c, p`,
		string(credentialID), digest[:], at,
	).Scan(
		&accountID, &principal, &state, &name, &tenantID, &organizationID,
		&accountCreatedAt, &accountUpdatedAt,
		&storedCredentialID, &credentialCreatedAt, &expiresAt, &lastUsedAt, &supersededAt, &revokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return ServiceAccount{}, APICredential{}, apiCredentialAuthenticationFailure()
	}
	if err != nil {
		return ServiceAccount{}, APICredential{}, serviceAccountRepositoryFailure(err)
	}
	binding, err := NewServiceAccountBinding(tenantID, organizationID)
	if err != nil {
		return ServiceAccount{}, APICredential{}, serviceAccountStoredInvalidFailure()
	}
	account, err := rehydrateServiceAccount(ServiceAccountID(accountID), PrincipalType(principal), LifecycleState(state), name, binding, accountCreatedAt, accountUpdatedAt)
	if err != nil {
		return ServiceAccount{}, APICredential{}, err
	}
	credential, err := rehydrateAPICredential(
		APICredentialID(storedCredentialID), ServiceAccountID(accountID), credentialCreatedAt, expiresAt,
		optionalTime(lastUsedAt), optionalTime(supersededAt), optionalTime(revokedAt),
	)
	if err != nil {
		return ServiceAccount{}, APICredential{}, err
	}
	return account, credential, nil
}

func (repository *PostgresServiceAccountRepository) RotateAPICredential(
	ctx context.Context,
	accountID ServiceAccountID,
	currentID APICredentialID,
	next APICredential,
	nextDigest [sha256.Size]byte,
	rotatedAt time.Time,
) (APICredential, error) {
	if repository == nil || repository.pool == nil || !accountID.Valid() || !currentID.Valid() || next.validate() != nil ||
		next.serviceAccountID != accountID || rotatedAt.IsZero() {
		return APICredential{}, apiCredentialInvalidFailure()
	}
	tx, err := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return APICredential{}, serviceAccountRepositoryFailure(err)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	var createdAt, expiresAt time.Time
	var lastUsedAt *time.Time
	err = tx.QueryRow(
		ctx,
		`UPDATE omnexa_identity.api_credentials c
		 SET superseded_at = $3
		 FROM omnexa_identity.principals p
		 WHERE c.id = $1
		   AND c.service_account_id = $2
		   AND p.id = c.service_account_id
		   AND p.principal_type = 'service_account'
		   AND p.lifecycle_state = 'active'
		   AND c.revoked_at IS NULL
		   AND c.superseded_at IS NULL
		   AND c.expires_at > $3
		 RETURNING c.created_at, c.expires_at, c.last_used_at`,
		string(currentID), string(accountID), rotatedAt.UTC(),
	).Scan(&createdAt, &expiresAt, &lastUsedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return APICredential{}, apiCredentialConflictFailure()
	}
	if err != nil {
		return APICredential{}, serviceAccountRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.api_credentials
		 (id, service_account_id, secret_digest, created_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		string(next.id), string(accountID), nextDigest[:], next.createdAt, next.expiresAt,
	); err != nil {
		return APICredential{}, serviceAccountRepositoryFailure(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return APICredential{}, serviceAccountRepositoryFailure(err)
	}
	return rehydrateAPICredential(currentID, accountID, createdAt, expiresAt, optionalTime(lastUsedAt), rotatedAt.UTC(), time.Time{})
}

func (repository *PostgresServiceAccountRepository) RevokeAPICredential(
	ctx context.Context,
	accountID ServiceAccountID,
	credentialID APICredentialID,
	revokedAt time.Time,
) (APICredential, error) {
	if repository == nil || repository.pool == nil || !accountID.Valid() || !credentialID.Valid() || revokedAt.IsZero() {
		return APICredential{}, apiCredentialInvalidFailure()
	}
	var createdAt, expiresAt time.Time
	var lastUsedAt *time.Time
	err := repository.pool.QueryRow(
		ctx,
		`UPDATE omnexa_identity.api_credentials c
		 SET revoked_at = $3
		 FROM omnexa_identity.principals p
		 WHERE c.id = $1
		   AND c.service_account_id = $2
		   AND p.id = c.service_account_id
		   AND p.principal_type = 'service_account'
		   AND p.lifecycle_state = 'active'
		   AND c.revoked_at IS NULL
		   AND c.superseded_at IS NULL
		   AND c.expires_at > $3
		 RETURNING c.created_at, c.expires_at, c.last_used_at`,
		string(credentialID), string(accountID), revokedAt.UTC(),
	).Scan(&createdAt, &expiresAt, &lastUsedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return APICredential{}, apiCredentialConflictFailure()
	}
	if err != nil {
		return APICredential{}, serviceAccountRepositoryFailure(err)
	}
	return rehydrateAPICredential(credentialID, accountID, createdAt, expiresAt, optionalTime(lastUsedAt), time.Time{}, revokedAt.UTC())
}

func optionalTime(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return value.UTC()
}

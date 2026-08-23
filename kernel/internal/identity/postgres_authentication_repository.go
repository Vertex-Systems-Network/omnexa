package identity

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresAuthenticationRepository persists P02.04 authentication/session state
// only in kernel.identity-owned tables. Tenant/organization references are opaque
// context hints and are never resolved or written outside this schema.
type PostgresAuthenticationRepository struct{ pool *pgxpool.Pool }

// NewPostgresAuthenticationRepository returns the owner-bounded P02.04 repository.
func NewPostgresAuthenticationRepository(pool *pgxpool.Pool) (*PostgresAuthenticationRepository, error) {
	if pool == nil {
		return nil, repositoryInvalidFailure()
	}
	return &PostgresAuthenticationRepository{pool: pool}, nil
}

func (repository *PostgresAuthenticationRepository) EnrollPassword(
	ctx context.Context,
	userID UserID,
	passwordHash string,
	changedAt time.Time,
) (uint64, error) {
	if repository == nil || repository.pool == nil || !userID.Valid() || passwordHash == "" || changedAt.IsZero() {
		return 0, repositoryInvalidFailure()
	}
	var version uint64
	err := repository.pool.QueryRow(
		ctx,
		`INSERT INTO omnexa_identity.password_credentials (
			principal_id, password_hash, credential_version, created_at, updated_at
		 )
		 SELECT p.id, $2, 1, $3, $3
		 FROM omnexa_identity.principals p
		 WHERE p.id = $1
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state IN ('provisioned', 'active')
		 RETURNING credential_version`,
		string(userID),
		passwordHash,
		changedAt.UTC(),
	).Scan(&version)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, credentialConflictFailure()
	}
	if err != nil {
		return 0, authenticationRepositoryFailure(err)
	}
	return version, nil
}

func (repository *PostgresAuthenticationRepository) AuthenticationSnapshot(
	ctx context.Context,
	userID UserID,
) (authenticationSnapshot, error) {
	if repository == nil || repository.pool == nil || !userID.Valid() {
		return authenticationSnapshot{}, credentialNotFoundFailure()
	}
	var snapshot authenticationSnapshot
	var principalID string
	var state string
	var version uint64
	err := repository.pool.QueryRow(
		ctx,
		`SELECT p.id, p.lifecycle_state, p.updated_at, pc.password_hash, pc.credential_version
		 FROM omnexa_identity.principals p
		 JOIN omnexa_identity.users u ON u.principal_id = p.id
		 JOIN omnexa_identity.password_credentials pc ON pc.principal_id = p.id
		 WHERE p.id = $1 AND p.principal_type = 'human_user'`,
		string(userID),
	).Scan(&principalID, &state, &snapshot.userUpdatedAt, &snapshot.passwordHash, &version)
	if errors.Is(err, pgx.ErrNoRows) {
		return authenticationSnapshot{}, credentialNotFoundFailure()
	}
	if err != nil {
		return authenticationSnapshot{}, authenticationRepositoryFailure(err)
	}
	snapshot.principalID = UserID(principalID)
	snapshot.state = LifecycleState(state)
	snapshot.credentialVersion = version
	if !snapshot.principalID.Valid() || !snapshot.state.Valid() || snapshot.userUpdatedAt.IsZero() || snapshot.passwordHash == "" || version == 0 {
		return authenticationSnapshot{}, invalidStoredUserFailure()
	}
	return snapshot, nil
}

func (repository *PostgresAuthenticationRepository) ChangePassword(
	ctx context.Context,
	userID UserID,
	expectedVersion uint64,
	passwordHash string,
	changedAt time.Time,
) (uint64, error) {
	if repository == nil || repository.pool == nil || !userID.Valid() || expectedVersion == 0 || passwordHash == "" || changedAt.IsZero() {
		return 0, repositoryInvalidFailure()
	}
	tx, err := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, authenticationRepositoryFailure(err)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	var version uint64
	err = tx.QueryRow(
		ctx,
		`UPDATE omnexa_identity.password_credentials pc
		 SET password_hash = $3,
		     credential_version = pc.credential_version + 1,
		     updated_at = $4
		 FROM omnexa_identity.principals p
		 WHERE pc.principal_id = $1
		   AND p.id = pc.principal_id
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pc.credential_version = $2
		   AND pc.updated_at <= $4
		 RETURNING pc.credential_version`,
		string(userID),
		expectedVersion,
		passwordHash,
		changedAt.UTC(),
	).Scan(&version)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, credentialConflictFailure()
	}
	if err != nil {
		return 0, authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`UPDATE omnexa_identity.sessions
		 SET revoked_at = COALESCE(revoked_at, $2)
		 WHERE principal_id = $1`,
		string(userID),
		changedAt.UTC(),
	); err != nil {
		return 0, authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`UPDATE omnexa_identity.access_credentials a
		 SET revoked_at = COALESCE(a.revoked_at, $2)
		 FROM omnexa_identity.sessions s
		 WHERE a.session_id = s.id AND s.principal_id = $1`,
		string(userID),
		changedAt.UTC(),
	); err != nil {
		return 0, authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`UPDATE omnexa_identity.refresh_credentials r
		 SET revoked_at = COALESCE(r.revoked_at, $2)
		 FROM omnexa_identity.sessions s
		 WHERE r.session_id = s.id AND s.principal_id = $1`,
		string(userID),
		changedAt.UTC(),
	); err != nil {
		return 0, authenticationRepositoryFailure(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, authenticationRepositoryFailure(err)
	}
	return version, nil
}

func (repository *PostgresAuthenticationRepository) CreateSession(
	ctx context.Context,
	record sessionRecord,
	accessDigest [sha256.Size]byte,
	accessExpiresAt time.Time,
	refreshDigest [sha256.Size]byte,
	refreshExpiresAt time.Time,
) error {
	if repository == nil || repository.pool == nil || !validSessionRecord(record) || accessExpiresAt.IsZero() || refreshExpiresAt.IsZero() {
		return repositoryInvalidFailure()
	}
	tx, err := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return authenticationRepositoryFailure(err)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	commandTag, err := tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.sessions (
			id, principal_id, credential_version, device_label,
			tenant_context_hint, organization_context_hint,
			created_at, refreshed_at, expires_at
		 )
		 SELECT $1, p.id, $3, $4, NULLIF($5, '')::uuid, NULLIF($6, '')::uuid, $7, $7, $8
		 FROM omnexa_identity.principals p
		 JOIN omnexa_identity.password_credentials pc ON pc.principal_id = p.id
		 WHERE p.id = $2
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pc.credential_version = $3
		   AND p.updated_at <= $7`,
		string(record.session.id),
		string(record.session.principalID),
		record.credentialVersion,
		record.session.deviceLabel,
		record.session.context.tenantID,
		record.session.context.organizationID,
		record.session.createdAt.UTC(),
		record.session.expiresAt.UTC(),
	)
	if err != nil {
		return authenticationRepositoryFailure(err)
	}
	if commandTag.RowsAffected() != 1 {
		return sessionConflictFailure()
	}
	if _, err = tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.access_credentials
			(session_id, secret_digest, created_at, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		string(record.session.id),
		accessDigest[:],
		record.session.createdAt.UTC(),
		accessExpiresAt.UTC(),
	); err != nil {
		return authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.refresh_credentials
			(session_id, secret_digest, created_at, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		string(record.session.id),
		refreshDigest[:],
		record.session.createdAt.UTC(),
		refreshExpiresAt.UTC(),
	); err != nil {
		return authenticationRepositoryFailure(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return authenticationRepositoryFailure(err)
	}
	return nil
}

func (repository *PostgresAuthenticationRepository) AccessSession(
	ctx context.Context,
	digest [sha256.Size]byte,
	at time.Time,
) (sessionRecord, error) {
	if repository == nil || repository.pool == nil || at.IsZero() {
		return sessionRecord{}, sessionFailure()
	}
	return repository.sessionByCredential(
		ctx,
		`SELECT s.id, s.principal_id, s.credential_version, s.device_label,
		        COALESCE(s.tenant_context_hint::text, ''), COALESCE(s.organization_context_hint::text, ''),
		        s.created_at, s.refreshed_at, s.expires_at,
		        COALESCE(s.revoked_at, '0001-01-01T00:00:00Z'::timestamptz),
		        (p.lifecycle_state <> 'active' OR pc.credential_version <> s.credential_version OR s.created_at < p.updated_at)
		 FROM omnexa_identity.access_credentials c
		 JOIN omnexa_identity.sessions s ON s.id = c.session_id
		 JOIN omnexa_identity.principals p ON p.id = s.principal_id
		 JOIN omnexa_identity.password_credentials pc ON pc.principal_id = s.principal_id
		 WHERE c.secret_digest = $1
		   AND c.revoked_at IS NULL
		   AND c.expires_at > $2
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $2
		   AND p.lifecycle_state = 'active'
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at`,
		digest,
		at,
	)
}

func (repository *PostgresAuthenticationRepository) RefreshSession(
	ctx context.Context,
	digest [sha256.Size]byte,
	at time.Time,
) (sessionRecord, error) {
	if repository == nil || repository.pool == nil || at.IsZero() {
		return sessionRecord{}, sessionFailure()
	}
	return repository.sessionByCredential(
		ctx,
		`SELECT s.id, s.principal_id, s.credential_version, s.device_label,
		        COALESCE(s.tenant_context_hint::text, ''), COALESCE(s.organization_context_hint::text, ''),
		        s.created_at, s.refreshed_at, s.expires_at,
		        COALESCE(s.revoked_at, '0001-01-01T00:00:00Z'::timestamptz),
		        (p.lifecycle_state <> 'active' OR pc.credential_version <> s.credential_version OR s.created_at < p.updated_at)
		 FROM omnexa_identity.refresh_credentials c
		 JOIN omnexa_identity.sessions s ON s.id = c.session_id
		 JOIN omnexa_identity.principals p ON p.id = s.principal_id
		 JOIN omnexa_identity.password_credentials pc ON pc.principal_id = s.principal_id
		 WHERE c.secret_digest = $1
		   AND c.revoked_at IS NULL
		   AND c.consumed_at IS NULL
		   AND c.expires_at > $2
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $2
		   AND p.lifecycle_state = 'active'
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at`,
		digest,
		at,
	)
}

func (repository *PostgresAuthenticationRepository) sessionByCredential(
	ctx context.Context,
	query string,
	digest [sha256.Size]byte,
	at time.Time,
) (sessionRecord, error) {
	var record sessionRecord
	var sessionID, principalID, tenantHint, organizationHint string
	var revokedAt time.Time
	err := repository.pool.QueryRow(ctx, query, digest[:], at.UTC()).Scan(
		&sessionID,
		&principalID,
		&record.credentialVersion,
		&record.session.deviceLabel,
		&tenantHint,
		&organizationHint,
		&record.session.createdAt,
		&record.session.refreshedAt,
		&record.session.expiresAt,
		&revokedAt,
		&record.session.invalidated,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return sessionRecord{}, sessionFailure()
	}
	if err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	contextValue, contextErr := NewSessionContext(tenantHint, organizationHint)
	if contextErr != nil {
		return sessionRecord{}, sessionFailure()
	}
	record.session.id = SessionID(sessionID)
	record.session.principalID = UserID(principalID)
	record.session.context = contextValue
	if revokedAt.Year() > 1 {
		record.session.revokedAt = revokedAt.UTC()
	}
	if !validSessionRecord(record) {
		return sessionRecord{}, sessionFailure()
	}
	return record, nil
}

func (repository *PostgresAuthenticationRepository) RotateRefresh(
	ctx context.Context,
	sessionID SessionID,
	expectedRefreshDigest [sha256.Size]byte,
	accessDigest [sha256.Size]byte,
	accessExpiresAt time.Time,
	refreshDigest [sha256.Size]byte,
	refreshExpiresAt time.Time,
	rotatedAt time.Time,
) (sessionRecord, error) {
	if repository == nil || repository.pool == nil || !sessionID.Valid() || accessExpiresAt.IsZero() || refreshExpiresAt.IsZero() || rotatedAt.IsZero() {
		return sessionRecord{}, sessionFailure()
	}
	tx, err := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	var record sessionRecord
	var storedSessionID, principalID, tenantHint, organizationHint string
	var revokedAt time.Time
	err = tx.QueryRow(
		ctx,
		`UPDATE omnexa_identity.refresh_credentials r
		 SET consumed_at = $3
		 FROM omnexa_identity.sessions s,
		      omnexa_identity.principals p,
		      omnexa_identity.password_credentials pc
		 WHERE r.session_id = $1
		   AND r.secret_digest = $2
		   AND r.session_id = s.id
		   AND p.id = s.principal_id
		   AND pc.principal_id = s.principal_id
		   AND r.revoked_at IS NULL
		   AND r.consumed_at IS NULL
		   AND r.expires_at > $3
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $3
		   AND p.lifecycle_state = 'active'
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at
		 RETURNING s.id, s.principal_id, s.credential_version, s.device_label,
		           COALESCE(s.tenant_context_hint::text, ''), COALESCE(s.organization_context_hint::text, ''),
		           s.created_at, s.refreshed_at, s.expires_at,
		           COALESCE(s.revoked_at, '0001-01-01T00:00:00Z'::timestamptz)`,
		string(sessionID),
		expectedRefreshDigest[:],
		rotatedAt.UTC(),
	).Scan(
		&storedSessionID,
		&principalID,
		&record.credentialVersion,
		&record.session.deviceLabel,
		&tenantHint,
		&organizationHint,
		&record.session.createdAt,
		&record.session.refreshedAt,
		&record.session.expiresAt,
		&revokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return sessionRecord{}, sessionFailure()
	}
	if err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`UPDATE omnexa_identity.access_credentials
		 SET revoked_at = COALESCE(revoked_at, $2)
		 WHERE session_id = $1`,
		string(sessionID),
		rotatedAt.UTC(),
	); err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.access_credentials
			(session_id, secret_digest, created_at, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		string(sessionID),
		accessDigest[:],
		rotatedAt.UTC(),
		accessExpiresAt.UTC(),
	); err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.refresh_credentials
			(session_id, secret_digest, created_at, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		string(sessionID),
		refreshDigest[:],
		rotatedAt.UTC(),
		refreshExpiresAt.UTC(),
	); err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`UPDATE omnexa_identity.sessions SET refreshed_at = $2 WHERE id = $1`,
		string(sessionID),
		rotatedAt.UTC(),
	); err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}

	contextValue, contextErr := NewSessionContext(tenantHint, organizationHint)
	if contextErr != nil {
		return sessionRecord{}, sessionFailure()
	}
	record.session.id = SessionID(storedSessionID)
	record.session.principalID = UserID(principalID)
	record.session.context = contextValue
	record.session.refreshedAt = rotatedAt.UTC()
	if revokedAt.Year() > 1 {
		record.session.revokedAt = revokedAt.UTC()
	}
	if !validSessionRecord(record) {
		return sessionRecord{}, sessionFailure()
	}
	return record, nil
}

func (repository *PostgresAuthenticationRepository) RevokeSession(
	ctx context.Context,
	userID UserID,
	sessionID SessionID,
	revokedAt time.Time,
) (sessionRecord, error) {
	if repository == nil || repository.pool == nil || !userID.Valid() || !sessionID.Valid() || revokedAt.IsZero() {
		return sessionRecord{}, sessionFailure()
	}
	tx, err := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	var record sessionRecord
	var storedSessionID, principalID, tenantHint, organizationHint string
	var storedRevokedAt time.Time
	err = tx.QueryRow(
		ctx,
		`UPDATE omnexa_identity.sessions
		 SET revoked_at = COALESCE(revoked_at, $3)
		 WHERE id = $1 AND principal_id = $2
		 RETURNING id, principal_id, credential_version, device_label,
		           COALESCE(tenant_context_hint::text, ''), COALESCE(organization_context_hint::text, ''),
		           created_at, refreshed_at, expires_at, revoked_at`,
		string(sessionID),
		string(userID),
		revokedAt.UTC(),
	).Scan(
		&storedSessionID,
		&principalID,
		&record.credentialVersion,
		&record.session.deviceLabel,
		&tenantHint,
		&organizationHint,
		&record.session.createdAt,
		&record.session.refreshedAt,
		&record.session.expiresAt,
		&storedRevokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return sessionRecord{}, sessionFailure()
	}
	if err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`UPDATE omnexa_identity.access_credentials
		 SET revoked_at = COALESCE(revoked_at, $2)
		 WHERE session_id = $1`,
		string(sessionID),
		storedRevokedAt.UTC(),
	); err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	if _, err = tx.Exec(
		ctx,
		`UPDATE omnexa_identity.refresh_credentials
		 SET revoked_at = COALESCE(revoked_at, $2)
		 WHERE session_id = $1`,
		string(sessionID),
		storedRevokedAt.UTC(),
	); err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return sessionRecord{}, authenticationRepositoryFailure(err)
	}

	contextValue, contextErr := NewSessionContext(tenantHint, organizationHint)
	if contextErr != nil {
		return sessionRecord{}, sessionFailure()
	}
	record.session.id = SessionID(storedSessionID)
	record.session.principalID = UserID(principalID)
	record.session.context = contextValue
	record.session.revokedAt = storedRevokedAt.UTC()
	if !validSessionRecord(record) {
		return sessionRecord{}, sessionFailure()
	}
	return record, nil
}

func (repository *PostgresAuthenticationRepository) ListSessions(
	ctx context.Context,
	userID UserID,
) ([]Session, error) {
	if repository == nil || repository.pool == nil || !userID.Valid() {
		return nil, repositoryInvalidFailure()
	}
	rows, err := repository.pool.Query(
		ctx,
		`SELECT s.id, s.principal_id, s.device_label,
		        COALESCE(s.tenant_context_hint::text, ''), COALESCE(s.organization_context_hint::text, ''),
		        s.created_at, s.refreshed_at, s.expires_at,
		        COALESCE(s.revoked_at, '0001-01-01T00:00:00Z'::timestamptz),
		        (p.lifecycle_state <> 'active' OR pc.credential_version <> s.credential_version OR s.created_at < p.updated_at)
		 FROM omnexa_identity.sessions s
		 JOIN omnexa_identity.principals p ON p.id = s.principal_id
		 JOIN omnexa_identity.password_credentials pc ON pc.principal_id = s.principal_id
		 WHERE s.principal_id = $1
		 ORDER BY s.created_at DESC, s.id DESC`,
		string(userID),
	)
	if err != nil {
		return nil, authenticationRepositoryFailure(err)
	}
	defer rows.Close()

	sessions := make([]Session, 0)
	for rows.Next() {
		var session Session
		var sessionID, principalID, tenantHint, organizationHint string
		var revokedAt time.Time
		if err = rows.Scan(
			&sessionID,
			&principalID,
			&session.deviceLabel,
			&tenantHint,
			&organizationHint,
			&session.createdAt,
			&session.refreshedAt,
			&session.expiresAt,
			&revokedAt,
			&session.invalidated,
		); err != nil {
			return nil, authenticationRepositoryFailure(err)
		}
		contextValue, contextErr := NewSessionContext(tenantHint, organizationHint)
		if contextErr != nil {
			return nil, sessionFailure()
		}
		session.id = SessionID(sessionID)
		session.principalID = UserID(principalID)
		session.context = contextValue
		if revokedAt.Year() > 1 {
			session.revokedAt = revokedAt.UTC()
		}
		if !session.id.Valid() || !session.principalID.Valid() || session.createdAt.IsZero() || session.expiresAt.IsZero() {
			return nil, sessionFailure()
		}
		sessions = append(sessions, session)
	}
	if err = rows.Err(); err != nil {
		return nil, authenticationRepositoryFailure(err)
	}
	return sessions, nil
}

func validSessionRecord(record sessionRecord) bool {
	return record.session.id.Valid() &&
		record.session.principalID.Valid() &&
		record.credentialVersion > 0 &&
		record.session.context.valid() &&
		validDeviceLabel(record.session.deviceLabel) &&
		!record.session.createdAt.IsZero() &&
		!record.session.refreshedAt.IsZero() &&
		!record.session.expiresAt.IsZero() &&
		!record.session.expiresAt.Before(record.session.refreshedAt) &&
		!record.session.refreshedAt.Before(record.session.createdAt)
}

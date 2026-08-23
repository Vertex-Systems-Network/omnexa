package identity

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStrongAuthenticationRepository persists only P02.07 kernel.identity
// strong-authentication state. Raw challenge/recovery material and authenticator
// private keys never cross this repository boundary.
type PostgresStrongAuthenticationRepository struct{ pool *pgxpool.Pool }

// NewPostgresStrongAuthenticationRepository returns the owner-bounded P02.07 repository.
func NewPostgresStrongAuthenticationRepository(pool *pgxpool.Pool) (*PostgresStrongAuthenticationRepository, error) {
	if pool == nil {
		return nil, repositoryInvalidFailure()
	}
	return &PostgresStrongAuthenticationRepository{pool: pool}, nil
}

func (repository *PostgresStrongAuthenticationRepository) CreatePasskeyEnrollment(
	ctx context.Context,
	factor StrongFactor,
	challenge strongChallengeRecord,
) error {
	if repository == nil || repository.pool == nil || !factor.valid() || factor.state != StrongFactorPending ||
		!validStrongChallenge(challenge) || challenge.factor != factor.id || challenge.principal != factor.principal ||
		challenge.purpose != challengePurposeEnrollment {
		return repositoryInvalidFailure()
	}
	tx, beginErr := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if beginErr != nil {
		return authenticationRepositoryFailure(beginErr)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	factorTag, factorErr := tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.mfa_factors
			(id, principal_id, factor_type, label, lifecycle_state, created_at)
		 SELECT $1, s.principal_id, $3, $4, $5, $6
		 FROM omnexa_identity.sessions s
		 JOIN omnexa_identity.principals p ON p.id = s.principal_id
		 JOIN omnexa_identity.password_credentials pc ON pc.principal_id = s.principal_id
		 WHERE s.id = $2
		   AND s.principal_id = $7
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $6
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at`,
		string(factor.id),
		string(challenge.session),
		string(factor.factorType),
		factor.label,
		string(factor.state),
		factor.createdAt.UTC(),
		string(factor.principal),
	)
	if factorErr != nil {
		return authenticationRepositoryFailure(factorErr)
	}
	if factorTag.RowsAffected() != 1 {
		return strongAuthInvalidFailure()
	}
	if _, challengeErr := tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.authentication_challenges
			(id, principal_id, session_id, factor_id, purpose, secret_digest, created_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(challenge.id),
		string(challenge.principal),
		string(challenge.session),
		string(challenge.factor),
		challenge.purpose,
		challenge.digest[:],
		challenge.createdAt.UTC(),
		challenge.expiresAt.UTC(),
	); challengeErr != nil {
		return authenticationRepositoryFailure(challengeErr)
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		return authenticationRepositoryFailure(commitErr)
	}
	return nil
}

func (repository *PostgresStrongAuthenticationRepository) ConsumeEnrollmentChallenge(
	ctx context.Context,
	principalID UserID,
	sessionID SessionID,
	challengeID ChallengeID,
	digest [sha256.Size]byte,
	at time.Time,
) (StrongFactor, error) {
	if repository == nil || repository.pool == nil || !principalID.Valid() || !sessionID.Valid() || !challengeID.Valid() || at.IsZero() {
		return StrongFactor{}, strongAuthChallengeFailure()
	}
	var factorID, storedPrincipal, factorType, label, state string
	var createdAt, verifiedAt, revokedAt time.Time
	queryErr := repository.pool.QueryRow(
		ctx,
		`UPDATE omnexa_identity.authentication_challenges c
		 SET consumed_at = $6
		 FROM omnexa_identity.mfa_factors f,
		      omnexa_identity.sessions s,
		      omnexa_identity.principals p,
		      omnexa_identity.password_credentials pc
		 WHERE c.id = $1
		   AND c.principal_id = $2
		   AND c.session_id = $3
		   AND c.secret_digest = $4
		   AND c.purpose = $5
		   AND c.consumed_at IS NULL
		   AND c.revoked_at IS NULL
		   AND c.expires_at > $6
		   AND f.id = c.factor_id
		   AND f.principal_id = c.principal_id
		   AND f.lifecycle_state = 'pending'
		   AND s.id = c.session_id
		   AND s.principal_id = c.principal_id
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $6
		   AND p.id = s.principal_id
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pc.principal_id = s.principal_id
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at
		 RETURNING f.id, f.principal_id, f.factor_type, f.label, f.lifecycle_state,
		           f.created_at, COALESCE(f.verified_at, '0001-01-01T00:00:00Z'::timestamptz),
		           COALESCE(f.revoked_at, '0001-01-01T00:00:00Z'::timestamptz)`,
		string(challengeID),
		string(principalID),
		string(sessionID),
		digest[:],
		challengePurposeEnrollment,
		at.UTC(),
	).Scan(&factorID, &storedPrincipal, &factorType, &label, &state, &createdAt, &verifiedAt, &revokedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return StrongFactor{}, strongAuthChallengeFailure()
	}
	if queryErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(queryErr)
	}
	return storedStrongFactor(factorID, storedPrincipal, factorType, label, state, createdAt, verifiedAt, revokedAt)
}

func (repository *PostgresStrongAuthenticationRepository) ActivatePasskey(
	ctx context.Context,
	principalID UserID,
	sessionID SessionID,
	factorID FactorID,
	credential VerifiedPasskeyCredential,
	at time.Time,
) (StrongFactor, error) {
	if repository == nil || repository.pool == nil || !principalID.Valid() || !sessionID.Valid() || !factorID.Valid() || !credential.valid() || at.IsZero() {
		return StrongFactor{}, strongAuthInvalidFailure()
	}
	tx, beginErr := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if beginErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(beginErr)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	var storedID, storedPrincipal, factorType, label, state string
	var createdAt, verifiedAt, revokedAt time.Time
	queryErr := tx.QueryRow(
		ctx,
		`UPDATE omnexa_identity.mfa_factors f
		 SET lifecycle_state = 'active', verified_at = $4
		 FROM omnexa_identity.sessions s,
		      omnexa_identity.principals p,
		      omnexa_identity.password_credentials pc
		 WHERE f.id = $1
		   AND f.principal_id = $2
		   AND f.lifecycle_state = 'pending'
		   AND s.id = $3
		   AND s.principal_id = f.principal_id
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $4
		   AND p.id = s.principal_id
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pc.principal_id = s.principal_id
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at
		 RETURNING f.id, f.principal_id, f.factor_type, f.label, f.lifecycle_state,
		           f.created_at, f.verified_at,
		           COALESCE(f.revoked_at, '0001-01-01T00:00:00Z'::timestamptz)`,
		string(factorID), string(principalID), string(sessionID), at.UTC(),
	).Scan(&storedID, &storedPrincipal, &factorType, &label, &state, &createdAt, &verifiedAt, &revokedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return StrongFactor{}, strongFactorConflictFailure()
	}
	if queryErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(queryErr)
	}
	factor, factorErr := storedStrongFactor(storedID, storedPrincipal, factorType, label, state, createdAt, verifiedAt, revokedAt)
	if factorErr != nil {
		return StrongFactor{}, factorErr
	}
	if _, insertErr := tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.passkey_credentials
			(factor_id, credential_id, public_key, counter_supported, sign_count, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $6)`,
		string(factorID), credential.credentialID, credential.publicKey,
		credential.counterSupported, credential.signCount, at.UTC(),
	); insertErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(insertErr)
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(commitErr)
	}
	return factor, nil
}

func (repository *PostgresStrongAuthenticationRepository) CreateAssertionChallenge(
	ctx context.Context,
	principalID UserID,
	sessionID SessionID,
	factorID FactorID,
	challenge strongChallengeRecord,
	at time.Time,
) (StrongFactor, error) {
	if repository == nil || repository.pool == nil || !principalID.Valid() || !sessionID.Valid() || !factorID.Valid() ||
		!validStrongChallenge(challenge) || challenge.principal != principalID || challenge.session != sessionID ||
		challenge.factor != factorID || challenge.purpose != challengePurposeAssertion || at.IsZero() {
		return StrongFactor{}, repositoryInvalidFailure()
	}
	tx, beginErr := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if beginErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(beginErr)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	factor, factorErr := loadCurrentActivePasskeyFactor(ctx, tx, principalID, sessionID, factorID, at.UTC())
	if factorErr != nil {
		return StrongFactor{}, factorErr
	}
	if _, insertErr := tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.authentication_challenges
			(id, principal_id, session_id, factor_id, purpose, secret_digest, created_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(challenge.id), string(principalID), string(sessionID), string(factorID),
		challenge.purpose, challenge.digest[:], challenge.createdAt.UTC(), challenge.expiresAt.UTC(),
	); insertErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(insertErr)
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(commitErr)
	}
	return factor, nil
}

func (repository *PostgresStrongAuthenticationRepository) ConsumeAssertionChallenge(
	ctx context.Context,
	principalID UserID,
	sessionID SessionID,
	challengeID ChallengeID,
	factorID FactorID,
	digest [sha256.Size]byte,
	at time.Time,
) (StrongFactor, passkeyCredentialRecord, error) {
	if repository == nil || repository.pool == nil || !principalID.Valid() || !sessionID.Valid() || !challengeID.Valid() || !factorID.Valid() || at.IsZero() {
		return StrongFactor{}, passkeyCredentialRecord{}, strongAuthChallengeFailure()
	}
	tx, beginErr := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if beginErr != nil {
		return StrongFactor{}, passkeyCredentialRecord{}, authenticationRepositoryFailure(beginErr)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	consumeTag, consumeErr := tx.Exec(
		ctx,
		`UPDATE omnexa_identity.authentication_challenges c
		 SET consumed_at = $7
		 FROM omnexa_identity.sessions s,
		      omnexa_identity.principals p,
		      omnexa_identity.password_credentials pc,
		      omnexa_identity.mfa_factors f
		 WHERE c.id = $1
		   AND c.principal_id = $2
		   AND c.session_id = $3
		   AND c.factor_id = $4
		   AND c.secret_digest = $5
		   AND c.purpose = $6
		   AND c.consumed_at IS NULL
		   AND c.revoked_at IS NULL
		   AND c.expires_at > $7
		   AND s.id = c.session_id
		   AND s.principal_id = c.principal_id
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $7
		   AND p.id = s.principal_id
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pc.principal_id = s.principal_id
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at
		   AND f.id = c.factor_id
		   AND f.principal_id = c.principal_id
		   AND f.lifecycle_state = 'active'`,
		string(challengeID), string(principalID), string(sessionID), string(factorID),
		digest[:], challengePurposeAssertion, at.UTC(),
	)
	if consumeErr != nil {
		return StrongFactor{}, passkeyCredentialRecord{}, authenticationRepositoryFailure(consumeErr)
	}
	if consumeTag.RowsAffected() != 1 {
		return StrongFactor{}, passkeyCredentialRecord{}, strongAuthChallengeFailure()
	}
	factor, factorErr := loadCurrentActivePasskeyFactor(ctx, tx, principalID, sessionID, factorID, at.UTC())
	if factorErr != nil {
		return StrongFactor{}, passkeyCredentialRecord{}, factorErr
	}
	var credentialID, publicKey []byte
	var counterSupported bool
	var signCount uint32
	credentialErr := tx.QueryRow(
		ctx,
		`SELECT credential_id, public_key, counter_supported, sign_count
		 FROM omnexa_identity.passkey_credentials
		 WHERE factor_id = $1`,
		string(factorID),
	).Scan(&credentialID, &publicKey, &counterSupported, &signCount)
	if errors.Is(credentialErr, pgx.ErrNoRows) {
		return StrongFactor{}, passkeyCredentialRecord{}, strongFactorNotFoundFailure()
	}
	if credentialErr != nil {
		return StrongFactor{}, passkeyCredentialRecord{}, authenticationRepositoryFailure(credentialErr)
	}
	record := passkeyCredentialRecord{
		factor:           factorID,
		credentialID:     append([]byte(nil), credentialID...),
		publicKey:        append([]byte(nil), publicKey...),
		counterSupported: counterSupported,
		signCount:        signCount,
	}
	if !record.valid() {
		return StrongFactor{}, passkeyCredentialRecord{}, invalidStoredUserFailure()
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		return StrongFactor{}, passkeyCredentialRecord{}, authenticationRepositoryFailure(commitErr)
	}
	return factor, record, nil
}

func (repository *PostgresStrongAuthenticationRepository) AdvancePasskeyCounter(
	ctx context.Context,
	principalID UserID,
	sessionID SessionID,
	factorID FactorID,
	expected uint32,
	next uint32,
	counterSupported bool,
	at time.Time,
) error {
	if repository == nil || repository.pool == nil || !principalID.Valid() || !sessionID.Valid() || !factorID.Valid() || at.IsZero() ||
		(counterSupported && next <= expected) || (!counterSupported && (expected != 0 || next != 0)) {
		return passkeyVerificationFailure()
	}
	commandTag, updateErr := repository.pool.Exec(
		ctx,
		`UPDATE omnexa_identity.passkey_credentials pc
		 SET sign_count = $5, updated_at = $7
		 FROM omnexa_identity.mfa_factors f,
		      omnexa_identity.sessions s,
		      omnexa_identity.principals p,
		      omnexa_identity.password_credentials pwd
		 WHERE pc.factor_id = $1
		   AND f.id = pc.factor_id
		   AND f.principal_id = $2
		   AND f.lifecycle_state = 'active'
		   AND s.id = $3
		   AND s.principal_id = f.principal_id
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $7
		   AND p.id = s.principal_id
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pwd.principal_id = s.principal_id
		   AND pwd.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at
		   AND pc.sign_count = $4
		   AND pc.counter_supported = $6`,
		string(factorID), string(principalID), string(sessionID), expected, next, counterSupported, at.UTC(),
	)
	if updateErr != nil {
		return authenticationRepositoryFailure(updateErr)
	}
	if commandTag.RowsAffected() != 1 {
		return passkeyVerificationFailure()
	}
	return nil
}

func (repository *PostgresStrongAuthenticationRepository) ReplaceRecoveryCodes(
	ctx context.Context,
	principalID UserID,
	sessionID SessionID,
	setID RecoverySetID,
	digests [][sha256.Size]byte,
	at time.Time,
) error {
	if repository == nil || repository.pool == nil || !principalID.Valid() || !sessionID.Valid() || !setID.Valid() ||
		len(digests) < minRecoveryCodes || len(digests) > maxRecoveryCodes || at.IsZero() {
		return repositoryInvalidFailure()
	}
	tx, beginErr := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if beginErr != nil {
		return authenticationRepositoryFailure(beginErr)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	var currentSession int
	validationErr := tx.QueryRow(
		ctx,
		`SELECT 1
		 FROM omnexa_identity.sessions s
		 JOIN omnexa_identity.principals p ON p.id = s.principal_id
		 JOIN omnexa_identity.password_credentials pc ON pc.principal_id = s.principal_id
		 WHERE s.id = $1
		   AND s.principal_id = $2
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $3
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at
		 FOR UPDATE OF s`,
		string(sessionID), string(principalID), at.UTC(),
	).Scan(&currentSession)
	if errors.Is(validationErr, pgx.ErrNoRows) {
		return strongAuthInvalidFailure()
	}
	if validationErr != nil {
		return authenticationRepositoryFailure(validationErr)
	}
	if _, revokeErr := tx.Exec(
		ctx,
		`UPDATE omnexa_identity.recovery_code_sets
		 SET revoked_at = COALESCE(revoked_at, $2)
		 WHERE principal_id = $1 AND revoked_at IS NULL`,
		string(principalID), at.UTC(),
	); revokeErr != nil {
		return authenticationRepositoryFailure(revokeErr)
	}
	if _, insertErr := tx.Exec(
		ctx,
		`INSERT INTO omnexa_identity.recovery_code_sets (id, principal_id, created_at)
		 VALUES ($1, $2, $3)`,
		string(setID), string(principalID), at.UTC(),
	); insertErr != nil {
		return authenticationRepositoryFailure(insertErr)
	}
	for index, digest := range digests {
		if _, insertErr := tx.Exec(
			ctx,
			`INSERT INTO omnexa_identity.recovery_codes
				(set_id, code_index, secret_digest, created_at)
			 VALUES ($1, $2, $3, $4)`,
			string(setID), index, digest[:], at.UTC(),
		); insertErr != nil {
			return authenticationRepositoryFailure(insertErr)
		}
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		return authenticationRepositoryFailure(commitErr)
	}
	return nil
}

func (repository *PostgresStrongAuthenticationRepository) ConsumeRecoveryCode(
	ctx context.Context,
	principalID UserID,
	sessionID SessionID,
	digest [sha256.Size]byte,
	at time.Time,
) error {
	if repository == nil || repository.pool == nil || !principalID.Valid() || !sessionID.Valid() || at.IsZero() {
		return recoveryCodeFailure()
	}
	commandTag, consumeErr := repository.pool.Exec(
		ctx,
		`UPDATE omnexa_identity.recovery_codes rc
		 SET consumed_at = $4
		 FROM omnexa_identity.recovery_code_sets rs,
		      omnexa_identity.sessions s,
		      omnexa_identity.principals p,
		      omnexa_identity.password_credentials pc
		 WHERE rc.secret_digest = $1
		   AND rc.consumed_at IS NULL
		   AND rs.id = rc.set_id
		   AND rs.principal_id = $2
		   AND rs.revoked_at IS NULL
		   AND s.id = $3
		   AND s.principal_id = rs.principal_id
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $4
		   AND p.id = s.principal_id
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pc.principal_id = s.principal_id
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at`,
		digest[:], string(principalID), string(sessionID), at.UTC(),
	)
	if consumeErr != nil {
		return authenticationRepositoryFailure(consumeErr)
	}
	if commandTag.RowsAffected() != 1 {
		return recoveryCodeFailure()
	}
	return nil
}

func (repository *PostgresStrongAuthenticationRepository) RevokeFactor(
	ctx context.Context,
	principalID UserID,
	sessionID SessionID,
	factorID FactorID,
	invalidateSessions bool,
	at time.Time,
) (StrongFactor, error) {
	if repository == nil || repository.pool == nil || !principalID.Valid() || !sessionID.Valid() || !factorID.Valid() || at.IsZero() {
		return StrongFactor{}, strongAuthInvalidFailure()
	}
	tx, beginErr := repository.pool.BeginTx(ctx, pgx.TxOptions{})
	if beginErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(beginErr)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	var storedID, storedPrincipal, factorType, label, state string
	var createdAt, verifiedAt, revokedAt time.Time
	queryErr := tx.QueryRow(
		ctx,
		`UPDATE omnexa_identity.mfa_factors f
		 SET lifecycle_state = 'revoked', revoked_at = COALESCE(f.revoked_at, $4)
		 FROM omnexa_identity.sessions s,
		      omnexa_identity.principals p,
		      omnexa_identity.password_credentials pc
		 WHERE f.id = $1
		   AND f.principal_id = $2
		   AND f.lifecycle_state = 'active'
		   AND s.id = $3
		   AND s.principal_id = f.principal_id
		   AND s.revoked_at IS NULL
		   AND s.expires_at > $4
		   AND p.id = s.principal_id
		   AND p.principal_type = 'human_user'
		   AND p.lifecycle_state = 'active'
		   AND pc.principal_id = s.principal_id
		   AND pc.credential_version = s.credential_version
		   AND s.created_at >= p.updated_at
		 RETURNING f.id, f.principal_id, f.factor_type, f.label, f.lifecycle_state,
		           f.created_at, f.verified_at, f.revoked_at`,
		string(factorID), string(principalID), string(sessionID), at.UTC(),
	).Scan(&storedID, &storedPrincipal, &factorType, &label, &state, &createdAt, &verifiedAt, &revokedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return StrongFactor{}, strongFactorConflictFailure()
	}
	if queryErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(queryErr)
	}
	factor, factorErr := storedStrongFactor(storedID, storedPrincipal, factorType, label, state, createdAt, verifiedAt, revokedAt)
	if factorErr != nil {
		return StrongFactor{}, factorErr
	}
	if _, revokeErr := tx.Exec(
		ctx,
		`UPDATE omnexa_identity.authentication_challenges
		 SET revoked_at = COALESCE(revoked_at, $2)
		 WHERE factor_id = $1 AND consumed_at IS NULL AND revoked_at IS NULL`,
		string(factorID), at.UTC(),
	); revokeErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(revokeErr)
	}
	if invalidateSessions {
		if invalidateErr := revokeHumanSessions(ctx, tx, principalID, at.UTC()); invalidateErr != nil {
			return StrongFactor{}, invalidateErr
		}
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(commitErr)
	}
	return factor, nil
}

func (repository *PostgresStrongAuthenticationRepository) ListFactors(
	ctx context.Context,
	principalID UserID,
	sessionID SessionID,
	at time.Time,
) ([]StrongFactor, error) {
	if repository == nil || repository.pool == nil || !principalID.Valid() || !sessionID.Valid() || at.IsZero() {
		return nil, strongAuthInvalidFailure()
	}
	rows, queryErr := repository.pool.Query(
		ctx,
		`SELECT f.id, f.principal_id, f.factor_type, f.label, f.lifecycle_state, f.created_at,
		        COALESCE(f.verified_at, '0001-01-01T00:00:00Z'::timestamptz),
		        COALESCE(f.revoked_at, '0001-01-01T00:00:00Z'::timestamptz)
		 FROM omnexa_identity.mfa_factors f
		 WHERE f.principal_id = $1
		   AND EXISTS (
		       SELECT 1
		       FROM omnexa_identity.sessions s
		       JOIN omnexa_identity.principals p ON p.id = s.principal_id
		       JOIN omnexa_identity.password_credentials pc ON pc.principal_id = s.principal_id
		       WHERE s.id = $2
		         AND s.principal_id = f.principal_id
		         AND s.revoked_at IS NULL
		         AND s.expires_at > $3
		         AND p.principal_type = 'human_user'
		         AND p.lifecycle_state = 'active'
		         AND pc.credential_version = s.credential_version
		         AND s.created_at >= p.updated_at
		   )
		 ORDER BY f.created_at, f.id`,
		string(principalID), string(sessionID), at.UTC(),
	)
	if queryErr != nil {
		return nil, authenticationRepositoryFailure(queryErr)
	}
	defer rows.Close()
	factors := make([]StrongFactor, 0)
	for rows.Next() {
		var factorID, storedPrincipal, factorType, label, state string
		var createdAt, verifiedAt, revokedAt time.Time
		if scanErr := rows.Scan(&factorID, &storedPrincipal, &factorType, &label, &state, &createdAt, &verifiedAt, &revokedAt); scanErr != nil {
			return nil, authenticationRepositoryFailure(scanErr)
		}
		factor, factorErr := storedStrongFactor(factorID, storedPrincipal, factorType, label, state, createdAt, verifiedAt, revokedAt)
		if factorErr != nil {
			return nil, factorErr
		}
		factors = append(factors, factor)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, authenticationRepositoryFailure(rowsErr)
	}
	return factors, nil
}

func loadCurrentActivePasskeyFactor(
	ctx context.Context,
	tx pgx.Tx,
	principalID UserID,
	sessionID SessionID,
	factorID FactorID,
	at time.Time,
) (StrongFactor, error) {
	var storedID, storedPrincipal, factorType, label, state string
	var createdAt, verifiedAt, revokedAt time.Time
	queryErr := tx.QueryRow(
		ctx,
		`SELECT f.id, f.principal_id, f.factor_type, f.label, f.lifecycle_state, f.created_at,
		        f.verified_at, COALESCE(f.revoked_at, '0001-01-01T00:00:00Z'::timestamptz)
		 FROM omnexa_identity.mfa_factors f
		 JOIN omnexa_identity.passkey_credentials pk ON pk.factor_id = f.id
		 WHERE f.id = $1
		   AND f.principal_id = $2
		   AND f.factor_type = 'passkey'
		   AND f.lifecycle_state = 'active'
		   AND EXISTS (
		       SELECT 1
		       FROM omnexa_identity.sessions s
		       JOIN omnexa_identity.principals p ON p.id = s.principal_id
		       JOIN omnexa_identity.password_credentials pc ON pc.principal_id = s.principal_id
		       WHERE s.id = $3
		         AND s.principal_id = f.principal_id
		         AND s.revoked_at IS NULL
		         AND s.expires_at > $4
		         AND p.principal_type = 'human_user'
		         AND p.lifecycle_state = 'active'
		         AND pc.credential_version = s.credential_version
		         AND s.created_at >= p.updated_at
		   )
		 FOR UPDATE OF f`,
		string(factorID), string(principalID), string(sessionID), at.UTC(),
	).Scan(&storedID, &storedPrincipal, &factorType, &label, &state, &createdAt, &verifiedAt, &revokedAt)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		return StrongFactor{}, strongFactorNotFoundFailure()
	}
	if queryErr != nil {
		return StrongFactor{}, authenticationRepositoryFailure(queryErr)
	}
	return storedStrongFactor(storedID, storedPrincipal, factorType, label, state, createdAt, verifiedAt, revokedAt)
}

func revokeHumanSessions(ctx context.Context, tx pgx.Tx, principalID UserID, at time.Time) error {
	if _, sessionErr := tx.Exec(
		ctx,
		`UPDATE omnexa_identity.sessions
		 SET revoked_at = COALESCE(revoked_at, $2)
		 WHERE principal_id = $1`,
		string(principalID), at.UTC(),
	); sessionErr != nil {
		return authenticationRepositoryFailure(sessionErr)
	}
	if _, accessErr := tx.Exec(
		ctx,
		`UPDATE omnexa_identity.access_credentials a
		 SET revoked_at = COALESCE(a.revoked_at, $2)
		 FROM omnexa_identity.sessions s
		 WHERE a.session_id = s.id AND s.principal_id = $1`,
		string(principalID), at.UTC(),
	); accessErr != nil {
		return authenticationRepositoryFailure(accessErr)
	}
	if _, refreshErr := tx.Exec(
		ctx,
		`UPDATE omnexa_identity.refresh_credentials r
		 SET revoked_at = COALESCE(r.revoked_at, $2)
		 FROM omnexa_identity.sessions s
		 WHERE r.session_id = s.id AND s.principal_id = $1`,
		string(principalID), at.UTC(),
	); refreshErr != nil {
		return authenticationRepositoryFailure(refreshErr)
	}
	return nil
}

func storedStrongFactor(
	id string,
	principal string,
	factorType string,
	label string,
	state string,
	createdAt time.Time,
	verifiedAt time.Time,
	revokedAt time.Time,
) (StrongFactor, error) {
	factor := StrongFactor{
		id:         FactorID(id),
		principal:  UserID(principal),
		factorType: StrongFactorType(factorType),
		label:      label,
		state:      StrongFactorState(state),
		createdAt:  createdAt.UTC(),
	}
	if verifiedAt.Year() > 1 {
		factor.verifiedAt = verifiedAt.UTC()
	}
	if revokedAt.Year() > 1 {
		factor.revokedAt = revokedAt.UTC()
	}
	if !factor.valid() {
		return StrongFactor{}, invalidStoredUserFailure()
	}
	return factor, nil
}

func validStrongChallenge(challenge strongChallengeRecord) bool {
	return challenge.id.Valid() && challenge.principal.Valid() && challenge.session.Valid() && challenge.factor.Valid() &&
		(challenge.purpose == challengePurposeEnrollment || challenge.purpose == challengePurposeAssertion) &&
		!challenge.createdAt.IsZero() && challenge.expiresAt.After(challenge.createdAt)
}

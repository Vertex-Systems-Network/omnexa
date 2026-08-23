package identity

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const dummyPassword = "omnexa-p02-04-synthetic-missing-account-password"

// AuthenticationService owns the P02.04 human password and interactive-session
// lifecycle. Successful authentication remains identity proof only; it never
// evaluates or grants P02.05+ roles, permissions, policy, or business authority.
type AuthenticationService struct {
	repository AuthenticationRepository
	hasher     PasswordHasher
	policy     SessionPolicy
	contexts   ContextReauthorizer
	audit      SecurityAuditHook
	random     io.Reader
	dummyHash  string
}

// NewAuthenticationService constructs the governed P02.04 service. A nil context
// reauthorizer is valid only for sessions carrying no tenant/organization context.
func NewAuthenticationService(
	repository AuthenticationRepository,
	hasher PasswordHasher,
	policy SessionPolicy,
	contexts ContextReauthorizer,
	audit SecurityAuditHook,
) (*AuthenticationService, error) {
	return newAuthenticationService(repository, hasher, policy, contexts, audit, rand.Reader)
}

func newAuthenticationService(
	repository AuthenticationRepository,
	hasher PasswordHasher,
	policy SessionPolicy,
	contexts ContextReauthorizer,
	audit SecurityAuditHook,
	random io.Reader,
) (*AuthenticationService, error) {
	if repository == nil || hasher == nil || !policy.valid() || random == nil {
		return nil, repositoryInvalidFailure()
	}
	dummyHash, err := hasher.Hash(dummyPassword)
	if err != nil || dummyHash == "" {
		return nil, passwordHashFailure(nonNilCause(err))
	}
	return &AuthenticationService{
		repository: repository,
		hasher:     hasher,
		policy:     policy,
		contexts:   contexts,
		audit:      audit,
		random:     random,
		dummyHash:  dummyHash,
	}, nil
}

// EnrollPassword establishes the first password credential through this explicit
// authentication capability. It is not an ordinary profile update or reset flow.
func (service *AuthenticationService) EnrollPassword(
	ctx context.Context,
	userID UserID,
	password string,
	at time.Time,
) error {
	if service == nil || !userID.Valid() || at.IsZero() || !validPassword(password) {
		return passwordInvalidFailure()
	}
	hash, err := service.hasher.Hash(password)
	if err != nil {
		return err
	}
	_, err = service.repository.EnrollPassword(ctx, userID, hash, at.UTC())
	service.recordAudit(SecurityAuditPasswordEnrolled, userID, "", err == nil, at)
	return err
}

// AuthenticatePassword verifies human-user credentials with disclosure-safe
// failure behavior. Missing credentials execute the same password-hash boundary
// against a synthetic dummy representation before returning the generic failure.
func (service *AuthenticationService) AuthenticatePassword(
	ctx context.Context,
	userID UserID,
	password string,
	at time.Time,
) (Authentication, error) {
	if service == nil || at.IsZero() {
		return Authentication{}, authenticationFailure()
	}
	if !userID.Valid() {
		_, _ = service.hasher.Verify(password, service.dummyHash)
		service.recordAudit(SecurityAuditAuthenticationFail, "", "", false, at)
		return Authentication{}, authenticationFailure()
	}

	snapshot, err := service.repository.AuthenticationSnapshot(ctx, userID)
	if err != nil {
		if failure.IsCode(err, codeCredentialNotFound) {
			_, _ = service.hasher.Verify(password, service.dummyHash)
			service.recordAudit(SecurityAuditAuthenticationFail, userID, "", false, at)
			return Authentication{}, authenticationFailure()
		}
		return Authentication{}, err
	}
	verified, verifyErr := service.hasher.Verify(password, snapshot.passwordHash)
	if verifyErr != nil || !verified || snapshot.state != LifecycleActive || at.UTC().Before(snapshot.userUpdatedAt.UTC()) {
		service.recordAudit(SecurityAuditAuthenticationFail, userID, "", false, at)
		return Authentication{}, authenticationFailure()
	}

	authentication := Authentication{
		principalID:       userID,
		credentialVersion: snapshot.credentialVersion,
		authenticatedAt:   at.UTC(),
	}
	service.recordAudit(SecurityAuditAuthenticationOK, userID, "", true, at)
	return authentication, nil
}

// ChangePassword replaces the authenticated user's password, increments the
// credential version, and transactionally revokes all existing sessions.
func (service *AuthenticationService) ChangePassword(
	ctx context.Context,
	authentication Authentication,
	newPassword string,
	at time.Time,
) error {
	if service == nil || at.IsZero() || !validPassword(newPassword) {
		return passwordInvalidFailure()
	}
	snapshot, err := service.currentAuthentication(ctx, authentication, at)
	if err != nil {
		return err
	}
	hash, err := service.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	_, err = service.repository.ChangePassword(
		ctx,
		authentication.principalID,
		snapshot.credentialVersion,
		hash,
		at.UTC(),
	)
	service.recordAudit(SecurityAuditPasswordChanged, authentication.principalID, "", err == nil, at)
	return err
}

// IssueSession creates one bounded interactive session only from a current
// authentication proof. Non-empty tenant/organization context is reauthorized
// before any session/token state is persisted.
func (service *AuthenticationService) IssueSession(
	ctx context.Context,
	authentication Authentication,
	sessionContext SessionContext,
	deviceLabel string,
	at time.Time,
) (IssuedSession, error) {
	if service == nil || at.IsZero() || !sessionContext.valid() || !validDeviceLabel(deviceLabel) {
		return IssuedSession{}, sessionFailure()
	}
	snapshot, err := service.currentAuthentication(ctx, authentication, at)
	if err != nil {
		return IssuedSession{}, err
	}
	if err = service.reauthorize(ctx, authentication.principalID, sessionContext); err != nil {
		return IssuedSession{}, err
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return IssuedSession{}, err
	}
	secrets, err := generateSecrets(service.random)
	if err != nil {
		return IssuedSession{}, err
	}
	instant := at.UTC()
	session := Session{
		id:          sessionID,
		principalID: authentication.principalID,
		deviceLabel: deviceLabel,
		context:     sessionContext,
		createdAt:   instant,
		refreshedAt: instant,
		expiresAt:   instant.Add(service.policy.SessionLifetime),
	}
	accessExpiresAt := instant.Add(service.policy.AccessLifetime)
	refreshExpiresAt := instant.Add(service.policy.RefreshLifetime)
	record := sessionRecord{session: session, credentialVersion: snapshot.credentialVersion}
	if err = service.repository.CreateSession(
		ctx,
		record,
		secrets.accessDigest,
		accessExpiresAt,
		secrets.refreshDigest,
		refreshExpiresAt,
	); err != nil {
		service.recordAudit(SecurityAuditSessionIssued, authentication.principalID, sessionID, false, at)
		return IssuedSession{}, err
	}
	service.recordAudit(SecurityAuditSessionIssued, authentication.principalID, sessionID, true, at)
	return IssuedSession{
		session:          session,
		accessSecret:     secrets.access,
		accessExpiresAt:  accessExpiresAt,
		refreshSecret:    secrets.refresh,
		refreshExpiresAt: refreshExpiresAt,
	}, nil
}

// ValidateAccess resolves the opaque credential against current account/password
// state and then reauthorizes the stored context. The result authenticates a
// session only; callers must still perform later authorization/policy checks.
func (service *AuthenticationService) ValidateAccess(
	ctx context.Context,
	secret AccessSecret,
	at time.Time,
) (AuthenticatedSession, error) {
	if service == nil || at.IsZero() {
		return AuthenticatedSession{}, sessionFailure()
	}
	digest, err := digestAccessSecret(secret)
	if err != nil {
		service.recordAudit(SecurityAuditAccessValidated, "", "", false, at)
		return AuthenticatedSession{}, sessionFailure()
	}
	record, err := service.repository.AccessSession(ctx, digest, at.UTC())
	if err != nil {
		service.recordAudit(SecurityAuditAccessValidated, "", "", false, at)
		return AuthenticatedSession{}, sessionFailure()
	}
	if err = service.reauthorize(ctx, record.session.principalID, record.session.context); err != nil {
		service.recordAudit(SecurityAuditAccessValidated, record.session.principalID, record.session.id, false, at)
		return AuthenticatedSession{}, sessionFailure()
	}
	service.recordAudit(SecurityAuditAccessValidated, record.session.principalID, record.session.id, true, at)
	return AuthenticatedSession{session: record.session}, nil
}

// RotateRefresh consumes exactly one current refresh credential, reauthorizes its
// context first, revokes prior access credentials, and returns fresh opaque tokens.
func (service *AuthenticationService) RotateRefresh(
	ctx context.Context,
	secret RefreshSecret,
	at time.Time,
) (IssuedSession, error) {
	if service == nil || at.IsZero() {
		return IssuedSession{}, sessionFailure()
	}
	digest, err := digestRefreshSecret(secret)
	if err != nil {
		return IssuedSession{}, sessionFailure()
	}
	record, err := service.repository.RefreshSession(ctx, digest, at.UTC())
	if err != nil {
		return IssuedSession{}, sessionFailure()
	}
	if err = service.reauthorize(ctx, record.session.principalID, record.session.context); err != nil {
		service.recordAudit(SecurityAuditRefreshRotated, record.session.principalID, record.session.id, false, at)
		return IssuedSession{}, sessionFailure()
	}
	secrets, err := generateSecrets(service.random)
	if err != nil {
		return IssuedSession{}, err
	}
	instant := at.UTC()
	accessExpiresAt := minTime(instant.Add(service.policy.AccessLifetime), record.session.expiresAt)
	refreshExpiresAt := minTime(instant.Add(service.policy.RefreshLifetime), record.session.expiresAt)
	if !instant.Before(accessExpiresAt) || !instant.Before(refreshExpiresAt) {
		return IssuedSession{}, sessionFailure()
	}
	rotated, err := service.repository.RotateRefresh(
		ctx,
		record.session.id,
		digest,
		secrets.accessDigest,
		accessExpiresAt,
		secrets.refreshDigest,
		refreshExpiresAt,
		instant,
	)
	if err != nil {
		service.recordAudit(SecurityAuditRefreshRotated, record.session.principalID, record.session.id, false, at)
		return IssuedSession{}, sessionFailure()
	}
	service.recordAudit(SecurityAuditRefreshRotated, rotated.session.principalID, rotated.session.id, true, at)
	return IssuedSession{
		session:          rotated.session,
		accessSecret:     secrets.access,
		accessExpiresAt:  accessExpiresAt,
		refreshSecret:    secrets.refresh,
		refreshExpiresAt: refreshExpiresAt,
	}, nil
}

// RevokeSession explicitly terminates one session owned by the currently
// authenticated human user. Repeated revocation remains deterministic/idempotent.
func (service *AuthenticationService) RevokeSession(
	ctx context.Context,
	authentication Authentication,
	sessionID SessionID,
	at time.Time,
) (Session, error) {
	if service == nil || !sessionID.Valid() || at.IsZero() {
		return Session{}, sessionFailure()
	}
	if _, err := service.currentAuthentication(ctx, authentication, at); err != nil {
		return Session{}, err
	}
	record, err := service.repository.RevokeSession(ctx, authentication.principalID, sessionID, at.UTC())
	service.recordAudit(SecurityAuditSessionRevoked, authentication.principalID, sessionID, err == nil, at)
	if err != nil {
		return Session{}, sessionFailure()
	}
	return record.session, nil
}

// ListSessions returns only safe device/session inventory for the currently
// authenticated user. It includes invalidated/revoked state but never secrets.
func (service *AuthenticationService) ListSessions(
	ctx context.Context,
	authentication Authentication,
	at time.Time,
) ([]Session, error) {
	if service == nil || at.IsZero() {
		return nil, authenticationInvalidFailure()
	}
	if _, err := service.currentAuthentication(ctx, authentication, at); err != nil {
		return nil, err
	}
	return service.repository.ListSessions(ctx, authentication.principalID)
}

func (service *AuthenticationService) currentAuthentication(
	ctx context.Context,
	authentication Authentication,
	at time.Time,
) (authenticationSnapshot, error) {
	if !authentication.valid() || at.IsZero() || at.UTC().Before(authentication.authenticatedAt.UTC()) {
		return authenticationSnapshot{}, authenticationInvalidFailure()
	}
	snapshot, err := service.repository.AuthenticationSnapshot(ctx, authentication.principalID)
	if err != nil || snapshot.state != LifecycleActive || snapshot.credentialVersion != authentication.credentialVersion {
		return authenticationSnapshot{}, authenticationInvalidFailure()
	}
	if authentication.authenticatedAt.UTC().Before(snapshot.userUpdatedAt.UTC()) {
		return authenticationSnapshot{}, authenticationInvalidFailure()
	}
	return snapshot, nil
}

func (service *AuthenticationService) reauthorize(
	ctx context.Context,
	userID UserID,
	sessionContext SessionContext,
) error {
	if !sessionContext.valid() {
		return sessionContextFailure()
	}
	if sessionContext.Empty() {
		return nil
	}
	if service.contexts == nil || service.contexts.Reauthorize(ctx, userID, sessionContext) != nil {
		return sessionContextFailure()
	}
	return nil
}

func (service *AuthenticationService) recordAudit(
	action SecurityAuditAction,
	userID UserID,
	sessionID SessionID,
	succeeded bool,
	at time.Time,
) {
	if service == nil || service.audit == nil || at.IsZero() {
		return
	}
	service.audit.RecordSecurityEvent(SecurityAuditEvent{
		Action:      action,
		PrincipalID: userID,
		SessionID:   sessionID,
		Succeeded:   succeeded,
		OccurredAt:  at.UTC(),
	})
}

func minTime(left, right time.Time) time.Time {
	if left.Before(right) {
		return left
	}
	return right
}

func nonNilCause(err error) error {
	if err != nil {
		return err
	}
	return errors.New("password hash initialization failed")
}

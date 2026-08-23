package identity

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	accessSecretPrefix  = "ox_at_"
	refreshSecretPrefix = "ox_rt_"
	secretBytes         = 32
	maxDeviceLabelRunes = 128
)

// SessionID is the canonical UUIDv7 identifier for one interactive human session.
type SessionID string

// Valid reports whether the identifier is a canonical UUIDv7 value.
func (id SessionID) Valid() bool { return validUUIDv7(string(id)) }

// SessionPolicy defines bounded interactive credential lifetimes. Access credentials
// must always be shorter-lived than refresh credentials and the enclosing session.
type SessionPolicy struct {
	AccessLifetime  time.Duration
	RefreshLifetime time.Duration
	SessionLifetime time.Duration
}

// DefaultSessionPolicy returns the governed interactive-session baseline.
func DefaultSessionPolicy() SessionPolicy {
	return SessionPolicy{
		AccessLifetime:  15 * time.Minute,
		RefreshLifetime: 7 * 24 * time.Hour,
		SessionLifetime: 30 * 24 * time.Hour,
	}
}

func (policy SessionPolicy) valid() bool {
	return policy.AccessLifetime > 0 &&
		policy.RefreshLifetime > policy.AccessLifetime &&
		policy.SessionLifetime >= policy.RefreshLifetime
}

// SessionContext carries only caller-selected tenant/organization references. It is
// never authorization authority: a ContextReauthorizer must re-resolve non-empty
// context before session issuance, access validation, and refresh rotation.
type SessionContext struct {
	tenantID       string
	organizationID string
}

// NewSessionContext validates opaque P02.02/P02.03 UUIDv7 references without
// importing or writing those owner domains. Empty/empty creates a platform context.
func NewSessionContext(tenantID, organizationID string) (SessionContext, error) {
	context := SessionContext{
		tenantID:       strings.TrimSpace(tenantID),
		organizationID: strings.TrimSpace(organizationID),
	}
	if !context.valid() {
		return SessionContext{}, sessionContextFailure()
	}
	return context, nil
}

func (sessionContext SessionContext) valid() bool {
	if sessionContext.tenantID == "" {
		return sessionContext.organizationID == ""
	}
	if !validUUIDv7(sessionContext.tenantID) {
		return false
	}
	return sessionContext.organizationID == "" || validUUIDv7(sessionContext.organizationID)
}

// Empty reports whether no tenant/organization selection is carried.
func (sessionContext SessionContext) Empty() bool {
	return sessionContext.tenantID == "" && sessionContext.organizationID == ""
}

// TenantID returns the opaque tenant reference; the value is not authorization proof.
func (sessionContext SessionContext) TenantID() string { return sessionContext.tenantID }

// OrganizationID returns the opaque organization reference; the value is not authorization proof.
func (sessionContext SessionContext) OrganizationID() string { return sessionContext.organizationID }

// ContextReauthorizer is an inversion boundary for current P02.02/P02.03 relationship
// state. Implementations must fail closed when the requested context is stale, revoked,
// cross-tenant, or otherwise no longer current. Success is context validation only and
// does not evaluate P02.05+ roles, permissions, or policy.
type ContextReauthorizer interface {
	Reauthorize(context.Context, UserID, SessionContext) error
}

// AccessSecret is a restricted opaque access credential. Its String/JSON forms are
// always redacted; Reveal is intentionally explicit for a transport adapter.
type AccessSecret struct{ value string }

// String never reveals authentication material.
func (secret AccessSecret) String() string { return "[REDACTED]" }

// GoString never reveals authentication material.
func (secret AccessSecret) GoString() string { return "[REDACTED]" }

// MarshalJSON keeps accidental structured logging/serialization disclosure-safe.
func (secret AccessSecret) MarshalJSON() ([]byte, error) { return json.Marshal("[REDACTED]") }

// Reveal returns the credential only for an explicit transport/storage boundary.
func (secret AccessSecret) Reveal() string { return secret.value }

func (secret AccessSecret) valid() bool { return validOpaqueSecret(secret.value, accessSecretPrefix) }

// RefreshSecret is a restricted opaque refresh credential. Its String/JSON forms
// are always redacted; Reveal is intentionally explicit for a transport adapter.
type RefreshSecret struct{ value string }

// String never reveals authentication material.
func (secret RefreshSecret) String() string { return "[REDACTED]" }

// GoString never reveals authentication material.
func (secret RefreshSecret) GoString() string { return "[REDACTED]" }

// MarshalJSON keeps accidental structured logging/serialization disclosure-safe.
func (secret RefreshSecret) MarshalJSON() ([]byte, error) { return json.Marshal("[REDACTED]") }

// Reveal returns the credential only for an explicit transport/storage boundary.
func (secret RefreshSecret) Reveal() string { return secret.value }

func (secret RefreshSecret) valid() bool { return validOpaqueSecret(secret.value, refreshSecretPrefix) }

// ParseAccessSecret accepts only the canonical opaque access-token shape and never
// includes rejected material in returned failures.
func ParseAccessSecret(value string) (AccessSecret, error) {
	secret := AccessSecret{value: value}
	if !secret.valid() {
		return AccessSecret{}, sessionCredentialFailure()
	}
	return secret, nil
}

// ParseRefreshSecret accepts only the canonical opaque refresh-token shape and never
// includes rejected material in returned failures.
func ParseRefreshSecret(value string) (RefreshSecret, error) {
	secret := RefreshSecret{value: value}
	if !secret.valid() {
		return RefreshSecret{}, sessionCredentialFailure()
	}
	return secret, nil
}

// Authentication is a non-authorizing proof produced only after successful password
// verification. Private fields prevent callers from minting a valid proof directly.
type Authentication struct {
	principalID       UserID
	credentialVersion uint64
	authenticatedAt   time.Time
}

// PrincipalID returns the authenticated human User identifier.
func (authentication Authentication) PrincipalID() UserID { return authentication.principalID }

// AuthenticatedAt returns the UTC authentication instant.
func (authentication Authentication) AuthenticatedAt() time.Time {
	return authentication.authenticatedAt
}

func (authentication Authentication) valid() bool {
	return authentication.principalID.Valid() &&
		authentication.credentialVersion > 0 &&
		!authentication.authenticatedAt.IsZero()
}

// Session is the classification-safe interactive-session inventory projection. It
// contains no password hash, access secret, refresh secret, or secret digest.
type Session struct {
	id          SessionID
	principalID UserID
	deviceLabel string
	context     SessionContext
	createdAt   time.Time
	refreshedAt time.Time
	expiresAt   time.Time
	revokedAt   time.Time
	invalidated bool
}

// ID returns the canonical session identifier.
func (session Session) ID() SessionID { return session.id }

// PrincipalID returns the authenticated human User identifier.
func (session Session) PrincipalID() UserID { return session.principalID }

// DeviceLabel returns bounded caller-supplied inventory metadata.
func (session Session) DeviceLabel() string { return session.deviceLabel }

// Context returns non-authorizing tenant/organization hints that require current reauthorization.
func (session Session) Context() SessionContext { return session.context }

// CreatedAt returns the session creation instant.
func (session Session) CreatedAt() time.Time { return session.createdAt }

// RefreshedAt returns the most recent successful refresh rotation instant.
func (session Session) RefreshedAt() time.Time { return session.refreshedAt }

// ExpiresAt returns the absolute session expiry instant.
func (session Session) ExpiresAt() time.Time { return session.expiresAt }

// RevokedAt returns zero when the session has not been explicitly revoked.
func (session Session) RevokedAt() time.Time { return session.revokedAt }

// Active reports whether this inventory projection is currently usable before the
// mandatory context reauthorization step. It is not an authorization decision.
func (session Session) Active(at time.Time) bool {
	if at.IsZero() {
		return false
	}
	instant := at.UTC()
	return !session.invalidated && session.revokedAt.IsZero() && instant.Before(session.expiresAt)
}

// IssuedSession contains the one-time returned access/refresh credentials plus the
// safe session inventory projection. Raw credentials are not persisted by this value.
type IssuedSession struct {
	session          Session
	accessSecret     AccessSecret
	accessExpiresAt  time.Time
	refreshSecret    RefreshSecret
	refreshExpiresAt time.Time
}

// Session returns the safe session projection.
func (issued IssuedSession) Session() Session { return issued.session }

// AccessSecret returns the restricted access credential.
func (issued IssuedSession) AccessSecret() AccessSecret { return issued.accessSecret }

// AccessExpiresAt returns access-credential expiry.
func (issued IssuedSession) AccessExpiresAt() time.Time { return issued.accessExpiresAt }

// RefreshSecret returns the restricted refresh credential.
func (issued IssuedSession) RefreshSecret() RefreshSecret { return issued.refreshSecret }

// RefreshExpiresAt returns refresh-credential expiry.
func (issued IssuedSession) RefreshExpiresAt() time.Time { return issued.refreshExpiresAt }

// AuthenticatedSession is the current, revalidated authentication result for one
// access credential. It grants no P02.05+ role, permission, or business authority.
type AuthenticatedSession struct{ session Session }

// Session returns the safe session projection used for later policy evaluation.
func (authenticated AuthenticatedSession) Session() Session { return authenticated.session }

// SecurityAuditAction is a bounded, secret-free lifecycle action vocabulary for P02.04 hooks.
type SecurityAuditAction string

const (
	SecurityAuditPasswordEnrolled   SecurityAuditAction = "identity.password.enrolled"
	SecurityAuditPasswordChanged    SecurityAuditAction = "identity.password.changed"
	SecurityAuditAuthenticationOK   SecurityAuditAction = "identity.authentication.succeeded"
	SecurityAuditAuthenticationFail SecurityAuditAction = "identity.authentication.failed"
	SecurityAuditSessionIssued      SecurityAuditAction = "identity.session.issued"
	SecurityAuditAccessValidated    SecurityAuditAction = "identity.session.access_validated"
	SecurityAuditRefreshRotated     SecurityAuditAction = "identity.session.refresh_rotated"
	SecurityAuditSessionRevoked     SecurityAuditAction = "identity.session.revoked"
)

// SecurityAuditEvent intentionally has no field capable of carrying passwords,
// password hashes, token material, secret digests, or authorization decisions.
type SecurityAuditEvent struct {
	Action      SecurityAuditAction
	PrincipalID UserID
	SessionID   SessionID
	Succeeded   bool
	OccurredAt  time.Time
}

// SecurityAuditHook receives classification-safe lifecycle facts. Durable P02 exit
// audit product behavior remains owned by later P02.10; this hook grants no authority.
type SecurityAuditHook interface {
	RecordSecurityEvent(SecurityAuditEvent)
}

type generatedSecrets struct {
	access        AccessSecret
	accessDigest  [sha256.Size]byte
	refresh       RefreshSecret
	refreshDigest [sha256.Size]byte
}

func generateSessionID() (SessionID, error) {
	identifier, err := uuid.NewV7()
	if err != nil {
		return "", identifierFailure(err)
	}
	return SessionID(identifier.String()), nil
}

func generateSecrets(reader io.Reader) (generatedSecrets, error) {
	access, accessDigest, err := generateSecret(reader, accessSecretPrefix)
	if err != nil {
		return generatedSecrets{}, err
	}
	refresh, refreshDigest, err := generateSecret(reader, refreshSecretPrefix)
	if err != nil {
		return generatedSecrets{}, err
	}
	return generatedSecrets{
		access:        AccessSecret{value: access},
		accessDigest:  accessDigest,
		refresh:       RefreshSecret{value: refresh},
		refreshDigest: refreshDigest,
	}, nil
}

func generateSecret(reader io.Reader, prefix string) (string, [sha256.Size]byte, error) {
	if reader == nil {
		reader = rand.Reader
	}
	payload := make([]byte, secretBytes)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return "", [sha256.Size]byte{}, sessionSecretGenerationFailure(err)
	}
	value := prefix + base64.RawURLEncoding.EncodeToString(payload)
	return value, sha256.Sum256([]byte(value)), nil
}

func digestAccessSecret(secret AccessSecret) ([sha256.Size]byte, error) {
	if !secret.valid() {
		return [sha256.Size]byte{}, sessionCredentialFailure()
	}
	return sha256.Sum256([]byte(secret.value)), nil
}

func digestRefreshSecret(secret RefreshSecret) ([sha256.Size]byte, error) {
	if !secret.valid() {
		return [sha256.Size]byte{}, sessionCredentialFailure()
	}
	return sha256.Sum256([]byte(secret.value)), nil
}

func validOpaqueSecret(value, prefix string) bool {
	if !strings.HasPrefix(value, prefix) {
		return false
	}
	encoded := strings.TrimPrefix(value, prefix)
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	return err == nil && len(decoded) == secretBytes
}

func secureEqual(left, right []byte) bool {
	return len(left) == len(right) && subtle.ConstantTimeCompare(left, right) == 1
}

func validUUIDv7(value string) bool {
	parsed, err := uuid.Parse(value)
	return err == nil && parsed.Version() == 7 && strings.ToLower(parsed.String()) == value
}

func validDeviceLabel(value string) bool {
	trimmed := strings.TrimSpace(value)
	return trimmed == value && len([]rune(value)) <= maxDeviceLabelRunes && !strings.ContainsAny(value, "\r\n\x00")
}

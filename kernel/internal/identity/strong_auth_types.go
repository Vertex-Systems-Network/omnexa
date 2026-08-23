package identity

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	challengePrefix          = "ox_ch_"
	recoveryPrefix           = "ox_rc_"
	strongAuthRandomBytes    = 32
	maxStrongAuthLabelRunes  = 128
	maxPasskeyCredentialSize = 1024
	maxPasskeyPublicKeySize  = 4096
	maxPasskeyResponseSize   = 64 * 1024
	minRecoveryCodes         = 4
	maxRecoveryCodes         = 16
)

// StrongAuthenticationPolicy is the fixed P02.07 policy surface. Tenant-specific
// settings remain owned by P02.09 and are intentionally not read here.
type StrongAuthenticationPolicy struct {
	ChallengeLifetime                 time.Duration
	StepUpLifetime                    time.Duration
	RecoveryCodeCount                 int
	InvalidateSessionsOnFactorRemoval bool
}

// DefaultStrongAuthenticationPolicy returns the bounded platform baseline.
func DefaultStrongAuthenticationPolicy() StrongAuthenticationPolicy {
	return StrongAuthenticationPolicy{
		ChallengeLifetime:                 5 * time.Minute,
		StepUpLifetime:                    10 * time.Minute,
		RecoveryCodeCount:                 8,
		InvalidateSessionsOnFactorRemoval: true,
	}
}

func (policy StrongAuthenticationPolicy) valid() bool {
	return policy.ChallengeLifetime > 0 &&
		policy.StepUpLifetime > 0 &&
		policy.RecoveryCodeCount >= minRecoveryCodes &&
		policy.RecoveryCodeCount <= maxRecoveryCodes
}

// FactorID is the UUIDv7 identifier for one human-user strong-authentication factor.
type FactorID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id FactorID) Valid() bool { return validUUIDv7(string(id)) }

// ChallengeID is the UUIDv7 identifier for one bounded, one-time authentication challenge.
type ChallengeID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id ChallengeID) Valid() bool { return validUUIDv7(string(id)) }

// RecoverySetID is the UUIDv7 identifier for one recovery-code generation.
type RecoverySetID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id RecoverySetID) Valid() bool { return validUUIDv7(string(id)) }

// StrongFactorType is the P02.07 factor vocabulary. Passkeys are the only persisted
// factor type in this package; service-account/API credentials remain P02.08.
type StrongFactorType string

const StrongFactorPasskey StrongFactorType = "passkey"

func (factorType StrongFactorType) valid() bool { return factorType == StrongFactorPasskey }

// StrongFactorState is the deterministic enrollment/removal lifecycle.
type StrongFactorState string

const (
	StrongFactorPending StrongFactorState = "pending"
	StrongFactorActive  StrongFactorState = "active"
	StrongFactorRevoked StrongFactorState = "revoked"
)

func (state StrongFactorState) valid() bool {
	switch state {
	case StrongFactorPending, StrongFactorActive, StrongFactorRevoked:
		return true
	default:
		return false
	}
}

// StrongFactor is the classification-safe inventory projection for a human-user factor.
type StrongFactor struct {
	id         FactorID
	principal  UserID
	factorType StrongFactorType
	label      string
	state      StrongFactorState
	createdAt  time.Time
	verifiedAt time.Time
	revokedAt  time.Time
}

func (factor StrongFactor) ID() FactorID             { return factor.id }
func (factor StrongFactor) PrincipalID() UserID      { return factor.principal }
func (factor StrongFactor) Type() StrongFactorType   { return factor.factorType }
func (factor StrongFactor) Label() string            { return factor.label }
func (factor StrongFactor) State() StrongFactorState { return factor.state }
func (factor StrongFactor) CreatedAt() time.Time     { return factor.createdAt }
func (factor StrongFactor) VerifiedAt() time.Time    { return factor.verifiedAt }
func (factor StrongFactor) RevokedAt() time.Time     { return factor.revokedAt }
func (factor StrongFactor) Active() bool             { return factor.state == StrongFactorActive }
func (factor StrongFactor) valid() bool {
	if !factor.id.Valid() || !factor.principal.Valid() || !factor.factorType.valid() || !factor.state.valid() || !validStrongAuthLabel(factor.label) || factor.createdAt.IsZero() {
		return false
	}
	if factor.state == StrongFactorPending {
		return factor.verifiedAt.IsZero() && factor.revokedAt.IsZero()
	}
	if factor.verifiedAt.IsZero() || factor.verifiedAt.Before(factor.createdAt) {
		return false
	}
	if factor.state == StrongFactorActive {
		return factor.revokedAt.IsZero()
	}
	return !factor.revokedAt.IsZero() && !factor.revokedAt.Before(factor.verifiedAt)
}

// ChallengeSecret is RESTRICTED authentication material. Its ordinary string and
// JSON representations are always redacted; Reveal is explicit for protocol adapters.
type ChallengeSecret struct{ value string }

func (secret ChallengeSecret) String() string               { return "[REDACTED]" }
func (secret ChallengeSecret) GoString() string             { return "[REDACTED]" }
func (secret ChallengeSecret) MarshalJSON() ([]byte, error) { return json.Marshal("[REDACTED]") }
func (secret ChallengeSecret) Reveal() string               { return secret.value }
func (secret ChallengeSecret) valid() bool {
	return validRestrictedValue(secret.value, challengePrefix)
}

// ParseChallengeSecret accepts only the canonical high-entropy challenge shape.
func ParseChallengeSecret(value string) (ChallengeSecret, error) {
	secret := ChallengeSecret{value: value}
	if !secret.valid() {
		return ChallengeSecret{}, strongAuthChallengeFailure()
	}
	return secret, nil
}

// RecoveryCode is RESTRICTED one-time recovery material. The raw value is returned
// only at issuance/transport and never persisted by the identity repositories.
type RecoveryCode struct{ value string }

func (code RecoveryCode) String() string               { return "[REDACTED]" }
func (code RecoveryCode) GoString() string             { return "[REDACTED]" }
func (code RecoveryCode) MarshalJSON() ([]byte, error) { return json.Marshal("[REDACTED]") }
func (code RecoveryCode) Reveal() string               { return code.value }
func (code RecoveryCode) valid() bool                  { return validRestrictedValue(code.value, recoveryPrefix) }

// ParseRecoveryCode accepts only the canonical high-entropy recovery-code shape.
func ParseRecoveryCode(value string) (RecoveryCode, error) {
	code := RecoveryCode{value: value}
	if !code.valid() {
		return RecoveryCode{}, recoveryCodeFailure()
	}
	return code, nil
}

// PasskeyResponse is an opaque protocol response. It is treated as authentication
// material and cannot be accidentally serialized into telemetry.
type PasskeyResponse struct{ value []byte }

// NewPasskeyResponse validates and copies a bounded protocol response.
func NewPasskeyResponse(value []byte) (PasskeyResponse, error) {
	if len(value) == 0 || len(value) > maxPasskeyResponseSize {
		return PasskeyResponse{}, passkeyVerificationFailure()
	}
	return PasskeyResponse{value: append([]byte(nil), value...)}, nil
}

func (response PasskeyResponse) String() string               { return "[REDACTED]" }
func (response PasskeyResponse) GoString() string             { return "[REDACTED]" }
func (response PasskeyResponse) MarshalJSON() ([]byte, error) { return json.Marshal("[REDACTED]") }

// Reveal returns a copy only for the injected protocol verifier boundary.
func (response PasskeyResponse) Reveal() []byte { return append([]byte(nil), response.value...) }
func (response PasskeyResponse) valid() bool {
	return len(response.value) > 0 && len(response.value) <= maxPasskeyResponseSize
}

// VerifiedPasskeyCredential is the normalized public credential output of an
// approved passkey/WebAuthn verifier. No private key or authenticator secret exists here.
type VerifiedPasskeyCredential struct {
	credentialID     []byte
	publicKey        []byte
	counterSupported bool
	signCount        uint32
}

// NewVerifiedPasskeyCredential constructs verifier output after external protocol verification.
func NewVerifiedPasskeyCredential(credentialID, publicKey []byte, counterSupported bool, signCount uint32) (VerifiedPasskeyCredential, error) {
	if !validPasskeyCredentialMaterial(credentialID, publicKey) || (!counterSupported && signCount != 0) {
		return VerifiedPasskeyCredential{}, passkeyVerificationFailure()
	}
	return VerifiedPasskeyCredential{
		credentialID:     append([]byte(nil), credentialID...),
		publicKey:        append([]byte(nil), publicKey...),
		counterSupported: counterSupported,
		signCount:        signCount,
	}, nil
}

func (credential VerifiedPasskeyCredential) CredentialID() []byte {
	return append([]byte(nil), credential.credentialID...)
}
func (credential VerifiedPasskeyCredential) PublicKey() []byte {
	return append([]byte(nil), credential.publicKey...)
}
func (credential VerifiedPasskeyCredential) CounterSupported() bool {
	return credential.counterSupported
}
func (credential VerifiedPasskeyCredential) SignCount() uint32 { return credential.signCount }
func (credential VerifiedPasskeyCredential) valid() bool {
	return validPasskeyCredentialMaterial(credential.credentialID, credential.publicKey) &&
		(credential.counterSupported || credential.signCount == 0)
}

// PasskeyRegistrationVerification is the bounded input to an approved verifier.
type PasskeyRegistrationVerification struct {
	principal UserID
	session   SessionID
	challenge ChallengeSecret
	response  PasskeyResponse
}

func (verification PasskeyRegistrationVerification) PrincipalID() UserID {
	return verification.principal
}
func (verification PasskeyRegistrationVerification) SessionID() SessionID {
	return verification.session
}
func (verification PasskeyRegistrationVerification) Challenge() ChallengeSecret {
	return verification.challenge
}
func (verification PasskeyRegistrationVerification) Response() PasskeyResponse {
	return verification.response
}

// PasskeyAssertionVerification carries the stored public credential plus current challenge.
type PasskeyAssertionVerification struct {
	principal         UserID
	session           SessionID
	challenge         ChallengeSecret
	response          PasskeyResponse
	credentialID      []byte
	publicKey         []byte
	counterSupported  bool
	previousSignCount uint32
}

func (verification PasskeyAssertionVerification) PrincipalID() UserID  { return verification.principal }
func (verification PasskeyAssertionVerification) SessionID() SessionID { return verification.session }
func (verification PasskeyAssertionVerification) Challenge() ChallengeSecret {
	return verification.challenge
}
func (verification PasskeyAssertionVerification) Response() PasskeyResponse {
	return verification.response
}
func (verification PasskeyAssertionVerification) CredentialID() []byte {
	return append([]byte(nil), verification.credentialID...)
}
func (verification PasskeyAssertionVerification) PublicKey() []byte {
	return append([]byte(nil), verification.publicKey...)
}
func (verification PasskeyAssertionVerification) CounterSupported() bool {
	return verification.counterSupported
}
func (verification PasskeyAssertionVerification) PreviousSignCount() uint32 {
	return verification.previousSignCount
}

// PasskeyAssertionResult is verifier output after cryptographic assertion validation.
type PasskeyAssertionResult struct {
	counterSupported bool
	signCount        uint32
}

// NewPasskeyAssertionResult constructs normalized assertion counter evidence.
func NewPasskeyAssertionResult(counterSupported bool, signCount uint32) (PasskeyAssertionResult, error) {
	if !counterSupported && signCount != 0 {
		return PasskeyAssertionResult{}, passkeyVerificationFailure()
	}
	return PasskeyAssertionResult{counterSupported: counterSupported, signCount: signCount}, nil
}

func (result PasskeyAssertionResult) CounterSupported() bool { return result.counterSupported }
func (result PasskeyAssertionResult) SignCount() uint32      { return result.signCount }

// PasskeyVerifier owns approved protocol/cryptographic verification. P02.07 does
// not implement custom WebAuthn cryptography inside the identity domain service.
type PasskeyVerifier interface {
	VerifyRegistration(context.Context, PasskeyRegistrationVerification) (VerifiedPasskeyCredential, error)
	VerifyAssertion(context.Context, PasskeyAssertionVerification) (PasskeyAssertionResult, error)
}

// PasskeyEnrollment contains one pending factor and one-time registration challenge.
type PasskeyEnrollment struct {
	factor    StrongFactor
	challenge ChallengeID
	secret    ChallengeSecret
	expiresAt time.Time
}

func (enrollment PasskeyEnrollment) Factor() StrongFactor       { return enrollment.factor }
func (enrollment PasskeyEnrollment) ChallengeID() ChallengeID   { return enrollment.challenge }
func (enrollment PasskeyEnrollment) Challenge() ChallengeSecret { return enrollment.secret }
func (enrollment PasskeyEnrollment) ExpiresAt() time.Time       { return enrollment.expiresAt }

// PasskeyAssertionChallenge contains a one-time session-bound assertion challenge.
type PasskeyAssertionChallenge struct {
	factor    StrongFactor
	challenge ChallengeID
	secret    ChallengeSecret
	expiresAt time.Time
}

func (assertion PasskeyAssertionChallenge) Factor() StrongFactor       { return assertion.factor }
func (assertion PasskeyAssertionChallenge) ChallengeID() ChallengeID   { return assertion.challenge }
func (assertion PasskeyAssertionChallenge) Challenge() ChallengeSecret { return assertion.secret }
func (assertion PasskeyAssertionChallenge) ExpiresAt() time.Time       { return assertion.expiresAt }

// StrongAuthentication is in-memory second-factor evidence. It is session-bound,
// short-lived, deliberately non-serializable, and grants no authorization capability.
type StrongAuthentication struct {
	principal       UserID
	session         SessionID
	method          string
	factor          FactorID
	authenticatedAt time.Time
}

func (authentication StrongAuthentication) String() string   { return "[REDACTED]" }
func (authentication StrongAuthentication) GoString() string { return "[REDACTED]" }
func (authentication StrongAuthentication) MarshalJSON() ([]byte, error) {
	return json.Marshal("[REDACTED]")
}
func (authentication StrongAuthentication) PrincipalID() UserID  { return authentication.principal }
func (authentication StrongAuthentication) SessionID() SessionID { return authentication.session }
func (authentication StrongAuthentication) AuthenticatedAt() time.Time {
	return authentication.authenticatedAt
}
func (authentication StrongAuthentication) valid() bool {
	if !authentication.principal.Valid() || !authentication.session.Valid() || authentication.authenticatedAt.IsZero() {
		return false
	}
	switch authentication.method {
	case "passkey":
		return authentication.factor.Valid()
	case "recovery_code":
		return authentication.factor == ""
	default:
		return false
	}
}

// RecoveryCodeBundle is the only one-time raw recovery-code issuance surface.
type RecoveryCodeBundle struct {
	setID RecoverySetID
	codes []RecoveryCode
}

func (bundle RecoveryCodeBundle) String() string               { return "[REDACTED]" }
func (bundle RecoveryCodeBundle) GoString() string             { return "[REDACTED]" }
func (bundle RecoveryCodeBundle) MarshalJSON() ([]byte, error) { return json.Marshal("[REDACTED]") }
func (bundle RecoveryCodeBundle) SetID() RecoverySetID         { return bundle.setID }
func (bundle RecoveryCodeBundle) Codes() []RecoveryCode {
	return append([]RecoveryCode(nil), bundle.codes...)
}

// P02.07 classification-safe security action vocabulary.
const (
	SecurityAuditStrongEnrollmentStarted  SecurityAuditAction = "identity.strong_auth.enrollment_started"
	SecurityAuditStrongEnrollmentVerified SecurityAuditAction = "identity.strong_auth.enrollment_verified"
	SecurityAuditStrongAssertionVerified  SecurityAuditAction = "identity.strong_auth.assertion_verified"
	SecurityAuditStrongChallengeFailed    SecurityAuditAction = "identity.strong_auth.challenge_failed"
	SecurityAuditStrongRecoveryIssued     SecurityAuditAction = "identity.strong_auth.recovery_issued"
	SecurityAuditStrongRecoveryUsed       SecurityAuditAction = "identity.strong_auth.recovery_used"
	SecurityAuditStrongFactorRemoved      SecurityAuditAction = "identity.strong_auth.factor_removed"
	SecurityAuditStrongStepUpRequired     SecurityAuditAction = "identity.strong_auth.step_up_required"
)

type strongChallengeRecord struct {
	id        ChallengeID
	principal UserID
	session   SessionID
	factor    FactorID
	purpose   string
	digest    [sha256.Size]byte
	createdAt time.Time
	expiresAt time.Time
}

type passkeyCredentialRecord struct {
	factor           FactorID
	credentialID     []byte
	publicKey        []byte
	counterSupported bool
	signCount        uint32
}

func (record passkeyCredentialRecord) valid() bool {
	return record.factor.Valid() && validPasskeyCredentialMaterial(record.credentialID, record.publicKey) &&
		(record.counterSupported || record.signCount == 0)
}

func generateStrongAuthID() (string, error) {
	identifier, generationErr := uuid.NewV7()
	if generationErr != nil {
		return "", identifierFailure(generationErr)
	}
	return identifier.String(), nil
}

func generateRestrictedValue(reader io.Reader, prefix string) (string, [sha256.Size]byte, error) {
	if reader == nil {
		reader = rand.Reader
	}
	payload := make([]byte, strongAuthRandomBytes)
	if _, readErr := io.ReadFull(reader, payload); readErr != nil {
		return "", [sha256.Size]byte{}, strongAuthSecretGenerationFailure(readErr)
	}
	value := prefix + base64.RawURLEncoding.EncodeToString(payload)
	return value, sha256.Sum256([]byte(value)), nil
}

func digestChallenge(secret ChallengeSecret) ([sha256.Size]byte, error) {
	if !secret.valid() {
		return [sha256.Size]byte{}, strongAuthChallengeFailure()
	}
	return sha256.Sum256([]byte(secret.value)), nil
}

func digestRecovery(code RecoveryCode) ([sha256.Size]byte, error) {
	if !code.valid() {
		return [sha256.Size]byte{}, recoveryCodeFailure()
	}
	return sha256.Sum256([]byte(code.value)), nil
}

func validRestrictedValue(value, prefix string) bool {
	if !strings.HasPrefix(value, prefix) {
		return false
	}
	decoded, decodeErr := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(value, prefix))
	return decodeErr == nil && len(decoded) == strongAuthRandomBytes
}

func validStrongAuthLabel(value string) bool {
	return value != "" && strings.TrimSpace(value) == value && utf8.RuneCountInString(value) <= maxStrongAuthLabelRunes && !strings.ContainsAny(value, "\r\n\x00")
}

func validPasskeyCredentialMaterial(credentialID, publicKey []byte) bool {
	return len(credentialID) > 0 && len(credentialID) <= maxPasskeyCredentialSize &&
		len(publicKey) > 0 && len(publicKey) <= maxPasskeyPublicKeySize
}

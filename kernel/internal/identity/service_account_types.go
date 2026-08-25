package identity

import (
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
	serviceCredentialPrefix      = "ox_sa_"
	serviceCredentialSecretBytes = 32
	maxServiceAccountNameRunes   = 128
)

// ServiceAccountID is the stable UUIDv7 identity of one non-human service principal.
type ServiceAccountID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id ServiceAccountID) Valid() bool { return validUUIDv7(string(id)) }

// ServiceAccountBinding is the exact tenant/organization ceiling carried by one
// service account. It is an authentication scope constraint, not authorization
// proof; current RBAC still evaluates every governed capability separately.
type ServiceAccountBinding struct {
	tenantID       string
	organizationID string
}

// NewServiceAccountBinding creates one exact tenant or organization binding from
// opaque UUIDv7 references. A service account is always tenant-bound.
func NewServiceAccountBinding(tenantID, organizationID string) (ServiceAccountBinding, error) {
	binding := ServiceAccountBinding{
		tenantID:       strings.TrimSpace(tenantID),
		organizationID: strings.TrimSpace(organizationID),
	}
	if !binding.Valid() {
		return ServiceAccountBinding{}, serviceAccountBindingFailure()
	}
	return binding, nil
}

// TenantID returns the opaque tenant reference. It never grants tenant authority.
func (binding ServiceAccountBinding) TenantID() string { return binding.tenantID }

// OrganizationID returns the exact optional organization reference.
func (binding ServiceAccountBinding) OrganizationID() string { return binding.organizationID }

// Valid reports whether the binding is a canonical exact tenant/org scope.
func (binding ServiceAccountBinding) Valid() bool {
	return validUUIDv7(binding.tenantID) &&
		(binding.organizationID == "" || validUUIDv7(binding.organizationID))
}

// Equal reports exact binding equality. Tenant scope never implies organization scope.
func (binding ServiceAccountBinding) Equal(other ServiceAccountBinding) bool {
	return binding.Valid() && other.Valid() &&
		binding.tenantID == other.tenantID &&
		binding.organizationID == other.organizationID
}

// ServiceAccount is the kernel.identity-owned non-human principal. Its binding
// limits authentication context but carries no role/permission authority itself.
type ServiceAccount struct {
	id        ServiceAccountID
	principal PrincipalType
	state     LifecycleState
	name      string
	binding   ServiceAccountBinding
	createdAt time.Time
	updatedAt time.Time
}

// SafeServiceAccount is suitable for inventory/diagnostics and contains no credential material.
type SafeServiceAccount struct {
	ID            ServiceAccountID
	PrincipalType PrincipalType
	State         LifecycleState
	Name          string
	TenantID      string
	OrganizationID string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewServiceAccount creates one provisioned non-human principal. It does not issue
// a credential and does not create a fake human User record.
func NewServiceAccount(name string, binding ServiceAccountBinding) (ServiceAccount, error) {
	identifier, err := uuid.NewV7()
	if err != nil {
		return ServiceAccount{}, identifierFailure(err)
	}
	return newServiceAccountAt(ServiceAccountID(identifier.String()), name, binding, time.Now().UTC())
}

func newServiceAccountAt(id ServiceAccountID, name string, binding ServiceAccountBinding, createdAt time.Time) (ServiceAccount, error) {
	normalizedName, err := normalizeServiceAccountName(name)
	if err != nil || !id.Valid() || !binding.Valid() || createdAt.IsZero() {
		return ServiceAccount{}, serviceAccountInvalidFailure()
	}
	instant := createdAt.UTC()
	return ServiceAccount{
		id:        id,
		principal: PrincipalTypeServiceAccount,
		state:     LifecycleProvisioned,
		name:      normalizedName,
		binding:   binding,
		createdAt: instant,
		updatedAt: instant,
	}, nil
}

func rehydrateServiceAccount(
	id ServiceAccountID,
	principal PrincipalType,
	state LifecycleState,
	name string,
	binding ServiceAccountBinding,
	createdAt time.Time,
	updatedAt time.Time,
) (ServiceAccount, error) {
	account := ServiceAccount{
		id:        id,
		principal: principal,
		state:     state,
		name:      name,
		binding:   binding,
		createdAt: createdAt.UTC(),
		updatedAt: updatedAt.UTC(),
	}
	if account.validate() != nil {
		return ServiceAccount{}, serviceAccountStoredInvalidFailure()
	}
	return account, nil
}

func (account ServiceAccount) validate() error {
	if !account.id.Valid() || account.principal != PrincipalTypeServiceAccount || !account.state.Valid() || !account.binding.Valid() {
		return serviceAccountInvalidFailure()
	}
	if _, err := normalizeServiceAccountName(account.name); err != nil {
		return serviceAccountInvalidFailure()
	}
	if account.createdAt.IsZero() || account.updatedAt.IsZero() || account.updatedAt.Before(account.createdAt) {
		return serviceAccountInvalidFailure()
	}
	return nil
}

// ID returns the stable service-account principal identifier.
func (account ServiceAccount) ID() ServiceAccountID { return account.id }

// PrincipalType returns service_account.
func (account ServiceAccount) PrincipalType() PrincipalType { return account.principal }

// State returns the shared principal lifecycle state.
func (account ServiceAccount) State() LifecycleState { return account.state }

// Name returns classification-safe inventory metadata, not authority.
func (account ServiceAccount) Name() string { return account.name }

// Binding returns the account's exact authentication scope ceiling.
func (account ServiceAccount) Binding() ServiceAccountBinding { return account.binding }

// CreatedAt returns the UTC creation instant.
func (account ServiceAccount) CreatedAt() time.Time { return account.createdAt }

// UpdatedAt returns the UTC last lifecycle-change instant.
func (account ServiceAccount) UpdatedAt() time.Time { return account.updatedAt }

// Safe returns a credential-free inventory projection.
func (account ServiceAccount) Safe() SafeServiceAccount {
	return SafeServiceAccount{
		ID:             account.id,
		PrincipalType:  account.principal,
		State:          account.state,
		Name:           account.name,
		TenantID:       account.binding.tenantID,
		OrganizationID: account.binding.organizationID,
		CreatedAt:      account.createdAt,
		UpdatedAt:      account.updatedAt,
	}
}

// Transition applies the same canonical principal lifecycle semantics as User
// without turning the service account into a User.
func (account ServiceAccount) Transition(next LifecycleState, changedAt time.Time) (ServiceAccount, error) {
	if account.validate() != nil || !next.Valid() || !transitionAllowed(account.state, next) || changedAt.IsZero() {
		return ServiceAccount{}, serviceAccountTransitionFailure()
	}
	instant := changedAt.UTC()
	if instant.Before(account.updatedAt) {
		return ServiceAccount{}, serviceAccountTransitionFailure()
	}
	updated := account
	updated.state = next
	updated.updatedAt = instant
	return updated, nil
}

// APICredentialID identifies one rotatable service-account credential inventory record.
type APICredentialID string

// Valid reports whether id is a canonical UUIDv7 identifier.
func (id APICredentialID) Valid() bool { return validUUIDv7(string(id)) }

// APICredential contains safe lifecycle metadata only. The raw secret and its
// verifier digest are deliberately absent.
type APICredential struct {
	id               APICredentialID
	serviceAccountID ServiceAccountID
	createdAt        time.Time
	expiresAt        time.Time
	lastUsedAt       time.Time
	supersededAt     time.Time
	revokedAt        time.Time
}

func newAPICredential(
	id APICredentialID,
	serviceAccountID ServiceAccountID,
	createdAt time.Time,
	expiresAt time.Time,
) (APICredential, error) {
	credential := APICredential{
		id:               id,
		serviceAccountID: serviceAccountID,
		createdAt:        createdAt.UTC(),
		expiresAt:        expiresAt.UTC(),
	}
	if credential.validate() != nil {
		return APICredential{}, apiCredentialInvalidFailure()
	}
	return credential, nil
}

func rehydrateAPICredential(
	id APICredentialID,
	serviceAccountID ServiceAccountID,
	createdAt time.Time,
	expiresAt time.Time,
	lastUsedAt time.Time,
	supersededAt time.Time,
	revokedAt time.Time,
) (APICredential, error) {
	credential := APICredential{
		id:               id,
		serviceAccountID: serviceAccountID,
		createdAt:        createdAt.UTC(),
		expiresAt:        expiresAt.UTC(),
		lastUsedAt:       utcOrZero(lastUsedAt),
		supersededAt:     utcOrZero(supersededAt),
		revokedAt:        utcOrZero(revokedAt),
	}
	if credential.validate() != nil {
		return APICredential{}, apiCredentialStoredInvalidFailure()
	}
	return credential, nil
}

func (credential APICredential) validate() error {
	if !credential.id.Valid() || !credential.serviceAccountID.Valid() || credential.createdAt.IsZero() || credential.expiresAt.IsZero() || !credential.createdAt.Before(credential.expiresAt) {
		return apiCredentialInvalidFailure()
	}
	for _, value := range []time.Time{credential.lastUsedAt, credential.supersededAt, credential.revokedAt} {
		if !value.IsZero() && value.Before(credential.createdAt) {
			return apiCredentialInvalidFailure()
		}
	}
	return nil
}

// ID returns the safe credential inventory identifier.
func (credential APICredential) ID() APICredentialID { return credential.id }

// ServiceAccountID returns the owning non-human principal identifier.
func (credential APICredential) ServiceAccountID() ServiceAccountID { return credential.serviceAccountID }

// CreatedAt returns the issuance instant.
func (credential APICredential) CreatedAt() time.Time { return credential.createdAt }

// ExpiresAt returns the absolute expiry instant.
func (credential APICredential) ExpiresAt() time.Time { return credential.expiresAt }

// LastUsedAt returns zero until a successful verification is recorded.
func (credential APICredential) LastUsedAt() time.Time { return credential.lastUsedAt }

// SupersededAt returns zero until rotation invalidates this credential.
func (credential APICredential) SupersededAt() time.Time { return credential.supersededAt }

// RevokedAt returns zero until explicit revocation.
func (credential APICredential) RevokedAt() time.Time { return credential.revokedAt }

// Active reports local lifecycle usability before current principal/RBAC checks.
func (credential APICredential) Active(at time.Time) bool {
	if at.IsZero() || credential.validate() != nil {
		return false
	}
	instant := at.UTC()
	return instant.Before(credential.expiresAt) && credential.supersededAt.IsZero() && credential.revokedAt.IsZero()
}

// APICredentialSecret is RESTRICTED authentication material. String, GoString,
// and JSON forms are always redacted; Reveal is explicit for the one-time issuance
// or inbound authentication boundary only.
type APICredentialSecret struct {
	credentialID APICredentialID
	value        string
}

func (secret APICredentialSecret) String() string { return "[REDACTED]" }
func (secret APICredentialSecret) GoString() string { return "[REDACTED]" }
func (secret APICredentialSecret) MarshalJSON() ([]byte, error) { return json.Marshal("[REDACTED]") }
func (secret APICredentialSecret) Reveal() string { return secret.value }
func (secret APICredentialSecret) CredentialID() APICredentialID { return secret.credentialID }

func (secret APICredentialSecret) valid() bool {
	if !secret.credentialID.Valid() {
		return false
	}
	parsed, err := parseAPICredentialSecret(secret.value)
	return err == nil && parsed.credentialID == secret.credentialID
}

// ParseAPICredentialSecret identifies a credential from its opaque token while
// keeping the secret redacted in ordinary formatting/serialization.
func ParseAPICredentialSecret(value string) (APICredentialSecret, error) {
	return parseAPICredentialSecret(value)
}

func parseAPICredentialSecret(value string) (APICredentialSecret, error) {
	if !strings.HasPrefix(value, serviceCredentialPrefix) {
		return APICredentialSecret{}, apiCredentialAuthenticationFailure()
	}
	remainder := strings.TrimPrefix(value, serviceCredentialPrefix)
	identifierText, encoded, found := strings.Cut(remainder, "_")
	identifier := APICredentialID(identifierText)
	if !found || !identifier.Valid() || encoded == "" {
		return APICredentialSecret{}, apiCredentialAuthenticationFailure()
	}
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil || len(decoded) != serviceCredentialSecretBytes {
		return APICredentialSecret{}, apiCredentialAuthenticationFailure()
	}
	canonical := serviceCredentialPrefix + string(identifier) + "_" + encoded
	if canonical != value {
		return APICredentialSecret{}, apiCredentialAuthenticationFailure()
	}
	return APICredentialSecret{credentialID: identifier, value: value}, nil
}

// IssuedAPICredential is the only value that pairs safe metadata with a raw
// credential secret. Persistence APIs accept only the digest, never this value.
type IssuedAPICredential struct {
	credential APICredential
	secret     APICredentialSecret
}

func (issued IssuedAPICredential) Credential() APICredential { return issued.credential }
func (issued IssuedAPICredential) Secret() APICredentialSecret { return issued.secret }

// AuthenticatedServiceAccount proves one current API credential and exact account
// binding. It contains no permission grant; authorization remains a separate check.
type AuthenticatedServiceAccount struct {
	account    ServiceAccount
	credential APICredential
}

func (authenticated AuthenticatedServiceAccount) ServiceAccount() ServiceAccount { return authenticated.account }
func (authenticated AuthenticatedServiceAccount) Credential() APICredential { return authenticated.credential }
func (authenticated AuthenticatedServiceAccount) Valid() bool {
	return authenticated.account.validate() == nil &&
		authenticated.account.state == LifecycleActive &&
		authenticated.credential.validate() == nil &&
		authenticated.credential.serviceAccountID == authenticated.account.id
}

func generateAPICredential(
	reader io.Reader,
	serviceAccountID ServiceAccountID,
	createdAt time.Time,
	expiresAt time.Time,
) (IssuedAPICredential, [sha256.Size]byte, error) {
	if reader == nil {
		reader = rand.Reader
	}
	identifier, err := uuid.NewV7()
	if err != nil {
		return IssuedAPICredential{}, [sha256.Size]byte{}, identifierFailure(err)
	}
	credential, err := newAPICredential(APICredentialID(identifier.String()), serviceAccountID, createdAt, expiresAt)
	if err != nil {
		return IssuedAPICredential{}, [sha256.Size]byte{}, err
	}
	payload := make([]byte, serviceCredentialSecretBytes)
	if _, err = io.ReadFull(reader, payload); err != nil {
		return IssuedAPICredential{}, [sha256.Size]byte{}, apiCredentialSecretGenerationFailure(err)
	}
	value := serviceCredentialPrefix + string(credential.id) + "_" + base64.RawURLEncoding.EncodeToString(payload)
	secret := APICredentialSecret{credentialID: credential.id, value: value}
	digest := sha256.Sum256([]byte(value))
	return IssuedAPICredential{credential: credential, secret: secret}, digest, nil
}

func digestAPICredentialSecret(secret APICredentialSecret) ([sha256.Size]byte, error) {
	if !secret.valid() {
		return [sha256.Size]byte{}, apiCredentialAuthenticationFailure()
	}
	return sha256.Sum256([]byte(secret.value)), nil
}

func normalizeServiceAccountName(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || utf8.RuneCountInString(trimmed) > maxServiceAccountNameRunes || strings.ContainsAny(trimmed, "\r\n\x00") {
		return "", serviceAccountInvalidFailure()
	}
	return trimmed, nil
}

func utcOrZero(value time.Time) time.Time {
	if value.IsZero() {
		return time.Time{}
	}
	return value.UTC()
}

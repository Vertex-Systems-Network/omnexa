// Package identity implements the P02.01 human principal and User identity foundation.
package identity

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const maxPrimaryEmailRunes = 320

// PrincipalType is the canonical security principal vocabulary. P02.01 only
// instantiates HumanUser; later types remain vocabulary until their own package
// is authorized.
type PrincipalType string

const (
	PrincipalTypeHumanUser       PrincipalType = "human_user"
	PrincipalTypeServiceAccount  PrincipalType = "service_account"
	PrincipalTypeWorkload        PrincipalType = "workload"
	PrincipalTypeDevice          PrincipalType = "device"
	PrincipalTypeIntegration     PrincipalType = "integration"
	PrincipalTypeSupportOperator PrincipalType = "support_operator"
	PrincipalTypeAIAgent         PrincipalType = "ai_agent"
)

// Valid reports whether principalType belongs to the frozen minimum security
// vocabulary. Validity does not mean the type is authorized for construction by
// this package.
func (principalType PrincipalType) Valid() bool {
	switch principalType {
	case PrincipalTypeHumanUser,
		PrincipalTypeServiceAccount,
		PrincipalTypeWorkload,
		PrincipalTypeDevice,
		PrincipalTypeIntegration,
		PrincipalTypeSupportOperator,
		PrincipalTypeAIAgent:
		return true
	default:
		return false
	}
}

// LifecycleState is the transport-neutral User lifecycle vocabulary. It does
// not represent authentication/session state or authorization.
type LifecycleState string

const (
	LifecycleProvisioned LifecycleState = "provisioned"
	LifecycleActive      LifecycleState = "active"
	LifecycleSuspended   LifecycleState = "suspended"
	LifecycleDisabled    LifecycleState = "disabled"
)

// Valid reports whether state is a canonical P02.01 lifecycle state.
func (state LifecycleState) Valid() bool {
	switch state {
	case LifecycleProvisioned, LifecycleActive, LifecycleSuspended, LifecycleDisabled:
		return true
	default:
		return false
	}
}

// UserID is the stable UUIDv7 identity of a human User principal.
type UserID string

// Valid reports whether id is a standards-compliant UUIDv7 identifier.
func (id UserID) Valid() bool {
	parsed, err := uuid.Parse(string(id))
	return err == nil && parsed.Version() == 7
}

// User is the immutable P02.01 authentication identity record. It is explicitly
// not a business Person and carries no tenant, credential, session or permission
// authority.
type User struct {
	id           UserID
	principal    PrincipalType
	state        LifecycleState
	primaryEmail string
	createdAt    time.Time
	updatedAt    time.Time
}

// SafeUser is a classification-safe projection suitable for ordinary diagnostics.
// The CONFIDENTIAL/PII primary email is intentionally absent.
type SafeUser struct {
	ID            UserID
	PrincipalType PrincipalType
	State         LifecycleState
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewUser creates one provisioned human User with a UUIDv7 identity and UTC
// timestamps. primaryEmail is a CONFIDENTIAL/PII identity attribute only; it is
// not a credential, login proof or authorization signal.
func NewUser(primaryEmail string) (User, error) {
	identifier, err := uuid.NewV7()
	if err != nil {
		return User{}, identifierFailure(err)
	}
	return newUserAt(UserID(identifier.String()), primaryEmail, time.Now().UTC())
}

func newUserAt(id UserID, primaryEmail string, createdAt time.Time) (User, error) {
	email, err := normalizePrimaryEmail(primaryEmail)
	if err != nil || !id.Valid() || createdAt.IsZero() {
		return User{}, invalidUserFailure()
	}
	instant := createdAt.UTC()
	return User{
		id:           id,
		principal:    PrincipalTypeHumanUser,
		state:        LifecycleProvisioned,
		primaryEmail: email,
		createdAt:    instant,
		updatedAt:    instant,
	}, nil
}

func rehydrateUser(
	id UserID,
	principal PrincipalType,
	state LifecycleState,
	primaryEmail string,
	createdAt time.Time,
	updatedAt time.Time,
) (User, error) {
	email, err := normalizePrimaryEmail(primaryEmail)
	if err != nil {
		return User{}, invalidStoredUserFailure()
	}
	user := User{
		id:           id,
		principal:    principal,
		state:        state,
		primaryEmail: email,
		createdAt:    createdAt.UTC(),
		updatedAt:    updatedAt.UTC(),
	}
	if err := user.validate(); err != nil {
		return User{}, invalidStoredUserFailure()
	}
	return user, nil
}

// ID returns the stable UUIDv7 User/principal identifier.
func (user User) ID() UserID { return user.id }

// PrincipalType returns human_user. Other principal types are not constructible
// as User records in P02.01.
func (user User) PrincipalType() PrincipalType { return user.principal }

// State returns the non-authoritative User lifecycle state.
func (user User) State() LifecycleState { return user.state }

// PrimaryEmail returns a CONFIDENTIAL/PII identity attribute. Callers must not
// place the value in ordinary logs, traces or safe error projections.
func (user User) PrimaryEmail() string { return user.primaryEmail }

// CreatedAt returns the canonical UTC creation instant.
func (user User) CreatedAt() time.Time { return user.createdAt }

// UpdatedAt returns the canonical UTC last lifecycle-change instant.
func (user User) UpdatedAt() time.Time { return user.updatedAt }

// Safe returns an ordinary-diagnostic projection with no CONFIDENTIAL/PII field.
func (user User) Safe() SafeUser {
	return SafeUser{
		ID:            user.id,
		PrincipalType: user.principal,
		State:         user.state,
		CreatedAt:     user.createdAt,
		UpdatedAt:     user.updatedAt,
	}
}

// Transition returns a new User after validating one explicit lifecycle change.
// Same-state transitions and resurrection from disabled fail closed.
func (user User) Transition(next LifecycleState, changedAt time.Time) (User, error) {
	if err := user.validate(); err != nil || !next.Valid() || !transitionAllowed(user.state, next) {
		return User{}, transitionFailure()
	}
	instant := changedAt.UTC()
	if changedAt.IsZero() || instant.Before(user.updatedAt) {
		return User{}, transitionFailure()
	}
	updated := user
	updated.state = next
	updated.updatedAt = instant
	return updated, nil
}

func (user User) validate() error {
	if !user.id.Valid() || user.principal != PrincipalTypeHumanUser || !user.state.Valid() {
		return invalidUserFailure()
	}
	if _, err := normalizePrimaryEmail(user.primaryEmail); err != nil {
		return invalidUserFailure()
	}
	if user.createdAt.IsZero() || user.updatedAt.IsZero() || user.updatedAt.Before(user.createdAt) {
		return invalidUserFailure()
	}
	return nil
}

func transitionAllowed(from, to LifecycleState) bool {
	switch from {
	case LifecycleProvisioned:
		return to == LifecycleActive || to == LifecycleDisabled
	case LifecycleActive:
		return to == LifecycleSuspended || to == LifecycleDisabled
	case LifecycleSuspended:
		return to == LifecycleActive || to == LifecycleDisabled
	case LifecycleDisabled:
		return false
	default:
		return false
	}
}

func normalizePrimaryEmail(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || utf8.RuneCountInString(trimmed) > maxPrimaryEmailRunes || strings.ContainsAny(trimmed, " \t\r\n\x00") {
		return "", invalidUserFailure()
	}
	local, domain, found := strings.Cut(trimmed, "@")
	if !found || local == "" || domain == "" || strings.Contains(domain, "@") {
		return "", invalidUserFailure()
	}
	return strings.ToLower(trimmed), nil
}

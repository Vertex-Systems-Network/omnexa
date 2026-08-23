package identity

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/google/uuid"
)

const fixedUserID UserID = "01890f3e-7b9a-7cc0-98c4-dc0c0c07398f"

func TestNewUserUsesUUIDv7UTCAndSafeProjectionOmitsPII(t *testing.T) {
	user, err := NewUser(" Alice@Example.COM ")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}
	parsed, err := uuid.Parse(string(user.ID()))
	if err != nil || parsed.Version() != 7 {
		t.Fatalf("NewUser() ID = %q, want UUIDv7", user.ID())
	}
	if user.PrincipalType() != PrincipalTypeHumanUser || user.State() != LifecycleProvisioned {
		t.Fatalf("NewUser() principal/state = %q/%q", user.PrincipalType(), user.State())
	}
	if user.PrimaryEmail() != "alice@example.com" {
		t.Fatalf("NewUser() normalized email = %q", user.PrimaryEmail())
	}
	if user.CreatedAt().Location() != time.UTC || user.UpdatedAt().Location() != time.UTC {
		t.Fatalf("NewUser() timestamps must be UTC")
	}

	encoded, err := json.Marshal(user.Safe())
	if err != nil {
		t.Fatalf("json.Marshal(Safe()) error = %v", err)
	}
	payload := strings.ToLower(string(encoded))
	if strings.Contains(payload, "alice@example.com") || strings.Contains(payload, "primary_email") {
		t.Fatalf("Safe() leaked CONFIDENTIAL/PII identity data: %s", payload)
	}
}

func TestUserLifecycleTransitionsAreDeterministicAndFailClosed(t *testing.T) {
	createdAt := time.Date(2026, time.August, 23, 10, 0, 0, 0, time.UTC)
	user, err := newUserAt(fixedUserID, "person@example.com", createdAt)
	if err != nil {
		t.Fatalf("newUserAt() error = %v", err)
	}

	active, err := user.Transition(LifecycleActive, createdAt.Add(time.Minute))
	if err != nil {
		t.Fatalf("provisioned -> active error = %v", err)
	}
	suspended, err := active.Transition(LifecycleSuspended, createdAt.Add(2*time.Minute))
	if err != nil {
		t.Fatalf("active -> suspended error = %v", err)
	}
	reactivated, err := suspended.Transition(LifecycleActive, createdAt.Add(3*time.Minute))
	if err != nil {
		t.Fatalf("suspended -> active error = %v", err)
	}
	disabled, err := reactivated.Transition(LifecycleDisabled, createdAt.Add(4*time.Minute))
	if err != nil {
		t.Fatalf("active -> disabled error = %v", err)
	}

	if _, err := disabled.Transition(LifecycleActive, createdAt.Add(5*time.Minute)); !failure.IsCode(err, codeTransitionInvalid) {
		t.Fatalf("disabled -> active error = %v, want %s", err, codeTransitionInvalid)
	}
	if _, err := active.Transition(LifecycleActive, createdAt.Add(2*time.Minute)); !failure.IsCode(err, codeTransitionInvalid) {
		t.Fatalf("same-state transition error = %v, want %s", err, codeTransitionInvalid)
	}
	if _, err := active.Transition(LifecycleSuspended, createdAt); !failure.IsCode(err, codeTransitionInvalid) {
		t.Fatalf("backdated transition error = %v, want %s", err, codeTransitionInvalid)
	}
}

func TestUserRejectsUnsafeIdentityDataWithoutDisclosure(t *testing.T) {
	inputs := []string{
		"",
		"not-an-email",
		"two@@example.com",
		"victim@example.com\ncredential=secret",
	}
	for _, input := range inputs {
		_, err := newUserAt(fixedUserID, input, time.Now().UTC())
		if !failure.IsCode(err, codeUserInvalid) {
			t.Fatalf("newUserAt(%q) error = %v, want %s", input, err, codeUserInvalid)
		}
		if err != nil && input != "" && strings.Contains(err.Error(), input) {
			t.Fatalf("error disclosed rejected identity input %q", input)
		}
	}
}

func TestPrincipalVocabularyDoesNotCreateFakeHumanUsers(t *testing.T) {
	principalTypes := []PrincipalType{
		PrincipalTypeHumanUser,
		PrincipalTypeServiceAccount,
		PrincipalTypeWorkload,
		PrincipalTypeDevice,
		PrincipalTypeIntegration,
		PrincipalTypeSupportOperator,
		PrincipalTypeAIAgent,
	}
	for _, principalType := range principalTypes {
		if !principalType.Valid() {
			t.Fatalf("PrincipalType(%q).Valid() = false", principalType)
		}
	}

	instant := time.Date(2026, time.August, 23, 10, 0, 0, 0, time.UTC)
	_, err := rehydrateUser(
		fixedUserID,
		PrincipalTypeServiceAccount,
		LifecycleProvisioned,
		"service@example.com",
		instant,
		instant,
	)
	if !failure.IsCode(err, codePersistenceInvalid) {
		t.Fatalf("rehydrate service account as User error = %v, want %s", err, codePersistenceInvalid)
	}
}

func TestUserIdentityRemainsSeparateFromBusinessPersonAndAuthority(t *testing.T) {
	user, err := newUserAt(
		fixedUserID,
		"identity@example.com",
		time.Date(2026, time.August, 23, 10, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("newUserAt() error = %v", err)
	}
	if user.PrincipalType() != PrincipalTypeHumanUser {
		t.Fatalf("User principal type = %q", user.PrincipalType())
	}
	if user.Safe().ID != user.ID() {
		t.Fatalf("Safe User identity changed canonical ID")
	}
}

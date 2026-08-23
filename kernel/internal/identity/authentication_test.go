package identity

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	fixedTenantContextID       = "01890f3e-7b9a-7cc0-98c4-dc0c0c073991"
	fixedOrganizationContextID = "01890f3e-7b9a-7cc0-98c4-dc0c0c073992"
)

func TestPBKDF2PasswordHasherUsesGovernedAdaptiveOneWayRepresentation(t *testing.T) {
	hasher := NewPBKDF2PasswordHasher()
	password := "Pāssword\x00🔐-synthetic"
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if !strings.HasPrefix(hash, "$pbkdf2-sha256$i=600000$") {
		t.Fatalf("Hash() does not encode governed algorithm/work factor")
	}
	if strings.Contains(hash, password) {
		t.Fatalf("Hash() disclosed plaintext password")
	}
	verified, err := hasher.Verify(password, hash)
	if err != nil || !verified {
		t.Fatalf("Verify(correct) = %v, %v", verified, err)
	}
	verified, err = hasher.Verify("wrong-synthetic-password", hash)
	if err != nil || verified {
		t.Fatalf("Verify(wrong) = %v, %v", verified, err)
	}
}

func TestOpaqueSessionSecretsRedactStringJSONAndRejectedInput(t *testing.T) {
	reader := bytes.NewReader(bytes.Repeat([]byte{0x5a}, secretBytes*2))
	generated, err := generateSecrets(reader)
	if err != nil {
		t.Fatalf("generateSecrets() error = %v", err)
	}
	accessRaw := generated.access.Reveal()
	refreshRaw := generated.refresh.Reveal()
	if accessRaw == refreshRaw || !strings.HasPrefix(accessRaw, accessSecretPrefix) || !strings.HasPrefix(refreshRaw, refreshSecretPrefix) {
		t.Fatalf("generated credential shape is invalid")
	}
	if generated.access.String() != "[REDACTED]" || generated.refresh.String() != "[REDACTED]" {
		t.Fatalf("String() disclosed a credential")
	}
	encoded, err := json.Marshal(struct {
		Access  AccessSecret
		Refresh RefreshSecret
	}{generated.access, generated.refresh})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	payload := string(encoded)
	if strings.Contains(payload, accessRaw) || strings.Contains(payload, refreshRaw) {
		t.Fatalf("JSON disclosed restricted credentials")
	}
	if _, err = ParseAccessSecret(accessRaw); err != nil {
		t.Fatalf("ParseAccessSecret(valid) error = %v", err)
	}
	rejected := "ox_at_not-a-valid-secret-material"
	if _, err = ParseAccessSecret(rejected); !failure.IsCode(err, codeSessionCredentialInvalid) || strings.Contains(err.Error(), rejected) {
		t.Fatalf("ParseAccessSecret(rejected) error = %v", err)
	}
}

func TestSessionPolicyAndContextAreBoundedAndNonAuthorizing(t *testing.T) {
	policy := DefaultSessionPolicy()
	if !policy.valid() || !(policy.AccessLifetime < policy.RefreshLifetime && policy.RefreshLifetime <= policy.SessionLifetime) {
		t.Fatalf("DefaultSessionPolicy() lifetimes are not strictly bounded")
	}
	platform, err := NewSessionContext("", "")
	if err != nil || !platform.Empty() {
		t.Fatalf("NewSessionContext(platform) = %+v, %v", platform, err)
	}
	scoped, err := NewSessionContext(fixedTenantContextID, fixedOrganizationContextID)
	if err != nil || scoped.Empty() {
		t.Fatalf("NewSessionContext(scoped) error = %v", err)
	}
	if scoped.TenantID() != fixedTenantContextID || scoped.OrganizationID() != fixedOrganizationContextID {
		t.Fatalf("NewSessionContext(scoped) changed references")
	}
	if _, err = NewSessionContext("", fixedOrganizationContextID); !failure.IsCode(err, codeSessionContextInvalid) {
		t.Fatalf("organization without tenant error = %v", err)
	}
}

func TestAuthenticationFailuresAreDisclosureSafeAndEqual(t *testing.T) {
	instant := time.Date(2026, time.August, 23, 14, 0, 0, 0, time.UTC)
	hasher := &fakePasswordHasher{}
	repository := &fakeAuthenticationRepository{
		snapshot: authenticationSnapshot{
			principalID:       fixedUserID,
			state:             LifecycleActive,
			userUpdatedAt:     instant.Add(-time.Minute),
			passwordHash:      "hash:correct-password",
			credentialVersion: 1,
		},
	}
	service, err := newAuthenticationService(repository, hasher, DefaultSessionPolicy(), nil, nil, bytes.NewReader(bytes.Repeat([]byte{1}, 256)))
	if err != nil {
		t.Fatalf("newAuthenticationService() error = %v", err)
	}

	_, wrongErr := service.AuthenticatePassword(context.Background(), fixedUserID, "wrong-password", instant)
	repository.missing = true
	_, missingErr := service.AuthenticatePassword(context.Background(), UserID(fixedTenantContextID), "wrong-password", instant)
	if !failure.IsCode(wrongErr, codeAuthenticationFailed) || !failure.IsCode(missingErr, codeAuthenticationFailed) {
		t.Fatalf("wrong/missing codes = %v / %v", wrongErr, missingErr)
	}
	if wrongErr.Error() != missingErr.Error() {
		t.Fatalf("authentication failures differ: %q vs %q", wrongErr.Error(), missingErr.Error())
	}
	if hasher.verifyCalls < 2 {
		t.Fatalf("missing credential skipped password verification boundary")
	}
}

func TestSessionContextIsReauthorizedOnIssueAccessAndRefresh(t *testing.T) {
	instant := time.Date(2026, time.August, 23, 14, 0, 0, 0, time.UTC)
	hasher := &fakePasswordHasher{}
	repository := &fakeAuthenticationRepository{
		snapshot: authenticationSnapshot{
			principalID:       fixedUserID,
			state:             LifecycleActive,
			userUpdatedAt:     instant.Add(-time.Minute),
			passwordHash:      "hash:correct-password",
			credentialVersion: 1,
		},
	}
	reauthorizer := &fakeContextReauthorizer{}
	audit := &captureSecurityAudit{}
	service, err := newAuthenticationService(
		repository,
		hasher,
		DefaultSessionPolicy(),
		reauthorizer,
		audit,
		bytes.NewReader(bytes.Repeat([]byte{0x42}, 1024)),
	)
	if err != nil {
		t.Fatalf("newAuthenticationService() error = %v", err)
	}
	authentication, err := service.AuthenticatePassword(context.Background(), fixedUserID, "correct-password", instant)
	if err != nil {
		t.Fatalf("AuthenticatePassword() error = %v", err)
	}
	sessionContext, err := NewSessionContext(fixedTenantContextID, fixedOrganizationContextID)
	if err != nil {
		t.Fatalf("NewSessionContext() error = %v", err)
	}
	issued, err := service.IssueSession(context.Background(), authentication, sessionContext, "synthetic-browser", instant.Add(time.Second))
	if err != nil {
		t.Fatalf("IssueSession() error = %v", err)
	}
	if reauthorizer.calls != 1 {
		t.Fatalf("IssueSession() context reauthorization calls = %d, want 1", reauthorizer.calls)
	}

	repository.accessRecord = repository.createdRecord
	access, err := service.ValidateAccess(context.Background(), issued.AccessSecret(), instant.Add(2*time.Second))
	if err != nil || access.Session().ID() != issued.Session().ID() {
		t.Fatalf("ValidateAccess() session/error = %q / %v", access.Session().ID(), err)
	}
	if reauthorizer.calls != 2 {
		t.Fatalf("ValidateAccess() context reauthorization calls = %d, want 2", reauthorizer.calls)
	}

	repository.refreshRecord = repository.createdRecord
	reauthorizer.reject = true
	_, err = service.RotateRefresh(context.Background(), issued.RefreshSecret(), instant.Add(3*time.Second))
	if !failure.IsCode(err, codeSessionInvalid) {
		t.Fatalf("RotateRefresh(stale context) error = %v, want %s", err, codeSessionInvalid)
	}
	if repository.rotateCalls != 0 {
		t.Fatalf("stale context reached refresh mutation")
	}
	if len(audit.events) == 0 {
		t.Fatalf("security lifecycle hooks were not emitted")
	}
	auditJSON, err := json.Marshal(audit.events)
	if err != nil {
		t.Fatalf("json.Marshal(audit) error = %v", err)
	}
	if strings.Contains(string(auditJSON), issued.AccessSecret().Reveal()) || strings.Contains(string(auditJSON), issued.RefreshSecret().Reveal()) {
		t.Fatalf("audit hook leaked restricted session material")
	}
}

type fakePasswordHasher struct{ verifyCalls int }

func (*fakePasswordHasher) Hash(password string) (string, error) {
	if !validPassword(password) {
		return "", passwordInvalidFailure()
	}
	return "hash:" + password, nil
}

func (hasher *fakePasswordHasher) Verify(password, encoded string) (bool, error) {
	hasher.verifyCalls++
	return encoded == "hash:"+password, nil
}

type fakeContextReauthorizer struct {
	calls  int
	reject bool
}

func (reauthorizer *fakeContextReauthorizer) Reauthorize(_ context.Context, userID UserID, sessionContext SessionContext) error {
	reauthorizer.calls++
	if reauthorizer.reject || userID != fixedUserID || sessionContext.Empty() {
		return errors.New("synthetic context rejected")
	}
	return nil
}

type captureSecurityAudit struct{ events []SecurityAuditEvent }

func (capture *captureSecurityAudit) RecordSecurityEvent(event SecurityAuditEvent) {
	capture.events = append(capture.events, event)
}

type fakeAuthenticationRepository struct {
	snapshot      authenticationSnapshot
	missing       bool
	createdRecord sessionRecord
	accessRecord  sessionRecord
	refreshRecord sessionRecord
	rotateCalls   int
}

func (*fakeAuthenticationRepository) EnrollPassword(context.Context, UserID, string, time.Time) (uint64, error) {
	return 1, nil
}

func (repository *fakeAuthenticationRepository) AuthenticationSnapshot(context.Context, UserID) (authenticationSnapshot, error) {
	if repository.missing {
		return authenticationSnapshot{}, credentialNotFoundFailure()
	}
	return repository.snapshot, nil
}

func (*fakeAuthenticationRepository) ChangePassword(context.Context, UserID, uint64, string, time.Time) (uint64, error) {
	return 2, nil
}

func (repository *fakeAuthenticationRepository) CreateSession(
	_ context.Context,
	record sessionRecord,
	_ [sha256.Size]byte,
	_ time.Time,
	_ [sha256.Size]byte,
	_ time.Time,
) error {
	repository.createdRecord = record
	return nil
}

func (repository *fakeAuthenticationRepository) AccessSession(context.Context, [sha256.Size]byte, time.Time) (sessionRecord, error) {
	if !repository.accessRecord.session.id.Valid() {
		return sessionRecord{}, sessionFailure()
	}
	return repository.accessRecord, nil
}

func (repository *fakeAuthenticationRepository) RefreshSession(context.Context, [sha256.Size]byte, time.Time) (sessionRecord, error) {
	if !repository.refreshRecord.session.id.Valid() {
		return sessionRecord{}, sessionFailure()
	}
	return repository.refreshRecord, nil
}

func (repository *fakeAuthenticationRepository) RotateRefresh(
	_ context.Context,
	_ SessionID,
	_ [sha256.Size]byte,
	_ [sha256.Size]byte,
	_ time.Time,
	_ [sha256.Size]byte,
	_ time.Time,
	at time.Time,
) (sessionRecord, error) {
	repository.rotateCalls++
	record := repository.refreshRecord
	record.session.refreshedAt = at
	return record, nil
}

func (repository *fakeAuthenticationRepository) RevokeSession(context.Context, UserID, SessionID, time.Time) (sessionRecord, error) {
	return repository.createdRecord, nil
}

func (*fakeAuthenticationRepository) ListSessions(context.Context, UserID) ([]Session, error) {
	return nil, nil
}

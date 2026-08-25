package identity

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestServiceAccountCredentialLifecycleAndRedaction(t *testing.T) {
	repository := newMemoryServiceAccountRepository()
	audit := &captureServiceAccountAudit{}
	service, err := newServiceAccountService(repository, audit, &serviceSequenceReader{})
	if err != nil {
		t.Fatalf("newServiceAccountService() error = %v", err)
	}
	createdAt := time.Date(2026, time.August, 25, 12, 0, 0, 0, time.UTC)
	binding := mustServiceBinding(t, "01890f3e-7b9a-7cc0-98c4-dc0c0c073991", "")
	account, err := service.Create(context.Background(), "build automation", binding, createdAt)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if account.PrincipalType() != PrincipalTypeServiceAccount || account.State() != LifecycleProvisioned {
		t.Fatalf("created service principal = %q/%q", account.PrincipalType(), account.State())
	}
	activeAt := createdAt.Add(time.Second)
	account, err = service.Transition(context.Background(), account.ID(), LifecycleProvisioned, LifecycleActive, activeAt)
	if err != nil {
		t.Fatalf("Transition(active) error = %v", err)
	}

	expiresAt := activeAt.Add(time.Hour)
	issued, err := service.IssueCredential(context.Background(), account.ID(), expiresAt, activeAt.Add(time.Second))
	if err != nil {
		t.Fatalf("IssueCredential() error = %v", err)
	}
	raw := issued.Secret().Reveal()
	if raw == "" || !strings.HasPrefix(raw, serviceCredentialPrefix) {
		t.Fatal("issued credential did not contain canonical opaque secret")
	}
	if issued.Secret().String() != "[REDACTED]" || issued.Secret().GoString() != "[REDACTED]" {
		t.Fatal("API credential formatting exposed raw secret")
	}
	encoded, _ := json.Marshal(issued.Secret())
	if string(encoded) != `"[REDACTED]"` || strings.Contains(string(encoded), raw) {
		t.Fatalf("API credential JSON was not redacted: %s", encoded)
	}
	parsed, parseErr := ParseAPICredentialSecret(raw)
	if parseErr != nil || parsed.CredentialID() != issued.Credential().ID() {
		t.Fatalf("ParseAPICredentialSecret() = %q/%v", parsed.CredentialID(), parseErr)
	}
	stored := repository.credentials[issued.Credential().ID()]
	if stored.digest != sha256.Sum256([]byte(raw)) {
		t.Fatal("stored verifier digest does not match issued credential")
	}
	if strings.Contains(strings.Join(repository.persistenceStrings(), "\n"), raw) {
		t.Fatal("raw API credential leaked into persistence projection")
	}

	verifiedAt := activeAt.Add(2 * time.Second)
	authenticated, err := service.VerifyCredential(context.Background(), parsed, binding, verifiedAt)
	if err != nil || !authenticated.Valid() || authenticated.ServiceAccount().ID() != account.ID() {
		t.Fatalf("VerifyCredential() = %+v/%v", authenticated, err)
	}
	wrongBinding := mustServiceBinding(t, "01890f3e-7b9a-7cc0-98c4-dc0c0c073992", "")
	if _, err = service.VerifyCredential(context.Background(), parsed, wrongBinding, verifiedAt.Add(time.Second)); !failure.IsCode(err, codeAPICredentialAuthentication) {
		t.Fatalf("wrong-binding VerifyCredential() error = %v", err)
	}

	rotatedAt := verifiedAt.Add(2 * time.Second)
	rotated, err := service.RotateCredential(context.Background(), authenticated, rotatedAt.Add(time.Hour), rotatedAt)
	if err != nil {
		t.Fatalf("RotateCredential() error = %v", err)
	}
	if _, err = service.VerifyCredential(context.Background(), parsed, binding, rotatedAt.Add(time.Second)); !failure.IsCode(err, codeAPICredentialAuthentication) {
		t.Fatalf("superseded credential VerifyCredential() error = %v", err)
	}
	rotatedParsed, _ := ParseAPICredentialSecret(rotated.Secret().Reveal())
	rotatedAuth, err := service.VerifyCredential(context.Background(), rotatedParsed, binding, rotatedAt.Add(time.Second))
	if err != nil {
		t.Fatalf("rotated VerifyCredential() error = %v", err)
	}
	revokedAt := rotatedAt.Add(2 * time.Second)
	if _, err = service.RevokeCredential(context.Background(), rotatedAuth, revokedAt); err != nil {
		t.Fatalf("RevokeCredential() error = %v", err)
	}
	if _, err = service.VerifyCredential(context.Background(), rotatedParsed, binding, revokedAt.Add(time.Second)); !failure.IsCode(err, codeAPICredentialAuthentication) {
		t.Fatalf("revoked credential VerifyCredential() error = %v", err)
	}

	if _, err = service.VerifyCredential(context.Background(), rotatedParsed, binding, rotated.Credential().ExpiresAt()); !failure.IsCode(err, codeAPICredentialAuthentication) {
		t.Fatalf("expired credential VerifyCredential() error = %v", err)
	}
	for _, event := range audit.events {
		blob, marshalErr := json.Marshal(event)
		if marshalErr != nil {
			t.Fatalf("audit marshal error = %v", marshalErr)
		}
		if strings.Contains(string(blob), raw) || strings.Contains(string(blob), rotated.Secret().Reveal()) {
			t.Fatal("raw API credential leaked into audit event")
		}
	}
}

func TestServiceAccountLifecycleDoesNotConstructHumanUser(t *testing.T) {
	binding := mustServiceBinding(t, "01890f3e-7b9a-7cc0-98c4-dc0c0c073991", "01890f3e-7b9a-7cc0-98c4-dc0c0c0739a1")
	account, err := newServiceAccountAt(
		ServiceAccountID("01890f3e-7b9a-7cc0-98c4-dc0c0c0739b1"),
		"deploy agent",
		binding,
		time.Unix(1_700_000_000, 0).UTC(),
	)
	if err != nil {
		t.Fatalf("newServiceAccountAt() error = %v", err)
	}
	if account.PrincipalType() != PrincipalTypeServiceAccount || account.Safe().PrincipalType == PrincipalTypeHumanUser {
		t.Fatal("service account was represented as a human User")
	}
	if _, err = account.Transition(LifecycleActive, account.CreatedAt().Add(time.Second)); err != nil {
		t.Fatalf("service-account transition error = %v", err)
	}
}

func mustServiceBinding(t *testing.T, tenantID, organizationID string) ServiceAccountBinding {
	t.Helper()
	binding, err := NewServiceAccountBinding(tenantID, organizationID)
	if err != nil {
		t.Fatalf("NewServiceAccountBinding() error = %v", err)
	}
	return binding
}

type serviceSequenceReader struct{ next byte }

func (reader *serviceSequenceReader) Read(target []byte) (int, error) {
	for index := range target {
		reader.next++
		if reader.next == 0 {
			reader.next = 1
		}
		target[index] = reader.next
	}
	return len(target), nil
}

var _ io.Reader = (*serviceSequenceReader)(nil)

type captureServiceAccountAudit struct{ events []ServiceAccountAuditEvent }

func (capture *captureServiceAccountAudit) RecordServiceAccountEvent(event ServiceAccountAuditEvent) {
	capture.events = append(capture.events, event)
}

type memoryServiceCredential struct {
	credential APICredential
	digest     [sha256.Size]byte
}

type memoryServiceAccountRepository struct {
	accounts    map[ServiceAccountID]ServiceAccount
	credentials map[APICredentialID]memoryServiceCredential
}

func newMemoryServiceAccountRepository() *memoryServiceAccountRepository {
	return &memoryServiceAccountRepository{
		accounts:    make(map[ServiceAccountID]ServiceAccount),
		credentials: make(map[APICredentialID]memoryServiceCredential),
	}
}

func (repository *memoryServiceAccountRepository) CreateServiceAccount(_ context.Context, account ServiceAccount) error {
	if _, exists := repository.accounts[account.ID()]; exists {
		return serviceAccountConflictFailure()
	}
	repository.accounts[account.ID()] = account
	return nil
}

func (repository *memoryServiceAccountRepository) GetServiceAccount(_ context.Context, id ServiceAccountID) (ServiceAccount, error) {
	account, exists := repository.accounts[id]
	if !exists {
		return ServiceAccount{}, serviceAccountNotFoundFailure()
	}
	return account, nil
}

func (repository *memoryServiceAccountRepository) TransitionServiceAccount(_ context.Context, id ServiceAccountID, from, to LifecycleState, at time.Time) (ServiceAccount, error) {
	account, exists := repository.accounts[id]
	if !exists || account.State() != from {
		return ServiceAccount{}, serviceAccountConflictFailure()
	}
	updated, err := account.Transition(to, at)
	if err != nil {
		return ServiceAccount{}, err
	}
	repository.accounts[id] = updated
	return updated, nil
}

func (repository *memoryServiceAccountRepository) CreateAPICredential(_ context.Context, credential APICredential, digest [sha256.Size]byte) error {
	account, exists := repository.accounts[credential.ServiceAccountID()]
	if !exists || account.State() != LifecycleActive {
		return apiCredentialAuthenticationFailure()
	}
	repository.credentials[credential.ID()] = memoryServiceCredential{credential: credential, digest: digest}
	return nil
}

func (repository *memoryServiceAccountRepository) AuthenticateAPICredential(_ context.Context, id APICredentialID, digest [sha256.Size]byte, at time.Time) (ServiceAccount, APICredential, error) {
	stored, exists := repository.credentials[id]
	if !exists || stored.digest != digest || !stored.credential.Active(at) {
		return ServiceAccount{}, APICredential{}, apiCredentialAuthenticationFailure()
	}
	account, exists := repository.accounts[stored.credential.ServiceAccountID()]
	if !exists || account.State() != LifecycleActive {
		return ServiceAccount{}, APICredential{}, apiCredentialAuthenticationFailure()
	}
	stored.credential.lastUsedAt = at.UTC()
	repository.credentials[id] = stored
	return account, stored.credential, nil
}

func (repository *memoryServiceAccountRepository) RotateAPICredential(_ context.Context, accountID ServiceAccountID, currentID APICredentialID, next APICredential, digest [sha256.Size]byte, at time.Time) (APICredential, error) {
	stored, exists := repository.credentials[currentID]
	if !exists || stored.credential.ServiceAccountID() != accountID || !stored.credential.Active(at) {
		return APICredential{}, apiCredentialConflictFailure()
	}
	stored.credential.supersededAt = at.UTC()
	repository.credentials[currentID] = stored
	repository.credentials[next.ID()] = memoryServiceCredential{credential: next, digest: digest}
	return stored.credential, nil
}

func (repository *memoryServiceAccountRepository) RevokeAPICredential(_ context.Context, accountID ServiceAccountID, id APICredentialID, at time.Time) (APICredential, error) {
	stored, exists := repository.credentials[id]
	if !exists || stored.credential.ServiceAccountID() != accountID || !stored.credential.Active(at) {
		return APICredential{}, apiCredentialConflictFailure()
	}
	stored.credential.revokedAt = at.UTC()
	repository.credentials[id] = stored
	return stored.credential, nil
}

func (repository *memoryServiceAccountRepository) persistenceStrings() []string {
	values := make([]string, 0, len(repository.credentials))
	for id, stored := range repository.credentials {
		values = append(values, string(id), string(stored.credential.ServiceAccountID()))
	}
	return values
}

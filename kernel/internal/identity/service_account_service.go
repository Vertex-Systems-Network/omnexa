package identity

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	"github.com/google/uuid"
)

// ServiceAccountAuditAction is a bounded secret-free P02.08 lifecycle vocabulary.
type ServiceAccountAuditAction string

const (
	ServiceAccountAuditCreated            ServiceAccountAuditAction = "identity.service_account.created"
	ServiceAccountAuditTransitioned       ServiceAccountAuditAction = "identity.service_account.transitioned"
	ServiceAccountAuditCredentialIssued   ServiceAccountAuditAction = "identity.service_account.credential_issued"
	ServiceAccountAuditCredentialVerified ServiceAccountAuditAction = "identity.service_account.credential_verified"
	ServiceAccountAuditCredentialRotated  ServiceAccountAuditAction = "identity.service_account.credential_rotated"
	ServiceAccountAuditCredentialRevoked  ServiceAccountAuditAction = "identity.service_account.credential_revoked"
)

// ServiceAccountAuditEvent intentionally has no field capable of carrying raw
// API credentials, verifier digests, authorization decisions, or business data.
type ServiceAccountAuditEvent struct {
	Action           ServiceAccountAuditAction
	ServiceAccountID ServiceAccountID
	CredentialID     APICredentialID
	Succeeded        bool
	OccurredAt       time.Time
}

// ServiceAccountAuditHook receives classification-safe lifecycle facts. Durable
// identity/permission audit product behavior remains owned by P02.10.
type ServiceAccountAuditHook interface {
	RecordServiceAccountEvent(ServiceAccountAuditEvent)
}

// ServiceAccountService owns the P02.08 non-human principal and API credential
// lifecycle. Successful credential verification proves identity and exact binding
// only; authorization must still run through kernel.authorization.
type ServiceAccountService struct {
	repository ServiceAccountRepository
	audit      ServiceAccountAuditHook
	random     io.Reader
}

// NewServiceAccountService creates the governed service-account capability.
func NewServiceAccountService(repository ServiceAccountRepository, audit ServiceAccountAuditHook) (*ServiceAccountService, error) {
	return newServiceAccountService(repository, audit, rand.Reader)
}

func newServiceAccountService(repository ServiceAccountRepository, audit ServiceAccountAuditHook, random io.Reader) (*ServiceAccountService, error) {
	if repository == nil || random == nil {
		return nil, repositoryInvalidFailure()
	}
	return &ServiceAccountService{repository: repository, audit: audit, random: random}, nil
}

// Create creates one provisioned service principal. The account is distinct from
// User and receives no credential or permission automatically.
func (service *ServiceAccountService) Create(
	ctx context.Context,
	name string,
	binding ServiceAccountBinding,
	at time.Time,
) (ServiceAccount, error) {
	if service == nil || service.repository == nil || at.IsZero() || !binding.Valid() {
		return ServiceAccount{}, serviceAccountInvalidFailure()
	}
	identifier, err := uuid.NewV7()
	if err != nil {
		return ServiceAccount{}, identifierFailure(err)
	}
	account, err := newServiceAccountAt(ServiceAccountID(identifier.String()), name, binding, at.UTC())
	if err != nil {
		return ServiceAccount{}, err
	}
	err = service.repository.CreateServiceAccount(ctx, account)
	service.recordAudit(ServiceAccountAuditCreated, account.ID(), "", err == nil, at)
	if err != nil {
		return ServiceAccount{}, err
	}
	return account, nil
}

// Transition applies one canonical principal lifecycle transition. Suspended or
// disabled service accounts fail credential authentication immediately because
// every verification re-reads current principal state.
func (service *ServiceAccountService) Transition(
	ctx context.Context,
	accountID ServiceAccountID,
	from LifecycleState,
	to LifecycleState,
	at time.Time,
) (ServiceAccount, error) {
	if service == nil || service.repository == nil || !accountID.Valid() || at.IsZero() {
		return ServiceAccount{}, serviceAccountTransitionFailure()
	}
	account, err := service.repository.TransitionServiceAccount(ctx, accountID, from, to, at.UTC())
	service.recordAudit(ServiceAccountAuditTransitioned, accountID, "", err == nil, at)
	if err != nil {
		return ServiceAccount{}, err
	}
	return account, nil
}

// IssueCredential creates one bounded opaque credential for an active service
// account. The returned raw secret is one-time presentation material; persistence
// receives only its SHA-256 digest.
func (service *ServiceAccountService) IssueCredential(
	ctx context.Context,
	accountID ServiceAccountID,
	expiresAt time.Time,
	at time.Time,
) (IssuedAPICredential, error) {
	if service == nil || service.repository == nil || !accountID.Valid() || at.IsZero() || expiresAt.IsZero() || !at.UTC().Before(expiresAt.UTC()) {
		return IssuedAPICredential{}, apiCredentialInvalidFailure()
	}
	account, err := service.repository.GetServiceAccount(ctx, accountID)
	if err != nil || account.State() != LifecycleActive {
		service.recordAudit(ServiceAccountAuditCredentialIssued, accountID, "", false, at)
		return IssuedAPICredential{}, apiCredentialAuthenticationFailure()
	}
	issued, digest, err := generateAPICredential(service.random, accountID, at.UTC(), expiresAt.UTC())
	if err != nil {
		return IssuedAPICredential{}, err
	}
	if err = service.repository.CreateAPICredential(ctx, issued.Credential(), digest); err != nil {
		service.recordAudit(ServiceAccountAuditCredentialIssued, accountID, issued.Credential().ID(), false, at)
		return IssuedAPICredential{}, err
	}
	service.recordAudit(ServiceAccountAuditCredentialIssued, accountID, issued.Credential().ID(), true, at)
	return issued, nil
}

// VerifyCredential authenticates one opaque credential against current persisted
// principal/credential state and an exact caller-requested binding. It returns no
// permission/capability grant; RBAC remains mandatory after this step.
func (service *ServiceAccountService) VerifyCredential(
	ctx context.Context,
	secret APICredentialSecret,
	requestedBinding ServiceAccountBinding,
	at time.Time,
) (AuthenticatedServiceAccount, error) {
	if service == nil || service.repository == nil || at.IsZero() || !requestedBinding.Valid() {
		return AuthenticatedServiceAccount{}, apiCredentialAuthenticationFailure()
	}
	digest, err := digestAPICredentialSecret(secret)
	if err != nil {
		service.recordAudit(ServiceAccountAuditCredentialVerified, "", "", false, at)
		return AuthenticatedServiceAccount{}, apiCredentialAuthenticationFailure()
	}
	account, credential, err := service.repository.AuthenticateAPICredential(ctx, secret.CredentialID(), digest, at.UTC())
	if err != nil || account.State() != LifecycleActive || !account.Binding().Equal(requestedBinding) || !credential.Active(at) {
		service.recordAudit(ServiceAccountAuditCredentialVerified, account.ID(), secret.CredentialID(), false, at)
		return AuthenticatedServiceAccount{}, apiCredentialAuthenticationFailure()
	}
	authenticated := AuthenticatedServiceAccount{account: account, credential: credential}
	if !authenticated.Valid() {
		service.recordAudit(ServiceAccountAuditCredentialVerified, account.ID(), credential.ID(), false, at)
		return AuthenticatedServiceAccount{}, apiCredentialAuthenticationFailure()
	}
	service.recordAudit(ServiceAccountAuditCredentialVerified, account.ID(), credential.ID(), true, at)
	return authenticated, nil
}

// RotateCredential atomically supersedes the credential represented by a current
// authentication proof and returns a new one-time secret. The old credential is
// invalid from rotatedAt onward even if its original expiry has not elapsed.
func (service *ServiceAccountService) RotateCredential(
	ctx context.Context,
	authenticated AuthenticatedServiceAccount,
	expiresAt time.Time,
	rotatedAt time.Time,
) (IssuedAPICredential, error) {
	if service == nil || service.repository == nil || rotatedAt.IsZero() || expiresAt.IsZero() ||
		!authenticated.Valid() || !authenticated.credential.Active(rotatedAt) ||
		!rotatedAt.UTC().Before(expiresAt.UTC()) {
		return IssuedAPICredential{}, apiCredentialAuthenticationFailure()
	}
	issued, digest, err := generateAPICredential(
		service.random,
		authenticated.account.ID(),
		rotatedAt.UTC(),
		expiresAt.UTC(),
	)
	if err != nil {
		return IssuedAPICredential{}, err
	}
	oldCredential, err := service.repository.RotateAPICredential(
		ctx,
		authenticated.account.ID(),
		authenticated.credential.ID(),
		issued.Credential(),
		digest,
		rotatedAt.UTC(),
	)
	service.recordAudit(ServiceAccountAuditCredentialRotated, authenticated.account.ID(), authenticated.credential.ID(), err == nil, rotatedAt)
	if err != nil || oldCredential.SupersededAt().IsZero() {
		return IssuedAPICredential{}, apiCredentialConflictFailure()
	}
	return issued, nil
}

// RevokeCredential explicitly invalidates the credential represented by a current
// authentication proof. Repeated/replayed revocation fails closed.
func (service *ServiceAccountService) RevokeCredential(
	ctx context.Context,
	authenticated AuthenticatedServiceAccount,
	revokedAt time.Time,
) (APICredential, error) {
	if service == nil || service.repository == nil || revokedAt.IsZero() || !authenticated.Valid() || !authenticated.credential.Active(revokedAt) {
		return APICredential{}, apiCredentialAuthenticationFailure()
	}
	credential, err := service.repository.RevokeAPICredential(
		ctx,
		authenticated.account.ID(),
		authenticated.credential.ID(),
		revokedAt.UTC(),
	)
	service.recordAudit(ServiceAccountAuditCredentialRevoked, authenticated.account.ID(), authenticated.credential.ID(), err == nil, revokedAt)
	if err != nil {
		return APICredential{}, err
	}
	return credential, nil
}

func (service *ServiceAccountService) recordAudit(
	action ServiceAccountAuditAction,
	accountID ServiceAccountID,
	credentialID APICredentialID,
	succeeded bool,
	at time.Time,
) {
	if service == nil || service.audit == nil || at.IsZero() {
		return
	}
	service.audit.RecordServiceAccountEvent(ServiceAccountAuditEvent{
		Action:           action,
		ServiceAccountID: accountID,
		CredentialID:     credentialID,
		Succeeded:        succeeded,
		OccurredAt:       at.UTC(),
	})
}

package identity

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"time"
)

const (
	strongMethodPasskey      = "passkey"
	strongMethodRecoveryCode = "recovery_code"
)

// StrongAuthenticationService owns P02.07 human-user strong-authentication state.
// Its outputs prove authentication strength only; authorization remains an explicit
// P02.05/P02.06 responsibility at the protected capability boundary.
type StrongAuthenticationService struct {
	repository StrongAuthenticationRepository
	verifier   PasskeyVerifier
	policy     StrongAuthenticationPolicy
	audit      SecurityAuditHook
	random     io.Reader
}

// NewStrongAuthenticationService constructs the governed P02.07 service.
func NewStrongAuthenticationService(
	repository StrongAuthenticationRepository,
	verifier PasskeyVerifier,
	policy StrongAuthenticationPolicy,
	audit SecurityAuditHook,
) (*StrongAuthenticationService, error) {
	return newStrongAuthenticationService(repository, verifier, policy, audit, rand.Reader)
}

func newStrongAuthenticationService(
	repository StrongAuthenticationRepository,
	verifier PasskeyVerifier,
	policy StrongAuthenticationPolicy,
	audit SecurityAuditHook,
	random io.Reader,
) (*StrongAuthenticationService, error) {
	if repository == nil || verifier == nil || !policy.valid() || random == nil {
		return nil, repositoryInvalidFailure()
	}
	return &StrongAuthenticationService{
		repository: repository,
		verifier:   verifier,
		policy:     policy,
		audit:      audit,
		random:     random,
	}, nil
}

// BeginPasskeyEnrollment creates a pending factor and one-time challenge bound to
// the exact authenticated human-user session.
func (service *StrongAuthenticationService) BeginPasskeyEnrollment(
	ctx context.Context,
	authenticated AuthenticatedSession,
	label string,
	at time.Time,
) (PasskeyEnrollment, error) {
	session, validationErr := currentStrongAuthSession(authenticated, at)
	if service == nil || validationErr != nil || !validStrongAuthLabel(label) {
		return PasskeyEnrollment{}, strongAuthInvalidFailure()
	}
	factorValue, factorErr := generateStrongAuthID()
	if factorErr != nil {
		return PasskeyEnrollment{}, factorErr
	}
	challengeValue, challengeErr := generateStrongAuthID()
	if challengeErr != nil {
		return PasskeyEnrollment{}, challengeErr
	}
	rawChallenge, challengeDigest, generationErr := generateRestrictedValue(service.random, challengePrefix)
	if generationErr != nil {
		return PasskeyEnrollment{}, generationErr
	}
	instant := at.UTC()
	factor := StrongFactor{
		id:         FactorID(factorValue),
		principal:  session.principalID,
		factorType: StrongFactorPasskey,
		label:      label,
		state:      StrongFactorPending,
		createdAt:  instant,
	}
	challenge := strongChallengeRecord{
		id:        ChallengeID(challengeValue),
		principal: session.principalID,
		session:   session.id,
		factor:    factor.id,
		purpose:   challengePurposeEnrollment,
		digest:    challengeDigest,
		createdAt: instant,
		expiresAt: instant.Add(service.policy.ChallengeLifetime),
	}
	if persistErr := service.repository.CreatePasskeyEnrollment(ctx, factor, challenge); persistErr != nil {
		service.recordStrongAudit(SecurityAuditStrongEnrollmentStarted, session, false, at)
		return PasskeyEnrollment{}, persistErr
	}
	service.recordStrongAudit(SecurityAuditStrongEnrollmentStarted, session, true, at)
	return PasskeyEnrollment{
		factor:    factor,
		challenge: challenge.id,
		secret:    ChallengeSecret{value: rawChallenge},
		expiresAt: challenge.expiresAt,
	}, nil
}

// CompletePasskeyEnrollment consumes the exact session-bound challenge before
// invoking the approved protocol verifier, preventing concurrent replay.
func (service *StrongAuthenticationService) CompletePasskeyEnrollment(
	ctx context.Context,
	authenticated AuthenticatedSession,
	challengeID ChallengeID,
	challengeSecret ChallengeSecret,
	response PasskeyResponse,
	at time.Time,
) (StrongAuthentication, error) {
	session, validationErr := currentStrongAuthSession(authenticated, at)
	if service == nil || validationErr != nil || !challengeID.Valid() || !response.valid() {
		return StrongAuthentication{}, strongAuthChallengeFailure()
	}
	digest, digestErr := digestChallenge(challengeSecret)
	if digestErr != nil {
		service.recordStrongAudit(SecurityAuditStrongChallengeFailed, session, false, at)
		return StrongAuthentication{}, strongAuthChallengeFailure()
	}
	factor, consumeErr := service.repository.ConsumeEnrollmentChallenge(
		ctx, session.principalID, session.id, challengeID, digest, at.UTC(),
	)
	if consumeErr != nil {
		service.recordStrongAudit(SecurityAuditStrongChallengeFailed, session, false, at)
		return StrongAuthentication{}, strongAuthChallengeFailure()
	}
	verified, verifyErr := service.verifier.VerifyRegistration(ctx, PasskeyRegistrationVerification{
		principal: session.principalID,
		session:   session.id,
		challenge: challengeSecret,
		response:  response,
	})
	if verifyErr != nil || !verified.valid() {
		service.recordStrongAudit(SecurityAuditStrongEnrollmentVerified, session, false, at)
		return StrongAuthentication{}, passkeyVerificationFailure()
	}
	activated, activateErr := service.repository.ActivatePasskey(ctx, session.principalID, session.id, factor.id, verified, at.UTC())
	if activateErr != nil {
		service.recordStrongAudit(SecurityAuditStrongEnrollmentVerified, session, false, at)
		return StrongAuthentication{}, activateErr
	}
	service.recordStrongAudit(SecurityAuditStrongEnrollmentVerified, session, true, at)
	return StrongAuthentication{
		principal:       session.principalID,
		session:         session.id,
		method:          strongMethodPasskey,
		factor:          activated.id,
		authenticatedAt: at.UTC(),
	}, nil
}

// BeginPasskeyAssertion creates a one-time assertion challenge for one active factor.
func (service *StrongAuthenticationService) BeginPasskeyAssertion(
	ctx context.Context,
	authenticated AuthenticatedSession,
	factorID FactorID,
	at time.Time,
) (PasskeyAssertionChallenge, error) {
	session, validationErr := currentStrongAuthSession(authenticated, at)
	if service == nil || validationErr != nil || !factorID.Valid() {
		return PasskeyAssertionChallenge{}, strongAuthInvalidFailure()
	}
	challengeValue, challengeErr := generateStrongAuthID()
	if challengeErr != nil {
		return PasskeyAssertionChallenge{}, challengeErr
	}
	rawChallenge, challengeDigest, generationErr := generateRestrictedValue(service.random, challengePrefix)
	if generationErr != nil {
		return PasskeyAssertionChallenge{}, generationErr
	}
	instant := at.UTC()
	challenge := strongChallengeRecord{
		id:        ChallengeID(challengeValue),
		principal: session.principalID,
		session:   session.id,
		factor:    factorID,
		purpose:   challengePurposeAssertion,
		digest:    challengeDigest,
		createdAt: instant,
		expiresAt: instant.Add(service.policy.ChallengeLifetime),
	}
	factor, persistErr := service.repository.CreateAssertionChallenge(
		ctx, session.principalID, session.id, factorID, challenge, instant,
	)
	if persistErr != nil {
		return PasskeyAssertionChallenge{}, persistErr
	}
	return PasskeyAssertionChallenge{
		factor:    factor,
		challenge: challenge.id,
		secret:    ChallengeSecret{value: rawChallenge},
		expiresAt: challenge.expiresAt,
	}, nil
}

// CompletePasskeyAssertion consumes the exact challenge and validates verifier
// counter evidence before returning short-lived, non-authorizing step-up proof.
func (service *StrongAuthenticationService) CompletePasskeyAssertion(
	ctx context.Context,
	authenticated AuthenticatedSession,
	factorID FactorID,
	challengeID ChallengeID,
	challengeSecret ChallengeSecret,
	response PasskeyResponse,
	at time.Time,
) (StrongAuthentication, error) {
	session, validationErr := currentStrongAuthSession(authenticated, at)
	if service == nil || validationErr != nil || !factorID.Valid() || !challengeID.Valid() || !response.valid() {
		return StrongAuthentication{}, strongAuthChallengeFailure()
	}
	digest, digestErr := digestChallenge(challengeSecret)
	if digestErr != nil {
		service.recordStrongAudit(SecurityAuditStrongChallengeFailed, session, false, at)
		return StrongAuthentication{}, strongAuthChallengeFailure()
	}
	factor, credential, consumeErr := service.repository.ConsumeAssertionChallenge(
		ctx, session.principalID, session.id, challengeID, factorID, digest, at.UTC(),
	)
	if consumeErr != nil || !factor.Active() || !credential.valid() {
		service.recordStrongAudit(SecurityAuditStrongChallengeFailed, session, false, at)
		return StrongAuthentication{}, strongAuthChallengeFailure()
	}
	result, verifyErr := service.verifier.VerifyAssertion(ctx, PasskeyAssertionVerification{
		principal:         session.principalID,
		session:           session.id,
		challenge:         challengeSecret,
		response:          response,
		credentialID:      credential.credentialID,
		publicKey:         credential.publicKey,
		counterSupported:  credential.counterSupported,
		previousSignCount: credential.signCount,
	})
	if verifyErr != nil || result.counterSupported != credential.counterSupported ||
		(result.counterSupported && result.signCount <= credential.signCount) ||
		(!result.counterSupported && result.signCount != 0) {
		service.recordStrongAudit(SecurityAuditStrongAssertionVerified, session, false, at)
		return StrongAuthentication{}, passkeyVerificationFailure()
	}
	if counterErr := service.repository.AdvancePasskeyCounter(
		ctx,
		session.principalID,
		session.id,
		factor.id,
		credential.signCount,
		result.signCount,
		result.counterSupported,
		at.UTC(),
	); counterErr != nil {
		service.recordStrongAudit(SecurityAuditStrongAssertionVerified, session, false, at)
		return StrongAuthentication{}, counterErr
	}
	service.recordStrongAudit(SecurityAuditStrongAssertionVerified, session, true, at)
	return StrongAuthentication{
		principal:       session.principalID,
		session:         session.id,
		method:          strongMethodPasskey,
		factor:          factor.id,
		authenticatedAt: at.UTC(),
	}, nil
}

// IssueRecoveryCodes replaces prior recovery material only after recent strong
// authentication on the exact current session. Raw values are returned once only.
func (service *StrongAuthenticationService) IssueRecoveryCodes(
	ctx context.Context,
	authenticated AuthenticatedSession,
	strongAuthentication StrongAuthentication,
	at time.Time,
) (RecoveryCodeBundle, error) {
	session, validationErr := currentStrongAuthSession(authenticated, at)
	if service == nil || validationErr != nil || service.RequireStepUp(authenticated, strongAuthentication, at) != nil {
		return RecoveryCodeBundle{}, strongStepUpFailure()
	}
	setValue, setErr := generateStrongAuthID()
	if setErr != nil {
		return RecoveryCodeBundle{}, setErr
	}
	codes := make([]RecoveryCode, 0, service.policy.RecoveryCodeCount)
	digests := make([][sha256.Size]byte, 0, service.policy.RecoveryCodeCount)
	for range service.policy.RecoveryCodeCount {
		rawCode, digest, generationErr := generateRestrictedValue(service.random, recoveryPrefix)
		if generationErr != nil {
			return RecoveryCodeBundle{}, generationErr
		}
		codes = append(codes, RecoveryCode{value: rawCode})
		digests = append(digests, digest)
	}
	setID := RecoverySetID(setValue)
	if persistErr := service.repository.ReplaceRecoveryCodes(
		ctx, session.principalID, session.id, setID, digests, at.UTC(),
	); persistErr != nil {
		service.recordStrongAudit(SecurityAuditStrongRecoveryIssued, session, false, at)
		return RecoveryCodeBundle{}, persistErr
	}
	service.recordStrongAudit(SecurityAuditStrongRecoveryIssued, session, true, at)
	return RecoveryCodeBundle{setID: setID, codes: codes}, nil
}

// UseRecoveryCode atomically consumes one code and returns short-lived session-bound
// authentication-strength evidence. Replayed codes fail closed.
func (service *StrongAuthenticationService) UseRecoveryCode(
	ctx context.Context,
	authenticated AuthenticatedSession,
	code RecoveryCode,
	at time.Time,
) (StrongAuthentication, error) {
	session, validationErr := currentStrongAuthSession(authenticated, at)
	if service == nil || validationErr != nil {
		return StrongAuthentication{}, recoveryCodeFailure()
	}
	digest, digestErr := digestRecovery(code)
	if digestErr != nil {
		service.recordStrongAudit(SecurityAuditStrongRecoveryUsed, session, false, at)
		return StrongAuthentication{}, recoveryCodeFailure()
	}
	if consumeErr := service.repository.ConsumeRecoveryCode(ctx, session.principalID, session.id, digest, at.UTC()); consumeErr != nil {
		service.recordStrongAudit(SecurityAuditStrongRecoveryUsed, session, false, at)
		return StrongAuthentication{}, recoveryCodeFailure()
	}
	service.recordStrongAudit(SecurityAuditStrongRecoveryUsed, session, true, at)
	return StrongAuthentication{
		principal:       session.principalID,
		session:         session.id,
		method:          strongMethodRecoveryCode,
		authenticatedAt: at.UTC(),
	}, nil
}

// RemoveFactor is a privileged authentication-policy mutation. It requires recent
// strong authentication and applies the fixed P02.07 session-invalidation policy.
func (service *StrongAuthenticationService) RemoveFactor(
	ctx context.Context,
	authenticated AuthenticatedSession,
	strongAuthentication StrongAuthentication,
	factorID FactorID,
	at time.Time,
) (StrongFactor, error) {
	session, validationErr := currentStrongAuthSession(authenticated, at)
	if service == nil || validationErr != nil || !factorID.Valid() || service.RequireStepUp(authenticated, strongAuthentication, at) != nil {
		return StrongFactor{}, strongStepUpFailure()
	}
	factor, revokeErr := service.repository.RevokeFactor(
		ctx,
		session.principalID,
		session.id,
		factorID,
		service.policy.InvalidateSessionsOnFactorRemoval,
		at.UTC(),
	)
	service.recordStrongAudit(SecurityAuditStrongFactorRemoved, session, revokeErr == nil, at)
	if revokeErr != nil {
		return StrongFactor{}, revokeErr
	}
	return factor, nil
}

// ListFactors returns classification-safe factor inventory for one current session.
func (service *StrongAuthenticationService) ListFactors(
	ctx context.Context,
	authenticated AuthenticatedSession,
	at time.Time,
) ([]StrongFactor, error) {
	session, validationErr := currentStrongAuthSession(authenticated, at)
	if service == nil || validationErr != nil {
		return nil, strongAuthInvalidFailure()
	}
	return service.repository.ListFactors(ctx, session.principalID, session.id, at.UTC())
}

// RequireStepUp validates only authentication strength/freshness. It deliberately
// returns no permission/policy decision and cannot replace P02.05/P02.06 authorization.
func (service *StrongAuthenticationService) RequireStepUp(
	authenticated AuthenticatedSession,
	strongAuthentication StrongAuthentication,
	at time.Time,
) error {
	if service == nil || !strongAuthentication.valid() || at.IsZero() {
		return strongStepUpFailure()
	}
	session, sessionErr := currentStrongAuthSession(authenticated, at)
	if sessionErr != nil || strongAuthentication.principal != session.principalID || strongAuthentication.session != session.id {
		return strongStepUpFailure()
	}
	instant := at.UTC()
	if instant.Before(strongAuthentication.authenticatedAt) || instant.Sub(strongAuthentication.authenticatedAt) > service.policy.StepUpLifetime {
		service.recordStrongAudit(SecurityAuditStrongStepUpRequired, session, false, at)
		return strongStepUpFailure()
	}
	service.recordStrongAudit(SecurityAuditStrongStepUpRequired, session, true, at)
	return nil
}

func currentStrongAuthSession(authenticated AuthenticatedSession, at time.Time) (Session, error) {
	if at.IsZero() {
		return Session{}, strongAuthInvalidFailure()
	}
	session := authenticated.session
	if !session.id.Valid() || !session.principalID.Valid() || !session.Active(at.UTC()) {
		return Session{}, strongAuthInvalidFailure()
	}
	return session, nil
}

func (service *StrongAuthenticationService) recordStrongAudit(
	action SecurityAuditAction,
	session Session,
	succeeded bool,
	at time.Time,
) {
	if service == nil || service.audit == nil || at.IsZero() {
		return
	}
	service.audit.RecordSecurityEvent(SecurityAuditEvent{
		Action:      action,
		PrincipalID: session.principalID,
		SessionID:   session.id,
		Succeeded:   succeeded,
		OccurredAt:  at.UTC(),
	})
}

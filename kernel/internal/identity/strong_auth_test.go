package identity

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	fixedStrongSessionID      SessionID = "01890f3e-7b9a-7cc0-98c4-dc0c0c073994"
	fixedOtherStrongSessionID SessionID = "01890f3e-7b9a-7cc0-98c4-dc0c0c073995"
	fixedOtherStrongUserID    UserID    = "01890f3e-7b9a-7cc0-98c4-dc0c0c073996"
)

func TestStrongAuthenticationRestrictedMaterialRedactsAndPolicyIsBounded(t *testing.T) {
	policy := DefaultStrongAuthenticationPolicy()
	if !policy.valid() || policy.ChallengeLifetime != 5*time.Minute || policy.StepUpLifetime != 10*time.Minute || policy.RecoveryCodeCount != 8 || !policy.InvalidateSessionsOnFactorRemoval {
		t.Fatalf("DefaultStrongAuthenticationPolicy() = %+v", policy)
	}
	reader := &incrementingReader{}
	challengeRaw, _, challengeErr := generateRestrictedValue(reader, challengePrefix)
	if challengeErr != nil {
		t.Fatalf("generate challenge error = %v", challengeErr)
	}
	recoveryRaw, _, recoveryErr := generateRestrictedValue(reader, recoveryPrefix)
	if recoveryErr != nil {
		t.Fatalf("generate recovery error = %v", recoveryErr)
	}
	challenge, parseChallengeErr := ParseChallengeSecret(challengeRaw)
	if parseChallengeErr != nil {
		t.Fatalf("ParseChallengeSecret() error = %v", parseChallengeErr)
	}
	recovery, parseRecoveryErr := ParseRecoveryCode(recoveryRaw)
	if parseRecoveryErr != nil {
		t.Fatalf("ParseRecoveryCode() error = %v", parseRecoveryErr)
	}
	payload, marshalErr := json.Marshal(struct {
		Challenge ChallengeSecret
		Recovery  RecoveryCode
	}{challenge, recovery})
	if marshalErr != nil {
		t.Fatalf("json.Marshal() error = %v", marshalErr)
	}
	if challenge.String() != "[REDACTED]" || recovery.String() != "[REDACTED]" || strings.Contains(string(payload), challengeRaw) || strings.Contains(string(payload), recoveryRaw) {
		t.Fatalf("restricted strong-auth material leaked through ordinary representations")
	}
	rejected := "ox_ch_not-valid-restricted-material"
	if _, rejectedErr := ParseChallengeSecret(rejected); !failure.IsCode(rejectedErr, codeStrongChallengeInvalid) || strings.Contains(rejectedErr.Error(), rejected) {
		t.Fatalf("invalid challenge error = %v", rejectedErr)
	}
}

func TestPasskeyLifecycleChallengeBindingReplayRecoveryAndStepUpFailClosed(t *testing.T) {
	instant := time.Date(2026, time.August, 24, 2, 0, 0, 0, time.UTC)
	repository := &fakeStrongAuthenticationRepository{}
	verifier := &fakePasskeyVerifier{}
	audit := &captureSecurityAudit{}
	service, serviceErr := newStrongAuthenticationService(
		repository,
		verifier,
		DefaultStrongAuthenticationPolicy(),
		audit,
		&incrementingReader{},
	)
	if serviceErr != nil {
		t.Fatalf("newStrongAuthenticationService() error = %v", serviceErr)
	}
	authenticated := strongAuthenticatedSession(fixedUserID, fixedStrongSessionID, instant)
	otherPrincipal := strongAuthenticatedSession(fixedOtherStrongUserID, fixedStrongSessionID, instant)
	otherSession := strongAuthenticatedSession(fixedUserID, fixedOtherStrongSessionID, instant)

	enrollment, beginErr := service.BeginPasskeyEnrollment(context.Background(), authenticated, "security-key", instant)
	if beginErr != nil || enrollment.Factor().State() != StrongFactorPending || enrollment.Challenge().String() != "[REDACTED]" {
		t.Fatalf("BeginPasskeyEnrollment() factor/error = %+v/%v", enrollment.Factor(), beginErr)
	}
	response, responseErr := NewPasskeyResponse([]byte("synthetic-registration-response"))
	if responseErr != nil {
		t.Fatalf("NewPasskeyResponse() error = %v", responseErr)
	}
	if _, wrongErr := service.CompletePasskeyEnrollment(
		context.Background(), otherPrincipal, enrollment.ChallengeID(), enrollment.Challenge(), response, instant.Add(time.Second),
	); !failure.IsCode(wrongErr, codeStrongChallengeInvalid) {
		t.Fatalf("wrong-principal enrollment error = %v", wrongErr)
	}
	if _, wrongSessionErr := service.CompletePasskeyEnrollment(
		context.Background(), otherSession, enrollment.ChallengeID(), enrollment.Challenge(), response, instant.Add(time.Second),
	); !failure.IsCode(wrongSessionErr, codeStrongChallengeInvalid) {
		t.Fatalf("wrong-session enrollment error = %v", wrongSessionErr)
	}
	strongProof, completeErr := service.CompletePasskeyEnrollment(
		context.Background(), authenticated, enrollment.ChallengeID(), enrollment.Challenge(), response, instant.Add(2*time.Second),
	)
	if completeErr != nil || !strongProof.valid() || strongProof.PrincipalID() != fixedUserID || strongProof.SessionID() != fixedStrongSessionID {
		t.Fatalf("CompletePasskeyEnrollment() proof/error = %+v/%v", strongProof, completeErr)
	}
	if _, replayErr := service.CompletePasskeyEnrollment(
		context.Background(), authenticated, enrollment.ChallengeID(), enrollment.Challenge(), response, instant.Add(3*time.Second),
	); !failure.IsCode(replayErr, codeStrongChallengeInvalid) {
		t.Fatalf("replayed enrollment challenge error = %v", replayErr)
	}

	assertion, assertionErr := service.BeginPasskeyAssertion(context.Background(), authenticated, enrollment.Factor().ID(), instant.Add(4*time.Second))
	if assertionErr != nil {
		t.Fatalf("BeginPasskeyAssertion() error = %v", assertionErr)
	}
	assertionResponse, assertionResponseErr := NewPasskeyResponse([]byte("synthetic-assertion-response"))
	if assertionResponseErr != nil {
		t.Fatalf("NewPasskeyResponse(assertion) error = %v", assertionResponseErr)
	}
	verified, verifyErr := service.CompletePasskeyAssertion(
		context.Background(), authenticated, enrollment.Factor().ID(), assertion.ChallengeID(), assertion.Challenge(), assertionResponse, instant.Add(5*time.Second),
	)
	if verifyErr != nil || service.RequireStepUp(authenticated, verified, instant.Add(6*time.Second)) != nil {
		t.Fatalf("CompletePasskeyAssertion()/RequireStepUp() error = %v", verifyErr)
	}
	if _, replayErr := service.CompletePasskeyAssertion(
		context.Background(), authenticated, enrollment.Factor().ID(), assertion.ChallengeID(), assertion.Challenge(), assertionResponse, instant.Add(7*time.Second),
	); !failure.IsCode(replayErr, codeStrongChallengeInvalid) {
		t.Fatalf("replayed assertion challenge error = %v", replayErr)
	}
	if stepErr := service.RequireStepUp(otherSession, verified, instant.Add(6*time.Second)); !failure.IsCode(stepErr, codeStrongStepUpRequired) {
		t.Fatalf("wrong-session step-up error = %v", stepErr)
	}
	if staleErr := service.RequireStepUp(authenticated, verified, instant.Add(11*time.Minute)); !failure.IsCode(staleErr, codeStrongStepUpRequired) {
		t.Fatalf("stale step-up error = %v", staleErr)
	}

	bundle, recoveryIssueErr := service.IssueRecoveryCodes(context.Background(), authenticated, verified, instant.Add(8*time.Second))
	if recoveryIssueErr != nil || len(bundle.Codes()) != DefaultStrongAuthenticationPolicy().RecoveryCodeCount {
		t.Fatalf("IssueRecoveryCodes() count/error = %d/%v", len(bundle.Codes()), recoveryIssueErr)
	}
	codes := bundle.Codes()
	rawCode := codes[0].Reveal()
	bundleJSON, marshalErr := json.Marshal(bundle)
	if marshalErr != nil || strings.Contains(string(bundleJSON), rawCode) {
		t.Fatalf("recovery bundle leaked raw material")
	}
	if repository.persistedRecoveryRaw != "" {
		t.Fatalf("fake repository received raw recovery material")
	}
	recoveryProof, recoveryErr := service.UseRecoveryCode(context.Background(), authenticated, codes[0], instant.Add(9*time.Second))
	if recoveryErr != nil || service.RequireStepUp(authenticated, recoveryProof, instant.Add(10*time.Second)) != nil {
		t.Fatalf("UseRecoveryCode()/RequireStepUp() error = %v", recoveryErr)
	}
	if _, replayRecoveryErr := service.UseRecoveryCode(context.Background(), authenticated, codes[0], instant.Add(11*time.Second)); !failure.IsCode(replayRecoveryErr, codeRecoveryCodeInvalid) {
		t.Fatalf("replayed recovery code error = %v", replayRecoveryErr)
	}

	removed, removeErr := service.RemoveFactor(
		context.Background(), authenticated, recoveryProof, enrollment.Factor().ID(), instant.Add(12*time.Second),
	)
	if removeErr != nil || removed.State() != StrongFactorRevoked || !repository.sessionsInvalidated {
		t.Fatalf("RemoveFactor() state/error/invalidation = %q/%v/%v", removed.State(), removeErr, repository.sessionsInvalidated)
	}

	auditJSON, marshalAuditErr := json.Marshal(audit.events)
	if marshalAuditErr != nil {
		t.Fatalf("json.Marshal(audit) error = %v", marshalAuditErr)
	}
	auditPayload := string(auditJSON)
	for _, restricted := range []string{enrollment.Challenge().Reveal(), assertion.Challenge().Reveal(), rawCode, string(response.Reveal()), string(assertionResponse.Reveal())} {
		if strings.Contains(auditPayload, restricted) {
			t.Fatalf("audit leaked restricted authentication material")
		}
	}
}

func TestExpiredChallengeAndVerifierFailureFailClosed(t *testing.T) {
	instant := time.Date(2026, time.August, 24, 3, 0, 0, 0, time.UTC)
	repository := &fakeStrongAuthenticationRepository{}
	verifier := &fakePasskeyVerifier{rejectRegistration: true}
	service, serviceErr := newStrongAuthenticationService(repository, verifier, DefaultStrongAuthenticationPolicy(), nil, &incrementingReader{})
	if serviceErr != nil {
		t.Fatalf("newStrongAuthenticationService() error = %v", serviceErr)
	}
	authenticated := strongAuthenticatedSession(fixedUserID, fixedStrongSessionID, instant)
	enrollment, beginErr := service.BeginPasskeyEnrollment(context.Background(), authenticated, "passkey", instant)
	if beginErr != nil {
		t.Fatalf("BeginPasskeyEnrollment() error = %v", beginErr)
	}
	response, responseErr := NewPasskeyResponse([]byte("synthetic-response"))
	if responseErr != nil {
		t.Fatalf("NewPasskeyResponse() error = %v", responseErr)
	}
	if _, expiredErr := service.CompletePasskeyEnrollment(
		context.Background(), authenticated, enrollment.ChallengeID(), enrollment.Challenge(), response, instant.Add(6*time.Minute),
	); !failure.IsCode(expiredErr, codeStrongChallengeInvalid) {
		t.Fatalf("expired enrollment challenge error = %v", expiredErr)
	}

	repository = &fakeStrongAuthenticationRepository{}
	service, serviceErr = newStrongAuthenticationService(repository, verifier, DefaultStrongAuthenticationPolicy(), nil, &incrementingReader{})
	if serviceErr != nil {
		t.Fatalf("newStrongAuthenticationService(second) error = %v", serviceErr)
	}
	enrollment, beginErr = service.BeginPasskeyEnrollment(context.Background(), authenticated, "passkey", instant)
	if beginErr != nil {
		t.Fatalf("BeginPasskeyEnrollment(second) error = %v", beginErr)
	}
	if _, verifierErr := service.CompletePasskeyEnrollment(
		context.Background(), authenticated, enrollment.ChallengeID(), enrollment.Challenge(), response, instant.Add(time.Second),
	); !failure.IsCode(verifierErr, codePasskeyVerification) {
		t.Fatalf("rejected verifier error = %v", verifierErr)
	}
	if _, replayErr := service.CompletePasskeyEnrollment(
		context.Background(), authenticated, enrollment.ChallengeID(), enrollment.Challenge(), response, instant.Add(2*time.Second),
	); !failure.IsCode(replayErr, codeStrongChallengeInvalid) {
		t.Fatalf("verifier failure did not consume challenge: %v", replayErr)
	}
}

type incrementingReader struct{ next byte }

func (reader *incrementingReader) Read(target []byte) (int, error) {
	if len(target) == 0 {
		return 0, nil
	}
	for index := range target {
		reader.next++
		target[index] = reader.next
	}
	return len(target), nil
}

var _ io.Reader = (*incrementingReader)(nil)

func strongAuthenticatedSession(principal UserID, sessionID SessionID, at time.Time) AuthenticatedSession {
	return AuthenticatedSession{session: Session{
		id:          sessionID,
		principalID: principal,
		createdAt:   at.Add(-time.Minute),
		refreshedAt: at.Add(-time.Minute),
		expiresAt:   at.Add(time.Hour),
	}}
}

type fakePasskeyVerifier struct {
	rejectRegistration bool
	rejectAssertion    bool
}

func (verifier *fakePasskeyVerifier) VerifyRegistration(_ context.Context, verification PasskeyRegistrationVerification) (VerifiedPasskeyCredential, error) {
	if verifier.rejectRegistration || !verification.PrincipalID().Valid() || !verification.SessionID().Valid() || !verification.Challenge().valid() || !verification.Response().valid() {
		return VerifiedPasskeyCredential{}, errors.New("synthetic registration rejection")
	}
	return NewVerifiedPasskeyCredential([]byte("credential-id"), []byte("synthetic-public-key"), true, 1)
}

func (verifier *fakePasskeyVerifier) VerifyAssertion(_ context.Context, verification PasskeyAssertionVerification) (PasskeyAssertionResult, error) {
	if verifier.rejectAssertion || !verification.PrincipalID().Valid() || !verification.SessionID().Valid() || !verification.Challenge().valid() || !verification.Response().valid() {
		return PasskeyAssertionResult{}, errors.New("synthetic assertion rejection")
	}
	if verification.CounterSupported() {
		return NewPasskeyAssertionResult(true, verification.PreviousSignCount()+1)
	}
	return NewPasskeyAssertionResult(false, 0)
}

type fakeStrongAuthenticationRepository struct {
	factor               StrongFactor
	challenge            strongChallengeRecord
	challengeConsumed    bool
	credential           passkeyCredentialRecord
	recoveryDigests      map[[sha256.Size]byte]bool
	persistedRecoveryRaw string
	sessionsInvalidated  bool
}

func (repository *fakeStrongAuthenticationRepository) CreatePasskeyEnrollment(_ context.Context, factor StrongFactor, challenge strongChallengeRecord) error {
	repository.factor = factor
	repository.challenge = challenge
	repository.challengeConsumed = false
	return nil
}

func (repository *fakeStrongAuthenticationRepository) ConsumeEnrollmentChallenge(
	_ context.Context,
	principal UserID,
	session SessionID,
	challengeID ChallengeID,
	digest [sha256.Size]byte,
	at time.Time,
) (StrongFactor, error) {
	if repository.challengeConsumed || repository.challenge.id != challengeID || repository.challenge.principal != principal || repository.challenge.session != session || repository.challenge.digest != digest || repository.challenge.purpose != challengePurposeEnrollment || !at.Before(repository.challenge.expiresAt) {
		return StrongFactor{}, strongAuthChallengeFailure()
	}
	repository.challengeConsumed = true
	return repository.factor, nil
}

func (repository *fakeStrongAuthenticationRepository) ActivatePasskey(
	_ context.Context,
	principal UserID,
	session SessionID,
	factorID FactorID,
	credential VerifiedPasskeyCredential,
	at time.Time,
) (StrongFactor, error) {
	if principal != repository.factor.principal || !session.Valid() || factorID != repository.factor.id || repository.factor.state != StrongFactorPending {
		return StrongFactor{}, strongFactorConflictFailure()
	}
	repository.factor.state = StrongFactorActive
	repository.factor.verifiedAt = at
	repository.credential = passkeyCredentialRecord{
		factor:           factorID,
		credentialID:     credential.CredentialID(),
		publicKey:        credential.PublicKey(),
		counterSupported: credential.CounterSupported(),
		signCount:        credential.SignCount(),
	}
	return repository.factor, nil
}

func (repository *fakeStrongAuthenticationRepository) CreateAssertionChallenge(
	_ context.Context,
	principal UserID,
	session SessionID,
	factorID FactorID,
	challenge strongChallengeRecord,
	_ time.Time,
) (StrongFactor, error) {
	if principal != repository.factor.principal || !session.Valid() || factorID != repository.factor.id || !repository.factor.Active() {
		return StrongFactor{}, strongFactorNotFoundFailure()
	}
	repository.challenge = challenge
	repository.challengeConsumed = false
	return repository.factor, nil
}

func (repository *fakeStrongAuthenticationRepository) ConsumeAssertionChallenge(
	_ context.Context,
	principal UserID,
	session SessionID,
	challengeID ChallengeID,
	factorID FactorID,
	digest [sha256.Size]byte,
	at time.Time,
) (StrongFactor, passkeyCredentialRecord, error) {
	if repository.challengeConsumed || repository.challenge.id != challengeID || repository.challenge.principal != principal || repository.challenge.session != session || repository.challenge.factor != factorID || repository.challenge.digest != digest || repository.challenge.purpose != challengePurposeAssertion || !at.Before(repository.challenge.expiresAt) || !repository.factor.Active() {
		return StrongFactor{}, passkeyCredentialRecord{}, strongAuthChallengeFailure()
	}
	repository.challengeConsumed = true
	return repository.factor, repository.credential, nil
}

func (repository *fakeStrongAuthenticationRepository) AdvancePasskeyCounter(
	_ context.Context,
	principal UserID,
	session SessionID,
	factorID FactorID,
	expected uint32,
	next uint32,
	counterSupported bool,
	_ time.Time,
) error {
	if principal != repository.factor.principal || !session.Valid() || factorID != repository.factor.id || expected != repository.credential.signCount || counterSupported != repository.credential.counterSupported {
		return passkeyVerificationFailure()
	}
	repository.credential.signCount = next
	return nil
}

func (repository *fakeStrongAuthenticationRepository) ReplaceRecoveryCodes(
	_ context.Context,
	principal UserID,
	session SessionID,
	setID RecoverySetID,
	digests [][sha256.Size]byte,
	_ time.Time,
) error {
	if principal != repository.factor.principal || !session.Valid() || !setID.Valid() {
		return strongAuthInvalidFailure()
	}
	repository.recoveryDigests = make(map[[sha256.Size]byte]bool, len(digests))
	for _, digest := range digests {
		repository.recoveryDigests[digest] = true
	}
	return nil
}

func (repository *fakeStrongAuthenticationRepository) ConsumeRecoveryCode(
	_ context.Context,
	principal UserID,
	session SessionID,
	digest [sha256.Size]byte,
	_ time.Time,
) error {
	if principal != repository.factor.principal || !session.Valid() || !repository.recoveryDigests[digest] {
		return recoveryCodeFailure()
	}
	delete(repository.recoveryDigests, digest)
	return nil
}

func (repository *fakeStrongAuthenticationRepository) RevokeFactor(
	_ context.Context,
	principal UserID,
	session SessionID,
	factorID FactorID,
	invalidateSessions bool,
	at time.Time,
) (StrongFactor, error) {
	if principal != repository.factor.principal || !session.Valid() || factorID != repository.factor.id || !repository.factor.Active() {
		return StrongFactor{}, strongFactorConflictFailure()
	}
	repository.factor.state = StrongFactorRevoked
	repository.factor.revokedAt = at
	repository.sessionsInvalidated = invalidateSessions
	return repository.factor, nil
}

func (repository *fakeStrongAuthenticationRepository) ListFactors(_ context.Context, principal UserID, session SessionID, _ time.Time) ([]StrongFactor, error) {
	if principal != repository.factor.principal || !session.Valid() {
		return nil, strongAuthInvalidFailure()
	}
	return []StrongFactor{repository.factor}, nil
}

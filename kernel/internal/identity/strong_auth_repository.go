package identity

import (
	"context"
	"crypto/sha256"
	"time"
)

const (
	challengePurposeEnrollment = "passkey_enrollment"
	challengePurposeAssertion  = "passkey_assertion"
)

// StrongAuthenticationRepository is the P02.07 owner-bounded persistence contract.
// Raw challenges, recovery codes, passkey responses and authenticator private keys
// are intentionally absent.
type StrongAuthenticationRepository interface {
	CreatePasskeyEnrollment(context.Context, StrongFactor, strongChallengeRecord) error
	ConsumeEnrollmentChallenge(context.Context, UserID, SessionID, ChallengeID, [sha256.Size]byte, time.Time) (StrongFactor, error)
	ActivatePasskey(context.Context, UserID, SessionID, FactorID, VerifiedPasskeyCredential, time.Time) (StrongFactor, error)
	CreateAssertionChallenge(context.Context, UserID, SessionID, FactorID, strongChallengeRecord, time.Time) (StrongFactor, error)
	ConsumeAssertionChallenge(context.Context, UserID, SessionID, ChallengeID, FactorID, [sha256.Size]byte, time.Time) (StrongFactor, passkeyCredentialRecord, error)
	AdvancePasskeyCounter(context.Context, UserID, SessionID, FactorID, uint32, uint32, bool, time.Time) error
	ReplaceRecoveryCodes(context.Context, UserID, SessionID, RecoverySetID, [][sha256.Size]byte, time.Time) error
	ConsumeRecoveryCode(context.Context, UserID, SessionID, [sha256.Size]byte, time.Time) error
	RevokeFactor(context.Context, UserID, SessionID, FactorID, bool, time.Time) (StrongFactor, error)
	ListFactors(context.Context, UserID, SessionID, time.Time) ([]StrongFactor, error)
}

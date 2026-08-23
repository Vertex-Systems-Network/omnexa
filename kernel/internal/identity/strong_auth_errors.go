package identity

import "github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"

const (
	codeStrongAuthInvalid      failure.Code = "identity.strong_auth.invalid"
	codeStrongFactorNotFound   failure.Code = "identity.strong_auth.factor_not_found"
	codeStrongFactorConflict   failure.Code = "identity.strong_auth.factor_conflict"
	codeStrongChallengeInvalid failure.Code = "identity.strong_auth.challenge_invalid"
	codePasskeyVerification    failure.Code = "identity.strong_auth.passkey_verification_failed"
	codeRecoveryCodeInvalid    failure.Code = "identity.strong_auth.recovery_invalid"
	codeStrongStepUpRequired   failure.Code = "identity.strong_auth.step_up_required"
	codeStrongAuthSecretFailed failure.Code = "identity.strong_auth.secret_failed"
)

func strongAuthInvalidFailure() error {
	return classifiedFailure(codeStrongAuthInvalid, failure.CategoryAuthentication, "strong authentication is invalid", false)
}

func strongFactorNotFoundFailure() error {
	return classifiedFailure(codeStrongFactorNotFound, failure.CategoryNotFound, "strong authentication factor was not found", false)
}

func strongFactorConflictFailure() error {
	return classifiedFailure(codeStrongFactorConflict, failure.CategoryConflict, "strong authentication factor conflicts with current state", false)
}

func strongAuthChallengeFailure() error {
	return classifiedFailure(codeStrongChallengeInvalid, failure.CategoryAuthentication, "strong authentication challenge is invalid", false)
}

func passkeyVerificationFailure() error {
	return classifiedFailure(codePasskeyVerification, failure.CategoryAuthentication, "passkey verification failed", false)
}

func recoveryCodeFailure() error {
	return classifiedFailure(codeRecoveryCodeInvalid, failure.CategoryAuthentication, "recovery authentication failed", false)
}

func strongStepUpFailure() error {
	return classifiedFailure(codeStrongStepUpRequired, failure.CategoryAuthentication, "strong authentication is required", false)
}

func strongAuthSecretGenerationFailure(cause error) error {
	return wrappedFailure(cause, codeStrongAuthSecretFailed, failure.CategoryInternal, "strong authentication material generation failed", false)
}

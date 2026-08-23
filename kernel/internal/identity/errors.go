package identity

import (
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	codeUserInvalid        failure.Code = "identity.user.invalid"
	codeTransitionInvalid  failure.Code = "identity.user.transition_invalid"
	codeIdentifierFailed   failure.Code = "identity.identifier.failed"
	codeUserNotFound       failure.Code = "identity.user.not_found"
	codeUserConflict       failure.Code = "identity.user.conflict"
	codeRepositoryInvalid  failure.Code = "identity.repository.invalid"
	codeRepositoryFailure  failure.Code = "identity.repository.failure"
	codePersistenceInvalid failure.Code = "identity.persistence.invalid"

	codePasswordInvalid          failure.Code = "identity.password.invalid"
	codePasswordHashFailed       failure.Code = "identity.password.hash_failed"
	codePasswordHashInvalid      failure.Code = "identity.password.hash_invalid"
	codeCredentialNotFound       failure.Code = "identity.credential.not_found"
	codeCredentialConflict       failure.Code = "identity.credential.conflict"
	codeAuthenticationFailed     failure.Code = "identity.authentication.failed"
	codeAuthenticationInvalid    failure.Code = "identity.authentication.invalid"
	codeSessionInvalid           failure.Code = "identity.session.invalid"
	codeSessionCredentialInvalid failure.Code = "identity.session.credential_invalid"
	codeSessionContextInvalid    failure.Code = "identity.session.context_invalid"
	codeSessionSecretFailed      failure.Code = "identity.session.secret_failed"
	codeSessionConflict          failure.Code = "identity.session.conflict"
)

func invalidUserFailure() error {
	return classifiedFailure(codeUserInvalid, failure.CategoryValidation, "user identity is invalid", false)
}

func invalidStoredUserFailure() error {
	return classifiedFailure(codePersistenceInvalid, failure.CategoryInvariant, "stored user identity is invalid", false)
}

func transitionFailure() error {
	return classifiedFailure(codeTransitionInvalid, failure.CategoryConflict, "user lifecycle transition is invalid", false)
}

func identifierFailure(cause error) error {
	return wrappedFailure(cause, codeIdentifierFailed, failure.CategoryInternal, "identity identifier generation failed", false)
}

func userNotFoundFailure() error {
	return classifiedFailure(codeUserNotFound, failure.CategoryNotFound, "user identity was not found", false)
}

func userConflictFailure() error {
	return classifiedFailure(codeUserConflict, failure.CategoryConflict, "user identity conflicts with current state", false)
}

func repositoryInvalidFailure() error {
	return classifiedFailure(codeRepositoryInvalid, failure.CategoryValidation, "identity repository is invalid", false)
}

func passwordInvalidFailure() error {
	return classifiedFailure(codePasswordInvalid, failure.CategoryValidation, "password input is invalid", false)
}

func passwordHashFailure(cause error) error {
	return wrappedFailure(cause, codePasswordHashFailed, failure.CategoryInternal, "password hashing failed", false)
}

func passwordHashInvalidFailure() error {
	return classifiedFailure(codePasswordHashInvalid, failure.CategoryInvariant, "stored password representation is invalid", false)
}

func credentialNotFoundFailure() error {
	return classifiedFailure(codeCredentialNotFound, failure.CategoryNotFound, "authentication credential was not found", false)
}

func credentialConflictFailure() error {
	return classifiedFailure(codeCredentialConflict, failure.CategoryConflict, "authentication credential conflicts with current state", false)
}

func authenticationFailure() error {
	return classifiedFailure(codeAuthenticationFailed, failure.CategoryAuthentication, "authentication failed", false)
}

func authenticationInvalidFailure() error {
	return classifiedFailure(codeAuthenticationInvalid, failure.CategoryAuthentication, "authentication proof is no longer current", false)
}

func sessionFailure() error {
	return classifiedFailure(codeSessionInvalid, failure.CategoryAuthentication, "session is invalid", false)
}

func sessionCredentialFailure() error {
	return classifiedFailure(codeSessionCredentialInvalid, failure.CategoryAuthentication, "session credential is invalid", false)
}

func sessionContextFailure() error {
	return classifiedFailure(codeSessionContextInvalid, failure.CategoryAuthentication, "session context is invalid", false)
}

func sessionSecretGenerationFailure(cause error) error {
	return wrappedFailure(cause, codeSessionSecretFailed, failure.CategoryInternal, "session credential generation failed", false)
}

func sessionConflictFailure() error {
	return classifiedFailure(codeSessionConflict, failure.CategoryConflict, "session conflicts with current state", false)
}

func repositoryFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryTimeout, "identity repository operation timed out", true)
	}
	if errors.Is(cause, context.Canceled) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryUnavailable, "identity repository operation was canceled", false)
	}
	var postgresError *pgconn.PgError
	if errors.As(cause, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return userConflictFailure()
		case "23503", "23514":
			return wrappedFailure(cause, codePersistenceInvalid, failure.CategoryInvariant, "identity persistence invariant failed", false)
		}
	}
	return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryDependency, "identity repository operation failed", false)
}

func authenticationRepositoryFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryTimeout, "identity repository operation timed out", true)
	}
	if errors.Is(cause, context.Canceled) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryUnavailable, "identity repository operation was canceled", false)
	}
	var postgresError *pgconn.PgError
	if errors.As(cause, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return credentialConflictFailure()
		case "23503", "23514":
			return wrappedFailure(cause, codePersistenceInvalid, failure.CategoryInvariant, "authentication persistence invariant failed", false)
		}
	}
	return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryDependency, "identity repository operation failed", false)
}

func classifiedFailure(code failure.Code, category failure.Category, title string, retryable bool) error {
	value, err := failure.New(code, category, title, failure.WithRetryable(retryable))
	if err != nil {
		return errors.New("identity failure could not be classified safely")
	}
	return value
}

func wrappedFailure(cause error, code failure.Code, category failure.Category, title string, retryable bool) error {
	value, err := failure.Wrap(cause, code, category, title, failure.WithRetryable(retryable))
	if err != nil {
		return errors.New("identity failure could not be classified safely")
	}
	return value
}

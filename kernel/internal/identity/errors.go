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
	return wrappedFailure(cause, codeIdentifierFailed, failure.CategoryInternal, "user identifier generation failed", false)
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

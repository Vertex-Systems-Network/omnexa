package authorization

import (
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	codePermissionInvalid          failure.Code = "authorization.permission.invalid"
	codeRoleInvalid                failure.Code = "authorization.role.invalid"
	codeRoleNotFound               failure.Code = "authorization.role.not_found"
	codeRoleConflict               failure.Code = "authorization.role.conflict"
	codeAssignmentInvalid          failure.Code = "authorization.assignment.invalid"
	codeAssignmentNotFound         failure.Code = "authorization.assignment.not_found"
	codeAssignmentConflict         failure.Code = "authorization.assignment.conflict"
	codeSubjectInvalid             failure.Code = "authorization.subject.invalid"
	codeScopeDenied                failure.Code = "authorization.scope.denied"
	codeDenied                     failure.Code = "authorization.permission.denied"
	codeMutationMetadataInvalid    failure.Code = "authorization.mutation_metadata.invalid"
	codeIdentifierFailed           failure.Code = "authorization.identifier.failed"
	codeRepositoryInvalid          failure.Code = "authorization.repository.invalid"
	codeRepositoryFailure          failure.Code = "authorization.repository.failure"
	codePersistenceInvariantFailed failure.Code = "authorization.persistence.invariant_failed"
	codeServiceInvalid             failure.Code = "authorization.service.invalid"
)

func invalidPermissionFailure() error {
	return classifiedFailure(codePermissionInvalid, failure.CategoryValidation, "permission is invalid", false)
}

func invalidRoleFailure() error {
	return classifiedFailure(codeRoleInvalid, failure.CategoryValidation, "role is invalid", false)
}

func invalidStoredRoleFailure() error {
	return classifiedFailure(codePersistenceInvariantFailed, failure.CategoryInvariant, "stored role is invalid", false)
}

func roleNotFoundFailure() error {
	return classifiedFailure(codeRoleNotFound, failure.CategoryNotFound, "role was not found", false)
}

func roleConflictFailure() error {
	return classifiedFailure(codeRoleConflict, failure.CategoryConflict, "role conflicts with current state", false)
}

func invalidAssignmentFailure() error {
	return classifiedFailure(codeAssignmentInvalid, failure.CategoryValidation, "role assignment is invalid", false)
}

func invalidStoredAssignmentFailure() error {
	return classifiedFailure(codePersistenceInvariantFailed, failure.CategoryInvariant, "stored role assignment is invalid", false)
}

func assignmentNotFoundFailure() error {
	return classifiedFailure(codeAssignmentNotFound, failure.CategoryNotFound, "role assignment was not found", false)
}

func assignmentConflictFailure() error {
	return classifiedFailure(codeAssignmentConflict, failure.CategoryConflict, "role assignment conflicts with current state", false)
}

func invalidSubjectFailure() error {
	return classifiedFailure(codeSubjectInvalid, failure.CategoryAuthorization, "trusted authorization subject could not be established", false)
}

func scopeDeniedFailure() error {
	return classifiedFailure(codeScopeDenied, failure.CategoryAuthorization, "authorization scope is not permitted", false)
}

func deniedFailure() error {
	return classifiedFailure(codeDenied, failure.CategoryAuthorization, "authorization denied", false)
}

func invalidMutationMetadataFailure() error {
	return classifiedFailure(codeMutationMetadataInvalid, failure.CategoryValidation, "authorization mutation metadata is invalid", false)
}

func identifierFailure(cause error) error {
	return wrappedFailure(cause, codeIdentifierFailed, failure.CategoryInternal, "authorization identifier generation failed", false)
}

func repositoryInvalidFailure() error {
	return classifiedFailure(codeRepositoryInvalid, failure.CategoryValidation, "authorization repository is invalid", false)
}

func serviceInvalidFailure() error {
	return classifiedFailure(codeServiceInvalid, failure.CategoryValidation, "authorization service is invalid", false)
}

func repositoryFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryTimeout, "authorization repository operation timed out", true)
	}
	if errors.Is(cause, context.Canceled) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryUnavailable, "authorization repository operation was canceled", false)
	}
	return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryDependency, "authorization repository operation failed", false)
}

func rolePersistenceFailure(cause error) error {
	var postgresError *pgconn.PgError
	if errors.As(cause, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return roleConflictFailure()
		case "23503", "23514":
			return wrappedFailure(cause, codePersistenceInvariantFailed, failure.CategoryInvariant, "role persistence invariant failed", false)
		}
	}
	return repositoryFailure(cause)
}

func assignmentPersistenceFailure(cause error) error {
	var postgresError *pgconn.PgError
	if errors.As(cause, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return assignmentConflictFailure()
		case "23503", "23514":
			return wrappedFailure(cause, codePersistenceInvariantFailed, failure.CategoryInvariant, "role assignment persistence invariant failed", false)
		}
	}
	return repositoryFailure(cause)
}

func classifiedFailure(code failure.Code, category failure.Category, title string, retryable bool) error {
	value, constructionErr := failure.New(code, category, title, failure.WithRetryable(retryable))
	if constructionErr != nil {
		return errors.New("authorization failure could not be classified safely")
	}
	return value
}

func wrappedFailure(cause error, code failure.Code, category failure.Category, title string, retryable bool) error {
	value, constructionErr := failure.Wrap(cause, code, category, title, failure.WithRetryable(retryable))
	if constructionErr != nil {
		return errors.New("authorization failure could not be classified safely")
	}
	return value
}

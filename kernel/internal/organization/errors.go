package organization

import (
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	codeNodeInvalid                failure.Code = "organization.node.invalid"
	codeNodeNotFound               failure.Code = "organization.node.not_found"
	codeHierarchyParentInvalid     failure.Code = "organization.hierarchy.parent_invalid"
	codeHierarchyTransitionInvalid failure.Code = "organization.hierarchy.transition_invalid"
	codeHierarchyCycle             failure.Code = "organization.hierarchy.cycle"
	codeMembershipInvalid          failure.Code = "organization.membership.invalid"
	codeMembershipConflict         failure.Code = "organization.membership.conflict"
	codeContextUntrusted           failure.Code = "organization.context.untrusted"
	codeScopeDenied                failure.Code = "organization.scope.denied"
	codeIdentifierFailed           failure.Code = "organization.identifier.failed"
	codeRepositoryInvalid          failure.Code = "organization.repository.invalid"
	codeRepositoryFailure          failure.Code = "organization.repository.failure"
	codePersistenceInvariantFailed failure.Code = "organization.persistence.invariant_failed"
)

func invalidNodeFailure() error {
	return classifiedFailure(codeNodeInvalid, failure.CategoryValidation, "organization hierarchy node is invalid", false)
}

func invalidStoredNodeFailure() error {
	return classifiedFailure(codePersistenceInvariantFailed, failure.CategoryInvariant, "stored organization hierarchy node is invalid", false)
}

func nodeNotFoundFailure() error {
	return classifiedFailure(codeNodeNotFound, failure.CategoryNotFound, "organization hierarchy node was not found", false)
}

func hierarchyParentInvalidFailure() error {
	return classifiedFailure(codeHierarchyParentInvalid, failure.CategoryConflict, "organization hierarchy parent is invalid", false)
}

func hierarchyTransitionFailure() error {
	return classifiedFailure(codeHierarchyTransitionInvalid, failure.CategoryConflict, "organization hierarchy transition is invalid", false)
}

func hierarchyCycleFailure() error {
	return classifiedFailure(codeHierarchyCycle, failure.CategoryConflict, "organization hierarchy cycle is forbidden", false)
}

func invalidMembershipFailure() error {
	return classifiedFailure(codeMembershipInvalid, failure.CategoryValidation, "organization membership is invalid", false)
}

func invalidStoredMembershipFailure() error {
	return classifiedFailure(codePersistenceInvariantFailed, failure.CategoryInvariant, "stored organization membership is invalid", false)
}

func membershipConflictFailure() error {
	return classifiedFailure(codeMembershipConflict, failure.CategoryConflict, "organization membership conflicts with current state", false)
}

func contextUntrustedFailure() error {
	return classifiedFailure(codeContextUntrusted, failure.CategoryAuthorization, "organization scoped context could not be established", false)
}

func scopeDeniedFailure() error {
	return classifiedFailure(codeScopeDenied, failure.CategoryAuthorization, "organization scope does not match trusted context", false)
}

func identifierFailure(cause error) error {
	return wrappedFailure(cause, codeIdentifierFailed, failure.CategoryInternal, "organization identifier generation failed", false)
}

func repositoryInvalidFailure() error {
	return classifiedFailure(codeRepositoryInvalid, failure.CategoryValidation, "organization repository is invalid", false)
}

func repositoryFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryTimeout, "organization repository operation timed out", true)
	}
	if errors.Is(cause, context.Canceled) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryUnavailable, "organization repository operation was canceled", false)
	}
	return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryDependency, "organization repository operation failed", false)
}

func nodePersistenceFailure(cause error) error {
	var postgresError *pgconn.PgError
	if errors.As(cause, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return hierarchyTransitionFailure()
		case "23503", "23514":
			return wrappedFailure(cause, codePersistenceInvariantFailed, failure.CategoryInvariant, "organization hierarchy persistence invariant failed", false)
		}
	}
	return repositoryFailure(cause)
}

func membershipPersistenceFailure(cause error) error {
	var postgresError *pgconn.PgError
	if errors.As(cause, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return membershipConflictFailure()
		case "23503", "23514":
			return wrappedFailure(cause, codePersistenceInvariantFailed, failure.CategoryInvariant, "organization membership persistence invariant failed", false)
		}
	}
	return repositoryFailure(cause)
}

func classifiedFailure(code failure.Code, category failure.Category, title string, retryable bool) error {
	value, constructionErr := failure.New(code, category, title, failure.WithRetryable(retryable))
	if constructionErr != nil {
		return errors.New("organization failure could not be classified safely")
	}
	return value
}

func wrappedFailure(cause error, code failure.Code, category failure.Category, title string, retryable bool) error {
	value, constructionErr := failure.Wrap(cause, code, category, title, failure.WithRetryable(retryable))
	if constructionErr != nil {
		return errors.New("organization failure could not be classified safely")
	}
	return value
}

package tenancy

import (
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	codeTenantInvalid              failure.Code = "tenancy.tenant.invalid"
	codeTenantTransitionInvalid    failure.Code = "tenancy.tenant.transition_invalid"
	codeTenantNotFound             failure.Code = "tenancy.tenant.not_found"
	codeTenantConflict             failure.Code = "tenancy.tenant.conflict"
	codeMembershipInvalid          failure.Code = "tenancy.membership.invalid"
	codeMembershipTransition       failure.Code = "tenancy.membership.transition_invalid"
	codeMembershipConflict         failure.Code = "tenancy.membership.conflict"
	codeContextUntrusted           failure.Code = "tenancy.context.untrusted"
	codeCrossTenantDenied          failure.Code = "tenancy.scope.cross_tenant_denied"
	codeIdentifierFailed           failure.Code = "tenancy.identifier.failed"
	codeRepositoryInvalid          failure.Code = "tenancy.repository.invalid"
	codeRepositoryFailure          failure.Code = "tenancy.repository.failure"
	codePersistenceInvariantFailed failure.Code = "tenancy.persistence.invariant_failed"
)

func invalidTenantFailure() error {
	return classifiedFailure(codeTenantInvalid, failure.CategoryValidation, "tenant is invalid", false)
}

func invalidStoredTenantFailure() error {
	return classifiedFailure(codePersistenceInvariantFailed, failure.CategoryInvariant, "stored tenant is invalid", false)
}

func tenantTransitionFailure() error {
	return classifiedFailure(codeTenantTransitionInvalid, failure.CategoryConflict, "tenant lifecycle transition is invalid", false)
}

func tenantNotFoundFailure() error {
	return classifiedFailure(codeTenantNotFound, failure.CategoryNotFound, "tenant was not found", false)
}

func tenantConflictFailure() error {
	return classifiedFailure(codeTenantConflict, failure.CategoryConflict, "tenant conflicts with current state", false)
}

func invalidMembershipFailure() error {
	return classifiedFailure(codeMembershipInvalid, failure.CategoryValidation, "tenant membership is invalid", false)
}

func invalidStoredMembershipFailure() error {
	return classifiedFailure(codePersistenceInvariantFailed, failure.CategoryInvariant, "stored tenant membership is invalid", false)
}

func membershipTransitionFailure() error {
	return classifiedFailure(codeMembershipTransition, failure.CategoryConflict, "tenant membership transition is invalid", false)
}

func membershipConflictFailure() error {
	return classifiedFailure(codeMembershipConflict, failure.CategoryConflict, "tenant membership conflicts with current state", false)
}

func contextUntrustedFailure() error {
	return classifiedFailure(codeContextUntrusted, failure.CategoryAuthorization, "trusted tenant context could not be established", false)
}

func crossTenantDeniedFailure() error {
	return classifiedFailure(codeCrossTenantDenied, failure.CategoryAuthorization, "tenant scope does not match trusted context", false)
}

func identifierFailure(cause error) error {
	return wrappedFailure(cause, codeIdentifierFailed, failure.CategoryInternal, "tenancy identifier generation failed", false)
}

func repositoryInvalidFailure() error {
	return classifiedFailure(codeRepositoryInvalid, failure.CategoryValidation, "tenancy repository is invalid", false)
}

func repositoryFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryTimeout, "tenancy repository operation timed out", true)
	}
	if errors.Is(cause, context.Canceled) {
		return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryUnavailable, "tenancy repository operation was canceled", false)
	}
	return wrappedFailure(cause, codeRepositoryFailure, failure.CategoryDependency, "tenancy repository operation failed", false)
}

func tenantPersistenceFailure(cause error) error {
	var postgresError *pgconn.PgError
	if errors.As(cause, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return tenantConflictFailure()
		case "23503", "23514":
			return wrappedFailure(cause, codePersistenceInvariantFailed, failure.CategoryInvariant, "tenant persistence invariant failed", false)
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
			return wrappedFailure(cause, codePersistenceInvariantFailed, failure.CategoryInvariant, "tenant membership persistence invariant failed", false)
		}
	}
	return repositoryFailure(cause)
}

func classifiedFailure(code failure.Code, category failure.Category, title string, retryable bool) error {
	value, constructionErr := failure.New(code, category, title, failure.WithRetryable(retryable))
	if constructionErr != nil {
		return errors.New("tenancy failure could not be classified safely")
	}
	return value
}

func wrappedFailure(cause error, code failure.Code, category failure.Category, title string, retryable bool) error {
	value, constructionErr := failure.Wrap(cause, code, category, title, failure.WithRetryable(retryable))
	if constructionErr != nil {
		return errors.New("tenancy failure could not be classified safely")
	}
	return value
}

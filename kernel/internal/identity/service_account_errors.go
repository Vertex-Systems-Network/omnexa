package identity

import (
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	codeServiceAccountInvalid        failure.Code = "identity.service_account.invalid"
	codeServiceAccountStoredInvalid  failure.Code = "identity.service_account.stored_invalid"
	codeServiceAccountNotFound       failure.Code = "identity.service_account.not_found"
	codeServiceAccountConflict       failure.Code = "identity.service_account.conflict"
	codeServiceAccountTransition     failure.Code = "identity.service_account.transition_invalid"
	codeServiceAccountBindingInvalid failure.Code = "identity.service_account.binding_invalid"
	codeAPICredentialInvalid         failure.Code = "identity.api_credential.invalid"
	codeAPICredentialStoredInvalid   failure.Code = "identity.api_credential.stored_invalid"
	codeAPICredentialAuthentication  failure.Code = "identity.api_credential.authentication_failed"
	codeAPICredentialConflict        failure.Code = "identity.api_credential.conflict"
	codeAPICredentialSecretFailed    failure.Code = "identity.api_credential.secret_failed"
	codeServiceAccountRepository     failure.Code = "identity.service_account.repository_failure"
)

func serviceAccountInvalidFailure() error {
	return classifiedFailure(codeServiceAccountInvalid, failure.CategoryValidation, "service account is invalid", false)
}

func serviceAccountStoredInvalidFailure() error {
	return classifiedFailure(codeServiceAccountStoredInvalid, failure.CategoryInvariant, "stored service account is invalid", false)
}

func serviceAccountNotFoundFailure() error {
	return classifiedFailure(codeServiceAccountNotFound, failure.CategoryNotFound, "service account was not found", false)
}

func serviceAccountConflictFailure() error {
	return classifiedFailure(codeServiceAccountConflict, failure.CategoryConflict, "service account conflicts with current state", false)
}

func serviceAccountTransitionFailure() error {
	return classifiedFailure(codeServiceAccountTransition, failure.CategoryConflict, "service account lifecycle transition is invalid", false)
}

func serviceAccountBindingFailure() error {
	return classifiedFailure(codeServiceAccountBindingInvalid, failure.CategoryValidation, "service account binding is invalid", false)
}

func apiCredentialInvalidFailure() error {
	return classifiedFailure(codeAPICredentialInvalid, failure.CategoryValidation, "API credential is invalid", false)
}

func apiCredentialStoredInvalidFailure() error {
	return classifiedFailure(codeAPICredentialStoredInvalid, failure.CategoryInvariant, "stored API credential is invalid", false)
}

func apiCredentialAuthenticationFailure() error {
	return classifiedFailure(codeAPICredentialAuthentication, failure.CategoryAuthentication, "API credential authentication failed", false)
}

func apiCredentialConflictFailure() error {
	return classifiedFailure(codeAPICredentialConflict, failure.CategoryConflict, "API credential conflicts with current state", false)
}

func apiCredentialSecretGenerationFailure(cause error) error {
	return wrappedFailure(cause, codeAPICredentialSecretFailed, failure.CategoryInternal, "API credential generation failed", false)
}

func serviceAccountRepositoryFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) {
		return wrappedFailure(cause, codeServiceAccountRepository, failure.CategoryTimeout, "service account repository operation timed out", true)
	}
	if errors.Is(cause, context.Canceled) {
		return wrappedFailure(cause, codeServiceAccountRepository, failure.CategoryUnavailable, "service account repository operation was canceled", false)
	}
	var postgresError *pgconn.PgError
	if errors.As(cause, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return serviceAccountConflictFailure()
		case "23503", "23514":
			return wrappedFailure(cause, codeServiceAccountRepository, failure.CategoryInvariant, "service account persistence invariant failed", false)
		}
	}
	return wrappedFailure(cause, codeServiceAccountRepository, failure.CategoryDependency, "service account repository operation failed", false)
}

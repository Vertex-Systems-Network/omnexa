package authorization

import (
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeContextInvalid               failure.Code = "authorization.context.invalid"
	codeRelationshipResolutionFailed failure.Code = "authorization.relationship_resolution.failed"
	codeContextConstraintFailed      failure.Code = "authorization.context_constraint.failed"
)

func invalidContextFailure() error {
	return classifiedFailure(codeContextInvalid, failure.CategoryAuthorization, "authorization context is invalid", false)
}

func relationshipResolutionFailure(cause error) error {
	return contextualDependencyFailure(cause, codeRelationshipResolutionFailed, "authorization relationship resolution failed")
}

func constraintEvaluationFailure(cause error) error {
	return contextualDependencyFailure(cause, codeContextConstraintFailed, "authorization context constraint evaluation failed")
}

func contextualDependencyFailure(cause error, code failure.Code, title string) error {
	switch {
	case errors.Is(cause, context.DeadlineExceeded):
		return wrappedFailure(cause, code, failure.CategoryTimeout, title, true)
	case errors.Is(cause, context.Canceled):
		return wrappedFailure(cause, code, failure.CategoryUnavailable, title, false)
	default:
		return wrappedFailure(cause, code, failure.CategoryDependency, title, false)
	}
}

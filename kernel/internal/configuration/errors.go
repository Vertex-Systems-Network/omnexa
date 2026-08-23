package configuration

import (
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeRegistryInvalid     failure.Code = "configuration.registry.invalid"
	codeDefinitionInvalid   failure.Code = "configuration.definition.invalid"
	codeDefinitionDuplicate failure.Code = "configuration.definition.duplicate"
	codeDefinitionUnknown   failure.Code = "configuration.definition.unknown"
	codeContextInvalid      failure.Code = "configuration.context.invalid"
	codeEvaluatorInvalid    failure.Code = "configuration.evaluator.invalid"
	codeEvaluationCanceled  failure.Code = "configuration.evaluation.canceled"
	codeEvaluationDeadline  failure.Code = "configuration.evaluation.deadline"
)

func safeFailure(code failure.Code, title string) error {
	return classifiedFailure(code, failure.CategoryValidation, title)
}

func classifiedFailure(code failure.Code, category failure.Category, title string) error {
	value, err := failure.New(code, category, title)
	if err != nil {
		return errors.New("runtime configuration failure could not be classified safely")
	}
	return value
}

func evaluationContextFailure(cause error) error {
	code := codeEvaluationCanceled
	category := failure.CategoryUnavailable
	title := "runtime configuration evaluation was canceled"
	if errors.Is(cause, context.DeadlineExceeded) {
		code = codeEvaluationDeadline
		category = failure.CategoryTimeout
		title = "runtime configuration evaluation deadline was exceeded"
	}
	value, err := failure.Wrap(cause, code, category, title)
	if err != nil {
		return errors.New("runtime configuration evaluation failure could not be classified safely")
	}
	return value
}

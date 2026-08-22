// Package observability implements the P01.07 structured logging and OpenTelemetry baseline.
package observability

import (
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeConfigurationInvalid failure.Code = "observability.configuration.invalid"
	codeCorrelationInvalid   failure.Code = "observability.correlation.invalid"
	codeLifecycleCanceled    failure.Code = "observability.lifecycle.canceled"
	codeLifecycleTimeout     failure.Code = "observability.lifecycle.timeout"
	codeLifecycleFailed      failure.Code = "observability.lifecycle.failed"
)

func safeFailure(code failure.Code, category failure.Category, title string, options ...failure.Option) error {
	value, err := failure.New(code, category, title, options...)
	if err != nil {
		return errors.New("observability failure could not be classified safely")
	}
	return value
}

func safeWrappedFailure(cause error, code failure.Code, category failure.Category, title string, options ...failure.Option) error {
	if cause == nil {
		return safeFailure(code, category, title, options...)
	}
	value, err := failure.Wrap(cause, code, category, title, options...)
	if err != nil {
		return errors.New("observability failure could not be classified safely")
	}
	return value
}

// Package cache implements the P01.05 Redis-compatible cache abstraction.
package cache

import (
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeConfigurationInvalid failure.Code = "cache.configuration.invalid"
	codeConnectionUnavailable failure.Code = "cache.connection.unavailable"
	codeOperationFailed       failure.Code = "cache.operation.failed"
	codeKeyInvalid            failure.Code = "cache.key.invalid"
	codeValueInvalid          failure.Code = "cache.value.invalid"
	codeSerializationFailed   failure.Code = "cache.serialization.failed"
)

func safeFailure(code failure.Code, category failure.Category, title string, options ...failure.Option) error {
	value, err := failure.New(code, category, title, options...)
	if err != nil {
		return errors.New("cache failure could not be classified safely")
	}
	return value
}

func safeWrappedFailure(cause error, code failure.Code, category failure.Category, title string, options ...failure.Option) error {
	if cause == nil {
		return safeFailure(code, category, title, options...)
	}
	value, err := failure.Wrap(cause, code, category, title, options...)
	if err != nil {
		return errors.New("cache failure could not be classified safely")
	}
	return value
}

// Package storage implements the P01.06 S3-compatible object/file storage abstraction.
package storage

import (
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeConfigurationInvalid  failure.Code = "storage.configuration.invalid"
	codeConnectionCanceled    failure.Code = "storage.connection.canceled"
	codeConnectionTimeout     failure.Code = "storage.connection.timeout"
	codeConnectionUnavailable failure.Code = "storage.connection.unavailable"
	codeOperationCanceled     failure.Code = "storage.operation.canceled"
	codeOperationTimeout      failure.Code = "storage.operation.timeout"
	codeOperationFailed       failure.Code = "storage.operation.failed"
	codeObjectNotFound        failure.Code = "storage.object.not_found"
	codeKeyInvalid            failure.Code = "storage.key.invalid"
	codeMetadataInvalid       failure.Code = "storage.metadata.invalid"
	codeIntegrityFailed       failure.Code = "storage.integrity.failed"
	codeLengthInvalid         failure.Code = "storage.length.invalid"
)

func safeFailure(code failure.Code, category failure.Category, title string, options ...failure.Option) error {
	value, err := failure.New(code, category, title, options...)
	if err != nil {
		return errors.New("storage failure could not be classified safely")
	}
	return value
}

func safeWrappedFailure(cause error, code failure.Code, category failure.Category, title string, options ...failure.Option) error {
	if cause == nil {
		return safeFailure(code, category, title, options...)
	}
	value, err := failure.Wrap(cause, code, category, title, options...)
	if err != nil {
		return errors.New("storage failure could not be classified safely")
	}
	return value
}

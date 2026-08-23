package audit

import (
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeRecordInvalid      failure.Code = "audit.record.invalid"
	codeRecordProhibited   failure.Code = "audit.record.prohibited"
	codeRecordIntegrity    failure.Code = "audit.record.integrity"
	codeIdentifierFailed   failure.Code = "audit.identifier.failed"
	codeWriterInvalid      failure.Code = "audit.writer.invalid"
	codeRequirementInvalid failure.Code = "audit.requirement.invalid"
	codeSinkInvalid        failure.Code = "audit.sink.invalid"
	codeSinkFull           failure.Code = "audit.sink.full"
	codeSinkRequiredFailed failure.Code = "audit.sink.required_failed"
	codeDeliveryCanceled   failure.Code = "audit.delivery.canceled"
	codeDeliveryDeadline   failure.Code = "audit.delivery.deadline"
)

func invalidRecordFailure() error {
	return classifiedFailure(codeRecordInvalid, failure.CategoryValidation, "audit record is invalid", false)
}

func prohibitedFieldFailure() error {
	return classifiedFailure(codeRecordProhibited, failure.CategoryValidation, "audit record contains prohibited protected data", false)
}

func integrityFailure() error {
	return classifiedFailure(codeRecordIntegrity, failure.CategoryInvariant, "audit record integrity verification failed", false)
}

func classifiedFailure(code failure.Code, category failure.Category, title string, retryable bool) error {
	value, err := failure.New(code, category, title, failure.WithRetryable(retryable))
	if err != nil {
		return errors.New("audit failure could not be classified safely")
	}
	return value
}

func wrappedFailure(cause error, code failure.Code, category failure.Category, title string, retryable bool) error {
	value, err := failure.Wrap(cause, code, category, title, failure.WithRetryable(retryable))
	if err != nil {
		return errors.New("audit failure could not be classified safely")
	}
	return value
}

func contextFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) {
		return wrappedFailure(cause, codeDeliveryDeadline, failure.CategoryTimeout, "audit delivery deadline was exceeded", false)
	}
	return wrappedFailure(cause, codeDeliveryCanceled, failure.CategoryUnavailable, "audit delivery was canceled", false)
}

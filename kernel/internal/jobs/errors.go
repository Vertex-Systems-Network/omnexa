package jobs

import (
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeExecutorInvalid        failure.Code = "jobs.executor.invalid"
	codeExecutorStopping       failure.Code = "jobs.executor.stopping"
	codeQueueFull              failure.Code = "jobs.queue.full"
	codeRequestInvalid         failure.Code = "jobs.request.invalid"
	codeTypeInvalid            failure.Code = "jobs.type.invalid"
	codeTypeDuplicate          failure.Code = "jobs.type.duplicate"
	codeTypeUnknown            failure.Code = "jobs.type.unknown"
	codeRegistryFrozen         failure.Code = "jobs.registry.frozen"
	codeRegistryInvalid        failure.Code = "jobs.registry.invalid"
	codeIDGenerationFailed     failure.Code = "jobs.execution.id_generation_failed"
	codeExecutionFailed        failure.Code = "jobs.execution.failed"
	codeExecutionCanceled      failure.Code = "jobs.execution.canceled"
	codeExecutionDeadline      failure.Code = "jobs.execution.deadline"
	codeIdempotencyConflict    failure.Code = "jobs.idempotency.conflict"
	codeIdempotencyInProgress  failure.Code = "jobs.idempotency.in_progress"
	codeRetryRequiresIdempotency failure.Code = "jobs.retry.idempotency_required"
	codeScheduleInvalid        failure.Code = "jobs.schedule.invalid"
)

func classifiedFailure(code failure.Code, category failure.Category, title string, retryable bool) error {
	value, err := failure.New(code, category, title, failure.WithRetryable(retryable))
	if err != nil {
		return errors.New("job failure could not be classified safely")
	}
	return value
}

func wrappedFailure(cause error, code failure.Code, category failure.Category, title string, retryable bool) error {
	if cause == nil {
		return classifiedFailure(code, category, title, retryable)
	}
	value, err := failure.Wrap(cause, code, category, title, failure.WithRetryable(retryable))
	if err != nil {
		return errors.New("job failure could not be classified safely")
	}
	return value
}

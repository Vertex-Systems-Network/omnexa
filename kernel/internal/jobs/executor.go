package jobs

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/observability"
	"github.com/google/uuid"
)

const (
	minWorkers       = 1
	maxWorkers       = 64
	minQueueCapacity = 1
	maxQueueCapacity = 4096
)

// Settings bounds the process-local executor. It does not describe a durable
// broker, distributed consumer group, or workflow runtime.
type Settings struct {
	Workers       int
	QueueCapacity int
}

// Future is an in-memory completion handle for one queued request.
type Future struct {
	completion *completion
}

// Wait returns the repeatable terminal outcome without changing execution
// ownership. A caller that stops waiting does not silently grant permission to
// replay work.
func (future *Future) Wait(ctx context.Context) (Result, error) {
	if future == nil || future.completion == nil || future.completion.done == nil {
		return Result{}, classifiedFailure(codeRequestInvalid, failure.CategoryValidation, "job completion handle is invalid", false)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-future.completion.done:
		return future.completion.outcome.result, future.completion.outcome.err
	case <-ctx.Done():
		return Result{}, contextFailure(ctx.Err())
	}
}

// Executor provides bounded in-memory execution and queueing for kernel-local
// jobs. It intentionally has no persistence or delivery guarantee across a
// process restart.
type Executor struct {
	registry     *Registry
	logger       *observability.Logger
	settings     Settings
	queue        chan queuedRequest
	rootCtx      context.Context
	rootCancel   context.CancelFunc
	workers      sync.WaitGroup
	work         sync.WaitGroup
	mu           sync.RWMutex
	accepting    bool
	shutdown     sync.Once
	shutdownDone chan struct{}
	idempotency  memoryIdempotency
}

type queuedRequest struct {
	ctx        context.Context
	request    Request
	completion *completion
}

type completion struct {
	done    chan struct{}
	outcome outcome
}

type outcome struct {
	result Result
	err    error
}

type idempotencyEntry struct {
	fingerprint string
	completed   bool
	result      Result
}

type memoryIdempotency struct {
	mu      sync.Mutex
	entries map[string]idempotencyEntry
}

// NewExecutor freezes the supplied registry and starts a bounded worker pool.
func NewExecutor(registry *Registry, logger *observability.Logger, settings Settings) (*Executor, error) {
	if registry == nil || settings.Workers < minWorkers || settings.Workers > maxWorkers || settings.QueueCapacity < minQueueCapacity || settings.QueueCapacity > maxQueueCapacity {
		return nil, classifiedFailure(codeExecutorInvalid, failure.CategoryValidation, "job executor configuration is invalid", false)
	}
	if err := registry.freeze(); err != nil {
		return nil, err
	}

	rootCtx, rootCancel := context.WithCancel(context.Background())
	executor := &Executor{
		registry:     registry,
		logger:       logger,
		settings:     settings,
		queue:        make(chan queuedRequest, settings.QueueCapacity),
		rootCtx:      rootCtx,
		rootCancel:   rootCancel,
		accepting:    true,
		shutdownDone: make(chan struct{}),
		idempotency:  memoryIdempotency{entries: make(map[string]idempotencyEntry)},
	}
	for index := 0; index < settings.Workers; index++ {
		executor.workers.Add(1)
		go executor.worker()
	}
	return executor, nil
}

// Execute runs one request synchronously using the same registry, retry,
// idempotency, context, and observability contracts as queued work. Admission is
// lifecycle-aware: once shutdown starts, no new synchronous work is accepted.
func (executor *Executor) Execute(ctx context.Context, request Request) (Result, error) {
	if executor == nil || executor.registry == nil {
		return Result{}, classifiedFailure(codeExecutorInvalid, failure.CategoryInvariant, "job executor is invalid", false)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	executor.mu.RLock()
	if !executor.accepting {
		executor.mu.RUnlock()
		return Result{}, classifiedFailure(codeExecutorStopping, failure.CategoryUnavailable, "job executor is not accepting work", false)
	}
	executor.work.Add(1)
	executor.mu.RUnlock()
	defer executor.work.Done()

	executionCtx, cancel := linkedContext(ctx, executor.rootCtx)
	defer cancel()
	return executor.execute(executionCtx, request)
}

// execute runs work that has already passed executor admission. Queue workers
// use this path so requests accepted before shutdown can drain normally.
func (executor *Executor) execute(ctx context.Context, request Request) (Result, error) {
	policy, err := validateRequest(request)
	if err != nil {
		return Result{}, err
	}
	handler, err := executor.registry.lookup(request.Type)
	if err != nil {
		return Result{}, err
	}
	if ctx.Err() != nil {
		result := canceledResult("", request.Type, 0, ctx.Err())
		executor.logResult(ctx, result)
		return result, contextFailure(ctx.Err())
	}

	executionUUID, err := uuid.NewV7()
	if err != nil {
		return Result{}, wrappedFailure(err, codeIDGenerationFailed, failure.CategoryInternal, "job execution identifier generation failed", false)
	}
	executionID := executionUUID.String()

	if request.Idempotency != nil {
		duplicate, reserved, reserveErr := executor.idempotency.reserve(request.Type, *request.Idempotency)
		if reserveErr != nil {
			return Result{}, reserveErr
		}
		if duplicate {
			reserved.Duplicate = true
			executor.logResult(ctx, reserved)
			return reserved, errorForResult(reserved)
		}
	}

	result, executionErr := executor.executeAttempts(ctx, handler, request, policy, executionID)
	if request.Idempotency != nil {
		executor.idempotency.complete(request.Type, *request.Idempotency, result)
	}
	executor.logResult(ctx, result)
	return result, executionErr
}

// Enqueue submits one request to the bounded process-local queue. It fails
// closed when the executor is stopping or the queue is full.
func (executor *Executor) Enqueue(ctx context.Context, request Request) (*Future, error) {
	if executor == nil || executor.registry == nil {
		return nil, classifiedFailure(codeExecutorInvalid, failure.CategoryInvariant, "job executor is invalid", false)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if _, err := validateRequest(request); err != nil {
		return nil, err
	}
	if _, err := executor.registry.lookup(request.Type); err != nil {
		return nil, err
	}
	if ctx.Err() != nil {
		return nil, contextFailure(ctx.Err())
	}

	completion := &completion{done: make(chan struct{})}
	queued := queuedRequest{ctx: ctx, request: request, completion: completion}

	executor.mu.RLock()
	defer executor.mu.RUnlock()
	if !executor.accepting {
		return nil, classifiedFailure(codeExecutorStopping, failure.CategoryUnavailable, "job executor is not accepting work", false)
	}

	// Accepted work is registered before it becomes visible to a worker. Holding
	// the read lock guarantees Shutdown cannot begin waiting until this admission
	// decision is complete.
	executor.work.Add(1)
	select {
	case executor.queue <- queued:
		return &Future{completion: completion}, nil
	case <-ctx.Done():
		executor.work.Done()
		return nil, contextFailure(ctx.Err())
	default:
		executor.work.Done()
		return nil, classifiedFailure(codeQueueFull, failure.CategoryUnavailable, "job queue is full", true)
	}
}

// Shutdown stops accepting work, drains work accepted before shutdown while the
// supplied context remains valid, then cancels active work if the shutdown
// deadline is reached.
func (executor *Executor) Shutdown(ctx context.Context) error {
	if executor == nil {
		return classifiedFailure(codeExecutorInvalid, failure.CategoryInvariant, "job executor is invalid", false)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	executor.shutdown.Do(func() {
		executor.mu.Lock()
		executor.accepting = false
		close(executor.queue)
		executor.mu.Unlock()

		go func() {
			executor.work.Wait()
			executor.workers.Wait()
			executor.rootCancel()
			close(executor.shutdownDone)
		}()
	})

	select {
	case <-executor.shutdownDone:
		return nil
	case <-ctx.Done():
		executor.rootCancel()
		return contextFailure(ctx.Err())
	}
}

func (executor *Executor) worker() {
	defer executor.workers.Done()
	for queued := range executor.queue {
		executionCtx, cancel := linkedContext(queued.ctx, executor.rootCtx)
		result, err := executor.execute(executionCtx, queued.request)
		cancel()
		queued.completion.outcome = outcome{result: result, err: err}
		close(queued.completion.done)
		executor.work.Done()
	}
}

func (executor *Executor) executeAttempts(ctx context.Context, handler Handler, request Request, policy RetryPolicy, executionID string) (Result, error) {
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			result := canceledResult(executionID, request.Type, attempt-1, ctx.Err())
			return result, contextFailure(ctx.Err())
		}

		invocation := Invocation{
			ExecutionID: executionID,
			Type:        request.Type,
			Attempt:     attempt,
			Payload:     request.Payload,
		}
		if request.Idempotency != nil {
			invocation.IdempotencyKey = request.Idempotency.Key
		}

		handlerErr := runHandler(ctx, handler, invocation)
		if handlerErr == nil {
			return Result{
				ExecutionID: executionID,
				Type:        request.Type,
				State:       StateSucceeded,
				Reason:      ReasonOK,
				Attempts:    attempt,
			}, nil
		}
		if ctx.Err() != nil {
			result := canceledResult(executionID, request.Type, attempt, ctx.Err())
			return result, contextFailure(ctx.Err())
		}
		if attempt == policy.MaxAttempts {
			result := Result{
				ExecutionID: executionID,
				Type:        request.Type,
				State:       StateFailed,
				Reason:      ReasonHandlerFailed,
				Attempts:    attempt,
			}
			return result, wrappedFailure(handlerErr, codeExecutionFailed, failure.CategoryInternal, "job execution failed", false)
		}

		delay := policy.backoffAfter(attempt)
		if err := waitBackoff(ctx, delay); err != nil {
			result := canceledResult(executionID, request.Type, attempt, err)
			return result, contextFailure(err)
		}
	}

	return Result{}, classifiedFailure(codeExecutionFailed, failure.CategoryInvariant, "job execution did not reach a terminal state", false)
}

func validateRequest(request Request) (RetryPolicy, error) {
	if !request.Type.valid() {
		return RetryPolicy{}, classifiedFailure(codeRequestInvalid, failure.CategoryValidation, "job request is invalid", false)
	}
	policy, ok := request.Retry.normalize()
	if !ok {
		return RetryPolicy{}, classifiedFailure(codeRequestInvalid, failure.CategoryValidation, "job retry policy is invalid", false)
	}
	if request.Idempotency != nil && !request.Idempotency.valid() {
		return RetryPolicy{}, classifiedFailure(codeRequestInvalid, failure.CategoryValidation, "job idempotency contract is invalid", false)
	}
	if policy.MaxAttempts > 1 && request.Idempotency == nil {
		return RetryPolicy{}, classifiedFailure(codeRetryRequiresIdempotency, failure.CategoryValidation, "job retry requires an explicit idempotency contract", false)
	}
	return policy, nil
}

func runHandler(ctx context.Context, handler Handler, invocation Invocation) (err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("job handler failed safely")
		}
	}()
	return handler(ctx, invocation)
}

func waitBackoff(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func canceledResult(executionID string, jobType Type, attempts int, err error) Result {
	state := StateCanceled
	reason := ReasonCanceled
	if errors.Is(err, context.DeadlineExceeded) {
		state = StateDeadline
		reason = ReasonDeadline
	}
	return Result{
		ExecutionID: executionID,
		Type:        jobType,
		State:       state,
		Reason:      reason,
		Attempts:    attempts,
	}
}

func contextFailure(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return wrappedFailure(err, codeExecutionDeadline, failure.CategoryTimeout, "job execution deadline exceeded", false)
	}
	return wrappedFailure(err, codeExecutionCanceled, failure.CategoryUnavailable, "job execution was canceled", false)
}

func errorForResult(result Result) error {
	switch result.State {
	case StateSucceeded:
		return nil
	case StateDeadline:
		return classifiedFailure(codeExecutionDeadline, failure.CategoryTimeout, "job execution deadline exceeded", false)
	case StateCanceled:
		return classifiedFailure(codeExecutionCanceled, failure.CategoryUnavailable, "job execution was canceled", false)
	case StateFailed:
		return classifiedFailure(codeExecutionFailed, failure.CategoryInternal, "job execution failed", false)
	default:
		return classifiedFailure(codeExecutionFailed, failure.CategoryInvariant, "job execution result is invalid", false)
	}
}

func linkedContext(parent, shutdown context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	linked, cancel := context.WithCancel(parent)
	stop := context.AfterFunc(shutdown, cancel)
	return linked, func() {
		stop()
		cancel()
	}
}

func (guard *memoryIdempotency) reserve(jobType Type, contract Idempotency) (bool, Result, error) {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	key := string(jobType) + "\x00" + contract.Key
	entry, exists := guard.entries[key]
	if exists {
		if entry.fingerprint != contract.Fingerprint {
			return false, Result{}, classifiedFailure(codeIdempotencyConflict, failure.CategoryConflict, "job idempotency key conflicts with different input", false)
		}
		if !entry.completed {
			return false, Result{}, classifiedFailure(codeIdempotencyInProgress, failure.CategoryConflict, "job idempotency key is already in progress", false)
		}
		return true, entry.result, nil
	}
	guard.entries[key] = idempotencyEntry{fingerprint: contract.Fingerprint}
	return false, Result{}, nil
}

func (guard *memoryIdempotency) complete(jobType Type, contract Idempotency, result Result) {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	key := string(jobType) + "\x00" + contract.Key
	entry, exists := guard.entries[key]
	if !exists || entry.fingerprint != contract.Fingerprint {
		return
	}
	entry.completed = true
	entry.result = result
	guard.entries[key] = entry
}

func (executor *Executor) logResult(ctx context.Context, result Result) {
	if executor == nil || executor.logger == nil {
		return
	}
	attrs := []slog.Attr{
		slog.String("job_type", string(result.Type)),
		slog.String("execution_id", result.ExecutionID),
		slog.String("state", string(result.State)),
		slog.String("reason", string(result.Reason)),
		slog.Int("attempts", result.Attempts),
		slog.Bool("duplicate", result.Duplicate),
	}
	if result.State == StateSucceeded {
		executor.logger.Info(ctx, "job execution completed", attrs...)
		return
	}
	executor.logger.Warn(ctx, "job execution completed", attrs...)
}

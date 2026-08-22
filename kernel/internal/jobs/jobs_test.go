package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/observability"
	"github.com/google/uuid"
)

func TestExecutionIDUsesUUIDv7(t *testing.T) {
	executor := newTestExecutor(t, "kernel.maintenance", func(context.Context, Invocation) error { return nil }, Settings{Workers: 1, QueueCapacity: 1}, nil)

	result, err := executor.Execute(context.Background(), Request{Type: "kernel.maintenance"})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	parsed, err := uuid.Parse(result.ExecutionID)
	if err != nil {
		t.Fatalf("parse execution id: %v", err)
	}
	if parsed.Version() != uuid.Version(7) {
		t.Fatalf("execution id version = %d, want 7", parsed.Version())
	}
	if result.State != StateSucceeded || result.Reason != ReasonOK || result.Attempts != 1 {
		t.Fatalf("result = %+v", result)
	}
}

func TestRegistryIsDeterministicFrozenAndUnknownJobsFailSafely(t *testing.T) {
	registry := NewRegistry()
	mustRegisterJob(t, registry, "zeta.task", func(context.Context, Invocation) error { return nil })
	mustRegisterJob(t, registry, "alpha.task", func(context.Context, Invocation) error { return nil })

	types := registry.Types()
	if len(types) != 2 || types[0] != "alpha.task" || types[1] != "zeta.task" {
		t.Fatalf("registry types = %v", types)
	}
	assertFailureCode(t, registry.Register("alpha.task", func(context.Context, Invocation) error { return nil }), codeTypeDuplicate)
	assertFailureCode(t, registry.Register("../unsafe", func(context.Context, Invocation) error { return nil }), codeTypeInvalid)

	executor, err := NewExecutor(registry, nil, Settings{Workers: 1, QueueCapacity: 1})
	if err != nil {
		t.Fatalf("new executor: %v", err)
	}
	shutdownExecutor(t, executor)
	assertFailureCode(t, registry.Register("later.task", func(context.Context, Invocation) error { return nil }), codeRegistryFrozen)

	_, err = executor.Execute(context.Background(), Request{Type: "missing.task"})
	assertFailureCode(t, err, codeTypeUnknown)
}

func TestRetryRequiresIdempotencyAndIsBounded(t *testing.T) {
	var calls atomic.Int32
	executor := newTestExecutor(t, "retry.task", func(context.Context, Invocation) error {
		if calls.Add(1) < 3 {
			return errors.New("synthetic handler failure")
		}
		return nil
	}, Settings{Workers: 1, QueueCapacity: 1}, nil)

	policy := RetryPolicy{MaxAttempts: 3, InitialBackoff: time.Millisecond, MaxBackoff: 2 * time.Millisecond}
	_, err := executor.Execute(context.Background(), Request{Type: "retry.task", Retry: policy})
	assertFailureCode(t, err, codeRetryRequiresIdempotency)

	contract := &Idempotency{Key: "retry-key-1", Fingerprint: "sha256:retry-input-1"}
	result, err := executor.Execute(context.Background(), Request{Type: "retry.task", Retry: policy, Idempotency: contract})
	if err != nil {
		t.Fatalf("retry execute: %v", err)
	}
	if result.State != StateSucceeded || result.Attempts != 3 || calls.Load() != 3 {
		t.Fatalf("retry result = %+v calls=%d", result, calls.Load())
	}
	if policy.backoffAfter(1) != time.Millisecond || policy.backoffAfter(2) != 2*time.Millisecond || policy.backoffAfter(3) != 0 {
		t.Fatalf("unexpected backoff sequence")
	}
}

func TestIdempotencyDuplicateAndConflictContracts(t *testing.T) {
	var calls atomic.Int32
	executor := newTestExecutor(t, "idempotent.task", func(context.Context, Invocation) error {
		calls.Add(1)
		return nil
	}, Settings{Workers: 1, QueueCapacity: 1}, nil)

	contract := &Idempotency{Key: "operation-key-1", Fingerprint: "sha256:input-a"}
	request := Request{Type: "idempotent.task", Idempotency: contract}
	first, err := executor.Execute(context.Background(), request)
	if err != nil {
		t.Fatalf("first execute: %v", err)
	}
	second, err := executor.Execute(context.Background(), request)
	if err != nil {
		t.Fatalf("duplicate execute: %v", err)
	}
	if !second.Duplicate || second.ExecutionID != first.ExecutionID || calls.Load() != 1 {
		t.Fatalf("duplicate result = %+v first=%+v calls=%d", second, first, calls.Load())
	}

	conflicting := Request{
		Type:        "idempotent.task",
		Idempotency: &Idempotency{Key: contract.Key, Fingerprint: "sha256:input-b"},
	}
	_, err = executor.Execute(context.Background(), conflicting)
	assertFailureCode(t, err, codeIdempotencyConflict)
}

func TestConcurrentDuplicateFailsClosedWhileOriginalIsInProgress(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	executor := newTestExecutor(t, "guarded.task", func(context.Context, Invocation) error {
		select {
		case started <- struct{}{}:
		default:
		}
		<-release
		return nil
	}, Settings{Workers: 1, QueueCapacity: 1}, nil)

	request := Request{
		Type:        "guarded.task",
		Idempotency: &Idempotency{Key: "guarded-key", Fingerprint: "sha256:guarded-input"},
	}
	firstDone := make(chan error, 1)
	go func() {
		_, err := executor.Execute(context.Background(), request)
		firstDone <- err
	}()
	<-started

	_, err := executor.Execute(context.Background(), request)
	assertFailureCode(t, err, codeIdempotencyInProgress)
	close(release)
	if err := <-firstDone; err != nil {
		t.Fatalf("original execute: %v", err)
	}
}

func TestCancellationDuringRetryBackoffIsBounded(t *testing.T) {
	var calls atomic.Int32
	executor := newTestExecutor(t, "cancel.task", func(context.Context, Invocation) error {
		calls.Add(1)
		return errors.New("retryable synthetic failure")
	}, Settings{Workers: 1, QueueCapacity: 1}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	started := time.Now()
	result, err := executor.Execute(ctx, Request{
		Type: "cancel.task",
		Retry: RetryPolicy{MaxAttempts: 3, InitialBackoff: 100 * time.Millisecond, MaxBackoff: 100 * time.Millisecond},
		Idempotency: &Idempotency{Key: "cancel-key", Fingerprint: "sha256:cancel-input"},
	})
	if elapsed := time.Since(started); elapsed > 500*time.Millisecond {
		t.Fatalf("canceled retry took %s", elapsed)
	}
	assertFailureCode(t, err, codeExecutionDeadline)
	if result.State != StateDeadline || result.Reason != ReasonDeadline || result.Attempts != 1 || calls.Load() != 1 {
		t.Fatalf("canceled retry result = %+v calls=%d", result, calls.Load())
	}
}

func TestWorkerConcurrencyAndQueueCapacityAreBounded(t *testing.T) {
	var active atomic.Int32
	var maximum atomic.Int32
	started := make(chan struct{}, 2)
	release := make(chan struct{})

	handler := func(context.Context, Invocation) error {
		current := active.Add(1)
		for {
			observed := maximum.Load()
			if current <= observed || maximum.CompareAndSwap(observed, current) {
				break
			}
		}
		started <- struct{}{}
		<-release
		active.Add(-1)
		return nil
	}
	executor := newTestExecutor(t, "bounded.task", handler, Settings{Workers: 2, QueueCapacity: 2}, nil)

	first, err := executor.Enqueue(context.Background(), Request{Type: "bounded.task"})
	if err != nil {
		t.Fatalf("enqueue first: %v", err)
	}
	second, err := executor.Enqueue(context.Background(), Request{Type: "bounded.task"})
	if err != nil {
		t.Fatalf("enqueue second: %v", err)
	}
	<-started
	<-started

	third, err := executor.Enqueue(context.Background(), Request{Type: "bounded.task"})
	if err != nil {
		t.Fatalf("enqueue third: %v", err)
	}
	fourth, err := executor.Enqueue(context.Background(), Request{Type: "bounded.task"})
	if err != nil {
		t.Fatalf("enqueue fourth: %v", err)
	}
	_, err = executor.Enqueue(context.Background(), Request{Type: "bounded.task"})
	assertFailureCode(t, err, codeQueueFull)

	close(release)
	for index, future := range []*Future{first, second, third, fourth} {
		result, waitErr := future.Wait(context.Background())
		if waitErr != nil || result.State != StateSucceeded {
			t.Fatalf("future %d result=%+v err=%v", index, result, waitErr)
		}
	}
	if maximum.Load() != 2 {
		t.Fatalf("maximum concurrent handlers = %d, want 2", maximum.Load())
	}
}

func TestGracefulShutdownDrainsQueuedWorkAndStopsAcceptance(t *testing.T) {
	var calls atomic.Int32
	started := make(chan struct{})
	release := make(chan struct{})
	executor := newTestExecutor(t, "drain.task", func(context.Context, Invocation) error {
		call := calls.Add(1)
		if call == 1 {
			close(started)
			<-release
		}
		return nil
	}, Settings{Workers: 1, QueueCapacity: 2}, nil)

	first, err := executor.Enqueue(context.Background(), Request{Type: "drain.task"})
	if err != nil {
		t.Fatalf("enqueue first: %v", err)
	}
	<-started
	second, err := executor.Enqueue(context.Background(), Request{Type: "drain.task"})
	if err != nil {
		t.Fatalf("enqueue second: %v", err)
	}

	shutdownDone := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		shutdownDone <- executor.Shutdown(ctx)
	}()
	close(release)

	if err := <-shutdownDone; err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	for _, future := range []*Future{first, second} {
		result, waitErr := future.Wait(context.Background())
		if waitErr != nil || result.State != StateSucceeded {
			t.Fatalf("drained result=%+v err=%v", result, waitErr)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("drained calls = %d, want 2", calls.Load())
	}
	_, err = executor.Enqueue(context.Background(), Request{Type: "drain.task"})
	assertFailureCode(t, err, codeExecutorStopping)
}

func TestShutdownDeadlineCancelsActiveHandler(t *testing.T) {
	started := make(chan struct{})
	observedCancellation := make(chan struct{})
	executor := newTestExecutor(t, "shutdown.task", func(ctx context.Context, _ Invocation) error {
		close(started)
		<-ctx.Done()
		close(observedCancellation)
		return ctx.Err()
	}, Settings{Workers: 1, QueueCapacity: 1}, nil)

	future, err := executor.Enqueue(context.Background(), Request{Type: "shutdown.task"})
	if err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	startedShutdown := time.Now()
	err = executor.Shutdown(ctx)
	if elapsed := time.Since(startedShutdown); elapsed > 500*time.Millisecond {
		t.Fatalf("bounded shutdown took %s", elapsed)
	}
	assertFailureCode(t, err, codeExecutionDeadline)

	select {
	case <-observedCancellation:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("active handler did not observe shutdown cancellation")
	}
	result, waitErr := future.Wait(context.Background())
	assertFailureCode(t, waitErr, codeExecutionCanceled)
	if result.State != StateCanceled {
		t.Fatalf("shutdown-canceled result = %+v", result)
	}
}

func TestObservabilityPropagatesCorrelationWithoutPrivateJobData(t *testing.T) {
	logger, capture := captureJobLogger()
	seenCorrelation := make(chan string, 1)
	executor := newTestExecutor(t, "observe.task", func(ctx context.Context, invocation Invocation) error {
		value, _ := observability.CorrelationIDFromContext(ctx)
		seenCorrelation <- value
		if invocation.Payload == nil || invocation.IdempotencyKey == "" {
			return nil
		}
		return errors.New("private-handler-diagnostic-marker")
	}, Settings{Workers: 1, QueueCapacity: 1}, logger)

	ctx, err := observability.WithCorrelationID(context.Background(), "corr-job-123")
	if err != nil {
		t.Fatalf("correlation context: %v", err)
	}
	result, err := executor.Execute(ctx, Request{
		Type:        "observe.task",
		Payload:     "private-payload-marker",
		Idempotency: &Idempotency{Key: "private-idempotency-marker", Fingerprint: "private-fingerprint-marker"},
	})
	assertFailureCode(t, err, codeExecutionFailed)
	if result.State != StateFailed || <-seenCorrelation != "corr-job-123" {
		t.Fatalf("observability result = %+v", result)
	}

	records := capture.Records()
	if len(records) != 1 {
		t.Fatalf("captured records = %d, want 1", len(records))
	}
	serialized := fmt.Sprint(records[0].Attributes)
	for _, private := range []string{"private-handler-diagnostic-marker", "private-payload-marker", "private-idempotency-marker", "private-fingerprint-marker"} {
		if strings.Contains(serialized, private) || strings.Contains(records[0].Message, private) {
			t.Fatalf("observability leaked private job data %q: %s", private, serialized)
		}
	}
	if records[0].Attributes["correlation_id"] != "corr-job-123" || records[0].Attributes["job_type"] != "observe.task" {
		t.Fatalf("captured safe attributes = %+v", records[0].Attributes)
	}
	if records[0].Severity != slog.LevelWarn {
		t.Fatalf("failure severity = %v, want WARN", records[0].Severity)
	}
}

func TestSchedulesAreUTCNormalizedAndDeterministic(t *testing.T) {
	location := time.FixedZone("fixture", 2*60*60)
	start := time.Date(2026, 8, 23, 10, 0, 0, 0, location)

	oneShot, err := NewOneShot(start)
	if err != nil {
		t.Fatalf("one shot: %v", err)
	}
	if oneShot.StartAt.Location() != time.UTC {
		t.Fatalf("one-shot location = %v", oneShot.StartAt.Location())
	}
	if next, ok := oneShot.Next(start.Add(-time.Second)); !ok || !next.Equal(start) {
		t.Fatalf("one-shot next=%v ok=%v", next, ok)
	}
	if _, ok := oneShot.Next(start); ok {
		t.Fatal("one-shot schedule must not repeat")
	}

	recurring, err := NewRecurring(start, 10*time.Minute)
	if err != nil {
		t.Fatalf("recurring: %v", err)
	}
	if next, ok := recurring.Next(start); !ok || !next.Equal(start.Add(10*time.Minute)) {
		t.Fatalf("recurring first next=%v ok=%v", next, ok)
	}
	if next, ok := recurring.Next(start.Add(25 * time.Minute)); !ok || !next.Equal(start.Add(30*time.Minute)) {
		t.Fatalf("recurring later next=%v ok=%v", next, ok)
	}
	if _, err := NewRecurring(start, 100*time.Millisecond); err == nil {
		t.Fatal("expected invalid recurring interval")
	}
	if _, err := NewOneShot(time.Time{}); err == nil {
		t.Fatal("expected invalid one-shot schedule")
	}
}

func TestPanickingHandlerFailsSafely(t *testing.T) {
	executor := newTestExecutor(t, "panic.task", func(context.Context, Invocation) error {
		panic("private-panic-marker")
	}, Settings{Workers: 1, QueueCapacity: 1}, nil)

	result, err := executor.Execute(context.Background(), Request{Type: "panic.task"})
	assertFailureCode(t, err, codeExecutionFailed)
	if result.State != StateFailed || strings.Contains(err.Error(), "private-panic-marker") {
		t.Fatalf("panic result=%+v err=%v", result, err)
	}
}

func TestConcurrentExecutionsAreRaceSafe(t *testing.T) {
	var calls atomic.Int32
	executor := newTestExecutor(t, "race.task", func(context.Context, Invocation) error {
		calls.Add(1)
		return nil
	}, Settings{Workers: 4, QueueCapacity: 8}, nil)

	var workers sync.WaitGroup
	for index := 0; index < 24; index++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			result, err := executor.Execute(context.Background(), Request{Type: "race.task"})
			if err != nil || result.State != StateSucceeded {
				t.Errorf("concurrent result=%+v err=%v", result, err)
			}
		}()
	}
	workers.Wait()
	if calls.Load() != 24 {
		t.Fatalf("concurrent calls = %d, want 24", calls.Load())
	}
}

func newTestExecutor(t *testing.T, jobType Type, handler Handler, settings Settings, logger *observability.Logger) *Executor {
	t.Helper()
	registry := NewRegistry()
	mustRegisterJob(t, registry, jobType, handler)
	executor, err := NewExecutor(registry, logger, settings)
	if err != nil {
		t.Fatalf("new executor: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = executor.Shutdown(ctx)
	})
	return executor
}

func shutdownExecutor(t *testing.T, executor *Executor) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := executor.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown executor: %v", err)
	}
}

func mustRegisterJob(t *testing.T, registry *Registry, jobType Type, handler Handler) {
	t.Helper()
	if err := registry.Register(jobType, handler); err != nil {
		t.Fatalf("register %s: %v", jobType, err)
	}
}

func assertFailureCode(t *testing.T, err error, code failure.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected failure code %s", code)
	}
	classified, ok := failure.As(err)
	if !ok {
		t.Fatalf("error is not an Omnexa failure: %v", err)
	}
	if classified.Code() != code {
		t.Fatalf("failure code = %s, want %s", classified.Code(), code)
	}
}

func captureJobLogger() (*observability.Logger, *observability.Capture) {
	settings := observability.Settings{
		Enabled:         true,
		ServiceName:     "omnexa-kernel",
		ServiceVersion:  "test",
		Environment:     config.EnvironmentTest,
		LogLevel:        slog.LevelDebug,
		ExportInterval:  time.Second,
		ExportTimeout:   time.Second,
		ShutdownTimeout: time.Second,
	}
	return observability.NewCaptureLogger(settings)
}

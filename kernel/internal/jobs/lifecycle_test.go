package jobs

import (
	"context"
	"testing"
	"time"
)

func TestFutureWaitReturnsSameTerminalOutcomeRepeatedly(t *testing.T) {
	executor := newTestExecutor(t, "future.task", func(context.Context, Invocation) error { return nil }, Settings{Workers: 1, QueueCapacity: 1}, nil)

	future, err := executor.Enqueue(context.Background(), Request{Type: "future.task"})
	if err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	first, firstErr := future.Wait(context.Background())
	if firstErr != nil {
		t.Fatalf("first wait: %v", firstErr)
	}
	second, secondErr := future.Wait(context.Background())
	if secondErr != nil {
		t.Fatalf("second wait: %v", secondErr)
	}
	if first != second || first.ExecutionID == "" || first.State != StateSucceeded {
		t.Fatalf("repeatable future outcomes differ: first=%+v second=%+v", first, second)
	}
}

func TestExecuteRejectsNewWorkAfterShutdown(t *testing.T) {
	executor := newTestExecutor(t, "stopped.task", func(context.Context, Invocation) error { return nil }, Settings{Workers: 1, QueueCapacity: 1}, nil)
	shutdownExecutor(t, executor)

	_, err := executor.Execute(context.Background(), Request{Type: "stopped.task"})
	assertFailureCode(t, err, codeExecutorStopping)

	_, err = executor.Enqueue(context.Background(), Request{Type: "stopped.task"})
	assertFailureCode(t, err, codeExecutorStopping)
}

func TestGracefulShutdownWaitsForAcceptedSynchronousWork(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	executor := newTestExecutor(t, "direct-drain.task", func(context.Context, Invocation) error {
		close(started)
		<-release
		return nil
	}, Settings{Workers: 1, QueueCapacity: 1}, nil)

	executionDone := make(chan error, 1)
	go func() {
		_, err := executor.Execute(context.Background(), Request{Type: "direct-drain.task"})
		executionDone <- err
	}()
	<-started

	shutdownDone := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		shutdownDone <- executor.Shutdown(ctx)
	}()

	deadline := time.Now().Add(500 * time.Millisecond)
	for {
		executor.mu.RLock()
		accepting := executor.accepting
		executor.mu.RUnlock()
		if !accepting {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("shutdown did not stop admission")
		}
		time.Sleep(time.Millisecond)
	}

	select {
	case err := <-shutdownDone:
		t.Fatalf("shutdown returned before accepted synchronous work drained: %v", err)
	default:
	}

	close(release)
	if err := <-executionDone; err != nil {
		t.Fatalf("synchronous execution: %v", err)
	}
	if err := <-shutdownDone; err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestShutdownDeadlineCancelsAcceptedSynchronousWork(t *testing.T) {
	started := make(chan struct{})
	observedCancellation := make(chan struct{})
	executor := newTestExecutor(t, "direct-cancel.task", func(ctx context.Context, _ Invocation) error {
		close(started)
		<-ctx.Done()
		close(observedCancellation)
		return ctx.Err()
	}, Settings{Workers: 1, QueueCapacity: 1}, nil)

	executionDone := make(chan error, 1)
	go func() {
		_, err := executor.Execute(context.Background(), Request{Type: "direct-cancel.task"})
		executionDone <- err
	}()
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	startedShutdown := time.Now()
	err := executor.Shutdown(ctx)
	if elapsed := time.Since(startedShutdown); elapsed > 500*time.Millisecond {
		t.Fatalf("bounded shutdown took %s", elapsed)
	}
	assertFailureCode(t, err, codeExecutionDeadline)

	select {
	case <-observedCancellation:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("synchronous handler did not observe shutdown cancellation")
	}
	assertFailureCode(t, <-executionDone, codeExecutionCanceled)
}

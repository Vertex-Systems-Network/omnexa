package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestRetryPolicyRejectsUnboundedOrAmbiguousValues(t *testing.T) {
	valid := testRetryPolicy()
	tests := []RetryPolicy{
		{Version: 1, MaxAttempts: 3, InitialBackoff: time.Second, MaxBackoff: time.Minute},
		{ID: "policy.v1", MaxAttempts: 3, InitialBackoff: time.Second, MaxBackoff: time.Minute},
		{ID: "policy.v1", Version: 1, InitialBackoff: time.Second, MaxBackoff: time.Minute},
		{ID: "policy.v1", Version: 1, MaxAttempts: MaxRetryAttempts + 1, InitialBackoff: time.Second, MaxBackoff: time.Minute},
		{ID: "policy.v1", Version: 1, MaxAttempts: 3, InitialBackoff: 0, MaxBackoff: time.Minute},
		{ID: "policy.v1", Version: 1, MaxAttempts: 3, InitialBackoff: time.Minute, MaxBackoff: time.Second},
		{ID: "policy.v1", Version: 1, MaxAttempts: 3, InitialBackoff: time.Second, MaxBackoff: MaxRetryBackoff + time.Second},
	}
	for index, policy := range tests {
		if err := policy.Validate(); err == nil {
			t.Fatalf("invalid policy %d unexpectedly passed validation: %+v", index, policy)
		}
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid policy failed validation: %v", err)
	}
}

func TestRetryPolicyDelayIsDeterministicCappedExponential(t *testing.T) {
	policy := RetryPolicy{
		ID:             "events.default.v1",
		Version:        1,
		MaxAttempts:    6,
		InitialBackoff: 2 * time.Second,
		MaxBackoff:     5 * time.Second,
	}
	want := []time.Duration{2 * time.Second, 4 * time.Second, 5 * time.Second, 5 * time.Second, 5 * time.Second}
	for attempt, expected := range want {
		delay, err := policy.DelayAfter(uint32(attempt + 1))
		if err != nil {
			t.Fatalf("attempt %d delay failed: %v", attempt+1, err)
		}
		if delay != expected {
			t.Fatalf("attempt %d delay = %s, want %s", attempt+1, delay, expected)
		}
	}
	if _, err := policy.DelayAfter(policy.MaxAttempts); err == nil {
		t.Fatal("exhausted attempt unexpectedly returned another delay")
	}
}

func TestDecideRetrySchedulesStructuredSafeRetryUsingUTC(t *testing.T) {
	policy := testRetryPolicy()
	failureTime := time.Date(2026, 9, 5, 22, 0, 0, 0, time.FixedZone("local", 5*60*60))
	decision, err := DecideRetry(
		policy,
		2,
		failureTime,
		mustRetryFailure(t, failure.CategoryDependency, true),
		true,
	)
	if err != nil {
		t.Fatalf("retry decision failed: %v", err)
	}
	if decision.Disposition != RetryDispositionRetryable || !decision.ConsumeAttempt {
		t.Fatalf("unexpected retry decision: %+v", decision)
	}
	if decision.Delay != 4*time.Second {
		t.Fatalf("delay = %s, want 4s", decision.Delay)
	}
	if decision.NextEligibleAt.Location() != time.UTC {
		t.Fatalf("next eligibility must be UTC, got %s", decision.NextEligibleAt.Location())
	}
	want := failureTime.UTC().Add(4 * time.Second)
	if !decision.NextEligibleAt.Equal(want) {
		t.Fatalf("next eligibility = %s, want %s", decision.NextEligibleAt, want)
	}
	if decision.PolicyID != policy.ID || decision.PolicyVersion != policy.Version || decision.Attempt != 2 {
		t.Fatalf("policy snapshot metadata drifted: %+v", decision)
	}
}

func TestDecideRetryFailsClosedWhenOperationIsUnsafe(t *testing.T) {
	decision, err := DecideRetry(
		testRetryPolicy(),
		1,
		time.Now(),
		mustRetryFailure(t, failure.CategoryTimeout, true),
		false,
	)
	if err != nil {
		t.Fatalf("retry decision failed: %v", err)
	}
	if decision.Disposition != RetryDispositionTerminal || decision.TerminalReason != RetryTerminalReasonUnsafeOperation {
		t.Fatalf("unsafe operation must be terminal: %+v", decision)
	}
}

func TestDecideRetryNeverUpgradesStructuredNonRetryableFailure(t *testing.T) {
	decision, err := DecideRetry(
		testRetryPolicy(),
		1,
		time.Now(),
		mustRetryFailure(t, failure.CategoryDependency, false),
		true,
	)
	if err != nil {
		t.Fatalf("retry decision failed: %v", err)
	}
	if decision.Disposition != RetryDispositionTerminal || decision.TerminalReason != RetryTerminalReasonNonRetryable {
		t.Fatalf("non-retryable structured failure was upgraded: %+v", decision)
	}
}

func TestDecideRetryDoesNotAutoRetryNonCandidateCategory(t *testing.T) {
	decision, err := DecideRetry(
		testRetryPolicy(),
		1,
		time.Now(),
		mustRetryFailure(t, failure.CategoryAuthorization, true),
		true,
	)
	if err != nil {
		t.Fatalf("retry decision failed: %v", err)
	}
	if decision.Disposition != RetryDispositionTerminal || decision.TerminalReason != RetryTerminalReasonNonRetryable {
		t.Fatalf("authorization failure must remain terminal: %+v", decision)
	}
}

func TestDecideRetryUnknownFailureDefaultsTerminal(t *testing.T) {
	decision, err := DecideRetry(testRetryPolicy(), 1, time.Now(), errors.New("opaque provider text"), true)
	if err != nil {
		t.Fatalf("retry decision failed: %v", err)
	}
	if decision.Disposition != RetryDispositionTerminal || decision.TerminalReason != RetryTerminalReasonUnknownFailure || !decision.ConsumeAttempt {
		t.Fatalf("unknown failure did not fail closed: %+v", decision)
	}
}

func TestDecideRetryCancellationIsInterruptedWithoutAttemptConsumption(t *testing.T) {
	decision, err := DecideRetry(testRetryPolicy(), 2, time.Time{}, context.Canceled, true)
	if err != nil {
		t.Fatalf("interruption decision failed: %v", err)
	}
	if decision.Disposition != RetryDispositionInterrupted || decision.ConsumeAttempt || decision.TerminalReason != RetryTerminalReasonNone {
		t.Fatalf("cancellation must remain interruption: %+v", decision)
	}
	if !decision.NextEligibleAt.IsZero() || decision.Delay != 0 {
		t.Fatalf("interruption must not schedule retry: %+v", decision)
	}
}

func TestDecideRetryExhaustionBecomesTerminal(t *testing.T) {
	policy := testRetryPolicy()
	decision, err := DecideRetry(
		policy,
		policy.MaxAttempts,
		time.Now(),
		mustRetryFailure(t, failure.CategoryUnavailable, true),
		true,
	)
	if err != nil {
		t.Fatalf("retry decision failed: %v", err)
	}
	if decision.Disposition != RetryDispositionTerminal || decision.TerminalReason != RetryTerminalReasonAttemptsExhausted || !decision.ConsumeAttempt {
		t.Fatalf("exhaustion must become terminal: %+v", decision)
	}
}

func TestDecideRetryRequiresAuthoritativeFailureTimeOnlyWhenScheduling(t *testing.T) {
	_, err := DecideRetry(
		testRetryPolicy(),
		1,
		time.Time{},
		mustRetryFailure(t, failure.CategoryRateLimit, true),
		true,
	)
	if err == nil {
		t.Fatal("retry scheduling without authoritative failure time unexpectedly succeeded")
	}
}

func testRetryPolicy() RetryPolicy {
	return RetryPolicy{
		ID:             "events.default.v1",
		Version:        1,
		MaxAttempts:    4,
		InitialBackoff: 2 * time.Second,
		MaxBackoff:     30 * time.Second,
	}
}

func mustRetryFailure(t *testing.T, category failure.Category, retryable bool) *failure.Error {
	t.Helper()
	classified, err := failure.New(
		"events.test.failure",
		category,
		"synthetic event failure",
		failure.WithRetryable(retryable),
	)
	if err != nil {
		t.Fatalf("construct failure: %v", err)
	}
	return classified
}

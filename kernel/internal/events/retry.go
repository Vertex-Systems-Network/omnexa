package events

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	// MaxRetryAttempts is the hard V1 ceiling for one policy snapshot, including
	// the first handler attempt.
	MaxRetryAttempts uint32 = 100
	// MaxRetryBackoff prevents caller-controlled policies from creating
	// operationally unbounded eligibility windows.
	MaxRetryBackoff       = 24 * time.Hour
	maxRetryPolicyIDBytes = 128
)

const (
	codeRetryPolicyInvalid  failure.Code = "events.retry.policy_invalid"
	codeRetryAttemptInvalid failure.Code = "events.retry.attempt_invalid"
	codeRetryFailureInvalid failure.Code = "events.retry.failure_invalid"
	codeRetryTimeInvalid    failure.Code = "events.retry.time_invalid"
)

var retryPolicyIDPattern = regexp.MustCompile(`^[a-z][a-z0-9._-]*$`)

// RetryDisposition is the provider-neutral P04.06 disposition of one
// authoritative failed delivery.
type RetryDisposition string

const (
	RetryDispositionRetryable   RetryDisposition = "retryable"
	RetryDispositionTerminal    RetryDisposition = "terminal"
	RetryDispositionInterrupted RetryDisposition = "interrupted"
)

// RetryTerminalReason is a stable bounded terminal reason. It intentionally
// carries no raw error/provider text.
type RetryTerminalReason string

const (
	RetryTerminalReasonNone              RetryTerminalReason = ""
	RetryTerminalReasonNonRetryable      RetryTerminalReason = "non_retryable"
	RetryTerminalReasonUnsafeOperation   RetryTerminalReason = "unsafe_operation"
	RetryTerminalReasonAttemptsExhausted RetryTerminalReason = "attempts_exhausted"
	RetryTerminalReasonUnknownFailure    RetryTerminalReason = "unknown_failure"
)

// RetryPolicy is the immutable policy snapshot used for one processing
// identity. Exact runtime persistence is intentionally outside Wave 1.
type RetryPolicy struct {
	ID             string
	Version        uint32
	MaxAttempts    uint32
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// Validate rejects unbounded or ambiguous policy inputs before any retry
// decision can be claimed.
func (policy RetryPolicy) Validate() error {
	if policy.ID == "" || len(policy.ID) > maxRetryPolicyIDBytes || !retryPolicyIDPattern.MatchString(policy.ID) {
		return classifiedFailure(codeRetryPolicyInvalid, failure.CategoryValidation, "event retry policy id is invalid")
	}
	if policy.Version == 0 {
		return classifiedFailure(codeRetryPolicyInvalid, failure.CategoryValidation, "event retry policy version must be positive")
	}
	if policy.MaxAttempts == 0 || policy.MaxAttempts > MaxRetryAttempts {
		return classifiedFailure(codeRetryPolicyInvalid, failure.CategoryValidation, "event retry max attempts is outside the bounded range")
	}
	if policy.InitialBackoff <= 0 || policy.MaxBackoff <= 0 {
		return classifiedFailure(codeRetryPolicyInvalid, failure.CategoryValidation, "event retry backoff durations must be positive")
	}
	if policy.InitialBackoff > policy.MaxBackoff || policy.MaxBackoff > MaxRetryBackoff {
		return classifiedFailure(codeRetryPolicyInvalid, failure.CategoryValidation, "event retry backoff durations exceed the bounded policy range")
	}
	return nil
}

// DelayAfter returns the deterministic capped exponential delay after one
// authoritative failed attempt. Attempt numbering starts at 1.
func (policy RetryPolicy) DelayAfter(attempt uint32) (time.Duration, error) {
	if err := policy.Validate(); err != nil {
		return 0, err
	}
	if attempt == 0 || attempt >= policy.MaxAttempts {
		return 0, classifiedFailure(codeRetryAttemptInvalid, failure.CategoryValidation, "event retry attempt has no eligible next automatic attempt")
	}

	delay := policy.InitialBackoff
	for current := uint32(1); current < attempt; current++ {
		if delay >= policy.MaxBackoff || delay > policy.MaxBackoff-delay {
			return policy.MaxBackoff, nil
		}
		delay += delay
		if delay >= policy.MaxBackoff {
			return policy.MaxBackoff, nil
		}
	}
	return delay, nil
}

// RetryDecision contains only bounded policy/disposition evidence. It contains
// no raw error, stack, provider payload, or authorization decision.
type RetryDecision struct {
	Disposition    RetryDisposition
	TerminalReason RetryTerminalReason
	ConsumeAttempt bool
	Attempt        uint32
	PolicyID       string
	PolicyVersion  uint32
	Delay          time.Duration
	NextEligibleAt time.Time
}

// DecideRetry classifies one authoritative failed delivery. Structured
// retryability is eligibility evidence only: operationRetrySafe must also be
// true, and only dependency/rate-limit/timeout/unavailable categories are V1
// automatic-retry candidates. Unknown failures fail closed to terminal.
func DecideRetry(
	policy RetryPolicy,
	attempt uint32,
	failureTime time.Time,
	err error,
	operationRetrySafe bool,
) (RetryDecision, error) {
	if err := policy.Validate(); err != nil {
		return RetryDecision{}, err
	}
	if attempt == 0 || attempt > policy.MaxAttempts {
		return RetryDecision{}, classifiedFailure(codeRetryAttemptInvalid, failure.CategoryValidation, "event retry attempt is outside the policy snapshot")
	}
	if err == nil {
		return RetryDecision{}, classifiedFailure(codeRetryFailureInvalid, failure.CategoryValidation, "event retry decision requires a failed delivery")
	}

	base := RetryDecision{
		Attempt:       attempt,
		PolicyID:      policy.ID,
		PolicyVersion: policy.Version,
	}

	if errors.Is(err, context.Canceled) {
		base.Disposition = RetryDispositionInterrupted
		base.ConsumeAttempt = false
		return base, nil
	}

	var structured *failure.Error
	if !errors.As(err, &structured) {
		base.Disposition = RetryDispositionTerminal
		base.TerminalReason = RetryTerminalReasonUnknownFailure
		base.ConsumeAttempt = true
		return base, nil
	}

	base.ConsumeAttempt = true
	if !structured.Retryable() || !retryCandidateCategory(structured.Category()) {
		base.Disposition = RetryDispositionTerminal
		base.TerminalReason = RetryTerminalReasonNonRetryable
		return base, nil
	}
	if !operationRetrySafe {
		base.Disposition = RetryDispositionTerminal
		base.TerminalReason = RetryTerminalReasonUnsafeOperation
		return base, nil
	}
	if attempt >= policy.MaxAttempts {
		base.Disposition = RetryDispositionTerminal
		base.TerminalReason = RetryTerminalReasonAttemptsExhausted
		return base, nil
	}
	if failureTime.IsZero() {
		return RetryDecision{}, classifiedFailure(codeRetryTimeInvalid, failure.CategoryValidation, "event retry authoritative failure time is required")
	}

	delay, delayErr := policy.DelayAfter(attempt)
	if delayErr != nil {
		return RetryDecision{}, delayErr
	}
	failureUTC := failureTime.UTC()
	next := failureUTC.Add(delay)
	if next.Year() < 1 || next.Year() > 9999 || next.Before(failureUTC) {
		return RetryDecision{}, classifiedFailure(codeRetryTimeInvalid, failure.CategoryValidation, "event retry eligibility time overflowed the supported UTC range")
	}

	base.Disposition = RetryDispositionRetryable
	base.Delay = delay
	base.NextEligibleAt = next
	return base, nil
}

func retryCandidateCategory(category failure.Category) bool {
	switch category {
	case failure.CategoryDependency,
		failure.CategoryRateLimit,
		failure.CategoryTimeout,
		failure.CategoryUnavailable:
		return true
	default:
		return false
	}
}

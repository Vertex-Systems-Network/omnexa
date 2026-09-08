package events

import (
	"regexp"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeRetryStateInvalid   failure.Code = "events.retry.state_invalid"
	codeRetryStateMalformed failure.Code = "events.retry.state_malformed"
)

var retryClaimTokenPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-57][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// RetryState is the finite durable P04.06 state of one exact processing
// identity. It is failure/scheduling evidence only and is never authorization,
// checkpoint progress, inbox completion, or an exactly-once claim.
type RetryState string

const (
	RetryStateScheduled   RetryState = "retry_scheduled"
	RetryStateQuarantined RetryState = "quarantined"
	RetryStateResolved    RetryState = "resolved"
)

func (state RetryState) valid() bool {
	switch state {
	case RetryStateScheduled, RetryStateQuarantined, RetryStateResolved:
		return true
	default:
		return false
	}
}

// RetryFailureEvidence is the bounded safe projection retained for diagnostics.
// Empty code/category represents an unknown/unstructured failure. Raw error text,
// payloads, stacks, SQL and provider responses are intentionally absent.
type RetryFailureEvidence struct {
	Code      failure.Code
	Category  failure.Category
	Retryable bool
}

func (evidence RetryFailureEvidence) validate() error {
	if evidence.Code == "" && evidence.Category == "" {
		if evidence.Retryable {
			return classifiedFailure(codeRetryStateInvalid, failure.CategoryValidation, "unknown retry failure evidence cannot be marked retryable")
		}
		return nil
	}
	if !evidence.Code.Valid() || !evidence.Category.Valid() {
		return classifiedFailure(codeRetryStateInvalid, failure.CategoryValidation, "retry failure evidence is malformed")
	}
	return nil
}

// RetryStateRecord is the provider-neutral durable row contract reserved by
// kernel.events migration 3. Identity reuses the accepted P04.05 consumer scope;
// Position additionally binds the exact ordered P04.03 delivery.
type RetryStateRecord struct {
	Identity         InboxIdentity
	Fingerprint      InboxFingerprint
	Position         uint64
	Policy           RetryPolicy
	AttemptsConsumed uint32
	State            RetryState
	NextEligibleAt   time.Time
	TerminalReason   RetryTerminalReason
	Failure          RetryFailureEvidence
	ClaimToken       string
	ClaimExpiresAt   time.Time
	Revision         uint64
}

// Validate rejects inconsistent scheduling/quarantine/claim evidence before a
// persistence adapter can represent it as authoritative state.
func (record RetryStateRecord) Validate() error {
	if err := record.Identity.validate(); err != nil {
		return err
	}
	if record.Fingerprint == (InboxFingerprint{}) {
		return classifiedFailure(codeRetryStateInvalid, failure.CategoryValidation, "retry state canonical fingerprint is required")
	}
	if record.Position == 0 || record.Revision == 0 {
		return classifiedFailure(codeRetryStateInvalid, failure.CategoryValidation, "retry state position and revision must be positive")
	}
	if err := record.Policy.Validate(); err != nil {
		return err
	}
	if record.AttemptsConsumed == 0 || record.AttemptsConsumed > record.Policy.MaxAttempts {
		return classifiedFailure(codeRetryStateInvalid, failure.CategoryValidation, "retry state consumed attempts are outside the policy snapshot")
	}
	if !record.State.valid() {
		return classifiedFailure(codeRetryStateInvalid, failure.CategoryValidation, "retry durable state is invalid")
	}
	if err := record.Failure.validate(); err != nil {
		return err
	}
	if err := validateRetryClaim(record); err != nil {
		return err
	}

	switch record.State {
	case RetryStateScheduled:
		if record.AttemptsConsumed >= record.Policy.MaxAttempts {
			return classifiedFailure(codeRetryStateMalformed, failure.CategoryInvariant, "scheduled retry state has exhausted its policy budget")
		}
		if !validRetryUTC(record.NextEligibleAt) || record.TerminalReason != RetryTerminalReasonNone {
			return classifiedFailure(codeRetryStateMalformed, failure.CategoryInvariant, "scheduled retry state has inconsistent eligibility or terminal evidence")
		}
		if record.ClaimToken != "" && !record.ClaimExpiresAt.After(record.NextEligibleAt) {
			return classifiedFailure(codeRetryStateMalformed, failure.CategoryInvariant, "retry claim expiry must be later than scheduled eligibility")
		}
	case RetryStateQuarantined:
		if !record.NextEligibleAt.IsZero() || record.ClaimToken != "" || !record.ClaimExpiresAt.IsZero() || !retryTerminalReasonValid(record.TerminalReason) {
			return classifiedFailure(codeRetryStateMalformed, failure.CategoryInvariant, "quarantined retry state has inconsistent scheduling, claim, or terminal evidence")
		}
	case RetryStateResolved:
		if !record.NextEligibleAt.IsZero() || record.TerminalReason != RetryTerminalReasonNone || record.ClaimToken != "" || !record.ClaimExpiresAt.IsZero() {
			return classifiedFailure(codeRetryStateMalformed, failure.CategoryInvariant, "resolved retry state must not retain active failure scheduling or claim evidence")
		}
	default:
		return classifiedFailure(codeRetryStateInvalid, failure.CategoryValidation, "retry durable state is invalid")
	}
	return nil
}

func validateRetryClaim(record RetryStateRecord) error {
	hasToken := record.ClaimToken != ""
	hasExpiry := !record.ClaimExpiresAt.IsZero()
	if hasToken != hasExpiry {
		return classifiedFailure(codeRetryStateMalformed, failure.CategoryInvariant, "retry claim token and expiry must be present together")
	}
	if !hasToken {
		return nil
	}
	if record.State != RetryStateScheduled || !retryClaimTokenPattern.MatchString(record.ClaimToken) || !validRetryUTC(record.ClaimExpiresAt) {
		return classifiedFailure(codeRetryStateMalformed, failure.CategoryInvariant, "retry claim evidence is malformed or attached to an inactive state")
	}
	return nil
}

func retryTerminalReasonValid(reason RetryTerminalReason) bool {
	switch reason {
	case RetryTerminalReasonNonRetryable,
		RetryTerminalReasonUnsafeOperation,
		RetryTerminalReasonAttemptsExhausted,
		RetryTerminalReasonUnknownFailure:
		return true
	default:
		return false
	}
}

func validRetryUTC(value time.Time) bool {
	return !value.IsZero() && value.Location() == time.UTC && value.Year() >= 1 && value.Year() <= 9999
}

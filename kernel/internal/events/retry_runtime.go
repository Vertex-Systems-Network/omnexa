package events

import (
	"errors"
	"math"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeRetryRuntimeInvalid   failure.Code = "events.retry.runtime_invalid"
	codeRetryRuntimeMalformed failure.Code = "events.retry.runtime_malformed"
	codeRetryRuntimeExhausted failure.Code = "events.retry.runtime_exhausted"
)

// RetryRuntimeOutcome is the provider-neutral result of composing one retry
// state transition. Persist distinguishes an interruption before any retry row
// exists from an interruption after a durable claim, which must release that
// claim through the accepted persistence CAS boundary.
type RetryRuntimeOutcome string

const (
	RetryRuntimeInterrupted RetryRuntimeOutcome = "interrupted"
	RetryRuntimeScheduled   RetryRuntimeOutcome = "scheduled"
	RetryRuntimeQuarantined RetryRuntimeOutcome = "quarantined"
	RetryRuntimeResolved    RetryRuntimeOutcome = "resolved"
)

// RetryRuntimeTransition is a pure state-machine result. It never writes state,
// advances a checkpoint, invokes a handler, sleeps until due time, or selects a
// provider. Callers persist Record only when Persist is true through the already
// accepted retry-state storage boundary.
type RetryRuntimeTransition struct {
	Outcome RetryRuntimeOutcome
	Persist bool
	Record  RetryStateRecord
}

// RetryExecutionEligibility describes what an execution adapter may do with an
// already-authoritative retry-state record. Scheduled work without a durable
// claim is never executable here; authoritative due-time/lease arbitration stays
// with the persistence adapter rather than caller wall-clock time.
type RetryExecutionEligibility string

const (
	RetryExecutionAwaitingClaim RetryExecutionEligibility = "awaiting_claim"
	RetryExecutionClaimed       RetryExecutionEligibility = "claimed"
	RetryExecutionQuarantined   RetryExecutionEligibility = "quarantined"
	RetryExecutionResolved      RetryExecutionEligibility = "resolved"
)

// RetryExecutionDirective is bounded execution/eligibility state. NextAttempt
// is populated only for scheduled records and counts the authoritative handler
// attempt that a successful durable claim would permit.
type RetryExecutionDirective struct {
	Eligibility    RetryExecutionEligibility
	NextAttempt    uint32
	NextEligibleAt time.Time
}

// ComposeInitialRetryFailure converts the first authoritative failed delivery
// into the durable retry/quarantine state that must be stored before the caller
// may represent the delivery as deferred or terminal. Cancellation/interruption
// before an authoritative result creates no retry row and consumes no attempt.
func ComposeInitialRetryFailure(
	inbox InboxRecord,
	position uint64,
	policy RetryPolicy,
	failureTime time.Time,
	failureErr error,
	operationRetrySafe bool,
) (RetryRuntimeTransition, error) {
	if err := inbox.validate(); err != nil || position == 0 {
		return RetryRuntimeTransition{}, classifiedFailure(codeRetryRuntimeInvalid, failure.CategoryValidation, "event retry runtime initial boundary is invalid")
	}

	decision, err := DecideRetry(policy, 1, failureTime, failureErr, operationRetrySafe)
	if err != nil {
		return RetryRuntimeTransition{}, err
	}
	if decision.Disposition == RetryDispositionInterrupted {
		if decision.ConsumeAttempt {
			return RetryRuntimeTransition{}, classifiedFailure(codeRetryRuntimeMalformed, failure.CategoryInvariant, "event retry interruption consumed an attempt")
		}
		return RetryRuntimeTransition{Outcome: RetryRuntimeInterrupted}, nil
	}
	if !decision.ConsumeAttempt {
		return RetryRuntimeTransition{}, classifiedFailure(codeRetryRuntimeMalformed, failure.CategoryInvariant, "event retry authoritative failure did not consume an attempt")
	}

	record := RetryStateRecord{
		Identity:         inbox.Identity,
		Fingerprint:      inbox.Fingerprint,
		Position:         position,
		Policy:           policy,
		AttemptsConsumed: 1,
		Failure:          retryRuntimeFailureEvidence(failureErr),
		Revision:         1,
	}

	outcome, err := applyRetryDecisionToRecord(&record, decision)
	if err != nil {
		return RetryRuntimeTransition{}, err
	}
	if validationErr := record.Validate(); validationErr != nil {
		return RetryRuntimeTransition{}, validationErr
	}
	return RetryRuntimeTransition{Outcome: outcome, Persist: true, Record: record}, nil
}

// RetryExecutionDirectiveFor interprets one accepted durable record without
// making a wall-clock due decision. A scheduled unclaimed row remains awaiting
// claim even if its NextEligibleAt is in the past; only the durable claim/CAS
// adapter may turn it into executable state using authoritative time.
func RetryExecutionDirectiveFor(record RetryStateRecord) (RetryExecutionDirective, error) {
	if err := record.Validate(); err != nil {
		return RetryExecutionDirective{}, err
	}

	switch record.State {
	case RetryStateScheduled:
		if record.AttemptsConsumed >= record.Policy.MaxAttempts {
			return RetryExecutionDirective{}, classifiedFailure(codeRetryRuntimeMalformed, failure.CategoryInvariant, "event retry scheduled state has no remaining attempt")
		}
		directive := RetryExecutionDirective{
			Eligibility:    RetryExecutionAwaitingClaim,
			NextAttempt:    record.AttemptsConsumed + 1,
			NextEligibleAt: record.NextEligibleAt,
		}
		if record.ClaimToken != "" {
			directive.Eligibility = RetryExecutionClaimed
		}
		return directive, nil
	case RetryStateQuarantined:
		return RetryExecutionDirective{Eligibility: RetryExecutionQuarantined}, nil
	case RetryStateResolved:
		return RetryExecutionDirective{Eligibility: RetryExecutionResolved}, nil
	default:
		return RetryExecutionDirective{}, classifiedFailure(codeRetryRuntimeMalformed, failure.CategoryInvariant, "event retry runtime received unknown durable state")
	}
}

// ComposeClaimedRetryResult composes the state that the accepted CAS persistence
// boundary must commit after one durably claimed retry execution. The claim is
// always cleared in the candidate state. A successful attempt resolves retry
// scheduling; interruption releases the claim without consuming budget; a failed
// authoritative attempt is reclassified through the accepted retry policy.
func ComposeClaimedRetryResult(
	claimed RetryStateRecord,
	failureTime time.Time,
	failureErr error,
	operationRetrySafe bool,
) (RetryRuntimeTransition, error) {
	directive, err := RetryExecutionDirectiveFor(claimed)
	if err != nil {
		return RetryRuntimeTransition{}, err
	}
	if directive.Eligibility != RetryExecutionClaimed {
		return RetryRuntimeTransition{}, classifiedFailure(codeRetryRuntimeInvalid, failure.CategoryConflict, "event retry execution requires an authoritative claim")
	}
	if claimed.Revision == math.MaxUint64 {
		return RetryRuntimeTransition{}, classifiedFailure(codeRetryRuntimeExhausted, failure.CategoryInvariant, "event retry state revision is exhausted")
	}

	if failureErr == nil {
		next := claimed
		next.AttemptsConsumed = directive.NextAttempt
		next.State = RetryStateResolved
		next.NextEligibleAt = time.Time{}
		next.TerminalReason = RetryTerminalReasonNone
		next.ClaimToken = ""
		next.ClaimExpiresAt = time.Time{}
		next.Revision++
		if validationErr := next.Validate(); validationErr != nil {
			return RetryRuntimeTransition{}, validationErr
		}
		return RetryRuntimeTransition{Outcome: RetryRuntimeResolved, Persist: true, Record: next}, nil
	}

	decision, err := DecideRetry(
		claimed.Policy,
		directive.NextAttempt,
		failureTime,
		failureErr,
		operationRetrySafe,
	)
	if err != nil {
		return RetryRuntimeTransition{}, err
	}

	if decision.Disposition == RetryDispositionInterrupted {
		if decision.ConsumeAttempt {
			return RetryRuntimeTransition{}, classifiedFailure(codeRetryRuntimeMalformed, failure.CategoryInvariant, "event retry interruption consumed an attempt")
		}
		next := claimed
		next.ClaimToken = ""
		next.ClaimExpiresAt = time.Time{}
		next.Revision++
		if validationErr := next.Validate(); validationErr != nil {
			return RetryRuntimeTransition{}, validationErr
		}
		return RetryRuntimeTransition{Outcome: RetryRuntimeInterrupted, Persist: true, Record: next}, nil
	}
	if !decision.ConsumeAttempt {
		return RetryRuntimeTransition{}, classifiedFailure(codeRetryRuntimeMalformed, failure.CategoryInvariant, "event retry authoritative failure did not consume an attempt")
	}

	next := claimed
	next.AttemptsConsumed = directive.NextAttempt
	next.Failure = retryRuntimeFailureEvidence(failureErr)
	next.ClaimToken = ""
	next.ClaimExpiresAt = time.Time{}
	next.Revision++
	outcome, err := applyRetryDecisionToRecord(&next, decision)
	if err != nil {
		return RetryRuntimeTransition{}, err
	}
	if validationErr := next.Validate(); validationErr != nil {
		return RetryRuntimeTransition{}, validationErr
	}
	return RetryRuntimeTransition{Outcome: outcome, Persist: true, Record: next}, nil
}

func applyRetryDecisionToRecord(record *RetryStateRecord, decision RetryDecision) (RetryRuntimeOutcome, error) {
	if record == nil {
		return "", classifiedFailure(codeRetryRuntimeInvalid, failure.CategoryValidation, "event retry runtime transition target is invalid")
	}

	switch decision.Disposition {
	case RetryDispositionRetryable:
		if decision.TerminalReason != RetryTerminalReasonNone || decision.NextEligibleAt.IsZero() || decision.Delay <= 0 {
			return "", classifiedFailure(codeRetryRuntimeMalformed, failure.CategoryInvariant, "event retry scheduling decision is malformed")
		}
		record.State = RetryStateScheduled
		record.NextEligibleAt = decision.NextEligibleAt
		record.TerminalReason = RetryTerminalReasonNone
		return RetryRuntimeScheduled, nil
	case RetryDispositionTerminal:
		if !retryTerminalReasonValid(decision.TerminalReason) || !decision.NextEligibleAt.IsZero() || decision.Delay != 0 {
			return "", classifiedFailure(codeRetryRuntimeMalformed, failure.CategoryInvariant, "event retry terminal decision is malformed")
		}
		record.State = RetryStateQuarantined
		record.NextEligibleAt = time.Time{}
		record.TerminalReason = decision.TerminalReason
		return RetryRuntimeQuarantined, nil
	default:
		return "", classifiedFailure(codeRetryRuntimeMalformed, failure.CategoryInvariant, "event retry runtime cannot persist this disposition")
	}
}

func retryRuntimeFailureEvidence(err error) RetryFailureEvidence {
	var structured *failure.Error
	if !errors.As(err, &structured) {
		return RetryFailureEvidence{}
	}
	return RetryFailureEvidence{
		Code:      structured.Code(),
		Category:  structured.Category(),
		Retryable: structured.Retryable(),
	}
}

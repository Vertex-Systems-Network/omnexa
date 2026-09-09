package events

import (
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeRetryInboxInvalid    failure.Code = "events.retry.inbox_invalid"
	codeRetryInboxConflict   failure.Code = "events.retry.inbox_conflict"
	codeRetryInboxConcurrent failure.Code = "events.retry.inbox_concurrent"
	codeRetryInboxMalformed  failure.Code = "events.retry.inbox_malformed"
)

// ComposeClaimedRetryInboxResult composes an already-authoritative P04.05 inbox
// result with one durably claimed P04.06 retry attempt. InboxApplied and
// InboxAlreadyApplied are both success-like outcomes: stale retry scheduling is
// resolved through the accepted claimed-retry success path. This seam has no
// protected-mutation callback, so it cannot re-run a mutation that P04.05 has
// already classified as applied.
//
// Conflict, concurrent and unknown inbox outcomes fail closed and never resolve,
// schedule or quarantine retry state. A durable retry claim remains mandatory;
// inbox outcome evidence is not execution authority.
func ComposeClaimedRetryInboxResult(
	claimed RetryStateRecord,
	result InboxApplyResult,
) (RetryRuntimeTransition, error) {
	directive, err := RetryExecutionDirectiveFor(claimed)
	if err != nil {
		return RetryRuntimeTransition{}, err
	}
	if directive.Eligibility != RetryExecutionClaimed {
		return RetryRuntimeTransition{}, classifiedFailure(
			codeRetryInboxInvalid,
			failure.CategoryConflict,
			"event retry inbox composition requires an authoritative retry claim",
		)
	}

	switch result {
	case InboxApplied, InboxAlreadyApplied:
		return ComposeClaimedRetryResult(claimed, time.Time{}, nil, true)
	case InboxConflict:
		return RetryRuntimeTransition{}, classifiedFailure(
			codeRetryInboxConflict,
			failure.CategoryConflict,
			"event retry inbox result conflicts with committed canonical evidence",
		)
	case InboxConcurrent:
		return RetryRuntimeTransition{}, classifiedFailure(
			codeRetryInboxConcurrent,
			failure.CategoryConflict,
			"event retry inbox result is still being resolved concurrently",
		)
	case InboxApplyResultUnknown:
		fallthrough
	default:
		return RetryRuntimeTransition{}, classifiedFailure(
			codeRetryInboxMalformed,
			failure.CategoryInvariant,
			"event retry inbox result is invalid",
		)
	}
}

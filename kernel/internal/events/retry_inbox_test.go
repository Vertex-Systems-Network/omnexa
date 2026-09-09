package events

import (
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestComposeClaimedRetryInboxResultAppliedAndAlreadyAppliedResolve(t *testing.T) {
	for _, result := range []InboxApplyResult{InboxApplied, InboxAlreadyApplied} {
		t.Run(retryInboxResultName(result), func(t *testing.T) {
			claimed := testClaimedRetryRuntimeState(t)
			transition, err := ComposeClaimedRetryInboxResult(claimed, result)
			if err != nil {
				t.Fatalf("compose inbox success result: %v", err)
			}
			if transition.Outcome != RetryRuntimeResolved || !transition.Persist {
				t.Fatalf("inbox success did not resolve retry scheduling: %+v", transition)
			}
			next := transition.Record
			if next.State != RetryStateResolved {
				t.Fatalf("inbox success state = %q, want resolved", next.State)
			}
			if next.AttemptsConsumed != claimed.AttemptsConsumed+1 || next.Revision != claimed.Revision+1 {
				t.Fatalf("inbox success attempt/revision drifted: before=%+v after=%+v", claimed, next)
			}
			if !next.NextEligibleAt.IsZero() || next.TerminalReason != RetryTerminalReasonNone {
				t.Fatalf("resolved inbox result retained retry/quarantine evidence: %+v", next)
			}
			if next.ClaimToken != "" || !next.ClaimExpiresAt.IsZero() {
				t.Fatalf("resolved inbox result retained a claim: %+v", next)
			}
			if next.Failure != claimed.Failure {
				t.Fatalf("resolved inbox result rewrote bounded historical failure evidence: before=%+v after=%+v", claimed.Failure, next.Failure)
			}
		})
	}
}

func TestComposeClaimedRetryInboxResultFailsClosedForNonSuccessOutcomes(t *testing.T) {
	tests := []struct {
		name   string
		result InboxApplyResult
		code   failure.Code
	}{
		{"conflict", InboxConflict, codeRetryInboxConflict},
		{"concurrent", InboxConcurrent, codeRetryInboxConcurrent},
		{"unknown", InboxApplyResultUnknown, codeRetryInboxMalformed},
		{"future-invalid", InboxApplyResult(255), codeRetryInboxMalformed},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			claimed := testClaimedRetryRuntimeState(t)
			transition, err := ComposeClaimedRetryInboxResult(claimed, test.result)
			if !failure.IsCode(err, test.code) {
				t.Fatalf("error = %v, want code %s", err, test.code)
			}
			if transition != (RetryRuntimeTransition{}) {
				t.Fatalf("failed inbox outcome created retry state: %+v", transition)
			}
		})
	}
}

func TestComposeClaimedRetryInboxResultRequiresAuthoritativeClaim(t *testing.T) {
	unclaimed := testRetryStateRecord(t, RetryStateScheduled)
	transition, err := ComposeClaimedRetryInboxResult(unclaimed, InboxAlreadyApplied)
	if !failure.IsCode(err, codeRetryInboxInvalid) {
		t.Fatalf("unclaimed already-applied error = %v, want %s", err, codeRetryInboxInvalid)
	}
	if transition != (RetryRuntimeTransition{}) {
		t.Fatalf("unclaimed already-applied result resolved retry state: %+v", transition)
	}
}

func retryInboxResultName(result InboxApplyResult) string {
	switch result {
	case InboxApplied:
		return "applied"
	case InboxAlreadyApplied:
		return "already-applied"
	default:
		return "other"
	}
}

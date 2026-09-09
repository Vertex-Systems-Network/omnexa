package events

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestComposeInitialRetryFailureSchedulesDurableStateBeforeDeferral(t *testing.T) {
	policy := testRetryPolicy()
	failureTime := time.Date(2026, 9, 9, 12, 0, 0, 0, time.FixedZone("local", 5*60*60))
	transition, err := ComposeInitialRetryFailure(
		testRetryRuntimeInbox(t),
		9,
		policy,
		failureTime,
		mustRetryFailure(t, failure.CategoryDependency, true),
		true,
	)
	if err != nil {
		t.Fatalf("compose initial retry failure: %v", err)
	}
	if transition.Outcome != RetryRuntimeScheduled || !transition.Persist {
		t.Fatalf("unexpected initial transition: %+v", transition)
	}
	record := transition.Record
	if record.State != RetryStateScheduled || record.AttemptsConsumed != 1 || record.Revision != 1 {
		t.Fatalf("unexpected scheduled record: %+v", record)
	}
	if record.NextEligibleAt.Location() != time.UTC || !record.NextEligibleAt.Equal(failureTime.UTC().Add(policy.InitialBackoff)) {
		t.Fatalf("next eligibility = %s, want %s", record.NextEligibleAt, failureTime.UTC().Add(policy.InitialBackoff))
	}
	if record.Failure.Code != "events.test.failure" || record.Failure.Category != failure.CategoryDependency || !record.Failure.Retryable {
		t.Fatalf("structured failure evidence drifted: %+v", record.Failure)
	}
	if record.ClaimToken != "" || !record.ClaimExpiresAt.IsZero() {
		t.Fatalf("initial scheduled state unexpectedly contains claim evidence: %+v", record)
	}
}

func TestComposeInitialRetryFailureInterruptionCreatesNoDurableState(t *testing.T) {
	transition, err := ComposeInitialRetryFailure(
		testRetryRuntimeInbox(t),
		9,
		testRetryPolicy(),
		time.Time{},
		context.Canceled,
		true,
	)
	if err != nil {
		t.Fatalf("compose interrupted initial attempt: %v", err)
	}
	if transition.Outcome != RetryRuntimeInterrupted || transition.Persist || transition.Record != (RetryStateRecord{}) {
		t.Fatalf("interruption created false durable state: %+v", transition)
	}
}

func TestComposeInitialRetryFailureUnknownErrorFailsClosedToQuarantine(t *testing.T) {
	transition, err := ComposeInitialRetryFailure(
		testRetryRuntimeInbox(t),
		9,
		testRetryPolicy(),
		time.Time{},
		errors.New("opaque provider detail"),
		true,
	)
	if err != nil {
		t.Fatalf("compose unknown initial failure: %v", err)
	}
	if transition.Outcome != RetryRuntimeQuarantined || !transition.Persist {
		t.Fatalf("unknown failure did not become durable quarantine: %+v", transition)
	}
	if transition.Record.State != RetryStateQuarantined || transition.Record.TerminalReason != RetryTerminalReasonUnknownFailure {
		t.Fatalf("unexpected quarantine state: %+v", transition.Record)
	}
	if transition.Record.Failure != (RetryFailureEvidence{}) {
		t.Fatalf("unknown failure leaked structured evidence: %+v", transition.Record.Failure)
	}
}

func TestRetryExecutionDirectiveRequiresAuthoritativeClaimBeforeExecution(t *testing.T) {
	record := testRetryStateRecord(t, RetryStateScheduled)
	directive, err := RetryExecutionDirectiveFor(record)
	if err != nil {
		t.Fatalf("inspect scheduled record: %v", err)
	}
	if directive.Eligibility != RetryExecutionAwaitingClaim || directive.NextAttempt != record.AttemptsConsumed+1 {
		t.Fatalf("unclaimed retry became executable: %+v", directive)
	}
	if !directive.NextEligibleAt.Equal(record.NextEligibleAt) {
		t.Fatalf("eligibility evidence drifted: %+v", directive)
	}

	claimed := record
	claimed.ClaimToken = "01990f6e-1f30-7000-8000-000000000901"
	claimed.ClaimExpiresAt = claimed.NextEligibleAt.Add(time.Minute)
	directive, err = RetryExecutionDirectiveFor(claimed)
	if err != nil {
		t.Fatalf("inspect claimed record: %v", err)
	}
	if directive.Eligibility != RetryExecutionClaimed || directive.NextAttempt != claimed.AttemptsConsumed+1 {
		t.Fatalf("authoritative claim did not permit the next attempt: %+v", directive)
	}
}

func TestRetryExecutionDirectiveStopsQuarantinedAndResolvedState(t *testing.T) {
	for _, test := range []struct {
		state RetryState
		want  RetryExecutionEligibility
	}{
		{RetryStateQuarantined, RetryExecutionQuarantined},
		{RetryStateResolved, RetryExecutionResolved},
	} {
		directive, err := RetryExecutionDirectiveFor(testRetryStateRecord(t, test.state))
		if err != nil {
			t.Fatalf("inspect %s record: %v", test.state, err)
		}
		if directive.Eligibility != test.want || directive.NextAttempt != 0 || !directive.NextEligibleAt.IsZero() {
			t.Fatalf("unexpected %s directive: %+v", test.state, directive)
		}
	}
}

func TestComposeClaimedRetryResultSchedulesNextFailureAndClearsClaim(t *testing.T) {
	claimed := testClaimedRetryRuntimeState(t)
	failureTime := time.Date(2026, 9, 9, 8, 0, 0, 0, time.UTC)
	transition, err := ComposeClaimedRetryResult(
		claimed,
		failureTime,
		mustRetryFailure(t, failure.CategoryRateLimit, true),
		true,
	)
	if err != nil {
		t.Fatalf("compose claimed retry failure: %v", err)
	}
	if transition.Outcome != RetryRuntimeScheduled || !transition.Persist {
		t.Fatalf("unexpected retry transition: %+v", transition)
	}
	next := transition.Record
	if next.AttemptsConsumed != claimed.AttemptsConsumed+1 || next.Revision != claimed.Revision+1 {
		t.Fatalf("attempt/revision did not advance exactly once: %+v", next)
	}
	wantEligibility := failureTime.Add(8 * time.Second)
	if !next.NextEligibleAt.Equal(wantEligibility) {
		t.Fatalf("next eligibility = %s, want %s", next.NextEligibleAt, wantEligibility)
	}
	if next.Failure.Category != failure.CategoryRateLimit || !next.Failure.Retryable {
		t.Fatalf("new failure evidence not retained safely: %+v", next.Failure)
	}
	if next.ClaimToken != "" || !next.ClaimExpiresAt.IsZero() {
		t.Fatalf("transition retained stale claim evidence: %+v", next)
	}
}

func TestComposeClaimedRetryResultSuccessResolvesScheduling(t *testing.T) {
	claimed := testClaimedRetryRuntimeState(t)
	transition, err := ComposeClaimedRetryResult(claimed, time.Time{}, nil, true)
	if err != nil {
		t.Fatalf("compose claimed success: %v", err)
	}
	if transition.Outcome != RetryRuntimeResolved || !transition.Persist {
		t.Fatalf("success did not resolve retry state: %+v", transition)
	}
	next := transition.Record
	if next.State != RetryStateResolved || next.AttemptsConsumed != claimed.AttemptsConsumed+1 || next.Revision != claimed.Revision+1 {
		t.Fatalf("unexpected resolved state: %+v", next)
	}
	if !next.NextEligibleAt.IsZero() || next.TerminalReason != RetryTerminalReasonNone || next.ClaimToken != "" || !next.ClaimExpiresAt.IsZero() {
		t.Fatalf("resolved state retained active scheduling/claim evidence: %+v", next)
	}
	if next.Failure != claimed.Failure {
		t.Fatalf("resolved state discarded bounded historical failure evidence: before=%+v after=%+v", claimed.Failure, next.Failure)
	}
}

func TestComposeClaimedRetryResultInterruptionReleasesClaimWithoutConsumingBudget(t *testing.T) {
	claimed := testClaimedRetryRuntimeState(t)
	transition, err := ComposeClaimedRetryResult(claimed, time.Time{}, context.Canceled, true)
	if err != nil {
		t.Fatalf("compose claimed interruption: %v", err)
	}
	if transition.Outcome != RetryRuntimeInterrupted || !transition.Persist {
		t.Fatalf("claimed interruption must persist claim release: %+v", transition)
	}
	next := transition.Record
	if next.State != RetryStateScheduled || next.AttemptsConsumed != claimed.AttemptsConsumed {
		t.Fatalf("interruption consumed retry budget: %+v", next)
	}
	if !next.NextEligibleAt.Equal(claimed.NextEligibleAt) || next.Failure != claimed.Failure {
		t.Fatalf("interruption rewrote prior scheduling/failure evidence: %+v", next)
	}
	if next.Revision != claimed.Revision+1 || next.ClaimToken != "" || !next.ClaimExpiresAt.IsZero() {
		t.Fatalf("interruption did not release claim through a new revision: %+v", next)
	}
}

func TestComposeClaimedRetryResultExhaustionQuarantines(t *testing.T) {
	claimed := testClaimedRetryRuntimeState(t)
	claimed.AttemptsConsumed = claimed.Policy.MaxAttempts - 1
	if err := claimed.Validate(); err != nil {
		t.Fatalf("prepare last eligible claimed state: %v", err)
	}
	transition, err := ComposeClaimedRetryResult(
		claimed,
		time.Now(),
		mustRetryFailure(t, failure.CategoryUnavailable, true),
		true,
	)
	if err != nil {
		t.Fatalf("compose exhausted retry: %v", err)
	}
	if transition.Outcome != RetryRuntimeQuarantined || transition.Record.State != RetryStateQuarantined {
		t.Fatalf("exhaustion did not quarantine: %+v", transition)
	}
	if transition.Record.TerminalReason != RetryTerminalReasonAttemptsExhausted || transition.Record.AttemptsConsumed != claimed.Policy.MaxAttempts {
		t.Fatalf("exhaustion evidence drifted: %+v", transition.Record)
	}
}

func TestComposeClaimedRetryResultRejectsUnclaimedOrExhaustedRevision(t *testing.T) {
	unclaimed := testRetryStateRecord(t, RetryStateScheduled)
	if _, err := ComposeClaimedRetryResult(unclaimed, time.Time{}, nil, true); err == nil {
		t.Fatal("unclaimed state unexpectedly executed")
	}

	claimed := testClaimedRetryRuntimeState(t)
	claimed.Revision = math.MaxUint64
	if err := claimed.Validate(); err != nil {
		t.Fatalf("provider-neutral record validation unexpectedly rejected max revision: %v", err)
	}
	if _, err := ComposeClaimedRetryResult(claimed, time.Time{}, nil, true); err == nil {
		t.Fatal("revision overflow unexpectedly composed")
	}
}

func testRetryRuntimeInbox(t *testing.T) InboxRecord {
	t.Helper()
	envelope := testEnvelope(t)
	record, err := NewInboxRecord(testInboxBinding(envelope, "billing.projection"), envelope)
	if err != nil {
		t.Fatalf("build retry runtime inbox evidence: %v", err)
	}
	return record
}

func testClaimedRetryRuntimeState(t *testing.T) RetryStateRecord {
	t.Helper()
	record := testRetryStateRecord(t, RetryStateScheduled)
	record.ClaimToken = "01990f6e-1f30-7000-8000-000000000901"
	record.ClaimExpiresAt = record.NextEligibleAt.Add(time.Minute)
	if err := record.Validate(); err != nil {
		t.Fatalf("prepare claimed retry state: %v", err)
	}
	return record
}

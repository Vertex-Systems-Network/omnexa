package events

import (
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestRetryStateScheduledAcceptsBoundedUTCEligibilityAndOptionalClaim(t *testing.T) {
	record := testRetryStateRecord(t, RetryStateScheduled)
	if err := record.Validate(); err != nil {
		t.Fatalf("valid scheduled state failed: %v", err)
	}

	record.ClaimToken = "01990f6e-1f30-4000-8000-000000000901"
	record.ClaimExpiresAt = record.NextEligibleAt.Add(30 * time.Second)
	if err := record.Validate(); err != nil {
		t.Fatalf("valid claimed scheduled state failed: %v", err)
	}
}

func TestRetryStateScheduledRejectsMalformedBudgetEligibilityAndClaimEvidence(t *testing.T) {
	base := testRetryStateRecord(t, RetryStateScheduled)
	tests := []struct {
		name   string
		mutate func(*RetryStateRecord)
	}{
		{"missing fingerprint", func(record *RetryStateRecord) { record.Fingerprint = InboxFingerprint{} }},
		{"zero position", func(record *RetryStateRecord) { record.Position = 0 }},
		{"zero revision", func(record *RetryStateRecord) { record.Revision = 0 }},
		{"zero attempts", func(record *RetryStateRecord) { record.AttemptsConsumed = 0 }},
		{"attempts above policy", func(record *RetryStateRecord) { record.AttemptsConsumed = record.Policy.MaxAttempts + 1 }},
		{"exhausted scheduled state", func(record *RetryStateRecord) { record.AttemptsConsumed = record.Policy.MaxAttempts }},
		{"missing eligibility", func(record *RetryStateRecord) { record.NextEligibleAt = time.Time{} }},
		{"non utc eligibility", func(record *RetryStateRecord) { record.NextEligibleAt = time.Date(2026, 9, 7, 13, 0, 0, 0, time.FixedZone("local", 5*60*60)) }},
		{"terminal reason on scheduled state", func(record *RetryStateRecord) { record.TerminalReason = RetryTerminalReasonNonRetryable }},
		{"claim token without expiry", func(record *RetryStateRecord) { record.ClaimToken = "01990f6e-1f30-4000-8000-000000000901" }},
		{"claim expiry without token", func(record *RetryStateRecord) { record.ClaimExpiresAt = record.NextEligibleAt.Add(time.Minute) }},
		{"malformed claim token", func(record *RetryStateRecord) {
			record.ClaimToken = "not-a-uuid"
			record.ClaimExpiresAt = record.NextEligibleAt.Add(time.Minute)
		}},
		{"claim expires before eligibility", func(record *RetryStateRecord) {
			record.ClaimToken = "01990f6e-1f30-4000-8000-000000000901"
			record.ClaimExpiresAt = record.NextEligibleAt
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			record := base
			test.mutate(&record)
			if err := record.Validate(); err == nil {
				t.Fatalf("malformed state unexpectedly passed: %+v", record)
			}
		})
	}
}

func TestRetryStateQuarantineRequiresFiniteTerminalEvidenceAndNoActiveClaim(t *testing.T) {
	record := testRetryStateRecord(t, RetryStateQuarantined)
	if err := record.Validate(); err != nil {
		t.Fatalf("valid quarantine failed: %v", err)
	}

	for _, mutate := range []func(*RetryStateRecord){
		func(candidate *RetryStateRecord) { candidate.TerminalReason = RetryTerminalReasonNone },
		func(candidate *RetryStateRecord) { candidate.TerminalReason = RetryTerminalReason("future_unaccepted_reason") },
		func(candidate *RetryStateRecord) { candidate.NextEligibleAt = time.Date(2026, 9, 7, 8, 0, 0, 0, time.UTC) },
		func(candidate *RetryStateRecord) {
			candidate.ClaimToken = "01990f6e-1f30-4000-8000-000000000901"
			candidate.ClaimExpiresAt = time.Date(2026, 9, 7, 8, 1, 0, 0, time.UTC)
		},
	} {
		candidate := record
		mutate(&candidate)
		if err := candidate.Validate(); err == nil {
			t.Fatalf("inconsistent quarantine unexpectedly passed: %+v", candidate)
		}
	}
}

func TestRetryStateResolvedCarriesNoActiveSchedulingQuarantineOrClaimEvidence(t *testing.T) {
	record := testRetryStateRecord(t, RetryStateResolved)
	if err := record.Validate(); err != nil {
		t.Fatalf("valid resolved state failed: %v", err)
	}

	candidate := record
	candidate.NextEligibleAt = time.Date(2026, 9, 7, 8, 0, 0, 0, time.UTC)
	if err := candidate.Validate(); err == nil {
		t.Fatal("resolved state retained next eligibility")
	}
	candidate = record
	candidate.TerminalReason = RetryTerminalReasonAttemptsExhausted
	if err := candidate.Validate(); err == nil {
		t.Fatal("resolved state retained terminal reason")
	}
}

func TestRetryFailureEvidenceIsStructuredOrExplicitlyUnknown(t *testing.T) {
	record := testRetryStateRecord(t, RetryStateScheduled)
	record.Failure = RetryFailureEvidence{}
	if err := record.Validate(); err != nil {
		t.Fatalf("explicit unknown failure evidence failed: %v", err)
	}

	record.Failure = RetryFailureEvidence{Retryable: true}
	if err := record.Validate(); err == nil {
		t.Fatal("unknown failure evidence was marked retryable")
	}

	record.Failure = RetryFailureEvidence{Code: "events.test.failure", Category: failure.Category("not-a-category")}
	if err := record.Validate(); err == nil {
		t.Fatal("malformed structured failure evidence passed")
	}
}

func testRetryStateRecord(t *testing.T, state RetryState) RetryStateRecord {
	t.Helper()
	envelope := testEnvelope(t)
	inbox, err := NewInboxRecord(testInboxBinding(envelope, "billing.projection"), envelope)
	if err != nil {
		t.Fatalf("build inbox evidence: %v", err)
	}

	record := RetryStateRecord{
		Identity:         inbox.Identity,
		Fingerprint:      inbox.Fingerprint,
		Position:         7,
		Policy:           testRetryPolicy(),
		AttemptsConsumed: 2,
		State:            state,
		Failure: RetryFailureEvidence{
			Code:      "events.test.failure",
			Category:  failure.CategoryDependency,
			Retryable: state == RetryStateScheduled,
		},
		Revision: 1,
	}

	switch state {
	case RetryStateScheduled:
		record.NextEligibleAt = time.Date(2026, 9, 7, 8, 0, 0, 0, time.UTC)
	case RetryStateQuarantined:
		record.TerminalReason = RetryTerminalReasonNonRetryable
	case RetryStateResolved:
		// Resolved state retains only historical identity/policy/attempt/failure evidence.
	}
	return record
}

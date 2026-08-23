package audit

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/observability"
	"github.com/google/uuid"
)

func TestRecordEnvelopeUsesUUIDv7UTCAndIsImmutable(t *testing.T) {
	t.Parallel()
	input := validInput()
	input.Fields = []Field{{
		Key:            "subject_hint",
		Value:          "customer@example.test",
		Classification: ClassificationConfidential,
		Tags:           []HandlingTag{TagPII},
	}}

	record, err := NewRecord(input)
	if err != nil {
		t.Fatalf("NewRecord() error = %v", err)
	}
	parsed, err := uuid.Parse(record.ID())
	if err != nil || parsed.Version() != 7 {
		t.Fatalf("record ID = %q, want UUIDv7", record.ID())
	}
	if record.OccurredAt().Location() != time.UTC {
		t.Fatalf("OccurredAt() location = %v, want UTC", record.OccurredAt().Location())
	}
	if err := record.Verify(); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if len(record.IntegrityDigest()) != 64 {
		t.Fatalf("IntegrityDigest() length = %d, want 64", len(record.IntegrityDigest()))
	}

	fields := record.Fields()
	fields[0].Value = "mutated"
	fields[0].Tags[0] = TagLegalRecord
	if got := record.Fields()[0].Value; got != "customer@example.test" {
		t.Fatalf("record field mutated through defensive copy: %q", got)
	}
	if got := record.Fields()[0].Tags[0]; got != TagPII {
		t.Fatalf("record tags mutated through defensive copy: %q", got)
	}
}

func TestRecordRejectsSecretsAndUnsafeClassification(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		mutate func(*RecordInput)
		code  failure.Code
	}{
		{
			name: "sensitive key",
			mutate: func(input *RecordInput) {
				input.Fields = []Field{{Key: "access_token", Value: "value", Classification: ClassificationInternal}}
			},
			code: codeRecordProhibited,
		},
		{
			name: "auth secret tag",
			mutate: func(input *RecordInput) {
				input.Fields = []Field{{Key: "auth_material", Value: "value", Classification: ClassificationRestricted, Tags: []HandlingTag{TagAuthSecret}}}
				input.Classification = ClassificationRestricted
			},
			code: codeRecordProhibited,
		},
		{
			name: "classification downgrade",
			mutate: func(input *RecordInput) {
				input.Fields = []Field{{Key: "customer_hint", Value: "value", Classification: ClassificationRestricted}}
				input.Classification = ClassificationConfidential
			},
			code: codeRecordInvalid,
		},
		{
			name: "ambiguous platform scope",
			mutate: func(input *RecordInput) {
				input.Scope = Scope{Platform: true, TenantID: mustUUIDv7(t)}
			},
			code: codeRecordInvalid,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := validInput()
			test.mutate(&input)
			_, err := NewRecord(input)
			if !failure.IsCode(err, test.code) {
				t.Fatalf("NewRecord() error = %v, want code %s", err, test.code)
			}
		})
	}
}

func TestMemorySinkIsAppendOnlyDeterministicAndRejectsTamper(t *testing.T) {
	t.Parallel()
	sink, err := NewMemorySink(2)
	if err != nil {
		t.Fatalf("NewMemorySink() error = %v", err)
	}
	first, err := NewRecord(validInput())
	if err != nil {
		t.Fatalf("NewRecord(first) error = %v", err)
	}
	secondInput := validInput()
	secondInput.Action = "kernel.audit.second"
	second, err := NewRecord(secondInput)
	if err != nil {
		t.Fatalf("NewRecord(second) error = %v", err)
	}
	if err := sink.Append(context.Background(), first); err != nil {
		t.Fatalf("Append(first) error = %v", err)
	}
	if err := sink.Append(context.Background(), second); err != nil {
		t.Fatalf("Append(second) error = %v", err)
	}
	if got := sink.Snapshot(); len(got) != 2 || got[0].ID() != first.ID() || got[1].ID() != second.ID() {
		t.Fatalf("Snapshot() did not preserve deterministic append order")
	}

	third, err := NewRecord(validInput())
	if err != nil {
		t.Fatalf("NewRecord(third) error = %v", err)
	}
	if err := sink.Append(context.Background(), third); !failure.IsCode(err, codeSinkFull) {
		t.Fatalf("Append(full) error = %v, want %s", err, codeSinkFull)
	}

	tampered := first
	tampered.reason = "tampered after sealing"
	tamperSink, err := NewMemorySink(1)
	if err != nil {
		t.Fatalf("NewMemorySink(tamper) error = %v", err)
	}
	if err := tamperSink.Append(context.Background(), tampered); !failure.IsCode(err, codeRecordIntegrity) {
		t.Fatalf("Append(tampered) error = %v, want %s", err, codeRecordIntegrity)
	}
}

func TestRequiredSinkFailureIsExplicitAndBestEffortIsDegraded(t *testing.T) {
	t.Parallel()
	sinkFailure := errors.New("private backend diagnostic")
	writer, err := NewWriter(failingSink{err: sinkFailure}, nil)
	if err != nil {
		t.Fatalf("NewWriter() error = %v", err)
	}

	receipt, err := writer.Write(context.Background(), RequirementRequired, validInput())
	if !failure.IsCode(err, codeSinkRequiredFailed) {
		t.Fatalf("Write(required) error = %v, want %s", err, codeSinkRequiredFailed)
	}
	if receipt.Status != DeliveryFailed || receipt.Reason != DeliveryReasonSinkFailure || receipt.RecordID == "" {
		t.Fatalf("Write(required) receipt = %+v", receipt)
	}

	receipt, err = writer.Write(context.Background(), RequirementBestEffort, validInput())
	if err != nil {
		t.Fatalf("Write(best-effort) error = %v", err)
	}
	if receipt.Status != DeliveryDegraded || receipt.Reason != DeliveryReasonSinkFailure || receipt.RecordID == "" {
		t.Fatalf("Write(best-effort) receipt = %+v", receipt)
	}

	health := writer.Health()
	if health.Submitted != 2 || health.Recorded != 0 || health.Degraded != 1 || health.Failed != 1 {
		t.Fatalf("Health() = %+v", health)
	}
}

func TestImpersonationReasonApprovalMetadataRepresentedWithoutAuthority(t *testing.T) {
	t.Parallel()
	input := validInput()
	input.Actor.ImpersonatorReference = "support-session-42"
	input.Privileged = true
	input.Reason = "approved emergency diagnostic action"
	input.Approval = &Approval{Kind: "change_ticket", Reference: "chg-1001"}

	record, err := NewRecord(input)
	if err != nil {
		t.Fatalf("NewRecord() error = %v", err)
	}
	if !record.Privileged() || record.Actor().ImpersonatorReference != "support-session-42" || record.Reason() != input.Reason {
		t.Fatalf("privileged/impersonation metadata was not preserved")
	}
	approval, ok := record.Approval()
	if !ok || approval != *input.Approval {
		t.Fatalf("Approval() = %+v, %v", approval, ok)
	}
}

func TestAuditObservabilityExposesNoProtectedPayload(t *testing.T) {
	t.Parallel()
	logger, capture := observability.NewCaptureLogger(observability.Settings{ServiceName: "audit-test"})
	writer, err := NewWriter(failingSink{err: errors.New("token=backend-secret")}, logger)
	if err != nil {
		t.Fatalf("NewWriter() error = %v", err)
	}
	input := validInput()
	input.Reason = "private-reason-marker"
	input.Fields = []Field{{
		Key:            "customer_hint",
		Value:          "private-customer-marker",
		Classification: ClassificationConfidential,
		Tags:           []HandlingTag{TagPII},
	}}

	receipt, err := writer.Write(context.Background(), RequirementBestEffort, input)
	if err != nil || receipt.Status != DeliveryDegraded {
		t.Fatalf("Write() = %+v, %v", receipt, err)
	}
	records := capture.Records()
	if len(records) != 1 {
		t.Fatalf("captured records = %d, want 1", len(records))
	}
	serialized := fmt.Sprint(records[0].Message, records[0].Attributes)
	for _, forbidden := range []string{"private-reason-marker", "private-customer-marker", "backend-secret", input.Actor.Reference, input.Target.Reference} {
		if strings.Contains(serialized, forbidden) {
			t.Fatalf("observability leaked protected audit content %q in %q", forbidden, serialized)
		}
	}
}

func TestCallerCancellationAndSinkPanicFailSafely(t *testing.T) {
	t.Parallel()
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	memory, err := NewMemorySink(1)
	if err != nil {
		t.Fatalf("NewMemorySink() error = %v", err)
	}
	writer, err := NewWriter(memory, nil)
	if err != nil {
		t.Fatalf("NewWriter(memory) error = %v", err)
	}
	if _, err := writer.Write(canceled, RequirementRequired, validInput()); !failure.IsCode(err, codeDeliveryCanceled) {
		t.Fatalf("Write(canceled) error = %v, want %s", err, codeDeliveryCanceled)
	}

	panicWriter, err := NewWriter(failingSink{panic: true}, nil)
	if err != nil {
		t.Fatalf("NewWriter(panic) error = %v", err)
	}
	receipt, err := panicWriter.Write(context.Background(), RequirementRequired, validInput())
	if !failure.IsCode(err, codeSinkRequiredFailed) || receipt.Status != DeliveryFailed {
		t.Fatalf("Write(panic) = %+v, %v", receipt, err)
	}
}

func TestConcurrentAuditWritesAreRaceSafe(t *testing.T) {
	sink, err := NewMemorySink(128)
	if err != nil {
		t.Fatalf("NewMemorySink() error = %v", err)
	}
	writer, err := NewWriter(sink, nil)
	if err != nil {
		t.Fatalf("NewWriter() error = %v", err)
	}

	const writes = 64
	errs := make(chan error, writes)
	var group sync.WaitGroup
	for index := 0; index < writes; index++ {
		group.Add(1)
		go func() {
			defer group.Done()
			receipt, writeErr := writer.Write(context.Background(), RequirementRequired, validInput())
			if writeErr != nil {
				errs <- writeErr
				return
			}
			if receipt.Status != DeliveryRecorded {
				errs <- fmt.Errorf("unexpected delivery status %q", receipt.Status)
			}
		}()
	}
	group.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("concurrent Write() error = %v", err)
	}
	if sink.Len() != writes {
		t.Fatalf("MemorySink.Len() = %d, want %d", sink.Len(), writes)
	}
	health := writer.Health()
	if health.Submitted != writes || health.Recorded != writes || health.Degraded != 0 || health.Failed != 0 {
		t.Fatalf("Health() = %+v", health)
	}
}

type failingSink struct {
	err   error
	panic bool
}

func (sink failingSink) Append(context.Context, Record) error {
	if sink.panic {
		panic("private sink panic")
	}
	return sink.err
}

func validInput() RecordInput {
	return RecordInput{
		Classification: ClassificationConfidential,
		Actor: Actor{
			Kind:      "service",
			Reference: "service-123",
		},
		Action: "kernel.audit.test",
		Target: Target{
			Kind:      "kernel.resource",
			Reference: "resource-123",
		},
		Scope:         Scope{Platform: true},
		Outcome:       OutcomeSucceeded,
		CorrelationID: "correlation-123",
	}
}

func mustUUIDv7(t *testing.T) string {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("uuid.NewV7() error = %v", err)
	}
	return value.String()
}

package identity

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/google/uuid"
)

func TestAuditAdapterRecordsSafeIdentityAndServiceAccountEvents(t *testing.T) {
	sink, err := audit.NewMemorySink(16)
	if err != nil {
		t.Fatalf("audit.NewMemorySink() error = %v", err)
	}
	writer, err := audit.NewWriter(sink, nil)
	if err != nil {
		t.Fatalf("audit.NewWriter() error = %v", err)
	}
	adapter, err := NewAuditAdapter(writer)
	if err != nil {
		t.Fatalf("NewAuditAdapter() error = %v", err)
	}
	userID := UserID(mustAuditUUIDv7(t))
	sessionID := SessionID(mustAuditUUIDv7(t))
	accountID := ServiceAccountID(mustAuditUUIDv7(t))
	credentialID := APICredentialID(mustAuditUUIDv7(t))
	at := time.Date(2026, time.August, 26, 2, 0, 0, 0, time.UTC)

	adapter.RecordSecurityEvent(SecurityAuditEvent{
		Action:      SecurityAuditPasswordChanged,
		PrincipalID: userID,
		SessionID:   sessionID,
		Succeeded:   true,
		OccurredAt:  at,
	})
	adapter.RecordServiceAccountEvent(ServiceAccountAuditEvent{
		Action:           ServiceAccountAuditCredentialVerified,
		ServiceAccountID: accountID,
		CredentialID:     credentialID,
		Succeeded:        false,
		OccurredAt:       at.Add(time.Second),
	})

	records := sink.Snapshot()
	if len(records) != 2 {
		t.Fatalf("audit records = %d, want 2", len(records))
	}
	if records[0].Action() != string(SecurityAuditPasswordChanged) || records[0].Actor().Reference != string(userID) {
		t.Fatalf("human audit attribution = %#v/%q", records[0].Actor(), records[0].Action())
	}
	if records[0].Target().Kind != "identity.session" || records[0].Scope().Platform != true || records[0].Outcome() != audit.OutcomeSucceeded {
		t.Fatalf("human audit target/scope/outcome = %#v/%#v/%q", records[0].Target(), records[0].Scope(), records[0].Outcome())
	}
	if records[1].Action() != string(ServiceAccountAuditCredentialVerified) || records[1].Actor().Kind != identityAuditActorServiceAccount {
		t.Fatalf("service account audit attribution = %#v/%q", records[1].Actor(), records[1].Action())
	}
	if records[1].Target().Reference != string(credentialID) || records[1].Outcome() != audit.OutcomeFailed {
		t.Fatalf("service account audit target/outcome = %#v/%q", records[1].Target(), records[1].Outcome())
	}
	for _, record := range records {
		if record.Classification() != audit.ClassificationConfidential || len(record.Fields()) != 0 {
			t.Fatalf("audit record leaked supplemental fields or class = %q/%#v", record.Classification(), record.Fields())
		}
		if err = record.Verify(); err != nil {
			t.Fatalf("record.Verify() error = %v", err)
		}
	}
}

func TestAuditAdapterBestEffortFailureIsExplicitInWriterHealth(t *testing.T) {
	writer, err := audit.NewWriter(identityFailingSink{}, nil)
	if err != nil {
		t.Fatalf("audit.NewWriter() error = %v", err)
	}
	adapter, err := NewAuditAdapter(writer)
	if err != nil {
		t.Fatalf("NewAuditAdapter() error = %v", err)
	}
	adapter.RecordSecurityEvent(SecurityAuditEvent{
		Action:     SecurityAuditAuthenticationFail,
		Succeeded:  false,
		OccurredAt: time.Date(2026, time.August, 26, 2, 1, 0, 0, time.UTC),
	})
	health := writer.Health()
	if health.Submitted != 1 || health.Degraded != 1 || health.Failed != 0 {
		t.Fatalf("writer health = %#v", health)
	}
}

type identityFailingSink struct{}

func (identityFailingSink) Append(context.Context, audit.Record) error {
	return errors.New("synthetic sink failure")
}

func mustAuditUUIDv7(t *testing.T) string {
	t.Helper()
	identifier, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("uuid.NewV7() error = %v", err)
	}
	return identifier.String()
}

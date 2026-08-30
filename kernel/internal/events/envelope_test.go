package events

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/google/uuid"
)

func TestEnvelopeRoundTripIsDeterministic(t *testing.T) {
	tenantID := tenancy.TenantID(testUUIDv7(t))
	correlationID := CorrelationID(testUUIDv7(t))
	occurredAt := time.Date(2026, time.August, 30, 12, 0, 0, 123456000, time.UTC)

	envelope, createErr := New(Params{
		Type:           EventType("commerce.order.created.v1"),
		Producer:       Producer("urn:omnexa:module:commerce.orders"),
		OccurredAt:     occurredAt,
		TenantID:       tenantID,
		CorrelationID:  correlationID,
		Classification: DataClassInternal,
		DataSchema:     SchemaID("urn:omnexa:event-schema:commerce.order.created:v1"),
		Subject:        "order/ord_123",
		Data: map[string]any{
			"order_id":    "ord_123",
			"total_minor": 12500,
		},
	})
	if createErr != nil {
		t.Fatalf("New() error = %v", createErr)
	}
	if !envelope.ID.Valid() {
		t.Fatalf("generated event id is not UUIDv7: %q", envelope.ID)
	}
	if envelope.EventVersion != EnvelopeVersion {
		t.Fatalf("event version = %d, want %d", envelope.EventVersion, EnvelopeVersion)
	}
	if !envelope.OccurredAt.Equal(occurredAt) || envelope.OccurredAt.Location() != time.UTC {
		t.Fatalf("occurred_at = %v, want canonical UTC %v", envelope.OccurredAt, occurredAt)
	}
	if tenantErr := envelope.ValidateForTenant(tenantID); tenantErr != nil {
		t.Fatalf("ValidateForTenant() error = %v", tenantErr)
	}

	first, marshalErr := envelope.Marshal()
	if marshalErr != nil {
		t.Fatalf("Marshal() error = %v", marshalErr)
	}
	parsed, parseErr := Parse(first)
	if parseErr != nil {
		t.Fatalf("Parse() error = %v", parseErr)
	}
	second, secondMarshalErr := parsed.Marshal()
	if secondMarshalErr != nil {
		t.Fatalf("parsed Marshal() error = %v", secondMarshalErr)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("round-trip bytes differ\nfirst:  %s\nsecond: %s", first, second)
	}
}

func TestEnvelopeRejectsSelfCausationAndFutureEnvelopeVersion(t *testing.T) {
	envelope := testEnvelope(t)
	envelope.CausationID = CausationID(envelope.ID)
	assertFailureCode(t, envelope.Validate(), codeEnvelopeInvalid)

	envelope = testEnvelope(t)
	envelope.Type = EventType("commerce.order.created.v2")
	envelope.EventVersion = 2
	assertFailureCode(t, envelope.Validate(), codeEnvelopeInvalid)
}

func TestEnvelopeTenantBoundaryFailsClosed(t *testing.T) {
	envelope := testEnvelope(t)
	otherTenant := tenancy.TenantID(testUUIDv7(t))
	assertFailureCode(t, envelope.ValidateForTenant(otherTenant), codeTenantMismatch)

	envelope.TenantID = ""
	assertFailureCode(t, envelope.ValidateForTenant(otherTenant), codeTenantMismatch)
}

func TestEnvelopeRejectsSecretLikeAndUnboundedPayloads(t *testing.T) {
	base := testParams(t)

	base.Data = map[string]any{"access_token": "not-allowed"}
	_, err := New(base)
	assertFailureCode(t, err, codePayloadInvalid)

	base = testParams(t)
	base.Data = map[string]any{"message": "Bearer should-not-enter-event-payloads"}
	_, err = New(base)
	assertFailureCode(t, err, codePayloadInvalid)

	base = testParams(t)
	base.Data = map[string]any{"message": strings.Repeat("x", maxStringRunes+1)}
	_, err = New(base)
	assertFailureCode(t, err, codePayloadInvalid)

	base = testParams(t)
	items := make([]any, maxCollectionSize+1)
	base.Data = map[string]any{"items": items}
	_, err = New(base)
	assertFailureCode(t, err, codePayloadInvalid)

	base = testParams(t)
	var nested any = "leaf"
	for range maxPayloadDepth + 1 {
		nested = map[string]any{"child": nested}
	}
	base.Data = map[string]any{"root": nested}
	_, err = New(base)
	assertFailureCode(t, err, codePayloadInvalid)
}

func TestEnvelopeRejectsInvalidActorAndCorrelationIdentity(t *testing.T) {
	base := testParams(t)
	base.ActorType = ActorServiceAccount
	_, err := New(base)
	assertFailureCode(t, err, codeEnvelopeInvalid)

	base = testParams(t)
	base.CorrelationID = CorrelationID(uuid.NewString())
	_, err = New(base)
	assertFailureCode(t, err, codeEnvelopeInvalid)
}

func testEnvelope(t *testing.T) Envelope {
	t.Helper()
	envelope, err := New(testParams(t))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return envelope
}

func testParams(t *testing.T) Params {
	t.Helper()
	return Params{
		Type:           EventType("commerce.order.created.v1"),
		Producer:       Producer("urn:omnexa:module:commerce.orders"),
		OccurredAt:     time.Date(2026, time.August, 30, 12, 0, 0, 0, time.UTC),
		TenantID:       tenancy.TenantID(testUUIDv7(t)),
		CorrelationID:  CorrelationID(testUUIDv7(t)),
		Classification: DataClassInternal,
		DataSchema:     SchemaID("urn:omnexa:event-schema:commerce.order.created:v1"),
		Data:           map[string]any{"order_id": "ord_123"},
	}
}

func testUUIDv7(t *testing.T) string {
	t.Helper()
	identifier, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("uuid.NewV7() error = %v", err)
	}
	return identifier.String()
}

func assertFailureCode(t *testing.T, err error, want failure.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected failure code %q, got nil", want)
	}
	var structured *failure.Error
	if !errors.As(err, &structured) {
		t.Fatalf("error %T is not structured failure: %v", err, err)
	}
	if structured.Code() != want {
		t.Fatalf("failure code = %q, want %q", structured.Code(), want)
	}
}

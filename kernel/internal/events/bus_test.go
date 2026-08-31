package events

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

func TestPublisherValidatesAndPreservesEnvelope(t *testing.T) {
	envelope := testEnvelope(t)
	var received Envelope
	publisher, err := NewPublisher(func(_ context.Context, candidate Envelope) error {
		received = candidate
		return nil
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}

	result, publishErr := publisher.Publish(context.Background(), envelope)
	if publishErr != nil {
		t.Fatalf("Publish() error = %v", publishErr)
	}
	if !result.Accepted {
		t.Fatal("Publish() did not report provider-neutral acceptance")
	}
	if received.ID != envelope.ID || received.CorrelationID != envelope.CorrelationID || received.CausationID != envelope.CausationID || received.TenantID != envelope.TenantID || received.Classification != envelope.Classification {
		t.Fatal("publish boundary mutated canonical envelope identity or context metadata")
	}
}

func TestPublisherRejectsInvalidEnvelopeBeforeTransport(t *testing.T) {
	called := false
	publisher, err := NewPublisher(func(_ context.Context, _ Envelope) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	envelope := testEnvelope(t)
	envelope.SpecVersion = ""

	_, publishErr := publisher.Publish(context.Background(), envelope)
	assertFailureCode(t, publishErr, codeEnvelopeInvalid)
	if called {
		t.Fatal("transport was called for invalid envelope")
	}
}

func TestPublisherFailureIsSafeAndDoesNotClaimAcceptance(t *testing.T) {
	const privateCause = "restricted payload and provider topology must remain private"
	publisher, err := NewPublisher(func(_ context.Context, _ Envelope) error {
		return errors.New(privateCause)
	})
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	result, publishErr := publisher.Publish(context.Background(), testEnvelope(t))
	assertFailureCode(t, publishErr, codePublishFailed)
	if result.Accepted {
		t.Fatal("failed publish claimed acceptance")
	}
	if publishErr != nil && strings.Contains(publishErr.Error(), privateCause) {
		t.Fatalf("public publish failure leaked internal cause: %v", publishErr)
	}
}

func TestRegistryRegistersDeterministicallyAndRetainsOwner(t *testing.T) {
	registry := NewRegistry()
	registration := Registration{
		Owner:        Producer("urn:omnexa:module:commerce.orders"),
		ConsumerID:   "projection.order_summary",
		EventTypes:   []EventType{"commerce.order.updated.v1", "commerce.order.created.v1"},
		TenantScoped: true,
		Handler:      func(context.Context, Envelope) error { return nil },
	}
	if err := registry.Register(registration); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	resolved, err := registry.Resolve("commerce.order.created.v1")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if resolved.Owner != registration.Owner || resolved.ConsumerID != registration.ConsumerID || !resolved.TenantScoped {
		t.Fatalf("resolved registration lost ownership semantics: %#v", resolved)
	}
	if len(resolved.EventTypes) != 2 || resolved.EventTypes[0] != "commerce.order.created.v1" || resolved.EventTypes[1] != "commerce.order.updated.v1" {
		t.Fatalf("event types are not normalized deterministically: %#v", resolved.EventTypes)
	}
}

func TestRegistryRejectsMalformedDuplicateAndConflictingRegistrations(t *testing.T) {
	registry := NewRegistry()
	handler := func(context.Context, Envelope) error { return nil }

	assertFailureCode(t, registry.Register(Registration{
		Owner:      "",
		ConsumerID: "bad consumer id",
		EventTypes: []EventType{"commerce.order.created.v1"},
		Handler:    handler,
	}), codeRegistrationInvalid)

	assertFailureCode(t, registry.Register(Registration{
		Owner:      Producer("urn:omnexa:module:commerce.orders"),
		ConsumerID: "orders.primary",
		EventTypes: []EventType{"commerce.order.created.v1", "commerce.order.created.v1"},
		Handler:    handler,
	}), codeRegistrationInvalid)

	if err := registry.Register(Registration{
		Owner:      Producer("urn:omnexa:module:commerce.orders"),
		ConsumerID: "orders.primary",
		EventTypes: []EventType{"commerce.order.created.v1"},
		Handler:    handler,
	}); err != nil {
		t.Fatalf("first Register() error = %v", err)
	}
	assertFailureCode(t, registry.Register(Registration{
		Owner:      Producer("urn:omnexa:module:analytics.orders"),
		ConsumerID: "orders.conflict",
		EventTypes: []EventType{"commerce.order.created.v1"},
		Handler:    handler,
	}), codeRegistrationConflict)

	resolved, err := registry.Resolve("commerce.order.created.v1")
	if err != nil {
		t.Fatalf("Resolve() after conflict error = %v", err)
	}
	if resolved.Owner != Producer("urn:omnexa:module:commerce.orders") || resolved.ConsumerID != "orders.primary" {
		t.Fatal("conflicting registration silently replaced original owner")
	}
}

func TestTenantScopedInvocationFailsClosedBeforeHandler(t *testing.T) {
	registry := NewRegistry()
	called := false
	if err := registry.Register(Registration{
		Owner:        Producer("urn:omnexa:module:commerce.orders"),
		ConsumerID:   "orders.tenant_projection",
		EventTypes:   []EventType{"commerce.order.created.v1"},
		TenantScoped: true,
		Handler: func(context.Context, Envelope) error {
			called = true
			return nil
		},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	otherTenant := tenancy.TenantID(testUUIDv7(t))
	invokeErr := registry.Invoke(context.Background(), testEnvelope(t), otherTenant)
	assertFailureCode(t, invokeErr, codeTenantMismatch)
	if called {
		t.Fatal("tenant-mismatched event reached handler")
	}
}

func TestHandlerFailureIsClassifiedSeparatelyAndCauseIsPrivate(t *testing.T) {
	const privateCause = "private consumer diagnostic"
	registry := NewRegistry()
	if err := registry.Register(Registration{
		Owner:      Producer("urn:omnexa:module:commerce.orders"),
		ConsumerID: "orders.failure_projection",
		EventTypes: []EventType{"commerce.order.created.v1"},
		Handler: func(context.Context, Envelope) error {
			return errors.New(privateCause)
		},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	invokeErr := registry.Invoke(context.Background(), testEnvelope(t), "")
	assertFailureCode(t, invokeErr, codeHandlerFailed)
	if invokeErr != nil && strings.Contains(invokeErr.Error(), privateCause) {
		t.Fatalf("public handler failure leaked internal cause: %v", invokeErr)
	}
}

func TestRegistryAllowsDuplicateDeliveryAndCreatesNoOrderingGuarantee(t *testing.T) {
	registry := NewRegistry()
	calls := 0
	if err := registry.Register(Registration{
		Owner:      Producer("urn:omnexa:module:commerce.orders"),
		ConsumerID: "orders.audit_projection",
		EventTypes: []EventType{"commerce.order.created.v1"},
		Handler: func(context.Context, Envelope) error {
			calls++
			return nil
		},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	envelope := testEnvelope(t)
	if err := registry.Invoke(context.Background(), envelope, ""); err != nil {
		t.Fatalf("first Invoke() error = %v", err)
	}
	if err := registry.Invoke(context.Background(), envelope, ""); err != nil {
		t.Fatalf("duplicate Invoke() error = %v", err)
	}
	if calls != 2 {
		t.Fatalf("duplicate delivery was hidden/deduplicated: calls = %d, want 2", calls)
	}
}

func TestUnknownRouteFailsClosed(t *testing.T) {
	registry := NewRegistry()
	_, err := registry.Resolve("commerce.order.created.v1")
	assertFailureCode(t, err, codeRouteUnknown)
}

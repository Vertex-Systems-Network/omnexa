package events

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

var consumerPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$`)

const (
	codePublishInvalid       failure.Code = "events.publish.invalid"
	codePublishFailed        failure.Code = "events.publish.failed"
	codeRegistrationInvalid  failure.Code = "events.registration.invalid"
	codeRegistrationConflict failure.Code = "events.registration.conflict"
	codeRouteUnknown         failure.Code = "events.route.unknown"
)

// PublishResult reports only whether the provider-neutral publish boundary accepted
// the operation. It never represents durable delivery, handler execution or a
// downstream business mutation.
type PublishResult struct {
	Accepted bool
}

// PublishTransport is the narrow adapter seam for a later, separately governed
// transport implementation. P04.02 intentionally owns no broker or durability.
type PublishTransport func(context.Context, Envelope) error

// Publisher validates canonical P04.01 envelopes before crossing the adapter seam.
type Publisher struct {
	transport PublishTransport
}

func NewPublisher(transport PublishTransport) (*Publisher, error) {
	if transport == nil {
		return nil, classifiedFailure(codePublishInvalid, failure.CategoryValidation, "event publisher configuration is invalid")
	}
	return &Publisher{transport: transport}, nil
}

func (publisher *Publisher) Publish(ctx context.Context, envelope Envelope) (PublishResult, error) {
	if publisher == nil || publisher.transport == nil || ctx == nil {
		return PublishResult{}, classifiedFailure(codePublishInvalid, failure.CategoryValidation, "event publish request is invalid")
	}
	if err := envelope.Validate(); err != nil {
		return PublishResult{}, err
	}
	if err := publisher.transport(ctx, envelope); err != nil {
		return PublishResult{}, wrappedFailure(err, codePublishFailed, failure.CategoryUnavailable, "event publish operation failed")
	}
	return PublishResult{Accepted: true}, nil
}

// Handler is a provider-neutral consumer callback. Receipt of an Envelope never
// grants authority; protected mutations must still use already-authorized public
// capabilities owned by the target module.
type Handler func(context.Context, Envelope) error

// Registration binds explicit event contracts to exactly one stable consumer and
// owning module/component identity.
type Registration struct {
	Owner        Producer
	ConsumerID   string
	EventTypes   []EventType
	TenantScoped bool
	Handler      Handler
}

func (registration Registration) validate() error {
	if !registration.Owner.Valid() || !validConsumerID(registration.ConsumerID) || registration.Handler == nil || len(registration.EventTypes) == 0 {
		return classifiedFailure(codeRegistrationInvalid, failure.CategoryValidation, "event subscription registration is invalid")
	}
	seen := make(map[EventType]struct{}, len(registration.EventTypes))
	for _, eventType := range registration.EventTypes {
		if !eventType.Valid() {
			return classifiedFailure(codeRegistrationInvalid, failure.CategoryValidation, "event subscription registration is invalid")
		}
		if _, exists := seen[eventType]; exists {
			return classifiedFailure(codeRegistrationInvalid, failure.CategoryValidation, "event subscription registration is invalid")
		}
		seen[eventType] = struct{}{}
	}
	return nil
}

// Registry is an in-process contract registry only. It is not a worker loop,
// checkpoint store, retry scheduler, durable consumer or broker adapter.
type Registry struct {
	mu     sync.RWMutex
	routes map[EventType]Registration
}

func NewRegistry() *Registry {
	return &Registry{routes: make(map[EventType]Registration)}
}

// Register atomically rejects malformed, duplicate or conflicting routes. No
// prior owner may be silently replaced by a later registration.
func (registry *Registry) Register(registration Registration) error {
	if registry == nil {
		return classifiedFailure(codeRegistrationInvalid, failure.CategoryValidation, "event subscription registration is invalid")
	}
	if err := registration.validate(); err != nil {
		return err
	}

	normalized := registration
	normalized.ConsumerID = strings.TrimSpace(registration.ConsumerID)
	normalized.EventTypes = append([]EventType(nil), registration.EventTypes...)
	sort.Slice(normalized.EventTypes, func(i, j int) bool { return normalized.EventTypes[i] < normalized.EventTypes[j] })

	registry.mu.Lock()
	defer registry.mu.Unlock()
	for _, eventType := range normalized.EventTypes {
		if _, exists := registry.routes[eventType]; exists {
			return classifiedFailure(codeRegistrationConflict, failure.CategoryConflict, "event subscription registration conflicts with an existing route")
		}
	}
	for _, eventType := range normalized.EventTypes {
		registry.routes[eventType] = normalized
	}
	return nil
}

// Resolve returns a defensive registration snapshot for one explicit event type.
func (registry *Registry) Resolve(eventType EventType) (Registration, error) {
	if registry == nil || !eventType.Valid() {
		return Registration{}, classifiedFailure(codeRouteUnknown, failure.CategoryNotFound, "event subscription route is not registered")
	}
	registry.mu.RLock()
	registration, exists := registry.routes[eventType]
	registry.mu.RUnlock()
	if !exists {
		return Registration{}, classifiedFailure(codeRouteUnknown, failure.CategoryNotFound, "event subscription route is not registered")
	}
	registration.EventTypes = append([]EventType(nil), registration.EventTypes...)
	return registration, nil
}

// Invoke validates the canonical envelope before handler execution. Tenant-scoped
// registrations additionally require an exact trusted tenant match. Transport and
// envelope metadata never synthesize authority.
func (registry *Registry) Invoke(ctx context.Context, envelope Envelope, trustedTenant tenancy.TenantID) error {
	if ctx == nil {
		return classifiedFailure(codeRegistrationInvalid, failure.CategoryValidation, "event handler invocation context is invalid")
	}
	registration, err := registry.Resolve(envelope.Type)
	if err != nil {
		return err
	}
	if registration.TenantScoped {
		if err := envelope.ValidateForTenant(trustedTenant); err != nil {
			return err
		}
	} else if err := envelope.Validate(); err != nil {
		return err
	}
	if err := registration.Handler(ctx, envelope); err != nil {
		return wrappedFailure(err, codePublishFailed, failure.CategoryInternal, "event handler invocation failed")
	}
	return nil
}

func validConsumerID(value string) bool {
	trimmed := strings.TrimSpace(value)
	return trimmed != "" && trimmed == value && len(value) <= 128 && consumerPattern.MatchString(value)
}

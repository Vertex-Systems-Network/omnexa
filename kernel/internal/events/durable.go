package events

import (
	"context"
	"errors"
	"math"
	"regexp"
	"strings"
	"sync"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

const maxDurableScopeLength = 128

var durableScopePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$`)

const (
	codeDurableInvalid          failure.Code = "events.durable.invalid"
	codeDurableBindingConflict  failure.Code = "events.durable.binding_conflict"
	codeDurableInterrupted      failure.Code = "events.durable.interrupted"
	codeCheckpointBindFailed    failure.Code = "events.checkpoint.bind_failed"
	codeCheckpointMalformed     failure.Code = "events.checkpoint.malformed"
	codeCheckpointStale         failure.Code = "events.checkpoint.stale"
	codeCheckpointConflict      failure.Code = "events.checkpoint.conflict"
	codeCheckpointReadFailed    failure.Code = "events.checkpoint.read_failed"
	codeCheckpointWriteFailed   failure.Code = "events.checkpoint.write_failed"
	codeCheckpointExhausted     failure.Code = "events.checkpoint.exhausted"
)

// DurableScope is an explicit provider-neutral ordering/checkpoint scope. Stream
// and Partition are logical Omnexa identifiers rather than broker-specific
// offsets, topics or partition types. TenantID is required only when the bound
// P04.02 registration is tenant-scoped.
type DurableScope struct {
	Stream    string
	Partition string
	TenantID  tenancy.TenantID
}

func (scope DurableScope) validBasic() bool {
	return validDurableScopePart(scope.Stream) &&
		validDurableScopePart(scope.Partition) &&
		(scope.TenantID == "" || scope.TenantID.Valid())
}

// DurableBinding binds one durable route to the accepted P04.02 owner and stable
// consumer identity. It carries no authorization and contains no provider state.
type DurableBinding struct {
	Owner      Producer
	ConsumerID string
	EventType  EventType
	Scope      DurableScope
}

func (binding DurableBinding) validateBasic() error {
	if !binding.Owner.Valid() || !validConsumerID(binding.ConsumerID) || !binding.EventType.Valid() || !binding.Scope.validBasic() {
		return classifiedFailure(codeDurableInvalid, failure.CategoryValidation, "durable event binding is invalid")
	}
	return nil
}

func (binding DurableBinding) validateForRegistration(registration Registration) error {
	if err := binding.validateBasic(); err != nil {
		return err
	}
	if registration.Owner != binding.Owner || registration.ConsumerID != binding.ConsumerID {
		return classifiedFailure(codeDurableBindingConflict, failure.CategoryConflict, "durable event binding conflicts with the registered consumer identity")
	}
	if registration.TenantScoped {
		if !binding.Scope.TenantID.Valid() {
			return classifiedFailure(codeDurableInvalid, failure.CategoryValidation, "tenant-scoped durable event binding requires a trusted tenant scope")
		}
	} else if binding.Scope.TenantID != "" {
		return classifiedFailure(codeDurableInvalid, failure.CategoryValidation, "non-tenant durable event binding cannot invent a tenant scope")
	}
	return nil
}

// Checkpoint records consumption progress only. EventID makes a position
// auditable/conflict-detectable without copying any event payload or secret into
// checkpoint state.
type Checkpoint struct {
	Position uint64
	EventID  EventID
}

func (checkpoint Checkpoint) valid() bool {
	return checkpoint.Position > 0 && checkpoint.EventID.Valid()
}

// BindingResult separates deterministic identity/scope conflict from a storage
// adapter failure, preventing provider error text from becoming public behavior.
type BindingResult uint8

const (
	BindingResultUnknown BindingResult = iota
	BindingAccepted
	BindingUnchanged
	BindingConflict
)

// AdvanceResult separates monotonic checkpoint outcomes from provider/storage
// failures. Stale/regressive and concurrent/conflicting writes are deterministic
// contract outcomes rather than leaked backend diagnostics.
type AdvanceResult uint8

const (
	AdvanceResultUnknown AdvanceResult = iota
	CheckpointAdvanced
	CheckpointStale
	CheckpointConflict
)

// CheckpointStore is the narrow P04.03 provider-neutral persistence seam. A later
// adapter may persist this state, but no concrete database/broker or migration is
// selected by this package.
type CheckpointStore interface {
	Bind(context.Context, DurableBinding) (BindingResult, error)
	Load(context.Context, DurableBinding) (Checkpoint, bool, error)
	Advance(context.Context, DurableBinding, uint64, Checkpoint) (AdvanceResult, error)
}

// Delivery is one explicitly ordered delivery inside one DurableBinding scope.
// Position starts at 1 and must advance contiguously inside that scope; no ordering
// is inferred across different scopes.
type Delivery struct {
	Position uint64
	Envelope Envelope
}

// DurableConsumer composes the accepted P04.02 Registry with a provider-neutral
// checkpoint store. Handler success is the acknowledgement boundary; checkpoint
// advancement happens afterwards and is not a business transaction or an
// exactly-once guarantee.
type DurableConsumer struct {
	registry     *Registry
	store        CheckpointStore
	binding      DurableBinding
	tenantScoped bool
}

// NewDurableConsumer binds an already-registered P04.02 route to one explicit
// durable scope. Re-creating the consumer with the exact binding is allowed for
// restart/resume; conflicting owner/scope rebinding fails closed.
func NewDurableConsumer(ctx context.Context, registry *Registry, store CheckpointStore, binding DurableBinding) (*DurableConsumer, error) {
	if registry == nil || store == nil {
		return nil, classifiedFailure(codeDurableInvalid, failure.CategoryValidation, "durable event consumer configuration is invalid")
	}
	if err := durableContextError(ctx); err != nil {
		return nil, err
	}
	registration, err := registry.Resolve(binding.EventType)
	if err != nil {
		return nil, err
	}
	if err := binding.validateForRegistration(registration); err != nil {
		return nil, err
	}

	result, bindErr := store.Bind(ctx, binding)
	if bindErr != nil {
		return nil, wrappedFailure(bindErr, codeCheckpointBindFailed, failure.CategoryUnavailable, "durable checkpoint binding failed")
	}
	switch result {
	case BindingAccepted, BindingUnchanged:
		return &DurableConsumer{
			registry:     registry,
			store:        store,
			binding:      binding,
			tenantScoped: registration.TenantScoped,
		}, nil
	case BindingConflict:
		return nil, classifiedFailure(codeDurableBindingConflict, failure.CategoryConflict, "durable event binding conflicts with existing checkpoint scope")
	default:
		return nil, classifiedFailure(codeCheckpointMalformed, failure.CategoryInvariant, "durable checkpoint store returned an invalid binding result")
	}
}

// Binding returns the immutable logical binding used by this consumer.
func (consumer *DurableConsumer) Binding() (DurableBinding, error) {
	if consumer == nil || consumer.registry == nil || consumer.store == nil {
		return DurableBinding{}, classifiedFailure(codeDurableInvalid, failure.CategoryValidation, "durable event consumer is invalid")
	}
	return consumer.binding, nil
}

// LastCheckpoint reads and validates the last accepted progress for this exact
// binding. Missing state means no event has reached the acknowledgement point yet.
func (consumer *DurableConsumer) LastCheckpoint(ctx context.Context) (Checkpoint, bool, error) {
	if consumer == nil || consumer.registry == nil || consumer.store == nil {
		return Checkpoint{}, false, classifiedFailure(codeDurableInvalid, failure.CategoryValidation, "durable event consumer is invalid")
	}
	if err := durableContextError(ctx); err != nil {
		return Checkpoint{}, false, err
	}
	checkpoint, exists, err := consumer.store.Load(ctx, consumer.binding)
	if err != nil {
		return Checkpoint{}, false, wrappedFailure(err, codeCheckpointReadFailed, failure.CategoryUnavailable, "durable checkpoint read failed")
	}
	if !exists {
		if checkpoint != (Checkpoint{}) {
			return Checkpoint{}, false, classifiedFailure(codeCheckpointMalformed, failure.CategoryInvariant, "durable checkpoint store returned malformed missing state")
		}
		return Checkpoint{}, false, nil
	}
	if !checkpoint.valid() {
		return Checkpoint{}, false, classifiedFailure(codeCheckpointMalformed, failure.CategoryInvariant, "durable checkpoint state is malformed")
	}
	return checkpoint, true, nil
}

// ResumePosition returns the next contiguous position for this binding. Position 1
// is the documented initial boundary convention when no checkpoint exists.
func (consumer *DurableConsumer) ResumePosition(ctx context.Context) (uint64, error) {
	checkpoint, exists, err := consumer.LastCheckpoint(ctx)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 1, nil
	}
	if checkpoint.Position == math.MaxUint64 {
		return 0, classifiedFailure(codeCheckpointExhausted, failure.CategoryInvariant, "durable checkpoint position is exhausted")
	}
	return checkpoint.Position + 1, nil
}

// Process validates one delivery, invokes the accepted P04.02 handler, then
// advances the checkpoint exactly one position. Handler failure/cancellation never
// advances progress. If checkpoint persistence fails after handler success, the
// same delivery may be observed again on restart by design.
func (consumer *DurableConsumer) Process(ctx context.Context, delivery Delivery) error {
	if consumer == nil || consumer.registry == nil || consumer.store == nil || delivery.Position == 0 {
		return classifiedFailure(codeDurableInvalid, failure.CategoryValidation, "durable event delivery is invalid")
	}
	if err := durableContextError(ctx); err != nil {
		return err
	}
	if delivery.Envelope.Type != consumer.binding.EventType {
		return classifiedFailure(codeDurableBindingConflict, failure.CategoryConflict, "durable event delivery conflicts with the bound event route")
	}
	if consumer.tenantScoped {
		if err := delivery.Envelope.ValidateForTenant(consumer.binding.Scope.TenantID); err != nil {
			return err
		}
	} else if err := delivery.Envelope.Validate(); err != nil {
		return err
	}

	checkpoint, exists, err := consumer.LastCheckpoint(ctx)
	if err != nil {
		return err
	}
	expectedPosition := uint64(0)
	if exists {
		expectedPosition = checkpoint.Position
		if delivery.Position <= checkpoint.Position {
			return classifiedFailure(codeCheckpointStale, failure.CategoryConflict, "durable checkpoint advancement is stale or regressive")
		}
		if checkpoint.Position == math.MaxUint64 || delivery.Position != checkpoint.Position+1 {
			return classifiedFailure(codeCheckpointConflict, failure.CategoryConflict, "durable checkpoint advancement conflicts with contiguous progress")
		}
	} else if delivery.Position != 1 {
		return classifiedFailure(codeCheckpointConflict, failure.CategoryConflict, "durable checkpoint advancement cannot skip unacknowledged work")
	}

	trustedTenant := tenancy.TenantID("")
	if consumer.tenantScoped {
		trustedTenant = consumer.binding.Scope.TenantID
	}
	if invokeErr := consumer.registry.Invoke(ctx, delivery.Envelope, trustedTenant); invokeErr != nil {
		return invokeErr
	}
	if err := durableContextError(ctx); err != nil {
		return err
	}

	next := Checkpoint{Position: delivery.Position, EventID: delivery.Envelope.ID}
	result, advanceErr := consumer.store.Advance(ctx, consumer.binding, expectedPosition, next)
	if advanceErr != nil {
		return wrappedFailure(advanceErr, codeCheckpointWriteFailed, failure.CategoryUnavailable, "durable checkpoint write failed")
	}
	switch result {
	case CheckpointAdvanced:
		return nil
	case CheckpointStale:
		return classifiedFailure(codeCheckpointStale, failure.CategoryConflict, "durable checkpoint advancement is stale or regressive")
	case CheckpointConflict:
		return classifiedFailure(codeCheckpointConflict, failure.CategoryConflict, "durable checkpoint advancement conflicts with accepted progress")
	default:
		return classifiedFailure(codeCheckpointMalformed, failure.CategoryInvariant, "durable checkpoint store returned an invalid advancement result")
	}
}

// MemoryCheckpointStore is a deterministic reference adapter for focused P04.03
// tests and single-process restart simulation. It is not process-durable storage
// and must not be represented as a production persistence guarantee.
type MemoryCheckpointStore struct {
	mu             sync.RWMutex
	consumerOwners map[string]Producer
	scopes         map[string]DurableScope
	checkpoints    map[string]Checkpoint
}

func NewMemoryCheckpointStore() *MemoryCheckpointStore {
	return &MemoryCheckpointStore{
		consumerOwners: make(map[string]Producer),
		scopes:         make(map[string]DurableScope),
		checkpoints:    make(map[string]Checkpoint),
	}
}

func (store *MemoryCheckpointStore) Bind(ctx context.Context, binding DurableBinding) (BindingResult, error) {
	if err := memoryStoreRequestError(ctx, store, binding); err != nil {
		return BindingResultUnknown, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.consumerOwners == nil || store.scopes == nil || store.checkpoints == nil {
		return BindingResultUnknown, errors.New("durable checkpoint store is not initialized")
	}
	if owner, exists := store.consumerOwners[binding.ConsumerID]; exists && owner != binding.Owner {
		return BindingConflict, nil
	}
	logicalKey := bindingLogicalKey(binding)
	if scope, exists := store.scopes[logicalKey]; exists {
		if scope != binding.Scope {
			return BindingConflict, nil
		}
		return BindingUnchanged, nil
	}
	store.consumerOwners[binding.ConsumerID] = binding.Owner
	store.scopes[logicalKey] = binding.Scope
	return BindingAccepted, nil
}

func (store *MemoryCheckpointStore) Load(ctx context.Context, binding DurableBinding) (Checkpoint, bool, error) {
	if err := memoryStoreRequestError(ctx, store, binding); err != nil {
		return Checkpoint{}, false, err
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	if store.consumerOwners == nil || store.scopes == nil || store.checkpoints == nil {
		return Checkpoint{}, false, errors.New("durable checkpoint store is not initialized")
	}
	if !store.bindingMatchesLocked(binding) {
		return Checkpoint{}, false, errors.New("durable checkpoint binding is not registered")
	}
	checkpoint, exists := store.checkpoints[bindingCheckpointKey(binding)]
	return checkpoint, exists, nil
}

func (store *MemoryCheckpointStore) Advance(ctx context.Context, binding DurableBinding, expectedPosition uint64, next Checkpoint) (AdvanceResult, error) {
	if err := memoryStoreRequestError(ctx, store, binding); err != nil {
		return AdvanceResultUnknown, err
	}
	if !next.valid() {
		return AdvanceResultUnknown, errors.New("durable checkpoint candidate is invalid")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.consumerOwners == nil || store.scopes == nil || store.checkpoints == nil {
		return AdvanceResultUnknown, errors.New("durable checkpoint store is not initialized")
	}
	if !store.bindingMatchesLocked(binding) {
		return AdvanceResultUnknown, errors.New("durable checkpoint binding is not registered")
	}

	key := bindingCheckpointKey(binding)
	current, exists := store.checkpoints[key]
	if !exists {
		if expectedPosition != 0 || next.Position != 1 {
			return CheckpointConflict, nil
		}
		store.checkpoints[key] = next
		return CheckpointAdvanced, nil
	}
	if !current.valid() {
		return AdvanceResultUnknown, errors.New("durable checkpoint store contains malformed state")
	}
	if next.Position <= current.Position {
		return CheckpointStale, nil
	}
	if expectedPosition != current.Position || current.Position == math.MaxUint64 || next.Position != current.Position+1 {
		return CheckpointConflict, nil
	}
	store.checkpoints[key] = next
	return CheckpointAdvanced, nil
}

func (store *MemoryCheckpointStore) bindingMatchesLocked(binding DurableBinding) bool {
	owner, ownerExists := store.consumerOwners[binding.ConsumerID]
	if !ownerExists || owner != binding.Owner {
		return false
	}
	scope, scopeExists := store.scopes[bindingLogicalKey(binding)]
	return scopeExists && scope == binding.Scope
}

func memoryStoreRequestError(ctx context.Context, store *MemoryCheckpointStore, binding DurableBinding) error {
	if store == nil {
		return errors.New("durable checkpoint store is nil")
	}
	if ctx == nil {
		return errors.New("durable checkpoint store context is invalid")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := binding.validateBasic(); err != nil {
		return errors.New("durable checkpoint binding is invalid")
	}
	return nil
}

func bindingLogicalKey(binding DurableBinding) string {
	return strings.Join([]string{
		binding.ConsumerID,
		string(binding.EventType),
		string(binding.Scope.TenantID),
	}, "\x1f")
}

func bindingCheckpointKey(binding DurableBinding) string {
	return strings.Join([]string{
		string(binding.Owner),
		binding.ConsumerID,
		string(binding.EventType),
		binding.Scope.Stream,
		binding.Scope.Partition,
		string(binding.Scope.TenantID),
	}, "\x1f")
}

func validDurableScopePart(value string) bool {
	trimmed := strings.TrimSpace(value)
	return trimmed != "" && trimmed == value && len(value) <= maxDurableScopeLength && durableScopePattern.MatchString(value)
}

func durableContextError(ctx context.Context) error {
	if ctx == nil {
		return classifiedFailure(codeDurableInvalid, failure.CategoryValidation, "durable event context is invalid")
	}
	if err := ctx.Err(); err != nil {
		return wrappedFailure(err, codeDurableInterrupted, failure.CategoryUnavailable, "durable event consumption was interrupted")
	}
	return nil
}

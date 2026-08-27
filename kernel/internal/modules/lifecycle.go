package modules

import (
	"context"
	"errors"
	"sort"
	"sync"
)

// LifecycleState is durable lifecycle metadata only. It does not grant
// capabilities, permissions, tenant scope, private-package access, or database
// authority.
type LifecycleState string

const (
	LifecycleAvailable        LifecycleState = "available"
	LifecycleInstalled        LifecycleState = "installed"
	LifecycleEnabled          LifecycleState = "enabled"
	LifecycleDisabled         LifecycleState = "disabled"
	LifecycleSuspended        LifecycleState = "suspended"
	LifecycleArchived         LifecycleState = "archived"
	LifecycleDetached         LifecycleState = "detached"
	LifecycleRecoveryRequired LifecycleState = "recovery_required"
	LifecyclePurged           LifecycleState = "purged"
)

// LifecycleAction is an explicit state-machine mutation request.
type LifecycleAction string

const (
	LifecycleInstall LifecycleAction = "install"
	LifecycleEnable  LifecycleAction = "enable"
	LifecycleDisable LifecycleAction = "disable"
	LifecycleSuspend LifecycleAction = "suspend"
	LifecycleResume  LifecycleAction = "resume"
	LifecycleArchive LifecycleAction = "archive"
	LifecycleRestore LifecycleAction = "restore"
	LifecycleDetach  LifecycleAction = "detach"
	LifecyclePurge   LifecycleAction = "purge"
	LifecycleUpgrade LifecycleAction = "upgrade"
	LifecycleRecover LifecycleAction = "recover"
)

// LifecycleRecord is caller-owned lifecycle state. Registry discovery remains a
// separate immutable metadata boundary and never silently becomes lifecycle
// truth. Revision is used for optimistic concurrency.
type LifecycleRecord struct {
	ModuleID        string          `json:"module_id"`
	Version         string          `json:"version,omitempty"`
	State           LifecycleState  `json:"state"`
	Revision        uint64          `json:"revision"`
	PreviousState   LifecycleState  `json:"previous_state,omitempty"`
	RecoveryState   LifecycleState  `json:"recovery_state,omitempty"`
	FailedAction    LifecycleAction `json:"failed_action,omitempty"`
	FailureCode     string          `json:"failure_code,omitempty"`
	LastOperationID string          `json:"last_operation_id,omitempty"`
	LastAction      LifecycleAction `json:"last_action,omitempty"`
}

// LifecycleRequest requires a stable operation ID. Callers retrying the same
// logical mutation must reuse the same ID; reuse with a different action fails.
type LifecycleRequest struct {
	ModuleID    string
	Action      LifecycleAction
	OperationID string
}

// LifecycleFailureRequest records partial external-hook/side-effect failure in
// a recoverable, explicit state. The state machine itself does not execute
// module code, migrations, filesystem actions, network calls, or package hooks.
type LifecycleFailureRequest struct {
	ModuleID     string
	FailedAction LifecycleAction
	OperationID  string
	FailureCode  string
}

// LifecycleResult reports the committed record and whether the request was an
// idempotent replay of an already committed operation.
type LifecycleResult struct {
	Record   LifecycleRecord
	Replayed bool
}

// LifecycleDiagnostic is stable and safe to expose. It intentionally omits raw
// manifests, authorization errors, audit payloads, filesystem paths, and data.
type LifecycleDiagnostic struct {
	Code         string `json:"code"`
	ModuleID     string `json:"module_id,omitempty"`
	DependencyID string `json:"dependency_id,omitempty"`
}

// LifecycleError is a fail-closed lifecycle error.
type LifecycleError struct {
	Diagnostic LifecycleDiagnostic
}

func (e *LifecycleError) Error() string { return "module lifecycle transition failed" }

// LifecycleStore is the persistence boundary for P03.04. Durable deployments
// may provide a database-backed CAS implementation without changing state-
// machine semantics. P03.04 does not introduce a schema/migration here.
type LifecycleStore interface {
	Load(ctx context.Context, moduleID string) (LifecycleRecord, bool, error)
	CompareAndSwap(ctx context.Context, moduleID string, expectedRevision uint64, next LifecycleRecord) error
}

// ErrLifecycleConflict signals an optimistic-concurrency conflict.
var ErrLifecycleConflict = errors.New("module lifecycle revision conflict")

// MemoryLifecycleStore is a deterministic CAS store useful for kernel wiring,
// tests, and non-durable contexts. Durable runtime state must use an adapter
// appropriate to the owning deployment boundary.
type MemoryLifecycleStore struct {
	mu      sync.Mutex
	records map[string]LifecycleRecord
}

func NewMemoryLifecycleStore() *MemoryLifecycleStore {
	return &MemoryLifecycleStore{records: make(map[string]LifecycleRecord)}
}

func (s *MemoryLifecycleStore) Load(_ context.Context, moduleID string) (LifecycleRecord, bool, error) {
	if s == nil {
		return LifecycleRecord{}, false, errors.New("lifecycle store unavailable")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.records[moduleID]
	return record, ok, nil
}

func (s *MemoryLifecycleStore) CompareAndSwap(_ context.Context, moduleID string, expectedRevision uint64, next LifecycleRecord) error {
	if s == nil {
		return errors.New("lifecycle store unavailable")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.records[moduleID]
	if !ok {
		if expectedRevision != 0 {
			return ErrLifecycleConflict
		}
	} else if current.Revision != expectedRevision {
		return ErrLifecycleConflict
	}
	s.records[moduleID] = next
	return nil
}

// LifecycleAuthorizer preserves the P02 authorization boundary. The lifecycle
// package never interprets identity/role claims itself.
type LifecycleAuthorizer interface {
	AuthorizeLifecycle(ctx context.Context, moduleID string, action LifecycleAction) error
}

// LifecycleAuditEvent is the pre-commit protected-mutation audit envelope. A
// required audit failure prevents lifecycle state from changing.
type LifecycleAuditEvent struct {
	ModuleID    string
	Action      LifecycleAction
	OperationID string
	FromState   LifecycleState
	ToState     LifecycleState
	FromVersion string
	ToVersion   string
	FailureCode string
}

// LifecycleAuditor preserves the P01.11 protected audit boundary. The state
// machine requires successful audit delivery before committing a mutation.
type LifecycleAuditor interface {
	RecordLifecycle(ctx context.Context, event LifecycleAuditEvent) error
}

// LifecycleUpgradeCoordinator is the narrow forward-compatible seam for the
// later P03.09 migration-ownership contract. P03.04 coordinates but does not
// execute or invent migrations. Upgrade fails closed when no coordinator exists.
type LifecycleUpgradeCoordinator interface {
	ValidateUpgrade(ctx context.Context, moduleID, fromVersion, toVersion string) error
}

// LifecycleManager applies deterministic lifecycle transitions against exact
// registry metadata and an optimistic-concurrency store.
type LifecycleManager struct {
	Registry           Registry
	Platform           PlatformSnapshot
	Store              LifecycleStore
	Authorizer         LifecycleAuthorizer
	Auditor            LifecycleAuditor
	UpgradeCoordinator LifecycleUpgradeCoordinator
}

// Apply authorizes, audits, and atomically commits one lifecycle mutation.
func (m LifecycleManager) Apply(ctx context.Context, request LifecycleRequest) (LifecycleResult, error) {
	if !validLifecycleOperationID(request.OperationID) {
		return LifecycleResult{}, lifecycleErr("lifecycle.operation.invalid", request.ModuleID, "")
	}
	if !validLifecycleAction(request.Action) {
		return LifecycleResult{}, lifecycleErr("lifecycle.action.invalid", request.ModuleID, "")
	}
	meta, ok := m.Registry.Lookup(request.ModuleID)
	if !ok {
		return LifecycleResult{}, lifecycleErr("lifecycle.module.not_discovered", request.ModuleID, "")
	}
	current, found, loadErr := m.loadRecord(ctx, request.ModuleID)
	if loadErr != nil {
		return LifecycleResult{}, loadErr
	}
	if !found {
		current = LifecycleRecord{ModuleID: request.ModuleID, State: LifecycleAvailable}
	}
	if current.ModuleID != request.ModuleID {
		return LifecycleResult{}, lifecycleErr("lifecycle.store.identity_mismatch", request.ModuleID, "")
	}

	// Retry/idempotency never bypasses current authorization. A replay is not a
	// new mutation, so it does not emit another mutation audit event.
	if current.LastOperationID == request.OperationID {
		if err := m.authorize(ctx, request.ModuleID, request.Action); err != nil {
			return LifecycleResult{}, err
		}
		if current.LastAction != request.Action {
			return LifecycleResult{}, lifecycleErr("lifecycle.operation.reused", request.ModuleID, "")
		}
		return LifecycleResult{Record: current, Replayed: true}, nil
	}

	next, err := m.plan(ctx, current, meta, request)
	if err != nil {
		return LifecycleResult{}, err
	}
	if err := m.authorize(ctx, request.ModuleID, request.Action); err != nil {
		return LifecycleResult{}, err
	}
	if err := m.audit(ctx, LifecycleAuditEvent{
		ModuleID: request.ModuleID, Action: request.Action, OperationID: request.OperationID,
		FromState: current.State, ToState: next.State, FromVersion: current.Version, ToVersion: next.Version,
	}); err != nil {
		return LifecycleResult{}, err
	}
	next.Revision = current.Revision + 1
	next.LastOperationID = request.OperationID
	next.LastAction = request.Action
	if err := m.storeCAS(ctx, request.ModuleID, current.Revision, next); err != nil {
		return LifecycleResult{}, err
	}
	return LifecycleResult{Record: next}, nil
}

// MarkRecoveryRequired records a recoverable partial-failure state after an
// external lifecycle hook/side effect fails. A failed initial install is
// representable from the implicit available state without first committing an
// installed record.
func (m LifecycleManager) MarkRecoveryRequired(ctx context.Context, request LifecycleFailureRequest) (LifecycleRecord, error) {
	if !validLifecycleOperationID(request.OperationID) || !validLifecycleFailureCode(request.FailureCode) || !validLifecycleAction(request.FailedAction) {
		return LifecycleRecord{}, lifecycleErr("lifecycle.failure.invalid", request.ModuleID, "")
	}
	if _, ok := m.Registry.Lookup(request.ModuleID); !ok {
		return LifecycleRecord{}, lifecycleErr("lifecycle.module.not_discovered", request.ModuleID, "")
	}
	current, found, err := m.loadRecord(ctx, request.ModuleID)
	if err != nil {
		return LifecycleRecord{}, err
	}
	if !found {
		if request.FailedAction != LifecycleInstall {
			return LifecycleRecord{}, lifecycleErr("lifecycle.failure.state_invalid", request.ModuleID, "")
		}
		current = LifecycleRecord{ModuleID: request.ModuleID, State: LifecycleAvailable}
	} else if current.ModuleID != request.ModuleID {
		return LifecycleRecord{}, lifecycleErr("lifecycle.store.identity_mismatch", request.ModuleID, "")
	}

	if current.LastOperationID == request.OperationID {
		if err := m.authorize(ctx, request.ModuleID, request.FailedAction); err != nil {
			return LifecycleRecord{}, err
		}
		if current.LastAction != request.FailedAction || current.State != LifecycleRecoveryRequired || current.FailureCode != request.FailureCode {
			return LifecycleRecord{}, lifecycleErr("lifecycle.operation.reused", request.ModuleID, "")
		}
		return current, nil
	}
	if current.State == LifecyclePurged || current.State == LifecycleRecoveryRequired {
		return LifecycleRecord{}, lifecycleErr("lifecycle.failure.state_invalid", request.ModuleID, "")
	}
	if current.State == LifecycleAvailable && request.FailedAction != LifecycleInstall {
		return LifecycleRecord{}, lifecycleErr("lifecycle.failure.state_invalid", request.ModuleID, "")
	}
	if err := m.authorize(ctx, request.ModuleID, request.FailedAction); err != nil {
		return LifecycleRecord{}, err
	}

	next := current
	next.State = LifecycleRecoveryRequired
	next.RecoveryState = current.State
	next.FailedAction = request.FailedAction
	next.FailureCode = request.FailureCode
	next.Revision = current.Revision + 1
	next.LastOperationID = request.OperationID
	next.LastAction = request.FailedAction
	if err := m.audit(ctx, LifecycleAuditEvent{
		ModuleID: request.ModuleID, Action: request.FailedAction, OperationID: request.OperationID,
		FromState: current.State, ToState: next.State, FromVersion: current.Version, ToVersion: next.Version,
		FailureCode: request.FailureCode,
	}); err != nil {
		return LifecycleRecord{}, err
	}
	if err := m.storeCAS(ctx, request.ModuleID, current.Revision, next); err != nil {
		return LifecycleRecord{}, err
	}
	return next, nil
}

func (m LifecycleManager) plan(ctx context.Context, current LifecycleRecord, meta RegistryRecord, request LifecycleRequest) (LifecycleRecord, error) {
	next := current
	next.ModuleID = request.ModuleID
	next.FailureCode = ""
	next.FailedAction = ""
	next.RecoveryState = ""

	switch request.Action {
	case LifecycleInstall:
		if current.State != LifecycleAvailable && current.State != LifecycleDetached && current.State != LifecyclePurged {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		if err := m.requireGraphEligibility(request.ModuleID); err != nil {
			return LifecycleRecord{}, err
		}
		if err := m.requireDependencies(ctx, request.ModuleID, false); err != nil {
			return LifecycleRecord{}, err
		}
		next.State = LifecycleInstalled
		next.Version = meta.Version
		next.PreviousState = ""

	case LifecycleEnable:
		if current.State != LifecycleInstalled && current.State != LifecycleDisabled {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		if current.Version != meta.Version {
			return LifecycleRecord{}, lifecycleErr("lifecycle.version.mismatch", request.ModuleID, "")
		}
		if err := m.requireGraphEligibility(request.ModuleID); err != nil {
			return LifecycleRecord{}, err
		}
		if err := m.requireDependencies(ctx, request.ModuleID, true); err != nil {
			return LifecycleRecord{}, err
		}
		next.State = LifecycleEnabled
		next.PreviousState = ""

	case LifecycleDisable:
		if current.State != LifecycleEnabled {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		blocker, err := m.reverseBlocker(ctx, request.ModuleID, true)
		if err != nil {
			return LifecycleRecord{}, err
		}
		if blocker != "" {
			return LifecycleRecord{}, lifecycleErr("lifecycle.reverse_dependency.active", request.ModuleID, blocker)
		}
		next.State = LifecycleDisabled

	case LifecycleSuspend:
		if current.State != LifecycleEnabled {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		blocker, err := m.reverseBlocker(ctx, request.ModuleID, true)
		if err != nil {
			return LifecycleRecord{}, err
		}
		if blocker != "" {
			return LifecycleRecord{}, lifecycleErr("lifecycle.reverse_dependency.active", request.ModuleID, blocker)
		}
		next.PreviousState = current.State
		next.State = LifecycleSuspended

	case LifecycleResume:
		if current.State != LifecycleSuspended || current.PreviousState != LifecycleEnabled {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		if current.Version != meta.Version {
			return LifecycleRecord{}, lifecycleErr("lifecycle.version.mismatch", request.ModuleID, "")
		}
		if err := m.requireGraphEligibility(request.ModuleID); err != nil {
			return LifecycleRecord{}, err
		}
		if err := m.requireDependencies(ctx, request.ModuleID, true); err != nil {
			return LifecycleRecord{}, err
		}
		next.State = current.PreviousState
		next.PreviousState = ""

	case LifecycleArchive:
		if current.State != LifecycleInstalled && current.State != LifecycleDisabled {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		blocker, err := m.reverseBlocker(ctx, request.ModuleID, false)
		if err != nil {
			return LifecycleRecord{}, err
		}
		if blocker != "" {
			return LifecycleRecord{}, lifecycleErr("lifecycle.reverse_dependency.present", request.ModuleID, blocker)
		}
		next.PreviousState = current.State
		next.State = LifecycleArchived

	case LifecycleRestore:
		if current.State != LifecycleArchived || (current.PreviousState != LifecycleInstalled && current.PreviousState != LifecycleDisabled) {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		next.State = current.PreviousState
		next.PreviousState = ""

	case LifecycleDetach:
		if current.State != LifecycleInstalled && current.State != LifecycleDisabled && current.State != LifecycleArchived {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		blocker, err := m.reverseBlocker(ctx, request.ModuleID, false)
		if err != nil {
			return LifecycleRecord{}, err
		}
		if blocker != "" {
			return LifecycleRecord{}, lifecycleErr("lifecycle.reverse_dependency.present", request.ModuleID, blocker)
		}
		next.PreviousState = current.State
		next.State = LifecycleDetached

	case LifecyclePurge:
		if current.State != LifecycleDetached {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		blocker, err := m.reverseBlocker(ctx, request.ModuleID, false)
		if err != nil {
			return LifecycleRecord{}, err
		}
		if blocker != "" {
			return LifecycleRecord{}, lifecycleErr("lifecycle.reverse_dependency.present", request.ModuleID, blocker)
		}
		next.State = LifecyclePurged

	case LifecycleUpgrade:
		if !upgradeableState(current.State) || current.Version == "" {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		fromVersion, ok := parseStrictSemVer(current.Version)
		if !ok {
			return LifecycleRecord{}, lifecycleErr("lifecycle.upgrade.current_version_invalid", request.ModuleID, "")
		}
		toVersion, ok := parseStrictSemVer(meta.Version)
		if !ok || compareStrictSemVer(toVersion, fromVersion) <= 0 {
			return LifecycleRecord{}, lifecycleErr("lifecycle.upgrade.target_invalid", request.ModuleID, "")
		}
		if err := m.requireGraphEligibility(request.ModuleID); err != nil {
			return LifecycleRecord{}, err
		}
		if err := m.requireDependencies(ctx, request.ModuleID, current.State == LifecycleEnabled); err != nil {
			return LifecycleRecord{}, err
		}
		if m.UpgradeCoordinator == nil {
			return LifecycleRecord{}, lifecycleErr("lifecycle.upgrade.coordinator_unavailable", request.ModuleID, "")
		}
		if err := m.UpgradeCoordinator.ValidateUpgrade(ctx, request.ModuleID, current.Version, meta.Version); err != nil {
			return LifecycleRecord{}, lifecycleErr("lifecycle.upgrade.not_ready", request.ModuleID, "")
		}
		next.Version = meta.Version

	case LifecycleRecover:
		if current.State != LifecycleRecoveryRequired || !stableRecoveryState(current.RecoveryState) {
			return LifecycleRecord{}, invalidTransition(request.ModuleID)
		}
		next.State = current.RecoveryState
		next.RecoveryState = ""
		next.FailedAction = ""
		next.FailureCode = ""

	default:
		return LifecycleRecord{}, lifecycleErr("lifecycle.action.invalid", request.ModuleID, "")
	}
	return next, nil
}

func (m LifecycleManager) requireGraphEligibility(moduleID string) error {
	resolution, err := ResolveDependencies(m.Registry, m.Platform, nil)
	if err != nil {
		return lifecycleErr("lifecycle.dependency.graph_invalid", moduleID, "")
	}
	for _, id := range resolution.Order {
		if id == moduleID {
			return nil
		}
	}
	return lifecycleErr("lifecycle.dependency.order_missing", moduleID, "")
}

func (m LifecycleManager) requireDependencies(ctx context.Context, moduleID string, requireEnabled bool) error {
	snapshot, ok := m.Registry.manifestSnapshot(moduleID)
	if !ok {
		return lifecycleErr("lifecycle.registry.snapshot_missing", moduleID, "")
	}
	for _, dependency := range snapshot.RequiredDependencies {
		meta, ok := m.Registry.Lookup(dependency.ID)
		if !ok {
			return lifecycleErr("lifecycle.dependency.not_discovered", moduleID, dependency.ID)
		}
		record, found, err := m.loadRecord(ctx, dependency.ID)
		if err != nil {
			return err
		}
		state := LifecycleAvailable
		if found {
			if record.ModuleID != dependency.ID {
				return lifecycleErr("lifecycle.store.identity_mismatch", moduleID, dependency.ID)
			}
			if record.Version != meta.Version {
				return lifecycleErr("lifecycle.dependency.version_mismatch", moduleID, dependency.ID)
			}
			state = record.State
		}
		if requireEnabled {
			if state != LifecycleEnabled {
				return lifecycleErr("lifecycle.dependency.not_enabled", moduleID, dependency.ID)
			}
			continue
		}
		if !installedDependencyState(state) {
			return lifecycleErr("lifecycle.dependency.not_installed", moduleID, dependency.ID)
		}
	}
	return nil
}

func (m LifecycleManager) reverseBlocker(ctx context.Context, providerID string, activeOnly bool) (string, error) {
	blockers := make([]string, 0)
	for _, candidate := range m.Registry.List() {
		if candidate.ID == providerID {
			continue
		}
		snapshot, ok := m.Registry.manifestSnapshot(candidate.ID)
		if !ok {
			return "", lifecycleErr("lifecycle.registry.snapshot_missing", providerID, candidate.ID)
		}
		if !requiredDependency(snapshot, providerID) {
			continue
		}
		record, found, err := m.loadRecord(ctx, candidate.ID)
		if err != nil {
			return "", err
		}
		if !found {
			continue
		}
		if record.ModuleID != candidate.ID {
			return "", lifecycleErr("lifecycle.store.identity_mismatch", providerID, candidate.ID)
		}
		if activeOnly {
			if record.State == LifecycleEnabled {
				blockers = append(blockers, candidate.ID)
			}
			continue
		}
		if installedDependencyState(record.State) || record.State == LifecycleRecoveryRequired {
			blockers = append(blockers, candidate.ID)
		}
	}
	sort.Strings(blockers)
	if len(blockers) == 0 {
		return "", nil
	}
	return blockers[0], nil
}

func requiredDependency(snapshot validatedManifestSnapshot, id string) bool {
	for _, dependency := range snapshot.RequiredDependencies {
		if dependency.ID == id {
			return true
		}
	}
	return false
}

func (m LifecycleManager) loadRecord(ctx context.Context, moduleID string) (LifecycleRecord, bool, error) {
	if m.Store == nil {
		return LifecycleRecord{}, false, lifecycleErr("lifecycle.store.unavailable", moduleID, "")
	}
	record, ok, err := m.Store.Load(ctx, moduleID)
	if err != nil {
		return LifecycleRecord{}, false, lifecycleErr("lifecycle.store.read_failed", moduleID, "")
	}
	return record, ok, nil
}

func (m LifecycleManager) storeCAS(ctx context.Context, moduleID string, revision uint64, next LifecycleRecord) error {
	if m.Store == nil {
		return lifecycleErr("lifecycle.store.unavailable", moduleID, "")
	}
	if err := m.Store.CompareAndSwap(ctx, moduleID, revision, next); err != nil {
		if errors.Is(err, ErrLifecycleConflict) {
			return lifecycleErr("lifecycle.concurrent_conflict", moduleID, "")
		}
		return lifecycleErr("lifecycle.store.write_failed", moduleID, "")
	}
	return nil
}

func (m LifecycleManager) authorize(ctx context.Context, moduleID string, action LifecycleAction) error {
	if m.Authorizer == nil {
		return lifecycleErr("lifecycle.authorization.unavailable", moduleID, "")
	}
	if err := m.Authorizer.AuthorizeLifecycle(ctx, moduleID, action); err != nil {
		return lifecycleErr("lifecycle.authorization.denied", moduleID, "")
	}
	return nil
}

func (m LifecycleManager) audit(ctx context.Context, event LifecycleAuditEvent) error {
	if m.Auditor == nil {
		return lifecycleErr("lifecycle.audit.unavailable", event.ModuleID, "")
	}
	if err := m.Auditor.RecordLifecycle(ctx, event); err != nil {
		return lifecycleErr("lifecycle.audit.failed", event.ModuleID, "")
	}
	return nil
}

func lifecycleErr(code, moduleID, dependencyID string) error {
	return &LifecycleError{Diagnostic: LifecycleDiagnostic{Code: code, ModuleID: moduleID, DependencyID: dependencyID}}
}

func invalidTransition(moduleID string) error {
	return lifecycleErr("lifecycle.transition.invalid", moduleID, "")
}

func validLifecycleAction(action LifecycleAction) bool {
	switch action {
	case LifecycleInstall, LifecycleEnable, LifecycleDisable, LifecycleSuspend, LifecycleResume,
		LifecycleArchive, LifecycleRestore, LifecycleDetach, LifecyclePurge, LifecycleUpgrade, LifecycleRecover:
		return true
	default:
		return false
	}
}

func validLifecycleOperationID(value string) bool {
	if value == "" || len(value) > 128 {
		return false
	}
	for i := 0; i < len(value); i++ {
		c := value[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '.' || c == '-' || c == '_' || c == ':' {
			continue
		}
		return false
	}
	return true
}

func validLifecycleFailureCode(value string) bool {
	return validLifecycleOperationID(value)
}

func installedDependencyState(state LifecycleState) bool {
	switch state {
	case LifecycleInstalled, LifecycleEnabled, LifecycleDisabled, LifecycleSuspended, LifecycleArchived:
		return true
	default:
		return false
	}
}

func upgradeableState(state LifecycleState) bool {
	return installedDependencyState(state)
}

func stableRecoveryState(state LifecycleState) bool {
	return state == LifecycleAvailable || installedDependencyState(state) || state == LifecycleDetached
}

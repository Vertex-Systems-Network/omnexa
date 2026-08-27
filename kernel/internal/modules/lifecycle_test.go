package modules

import (
	"context"
	"errors"
	"sync"
	"testing"
)

type lifecycleAllowAuthorizer struct {
	deny LifecycleAction
}

func (a lifecycleAllowAuthorizer) AuthorizeLifecycle(_ context.Context, _ string, action LifecycleAction) error {
	if action == a.deny {
		return errors.New("denied")
	}
	return nil
}

type lifecycleAuditRecorder struct {
	mu     sync.Mutex
	events []LifecycleAuditEvent
	fail   bool
}

func (a *lifecycleAuditRecorder) RecordLifecycle(_ context.Context, event LifecycleAuditEvent) error {
	if a.fail {
		return errors.New("audit unavailable")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, event)
	return nil
}

type lifecycleUpgradeAllow struct {
	called bool
}

func (u *lifecycleUpgradeAllow) ValidateUpgrade(_ context.Context, _, _, _ string) error {
	u.called = true
	return nil
}

type lifecycleConflictStore struct {
	base LifecycleStore
}

func (s lifecycleConflictStore) Load(ctx context.Context, moduleID string) (LifecycleRecord, bool, error) {
	return s.base.Load(ctx, moduleID)
}

func (lifecycleConflictStore) CompareAndSwap(context.Context, string, uint64, LifecycleRecord) error {
	return ErrLifecycleConflict
}

func lifecycleManagerFor(registry Registry, store LifecycleStore, auditor LifecycleAuditor) LifecycleManager {
	return LifecycleManager{
		Registry:   registry,
		Platform:   PlatformSnapshot{Version: "1.0.0"},
		Store:      store,
		Authorizer: lifecycleAllowAuthorizer{},
		Auditor:    auditor,
	}
}

func lifecycleCode(t *testing.T, err error) string {
	t.Helper()
	var lifecycleErrValue *LifecycleError
	if !errors.As(err, &lifecycleErrValue) {
		t.Fatalf("expected LifecycleError, got %T: %v", err, err)
	}
	return lifecycleErrValue.Diagnostic.Code
}

func applyLifecycle(t *testing.T, manager LifecycleManager, moduleID string, action LifecycleAction, operationID string) LifecycleResult {
	t.Helper()
	result, err := manager.Apply(context.Background(), LifecycleRequest{
		ModuleID: moduleID, Action: action, OperationID: operationID,
	})
	if err != nil {
		t.Fatalf("Apply(%s, %s) error = %v", moduleID, action, err)
	}
	return result
}

func TestLifecycleStateMachineSupportsNonDestructiveDisableAndReenable(t *testing.T) {
	registry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "1.0.0"))
	store := NewMemoryLifecycleStore()
	audit := &lifecycleAuditRecorder{}
	manager := lifecycleManagerFor(registry, store, audit)

	if _, err := manager.Apply(context.Background(), LifecycleRequest{ModuleID: "omnexa.alpha", Action: LifecycleEnable, OperationID: "op-invalid"}); lifecycleCode(t, err) != "lifecycle.transition.invalid" {
		t.Fatalf("unexpected invalid-transition error: %v", err)
	}
	installed := applyLifecycle(t, manager, "omnexa.alpha", LifecycleInstall, "op-install")
	if installed.Record.State != LifecycleInstalled || installed.Record.Version != "1.0.0" {
		t.Fatalf("unexpected install record: %#v", installed.Record)
	}
	enabled := applyLifecycle(t, manager, "omnexa.alpha", LifecycleEnable, "op-enable")
	if enabled.Record.State != LifecycleEnabled {
		t.Fatalf("expected enabled, got %#v", enabled.Record)
	}
	disabled := applyLifecycle(t, manager, "omnexa.alpha", LifecycleDisable, "op-disable")
	if disabled.Record.State != LifecycleDisabled || disabled.Record.Version != "1.0.0" {
		t.Fatalf("disable must preserve version/data-bearing state: %#v", disabled.Record)
	}
	reenabled := applyLifecycle(t, manager, "omnexa.alpha", LifecycleEnable, "op-reenable")
	if reenabled.Record.State != LifecycleEnabled || reenabled.Record.Version != "1.0.0" {
		t.Fatalf("unexpected re-enable record: %#v", reenabled.Record)
	}
	if len(audit.events) != 4 {
		t.Fatalf("expected one audit event per mutation, got %d", len(audit.events))
	}
}

func TestLifecycleDependencyPreconditionsAndReverseProtection(t *testing.T) {
	provider := simpleManifestV2("omnexa.provider", "1.0.0")
	consumer := simpleManifestV2("omnexa.consumer", "1.0.0")
	consumer.Dependencies = []DependencyRequirement{{ID: provider.ID, Constraint: "=1.0.0"}}
	registry := discoverV2Registry(t, consumer, provider)
	store := NewMemoryLifecycleStore()
	manager := lifecycleManagerFor(registry, store, &lifecycleAuditRecorder{})

	if _, err := manager.Apply(context.Background(), LifecycleRequest{ModuleID: consumer.ID, Action: LifecycleInstall, OperationID: "consumer-too-early"}); lifecycleCode(t, err) != "lifecycle.dependency.not_installed" {
		t.Fatalf("unexpected dependency-install error: %v", err)
	}
	applyLifecycle(t, manager, provider.ID, LifecycleInstall, "provider-install")
	applyLifecycle(t, manager, consumer.ID, LifecycleInstall, "consumer-install")
	if _, err := manager.Apply(context.Background(), LifecycleRequest{ModuleID: consumer.ID, Action: LifecycleEnable, OperationID: "consumer-enable-too-early"}); lifecycleCode(t, err) != "lifecycle.dependency.not_enabled" {
		t.Fatalf("unexpected dependency-enable error: %v", err)
	}
	applyLifecycle(t, manager, provider.ID, LifecycleEnable, "provider-enable")
	applyLifecycle(t, manager, consumer.ID, LifecycleEnable, "consumer-enable")

	if _, err := manager.Apply(context.Background(), LifecycleRequest{ModuleID: provider.ID, Action: LifecycleDisable, OperationID: "provider-disable-blocked"}); lifecycleCode(t, err) != "lifecycle.reverse_dependency.active" {
		t.Fatalf("unexpected reverse-active error: %v", err)
	}
	applyLifecycle(t, manager, consumer.ID, LifecycleDisable, "consumer-disable")
	applyLifecycle(t, manager, provider.ID, LifecycleDisable, "provider-disable")
	if _, err := manager.Apply(context.Background(), LifecycleRequest{ModuleID: provider.ID, Action: LifecycleDetach, OperationID: "provider-detach-blocked"}); lifecycleCode(t, err) != "lifecycle.reverse_dependency.present" {
		t.Fatalf("unexpected reverse-present error: %v", err)
	}
	applyLifecycle(t, manager, consumer.ID, LifecycleDetach, "consumer-detach")
	applyLifecycle(t, manager, provider.ID, LifecycleDetach, "provider-detach")
}

func TestLifecycleSuspendArchiveAndRestorePreserveHistoricalState(t *testing.T) {
	registry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "1.0.0"))
	manager := lifecycleManagerFor(registry, NewMemoryLifecycleStore(), &lifecycleAuditRecorder{})
	applyLifecycle(t, manager, "omnexa.alpha", LifecycleInstall, "install")
	applyLifecycle(t, manager, "omnexa.alpha", LifecycleEnable, "enable")
	suspended := applyLifecycle(t, manager, "omnexa.alpha", LifecycleSuspend, "suspend")
	if suspended.Record.State != LifecycleSuspended || suspended.Record.PreviousState != LifecycleEnabled {
		t.Fatalf("suspend must preserve resume state: %#v", suspended.Record)
	}
	resumed := applyLifecycle(t, manager, "omnexa.alpha", LifecycleResume, "resume")
	if resumed.Record.State != LifecycleEnabled || resumed.Record.PreviousState != "" {
		t.Fatalf("unexpected resume record: %#v", resumed.Record)
	}
	applyLifecycle(t, manager, "omnexa.alpha", LifecycleDisable, "disable")
	archived := applyLifecycle(t, manager, "omnexa.alpha", LifecycleArchive, "archive")
	if archived.Record.State != LifecycleArchived || archived.Record.PreviousState != LifecycleDisabled {
		t.Fatalf("archive must preserve prior state: %#v", archived.Record)
	}
	restored := applyLifecycle(t, manager, "omnexa.alpha", LifecycleRestore, "restore")
	if restored.Record.State != LifecycleDisabled {
		t.Fatalf("restore should return to historical state: %#v", restored.Record)
	}
}

func TestLifecyclePurgeRequiresAuthorizationAuditAndDetachedState(t *testing.T) {
	registry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "1.0.0"))
	store := NewMemoryLifecycleStore()
	audit := &lifecycleAuditRecorder{}
	manager := lifecycleManagerFor(registry, store, audit)
	applyLifecycle(t, manager, "omnexa.alpha", LifecycleInstall, "install")
	applyLifecycle(t, manager, "omnexa.alpha", LifecycleDetach, "detach")

	denied := manager
	denied.Authorizer = lifecycleAllowAuthorizer{deny: LifecyclePurge}
	if _, err := denied.Apply(context.Background(), LifecycleRequest{ModuleID: "omnexa.alpha", Action: LifecyclePurge, OperationID: "purge-denied"}); lifecycleCode(t, err) != "lifecycle.authorization.denied" {
		t.Fatalf("unexpected purge authorization error: %v", err)
	}
	record, _, _ := store.Load(context.Background(), "omnexa.alpha")
	if record.State != LifecycleDetached {
		t.Fatalf("denied purge mutated state: %#v", record)
	}

	failingAudit := manager
	failingAudit.Auditor = &lifecycleAuditRecorder{fail: true}
	if _, err := failingAudit.Apply(context.Background(), LifecycleRequest{ModuleID: "omnexa.alpha", Action: LifecyclePurge, OperationID: "purge-audit-fail"}); lifecycleCode(t, err) != "lifecycle.audit.failed" {
		t.Fatalf("unexpected purge audit error: %v", err)
	}
	record, _, _ = store.Load(context.Background(), "omnexa.alpha")
	if record.State != LifecycleDetached {
		t.Fatalf("audit-failed purge mutated state: %#v", record)
	}

	purged := applyLifecycle(t, manager, "omnexa.alpha", LifecyclePurge, "purge-ok")
	if purged.Record.State != LifecyclePurged {
		t.Fatalf("expected purged state, got %#v", purged.Record)
	}
}

func TestLifecycleOperationReplayAndConcurrentConflictAreExplicit(t *testing.T) {
	registry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "1.0.0"))
	store := NewMemoryLifecycleStore()
	manager := lifecycleManagerFor(registry, store, &lifecycleAuditRecorder{})
	first := applyLifecycle(t, manager, "omnexa.alpha", LifecycleInstall, "same-op")
	second := applyLifecycle(t, manager, "omnexa.alpha", LifecycleInstall, "same-op")
	if first.Replayed || !second.Replayed || second.Record.Revision != first.Record.Revision {
		t.Fatalf("unexpected replay results: first=%#v second=%#v", first, second)
	}
	if _, err := manager.Apply(context.Background(), LifecycleRequest{ModuleID: "omnexa.alpha", Action: LifecycleEnable, OperationID: "same-op"}); lifecycleCode(t, err) != "lifecycle.operation.reused" {
		t.Fatalf("unexpected operation reuse error: %v", err)
	}

	conflicting := manager
	conflicting.Store = lifecycleConflictStore{base: store}
	if _, err := conflicting.Apply(context.Background(), LifecycleRequest{ModuleID: "omnexa.alpha", Action: LifecycleEnable, OperationID: "conflict"}); lifecycleCode(t, err) != "lifecycle.concurrent_conflict" {
		t.Fatalf("unexpected concurrency error: %v", err)
	}
	record, _, _ := store.Load(context.Background(), "omnexa.alpha")
	if record.State != LifecycleInstalled {
		t.Fatalf("CAS conflict mutated state: %#v", record)
	}
}

func TestLifecycleFailureAndRecoveryPreserveStableState(t *testing.T) {
	registry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "1.0.0"))
	store := NewMemoryLifecycleStore()
	manager := lifecycleManagerFor(registry, store, &lifecycleAuditRecorder{})
	applyLifecycle(t, manager, "omnexa.alpha", LifecycleInstall, "install")

	failed, err := manager.MarkRecoveryRequired(context.Background(), LifecycleFailureRequest{
		ModuleID: "omnexa.alpha", FailedAction: LifecycleEnable, OperationID: "enable-failure", FailureCode: "hook.timeout",
	})
	if err != nil {
		t.Fatalf("MarkRecoveryRequired() error = %v", err)
	}
	if failed.State != LifecycleRecoveryRequired || failed.RecoveryState != LifecycleInstalled || failed.FailureCode != "hook.timeout" {
		t.Fatalf("unexpected recovery-required record: %#v", failed)
	}
	recovered := applyLifecycle(t, manager, "omnexa.alpha", LifecycleRecover, "recover")
	if recovered.Record.State != LifecycleInstalled || recovered.Record.FailureCode != "" || recovered.Record.RecoveryState != "" {
		t.Fatalf("unexpected recovered record: %#v", recovered.Record)
	}
}

func TestLifecycleUpgradeUsesFutureCoordinatorWithoutExecutingMigrations(t *testing.T) {
	store := NewMemoryLifecycleStore()
	oldRegistry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "1.0.0"))
	oldManager := lifecycleManagerFor(oldRegistry, store, &lifecycleAuditRecorder{})
	applyLifecycle(t, oldManager, "omnexa.alpha", LifecycleInstall, "install")

	newRegistry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "2.0.0"))
	newManager := lifecycleManagerFor(newRegistry, store, &lifecycleAuditRecorder{})
	if _, err := newManager.Apply(context.Background(), LifecycleRequest{ModuleID: "omnexa.alpha", Action: LifecycleUpgrade, OperationID: "upgrade-no-coordinator"}); lifecycleCode(t, err) != "lifecycle.upgrade.coordinator_unavailable" {
		t.Fatalf("unexpected missing-coordinator error: %v", err)
	}
	coordinator := &lifecycleUpgradeAllow{}
	newManager.UpgradeCoordinator = coordinator
	upgraded := applyLifecycle(t, newManager, "omnexa.alpha", LifecycleUpgrade, "upgrade-ok")
	if !coordinator.called || upgraded.Record.Version != "2.0.0" || upgraded.Record.State != LifecycleInstalled {
		t.Fatalf("unexpected upgrade result: called=%v record=%#v", coordinator.called, upgraded.Record)
	}
}

func TestLifecycleGraphFailureDoesNotMutateExistingState(t *testing.T) {
	consumer := simpleManifestV2("omnexa.consumer", "1.0.0")
	consumer.Dependencies = []DependencyRequirement{{ID: "omnexa.missing", Constraint: "=1.0.0"}}
	registry := discoverV2Registry(t, consumer)
	store := NewMemoryLifecycleStore()
	manager := lifecycleManagerFor(registry, store, &lifecycleAuditRecorder{})
	if _, err := manager.Apply(context.Background(), LifecycleRequest{ModuleID: consumer.ID, Action: LifecycleInstall, OperationID: "install"}); lifecycleCode(t, err) != "lifecycle.dependency.graph_invalid" {
		t.Fatalf("unexpected graph error: %v", err)
	}
	if _, found, _ := store.Load(context.Background(), consumer.ID); found {
		t.Fatal("failed graph eligibility must not create lifecycle state")
	}
}

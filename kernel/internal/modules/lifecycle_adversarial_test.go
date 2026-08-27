package modules

import (
	"context"
	"errors"
	"testing"
)

type lifecycleSelectiveReadErrorStore struct {
	base     LifecycleStore
	moduleID string
}

func (s lifecycleSelectiveReadErrorStore) Load(ctx context.Context, moduleID string) (LifecycleRecord, bool, error) {
	if moduleID == s.moduleID {
		return LifecycleRecord{}, false, errors.New("synthetic read failure")
	}
	return s.base.Load(ctx, moduleID)
}

func (s lifecycleSelectiveReadErrorStore) CompareAndSwap(ctx context.Context, moduleID string, expectedRevision uint64, next LifecycleRecord) error {
	return s.base.CompareAndSwap(ctx, moduleID, expectedRevision, next)
}

func seedLifecycleRecord(t *testing.T, store LifecycleStore, record LifecycleRecord) {
	t.Helper()
	if record.Revision == 0 {
		record.Revision = 1
	}
	if err := store.CompareAndSwap(context.Background(), record.ModuleID, 0, record); err != nil {
		t.Fatalf("seed lifecycle record: %v", err)
	}
}

func TestLifecycleReplayStillRequiresCurrentAuthorization(t *testing.T) {
	registry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "1.0.0"))
	store := NewMemoryLifecycleStore()
	manager := lifecycleManagerFor(registry, store, &lifecycleAuditRecorder{})
	applyLifecycle(t, manager, "omnexa.alpha", LifecycleInstall, "install-once")

	manager.Authorizer = lifecycleAllowAuthorizer{deny: LifecycleInstall}
	_, err := manager.Apply(context.Background(), LifecycleRequest{
		ModuleID: "omnexa.alpha", Action: LifecycleInstall, OperationID: "install-once",
	})
	if lifecycleCode(t, err) != "lifecycle.authorization.denied" {
		t.Fatalf("replay must reauthorize current caller, got %v", err)
	}
}

func TestLifecycleAuthorizationPrecedesUpgradeCoordinator(t *testing.T) {
	store := NewMemoryLifecycleStore()
	oldRegistry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "1.0.0"))
	oldManager := lifecycleManagerFor(oldRegistry, store, &lifecycleAuditRecorder{})
	applyLifecycle(t, oldManager, "omnexa.alpha", LifecycleInstall, "install-before-upgrade")

	newRegistry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "2.0.0"))
	coordinator := &lifecycleUpgradeAllow{}
	manager := lifecycleManagerFor(newRegistry, store, &lifecycleAuditRecorder{})
	manager.Authorizer = lifecycleAllowAuthorizer{deny: LifecycleUpgrade}
	manager.UpgradeCoordinator = coordinator

	_, err := manager.Apply(context.Background(), LifecycleRequest{
		ModuleID: "omnexa.alpha", Action: LifecycleUpgrade, OperationID: "unauthorized-upgrade",
	})
	if lifecycleCode(t, err) != "lifecycle.authorization.denied" {
		t.Fatalf("unauthorized upgrade must fail at authorization boundary, got %v", err)
	}
	if coordinator.called {
		t.Fatal("unauthorized upgrade must not invoke upgrade coordinator")
	}
}

func TestLifecycleEnableRejectsStaleInstalledModuleVersion(t *testing.T) {
	registry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "2.0.0"))
	store := NewMemoryLifecycleStore()
	seedLifecycleRecord(t, store, LifecycleRecord{
		ModuleID: "omnexa.alpha", Version: "1.0.0", State: LifecycleDisabled,
	})
	manager := lifecycleManagerFor(registry, store, &lifecycleAuditRecorder{})

	_, err := manager.Apply(context.Background(), LifecycleRequest{
		ModuleID: "omnexa.alpha", Action: LifecycleEnable, OperationID: "enable-stale",
	})
	if lifecycleCode(t, err) != "lifecycle.version.mismatch" {
		t.Fatalf("stale installed version must require upgrade, got %v", err)
	}
}

func TestLifecycleDependencyRequiresInstalledVersionBoundToResolverRegistry(t *testing.T) {
	provider := simpleManifestV2("omnexa.provider", "2.0.0")
	consumer := simpleManifestV2("omnexa.consumer", "1.0.0")
	consumer.Dependencies = []DependencyRequirement{{ID: provider.ID, Constraint: "=2.0.0"}}
	registry := discoverV2Registry(t, consumer, provider)
	store := NewMemoryLifecycleStore()
	seedLifecycleRecord(t, store, LifecycleRecord{
		ModuleID: provider.ID, Version: "1.0.0", State: LifecycleEnabled,
	})
	seedLifecycleRecord(t, store, LifecycleRecord{
		ModuleID: consumer.ID, Version: "1.0.0", State: LifecycleInstalled,
	})
	manager := lifecycleManagerFor(registry, store, &lifecycleAuditRecorder{})

	_, err := manager.Apply(context.Background(), LifecycleRequest{
		ModuleID: consumer.ID, Action: LifecycleEnable, OperationID: "enable-consumer",
	})
	if lifecycleCode(t, err) != "lifecycle.dependency.version_mismatch" {
		t.Fatalf("dependency state must bind exact installed version to resolver registry, got %v", err)
	}
}

func TestLifecycleReverseDependencyReadFailureFailsClosed(t *testing.T) {
	provider := simpleManifestV2("omnexa.provider", "1.0.0")
	consumer := simpleManifestV2("omnexa.consumer", "1.0.0")
	consumer.Dependencies = []DependencyRequirement{{ID: provider.ID, Constraint: "=1.0.0"}}
	registry := discoverV2Registry(t, consumer, provider)
	base := NewMemoryLifecycleStore()
	seedLifecycleRecord(t, base, LifecycleRecord{
		ModuleID: provider.ID, Version: "1.0.0", State: LifecycleEnabled,
	})
	seedLifecycleRecord(t, base, LifecycleRecord{
		ModuleID: consumer.ID, Version: "1.0.0", State: LifecycleEnabled,
	})
	manager := lifecycleManagerFor(registry, lifecycleSelectiveReadErrorStore{base: base, moduleID: consumer.ID}, &lifecycleAuditRecorder{})

	_, err := manager.Apply(context.Background(), LifecycleRequest{
		ModuleID: provider.ID, Action: LifecycleDisable, OperationID: "disable-provider",
	})
	if lifecycleCode(t, err) != "lifecycle.store.read_failed" {
		t.Fatalf("reverse-dependency read failure must fail closed, got %v", err)
	}
	record, _, _ := base.Load(context.Background(), provider.ID)
	if record.State != LifecycleEnabled {
		t.Fatalf("read failure mutated provider state: %#v", record)
	}
}

func TestLifecycleFailedInstallCanEnterAndRecoverFromRecoveryRequired(t *testing.T) {
	registry := discoverV2Registry(t, simpleManifestV2("omnexa.alpha", "1.0.0"))
	store := NewMemoryLifecycleStore()
	manager := lifecycleManagerFor(registry, store, &lifecycleAuditRecorder{})

	failed, err := manager.MarkRecoveryRequired(context.Background(), LifecycleFailureRequest{
		ModuleID: "omnexa.alpha", FailedAction: LifecycleInstall, OperationID: "install-failed", FailureCode: "hook.timeout",
	})
	if err != nil {
		t.Fatalf("failed install must be representable as recovery-required: %v", err)
	}
	if failed.State != LifecycleRecoveryRequired || failed.RecoveryState != LifecycleAvailable {
		t.Fatalf("unexpected failed-install recovery state: %#v", failed)
	}
	recovered := applyLifecycle(t, manager, "omnexa.alpha", LifecycleRecover, "recover-install")
	if recovered.Record.State != LifecycleAvailable {
		t.Fatalf("failed install recovery must restore available state: %#v", recovered.Record)
	}
}

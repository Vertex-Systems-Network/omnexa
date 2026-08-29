package modules

import (
	"context"
	"reflect"
	"testing"
)

// TestP03ExitReferenceModuleAggregate composes the already-accepted P03 runtime
// contracts into the canonical EX-01..EX-07 phase-exit proof. The subtests call
// retained tests directly so this aggregate cannot replace or weaken their
// independent execution in the package suite.
func TestP03ExitReferenceModuleAggregate(t *testing.T) {
	t.Run("EX-01-required-dependency-enforcement", func(t *testing.T) {
		TestResolveDependenciesRejectsMissingAndIncompatibleRequiredDependency(t)
		TestLifecycleDependencyPreconditionsAndReverseProtection(t)
	})
	t.Run("EX-02-optional-dependency-degradation", func(t *testing.T) {
		TestResolveDependenciesDegradesOptionalDependencyWithoutGlobalFailure(t)
		TestModuleHealthReporterOptionalDependencySelectiveDegradation(t)
	})
	t.Run("EX-03-safe-disable-reenable", func(t *testing.T) {
		TestLifecycleStateMachineSupportsNonDestructiveDisableAndReenable(t)
	})
	t.Run("EX-04-upgrade-migration-path", func(t *testing.T) {
		TestLifecycleUpgradeUsesFutureCoordinatorWithoutExecutingMigrations(t)
		TestMigrationOwnershipRegistryDeterministicFreshInstallAndUpgradePlan(t)
	})
	t.Run("EX-05-forbidden-dependency-detection", func(t *testing.T) {
		TestResolveDependenciesRejectsRequiredCycleButIgnoresOptionalCycleForGlobalOrder(t)
		TestResolveDependenciesRejectsUndeclaredPrivateAndKernelToBusinessObservations(t)
	})
	t.Run("EX-06-health-state-accuracy", func(t *testing.T) {
		TestModuleHealthReporterHealthyAndDeterministic(t)
		TestModuleHealthReporterRequiredDependencyFailsClosed(t)
		TestModuleHealthReporterMigrationPendingAndInconsistentNeverHealthy(t)
	})
	t.Run("EX-07-unrelated-module-isolation", func(t *testing.T) {
		TestP03ExitUnrelatedModuleIsolationAcrossLifecycleOperations(t)
		TestModuleHealthReporterFailureIsolationAndLifecycleChanges(t)
		TestLifecycleGraphFailureDoesNotMutateExistingState(t)
	})
}

func TestP03ExitUnrelatedModuleIsolationAcrossLifecycleOperations(t *testing.T) {
	alpha := simpleManifestV2("omnexa.alpha", "1.0.0")
	beta := simpleManifestV2("omnexa.beta", "1.0.0")
	registry := discoverV2Registry(t, alpha, beta)
	store := NewMemoryLifecycleStore()
	manager := lifecycleManagerFor(registry, store, &lifecycleAuditRecorder{})

	applyLifecycle(t, manager, beta.ID, LifecycleInstall, "beta-install")
	applyLifecycle(t, manager, beta.ID, LifecycleEnable, "beta-enable")
	baseline, found, err := store.Load(context.Background(), beta.ID)
	if err != nil || !found {
		t.Fatalf("load unrelated beta baseline: found=%v err=%v", found, err)
	}
	assertBetaUnchanged := func(stage string) {
		t.Helper()
		current, currentFound, loadErr := store.Load(context.Background(), beta.ID)
		if loadErr != nil || !currentFound {
			t.Fatalf("%s: load unrelated beta: found=%v err=%v", stage, currentFound, loadErr)
		}
		if !reflect.DeepEqual(current, baseline) {
			t.Fatalf("%s mutated unrelated beta state:\n got=%#v\nwant=%#v", stage, current, baseline)
		}
	}

	applyLifecycle(t, manager, alpha.ID, LifecycleInstall, "alpha-install")
	assertBetaUnchanged("install")

	failed, err := manager.MarkRecoveryRequired(context.Background(), LifecycleFailureRequest{
		ModuleID: alpha.ID, FailedAction: LifecycleEnable, OperationID: "alpha-enable-failure", FailureCode: "hook.timeout",
	})
	if err != nil || failed.State != LifecycleRecoveryRequired {
		t.Fatalf("mark alpha recovery required: state=%q err=%v", failed.State, err)
	}
	assertBetaUnchanged("failed-operation")
	applyLifecycle(t, manager, alpha.ID, LifecycleRecover, "alpha-recover")
	assertBetaUnchanged("recover")

	applyLifecycle(t, manager, alpha.ID, LifecycleEnable, "alpha-enable")
	assertBetaUnchanged("enable")
	applyLifecycle(t, manager, alpha.ID, LifecycleSuspend, "alpha-suspend")
	assertBetaUnchanged("suspend")
	applyLifecycle(t, manager, alpha.ID, LifecycleResume, "alpha-resume")
	assertBetaUnchanged("resume")
	applyLifecycle(t, manager, alpha.ID, LifecycleDisable, "alpha-disable")
	assertBetaUnchanged("disable")
	applyLifecycle(t, manager, alpha.ID, LifecycleArchive, "alpha-archive")
	assertBetaUnchanged("archive")
	applyLifecycle(t, manager, alpha.ID, LifecycleRestore, "alpha-restore")
	assertBetaUnchanged("restore")
	applyLifecycle(t, manager, alpha.ID, LifecycleDetach, "alpha-detach")
	assertBetaUnchanged("detach")
	applyLifecycle(t, manager, alpha.ID, LifecyclePurge, "alpha-purge")
	assertBetaUnchanged("purge")
}

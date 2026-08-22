package operations

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/buildinfo"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/observability"
)

func TestLivenessReadinessAreDistinct(t *testing.T) {
	manager := NewManager(nil)
	mustRegister(t, manager, Dependency{Name: "database", Criticality: CriticalityRequired, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { return errors.New("database unavailable") }})
	starting := manager.Evaluate(context.Background())
	if starting.Liveness != StateHealthy || starting.Readiness != StateStarting {
		t.Fatalf("starting report = %+v", starting)
	}
	if !manager.MarkReady() { t.Fatal("expected startup -> ready transition") }
	ready := manager.Evaluate(context.Background())
	if ready.Liveness != StateHealthy || ready.Readiness != StateUnready {
		t.Fatalf("required failure report = %+v", ready)
	}
	if got := ready.Dependencies[0]; got.State != StateUnready || got.Reason != ReasonFailed { t.Fatalf("result = %+v", got) }
}

func TestOptionalDependencyProducesDegradedReadiness(t *testing.T) {
	manager := NewManager(nil)
	mustRegister(t, manager, Dependency{Name: "telemetry-backend", Criticality: CriticalityOptional, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { return errors.New("optional backend unavailable") }})
	manager.MarkReady()
	report := manager.Evaluate(context.Background())
	if report.Liveness != StateHealthy || report.Readiness != StateDegraded { t.Fatalf("report = %+v", report) }
	if got := report.Dependencies[0]; got.State != StateDegraded || got.Reason != ReasonFailed { t.Fatalf("result = %+v", got) }
}

func TestRequiredAndSecurityCriticalFailuresFailClosed(t *testing.T) {
	for _, criticality := range []Criticality{CriticalityRequired, CriticalitySecurityCritical} {
		t.Run(string(criticality), func(t *testing.T) {
			manager := NewManager(nil)
			mustRegister(t, manager, Dependency{Name: "policy", Criticality: criticality, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { return errors.New("dependency rejected request") }})
			manager.MarkReady()
			report := manager.Evaluate(context.Background())
			if report.Readiness != StateUnready || report.Liveness != StateHealthy { t.Fatalf("report = %+v", report) }
		})
	}
}

func TestDependencyChecksAreTimeoutBoundedEvenWhenProbeIgnoresContext(t *testing.T) {
	release := make(chan struct{})
	manager := NewManager(nil)
	mustRegister(t, manager, Dependency{Name: "slow-provider", Criticality: CriticalityRequired, Timeout: 20 * time.Millisecond, Check: func(context.Context) error { <-release; return nil }})
	manager.MarkReady()
	started := time.Now()
	report := manager.Evaluate(context.Background())
	elapsed := time.Since(started)
	close(release)
	if elapsed > 500*time.Millisecond { t.Fatalf("evaluation took %s", elapsed) }
	if report.Readiness != StateUnready || report.Dependencies[0].Reason != ReasonTimeout { t.Fatalf("report = %+v", report) }
}

func TestCallerCancellationIsBoundedAndClassified(t *testing.T) {
	release := make(chan struct{})
	manager := NewManager(nil)
	mustRegister(t, manager, Dependency{Name: "provider", Criticality: CriticalityRequired, Timeout: time.Second, Check: func(context.Context) error { <-release; return nil }})
	manager.MarkReady()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := time.Now()
	report := manager.Evaluate(ctx)
	close(release)
	if time.Since(started) > 500*time.Millisecond { t.Fatal("canceled evaluation was not bounded") }
	if got := report.Dependencies[0]; got.Reason != ReasonCanceled || got.State != StateUnready { t.Fatalf("result = %+v", got) }
}

func TestDiagnosticSummaryAndObservabilityDoNotExposeProbeErrors(t *testing.T) {
	const secret = "postgres://admin:super-secret@db.internal/omnexa?token=abc123"
	logger, capture := captureLogger(t)
	manager := NewManager(logger)
	mustRegister(t, manager, Dependency{Name: "database", Criticality: CriticalityRequired, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { return errors.New(secret + " object-key=tenant/private/file.pdf") }})
	manager.MarkReady()
	report := manager.Evaluate(context.Background())
	encoded, err := json.Marshal(report)
	if err != nil { t.Fatalf("marshal report: %v", err) }
	for _, forbidden := range []string{"super-secret", "abc123", "tenant/private"} {
		if strings.Contains(string(encoded), forbidden) { t.Fatalf("report leaked %q: %s", forbidden, encoded) }
	}
	if !strings.Contains(string(encoded), `"reason":"failed"`) { t.Fatalf("missing stable reason: %s", encoded) }
	records := capture.Records()
	if len(records) != 1 { t.Fatalf("captured records = %d", len(records)) }
	serialized := records[0].Message
	for key, value := range records[0].Attributes { serialized += key + stringify(value) }
	for _, forbidden := range []string{"super-secret", "abc123", "tenant/private"} {
		if strings.Contains(serialized, forbidden) { t.Fatalf("log leaked %q: %s", forbidden, serialized) }
	}
	if records[0].Severity != slog.LevelWarn { t.Fatalf("severity = %v", records[0].Severity) }
}

func TestBuildIdentityIsIncludedWithoutMachineMetadata(t *testing.T) {
	oldVersion, oldCommit := buildinfo.Version, buildinfo.Commit
	buildinfo.Version, buildinfo.Commit = "1.2.3", "abc123"
	t.Cleanup(func() { buildinfo.Version, buildinfo.Commit = oldVersion, oldCommit })
	manager := NewManager(nil)
	manager.MarkReady()
	report := manager.Evaluate(context.Background())
	if report.Build.Version != "1.2.3" || report.Build.Commit != "abc123" { t.Fatalf("build = %+v", report.Build) }
}

func TestRegistryValidationDuplicateAndFreezeAreSafeFailures(t *testing.T) {
	manager := NewManager(nil)
	valid := Dependency{Name: "database", Criticality: CriticalityRequired, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { return nil }}
	mustRegister(t, manager, valid)
	assertFailureCode(t, manager.Register(valid), codeDependencyDuplicate)
	assertFailureCode(t, manager.Register(Dependency{Name: "../secret", Criticality: CriticalityRequired, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { return nil }}), codeDependencyInvalid)
	manager.MarkReady()
	assertFailureCode(t, manager.Register(Dependency{Name: "cache", Criticality: CriticalityOptional, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { return nil }}), codeRegistryFrozen)
}

func TestDependencyResultsAreDeterministicallyOrdered(t *testing.T) {
	manager := NewManager(nil)
	for _, name := range []string{"storage", "cache", "database"} {
		mustRegister(t, manager, Dependency{Name: name, Criticality: CriticalityRequired, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { return nil }})
	}
	manager.MarkReady()
	report := manager.Evaluate(context.Background())
	got := []string{report.Dependencies[0].Name, report.Dependencies[1].Name, report.Dependencies[2].Name}
	want := []string{"cache", "database", "storage"}
	for index := range want { if got[index] != want[index] { t.Fatalf("order = %v, want %v", got, want) } }
	if report.Readiness != StateHealthy { t.Fatalf("readiness = %q", report.Readiness) }
}

func TestLifecycleStoppingAndFailedSemantics(t *testing.T) {
	manager := NewManager(nil)
	manager.MarkReady()
	if !manager.MarkStopping() { t.Fatal("expected stopping transition") }
	stopping := manager.Evaluate(context.Background())
	if stopping.Liveness != StateHealthy || stopping.Readiness != StateUnready { t.Fatalf("stopping = %+v", stopping) }
	if manager.MarkReady() { t.Fatal("stopping must not return to ready") }
	if !manager.MarkFailed() { t.Fatal("expected failed transition") }
	failed := manager.Evaluate(context.Background())
	if failed.Liveness != StateUnready || failed.Readiness != StateUnready { t.Fatalf("failed = %+v", failed) }
}

func TestConcurrentEvaluateAndLifecycleReadsAreRaceSafe(t *testing.T) {
	manager := NewManager(nil)
	mustRegister(t, manager, Dependency{Name: "database", Criticality: CriticalityRequired, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { return nil }})
	var workers sync.WaitGroup
	for index := 0; index < 8; index++ {
		workers.Add(1)
		go func() { defer workers.Done(); for iteration := 0; iteration < 20; iteration++ { _ = manager.Liveness(); _ = manager.Evaluate(context.Background()) } }()
	}
	manager.MarkReady()
	workers.Wait()
}

func TestPanickingProbeFailsSafely(t *testing.T) {
	manager := NewManager(nil)
	mustRegister(t, manager, Dependency{Name: "provider", Criticality: CriticalityRequired, Timeout: 100 * time.Millisecond, Check: func(context.Context) error { panic("secret panic payload") }})
	manager.MarkReady()
	report := manager.Evaluate(context.Background())
	if got := report.Dependencies[0]; got.Reason != ReasonFailed || got.State != StateUnready { t.Fatalf("result = %+v", got) }
	encoded, _ := json.Marshal(report)
	if strings.Contains(string(encoded), "secret panic payload") { t.Fatalf("panic leaked: %s", encoded) }
}

func mustRegister(t *testing.T, manager *Manager, dependency Dependency) {
	t.Helper()
	if err := manager.Register(dependency); err != nil { t.Fatalf("register %s: %v", dependency.Name, err) }
}

func assertFailureCode(t *testing.T, err error, code failure.Code) {
	t.Helper()
	if err == nil { t.Fatalf("expected failure code %s", code) }
	classified, ok := failure.As(err)
	if !ok { t.Fatalf("not Omnexa failure: %v", err) }
	if classified.Code() != code { t.Fatalf("code = %s, want %s", classified.Code(), code) }
}

func captureLogger(t *testing.T) (*observability.Logger, *observability.Capture) {
	t.Helper()
	settings := observability.Settings{Enabled: true, ServiceName: "omnexa-kernel", ServiceVersion: "test", Environment: config.EnvironmentTest, LogLevel: slog.LevelDebug, ExportInterval: time.Second, ExportTimeout: time.Second, ShutdownTimeout: time.Second}
	return observability.NewCaptureLogger(settings)
}

func stringify(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil { return "" }
	return string(encoded)
}

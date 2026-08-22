// Package operations implements the bounded P01.08 health, readiness and diagnostics foundation.
package operations

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/buildinfo"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/observability"
)

const (
	minCheckTimeout = 10 * time.Millisecond
	maxCheckTimeout = 30 * time.Second
)

var dependencyNamePattern = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)

type State string

const (
	StateStarting State = "starting"
	StateHealthy  State = "healthy"
	StateDegraded State = "degraded"
	StateUnready  State = "unready"
)

func (state State) Valid() bool {
	switch state {
	case StateStarting, StateHealthy, StateDegraded, StateUnready:
		return true
	default:
		return false
	}
}

type Lifecycle string

const (
	LifecycleStarting Lifecycle = "starting"
	LifecycleReady    Lifecycle = "ready"
	LifecycleStopping Lifecycle = "stopping"
	LifecycleFailed   Lifecycle = "failed"
)

type Criticality string

const (
	CriticalityRequired         Criticality = "required"
	CriticalityOptional         Criticality = "optional"
	CriticalitySecurityCritical Criticality = "security_critical"
)

func (criticality Criticality) Valid() bool {
	switch criticality {
	case CriticalityRequired, CriticalityOptional, CriticalitySecurityCritical:
		return true
	default:
		return false
	}
}

type ResultReason string

const (
	ReasonOK       ResultReason = "ok"
	ReasonFailed   ResultReason = "failed"
	ReasonTimeout  ResultReason = "timeout"
	ReasonCanceled ResultReason = "canceled"
)

type Check func(context.Context) error

type Dependency struct {
	Name        string
	Criticality Criticality
	Timeout     time.Duration
	Check       Check
}

type DependencyResult struct {
	Name        string       `json:"name"`
	Criticality Criticality  `json:"criticality"`
	State       State        `json:"state"`
	Reason      ResultReason `json:"reason"`
}

type BuildIdentity struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

type Report struct {
	Build        BuildIdentity      `json:"build"`
	Lifecycle    Lifecycle          `json:"lifecycle"`
	Liveness     State              `json:"liveness"`
	Readiness    State              `json:"readiness"`
	Dependencies []DependencyResult `json:"dependencies"`
}

type Manager struct {
	mu           sync.RWMutex
	dependencies map[string]Dependency
	lifecycle    Lifecycle
	logger       *observability.Logger
}

func NewManager(logger *observability.Logger) *Manager {
	return &Manager{dependencies: make(map[string]Dependency), lifecycle: LifecycleStarting, logger: logger}
}

func (manager *Manager) Register(dependency Dependency) error {
	if manager == nil {
		return safeFailure(codeRegistryInvalid, "health dependency registry is invalid")
	}
	if !dependencyNamePattern.MatchString(dependency.Name) || !dependency.Criticality.Valid() || dependency.Check == nil || dependency.Timeout < minCheckTimeout || dependency.Timeout > maxCheckTimeout {
		return safeFailure(codeDependencyInvalid, "health dependency registration is invalid")
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if manager.lifecycle != LifecycleStarting {
		return safeFailure(codeRegistryFrozen, "health dependency registry is frozen")
	}
	if _, exists := manager.dependencies[dependency.Name]; exists {
		return safeFailure(codeDependencyDuplicate, "health dependency is already registered")
	}
	manager.dependencies[dependency.Name] = dependency
	return nil
}

func (manager *Manager) MarkReady() bool {
	if manager == nil {
		return false
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if manager.lifecycle != LifecycleStarting {
		return false
	}
	manager.lifecycle = LifecycleReady
	return true
}

func (manager *Manager) MarkStopping() bool {
	if manager == nil {
		return false
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if manager.lifecycle == LifecycleStopping || manager.lifecycle == LifecycleFailed {
		return false
	}
	manager.lifecycle = LifecycleStopping
	return true
}

func (manager *Manager) MarkFailed() bool {
	if manager == nil {
		return false
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if manager.lifecycle == LifecycleFailed {
		return false
	}
	manager.lifecycle = LifecycleFailed
	return true
}

func (manager *Manager) Lifecycle() Lifecycle {
	if manager == nil {
		return LifecycleFailed
	}
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.lifecycle
}

func (manager *Manager) Liveness() State { return livenessFor(manager.Lifecycle()) }

func (manager *Manager) Evaluate(ctx context.Context) Report {
	if ctx == nil {
		ctx = context.Background()
	}
	dependencies, lifecycle := manager.snapshot()
	results := evaluateDependencies(ctx, dependencies)
	report := Report{Build: currentBuildIdentity(), Lifecycle: lifecycle, Liveness: livenessFor(lifecycle), Readiness: readinessFor(lifecycle, results), Dependencies: results}
	manager.logEvaluation(ctx, report)
	return report
}

func (manager *Manager) snapshot() ([]Dependency, Lifecycle) {
	if manager == nil {
		return nil, LifecycleFailed
	}
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	dependencies := make([]Dependency, 0, len(manager.dependencies))
	for _, dependency := range manager.dependencies {
		dependencies = append(dependencies, dependency)
	}
	sort.Slice(dependencies, func(left, right int) bool { return dependencies[left].Name < dependencies[right].Name })
	return dependencies, manager.lifecycle
}

func evaluateDependencies(parent context.Context, dependencies []Dependency) []DependencyResult {
	if len(dependencies) == 0 {
		return []DependencyResult{}
	}
	results := make(chan DependencyResult, len(dependencies))
	for _, dependency := range dependencies {
		dependency := dependency
		go func() { results <- evaluateDependency(parent, dependency) }()
	}
	ordered := make([]DependencyResult, 0, len(dependencies))
	for range dependencies {
		ordered = append(ordered, <-results)
	}
	sort.Slice(ordered, func(left, right int) bool { return ordered[left].Name < ordered[right].Name })
	return ordered
}

func evaluateDependency(parent context.Context, dependency Dependency) DependencyResult {
	bounded, cancel := context.WithTimeout(parent, dependency.Timeout)
	defer cancel()
	completed := make(chan error, 1)
	go func() { completed <- runCheckSafely(bounded, dependency.Check) }()
	var reason ResultReason
	select {
	case err := <-completed:
		reason = classifyResultReason(parent, bounded, err)
	case <-bounded.Done():
		reason = classifyResultReason(parent, bounded, bounded.Err())
	}
	return DependencyResult{Name: dependency.Name, Criticality: dependency.Criticality, State: stateForDependency(dependency.Criticality, reason), Reason: reason}
}

func runCheckSafely(ctx context.Context, check Check) (err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("health dependency check failed safely")
		}
	}()
	return check(ctx)
}

func classifyResultReason(parent, bounded context.Context, err error) ResultReason {
	if err == nil {
		return ReasonOK
	}
	if errors.Is(parent.Err(), context.DeadlineExceeded) {
		return ReasonTimeout
	}
	if errors.Is(parent.Err(), context.Canceled) {
		return ReasonCanceled
	}
	if errors.Is(bounded.Err(), context.DeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
		return ReasonTimeout
	}
	if errors.Is(bounded.Err(), context.Canceled) || errors.Is(err, context.Canceled) {
		return ReasonCanceled
	}
	return ReasonFailed
}

func stateForDependency(criticality Criticality, reason ResultReason) State {
	if reason == ReasonOK {
		return StateHealthy
	}
	if criticality == CriticalityOptional {
		return StateDegraded
	}
	return StateUnready
}

func readinessFor(lifecycle Lifecycle, results []DependencyResult) State {
	switch lifecycle {
	case LifecycleStarting:
		return StateStarting
	case LifecycleStopping, LifecycleFailed:
		return StateUnready
	case LifecycleReady:
		readiness := StateHealthy
		for _, result := range results {
			if result.State == StateUnready {
				return StateUnready
			}
			if result.State == StateDegraded {
				readiness = StateDegraded
			}
		}
		return readiness
	default:
		return StateUnready
	}
}

func livenessFor(lifecycle Lifecycle) State {
	if lifecycle == LifecycleFailed {
		return StateUnready
	}
	return StateHealthy
}

func currentBuildIdentity() BuildIdentity {
	identity := buildinfo.Current()
	return BuildIdentity{Version: identity.Version, Commit: identity.Commit}
}

func (manager *Manager) logEvaluation(ctx context.Context, report Report) {
	if manager == nil || manager.logger == nil {
		return
	}
	attrs := []slog.Attr{slog.String("lifecycle", string(report.Lifecycle)), slog.String("liveness", string(report.Liveness)), slog.String("readiness", string(report.Readiness)), slog.Int("dependency_count", len(report.Dependencies))}
	if report.Readiness == StateUnready {
		manager.logger.Warn(ctx, "health evaluation completed", attrs...)
		return
	}
	manager.logger.Info(ctx, "health evaluation completed", attrs...)
}

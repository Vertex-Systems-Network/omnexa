package modules

import "sort"

// PlatformSnapshot is explicit resolver input describing the platform version
// and kernel capabilities available to the discovered module set.
type PlatformSnapshot struct {
	Version      string
	Capabilities []string
}

// DependencyObservation is a pure validation hook for dependency uses observed
// by a caller-owned static analysis step. The resolver itself performs no source
// scanning, filesystem access, network access, or package execution.
type DependencyObservation struct {
	ConsumerID string
	ProviderID string
	Private    bool
}

// OptionalDependencyStatus is selective-degradation metadata only. It grants no
// authority and does not participate in required-edge ordering.
type OptionalDependencyStatus struct {
	ModuleID     string `json:"module_id"`
	DependencyID string `json:"dependency_id"`
	State        string `json:"state"`
	Reason       string `json:"reason,omitempty"`
}

const (
	optionalAvailable    = "available"
	optionalMissing      = "missing"
	optionalIncompatible = "incompatible"
	optionalUnresolved   = "unresolved"
)

// DependencyResolution is deterministic eligibility metadata. Order contains
// dependencies before dependents and is derived from required edges only.
type DependencyResolution struct {
	Order    []string                   `json:"order"`
	Optional []OptionalDependencyStatus `json:"optional"`
}

// ResolutionDiagnostic is a stable, classification-safe resolver failure.
type ResolutionDiagnostic struct {
	Code               string `json:"code"`
	ModuleID           string `json:"module_id,omitempty"`
	DependencyID       string `json:"dependency_id,omitempty"`
	PlatformDependency string `json:"platform_dependency,omitempty"`
}

// ResolutionErrors is sorted deterministically and never includes raw manifest
// payloads, constraints, filesystem paths, or secret material.
type ResolutionErrors struct {
	diagnostics []ResolutionDiagnostic
}

func (e *ResolutionErrors) Error() string {
	return "module dependency resolution failed"
}

// Diagnostics returns a copy so callers cannot mutate resolver evidence.
func (e *ResolutionErrors) Diagnostics() []ResolutionDiagnostic {
	if e == nil {
		return nil
	}
	result := make([]ResolutionDiagnostic, len(e.diagnostics))
	copy(result, e.diagnostics)
	return result
}

// ResolveDependencies validates the exact registry-bound manifest snapshot and
// returns deterministic required-edge order plus optional degradation metadata.
// Resolver output never grants capabilities, permissions, tenant scope, private
// package access, or database authority.
func ResolveDependencies(registry Registry, platform PlatformSnapshot, observations []DependencyObservation) (DependencyResolution, error) {
	diagnostics := make([]ResolutionDiagnostic, 0)
	optional := make([]OptionalDependencyStatus, 0)

	_, platformVersionValid := parseStrictSemVer(platform.Version)
	if !platformVersionValid {
		diagnostics = append(diagnostics, ResolutionDiagnostic{Code: "resolver.platform.version_invalid"})
	}
	platformCapabilities, platformDiagnostics := validatePlatformSnapshot(platform)
	diagnostics = append(diagnostics, platformDiagnostics...)

	records := registry.List()
	nodes := make([]string, 0, len(records))
	recordByID := make(map[string]RegistryRecord, len(records))
	indegree := make(map[string]int, len(records))
	edges := make(map[string]map[string]struct{}, len(records))

	for _, record := range records {
		nodes = append(nodes, record.ID)
		recordByID[record.ID] = record
		indegree[record.ID] = 0
	}
	sort.Strings(nodes)

	for _, moduleID := range nodes {
		record := recordByID[moduleID]
		snapshot, ok := registry.manifestSnapshot(moduleID)
		if !ok {
			diagnostics = append(diagnostics, ResolutionDiagnostic{
				Code:     "resolver.registry.snapshot_missing",
				ModuleID: moduleID,
			})
			continue
		}
		if snapshot.ID != record.ID || snapshot.Version != record.Version || snapshot.Owner != record.Owner {
			diagnostics = append(diagnostics, ResolutionDiagnostic{
				Code:     "resolver.registry.snapshot_mismatch",
				ModuleID: moduleID,
			})
			continue
		}

		if platformVersionValid {
			compatible, valid := matchesDependencyConstraint(platform.Version, snapshot.RequiredPlatformVersion)
			if !valid {
				diagnostics = append(diagnostics, ResolutionDiagnostic{
					Code:     "resolver.platform.requirement_invalid",
					ModuleID: moduleID,
				})
			} else if !compatible {
				diagnostics = append(diagnostics, ResolutionDiagnostic{
					Code:     "resolver.platform.version_incompatible",
					ModuleID: moduleID,
				})
			}
		}
		for _, platformDependency := range snapshot.PlatformDependencies {
			if _, ok := platformCapabilities[platformDependency]; !ok {
				diagnostics = append(diagnostics, ResolutionDiagnostic{
					Code:               "resolver.platform.dependency_missing",
					ModuleID:           moduleID,
					PlatformDependency: platformDependency,
				})
			}
		}

		if snapshot.SchemaVersion == SchemaVersion {
			for _, requirement := range snapshot.RequiredDependencies {
				diagnostics = append(diagnostics, ResolutionDiagnostic{
					Code:         "resolver.dependency.version_contract_missing",
					ModuleID:     moduleID,
					DependencyID: requirement.ID,
				})
			}
			for _, requirement := range snapshot.OptionalDependencies {
				optional = append(optional, OptionalDependencyStatus{
					ModuleID:     moduleID,
					DependencyID: requirement.ID,
					State:        optionalUnresolved,
					Reason:       "version_contract_missing",
				})
			}
			continue
		}

		for _, requirement := range snapshot.RequiredDependencies {
			provider, exists := recordByID[requirement.ID]
			if !exists {
				diagnostics = append(diagnostics, ResolutionDiagnostic{
					Code:         "resolver.dependency.required_missing",
					ModuleID:     moduleID,
					DependencyID: requirement.ID,
				})
				continue
			}
			compatible, valid := matchesDependencyConstraint(provider.Version, requirement.Constraint)
			if !valid {
				code := "resolver.dependency.constraint_invalid"
				if _, ok := parseStrictSemVer(provider.Version); !ok {
					code = "resolver.dependency.version_invalid"
				}
				diagnostics = append(diagnostics, ResolutionDiagnostic{
					Code:         code,
					ModuleID:     moduleID,
					DependencyID: requirement.ID,
				})
				continue
			}
			if !compatible {
				diagnostics = append(diagnostics, ResolutionDiagnostic{
					Code:         "resolver.dependency.required_incompatible",
					ModuleID:     moduleID,
					DependencyID: requirement.ID,
				})
				continue
			}
			addRequiredEdge(edges, indegree, requirement.ID, moduleID)
		}

		for _, requirement := range snapshot.OptionalDependencies {
			status := OptionalDependencyStatus{
				ModuleID:     moduleID,
				DependencyID: requirement.ID,
				State:        optionalAvailable,
			}
			provider, exists := recordByID[requirement.ID]
			if !exists {
				status.State = optionalMissing
				status.Reason = "dependency_missing"
				optional = append(optional, status)
				continue
			}
			compatible, valid := matchesDependencyConstraint(provider.Version, requirement.Constraint)
			if !valid {
				status.State = optionalIncompatible
				if _, ok := parseStrictSemVer(provider.Version); !ok {
					status.Reason = "version_invalid"
				} else {
					status.Reason = "constraint_invalid"
				}
			} else if !compatible {
				status.State = optionalIncompatible
				status.Reason = "version_incompatible"
			}
			optional = append(optional, status)
		}
	}

	diagnostics = append(diagnostics, validateDependencyObservations(registry, observations)...)

	order, cyclic := deterministicRequiredOrder(nodes, edges, indegree)
	for _, moduleID := range cyclic {
		diagnostics = append(diagnostics, ResolutionDiagnostic{
			Code:     "resolver.graph.required_cycle",
			ModuleID: moduleID,
		})
	}

	sortOptionalDependencyStatus(optional)
	if len(diagnostics) > 0 {
		return DependencyResolution{}, newResolutionErrors(diagnostics)
	}
	return DependencyResolution{
		Order:    order,
		Optional: optional,
	}, nil
}

func validatePlatformSnapshot(platform PlatformSnapshot) (map[string]struct{}, []ResolutionDiagnostic) {
	capabilities := make(map[string]struct{}, len(platform.Capabilities))
	diagnostics := make([]ResolutionDiagnostic, 0)
	if len(platform.Capabilities) > MaxDeclarationItems {
		diagnostics = append(diagnostics, ResolutionDiagnostic{Code: "resolver.platform.capabilities_too_many"})
		return capabilities, diagnostics
	}
	for _, capability := range platform.Capabilities {
		if !validBoundedPattern(capability, platformIDPattern, maxIdentifierLength) {
			diagnostics = append(diagnostics, ResolutionDiagnostic{Code: "resolver.platform.capability_invalid"})
			continue
		}
		if _, exists := capabilities[capability]; exists {
			diagnostics = append(diagnostics, ResolutionDiagnostic{
				Code:               "resolver.platform.capability_duplicate",
				PlatformDependency: capability,
			})
			continue
		}
		capabilities[capability] = struct{}{}
	}
	return capabilities, diagnostics
}

func addRequiredEdge(edges map[string]map[string]struct{}, indegree map[string]int, dependencyID, moduleID string) {
	dependents, ok := edges[dependencyID]
	if !ok {
		dependents = make(map[string]struct{})
		edges[dependencyID] = dependents
	}
	if _, exists := dependents[moduleID]; exists {
		return
	}
	dependents[moduleID] = struct{}{}
	indegree[moduleID]++
}

func deterministicRequiredOrder(nodes []string, edges map[string]map[string]struct{}, originalIndegree map[string]int) ([]string, []string) {
	indegree := make(map[string]int, len(originalIndegree))
	for id, count := range originalIndegree {
		indegree[id] = count
	}
	ready := make([]string, 0)
	for _, id := range nodes {
		if indegree[id] == 0 {
			ready = append(ready, id)
		}
	}
	sort.Strings(ready)

	order := make([]string, 0, len(nodes))
	for len(ready) > 0 {
		id := ready[0]
		ready = ready[1:]
		order = append(order, id)

		dependents := make([]string, 0, len(edges[id]))
		for dependent := range edges[id] {
			dependents = append(dependents, dependent)
		}
		sort.Strings(dependents)
		for _, dependent := range dependents {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				ready = append(ready, dependent)
				sort.Strings(ready)
			}
		}
	}

	if len(order) == len(nodes) {
		return order, nil
	}
	cyclic := make([]string, 0)
	for _, id := range nodes {
		if indegree[id] > 0 {
			cyclic = append(cyclic, id)
		}
	}
	sort.Strings(cyclic)
	return order, cyclic
}

func validateDependencyObservations(registry Registry, observations []DependencyObservation) []ResolutionDiagnostic {
	diagnostics := make([]ResolutionDiagnostic, 0)
	for _, observation := range observations {
		providerValid := validBoundedPattern(observation.ProviderID, moduleIDPattern, maxIdentifierLength)
		if validBoundedPattern(observation.ConsumerID, platformIDPattern, maxIdentifierLength) {
			if providerValid {
				diagnostics = append(diagnostics, ResolutionDiagnostic{
					Code:         "resolver.dependency.forbidden",
					ModuleID:     observation.ConsumerID,
					DependencyID: observation.ProviderID,
				})
			} else {
				diagnostics = append(diagnostics, ResolutionDiagnostic{Code: "resolver.observation.invalid"})
			}
			continue
		}
		if !validBoundedPattern(observation.ConsumerID, moduleIDPattern, maxIdentifierLength) || !providerValid {
			diagnostics = append(diagnostics, ResolutionDiagnostic{Code: "resolver.observation.invalid"})
			continue
		}

		snapshot, ok := registry.manifestSnapshot(observation.ConsumerID)
		if !ok {
			diagnostics = append(diagnostics, ResolutionDiagnostic{
				Code:     "resolver.observation.consumer_missing",
				ModuleID: observation.ConsumerID,
			})
			continue
		}
		if observation.Private {
			diagnostics = append(diagnostics, ResolutionDiagnostic{
				Code:         "resolver.dependency.private_forbidden",
				ModuleID:     observation.ConsumerID,
				DependencyID: observation.ProviderID,
			})
			continue
		}
		if !snapshot.declaresModuleDependency(observation.ProviderID) {
			diagnostics = append(diagnostics, ResolutionDiagnostic{
				Code:         "resolver.dependency.undeclared",
				ModuleID:     observation.ConsumerID,
				DependencyID: observation.ProviderID,
			})
		}
	}
	return diagnostics
}

func sortOptionalDependencyStatus(values []OptionalDependencyStatus) {
	sort.SliceStable(values, func(i, j int) bool {
		if values[i].ModuleID != values[j].ModuleID {
			return values[i].ModuleID < values[j].ModuleID
		}
		if values[i].DependencyID != values[j].DependencyID {
			return values[i].DependencyID < values[j].DependencyID
		}
		if values[i].State != values[j].State {
			return values[i].State < values[j].State
		}
		return values[i].Reason < values[j].Reason
	})
}

func newResolutionErrors(diagnostics []ResolutionDiagnostic) error {
	sort.SliceStable(diagnostics, func(i, j int) bool {
		if diagnostics[i].ModuleID != diagnostics[j].ModuleID {
			return diagnostics[i].ModuleID < diagnostics[j].ModuleID
		}
		if diagnostics[i].DependencyID != diagnostics[j].DependencyID {
			return diagnostics[i].DependencyID < diagnostics[j].DependencyID
		}
		if diagnostics[i].PlatformDependency != diagnostics[j].PlatformDependency {
			return diagnostics[i].PlatformDependency < diagnostics[j].PlatformDependency
		}
		return diagnostics[i].Code < diagnostics[j].Code
	})
	return &ResolutionErrors{diagnostics: append([]ResolutionDiagnostic(nil), diagnostics...)}
}

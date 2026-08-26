package modules

import (
	"errors"
	"fmt"
	"sort"
)

// DiscoverySource is an explicit, caller-approved source of manifest bytes.
// P03.02 discovery never scans filesystems or networks; callers must supply the
// exact bounded source set to inspect.
type DiscoverySource struct {
	ID        string
	Version   string
	Manifests [][]byte
}

// RegistryRecord is discovered/available metadata only. It intentionally
// contains no installed/enabled lifecycle state and grants no authority.
type RegistryRecord struct {
	ID            string `json:"id"`
	Version       string `json:"version"`
	Owner         string `json:"owner"`
	SourceID      string `json:"source_id"`
	SourceVersion string `json:"source_version"`
}

// DiscoveryDiagnostic is a stable operator-safe discovery result. It never
// includes filesystem paths, raw manifest content, secret references or parser
// text.
type DiscoveryDiagnostic struct {
	Code     string `json:"code"`
	SourceID string `json:"source_id,omitempty"`
	ModuleID string `json:"module_id,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// DiscoveryErrors contains deterministic fail-closed diagnostics.
type DiscoveryErrors struct {
	diagnostics []DiscoveryDiagnostic
}

func (e *DiscoveryErrors) Error() string {
	return "module discovery failed"
}

// Diagnostics returns a copy so callers cannot mutate discovery evidence.
func (e *DiscoveryErrors) Diagnostics() []DiscoveryDiagnostic {
	if e == nil {
		return nil
	}
	result := make([]DiscoveryDiagnostic, len(e.diagnostics))
	copy(result, e.diagnostics)
	return result
}

// Registry is an immutable-by-convention deterministic snapshot. Construction
// is available only through Discover.
type Registry struct {
	records []RegistryRecord
	byID    map[string]RegistryRecord
}

// List returns a stable ID/version/source ordered copy.
func (r Registry) List() []RegistryRecord {
	result := make([]RegistryRecord, len(r.records))
	copy(result, r.records)
	return result
}

// Lookup returns discovered metadata for exactly one module ID.
func (r Registry) Lookup(id string) (RegistryRecord, bool) {
	if r.byID == nil {
		return RegistryRecord{}, false
	}
	record, ok := r.byID[id]
	return record, ok
}

// Discover builds one deterministic registry from the exact explicit source set.
// Every manifest must pass the existing P03.01 ParseManifest validation contract.
// Any malformed source, malformed manifest, duplicate module identity or
// conflicting module version fails the entire discovery operation.
func Discover(sources []DiscoverySource) (Registry, error) {
	candidates := make([]RegistryRecord, 0)
	diagnostics := make([]DiscoveryDiagnostic, 0)

	for _, source := range sources {
		safeSourceID := discoverySafeSourceID(source.ID)
		if !validBoundedPattern(source.ID, identifierPattern, maxIdentifierLength) || !semverPattern.MatchString(source.Version) || len(source.Version) > maxIdentifierLength {
			diagnostics = append(diagnostics, DiscoveryDiagnostic{
				Code:     "discovery.source.invalid",
				SourceID: safeSourceID,
			})
			continue
		}

		for _, payload := range source.Manifests {
			manifest, err := ParseManifest(payload)
			if err != nil {
				diagnostics = append(diagnostics, DiscoveryDiagnostic{
					Code:     "discovery.manifest.invalid",
					SourceID: source.ID,
					Reason:   discoveryReason(err),
				})
				continue
			}
			candidates = append(candidates, RegistryRecord{
				ID:            manifest.ID,
				Version:       manifest.Version,
				Owner:         manifest.Owner,
				SourceID:      source.ID,
				SourceVersion: source.Version,
			})
		}
	}

	if len(diagnostics) > 0 {
		return Registry{}, newDiscoveryErrors(diagnostics)
	}

	sortRegistryRecords(candidates)
	for start := 0; start < len(candidates); {
		end := start + 1
		for end < len(candidates) && candidates[end].ID == candidates[start].ID {
			end++
		}
		if end-start > 1 {
			code := "discovery.module.duplicate"
			version := candidates[start].Version
			for index := start + 1; index < end; index++ {
				if candidates[index].Version != version {
					code = "discovery.module.version_conflict"
					break
				}
			}
			for index := start; index < end; index++ {
				diagnostics = append(diagnostics, DiscoveryDiagnostic{
					Code:     code,
					SourceID: candidates[index].SourceID,
					ModuleID: candidates[index].ID,
				})
			}
		}
		start = end
	}

	if len(diagnostics) > 0 {
		return Registry{}, newDiscoveryErrors(diagnostics)
	}

	byID := make(map[string]RegistryRecord, len(candidates))
	for _, record := range candidates {
		byID[record.ID] = record
	}
	return Registry{records: candidates, byID: byID}, nil
}

func sortRegistryRecords(records []RegistryRecord) {
	sort.SliceStable(records, func(i, j int) bool {
		if records[i].ID != records[j].ID {
			return records[i].ID < records[j].ID
		}
		if records[i].Version != records[j].Version {
			return records[i].Version < records[j].Version
		}
		if records[i].SourceID != records[j].SourceID {
			return records[i].SourceID < records[j].SourceID
		}
		return records[i].SourceVersion < records[j].SourceVersion
	})
}

func newDiscoveryErrors(diagnostics []DiscoveryDiagnostic) error {
	sort.SliceStable(diagnostics, func(i, j int) bool {
		if diagnostics[i].ModuleID != diagnostics[j].ModuleID {
			return diagnostics[i].ModuleID < diagnostics[j].ModuleID
		}
		if diagnostics[i].SourceID != diagnostics[j].SourceID {
			return diagnostics[i].SourceID < diagnostics[j].SourceID
		}
		if diagnostics[i].Code != diagnostics[j].Code {
			return diagnostics[i].Code < diagnostics[j].Code
		}
		return diagnostics[i].Reason < diagnostics[j].Reason
	})
	return &DiscoveryErrors{diagnostics: append([]DiscoveryDiagnostic(nil), diagnostics...)}
}

func discoveryReason(err error) string {
	var parseErr *ParseError
	if errors.As(err, &parseErr) {
		return parseErr.Code
	}

	var validationErr *ValidationErrors
	if errors.As(err, &validationErr) {
		issues := validationErr.Issues()
		if len(issues) > 0 {
			return issues[0].Code
		}
	}

	return "manifest.invalid"
}

func discoverySafeSourceID(id string) string {
	if validBoundedPattern(id, identifierPattern, maxIdentifierLength) {
		return id
	}
	return ""
}

func (d DiscoveryDiagnostic) String() string {
	return fmt.Sprintf("%s:%s:%s", d.Code, d.SourceID, d.ModuleID)
}

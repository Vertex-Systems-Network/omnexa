package modules

import "sort"

// PackageTrustContractVersion versions the metadata-only P03.11 trust hook
// contract. It does not version or imply a trust-root/certification decision.
const PackageTrustContractVersion = 1

// PackageTrustAuthority makes the current decision boundary explicit.
type PackageTrustAuthority string

const (
	// PackageTrustMetadataOnly means the profile is declarative evidence only.
	// Trust roots, signature verification and certification belong to a future
	// separately-authorized XTRUST-100 boundary.
	PackageTrustMetadataOnly PackageTrustAuthority = "metadata_only"
)

// PublisherIdentityHook is optional publisher identity metadata. Presence does
// not prove publisher ownership or trust.
type PublisherIdentityHook struct {
	ContractVersion int    `json:"contract_version"`
	Identity        string `json:"identity"`
}

// PackageProvenanceHook is an optional bounded package signature/provenance
// reference. P03.11 records the reference but does not dereference or verify it.
type PackageProvenanceHook struct {
	ContractVersion int    `json:"contract_version"`
	Reference       string `json:"reference"`
}

// SBOMIdentityHook is an optional bounded SBOM identity/reference. P03.11 does
// not fetch, parse, attest or enforce the referenced SBOM.
type SBOMIdentityHook struct {
	ContractVersion int    `json:"contract_version"`
	Reference       string `json:"reference"`
}

// PackageDeclaredScopeProfile is the already-validated declared package scope.
// SecretReferences contains symbolic secret names only, never secret values or
// the underlying secret:// locator.
type PackageDeclaredScopeProfile struct {
	ContractVersion      int      `json:"contract_version"`
	CapabilitiesProvided []string `json:"capabilities_provided"`
	CapabilitiesConsumed []string `json:"capabilities_consumed"`
	Permissions          []string `json:"permissions"`
	DataClassification   []string `json:"data_classification"`
	SecretReferences     []string `json:"secret_references"`
	NetworkDestinations  []string `json:"network_destinations"`
	ExposedEndpoints     []string `json:"exposed_endpoints"`
	FileTypes            []string `json:"file_types"`
	PrivilegedOperations []string `json:"privileged_operations"`
	AICapabilities       []string `json:"ai_capabilities"`
}

// PackageTrustProfile is a deterministic, non-authoritative registry-bound view
// of package identity/provenance/SBOM/scope declarations. It intentionally has
// no trusted/certified/signature-valid field and grants no runtime authority.
type PackageTrustProfile struct {
	ContractVersion int                         `json:"contract_version"`
	ModuleID        string                      `json:"module_id"`
	ModuleVersion   string                      `json:"module_version"`
	Owner           string                      `json:"owner"`
	Authority       PackageTrustAuthority       `json:"authority"`
	Publisher       *PublisherIdentityHook      `json:"publisher,omitempty"`
	Provenance      *PackageProvenanceHook      `json:"provenance,omitempty"`
	SBOM            *SBOMIdentityHook           `json:"sbom,omitempty"`
	Scope           PackageDeclaredScopeProfile `json:"scope"`
}

// PackageTrustDiagnostic is a stable, value-bounded profile construction
// failure. It excludes raw manifests, package code, secret locators and parser
// text.
type PackageTrustDiagnostic struct {
	Code     string `json:"code"`
	ModuleID string `json:"module_id,omitempty"`
}

// PackageTrustError is the fail-closed P03.11 registry/profile error.
type PackageTrustError struct {
	Diagnostic PackageTrustDiagnostic
}

func (e *PackageTrustError) Error() string { return "package trust profile construction failed" }

// BuildPackageTrustProfiles returns a deterministic copy-only metadata view for
// every discovered module. It never loads package files, executes hooks, makes
// network calls or verifies trust roots/signatures.
func BuildPackageTrustProfiles(registry Registry) ([]PackageTrustProfile, error) {
	records := registry.List()
	profiles := make([]PackageTrustProfile, 0, len(records))
	for _, record := range records {
		profile, err := packageTrustProfile(registry, record)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	sort.Slice(profiles, func(i, j int) bool {
		if profiles[i].ModuleID != profiles[j].ModuleID {
			return profiles[i].ModuleID < profiles[j].ModuleID
		}
		return profiles[i].ModuleVersion < profiles[j].ModuleVersion
	})
	return profiles, nil
}

// PackageTrustProfileFor returns one registry-bound metadata-only profile.
func PackageTrustProfileFor(registry Registry, moduleID string) (PackageTrustProfile, bool, error) {
	record, ok := registry.Lookup(moduleID)
	if !ok {
		return PackageTrustProfile{}, false, nil
	}
	profile, err := packageTrustProfile(registry, record)
	if err != nil {
		return PackageTrustProfile{}, false, err
	}
	return profile, true, nil
}

func packageTrustProfile(registry Registry, record RegistryRecord) (PackageTrustProfile, error) {
	snapshot, ok := registry.manifestSnapshot(record.ID)
	if !ok {
		return PackageTrustProfile{}, packageTrustError("package.trust.registry.snapshot_missing", record.ID)
	}
	if snapshot.ID != record.ID || snapshot.Version != record.Version || snapshot.Owner != record.Owner {
		return PackageTrustProfile{}, packageTrustError("package.trust.registry.snapshot_mismatch", record.ID)
	}

	profile := PackageTrustProfile{
		ContractVersion: PackageTrustContractVersion,
		ModuleID:        record.ID,
		ModuleVersion:   record.Version,
		Owner:           record.Owner,
		Authority:       PackageTrustMetadataOnly,
		Scope: PackageDeclaredScopeProfile{
			ContractVersion:      PackageTrustContractVersion,
			CapabilitiesProvided: sortedStringCopy(snapshot.CapabilitiesProvided),
			CapabilitiesConsumed: sortedStringCopy(snapshot.CapabilitiesConsumed),
			Permissions:          sortedStringCopy(snapshot.Permissions),
			DataClassification:   sortedStringCopy(snapshot.DataClassification),
			SecretReferences:     sortedSecretReferenceNames(snapshot.Security.SecretReferences),
			NetworkDestinations:  sortedStringCopy(snapshot.Security.NetworkDestinations),
			ExposedEndpoints:     sortedStringCopy(snapshot.Security.ExposedEndpoints),
			FileTypes:            sortedStringCopy(snapshot.Security.FileTypes),
			PrivilegedOperations: sortedStringCopy(snapshot.Security.PrivilegedOperations),
			AICapabilities:       sortedStringCopy(snapshot.Security.AICapabilities),
		},
	}
	if snapshot.Publisher != "" {
		profile.Publisher = &PublisherIdentityHook{ContractVersion: PackageTrustContractVersion, Identity: snapshot.Publisher}
	}
	if snapshot.ProvenanceRef != "" {
		profile.Provenance = &PackageProvenanceHook{ContractVersion: PackageTrustContractVersion, Reference: snapshot.ProvenanceRef}
	}
	if snapshot.SBOMRef != "" {
		profile.SBOM = &SBOMIdentityHook{ContractVersion: PackageTrustContractVersion, Reference: snapshot.SBOMRef}
	}
	return profile, nil
}

func sortedStringCopy(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}

func sortedSecretReferenceNames(values []SecretReference) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.Name)
	}
	sort.Strings(result)
	return result
}

func packageTrustError(code, moduleID string) error {
	return &PackageTrustError{Diagnostic: PackageTrustDiagnostic{Code: code, ModuleID: moduleID}}
}

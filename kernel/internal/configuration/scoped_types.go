package configuration

import (
	"context"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/authorization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

const (
	maxScopedCorrelationRunes = 128
	maxScopedReasonRunes      = 512
)

const (
	PermissionSettingRead   authorization.PermissionID = "configuration.setting.read"
	PermissionSettingManage authorization.PermissionID = "configuration.setting.manage"
)

// DataClassification is the subset of the frozen classification vocabulary that
// the generic runtime setting surface may return. RESTRICTED values remain outside
// this generic registry because P02.09 is not a secrets-management product.
type DataClassification string

const (
	DataPublic       DataClassification = "PUBLIC"
	DataInternal     DataClassification = "INTERNAL"
	DataConfidential DataClassification = "CONFIDENTIAL"
)

func (classification DataClassification) valid() bool {
	switch classification {
	case DataPublic, DataInternal, DataConfidential:
		return true
	default:
		return false
	}
}

// SettingPolicy binds one existing P01.10 definition to P02.09 scoped behavior.
// It describes handling and scope only; it never grants authorization.
type SettingPolicy struct {
	Key                       Key
	Classification            DataClassification
	AllowOrganizationOverride bool
	ProtectedRead             bool
	SecuritySignificant       bool
}

func validateSettingPolicy(registry *Registry, policy SettingPolicy) error {
	if registry == nil || !policy.Classification.valid() {
		return scopedPolicyInvalidFailure()
	}
	if _, ok := registry.Definition(policy.Key); !ok {
		return scopedPolicyInvalidFailure()
	}
	if policy.Classification != DataPublic && !policy.ProtectedRead {
		return scopedPolicyInvalidFailure()
	}
	if policy.SecuritySignificant && !policy.ProtectedRead {
		return scopedPolicyInvalidFailure()
	}
	return nil
}

// ScopedKind is the only P02.09 setting-scope vocabulary.
type ScopedKind string

const (
	ScopedTenant       ScopedKind = "tenant"
	ScopedOrganization ScopedKind = "organization"
)

// TrustedSettingScope can only be derived from accepted P02 trusted contexts.
// Its private fields prevent client-supplied tenant/org identifiers from becoming
// setting authority.
type TrustedSettingScope struct {
	kind           ScopedKind
	tenantID       tenancy.TenantID
	organizationID organization.NodeID
	subject        authorization.Subject
}

// ScopeFromTenantContext derives exact tenant setting scope from current trusted tenancy.
func ScopeFromTenantContext(trusted tenancy.TrustedContext) (TrustedSettingScope, error) {
	subject, err := authorization.SubjectFromTenantContext(trusted)
	if err != nil {
		return TrustedSettingScope{}, scopedContextInvalidFailure()
	}
	scope := TrustedSettingScope{kind: ScopedTenant, tenantID: trusted.TenantID(), subject: subject}
	if !scope.valid() {
		return TrustedSettingScope{}, scopedContextInvalidFailure()
	}
	return scope, nil
}

// ScopeFromOrganizationContext derives exact organization setting scope from a
// current trusted organization relationship. It does not create authorization.
func ScopeFromOrganizationContext(scoped organization.ScopedContext) (TrustedSettingScope, error) {
	subject, err := authorization.SubjectFromOrganizationContext(scoped)
	if err != nil {
		return TrustedSettingScope{}, scopedContextInvalidFailure()
	}
	scope := TrustedSettingScope{
		kind:           ScopedOrganization,
		tenantID:       scoped.TenantID(),
		organizationID: scoped.NodeID(),
		subject:        subject,
	}
	if !scope.valid() {
		return TrustedSettingScope{}, scopedContextInvalidFailure()
	}
	return scope, nil
}

func (scope TrustedSettingScope) valid() bool {
	if !scope.tenantID.Valid() || !scope.subject.Valid() {
		return false
	}
	if scope.subject.Scope().TenantID() != scope.tenantID {
		return false
	}
	switch scope.kind {
	case ScopedTenant:
		return scope.organizationID == "" && scope.subject.Scope().Kind() == authorization.ScopeTenant
	case ScopedOrganization:
		return scope.organizationID.Valid() &&
			scope.subject.Scope().Kind() == authorization.ScopeOrganization &&
			scope.subject.Scope().OrganizationID() == scope.organizationID
	default:
		return false
	}
}

func (scope TrustedSettingScope) evaluationContext() EvaluationContext {
	result := EvaluationContext{TenantID: string(scope.tenantID)}
	if scope.kind == ScopedOrganization {
		result.OrganizationID = string(scope.organizationID)
	}
	return result
}

// Kind returns tenant or organization without exposing a constructor.
func (scope TrustedSettingScope) Kind() ScopedKind { return scope.kind }

// TenantID returns the trusted tenant isolation boundary.
func (scope TrustedSettingScope) TenantID() tenancy.TenantID { return scope.tenantID }

// OrganizationID returns the exact organization or the zero value for tenant scope.
func (scope TrustedSettingScope) OrganizationID() organization.NodeID { return scope.organizationID }

// ScopedEvaluation is a classification-aware P02.09 projection. RESTRICTED
// values cannot enter this surface because such policies are rejected at construction.
type ScopedEvaluation struct {
	Classification DataClassification
	Evaluation     Evaluation
}

// SettingMutation is value-free mutation metadata safe for ordinary return paths.
type SettingMutation struct {
	Key            Key
	Scope          ScopedKind
	Classification DataClassification
	Revision       uint64
}

// SettingMutationMetadata is bounded required metadata for protected writes.
type SettingMutationMetadata struct {
	CorrelationID string
	Reason        string
}

func (metadata SettingMutationMetadata) validate() error {
	if !validScopedReference(metadata.CorrelationID, maxScopedCorrelationRunes) ||
		strings.TrimSpace(metadata.Reason) == "" || !validScopedText(metadata.Reason, maxScopedReasonRunes) {
		return scopedMutationInvalidFailure()
	}
	return nil
}

func validScopedReference(value string, maxRunes int) bool {
	if strings.TrimSpace(value) == "" || utf8.RuneCountInString(value) > maxRunes {
		return false
	}
	for _, char := range value {
		if unicode.IsSpace(char) || unicode.IsControl(char) {
			return false
		}
	}
	return true
}

func validScopedText(value string, maxRunes int) bool {
	if utf8.RuneCountInString(value) > maxRunes {
		return false
	}
	for _, char := range value {
		if unicode.IsControl(char) && char != '\t' {
			return false
		}
	}
	return true
}

// settingAuthorizer is intentionally narrow; *authorization.Service satisfies it.
type settingAuthorizer interface {
	Require(context.Context, authorization.Subject, authorization.PermissionID) error
}

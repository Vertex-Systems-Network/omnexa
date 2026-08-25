package configuration

import (
	"context"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/audit"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/authorization"
)

const scopedActorKindUser = "user"

// ScopedService is the P02.09 trusted tenant/organization setting boundary.
// Existing P01.10 evaluation remains non-authoritative beneath this service.
type ScopedService struct {
	registry   *Registry
	repository scopedSettingRepository
	authorizer settingAuthorizer
	audit      *audit.Writer
	policies   map[Key]SettingPolicy
	evaluator  *Evaluator
	now        func() time.Time
}

func NewScopedService(
	registry *Registry,
	repository scopedSettingRepository,
	authorizer settingAuthorizer,
	auditWriter *audit.Writer,
	policies []SettingPolicy,
	options EvaluatorOptions,
) (*ScopedService, error) {
	return newScopedServiceWithClock(registry, repository, authorizer, auditWriter, policies, options, func() time.Time { return time.Now().UTC() })
}

func newScopedServiceWithClock(
	registry *Registry,
	repository scopedSettingRepository,
	authorizer settingAuthorizer,
	auditWriter *audit.Writer,
	policies []SettingPolicy,
	options EvaluatorOptions,
	now func() time.Time,
) (*ScopedService, error) {
	if registry == nil || repository == nil || authorizer == nil || auditWriter == nil || now == nil || len(policies) == 0 {
		return nil, scopedServiceInvalidFailure()
	}
	policyMap := make(map[Key]SettingPolicy, len(policies))
	for _, policy := range policies {
		if err := validateSettingPolicy(registry, policy); err != nil {
			return nil, err
		}
		if _, exists := policyMap[policy.Key]; exists {
			return nil, scopedPolicyInvalidFailure()
		}
		policyMap[policy.Key] = policy
	}
	evaluator, err := NewEvaluator(registry, scopedProvider{repository: repository, policies: policyMap}, options)
	if err != nil {
		return nil, err
	}
	return &ScopedService{
		registry:   registry,
		repository: repository,
		authorizer: authorizer,
		audit:      auditWriter,
		policies:   policyMap,
		evaluator:  evaluator,
		now:        now,
	}, nil
}

// Resolve evaluates one setting only from a trusted tenant/org context.
func (service *ScopedService) Resolve(ctx context.Context, scope TrustedSettingScope, key Key) (ScopedEvaluation, error) {
	policy, err := service.policyAndScope(scope, key)
	if err != nil {
		return ScopedEvaluation{}, err
	}
	if policy.ProtectedRead {
		if err = service.authorizer.Require(ctx, scope.subject, PermissionSettingRead); err != nil {
			return ScopedEvaluation{}, err
		}
	}
	evaluation, err := service.evaluator.Evaluate(ctx, key, scope.evaluationContext())
	if err != nil {
		return ScopedEvaluation{}, err
	}
	return ScopedEvaluation{Classification: policy.Classification, Evaluation: evaluation}, nil
}

// Upsert writes one exact tenant/org override after server-side authorization.
// The returned mutation and audit record contain no setting value.
func (service *ScopedService) Upsert(
	ctx context.Context,
	scope TrustedSettingScope,
	key Key,
	value Value,
	metadata SettingMutationMetadata,
) (SettingMutation, error) {
	policy, err := service.policyAndScope(scope, key)
	if err != nil {
		return SettingMutation{}, err
	}
	if err = metadata.validate(); err != nil {
		return SettingMutation{}, err
	}
	definition, _ := service.registry.Definition(key)
	if !value.valid() || value.Kind() != definition.Kind {
		return SettingMutation{}, scopedValueInvalidFailure()
	}
	if scope.kind == ScopedOrganization && !policy.AllowOrganizationOverride {
		return SettingMutation{}, scopedPolicyInvalidFailure()
	}
	if err = service.authorizer.Require(ctx, scope.subject, PermissionSettingManage); err != nil {
		return SettingMutation{}, err
	}

	revision, err := service.repository.upsert(ctx, key, value, scope.tenantID, scope.organizationID, service.now())
	if err != nil {
		return SettingMutation{}, err
	}
	if _, err = service.evaluator.InvalidateKey(key); err != nil {
		return SettingMutation{}, err
	}
	if err = service.auditMutation(ctx, scope, policy, key, metadata); err != nil {
		return SettingMutation{}, err
	}
	return SettingMutation{Key: key, Scope: scope.kind, Classification: policy.Classification, Revision: revision}, nil
}

func (service *ScopedService) policyAndScope(scope TrustedSettingScope, key Key) (SettingPolicy, error) {
	if service == nil || service.registry == nil || service.repository == nil || service.authorizer == nil || service.audit == nil || service.evaluator == nil || service.now == nil {
		return SettingPolicy{}, scopedServiceInvalidFailure()
	}
	if !scope.valid() {
		return SettingPolicy{}, scopedContextInvalidFailure()
	}
	policy, ok := service.policies[key]
	if !ok {
		return SettingPolicy{}, scopedPolicyInvalidFailure()
	}
	if _, ok = service.registry.Definition(key); !ok {
		return SettingPolicy{}, scopedPolicyInvalidFailure()
	}
	return policy, nil
}

func (service *ScopedService) auditMutation(
	ctx context.Context,
	scope TrustedSettingScope,
	policy SettingPolicy,
	key Key,
	metadata SettingMutationMetadata,
) error {
	recordScope := audit.Scope{TenantID: string(scope.tenantID)}
	if scope.kind == ScopedOrganization {
		recordScope.OrganizationID = string(scope.organizationID)
	}
	_, err := service.audit.Write(ctx, audit.RequirementRequired, audit.RecordInput{
		Classification: auditClassification(policy.Classification),
		Actor:          audit.Actor{Kind: scopedActorKindUser, Reference: string(scope.subject.PrincipalID())},
		Action:         "configuration.setting.upsert",
		Target:         audit.Target{Kind: "configuration.setting", Reference: string(key)},
		Scope:          recordScope,
		Outcome:        audit.OutcomeSucceeded,
		CorrelationID:  metadata.CorrelationID,
		Reason:         metadata.Reason,
		Privileged:     true,
	})
	return err
}

func auditClassification(classification DataClassification) audit.Classification {
	switch classification {
	case DataPublic:
		return audit.ClassificationPublic
	case DataInternal:
		return audit.ClassificationInternal
	case DataConfidential:
		return audit.ClassificationConfidential
	default:
		return audit.ClassificationInternal
	}
}

var _ settingAuthorizer = (*authorization.Service)(nil)

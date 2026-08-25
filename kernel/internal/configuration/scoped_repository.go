package configuration

import (
	"context"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

type scopedSettingRepository interface {
	resolveExact(context.Context, Key, tenancy.TenantID, organization.NodeID) (ProviderResult, error)
	upsert(context.Context, Key, Value, tenancy.TenantID, organization.NodeID, time.Time) (uint64, error)
}

type scopedProvider struct {
	repository scopedSettingRepository
	policies   map[Key]SettingPolicy
}

func (provider scopedProvider) Resolve(ctx context.Context, key Key, scope EvaluationContext) (ProviderResult, error) {
	if provider.repository == nil || scope.UserID != "" {
		return ProviderResult{}, scopedContextInvalidFailure()
	}
	policy, ok := provider.policies[key]
	if !ok {
		return ProviderResult{}, ErrProviderValueNotFound
	}
	if err := scope.validate(); err != nil || scope.TenantID == "" {
		return ProviderResult{}, scopedContextInvalidFailure()
	}
	tenantID := tenancy.TenantID(scope.TenantID)
	if !tenantID.Valid() {
		return ProviderResult{}, scopedContextInvalidFailure()
	}

	if scope.OrganizationID != "" {
		organizationID := organization.NodeID(scope.OrganizationID)
		if !organizationID.Valid() {
			return ProviderResult{}, scopedContextInvalidFailure()
		}
		if policy.AllowOrganizationOverride {
			result, err := provider.repository.resolveExact(ctx, key, tenantID, organizationID)
			if err == nil {
				return result, nil
			}
			if !errors.Is(err, ErrProviderValueNotFound) {
				return ProviderResult{}, err
			}
		}
	}
	return provider.repository.resolveExact(ctx, key, tenantID, "")
}

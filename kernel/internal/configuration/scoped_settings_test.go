package configuration

import (
	"context"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
)

func TestScopedPolicyRejectsRestrictedAndUnprotectedSensitiveValues(t *testing.T) {
	t.Parallel()
	registry := mustRegistry(t, scopedStringDefinition())

	cases := []SettingPolicy{
		{Key: "platform.ui.locale", Classification: DataClassification("RESTRICTED"), ProtectedRead: true},
		{Key: "platform.ui.locale", Classification: DataConfidential, ProtectedRead: false},
		{Key: "platform.ui.locale", Classification: DataInternal, SecuritySignificant: true, ProtectedRead: false},
	}
	for _, policy := range cases {
		if err := validateSettingPolicy(registry, policy); err == nil {
			t.Fatalf("validateSettingPolicy(%+v) unexpectedly succeeded", policy)
		}
	}
}

func TestScopedProviderUsesExactOrganizationThenTenantPrecedence(t *testing.T) {
	t.Parallel()
	key := Key("platform.ui.locale")
	tenantID := tenancy.TenantID("018f47a6-7b7e-7c14-a847-0af56b2a44fe")
	organizationID := organization.NodeID("018f47b1-f3a8-79da-8d9d-4a64d65c2fd5")
	repository := newFakeScopedRepository()
	repository.seed(key, tenantID, "", StringValue("tenant"), 1)
	repository.seed(key, tenantID, organizationID, StringValue("organization"), 2)
	provider := scopedProvider{
		repository: repository,
		policies: map[Key]SettingPolicy{
			key: {Key: key, Classification: DataPublic, AllowOrganizationOverride: true},
		},
	}

	result, err := provider.Resolve(context.Background(), key, EvaluationContext{
		TenantID: string(tenantID), OrganizationID: string(organizationID),
	})
	if err != nil {
		t.Fatalf("Resolve(org) error = %v", err)
	}
	got, _ := result.Value.String()
	if got != "organization" || result.Revision != 2 {
		t.Fatalf("organization result = %q rev=%d", got, result.Revision)
	}

	otherOrganization := organization.NodeID("018f47b1-f3a8-79da-8d9d-4a64d65c2fd6")
	result, err = provider.Resolve(context.Background(), key, EvaluationContext{
		TenantID: string(tenantID), OrganizationID: string(otherOrganization),
	})
	if err != nil {
		t.Fatalf("Resolve(tenant fallback) error = %v", err)
	}
	got, _ = result.Value.String()
	if got != "tenant" || result.Revision != 1 {
		t.Fatalf("tenant fallback = %q rev=%d", got, result.Revision)
	}
}

func TestScopedProviderNeverUsesOrganizationOverrideWhenPolicyDisallowsIt(t *testing.T) {
	t.Parallel()
	key := Key("platform.security.login_notice")
	tenantID := tenancy.TenantID("018f47a6-7b7e-7c14-a847-0af56b2a44fe")
	organizationID := organization.NodeID("018f47b1-f3a8-79da-8d9d-4a64d65c2fd5")
	repository := newFakeScopedRepository()
	repository.seed(key, tenantID, "", StringValue("tenant-safe"), 4)
	repository.seed(key, tenantID, organizationID, StringValue("should-not-apply"), 9)
	provider := scopedProvider{
		repository: repository,
		policies: map[Key]SettingPolicy{
			key: {Key: key, Classification: DataConfidential, ProtectedRead: true},
		},
	}
	result, err := provider.Resolve(context.Background(), key, EvaluationContext{
		TenantID: string(tenantID), OrganizationID: string(organizationID),
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	got, _ := result.Value.String()
	if got != "tenant-safe" || result.Revision != 4 {
		t.Fatalf("disallowed org override affected result = %q rev=%d", got, result.Revision)
	}
}

func TestScopedProviderRejectsUserOrMissingTenantScope(t *testing.T) {
	t.Parallel()
	key := Key("platform.ui.locale")
	provider := scopedProvider{
		repository: newFakeScopedRepository(),
		policies:   map[Key]SettingPolicy{key: {Key: key, Classification: DataPublic}},
	}
	if _, err := provider.Resolve(context.Background(), key, EvaluationContext{UserID: "018f47c2-1111-7abc-8def-1234567890ab"}); err == nil {
		t.Fatal("user-scoped setting resolution unexpectedly succeeded")
	}
	if _, err := provider.Resolve(context.Background(), key, EvaluationContext{}); err == nil {
		t.Fatal("global setting resolution unexpectedly succeeded")
	}
}

func TestScopedValueRoundTripPreservesSupportedKinds(t *testing.T) {
	t.Parallel()
	values := []Value{
		BoolValue(true),
		StringValue("safe-value"),
		IntValue(-42),
		DurationValue(15 * time.Minute),
	}
	for _, original := range values {
		encoded, err := encodeScopedValue(original)
		if err != nil {
			t.Fatalf("encodeScopedValue(%v) error = %v", original.Kind(), err)
		}
		decoded, err := decodeScopedValue(original.Kind(), encoded)
		if err != nil {
			t.Fatalf("decodeScopedValue(%v) error = %v", original.Kind(), err)
		}
		if !original.equal(decoded) {
			t.Fatalf("round trip kind %s changed value", original.Kind())
		}
	}
}

func scopedStringDefinition() Definition {
	return Definition{
		Key: "platform.ui.locale", Description: "Synthetic scoped setting.", Owner: "kernel.configuration",
		Kind: KindString, Class: ClassRuntimeConfig, Version: 1,
		Default: StringValue("en"), Fallback: StringValue("en"),
	}
}

type fakeScopedRepository struct {
	values map[fakeScopedAddress]ProviderResult
}

type fakeScopedAddress struct {
	key            Key
	tenantID       tenancy.TenantID
	organizationID organization.NodeID
}

func newFakeScopedRepository() *fakeScopedRepository {
	return &fakeScopedRepository{values: make(map[fakeScopedAddress]ProviderResult)}
}

func (repository *fakeScopedRepository) seed(key Key, tenantID tenancy.TenantID, organizationID organization.NodeID, value Value, revision uint64) {
	repository.values[fakeScopedAddress{key: key, tenantID: tenantID, organizationID: organizationID}] = ProviderResult{Value: value, Revision: revision}
}

func (repository *fakeScopedRepository) resolveExact(_ context.Context, key Key, tenantID tenancy.TenantID, organizationID organization.NodeID) (ProviderResult, error) {
	result, ok := repository.values[fakeScopedAddress{key: key, tenantID: tenantID, organizationID: organizationID}]
	if !ok {
		return ProviderResult{}, ErrProviderValueNotFound
	}
	return result, nil
}

func (repository *fakeScopedRepository) upsert(_ context.Context, key Key, value Value, tenantID tenancy.TenantID, organizationID organization.NodeID, _ time.Time) (uint64, error) {
	address := fakeScopedAddress{key: key, tenantID: tenantID, organizationID: organizationID}
	revision := uint64(1)
	if current, ok := repository.values[address]; ok {
		revision = current.Revision + 1
	}
	repository.values[address] = ProviderResult{Value: value, Revision: revision}
	return revision, nil
}

var _ scopedSettingRepository = (*fakeScopedRepository)(nil)

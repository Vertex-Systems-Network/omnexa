package modules

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestCapabilityRegistryRetainsV1AndV2ValidatedDeclarations(t *testing.T) {
	v1 := validManifest()
	v1.CapabilitiesProvided = []string{"inventory.reserve-stock.v1"}
	v1.CapabilitiesConsumed = []string{"catalog.product.read.v2"}

	v1Snapshot, err := parseValidatedManifest(manifestBytes(t, v1))
	if err != nil {
		t.Fatalf("parse v1 snapshot: %v", err)
	}
	if !reflect.DeepEqual(v1Snapshot.CapabilitiesProvided, v1.CapabilitiesProvided) ||
		!reflect.DeepEqual(v1Snapshot.CapabilitiesConsumed, v1.CapabilitiesConsumed) {
		t.Fatalf("v1 capability declarations missing from validated snapshot: %#v", v1Snapshot)
	}

	v2 := validManifestV2()
	v2.CapabilitiesProvided = []string{"inventory.reserve-stock.v3"}
	v2.CapabilitiesConsumed = []string{"catalog.product.read.v4"}
	v2Snapshot, err := parseValidatedManifest(manifestV2Bytes(t, v2))
	if err != nil {
		t.Fatalf("parse v2 snapshot: %v", err)
	}
	if !reflect.DeepEqual(v2Snapshot.CapabilitiesProvided, v2.CapabilitiesProvided) ||
		!reflect.DeepEqual(v2Snapshot.CapabilitiesConsumed, v2.CapabilitiesConsumed) {
		t.Fatalf("v2 capability declarations missing from validated snapshot: %#v", v2Snapshot)
	}

	clone := v2Snapshot.clone()
	clone.CapabilitiesProvided[0] = "mutated.value.v1"
	if v2Snapshot.CapabilitiesProvided[0] != "inventory.reserve-stock.v3" {
		t.Fatal("validated snapshot clone leaked capability slice mutation")
	}
}

func TestCapabilityRegistryDeterministicProviderConsumerLookup(t *testing.T) {
	ctx := context.Background()
	moduleRegistry := capabilityModuleRegistry(t, "inventory.reserve-stock.v1", "inventory.reserve-stock.v1")
	store := enabledCapabilityLifecycleStore(t, "omnexa.inventory")

	registry, err := BindCapabilityRegistry(moduleRegistry, store, []CapabilityRegistration{
		validCapabilityRegistration("inventory.reserve-stock.v1"),
	})
	if err != nil {
		t.Fatalf("BindCapabilityRegistry() error = %v", err)
	}

	providers, err := registry.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(providers) != 1 {
		t.Fatalf("expected one provider, got %#v", providers)
	}
	provider := providers[0]
	if provider.ID != "inventory.reserve-stock" ||
		provider.Major != 1 ||
		provider.ProviderModuleID != "omnexa.inventory" ||
		provider.ProviderOwner != "inventory.team" ||
		!provider.Available ||
		provider.LifecycleState != LifecycleEnabled {
		t.Fatalf("unexpected provider record: %#v", provider)
	}
	if provider.AuthorizationRef != "authorization.inventory.reserve-stock.v1" ||
		provider.ScopeRef != "scope.tenant-organization" ||
		provider.ContractRef != "contract.inventory.reserve-stock.v1" {
		t.Fatalf("descriptive metadata refs were not retained: %#v", provider)
	}

	consumers := registry.Consumers()
	if len(consumers) != 1 ||
		consumers[0].ConsumerModuleID != "omnexa.checkout" ||
		consumers[0].ID != "inventory.reserve-stock" ||
		consumers[0].Major != 1 {
		t.Fatalf("unexpected consumers: %#v", consumers)
	}

	lookup, err := registry.Lookup(ctx, CapabilityQuery{
		ID:    "inventory.reserve-stock",
		Major: 1,
		Owner: "inventory.team",
	})
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}
	if !reflect.DeepEqual(lookup, provider) {
		t.Fatalf("lookup mismatch:\nlist=%#v\nlookup=%#v", provider, lookup)
	}

	resolved, err := registry.ResolveConsumer(ctx, "omnexa.checkout", "inventory.reserve-stock.v1")
	if err != nil {
		t.Fatalf("ResolveConsumer() error = %v", err)
	}
	if !reflect.DeepEqual(resolved, provider) {
		t.Fatalf("consumer resolution mismatch:\nprovider=%#v\nresolved=%#v", provider, resolved)
	}

	copyConsumers := registry.Consumers()
	copyConsumers[0].ConsumerOwner = "mutated.team"
	if registry.Consumers()[0].ConsumerOwner != "checkout.team" {
		t.Fatal("consumer list mutation escaped immutable-by-copy boundary")
	}
}

func TestCapabilityRegistryRejectsUndeclaredOwnerDuplicateAndCollision(t *testing.T) {
	baseRegistry := capabilityModuleRegistry(t, "inventory.reserve-stock.v1", "")
	store := NewMemoryLifecycleStore()

	t.Run("undeclared registration", func(t *testing.T) {
		registration := validCapabilityRegistration("inventory.unknown.v1")
		_, err := BindCapabilityRegistry(baseRegistry, store, []CapabilityRegistration{registration})
		assertCapabilityErrorCode(t, err, "module.capability.registration_undeclared")
	})

	t.Run("owner mismatch", func(t *testing.T) {
		registration := validCapabilityRegistration("inventory.reserve-stock.v1")
		registration.Owner = "other.team"
		_, err := BindCapabilityRegistry(baseRegistry, store, []CapabilityRegistration{registration})
		assertCapabilityErrorCode(t, err, "module.capability.owner_mismatch")
	})

	t.Run("duplicate registration", func(t *testing.T) {
		registration := validCapabilityRegistration("inventory.reserve-stock.v1")
		_, err := BindCapabilityRegistry(baseRegistry, store, []CapabilityRegistration{registration, registration})
		assertCapabilityErrorCode(t, err, "module.capability.registration_duplicate")
	})

	t.Run("provider declaration collision", func(t *testing.T) {
		first := capabilityProviderManifest("omnexa.inventory", "inventory.team", "inventory.reserve-stock.v1")
		second := capabilityProviderManifest("omnexa.inventory-alt", "inventory.team", "inventory.reserve-stock.v1")
		registry, err := Discover([]DiscoverySource{
			{ID: "repo.alpha", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, first)}},
			{ID: "repo.beta", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, second)}},
		})
		if err != nil {
			t.Fatalf("Discover() error = %v", err)
		}
		_, err = BindCapabilityRegistry(registry, store, []CapabilityRegistration{
			validCapabilityRegistration("inventory.reserve-stock.v1"),
		})
		assertCapabilityErrorCode(t, err, "module.capability.declaration_collision")
	})

	t.Run("stable capability owner conflict across majors", func(t *testing.T) {
		first := capabilityProviderManifest("omnexa.inventory", "inventory.team", "inventory.reserve-stock.v1")
		second := capabilityProviderManifest("omnexa.fulfillment", "fulfillment.team", "inventory.reserve-stock.v2")
		registry, err := Discover([]DiscoverySource{
			{ID: "repo.alpha", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, first)}},
			{ID: "repo.beta", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, second)}},
		})
		if err != nil {
			t.Fatalf("Discover() error = %v", err)
		}
		_, err = BindCapabilityRegistry(registry, store, []CapabilityRegistration{
			validCapabilityRegistration("inventory.reserve-stock.v1"),
			{
				ModuleID:         "omnexa.fulfillment",
				Owner:            "fulfillment.team",
				Declaration:      "inventory.reserve-stock.v2",
				AuthorizationRef: "authorization.inventory.reserve-stock.v2",
				ScopeRef:         "scope.tenant-organization",
				ContractRef:      "contract.inventory.reserve-stock.v2",
			},
		})
		assertCapabilityErrorCode(t, err, "module.capability.owner_conflict")
	})
}

func TestCapabilityRegistryRejectsInvalidCapabilityContractAndMetadataRefs(t *testing.T) {
	store := NewMemoryLifecycleStore()

	t.Run("provided declaration without major", func(t *testing.T) {
		provider := capabilityProviderManifest("omnexa.inventory", "inventory.team", "inventory.reserve-stock")
		registry, err := Discover([]DiscoverySource{{
			ID: "repo.alpha", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, provider)},
		}})
		if err != nil {
			t.Fatalf("Discover() error = %v", err)
		}
		_, err = BindCapabilityRegistry(registry, store, nil)
		assertCapabilityErrorCode(t, err, "module.capability.declaration_invalid")
	})

	t.Run("metadata reference invalid", func(t *testing.T) {
		registry := capabilityModuleRegistry(t, "inventory.reserve-stock.v1", "")
		registration := validCapabilityRegistration("inventory.reserve-stock.v1")
		registration.AuthorizationRef = "https://not-an-authority-grant.example.test"
		_, err := BindCapabilityRegistry(registry, store, []CapabilityRegistration{registration})
		assertCapabilityErrorCode(t, err, "module.capability.metadata_ref_invalid")
	})
}

func TestCapabilityRegistryMajorCompatibilityFailsClosed(t *testing.T) {
	ctx := context.Background()
	moduleRegistry := capabilityModuleRegistry(t, "inventory.reserve-stock.v1", "inventory.reserve-stock.v2")
	store := enabledCapabilityLifecycleStore(t, "omnexa.inventory")
	registry, err := BindCapabilityRegistry(moduleRegistry, store, []CapabilityRegistration{
		validCapabilityRegistration("inventory.reserve-stock.v1"),
	})
	if err != nil {
		t.Fatalf("BindCapabilityRegistry() error = %v", err)
	}

	_, err = registry.ResolveConsumer(ctx, "omnexa.checkout", "inventory.reserve-stock.v2")
	assertCapabilityErrorCode(t, err, "module.capability.version_incompatible")

	_, err = registry.ResolveConsumer(ctx, "omnexa.checkout", "inventory.reserve-stock.v1")
	assertCapabilityErrorCode(t, err, "module.capability.consumer_undeclared")
}

func TestCapabilityRegistryLifecycleAvailabilityIsEnabledOnlyAndNonDestructive(t *testing.T) {
	ctx := context.Background()
	states := []LifecycleState{
		LifecycleAvailable,
		LifecycleInstalled,
		LifecycleDisabled,
		LifecycleSuspended,
		LifecycleArchived,
		LifecycleDetached,
		LifecycleRecoveryRequired,
		LifecyclePurged,
	}

	for _, state := range states {
		t.Run(string(state), func(t *testing.T) {
			moduleRegistry := capabilityModuleRegistry(t, "inventory.reserve-stock.v1", "inventory.reserve-stock.v1")
			store := NewMemoryLifecycleStore()
			if state != LifecycleAvailable {
				if err := store.CompareAndSwap(ctx, "omnexa.inventory", 0, LifecycleRecord{
					ModuleID: "omnexa.inventory",
					Version:  "1.2.3",
					State:    state,
					Revision: 1,
				}); err != nil {
					t.Fatalf("seed lifecycle: %v", err)
				}
			}
			registry, err := BindCapabilityRegistry(moduleRegistry, store, []CapabilityRegistration{
				validCapabilityRegistration("inventory.reserve-stock.v1"),
			})
			if err != nil {
				t.Fatalf("BindCapabilityRegistry() error = %v", err)
			}

			record, err := registry.Lookup(ctx, CapabilityQuery{
				ID: "inventory.reserve-stock", Major: 1, Owner: "inventory.team",
			})
			if err != nil {
				t.Fatalf("Lookup() error = %v", err)
			}
			if record.Available {
				t.Fatalf("lifecycle %s must not report active capability: %#v", state, record)
			}
			if record.ID != "inventory.reserve-stock" || record.Declaration != "inventory.reserve-stock.v1" {
				t.Fatalf("historical capability identity was lost for lifecycle %s: %#v", state, record)
			}

			_, err = registry.ResolveConsumer(ctx, "omnexa.checkout", "inventory.reserve-stock.v1")
			assertCapabilityErrorCode(t, err, "module.capability.provider_unavailable")
		})
	}

	t.Run("enabled", func(t *testing.T) {
		moduleRegistry := capabilityModuleRegistry(t, "inventory.reserve-stock.v1", "inventory.reserve-stock.v1")
		store := enabledCapabilityLifecycleStore(t, "omnexa.inventory")
		registry, err := BindCapabilityRegistry(moduleRegistry, store, []CapabilityRegistration{
			validCapabilityRegistration("inventory.reserve-stock.v1"),
		})
		if err != nil {
			t.Fatalf("BindCapabilityRegistry() error = %v", err)
		}
		record, err := registry.ResolveConsumer(ctx, "omnexa.checkout", "inventory.reserve-stock.v1")
		if err != nil {
			t.Fatalf("ResolveConsumer() error = %v", err)
		}
		if !record.Available || record.LifecycleState != LifecycleEnabled {
			t.Fatalf("enabled provider must be available: %#v", record)
		}
	})
}

func TestCapabilityRegistryMetadataSurfaceHasNoInvocationOrRawScopeAuthority(t *testing.T) {
	for _, value := range []any{CapabilityRegistration{}, CapabilityRecord{}, CapabilityConsumer{}} {
		typ := reflect.TypeOf(value)
		for _, forbidden := range []string{
			"Handler",
			"Invoke",
			"Permission",
			"TenantID",
			"OrganizationID",
			"DatabaseTable",
			"Secret",
		} {
			if _, ok := typ.FieldByName(forbidden); ok {
				t.Fatalf("%s must not expose forbidden authority field %s", typ.Name(), forbidden)
			}
		}
	}
}

func capabilityModuleRegistry(t *testing.T, provided, consumed string) Registry {
	t.Helper()
	provider := capabilityProviderManifest("omnexa.inventory", "inventory.team", provided)

	consumer := validManifest()
	consumer.ID = "omnexa.checkout"
	consumer.Name = "Checkout"
	consumer.Owner = "checkout.team"
	consumer.Dependencies = nil
	consumer.OptionalDependencies = nil
	consumer.CapabilitiesProvided = nil
	if consumed == "" {
		consumer.CapabilitiesConsumed = nil
	} else {
		consumer.CapabilitiesConsumed = []string{consumed}
	}

	registry, err := Discover([]DiscoverySource{
		{ID: "repo.checkout", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, consumer)}},
		{ID: "repo.inventory", Version: "1.0.0", Manifests: [][]byte{manifestBytes(t, provider)}},
	})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	return registry
}

func capabilityProviderManifest(moduleID, owner, provided string) Manifest {
	manifest := validManifest()
	manifest.ID = moduleID
	manifest.Owner = owner
	manifest.Dependencies = nil
	manifest.OptionalDependencies = nil
	if provided == "" {
		manifest.CapabilitiesProvided = nil
	} else {
		manifest.CapabilitiesProvided = []string{provided}
	}
	manifest.CapabilitiesConsumed = nil
	return manifest
}

func validCapabilityRegistration(declaration string) CapabilityRegistration {
	return CapabilityRegistration{
		ModuleID:         "omnexa.inventory",
		Owner:            "inventory.team",
		Declaration:      declaration,
		AuthorizationRef: "authorization.inventory.reserve-stock.v1",
		ScopeRef:         "scope.tenant-organization",
		ContractRef:      "contract.inventory.reserve-stock.v1",
	}
}

func enabledCapabilityLifecycleStore(t *testing.T, moduleID string) *MemoryLifecycleStore {
	t.Helper()
	store := NewMemoryLifecycleStore()
	if err := store.CompareAndSwap(context.Background(), moduleID, 0, LifecycleRecord{
		ModuleID: moduleID,
		Version:  "1.2.3",
		State:    LifecycleEnabled,
		Revision: 1,
	}); err != nil {
		t.Fatalf("seed lifecycle: %v", err)
	}
	return store
}

func assertCapabilityErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	var capabilityErr *CapabilityError
	if !errors.As(err, &capabilityErr) {
		t.Fatalf("expected CapabilityError %s, got %T: %v", code, err, err)
	}
	if capabilityErr.Diagnostic.Code != code {
		t.Fatalf("unexpected capability error: got %#v want code %s", capabilityErr.Diagnostic, code)
	}
}

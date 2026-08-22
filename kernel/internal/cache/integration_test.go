package cache

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestUnavailableConnectionIsBoundedAndSafe(t *testing.T) {
	const secretLikeAddress = "127.0.0.1:1"
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			cacheAddressKey:        secretLikeAddress,
			cacheConnectTimeoutKey: "100ms",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}

	started := time.Now()
	_, err = NewStore(context.Background(), resolved)
	if err == nil {
		t.Fatal("NewStore() error = nil, want unavailable connection failure")
	}
	if time.Since(started) > 2*time.Second {
		t.Fatalf("unavailable connection was not bounded: %v", time.Since(started))
	}
	assertFailureCode(t, err, codeConnectionUnavailable)
	if strings.Contains(err.Error(), secretLikeAddress) {
		t.Fatalf("public failure leaked provider address: %q", err)
	}
}

func TestValkeyFoundationIntegration(t *testing.T) {
	address := os.Getenv("P01_05_TEST_CACHE_ADDRESS")
	if address == "" {
		t.Skip("P01_05_TEST_CACHE_ADDRESS is not set")
	}
	store := integrationStore(t, address, "p0105")
	defer store.Close()

	ctx := context.Background()
	if err := store.client.Do(ctx, store.client.B().Flushdb().Build()).Error(); err != nil {
		t.Fatalf("initial FLUSHDB error = %v", err)
	}

	t.Run("miss_is_distinct_from_provider_failure", func(t *testing.T) {
		entry, err := store.Get(ctx, Key{Namespace: "kernel.cache", Version: 1, Name: "missing"})
		if err != nil {
			t.Fatalf("Get(missing) error = %v", err)
		}
		if entry.Found || entry.Value != nil {
			t.Fatalf("missing entry = %#v, want explicit miss", entry)
		}
	})

	t.Run("set_get_and_expiry", func(t *testing.T) {
		key := Key{Namespace: "kernel.cache", Version: 1, Name: "expiry"}
		if err := store.Set(ctx, key, []byte("synthetic-value"), 150*time.Millisecond); err != nil {
			t.Fatalf("Set() error = %v", err)
		}
		entry, err := store.Get(ctx, key)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if !entry.Found || string(entry.Value) != "synthetic-value" {
			t.Fatalf("entry = %#v", entry)
		}

		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			entry, err = store.Get(ctx, key)
			if err != nil {
				t.Fatalf("Get(after expiry) error = %v", err)
			}
			if !entry.Found {
				return
			}
			time.Sleep(25 * time.Millisecond)
		}
		t.Fatal("cache entry did not expire within bounded test window")
	})

	t.Run("delete_reports_existing_state", func(t *testing.T) {
		key := Key{Namespace: "kernel.cache", Version: 1, Name: "delete"}
		if err := store.Set(ctx, key, []byte("delete-me"), time.Minute); err != nil {
			t.Fatalf("Set() error = %v", err)
		}
		deleted, err := store.Delete(ctx, key)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		if !deleted {
			t.Fatal("Delete() = false, want true for existing key")
		}
		deleted, err = store.Delete(ctx, key)
		if err != nil {
			t.Fatalf("Delete(second) error = %v", err)
		}
		if deleted {
			t.Fatal("Delete(second) = true, want false")
		}
	})

	t.Run("value_and_ttl_bounds_fail_closed", func(t *testing.T) {
		key := Key{Namespace: "kernel.cache", Version: 1, Name: "bounds"}
		oversized := make([]byte, store.settings.MaxValueBytes+1)
		if err := store.Set(ctx, key, oversized, time.Minute); err == nil {
			t.Fatal("Set(oversized) error = nil")
		} else {
			assertFailureCode(t, err, codeValueInvalid)
		}
		if err := store.Set(ctx, key, []byte("value"), 0); err == nil {
			t.Fatal("Set(zero ttl) error = nil")
		} else {
			assertFailureCode(t, err, codeValueInvalid)
		}
		if err := store.Set(ctx, key, []byte("value"), store.settings.MaxTTL+time.Millisecond); err == nil {
			t.Fatal("Set(long ttl) error = nil")
		} else {
			assertFailureCode(t, err, codeValueInvalid)
		}
	})

	t.Run("caller_cancellation_is_not_a_miss", func(t *testing.T) {
		canceled, cancel := context.WithCancel(ctx)
		cancel()
		_, err := store.Get(canceled, Key{Namespace: "kernel.cache", Version: 1, Name: "cancel"})
		if err == nil {
			t.Fatal("Get(canceled) error = nil")
		}
		assertFailureCode(t, err, codeOperationFailed)
	})

	t.Run("flush_proves_cache_is_non_authoritative", func(t *testing.T) {
		key := Key{Namespace: "kernel.cache", Version: 1, Name: "flush"}
		if err := store.Set(ctx, key, []byte("disposable"), time.Minute); err != nil {
			t.Fatalf("Set() error = %v", err)
		}
		if err := store.client.Do(ctx, store.client.B().Flushdb().Build()).Error(); err != nil {
			t.Fatalf("FLUSHDB error = %v", err)
		}
		entry, err := store.Get(ctx, key)
		if err != nil {
			t.Fatalf("Get(after flush) error = %v", err)
		}
		if entry.Found {
			t.Fatalf("entry survived provider flush: %#v", entry)
		}
	})
}

func TestValkeyReconnectAfterProviderRestart(t *testing.T) {
	if os.Getenv("P01_05_AFTER_RESTART") != "1" {
		t.Skip("provider restart phase is not active")
	}
	address := os.Getenv("P01_05_TEST_CACHE_ADDRESS")
	if address == "" {
		t.Fatal("P01_05_TEST_CACHE_ADDRESS is required after restart")
	}

	store := integrationStore(t, address, "p0105restart")
	defer store.Close()
	key := Key{Namespace: "kernel.cache", Version: 1, Name: "restart"}
	if err := store.Set(context.Background(), key, []byte("after-restart"), time.Minute); err != nil {
		t.Fatalf("Set(after restart) error = %v", err)
	}
	entry, err := store.Get(context.Background(), key)
	if err != nil {
		t.Fatalf("Get(after restart) error = %v", err)
	}
	if !entry.Found || string(entry.Value) != "after-restart" {
		t.Fatalf("entry after restart = %#v", entry)
	}
}

func integrationStore(t *testing.T, address, prefix string) *Store {
	t.Helper()
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			cacheAddressKey:          address,
			cacheConnectTimeoutKey:   "1s",
			cacheOperationTimeoutKey: "1s",
			cacheKeyPrefixKey:        prefix,
			cacheMaxValueBytesKey:    "4096",
			cacheMaxTTLKey:           "1h",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}
	store, err := NewStore(context.Background(), resolved)
	if err != nil {
		if structured, ok := failure.As(err); ok {
			t.Fatalf("NewStore() failure code=%s category=%s", structured.Code(), structured.Category())
		}
		t.Fatalf("NewStore() error = %v", err)
	}
	return store
}

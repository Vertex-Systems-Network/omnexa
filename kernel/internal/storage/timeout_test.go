package storage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestStorageConnectionVerificationTimeoutIsBoundedAndDistinct(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resolved := mustStorageConfiguration(t, map[string]string{
		storageEndpointKey:       server.URL,
		storageConnectTimeoutKey: "50ms",
	})
	started := time.Now()
	_, err := NewStore(context.Background(), resolved)
	if err == nil {
		t.Fatal("NewStore(timeout) error = nil")
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("NewStore(timeout) took %v, want bounded failure", elapsed)
	}
	assertStorageFailureCode(t, err, codeConnectionTimeout)
	assertPublicStorageFailureDoesNotContain(t, err, server.URL)
}

func TestStorageOperationTimeoutAndCallerCancellationAreDistinct(t *testing.T) {
	const bucketPath = "/omnexa-test-bucket"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodHead && strings.TrimSuffix(request.URL.Path, "/") == bucketPath {
			writer.WriteHeader(http.StatusOK)
			return
		}
		time.Sleep(200 * time.Millisecond)
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resolved := mustStorageConfiguration(t, map[string]string{
		storageEndpointKey:         server.URL,
		storageConnectTimeoutKey:   "500ms",
		storageOperationTimeoutKey: "50ms",
	})
	store, err := NewStore(context.Background(), resolved)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	key := Key{Namespace: "kernel.storage", Version: 1, Path: "timeout/object.bin"}
	started := time.Now()
	_, err = store.Head(context.Background(), key)
	if err == nil {
		t.Fatal("Head(timeout) error = nil")
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("Head(timeout) took %v, want bounded failure", elapsed)
	}
	assertStorageFailureCode(t, err, codeOperationTimeout)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = store.Head(ctx, key)
	if err == nil {
		t.Fatal("Head(canceled) error = nil")
	}
	assertStorageFailureCode(t, err, codeOperationCanceled)
}

package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	testStorageEndpointEnv = "P01_06_TEST_S3_ENDPOINT"
	testStorageBucketEnv   = "P01_06_TEST_S3_BUCKET"
)

func TestUnavailableStorageConnectionIsBoundedAndSafe(t *testing.T) {
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			storageEndpointKey:       "http://127.0.0.1:1",
			storageAccessKeyKey:      "synthetic-unavailable-access",
			storageSecretKeyKey:      "synthetic-unavailable-secret",
			storageBucketKey:         "omnexa-unavailable-bucket",
			storageConnectTimeoutKey: "150ms",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}

	started := time.Now()
	_, err = NewStore(context.Background(), resolved)
	if err == nil {
		t.Fatal("NewStore() error = nil, want unavailable provider failure")
	}
	if elapsed := time.Since(started); elapsed > 2*time.Second {
		t.Fatalf("unavailable provider took %v, want bounded failure", elapsed)
	}
	assertStorageFailureCode(t, err, codeConnectionUnavailable)
	structured, ok := failure.As(err)
	if !ok {
		t.Fatalf("failure.As() = false for %v", err)
	}
	public := structured.Public()
	joined := strings.Join([]string{public.Title, public.Detail, err.Error()}, " ")
	for _, secret := range []string{"127.0.0.1:1", "synthetic-unavailable-access", "synthetic-unavailable-secret"} {
		if strings.Contains(joined, secret) {
			t.Fatalf("public connection failure leaked %q", secret)
		}
	}
}

func TestS3CompatibleStorageFoundationIntegration(t *testing.T) {
	store := liveStorageStore(t)
	defer store.Close()

	t.Run("put_head_open_delete_and_missing", func(t *testing.T) {
		key := uniqueStorageKey(t, "contract/object.txt")
		payload := []byte("synthetic P01.06 object storage contract")
		info, err := store.Put(context.Background(), key, Upload{
			Body:          bytes.NewReader(payload),
			ContentLength: int64(len(payload)),
			ContentType:   "text/plain; charset=utf-8",
			FileName:      "object.txt",
			SHA256:        SHA256Hex(payload),
			Metadata:      map[string]string{"fixture-kind": "contract"},
		})
		if err != nil {
			t.Fatalf("Put() error = %v", err)
		}
		if info.ContentLength != int64(len(payload)) || info.SHA256 != SHA256Hex(payload) {
			t.Fatalf("Put() info = %#v", info)
		}

		head, err := store.Head(context.Background(), key)
		if err != nil {
			t.Fatalf("Head() error = %v", err)
		}
		if head.ContentLength != int64(len(payload)) || head.FileName != "object.txt" || head.Metadata["fixture-kind"] != "contract" {
			t.Fatalf("Head() info = %#v", head)
		}

		object, err := store.Open(context.Background(), key)
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		decoded, err := io.ReadAll(object.Body)
		closeErr := object.Body.Close()
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		if closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
		if !bytes.Equal(decoded, payload) {
			t.Fatalf("download = %q, want %q", decoded, payload)
		}

		if deleteErr := store.Delete(context.Background(), key); deleteErr != nil {
			t.Fatalf("Delete() error = %v", deleteErr)
		}
		if deleteErr := store.Delete(context.Background(), key); deleteErr != nil {
			t.Fatalf("Delete(missing) error = %v, want idempotent success", deleteErr)
		}
		_, err = store.Head(context.Background(), key)
		if err == nil {
			t.Fatal("Head(deleted) error = nil, want not found")
		}
		assertStorageFailureCode(t, err, codeObjectNotFound)
	})

	t.Run("large_object_streams_with_bounded_working_buffers", func(t *testing.T) {
		const size = int64(8 * 1024 * 1024)
		key := uniqueStorageKey(t, "streaming/large.bin")
		expectedHash := patternSHA256(t, size)
		uploadReader := &patternReader{remaining: size}
		_, err := store.Put(context.Background(), key, Upload{
			Body:          uploadReader,
			ContentLength: size,
			ContentType:   "application/octet-stream",
			FileName:      "large.bin",
			SHA256:        expectedHash,
		})
		if err != nil {
			t.Fatalf("Put(large) error = %v", err)
		}
		defer func() { _ = store.Delete(context.Background(), key) }()
		if uploadReader.maxRequest > 128*1024 {
			t.Fatalf("upload requested %d bytes in one read; want bounded streaming", uploadReader.maxRequest)
		}

		object, err := store.Open(context.Background(), key)
		if err != nil {
			t.Fatalf("Open(large) error = %v", err)
		}
		defer func() { _ = object.Body.Close() }()
		hasher := sha256.New()
		buffer := make([]byte, 32*1024)
		written, err := io.CopyBuffer(hasher, object.Body, buffer)
		if err != nil {
			t.Fatalf("streaming download error = %v", err)
		}
		if written != size {
			t.Fatalf("streamed bytes = %d, want %d", written, size)
		}
		if actual := hex.EncodeToString(hasher.Sum(nil)); actual != expectedHash {
			t.Fatalf("streamed SHA-256 = %s, want %s", actual, expectedHash)
		}
	})

	t.Run("caller_cancellation_is_not_timeout_or_missing_object", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := store.Head(ctx, uniqueStorageKey(t, "cancelled/object.bin"))
		if err == nil {
			t.Fatal("Head(cancelled) error = nil")
		}
		assertStorageFailureCode(t, err, codeOperationCanceled)
		structured, ok := failure.As(err)
		if !ok || structured.Category() != failure.CategoryTimeout {
			t.Fatalf("cancelled category = %v, want timeout category with cancellation-specific code", err)
		}
		if structured.Retryable() {
			t.Fatal("caller cancellation must not be marked retryable")
		}
	})

	t.Run("provider_side_tamper_is_detected_on_stream", func(t *testing.T) {
		key := uniqueStorageKey(t, "integrity/tampered.bin")
		rendered, err := RenderKey(store.settings.KeyPrefix, key)
		if err != nil {
			t.Fatalf("RenderKey() error = %v", err)
		}
		expected := []byte("expected-content")
		tampered := []byte("tampered-content")
		_, err = store.client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket:        aws.String(store.settings.Bucket),
			Key:           aws.String(rendered),
			Body:          bytes.NewReader(tampered),
			ContentLength: aws.Int64(int64(len(tampered))),
			ContentType:   aws.String("application/octet-stream"),
			Metadata: map[string]string{
				checksumMetadataKey: SHA256Hex(expected),
				fileNameMetadataKey: "tampered.bin",
			},
		})
		if err != nil {
			t.Fatalf("raw provider tamper fixture error = %v", err)
		}
		defer func() { _ = store.Delete(context.Background(), key) }()

		object, err := store.Open(context.Background(), key)
		if err != nil {
			t.Fatalf("Open(tampered) error = %v", err)
		}
		_, err = io.Copy(io.Discard, object.Body)
		_ = object.Body.Close()
		if err == nil {
			t.Fatal("tampered download error = nil")
		}
		assertStorageFailureCode(t, err, codeIntegrityFailed)
	})
}

func liveStorageStore(t *testing.T) *Store {
	t.Helper()
	endpoint := os.Getenv(testStorageEndpointEnv)
	bucket := os.Getenv(testStorageBucketEnv)
	if endpoint == "" || bucket == "" {
		t.Skipf("%s and %s are required for P01.06 integration evidence", testStorageEndpointEnv, testStorageBucketEnv)
	}
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			storageEndpointKey:         endpoint,
			storageAccessKeyKey:        "synthetic-p01-06-access",
			storageSecretKeyKey:        "synthetic-p01-06-secret",
			storageRegionKey:           "us-east-1",
			storageBucketKey:           bucket,
			storageUsePathStyleKey:     "true",
			storageConnectTimeoutKey:   "5s",
			storageOperationTimeoutKey: "30s",
			storageMaxObjectBytesKey:   "33554432",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}
	store, err := NewStore(context.Background(), resolved)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	return store
}

func uniqueStorageKey(t *testing.T, path string) Key {
	t.Helper()
	name := strings.ToLower(strings.ReplaceAll(t.Name(), "/", "-"))
	name = strings.ReplaceAll(name, "_", "-")
	return Key{
		Namespace: "kernel.storage",
		Version:   1,
		Path:      "integration/" + name + "/" + path,
	}
}

func patternSHA256(t *testing.T, size int64) string {
	t.Helper()
	hasher := sha256.New()
	reader := &patternReader{remaining: size}
	if _, err := io.Copy(hasher, reader); err != nil {
		t.Fatalf("pattern SHA-256 error = %v", err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

type patternReader struct {
	remaining  int64
	maxRequest int
}

func (reader *patternReader) Read(buffer []byte) (int, error) {
	if reader.remaining == 0 {
		return 0, io.EOF
	}
	if len(buffer) > reader.maxRequest {
		reader.maxRequest = len(buffer)
	}
	readBuffer := buffer
	if reader.remaining < int64(len(buffer)) {
		readBuffer = buffer[:reader.remaining]
	}
	clear(readBuffer)
	readCount := len(readBuffer)
	reader.remaining -= int64(readCount)
	return readCount, nil
}

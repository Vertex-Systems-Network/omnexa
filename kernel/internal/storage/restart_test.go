package storage

import (
	"bytes"
	"context"
	"os"
	"testing"
)

func TestS3CompatibleReconnectAfterProviderRestart(t *testing.T) {
	if os.Getenv("P01_06_AFTER_RESTART") != "1" {
		t.Skip("P01_06_AFTER_RESTART=1 is required for governed restart evidence")
	}
	store := liveStorageStore(t)
	defer store.Close()

	key := uniqueStorageKey(t, "restart/reconnected.txt")
	payload := []byte("synthetic provider restart recovery")
	_, err := store.Put(context.Background(), key, Upload{
		Body:          bytes.NewReader(payload),
		ContentLength: int64(len(payload)),
		ContentType:   "text/plain",
		FileName:      "reconnected.txt",
		SHA256:        SHA256Hex(payload),
	})
	if err != nil {
		t.Fatalf("Put(after restart) error = %v", err)
	}
	defer func() { _ = store.Delete(context.Background(), key) }()

	if _, err := store.Head(context.Background(), key); err != nil {
		t.Fatalf("Head(after restart) error = %v", err)
	}
}

package storage

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestS3CompatibleStoreSupportsConcurrentDistinctObjects(t *testing.T) {
	store := liveStorageStore(t)
	defer store.Close()

	const workers = 4
	errorsByWorker := make(chan error, workers)
	var group sync.WaitGroup
	group.Add(workers)
	for index := 0; index < workers; index++ {
		index := index
		go func() {
			defer group.Done()
			payload := []byte(fmt.Sprintf("synthetic-concurrent-object-%d", index))
			key := uniqueStorageKey(t, fmt.Sprintf("concurrent/object-%d.txt", index))
			if _, err := store.Put(context.Background(), key, Upload{
				Body:          bytes.NewReader(payload),
				ContentLength: int64(len(payload)),
				ContentType:   "text/plain",
				FileName:      fmt.Sprintf("object-%d.txt", index),
				SHA256:        SHA256Hex(payload),
			}); err != nil {
				errorsByWorker <- fmt.Errorf("worker %d put: %w", index, err)
				return
			}
			defer func() { _ = store.Delete(context.Background(), key) }()
			if _, err := store.Head(context.Background(), key); err != nil {
				errorsByWorker <- fmt.Errorf("worker %d head: %w", index, err)
			}
		}()
	}
	group.Wait()
	close(errorsByWorker)
	for err := range errorsByWorker {
		if err != nil {
			t.Fatal(err)
		}
	}
}

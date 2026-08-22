package storage

import (
	"context"
	"io"
	"sync"
)

type cancelReadCloser struct {
	body   io.ReadCloser
	cancel context.CancelFunc
	once   sync.Once
}

func (reader *cancelReadCloser) Read(buffer []byte) (int, error) {
	n, err := reader.body.Read(buffer)
	if err != nil {
		reader.release()
	}
	return n, err
}

func (reader *cancelReadCloser) Close() error {
	err := reader.body.Close()
	reader.release()
	return err
}

func (reader *cancelReadCloser) release() {
	reader.once.Do(reader.cancel)
}

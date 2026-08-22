package storage

import (
	"bytes"
	"io"
	"testing"
)

func TestIntegrityReaderValidatesZeroLengthObject(t *testing.T) {
	emptyHash := SHA256Hex(nil)
	reader := newIntegrityReader(bytes.NewReader(nil), 0, emptyHash)
	buffer := make([]byte, 1)
	n, err := reader.Read(buffer)
	if n != 0 || err != io.EOF {
		t.Fatalf("Read(empty) = (%d, %v), want (0, EOF)", n, err)
	}

	reader = newIntegrityReader(bytes.NewReader(nil), 0, SHA256Hex([]byte("wrong")))
	_, err = reader.Read(buffer)
	if err == nil {
		t.Fatal("Read(empty checksum mismatch) error = nil")
	}
	assertStorageFailureCode(t, err, codeIntegrityFailed)
}

func TestIntegrityReaderDetectsExtraByteAtDeclaredBoundary(t *testing.T) {
	payload := []byte("abcde")
	reader := newIntegrityReader(bytes.NewReader(append(append([]byte(nil), payload...), 'x')), int64(len(payload)), SHA256Hex(payload))
	buffer := make([]byte, len(payload))
	n, err := reader.Read(buffer)
	if n != len(payload) {
		t.Fatalf("Read(extra) n = %d, want %d", n, len(payload))
	}
	if err == nil || err == io.EOF {
		t.Fatalf("Read(extra) error = %v, want length failure", err)
	}
	assertStorageFailureCode(t, err, codeLengthInvalid)
}

func TestIntegrityReaderReportsChecksumFailureOnFinalChunk(t *testing.T) {
	payload := []byte("terminal-checksum")
	reader := newIntegrityReader(bytes.NewReader(payload), int64(len(payload)), SHA256Hex([]byte("different")))
	buffer := make([]byte, len(payload))
	n, err := reader.Read(buffer)
	if n != len(payload) {
		t.Fatalf("Read(checksum mismatch) n = %d, want %d", n, len(payload))
	}
	if err == nil || err == io.EOF {
		t.Fatalf("Read(checksum mismatch) error = %v, want integrity failure", err)
	}
	assertStorageFailureCode(t, err, codeIntegrityFailed)
}

func TestVerifiedReadCloserValidatesOnFinalChunk(t *testing.T) {
	payload := []byte("download-terminal")
	reader := newVerifiedReadCloser(io.NopCloser(bytes.NewReader(payload)), int64(len(payload)), SHA256Hex(payload))
	buffer := make([]byte, len(payload))
	n, err := reader.Read(buffer)
	if n != len(payload) || err != io.EOF {
		t.Fatalf("Read(valid download) = (%d, %v), want (%d, EOF)", n, err, len(payload))
	}

	reader = newVerifiedReadCloser(io.NopCloser(bytes.NewReader(payload)), int64(len(payload)), SHA256Hex([]byte("different")))
	n, err = reader.Read(buffer)
	if n != len(payload) {
		t.Fatalf("Read(corrupt download) n = %d, want %d", n, len(payload))
	}
	if err == nil || err == io.EOF {
		t.Fatalf("Read(corrupt download) error = %v, want integrity failure", err)
	}
	assertStorageFailureCode(t, err, codeIntegrityFailed)
}

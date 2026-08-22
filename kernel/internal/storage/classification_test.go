package storage

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestFailureClassificationDistinguishesCancellationTimeoutAndUnavailable(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		code      failure.Code
		category  failure.Category
		retryable bool
	}{
		{name: "connection canceled", err: classifyConnectionFailure(context.Canceled), code: codeConnectionCanceled, category: failure.CategoryTimeout},
		{name: "connection timeout", err: classifyConnectionFailure(context.DeadlineExceeded), code: codeConnectionTimeout, category: failure.CategoryTimeout, retryable: true},
		{name: "connection unavailable", err: classifyConnectionFailure(errors.New("provider unavailable")), code: codeConnectionUnavailable, category: failure.CategoryUnavailable, retryable: true},
		{name: "operation canceled", err: classifyOperationFailure(context.Canceled), code: codeOperationCanceled, category: failure.CategoryTimeout},
		{name: "operation timeout", err: classifyOperationFailure(context.DeadlineExceeded), code: codeOperationTimeout, category: failure.CategoryTimeout, retryable: true},
		{name: "operation unavailable", err: classifyOperationFailure(errors.New("provider unavailable")), code: codeOperationFailed, category: failure.CategoryUnavailable, retryable: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertStorageFailureCode(t, test.err, test.code)
			structured, ok := failure.As(test.err)
			if !ok {
				t.Fatalf("failure.As() = false for %v", test.err)
			}
			if structured.Category() != test.category {
				t.Fatalf("category = %q, want %q", structured.Category(), test.category)
			}
			if structured.Retryable() != test.retryable {
				t.Fatalf("retryable = %v, want %v", structured.Retryable(), test.retryable)
			}
		})
	}
}

func TestObjectInfoRejectsProviderSideBoundsAndMetadataViolations(t *testing.T) {
	contentType := "application/octet-stream"
	checksum := SHA256Hex([]byte("synthetic"))

	tooLarge := int64(1025)
	_, err := objectInfo("omnexa/kernel.storage/v1/too-large.bin", 1024, &tooLarge, &contentType, nil, nil, map[string]string{
		checksumMetadataKey: checksum,
	})
	if err == nil {
		t.Fatal("objectInfo(oversized provider object) error = nil")
	}
	assertStorageFailureCode(t, err, codeLengthInvalid)

	validLength := int64(1)
	_, err = objectInfo("omnexa/kernel.storage/v1/bad-metadata.bin", 1024, &validLength, &contentType, nil, nil, map[string]string{
		checksumMetadataKey: checksum,
		"bad key":           "value",
	})
	if err == nil {
		t.Fatal("objectInfo(invalid provider metadata) error = nil")
	}
	assertStorageFailureCode(t, err, codeIntegrityFailed)

	badContentType := "text/plain\nunsafe"
	_, err = objectInfo("omnexa/kernel.storage/v1/bad-type.bin", 1024, &validLength, &badContentType, nil, nil, map[string]string{
		checksumMetadataKey: checksum,
	})
	if err == nil {
		t.Fatal("objectInfo(invalid provider content type) error = nil")
	}
	assertStorageFailureCode(t, err, codeIntegrityFailed)
}

func TestDownloadStreamProviderErrorsRemainPublicSafe(t *testing.T) {
	providerDetail := "provider-internal-detail-" + t.Name()
	reader := newVerifiedReadCloser(&failingReadCloser{readErr: errors.New(providerDetail)}, 1, SHA256Hex([]byte("x")))
	_, err := reader.Read(make([]byte, 1))
	if err == nil {
		t.Fatal("Read(provider failure) error = nil")
	}
	assertStorageFailureCode(t, err, codeOperationFailed)
	assertPublicStorageFailureDoesNotContain(t, err, providerDetail)

	reader = newVerifiedReadCloser(&failingReadCloser{closeErr: errors.New(providerDetail)}, 1, SHA256Hex([]byte("x")))
	err = reader.Close()
	if err == nil {
		t.Fatal("Close(provider failure) error = nil")
	}
	assertStorageFailureCode(t, err, codeOperationFailed)
	assertPublicStorageFailureDoesNotContain(t, err, providerDetail)
}

func assertPublicStorageFailureDoesNotContain(t *testing.T, err error, providerDetail string) {
	t.Helper()
	structured, ok := failure.As(err)
	if !ok {
		t.Fatalf("failure.As() = false for %v", err)
	}
	public := structured.Public()
	joined := strings.Join([]string{string(public.Code), string(public.Category), public.Title, public.Detail, err.Error()}, " ")
	if strings.Contains(joined, providerDetail) {
		t.Fatalf("public storage failure leaked provider detail: %q", joined)
	}
}

type failingReadCloser struct {
	readErr  error
	closeErr error
}

func (reader *failingReadCloser) Read([]byte) (int, error) {
	if reader.readErr == nil {
		return 0, io.EOF
	}
	return 0, reader.readErr
}

func (reader *failingReadCloser) Close() error {
	return reader.closeErr
}

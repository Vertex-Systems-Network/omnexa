package storage

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestLoadConfigurationRedactsSensitiveStorageValues(t *testing.T) {
	const (
		endpoint  = "https://storage.internal.example"
		accessKey = "restricted-access-key"
		secretKey = "restricted-secret-key"
	)
	resolved, err := LoadConfiguration(config.Options{
		Environ: map[string]string{
			"OMNEXA_STORAGE_ENDPOINT":   endpoint,
			"OMNEXA_STORAGE_ACCESS_KEY": accessKey,
			"OMNEXA_STORAGE_SECRET_KEY": secretKey,
			"OMNEXA_STORAGE_BUCKET":     "omnexa-test-bucket",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}

	redacted := resolved.Redacted()
	for _, key := range []string{storageEndpointKey, storageAccessKeyKey, storageSecretKeyKey} {
		if redacted[key] != "<redacted>" {
			t.Fatalf("%s diagnostic = %q, want redacted", key, redacted[key])
		}
	}
	joined := strings.Join(storageMapValues(redacted), " ")
	for _, secret := range []string{endpoint, accessKey, secretKey} {
		if strings.Contains(joined, secret) {
			t.Fatalf("redacted configuration leaked sensitive provider value %q", secret)
		}
	}
}

func TestSettingsFromConfigAppliesBoundedDefaults(t *testing.T) {
	resolved := mustStorageConfiguration(t, map[string]string{})
	settings, err := SettingsFromConfig(resolved)
	if err != nil {
		t.Fatalf("SettingsFromConfig() error = %v", err)
	}
	if settings.Region != "us-east-1" {
		t.Fatalf("Region = %q, want us-east-1", settings.Region)
	}
	if !settings.UsePathStyle {
		t.Fatal("UsePathStyle = false, want true")
	}
	if settings.ConnectTimeout != 3*time.Second {
		t.Fatalf("ConnectTimeout = %v, want 3s", settings.ConnectTimeout)
	}
	if settings.OperationTimeout != 5*time.Second {
		t.Fatalf("OperationTimeout = %v, want 5s", settings.OperationTimeout)
	}
	if settings.KeyPrefix != "omnexa" {
		t.Fatalf("KeyPrefix = %q, want omnexa", settings.KeyPrefix)
	}
	if settings.MaxObjectBytes != 1073741824 {
		t.Fatalf("MaxObjectBytes = %d, want 1073741824", settings.MaxObjectBytes)
	}
}

func TestSettingsFromConfigRejectsUnsafeBounds(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
	}{
		{name: "endpoint scheme", overrides: map[string]string{storageEndpointKey: "ftp://storage.example"}},
		{name: "endpoint credentials", overrides: map[string]string{storageEndpointKey: "https://user:secret@storage.example"}},
		{name: "bucket traversal", overrides: map[string]string{storageBucketKey: "bad..bucket"}},
		{name: "region", overrides: map[string]string{storageRegionKey: "Bad Region"}},
		{name: "connect timeout", overrides: map[string]string{storageConnectTimeoutKey: "31s"}},
		{name: "operation timeout", overrides: map[string]string{storageOperationTimeoutKey: "121s"}},
		{name: "key prefix", overrides: map[string]string{storageKeyPrefixKey: "Bad Prefix"}},
		{name: "zero max bytes", overrides: map[string]string{storageMaxObjectBytesKey: "0"}},
		{name: "oversize max bytes", overrides: map[string]string{storageMaxObjectBytesKey: "5368709121"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolved := mustStorageConfiguration(t, test.overrides)
			_, err := SettingsFromConfig(resolved)
			if err == nil {
				t.Fatal("SettingsFromConfig() error = nil, want validation failure")
			}
			assertStorageFailureCode(t, err, codeConfigurationInvalid)
		})
	}
}

func TestRenderKeyIsDeterministicVersionedAndScoped(t *testing.T) {
	key := Key{Namespace: "kernel.storage", Version: 2, Path: "synthetic/report.txt"}
	first, err := RenderKey("omnexa", key)
	if err != nil {
		t.Fatalf("RenderKey() error = %v", err)
	}
	second, err := RenderKey("omnexa", key)
	if err != nil {
		t.Fatalf("RenderKey() second error = %v", err)
	}
	if first != "omnexa/kernel.storage/v2/synthetic/report.txt" || second != first {
		t.Fatalf("rendered keys = %q / %q", first, second)
	}

	otherVersion, _ := RenderKey("omnexa", Key{Namespace: "kernel.storage", Version: 3, Path: "synthetic/report.txt"})
	otherNamespace, _ := RenderKey("omnexa", Key{Namespace: "kernel.other", Version: 2, Path: "synthetic/report.txt"})
	otherPath, _ := RenderKey("omnexa", Key{Namespace: "kernel.storage", Version: 2, Path: "synthetic/other.txt"})
	for label, value := range map[string]string{"version": otherVersion, "namespace": otherNamespace, "path": otherPath} {
		if value == first {
			t.Fatalf("%s change collided with base object key", label)
		}
	}
}

func TestRenderKeyRejectsTraversalAndAmbiguousPaths(t *testing.T) {
	paths := []string{
		"../secret",
		"./secret",
		"folder/../secret",
		"folder/./secret",
		"/absolute",
		"trailing/",
		"double//slash",
		`windows\path`,
		"folder/bad segment",
	}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			_, err := RenderKey("omnexa", Key{Namespace: "kernel.storage", Version: 1, Path: path})
			if err == nil {
				t.Fatal("RenderKey() error = nil, want validation failure")
			}
			assertStorageFailureCode(t, err, codeKeyInvalid)
		})
	}
}

func TestValidateUploadAcceptsBoundedMetadataAndReservesIntegrityKeys(t *testing.T) {
	payload := []byte("synthetic-object")
	metadata, err := validateUpload(Upload{
		Body:          bytes.NewReader(payload),
		ContentLength: int64(len(payload)),
		ContentType:   "text/plain; charset=utf-8",
		FileName:      "report.txt",
		SHA256:        SHA256Hex(payload),
		Metadata:      map[string]string{"fixture-kind": "unit"},
	}, 1024)
	if err != nil {
		t.Fatalf("validateUpload() error = %v", err)
	}
	if metadata[checksumMetadataKey] != SHA256Hex(payload) {
		t.Fatalf("checksum metadata = %q", metadata[checksumMetadataKey])
	}
	if metadata[fileNameMetadataKey] != "report.txt" || metadata["fixture-kind"] != "unit" {
		t.Fatalf("metadata = %#v", metadata)
	}

	_, err = validateUpload(Upload{
		Body:          bytes.NewReader(payload),
		ContentLength: int64(len(payload)),
		SHA256:        SHA256Hex(payload),
		Metadata:      map[string]string{checksumMetadataKey: "override"},
	}, 1024)
	if err == nil {
		t.Fatal("validateUpload(reserved metadata) error = nil")
	}
	assertStorageFailureCode(t, err, codeMetadataInvalid)
}

func TestValidateUploadRejectsUntrustedMetadataAndLength(t *testing.T) {
	payload := []byte("synthetic")
	tests := []struct {
		name   string
		upload Upload
		code   failure.Code
	}{
		{name: "missing body", upload: Upload{ContentLength: 1, SHA256: SHA256Hex(payload)}, code: codeMetadataInvalid},
		{name: "negative length", upload: Upload{Body: bytes.NewReader(payload), ContentLength: -1, SHA256: SHA256Hex(payload)}, code: codeLengthInvalid},
		{name: "large length", upload: Upload{Body: bytes.NewReader(payload), ContentLength: 2048, SHA256: SHA256Hex(payload)}, code: codeLengthInvalid},
		{name: "bad checksum", upload: Upload{Body: bytes.NewReader(payload), ContentLength: int64(len(payload)), SHA256: "not-a-checksum"}, code: codeMetadataInvalid},
		{name: "bad type", upload: Upload{Body: bytes.NewReader(payload), ContentLength: int64(len(payload)), SHA256: SHA256Hex(payload), ContentType: "text/plain\nmalicious"}, code: codeMetadataInvalid},
		{name: "path filename", upload: Upload{Body: bytes.NewReader(payload), ContentLength: int64(len(payload)), SHA256: SHA256Hex(payload), FileName: "../report.txt"}, code: codeMetadataInvalid},
		{name: "newline filename", upload: Upload{Body: bytes.NewReader(payload), ContentLength: int64(len(payload)), SHA256: SHA256Hex(payload), FileName: "report\n.txt"}, code: codeMetadataInvalid},
		{name: "newline metadata", upload: Upload{Body: bytes.NewReader(payload), ContentLength: int64(len(payload)), SHA256: SHA256Hex(payload), Metadata: map[string]string{"fixture": "bad\nvalue"}}, code: codeMetadataInvalid},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := validateUpload(test.upload, 1024)
			if err == nil {
				t.Fatal("validateUpload() error = nil, want failure")
			}
			assertStorageFailureCode(t, err, test.code)
		})
	}
}

func TestIntegrityReaderStreamsAndDetectsMismatch(t *testing.T) {
	payload := bytes.Repeat([]byte("0123456789abcdef"), 4096)
	tracking := &trackingReader{reader: bytes.NewReader(payload)}
	reader := newIntegrityReader(tracking, int64(len(payload)), SHA256Hex(payload))
	if _, err := io.Copy(io.Discard, reader); err != nil {
		t.Fatalf("io.Copy(valid) error = %v", err)
	}
	if tracking.maxRequest > 64*1024 {
		t.Fatalf("provider requested %d bytes in one read; want bounded streaming", tracking.maxRequest)
	}

	reader = newIntegrityReader(bytes.NewReader(payload), int64(len(payload)), SHA256Hex([]byte("different")))
	_, err := io.Copy(io.Discard, reader)
	if err == nil {
		t.Fatal("io.Copy(checksum mismatch) error = nil")
	}
	assertStorageFailureCode(t, err, codeIntegrityFailed)
}

func TestIntegrityReaderDetectsDeclaredLengthMismatch(t *testing.T) {
	payload := []byte("short")
	reader := newIntegrityReader(bytes.NewReader(payload), int64(len(payload)+1), SHA256Hex(payload))
	_, err := io.Copy(io.Discard, reader)
	if err == nil {
		t.Fatal("short upload error = nil")
	}
	assertStorageFailureCode(t, err, codeLengthInvalid)

	reader = newIntegrityReader(bytes.NewReader(append(payload, 'x')), int64(len(payload)), SHA256Hex(payload))
	buffer := make([]byte, 16)
	err = nil
	for err == nil {
		_, err = reader.Read(buffer)
	}
	if err == io.EOF {
		t.Fatal("long upload returned EOF instead of length failure")
	}
	assertStorageFailureCode(t, err, codeLengthInvalid)
}

func TestVerifiedReadCloserDetectsCorruption(t *testing.T) {
	payload := []byte("synthetic-download")
	valid := newVerifiedReadCloser(io.NopCloser(bytes.NewReader(payload)), int64(len(payload)), SHA256Hex(payload))
	decoded, err := io.ReadAll(valid)
	if err != nil {
		t.Fatalf("ReadAll(valid) error = %v", err)
	}
	if !bytes.Equal(decoded, payload) {
		t.Fatalf("decoded = %#v, want %#v", decoded, payload)
	}

	corrupt := newVerifiedReadCloser(io.NopCloser(bytes.NewReader(payload)), int64(len(payload)), SHA256Hex([]byte("other")))
	_, err = io.ReadAll(corrupt)
	if err == nil {
		t.Fatal("ReadAll(corrupt) error = nil")
	}
	assertStorageFailureCode(t, err, codeIntegrityFailed)
}

func TestProviderFailurePublicProjectionDoesNotLeakCause(t *testing.T) {
	const secret = "https://access:secret@storage.internal.example/private-object"
	err := classifyOperationFailure(errors.New(secret))
	structured, ok := failure.As(err)
	if !ok {
		t.Fatalf("failure.As() = false for %v", err)
	}
	public := structured.Public()
	joined := strings.Join([]string{string(public.Code), string(public.Category), public.Title, public.Detail, err.Error()}, " ")
	if strings.Contains(joined, secret) || strings.Contains(joined, "access:secret") {
		t.Fatalf("public failure leaked provider secret: %q", joined)
	}
	if !structured.Retryable() {
		t.Fatal("provider operation failure must be retryable")
	}
}

func mustStorageConfiguration(t *testing.T, overrides map[string]string) config.Config {
	t.Helper()
	values := map[string]string{
		storageEndpointKey:  "http://127.0.0.1:9090",
		storageAccessKeyKey: "synthetic-access-key",
		storageSecretKeyKey: "synthetic-secret-key",
		storageBucketKey:    "omnexa-test-bucket",
	}
	for key, value := range overrides {
		values[key] = value
	}
	resolved, err := LoadConfiguration(config.Options{Overrides: values, Strict: true})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}
	return resolved
}

func assertStorageFailureCode(t *testing.T, err error, code failure.Code) {
	t.Helper()
	if !failure.IsCode(err, code) {
		actual, _ := failure.CodeOf(err)
		t.Fatalf("failure code = %q, want %q (error %v)", actual, code, err)
	}
}

func storageMapValues(values map[string]string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}

type trackingReader struct {
	reader     *bytes.Reader
	maxRequest int
}

func (reader *trackingReader) Read(buffer []byte) (int, error) {
	if len(buffer) > reader.maxRequest {
		reader.maxRequest = len(buffer)
	}
	return reader.reader.Read(buffer)
}

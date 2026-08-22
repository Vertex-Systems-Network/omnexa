package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	checksumMetadataKey   = "omnexa-sha256"
	fileNameMetadataKey   = "omnexa-file-name"
	maxMetadataEntries    = 16
	maxMetadataKeyBytes   = 64
	maxMetadataValueBytes = 512
	maxFileNameBytes      = 255
)

var metadataKeyPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// Upload is a streaming object write. ContentLength is mandatory so P01.06 can
// enforce explicit size bounds without buffering the object. SHA256 is the
// lowercase hexadecimal digest of the payload and is persisted as integrity
// metadata.
type Upload struct {
	Body          io.Reader
	ContentLength int64
	ContentType   string
	FileName      string
	SHA256        string
	Metadata      map[string]string
}

// ObjectInfo is safe provider-neutral object metadata. It deliberately exposes
// no signed URL, credentials or provider-internal request diagnostics.
type ObjectInfo struct {
	Key           string
	ContentLength int64
	ContentType   string
	FileName      string
	SHA256        string
	ETag          string
	LastModified  time.Time
	Metadata      map[string]string
}

// ObjectReader streams object content while verifying the P01.06 SHA-256
// integrity metadata at end-of-stream.
type ObjectReader struct {
	Info ObjectInfo
	Body io.ReadCloser
}

// SHA256Hex computes the canonical lowercase SHA-256 hex digest for synthetic
// fixtures and caller-side streaming preparation.
func SHA256Hex(value []byte) string {
	digest := sha256.Sum256(value)
	return hex.EncodeToString(digest[:])
}

func validateUpload(upload Upload, maxObjectBytes int64) (map[string]string, error) {
	if upload.Body == nil {
		return nil, invalidMetadata("upload body is required")
	}
	if upload.ContentLength < 0 || upload.ContentLength > maxObjectBytes {
		return nil, safeFailure(
			codeLengthInvalid,
			failure.CategoryValidation,
			"storage object length is invalid",
			failure.WithDetail("content length is outside the configured P01.06 bounds"),
		)
	}
	checksum, err := normalizeChecksum(upload.SHA256)
	if err != nil {
		return nil, err
	}
	if upload.ContentType != "" {
		if len(upload.ContentType) > 255 {
			return nil, invalidMetadata("content type exceeds the supported length")
		}
		if _, _, err := mime.ParseMediaType(upload.ContentType); err != nil {
			return nil, invalidMetadata("content type is invalid")
		}
	}
	if err := validateFileName(upload.FileName); err != nil {
		return nil, err
	}
	if len(upload.Metadata) > maxMetadataEntries {
		return nil, invalidMetadata("too many custom metadata entries")
	}

	metadata := make(map[string]string, len(upload.Metadata)+2)
	for key, value := range upload.Metadata {
		key = strings.ToLower(strings.TrimSpace(key))
		if len(key) > maxMetadataKeyBytes || !metadataKeyPattern.MatchString(key) {
			return nil, invalidMetadata("custom metadata key is invalid")
		}
		if key == checksumMetadataKey || key == fileNameMetadataKey {
			return nil, invalidMetadata("custom metadata uses a reserved P01.06 key")
		}
		if !validMetadataValue(value) {
			return nil, invalidMetadata("custom metadata value is invalid")
		}
		metadata[key] = value
	}
	metadata[checksumMetadataKey] = checksum
	if upload.FileName != "" {
		metadata[fileNameMetadataKey] = upload.FileName
	}
	return metadata, nil
}

func normalizeChecksum(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if len(value) != sha256.Size*2 {
		return "", invalidMetadata("SHA-256 checksum must contain exactly 64 hexadecimal characters")
	}
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != sha256.Size {
		return "", invalidMetadata("SHA-256 checksum is invalid")
	}
	return value, nil
}

func validateFileName(value string) error {
	if value == "" {
		return nil
	}
	if len(value) > maxFileNameBytes || !utf8.ValidString(value) {
		return invalidMetadata("file name is invalid")
	}
	if strings.ContainsAny(value, "/\\") || value == "." || value == ".." || containsControl(value) {
		return invalidMetadata("file name contains path or control characters")
	}
	return nil
}

func validMetadataValue(value string) bool {
	return len(value) <= maxMetadataValueBytes && utf8.ValidString(value) && !containsControl(value)
}

func containsControl(value string) bool {
	for _, r := range value {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func invalidMetadata(detail string) error {
	return safeFailure(
		codeMetadataInvalid,
		failure.CategoryValidation,
		"storage object metadata is invalid",
		failure.WithDetail(detail),
	)
}

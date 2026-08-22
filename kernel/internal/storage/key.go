package storage

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const maxRenderedObjectKeyBytes = 1024

var (
	objectNamespacePattern   = regexp.MustCompile(`^[a-z0-9][a-z0-9_.-]{0,127}$`)
	objectPathSegmentPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)
)

// Key is a provider-neutral, explicitly versioned storage identity. It carries
// no tenancy, authorization, CMS, media-library or business-domain semantics.
type Key struct {
	Namespace string
	Version   uint16
	Path      string
}

// RenderKey validates and renders a deterministic object key under the
// configured application prefix. Traversal and ambiguous empty segments fail
// closed rather than being normalized silently.
func RenderKey(prefix string, key Key) (string, error) {
	if !objectNamespacePattern.MatchString(prefix) {
		return "", invalidKey("storage key prefix is invalid")
	}
	if !objectNamespacePattern.MatchString(key.Namespace) {
		return "", invalidKey("storage key namespace is invalid")
	}
	if key.Version == 0 {
		return "", invalidKey("storage key version must be greater than zero")
	}
	if key.Path == "" || strings.HasPrefix(key.Path, "/") || strings.HasSuffix(key.Path, "/") || strings.Contains(key.Path, "//") || strings.Contains(key.Path, "\\") {
		return "", invalidKey("storage object path is invalid")
	}

	segments := strings.Split(key.Path, "/")
	for _, segment := range segments {
		if segment == "." || segment == ".." || !objectPathSegmentPattern.MatchString(segment) {
			return "", invalidKey("storage object path contains an invalid segment")
		}
	}

	var builder strings.Builder
	builder.Grow(len(prefix) + len(key.Namespace) + len(key.Path) + 16)
	builder.WriteString(prefix)
	builder.WriteByte('/')
	builder.WriteString(key.Namespace)
	builder.WriteString("/v")
	builder.WriteString(strconv.FormatUint(uint64(key.Version), 10))
	builder.WriteByte('/')
	builder.WriteString(key.Path)
	value := builder.String()
	if len(value) > maxRenderedObjectKeyBytes {
		return "", invalidKey("storage object key exceeds the supported length")
	}
	return value, nil
}

func invalidKey(detail string) error {
	return safeFailure(
		codeKeyInvalid,
		failure.CategoryValidation,
		"storage object key is invalid",
		failure.WithDetail(detail),
	)
}

package cache

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const maxRenderedKeyBytes = 512

var keySegmentPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_.-]{0,127}$`)

// Key is a provider-neutral, explicitly versioned cache identity. It does not
// infer tenancy, authorization, sessions, or business-domain ownership.
type Key struct {
	Namespace string
	Version   uint16
	Name      string
}

// RenderKey validates and renders a deterministic provider key under the
// configured application prefix.
func RenderKey(prefix string, key Key) (string, error) {
	if !keySegmentPattern.MatchString(prefix) {
		return "", invalidKey("cache key prefix is invalid")
	}
	if !keySegmentPattern.MatchString(key.Namespace) {
		return "", invalidKey("cache key namespace is invalid")
	}
	if key.Version == 0 {
		return "", invalidKey("cache key version must be greater than zero")
	}
	if !keySegmentPattern.MatchString(key.Name) {
		return "", invalidKey("cache key name is invalid")
	}

	var builder strings.Builder
	builder.Grow(len(prefix) + len(key.Namespace) + len(key.Name) + 16)
	builder.WriteString(prefix)
	builder.WriteByte(':')
	builder.WriteString(key.Namespace)
	builder.WriteString(":v")
	builder.WriteString(strconv.FormatUint(uint64(key.Version), 10))
	builder.WriteByte(':')
	builder.WriteString(key.Name)
	value := builder.String()
	if len(value) > maxRenderedKeyBytes {
		return "", invalidKey("cache key exceeds the supported length")
	}
	return value, nil
}

func invalidKey(detail string) error {
	return safeFailure(
		codeKeyInvalid,
		failure.CategoryValidation,
		"cache key is invalid",
		failure.WithDetail(detail),
	)
}

func keyDiagnostic(key Key) string {
	return fmt.Sprintf("namespace=%s version=%d", key.Namespace, key.Version)
}

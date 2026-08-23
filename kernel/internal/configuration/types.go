package configuration

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxDescriptionRunes = 512
	maxStringValueRunes = 4096
)

var (
	keyPattern    = regexp.MustCompile(`^[a-z][a-z0-9_]*(?:\.[a-z][a-z0-9_]*)+$`)
	ownerPattern  = regexp.MustCompile(`^[a-z][a-z0-9_]*(?:\.[a-z][a-z0-9_]*)+$`)
	uuidV7Pattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
)

// Key is a stable dot-separated runtime configuration identifier.
type Key string

// Kind is the typed representation of a runtime configuration value.
type Kind string

const (
	KindBool     Kind = "bool"
	KindString   Kind = "string"
	KindInt      Kind = "int"
	KindDuration Kind = "duration"
)

func (kind Kind) valid() bool {
	switch kind {
	case KindBool, KindString, KindInt, KindDuration:
		return true
	default:
		return false
	}
}

// Class declares the operational semantics of a definition. A kill switch is
// disable-only: true means the governed capability must stop or degrade safely.
type Class string

const (
	ClassRuntimeConfig Class = "runtime_config"
	ClassFeatureFlag   Class = "feature_flag"
	ClassKillSwitch    Class = "kill_switch"
)

func (class Class) valid() bool {
	switch class {
	case ClassRuntimeConfig, ClassFeatureFlag, ClassKillSwitch:
		return true
	default:
		return false
	}
}

// Source identifies where one effective evaluation value came from.
type Source string

const (
	SourceDefault  Source = "default"
	SourceProvider Source = "provider"
	SourceFallback Source = "fallback"
)

// Reason describes safe degraded evaluation state without exposing provider errors.
type Reason string

const (
	ReasonNone          Reason = "none"
	ReasonUnavailable   Reason = "provider_unavailable"
	ReasonTimeout       Reason = "provider_timeout"
	ReasonPanic         Reason = "provider_panic"
	ReasonInvalidResult Reason = "provider_invalid_result"
)

// Value is a small typed runtime value. It is intentionally not a generic secret
// container; sensitive credentials remain owned by the secrets/configuration policy.
type Value struct {
	kind          Kind
	boolValue     bool
	stringValue   string
	intValue      int64
	durationValue time.Duration
}

func BoolValue(value bool) Value {
	return Value{kind: KindBool, boolValue: value}
}

func StringValue(value string) Value {
	return Value{kind: KindString, stringValue: value}
}

func IntValue(value int64) Value {
	return Value{kind: KindInt, intValue: value}
}

func DurationValue(value time.Duration) Value {
	return Value{kind: KindDuration, durationValue: value}
}

func (value Value) Kind() Kind {
	return value.kind
}

func (value Value) Bool() (bool, bool) {
	return value.boolValue, value.kind == KindBool
}

func (value Value) String() (string, bool) {
	return value.stringValue, value.kind == KindString
}

func (value Value) Int() (int64, bool) {
	return value.intValue, value.kind == KindInt
}

func (value Value) Duration() (time.Duration, bool) {
	return value.durationValue, value.kind == KindDuration
}

func (value Value) equal(other Value) bool {
	if value.kind != other.kind {
		return false
	}
	switch value.kind {
	case KindBool:
		return value.boolValue == other.boolValue
	case KindString:
		return value.stringValue == other.stringValue
	case KindInt:
		return value.intValue == other.intValue
	case KindDuration:
		return value.durationValue == other.durationValue
	default:
		return false
	}
}

func (value Value) valid() bool {
	if !value.kind.valid() {
		return false
	}
	if value.kind == KindString {
		return utf8.RuneCountInString(value.stringValue) <= maxStringValueRunes && !hasControl(value.stringValue)
	}
	return true
}

// Definition declares one owner-bound runtime value contract. Version is the
// definition/schema revision, not an authorization or entitlement version.
type Definition struct {
	Key         Key
	Description string
	Owner       string
	Kind        Kind
	Class       Class
	Version     uint64
	Default     Value
	Fallback    Value
}

func validateDefinition(definition Definition) error {
	if !keyPattern.MatchString(string(definition.Key)) {
		return safeFailure(codeDefinitionInvalid, "runtime configuration definition is invalid")
	}
	if strings.TrimSpace(definition.Description) == "" || utf8.RuneCountInString(definition.Description) > maxDescriptionRunes || hasControl(definition.Description) {
		return safeFailure(codeDefinitionInvalid, "runtime configuration definition is invalid")
	}
	if !ownerPattern.MatchString(definition.Owner) || !definition.Kind.valid() || !definition.Class.valid() || definition.Version == 0 {
		return safeFailure(codeDefinitionInvalid, "runtime configuration definition is invalid")
	}
	if !definition.Default.valid() || !definition.Fallback.valid() || definition.Default.Kind() != definition.Kind || definition.Fallback.Kind() != definition.Kind {
		return safeFailure(codeDefinitionInvalid, "runtime configuration definition is invalid")
	}
	if (definition.Class == ClassFeatureFlag || definition.Class == ClassKillSwitch) && definition.Kind != KindBool {
		return safeFailure(codeDefinitionInvalid, "runtime configuration definition is invalid")
	}
	if definition.Class == ClassKillSwitch {
		fallback, ok := definition.Fallback.Bool()
		if !ok || !fallback {
			return safeFailure(codeDefinitionInvalid, "kill switch fallback must remain fail closed")
		}
	}
	return nil
}

// EvaluationContext carries future scope references as opaque metadata only.
// It does not authenticate a principal, establish tenancy, or grant authority.
type EvaluationContext struct {
	TenantID       string
	OrganizationID string
	UserID         string
}

func (scope EvaluationContext) validate() error {
	for _, value := range []string{scope.TenantID, scope.OrganizationID, scope.UserID} {
		if value != "" && !uuidV7Pattern.MatchString(value) {
			return safeFailure(codeContextInvalid, "runtime configuration evaluation context is invalid")
		}
	}
	return nil
}

func (scope EvaluationContext) cacheKey() string {
	return scope.TenantID + "\x1f" + scope.OrganizationID + "\x1f" + scope.UserID
}

// ProviderResult is one provider-supplied override and its monotonic revision.
type ProviderResult struct {
	Value    Value
	Revision uint64
}

// Provider resolves an exact definition/context pair. Implementations must honor context cancellation.
type Provider interface {
	Resolve(ctx context.Context, key Key, scope EvaluationContext) (ProviderResult, error)
}

// ErrProviderValueNotFound means no runtime override exists for the exact context.
var ErrProviderValueNotFound = errors.New("runtime configuration provider value not found")

// ChangeAction identifies material cache/evaluation transitions without carrying values or scope IDs.
type ChangeAction string

const (
	ChangeResolved    ChangeAction = "resolved"
	ChangeInvalidated ChangeAction = "invalidated"
)

// Change is classification-safe change metadata suitable for a later audit transport.
type Change struct {
	Key               Key
	Action            ChangeAction
	Source            Source
	DefinitionVersion uint64
	ProviderRevision  uint64
	At                time.Time
}

// Evaluation is the effective typed runtime value plus safe provenance metadata.
type Evaluation struct {
	Key               Key
	Value             Value
	Source            Source
	Reason            Reason
	DefinitionVersion uint64
	ProviderRevision  uint64
	EvaluatedAt       time.Time
	ExpiresAt         time.Time
	CacheHit          bool
	Degraded          bool
	Change            *Change
}

func hasControl(value string) bool {
	for _, char := range value {
		if char < 0x20 && char != '\t' {
			return true
		}
		if char == 0x7f {
			return true
		}
	}
	return false
}

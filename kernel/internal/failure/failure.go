// Package failure implements the P01.03 transport-neutral structured error contract.
package failure

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	// MaxViolations bounds public validation detail so malformed input cannot
	// create unbounded error payloads.
	MaxViolations       = 100
	maxPathRunes        = 256
	maxMessageRunes     = 512
	maxCorrelationRunes = 128
)

var (
	codePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*(?:\.[a-z][a-z0-9_]*)+$`)
	pathPattern = regexp.MustCompile(`^[A-Za-z0-9_\-\[\]\.]+$`)
)

// Code is a stable lowercase dot-separated machine error identifier.
type Code string

// Valid reports whether the code follows the frozen Omnexa code convention.
func (code Code) Valid() bool {
	return codePattern.MatchString(string(code))
}

// Category classifies a failure independently from presentation text.
type Category string

const (
	CategoryValidation     Category = "validation"
	CategoryAuthentication Category = "authentication"
	CategoryAuthorization  Category = "authorization"
	CategoryNotFound       Category = "not_found"
	CategoryConflict       Category = "conflict"
	CategoryRateLimit      Category = "rate_limit"
	CategoryDependency     Category = "dependency"
	CategoryTimeout        Category = "timeout"
	CategoryUnavailable    Category = "unavailable"
	CategoryInvariant      Category = "invariant"
	CategoryInternal       Category = "internal"
)

// Valid reports whether category is part of the frozen minimum category set.
func (category Category) Valid() bool {
	switch category {
	case CategoryValidation,
		CategoryAuthentication,
		CategoryAuthorization,
		CategoryNotFound,
		CategoryConflict,
		CategoryRateLimit,
		CategoryDependency,
		CategoryTimeout,
		CategoryUnavailable,
		CategoryInvariant,
		CategoryInternal:
		return true
	default:
		return false
	}
}

// Violation is safe, bounded validation detail using contract field names.
type Violation struct {
	Path    string
	Code    Code
	Message string
}

// Public is the transport-neutral safe projection of an Error. It intentionally
// contains no wrapped cause, stack, SQL, provider payload, or diagnostic text.
type Public struct {
	Code                Code
	Category            Category
	Title               string
	Detail              string
	Retryable           bool
	RequestID           string
	TraceID             string
	Violations          []Violation
	ViolationsTruncated bool
}

// Error is the canonical structured kernel failure. Fields are private so a
// caller cannot accidentally mutate classification or expose the internal cause.
type Error struct {
	code                Code
	category            Category
	title               string
	detail              string
	retryable           bool
	requestID           string
	traceID             string
	violations          []Violation
	violationsTruncated bool
	cause               error
}

// Option adds bounded metadata to a failure during construction.
type Option func(*Error) error

// New creates a classified failure after validating contract-safe metadata.
func New(code Code, category Category, title string, options ...Option) (*Error, error) {
	if !code.Valid() {
		return nil, fmt.Errorf("invalid Omnexa error code %q", code)
	}
	if !category.Valid() {
		return nil, fmt.Errorf("invalid Omnexa error category %q", category)
	}
	if err := validateSafeText("title", title, 256, true); err != nil {
		return nil, err
	}

	failure := &Error{code: code, category: category, title: title}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(failure); err != nil {
			return nil, err
		}
	}
	return failure, nil
}

// Wrap creates a classified failure while retaining a private lower-level cause.
func Wrap(cause error, code Code, category Category, title string, options ...Option) (*Error, error) {
	if cause == nil {
		return nil, errors.New("failure wrap requires a non-nil cause")
	}
	options = append(options, WithCause(cause))
	return New(code, category, title, options...)
}

// WithDetail adds safe public detail. Raw implementation/provider values should
// not be passed here; private diagnostic information belongs in the cause.
func WithDetail(detail string) Option {
	return func(target *Error) error {
		if err := validateSafeText("detail", detail, 1024, false); err != nil {
			return err
		}
		target.detail = detail
		return nil
	}
}

// WithRetryable sets explicit retryability. It never implies that retry is safe
// without the caller's idempotency/operation policy.
func WithRetryable(retryable bool) Option {
	return func(target *Error) error {
		target.retryable = retryable
		return nil
	}
}

// WithCorrelation attaches non-authoritative diagnostic identifiers.
func WithCorrelation(requestID, traceID string) Option {
	return func(target *Error) error {
		if err := validateCorrelation("request_id", requestID); err != nil {
			return err
		}
		if err := validateCorrelation("trace_id", traceID); err != nil {
			return err
		}
		target.requestID = requestID
		target.traceID = traceID
		return nil
	}
}

// WithViolations adds deterministic, deduplicated and bounded validation detail.
func WithViolations(violations ...Violation) Option {
	return func(target *Error) error {
		normalized, truncated, err := normalizeViolations(violations)
		if err != nil {
			return err
		}
		target.violations = normalized
		target.violationsTruncated = truncated
		return nil
	}
}

// WithCause retains private diagnostic causality. Cause text is never included
// by Error(), Public(), or other safe projections.
func WithCause(cause error) Option {
	return func(target *Error) error {
		if cause == nil {
			return errors.New("failure cause cannot be nil")
		}
		target.cause = cause
		return nil
	}
}

// Error implements error using only contract-safe fields.
func (failure *Error) Error() string {
	if failure == nil {
		return "<nil>"
	}
	return string(failure.code) + ": " + failure.title
}

// Unwrap retains standard errors.Is/errors.As traversal to the private cause.
func (failure *Error) Unwrap() error {
	if failure == nil {
		return nil
	}
	return failure.cause
}

// Code returns the stable machine code.
func (failure *Error) Code() Code {
	if failure == nil {
		return ""
	}
	return failure.code
}

// Category returns the stable failure category.
func (failure *Error) Category() Category {
	if failure == nil {
		return ""
	}
	return failure.category
}

// Retryable returns explicit retryability metadata.
func (failure *Error) Retryable() bool {
	return failure != nil && failure.retryable
}

// Public returns a defensive safe copy with no internal cause.
func (failure *Error) Public() Public {
	if failure == nil {
		return Public{}
	}
	violations := append([]Violation(nil), failure.violations...)
	return Public{
		Code:                failure.code,
		Category:            failure.category,
		Title:               failure.title,
		Detail:              failure.detail,
		Retryable:           failure.retryable,
		RequestID:           failure.requestID,
		TraceID:             failure.traceID,
		Violations:          violations,
		ViolationsTruncated: failure.violationsTruncated,
	}
}

// As returns the nearest structured Omnexa failure in an error chain.
func As(err error) (*Error, bool) {
	var target *Error
	if !errors.As(err, &target) {
		return nil, false
	}
	return target, true
}

// CodeOf returns the stable code of the nearest structured failure.
func CodeOf(err error) (Code, bool) {
	failure, ok := As(err)
	if !ok {
		return "", false
	}
	return failure.code, true
}

// IsCode reports whether the nearest structured failure has the requested code.
func IsCode(err error, code Code) bool {
	actual, ok := CodeOf(err)
	return ok && actual == code
}

func normalizeViolations(input []Violation) ([]Violation, bool, error) {
	if len(input) == 0 {
		return nil, false, nil
	}

	items := make([]Violation, 0, len(input))
	seen := make(map[string]struct{}, len(input))
	for _, violation := range input {
		if violation.Path == "" || utf8.RuneCountInString(violation.Path) > maxPathRunes || !pathPattern.MatchString(violation.Path) {
			return nil, false, fmt.Errorf("invalid validation path %q", violation.Path)
		}
		if !violation.Code.Valid() {
			return nil, false, fmt.Errorf("invalid validation error code %q", violation.Code)
		}
		if err := validateSafeText("validation message", violation.Message, maxMessageRunes, true); err != nil {
			return nil, false, err
		}
		key := violation.Path + "\x00" + string(violation.Code) + "\x00" + violation.Message
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		items = append(items, violation)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Path != items[j].Path {
			return items[i].Path < items[j].Path
		}
		if items[i].Code != items[j].Code {
			return items[i].Code < items[j].Code
		}
		return items[i].Message < items[j].Message
	})

	truncated := len(items) > MaxViolations
	if truncated {
		items = items[:MaxViolations]
	}
	return items, truncated, nil
}

func validateCorrelation(name, value string) error {
	if value == "" {
		return nil
	}
	if utf8.RuneCountInString(value) > maxCorrelationRunes || hasControl(value) {
		return fmt.Errorf("%s contains invalid diagnostic identifier data", name)
	}
	return nil
}

func validateSafeText(name, value string, maxRunes int, required bool) error {
	if required && strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", name)
	}
	if utf8.RuneCountInString(value) > maxRunes {
		return fmt.Errorf("%s exceeds %d characters", name, maxRunes)
	}
	if hasControl(value) {
		return fmt.Errorf("%s contains control characters", name)
	}
	return nil
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

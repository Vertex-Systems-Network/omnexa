package failure

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestCodeValidation(t *testing.T) {
	for _, code := range []Code{
		"validation.failed",
		"authentication.required",
		"commerce.order.invalid_state",
		"integration.provider.unavailable",
	} {
		if !code.Valid() {
			t.Errorf("Code(%q).Valid() = false, want true", code)
		}
	}

	for _, code := range []Code{
		"",
		"internal",
		"Validation.failed",
		"validation-failed",
		"validation..failed",
		"validation.failed!",
		"validation.123",
	} {
		if code.Valid() {
			t.Errorf("Code(%q).Valid() = true, want false", code)
		}
	}
}

func TestCategoryValidation(t *testing.T) {
	valid := []Category{
		CategoryValidation,
		CategoryAuthentication,
		CategoryAuthorization,
		CategoryNotFound,
		CategoryConflict,
		CategoryRateLimit,
		CategoryDependency,
		CategoryTimeout,
		CategoryUnavailable,
		CategoryInvariant,
		CategoryInternal,
	}
	for _, category := range valid {
		if !category.Valid() {
			t.Errorf("Category(%q).Valid() = false, want true", category)
		}
	}
	if Category("database").Valid() {
		t.Fatal("unknown category unexpectedly valid")
	}
}

func TestWrapPreservesCauseAndNeverPublishesPrivateCause(t *testing.T) {
	const secret = "password=never-publish SQL=select-secret /home/private/path"
	cause := errors.New(secret)

	failure, err := Wrap(
		cause,
		"dependency.unavailable",
		CategoryDependency,
		"Required dependency is unavailable",
		WithDetail("Try the operation again later."),
		WithRetryable(true),
		WithCorrelation("018f47b1-f3a8-79da-8d9d-4a64d65c2fd5", "6f1a4c0d88b4427e"),
	)
	if err != nil {
		t.Fatalf("Wrap() error = %v", err)
	}

	if !errors.Is(failure, cause) {
		t.Fatal("errors.Is() did not preserve wrapped cause")
	}
	var typed *Error
	if !errors.As(failure, &typed) || typed != failure {
		t.Fatal("errors.As() did not recover structured failure")
	}
	if !failure.Retryable() {
		t.Fatal("Retryable() = false, want true")
	}

	public := failure.Public()
	if public.Code != "dependency.unavailable" || public.Category != CategoryDependency {
		t.Fatalf("Public() classification = %#v", public)
	}
	if public.Title != "Required dependency is unavailable" || public.Detail != "Try the operation again later." {
		t.Fatalf("Public() safe text = %#v", public)
	}
	if public.RequestID == "" || public.TraceID == "" {
		t.Fatalf("Public() correlation metadata missing: %#v", public)
	}

	for name, text := range map[string]string{
		"Error()":  failure.Error(),
		"Public()": fmt.Sprintf("%+v", public),
	} {
		if strings.Contains(text, secret) || strings.Contains(text, "select-secret") || strings.Contains(text, "/home/private") {
			t.Fatalf("%s leaked private cause: %q", name, text)
		}
	}
}

func TestCodeLookupSurvivesOuterWrapping(t *testing.T) {
	failure, err := New("conflict.version_mismatch", CategoryConflict, "Version conflict")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	outer := fmt.Errorf("operation failed: %w", failure)

	if !IsCode(outer, "conflict.version_mismatch") {
		t.Fatal("IsCode() = false, want true")
	}
	code, ok := CodeOf(outer)
	if !ok || code != "conflict.version_mismatch" {
		t.Fatalf("CodeOf() = %q, %v", code, ok)
	}
	found, ok := As(outer)
	if !ok || found != failure {
		t.Fatal("As() did not return nearest structured failure")
	}
}

func TestRetryabilityIsExplicit(t *testing.T) {
	nonRetryable, err := New("validation.failed", CategoryValidation, "Validation failed")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if nonRetryable.Retryable() {
		t.Fatal("default Retryable() = true, want false")
	}

	retryable, err := New(
		"service.unavailable",
		CategoryUnavailable,
		"Service temporarily unavailable",
		WithRetryable(true),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if !retryable.Retryable() {
		t.Fatal("explicit Retryable() = false, want true")
	}
}

func TestViolationsAreDeterministicDeduplicatedAndBounded(t *testing.T) {
	violations := []Violation{
		{Path: "profile.email", Code: "validation.email.invalid", Message: "Enter a valid email address."},
		{Path: "account.name", Code: "validation.required", Message: "Name is required."},
		{Path: "profile.email", Code: "validation.email.invalid", Message: "Enter a valid email address."},
	}
	for i := 0; i < MaxViolations+10; i++ {
		violations = append(violations, Violation{
			Path:    fmt.Sprintf("items[%03d].name", i),
			Code:    "validation.required",
			Message: "Name is required.",
		})
	}

	failure, err := New(
		"validation.failed",
		CategoryValidation,
		"Validation failed",
		WithViolations(violations...),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	public := failure.Public()
	if len(public.Violations) != MaxViolations {
		t.Fatalf("len(Violations) = %d, want %d", len(public.Violations), MaxViolations)
	}
	if !public.ViolationsTruncated {
		t.Fatal("ViolationsTruncated = false, want true")
	}

	for i := 1; i < len(public.Violations); i++ {
		previous := public.Violations[i-1]
		current := public.Violations[i]
		if violationSortKey(previous) > violationSortKey(current) {
			t.Fatalf("violations not sorted at %d: %#v > %#v", i, previous, current)
		}
		if previous == current {
			t.Fatalf("duplicate violation retained: %#v", current)
		}
	}
}

func TestPublicReturnsDefensiveViolationCopy(t *testing.T) {
	failure, err := New(
		"validation.failed",
		CategoryValidation,
		"Validation failed",
		WithViolations(Violation{Path: "email", Code: "validation.required", Message: "Email is required."}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	first := failure.Public()
	first.Violations[0].Message = "mutated"
	second := failure.Public()
	if second.Violations[0].Message != "Email is required." {
		t.Fatalf("Public() leaked mutable slice state: %#v", second.Violations)
	}
}

func TestConstructorsRejectUnsafeOrInvalidContractData(t *testing.T) {
	tests := []struct {
		name string
		make func() error
	}{
		{
			name: "invalid code",
			make: func() error {
				_, err := New("INVALID", CategoryInternal, "Internal failure")
				return err
			},
		},
		{
			name: "invalid category",
			make: func() error {
				_, err := New("internal.failure", Category("database"), "Internal failure")
				return err
			},
		},
		{
			name: "empty title",
			make: func() error {
				_, err := New("internal.failure", CategoryInternal, "   ")
				return err
			},
		},
		{
			name: "control character",
			make: func() error {
				_, err := New("internal.failure", CategoryInternal, "unsafe\nheader")
				return err
			},
		},
		{
			name: "invalid correlation",
			make: func() error {
				_, err := New("internal.failure", CategoryInternal, "Internal failure", WithCorrelation("bad\nrequest", ""))
				return err
			},
		},
		{
			name: "invalid violation path",
			make: func() error {
				_, err := New("validation.failed", CategoryValidation, "Validation failed", WithViolations(
					Violation{Path: "database column", Code: "validation.required", Message: "Required."},
				))
				return err
			},
		},
		{
			name: "nil wrap cause",
			make: func() error {
				_, err := Wrap(nil, "internal.failure", CategoryInternal, "Internal failure")
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.make(); err == nil {
				t.Fatal("error = nil, want validation failure")
			}
		})
	}
}

func violationSortKey(violation Violation) string {
	return violation.Path + "\x00" + string(violation.Code) + "\x00" + violation.Message
}

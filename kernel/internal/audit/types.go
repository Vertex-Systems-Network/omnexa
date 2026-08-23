package audit

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxReferenceRunes   = 128
	maxActionRunes      = 128
	maxReasonRunes      = 512
	maxFieldKeyRunes    = 64
	maxFieldValueRunes  = 1024
	maxFields           = 32
	maxTagsPerField     = 16
)

var (
	stableTokenPattern = regexp.MustCompile(`^[a-z][a-z0-9_.:-]{0,127}$`)
	fieldKeyPattern    = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)
	uuidV7Pattern      = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
)

// Classification mirrors the frozen Omnexa confidentiality vocabulary for
// protected audit records. It does not create or lower a source data class.
type Classification string

const (
	ClassificationPublic       Classification = "PUBLIC"
	ClassificationInternal     Classification = "INTERNAL"
	ClassificationConfidential Classification = "CONFIDENTIAL"
	ClassificationRestricted   Classification = "RESTRICTED"
)

func (classification Classification) valid() bool {
	switch classification {
	case ClassificationPublic, ClassificationInternal, ClassificationConfidential, ClassificationRestricted:
		return true
	default:
		return false
	}
}

func (classification Classification) rank() int {
	switch classification {
	case ClassificationPublic:
		return 1
	case ClassificationInternal:
		return 2
	case ClassificationConfidential:
		return 3
	case ClassificationRestricted:
		return 4
	default:
		return 0
	}
}

// HandlingTag mirrors the frozen data-classification handling tags. Tags are
// descriptive obligations only and never authority or identity signals.
type HandlingTag string

const (
	TagPII               HandlingTag = "PII"
	TagAuthSecret        HandlingTag = "AUTH_SECRET"
	TagCryptoKey         HandlingTag = "CRYPTO_KEY"
	TagPaymentSensitive  HandlingTag = "PAYMENT_SENSITIVE"
	TagFinancialRecord   HandlingTag = "FINANCIAL_RECORD"
	TagEmployeeData      HandlingTag = "EMPLOYEE_DATA"
	TagLegalRecord       HandlingTag = "LEGAL_RECORD"
	TagHealthSensitive   HandlingTag = "HEALTH_SENSITIVE"
	TagChildData         HandlingTag = "CHILD_DATA"
	TagBiometric         HandlingTag = "BIOMETRIC"
	TagLocationPrecise   HandlingTag = "LOCATION_PRECISE"
	TagSecurityTelemetry HandlingTag = "SECURITY_TELEMETRY"
	TagCustomerContent   HandlingTag = "CUSTOMER_CONTENT"
	TagModelInput        HandlingTag = "MODEL_INPUT"
	TagModelOutput       HandlingTag = "MODEL_OUTPUT"
)

func (tag HandlingTag) valid() bool {
	switch tag {
	case TagPII, TagAuthSecret, TagCryptoKey, TagPaymentSensitive, TagFinancialRecord, TagEmployeeData,
		TagLegalRecord, TagHealthSensitive, TagChildData, TagBiometric, TagLocationPrecise,
		TagSecurityTelemetry, TagCustomerContent, TagModelInput, TagModelOutput:
		return true
	default:
		return false
	}
}

func (tag HandlingTag) prohibitedInGenericAudit() bool {
	switch tag {
	case TagAuthSecret, TagCryptoKey, TagPaymentSensitive:
		return true
	default:
		return false
	}
}

// Actor is descriptive audit metadata. Reference and ImpersonatorReference are
// opaque references and do not authenticate a principal or grant authority.
type Actor struct {
	Kind                  string
	Reference             string
	ImpersonatorReference string
}

// Target identifies the object class/reference affected by the audited action.
// It does not grant read or write authority over the referenced object.
type Target struct {
	Kind      string
	Reference string
}

// Scope represents either an explicitly platform-scoped audit fact or future
// tenant/organization UUIDv7 metadata without implementing P02 tenancy.
type Scope struct {
	Platform       bool
	TenantID       string
	OrganizationID string
}

// Outcome is the stable minimum result vocabulary for an audited action.
type Outcome string

const (
	OutcomeSucceeded Outcome = "succeeded"
	OutcomeFailed    Outcome = "failed"
	OutcomeDenied    Outcome = "denied"
)

func (outcome Outcome) valid() bool {
	switch outcome {
	case OutcomeSucceeded, OutcomeFailed, OutcomeDenied:
		return true
	default:
		return false
	}
}

// Approval is optional descriptive metadata for a privileged action. It is not
// an authorization decision and carries no permission semantics.
type Approval struct {
	Kind      string
	Reference string
}

// Field is bounded supplemental audit metadata. Generic audit rejects secret,
// key and payment-sensitive fields by default; specialized future handlers need
// explicit governance rather than bypassing this boundary.
type Field struct {
	Key            string
	Value          string
	Classification Classification
	Tags           []HandlingTag
}

// RecordInput is the caller-owned source used to construct one immutable audit
// record. NewRecord validates and defensively copies all mutable input.
type RecordInput struct {
	Classification Classification
	Actor          Actor
	Action         string
	Target         Target
	Scope          Scope
	Outcome        Outcome
	CorrelationID  string
	Reason         string
	Approval       *Approval
	Privileged     bool
	Fields         []Field
}

// Requirement determines whether an acknowledged audit append is mandatory for
// the guarded operation. Best-effort failure is still surfaced as degradation.
type Requirement string

const (
	RequirementRequired   Requirement = "required"
	RequirementBestEffort Requirement = "best_effort"
)

func (requirement Requirement) valid() bool {
	return requirement == RequirementRequired || requirement == RequirementBestEffort
}

// DeliveryStatus is a transport-only outcome. Recorded means the configured
// sink acknowledged append; it does not claim legal retention or cross-process durability.
type DeliveryStatus string

const (
	DeliveryRecorded DeliveryStatus = "recorded"
	DeliveryDegraded DeliveryStatus = "degraded"
	DeliveryFailed   DeliveryStatus = "failed"
)

// DeliveryReason is safe transport-health metadata and never contains a sink error body.
type DeliveryReason string

const (
	DeliveryReasonNone          DeliveryReason = "none"
	DeliveryReasonSinkFailure   DeliveryReason = "sink_failure"
	DeliveryReasonCanceled      DeliveryReason = "canceled"
	DeliveryReasonDeadline      DeliveryReason = "deadline_exceeded"
)

// Receipt is a classification-safe audit delivery projection.
type Receipt struct {
	RecordID string
	Status   DeliveryStatus
	Reason   DeliveryReason
}

// Health is a fixed-size payload-free transport-health projection.
type Health struct {
	Submitted uint64
	Recorded  uint64
	Degraded  uint64
	Failed    uint64
}

func validateInput(input RecordInput) error {
	if !input.Classification.valid() || !input.Outcome.valid() {
		return invalidRecordFailure()
	}
	if !validStableToken(input.Actor.Kind, maxReferenceRunes) || !validReference(input.Actor.Reference, maxReferenceRunes) {
		return invalidRecordFailure()
	}
	if input.Actor.ImpersonatorReference != "" && !validReference(input.Actor.ImpersonatorReference, maxReferenceRunes) {
		return invalidRecordFailure()
	}
	if !validStableToken(input.Action, maxActionRunes) || !validStableToken(input.Target.Kind, maxReferenceRunes) || !validReference(input.Target.Reference, maxReferenceRunes) {
		return invalidRecordFailure()
	}
	if err := validateScope(input.Scope); err != nil {
		return err
	}
	if !validReference(input.CorrelationID, maxReferenceRunes) {
		return invalidRecordFailure()
	}
	if input.Reason != "" && !validText(input.Reason, maxReasonRunes) {
		return invalidRecordFailure()
	}
	if (input.Privileged || input.Actor.ImpersonatorReference != "") && strings.TrimSpace(input.Reason) == "" {
		return invalidRecordFailure()
	}
	if input.Approval != nil {
		if !validStableToken(input.Approval.Kind, maxReferenceRunes) || !validReference(input.Approval.Reference, maxReferenceRunes) {
			return invalidRecordFailure()
		}
	}
	if len(input.Fields) > maxFields {
		return invalidRecordFailure()
	}
	seen := make(map[string]struct{}, len(input.Fields))
	for _, field := range input.Fields {
		if err := validateField(input.Classification, field); err != nil {
			return err
		}
		if _, exists := seen[field.Key]; exists {
			return invalidRecordFailure()
		}
		seen[field.Key] = struct{}{}
	}
	return nil
}

func validateScope(scope Scope) error {
	if scope.Platform {
		if scope.TenantID != "" || scope.OrganizationID != "" {
			return invalidRecordFailure()
		}
		return nil
	}
	if !uuidV7Pattern.MatchString(scope.TenantID) {
		return invalidRecordFailure()
	}
	if scope.OrganizationID != "" && !uuidV7Pattern.MatchString(scope.OrganizationID) {
		return invalidRecordFailure()
	}
	return nil
}

func validateField(recordClass Classification, field Field) error {
	if !fieldKeyPattern.MatchString(field.Key) || utf8.RuneCountInString(field.Key) > maxFieldKeyRunes || sensitiveKey(field.Key) {
		return prohibitedFieldFailure()
	}
	if !field.Classification.valid() || field.Classification.rank() > recordClass.rank() || !validText(field.Value, maxFieldValueRunes) {
		return invalidRecordFailure()
	}
	if len(field.Tags) > maxTagsPerField {
		return invalidRecordFailure()
	}
	seen := make(map[HandlingTag]struct{}, len(field.Tags))
	for _, tag := range field.Tags {
		if !tag.valid() {
			return invalidRecordFailure()
		}
		if tag.prohibitedInGenericAudit() {
			return prohibitedFieldFailure()
		}
		if _, exists := seen[tag]; exists {
			return invalidRecordFailure()
		}
		seen[tag] = struct{}{}
	}
	return nil
}

func normalizeFields(input []Field) []Field {
	fields := make([]Field, len(input))
	for index, field := range input {
		field.Tags = append([]HandlingTag(nil), field.Tags...)
		sort.Slice(field.Tags, func(i, j int) bool { return field.Tags[i] < field.Tags[j] })
		fields[index] = field
	}
	sort.Slice(fields, func(i, j int) bool { return fields[i].Key < fields[j].Key })
	return fields
}

func validStableToken(value string, maxRunes int) bool {
	return utf8.RuneCountInString(value) <= maxRunes && stableTokenPattern.MatchString(value)
}

func validReference(value string, maxRunes int) bool {
	if strings.TrimSpace(value) == "" || utf8.RuneCountInString(value) > maxRunes {
		return false
	}
	for _, char := range value {
		if unicode.IsSpace(char) || unicode.IsControl(char) {
			return false
		}
	}
	return true
}

func validText(value string, maxRunes int) bool {
	if utf8.RuneCountInString(value) > maxRunes {
		return false
	}
	for _, char := range value {
		if unicode.IsControl(char) && char != '\t' {
			return false
		}
	}
	return true
}

func sensitiveKey(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))
	markers := []string{
		"password", "passwd", "secret", "token", "authorization", "cookie",
		"api_key", "apikey", "access_key", "private_key", "credential",
	}
	for _, marker := range markers {
		if key == marker || strings.HasSuffix(key, "_"+marker) || strings.HasPrefix(key, marker+"_") {
			return true
		}
	}
	return false
}

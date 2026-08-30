// Package events implements the P04.01 provider-neutral event envelope contract.
package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/organization"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/tenancy"
	"github.com/google/uuid"
)

const (
	SpecVersion       = "1.0"
	EnvelopeVersion   = uint(1)
	DataContentType   = "application/json"
	maxPayloadBytes   = 64 * 1024
	maxPayloadDepth   = 16
	maxCollectionSize = 256
	maxStringRunes    = 4096
	maxMetadataRunes  = 512
)

var (
	producerPattern  = regexp.MustCompile(`^urn:omnexa:module:[a-z0-9.\-]+$`)
	eventTypePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:\.[a-z][a-z0-9_]*){2,}\.v([1-9][0-9]*)$`)
	schemaPattern    = regexp.MustCompile(`^urn:omnexa:event-schema:[a-z0-9._-]+:v[1-9][0-9]*$`)
	tracePattern     = regexp.MustCompile(`^[0-9a-f]{2}-[0-9a-f]{32}-[0-9a-f]{16}-[0-9a-f]{2}$`)
)

const (
	codeEnvelopeInvalid failure.Code = "events.envelope.invalid"
	codePayloadInvalid  failure.Code = "events.payload.invalid"
	codeTenantMismatch  failure.Code = "events.tenant.mismatch"
	codeIdentifierFail  failure.Code = "events.identifier.failed"
)

// DataClass is the frozen Omnexa data-classification vocabulary.
type DataClass string

const (
	DataClassPublic       DataClass = "PUBLIC"
	DataClassInternal     DataClass = "INTERNAL"
	DataClassConfidential DataClass = "CONFIDENTIAL"
	DataClassRestricted   DataClass = "RESTRICTED"
)

func (class DataClass) Valid() bool {
	switch class {
	case DataClassPublic, DataClassInternal, DataClassConfidential, DataClassRestricted:
		return true
	default:
		return false
	}
}

// ActorType is optional attribution metadata. It never grants authority.
type ActorType string

const (
	ActorUser           ActorType = "user"
	ActorServiceAccount ActorType = "service_account"
	ActorSystem         ActorType = "system"
	ActorAgent          ActorType = "agent"
)

func (actorType ActorType) Valid() bool {
	switch actorType {
	case ActorUser, ActorServiceAccount, ActorSystem, ActorAgent:
		return true
	default:
		return false
	}
}

// EventID, CorrelationID and CausationID are UUIDv7 event identities.
type EventID string
type CorrelationID string
type CausationID string

func (id EventID) Valid() bool       { return validUUIDv7(string(id)) }
func (id CorrelationID) Valid() bool { return validUUIDv7(string(id)) }
func (id CausationID) Valid() bool   { return validUUIDv7(string(id)) }

// Producer is the canonical module producer identity used as the CloudEvents source.
type Producer string

func (producer Producer) Valid() bool {
	return utf8.RuneCountInString(string(producer)) <= maxMetadataRunes && producerPattern.MatchString(string(producer))
}

// EventType is a stable versioned event contract identity.
type EventType string

func (eventType EventType) Valid() bool {
	version, ok := eventType.Version()
	return ok && version == EnvelopeVersion
}

// Version returns the explicit version encoded by the canonical event type suffix.
func (eventType EventType) Version() (uint, bool) {
	matches := eventTypePattern.FindStringSubmatch(string(eventType))
	if len(matches) != 2 || utf8.RuneCountInString(string(eventType)) > maxMetadataRunes {
		return 0, false
	}
	parsed, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil || parsed == 0 {
		return 0, false
	}
	return uint(parsed), true
}

// SchemaID identifies a versioned event payload contract without implementing a registry.
type SchemaID string

func (schema SchemaID) Valid() bool {
	return utf8.RuneCountInString(string(schema)) <= maxMetadataRunes && schemaPattern.MatchString(string(schema))
}

// Params contains the authoritative inputs for constructing one canonical event.
type Params struct {
	Type           EventType
	Producer       Producer
	OccurredAt     time.Time
	TenantID       tenancy.TenantID
	OrganizationID organization.NodeID
	CorrelationID  CorrelationID
	CausationID    CausationID
	Classification DataClass
	DataSchema     SchemaID
	Subject        string
	TraceParent    string
	ActorID        string
	ActorType      ActorType
	SubjectSequence uint64
	Data           any
}

// Envelope is the canonical P04.01 CloudEvents-compatible transport-neutral envelope.
// It carries identity and validated data only; it grants no authorization or delivery guarantee.
type Envelope struct {
	SpecVersion    string              `json:"specversion"`
	ID             EventID             `json:"id"`
	Source         Producer            `json:"source"`
	Type           EventType           `json:"type"`
	EventVersion   uint                `json:"eventversion"`
	OccurredAt     time.Time           `json:"time"`
	DataContentType string             `json:"datacontenttype"`
	DataSchema     SchemaID            `json:"dataschema"`
	TenantID       tenancy.TenantID    `json:"tenantid,omitempty"`
	OrganizationID organization.NodeID `json:"organizationid,omitempty"`
	CorrelationID  CorrelationID       `json:"correlationid"`
	CausationID    CausationID         `json:"causationid,omitempty"`
	Classification DataClass           `json:"classification"`
	Subject        string              `json:"subject,omitempty"`
	TraceParent    string              `json:"traceparent,omitempty"`
	ActorID        string              `json:"actorid,omitempty"`
	ActorType      ActorType           `json:"actortype,omitempty"`
	SubjectSequence uint64             `json:"subjectsequence,omitempty"`
	Data           json.RawMessage     `json:"data"`
}

// New creates one validated canonical envelope and generates a UUIDv7 event identity.
func New(params Params) (Envelope, error) {
	identifier, err := uuid.NewV7()
	if err != nil {
		return Envelope{}, wrappedFailure(err, codeIdentifierFail, failure.CategoryInternal, "event identifier generation failed")
	}
	payload, err := normalizePayload(params.Data)
	if err != nil {
		return Envelope{}, err
	}
	envelope := Envelope{
		SpecVersion:     SpecVersion,
		ID:              EventID(identifier.String()),
		Source:          params.Producer,
		Type:            params.Type,
		EventVersion:    EnvelopeVersion,
		OccurredAt:      params.OccurredAt.UTC(),
		DataContentType: DataContentType,
		DataSchema:      params.DataSchema,
		TenantID:        params.TenantID,
		OrganizationID:  params.OrganizationID,
		CorrelationID:   params.CorrelationID,
		CausationID:     params.CausationID,
		Classification:  params.Classification,
		Subject:         strings.TrimSpace(params.Subject),
		TraceParent:     strings.TrimSpace(params.TraceParent),
		ActorID:         strings.TrimSpace(params.ActorID),
		ActorType:       params.ActorType,
		SubjectSequence: params.SubjectSequence,
		Data:            payload,
	}
	if err := envelope.Validate(); err != nil {
		return Envelope{}, err
	}
	return envelope, nil
}

// Parse validates untrusted serialized input and returns a normalized v1 envelope.
func Parse(serialized []byte) (Envelope, error) {
	if len(serialized) == 0 || len(serialized) > maxPayloadBytes*2 {
		return Envelope{}, invalidEnvelopeFailure()
	}
	decoder := json.NewDecoder(bytes.NewReader(serialized))
	decoder.DisallowUnknownFields()
	var envelope Envelope
	if err := decoder.Decode(&envelope); err != nil {
		return Envelope{}, invalidEnvelopeFailure()
	}
	if decoder.More() {
		return Envelope{}, invalidEnvelopeFailure()
	}
	payload, err := normalizePayload(envelope.Data)
	if err != nil {
		return Envelope{}, err
	}
	envelope.Data = payload
	envelope.OccurredAt = envelope.OccurredAt.UTC()
	if err := envelope.Validate(); err != nil {
		return Envelope{}, err
	}
	return envelope, nil
}

// Marshal returns the canonical deterministic JSON representation of a validated envelope.
func (envelope Envelope) Marshal() ([]byte, error) {
	if err := envelope.Validate(); err != nil {
		return nil, err
	}
	serialized, err := json.Marshal(envelope)
	if err != nil {
		return nil, invalidEnvelopeFailure()
	}
	return serialized, nil
}

// Validate enforces bounded P04.01 identity, tenancy, classification and payload invariants.
func (envelope Envelope) Validate() error {
	version, typeValid := envelope.Type.Version()
	if envelope.SpecVersion != SpecVersion || envelope.EventVersion != EnvelopeVersion || !typeValid || version != envelope.EventVersion {
		return invalidEnvelopeFailure()
	}
	if !envelope.ID.Valid() || !envelope.Source.Valid() || !envelope.DataSchema.Valid() || !envelope.CorrelationID.Valid() {
		return invalidEnvelopeFailure()
	}
	if envelope.OccurredAt.IsZero() || envelope.DataContentType != DataContentType || !envelope.Classification.Valid() {
		return invalidEnvelopeFailure()
	}
	if envelope.CausationID != "" {
		if !envelope.CausationID.Valid() || string(envelope.CausationID) == string(envelope.ID) {
			return invalidEnvelopeFailure()
		}
	}
	if envelope.TenantID != "" && !envelope.TenantID.Valid() {
		return invalidEnvelopeFailure()
	}
	if envelope.OrganizationID != "" {
		if envelope.TenantID == "" || !envelope.OrganizationID.Valid() {
			return invalidEnvelopeFailure()
		}
	}
	if !validOptionalMetadata(envelope.Subject) || !validTraceParent(envelope.TraceParent) {
		return invalidEnvelopeFailure()
	}
	if envelope.ActorType == "" {
		if envelope.ActorID != "" {
			return invalidEnvelopeFailure()
		}
	} else if !envelope.ActorType.Valid() || !validUUIDv7(envelope.ActorID) {
		return invalidEnvelopeFailure()
	}
	if envelope.SubjectSequence > 0 && envelope.Subject == "" {
		return invalidEnvelopeFailure()
	}
	_, err := normalizePayload(envelope.Data)
	return err
}

// ValidateForTenant fails closed unless this event belongs to the exact trusted tenant.
func (envelope Envelope) ValidateForTenant(expected tenancy.TenantID) error {
	if err := envelope.Validate(); err != nil {
		return err
	}
	if !expected.Valid() || envelope.TenantID == "" || envelope.TenantID != expected {
		return tenantMismatchFailure()
	}
	return nil
}

func normalizePayload(value any) (json.RawMessage, error) {
	var raw []byte
	switch typed := value.(type) {
	case json.RawMessage:
		raw = append([]byte(nil), typed...)
	case []byte:
		raw = append([]byte(nil), typed...)
	default:
		encoded, err := json.Marshal(value)
		if err != nil {
			return nil, invalidPayloadFailure()
		}
		raw = encoded
	}
	if len(raw) == 0 || len(raw) > maxPayloadBytes {
		return nil, invalidPayloadFailure()
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, invalidPayloadFailure()
	}
	root, ok := decoded.(map[string]any)
	if !ok {
		return nil, invalidPayloadFailure()
	}
	if err := validatePayloadValue(root, 1); err != nil {
		return nil, err
	}
	canonical, err := json.Marshal(root)
	if err != nil || len(canonical) > maxPayloadBytes {
		return nil, invalidPayloadFailure()
	}
	return canonical, nil
}

func validatePayloadValue(value any, depth int) error {
	if depth > maxPayloadDepth {
		return invalidPayloadFailure()
	}
	switch typed := value.(type) {
	case map[string]any:
		if len(typed) > maxCollectionSize {
			return invalidPayloadFailure()
		}
		for key, nested := range typed {
			if !validPayloadKey(key) || secretLikeKey(key) {
				return invalidPayloadFailure()
			}
			if err := validatePayloadValue(nested, depth+1); err != nil {
				return err
			}
		}
	case []any:
		if len(typed) > maxCollectionSize {
			return invalidPayloadFailure()
		}
		for _, nested := range typed {
			if err := validatePayloadValue(nested, depth+1); err != nil {
				return err
			}
		}
	case string:
		if utf8.RuneCountInString(typed) > maxStringRunes || secretLikeValue(typed) {
			return invalidPayloadFailure()
		}
	case json.Number:
		if _, err := typed.Float64(); err != nil {
			return invalidPayloadFailure()
		}
	case bool, nil:
		return nil
	default:
		return invalidPayloadFailure()
	}
	return nil
}

func validPayloadKey(key string) bool {
	trimmed := strings.TrimSpace(key)
	return trimmed != "" && trimmed == key && utf8.RuneCountInString(key) <= 128 && !strings.ContainsAny(key, "\x00\r\n")
}

func secretLikeKey(key string) bool {
	normalized := strings.ToLower(strings.NewReplacer("-", "_", ".", "_", " ", "_").Replace(key))
	parts := strings.FieldsFunc(normalized, func(r rune) bool { return r == '_' })
	for _, part := range parts {
		switch part {
		case "password", "passwd", "secret", "token", "credential", "credentials", "privatekey", "apikey", "authorization":
			return true
		}
	}
	return strings.Contains(normalized, "private_key") || strings.Contains(normalized, "api_key") || strings.Contains(normalized, "client_secret") || strings.Contains(normalized, "access_token") || strings.Contains(normalized, "refresh_token")
}

func secretLikeValue(value string) bool {
	trimmed := strings.TrimSpace(value)
	upper := strings.ToUpper(trimmed)
	return strings.HasPrefix(trimmed, "Bearer ") || strings.Contains(upper, "-----BEGIN PRIVATE KEY-----") || strings.Contains(upper, "-----BEGIN RSA PRIVATE KEY-----")
}

func validOptionalMetadata(value string) bool {
	return utf8.RuneCountInString(value) <= maxMetadataRunes && !strings.ContainsAny(value, "\x00\r\n") && !secretLikeValue(value)
}

func validTraceParent(value string) bool {
	if value == "" {
		return true
	}
	return len(value) == 55 && tracePattern.MatchString(value)
}

func validUUIDv7(value string) bool {
	parsed, err := uuid.Parse(value)
	return err == nil && parsed.Version() == 7
}

func invalidEnvelopeFailure() error {
	return classifiedFailure(codeEnvelopeInvalid, failure.CategoryValidation, "event envelope is invalid")
}

func invalidPayloadFailure() error {
	return classifiedFailure(codePayloadInvalid, failure.CategoryValidation, "event payload is invalid")
}

func tenantMismatchFailure() error {
	return classifiedFailure(codeTenantMismatch, failure.CategoryAuthorization, "event tenant does not match trusted tenant context")
}

func classifiedFailure(code failure.Code, category failure.Category, title string) error {
	value, err := failure.New(code, category, title, failure.WithRetryable(false))
	if err != nil {
		return errors.New("event failure could not be classified safely")
	}
	return value
}

func wrappedFailure(cause error, code failure.Code, category failure.Category, title string) error {
	value, err := failure.Wrap(cause, code, category, title, failure.WithRetryable(false))
	if err != nil {
		return errors.New("event failure could not be classified safely")
	}
	return value
}

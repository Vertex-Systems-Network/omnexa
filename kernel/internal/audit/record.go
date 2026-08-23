package audit

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/google/uuid"
)

// Record is an immutable, classification-aware audit fact. All fields remain
// private so callers cannot mutate a record after its integrity digest is sealed.
type Record struct {
	id             string
	occurredAt     time.Time
	classification Classification
	actor          Actor
	action         string
	target         Target
	scope          Scope
	outcome        Outcome
	correlationID  string
	reason         string
	approval       Approval
	hasApproval    bool
	privileged     bool
	fields         []Field
	digest         [sha256.Size]byte
}

type canonicalRecord struct {
	ID             string         `json:"id"`
	OccurredAt     string         `json:"occurred_at"`
	Classification Classification `json:"classification"`
	Actor          Actor          `json:"actor"`
	Action         string         `json:"action"`
	Target         Target         `json:"target"`
	Scope          Scope          `json:"scope"`
	Outcome        Outcome        `json:"outcome"`
	CorrelationID  string         `json:"correlation_id"`
	Reason         string         `json:"reason,omitempty"`
	Approval       *Approval      `json:"approval,omitempty"`
	Privileged     bool           `json:"privileged"`
	Fields         []Field        `json:"fields,omitempty"`
}

// NewRecord validates, normalizes and seals one immutable audit record using a
// standards-compliant UUIDv7 identifier and a UTC timestamp.
func NewRecord(input RecordInput) (Record, error) {
	if err := validateInput(input); err != nil {
		return Record{}, err
	}
	identifier, err := uuid.NewV7()
	if err != nil {
		return Record{}, wrappedFailure(err, codeIdentifierFailed, failure.CategoryInternal, "audit record identifier generation failed", false)
	}
	return newRecordAt(input, identifier.String(), time.Now().UTC())
}

func newRecordAt(input RecordInput, identifier string, occurredAt time.Time) (Record, error) {
	if err := validateInput(input); err != nil {
		return Record{}, err
	}
	if !uuidV7Pattern.MatchString(identifier) || occurredAt.IsZero() {
		return Record{}, invalidRecordFailure()
	}

	record := Record{
		id:             identifier,
		occurredAt:     occurredAt.UTC(),
		classification: input.Classification,
		actor:          input.Actor,
		action:         input.Action,
		target:         input.Target,
		scope:          input.Scope,
		outcome:        input.Outcome,
		correlationID:  input.CorrelationID,
		reason:         input.Reason,
		privileged:     input.Privileged,
		fields:         normalizeFields(input.Fields),
	}
	if input.Approval != nil {
		record.approval = *input.Approval
		record.hasApproval = true
	}

	digest, err := record.computeDigest()
	if err != nil {
		return Record{}, err
	}
	record.digest = digest
	return record, nil
}

// ID returns the canonical immutable UUIDv7 audit-record identifier.
func (record Record) ID() string { return record.id }

// OccurredAt returns the canonical UTC record timestamp.
func (record Record) OccurredAt() time.Time { return record.occurredAt }

// Classification returns the record's confidentiality class.
func (record Record) Classification() Classification { return record.classification }

// Actor returns descriptive actor metadata only; it conveys no authority.
func (record Record) Actor() Actor { return record.actor }

// Action returns the stable caller-supplied audit action identifier.
func (record Record) Action() string { return record.action }

// Target returns descriptive target metadata only.
func (record Record) Target() Target { return record.target }

// Scope returns descriptive platform/tenant/organization metadata only.
func (record Record) Scope() Scope { return record.scope }

// Outcome returns the audited operation's stable result vocabulary.
func (record Record) Outcome() Outcome { return record.outcome }

// CorrelationID returns the governed diagnostic correlation identifier.
func (record Record) CorrelationID() string { return record.correlationID }

// Reason returns bounded caller-supplied reason metadata.
func (record Record) Reason() string { return record.reason }

// Approval returns optional descriptive approval metadata.
func (record Record) Approval() (Approval, bool) { return record.approval, record.hasApproval }

// Privileged reports whether the audited operation was marked privileged.
func (record Record) Privileged() bool { return record.privileged }

// Fields returns a defensive copy of bounded supplemental metadata.
func (record Record) Fields() []Field { return normalizeFields(record.fields) }

// IntegrityDigest returns the hexadecimal SHA-256 integrity fingerprint. It is
// evidence of mutation, not a signature, authorization token or retention proof.
func (record Record) IntegrityDigest() string { return hex.EncodeToString(record.digest[:]) }

// Verify checks structural invariants and the sealed digest before transport.
func (record Record) Verify() error {
	input := RecordInput{
		Classification: record.classification,
		Actor:          record.actor,
		Action:         record.action,
		Target:         record.target,
		Scope:          record.scope,
		Outcome:        record.outcome,
		CorrelationID:  record.correlationID,
		Reason:         record.reason,
		Privileged:     record.privileged,
		Fields:         record.fields,
	}
	if record.hasApproval {
		approval := record.approval
		input.Approval = &approval
	}
	if !uuidV7Pattern.MatchString(record.id) || record.occurredAt.IsZero() {
		return integrityFailure()
	}
	if err := validateInput(input); err != nil {
		return integrityFailure()
	}
	expected, err := record.computeDigest()
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare(expected[:], record.digest[:]) != 1 {
		return integrityFailure()
	}
	return nil
}

func (record Record) computeDigest() ([sha256.Size]byte, error) {
	canonical := canonicalRecord{
		ID:             record.id,
		OccurredAt:     record.occurredAt.UTC().Format(time.RFC3339Nano),
		Classification: record.classification,
		Actor:          record.actor,
		Action:         record.action,
		Target:         record.target,
		Scope:          record.scope,
		Outcome:        record.outcome,
		CorrelationID:  record.correlationID,
		Reason:         record.reason,
		Privileged:     record.privileged,
		Fields:         normalizeFields(record.fields),
	}
	if record.hasApproval {
		approval := record.approval
		canonical.Approval = &approval
	}
	encoded, err := json.Marshal(canonical)
	if err != nil {
		return [sha256.Size]byte{}, wrappedFailure(err, codeRecordIntegrity, failure.CategoryInternal, "audit record integrity calculation failed", false)
	}
	return sha256.Sum256(encoded), nil
}

func cloneRecord(record Record) Record {
	cloned := record
	cloned.fields = normalizeFields(record.fields)
	return cloned
}

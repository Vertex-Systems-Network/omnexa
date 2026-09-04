package events

import (
	"bytes"
	"context"
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5"
)

// PostgresInboxStore persists P04.05 processing identity and completion evidence
// through the exact caller-provided transaction. It owns no pool and cannot open,
// commit, or hide an independent transaction.
type PostgresInboxStore struct{}

func NewPostgresInboxStore() *PostgresInboxStore {
	return &PostgresInboxStore{}
}

// Claim establishes the processing identity in the caller transaction. A unique
// identity conflict is resolved by reading the authoritative row after PostgreSQL
// has resolved the competing transaction: exact completed evidence is duplicate,
// different fingerprint is conflict, and committed claimed state is concurrent.
func (store *PostgresInboxStore) Claim(
	ctx context.Context,
	tx OutboxTransaction,
	record InboxRecord,
) (InboxClaimResult, error) {
	if store == nil || tx == nil {
		return InboxClaimResultUnknown, classifiedFailure(codeInboxInvalid, failure.CategoryValidation, "event inbox postgres claim boundary is invalid")
	}
	if err := inboxContextError(ctx); err != nil {
		return InboxClaimResultUnknown, err
	}
	if err := record.validate(); err != nil {
		return InboxClaimResultUnknown, err
	}

	identity := record.Identity
	tag, err := tx.Exec(
		ctx,
		`INSERT INTO omnexa_events.consumer_inbox
			(event_id, owner, consumer_id, event_type, stream, partition_key, tenant_id, canonical_fingerprint, processing_state)
		 VALUES ($1::uuid, $2, $3, $4, $5, $6, $7::uuid, $8, 'claimed')
		 ON CONFLICT DO NOTHING`,
		string(identity.EventID),
		string(identity.Owner),
		identity.ConsumerID,
		string(identity.EventType),
		identity.Stream,
		identity.Partition,
		postgresInboxTenant(identity),
		record.Fingerprint[:],
	)
	if err != nil {
		return InboxClaimResultUnknown, err
	}
	if tag.RowsAffected() == 1 {
		return InboxClaimed, nil
	}
	if tag.RowsAffected() != 0 {
		return InboxClaimResultUnknown, classifiedFailure(codeInboxStateMalformed, failure.CategoryInvariant, "event inbox postgres claim returned an invalid write count")
	}

	var storedFingerprint []byte
	var state string
	err = tx.QueryRow(
		ctx,
		`SELECT canonical_fingerprint, processing_state
		 FROM omnexa_events.consumer_inbox
		 WHERE event_id = $1::uuid
		   AND owner = $2
		   AND consumer_id = $3
		   AND event_type = $4
		   AND stream = $5
		   AND partition_key = $6
		   AND tenant_id IS NOT DISTINCT FROM $7::uuid`,
		string(identity.EventID),
		string(identity.Owner),
		identity.ConsumerID,
		string(identity.EventType),
		identity.Stream,
		identity.Partition,
		postgresInboxTenant(identity),
	).Scan(&storedFingerprint, &state)
	if errors.Is(err, pgx.ErrNoRows) {
		return InboxClaimResultUnknown, classifiedFailure(codeInboxStateMalformed, failure.CategoryInvariant, "event inbox postgres conflicting identity disappeared during claim")
	}
	if err != nil {
		return InboxClaimResultUnknown, err
	}
	if len(storedFingerprint) != len(record.Fingerprint) {
		return InboxClaimResultUnknown, classifiedFailure(codeInboxStateMalformed, failure.CategoryInvariant, "stored event inbox fingerprint is malformed")
	}
	if !bytes.Equal(storedFingerprint, record.Fingerprint[:]) {
		return InboxIdentityConflict, nil
	}
	switch state {
	case "completed":
		return InboxAlreadyCompleted, nil
	case "claimed":
		return InboxConcurrentProcessing, nil
	default:
		return InboxClaimResultUnknown, classifiedFailure(codeInboxStateMalformed, failure.CategoryInvariant, "stored event inbox processing state is malformed")
	}
}

// Complete records local completion only after the protected mutation succeeded,
// still through the exact caller transaction. A missing, conflicting, or already
// transitioned row is an invariant failure rather than an invented success.
func (store *PostgresInboxStore) Complete(
	ctx context.Context,
	tx OutboxTransaction,
	record InboxRecord,
) error {
	if store == nil || tx == nil {
		return classifiedFailure(codeInboxInvalid, failure.CategoryValidation, "event inbox postgres completion boundary is invalid")
	}
	if err := inboxContextError(ctx); err != nil {
		return err
	}
	if err := record.validate(); err != nil {
		return err
	}

	identity := record.Identity
	tag, err := tx.Exec(
		ctx,
		`UPDATE omnexa_events.consumer_inbox
		 SET processing_state = 'completed',
		     completed_at = CURRENT_TIMESTAMP
		 WHERE event_id = $1::uuid
		   AND owner = $2
		   AND consumer_id = $3
		   AND event_type = $4
		   AND stream = $5
		   AND partition_key = $6
		   AND tenant_id IS NOT DISTINCT FROM $7::uuid
		   AND canonical_fingerprint = $8
		   AND processing_state = 'claimed'`,
		string(identity.EventID),
		string(identity.Owner),
		identity.ConsumerID,
		string(identity.EventType),
		identity.Stream,
		identity.Partition,
		postgresInboxTenant(identity),
		record.Fingerprint[:],
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 1 {
		return nil
	}
	if tag.RowsAffected() != 0 {
		return classifiedFailure(codeInboxStateMalformed, failure.CategoryInvariant, "event inbox postgres completion returned an invalid write count")
	}

	var storedFingerprint []byte
	var state string
	err = tx.QueryRow(
		ctx,
		`SELECT canonical_fingerprint, processing_state
		 FROM omnexa_events.consumer_inbox
		 WHERE event_id = $1::uuid
		   AND owner = $2
		   AND consumer_id = $3
		   AND event_type = $4
		   AND stream = $5
		   AND partition_key = $6
		   AND tenant_id IS NOT DISTINCT FROM $7::uuid`,
		string(identity.EventID),
		string(identity.Owner),
		identity.ConsumerID,
		string(identity.EventType),
		identity.Stream,
		identity.Partition,
		postgresInboxTenant(identity),
	).Scan(&storedFingerprint, &state)
	if errors.Is(err, pgx.ErrNoRows) {
		return classifiedFailure(codeInboxStateMalformed, failure.CategoryInvariant, "event inbox postgres completion identity is missing")
	}
	if err != nil {
		return err
	}
	if len(storedFingerprint) != len(record.Fingerprint) || !bytes.Equal(storedFingerprint, record.Fingerprint[:]) {
		return classifiedFailure(codeInboxStateMalformed, failure.CategoryInvariant, "event inbox postgres completion evidence conflicts with stored identity")
	}
	return classifiedFailure(codeInboxStateMalformed, failure.CategoryInvariant, "event inbox postgres completion state is not claimable: "+state)
}

func postgresInboxTenant(identity InboxIdentity) any {
	if identity.TenantID == "" {
		return nil
	}
	return string(identity.TenantID)
}

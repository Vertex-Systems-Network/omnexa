CREATE SCHEMA IF NOT EXISTS omnexa_events;

CREATE TABLE IF NOT EXISTS omnexa_events.consumer_inbox (
    event_id uuid NOT NULL,
    owner text NOT NULL,
    consumer_id text NOT NULL,
    event_type text NOT NULL,
    stream text NOT NULL,
    partition_key text NOT NULL,
    tenant_id uuid NULL,
    canonical_fingerprint bytea NOT NULL,
    processing_state text NOT NULL DEFAULT 'claimed',
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at timestamptz NULL,

    CONSTRAINT events_inbox_owner_valid CHECK (
        owner ~ '^urn:omnexa:module:[a-z0-9.-]+$'
        AND char_length(owner) <= 512
    ),
    CONSTRAINT events_inbox_consumer_valid CHECK (
        consumer_id ~ '^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$'
        AND char_length(consumer_id) <= 128
    ),
    CONSTRAINT events_inbox_event_type_valid CHECK (
        event_type ~ '^[a-z][a-z0-9]*(?:\.[a-z][a-z0-9_]*){2,}\.v[1-9][0-9]*$'
        AND char_length(event_type) <= 512
    ),
    CONSTRAINT events_inbox_stream_valid CHECK (
        stream ~ '^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$'
        AND char_length(stream) <= 128
    ),
    CONSTRAINT events_inbox_partition_valid CHECK (
        partition_key ~ '^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$'
        AND char_length(partition_key) <= 128
    ),
    CONSTRAINT events_inbox_fingerprint_size CHECK (octet_length(canonical_fingerprint) = 32),
    CONSTRAINT events_inbox_state_valid CHECK (processing_state IN ('claimed', 'completed')),
    CONSTRAINT events_inbox_completion_consistent CHECK (
        (processing_state = 'claimed' AND completed_at IS NULL)
        OR
        (processing_state = 'completed' AND completed_at IS NOT NULL)
    ),
    CONSTRAINT events_inbox_completion_time_order CHECK (
        completed_at IS NULL OR completed_at >= created_at
    )
);

-- NULL tenant scope is a legitimate non-tenant consumer scope. The expression
-- index makes NULL participate in processing identity without requiring a fake
-- tenant row. Canonical TenantID is UUIDv7, so the all-zero UUID is not a valid
-- trusted tenant identity and is safe only as this internal index sentinel.
CREATE UNIQUE INDEX IF NOT EXISTS consumer_inbox_processing_identity_uq
    ON omnexa_events.consumer_inbox (
        event_id,
        owner,
        consumer_id,
        event_type,
        stream,
        partition_key,
        COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid)
    );

CREATE INDEX IF NOT EXISTS consumer_inbox_completed_scope_idx
    ON omnexa_events.consumer_inbox (
        owner,
        consumer_id,
        tenant_id,
        completed_at,
        event_id
    )
    WHERE processing_state = 'completed';

COMMENT ON TABLE omnexa_events.consumer_inbox IS 'P04.05 consumer-side local idempotency evidence. A completed row proves only one accepted local processing identity completed inside its PostgreSQL transaction; it is not a checkpoint, retry/DLQ record, authorization grant, broker receipt, or end-to-end exactly-once proof.';
COMMENT ON COLUMN omnexa_events.consumer_inbox.canonical_fingerprint IS 'SHA-256 of the validated canonical P04.01 envelope used only to detect conflicting identity/content reuse; payload content is not retained here.';
COMMENT ON COLUMN omnexa_events.consumer_inbox.processing_state IS 'Transactional local application state. claimed must not be committed without the protected mutation and completion transition in the same transaction.';
COMMENT ON COLUMN omnexa_events.consumer_inbox.completed_at IS 'Local transaction completion evidence only; it does not advance a P04.03 checkpoint or prove external side effects.';

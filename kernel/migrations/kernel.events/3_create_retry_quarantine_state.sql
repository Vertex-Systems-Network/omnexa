CREATE SCHEMA IF NOT EXISTS omnexa_events;

CREATE TABLE IF NOT EXISTS omnexa_events.consumer_retry_state (
    event_id uuid NOT NULL,
    owner text NOT NULL,
    consumer_id text NOT NULL,
    event_type text NOT NULL,
    stream text NOT NULL,
    partition_key text NOT NULL,
    tenant_id uuid NULL,
    delivery_position bigint NOT NULL,
    canonical_fingerprint bytea NOT NULL,

    policy_id text NOT NULL,
    policy_version integer NOT NULL,
    max_attempts integer NOT NULL,
    initial_backoff_ms bigint NOT NULL,
    max_backoff_ms bigint NOT NULL,

    attempts_consumed integer NOT NULL,
    retry_state text NOT NULL,
    next_eligible_at timestamptz NULL,
    terminal_reason text NULL,

    failure_code text NULL,
    failure_category text NULL,
    failure_retryable boolean NULL,

    claim_token uuid NULL,
    claim_expires_at timestamptz NULL,
    revision bigint NOT NULL DEFAULT 1,

    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    quarantined_at timestamptz NULL,
    resolved_at timestamptz NULL,

    CONSTRAINT events_retry_owner_valid CHECK (
        owner ~ '^urn:omnexa:module:[a-z0-9.-]+$'
        AND char_length(owner) <= 512
    ),
    CONSTRAINT events_retry_consumer_valid CHECK (
        consumer_id ~ '^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$'
        AND char_length(consumer_id) <= 128
    ),
    CONSTRAINT events_retry_event_type_valid CHECK (
        event_type ~ '^[a-z][a-z0-9]*(?:\.[a-z][a-z0-9_]*){2,}\.v[1-9][0-9]*$'
        AND char_length(event_type) <= 512
    ),
    CONSTRAINT events_retry_stream_valid CHECK (
        stream ~ '^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$'
        AND char_length(stream) <= 128
    ),
    CONSTRAINT events_retry_partition_valid CHECK (
        partition_key ~ '^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$'
        AND char_length(partition_key) <= 128
    ),
    CONSTRAINT events_retry_position_positive CHECK (delivery_position > 0),
    CONSTRAINT events_retry_fingerprint_size CHECK (octet_length(canonical_fingerprint) = 32),
    CONSTRAINT events_retry_policy_id_valid CHECK (
        policy_id ~ '^[a-z][a-z0-9._-]*$'
        AND char_length(policy_id) <= 128
    ),
    CONSTRAINT events_retry_policy_version_positive CHECK (policy_version > 0),
    CONSTRAINT events_retry_attempt_policy_bounded CHECK (
        max_attempts BETWEEN 1 AND 100
        AND attempts_consumed BETWEEN 1 AND max_attempts
    ),
    CONSTRAINT events_retry_backoff_bounded CHECK (
        initial_backoff_ms > 0
        AND max_backoff_ms >= initial_backoff_ms
        AND max_backoff_ms <= 86400000
    ),
    CONSTRAINT events_retry_state_valid CHECK (
        retry_state IN ('retry_scheduled', 'quarantined', 'resolved')
    ),
    CONSTRAINT events_retry_revision_positive CHECK (revision > 0),
    CONSTRAINT events_retry_failure_evidence_consistent CHECK (
        (failure_code IS NULL AND failure_category IS NULL AND failure_retryable IS NULL)
        OR
        (
            failure_code ~ '^[a-z][a-z0-9_]*(?:\.[a-z][a-z0-9_]*)+$'
            AND char_length(failure_code) <= 512
            AND failure_category IN (
                'validation', 'authentication', 'authorization', 'not_found',
                'conflict', 'rate_limit', 'dependency', 'timeout', 'unavailable',
                'invariant', 'internal'
            )
            AND failure_retryable IS NOT NULL
        )
    ),
    CONSTRAINT events_retry_claim_consistent CHECK (
        (claim_token IS NULL AND claim_expires_at IS NULL)
        OR
        (
            retry_state = 'retry_scheduled'
            AND claim_token IS NOT NULL
            AND claim_expires_at IS NOT NULL
            AND next_eligible_at IS NOT NULL
            AND claim_expires_at > next_eligible_at
        )
    ),
    CONSTRAINT events_retry_state_evidence_consistent CHECK (
        (
            retry_state = 'retry_scheduled'
            AND attempts_consumed < max_attempts
            AND next_eligible_at IS NOT NULL
            AND terminal_reason IS NULL
            AND quarantined_at IS NULL
            AND resolved_at IS NULL
        )
        OR
        (
            retry_state = 'quarantined'
            AND next_eligible_at IS NULL
            AND terminal_reason IN ('non_retryable', 'unsafe_operation', 'attempts_exhausted', 'unknown_failure')
            AND claim_token IS NULL
            AND claim_expires_at IS NULL
            AND quarantined_at IS NOT NULL
            AND resolved_at IS NULL
        )
        OR
        (
            retry_state = 'resolved'
            AND next_eligible_at IS NULL
            AND terminal_reason IS NULL
            AND claim_token IS NULL
            AND claim_expires_at IS NULL
            AND resolved_at IS NOT NULL
        )
    ),
    CONSTRAINT events_retry_timestamp_order CHECK (
        updated_at >= created_at
        AND (quarantined_at IS NULL OR quarantined_at >= created_at)
        AND (resolved_at IS NULL OR resolved_at >= created_at)
    )
);

-- Processing identity is scoped to the exact durable consumer and ordered
-- delivery. NULL tenant is a legitimate non-tenant scope; the all-zero UUID is
-- only an internal index sentinel and is not a trusted TenantID.
CREATE UNIQUE INDEX IF NOT EXISTS consumer_retry_processing_identity_uq
    ON omnexa_events.consumer_retry_state (
        event_id,
        owner,
        consumer_id,
        event_type,
        stream,
        partition_key,
        COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid),
        delivery_position
    );

CREATE INDEX IF NOT EXISTS consumer_retry_due_idx
    ON omnexa_events.consumer_retry_state (
        next_eligible_at,
        owner,
        consumer_id,
        tenant_id,
        delivery_position
    )
    WHERE retry_state = 'retry_scheduled' AND claim_token IS NULL;

CREATE INDEX IF NOT EXISTS consumer_retry_claim_expiry_idx
    ON omnexa_events.consumer_retry_state (claim_expires_at, owner, consumer_id, tenant_id)
    WHERE retry_state = 'retry_scheduled' AND claim_token IS NOT NULL;

CREATE INDEX IF NOT EXISTS consumer_retry_quarantine_idx
    ON omnexa_events.consumer_retry_state (
        owner,
        consumer_id,
        tenant_id,
        quarantined_at,
        delivery_position
    )
    WHERE retry_state = 'quarantined';

COMMENT ON TABLE omnexa_events.consumer_retry_state IS 'P04.06 retry scheduling/quarantine evidence for one exact durable consumer delivery. It is not checkpoint progress, inbox completion, authorization, a business-success record, or broker DLQ authority.';
COMMENT ON COLUMN omnexa_events.consumer_retry_state.canonical_fingerprint IS 'SHA-256 canonical-envelope conflict evidence only. Event payload content is not retained in retry state.';
COMMENT ON COLUMN omnexa_events.consumer_retry_state.failure_code IS 'Bounded safe P01.03 machine failure code when structured evidence exists; raw error/provider/database text is forbidden.';
COMMENT ON COLUMN omnexa_events.consumer_retry_state.claim_token IS 'Reserved bounded one-at-a-time claim identity for a later P04.06 CAS/lease adapter; migration 3 alone grants no claim authority.';
COMMENT ON COLUMN omnexa_events.consumer_retry_state.revision IS 'Monotonic CAS revision reserved for later P04.06 persistence transitions.';

CREATE SCHEMA IF NOT EXISTS omnexa_events;

CREATE TABLE IF NOT EXISTS omnexa_events.transactional_outbox (
    event_id uuid PRIMARY KEY,
    owner text NOT NULL,
    tenant_id uuid NULL,
    envelope jsonb NOT NULL,
    publication_state text NOT NULL DEFAULT 'pending',
    revision bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at timestamptz NULL,

    CONSTRAINT events_outbox_owner_valid CHECK (
        owner ~ '^urn:omnexa:module:[a-z0-9.-]+$'
        AND char_length(owner) <= 512
    ),
    CONSTRAINT events_outbox_envelope_object CHECK (jsonb_typeof(envelope) = 'object'),
    CONSTRAINT events_outbox_envelope_size CHECK (octet_length(envelope::text) <= 131072),
    CONSTRAINT events_outbox_event_identity_matches CHECK (
        envelope ? 'id'
        AND (envelope ->> 'id')::uuid = event_id
    ),
    CONSTRAINT events_outbox_owner_identity_matches CHECK (
        envelope ? 'source'
        AND envelope ->> 'source' = owner
    ),
    CONSTRAINT events_outbox_tenant_identity_matches CHECK (
        CASE
            WHEN tenant_id IS NULL THEN NOT (envelope ? 'tenantid') OR nullif(envelope ->> 'tenantid', '') IS NULL
            ELSE envelope ? 'tenantid' AND (envelope ->> 'tenantid')::uuid = tenant_id
        END
    ),
    CONSTRAINT events_outbox_state_valid CHECK (publication_state IN ('pending', 'published')),
    CONSTRAINT events_outbox_revision_valid CHECK (revision > 0),
    CONSTRAINT events_outbox_publication_time_consistent CHECK (
        (publication_state = 'pending' AND published_at IS NULL)
        OR
        (publication_state = 'published' AND published_at IS NOT NULL)
    ),
    CONSTRAINT events_outbox_publication_time_order CHECK (
        published_at IS NULL OR published_at >= created_at
    )
);

CREATE INDEX IF NOT EXISTS transactional_outbox_pending_owner_tenant_idx
    ON omnexa_events.transactional_outbox (owner, tenant_id, created_at, event_id)
    WHERE publication_state = 'pending';

COMMENT ON SCHEMA omnexa_events IS 'kernel.events authoritative event reliability schema';
COMMENT ON TABLE omnexa_events.transactional_outbox IS 'Producer-side transactional outbox. Published records are retained; publication state does not prove consumer receipt or exactly-once mutation.';
COMMENT ON COLUMN omnexa_events.transactional_outbox.envelope IS 'Validated canonical P04.01 envelope preserved for faithful relay; must not be copied into logs or diagnostics.';
COMMENT ON COLUMN omnexa_events.transactional_outbox.revision IS 'Optimistic producer-side publication-state revision; not event identity, ordering authority, retry count, or consumer checkpoint.';
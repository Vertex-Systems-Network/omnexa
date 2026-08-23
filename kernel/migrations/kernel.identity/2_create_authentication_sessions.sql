CREATE TABLE omnexa_identity.password_credentials (
    principal_id uuid PRIMARY KEY
        REFERENCES omnexa_identity.users(principal_id) ON DELETE RESTRICT,
    password_hash text NOT NULL,
    credential_version bigint NOT NULL DEFAULT 1 CHECK (credential_version > 0),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT identity_password_hash_shape CHECK (
        length(password_hash) BETWEEN 64 AND 512
        AND password_hash LIKE '$pbkdf2-sha256$i=%'
    ),
    CONSTRAINT identity_password_timestamp_order CHECK (updated_at >= created_at)
);

COMMENT ON COLUMN omnexa_identity.password_credentials.password_hash IS
    'RESTRICTED AUTH_SECRET-derived one-way PBKDF2 representation; never log, trace, audit or expose.';

CREATE TABLE omnexa_identity.sessions (
    id uuid PRIMARY KEY,
    principal_id uuid NOT NULL
        REFERENCES omnexa_identity.users(principal_id) ON DELETE RESTRICT,
    credential_version bigint NOT NULL CHECK (credential_version > 0),
    device_label text NOT NULL DEFAULT '',
    tenant_context_hint uuid,
    organization_context_hint uuid,
    created_at timestamptz NOT NULL,
    refreshed_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    CONSTRAINT identity_session_id_uuidv7 CHECK (substring(id::text, 15, 1) = '7'),
    CONSTRAINT identity_session_device_label CHECK (
        length(device_label) <= 128
        AND device_label !~ E'[\\r\\n]'
    ),
    CONSTRAINT identity_session_context_shape CHECK (
        organization_context_hint IS NULL OR tenant_context_hint IS NOT NULL
    ),
    CONSTRAINT identity_session_timestamp_order CHECK (
        refreshed_at >= created_at
        AND expires_at > created_at
        AND (revoked_at IS NULL OR revoked_at >= created_at)
    )
);

COMMENT ON COLUMN omnexa_identity.sessions.tenant_context_hint IS
    'Non-authorizing P02.02 context reference; current relationship must be reauthorized before use.';
COMMENT ON COLUMN omnexa_identity.sessions.organization_context_hint IS
    'Non-authorizing P02.03 context reference; current relationship must be reauthorized before use.';

CREATE INDEX identity_sessions_principal_created_idx
    ON omnexa_identity.sessions (principal_id, created_at DESC);

CREATE TABLE omnexa_identity.access_credentials (
    session_id uuid NOT NULL
        REFERENCES omnexa_identity.sessions(id) ON DELETE RESTRICT,
    secret_digest bytea PRIMARY KEY,
    created_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    CONSTRAINT identity_access_digest_shape CHECK (octet_length(secret_digest) = 32),
    CONSTRAINT identity_access_timestamp_order CHECK (
        expires_at > created_at
        AND (revoked_at IS NULL OR revoked_at >= created_at)
    )
);

COMMENT ON COLUMN omnexa_identity.access_credentials.secret_digest IS
    'SHA-256 lookup digest of a high-entropy opaque access secret; raw secret is never persisted.';

CREATE INDEX identity_access_session_idx
    ON omnexa_identity.access_credentials (session_id, expires_at DESC);

CREATE TABLE omnexa_identity.refresh_credentials (
    session_id uuid NOT NULL
        REFERENCES omnexa_identity.sessions(id) ON DELETE RESTRICT,
    secret_digest bytea PRIMARY KEY,
    created_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL,
    consumed_at timestamptz,
    revoked_at timestamptz,
    CONSTRAINT identity_refresh_digest_shape CHECK (octet_length(secret_digest) = 32),
    CONSTRAINT identity_refresh_timestamp_order CHECK (
        expires_at > created_at
        AND (consumed_at IS NULL OR consumed_at >= created_at)
        AND (revoked_at IS NULL OR revoked_at >= created_at)
    )
);

COMMENT ON COLUMN omnexa_identity.refresh_credentials.secret_digest IS
    'SHA-256 lookup digest of a RESTRICTED high-entropy opaque refresh secret; raw secret is never persisted.';

CREATE INDEX identity_refresh_session_idx
    ON omnexa_identity.refresh_credentials (session_id, expires_at DESC);

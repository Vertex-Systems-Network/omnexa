ALTER TABLE omnexa_identity.sessions
    ADD CONSTRAINT identity_sessions_id_principal_unique UNIQUE (id, principal_id);

CREATE TABLE omnexa_identity.mfa_factors (
    id uuid PRIMARY KEY,
    principal_id uuid NOT NULL
        REFERENCES omnexa_identity.users(principal_id) ON DELETE RESTRICT,
    factor_type text NOT NULL CHECK (factor_type IN ('passkey')),
    label text NOT NULL,
    lifecycle_state text NOT NULL CHECK (lifecycle_state IN ('pending', 'active', 'revoked')),
    created_at timestamptz NOT NULL,
    verified_at timestamptz,
    revoked_at timestamptz,
    CONSTRAINT identity_mfa_factor_id_uuidv7 CHECK (substring(id::text, 15, 1) = '7'),
    CONSTRAINT identity_mfa_factor_label CHECK (
        length(label) BETWEEN 1 AND 128
        AND label !~ E'[\\r\\n]'
    ),
    CONSTRAINT identity_mfa_factor_lifecycle CHECK (
        (lifecycle_state = 'pending' AND verified_at IS NULL AND revoked_at IS NULL)
        OR (lifecycle_state = 'active' AND verified_at IS NOT NULL AND revoked_at IS NULL)
        OR (lifecycle_state = 'revoked' AND verified_at IS NOT NULL AND revoked_at IS NOT NULL)
    ),
    CONSTRAINT identity_mfa_factor_timestamp_order CHECK (
        (verified_at IS NULL OR verified_at >= created_at)
        AND (revoked_at IS NULL OR (verified_at IS NOT NULL AND revoked_at >= verified_at))
    ),
    UNIQUE (id, principal_id)
);

COMMENT ON TABLE omnexa_identity.mfa_factors IS
    'P02.07 human-user strong-authentication factor inventory; contains no factor secret material.';

CREATE INDEX identity_mfa_factors_principal_state_idx
    ON omnexa_identity.mfa_factors (principal_id, lifecycle_state, created_at);

CREATE TABLE omnexa_identity.passkey_credentials (
    factor_id uuid PRIMARY KEY
        REFERENCES omnexa_identity.mfa_factors(id) ON DELETE RESTRICT,
    credential_id bytea NOT NULL UNIQUE,
    public_key bytea NOT NULL,
    counter_supported boolean NOT NULL,
    sign_count bigint NOT NULL DEFAULT 0 CHECK (sign_count BETWEEN 0 AND 4294967295),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT identity_passkey_credential_id_shape CHECK (octet_length(credential_id) BETWEEN 1 AND 1024),
    CONSTRAINT identity_passkey_public_key_shape CHECK (octet_length(public_key) BETWEEN 1 AND 4096),
    CONSTRAINT identity_passkey_counter_shape CHECK (counter_supported OR sign_count = 0),
    CONSTRAINT identity_passkey_timestamp_order CHECK (updated_at >= created_at)
);

COMMENT ON COLUMN omnexa_identity.passkey_credentials.public_key IS
    'Public passkey/WebAuthn credential verification material only; authenticator private keys never enter Omnexa.';

CREATE TABLE omnexa_identity.authentication_challenges (
    id uuid PRIMARY KEY,
    principal_id uuid NOT NULL,
    session_id uuid NOT NULL,
    factor_id uuid NOT NULL,
    purpose text NOT NULL CHECK (purpose IN ('passkey_enrollment', 'passkey_assertion')),
    secret_digest bytea NOT NULL UNIQUE,
    created_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL,
    consumed_at timestamptz,
    revoked_at timestamptz,
    CONSTRAINT identity_auth_challenge_id_uuidv7 CHECK (substring(id::text, 15, 1) = '7'),
    CONSTRAINT identity_auth_challenge_digest_shape CHECK (octet_length(secret_digest) = 32),
    CONSTRAINT identity_auth_challenge_timestamp_order CHECK (
        expires_at > created_at
        AND (consumed_at IS NULL OR (consumed_at >= created_at AND consumed_at < expires_at))
        AND (revoked_at IS NULL OR revoked_at >= created_at)
    ),
    FOREIGN KEY (session_id, principal_id)
        REFERENCES omnexa_identity.sessions(id, principal_id) ON DELETE RESTRICT,
    FOREIGN KEY (factor_id, principal_id)
        REFERENCES omnexa_identity.mfa_factors(id, principal_id) ON DELETE RESTRICT
);

COMMENT ON COLUMN omnexa_identity.authentication_challenges.secret_digest IS
    'SHA-256 digest of one high-entropy RESTRICTED challenge; raw challenge is never persisted.';

CREATE INDEX identity_auth_challenges_session_expiry_idx
    ON omnexa_identity.authentication_challenges (principal_id, session_id, expires_at DESC);

CREATE TABLE omnexa_identity.recovery_code_sets (
    id uuid PRIMARY KEY,
    principal_id uuid NOT NULL
        REFERENCES omnexa_identity.users(principal_id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL,
    revoked_at timestamptz,
    CONSTRAINT identity_recovery_set_id_uuidv7 CHECK (substring(id::text, 15, 1) = '7'),
    CONSTRAINT identity_recovery_set_timestamp_order CHECK (revoked_at IS NULL OR revoked_at >= created_at)
);

CREATE UNIQUE INDEX identity_recovery_one_active_set_idx
    ON omnexa_identity.recovery_code_sets (principal_id)
    WHERE revoked_at IS NULL;

CREATE TABLE omnexa_identity.recovery_codes (
    set_id uuid NOT NULL
        REFERENCES omnexa_identity.recovery_code_sets(id) ON DELETE RESTRICT,
    code_index integer NOT NULL CHECK (code_index BETWEEN 0 AND 15),
    secret_digest bytea NOT NULL UNIQUE,
    created_at timestamptz NOT NULL,
    consumed_at timestamptz,
    PRIMARY KEY (set_id, code_index),
    CONSTRAINT identity_recovery_digest_shape CHECK (octet_length(secret_digest) = 32),
    CONSTRAINT identity_recovery_timestamp_order CHECK (consumed_at IS NULL OR consumed_at >= created_at)
);

COMMENT ON COLUMN omnexa_identity.recovery_codes.secret_digest IS
    'SHA-256 digest of one high-entropy RESTRICTED recovery code; raw recovery material is one-time issuance only.';

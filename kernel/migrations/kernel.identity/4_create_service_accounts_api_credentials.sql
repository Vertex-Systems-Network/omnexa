CREATE TABLE omnexa_identity.service_accounts (
    principal_id uuid PRIMARY KEY
        REFERENCES omnexa_identity.principals(id) ON DELETE RESTRICT,
    name text NOT NULL,
    tenant_id uuid NOT NULL,
    organization_id uuid,
    CONSTRAINT identity_service_account_name_valid CHECK (
        char_length(btrim(name)) BETWEEN 1 AND 128
        AND name !~ E'[\r\n]'
    ),
    CONSTRAINT identity_service_account_tenant_uuidv7 CHECK (
        substring(tenant_id::text, 15, 1) = '7'
    ),
    CONSTRAINT identity_service_account_organization_uuidv7 CHECK (
        organization_id IS NULL OR substring(organization_id::text, 15, 1) = '7'
    )
);

COMMENT ON TABLE omnexa_identity.service_accounts IS
    'P02.08 non-human service principals with an exact tenant/organization authentication binding; binding is not authorization authority.';
COMMENT ON COLUMN omnexa_identity.service_accounts.tenant_id IS
    'Opaque UUIDv7 tenant binding ceiling. Current kernel.authorization state still decides permission authority.';
COMMENT ON COLUMN omnexa_identity.service_accounts.organization_id IS
    'Optional exact organization binding ceiling; NULL means tenant scope and never implies child-organization authority.';

CREATE TABLE omnexa_identity.api_credentials (
    id uuid PRIMARY KEY,
    service_account_id uuid NOT NULL
        REFERENCES omnexa_identity.service_accounts(principal_id) ON DELETE RESTRICT,
    secret_digest bytea NOT NULL UNIQUE,
    created_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL,
    last_used_at timestamptz,
    superseded_at timestamptz,
    revoked_at timestamptz,
    CONSTRAINT identity_api_credential_id_uuidv7 CHECK (
        substring(id::text, 15, 1) = '7'
    ),
    CONSTRAINT identity_api_credential_digest_shape CHECK (
        octet_length(secret_digest) = 32
    ),
    CONSTRAINT identity_api_credential_timestamp_order CHECK (
        expires_at > created_at
        AND (last_used_at IS NULL OR last_used_at >= created_at)
        AND (superseded_at IS NULL OR superseded_at >= created_at)
        AND (revoked_at IS NULL OR revoked_at >= created_at)
    ),
    CONSTRAINT identity_api_credential_terminal_state CHECK (
        NOT (superseded_at IS NOT NULL AND revoked_at IS NOT NULL)
    )
);

CREATE INDEX identity_api_credentials_account_inventory_idx
    ON omnexa_identity.api_credentials (service_account_id, created_at DESC);

CREATE UNIQUE INDEX identity_api_credentials_one_active_digest_idx
    ON omnexa_identity.api_credentials (secret_digest)
    WHERE revoked_at IS NULL AND superseded_at IS NULL;

COMMENT ON TABLE omnexa_identity.api_credentials IS
    'P02.08 rotatable API credential inventory. Raw credentials are one-time issuance material and never persisted.';
COMMENT ON COLUMN omnexa_identity.api_credentials.secret_digest IS
    'SHA-256 verifier digest of high-entropy RESTRICTED API credential material; raw credential values never enter persistence.';

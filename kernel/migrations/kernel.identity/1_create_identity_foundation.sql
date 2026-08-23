CREATE SCHEMA IF NOT EXISTS omnexa_identity;

CREATE TABLE IF NOT EXISTS omnexa_identity.principals (
    id uuid PRIMARY KEY,
    principal_type text NOT NULL,
    lifecycle_state text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT identity_principal_type_valid CHECK (
        principal_type IN (
            'human_user',
            'service_account',
            'workload',
            'device',
            'integration',
            'support_operator',
            'ai_agent'
        )
    ),
    CONSTRAINT identity_lifecycle_state_valid CHECK (
        lifecycle_state IN ('provisioned', 'active', 'suspended', 'disabled')
    ),
    CONSTRAINT identity_principal_time_valid CHECK (updated_at >= created_at)
);

CREATE TABLE IF NOT EXISTS omnexa_identity.users (
    principal_id uuid PRIMARY KEY REFERENCES omnexa_identity.principals(id) ON DELETE RESTRICT,
    primary_email text NOT NULL,
    CONSTRAINT identity_user_email_length CHECK (char_length(primary_email) BETWEEN 3 AND 320)
);

COMMENT ON SCHEMA omnexa_identity IS 'kernel.identity authoritative identity schema';
COMMENT ON TABLE omnexa_identity.principals IS 'Security principal identity and lifecycle state; no tenant or authorization authority';
COMMENT ON TABLE omnexa_identity.users IS 'Human User identity attributes; primary_email is CONFIDENTIAL PII and is not a credential';

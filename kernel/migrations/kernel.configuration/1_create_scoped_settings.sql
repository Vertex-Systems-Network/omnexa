CREATE SCHEMA IF NOT EXISTS omnexa_configuration;

CREATE TABLE IF NOT EXISTS omnexa_configuration.setting_overrides (
    tenant_id uuid NOT NULL REFERENCES omnexa_tenancy.tenants(id) ON DELETE RESTRICT,
    organization_id uuid NULL,
    setting_key text NOT NULL,
    value_kind text NOT NULL,
    value_text text NOT NULL,
    revision bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT configuration_override_organization_same_tenant
        FOREIGN KEY (tenant_id, organization_id)
        REFERENCES omnexa_organization.nodes (tenant_id, id)
        ON DELETE RESTRICT,
    CONSTRAINT configuration_setting_key_valid CHECK (
        setting_key ~ '^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$'
    ),
    CONSTRAINT configuration_value_kind_valid CHECK (
        value_kind IN ('bool', 'string', 'int', 'duration')
    ),
    CONSTRAINT configuration_value_size_valid CHECK (char_length(value_text) <= 4096),
    CONSTRAINT configuration_revision_valid CHECK (revision > 0),
    CONSTRAINT configuration_override_time_valid CHECK (updated_at >= created_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS configuration_one_tenant_override_per_key
    ON omnexa_configuration.setting_overrides (tenant_id, setting_key)
    WHERE organization_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS configuration_one_organization_override_per_key
    ON omnexa_configuration.setting_overrides (tenant_id, organization_id, setting_key)
    WHERE organization_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS configuration_override_scope_lookup
    ON omnexa_configuration.setting_overrides (tenant_id, organization_id, setting_key);

COMMENT ON SCHEMA omnexa_configuration IS 'kernel.configuration authoritative scoped-setting schema';
COMMENT ON TABLE omnexa_configuration.setting_overrides IS 'Exact tenant/approved-organization runtime overrides; no global, user, authorization, entitlement or secret-management state';
COMMENT ON COLUMN omnexa_configuration.setting_overrides.value_text IS 'Typed generic runtime setting value; P02.09 policy rejects RESTRICTED/secret output surfaces';

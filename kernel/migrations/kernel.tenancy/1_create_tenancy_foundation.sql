CREATE SCHEMA IF NOT EXISTS omnexa_tenancy;

CREATE TABLE IF NOT EXISTS omnexa_tenancy.tenants (
    id uuid PRIMARY KEY,
    lifecycle_state text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT tenancy_lifecycle_state_valid CHECK (
        lifecycle_state IN ('provisioned', 'active', 'suspended', 'disabled')
    ),
    CONSTRAINT tenancy_tenant_time_valid CHECK (updated_at >= created_at)
);

CREATE TABLE IF NOT EXISTS omnexa_tenancy.tenant_memberships (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES omnexa_tenancy.tenants(id) ON DELETE RESTRICT,
    principal_id uuid NOT NULL REFERENCES omnexa_identity.principals(id) ON DELETE RESTRICT,
    relationship_state text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT tenancy_membership_state_valid CHECK (
        relationship_state IN ('active', 'revoked')
    ),
    CONSTRAINT tenancy_membership_time_valid CHECK (updated_at >= created_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS tenancy_one_active_membership_per_principal_tenant
    ON omnexa_tenancy.tenant_memberships (tenant_id, principal_id)
    WHERE relationship_state = 'active';

CREATE INDEX IF NOT EXISTS tenancy_memberships_principal_lookup
    ON omnexa_tenancy.tenant_memberships (principal_id, tenant_id);

COMMENT ON SCHEMA omnexa_tenancy IS 'kernel.tenancy authoritative tenant-isolation schema';
COMMENT ON TABLE omnexa_tenancy.tenants IS 'Primary Tenant isolation boundary and lifecycle state';
COMMENT ON TABLE omnexa_tenancy.tenant_memberships IS 'Minimal Tenant/User relationship; no organization, role, permission or session semantics';
COMMENT ON COLUMN omnexa_tenancy.tenant_memberships.principal_id IS 'Reference to kernel.identity principal identity; kernel.tenancy does not own or mutate identity state';

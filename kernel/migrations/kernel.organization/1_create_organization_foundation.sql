CREATE SCHEMA IF NOT EXISTS omnexa_organization;

CREATE TABLE IF NOT EXISTS omnexa_organization.nodes (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES omnexa_tenancy.tenants(id) ON DELETE RESTRICT,
    node_kind text NOT NULL,
    parent_id uuid NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT organization_node_kind_valid CHECK (
        node_kind IN ('organization', 'legal_entity', 'business_unit', 'branch', 'team', 'location')
    ),
    CONSTRAINT organization_root_shape_valid CHECK (
        (node_kind = 'organization' AND parent_id IS NULL)
        OR
        (node_kind <> 'organization' AND parent_id IS NOT NULL)
    ),
    CONSTRAINT organization_node_time_valid CHECK (updated_at >= created_at),
    CONSTRAINT organization_node_tenant_identity UNIQUE (tenant_id, id),
    CONSTRAINT organization_parent_same_tenant
        FOREIGN KEY (tenant_id, parent_id)
        REFERENCES omnexa_organization.nodes (tenant_id, id)
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS organization_nodes_parent_lookup
    ON omnexa_organization.nodes (tenant_id, parent_id, id);

CREATE TABLE IF NOT EXISTS omnexa_organization.scoped_memberships (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES omnexa_tenancy.tenants(id) ON DELETE RESTRICT,
    scope_id uuid NOT NULL,
    principal_id uuid NOT NULL REFERENCES omnexa_identity.principals(id) ON DELETE RESTRICT,
    relationship_state text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT organization_membership_scope_same_tenant
        FOREIGN KEY (tenant_id, scope_id)
        REFERENCES omnexa_organization.nodes (tenant_id, id)
        ON DELETE RESTRICT,
    CONSTRAINT organization_membership_state_valid CHECK (
        relationship_state IN ('active', 'revoked')
    ),
    CONSTRAINT organization_membership_time_valid CHECK (updated_at >= created_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS organization_one_active_membership_per_scope_principal
    ON omnexa_organization.scoped_memberships (tenant_id, scope_id, principal_id)
    WHERE relationship_state = 'active';

CREATE INDEX IF NOT EXISTS organization_memberships_principal_lookup
    ON omnexa_organization.scoped_memberships (principal_id, tenant_id, scope_id);

COMMENT ON SCHEMA omnexa_organization IS 'kernel.organization authoritative tenant-contained organization hierarchy schema';
COMMENT ON TABLE omnexa_organization.nodes IS 'Organization, Legal Entity, Business Unit, Branch, Team and Location access/operating scopes; not business Party Organization records';
COMMENT ON TABLE omnexa_organization.scoped_memberships IS 'Minimal User-to-organization-scope relationships; no roles, permissions, policies, sessions or employment semantics';
COMMENT ON COLUMN omnexa_organization.nodes.tenant_id IS 'Trusted Tenant isolation boundary inherited from kernel.tenancy';
COMMENT ON COLUMN omnexa_organization.scoped_memberships.principal_id IS 'Reference to kernel.identity User identity; kernel.organization does not own or mutate identity state';

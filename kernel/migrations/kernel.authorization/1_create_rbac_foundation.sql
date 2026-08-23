CREATE SCHEMA IF NOT EXISTS omnexa_authorization;

CREATE TABLE IF NOT EXISTS omnexa_authorization.permissions (
    permission_id text PRIMARY KEY,
    created_at timestamptz NOT NULL,
    CONSTRAINT authorization_permission_id_valid CHECK (
        permission_id ~ '^[a-z][a-z0-9]*(\.[a-z][a-z0-9_]*){2,7}$'
    )
);

INSERT INTO omnexa_authorization.permissions (permission_id, created_at)
VALUES
    ('authorization.role.read', TIMESTAMPTZ '2026-08-23 00:00:00+00'),
    ('authorization.role.manage', TIMESTAMPTZ '2026-08-23 00:00:00+00'),
    ('authorization.assignment.read', TIMESTAMPTZ '2026-08-23 00:00:00+00'),
    ('authorization.assignment.manage', TIMESTAMPTZ '2026-08-23 00:00:00+00')
ON CONFLICT (permission_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS omnexa_authorization.roles (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES omnexa_tenancy.tenants(id) ON DELETE RESTRICT,
    organization_id uuid NULL,
    name text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT authorization_role_organization_same_tenant
        FOREIGN KEY (tenant_id, organization_id)
        REFERENCES omnexa_organization.nodes (tenant_id, id)
        ON DELETE RESTRICT,
    CONSTRAINT authorization_role_name_valid CHECK (
        char_length(btrim(name)) BETWEEN 1 AND 128
        AND name !~ E'[\r\n]'
    ),
    CONSTRAINT authorization_role_time_valid CHECK (updated_at >= created_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS authorization_tenant_role_name_unique
    ON omnexa_authorization.roles (tenant_id, name)
    WHERE organization_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS authorization_organization_role_name_unique
    ON omnexa_authorization.roles (tenant_id, organization_id, name)
    WHERE organization_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS omnexa_authorization.role_permissions (
    role_id uuid NOT NULL REFERENCES omnexa_authorization.roles(id) ON DELETE CASCADE,
    permission_id text NOT NULL REFERENCES omnexa_authorization.permissions(permission_id) ON DELETE RESTRICT,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS omnexa_authorization.role_assignments (
    id uuid PRIMARY KEY,
    role_id uuid NOT NULL REFERENCES omnexa_authorization.roles(id) ON DELETE RESTRICT,
    principal_id uuid NOT NULL REFERENCES omnexa_identity.users(principal_id) ON DELETE RESTRICT,
    assignment_state text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT authorization_assignment_state_valid CHECK (
        assignment_state IN ('active', 'revoked')
    ),
    CONSTRAINT authorization_assignment_time_valid CHECK (updated_at >= created_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS authorization_one_active_role_assignment
    ON omnexa_authorization.role_assignments (role_id, principal_id)
    WHERE assignment_state = 'active';

CREATE INDEX IF NOT EXISTS authorization_assignment_principal_lookup
    ON omnexa_authorization.role_assignments (principal_id, role_id)
    WHERE assignment_state = 'active';

COMMENT ON SCHEMA omnexa_authorization IS 'kernel.authorization authoritative direct-RBAC schema';
COMMENT ON TABLE omnexa_authorization.permissions IS 'Stable capability-oriented P02.05 permission reference data; P03 module registration is not implemented here';
COMMENT ON TABLE omnexa_authorization.roles IS 'Exact tenant/organization scoped permission compositions; role names confer no bypass authority';
COMMENT ON TABLE omnexa_authorization.role_permissions IS 'Direct deterministic role permission composition only; no relationship or contextual policy semantics';
COMMENT ON TABLE omnexa_authorization.role_assignments IS 'Direct human User role grants with terminal revocation lifecycle';

ALTER TABLE omnexa_authorization.permissions
    ADD COLUMN source_kind text NOT NULL DEFAULT 'kernel',
    ADD COLUMN module_id text NULL,
    ADD COLUMN owner_name text NULL,
    ADD COLUMN capability_ref text NULL,
    ADD COLUMN available boolean NOT NULL DEFAULT TRUE,
    ADD COLUMN updated_at timestamptz NULL;

UPDATE omnexa_authorization.permissions
SET updated_at = created_at
WHERE updated_at IS NULL;

ALTER TABLE omnexa_authorization.permissions
    ALTER COLUMN updated_at SET NOT NULL,
    ADD CONSTRAINT authorization_permission_source_kind_valid
        CHECK (source_kind IN ('kernel', 'module')),
    ADD CONSTRAINT authorization_permission_owner_shape_valid
        CHECK (
            (source_kind = 'kernel' AND module_id IS NULL AND owner_name IS NULL)
            OR
            (
                source_kind = 'module'
                AND module_id ~ '^omnexa\.[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$'
                AND owner_name ~ '^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$'
            )
        ),
    ADD CONSTRAINT authorization_permission_capability_ref_valid
        CHECK (
            capability_ref IS NULL
            OR (
                char_length(capability_ref) BETWEEN 1 AND 128
                AND capability_ref !~ '[[:space:][:cntrl:]]'
            )
        ),
    ADD CONSTRAINT authorization_kernel_permission_available
        CHECK (source_kind <> 'kernel' OR available = TRUE),
    ADD CONSTRAINT authorization_permission_update_time_valid
        CHECK (updated_at >= created_at);

CREATE INDEX authorization_module_permission_lookup
    ON omnexa_authorization.permissions (module_id, available, permission_id)
    WHERE source_kind = 'module';

COMMENT ON TABLE omnexa_authorization.permissions IS 'Authorization-owned stable permission catalog. Kernel permissions remain available; P03.07 module permissions retain owner/module identity and lifecycle availability without creating grants.';
COMMENT ON COLUMN omnexa_authorization.permissions.source_kind IS 'kernel or module reference-data ownership; source kind never grants authority';
COMMENT ON COLUMN omnexa_authorization.permissions.available IS 'Module lifecycle availability precondition only; actual authorization remains deny-by-default role/policy evaluation';
COMMENT ON COLUMN omnexa_authorization.permissions.capability_ref IS 'Optional descriptive P03.06 capability association; never invocation authority';

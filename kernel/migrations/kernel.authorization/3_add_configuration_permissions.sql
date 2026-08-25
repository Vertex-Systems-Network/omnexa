INSERT INTO omnexa_authorization.permissions (permission_id, created_at)
VALUES
    ('configuration.setting.read', TIMESTAMPTZ '2026-08-25 00:00:00+00'),
    ('configuration.setting.manage', TIMESTAMPTZ '2026-08-25 00:00:00+00')
ON CONFLICT (permission_id) DO NOTHING;

COMMENT ON TABLE omnexa_authorization.permissions IS 'Stable capability-oriented kernel permission reference data including P02.09 configuration setting read/manage; permission existence never grants authority';

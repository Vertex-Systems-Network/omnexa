ALTER TABLE omnexa_authorization.role_assignments
    DROP CONSTRAINT role_assignments_principal_id_fkey;

ALTER TABLE omnexa_authorization.role_assignments
    ADD CONSTRAINT authorization_assignment_principal_fk
    FOREIGN KEY (principal_id)
    REFERENCES omnexa_identity.principals(id)
    ON DELETE RESTRICT;

COMMENT ON TABLE omnexa_authorization.role_assignments IS
    'Direct human User or P02.08 Service Account role grants with terminal revocation lifecycle; later principal kinds remain unauthorized by runtime type checks.';

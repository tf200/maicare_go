DROP INDEX IF EXISTS assigned_employee_one_coordinator_per_client_idx;

ALTER TABLE assigned_employee
    ALTER COLUMN role TYPE VARCHAR(100)
    USING role::text;

CREATE UNIQUE INDEX assigned_employee_one_coordinator_per_client_idx
    ON assigned_employee(client_id)
    WHERE role = 'coordinator';

DROP TYPE IF EXISTS client_involved_role_enum;

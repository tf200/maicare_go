CREATE TYPE client_involved_role_enum AS ENUM (
    'coordinator',
    'primary_counselor',
    'secondary_counselor',
    'behavioral_scientist',
    'case_manager',
    'specialist',
    'other'
);

-- Drop partial unique index before altering column type
DROP INDEX IF EXISTS assigned_employee_one_coordinator_per_client_idx;

-- Normalize any existing rows with unexpected roles to 'other'
UPDATE assigned_employee
SET role = 'other'
WHERE role NOT IN (
    'coordinator',
    'primary_counselor',
    'secondary_counselor',
    'behavioral_scientist',
    'case_manager',
    'specialist',
    'other'
);

-- Alter column type to client_involved_role_enum
ALTER TABLE assigned_employee
    ALTER COLUMN role TYPE client_involved_role_enum
    USING role::client_involved_role_enum;

-- Recreate partial unique index for single coordinator per client
CREATE UNIQUE INDEX assigned_employee_one_coordinator_per_client_idx
    ON assigned_employee(client_id)
    WHERE role = 'coordinator';

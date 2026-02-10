ALTER TYPE form_status_enum ADD VALUE IF NOT EXISTS 'in_review';
ALTER TYPE form_status_enum ADD VALUE IF NOT EXISTS 'archived';

CREATE TYPE intake_status_enum AS ENUM ('scheduled', 'in_progress', 'completed', 'cancelled');
CREATE TYPE urgency_level_enum AS ENUM ('low', 'medium', 'high', 'critical');

ALTER TABLE intake_forms
    ADD COLUMN status intake_status_enum NOT NULL DEFAULT 'scheduled',
    ADD COLUMN urgency_level urgency_level_enum NULL;

ALTER TABLE intake_forms
    ALTER COLUMN care_type DROP NOT NULL,
    ALTER COLUMN self_sufficiency DROP NOT NULL,
    ALTER COLUMN intake_conclusion DROP NOT NULL;

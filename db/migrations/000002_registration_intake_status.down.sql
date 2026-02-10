ALTER TABLE intake_forms
    ALTER COLUMN care_type SET NOT NULL,
    ALTER COLUMN self_sufficiency SET NOT NULL,
    ALTER COLUMN intake_conclusion SET NOT NULL;

ALTER TABLE intake_forms
    DROP COLUMN IF EXISTS urgency_level,
    DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS urgency_level_enum;
DROP TYPE IF EXISTS intake_status_enum;

-- NOTE: form_status_enum values cannot be removed in a down migration.

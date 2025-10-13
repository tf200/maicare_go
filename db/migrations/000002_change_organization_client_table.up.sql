ALTER TABLE client_details
DROP COLUMN organisation;

ALTER TABLE client_details
ADD COLUMN organization_id BIGINT NULL REFERENCES organisations(id) ON DELETE SET NULL;
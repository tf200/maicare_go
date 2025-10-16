ALTER TABLE client_details
DROP COLUMN organization_id;

ALTER TABLE client_details
ADD COLUMN organisation VARCHAR(100) NULL;
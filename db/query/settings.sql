-- name: GetAppOrganizationProfile :one
SELECT *
FROM app_organization_profile
WHERE singleton = TRUE
LIMIT 1;

-- name: UpdateAppOrganizationProfile :one
UPDATE app_organization_profile
SET
    name = $1,
    default_timezone = $2,
    email = $3,
    phone_number = $4,
    website = $5,
    hq_street = $6,
    hq_house_number = $7,
    hq_house_number_addition = $8,
    hq_postal_code = $9,
    hq_city = $10,
    updated_at = CURRENT_TIMESTAMP
WHERE singleton = TRUE
RETURNING *;

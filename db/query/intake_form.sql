


-- name: CreateIntakeFormToken :one
INSERT INTO intake_form_tokens(
    token,
    expires_at
) VALUES (
    $1,
    $2
)

RETURNING *;



-- name: GetIntakeFormToken :one
SELECT * FROM intake_form_tokens WHERE token = $1;

-- name: RevokedIntakeFormToken :one
UPDATE intake_form_tokens
SET
    is_revoked = true
WHERE token = $1
RETURNING *;


-- name: CreateIntakeForm :one
INSERT INTO intake_forms (
    intake_form_token,
    first_name,
    last_name,
    date_of_birth,
    phone_number,
    gender,
    place_of_birth,
    representative_first_name,
    representative_last_name,
    representative_phone_number,
    representative_email,
    representative_relationship,
    representative_address,
    attachement_ids
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;
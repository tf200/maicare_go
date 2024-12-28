-- name: CreateEmployeeProfile :one
INSERT INTO employee_profile (
    user_id,
    first_name,
    last_name,
    position,
    department,
    employee_number,
    employment_number,
    private_email_address,
    email_address,
    authentication_phone_number,
    private_phone_number,
    work_phone_number,
    date_of_birth,
    home_telephone_number,
    is_subcontractor,
    gender,
    location_id,
    has_borrowed,
    out_of_service,
    is_archived
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
) RETURNING *;



-- name: ListEmployeeProfile :many
SELECT 
    ep.*,
    u.profile_picture as profile_picture
FROM employee_profile ep
JOIN custom_user u ON ep.user_id = u.id
WHERE 
    (CASE 
        WHEN sqlc.narg('include_archived')::boolean IS NULL THEN true
        WHEN sqlc.narg('include_archived')::boolean = false THEN NOT ep.is_archived
        ELSE true
    END) AND
    (CASE 
        WHEN sqlc.narg('include_out_of_service')::boolean IS NULL THEN true
        WHEN sqlc.narg('include_out_of_service')::boolean = false THEN NOT COALESCE(ep.out_of_service, false)
        ELSE true
    END) AND
    (ep.department = sqlc.narg('department') OR sqlc.narg('department') IS NULL) AND
    (ep.position = sqlc.narg('position') OR sqlc.narg('position') IS NULL) AND
    (ep.location_id = sqlc.narg('location_id') OR sqlc.narg('location_id') IS NULL) AND
    (sqlc.narg('search')::TEXT IS NULL OR 
        ep.first_name ILIKE '%' || sqlc.narg('search') || '%' OR 
        ep.last_name ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY ep.created DESC
LIMIT $1 OFFSET $2;


-- name: CountEmployeeProfile :one
SELECT COUNT(*) 
FROM employee_profile ep
WHERE 
    (CASE 
        WHEN sqlc.narg('include_archived')::boolean IS NULL THEN true
        WHEN sqlc.narg('include_archived')::boolean = false THEN NOT ep.is_archived
        ELSE true
    END) AND
    (CASE 
        WHEN sqlc.narg('include_out_of_service')::boolean IS NULL THEN true
        WHEN sqlc.narg('include_out_of_service')::boolean = false THEN NOT COALESCE(ep.out_of_service, false)
        ELSE true
    END) AND
    (department = sqlc.narg('department') OR sqlc.narg('department') IS NULL) AND
    (position = sqlc.narg('position') OR sqlc.narg('position') IS NULL) AND
    (location_id = sqlc.narg('location_id') OR sqlc.narg('location_id') IS NULL);
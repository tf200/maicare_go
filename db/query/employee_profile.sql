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
    email,
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
ORDER BY ep.created_at DESC
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




-- name: GetEmployeeProfileByUserID :one
SELECT 
    cu.id as user_id,
    cu.email as email,
    ep.id as employee_id,
    ep.first_name,
    ep.last_name,
    cu.role_id
FROM custom_user cu
JOIN employee_profile ep ON ep.user_id = cu.id
WHERE cu.id = $1;


-- name: GetEmployeeProfileByID :one
SELECT 
    ep.*,
    cu.profile_picture as profile_picture,
    cu.role_id
FROM employee_profile ep
JOIN custom_user cu ON ep.user_id = cu.id
WHERE ep.id = $1;


-- name: UpdateEmployeeProfile :one
UPDATE employee_profile
SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    position = COALESCE(sqlc.narg('position'), position),
    department = COALESCE(sqlc.narg('department'), department),
    employee_number = COALESCE(sqlc.narg('employee_number'), employee_number),
    employment_number = COALESCE(sqlc.narg('employment_number'), employment_number),
    private_email_address = COALESCE(sqlc.narg('private_email_address'), private_email_address),
    email = COALESCE(sqlc.narg('email'), email),
    authentication_phone_number = COALESCE(sqlc.narg('authentication_phone_number'), authentication_phone_number),
    private_phone_number = COALESCE(sqlc.narg('private_phone_number'), private_phone_number),
    work_phone_number = COALESCE(sqlc.narg('work_phone_number'), work_phone_number),
    date_of_birth = COALESCE(sqlc.narg('date_of_birth'), date_of_birth),
    home_telephone_number = COALESCE(sqlc.narg('home_telephone_number'), home_telephone_number),
    is_subcontractor = COALESCE(sqlc.narg('is_subcontractor'), is_subcontractor),
    gender = COALESCE(sqlc.narg('gender'), gender),
    location_id = COALESCE(sqlc.narg('location_id'), location_id),
    has_borrowed = COALESCE(sqlc.narg('has_borrowed'), has_borrowed),
    out_of_service = COALESCE(sqlc.narg('out_of_service'), out_of_service),
    is_archived = COALESCE(sqlc.narg('is_archived'), is_archived)
WHERE id = sqlc.arg('id')
RETURNING *;
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


-- name: GetUserIDByEmployeeID :one
SELECT user_id FROM employee_profile
WHERE id = $1 LIMIT 1;

-- name: GetEmployeeProfileByUserID :one
SELECT 
    cu.id           AS user_id,
    cu.email        AS email,
    ep.id           AS employee_id,
    ep.first_name,
    ep.last_name,
    (
        SELECT COALESCE(json_agg(json_build_object(
            'id',       p.id,
            'name',     p.name,
            'resource', p.resource,
            'method',   p.method
        )), '[]'::json)
        FROM user_permissions up
        JOIN permissions p ON p.id = up.permission_id
        WHERE up.user_id = cu.id
    )::json AS permissions
FROM custom_user cu
JOIN employee_profile ep ON ep.user_id = cu.id
WHERE cu.id = $1;


-- name: GetEmployeeProfileByID :one
SELECT 
    ep.*,
    cu.profile_picture as profile_picture
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

-- name: SetEmployeeProfilePicture :one
UPDATE custom_user
SET profile_picture = $2
WHERE id = (
    SELECT user_id 
    FROM employee_profile
    WHERE employee_profile.id = $1
)
RETURNING *;


-- name: AddEmployeeContractDetails :one
UPDATE employee_profile
SET
    fixed_contract_hours = COALESCE(sqlc.narg('fixed_contract_hours'), fixed_contract_hours),
    variable_contract_hours = COALESCE(sqlc.narg('variable_contract_hours'), variable_contract_hours),
    contract_start_date = COALESCE(sqlc.narg('contract_start_date'), contract_start_date),
    contract_end_date = COALESCE(sqlc.narg('contract_end_date'), contract_end_date),
    contract_type = COALESCE(sqlc.narg('contract_type'), contract_type),
    contract_rate = COALESCE(sqlc.narg('contract_rate'), contract_rate)
WHERE id = $1
RETURNING *;

-- name: GetEmployeeContractDetails :one
SELECT
    fixed_contract_hours,
    variable_contract_hours,
    contract_start_date,
    contract_end_date,
    contract_type,
    contract_rate
FROM employee_profile
WHERE id = $1;


-- name: AddEducationToEmployeeProfile :one
INSERT INTO employee_education (
    employee_id,
    institution_name,
    degree,
    field_of_study,
    start_date,
    end_date
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: ListEducations :many
SELECT * FROM employee_education WHERE employee_id = $1;

-- name: UpdateEmployeeEducation :one
UPDATE employee_education
SET
    institution_name = COALESCE(sqlc.narg('institution_name'), institution_name),
    degree = COALESCE(sqlc.narg('degree'), degree),
    field_of_study = COALESCE(sqlc.narg('field_of_study'), field_of_study),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date)
WHERE id = $1
RETURNING *;


-- name: DeleteEmployeeEducation :one
DELETE FROM employee_education WHERE id = $1 RETURNING *;




-- name: AddEmployeeExperience :one
INSERT INTO employee_experience (
    employee_id,
    job_title,
    company_name,
    start_date,
    end_date,
    description
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: ListEmployeeExperience :many
SELECT * FROM employee_experience WHERE employee_id = $1;


-- name: UpdateEmployeeExperience :one
UPDATE employee_experience
SET
    job_title = COALESCE(sqlc.narg('job_title'), job_title),
    company_name = COALESCE(sqlc.narg('company_name'), company_name),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date),
    description = COALESCE(sqlc.narg('description'), description)
WHERE id = $1
RETURNING *;

-- name: DeleteEmployeeExperience :one
DELETE FROM employee_experience WHERE id = $1 RETURNING *;


-- name: AddEmployeeCertification :one
INSERT INTO certification (
    employee_id,
    name,
    issued_by,
    date_issued
) VALUES (
    $1, $2, $3, $4
) 
RETURNING *;


-- name: ListEmployeeCertifications :many
SELECT * FROM certification WHERE employee_id = $1;

-- name: UpdateEmployeeCertification :one
UPDATE certification
SET
    name = COALESCE(sqlc.narg('name'), name),
    issued_by = COALESCE(sqlc.narg('issued_by'), issued_by),
    date_issued = COALESCE(sqlc.narg('date_issued'), date_issued)
WHERE id = $1
RETURNING *;


-- name: DeleteEmployeeCertification :one
DELETE FROM certification WHERE id = $1 RETURNING *;


-- name: SearchEmployeesByNameOrEmail :many
SELECT
    id,
    first_name,
    last_name,
    email
FROM employee_profile
WHERE 
    first_name ILIKE '%' || @search || '%' OR
    last_name ILIKE '%' || @search || '%' OR
    email ILIKE '%' || @search || '%'
LIMIT 10;


-- name: GetEmployeeCounts :one
SELECT
    COUNT(*) FILTER (WHERE is_subcontractor IS NOT TRUE) AS total_employees,
    COUNT(*) FILTER (WHERE is_subcontractor = TRUE) AS total_subcontractors,
    COUNT(*) FILTER (WHERE is_archived = TRUE) AS total_archived,
    COUNT(*) FILTER (WHERE out_of_service = TRUE) AS total_out_of_service
FROM
    employee_profile;

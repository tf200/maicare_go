-- name: AddEmployeeContractDetails :one
UPDATE employee_profile
SET
    contract_hours = COALESCE(sqlc.narg('contract_hours'), contract_hours),
    contract_start_date = COALESCE(sqlc.narg('contract_start_date'), contract_start_date),
    contract_end_date = COALESCE(sqlc.narg('contract_end_date'), contract_end_date),
    contract_type = COALESCE(sqlc.narg('contract_type'), contract_type),
    contract_rate = COALESCE(sqlc.narg('contract_rate'), contract_rate)
WHERE id = $1
RETURNING *;

-- name: GetEmployeeContractDetails :one
SELECT
    contract_hours,
    contract_start_date,
    contract_end_date,
    contract_type,
    contract_rate,
    is_subcontractor
FROM employee_profile
WHERE id = $1;

-- name: UpdateEmployeeIsSubcontractor :one
UPDATE employee_profile
SET
    is_subcontractor = $2,
    contract_type = $3
WHERE id = $1
RETURNING *;
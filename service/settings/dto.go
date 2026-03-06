package settings

import (
	"time"

	"github.com/google/uuid"
)

type CreateDepartmentRequest struct {
	Name                     string     `json:"name" binding:"required"`
	Description              *string    `json:"description"`
	DepartmentHeadEmployeeID *uuid.UUID `json:"department_head_employee_id"`
}

type CreateDepartmentResponse struct {
	ID                       uuid.UUID  `json:"id"`
	Name                     string     `json:"name"`
	Description              *string    `json:"description"`
	DepartmentHeadEmployeeID *uuid.UUID `json:"department_head_employee_id"`
}

type UpdateDepartmentRequest struct {
	Name                     *string    `json:"name"`
	Description              *string    `json:"description"`
	DepartmentHeadEmployeeID *uuid.UUID `json:"department_head_employee_id"`
}

type UpdateDepartmentResponse struct {
	ID                       uuid.UUID  `json:"id"`
	Name                     string     `json:"name"`
	Description              *string    `json:"description"`
	DepartmentHeadEmployeeID *uuid.UUID `json:"department_head_employee_id"`
}

type ListDepartmentResponse struct {
	ID                       uuid.UUID  `json:"id"`
	Name                     string     `json:"name"`
	Description              *string    `json:"description"`
	DepartmentHeadEmployeeID *uuid.UUID `json:"department_head_employee_id"`
	EmployeeCount            int64      `json:"employee_count"`
}

type GetOrganizationProfileResponse struct {
	Name                  string    `json:"name"`
	DefaultTimezone       string    `json:"default_timezone"`
	Email                 *string   `json:"email"`
	PhoneNumber           *string   `json:"phone_number"`
	Website               *string   `json:"website"`
	HQStreet              *string   `json:"hq_street"`
	HQHouseNumber         *string   `json:"hq_house_number"`
	HQHouseNumberAddition *string   `json:"hq_house_number_addition"`
	HQPostalCode          *string   `json:"hq_postal_code"`
	HQCity                *string   `json:"hq_city"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type UpdateOrganizationProfileRequest struct {
	Name                  string  `json:"name" binding:"required"`
	DefaultTimezone       string  `json:"default_timezone" binding:"required"`
	Email                 *string `json:"email"`
	PhoneNumber           *string `json:"phone_number"`
	Website               *string `json:"website"`
	HQStreet              *string `json:"hq_street"`
	HQHouseNumber         *string `json:"hq_house_number"`
	HQHouseNumberAddition *string `json:"hq_house_number_addition"`
	HQPostalCode          *string `json:"hq_postal_code"`
	HQCity                *string `json:"hq_city"`
}

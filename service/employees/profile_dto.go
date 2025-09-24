package employees

import (
	"maicare_go/pagination"
	"time"
)

// CreateEmployeeProfileRequest represents the request for CreateEmployeeProfileApi
type CreateEmployeeProfileRequest struct {
	EmployeeNumber            *string `json:"employee_number" example:"123456"`
	EmploymentNumber          *string `json:"employment_number" example:"123456"`
	LocationID                *int64  `json:"location_id" example:"1"`
	IsSubcontractor           *bool   `json:"is_subcontractor" binding:"required" example:"false"`
	FirstName                 string  `json:"first_name" binding:"required" example:"fara"`
	LastName                  string  `json:"last_name" binding:"required" example:"joe"`
	DateOfBirth               *string `json:"date_of_birth" example:"2000-01-01"`
	Gender                    *string `json:"gender" example:"man"`
	Email                     string  `json:"email" binding:"required,email" example:"emai@exe.com"`
	PrivateEmailAddress       *string `json:"private_email_address" binding:"email" example:"joe@ex.com"`
	AuthenticationPhoneNumber *string `json:"authentication_phone_number" example:"1234567890"`
	WorkPhoneNumber           *string `json:"work_phone_number" example:"1234567890"`
	PrivatePhoneNumber        *string `json:"private_phone_number" example:"1234567890"`
	HomeTelephoneNumber       *string `json:"home_telephone_number" example:"1234567890"`
	RoleID                    int32   `json:"role_id" binding:"required" example:"1"`
	Position                  *string `json:"position" example:"developer"`
	Department                *string `json:"department" example:"IT"`
}

// CreateEmployeeProfileResponse represents the response for CreateEmployeeProfileApi
type CreateEmployeeProfileResponse struct {
	ID                        int64     `json:"id"`
	UserID                    int64     `json:"user_id"`
	FirstName                 string    `json:"first_name"`
	LastName                  string    `json:"last_name"`
	Position                  *string   `json:"position"`
	Department                *string   `json:"department"`
	EmployeeNumber            *string   `json:"employee_number"`
	EmploymentNumber          *string   `json:"employment_number"`
	PrivateEmailAddress       *string   `json:"private_email_address"`
	Email                     string    `json:"email"`
	AuthenticationPhoneNumber *string   `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string   `json:"private_phone_number"`
	WorkPhoneNumber           *string   `json:"work_phone_number"`
	DateOfBirth               time.Time `json:"date_of_birth"`
	HomeTelephoneNumber       *string   `json:"home_telephone_number"`
	CreatedAt                 time.Time `json:"created_at"`
	IsSubcontractor           *bool     `json:"is_subcontractor"`
	Gender                    *string   `json:"gender" binding:"oneof= male female not_specified"`
	LocationID                *int64    `json:"location_id"`
	HasBorrowed               bool      `json:"has_borrowed"`
	OutOfService              *bool     `json:"out_of_service"`
	IsArchived                bool      `json:"is_archived"`
}

// ListEmployeeRequest represents the request for ListEmployeeProfileApi
type ListEmployeeRequest struct {
	pagination.Request
	IncludeArchived     *bool   `form:"is_archived"`
	IncludeOutOfService *bool   `form:"out_of_service"`
	Department          *string `form:"department"`
	Position            *string `form:"position"`
	LocationID          *int64  `form:"location_id"`
	Search              *string `form:"search"`
}

// ListEmployeeResponse represents the response for ListEmployeeProfileApi
type ListEmployeeResponse struct {
	ID                        int64     `json:"id"`
	UserID                    int64     `json:"user_id"`
	FirstName                 string    `json:"first_name"`
	LastName                  string    `json:"last_name"`
	Position                  *string   `json:"position"`
	Department                *string   `json:"department"`
	EmployeeNumber            *string   `json:"employee_number"`
	EmploymentNumber          *string   `json:"employment_number"`
	PrivateEmailAddress       *string   `json:"private_email_address"`
	Email                     string    `json:"email"`
	AuthenticationPhoneNumber *string   `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string   `json:"private_phone_number"`
	WorkPhoneNumber           *string   `json:"work_phone_number"`
	DateOfBirth               time.Time `json:"date_of_birth"`
	HomeTelephoneNumber       *string   `json:"home_telephone_number"`
	CreatedAt                 time.Time `json:"created_at"`
	IsSubcontractor           *bool     `json:"is_subcontractor"`
	Gender                    *string   `json:"gender"`
	LocationID                *int64    `json:"location_id"`
	HasBorrowed               bool      `json:"has_borrowed"`
	OutOfService              *bool     `json:"out_of_service"`
	IsArchived                bool      `json:"is_archived"`
	ProfilePicture            *string   `json:"profile_picture"`
	Age                       int64     `json:"age"`
	RoleID                    *int32    `json:"role_id"`
	RoleName                  *string   `json:"role_name"`
}

// UpdateEmployeeIsSubcontractorRequest represents the request for UpdateEmployeeIsSubcontractorApi
type UpdateEmployeeIsSubcontractorRequest struct {
	IsSubcontractor *bool `json:"is_subcontractor" binding:"required"`
}

// UpdateEmployeeIsSubcontractorResponse represents the response for UpdateEmployeeIsSubcontractorApi
type UpdateEmployeeIsSubcontractorResponse struct {
	ID                int64     `json:"id"`
	IsSubcontractor   *bool     `json:"is_subcontractor"`
	ContractType      *string   `json:"contract_type"`
	ContractHours     *float64  `json:"contract_hours"`
	ContractRate      *float64  `json:"contract_rate"`
	ContractStartDate time.Time `json:"contract_start_date"`
	ContractEndDate   time.Time `json:"contract_end_date"`
}

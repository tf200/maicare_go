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

// GetEmployeeProfileResponse represents the response for GetEmployeeProfile
type GetEmployeeProfileResponse struct {
	UserID      int64        `json:"user_id"`
	Email       string       `json:"email"`
	EmployeeID  int64        `json:"employee_id"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	TwoFactor   bool         `json:"two_factor_enabled"`
	LastLogin   time.Time    `json:"last_login"`
	RoleID      int32        `json:"role_id"`
	Permissions []Permission `json:"permissions"`
}

// Permission represents a permission entity
type Permission struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Resource string `json:"resource"`
	Method   string `json:"method"`
}

// GetEmployeeProfileByIDResponse represents the response for GetEmployeeProfileByID
type GetEmployeeProfileByIDResponse struct {
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
	RoleID                    int32     `json:"role_id"`
	IsLoggedInUser            bool      `json:"is_logged_in_user"`
}

// UpdateEmployeeProfileRequest represents the request for UpdateEmployeeProfile
type UpdateEmployeeProfileRequest struct {
	FirstName                 *string `json:"first_name"`
	LastName                  *string `json:"last_name"`
	Position                  *string `json:"position"`
	Department                *string `json:"department"`
	EmployeeNumber            *string `json:"employee_number"`
	EmploymentNumber          *string `json:"employment_number"`
	PrivateEmailAddress       *string `json:"private_email_address"`
	Email                     *string `json:"email"`
	AuthenticationPhoneNumber *string `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string `json:"private_phone_number"`
	WorkPhoneNumber           *string `json:"work_phone_number"`
	DateOfBirth               *string `json:"date_of_birth"`
	HomeTelephoneNumber       *string `json:"home_telephone_number"`
	IsSubcontractor           *bool   `json:"is_subcontractor"`
	Gender                    *string `json:"gender"`
	LocationID                *int64  `json:"location_id"`
	HasBorrowed               *bool   `json:"has_borrowed"`
	OutOfService              *bool   `json:"out_of_service"`
	IsArchived                *bool   `json:"is_archived"`
}

// UpdateEmployeeProfileResponse represents the response for UpdateEmployeeProfile
type UpdateEmployeeProfileResponse struct {
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
}

// SetEmployeeProfilePictureRequest represents the request for SetEmployeeProfilePicture
type SetEmployeeProfilePictureRequest struct {
	AttachmentID string `json:"attachement_id" binding:"required"`
}

// SetEmployeeProfilePictureResponse represents the response for SetEmployeeProfilePicture
type SetEmployeeProfilePictureResponse struct {
	ID             int64   `json:"id"`
	Email          string  `json:"email"`
	ProfilePicture *string `json:"profile_picture"`
}

// GetEmployeeCountsResponse represents the response for GetEmployeeCounts
type GetEmployeeCountsResponse struct {
	TotalEmployees      int64 `json:"total_employees"`
	TotalSubcontractors int64 `json:"total_subcontractors"`
	TotalArchived       int64 `json:"total_archived"`
	TotalOutOfService   int64 `json:"total_out_of_service"`
}

// SearchEmployeesByNameOrEmailRequest represents the request for SearchEmployeesByNameOrEmail
type SearchEmployeesByNameOrEmailRequest struct {
	Search *string `form:"search" binding:"required"`
}

// SearchEmployeesByNameOrEmailResponse represents the response for SearchEmployeesByNameOrEmail
type SearchEmployeesByNameOrEmailResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// AddEmployeeContractDetailsRequest represents the request for AddEmployeeContractDetails
type AddEmployeeContractDetailsRequest struct {
	ContractHours     *float64  `json:"contract_hours" binding:"required"`
	ContractStartDate time.Time `json:"contract_start_date"`
	ContractEndDate   time.Time `json:"contract_end_date"`
	ContractRate      *float64  `json:"contract_rate"` // Optional field for contract rate
}

// AddEmployeeContractDetailsResponse represents the response for AddEmployeeContractDetails
type AddEmployeeContractDetailsResponse struct {
	ID                int64     `json:"id"`
	ContractHours     *float64  `json:"contract_hours"`
	ContractStartDate time.Time `json:"contract_start_date"`
	ContractEndDate   time.Time `json:"contract_end_date"`
	ContractRate      *float64  `json:"contract_rate"` // Optional field for contract rate
}

// GetEmployeeContractDetailsResponse represents the response for GetEmployeeContractDetails
type GetEmployeeContractDetailsResponse struct {
	ContractHours     *float64  `json:"contract_hours"`
	ContractStartDate time.Time `json:"contract_start_date"`
	ContractEndDate   time.Time `json:"contract_end_date"`
	ContractType      *string   `json:"contract_type"`
	ContractRate      *float64  `json:"contract_rate"` // Optional field for contract rate
	IsSubcontractor   *bool     `json:"is_subcontractor"`
}

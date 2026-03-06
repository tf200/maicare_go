package employees

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

// CreateEmployeeProfileRequest represents the request for CreateEmployeeProfileApi
type CreateEmployeeProfileRequest struct {
	EmployeeNumber      *string    `json:"employee_number" example:"123456"`
	EmploymentNumber    *string    `json:"employment_number" example:"123456"`
	LocationID          *uuid.UUID `json:"location_id" example:"1"`
	FirstName           string     `json:"first_name" binding:"required" example:"fara"`
	LastName            string     `json:"last_name" binding:"required" example:"joe"`
	Bsn                 string     `json:"bsn" binding:"required" example:"123456789"`
	Street              string     `json:"street" binding:"required" example:"Main St"`
	HouseNumber         string     `json:"house_number" binding:"required" example:"10"`
	HouseNumberAddition *string    `json:"house_number_addition" example:"A"`
	PostalCode          string     `json:"postal_code" binding:"required" example:"1234AB"`
	City                string     `json:"city" binding:"required" example:"Amsterdam"`
	Position            *string    `json:"position" example:"developer"`
	DepartmentID        *uuid.UUID `json:"department_id" example:"1"`
	ManagerEmployeeID   *uuid.UUID `json:"manager_employee_id" example:"1"`
	PrivateEmailAddress *string    `json:"private_email_address" binding:"email" example:"joe@ex.com"`
	WorkEmailAddress    string     `json:"work_email_address" binding:"required,email" example:"email@exe.com"`
	WorkPhoneNumber     *string    `json:"work_phone_number" example:"1234567890"`
	PrivatePhoneNumber  *string    `json:"private_phone_number" example:"1234567890"`
	DateOfBirth         *string    `json:"date_of_birth" example:"2000-01-01"`
	HomeTelephoneNumber *string    `json:"home_telephone_number" example:"1234567890"`
	Gender              string     `json:"gender" binding:"required,oneof=male female not_specified" example:"male"`
	ContractHours       *float64   `json:"contract_hours" example:"40"`
	ContractStartDate   *string    `json:"contract_start_date" example:"2000-01-01"`
	ContractEndDate     *string    `json:"contract_end_date" example:"2000-01-01"`
	ContractType        string     `json:"contract_type" binding:"required,oneof=loondienst ZZP none" example:"loondienst"`
	ContractRate        *float64   `json:"contract_rate" example:"100.00"`
	RoleID              uuid.UUID  `json:"role_id" binding:"required" example:"1"`
}

// CreateEmployeeProfileResponse represents the response for CreateEmployeeProfileApi
type CreateEmployeeProfileResponse struct {
	ID                  uuid.UUID  `json:"id"`
	UserID              uuid.UUID  `json:"user_id"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	Position            *string    `json:"position"`
	DepartmentID        *uuid.UUID `json:"department_id"`
	DepartmentName      *string    `json:"department_name"`
	ManagerEmployeeID   *uuid.UUID `json:"manager_employee_id"`
	ManagerFirstName    *string    `json:"manager_first_name"`
	ManagerLastName     *string    `json:"manager_last_name"`
	EmployeeNumber      *string    `json:"employee_number"`
	EmploymentNumber    *string    `json:"employment_number"`
	PrivateEmailAddress *string    `json:"private_email_address"`
	Email               string     `json:"email"`
	PrivatePhoneNumber  *string    `json:"private_phone_number"`
	WorkPhoneNumber     *string    `json:"work_phone_number"`
	DateOfBirth         time.Time  `json:"date_of_birth"`
	HomeTelephoneNumber *string    `json:"home_telephone_number"`
	CreatedAt           time.Time  `json:"created_at"`
	Gender              string     `json:"gender" binding:"oneof= male female not_specified"`
	LocationID          *uuid.UUID `json:"location_id"`
	HasBorrowed         bool       `json:"has_borrowed"`
	OutOfService        *bool      `json:"out_of_service"`
	IsArchived          bool       `json:"is_archived"`
}

// ListEmployeeRequest represents the request for ListEmployeeProfileApi
type ListEmployeeRequest struct {
	pagination.Request
	IncludeArchived     *bool      `form:"is_archived"`
	IncludeOutOfService *bool      `form:"out_of_service"`
	LocationID          *uuid.UUID `form:"location_id"`
	ContractType        *string    `form:"contract_type" binding:"omitempty,oneof=loondienst ZZP none"`
	Search              *string    `form:"search"`
}

// ListEmployeeResponse represents the response for ListEmployeeProfileApi
type ListEmployeeResponse struct {
	ID              uuid.UUID  `json:"id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Bsn             string     `json:"bsn"`
	ContractType    string     `json:"contract_type"`
	DepartmentName  *string    `json:"department_name"`
	LocationAddress string     `json:"location_address"`
	ContractEndDate *time.Time `json:"contract_end_date"`
}

// UpdateEmployeeIsSubcontractorRequest represents the request for UpdateEmployeeIsSubcontractorApi
type UpdateEmployeeIsSubcontractorRequest struct {
	IsSubcontractor *bool `json:"is_subcontractor" binding:"required"`
}

// UpdateEmployeeIsSubcontractorResponse represents the response for UpdateEmployeeIsSubcontractorApi
type UpdateEmployeeIsSubcontractorResponse struct {
	ID                uuid.UUID `json:"id"`
	IsSubcontractor   *bool     `json:"is_subcontractor"`
	ContractType      string    `json:"contract_type"`
	ContractHours     *float64  `json:"contract_hours"`
	ContractRate      *float64  `json:"contract_rate"`
	ContractStartDate time.Time `json:"contract_start_date"`
	ContractEndDate   time.Time `json:"contract_end_date"`
}

// GetEmployeeProfileResponse represents the response for GetEmployeeProfile
type GetEmployeeProfileResponse struct {
	UserID      uuid.UUID    `json:"user_id"`
	Email       string       `json:"email"`
	EmployeeID  uuid.UUID    `json:"employee_id"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	TwoFactor   bool         `json:"two_factor_enabled"`
	LastLogin   time.Time    `json:"last_login"`
	RoleID      *uuid.UUID   `json:"role_id"`
	Permissions []Permission `json:"permissions"`
}

type EmployeeRole struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BriefEducationDetail struct {
	InstitutionName string     `json:"institution_name"`
	Degree          string     `json:"degree"`
	FieldOfStudy    string     `json:"field_of_study"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
}

type BriefExperienceDetail struct {
	JobTitle    string     `json:"job_title"`
	CompanyName string     `json:"company_name"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type ActiveSessionDetail struct {
	ID        uuid.UUID `json:"id"`
	UserAgent string    `json:"user_agent"`
	ClientIP  string    `json:"client_ip"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type GetEmployeeProfileDetailsResponse struct {
	UserID              uuid.UUID               `json:"user_id"`
	EmployeeID          uuid.UUID               `json:"employee_id"`
	Email               string                  `json:"email"`
	FirstName           string                  `json:"first_name"`
	LastName            string                  `json:"last_name"`
	TwoFactorEnabled    bool                    `json:"two_factor_enabled"`
	LastLogin           time.Time               `json:"last_login"`
	Roles               []EmployeeRole          `json:"roles"`
	ActiveSessions      []ActiveSessionDetail   `json:"active_sessions"`
	Education           []BriefEducationDetail  `json:"education"`
	WorkExperience      []BriefExperienceDetail `json:"work_experience"`
	Street              string                  `json:"street"`
	HouseNumber         string                  `json:"house_number"`
	HouseNumberAddition *string                 `json:"house_number_addition"`
	PostalCode          string                  `json:"postal_code"`
	City                string                  `json:"city"`
	Position            *string                 `json:"position"`
	DepartmentID        *uuid.UUID              `json:"department_id"`
	DepartmentName      *string                 `json:"department_name"`
	ManagerEmployeeID   *uuid.UUID              `json:"manager_employee_id"`
	ManagerFirstName    *string                 `json:"manager_first_name"`
	ManagerLastName     *string                 `json:"manager_last_name"`
	EmployeeNumber      *string                 `json:"employee_number"`
	EmploymentNumber    *string                 `json:"employment_number"`
	PrivateEmailAddress *string                 `json:"private_email_address"`
	WorkEmailAddress    *string                 `json:"work_email_address"`
	PrivatePhoneNumber  *string                 `json:"private_phone_number"`
	WorkPhoneNumber     *string                 `json:"work_phone_number"`
	HomeTelephoneNumber *string                 `json:"home_telephone_number"`
	DateOfBirth         *time.Time              `json:"date_of_birth"`
	Gender              string                  `json:"gender"`
	LocationID          *uuid.UUID              `json:"location_id"`
	LocationName        *string                 `json:"location_name"`
	OrganisationName    *string                 `json:"organisation_name"`
	HasBorrowed         bool                    `json:"has_borrowed"`
	OutOfService        *bool                   `json:"out_of_service"`
	IsArchived          bool                    `json:"is_archived"`
	ContractType        string                  `json:"contract_type"`
	ContractHours       *float64                `json:"contract_hours"`
	ContractStartDate   *time.Time              `json:"contract_start_date"`
	ContractEndDate     *time.Time              `json:"contract_end_date"`
	ContractRate        *float64                `json:"contract_rate"`
}

// Permission represents a permission entity
type Permission struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Resource string    `json:"resource"`
	Method   string    `json:"method"`
}

// GetEmployeeProfileByIDResponse represents the response for GetEmployeeProfileByID
type GetEmployeeProfileByIDResponse struct {
	ID                        uuid.UUID  `json:"id"`
	UserID                    uuid.UUID  `json:"user_id"`
	FirstName                 string     `json:"first_name"`
	LastName                  string     `json:"last_name"`
	Position                  *string    `json:"position"`
	DepartmentID              *uuid.UUID `json:"department_id"`
	DepartmentName            *string    `json:"department_name"`
	ManagerEmployeeID         *uuid.UUID `json:"manager_employee_id"`
	ManagerFirstName          *string    `json:"manager_first_name"`
	ManagerLastName           *string    `json:"manager_last_name"`
	EmployeeNumber            *string    `json:"employee_number"`
	EmploymentNumber          *string    `json:"employment_number"`
	PrivateEmailAddress       *string    `json:"private_email_address"`
	Email                     string     `json:"email"`
	AuthenticationPhoneNumber *string    `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string    `json:"private_phone_number"`
	WorkPhoneNumber           *string    `json:"work_phone_number"`
	DateOfBirth               time.Time  `json:"date_of_birth"`
	HomeTelephoneNumber       *string    `json:"home_telephone_number"`
	CreatedAt                 time.Time  `json:"created_at"`
	IsSubcontractor           *bool      `json:"is_subcontractor"`
	Gender                    string     `json:"gender"`
	LocationID                *uuid.UUID `json:"location_id"`
	HasBorrowed               bool       `json:"has_borrowed"`
	OutOfService              *bool      `json:"out_of_service"`
	IsArchived                bool       `json:"is_archived"`
	ProfilePicture            *string    `json:"profile_picture"`
	RoleID                    *uuid.UUID `json:"role_id"`
	IsLoggedInUser            bool       `json:"is_logged_in_user"`
}

// UpdateEmployeeProfileRequest represents the request for UpdateEmployeeProfile
type UpdateEmployeeProfileRequest struct {
	FirstName                 *string    `json:"first_name"`
	LastName                  *string    `json:"last_name"`
	Position                  *string    `json:"position"`
	DepartmentID              *uuid.UUID `json:"department_id"`
	ManagerEmployeeID         *uuid.UUID `json:"manager_employee_id"`
	EmployeeNumber            *string    `json:"employee_number"`
	EmploymentNumber          *string    `json:"employment_number"`
	PrivateEmailAddress       *string    `json:"private_email_address"`
	Email                     *string    `json:"email"`
	AuthenticationPhoneNumber *string    `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string    `json:"private_phone_number"`
	WorkPhoneNumber           *string    `json:"work_phone_number"`
	DateOfBirth               *string    `json:"date_of_birth"`
	HomeTelephoneNumber       *string    `json:"home_telephone_number"`
	IsSubcontractor           *bool      `json:"is_subcontractor"`
	Gender                    *string    `json:"gender"`
	LocationID                *uuid.UUID `json:"location_id"`
	HasBorrowed               *bool      `json:"has_borrowed"`
	OutOfService              *bool      `json:"out_of_service"`
	IsArchived                *bool      `json:"is_archived"`
}

// UpdateEmployeeProfileResponse represents the response for UpdateEmployeeProfile
type UpdateEmployeeProfileResponse struct {
	ID                        uuid.UUID  `json:"id"`
	UserID                    uuid.UUID  `json:"user_id"`
	FirstName                 string     `json:"first_name"`
	LastName                  string     `json:"last_name"`
	Position                  *string    `json:"position"`
	DepartmentID              *uuid.UUID `json:"department_id"`
	DepartmentName            *string    `json:"department_name"`
	ManagerEmployeeID         *uuid.UUID `json:"manager_employee_id"`
	ManagerFirstName          *string    `json:"manager_first_name"`
	ManagerLastName           *string    `json:"manager_last_name"`
	EmployeeNumber            *string    `json:"employee_number"`
	EmploymentNumber          *string    `json:"employment_number"`
	PrivateEmailAddress       *string    `json:"private_email_address"`
	Email                     string     `json:"email"`
	AuthenticationPhoneNumber *string    `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string    `json:"private_phone_number"`
	WorkPhoneNumber           *string    `json:"work_phone_number"`
	DateOfBirth               time.Time  `json:"date_of_birth"`
	HomeTelephoneNumber       *string    `json:"home_telephone_number"`
	CreatedAt                 time.Time  `json:"created_at"`
	IsSubcontractor           *bool      `json:"is_subcontractor"`
	Gender                    string     `json:"gender"`
	LocationID                *uuid.UUID `json:"location_id"`
	HasBorrowed               bool       `json:"has_borrowed"`
	OutOfService              *bool      `json:"out_of_service"`
	IsArchived                bool       `json:"is_archived"`
}

// SetEmployeeProfilePictureRequest represents the request for SetEmployeeProfilePicture
type SetEmployeeProfilePictureRequest struct {
	AttachmentID string `json:"attachement_id" binding:"required"`
}

// SetEmployeeProfilePictureResponse represents the response for SetEmployeeProfilePicture
type SetEmployeeProfilePictureResponse struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	ProfilePicture *string   `json:"profile_picture"`
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
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
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
	ID                uuid.UUID `json:"id"`
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
	ContractType      string    `json:"contract_type"`
	ContractRate      *float64  `json:"contract_rate"` // Optional field for contract rate
	IsSubcontractor   *bool     `json:"is_subcontractor"`
}

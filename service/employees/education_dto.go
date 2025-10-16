package employees

import (
	"time"
)

// AddEducationToEmployeeProfileRequest represents the request for AddEducationToEmployeeProfile
type AddEducationToEmployeeProfileRequest struct {
	InstitutionName string `json:"institution_name" binding:"required"`
	Degree          string `json:"degree" binding:"required"`
	FieldOfStudy    string `json:"field_of_study" binding:"required"`
	StartDate       string `json:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate         string `json:"end_date" binding:"required" time_format:"2006-01-02"`
}

// AddEducationToEmployeeProfileResponse represents the response for AddEducationToEmployeeProfile
type AddEducationToEmployeeProfileResponse struct {
	ID              int64     `json:"id"`
	EmployeeID      int64     `json:"employee_id"`
	InstitutionName string    `json:"institution_name"`
	Degree          string    `json:"degree"`
	FieldOfStudy    string    `json:"field_of_study"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}

// ListEmployeeEducationResponse represents the response for ListEmployeeEducation
type ListEmployeeEducationResponse struct {
	ID              int64     `json:"id"`
	EmployeeID      int64     `json:"employee_id"`
	InstitutionName string    `json:"institution_name"`
	Degree          string    `json:"degree"`
	FieldOfStudy    string    `json:"field_of_study"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}

// UpdateEmployeeEducationRequest represents the request for UpdateEmployeeEducation
type UpdateEmployeeEducationRequest struct {
	InstitutionName *string `json:"institution_name"`
	Degree          *string `json:"degree"`
	FieldOfStudy    *string `json:"field_of_study"`
	StartDate       *string `json:"start_date" time_format:"2006-01-02"`
	EndDate         *string `json:"end_date" time_format:"2006-01-02"`
}

// UpdateEmployeeEducationResponse represents the response for UpdateEmployeeEducation
type UpdateEmployeeEducationResponse struct {
	ID              int64     `json:"id"`
	InstitutionName string    `json:"institution_name"`
	Degree          string    `json:"degree"`
	FieldOfStudy    string    `json:"field_of_study"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}

// DeleteEmployeeEducationResponse represents the response for DeleteEmployeeEducation
type DeleteEmployeeEducationResponse struct {
	ID              int64     `json:"id"`
	InstitutionName string    `json:"institution_name"`
	Degree          string    `json:"degree"`
	FieldOfStudy    string    `json:"field_of_study"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}
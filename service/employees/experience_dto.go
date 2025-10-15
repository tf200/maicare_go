package employees

import (
	"time"
)

// AddEmployeeExperienceRequest represents the request for AddEmployeeExperience
type AddEmployeeExperienceRequest struct {
	JobTitle    string  `json:"job_title" binding:"required"`
	CompanyName string  `json:"company_name" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate     string  `json:"end_date" binding:"required" time_format:"2006-01-02"`
	Description *string `json:"description"`
}

// AddEmployeeExperienceResponse represents the response for AddEmployeeExperience
type AddEmployeeExperienceResponse struct {
	ID          int64   `json:"id"`
	EmployeeID  int64   `json:"employee_id"`
	JobTitle    string  `json:"job_title"`
	CompanyName string  `json:"company_name"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

// ListEmployeeExperienceResponse represents the response for ListEmployeeExperience
type ListEmployeeExperienceResponse struct {
	ID          int64   `json:"id"`
	EmployeeID  int64   `json:"employee_id"`
	JobTitle    string  `json:"job_title"`
	CompanyName string  `json:"company_name"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

// UpdateEmployeeExperienceRequest represents the request for UpdateEmployeeExperience
type UpdateEmployeeExperienceRequest struct {
	JobTitle    *string `json:"job_title"`
	CompanyName *string `json:"company_name"`
	StartDate   *string `json:"start_date" time_format:"2006-01-02"`
	EndDate     *string `json:"end_date" time_format:"2006-01-02"`
	Description *string `json:"description"`
}

// UpdateEmployeeExperienceResponse represents the response for UpdateEmployeeExperience
type UpdateEmployeeExperienceResponse struct {
	ID          int64     `json:"id"`
	EmployeeID  int64     `json:"employee_id"`
	JobTitle    string    `json:"job_title"`
	CompanyName string    `json:"company_name"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// DeleteEmployeeExperienceResponse represents the response for DeleteEmployeeExperience
type DeleteEmployeeExperienceResponse struct {
	ID          int64     `json:"id"`
	EmployeeID  int64     `json:"employee_id"`
	JobTitle    string    `json:"job_title"`
	CompanyName string    `json:"company_name"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
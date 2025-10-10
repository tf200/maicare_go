package clientp

import (
	"maicare_go/pagination"
	"time"
)

// CreateProgressReportRequest defines the request payload for CreateProgressReport API
type CreateProgressReportRequest struct {
	EmployeeID     *int64    `json:"employee_id"`
	Title          *string   `json:"title"`
	Date           time.Time `json:"date"`
	ReportText     string    `json:"report_text" binding:"required"`
	Type           string    `json:"type" binding:"required,oneof=morning_report evening_report night_report shift_report one_to_one_report process_report contact_journal other"`
	EmotionalState string    `json:"emotional_state" binding:"required,oneof=normal excited happy sad angry anxious depressed"`
}

// CreateProgressReportResponse defines the response payload for CreateProgressReport API
type CreateProgressReportResponse struct {
	ID             int64     `json:"id"`
	ClientID       int64     `json:"client_id"`
	Date           time.Time `json:"date"`
	Title          *string   `json:"title"`
	ReportText     string    `json:"report_text"`
	EmployeeID     *int64    `json:"employee_id"`
	Type           string    `json:"type"`
	EmotionalState string    `json:"emotional_state"`
	CreatedAt      time.Time `json:"created_at"`
}

// ListProgressReportsRequest defines the request payload for ListProgressReports API
type ListProgressReportsRequest struct {
	pagination.Request
}

// ListProgressReportsResponse defines the response payload for ListProgressReports API
type ListProgressReportsResponse struct {
	ID                     int64     `json:"id"`
	ClientID               int64     `json:"client_id"`
	Date                   time.Time `json:"date"`
	Title                  *string   `json:"title"`
	ReportText             string    `json:"report_text"`
	EmployeeID             *int64    `json:"employee_id"`
	Type                   string    `json:"type"`
	EmotionalState         string    `json:"emotional_state"`
	CreatedAt              time.Time `json:"created_at"`
	EmployeeFirstName      string    `json:"employee_first_name"`
	EmployeeLastName       string    `json:"employee_last_name"`
	EmployeeProfilePicture *string   `json:"employee_profile_picture"`
}

// GetProgressReportResponse defines the response payload for GetProgressReport API
type GetProgressReportResponse struct {
	ID                     int64     `json:"id"`
	ClientID               int64     `json:"client_id"`
	Date                   time.Time `json:"date"`
	Title                  *string   `json:"title"`
	ReportText             string    `json:"report_text"`
	EmployeeID             *int64    `json:"employee_id"`
	Type                   string    `json:"type"`
	EmotionalState         string    `json:"emotional_state"`
	CreatedAt              time.Time `json:"created_at"`
	EmployeeFirstName      string    `json:"employee_first_name"`
	EmployeeLastName       string    `json:"employee_last_name"`
	EmployeeProfilePicture *string   `json:"employee_profile_picture"`
}

// UpdateProgressReportRequest defines the request payload for UpdateProgressReport API
type UpdateProgressReportRequest struct {
	ClientID       int64     `json:"client_id"`
	EmployeeID     *int64    `json:"employee_id"`
	Title          *string   `json:"title"`
	Date           time.Time `json:"date"`
	ReportText     *string   `json:"report_text"`
	Type           *string   `json:"type"`
	EmotionalState *string   `json:"emotional_state"`
}

// UpdateProgressReportResponse defines the response payload for UpdateProgressReport API
type UpdateProgressReportResponse struct {
	ID             int64     `json:"id"`
	ClientID       int64     `json:"client_id"`
	Date           time.Time `json:"date"`
	Title          *string   `json:"title"`
	ReportText     string    `json:"report_text"`
	EmployeeID     *int64    `json:"employee_id"`
	Type           string    `json:"type"`
	EmotionalState string    `json:"emotional_state"`
	CreatedAt      time.Time `json:"created_at"`
}

// GenerateAutoReportsRequest is the request format for the auto reports generation API
type GenerateAutoReportsRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

// GenerateAutoReportsResponse is the response format for the auto reports generation API
type GenerateAutoReportsResponse struct {
	Report string `json:"report"`
}

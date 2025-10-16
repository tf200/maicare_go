package employees

import (
	"time"
)

// AddEmployeeCertificationRequest represents the request for AddEmployeeCertification
type AddEmployeeCertificationRequest struct {
	Name       string `json:"name"`
	IssuedBy   string `json:"issued_by"`
	DateIssued string `json:"date_issued" time_format:"2006-01-02"`
}

// AddEmployeeCertificationResponse represents the response for AddEmployeeCertification
type AddEmployeeCertificationResponse struct {
	ID         int64     `json:"id"`
	EmployeeID int64     `json:"employee_id"`
	Name       string    `json:"name"`
	IssuedBy   string    `json:"issued_by"`
	DateIssued time.Time `json:"date_issued"`
	CreatedAt  time.Time `json:"created_at"`
}

// ListEmployeeCertificationResponse represents the response for ListEmployeeCertification
type ListEmployeeCertificationResponse struct {
	ID         int64     `json:"id"`
	EmployeeID int64     `json:"employee_id"`
	Name       string    `json:"name"`
	IssuedBy   string    `json:"issued_by"`
	DateIssued time.Time `json:"date_issued"`
}

// UpdateEmployeeCertificationRequest represents the request for UpdateEmployeeCertification
type UpdateEmployeeCertificationRequest struct {
	Name       *string `json:"name"`
	IssuedBy   *string `json:"issued_by"`
	DateIssued *string `json:"date_issued" time_format:"2006-01-02"`
}

// UpdateEmployeeCertificationResponse represents the response for UpdateEmployeeCertification
type UpdateEmployeeCertificationResponse struct {
	ID         int64     `json:"id"`
	EmployeeID int64     `json:"employee_id"`
	Name       string    `json:"name"`
	IssuedBy   string    `json:"issued_by"`
	DateIssued time.Time `json:"date_issued"`
	CreatedAt  time.Time `json:"created_at"`
}

// DeleteEmployeeCertificationResponse represents the response for DeleteEmployeeCertification
type DeleteEmployeeCertificationResponse struct {
	ID         int64     `json:"id"`
	EmployeeID int64     `json:"employee_id"`
	Name       string    `json:"name"`
	IssuedBy   string    `json:"issued_by"`
	DateIssued time.Time `json:"date_issued"`
	CreatedAt  time.Time `json:"created_at"`
}
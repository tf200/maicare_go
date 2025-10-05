package clientp

import (
	"maicare_go/pagination"
	"time"
)

// Contact represents a contact information.
type SenderContact struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"email"`
	PhoneNumber *string `json:"phone_number"`
}

// GetClientSenderResponse defines the request for getting a client sender
type GetClientSenderResponse struct {
	ID           int64           `json:"id"`
	Types        string          `json:"types"`
	Name         string          `json:"name"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	Kvknumber    *string         `json:"kvknumber"`
	Btwnumber    *string         `json:"btwnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	EmailAddress *string         `json:"email_address"`
	Contacts     []SenderContact `json:"contacts"`
	IsArchived   bool            `json:"is_archived"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// CreateClientEmergencyContactParams defines the request for creating a client emergency contact
type CreateClientEmergencyContactParams struct {
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	Email            *string `json:"email"`
	PhoneNumber      *string `json:"phone_number"`
	Address          *string `json:"address"`
	Relationship     *string `json:"relationship"`
	RelationStatus   *string `json:"relation_status"`
	MedicalReports   bool    `json:"medical_reports"`
	IncidentsReports bool    `json:"incidents_reports"`
	GoalsReports     bool    `json:"goals_reports"`
}

// CreateClientEmergencyContactResponse defines the response for creating a client emergency contact
type CreateClientEmergencyContactResponse struct {
	ID               int64     `json:"id"`
	ClientID         int64     `json:"client_id"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            *string   `json:"email"`
	PhoneNumber      *string   `json:"phone_number"`
	Address          *string   `json:"address"`
	Relationship     *string   `json:"relationship"`
	RelationStatus   *string   `json:"relation_status"`
	CreatedAt        time.Time `json:"created_at"`
	IsVerified       bool      `json:"is_verified"`
	MedicalReports   bool      `json:"medical_reports"`
	IncidentsReports bool      `json:"incidents_reports"`
	GoalsReports     bool      `json:"goals_reports"`
}

// ListClientEmergencyContactsRequest defines the request for listing client emergency contacts
type ListClientEmergencyContactsRequest struct {
	pagination.Request
	Search string `form:"search"`
}

// ListClientEmergencyContactsResponse defines the response for listing client emergency contacts
type ListClientEmergencyContactsResponse struct {
	ID               int64     `json:"id"`
	ClientID         int64     `json:"client_id"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            *string   `json:"email"`
	PhoneNumber      *string   `json:"phone_number"`
	Address          *string   `json:"address"`
	Relationship     *string   `json:"relationship"`
	RelationStatus   *string   `json:"relation_status"`
	CreatedAt        time.Time `json:"created_at"`
	IsVerified       bool      `json:"is_verified"`
	MedicalReports   bool      `json:"medical_reports"`
	IncidentsReports bool      `json:"incidents_reports"`
	GoalsReports     bool      `json:"goals_reports"`
}

// GetClientEmergencyContactResponse defines the response for getting a client emergency contact
type GetClientEmergencyContactResponse struct {
	ID               int64     `json:"id"`
	ClientID         int64     `json:"client_id"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            *string   `json:"email"`
	PhoneNumber      *string   `json:"phone_number"`
	Address          *string   `json:"address"`
	Relationship     *string   `json:"relationship"`
	RelationStatus   *string   `json:"relation_status"`
	CreatedAt        time.Time `json:"created_at"`
	IsVerified       bool      `json:"is_verified"`
	MedicalReports   bool      `json:"medical_reports"`
	IncidentsReports bool      `json:"incidents_reports"`
	GoalsReports     bool      `json:"goals_reports"`
}

// UpdateClientEmergencyContactParams defines the request for updating a client emergency contact
type UpdateClientEmergencyContactParams struct {
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	Email            *string `json:"email"`
	PhoneNumber      *string `json:"phone_number"`
	Address          *string `json:"address"`
	Relationship     *string `json:"relationship"`
	RelationStatus   *string `json:"relation_status"`
	MedicalReports   *bool   `json:"medical_reports"`
	IncidentsReports *bool   `json:"incidents_reports"`
	GoalsReports     *bool   `json:"goals_reports"`
}

// UpdateClientEmergencyContactResponse defines the response for updating a client emergency contact
type UpdateClientEmergencyContactResponse struct {
	ID               int64     `json:"id"`
	ClientID         int64     `json:"client_id"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            *string   `json:"email"`
	PhoneNumber      *string   `json:"phone_number"`
	Address          *string   `json:"address"`
	Relationship     *string   `json:"relationship"`
	RelationStatus   *string   `json:"relation_status"`
	CreatedAt        time.Time `json:"created_at"`
	IsVerified       bool      `json:"is_verified"`
	MedicalReports   bool      `json:"medical_reports"`
	IncidentsReports bool      `json:"incidents_reports"`
	GoalsReports     bool      `json:"goals_reports"`
}

// DeleteClientEmergencyContactResponse defines the response for deleting a client emergency contact
type DeleteClientEmergencyContactResponse struct {
	ID int64 `json:"id"`
}

// AssignEmployeeRequest defines the request for assigning an employee to a client
type AssignEmployeeRequest struct {
	EmployeeID int64     `json:"employee_id"`
	StartDate  time.Time `json:"start_date"`
	Role       string    `json:"role"`
}

// AssignEmployeeResponse defines the response for assigning an employee to a client
type AssignEmployeeResponse struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	EmployeeID int64     `json:"employee_id"`
	StartDate  time.Time `json:"start_date"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

// ListAssignedEmployeesRequest defines the request for listing assigned employees
type ListAssignedEmployeesRequest struct {
	pagination.Request
}

// ListAssignedEmployeesResponse defines the response for listing assigned employees
type ListAssignedEmployeesResponse struct {
	ID           int64     `json:"id"`
	ClientID     int64     `json:"client_id"`
	EmployeeID   int64     `json:"employee_id"`
	StartDate    time.Time `json:"start_date"`
	Role         string    `json:"role"`
	EmployeeName string    `json:"employee_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetAssignedEmployeeResponse defines the response for getting an assigned employee
type GetAssignedEmployeeResponse struct {
	ID           int64     `json:"id"`
	ClientID     int64     `json:"client_id"`
	EmployeeID   int64     `json:"employee_id"`
	StartDate    time.Time `json:"start_date"`
	Role         string    `json:"role"`
	EmployeeName string    `json:"employee_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// UpdateAssignedEmployeeRequest defines the request for updating an assigned employee
type UpdateAssignedEmployeeRequest struct {
	EmployeeID *int64    `json:"employee_id"`
	StartDate  time.Time `json:"start_date"`
	Role       *string   `json:"role"`
}

// UpdateAssignedEmployeeResponse defines the response for updating an assigned employee
type UpdateAssignedEmployeeResponse struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	EmployeeID int64     `json:"employee_id"`
	StartDate  time.Time `json:"start_date"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

// DeleteAssignedEmployeeResponse defines the response for deleting an assigned employee
type DeleteAssignedEmployeeResponse struct {
	ID int64 `json:"id"`
}

// GetClientRelatedEmailsResponse defines the response for getting client related emails
type GetClientRelatedEmailsResponse struct {
	Emails []string `json:"emails"`
}

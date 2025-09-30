package clientp

import "time"

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

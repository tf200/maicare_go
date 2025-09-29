package clientp

import (
	"maicare_go/pagination"
	"time"
)

type DiagnosisMedicationCreate struct {
	Name             string    `json:"name"`
	Dosage           string    `json:"dosage"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Notes            *string   `json:"notes"`
	SelfAdministered bool      `json:"self_administered"`
	AdministeredByID *int64    `json:"administered_by_id"`
	IsCritical       bool      `json:"is_critical"`
}

// CreateClientDiagnosisRequest defines the request for creating a client diagnosis
type CreateClientDiagnosisRequest struct {
	Title               *string                     `json:"title"`
	DiagnosisCode       string                      `json:"diagnosis_code"`
	Description         string                      `json:"description"`
	Severity            *string                     `json:"severity"`
	Status              string                      `json:"status"`
	DiagnosingClinician *string                     `json:"diagnosing_clinician"`
	Notes               *string                     `json:"notes"`
	Medications         []DiagnosisMedicationCreate `json:"medications"`
}

// CreateClientDiagnosisResponse defines the response for creating a client diagnosis
type CreateClientDiagnosisResponse struct {
	ID                  int64     `json:"id"`
	Title               *string   `json:"title"`
	ClientID            int64     `json:"client_id"`
	DiagnosisCode       string    `json:"diagnosis_code"`
	Description         string    `json:"description"`
	Severity            *string   `json:"severity"`
	Status              string    `json:"status"`
	DiagnosingClinician *string   `json:"diagnosing_clinician"`
	Notes               *string   `json:"notes"`
	CreatedAt           time.Time `json:"created_at"`
}

// ListClientDiagnosesRequest defines the request for listing client diagnoses
type ListClientDiagnosesRequest struct {
	pagination.Request
}

type DiagnosisMedicationList struct {
	ID               int64     `json:"id"`
	DiagnosisID      *int64    `json:"diagnosis_id"`
	Name             string    `json:"name"`
	Dosage           string    `json:"dosage"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Notes            *string   `json:"notes"`
	SelfAdministered bool      `json:"self_administered"`
	AdministeredByID *int64    `json:"administered_by_id"`
	IsCritical       bool      `json:"is_critical"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedAt        time.Time `json:"created_at"`
}

// ListClientDiagnosesResponse defines the response for listing client diagnoses
type ListClientDiagnosesResponse struct {
	ID                  int64                     `json:"id"`
	Title               *string                   `json:"title"`
	ClientID            int64                     `json:"client_id"`
	DiagnosisCode       string                    `json:"diagnosis_code"`
	Description         string                    `json:"description"`
	Severity            *string                   `json:"severity"`
	Status              string                    `json:"status"`
	DiagnosingClinician *string                   `json:"diagnosing_clinician"`
	Notes               *string                   `json:"notes"`
	CreatedAt           time.Time                 `json:"created_at"`
	Medications         []DiagnosisMedicationList `json:"medications"`
}

// GetClientDiagnosisResponse defines the response for getting a client diagnosis
type GetClientDiagnosisResponse struct {
	ID                  int64                     `json:"id"`
	Title               *string                   `json:"title"`
	ClientID            int64                     `json:"client_id"`
	DiagnosisCode       string                    `json:"diagnosis_code"`
	Description         string                    `json:"description"`
	DateOfDiagnosis     time.Time                 `json:"date_of_diagnosis"`
	Severity            *string                   `json:"severity"`
	Status              string                    `json:"status"`
	DiagnosingClinician *string                   `json:"diagnosing_clinician"`
	Notes               *string                   `json:"notes"`
	CreatedAt           time.Time                 `json:"created_at"`
	Medications         []DiagnosisMedicationList `json:"medications"`
}

// UpdateClientDiagnosisApi updates a client diagnosis
type UpdateClientDiagnosisRequest struct {
	Title               *string                     `json:"title"`
	DiagnosisCode       *string                     `json:"diagnosis_code"`
	Description         *string                     `json:"description"`
	Severity            *string                     `json:"severity"`
	Status              *string                     `json:"status"`
	DiagnosingClinician *string                     `json:"diagnosing_clinician"`
	Notes               *string                     `json:"notes"`
	MedicationIDs       []DiagnosisMedicationCreate `json:"medications"`
}

// UpdateClientDiagnosisApi updates a client diagnosis
type UpdateClientDiagnosisResponse struct {
	ID                  int64     `json:"id"`
	Title               *string   `json:"title"`
	ClientID            int64     `json:"client_id"`
	DiagnosisCode       string    `json:"diagnosis_code"`
	Description         string    `json:"description"`
	Severity            *string   `json:"severity"`
	Status              string    `json:"status"`
	DiagnosingClinician *string   `json:"diagnosing_clinician"`
	Notes               *string   `json:"notes"`
	CreatedAt           time.Time `json:"created_at"`
}

// DeleteClientDiagnosisResponse defines the response for deleting a client diagnosis
type DeleteClientDiagnosisResponse struct {
	ID int64 `json:"id"`
}

// CreateclientMedicationRequest defines the request for creating a client medication
type CreateclientMedicationRequest struct {
	Name             string    `json:"name"`
	Dosage           string    `json:"dosage"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Notes            *string   `json:"notes"`
	SelfAdministered bool      `json:"self_administered"`
	AdministeredByID *int64    `json:"administered_by_id"`
	IsCritical       bool      `json:"is_critical"`
}

// CreateClientMedicationResponse defines the response for creating a client medication
type CreateClientMedicationResponse struct {
	ID               int64     `json:"id"`
	DiagnosisID      *int64    `json:"diagnosis_id"`
	Name             string    `json:"name"`
	Dosage           string    `json:"dosage"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Notes            *string   `json:"notes"`
	SelfAdministered bool      `json:"self_administered"`
	AdministeredByID *int64    `json:"administered_by_id"`
	IsCritical       bool      `json:"is_critical"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedAt        time.Time `json:"created_at"`
}

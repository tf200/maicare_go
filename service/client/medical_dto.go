package clientp

import (
	"encoding/json"
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

// =====================
// Diagnoses
// =====================

type CreateClientDiagnosisRequest struct {
	CodeSystem          string     `json:"code_system"`
	Code                string     `json:"code"`
	Title               *string    `json:"title"`
	Description         *string    `json:"description"`
	Status              *string    `json:"status"`
	Severity            *string    `json:"severity"`
	DiagnosedOn         *time.Time `json:"diagnosed_on"`
	ResolvedOn          *time.Time `json:"resolved_on"`
	DiagnosingClinician *string    `json:"diagnosing_clinician"`
	Notes               *string    `json:"notes"`
}

type ClientDiagnosisResponse struct {
	ID                  uuid.UUID  `json:"id"`
	ClientID            uuid.UUID  `json:"client_id"`
	CodeSystem          string     `json:"code_system"`
	Code                string     `json:"code"`
	Title               *string    `json:"title"`
	Description         *string    `json:"description"`
	Status              string     `json:"status"`
	Severity            string     `json:"severity"`
	DiagnosedOn         *time.Time `json:"diagnosed_on"`
	ResolvedOn          *time.Time `json:"resolved_on"`
	DiagnosingClinician *string    `json:"diagnosing_clinician"`
	Notes               *string    `json:"notes"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type ListClientDiagnosesRequest struct {
	pagination.Request
}

type UpdateClientDiagnosisRequest struct {
	CodeSystem          *string    `json:"code_system"`
	Code                *string    `json:"code"`
	Title               *string    `json:"title"`
	Description         *string    `json:"description"`
	Status              *string    `json:"status"`
	Severity            *string    `json:"severity"`
	DiagnosedOn         *time.Time `json:"diagnosed_on"`
	ResolvedOn          *time.Time `json:"resolved_on"`
	DiagnosingClinician *string    `json:"diagnosing_clinician"`
	Notes               *string    `json:"notes"`
}

type DeleteClientDiagnosisResponse struct {
	ID uuid.UUID `json:"id"`
}

// =====================
// Medication Orders
// =====================

type CreateClientMedicationOrderRequest struct {
	DiagnosisID           *uuid.UUID      `json:"diagnosis_id"`
	MedicationName        string          `json:"medication_name"`
	DosageText            string          `json:"dosage_text"`
	DoseAmount            *float64        `json:"dose_amount"`
	DoseUnit              *string         `json:"dose_unit"`
	Route                 *string         `json:"route"`
	FrequencyText         *string         `json:"frequency_text"`
	Schedule              json.RawMessage `json:"schedule"`
	IsPrn                 bool            `json:"is_prn"`
	PrnIndication         *string         `json:"prn_indication"`
	MaxDosesPer24h        *int32          `json:"max_doses_per_24h"`
	StartDate             time.Time       `json:"start_date"`
	EndDate               *time.Time      `json:"end_date"`
	Status                *string         `json:"status"`
	AdminMode             *string         `json:"admin_mode"`
	ResponsibleEmployeeID *uuid.UUID      `json:"responsible_employee_id"`
	IsCritical            bool            `json:"is_critical"`
	Notes                 *string         `json:"notes"`
	SourceAttachmentUUID  *uuid.UUID      `json:"source_attachment_uuid"`
}

type ClientMedicationOrderResponse struct {
	ID                           uuid.UUID       `json:"id"`
	ClientID                     uuid.UUID       `json:"client_id"`
	DiagnosisID                  *uuid.UUID      `json:"diagnosis_id"`
	MedicationName               string          `json:"medication_name"`
	DosageText                   string          `json:"dosage_text"`
	DoseAmount                   *float64        `json:"dose_amount"`
	DoseUnit                     *string         `json:"dose_unit"`
	Route                        *string         `json:"route"`
	FrequencyText                *string         `json:"frequency_text"`
	Schedule                     json.RawMessage `json:"schedule"`
	IsPrn                        bool            `json:"is_prn"`
	PrnIndication                *string         `json:"prn_indication"`
	MaxDosesPer24h               *int32          `json:"max_doses_per_24h"`
	StartDate                    time.Time       `json:"start_date"`
	EndDate                      *time.Time      `json:"end_date"`
	Status                       string          `json:"status"`
	AdminMode                    string          `json:"admin_mode"`
	ResponsibleEmployeeID        *uuid.UUID      `json:"responsible_employee_id"`
	ResponsibleEmployeeFirstName *string         `json:"responsible_employee_first_name"`
	ResponsibleEmployeeLastName  *string         `json:"responsible_employee_last_name"`
	IsCritical                   bool            `json:"is_critical"`
	Notes                        *string         `json:"notes"`
	SourceAttachmentUUID         *uuid.UUID      `json:"source_attachment_uuid"`
	DiagnosisTitle               *string         `json:"diagnosis_title"`
	DiagnosisCodeSystem          *string         `json:"diagnosis_code_system"`
	DiagnosisCode                *string         `json:"diagnosis_code"`
	CreatedAt                    time.Time       `json:"created_at"`
	UpdatedAt                    time.Time       `json:"updated_at"`
}

type ListClientMedicationOrdersRequest struct {
	pagination.Request
	Status      *string    `json:"status" form:"status"`
	AdminMode   *string    `json:"admin_mode" form:"admin_mode"`
	DiagnosisID *uuid.UUID `json:"diagnosis_id" form:"diagnosis_id"`
	Search      *string    `json:"search" form:"search"`
}

type UpdateClientMedicationOrderRequest struct {
	DiagnosisID           *uuid.UUID      `json:"diagnosis_id"`
	MedicationName        *string         `json:"medication_name"`
	DosageText            *string         `json:"dosage_text"`
	DoseAmount            *float64        `json:"dose_amount"`
	DoseUnit              *string         `json:"dose_unit"`
	Route                 *string         `json:"route"`
	FrequencyText         *string         `json:"frequency_text"`
	Schedule              json.RawMessage `json:"schedule"`
	IsPrn                 *bool           `json:"is_prn"`
	PrnIndication         *string         `json:"prn_indication"`
	MaxDosesPer24h        *int32          `json:"max_doses_per_24h"`
	StartDate             *time.Time      `json:"start_date"`
	EndDate               *time.Time      `json:"end_date"`
	Status                *string         `json:"status"`
	AdminMode             *string         `json:"admin_mode"`
	ResponsibleEmployeeID *uuid.UUID      `json:"responsible_employee_id"`
	IsCritical            *bool           `json:"is_critical"`
	Notes                 *string         `json:"notes"`
	SourceAttachmentUUID  *uuid.UUID      `json:"source_attachment_uuid"`
}

type DeleteClientMedicationOrderResponse struct {
	ID uuid.UUID `json:"id"`
}

// =====================
// Overview
// =====================

type ClientMedicalOverviewResponse struct {
	Diagnoses        []ClientDiagnosisResponse       `json:"diagnoses"`
	MedicationOrders []ClientMedicationOrderResponse `json:"medication_orders"`
}

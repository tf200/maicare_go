package contract

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateContractTypeRequest defines the request for CreateContractType handler
type CreateContractTypeRequest struct {
	Name string `json:"name"`
}

// AttachmentDetail represents embedded attachment metadata in contract responses
type AttachmentDetail struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	DownloadURL string    `json:"download_url"`
}

// CreateContractTypeResponse defines the response for CreateContractType handler
type CreateContractTypeResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ListContractTypesResponse defines the response for ListContractTypes handler
type ListContractTypesResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// DeleteContractTypeResponse defines the response for DeleteContractType handler
type DeleteContractTypeResponse struct {
	ID uuid.UUID `json:"id"`
}

// CreateContractRequest defines the request for CreateContract handler
type CreateContractRequest struct {
	ClientID        uuid.UUID   `json:"client_id" binding:"required" example:"afc465cc-cddb-440b-9472-615bb07ec1d8"`
	TypeID          *uuid.UUID  `json:"type_id" example:"1"`
	StartDate       time.Time   `json:"start_date" binding:"required" example:"2023-01-01T00:00:00Z"`
	EndDate         time.Time   `json:"end_date" binding:"required" example:"2023-12-31T00:00:00Z"`
	ReminderPeriod  *int32      `json:"reminder_period" binding:"omitempty,gte=0" example:"30"`
	Vat             *int32      `json:"VAT" example:"21"`
	Price           float64     `json:"price" example:"100.50"`
	PriceTimeUnit   string      `json:"price_time_unit" binding:"required,oneof=minute hourly daily weekly" example:"weekly" enum:"minute,hourly,daily,weekly"`
	Hours           *float64    `json:"hours" example:"40"`
	HoursType       *string     `json:"hours_type" enum:"weekly,all_period" binding:"omitempty,oneof=weekly all_period" example:"weekly"`
	CareName        string      `json:"care_name" example:"Home Care"`
	CareType        string      `json:"care_type" binding:"required,oneof=ambulante accommodation" example:"ambulante" enum:"ambulante,accommodation"`
	SenderID        uuid.UUID   `json:"sender_id" binding:"required" example:"afc465cc-cddb-440b-9472-615bb07ec1d8"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act" binding:"required,oneof=WMO ZVW WLZ JW WPG" example:"WMO" enum:"WMO,ZVW,WLZ,JW,WPG"`
	FinancingOption string      `json:"financing_option" binding:"required,oneof=ZIN PGB" example:"ZIN" enum:"ZIN,PGB"`
}

// CreateContractResponse defines the response for CreateContract handler
type CreateContractResponse struct {
	ID              uuid.UUID          `json:"id"`
	TypeID          *uuid.UUID         `json:"type_id"`
	Status          string             `json:"status"`
	StartDate       time.Time          `json:"start_date"`
	EndDate         time.Time          `json:"end_date"`
	ReminderPeriod  int32              `json:"reminder_period"`
	Vat             *int32             `json:"VAT"`
	Price           float64            `json:"price"`
	PriceTimeUnit   string             `json:"price_time_unit"`
	Hours           *float64           `json:"hours"`
	HoursType       *string            `json:"hours_type"`
	CareName        string             `json:"care_name"`
	CareType        string             `json:"care_type"`
	ClientID        uuid.UUID          `json:"client_id"`
	SenderID        uuid.UUID          `json:"sender_id"`
	AttachmentIds   []uuid.UUID        `json:"attachment_ids"`
	Attachments     []AttachmentDetail `json:"attachments"`
	FinancingAct    string             `json:"financing_act"`
	FinancingOption string             `json:"financing_option"`
	DepartureReason *string            `json:"departure_reason"`
	DepartureReport *string            `json:"departure_report"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
}

// ListClientContractsRequest defines the request for ListClientContracts handler
type ListClientContractsRequest struct {
	pagination.Request
}

// ListClientContractsResponse defines the response for ListClientContracts handler
type ListClientContractsResponse struct {
	ID              uuid.UUID   `json:"id"`
	TypeID          *uuid.UUID  `json:"type_id"`
	Status          string      `json:"status"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  int32       `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           float64     `json:"price"`
	PriceTimeUnit   string      `json:"price_time_unit"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type"`
	ClientID        uuid.UUID   `json:"client_id"`
	ClientFirstName string      `json:"client_first_name"`
	ClientLastName  string      `json:"client_last_name"`
	SenderID        uuid.UUID   `json:"sender_id"`
	SenderName      *string     `json:"sender_name"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
	DepartureReason *string     `json:"departure_reason"`
	DepartureReport *string     `json:"departure_report"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`
}

// UpdateContractRequest defines the request for UpdateContract handler
type UpdateContractRequest struct {
	TypeID          *uuid.UUID  `json:"type_id"`
	StartDate       *time.Time  `json:"start_date"`
	EndDate         *time.Time  `json:"end_date"`
	ReminderPeriod  *int32      `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           *float64    `json:"price"`
	PriceTimeUnit   *string     `json:"price_time_unit"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        *string     `json:"care_name"`
	CareType        *string     `json:"care_type"`
	SenderID        *uuid.UUID  `json:"sender_id"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    *string     `json:"financing_act"`
	FinancingOption *string     `json:"financing_option"`
	Status          *string     `json:"status"`
}

// UpdateContractResponse defines the response for UpdateContract handler
type UpdateContractResponse struct {
	ID              uuid.UUID   `json:"id"`
	TypeID          *uuid.UUID  `json:"type_id"`
	Status          string      `json:"status"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  int32       `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           float64     `json:"price"`
	PriceFrequency  string      `json:"price_frequency"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type"`
	ClientID        uuid.UUID   `json:"client_id"`
	SenderID        uuid.UUID   `json:"sender_id"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
	DepartureReason *string     `json:"departure_reason"`
	DepartureReport *string     `json:"departure_report"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`
}

// UpdateContractStatusRequest defines the request for UpdateContractStatus handler
type UpdateContractStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=approved draft terminated stopped expired" example:"approved" enum:"approved,draft,terminated,stopped,expired"`
}

// UpdateContractStatusResponse defines the response for UpdateContractStatus handler
type UpdateContractStatusResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

// GetClientContractResponse defines the response for GetContract handler
type GetClientContractResponse struct {
	// Contract fields
	ID              uuid.UUID   `json:"id"`
	TypeID          *uuid.UUID  `json:"type_id"`
	TypeName        *string     `json:"type_name"`
	Status          string      `json:"status"`
	ApprovedAt      *time.Time  `json:"approved_at"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  int32       `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           float64     `json:"price"`
	PriceTimeUnit   string      `json:"price_time_unit"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
	DepartureReason *string     `json:"departure_reason"`
	DepartureReport *string     `json:"departure_report"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`

	// Client fields
	ClientID         uuid.UUID `json:"client_id"`
	ClientFirstName  string    `json:"client_first_name"`
	ClientLastName   string    `json:"client_last_name"`
	ClientFilenumber string    `json:"client_filenumber"`
	ClientBsn        *string   `json:"client_bsn"`

	// Sender fields
	SenderID                  uuid.UUID `json:"sender_id"`
	SenderName                string    `json:"sender_name"`
	SenderType                string    `json:"sender_type"`
	SenderStreet              *string   `json:"sender_street"`
	SenderHouseNumber         *string   `json:"sender_house_number"`
	SenderHouseNumberAddition *string   `json:"sender_house_number_addition"`
	SenderPostalCode          *string   `json:"sender_postal_code"`
	SenderCity                *string   `json:"sender_city"`
	SenderLand                *string   `json:"sender_land"`
	SenderKvknumber           *string   `json:"sender_kvknumber"`
	SenderBtwnumber           *string   `json:"sender_btwnumber"`
	SenderPhoneNumber         *string   `json:"sender_phone_number"`
	SenderClientNumber        *string   `json:"sender_client_number"`
	SenderEmailAddress        *string   `json:"sender_email_address"`
}

// ListContractsRequest defines the request for ListContracts handler
type ListContractsRequest struct {
	pagination.Request
	Search          *string    `form:"search" binding:"omitempty"`
	Status          []string   `form:"status" binding:"omitempty,dive,oneof=approved draft terminated stopped expired"`
	CareType        []string   `form:"care_type" binding:"omitempty,dive,oneof=ambulante accommodation"`
	FinancingAct    []string   `form:"financing_act" binding:"omitempty,dive,oneof=WMO ZVW WLZ JW WPG"`
	FinancingOption *string    `form:"financing_option" binding:"omitempty,oneof=ZIN PGB"`
	EndDateFrom     *time.Time `form:"end_date_from" time_format:"2006-01-02"`
	EndDateTo       *time.Time `form:"end_date_to" time_format:"2006-01-02"`
}

// ListContractsResponse defines the response for ListContracts handler
type ListContractsResponse struct {
	ID               uuid.UUID  `json:"id"`
	ClientID         uuid.UUID  `json:"client_id"`
	ClientFirstName  string     `json:"client_first_name"`
	ClientLastName   string     `json:"client_last_name"`
	ClientFilenumber string     `json:"client_filenumber"`
	SenderID         uuid.UUID  `json:"sender_id"`
	SenderName       string     `json:"sender_name"`
	CareName         string     `json:"care_name"`
	CareType         string     `json:"care_type"`
	Price            float64    `json:"price"`
	PriceTimeUnit    string     `json:"price_time_unit"`
	Hours            *float64   `json:"hours"`
	HoursType        *string    `json:"hours_type"`
	FinancingAct     string     `json:"financing_act"`
	FinancingOption  string     `json:"financing_option"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          time.Time  `json:"end_date"`
	DaysLeft         int32      `json:"days_left"`
	Status           string     `json:"status"`
	ApprovedAt       *time.Time `json:"approved_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// GetContractAuditLogResponse defines the response for GetContractAuditLog handler
type GetContractAuditLogResponse struct {
	AuditID            uuid.UUID          `json:"audit_id"`
	ContractID         uuid.UUID          `json:"contract_id"`
	Operation          string             `json:"operation"`
	ChangedBy          *uuid.UUID         `json:"changed_by"`
	ChangedAt          pgtype.Timestamptz `json:"changed_at"`
	OldValues          any                `json:"old_values"`
	NewValues          any                `json:"new_values"`
	ChangedFields      []string           `json:"changed_fields"`
	ChangedByFirstName *string            `json:"changed_by_first_name"`
	ChangedByLastName  *string            `json:"changed_by_last_name"`
}

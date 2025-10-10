package contract

import (
	"maicare_go/pagination"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateContractTypeRequest defines the request for CreateContractType handler
type CreateContractTypeRequest struct {
	Name string `json:"name"`
}

// CreateContractTypeResponse defines the response for CreateContractType handler
type CreateContractTypeResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ListContractTypesResponse defines the response for ListContractTypes handler
type ListContractTypesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// DeleteContractTypeResponse defines the response for DeleteContractType handler
type DeleteContractTypeResponse struct {
	ID int64 `json:"id"`
}

// CreateContractRequest defines the request for CreateContract handler
type CreateContractRequest struct {
	TypeID          *int64      `json:"type_id" example:"1"`
	StartDate       time.Time   `json:"start_date" example:"2023-01-01T00:00:00Z"`
	EndDate         time.Time   `json:"end_date" example:"2023-12-31T00:00:00Z"`
	ReminderPeriod  int32       `json:"reminder_period" example:"30"`
	Vat             *int32      `json:"VAT" example:"21"`
	Price           float64     `json:"price" example:"100.50"`
	PriceTimeUnit   string      `json:"price_time_unit" binding:"required,oneof=minute hourly daily weekly monthly yearly" example:"monthly" enum:"minute,hourly,daily,weekly,monthly,yearly"`
	Hours           *float64    `json:"hours" example:"40"`
	HoursType       *string     `json:"hours_type" enum:"weekly,all_period"`
	CareName        string      `json:"care_name" example:"Home Care"`
	CareType        string      `json:"care_type" binding:"required,oneof=ambulante accommodation" example:"ambulante" enum:"ambulante,accommodation"`
	SenderID        *int64      `json:"sender_id" example:"2"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act" binding:"required,oneof=WMO ZVW WLZ JW WPG" example:"WMO" enum:"WMO,ZVW,WLZ,JW,WPG"`
	FinancingOption string      `json:"financing_option" binding:"required,oneof=ZIN PGB" example:"ZIN" enum:"ZIN,PGB"`
}

// CreateContractResponse defines the response for CreateContract handler
type CreateContractResponse struct {
	ID              int64              `json:"id"`
	TypeID          *int64             `json:"type_id"`
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
	ClientID        int64              `json:"client_id"`
	SenderID        *int64             `json:"sender_id"`
	AttachmentIds   []uuid.UUID        `json:"attachment_ids"`
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
	ID              int64       `json:"id"`
	TypeID          *int64      `json:"type_id"`
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
	ClientID        int64       `json:"client_id"`
	ClientFirstName string      `json:"client_first_name"`
	ClientLastName  string      `json:"client_last_name"`
	SenderID        *int64      `json:"sender_id"`
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
	TypeID          *int64      `json:"type_id"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  *int32      `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           *float64    `json:"price"`
	PriceTimeUnit   *string     `json:"price_time_unit"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        *string     `json:"care_name"`
	CareType        *string     `json:"care_type"`
	SenderID        *int64      `json:"sender_id"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    *string     `json:"financing_act"`
	FinancingOption *string     `json:"financing_option"`
	Status          *string     `json:"status"`
}

// UpdateContractResponse defines the response for UpdateContract handler
type UpdateContractResponse struct {
	ID              int64       `json:"id"`
	TypeID          *int64      `json:"type_id"`
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
	ClientID        int64       `json:"client_id"`
	SenderID        *int64      `json:"sender_id"`
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
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

// GetClientContractResponse defines the response for GetContract handler
type GetClientContractResponse struct {
	ID              int64       `json:"id"`
	TypeID          *int64      `json:"type_id"`
	TypeName        string      `json:"type_name"`
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
	ClientID        int64       `json:"client_id"`
	ClientFirstName string      `json:"client_first_name"`
	ClientLastName  string      `json:"client_last_name"`
	SenderID        *int64      `json:"sender_id"`
	SenderName      *string     `json:"sender_name"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
	DepartureReason *string     `json:"departure_reason"`
	DepartureReport *string     `json:"departure_report"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`
}

// ListContractsRequest defines the request for ListContracts handler
type ListContractsRequest struct {
	pagination.Request
	Search          *string  `form:"search" binding:"omitempty"`
	Status          []string `form:"status" binding:"omitempty,dive,oneof=approved draft terminated stopped"`
	CareType        []string `form:"care_type" binding:"omitempty,dive,oneof=ambulante accommodation"`
	FinancingAct    []string `form:"financing_act" binding:"omitempty,dive,oneof=WMO ZVW WLZ JW WPG"`
	FinancingOption []string `form:"financing_option" binding:"omitempty,dive,oneof=ZIN PGB"`
}

// ListContractsResponse defines the response for ListContracts handler
type ListContractsResponse struct {
	ID              int64     `json:"id"`
	ClientID        int64     `json:"client_id"`
	Status          string    `json:"status"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Price           float64   `json:"price"`
	PriceTimeUnit   string    `json:"price_time_unit"`
	CareName        string    `json:"care_name"`
	CareType        string    `json:"care_type"`
	FinancingAct    string    `json:"financing_act"`
	FinancingOption string    `json:"financing_option"`
	CreatedAt       time.Time `json:"created_at"`
	SenderID        *int64    `json:"sender_id"`
	SenderName      *string   `json:"sender_name"`
	ClientFirstName string    `json:"client_first_name"`
	ClientLastName  string    `json:"client_last_name"`
}

// GetContractAuditLogResponse defines the response for GetContractAuditLog handler
type GetContractAuditLogResponse struct {
	AuditID            int64              `json:"audit_id"`
	ContractID         int64              `json:"contract_id"`
	Operation          string             `json:"operation"`
	ChangedBy          *int64             `json:"changed_by"`
	ChangedAt          pgtype.Timestamptz `json:"changed_at"`
	OldValues          any                `json:"old_values"`
	NewValues          any                `json:"new_values"`
	ChangedFields      []string           `json:"changed_fields"`
	ChangedByFirstName *string            `json:"changed_by_first_name"`
	ChangedByLastName  *string            `json:"changed_by_last_name"`
}
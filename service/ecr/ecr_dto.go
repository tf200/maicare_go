package ecr

import (
	"maicare_go/pagination"
	"time"

	"github.com/google/uuid"
)

// DischargeOverviewRequest defines the request for the DischargeOverview handler.
type DischargeOverviewRequest struct {
	pagination.Request
	FilterType string `form:"filter_type" binding:"required,oneof=urgent contract status_change all"`
}

// DischargeOverviewResponse defines the response for the DischargeOverview handler.
type DischargeOverviewResponse struct {
	ID                 int64     `json:"id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	CurrentStatus      *string   `json:"current_status"`
	ScheduledStatus    *string   `json:"scheduled_status"`
	StatusChangeReason *string   `json:"status_change_reason"`
	StatusChangeDate   time.Time `json:"status_change_date"`
	ContractEndDate    time.Time `json:"contract_end_date"`
	ContractStatus     *string   `json:"contract_status"`
	DepartureReason    *string   `json:"departure_reason"`
	FollowUpPlan       *string   `json:"follow_up_plan"`
	DischargeType      string    `json:"discharge_type"`
}

// TotalDischargeCountResponse defines the response for the TotalDischargeCount handler.
type TotalDischargeCountResponse struct {
	TotalDischargeCount int64 `json:"total_discharge_count"`
	UrgentCasesCount    int64 `json:"urgent_cases_count"`
	StatusChangeCount   int64 `json:"status_change_count"`
	ContractEndCount    int64 `json:"contract_end_count"`
}

// ListEmployeesByContractEndDateResponse defines the response for the ListEmployeesByContractEndDate handler.
type ListEmployeesByContractEndDateResponse struct {
	ID                int64     `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Position          *string   `json:"position"`
	Department        *string   `json:"department"`
	EmployeeNumber    *string   `json:"employee_number"`
	EmploymentNumber  *string   `json:"employment_number"`
	Email             string    `json:"email"`
	ContractStartDate time.Time `json:"contract_start_date"`
	ContractEndDate   time.Time `json:"contract_end_date"`
	ContractType      *string   `json:"contract_type"`
}

// ListLatestPaymentsResponse defines the response for the ListLatestPayments API.
type ListLatestPaymentsResponse struct {
	InvoiceID     int64     `json:"invoice_id"`
	InvoiceNumber string    `json:"invoice_number"`
	PaymentMethod *string   `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
	Amount        float64   `json:"amount"`
	PaymentDate   time.Time `json:"payment_date"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ListUpcomingAppointmentsResponse defines the response for the ListUpcomingAppointments API.
type ListUpcomingAppointmentsResponse struct {
	ID          uuid.UUID `json:"id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    *string   `json:"location"`
	Description *string   `json:"description"`
}
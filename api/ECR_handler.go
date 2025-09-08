package api

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ======================================================
// ADMIN DASHBOARD
// ======================================================

// CreateSchedueledClientStatusChangeRequest defines the request for the CreateSchedueledClientStatusChange handler.
type DischargeOverviewRequest struct {
	pagination.Request
	FilterType string `form:"filter_type" binding:"required,oneof=urgent contract status_change all"`
}

// GetParams returns the pagination params for the request.
type DischargeOverviewResponse struct {
	ID                 int64     `json:"id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	CurrentStatus      *string   `json:"current_status"`
	ScheduledStatus    string    `json:"scheduled_status"`
	StatusChangeReason *string   `json:"status_change_reason"`
	StatusChangeDate   time.Time `json:"status_change_date"`
	ContractEndDate    time.Time `json:"contract_end_date"`
	ContractStatus     *string   `json:"contract_status"`
	DepartureReason    *string   `json:"departure_reason"`
	FollowUpPlan       *string   `json:"follow_up_plan"`
	DischargeType      string    `json:"discharge_type"`
}

// DischargeOverviewApi returns a list of discharge overview.
// @Summary Returns a list of discharge overview.
// @Tags ECR
// @Produce json
// @Param page query integer false "page number"
// @Param page_size query integer false "number of items per page"
// @Param filter_type query string false "filter type" Enums(urgent, contract, status_change, all)
// @Success 200 {object} Response[pagination.Response[DischargeOverviewResponse]]
// @Failure 400 {object} Response[any]
// @Router /ecr/discharge_overview [get]
func (server *Server) DischargeOverviewApi(ctx *gin.Context) {
	var req DischargeOverviewRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "DischargeOverviewApi", "Failed to bind query parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("failed to bind query parameters")))
		return
	}

	params := req.GetParams()

	overview, err := server.store.DischargeOverview(ctx, db.DischargeOverviewParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DischargeOverviewApi", "Failed to get discharge overview", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get discharge overview")))
		return
	}

	overviewRes := []DischargeOverviewResponse{}
	for _, item := range overview {
		overviewRes = append(overviewRes, DischargeOverviewResponse{
			ID:                 item.ID,
			FirstName:          item.FirstName,
			LastName:           item.LastName,
			CurrentStatus:      item.CurrentStatus,
			ScheduledStatus:    item.ScheduledStatus,
			StatusChangeReason: item.StatusChangeReason,
			StatusChangeDate:   item.StatusChangeDate.Time,
			ContractEndDate:    item.ContractEndDate.Time,
			ContractStatus:     item.ContractStatus,
			DepartureReason:    item.DepartureReason,
			FollowUpPlan:       item.FollowUpPlan,
			DischargeType:      item.DischargeType,
		})
	}

	count, err := server.store.TotalDischargeCount(context.Background())
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DischargeOverviewApi", "Failed to get total discharge count", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get total discharge count")))
		return
	}

	pag := pagination.NewResponse(ctx, req.Request, overviewRes, count)
	res := SuccessResponse(pag, "Discharge overview Retrieved Successfully")
	ctx.JSON(http.StatusOK, res)

}

// TotalDischargeCountResponse defines the response for the TotalDischargeCount handler.
type TotalDischargeCountResponse struct {
	TotalDischargeCount int64 `json:"total_discharge_count"`
	UrgentCasesCount    int64 `json:"urgent_cases_count"`
	StatusChangeCount   int64 `json:"status_change_count"`
	ContractEndCount    int64 `json:"contract_end_count"`
}

// TotalDischargeCount returns the total discharge count.
// @Summary Returns the total discharge count.
// @Tags ECR
// @Produce json
// @Success 200 {object} Response[TotalDischargeCountResponse]
// @Failure 400 {object} Response[any]
// @Router /ecr/total_discharge_count [get]
func (server *Server) TotalDischargeCountApi(ctx *gin.Context) {
	// Initialize WaitGroup to track 4 Go routines
	var wg sync.WaitGroup
	wg.Add(4)

	// Variables to hold results and errors for each query
	var totalDischargeCount int64
	var totalDischargeErr error
	var urgentCasesCount int64
	var urgentCasesErr error
	var statusChangeCount int64
	var statusChangeErr error
	var contractEndingCount int64
	var contractEndingErr error

	// Get the request context to respect timeouts/cancellations
	reqCtx := ctx.Request.Context()

	// Go routine for TotalDischargeCount
	go func() {
		defer wg.Done() // Signal completion when done
		count, err := server.store.TotalDischargeCount(reqCtx)
		totalDischargeCount = count
		totalDischargeErr = err
	}()

	// Go routine for UrgentCasesCount
	go func() {
		defer wg.Done()
		count, err := server.store.UrgentCasesCount(reqCtx)
		urgentCasesCount = count
		urgentCasesErr = err
	}()

	// Go routine for StatusChangeCount
	go func() {
		defer wg.Done()
		count, err := server.store.StatusChangeCount(reqCtx)
		statusChangeCount = count
		statusChangeErr = err
	}()

	// Go routine for ContractEndCount
	go func() {
		defer wg.Done()
		count, err := server.store.ContractEndCount(reqCtx)
		contractEndingCount = count
		contractEndingErr = err
	}()

	// Wait for all Go routines to complete
	wg.Wait()

	// Check errors and return the first one encountered
	if totalDischargeErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(totalDischargeErr))
		return
	}
	if urgentCasesErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(urgentCasesErr))
		return
	}
	if statusChangeErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(statusChangeErr))
		return
	}
	if contractEndingErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(contractEndingErr))
		return
	}

	// All queries succeeded, construct and send the response
	res := SuccessResponse(TotalDischargeCountResponse{
		TotalDischargeCount: totalDischargeCount,
		UrgentCasesCount:    urgentCasesCount,
		StatusChangeCount:   statusChangeCount,
		ContractEndCount:    contractEndingCount,
	}, "Total Discharge Count Retrieved Successfully")
	ctx.JSON(http.StatusOK, res)
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

// ListEmployeesByContractEndDateApi handles the API request for listing employees by contract end date.
// @Summary Lists employees whose contracts are ending soon.
// @Description This endpoint retrieves a list of employees whose contracts are approaching their end date.
// @Tags ECR
// @Produce json
// @Success 200 {object} Response[[]ListEmployeesByContractEndDateResponse]
// @Failure 500 {object} Response[any]
// @Router /ecr/employee_ending_contract [get]
func (server *Server) ListEmployeesByContractEndDateApi(ctx *gin.Context) {
	employees, err := server.store.ListEmployeesByContractEndDate(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListEmployeesByContractEndDateApi", "Failed to list employees by contract end date", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := []ListEmployeesByContractEndDateResponse{}
	for _, emp := range employees {
		response = append(response, ListEmployeesByContractEndDateResponse{
			ID:                emp.ID,
			FirstName:         emp.FirstName,
			LastName:          emp.LastName,
			Position:          emp.Position,
			Department:        emp.Department,
			EmployeeNumber:    emp.EmployeeNumber,
			EmploymentNumber:  emp.EmploymentNumber,
			Email:             emp.Email,
			ContractStartDate: emp.ContractStartDate.Time,
			ContractEndDate:   emp.ContractEndDate.Time,
			ContractType:      emp.ContractType,
		})
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Employees listed successfully"))
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

// ListLatestPaymentsApi handles the API request for listing the latest payments.
// @Summary Lists the latest payments.
// @Description This endpoint retrieves a list of the latest payments.
// @Tags ECR
// @Produce json
// @Success 200 {object} Response[[]ListLatestPaymentsResponse]
// @Failure 500 {object} Response[any]
// @Router /ecr/latest_payments [get]
func (server *Server) ListLatestPaymentsApi(ctx *gin.Context) {
	latestPayments, err := server.store.ListLatestPayments(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListLatestPaymentsApi", "Failed to list latest payments", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := []ListLatestPaymentsResponse{}
	for _, payment := range latestPayments {
		response = append(response, ListLatestPaymentsResponse{
			InvoiceID:     payment.InvoiceID,
			InvoiceNumber: payment.InvoiceNumber,
			PaymentMethod: payment.PaymentMethod,
			PaymentStatus: payment.PaymentStatus,
			Amount:        payment.Amount,
			PaymentDate:   payment.PaymentDate.Time,
			UpdatedAt:     payment.UpdatedAt.Time,
		})
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Latest payments listed successfully"))
}

// ListUpcomingAppointmentsResponse defines the response for the ListUpcomingAppointments API.
type ListUpcomingAppointmentsResponse struct {
	ID          uuid.UUID `json:"id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    *string   `json:"location"`
	Description *string   `json:"description"`
}

// ListUpcomingAppointmentsApi handles the API request for listing upcoming appointments.
// @Summary Lists upcoming appointments.
// @Description This endpoint retrieves a list of upcoming appointments.
// @Tags ECR
// @Produce json
// @Success 200 {object} Response[[]ListUpcomingAppointmentsResponse]
// @Failure 500 {object} Response[any]
// @Router /ecr/upcoming_appointments [get]
func (server *Server) ListUpcomingAppointmentsApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListUpcomingAppointmentsApi", "Failed to get auth payload", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("failed to get auth payload")))
		return
	}

	appointments, err := server.store.ListUpcomingAppointments(ctx, &payload.EmployeeID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListUpcomingAppointmentsApi", "Failed to list upcoming appointments", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := []ListUpcomingAppointmentsResponse{}
	for _, appt := range appointments {
		response = append(response, ListUpcomingAppointmentsResponse{
			ID:          appt.ID,
			StartTime:   appt.StartTime.Time,
			EndTime:     appt.EndTime.Time,
			Location:    appt.Location,
			Description: appt.Description,
		})
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Upcoming appointments listed successfully"))
}

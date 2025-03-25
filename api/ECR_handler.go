package api

import (
	"context"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateSchedueledClientStatusChangeRequest defines the request for the CreateSchedueledClientStatusChange handler.
type DischargeOverviewRequest struct {
	pagination.Request
	FilterType string `form:"filter_type" binding:"required ,oneof=urgent contract status_change all"`
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
// @Success 200 {object} DischargeOverviewResponse
// @Failure 400 {object} Response[any]
// @Router /ecr/discharge_overview [get]
func (server *Server) DischargeOverviewApi(ctx *gin.Context) {
	var req DischargeOverviewRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	overview, err := server.store.DischargeOverview(ctx, db.DischargeOverviewParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	overviewRes := make([]DischargeOverviewResponse, len(overview))
	for i, item := range overview {
		overviewRes[i] = DischargeOverviewResponse{
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
		}
	}

	count, err := server.store.TotalDischargeCount(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
// @Success 200 {object} TotalDischargeCountResponse
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

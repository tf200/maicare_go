package api

import (
	"context"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateSchedueledClientStatusChangeRequest defines the request for the CreateSchedueledClientStatusChange handler.
type DischargeOverviewRequest struct {
	pagination.Request
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
}

// TotalDischargeCount returns the total discharge count.
// @Summary Returns the total discharge count.
// @Tags ECR
// @Produce json
// @Success 200 {object} TotalDischargeCountResponse
// @Failure 400 {object} Response[any]
// @Router /ecr/total_discharge_count [get]
func (server *Server) TotalDischargeCountApi(ctx *gin.Context) {
	count, err := server.store.TotalDischargeCount(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(TotalDischargeCountResponse{
		TotalDischargeCount: count,
	}, "Total Discharge Count Retrieved Successfully")
	ctx.JSON(http.StatusOK, res)
}

// UrgentCasesCountResponse defines the response for the UrgentCasesCount handler.
type UrgentCasesCountResponse struct {
	UrgentCasesCount int64 `json:"urgent_cases_count"`
}

// UrgentCasesCount returns the urgent cases count.
// @Summary Returns the urgent cases count.
// @Tags ECR
// @Produce json
// @Success 200 {object} UrgentCasesCountResponse
// @Failure 400 {object} Response[any]
// @Router /ecr/urgent_cases_count [get]
func (server *Server) UrgentCasesCountApi(ctx *gin.Context) {
	count, err := server.store.UrgentCasesCount(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(UrgentCasesCountResponse{
		UrgentCasesCount: count,
	}, "Urgent Cases Count Retrieved Successfully")
	ctx.JSON(http.StatusOK, res)
}

// StatusChangeCountResponse defines the response for the StatusChangeCount handler.
type StatusChangeCountResponse struct {
	StatusChangeCount int64 `json:"status_change_count"`
}

// StatusChangeCount returns the status change count.
// @Summary Returns the status change count.
// @Tags ECR
// @Produce json
// @Success 200 {object} StatusChangeCountResponse
// @Failure 400 {object} Response[any]
// @Router /ecr/status_change_count [get]
func (server *Server) StatusChangeCountApi(ctx *gin.Context) {
	count, err := server.store.StatusChangeCount(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(StatusChangeCountResponse{
		StatusChangeCount: count,
	}, "Status Change Count Retrieved Successfully")
	ctx.JSON(http.StatusOK, res)
}

// ContractEndCountResponse defines the response for the ContractEndCount handler.
type ContractEndCountResponse struct {
	ContractEndCount int64 `json:"contract_end_count"`
}

// ContractEndCount returns the contract end count.
// @Summary Returns the contract end count.
// @Tags ECR
// @Produce json
// @Success 200 {object} ContractEndCountResponse
// @Failure 400 {object} Response[any]
// @Router /ecr/contract_end_count [get]
func (server *Server) ContractEndCountApi(ctx *gin.Context) {
	count, err := server.store.ContractEndCount(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(ContractEndCountResponse{
		ContractEndCount: count,
	}, "Contract End Count Retrieved Successfully")
	ctx.JSON(http.StatusOK, res)
}

package ecr

import (
	"context"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
)

type ECRService interface {
	DischargeOverview(ctx *gin.Context, req DischargeOverviewRequest) (*pagination.Response[DischargeOverviewResponse], error)
	TotalDischargeCount(ctx *gin.Context) (*TotalDischargeCountResponse, error)
	ListEmployeesByContractEndDate(ctx context.Context) ([]ListEmployeesByContractEndDateResponse, error)
	ListLatestPayments(ctx context.Context) ([]ListLatestPaymentsResponse, error)
	ListUpcomingAppointments(ctx context.Context, employeeID int64) ([]ListUpcomingAppointmentsResponse, error)
}

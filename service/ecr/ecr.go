package ecr

import (
	"context"

	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ECRService interface {
	DischargeOverview(ctx *gin.Context, req DischargeOverviewRequest) (*pagination.Response[DischargeOverviewResponse], error)
	TotalDischargeCount(ctx *gin.Context) (*TotalDischargeCountResponse, error)
	ListEmployeesByContractEndDate(ctx context.Context) ([]ListEmployeesByContractEndDateResponse, error)
	ListLatestPayments(ctx context.Context) ([]ListLatestPaymentsResponse, error)
	ListUpcomingAppointments(ctx context.Context, employeeID uuid.UUID) ([]ListUpcomingAppointmentsResponse, error)
}

type ecrService struct {
	*deps.ServiceDependencies
}

func NewECRService(deps *deps.ServiceDependencies) ECRService {
	return &ecrService{
		ServiceDependencies: deps,
	}
}

package late_arrival

import (
	"context"

	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LateArrivalService interface {
	CreateLateArrival(ctx context.Context, employeeID uuid.UUID, req *CreateLateArrivalRequest) (*CreateLateArrivalResponse, error)
	CreateLateArrivalByAdmin(ctx context.Context, adminEmployeeID uuid.UUID, req *CreateLateArrivalByAdminRequest) (*CreateLateArrivalResponse, error)
	ListMyLateArrivals(ctx *gin.Context, employeeID uuid.UUID, req *ListMyLateArrivalsRequest) (*pagination.Response[LateArrivalListItem], error)
	ListLateArrivals(ctx *gin.Context, req *ListLateArrivalsRequest) (*pagination.Response[LateArrivalListItem], error)
}

type lateArrivalService struct {
	*deps.ServiceDependencies
}

func NewLateArrivalService(deps *deps.ServiceDependencies) LateArrivalService {
	return &lateArrivalService{
		ServiceDependencies: deps,
	}
}

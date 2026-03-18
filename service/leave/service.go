package leave

import (
	"context"

	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LeaveService interface {
	CreateLeaveRequest(ctx context.Context, employeeID uuid.UUID, req *CreateLeaveRequestRequest) (*CreateLeaveRequestResponse, error)
	UpdateLeaveRequest(ctx context.Context, employeeID, leaveRequestID uuid.UUID, req *UpdateLeaveRequestRequest) (*UpdateLeaveRequestResponse, error)
	UpdateLeaveRequestByAdmin(ctx context.Context, adminEmployeeID, leaveRequestID uuid.UUID, req *UpdateLeaveRequestAdminRequest) (*UpdateLeaveRequestResponse, error)
	ListMyLeaveRequests(ctx *gin.Context, employeeID uuid.UUID, req *ListMyLeaveRequestsRequest) (*pagination.Response[LeaveRequestListItem], error)
	ListLeaveRequests(ctx *gin.Context, req *ListLeaveRequestsRequest) (*pagination.Response[LeaveRequestListItem], error)
}

type leaveService struct {
	*deps.ServiceDependencies
}

func NewLeaveService(deps *deps.ServiceDependencies) LeaveService {
	return &leaveService{
		ServiceDependencies: deps,
	}
}

package schedule

import (
	"context"

	"maicare_go/async/aclient"
	"maicare_go/service/deps"

	"github.com/google/uuid"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, employeeID uuid.UUID, req *CreateScheduleRequest) (*CreateScheduleResponse, error)
	GetMonthlySchedulesByLocation(ctx context.Context, locationID int64, req *GetMonthlySchedulesByLocationRequest) ([]GetMonthlySchedulesByLocationResponse, error)
	GetDailySchedulesByLocation(ctx context.Context, locationID int64, req *GetDailySchedulesByLocationRequest) (*GetDailySchedulesByLocationResponse, error)
	GetScheduleByID(ctx context.Context, scheduleID uuid.UUID) (*GetScheduleByIdResponse, error)
	UpdateSchedule(ctx context.Context, scheduleID uuid.UUID, updaterEmployeeID uuid.UUID, req *UpdateScheduleRequest) (*UpdateScheduleResponse, error)
	DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error
}

type scheduleService struct {
	*deps.ServiceDependencies
	asynqClient aclient.AsynqClientInterface
}

func NewScheduleService(deps *deps.ServiceDependencies, asynqClient aclient.AsynqClientInterface) ScheduleService {
	return &scheduleService{
		ServiceDependencies: deps,
		asynqClient:         asynqClient,
	}
}

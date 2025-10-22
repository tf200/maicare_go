package schedule

import (
	"context"
	"maicare_go/async/aclient"
	"maicare_go/service/deps"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, employeeID int64, req *CreateScheduleRequest) (*CreateScheduleResponse, error)
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

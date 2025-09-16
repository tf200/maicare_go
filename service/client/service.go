package client

import (
	"context"
	db "maicare_go/db/sqlc"
	"maicare_go/service/deps"
)

type ClientService interface {
	CreateAppointmentCard(req CreateAppointmentCardRequest, ctx context.Context) (*db.AppointmentCard, error)
	GetAppointmentCard(ctx context.Context, clientID int64) (*db.GetAppointmentCardRow, error)
}

type clientService struct {
	*deps.ServiceDependencies
}

func NewClientService(deps *deps.ServiceDependencies) ClientService {
	return &clientService{
		ServiceDependencies: deps,
	}
}

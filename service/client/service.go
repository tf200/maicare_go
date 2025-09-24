package clientp

import (
	"context"
	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
)

type ClientService interface {
	// Client Details
	CreateClientDetails(req CreateClientDetailsRequest, ctx context.Context) (*CreateClientDetailsResponse, error)
	ListClientDetails(ctx *gin.Context, req ListClientsApiParams) (*pagination.Response[ListClientsApiResponse], error)
	GetClientsCount(ctx context.Context) (*GetClientsCountResponse, error)
	GetClientDetails(ctx context.Context, clientID int64) (*GetClientApiResponse, error)
	GetClientAddresses(ctx context.Context, clientID int64) (*GetClientAddressesApiResponse, error)
	UpdateClientDetails(ctx context.Context, req UpdateClientDetailsRequest, clientID int64) (*UpdateClientDetailsResponse, error)

	// Client Appointment Card
	CreateAppointmentCard(req CreateAppointmentCardRequest, clientID int64, ctx context.Context) (*CreateAppointmentCardResponse, error)
	GetAppointmentCard(ctx context.Context, clientID int64) (*GetAppointmentCardResponse, error)
}

type clientService struct {
	*deps.ServiceDependencies
}

func NewClientService(deps *deps.ServiceDependencies) ClientService {
	return &clientService{
		ServiceDependencies: deps,
	}
}

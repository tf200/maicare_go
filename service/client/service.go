package clientp

import (
	"context"
	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientService interface {
	// Client Details
	CreateClientDetails(req CreateClientDetailsRequest, ctx context.Context) (*CreateClientDetailsResponse, error)
	ListClientDetails(ctx *gin.Context, req ListClientsApiParams) (*pagination.Response[ListClientsApiResponse], error)
	GetClientsCount(ctx context.Context) (*GetClientsCountResponse, error)
	GetClientDetails(ctx context.Context, clientID int64) (*GetClientApiResponse, error)
	GetClientAddresses(ctx context.Context, clientID int64) (*GetClientAddressesApiResponse, error)
	UpdateClientDetails(ctx context.Context, req UpdateClientDetailsRequest, clientID int64) (*UpdateClientDetailsResponse, error)
	UpdateClientStatus(ctx context.Context, req UpdateClientStatusRequest, clientID int64) (*UpdateClientStatusResponse, error)
	ListStatusHistory(ctx context.Context, clientID int64) ([]ListStatusHistoryApiResponse, error)
	SetClientProfilePicture(ctx context.Context, req SetClientProfilePictureRequest, clientID int64) (*SetClientProfilePictureResponse, error)
	// Client Documents
	AddClientDocument(ctx context.Context, req AddClientDocumentApiRequest, clientID int64) (*AddClientDocumentApiResponse, error)
	ListClientDocuments(ctx *gin.Context, req ListClientDocumentsApiRequest, clientID int64) (*pagination.Response[ListClientDocumentsApiResponse], error)
	DeleteClientDocument(ctx context.Context, clientID int64, attachmentID uuid.UUID) (*DeleteClientDocumentApiResponse, error)
	GetMissingClientDocuments(ctx context.Context, clientID int64) (*GetMissingClientDocumentsApiResponse, error)

	// Client Appointment Card
	CreateAppointmentCard(req CreateAppointmentCardRequest, clientID int64, ctx context.Context) (*CreateAppointmentCardResponse, error)
	GetAppointmentCard(ctx context.Context, clientID int64) (*GetAppointmentCardResponse, error)
	UpdateAppointmentCard(req UpdateAppointmentCardRequest, clientID int64, ctx context.Context) (*UpdateAppointmentCardResponse, error)
	GenerateAppointmentCardDocumentApi(ctx context.Context, clientID int64) (*GenerateAppointmentCardDocumentApiResponse, error)

	// Client Incidents
	CreateIncident(ctx context.Context, req CreateIncidentRequest, clientID int64) (*CreateIncidentResponse, error)
	ListIncidents(ctx *gin.Context, req ListIncidentsRequest, clientID int64) (*pagination.Response[ListIncidentsResponse], error)
	GetIncident(ctx context.Context, incidentID int64) (*GetIncidentResponse, error)
	UpdateIncident(ctx context.Context, req UpdateIncidentRequest, incidentID int64) (*UpdateIncidentResponse, error)
	DeleteIncident(ctx context.Context, incidentID int64) error
	GenerateIncidentFile(ctx context.Context, incidentID int64) (*GenerateIncidentFileResponse, error)
	ConfirmIncident(ctx context.Context, incidentID int64) (*ConfirmIncidentResponse, error)
}

type clientService struct {
	*deps.ServiceDependencies
}

func NewClientService(deps *deps.ServiceDependencies) ClientService {
	return &clientService{
		ServiceDependencies: deps,
	}
}

package appointment

import (
	"context"
	"maicare_go/async/aclient"
	"maicare_go/service/deps"

	"github.com/google/uuid"
)

type AppointmentService interface {
	CreateAppointment(req *CreateAppointmentRequest, userID uuid.UUID, ctx context.Context) (*CreateAppointmentResponse, error)
	AddParticipantToAppointment(ctx context.Context, appointmentID uuid.UUID, req AddParticipantToAppointmentRequest) error
	AddClientToAppointment(ctx context.Context, appointmentID uuid.UUID, req AddClientToAppointmentRequest) error
	ListAppointmentsForEmployeeInRange(ctx context.Context, employeeID uuid.UUID, req ListAppointmentsForEmployeeInRangeRequest) ([]ListAppointmentsForEmployeeInRangeResponse, error)
	ListAppointmentsForClientInRange(ctx context.Context, clientID uuid.UUID, req ListAppointmentsForClientRequest) ([]ListAppointmentsForClientResponse, error)
	GetAppointment(ctx context.Context, appointmentID uuid.UUID) (*GetAppointmentResponse, error)
	UpdateAppointment(ctx context.Context, appointmentID uuid.UUID, req *UpdateAppointmentRequest) (*UpdateAppointmentResponse, error)
	DeleteAppointment(ctx context.Context, appointmentID uuid.UUID) error
	ConfirmAppointment(ctx context.Context, appointmentID uuid.UUID, employeeID uuid.UUID) error
}

type appointmentService struct {
	*deps.ServiceDependencies
	asynqClient aclient.AsynqClientInterface
}

func NewAppointmentService(deps *deps.ServiceDependencies, asynqClient aclient.AsynqClientInterface) AppointmentService {
	return &appointmentService{
		ServiceDependencies: deps,
		asynqClient:         asynqClient,
	}
}

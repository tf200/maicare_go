package appointment

import (
	"context"
	"maicare_go/service/deps"

	"github.com/google/uuid"
)

type AppointmentService interface {
	CreateAppointment(req *CreateAppointmentRequest, userID int64, ctx context.Context) (*CreateAppointmentResponse, error)
	AddParticipantToAppointment(ctx context.Context, appointmentID uuid.UUID, req AddParticipantToAppointmentRequest) error
	AddClientToAppointment(ctx context.Context, appointmentID uuid.UUID, req AddClientToAppointmentRequest) error
	ListAppointmentsForEmployeeInRange(ctx context.Context, employeeID int64, req ListAppointmentsForEmployeeInRangeRequest) ([]ListAppointmentsForEmployeeInRangeResponse, error)
	ListAppointmentsForClientInRange(ctx context.Context, clientID int64, req ListAppointmentsForClientRequest) ([]ListAppointmentsForClientResponse, error)
	GetAppointment(ctx context.Context, appointmentID uuid.UUID) (*GetAppointmentResponse, error)
	UpdateAppointment(ctx context.Context, appointmentID uuid.UUID, req *UpdateAppointmentRequest) (*UpdateAppointmentResponse, error)
	DeleteAppointment(ctx context.Context, appointmentID uuid.UUID) error
	ConfirmAppointment(ctx context.Context, appointmentID uuid.UUID, employeeID int64) error
}

type appointmentService struct {
	*deps.ServiceDependencies
}

func NewAppointmentService(deps *deps.ServiceDependencies) AppointmentService {
	return &appointmentService{
		ServiceDependencies: deps,
	}
}

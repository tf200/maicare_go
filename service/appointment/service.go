package appointment

import (
	"context"

	"maicare_go/async/aclient"
	"maicare_go/service/deps"

	"github.com/google/uuid"
)

type AppointmentService interface {
	CreateEvent(ctx context.Context, req *CreateEventRequest, employeeID uuid.UUID) (*EventResponse, error)
	ListEvents(ctx context.Context, req ListEventsRequest, employeeID uuid.UUID) ([]EventOccurrenceResponse, error)
	GetEvent(ctx context.Context, eventID uuid.UUID, employeeID uuid.UUID) (*EventResponse, error)
	UpdateEvent(ctx context.Context, eventID uuid.UUID, req *UpdateEventRequest, employeeID uuid.UUID) (*EventResponse, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID, req DeleteEventRequest, employeeID uuid.UUID) error
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

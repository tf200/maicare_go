package appointment

import "maicare_go/service/deps"

type AppointmentService interface{}

type appointmentService struct {
	*deps.ServiceDependencies
}

func NewAppointmentService(deps *deps.ServiceDependencies) AppointmentService {
	return &appointmentService{
		ServiceDependencies: deps,
	}
}

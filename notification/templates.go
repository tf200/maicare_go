package notification

const (
	// Notification Types
	TypeNewAppointment    = "new_appointment"
	TypeAppointmentUpdate = "appointment_update"
)

// Notifications Data Templates

type NewAppointmentData struct {
	AppointmentID int64  `json:"appointment_id"`
	CreatedBy     string `json:"created_by"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	Location      string `json:"location"`
}

package clientp

import (
	"github.com/google/uuid"
)

// PromoteIntakeToClientRequest represents a request to promote an intake form to a full client record
type PromoteIntakeToClientRequest struct {
	IntakeFormID uuid.UUID `json:"intake_form_id" binding:"required,uuid"`
}

// PromoteIntakeToClientResponse represents the response after successfully promoting an intake to client
type PromoteIntakeToClientResponse struct {
	ClientID                   uuid.UUID `json:"client_id"`
	IntakeFormID               uuid.UUID `json:"intake_form_id"`
	RegistrationFormID         uuid.UUID `json:"registration_form_id"`
	Message                    string    `json:"message"`
	MaturityAssessmentsCreated int       `json:"maturity_assessments_created"`
	EmergencyContactsCreated   int       `json:"emergency_contacts_created"`
}

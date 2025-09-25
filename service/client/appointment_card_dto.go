package clientp

import "time"

// CreateAppointmentCardRequest represents a request to create a new appointment card
type CreateAppointmentCardRequest struct {
	GeneralInformation     []string `json:"general_information"`
	ImportantContacts      []string `json:"important_contacts"`
	HouseholdInfo          []string `json:"household_info"`
	OrganizationAgreements []string `json:"organization_agreements"`
	YouthOfficerAgreements []string `json:"youth_officer_agreements"`
	TreatmentAgreements    []string `json:"treatment_agreements"`
	SmokingRules           []string `json:"smoking_rules"`
	Work                   []string `json:"work"`
	SchoolInternship       []string `json:"school_internship"`
	Travel                 []string `json:"travel"`
	Leave                  []string `json:"leave"`
}

// CreateAppointmentCardResponse represents a response to a create appointment card request
type CreateAppointmentCardResponse struct {
	ID                     int64     `json:"id"`
	ClientID               int64     `json:"client_id"`
	GeneralInformation     []string  `json:"general_information"`
	ImportantContacts      []string  `json:"important_contacts"`
	HouseholdInfo          []string  `json:"household_info"`
	OrganizationAgreements []string  `json:"organization_agreements"`
	YouthOfficerAgreements []string  `json:"youth_officer_agreements"`
	TreatmentAgreements    []string  `json:"treatment_agreements"`
	SmokingRules           []string  `json:"smoking_rules"`
	Work                   []string  `json:"work"`
	SchoolInternship       []string  `json:"school_internship"`
	Travel                 []string  `json:"travel"`
	Leave                  []string  `json:"leave"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	FileUrl                *string   `json:"file_url"`
}

// GetAppointmentCardResponse represents a response to a get appointment card request
type GetAppointmentCardResponse struct {
	ID                     int64     `json:"id"`
	ClientID               int64     `json:"client_id"`
	GeneralInformation     []string  `json:"general_information"`
	ImportantContacts      []string  `json:"important_contacts"`
	HouseholdInfo          []string  `json:"household_info"`
	OrganizationAgreements []string  `json:"organization_agreements"`
	YouthOfficerAgreements []string  `json:"youth_officer_agreements"`
	TreatmentAgreements    []string  `json:"treatment_agreements"`
	SmokingRules           []string  `json:"smoking_rules"`
	Work                   []string  `json:"work"`
	SchoolInternship       []string  `json:"school_internship"`
	Travel                 []string  `json:"travel"`
	Leave                  []string  `json:"leave"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	FileUrl                *string   `json:"file_url"`
}

// UpdateAppointmentCardRequest represents a request to update an appointment card
type UpdateAppointmentCardRequest struct {
	GeneralInformation     []string `json:"general_information"`
	ImportantContacts      []string `json:"important_contacts"`
	HouseholdInfo          []string `json:"household_info"`
	OrganizationAgreements []string `json:"organization_agreements"`
	YouthOfficerAgreements []string `json:"youth_officer_agreements"`
	TreatmentAgreements    []string `json:"treatment_agreements"`
	SmokingRules           []string `json:"smoking_rules"`
	Work                   []string `json:"work"`
	SchoolInternship       []string `json:"school_internship"`
	Travel                 []string `json:"travel"`
	Leave                  []string `json:"leave"`
}

// UpdateAppointmentCardResponse represents a response to an update appointment card request
type UpdateAppointmentCardResponse struct {
	ID                     int64     `json:"id"`
	ClientID               int64     `json:"client_id"`
	GeneralInformation     []string  `json:"general_information"`
	ImportantContacts      []string  `json:"important_contacts"`
	HouseholdInfo          []string  `json:"household_info"`
	OrganizationAgreements []string  `json:"organization_agreements"`
	YouthOfficerAgreements []string  `json:"youth_officer_agreements"`
	TreatmentAgreements    []string  `json:"treatment_agreements"`
	SmokingRules           []string  `json:"smoking_rules"`
	Work                   []string  `json:"work"`
	SchoolInternship       []string  `json:"school_internship"`
	Travel                 []string  `json:"travel"`
	Leave                  []string  `json:"leave"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// GenerateAppointmentCardDocumentApiResponse represents a response to generate appointment card document request
type GenerateAppointmentCardDocumentApiResponse struct {
	ClientID int64   `json:"client_id"`
	FileUrl  *string `json:"file_url"`
}

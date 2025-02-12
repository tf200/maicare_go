package api

import (
	db "maicare_go/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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
}

// CreateAppointmentCardApi creates a new appointment card
// @Summary Create a new appointment card
// @Description Create a new appointment card
// @Tags appointment_cards
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body CreateAppointmentCardRequest true "Request body"
// @Success 201 {object} Response[CreateAppointmentCardResponse]
// @Router /clients/{id}/appointment_cards [post]
func (server *Server) CreateAppointmentCardApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateAppointmentCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAppointmentCardParams{
		ClientID:               clientID,
		GeneralInformation:     req.GeneralInformation,
		ImportantContacts:      req.ImportantContacts,
		HouseholdInfo:          req.HouseholdInfo,
		OrganizationAgreements: req.OrganizationAgreements,
		YouthOfficerAgreements: req.YouthOfficerAgreements,
		TreatmentAgreements:    req.TreatmentAgreements,
		SmokingRules:           req.SmokingRules,
		Work:                   req.Work,
		SchoolInternship:       req.SchoolInternship,
		Travel:                 req.Travel,
		Leave:                  req.Leave,
	}

	appointmentCard, err := server.store.CreateAppointmentCard(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateAppointmentCardResponse{
		ID:                     appointmentCard.ID,
		ClientID:               appointmentCard.ClientID,
		GeneralInformation:     appointmentCard.GeneralInformation,
		ImportantContacts:      appointmentCard.ImportantContacts,
		HouseholdInfo:          appointmentCard.HouseholdInfo,
		OrganizationAgreements: appointmentCard.OrganizationAgreements,
		YouthOfficerAgreements: appointmentCard.YouthOfficerAgreements,
		TreatmentAgreements:    appointmentCard.TreatmentAgreements,
		SmokingRules:           appointmentCard.SmokingRules,
		Work:                   appointmentCard.Work,
		SchoolInternship:       appointmentCard.SchoolInternship,
		Travel:                 appointmentCard.Travel,
		Leave:                  appointmentCard.Leave,
		CreatedAt:              appointmentCard.CreatedAt.Time,
		UpdatedAt:              appointmentCard.UpdatedAt.Time,
	}, "Appointment card created successfully")

	ctx.JSON(http.StatusCreated, res)

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
}

// GetAppointmentCardApi retrieves an appointment card by client ID
// @Summary Get an appointment card by client ID
// @Description Get an appointment card by client ID
// @Tags appointment_cards
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[GetAppointmentCardResponse]
// @Router /clients/{id}/appointment_cards [get]
func (server *Server) GetAppointmentCardApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointmentCard, err := server.store.GetAppointmentCard(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetAppointmentCardResponse{
		ID:                     appointmentCard.ID,
		ClientID:               appointmentCard.ClientID,
		GeneralInformation:     appointmentCard.GeneralInformation,
		ImportantContacts:      appointmentCard.ImportantContacts,
		HouseholdInfo:          appointmentCard.HouseholdInfo,
		OrganizationAgreements: appointmentCard.OrganizationAgreements,
		YouthOfficerAgreements: appointmentCard.YouthOfficerAgreements,
		TreatmentAgreements:    appointmentCard.TreatmentAgreements,
		SmokingRules:           appointmentCard.SmokingRules,
		Work:                   appointmentCard.Work,
		SchoolInternship:       appointmentCard.SchoolInternship,
		Travel:                 appointmentCard.Travel,
		Leave:                  appointmentCard.Leave,
		CreatedAt:              appointmentCard.CreatedAt.Time,
		UpdatedAt:              appointmentCard.UpdatedAt.Time,
	}, "Appointment card retrieved successfully")

	ctx.JSON(http.StatusOK, res)
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

// UpdateAppointmentCardApi updates an appointment card by client ID
// @Summary Update an appointment card by client ID
// @Description Update an appointment card by client ID
// @Tags appointment_cards
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body UpdateAppointmentCardRequest true "Request body"
// @Success 200 {object} Response[UpdateAppointmentCardResponse]
// @Router /clients/{id}/appointment_cards [put]
func (server *Server) UpdateAppointmentCardApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateAppointmentCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAppointmentCardParams{
		ClientID:               clientID,
		GeneralInformation:     req.GeneralInformation,
		ImportantContacts:      req.ImportantContacts,
		HouseholdInfo:          req.HouseholdInfo,
		OrganizationAgreements: req.OrganizationAgreements,
		YouthOfficerAgreements: req.YouthOfficerAgreements,
		TreatmentAgreements:    req.TreatmentAgreements,
		SmokingRules:           req.SmokingRules,
		Work:                   req.Work,
		SchoolInternship:       req.SchoolInternship,
		Travel:                 req.Travel,
		Leave:                  req.Leave,
	}

	appointmentCard, err := server.store.UpdateAppointmentCard(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateAppointmentCardResponse{
		ID:                     appointmentCard.ID,
		ClientID:               appointmentCard.ClientID,
		GeneralInformation:     appointmentCard.GeneralInformation,
		ImportantContacts:      appointmentCard.ImportantContacts,
		HouseholdInfo:          appointmentCard.HouseholdInfo,
		OrganizationAgreements: appointmentCard.OrganizationAgreements,
		YouthOfficerAgreements: appointmentCard.YouthOfficerAgreements,
		TreatmentAgreements:    appointmentCard.TreatmentAgreements,
		SmokingRules:           appointmentCard.SmokingRules,
		Work:                   appointmentCard.Work,
		SchoolInternship:       appointmentCard.SchoolInternship,
		Travel:                 appointmentCard.Travel,
		Leave:                  appointmentCard.Leave,
		CreatedAt:              appointmentCard.CreatedAt.Time,
		UpdatedAt:              appointmentCard.UpdatedAt.Time,
	}, "Appointment card updated successfully")

	ctx.JSON(http.StatusOK, res)
}

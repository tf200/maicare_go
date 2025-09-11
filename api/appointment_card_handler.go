package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pdf"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	FileUrl                *string   `json:"file_url"`
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
		server.logBusinessEvent(LogLevelError, "CreateAppointmentCardApi", "Invalid client ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	var req CreateAppointmentCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateAppointmentCardApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
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
		server.logBusinessEvent(LogLevelError, "CreateAppointmentCardApi", "Failed to create appointment card", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create appointment card")))
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
		FileUrl:                appointmentCard.FileUrl,
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
	FileUrl                *string   `json:"file_url"`
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
		server.logBusinessEvent(LogLevelError, "GetAppointmentCardApi", "Invalid client ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	appointmentCard, err := server.store.GetAppointmentCard(ctx, clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusOK, SuccessResponse[any](nil, "No appointment card found for the given client ID"))
			return
		}
		server.logBusinessEvent(LogLevelError, "GetAppointmentCardApi", "Failed to retrieve appointment card", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to retrieve appointment card")))
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
		FileUrl:                appointmentCard.FileUrl,
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
		server.logBusinessEvent(LogLevelError, "UpdateAppointmentCardApi", "Invalid client ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	var req UpdateAppointmentCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateAppointmentCardApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
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
		server.logBusinessEvent(LogLevelError, "UpdateAppointmentCardApi", "Failed to update appointment card", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update appointment card")))
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

type GenerateAppointmentCardDocumentApiResponse struct {
	ClientID int64   `json:"client_id"`
	FileUrl  *string `json:"file_url"`
}

// GenerateAppointmentCardDocument generates an appointment card document by client ID
// @Summary Generate an appointment card document by client ID
// @Description Generate an appointment card document by client ID
// @Tags appointment_cards
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[UpdateAppointmentCardResponse]
// @Router /clients/{id}/appointment_cards/generate_document [post]
func (server *Server) GenerateAppointmentCardDocumentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateAppointmentCardDocumentApi", "Invalid client ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	appointmentCard, err := server.store.GetAppointmentCard(ctx, clientID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateAppointmentCardDocumentApi", "Failed to retrieve appointment card", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to retrieve appointment card")))
		return
	}

	if appointmentCard.FileUrl != nil && *appointmentCard.FileUrl != "" {
		err = server.b2Client.Delete(ctx, *appointmentCard.FileUrl)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "GenerateAppointmentCardDocumentApi", "Failed to delete existing appointment card document", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to delete existing appointment card document"))
			return
		}
	}

	pdfArg := pdf.AppointmentCard{
		ID:                     appointmentCard.ID,
		ClientName:             appointmentCard.FirstName + " " + appointmentCard.LastName,
		Date:                   appointmentCard.CreatedAt.Time.Format("02-01-2006"),
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
	}

	pdfUrl, err := pdf.GenerateAndUploadAppointmentCardPDF(ctx, pdfArg, server.b2Client)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateAppointmentCardDocumentApi", "Failed to generate and upload appointment card PDF", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to generate and upload appointment card PDF")))
		return
	}

	fileKey, err := server.store.UpdateAppointmentCardUrl(ctx, db.UpdateAppointmentCardUrlParams{
		ClientID: clientID,
		FileUrl:  &pdfUrl,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateAppointmentCardDocumentApi", "Failed to update appointment card with file URL", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update appointment card with file URL")))
		return
	}

	response := GenerateAppointmentCardDocumentApiResponse{
		FileUrl:  server.generateResponsePresignedURL(fileKey),
		ClientID: clientID,
	}

	res := SuccessResponse(response, "Appointment card document generated successfully")
	ctx.JSON(http.StatusOK, res)

}

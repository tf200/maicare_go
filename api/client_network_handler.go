package api

import (
	"log"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateClientEmergencyContactParams defines the request for creating a client emergency contact
type CreateClientEmergencyContactParams struct {
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	Email            *string `json:"email"`
	PhoneNumber      *string `json:"phone_number"`
	Address          *string `json:"address"`
	Relationship     *string `json:"relationship"`
	RelationStatus   *string `json:"relation_status"`
	MedicalReports   bool    `json:"medical_reports"`
	IncidentsReports bool    `json:"incidents_reports"`
	GoalsReports     bool    `json:"goals_reports"`
}

// CreateClientEmergencyContactResponse defines the response for creating a client emergency contact
type CreateClientEmergencyContactResponse struct {
	ID               int64     `json:"id"`
	ClientID         int64     `json:"client_id"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            *string   `json:"email"`
	PhoneNumber      *string   `json:"phone_number"`
	Address          *string   `json:"address"`
	Relationship     *string   `json:"relationship"`
	RelationStatus   *string   `json:"relation_status"`
	CreatedAt        time.Time `json:"created_at"`
	IsVerified       bool      `json:"is_verified"`
	MedicalReports   bool      `json:"medical_reports"`
	IncidentsReports bool      `json:"incidents_reports"`
	GoalsReports     bool      `json:"goals_reports"`
}

// CreateClientEmergencyContactApi creates a client emergency contact
// @Summary Create a client emergency contact
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body CreateClientEmergencyContactParams true "Client emergency contact data"
// @Success 201 {object} Response[CreateClientEmergencyContactResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/emergency_contacts [post]
func (server *Server) CreateClientEmergencyContactApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateClientEmergencyContactParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateEmemrgencyContactParams{
		ClientID:         clientID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		Address:          req.Address,
		Relationship:     req.Relationship,
		RelationStatus:   req.RelationStatus,
		MedicalReports:   req.MedicalReports,
		IncidentsReports: req.IncidentsReports,
		GoalsReports:     req.GoalsReports,
	}
	clientEmergencyContact, err := server.store.CreateEmemrgencyContact(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateClientEmergencyContactResponse{
		ID:               clientEmergencyContact.ID,
		ClientID:         clientEmergencyContact.ClientID,
		FirstName:        clientEmergencyContact.FirstName,
		LastName:         clientEmergencyContact.LastName,
		Email:            clientEmergencyContact.Email,
		PhoneNumber:      clientEmergencyContact.PhoneNumber,
		Address:          clientEmergencyContact.Address,
		Relationship:     clientEmergencyContact.Relationship,
		RelationStatus:   clientEmergencyContact.RelationStatus,
		CreatedAt:        clientEmergencyContact.CreatedAt.Time,
		IsVerified:       clientEmergencyContact.IsVerified,
		MedicalReports:   clientEmergencyContact.MedicalReports,
		IncidentsReports: clientEmergencyContact.IncidentsReports,
		GoalsReports:     clientEmergencyContact.GoalsReports,
	}, "Client emergency contact created successfully")
	ctx.JSON(http.StatusCreated, res)

}

// ListClientEmergencyContactsRequest defines the request for listing client emergency contacts
type ListClientEmergencyContactsRequest struct {
	pagination.Request
	Search string `form:"search"`
}

// ListClientEmergencyContactsResponse defines the response for listing client emergency contacts
type ListClientEmergencyContactsResponse struct {
	ID               int64     `json:"id"`
	ClientID         int64     `json:"client_id"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            *string   `json:"email"`
	PhoneNumber      *string   `json:"phone_number"`
	Address          *string   `json:"address"`
	Relationship     *string   `json:"relationship"`
	RelationStatus   *string   `json:"relation_status"`
	CreatedAt        time.Time `json:"created_at"`
	IsVerified       bool      `json:"is_verified"`
	MedicalReports   bool      `json:"medical_reports"`
	IncidentsReports bool      `json:"incidents_reports"`
	GoalsReports     bool      `json:"goals_reports"`
}

// ListClientEmergencyContactsApi lists all client emergency contacts
// @Summary List all client emergency contacts
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search query"
// @Success 200 {object} Response[ListClientEmergencyContactsResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/emergency_contacts [get]
func (server *Server) ListClientEmergencyContactsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	log.Printf("Processing request for client ID: %s", id)

	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("Failed to parse client ID: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListClientEmergencyContactsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Printf("Failed to bind query params: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()
	log.Printf("Query params: limit=%d, offset=%d", params.Limit, params.Offset)

	contacts, err := server.store.ListEmergencyContacts(ctx, db.ListEmergencyContactsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
		Search:   req.Search,
	})
	if err != nil {
		log.Printf("Failed to fetch contacts: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Handle case where no contacts are found
	if len(contacts) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListClientEmergencyContactsResponse{}, 0)
		res := SuccessResponse(pag, "No emergency contacts found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	contactsRes := make([]ListClientEmergencyContactsResponse, len(contacts))
	for i, contact := range contacts {
		contactsRes[i] = ListClientEmergencyContactsResponse{
			ID:               contact.ID,
			ClientID:         contact.ClientID,
			FirstName:        contact.FirstName,
			LastName:         contact.LastName,
			Email:            contact.Email,
			PhoneNumber:      contact.PhoneNumber,
			Address:          contact.Address,
			Relationship:     contact.Relationship,
			RelationStatus:   contact.RelationStatus,
			CreatedAt:        contact.CreatedAt.Time,
			IsVerified:       contact.IsVerified,
			MedicalReports:   contact.MedicalReports,
			IncidentsReports: contact.IncidentsReports,
			GoalsReports:     contact.GoalsReports,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, contactsRes, contacts[0].TotalCount)
	res := SuccessResponse(pag, "Client emergency contacts fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientEmergencyContactResponse defines the response for getting a client emergency contact
type GetClientEmergencyContactResponse struct {
	ID               int64     `json:"id"`
	ClientID         int64     `json:"client_id"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            *string   `json:"email"`
	PhoneNumber      *string   `json:"phone_number"`
	Address          *string   `json:"address"`
	Relationship     *string   `json:"relationship"`
	RelationStatus   *string   `json:"relation_status"`
	CreatedAt        time.Time `json:"created_at"`
	IsVerified       bool      `json:"is_verified"`
	MedicalReports   bool      `json:"medical_reports"`
	IncidentsReports bool      `json:"incidents_reports"`
	GoalsReports     bool      `json:"goals_reports"`
}

// GetClientEmergencyContactApi gets a client emergency contact
// @Summary Get a client emergency contact
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param contact_id path int true "Contact ID"
// @Success 200 {object} Response[GetClientEmergencyContactResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/emergency_contacts/{contact_id} [get]
func (server *Server) GetClientEmergencyContactApi(ctx *gin.Context) {
	id := ctx.Param("contact_id")
	contactID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contact, err := server.store.GetEmergencyContact(ctx, contactID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetClientEmergencyContactResponse{
		ID:               contact.ID,
		ClientID:         contact.ClientID,
		FirstName:        contact.FirstName,
		LastName:         contact.LastName,
		Email:            contact.Email,
		PhoneNumber:      contact.PhoneNumber,
		Address:          contact.Address,
		Relationship:     contact.Relationship,
		RelationStatus:   contact.RelationStatus,
		CreatedAt:        contact.CreatedAt.Time,
		IsVerified:       contact.IsVerified,
		MedicalReports:   contact.MedicalReports,
		IncidentsReports: contact.IncidentsReports,
		GoalsReports:     contact.GoalsReports,
	}, "Client emergency contact fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateClientEmergencyContactParams defines the request for updating a client emergency contact
type UpdateClientEmergencyContactParams struct {
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	Email            *string `json:"email"`
	PhoneNumber      *string `json:"phone_number"`
	Address          *string `json:"address"`
	Relationship     *string `json:"relationship"`
	RelationStatus   *string `json:"relation_status"`
	MedicalReports   *bool   `json:"medical_reports"`
	IncidentsReports *bool   `json:"incidents_reports"`
	GoalsReports     *bool   `json:"goals_reports"`
}

// UpdateClientEmergencyContactResponse defines the response for updating a client emergency contact
type UpdateClientEmergencyContactResponse struct {
	ID               int64     `json:"id"`
	ClientID         int64     `json:"client_id"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            *string   `json:"email"`
	PhoneNumber      *string   `json:"phone_number"`
	Address          *string   `json:"address"`
	Relationship     *string   `json:"relationship"`
	RelationStatus   *string   `json:"relation_status"`
	CreatedAt        time.Time `json:"created_at"`
	IsVerified       bool      `json:"is_verified"`
	MedicalReports   bool      `json:"medical_reports"`
	IncidentsReports bool      `json:"incidents_reports"`
	GoalsReports     bool      `json:"goals_reports"`
}

// UpdateClientEmergencyContactApi updates a client emergency contact
// @Summary Update a client emergency contact
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param contact_id path int true "Contact ID"
// @Param request body UpdateClientEmergencyContactParams true "Client emergency contact data"
// @Success 200 {object} Response[UpdateClientEmergencyContactResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/emergency_contacts/{contact_id} [put]
func (server *Server) UpdateClientEmergencyContactApi(ctx *gin.Context) {
	id := ctx.Param("contact_id")
	contactID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateClientEmergencyContactParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.UpdateEmergencyContactParams{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		Address:          req.Address,
		Relationship:     req.Relationship,
		RelationStatus:   req.RelationStatus,
		MedicalReports:   req.MedicalReports,
		IncidentsReports: req.IncidentsReports,
		GoalsReports:     req.GoalsReports,
		ID:               contactID,
	}
	contact, err := server.store.UpdateEmergencyContact(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateClientEmergencyContactResponse{
		ID:               contact.ID,
		ClientID:         contact.ClientID,
		FirstName:        contact.FirstName,
		LastName:         contact.LastName,
		Email:            contact.Email,
		PhoneNumber:      contact.PhoneNumber,
		Address:          contact.Address,
		Relationship:     contact.Relationship,
		RelationStatus:   contact.RelationStatus,
		CreatedAt:        contact.CreatedAt.Time,
		IsVerified:       contact.IsVerified,
		MedicalReports:   contact.MedicalReports,
		IncidentsReports: contact.IncidentsReports,
		GoalsReports:     contact.GoalsReports,
	}, "Client emergency contact updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteClientEmergencyContactResponse defines the response for deleting a client emergency contact
type DeleteClientEmergencyContactResponse struct {
	ID int64 `json:"id"`
}

// DeleteClientEmergencyContactApi deletes a client emergency contact
// @Summary Delete a client emergency contact
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param contact_id path int true "Contact ID"
// @Success 200 {object} Response[DeleteClientEmergencyContactResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/emergency_contacts/{contact_id} [delete]
func (server *Server) DeleteClientEmergencyContactApi(ctx *gin.Context) {
	id := ctx.Param("contact_id")
	contactID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.DeleteEmergencyContact(ctx, contactID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(DeleteClientEmergencyContactResponse{
		ID: contactID,
	}, "Client emergency contact deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

package api

import (
	"log"
	"maicare_go/async"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetClientSenderResponse defines the request for getting a client sender
type GetClientSenderResponse struct {
	ID           int64     `json:"id"`
	Types        string    `json:"types"`
	Name         string    `json:"name"`
	Address      *string   `json:"address"`
	PostalCode   *string   `json:"postal_code"`
	Place        *string   `json:"place"`
	Land         *string   `json:"land"`
	Kvknumber    *string   `json:"kvknumber"`
	Btwnumber    *string   `json:"btwnumber"`
	PhoneNumber  *string   `json:"phone_number"`
	ClientNumber *string   `json:"client_number"`
	EmailAddress *string   `json:"email_address"`
	Contacts     []Contact `json:"contacts"`
	IsArchived   bool      `json:"is_archived"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GetClientSenderApi gets a client sender
// @Summary Get a client sender
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[GetClientSenderResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/sender [get]
func (server *Server) GetClientSenderApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sender, err := server.store.GetClientSender(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var contactsResp []Contact
	if err := json.Unmarshal(sender.Contacts, &contactsResp); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetClientSenderResponse{
		ID:           sender.ID,
		Types:        sender.Types,
		Name:         sender.Name,
		Address:      sender.Address,
		PostalCode:   sender.PostalCode,
		Place:        sender.Place,
		Land:         sender.Land,
		Kvknumber:    sender.Kvknumber,
		Btwnumber:    sender.Btwnumber,
		PhoneNumber:  sender.PhoneNumber,
		ClientNumber: sender.ClientNumber,
		IsArchived:   sender.IsArchived,
		Contacts:     contactsResp,
	}, "Client Sender fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListClientEmergencyContactsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	contacts, err := server.store.ListEmergencyContacts(ctx, db.ListEmergencyContactsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
		Search:   req.Search,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

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

// AssignEmployeeRequest defines the request for assigning an employee to a client
type AssignEmployeeRequest struct {
	EmployeeID int64     `json:"employee_id"`
	StartDate  time.Time `json:"start_date"`
	Role       string    `json:"role"`
}

// AssignEmployeeResponse defines the response for assigning an employee to a client
type AssignEmployeeResponse struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	EmployeeID int64     `json:"employee_id"`
	StartDate  time.Time `json:"start_date"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

// AssignEmployeeApi assigns an employee to a client
// @Summary Assign an employee to a client
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body AssignEmployeeRequest true "Employee assignment data"
// @Success 201 {object} Response[AssignEmployeeResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees [post]
func (server *Server) AssignEmployeeApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req AssignEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AssignEmployeeParams{
		ClientID:   clientID,
		EmployeeID: req.EmployeeID,
		StartDate:  pgtype.Date{Time: req.StartDate, Valid: true},
		Role:       req.Role,
	}

	assign, err := server.store.AssignEmployee(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.asynqClient.EnqueueNotificationTask(ctx, async.NotificationPayload{
		RecipientUserIDs: []int64{assign.UserID},
		Type:             "employee_assigned",
		Data:             []byte(`{"client_id":` + strconv.FormatInt(clientID, 10) + `,"employee_id":` + strconv.FormatInt(assign.EmployeeID, 10) + `}`),
	})

	res := SuccessResponse(AssignEmployeeResponse{
		ID:         assign.ID,
		ClientID:   assign.ClientID,
		EmployeeID: assign.EmployeeID,
		StartDate:  assign.StartDate.Time,
		Role:       assign.Role,
		CreatedAt:  assign.CreatedAt.Time,
	}, "Employee assigned successfully")

	ctx.JSON(http.StatusCreated, res)

}

// ListAssignedEmployeesRequest defines the request for listing assigned employees
type ListAssignedEmployeesRequest struct {
	pagination.Request
}

// ListAssignedEmployeesResponse defines the response for listing assigned employees
type ListAssignedEmployeesResponse struct {
	ID           int64     `json:"id"`
	ClientID     int64     `json:"client_id"`
	EmployeeID   int64     `json:"employee_id"`
	StartDate    time.Time `json:"start_date"`
	Role         string    `json:"role"`
	EmployeeName string    `json:"employee_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// ListAssignedEmployeesApi lists all assigned employees
// @Summary List all assigned employees
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListAssignedEmployeesResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees [get]
func (server *Server) ListAssignedEmployeesApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListAssignedEmployeesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	assigns, err := server.store.ListAssignedEmployees(ctx, db.ListAssignedEmployeesParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(assigns) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListAssignedEmployeesResponse{}, 0)
		res := SuccessResponse(pag, "No assigned employees found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	assignsRes := make([]ListAssignedEmployeesResponse, len(assigns))
	for i, assign := range assigns {
		assignsRes[i] = ListAssignedEmployeesResponse{
			ID:           assign.ID,
			ClientID:     assign.ClientID,
			EmployeeID:   assign.EmployeeID,
			StartDate:    assign.StartDate.Time,
			Role:         assign.Role,
			EmployeeName: assign.EmployeeFirstName + " " + assign.EmployeeLastName,
			CreatedAt:    assign.CreatedAt.Time,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, assignsRes, assigns[0].TotalCount)
	res := SuccessResponse(pag, "Assigned employees fetched successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetAssignedEmployeeResponse defines the response for getting an assigned employee
type GetAssignedEmployeeResponse struct {
	ID           int64     `json:"id"`
	ClientID     int64     `json:"client_id"`
	EmployeeID   int64     `json:"employee_id"`
	StartDate    time.Time `json:"start_date"`
	Role         string    `json:"role"`
	EmployeeName string    `json:"employee_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetAssignedEmployeeApi gets an assigned employee
// @Summary Get an assigned employee
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param assign_id path int true "Assignment ID"
// @Success 200 {object} Response[GetAssignedEmployeeResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees/{assign_id} [get]
func (server *Server) GetAssignedEmployeeApi(ctx *gin.Context) {
	id := ctx.Param("assign_id")
	assignID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	assign, err := server.store.GetAssignedEmployee(ctx, assignID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetAssignedEmployeeResponse{
		ID:           assign.ID,
		ClientID:     assign.ClientID,
		EmployeeID:   assign.EmployeeID,
		StartDate:    assign.StartDate.Time,
		Role:         assign.Role,
		EmployeeName: assign.EmployeeFirstName + " " + assign.EmployeeLastName,
		CreatedAt:    assign.CreatedAt.Time,
	}, "Assigned employee fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateAssignedEmployeeRequest defines the request for updating an assigned employee
type UpdateAssignedEmployeeRequest struct {
	EmployeeID *int64    `json:"employee_id"`
	StartDate  time.Time `json:"start_date"`
	Role       *string   `json:"role"`
}

// UpdateAssignedEmployeeResponse defines the response for updating an assigned employee
type UpdateAssignedEmployeeResponse struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	EmployeeID int64     `json:"employee_id"`
	StartDate  time.Time `json:"start_date"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

// UpdateAssignedEmployeeApi updates an assigned employee
// @Summary Update an assigned employee
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param assign_id path int true "Assignment ID"
// @Param request body UpdateAssignedEmployeeRequest true "Assigned employee data"
// @Success 200 {object} Response[UpdateAssignedEmployeeResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees/{assign_id} [put]
func (server *Server) UpdateAssignedEmployeeApi(ctx *gin.Context) {
	id := ctx.Param("assign_id")
	assignID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateAssignedEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAssignedEmployeeParams{
		ID:         assignID,
		EmployeeID: req.EmployeeID,
		StartDate:  pgtype.Date{Time: req.StartDate, Valid: true},
		Role:       req.Role,
	}

	assign, err := server.store.UpdateAssignedEmployee(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateAssignedEmployeeResponse{
		ID:         assign.ID,
		ClientID:   assign.ClientID,
		EmployeeID: assign.EmployeeID,
		StartDate:  assign.StartDate.Time,
		Role:       assign.Role,
		CreatedAt:  assign.CreatedAt.Time,
	}, "Assigned employee updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteAssignedEmployeeResponse defines the response for deleting an assigned employee
type DeleteAssignedEmployeeResponse struct {
	ID int64 `json:"id"`
}

// DeleteAssignedEmployeeApi deletes an assigned employee
// @Summary Delete an assigned employee
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param assign_id path int true "Assignment ID"
// @Success 200 {object} Response[DeleteAssignedEmployeeResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees/{assign_id} [delete]
func (server *Server) DeleteAssignedEmployeeApi(ctx *gin.Context) {
	id := ctx.Param("assign_id")
	assignID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = server.store.DeleteAssignedEmployee(ctx, assignID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(DeleteAssignedEmployeeResponse{
		ID: assignID,
	}, "Assigned employee deleted successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetClientRelatedEmailsResponse defines the response for getting client related emails
type GetClientRelatedEmailsResponse struct {
	Emails []string `json:"emails"`
}

// GetClientRelatedEmailsApi gets client related emails
// @Summary Get client related emails
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[GetClientRelatedEmailsResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/related_emails [get]
func (server *Server) GetClientRelatedEmailsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	emails, err := server.store.GetClientRelatedEmails(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetClientRelatedEmailsResponse{
		Emails: emails,
	}, "Client related emails fetched successfully")
	ctx.JSON(http.StatusOK, res)

}

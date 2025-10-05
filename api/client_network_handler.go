package api

import (
	"log"
	_ "maicare_go/pagination" // for swagger
	clientp "maicare_go/service/client"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetClientSenderApi gets a client sender
// @Summary Get a client sender
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[clientp.GetClientSenderResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/sender [get]
func (server *Server) GetClientSenderApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sender, err := server.businessService.ClientService.GetClientSender(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(sender, "Client Sender fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateClientEmergencyContactApi creates a client emergency contact
// @Summary Create a client emergency contact
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.CreateClientEmergencyContactParams true "Client emergency contact data"
// @Success 201 {object} Response[clientp.CreateClientEmergencyContactResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/emergency_contacts [post]
func (server *Server) CreateClientEmergencyContactApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.CreateClientEmergencyContactParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	clientEmergencyContact, err := server.businessService.ClientService.CreateClientEmergencyContact(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(clientEmergencyContact, "Client emergency contact created successfully")
	ctx.JSON(http.StatusCreated, res)

}

// ListClientEmergencyContactsApi lists all client emergency contacts
// @Summary List all client emergency contacts
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search query"
// @Success 200 {object} Response[clientp.ListClientEmergencyContactsResponse]
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

	var req clientp.ListClientEmergencyContactsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := server.businessService.ClientService.ListClientEmergencyContacts(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client emergency contacts fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientEmergencyContactApi gets a client emergency contact
// @Summary Get a client emergency contact
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param contact_id path int true "Contact ID"
// @Success 200 {object} Response[clientp.GetClientEmergencyContactResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/emergency_contacts/{contact_id} [get]
func (server *Server) GetClientEmergencyContactApi(ctx *gin.Context) {
	id := ctx.Param("contact_id")
	contactID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contact, err := server.businessService.ClientService.GetClientEmergencyContact(ctx, contactID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(contact, "Client emergency contact fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateClientEmergencyContactApi updates a client emergency contact
// @Summary Update a client emergency contact
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param contact_id path int true "Contact ID"
// @Param request body clientp.UpdateClientEmergencyContactParams true "Client emergency contact data"
// @Success 200 {object} Response[clientp.UpdateClientEmergencyContactResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/emergency_contacts/{contact_id} [put]
func (server *Server) UpdateClientEmergencyContactApi(ctx *gin.Context) {
	id := ctx.Param("contact_id")
	contactID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateClientEmergencyContactParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contact, err := server.businessService.ClientService.UpdateClientEmergencyContact(ctx, req, contactID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(contact, "Client emergency contact updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteClientEmergencyContactApi deletes a client emergency contact
// @Summary Delete a client emergency contact
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param contact_id path int true "Contact ID"
// @Success 200 {object} Response[clientp.DeleteClientEmergencyContactResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/emergency_contacts/{contact_id} [delete]
func (server *Server) DeleteClientEmergencyContactApi(ctx *gin.Context) {
	id := ctx.Param("contact_id")
	contactID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	cid, err := server.businessService.ClientService.DeleteClientEmergencyContact(ctx, contactID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(cid, "Client emergency contact deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// AssignEmployeeApi assigns an employee to a client
// @Summary Assign an employee to a client
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.AssignEmployeeRequest true "Employee assignment data"
// @Success 201 {object} Response[clientp.AssignEmployeeResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees [post]
func (server *Server) AssignEmployeeApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req clientp.AssignEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := server.businessService.ClientService.AssignEmployeeToClient(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(result, "Employee assigned successfully")
	ctx.JSON(http.StatusCreated, res)
}

// ListAssignedEmployeesApi lists all assigned employees
// @Summary List all assigned employees
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListAssignedEmployeesResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees [get]
func (server *Server) ListAssignedEmployeesApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.ListAssignedEmployeesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := server.businessService.ClientService.ListAssignedEmployees(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(result, "Assigned employees fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetAssignedEmployeeApi gets an assigned employee
// @Summary Get an assigned employee
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param assign_id path int true "Assignment ID"
// @Success 200 {object} Response[clientp.GetAssignedEmployeeResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees/{assign_id} [get]
func (server *Server) GetAssignedEmployeeApi(ctx *gin.Context) {
	id := ctx.Param("assign_id")
	assignID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.GetAssignedEmployee(ctx, assignID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(result, "Assigned employee fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateAssignedEmployeeApi updates an assigned employee
// @Summary Update an assigned employee
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param assign_id path int true "Assignment ID"
// @Param request body clientp.UpdateAssignedEmployeeRequest true "Assigned employee data"
// @Success 200 {object} Response[clientp.UpdateAssignedEmployeeResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees/{assign_id} [put]
func (server *Server) UpdateAssignedEmployeeApi(ctx *gin.Context) {
	id := ctx.Param("assign_id")
	assignID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateAssignedEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.UpdateAssignedEmployee(ctx, req, assignID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Assigned employee updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteAssignedEmployeeApi deletes an assigned employee
// @Summary Delete an assigned employee
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Param assign_id path int true "Assignment ID"
// @Success 200 {object} Response[clientp.DeleteAssignedEmployeeResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/involved_employees/{assign_id} [delete]
func (server *Server) DeleteAssignedEmployeeApi(ctx *gin.Context) {
	id := ctx.Param("assign_id")
	assignID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.DeleteAssignedEmployee(ctx, assignID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Assigned employee deleted successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetClientRelatedEmailsApi gets client related emails
// @Summary Get client related emails
// @Tags client_network
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[clientp.GetClientRelatedEmailsResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/related_emails [get]
func (server *Server) GetClientRelatedEmailsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.GetClientRelatedEmail(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(result, "Client related emails fetched successfully")
	ctx.JSON(http.StatusOK, res)

}

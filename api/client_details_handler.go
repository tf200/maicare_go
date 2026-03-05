package api

import (
	"fmt"
	"net/http"
	"strings"

	_ "maicare_go/pagination" // for swagger
	clientp "maicare_go/service/client"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateClientApi creates a new client
// @Summary Create a new client
// @Tags clients
// @Accept json
// @Produce json
// @Param request body clientp.CreateClientDetailsRequest true "Client details"
// @Success 201 {object} Response[clientp.CreateClientDetailsResponse]
// @Failure 400,404,500 {object} Response[clientp.CreateClientDetailsResponse]
// @Router /clients [post]
func (server *Server) CreateClientApi(ctx *gin.Context) {
	var req clientp.CreateClientDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	client, err := server.businessService.ClientService.CreateClientDetails(req, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(client, "Client created successfully")
	ctx.JSON(http.StatusCreated, res)
}

// ListClientsApi lists clients
// @Summary List clients
// @Tags clients
// @Produce json
// @Param status query string false "Client status"
// @Param location_id query int false "Location ID"
// @Param search query string false "Search query"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListClientsApiResponse]]
// @Failure 400,404,500 {object} Response[clientp.ListClientsApiResponse]
// @Router /clients [get]
func (server *Server) ListClientsApi(ctx *gin.Context) {
	var req clientp.ListClientsApiParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid query parameters: %v", err)))
		return
	}

	result, err := server.businessService.ClientService.ListClientDetails(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list clients: %v", err)))
		return
	}

	res := SuccessResponse(result, "Clients fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListWaitingListClientsApi lists waiting list clients
// @Summary List waiting list clients
// @Tags clients
// @Produce json
// @Param search query string false "Search by client first_name, last_name, or sender_name (max 120 chars)"
// @Param placement query string false "Care type filter: protected_living|training_center|supported_independent_living|ambulatory_support|other"
// @Param sort_days query string false "Sort by days_in_waitlist: asc|desc (default: desc)"
// @Param page query int true "Page number (min 1)"
// @Param page_size query int true "Page size (min 5, max 100)"
// @Success 200 {object} Response[pagination.Response[clientp.ListWaitingListClientsResponse]]
// @Failure 400,404,500 {object} Response[clientp.ListWaitingListClientsResponse]
// @Router /clients/waiting-list [get]
func (server *Server) ListWaitingListClientsApi(ctx *gin.Context) {
	var req clientp.ListWaitingListClientsParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid query parameters: %v", err)))
		return
	}

	result, err := server.businessService.ClientService.ListWaitingListClients(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list waiting list clients: %v", err)))
		return
	}

	res := SuccessResponse(result, "Waiting list clients fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListInCareClientsApi lists clients that are in care or scheduled in care
// @Summary List in-care clients
// @Tags clients
// @Produce json
// @Param search query string false "Search by first or last name (max 120 chars)"
// @Param status query []string false "Status filter" Enums(in_care, scheduled_in_care)
// @Param sort_days_in_care query string false "Sort by days_in_care: asc|desc (default: desc)"
// @Param page query int true "Page number (min 1)"
// @Param page_size query int true "Page size (min 5, max 100)"
// @Success 200 {object} Response[pagination.Response[clientp.ListInCareClientsResponse]]
// @Failure 400,404,500 {object} Response[clientp.ListInCareClientsResponse]
// @Router /clients/in-care [get]
func (server *Server) ListInCareClientsApi(ctx *gin.Context) {
	var req clientp.ListInCareClientsParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid query parameters: %v", err)))
		return
	}

	result, err := server.businessService.ClientService.ListInCareClients(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list in-care clients: %v", err)))
		return
	}

	res := SuccessResponse(result, "In-care clients fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientsCountApi gets the count of clients
// @Summary Get the count of clients
// @Tags clients
// @Produce json
// @Success 200 {object} Response[clientp.GetClientsCountResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/counts [get]
func (server *Server) GetClientsCountApi(ctx *gin.Context) {
	clientCount, err := server.businessService.ClientService.GetClientsCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(clientCount, "Clients count fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientStatusCountsApi gets grouped client status counts
// @Summary Get grouped client status counts
// @Tags clients
// @Produce json
// @Success 200 {object} Response[clientp.GetClientStatusCountsResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/status-counts [get]
func (server *Server) GetClientStatusCountsApi(ctx *gin.Context) {
	counts, err := server.businessService.ClientService.GetClientStatusCounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(counts, "Client status counts fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientApi gets a client
// @Summary Get a client
// @Tags clients
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} Response[clientp.GetClientApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id} [get]
func (server *Server) GetClientApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID: %v", err)))
		return
	}

	response, err := server.businessService.ClientService.GetClientDetails(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get client details: %v", err)))
		return
	}

	res := SuccessResponse(response, "Client fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientAddressesApi gets a client
// @Summary Get a client addresses
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} Response[clientp.GetClientAddressesApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/addresses [get]
func (server *Server) GetClientAddressesApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	response, err := server.businessService.ClientService.GetClientAddresses(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get client addresses: %v", err)))
		return
	}

	res := SuccessResponse(response, "Client addresses fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateClientApi updates a client
// @Summary Update a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body clientp.UpdateClientDetailsRequest true "Client details"
// @Success 200 {object} Response[clientp.UpdateClientDetailsResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id} [put]
func (server *Server) UpdateClientApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID: %v", err)))
		return
	}
	var req clientp.UpdateClientDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedClient, err := server.businessService.ClientService.UpdateClientDetails(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(updatedClient, "Client updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateClientStatusApi updates a client
// @Summary Update a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body clientp.UpdateClientStatusRequest true "Client status"
// @Success 200 {object} Response[clientp.UpdateClientStatusResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/status [put]
func (server *Server) UpdateClientStatusApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateClientStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedClient, err := server.businessService.ClientService.UpdateClientStatus(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(updatedClient, "Client status updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// PutClientInCareApi moves a waiting-list client into care lifecycle
// @Summary Put client in care lifecycle
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body clientp.PutClientInCareRequest true "Put in care payload"
// @Success 200 {object} Response[clientp.PutClientInCareResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/put-in-care [put]
func (server *Server) PutClientInCareApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.PutClientInCareRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.PutClientInCare(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client moved to care lifecycle successfully")
	ctx.JSON(http.StatusOK, res)
}

// PutClientOutOfCareApi moves an in-care client out of care
// @Summary Put client out of care
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body clientp.PutClientOutOfCareRequest true "Put out of care payload"
// @Success 200 {object} Response[clientp.PutClientOutOfCareResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/put-out-of-care [put]
func (server *Server) PutClientOutOfCareApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.PutClientOutOfCareRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.PutClientOutOfCare(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client moved out of care successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListStatusHistoryApi lists status history of a client
// @Summary List status history of a client
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} Response[[]clientp.ListStatusHistoryApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/status_history [get]
func (server *Server) ListStatusHistoryApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	statusHistoryList, err := server.businessService.ClientService.ListStatusHistory(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(statusHistoryList, "Status history fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// AddClientDocumentApi adds a document to a client
// @Summary Add a document to a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body clientp.AddClientDocumentApiRequest true "Client document"
// @Success 201 {object} Response[clientp.AddClientDocumentApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/documents [post]
func (server *Server) AddClientDocumentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.AddClientDocumentApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.AddClientDocument(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client document added successfully")
	ctx.JSON(http.StatusCreated, res)
}

// ListClientDocumentsApi lists documents of a client
// @Summary List documents of a client
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListClientDocumentsApiResponse]]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/documents [get]
func (server *Server) ListClientDocumentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.ListClientDocumentsApiRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pag, err := server.businessService.ClientService.ListClientDocuments(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(pag, "Client documents fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteClientDocumentApi deletes a client document
// @Summary Delete a client document
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param document_id path int true "Document ID"
// @Param request body clientp.DeleteClientDocumentApiRequest true "Client document"
// @Success 200 {object} Response[clientp.DeleteClientDocumentApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/documents/{document_id} [delete]
func (server *Server) DeleteClientDocumentApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req clientp.DeleteClientDocumentApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.DeleteClientDocument(ctx, clientID, req.AttachmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client document deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetMissingClientDocumentsApi gets missing documents of a client
// @Summary Get missing documents of a client
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} Response[clientp.GetMissingClientDocumentsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/missing_documents [get]
func (server *Server) GetMissingClientDocumentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.GetMissingClientDocuments(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Missing client documents fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// RequestLocationTransferApi handles location transfer requests
// @Summary Request a location transfer for a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body clientp.LocationTransferRequest true "Location transfer request"
// @Success 200 {object} Response[any]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/location_transfer [post]
func (server *Server) RequestLocationTransferApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.LocationTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.businessService.ClientService.RequestLocationTransfer(ctx, clientID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Location transfer request created successfully")
	ctx.JSON(http.StatusOK, res)
}

// ApproveOrRejectClientLocationTransferApi handles approval or rejection of location transfer requests
// @Summary Approve or reject a location transfer request for a client
// @Tags clients
// @Accept json
// @Produce json
// @Param request body clientp.ApproveOrRejectLocationTransferRequest true "Approve or reject location transfer request"
// @Success 200 {object} Response[any]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/location_transfer/approve_reject [post]
func (server *Server) ApproveOrRejectClientLocationTransferApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	employeeID := payload.EmployeeID

	var req clientp.ApproveOrRejectLocationTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.businessService.ClientService.ApproveLocationTransfer(ctx, employeeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Location transfer request processed successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListLocationTransferRequestsApi lists location transfer requests for a client
// @Summary List location transfer requests for a client
// @Tags clients
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListLocationTransferRequestsResponse]]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/location_transfer [get]
func (server *Server) ListLocationTransferRequestsApi(ctx *gin.Context) {

	var req clientp.ListLocationTransferRequestsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pag, err := server.businessService.ClientService.ListLocationTransferRequests(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(pag, "Location transfer requests fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetGoalEvaluationApi gets a single goal evaluation by ID.
// @Summary Get goal evaluation by ID
// @Tags evaluations
// @Produce json
// @Param evaluation_id path string true "Evaluation ID"
// @Success 200 {object} Response[clientp.GoalEvaluationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /evaluations/{evaluation_id} [get]
func (server *Server) GetGoalEvaluationApi(ctx *gin.Context) {
	evaluationID, err := uuid.Parse(ctx.Param("evaluation_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid evaluation ID")))
		return
	}

	result, err := server.businessService.ClientService.GetGoalEvaluation(ctx, evaluationID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Goal evaluation fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListUpcomingEvaluationsApi lists upcoming evaluations for the logged-in coordinator
// @Summary List upcoming evaluations for coordinator
// @Tags evaluations
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListUpcomingEvaluationsResponse]]
// @Failure 400,500 {object} Response[any]
// @Router /evaluations/upcoming [get]
func (server *Server) ListUpcomingEvaluationsApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req clientp.ListUpcomingEvaluationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.ListUpcomingEvaluations(ctx, payload.EmployeeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Upcoming evaluations fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetGoalEvaluationBootstrapApi returns create-screen bootstrap data for a client
// @Summary Get goal evaluation bootstrap data
// @Tags evaluations
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} Response[clientp.GoalEvaluationBootstrapResponse]
// @Failure 400,500 {object} Response[any]
// @Router /clients/{id}/evaluations/bootstrap [get]
func (server *Server) GetGoalEvaluationBootstrapApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	result, err := server.businessService.ClientService.GetGoalEvaluationBootstrap(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Goal evaluation bootstrap fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientGoalsForEvaluationPageApi returns client goals data for the evaluation page.
// @Summary Get client goals for evaluation page
// @Tags evaluations
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} Response[clientp.GetClientGoalsForEvaluationPageResponse]
// @Failure 400,500 {object} Response[any]
// @Router /clients/{id}/goals [get]
func (server *Server) GetClientGoalsForEvaluationPageApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.GetClientGoalsForEvaluationPage(ctx, clientID, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client goals for evaluation page fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListClientSubmittedEvaluationsApi lists submitted evaluations for a client.
// @Summary List submitted evaluations by client
// @Tags evaluations
// @Produce json
// @Param id path string true "Client ID"
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListClientSubmittedEvaluationsResponse]]
// @Failure 400,500 {object} Response[any]
// @Router /clients/{id}/evaluations/submitted [get]
func (server *Server) ListClientSubmittedEvaluationsApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	var req clientp.ListClientSubmittedEvaluationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.ListClientSubmittedEvaluations(ctx, clientID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client submitted evaluations fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListGoalEvaluationHistoryApi lists completed evaluation history points for a goal.
// @Summary List goal evaluation history
// @Tags evaluations
// @Produce json
// @Param id path string true "Client ID"
// @Param goal_id path string true "Goal ID"
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListGoalEvaluationHistoryResponse]]
// @Failure 400,500 {object} Response[any]
// @Router /clients/{id}/goals/{goal_id}/history [get]
func (server *Server) ListGoalEvaluationHistoryApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	goalID, err := uuid.Parse(ctx.Param("goal_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid goal ID")))
		return
	}

	var req clientp.ListGoalEvaluationHistoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.ListGoalEvaluationHistory(ctx, clientID, goalID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Goal evaluation history fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListRecentSubmittedEvaluationsApi lists recently submitted evaluations for the logged-in user
// @Summary List recent submitted evaluations
// @Tags evaluations
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListRecentSubmittedEvaluationsResponse]]
// @Failure 400,500 {object} Response[any]
// @Router /evaluations/recent-submitted [get]
func (server *Server) ListRecentSubmittedEvaluationsApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req clientp.ListRecentSubmittedEvaluationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.ListRecentSubmittedEvaluations(ctx, payload.EmployeeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Recent submitted evaluations fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListRecentDraftEvaluationsApi lists recent draft evaluations for the logged-in user
// @Summary List recent draft evaluations
// @Tags evaluations
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListRecentDraftEvaluationsResponse]]
// @Failure 400,500 {object} Response[any]
// @Router /evaluations/recent-drafts [get]
func (server *Server) ListRecentDraftEvaluationsApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req clientp.ListRecentDraftEvaluationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.ListRecentDraftEvaluations(ctx, payload.EmployeeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Recent draft evaluations fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateGoalEvaluationApi creates or updates the current scheduled goal evaluation for a client
// @Summary Save or submit the current scheduled goal evaluation
// @Tags evaluations
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param request body clientp.CreateGoalEvaluationRequest true "Goal evaluation details"
// @Success 200 {object} Response[clientp.GoalEvaluationResponse]
// @Failure 400,500 {object} Response[any]
// @Router /clients/{id}/evaluations [post]
func (server *Server) CreateGoalEvaluationApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	employeeID := payload.EmployeeID

	var req clientp.CreateGoalEvaluationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.CreateGoalEvaluation(ctx, clientID, employeeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Goal evaluation saved successfully")
	ctx.JSON(http.StatusOK, res)
}

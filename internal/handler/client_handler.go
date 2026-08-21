package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterClientRoutes(
	rg *gin.RouterGroup,
	handler *ClientHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	clientsGroup := rg.Group("/clients")
	{
		clientsGroup.POST("", auth, requirePermission("CLIENT.CREATE"), handler.CreateClient)
		clientsGroup.GET("", auth, requirePermission("CLIENT.VIEW"), handler.ListClients)
		clientsGroup.GET("/waiting-list", auth, requirePermission("CLIENT.VIEW"), handler.ListWaitingListClients)
		clientsGroup.GET("/waiting-list/stats", auth, requirePermission("CLIENT.VIEW"), handler.GetWaitingListStats)
		clientsGroup.GET("/in-care", auth, requirePermission("CLIENT.VIEW"), handler.ListInCareClients)
		clientsGroup.GET("/counts", auth, requirePermission("CLIENT.VIEW"), handler.GetClientsCount)
		clientsGroup.GET("/incare/stats", auth, requirePermission("CLIENT.VIEW"), handler.GetInCareStats)
		clientsGroup.GET("/status-counts", auth, requirePermission("CLIENT.VIEW"), handler.GetClientStatusCounts)
		clientsGroup.GET("/:id", auth, requirePermission("CLIENT.VIEW"), handler.GetClient)
		clientsGroup.PUT("/:id", auth, requirePermission("CLIENT.UPDATE"), handler.UpdateClient)
		clientsGroup.GET("/:id/addresses", auth, requirePermission("CLIENT.VIEW"), handler.GetClientAddresses)
		clientsGroup.PUT("/:id/status", auth, requirePermission("CLIENT.STATUS.UPDATE"), handler.UpdateClientStatus)
		clientsGroup.PUT("/:id/put-in-care", auth, requirePermission("CLIENT.STATUS.UPDATE"), handler.PutClientInCare)
		clientsGroup.PUT("/:id/put-out-of-care", auth, requirePermission("CLIENT.STATUS.UPDATE"), handler.PutClientOutOfCare)
		clientsGroup.GET("/:id/status_history", auth, requirePermission("CLIENT.VIEW"), handler.ListStatusHistory)
		clientsGroup.POST("/:id/documents", auth, requirePermission("CLIENT.CREATE"), handler.AddClientDocument)
		clientsGroup.GET("/:id/documents", auth, requirePermission("CLIENT.VIEW"), handler.ListClientDocuments)
		clientsGroup.DELETE("/:id/documents/:doc_id", auth, requirePermission("CLIENT.VIEW"), handler.DeleteClientDocument)
		clientsGroup.GET("/:id/missing_documents", auth, requirePermission("CLIENT.CREATE"), handler.GetMissingClientDocuments)
		clientsGroup.GET("/:id/evaluations/bootstrap", auth, requirePermission(domain.PermClientEvaluationView.String()), handler.GetGoalEvaluationBootstrap)
		clientsGroup.POST("/:id/goals", auth, requirePermission("CLIENT.UPDATE"), handler.CreateClientGoal)
		clientsGroup.PATCH("/:id/goals/:goal_id", auth, requirePermission("CLIENT.UPDATE"), handler.UpdateClientGoal)
		clientsGroup.GET("/:id/goals", auth, requirePermission("CLIENT.VIEW"), handler.GetClientGoalsForEvaluationPage)
		clientsGroup.GET("/:id/goals/:goal_id/history", auth, requirePermission(domain.PermClientEvaluationView.String()), handler.ListGoalEvaluationHistory)
		clientsGroup.GET("/:id/evaluations/submitted", auth, requirePermission(domain.PermClientEvaluationView.String()), handler.ListClientSubmittedEvaluations)
		clientsGroup.POST("/:id/evaluations", auth, requirePermission(domain.PermClientEvaluationCreate.String()), handler.CreateGoalEvaluation)
		clientsGroup.POST("/:id/location_transfer", auth, requirePermission("CLIENT.UPDATE"), handler.RequestLocationTransfer)
		clientsGroup.POST("/location_transfer/approve_reject", auth, requirePermission("CLIENT.UPDATE"), handler.ApproveOrRejectLocationTransfer)
		clientsGroup.GET("/location_transfer", auth, requirePermission("CLIENT.VIEW"), handler.ListLocationTransferRequests)

		// Medical
		clientsGroup.GET("/:id/medical/overview", auth,
			requirePermission("CLIENT.DIAGNOSIS.VIEW"),
			requirePermission("CLIENT.MEDICATION.VIEW"),
			handler.GetClientMedicalOverview,
		)
		clientsGroup.POST("/:id/medical/diagnoses", auth, requirePermission("CLIENT.DIAGNOSIS.CREATE"), handler.CreateClientDiagnosis)
		clientsGroup.GET("/:id/medical/diagnoses", auth, requirePermission("CLIENT.DIAGNOSIS.VIEW"), handler.ListClientDiagnoses)
		clientsGroup.GET("/:id/medical/diagnoses/:diagnosis_id", auth, requirePermission("CLIENT.DIAGNOSIS.VIEW"), handler.GetClientDiagnosis)
		clientsGroup.PUT("/:id/medical/diagnoses/:diagnosis_id", auth, requirePermission("CLIENT.DIAGNOSIS.UPDATE"), handler.UpdateClientDiagnosis)
		clientsGroup.DELETE("/:id/medical/diagnoses/:diagnosis_id", auth, requirePermission("CLIENT.DIAGNOSIS.DELETE"), handler.DeleteClientDiagnosis)
		clientsGroup.POST("/:id/medical/medication-orders", auth, requirePermission("CLIENT.MEDICATION.CREATE"), handler.CreateClientMedicationOrder)
		clientsGroup.GET("/:id/medical/medication-orders", auth, requirePermission("CLIENT.MEDICATION.VIEW"), handler.ListClientMedicationOrders)
		clientsGroup.GET("/:id/medical/medication-orders/:order_id", auth, requirePermission("CLIENT.MEDICATION.VIEW"), handler.GetClientMedicationOrder)
		clientsGroup.PUT("/:id/medical/medication-orders/:order_id", auth, requirePermission("CLIENT.MEDICATION.UPDATE"), handler.UpdateClientMedicationOrder)
		clientsGroup.DELETE("/:id/medical/medication-orders/:order_id", auth, requirePermission("CLIENT.MEDICATION.DELETE"), handler.DeleteClientMedicationOrder)

		// Network
		clientsGroup.GET("/:id/sender", auth, requirePermission("CLIENT.VIEW"), handler.GetClientSender)
		clientsGroup.POST("/:id/emergency_contacts", auth, requirePermission("CLIENT.EMERGENCY_CONTACT.CREATE"), handler.CreateClientEmergencyContact)
		clientsGroup.GET("/:id/emergency_contacts", auth, requirePermission("CLIENT.EMERGENCY_CONTACT.VIEW"), handler.ListClientEmergencyContacts)
		clientsGroup.GET("/:id/emergency_contacts/:contact_id", auth, requirePermission("CLIENT.EMERGENCY_CONTACT.VIEW"), handler.GetClientEmergencyContact)
		clientsGroup.PUT("/:id/emergency_contacts/:contact_id", auth, requirePermission("CLIENT.EMERGENCY_CONTACT.UPDATE"), handler.UpdateClientEmergencyContact)
		clientsGroup.DELETE("/:id/emergency_contacts/:contact_id", auth, requirePermission("CLIENT.EMERGENCY_CONTACT.DELETE"), handler.DeleteClientEmergencyContact)
		clientsGroup.POST("/:id/involved_employees", auth, requirePermission("CLIENT.INVOLVED_EMPLOYEE.CREATE"), handler.CreateAssignedEmployee)
		clientsGroup.GET("/:id/involved_employees", auth, requirePermission("CLIENT.INVOLVED_EMPLOYEE.VIEW"), handler.ListAssignedEmployees)
		clientsGroup.GET("/:id/involved_employees/:assign_id", auth, requirePermission("CLIENT.INVOLVED_EMPLOYEE.VIEW"), handler.GetAssignedEmployee)
		clientsGroup.PUT("/:id/involved_employees/:assign_id", auth, requirePermission("CLIENT.INVOLVED_EMPLOYEE.UPDATE"), handler.UpdateAssignedEmployee)
		clientsGroup.DELETE("/:id/involved_employees/:assign_id", auth, requirePermission("CLIENT.INVOLVED_EMPLOYEE.DELETE"), handler.DeleteAssignedEmployee)
		clientsGroup.GET("/:id/related_emails", auth, requirePermission("CLIENT.VIEW"), handler.GetClientRelatedEmails)

		// Progress Reports
		clientsGroup.POST("/:id/progress_reports", auth, requirePermission("CLIENT.PROGRESS_REPORT.CREATE"), handler.CreateProgressReport)
		clientsGroup.GET("/:id/progress_reports", auth, requirePermission("CLIENT.PROGRESS_REPORT.VIEW"), handler.ListProgressReports)
		clientsGroup.GET("/:id/progress_reports/:report_id", auth, requirePermission("CLIENT.PROGRESS_REPORT.VIEW"), handler.GetProgressReport)
		clientsGroup.PUT("/:id/progress_reports/:report_id", auth, requirePermission("CLIENT.PROGRESS_REPORT.UPDATE"), handler.UpdateProgressReport)
		clientsGroup.DELETE("/:id/progress_reports/:report_id", auth, requirePermission("CLIENT.PROGRESS_REPORT.DELETE"), handler.DeleteProgressReport)
		clientsGroup.POST("/:id/ai_progress_reports", auth, requirePermission("CLIENT.AI_PROGRESS_REPORT.GENERATE"), handler.GenerateAutoReports)
		clientsGroup.POST("/:id/ai_progress_reports/confirm", auth, requirePermission("CLIENT.AI_PROGRESS_REPORT.CONFIRM"), handler.ConfirmAiProgressReport)
		clientsGroup.GET("/:id/ai_progress_reports", auth, requirePermission("CLIENT.AI_PROGRESS_REPORT.VIEW"), handler.ListAiGeneratedReports)

		// Appointment Cards
		clientsGroup.GET("/:id/appointment_cards", auth, requirePermission("APPOINTMENT_CARD.VIEW"), handler.GetAppointmentCard)
		clientsGroup.PUT("/:id/appointment_cards", auth, requirePermission("APPOINTMENT_CARD.UPDATE"), handler.UpdateAppointmentCard)
		clientsGroup.POST("/:id/appointment_cards/generate_document", auth, requirePermission("APPOINTMENT_CARD.GENERATE_DOCUMENT"), handler.GenerateAppointmentCardDocument)
	}
}

func RegisterEvaluationRoutes(
	rg *gin.RouterGroup,
	handler *ClientHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	evaluationsGroup := rg.Group("/evaluations")
	{
		evaluationsGroup.GET("/upcoming", auth, requirePermission(domain.PermClientEvaluationView.String()), handler.ListUpcomingEvaluations)
		evaluationsGroup.GET("/recent-submitted", auth, requirePermission(domain.PermClientEvaluationView.String()), handler.ListRecentSubmittedEvaluations)
		evaluationsGroup.GET("/recent-drafts", auth, requirePermission(domain.PermClientEvaluationView.String()), handler.ListRecentDraftEvaluations)
		evaluationsGroup.GET("/:evaluation_id", auth, requirePermission(domain.PermClientEvaluationView.String()), handler.GetGoalEvaluation)
	}
}

type ClientHandler struct {
	service domain.ClientService
}

func NewClientHandler(service domain.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

func getEmployeeIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	idStr := ctx.GetString("employee_id")
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("employee ID not found in context")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid employee ID")
	}
	return id, nil
}

// CreateClient creates a new client
// @Summary Create a new client
// @Tags clients
// @Accept json
// @Produce json
// @Param request body createClientRequest true "Client details"
// @Success 201 {object} httpapi.Envelope[createClientResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[createClientResponse]
// @Router /clients [post]
func (h *ClientHandler) CreateClient(ctx *gin.Context) {
	var req createClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	params, err := toCreateClientParams(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid date of birth format", ""))
		return
	}

	client, err := h.service.CreateClient(ctx.Request.Context(), params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create client", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toCreateClientResponse(client), "Client created successfully"))
}

// ListClients lists clients
// @Summary List clients
// @Tags clients
// @Produce json
// @Param status query string false "Client status"
// @Param location_id query int false "Location ID"
// @Param search query string false "Search query"
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[listClientsResponse]]
// @Failure 400,404,500 {object} httpapi.Envelope[listClientsResponse]
// @Router /clients [get]
func (h *ClientHandler) ListClients(ctx *gin.Context) {
	var req listClientsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	page, err := h.service.ListClients(ctx.Request.Context(), toListClientsParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list clients", ""))
		return
	}

	results := make([]listClientsResponse, len(page.Items))
	for i, item := range page.Items {
		results[i] = toListClientsResponse(item)
	}

	response := httpapi.NewPageResponse(ctx, req.PageRequest, results, page.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(response, "Clients fetched successfully"))
}

// ListWaitingListClients lists waiting list clients
// @Summary List waiting list clients
// @Tags clients
// @Produce json
// @Param search query string false "Search by client first_name, last_name, or sender_name (max 120 chars)"
// @Param placement query string false "Care type filter: protected_living|training_center|supported_independent_living|ambulatory_support|other"
// @Param sort_days query string false "Sort by days_in_waitlist: asc|desc (default: desc)"
// @Param page query int true "Page number (min 1)"
// @Param page_size query int true "Page size (min 5, max 100)"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[listWaitingListClientsResponse]]
// @Failure 400,404,500 {object} httpapi.Envelope[listWaitingListClientsResponse]
// @Router /clients/waiting-list [get]
func (h *ClientHandler) ListWaitingListClients(ctx *gin.Context) {
	var req listWaitingListClientsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	page, err := h.service.ListWaitingListClients(ctx.Request.Context(), toListWaitingListClientsParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list waiting list clients", ""))
		return
	}

	results := make([]listWaitingListClientsResponse, len(page.Items))
	for i, item := range page.Items {
		results[i] = toListWaitingListClientsResponse(item)
	}

	response := httpapi.NewPageResponse(ctx, req.PageRequest, results, page.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(response, "Waiting list clients fetched successfully"))
}

// ListInCareClients lists clients that are in care or scheduled in care
// @Summary List in-care clients
// @Tags clients
// @Produce json
// @Param search query string false "Search by first or last name (max 120 chars)"
// @Param status query []string false "Status filter" Enums(in_care, scheduled_in_care)
// @Param sort_days_in_care query string false "Sort by days_in_care: asc|desc (default: desc)"
// @Param page query int true "Page number (min 1)"
// @Param page_size query int true "Page size (min 5, max 100)"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[listInCareClientsResponse]]
// @Failure 400,404,500 {object} httpapi.Envelope[listInCareClientsResponse]
// @Router /clients/in-care [get]
func (h *ClientHandler) ListInCareClients(ctx *gin.Context) {
	var req listInCareClientsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	page, err := h.service.ListInCareClients(ctx.Request.Context(), toListInCareClientsParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list in-care clients", ""))
		return
	}

	results := make([]listInCareClientsResponse, len(page.Items))
	for i, item := range page.Items {
		results[i] = toListInCareClientsResponse(item)
	}

	response := httpapi.NewPageResponse(ctx, req.PageRequest, results, page.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(response, "In-care clients fetched successfully"))
}

// GetClientsCount gets the count of clients
// @Summary Get the count of clients
// @Tags clients
// @Produce json
// @Success 200 {object} httpapi.Envelope[getClientsCountResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/counts [get]
func (h *ClientHandler) GetClientsCount(ctx *gin.Context) {
	counts, err := h.service.GetClientCounts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get clients count", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetClientsCountResponse(counts), "Clients count fetched successfully"))
}

// GetClientStatusCounts gets grouped client status counts
// @Summary Get grouped client status counts
// @Tags clients
// @Produce json
// @Success 200 {object} httpapi.Envelope[getClientStatusCountsResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/status-counts [get]
func (h *ClientHandler) GetClientStatusCounts(ctx *gin.Context) {
	counts, err := h.service.GetClientStatusCounts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get client status counts", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetClientStatusCountsResponse(counts), "Client status counts fetched successfully"))
}

// GetInCareStats gets in-care statistics
// @Summary Get in-care statistics
// @Tags clients
// @Produce json
// @Success 200 {object} httpapi.Envelope[getInCareStatsResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/incare/stats [get]
func (h *ClientHandler) GetInCareStats(ctx *gin.Context) {
	stats, err := h.service.GetInCareStats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get in-care stats", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetInCareStatsResponse(stats), "In-care stats fetched successfully"))
}

// GetWaitingListStats gets waiting list statistics
// @Summary Get waiting list statistics
// @Tags clients
// @Produce json
// @Success 200 {object} httpapi.Envelope[getWaitingListStatsResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/waiting-list/stats [get]
func (h *ClientHandler) GetWaitingListStats(ctx *gin.Context) {
	stats, err := h.service.GetWaitingListStats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get waiting list stats", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetWaitingListStatsResponse(stats), "Waiting list stats fetched successfully"))
}

// GetClient gets a client by ID
// @Summary Get a client
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} httpapi.Envelope[getClientResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id} [get]
func (h *ClientHandler) GetClient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	detail, err := h.service.GetClientByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			ctx.JSON(http.StatusNotFound, httpapi.Fail("client not found", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get client", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetClientResponse(detail), "Client fetched successfully"))
}

// UpdateClient updates a client
// @Summary Update a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body updateClientRequest true "Client details"
// @Success 200 {object} httpapi.Envelope[updateClientResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id} [put]
func (h *ClientHandler) UpdateClient(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req updateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	client, err := h.service.UpdateClient(ctx.Request.Context(), id, toUpdateClientParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update client", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toUpdateClientResponse(client), "Client updated successfully"))
}

// GetClientAddresses gets a client's addresses
// @Summary Get a client addresses
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} httpapi.Envelope[getClientAddressesResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/addresses [get]
func (h *ClientHandler) GetClientAddresses(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	addresses, err := h.service.GetClientAddresses(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get client addresses", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetClientAddressesResponse(addresses), "Client addresses fetched successfully"))
}

// UpdateClientStatus updates a client status
// @Summary Update a client status
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body updateClientStatusRequest true "Status update details"
// @Success 200 {object} httpapi.Envelope[updateClientStatusResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/status [put]
func (h *ClientHandler) UpdateClientStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req updateClientStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.UpdateClientStatus(ctx.Request.Context(), id, toUpdateClientStatusParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update client status", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toUpdateClientStatusResponse(result), "Client status updated successfully"))
}

// PutClientInCare puts a client in care
// @Summary Put a client in care
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body putClientInCareRequest true "Put in care details"
// @Success 200 {object} httpapi.Envelope[putClientInCareResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/put-in-care [put]
func (h *ClientHandler) PutClientInCare(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req putClientInCareRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.PutClientInCare(ctx.Request.Context(), id, toPutClientInCareParams(req))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toPutClientInCareResponse(result), "Client moved to care lifecycle successfully"))
}

// PutClientOutOfCare puts a client out of care
// @Summary Put a client out of care
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body putClientOutOfCareRequest true "Put out of care details"
// @Success 200 {object} httpapi.Envelope[putClientOutOfCareResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/put-out-of-care [put]
func (h *ClientHandler) PutClientOutOfCare(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req putClientOutOfCareRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.PutClientOutOfCare(ctx.Request.Context(), id, toPutClientOutOfCareParams(req))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toPutClientOutOfCareResponse(result), "Client moved out of care successfully"))
}

// ListStatusHistory lists a client's status history
// @Summary List client status history
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} httpapi.Envelope[[]listStatusHistoryResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/status_history [get]
func (h *ClientHandler) ListStatusHistory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	histories, err := h.service.ListStatusHistory(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list status history", ""))
		return
	}

	results := make([]listStatusHistoryResponse, len(histories))
	for i, history := range histories {
		results[i] = toListStatusHistoryResponse(history)
	}

	ctx.JSON(http.StatusOK, httpapi.OK(results, "Status history fetched successfully"))
}

// AddClientDocument adds documents to a client
// @Summary Add documents to a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body addClientDocumentRequest true "Client documents"
// @Success 201 {object} httpapi.Envelope[addClientDocumentResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/documents [post]
func (h *ClientHandler) AddClientDocument(ctx *gin.Context) {
	idStr := ctx.Param("id")
	clientID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req addClientDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	results, err := h.service.AddClientDocument(ctx.Request.Context(), clientID, toAddClientDocumentParams(req))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toAddClientDocumentResponse(results), "Client documents added successfully"))
}

// ListClientDocuments lists documents of a client
// @Summary List documents of a client
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[listClientDocumentsResponse]]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/documents [get]
func (h *ClientHandler) ListClientDocuments(ctx *gin.Context) {
	idStr := ctx.Param("id")
	clientID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req listClientDocumentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	params := req.Params()
	result, err := h.service.ListClientDocuments(ctx.Request.Context(), domain.ListClientDocumentsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list client documents", ""))
		return
	}

	items := make([]listClientDocumentsResponse, 0, len(result.Documents))
	for _, doc := range result.Documents {
		items = append(items, toListClientDocumentsResponse(doc))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Client documents fetched successfully"))
}

// DeleteClientDocument deletes a client document
// @Summary Delete a client document
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Param doc_id path uuid true "Document ID"
// @Success 200 {object} httpapi.Envelope[deleteClientDocumentResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/documents/{doc_id} [delete]
func (h *ClientHandler) DeleteClientDocument(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}
	documentID, err := uuid.Parse(ctx.Param("doc_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid document ID", ""))
		return
	}

	result, err := h.service.DeleteClientDocument(ctx.Request.Context(), clientID, documentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete client document", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toDeleteClientDocumentResponse(result), "Client document deleted successfully"))
}

// GetMissingClientDocuments gets missing documents of a client
// @Summary Get missing documents of a client
// @Tags clients
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} httpapi.Envelope[getMissingClientDocumentsResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/missing_documents [get]
func (h *ClientHandler) GetMissingClientDocuments(ctx *gin.Context) {
	idStr := ctx.Param("id")
	clientID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	missingDocs, err := h.service.GetMissingClientDocuments(ctx.Request.Context(), clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get missing client documents", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(getMissingClientDocumentsResponse{MissingDocs: missingDocs}, "Missing client documents fetched successfully"))
}

// GetGoalEvaluationBootstrap gets goal evaluation bootstrap data
// @Summary Get goal evaluation bootstrap
// @Tags evaluations
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} httpapi.Envelope[getGoalEvaluationBootstrapResponse]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/evaluations/bootstrap [get]
func (h *ClientHandler) GetGoalEvaluationBootstrap(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	result, err := h.service.GetGoalEvaluationBootstrap(ctx.Request.Context(), clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get goal evaluation bootstrap", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetGoalEvaluationBootstrapResponse(result), "Goal evaluation bootstrap fetched successfully"))
}

// CreateClientGoal creates a manual goal on a client
// @Summary Create client goal
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param request body createClientGoalRequest true "Client goal payload"
// @Success 201 {object} httpapi.Envelope[createClientGoalResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/goals [post]
func (h *ClientHandler) CreateClientGoal(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req createClientGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.CreateClientGoal(ctx.Request.Context(), clientID, toCreateClientGoalParams(req))
	if err != nil {
		if errors.Is(err, domain.ErrClientGoalTitleRequired) {
			ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
			return
		}
		if errors.Is(err, domain.ErrClientGoalClientNotFound) || errors.Is(err, domain.ErrClientGoalTopicNotFound) {
			ctx.JSON(http.StatusNotFound, httpapi.Fail(err.Error(), ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create client goal", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toCreateClientGoalResponse(*result), "Client goal created successfully"))
}

// UpdateClientGoal updates a client goal
// @Summary Update client goal
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param goal_id path string true "Goal ID"
// @Param request body updateClientGoalRequest true "Client goal update payload"
// @Success 200 {object} httpapi.Envelope[updateClientGoalResponse]
// @Failure 400,404,409,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/goals/{goal_id} [patch]
func (h *ClientHandler) UpdateClientGoal(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	goalID, err := uuid.Parse(ctx.Param("goal_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid goal ID", ""))
		return
	}

	var req updateClientGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.UpdateClientGoal(ctx.Request.Context(), clientID, goalID, toUpdateClientGoalParams(req))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrClientGoalEmptyPatch), errors.Is(err, domain.ErrClientGoalTitleRequired):
			ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
			return
		case errors.Is(err, domain.ErrClientGoalClientNotFound), errors.Is(err, domain.ErrClientGoalGoalNotFound), errors.Is(err, domain.ErrClientGoalTopicNotFound):
			ctx.JSON(http.StatusNotFound, httpapi.Fail(err.Error(), ""))
			return
		case errors.Is(err, domain.ErrClientGoalDraftEvaluationExists):
			ctx.JSON(http.StatusConflict, httpapi.Fail(err.Error(), ""))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update client goal", ""))
			return
		}
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toUpdateClientGoalResponse(result), "Client goal updated successfully"))
}

// GetClientGoalsForEvaluationPage returns goals for evaluation page
// @Summary Get client goals for evaluation page
// @Tags evaluations
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} httpapi.Envelope[getClientGoalsForEvaluationPageResponse]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/goals [get]
func (h *ClientHandler) GetClientGoalsForEvaluationPage(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	// Get employeeID from auth context
	employeeID, err := getEmployeeIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.GetClientGoalsForEvaluationPage(ctx.Request.Context(), clientID, employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get goals for evaluation page", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetClientGoalsForEvaluationPageResponse(result), "Client goals for evaluation page fetched successfully"))
}

// ListGoalEvaluationHistory lists completed evaluation history for a goal
// @Summary List goal evaluation history
// @Tags evaluations
// @Produce json
// @Param id path string true "Client ID"
// @Param goal_id path string true "Goal ID"
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[listGoalEvaluationHistoryResponse]]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/goals/{goal_id}/history [get]
func (h *ClientHandler) ListGoalEvaluationHistory(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	goalID, err := uuid.Parse(ctx.Param("goal_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid goal ID", ""))
		return
	}

	var req listGoalEvaluationHistoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	params := req.Params()
	result, err := h.service.ListGoalEvaluationHistory(ctx.Request.Context(), domain.ListGoalEvaluationHistoryParams{
		ClientID: clientID,
		GoalID:   goalID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list goal evaluation history", ""))
		return
	}

	items := make([]listGoalEvaluationHistoryResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toListGoalEvaluationHistoryResponse(item))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Goal evaluation history fetched successfully"))
}

// ListClientSubmittedEvaluations lists submitted evaluations for a client
// @Summary List submitted evaluations by client
// @Tags evaluations
// @Produce json
// @Param id path string true "Client ID"
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[listClientSubmittedEvaluationsResponse]]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/evaluations/submitted [get]
func (h *ClientHandler) ListClientSubmittedEvaluations(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req listClientSubmittedEvaluationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	params := req.Params()
	result, err := h.service.ListClientSubmittedEvaluations(ctx.Request.Context(), domain.ListClientSubmittedEvaluationsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list submitted evaluations", ""))
		return
	}

	items := make([]listClientSubmittedEvaluationsResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toListClientSubmittedEvaluationsResponse(item))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Client submitted evaluations fetched successfully"))
}

// CreateGoalEvaluation creates or updates the current scheduled goal evaluation
// @Summary Save or submit the current scheduled goal evaluation
// @Tags evaluations
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param request body createGoalEvaluationRequest true "Goal evaluation details"
// @Success 200 {object} httpapi.Envelope[goalEvaluationResponse]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/evaluations [post]
func (h *ClientHandler) CreateGoalEvaluation(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req createGoalEvaluationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	employeeID, err := getEmployeeIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.CreateGoalEvaluation(ctx.Request.Context(), clientID, employeeID, toCreateGoalEvaluationParams(req))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	resp := toGoalEvaluationResponse(*result)
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Goal evaluation saved successfully"))
}

// RequestLocationTransfer handles location transfer requests
// @Summary Request a location transfer for a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param request body requestLocationTransferRequest true "Location transfer request"
// @Success 200 {object} httpapi.Envelope[any]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /clients/{id}/location_transfer [post]
func (h *ClientHandler) RequestLocationTransfer(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req requestLocationTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	err = h.service.RequestLocationTransfer(ctx.Request.Context(), clientID, toCreateLocationTransferParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Location transfer request created successfully"))
}

// ApproveOrRejectLocationTransfer handles approval or rejection of location transfer requests
// @Summary Approve or reject a location transfer request
// @Tags clients
// @Accept json
// @Produce json
// @Param request body approveOrRejectLocationTransferRequest true "Approve or reject location transfer request"
// @Success 200 {object} httpapi.Envelope[any]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /clients/location_transfer/approve_reject [post]
func (h *ClientHandler) ApproveOrRejectLocationTransfer(ctx *gin.Context) {
	employeeID, err := getEmployeeIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail(err.Error(), ""))
		return
	}

	var req approveOrRejectLocationTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	err = h.service.ApproveLocationTransfer(ctx.Request.Context(), employeeID, toApproveLocationTransferParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Location transfer request processed successfully"))
}

// ListLocationTransferRequests lists location transfer requests
// @Summary List location transfer requests
// @Tags clients
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[locationTransferResponse]]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /clients/location_transfer [get]
func (h *ClientHandler) ListLocationTransferRequests(ctx *gin.Context) {
	var req listLocationTransferRequestsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListLocationTransferRequests(ctx.Request.Context(), domain.ListLocationTransferParams{
		Limit:  pageParams.Limit,
		Offset: pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), ""))
		return
	}

	items := make([]locationTransferResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toLocationTransferResponse(item))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Location transfer requests fetched successfully"))
}

// GetGoalEvaluation gets a single goal evaluation by ID.
// @Summary Get goal evaluation by ID
// @Tags evaluations
// @Produce json
// @Param evaluation_id path string true "Evaluation ID"
// @Success 200 {object} httpapi.Envelope[goalEvaluationResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /evaluations/{evaluation_id} [get]
func (h *ClientHandler) GetGoalEvaluation(ctx *gin.Context) {
	evaluationID, err := uuid.Parse(ctx.Param("evaluation_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid evaluation ID", ""))
		return
	}

	result, err := h.service.GetGoalEvaluation(ctx.Request.Context(), evaluationID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, httpapi.Fail(err.Error(), ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGoalEvaluationResponse(*result), "Goal evaluation fetched successfully"))
}

// ListUpcomingEvaluations lists upcoming evaluations for the logged-in coordinator.
// @Summary List upcoming evaluations for coordinator
// @Tags evaluations
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[upcomingEvaluationResponse]]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /evaluations/upcoming [get]
func (h *ClientHandler) ListUpcomingEvaluations(ctx *gin.Context) {
	employeeID, err := getEmployeeIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail(err.Error(), ""))
		return
	}

	var req listUpcomingEvaluationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListUpcomingEvaluations(ctx.Request.Context(), domain.ListUpcomingEvaluationsParams{
		EmployeeID: employeeID,
		Limit:      pageParams.Limit,
		Offset:     pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), ""))
		return
	}

	items := make([]upcomingEvaluationResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toUpcomingEvaluationResponse(item))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Upcoming evaluations fetched successfully"))
}

// ListRecentSubmittedEvaluations lists recently submitted evaluations for the logged-in user.
// @Summary List recent submitted evaluations
// @Tags evaluations
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[recentSubmittedEvaluationResponse]]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /evaluations/recent-submitted [get]
func (h *ClientHandler) ListRecentSubmittedEvaluations(ctx *gin.Context) {
	employeeID, err := getEmployeeIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail(err.Error(), ""))
		return
	}

	var req listRecentSubmittedEvaluationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListRecentSubmittedEvaluations(ctx.Request.Context(), domain.ListRecentSubmittedEvaluationsParams{
		EmployeeID: employeeID,
		Limit:      pageParams.Limit,
		Offset:     pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), ""))
		return
	}

	items := make([]recentSubmittedEvaluationResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toRecentSubmittedEvaluationResponse(item))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Recent submitted evaluations fetched successfully"))
}

// ListRecentDraftEvaluations lists recent draft evaluations for the logged-in user.
// @Summary List recent draft evaluations
// @Tags evaluations
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[recentDraftEvaluationResponse]]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /evaluations/recent-drafts [get]
func (h *ClientHandler) ListRecentDraftEvaluations(ctx *gin.Context) {
	employeeID, err := getEmployeeIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail(err.Error(), ""))
		return
	}

	var req listRecentDraftEvaluationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListRecentDraftEvaluations(ctx.Request.Context(), domain.ListRecentDraftEvaluationsParams{
		EmployeeID: employeeID,
		Limit:      pageParams.Limit,
		Offset:     pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), ""))
		return
	}

	items := make([]recentDraftEvaluationResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toRecentDraftEvaluationResponse(item))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Recent draft evaluations fetched successfully"))
}

// =====================
// Medical - Overview
// =====================

// GetClientMedicalOverview returns diagnoses + active medication orders
// @Summary Get client medical overview
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} httpapi.Envelope[clientMedicalOverviewResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/overview [get]
func (h *ClientHandler) GetClientMedicalOverview(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	result, err := h.service.GetClientMedicalOverview(ctx.Request.Context(), clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get client medical overview", ""))
		return
	}

	diagnoses := make([]clientDiagnosisResponse, 0, len(result.Diagnoses))
	for _, d := range result.Diagnoses {
		diagnoses = append(diagnoses, toClientDiagnosisResponse(d))
	}

	orders := make([]clientMedicationOrderResponse, 0, len(result.MedicationOrders))
	for _, o := range result.MedicationOrders {
		orders = append(orders, toClientMedicationOrderResponse(o))
	}

	ctx.JSON(http.StatusOK, httpapi.OK(clientMedicalOverviewResponse{
		Diagnoses:        diagnoses,
		MedicationOrders: orders,
	}, "Client medical overview fetched successfully"))
}

// =====================
// Medical - Diagnoses
// =====================

// CreateClientDiagnosis creates a client diagnosis
// @Summary Create a client diagnosis
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body createClientDiagnosisRequest true "Client diagnosis data"
// @Success 201 {object} httpapi.Envelope[clientDiagnosisResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/diagnoses [post]
func (h *ClientHandler) CreateClientDiagnosis(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req createClientDiagnosisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	employeeID, _ := getEmployeeIDFromContext(ctx)

	result, err := h.service.CreateClientDiagnosis(ctx.Request.Context(), clientID, employeeID, toDomainCreateClientDiagnosisParams(req, clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create client diagnosis", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toClientDiagnosisResponse(*result), "Client diagnosis created successfully"))
}

// ListClientDiagnoses lists client diagnoses
// @Summary List client diagnoses
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[clientDiagnosisResponse]]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/diagnoses [get]
func (h *ClientHandler) ListClientDiagnoses(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req listClientDiagnosesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to bind query parameters", ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListClientDiagnoses(ctx.Request.Context(), domain.ListClientDiagnosesParams{
		ClientID: clientID,
		Limit:    pageParams.Limit,
		Offset:   pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list client diagnoses", ""))
		return
	}

	items := make([]clientDiagnosisResponse, 0, len(result.Items))
	for _, d := range result.Items {
		items = append(items, toClientDiagnosisResponse(d))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Client diagnoses fetched successfully"))
}

// GetClientDiagnosis gets a client diagnosis
// @Summary Get a client diagnosis
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param diagnosis_id path uuid true "Diagnosis ID"
// @Success 200 {object} httpapi.Envelope[clientDiagnosisResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/diagnoses/{diagnosis_id} [get]
func (h *ClientHandler) GetClientDiagnosis(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}
	diagnosisID, err := uuid.Parse(ctx.Param("diagnosis_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid diagnosis ID", ""))
		return
	}

	result, err := h.service.GetClientDiagnosis(ctx.Request.Context(), clientID, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get client diagnosis", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toClientDiagnosisResponse(*result), "Client diagnosis fetched successfully"))
}

// UpdateClientDiagnosis updates a client diagnosis
// @Summary Update a client diagnosis
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param diagnosis_id path uuid true "Diagnosis ID"
// @Param request body updateClientDiagnosisRequest true "Client diagnosis data"
// @Success 200 {object} httpapi.Envelope[clientDiagnosisResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/diagnoses/{diagnosis_id} [put]
func (h *ClientHandler) UpdateClientDiagnosis(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}
	diagnosisID, err := uuid.Parse(ctx.Param("diagnosis_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid diagnosis ID", ""))
		return
	}

	var req updateClientDiagnosisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	employeeID, _ := getEmployeeIDFromContext(ctx)

	result, err := h.service.UpdateClientDiagnosis(ctx.Request.Context(), clientID, diagnosisID, employeeID, toDomainUpdateClientDiagnosisParams(req, clientID, diagnosisID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update client diagnosis", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toClientDiagnosisResponse(*result), "Client diagnosis updated successfully"))
}

// DeleteClientDiagnosis deletes a client diagnosis
// @Summary Delete a client diagnosis
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param diagnosis_id path uuid true "Diagnosis ID"
// @Success 200 {object} httpapi.Envelope[deleteClientDiagnosisResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/diagnoses/{diagnosis_id} [delete]
func (h *ClientHandler) DeleteClientDiagnosis(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}
	diagnosisID, err := uuid.Parse(ctx.Param("diagnosis_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid diagnosis ID", ""))
		return
	}

	result, err := h.service.DeleteClientDiagnosis(ctx.Request.Context(), clientID, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete client diagnosis", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(deleteClientDiagnosisResponse{ID: result.ID}, "Client diagnosis deleted successfully"))
}

// =====================
// Medical - Medication Orders
// =====================

// CreateClientMedicationOrder creates a medication order
// @Summary Create a medication order
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body createClientMedicationOrderRequest true "Medication order data"
// @Success 201 {object} httpapi.Envelope[clientMedicationOrderResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/medication-orders [post]
func (h *ClientHandler) CreateClientMedicationOrder(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req createClientMedicationOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	employeeID, _ := getEmployeeIDFromContext(ctx)

	result, err := h.service.CreateClientMedicationOrder(ctx.Request.Context(), clientID, employeeID, toDomainCreateClientMedicationOrderParams(req, clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create medication order", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toClientMedicationOrderResponse(*result), "Medication order created successfully"))
}

// ListClientMedicationOrders lists medication orders
// @Summary List medication orders
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Param admin_mode query string false "Filter by admin mode"
// @Param diagnosis_id query uuid false "Filter by diagnosis"
// @Param search query string false "Search"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[clientMedicationOrderResponse]]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/medication-orders [get]
func (h *ClientHandler) ListClientMedicationOrders(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req listClientMedicationOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to bind query parameters", ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListClientMedicationOrders(ctx.Request.Context(), domain.ListClientMedicationOrdersParams{
		ClientID:    clientID,
		Status:      req.Status,
		AdminMode:   req.AdminMode,
		DiagnosisID: req.DiagnosisID,
		Search:      req.Search,
		Limit:       pageParams.Limit,
		Offset:      pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list medication orders", ""))
		return
	}

	items := make([]clientMedicationOrderResponse, 0, len(result.Items))
	for _, o := range result.Items {
		items = append(items, toClientMedicationOrderResponse(o))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Medication orders fetched successfully"))
}

// GetClientMedicationOrder gets a medication order
// @Summary Get a medication order
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param order_id path uuid true "Order ID"
// @Success 200 {object} httpapi.Envelope[clientMedicationOrderResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/medication-orders/{order_id} [get]
func (h *ClientHandler) GetClientMedicationOrder(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}
	orderID, err := uuid.Parse(ctx.Param("order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid order ID", ""))
		return
	}

	result, err := h.service.GetClientMedicationOrder(ctx.Request.Context(), clientID, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get medication order", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toClientMedicationOrderResponse(*result), "Medication order fetched successfully"))
}

// UpdateClientMedicationOrder updates a medication order
// @Summary Update a medication order
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param order_id path uuid true "Order ID"
// @Param request body updateClientMedicationOrderRequest true "Medication order data"
// @Success 200 {object} httpapi.Envelope[clientMedicationOrderResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/medication-orders/{order_id} [put]
func (h *ClientHandler) UpdateClientMedicationOrder(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}
	orderID, err := uuid.Parse(ctx.Param("order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid order ID", ""))
		return
	}

	var req updateClientMedicationOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	employeeID, _ := getEmployeeIDFromContext(ctx)

	result, err := h.service.UpdateClientMedicationOrder(ctx.Request.Context(), clientID, orderID, employeeID, toDomainUpdateClientMedicationOrderParams(req, clientID, orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update medication order", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toClientMedicationOrderResponse(*result), "Medication order updated successfully"))
}

// DeleteClientMedicationOrder deletes a medication order
// @Summary Delete a medication order
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param order_id path uuid true "Order ID"
// @Success 200 {object} httpapi.Envelope[deleteClientMedicationOrderResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/medical/medication-orders/{order_id} [delete]
func (h *ClientHandler) DeleteClientMedicationOrder(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}
	orderID, err := uuid.Parse(ctx.Param("order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid order ID", ""))
		return
	}

	result, err := h.service.DeleteClientMedicationOrder(ctx.Request.Context(), clientID, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrClientMedicationOrderNotFound) {
			ctx.JSON(http.StatusNotFound, httpapi.Fail("medication order not found", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete medication order", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(deleteClientMedicationOrderResponse{ID: result.ID}, "Medication order deleted successfully"))
}

// =====================
// Network - Sender
// =====================

// GetClientSender gets a client sender
// @Summary Get a client sender
// @Tags client_network
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} httpapi.Envelope[getClientSenderResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/sender [get]
func (h *ClientHandler) GetClientSender(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	sender, err := h.service.GetClientSender(ctx.Request.Context(), clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get client sender", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetClientSenderResponse(*sender), "Client sender fetched successfully"))
}

// =====================
// Network - Emergency Contacts
// =====================

// CreateClientEmergencyContact creates a client emergency contact
// @Summary Create a client emergency contact
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body createClientEmergencyContactRequest true "Emergency contact data"
// @Success 201 {object} httpapi.Envelope[clientEmergencyContactResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/emergency_contacts [post]
func (h *ClientHandler) CreateClientEmergencyContact(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req createClientEmergencyContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	result, err := h.service.CreateClientEmergencyContact(ctx.Request.Context(), clientID, toDomainCreateClientEmergencyContactParams(req, clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create client emergency contact", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toClientEmergencyContactResponse(*result), "Client emergency contact created successfully"))
}

// ListClientEmergencyContacts lists client emergency contacts
// @Summary List client emergency contacts
// @Tags client_network
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search query"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[clientEmergencyContactResponse]]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/emergency_contacts [get]
func (h *ClientHandler) ListClientEmergencyContacts(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req listClientEmergencyContactsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to bind query parameters", ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListClientEmergencyContacts(ctx.Request.Context(), domain.ListClientEmergencyContactsParams{
		ClientID: clientID,
		Search:   req.Search,
		Limit:    pageParams.Limit,
		Offset:   pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list client emergency contacts", ""))
		return
	}

	items := make([]clientEmergencyContactResponse, 0, len(result.Items))
	for _, c := range result.Items {
		items = append(items, toClientEmergencyContactResponse(c))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Client emergency contacts fetched successfully"))
}

// GetClientEmergencyContact gets a client emergency contact
// @Summary Get a client emergency contact
// @Tags client_network
// @Produce json
// @Param id path uuid true "Client ID"
// @Param contact_id path uuid true "Contact ID"
// @Success 200 {object} httpapi.Envelope[clientEmergencyContactResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/emergency_contacts/{contact_id} [get]
func (h *ClientHandler) GetClientEmergencyContact(ctx *gin.Context) {
	contactID, err := uuid.Parse(ctx.Param("contact_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid contact ID", ""))
		return
	}

	result, err := h.service.GetClientEmergencyContact(ctx.Request.Context(), contactID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get client emergency contact", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toClientEmergencyContactResponse(*result), "Client emergency contact fetched successfully"))
}

// UpdateClientEmergencyContact updates a client emergency contact
// @Summary Update a client emergency contact
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param contact_id path uuid true "Contact ID"
// @Param request body updateClientEmergencyContactRequest true "Emergency contact data"
// @Success 200 {object} httpapi.Envelope[clientEmergencyContactResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/emergency_contacts/{contact_id} [put]
func (h *ClientHandler) UpdateClientEmergencyContact(ctx *gin.Context) {
	contactID, err := uuid.Parse(ctx.Param("contact_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid contact ID", ""))
		return
	}

	var req updateClientEmergencyContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	result, err := h.service.UpdateClientEmergencyContact(ctx.Request.Context(), contactID, toDomainUpdateClientEmergencyContactParams(req, contactID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update client emergency contact", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toClientEmergencyContactResponse(*result), "Client emergency contact updated successfully"))
}

// DeleteClientEmergencyContact deletes a client emergency contact
// @Summary Delete a client emergency contact
// @Tags client_network
// @Produce json
// @Param id path uuid true "Client ID"
// @Param contact_id path uuid true "Contact ID"
// @Success 200 {object} httpapi.Envelope[deleteClientEmergencyContactResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/emergency_contacts/{contact_id} [delete]
func (h *ClientHandler) DeleteClientEmergencyContact(ctx *gin.Context) {
	contactID, err := uuid.Parse(ctx.Param("contact_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid contact ID", ""))
		return
	}

	result, err := h.service.DeleteClientEmergencyContact(ctx.Request.Context(), contactID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete client emergency contact", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(deleteClientEmergencyContactResponse{ID: result.ID}, "Client emergency contact deleted successfully"))
}

// =====================
// Network - Assigned Employees
// =====================

// CreateAssignedEmployee assigns an employee to a client
// @Summary Assign an employee to a client
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body createAssignedEmployeeRequest true "Employee assignment data"
// @Success 201 {object} httpapi.Envelope[assignedEmployeeResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/involved_employees [post]
func (h *ClientHandler) CreateAssignedEmployee(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req createAssignedEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	result, err := h.service.CreateAssignedEmployee(ctx.Request.Context(), clientID, toDomainCreateAssignedEmployeeParams(req, clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to assign employee", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toAssignedEmployeeResponse(*result), "Employee assigned successfully"))
}

// ListAssignedEmployees lists assigned employees
// @Summary List assigned employees
// @Tags client_network
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[assignedEmployeeResponse]]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/involved_employees [get]
func (h *ClientHandler) ListAssignedEmployees(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req listAssignedEmployeesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to bind query parameters", ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListAssignedEmployees(ctx.Request.Context(), domain.ListAssignedEmployeesParams{
		ClientID: clientID,
		Limit:    pageParams.Limit,
		Offset:   pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list assigned employees", ""))
		return
	}

	items := make([]assignedEmployeeResponse, 0, len(result.Items))
	for _, e := range result.Items {
		items = append(items, toAssignedEmployeeResponse(e))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Assigned employees fetched successfully"))
}

// GetAssignedEmployee gets an assigned employee
// @Summary Get an assigned employee
// @Tags client_network
// @Produce json
// @Param id path uuid true "Client ID"
// @Param assign_id path uuid true "Assignment ID"
// @Success 200 {object} httpapi.Envelope[assignedEmployeeResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/involved_employees/{assign_id} [get]
func (h *ClientHandler) GetAssignedEmployee(ctx *gin.Context) {
	assignID, err := uuid.Parse(ctx.Param("assign_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid assignment ID", ""))
		return
	}

	result, err := h.service.GetAssignedEmployee(ctx.Request.Context(), assignID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get assigned employee", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toAssignedEmployeeResponse(*result), "Assigned employee fetched successfully"))
}

// UpdateAssignedEmployee updates an assigned employee
// @Summary Update an assigned employee
// @Tags client_network
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param assign_id path uuid true "Assignment ID"
// @Param request body updateAssignedEmployeeRequest true "Assigned employee data"
// @Success 200 {object} httpapi.Envelope[assignedEmployeeResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/involved_employees/{assign_id} [put]
func (h *ClientHandler) UpdateAssignedEmployee(ctx *gin.Context) {
	assignID, err := uuid.Parse(ctx.Param("assign_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid assignment ID", ""))
		return
	}

	var req updateAssignedEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	result, err := h.service.UpdateAssignedEmployee(ctx.Request.Context(), assignID, toDomainUpdateAssignedEmployeeParams(req, assignID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update assigned employee", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toAssignedEmployeeResponse(*result), "Assigned employee updated successfully"))
}

// DeleteAssignedEmployee deletes an assigned employee
// @Summary Delete an assigned employee
// @Tags client_network
// @Produce json
// @Param id path uuid true "Client ID"
// @Param assign_id path uuid true "Assignment ID"
// @Success 200 {object} httpapi.Envelope[deleteAssignedEmployeeResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/involved_employees/{assign_id} [delete]
func (h *ClientHandler) DeleteAssignedEmployee(ctx *gin.Context) {
	assignID, err := uuid.Parse(ctx.Param("assign_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid assignment ID", ""))
		return
	}

	result, err := h.service.DeleteAssignedEmployee(ctx.Request.Context(), assignID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete assigned employee", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(deleteAssignedEmployeeResponse{ID: result.ID}, "Assigned employee deleted successfully"))
}

// =====================
// Network - Related Emails
// =====================

// GetClientRelatedEmails gets client related emails
// @Summary Get client related emails
// @Tags client_network
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} httpapi.Envelope[getClientRelatedEmailsResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/related_emails [get]
func (h *ClientHandler) GetClientRelatedEmails(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	result, err := h.service.GetClientRelatedEmails(ctx.Request.Context(), clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get client related emails", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(getClientRelatedEmailsResponse{Emails: result.Emails}, "Client related emails fetched successfully"))
}

// CreateProgressReport creates a new progress report for a client
// @Summary Create a new progress report for a client
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body createProgressReportRequest true "Progress Report Request"
// @Success 201 {object} httpapi.Envelope[progressReportResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/progress_reports [post]
func (h *ClientHandler) CreateProgressReport(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req createProgressReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	result, err := h.service.CreateProgressReport(ctx.Request.Context(), clientID, toDomainCreateProgressReportParams(req, clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create progress report", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toProgressReportResponse(*result), "Progress Report created successfully"))
}

// ListProgressReports lists all progress reports for a client
// @Summary List all progress reports for a client
// @Tags progress_reports
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param type query string false "Filter by report type"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[progressReportResponse]]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/progress_reports [get]
func (h *ClientHandler) ListProgressReports(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req listProgressReportsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to bind query parameters", ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListProgressReports(ctx.Request.Context(), domain.ListProgressReportsParams{
		ClientID: clientID,
		Type:     req.Type,
		Limit:    pageParams.Limit,
		Offset:   pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list progress reports", ""))
		return
	}

	items := make([]progressReportResponse, 0, len(result.Items))
	for _, r := range result.Items {
		items = append(items, toProgressReportResponse(r))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Progress reports retrieved successfully"))
}

// GetProgressReport retrieves a progress report for a client
// @Summary Retrieve a progress report for a client
// @Tags progress_reports
// @Produce json
// @Param report_id path uuid true "Progress Report ID"
// @Success 200 {object} httpapi.Envelope[progressReportResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/progress_reports/{report_id} [get]
func (h *ClientHandler) GetProgressReport(ctx *gin.Context) {
	reportID, err := uuid.Parse(ctx.Param("report_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid report ID", ""))
		return
	}

	result, err := h.service.GetProgressReport(ctx.Request.Context(), reportID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get progress report", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toProgressReportResponse(*result), "Progress Report retrieved successfully"))
}

// UpdateProgressReport updates a progress report for a client
// @Summary Update a progress report for a client
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param report_id path uuid true "Progress Report ID"
// @Param request body updateProgressReportRequest true "Progress Report Request"
// @Success 200 {object} httpapi.Envelope[progressReportResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/progress_reports/{report_id} [put]
func (h *ClientHandler) UpdateProgressReport(ctx *gin.Context) {
	reportID, err := uuid.Parse(ctx.Param("report_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid report ID", ""))
		return
	}

	var req updateProgressReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	result, err := h.service.UpdateProgressReport(ctx.Request.Context(), reportID, toDomainUpdateProgressReportParams(req, reportID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update progress report", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toProgressReportResponse(*result), "Progress Report updated successfully"))
}

// DeleteProgressReport deletes a progress report for a client
// @Summary Delete a progress report for a client
// @Tags progress_reports
// @Produce json
// @Param report_id path uuid true "Progress Report ID"
// @Success 200 {object} httpapi.Envelope[any]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/progress_reports/{report_id} [delete]
func (h *ClientHandler) DeleteProgressReport(ctx *gin.Context) {
	reportID, err := uuid.Parse(ctx.Param("report_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid report ID", ""))
		return
	}

	if err := h.service.DeleteProgressReport(ctx.Request.Context(), reportID); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete progress report", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK([]string{}, "Progress Report deleted successfully"))
}

// GenerateAutoReports generates auto reports
// @Summary Generate auto reports
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body generateAutoReportsRequest true "Request body"
// @Success 200 {object} httpapi.Envelope[generateAutoReportsResponse]
// @Router /clients/{id}/ai_progress_reports [post]
func (h *ClientHandler) GenerateAutoReports(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req generateAutoReportsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	report, err := h.service.GenerateAutoReports(ctx.Request.Context(), clientID, req.StartDate, req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to generate auto reports", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(generateAutoReportsResponse{Report: report}, "Auto reports generated successfully"))
}

// ConfirmAiProgressReport confirms an AI progress report
// @Summary Confirm a progress report for a client
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body confirmProgressReportRequest true "Progress Report Request"
// @Success 201 {object} httpapi.Envelope[aiGeneratedReportResponse]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/ai_progress_reports/confirm [post]
func (h *ClientHandler) ConfirmAiProgressReport(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req confirmProgressReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	result, err := h.service.ConfirmAiProgressReport(ctx.Request.Context(), clientID, req.ReportText, req.StartDate, req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to confirm AI progress report", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toAiGeneratedReportResponse(*result), "Progress Report created successfully"))
}

// ListAiGeneratedReports lists all AI generated reports for a client
// @Summary List all AI generated reports for a client
// @Tags progress_reports
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[aiGeneratedReportResponse]]
// @Failure 400,404 {object} httpapi.Envelope[any]
// @Router /clients/{id}/ai_progress_reports [get]
func (h *ClientHandler) ListAiGeneratedReports(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req listAiGeneratedReportsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to bind query parameters", ""))
		return
	}

	pageParams := req.Params()
	result, err := h.service.ListAiGeneratedReports(ctx.Request.Context(), domain.ListAiGeneratedReportsParams{
		ClientID: clientID,
		Limit:    pageParams.Limit,
		Offset:   pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list AI generated reports", ""))
		return
	}

	items := make([]aiGeneratedReportResponse, 0, len(result.Items))
	for _, r := range result.Items {
		items = append(items, toAiGeneratedReportResponse(r))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Progress reports retrieved successfully"))
}

// GetAppointmentCard retrieves an appointment card by client ID
func (h *ClientHandler) GetAppointmentCard(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	card, err := h.service.GetAppointmentCard(ctx.Request.Context(), clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get appointment card", ""))
		return
	}

	if card == nil {
		ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Appointment card not found"))
		return
	}

	resp := toGetAppointmentCardResponse(card)
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Appointment card retrieved successfully"))
}

// UpdateAppointmentCard creates or updates an appointment card for a client
func (h *ClientHandler) UpdateAppointmentCard(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req updateAppointmentCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	card, err := h.service.UpdateAppointmentCard(ctx.Request.Context(), clientID, toUpdateAppointmentCardParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update appointment card", ""))
		return
	}

	resp := toUpdateAppointmentCardResponse(card)
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Appointment card updated successfully"))
}

// GenerateAppointmentCardDocument generates a PDF document for the appointment card
func (h *ClientHandler) GenerateAppointmentCardDocument(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	pdfBytes, fileName, err := h.service.GenerateAppointmentCardDocument(ctx.Request.Context(), clientID)
	if err != nil {
		if err.Error() == "appointment card not found" {
			ctx.JSON(http.StatusNotFound, httpapi.Fail("appointment card not found", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to generate appointment card document", ""))
		return
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

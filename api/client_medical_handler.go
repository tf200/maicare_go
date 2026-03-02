package api

import (
	"net/http"

	clientp "maicare_go/service/client"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetClientMedicalOverviewApi returns diagnoses + active medication orders
// @Summary Get client medical overview
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} Response[clientp.ClientMedicalOverviewResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/overview [get]
func (server *Server) GetClientMedicalOverviewApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.GetClientMedicalOverview(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result, "Client medical overview fetched successfully"))
}

// CreateClientDiagnosisApi creates a client diagnosis
// @Summary Create a client diagnosis
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body clientp.CreateClientDiagnosisRequest true "Client diagnosis data"
// @Success 201 {object} Response[clientp.ClientDiagnosisResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/diagnoses [post]
func (server *Server) CreateClientDiagnosisApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.CreateClientDiagnosisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.CreateClientDiagnosis(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, SuccessResponse(result, "Client diagnosis created successfully"))
}

// ListClientDiagnosesApi lists client diagnoses
// @Summary List client diagnoses
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ClientDiagnosisResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/diagnoses [get]
func (server *Server) ListClientDiagnosesApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.ListClientDiagnosesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pag, err := server.businessService.ClientService.ListClientDiagnoses(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(pag, "Client diagnoses fetched successfully"))
}

// GetClientDiagnosisApi gets a client diagnosis
// @Summary Get a client diagnosis
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param diagnosis_id path uuid true "Diagnosis ID"
// @Success 200 {object} Response[clientp.ClientDiagnosisResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/diagnoses/{diagnosis_id} [get]
func (server *Server) GetClientDiagnosisApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	diagnosisID, err := uuid.Parse(ctx.Param("diagnosis_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.GetClientDiagnosis(ctx, clientID, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result, "Client diagnosis fetched successfully"))
}

// UpdateClientDiagnosisApi updates a client diagnosis
// @Summary Update a client diagnosis
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param diagnosis_id path uuid true "Diagnosis ID"
// @Param request body clientp.UpdateClientDiagnosisRequest true "Client diagnosis data"
// @Success 200 {object} Response[clientp.ClientDiagnosisResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/diagnoses/{diagnosis_id} [put]
func (server *Server) UpdateClientDiagnosisApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	diagnosisID, err := uuid.Parse(ctx.Param("diagnosis_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateClientDiagnosisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.UpdateClientDiagnosis(ctx, req, clientID, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result, "Client diagnosis updated successfully"))
}

// DeleteClientDiagnosisApi deletes a client diagnosis
// @Summary Delete a client diagnosis
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param diagnosis_id path uuid true "Diagnosis ID"
// @Success 200 {object} Response[clientp.DeleteClientDiagnosisResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/diagnoses/{diagnosis_id} [delete]
func (server *Server) DeleteClientDiagnosisApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	diagnosisID, err := uuid.Parse(ctx.Param("diagnosis_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.DeleteClientDiagnosis(ctx, clientID, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result, "Client diagnosis deleted successfully"))
}

// CreateClientMedicationOrderApi creates a medication order
// @Summary Create a medication order
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body clientp.CreateClientMedicationOrderRequest true "Medication order data"
// @Success 201 {object} Response[clientp.ClientMedicationOrderResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/medication-orders [post]
func (server *Server) CreateClientMedicationOrderApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.CreateClientMedicationOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.CreateClientMedicationOrder(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, SuccessResponse(result, "Medication order created successfully"))
}

// ListClientMedicationOrdersApi lists medication orders
// @Summary List medication orders
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Param admin_mode query string false "Filter by admin mode"
// @Param diagnosis_id query uuid false "Filter by diagnosis"
// @Param search query string false "Search"
// @Success 200 {object} Response[pagination.Response[clientp.ClientMedicationOrderResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/medication-orders [get]
func (server *Server) ListClientMedicationOrdersApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.ListClientMedicationOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pag, err := server.businessService.ClientService.ListClientMedicationOrders(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(pag, "Medication orders fetched successfully"))
}

// GetClientMedicationOrderApi gets a medication order
// @Summary Get a medication order
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param order_id path uuid true "Order ID"
// @Success 200 {object} Response[clientp.ClientMedicationOrderResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/medication-orders/{order_id} [get]
func (server *Server) GetClientMedicationOrderApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	orderID, err := uuid.Parse(ctx.Param("order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.GetClientMedicationOrder(ctx, clientID, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result, "Medication order fetched successfully"))
}

// UpdateClientMedicationOrderApi updates a medication order
// @Summary Update a medication order
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param order_id path uuid true "Order ID"
// @Param request body clientp.UpdateClientMedicationOrderRequest true "Medication order data"
// @Success 200 {object} Response[clientp.ClientMedicationOrderResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/medication-orders/{order_id} [put]
func (server *Server) UpdateClientMedicationOrderApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	orderID, err := uuid.Parse(ctx.Param("order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateClientMedicationOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.UpdateClientMedicationOrder(ctx, req, clientID, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result, "Medication order updated successfully"))
}

// DeleteClientMedicationOrderApi deletes a medication order
// @Summary Delete a medication order
// @Tags client_Medical
// @Produce json
// @Param id path uuid true "Client ID"
// @Param order_id path uuid true "Order ID"
// @Success 200 {object} Response[clientp.DeleteClientMedicationOrderResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medical/medication-orders/{order_id} [delete]
func (server *Server) DeleteClientMedicationOrderApi(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	orderID, err := uuid.Parse(ctx.Param("order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.DeleteClientMedicationOrder(ctx, clientID, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(result, "Medication order deleted successfully"))
}

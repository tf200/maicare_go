package handler

import (
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterContractRoutes(
	rg *gin.RouterGroup,
	handler *ContractHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	// Contract type routes
	contractTypesGroup := rg.Group("/contract_types")
	{
		contractTypesGroup.POST("", auth, requirePermission("CONTRACT_TYPE.CREATE"), handler.CreateContractType)
		contractTypesGroup.GET("", auth, requirePermission("CONTRACT_TYPE.VIEW"), handler.ListContractTypes)
		contractTypesGroup.DELETE("/:id", auth, requirePermission("CONTRACT_TYPE.DELETE"), handler.DeleteContractType)
	}

	// Client-scoped contract routes
	clientsGroup := rg.Group("/clients")
	{
		clientsGroup.GET("/:id/contracts", auth, requirePermission("CONTRACT.VIEW"), handler.ListClientContracts)
	}

	// Top-level contract routes
	contractsGroup := rg.Group("/contracts")
	{
		contractsGroup.POST("", auth, requirePermission("CONTRACT.CREATE"), handler.CreateContract)
		contractsGroup.GET("", auth, requirePermission("CONTRACT.VIEW"), handler.ListContracts)
		contractsGroup.GET("/:id", auth, requirePermission("CONTRACT.VIEW"), handler.GetContract)
		contractsGroup.PUT("/:id", auth, requirePermission("CONTRACT.UPDATE"), handler.UpdateContract)
		contractsGroup.PUT("/:id/status", auth, requirePermission("CONTRACT.UPDATE"), handler.UpdateContractStatus)
		contractsGroup.GET("/:id/audit", auth, requirePermission("CONTRACT.VIEW"), handler.GetContractAuditLog)
	}
}

type ContractHandler struct {
	service domain.ContractService
}

func NewContractHandler(service domain.ContractService) *ContractHandler {
	return &ContractHandler{service: service}
}

// CreateContractType creates a new contract type.
// @Summary Create a new contract type
// @Tags contracts
// @Accept json
// @Produce json
// @Param request body createContractTypeRequest true "Create Contract Type Request"
// @Success 200 {object} httpapi.Envelope[contractTypeResponse]
// @Router /contract_types [post]
func (h *ContractHandler) CreateContractType(ctx *gin.Context) {
	var req createContractTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.CreateContractType(ctx.Request.Context(), domain.CreateContractTypeParams{Name: req.Name})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create contract type", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toCreateContractTypeResult(*result), "Contract type created successfully"))
}

// ListContractTypes returns a list of contract types.
// @Summary List contract types
// @Tags contracts
// @Produce json
// @Success 200 {object} httpapi.Envelope[[]contractTypeResponse]
// @Router /contract_types [get]
func (h *ContractHandler) ListContractTypes(ctx *gin.Context) {
	types, err := h.service.ListContractTypes(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list contract types", err.Error()))
		return
	}

	res := make([]contractTypeResponse, len(types))
	for i, t := range types {
		res[i] = toCreateContractTypeResult(t)
	}

	ctx.JSON(http.StatusOK, httpapi.OK(res, "Contract types retrieved successfully"))
}

// DeleteContractType deletes a contract type.
// @Summary Delete a contract type
// @Tags contracts
// @Produce json
// @Param id path string true "Contract Type ID"
// @Success 200 {object} httpapi.Envelope[contractTypeResponse]
// @Router /contract_types/{id} [delete]
func (h *ContractHandler) DeleteContractType(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid contract type ID", err.Error()))
		return
	}

	result, err := h.service.DeleteContractType(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete contract type", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toCreateContractTypeResult(*result), "Contract type deleted successfully"))
}

// CreateContract creates a new contract.
// @Summary Create a new contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param request body createContractRequest true "Create Contract Request"
// @Success 200 {object} httpapi.Envelope[contractResponse]
// @Router /contracts [post]
func (h *ContractHandler) CreateContract(ctx *gin.Context) {
	var req createContractRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.CreateContract(ctx.Request.Context(), toCreateContractParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create contract", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toContractResponse(*result), "Contract created successfully"))
}

// ListClientContracts returns a list of contracts for a client.
// @Summary List contracts for a client
// @Tags contracts
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[[]clientContractSummaryResponse]]
// @Router /clients/{id}/contracts [get]
func (h *ContractHandler) ListClientContracts(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", err.Error()))
		return
	}

	var req httpapi.PageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid pagination params", err.Error()))
		return
	}

	result, err := h.service.ListClientContracts(ctx.Request.Context(), domain.ListClientContractsParams{
		ClientID: clientID,
		Limit:    req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list client contracts", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(httpapi.NewPageResponse(ctx, req, toContractClientListItemResponses(result.Contracts), result.TotalCount), "Contracts retrieved successfully"))
}

// GetContract returns a contract by ID.
// @Summary Get a contract by ID
// @Tags contracts
// @Produce json
// @Param id path string true "Contract ID"
// @Success 200 {object} httpapi.Envelope[contractDetailResponse]
// @Router /contracts/{id} [get]
func (h *ContractHandler) GetContract(ctx *gin.Context) {
	contractID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid contract ID", err.Error()))
		return
	}

	contract, err := h.service.GetContractByID(ctx.Request.Context(), contractID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get contract", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toContractDetailResponse(contract), "Contract retrieved successfully"))
}

// UpdateContract updates a contract.
// @Summary Update a contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Param request body updateContractRequest true "Update Contract Request"
// @Success 200 {object} httpapi.Envelope[updateContractResponse]
// @Router /contracts/{id} [put]
func (h *ContractHandler) UpdateContract(ctx *gin.Context) {
	contractID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid contract ID", err.Error()))
		return
	}

	employeeID, err := getEmployeeIDFromContextForContracts(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid employee ID", err.Error()))
		return
	}

	var req updateContractRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.UpdateContract(ctx.Request.Context(), toUpdateContractParams(contractID, req), employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update contract", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toUpdateContractResponse(*result), "Contract updated successfully"))
}

// UpdateContractStatus updates the status of a contract.
// @Summary Update the status of a contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Param request body updateContractStatusRequest true "Update Contract Status Request"
// @Success 200 {object} httpapi.Envelope[updateContractStatusResponse]
// @Router /contracts/{id}/status [put]
func (h *ContractHandler) UpdateContractStatus(ctx *gin.Context) {
	contractID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid contract ID", err.Error()))
		return
	}

	employeeID, err := getEmployeeIDFromContextForContracts(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid employee ID", err.Error()))
		return
	}

	var req updateContractStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.UpdateContractStatus(ctx.Request.Context(), domain.UpdateContractStatusParams{
		ContractID: contractID,
		Status:     req.Status,
	}, employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update contract status", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toUpdateContractStatusResponse(*result), "Contract status updated successfully"))
}

// ListContracts returns a list of contracts.
// @Summary List contracts
// @Tags contracts
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search by client or sender name"
// @Param status query []string false "Status"
// @Param care_type query []string false "Care type"
// @Param financing_act query []string false "Financing act"
// @Param financing_option query string false "Financing option"
// @Param end_date_from query string false "Filter contracts ending on/after date (YYYY-MM-DD)"
// @Param end_date_to query string false "Filter contracts ending on/before date (YYYY-MM-DD)"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[[]contractListItemResponse]]
// @Router /contracts [get]
func (h *ContractHandler) ListContracts(ctx *gin.Context) {
	var req listContractsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid query params", err.Error()))
		return
	}

	result, err := h.service.ListContracts(ctx.Request.Context(), toListContractsParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list contracts", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(httpapi.NewPageResponse(ctx, req.PageRequest, toContractListItemResponses(result.Contracts), result.TotalCount), "Contracts retrieved successfully"))
}

// GetContractAuditLog returns the audit logs for a contract.
// @Summary Get audit logs for a contract
// @Tags contracts
// @Produce json
// @Param id path string true "Contract ID"
// @Success 200 {object} httpapi.Envelope[[]contractAuditLogResponse]
// @Router /contracts/{id}/audit [get]
func (h *ContractHandler) GetContractAuditLog(ctx *gin.Context) {
	contractID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid contract ID", err.Error()))
		return
	}

	logs, err := h.service.GetContractAuditLog(ctx.Request.Context(), contractID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get contract audit logs", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toContractAuditLogResponses(logs), "Audit logs retrieved successfully"))
}

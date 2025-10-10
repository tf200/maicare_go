package api

import (
	"maicare_go/service/contract"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateContractTypeApi creates a new contract type
// @Summary Create a new contract type
// @Tags contracts
// @Accept json
// @Produce json
// @Param request body contractp.CreateContractTypeRequest true "Create Contract Type Request"
// @Success 200 {object} Response[contractp.CreateContractTypeResponse]
// @Router /contract_types [post]
func (server *Server) CreateContractTypeApi(ctx *gin.Context) {
	var req contract.CreateContractTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contractType, err := server.businessService.ContractService.CreateContractType(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(contractType, "Contract type created successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListContractTypesApi returns a list of contract types
// @Summary List contract types
// @Tags contracts
// @Accept json
// @Produce json
// @Success 200 {array} Response[[]contractp.ListContractTypesResponse]
// @Router /contract_types [get]
func (server *Server) ListContractTypesApi(ctx *gin.Context) {
	contractTypes, err := server.businessService.ContractService.ListContractTypes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(contractTypes, "Contract types retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteContractTypeApi deletes a contract type
// @Summary Delete a contract type
// @Tags contracts
// @Produce json
// @Param id path string true "Contract Type ID"
// @Success 200 {object} Response[contractp.DeleteContractTypeResponse]
// @Router /contract_types/{id} [delete]
func (server *Server) DeleteContractTypeApi(ctx *gin.Context) {
	id := ctx.Param("id")
	contractTypeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ContractService.DeleteContractType(ctx, contractTypeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Contract type deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateContractApi creates a new contract
// @Summary Create a new contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param request body contractp.CreateContractRequest true "Create Contract Request"
// @Success 200 {object} Response[contractp.CreateContractResponse]
// @Router /clients/{id}/contracts [post]
func (server *Server) CreateContractApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req contract.CreateContractRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contract, err := server.businessService.ContractService.CreateContract(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(contract, "Contract created successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListClientContractsApi returns a list of contracts for a client
// @Summary List contracts for a client
// @Tags contracts
// @Produce json
// @Param id path string true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[[]contractp.ListClientContractsResponse]]
// @Router /clients/{id}/contracts [get]
func (server *Server) ListClientContractsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req contract.ListClientContractsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contracts, err := server.businessService.ContractService.ListClientContracts(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(contracts, "Contracts retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateContractApi updates a contract
// @Summary Update a contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Param request body contractp.UpdateContractRequest true "Update Contract Request"
// @Success 200 {object} Response[contractp.UpdateContractResponse]
// @Router /contracts/{id} [put]
func (server *Server) UpdateContractApi(ctx *gin.Context) {
	id := ctx.Param("id")
	contractID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req contract.UpdateContractRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contract, err := server.businessService.ContractService.UpdateContract(ctx, req, contractID, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(contract, "Contract updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateContractStatusApi updates the status of a contract
// @Summary Update the status of a contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Param request body contractp.UpdateContractStatusRequest true "Update Contract Status Request"
// @Success 200 {object} Response[contractp.UpdateContractStatusResponse]
// @Router /contracts/{id}/status [put]
func (server *Server) UpdateContractStatusApi(ctx *gin.Context) {
	id := ctx.Param("id")
	contractID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req contract.UpdateContractStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedContract, err := server.businessService.ContractService.UpdateContractStatus(ctx, req, contractID, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(updatedContract, "Contract status updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientContractApi returns a contract by ID
// @Summary Get a contract by ID
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Success 200 {object} Response[contractp.GetClientContractResponse]
// @Router /clients/{id}/contracts/{contract_id} [get]
func (server *Server) GetClientContractApi(ctx *gin.Context) {
	id := ctx.Param("contract_id")
	contractID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contract, err := server.businessService.ContractService.GetClientContract(ctx, contractID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(contract, "Contract retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListContractsApi returns a list of contracts
// @Summary List contracts
// @Tags contracts
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search query"
// @Param status query []string false "Status" Enums(approved, draft, terminated, stopped)
// @Param care_type query []string false "Care type" Enums(ambulante, accommodation)
// @Param financing_act query []string false "Financing act" Enums(WMO, ZVW, WLZ, JW, WPG)
// @Param financing_option query []string false "Financing option" Enums(ZIN, PGB)
// @Success 200 {object} Response[pagination.Response[[]contractp.ListContractsResponse]]
// @Router /contracts [get]
func (server *Server) ListContractsApi(ctx *gin.Context) {
	var req contract.ListContractsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contracts, err := server.businessService.ContractService.ListContracts(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(contracts, "Contracts retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetContractAuditLogApi returns the audit logs for a contract
// @Summary Get audit logs for a contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Success 200 {array} Response[[]contractp.GetContractAuditLogResponse]
// @Router /contracts/{id}/audit [get]
func (server *Server) GetContractAuditLogApi(ctx *gin.Context) {
	id := ctx.Param("id")
	contractID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	auditLogs, err := server.businessService.ContractService.GetContractAuditLog(ctx, contractID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(auditLogs, "Audit logs retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

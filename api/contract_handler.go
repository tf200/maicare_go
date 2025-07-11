package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateContractTypeRequest defines the request for CreateContractType handler
type CreateContractTypeRequest struct {
	Name string `json:"name"`
}

// CreateContractTypeResponse defines the response for CreateContractType handler
type CreateContractTypeResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// CreateContractTypeApi creates a new contract type
// @Summary Create a new contract type
// @Tags contracts
// @Accept json
// @Produce json
// @Param request body CreateContractTypeRequest true "Create Contract Type Request"
// @Success 200 {object} CreateContractTypeResponse
// @Router /contract_types [post]
func (server *Server) CreateContractTypeApi(ctx *gin.Context) {
	var req CreateContractTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contractType, err := server.store.CreateContractType(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateContractTypeResponse{
		ID:   contractType.ID,
		Name: contractType.Name,
	}, "Contract type created successfully")
	ctx.JSON(http.StatusOK, res)

}

// ListContractTypesResponse defines the response for ListContractTypes handler
type ListContractTypesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ListContractTypesApi returns a list of contract types
// @Summary List contract types
// @Tags contracts
// @Accept json
// @Produce json
// @Success 200 {array} ListContractTypesResponse
// @Router /contract_types [get]
func (server *Server) ListContractTypesApi(ctx *gin.Context) {
	contractTypes, err := server.store.ListContractTypes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	contractTypesRes := make([]ListContractTypesResponse, len(contractTypes))

	if len(contractTypes) == 0 {
		res := SuccessResponse(contractTypesRes, "No contract types found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	for i, contractType := range contractTypes {
		contractTypesRes[i] = ListContractTypesResponse{
			ID:   contractType.ID,
			Name: contractType.Name,
		}
	}

	res := SuccessResponse(contractTypesRes, "Contract types retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

type DeleteContractTypeResponse struct {
	ID int64 `json:"id"`
}

// DeleteContractTypeApi deletes a contract type
// @Summary Delete a contract type
// @Tags contracts
// @Produce json
// @Param id path string true "Contract Type ID"
// @Success 200 {object} DeleteContractTypeResponse
// @Router /contract_types/{id} [delete]
func (server *Server) DeleteContractTypeApi(ctx *gin.Context) {
	id := ctx.Param("id")
	contractTypeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteContractType(ctx, contractTypeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(DeleteContractTypeResponse{ID: contractTypeID}, "Contract type deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateContractRequest defines the request for CreateContract handler
type CreateContractRequest struct {
	TypeID          *int64      `json:"type_id" example:"1"`
	StartDate       time.Time   `json:"start_date" example:"2023-01-01T00:00:00Z"`
	EndDate         time.Time   `json:"end_date" example:"2023-12-31T00:00:00Z"`
	ReminderPeriod  int32       `json:"reminder_period" example:"30"`
	Vat             *int32      `json:"VAT" example:"21"`
	Price           float64     `json:"price" example:"100.50"`
	PriceTimeUnit   string      `json:"price_time_unit" binding:"required,oneof=minute hourly daily weekly monthly yearly" example:"monthly" enum:"minute,hourly,daily,weekly,monthly,yearly"`
	Hours           *float64    `json:"hours" example:"40"`
	HoursType       *string     `json:"hours_type" enum:"weekly,all_period"`
	CareName        string      `json:"care_name" example:"Home Care"`
	CareType        string      `json:"care_type" binding:"required,oneof=ambulante accommodation" example:"ambulante" enum:"ambulante,accommodation"`
	SenderID        *int64      `json:"sender_id" example:"2"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act" binding:"required,oneof=WMO ZVW WLZ JW WPG" example:"WMO" enum:"WMO,ZVW,WLZ,JW,WPG"`
	FinancingOption string      `json:"financing_option" binding:"required,oneof=ZIN PGB" example:"ZIN" enum:"ZIN,PGB"`
}

// CreateContractResponse defines the response for CreateContract handler
type CreateContractResponse struct {
	ID              int64              `json:"id"`
	TypeID          *int64             `json:"type_id"`
	Status          string             `json:"status"`
	StartDate       time.Time          `json:"start_date"`
	EndDate         time.Time          `json:"end_date"`
	ReminderPeriod  int32              `json:"reminder_period"`
	Vat             *int32             `json:"VAT"`
	Price           float64            `json:"price"`
	PriceTimeUnit   string             `json:"price_time_unit"`
	Hours           *float64           `json:"hours"`
	HoursType       *string            `json:"hours_type"`
	CareName        string             `json:"care_name"`
	CareType        string             `json:"care_type"`
	ClientID        int64              `json:"client_id"`
	SenderID        *int64             `json:"sender_id"`
	AttachmentIds   []uuid.UUID        `json:"attachment_ids"`
	FinancingAct    string             `json:"financing_act"`
	FinancingOption string             `json:"financing_option"`
	DepartureReason *string            `json:"departure_reason"`
	DepartureReport *string            `json:"departure_report"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
}

// CreateContractApi creates a new contract
// @Summary Create a new contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param request body CreateContractRequest true "Create Contract Request"
// @Success 200 {object} CreateContractResponse
// @Router /clients/{id}/contracts [post]
func (server *Server) CreateContractApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateContractRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contract, err := server.store.CreateContract(ctx, db.CreateContractParams{
		TypeID:          req.TypeID,
		StartDate:       pgtype.Timestamptz{Time: req.StartDate, Valid: true},
		EndDate:         pgtype.Timestamptz{Time: req.EndDate, Valid: true},
		ReminderPeriod:  req.ReminderPeriod,
		Vat:             req.Vat,
		Price:           req.Price,
		PriceTimeUnit:   req.PriceTimeUnit,
		Hours:           req.Hours,
		HoursType:       req.HoursType,
		CareName:        req.CareName,
		CareType:        req.CareType,
		ClientID:        clientID,
		SenderID:        req.SenderID,
		AttachmentIds:   req.AttachmentIds,
		FinancingAct:    req.FinancingAct,
		FinancingOption: req.FinancingOption,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateContractResponse{
		ID:              contract.ID,
		TypeID:          contract.TypeID,
		Status:          contract.Status,
		StartDate:       contract.StartDate.Time,
		EndDate:         contract.EndDate.Time,
		ReminderPeriod:  contract.ReminderPeriod,
		Vat:             contract.Vat,
		Price:           contract.Price,
		PriceTimeUnit:   contract.PriceTimeUnit,
		Hours:           contract.Hours,
		HoursType:       contract.HoursType,
		CareName:        contract.CareName,
		CareType:        contract.CareType,
		ClientID:        contract.ClientID,
		SenderID:        contract.SenderID,
		AttachmentIds:   contract.AttachmentIds,
		FinancingAct:    contract.FinancingAct,
		FinancingOption: contract.FinancingOption,
		DepartureReason: contract.DepartureReason,
		DepartureReport: contract.DepartureReport,
		UpdatedAt:       contract.UpdatedAt,
		CreatedAt:       contract.CreatedAt,
	}, "Contract created successfully")
	ctx.JSON(http.StatusOK, res)

}

// ListClientContractsRequest defines the request for ListClientContracts handler
type ListClientContractsRequest struct {
	pagination.Request
}

// ListClientContractsResponse defines the response for ListClientContracts handler
type ListClientContractsResponse struct {
	ID              int64       `json:"id"`
	TypeID          *int64      `json:"type_id"`
	Status          string      `json:"status"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  int32       `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           float64     `json:"price"`
	PriceTimeUnit   string      `json:"price_time_unit"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type"`
	ClientID        int64       `json:"client_id"`
	SenderID        *int64      `json:"sender_id"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
	DepartureReason *string     `json:"departure_reason"`
	DepartureReport *string     `json:"departure_report"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`
}

// ListClientContractsApi returns a list of contracts for a client
// @Summary List contracts for a client
// @Tags contracts
// @Produce json
// @Param id path string true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} pagination.Response[ListClientContractsResponse]
// @Router /clients/{id}/contracts [get]
func (server *Server) ListClientContractsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListClientContractsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	contracts, err := server.store.ListClientContracts(ctx, db.ListClientContractsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(contracts) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListClientContractsResponse{}, 0)
		res := SuccessResponse(pag, "No contracts found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	totalCount := contracts[0].TotalCount

	contractsRes := make([]ListClientContractsResponse, len(contracts))
	for i, contract := range contracts {
		contractsRes[i] = ListClientContractsResponse{
			ID:              contract.ID,
			TypeID:          contract.TypeID,
			Status:          contract.Status,
			StartDate:       contract.StartDate.Time,
			EndDate:         contract.EndDate.Time,
			ReminderPeriod:  contract.ReminderPeriod,
			Vat:             contract.Vat,
			Price:           contract.Price,
			PriceTimeUnit:   contract.PriceTimeUnit,
			Hours:           contract.Hours,
			HoursType:       contract.HoursType,
			CareName:        contract.CareName,
			CareType:        contract.CareType,
			ClientID:        contract.ClientID,
			SenderID:        contract.SenderID,
			AttachmentIds:   contract.AttachmentIds,
			FinancingAct:    contract.FinancingAct,
			FinancingOption: contract.FinancingOption,
			DepartureReason: contract.DepartureReason,
			DepartureReport: contract.DepartureReport,
			UpdatedAt:       contract.UpdatedAt.Time,
			CreatedAt:       contract.CreatedAt.Time,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, contractsRes, totalCount)
	res := SuccessResponse(pag, "Contracts retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateContractRequest defines the request for UpdateContract handler
type UpdateContractRequest struct {
	TypeID          *int64      `json:"type_id"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  *int32      `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           *float64    `json:"price"`
	PriceTimeUnit   *string     `json:"price_time_unit"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        *string     `json:"care_name"`
	CareType        *string     `json:"care_type"`
	SenderID        *int64      `json:"sender_id"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    *string     `json:"financing_act"`
	FinancingOption *string     `json:"financing_option"`
	Status          *string     `json:"status"`
}

// UpdateContractResponse defines the response for UpdateContract handler
type UpdateContractResponse struct {
	ID              int64       `json:"id"`
	TypeID          *int64      `json:"type_id"`
	Status          string      `json:"status"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  int32       `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           float64     `json:"price"`
	PriceFrequency  string      `json:"price_frequency"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type"`
	ClientID        int64       `json:"client_id"`
	SenderID        *int64      `json:"sender_id"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
	DepartureReason *string     `json:"departure_reason"`
	DepartureReport *string     `json:"departure_report"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`
}

// UpdateContractApi updates a contract
// @Summary Update a contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Param request body UpdateContractRequest true "Update Contract Request"
// @Success 200 {object} UpdateContractResponse
// @Router /contracts/{id} [put]
func (server *Server) UpdateContractApi(ctx *gin.Context) {
	id := ctx.Param("id")
	contractID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateContractRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contract, err := server.store.UpdateContract(ctx, db.UpdateContractParams{
		ID:              contractID,
		TypeID:          req.TypeID,
		StartDate:       pgtype.Timestamptz{Time: req.StartDate, Valid: true},
		EndDate:         pgtype.Timestamptz{Time: req.EndDate, Valid: true},
		ReminderPeriod:  req.ReminderPeriod,
		VAT:             req.Vat,
		Price:           req.Price,
		PriceTimeUnit:   req.PriceTimeUnit,
		Hours:           req.Hours,
		HoursType:       req.HoursType,
		CareName:        req.CareName,
		CareType:        req.CareType,
		SenderID:        req.SenderID,
		AttachmentIds:   req.AttachmentIds,
		FinancingAct:    req.FinancingAct,
		FinancingOption: req.FinancingOption,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateContractResponse{
		ID:              contract.ID,
		TypeID:          contract.TypeID,
		Status:          contract.Status,
		StartDate:       contract.StartDate.Time,
		EndDate:         contract.EndDate.Time,
		ReminderPeriod:  contract.ReminderPeriod,
		Vat:             contract.Vat,
		Price:           contract.Price,
		PriceFrequency:  contract.PriceTimeUnit,
		Hours:           contract.Hours,
		HoursType:       contract.HoursType,
		CareName:        contract.CareName,
		CareType:        contract.CareType,
		ClientID:        contract.ClientID,
		SenderID:        contract.SenderID,
		AttachmentIds:   contract.AttachmentIds,
		FinancingAct:    contract.FinancingAct,
		FinancingOption: contract.FinancingOption,
		DepartureReason: contract.DepartureReason,
		DepartureReport: contract.DepartureReport,
		UpdatedAt:       contract.UpdatedAt.Time,
		CreatedAt:       contract.CreatedAt.Time,
	}, "Contract updated successfully")
	ctx.JSON(http.StatusOK, res)
}

type UpdateContractStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=approved draft terminated stopped expired" example:"approved" enum:"approved,draft,terminated,stopped,expired"`
}

// UpdateContractStatusResponse defines the response for UpdateContractStatus handler
type UpdateContractStatusResponse struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

// UpdateContractStatusApi updates the status of a contract
// @Summary Update the status of a contract
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Param request body UpdateContractStatusRequest true "Update Contract Status Request"
// @Success 200 {object} UpdateContractStatusResponse
// @Router /contracts/{id}/status [put]
func (server *Server) UpdateContractStatusApi(ctx *gin.Context) {
	id := ctx.Param("id")
	contractID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req UpdateContractStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Fetch the current contract to validate status change
	contract, err := server.store.GetClientContract(ctx, contractID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.Status == "approved" && contract.EndDate.Time.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("cannot approve that has alredy ended")))
		return
	}

	// Validate status change

	updatedContract, err := server.store.UpdateContractStatus(ctx, db.UpdateContractStatusParams{
		ID:     contractID,
		Status: req.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateContractStatusResponse{
		ID:     updatedContract.ID,
		Status: updatedContract.Status,
	}, "Contract status updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientContractResponse defines the response for GetContract handler
type GetClientContractResponse struct {
	ID              int64       `json:"id"`
	TypeID          *int64      `json:"type_id"`
	TypeName        string      `json:"type_name"`
	Status          string      `json:"status"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  int32       `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           float64     `json:"price"`
	PriceTimeUnit   string      `json:"price_time_unit"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type"`
	ClientID        int64       `json:"client_id"`
	ClientFirstName string      `json:"client_first_name"`
	ClientLastName  string      `json:"client_last_name"`
	SenderID        *int64      `json:"sender_id"`
	SenderName      *string     `json:"sender_name"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
	DepartureReason *string     `json:"departure_reason"`
	DepartureReport *string     `json:"departure_report"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`
}

// GetClientContractApi returns a contract by ID
// @Summary Get a contract by ID
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Success 200 {object} GetClientContractResponse
// @Router /clients/{id}/contracts/{contract_id} [get]
func (server *Server) GetClientContractApi(ctx *gin.Context) {
	id := ctx.Param("contract_id")
	contractID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contract, err := server.store.GetClientContract(ctx, contractID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetClientContractResponse{
		ID:              contract.ID,
		TypeID:          contract.TypeID,
		TypeName:        contract.ContractTypeName,
		Status:          contract.Status,
		StartDate:       contract.StartDate.Time,
		EndDate:         contract.EndDate.Time,
		ReminderPeriod:  contract.ReminderPeriod,
		Vat:             contract.Vat,
		Price:           contract.Price,
		PriceTimeUnit:   contract.PriceTimeUnit,
		Hours:           contract.Hours,
		HoursType:       contract.HoursType,
		CareName:        contract.CareName,
		CareType:        contract.CareType,
		ClientID:        contract.ClientID,
		ClientFirstName: contract.ClientFirstName,
		ClientLastName:  contract.ClientLastName,
		SenderID:        contract.SenderID,
		SenderName:      contract.SenderName,
		AttachmentIds:   contract.AttachmentIds,
		FinancingAct:    contract.FinancingAct,
		FinancingOption: contract.FinancingOption,
		DepartureReason: contract.DepartureReason,
		DepartureReport: contract.DepartureReport,
		UpdatedAt:       contract.UpdatedAt.Time,
		CreatedAt:       contract.CreatedAt.Time,
	}, "Contract retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListContractsRequest defines the request for ListContracts handler
type ListContractsRequest struct {
	pagination.Request
	Search          *string  `form:"search" binding:"omitempty"`
	Status          []string `form:"status" binding:"omitempty,dive,oneof=approved draft terminated stopped"`
	CareType        []string `form:"care_type" binding:"omitempty,dive,oneof=ambulante accommodation"`
	FinancingAct    []string `form:"financing_act" binding:"omitempty,dive,oneof=WMO ZVW WLZ JW WPG"`
	FinancingOption []string `form:"financing_option" binding:"omitempty,dive,oneof=ZIN PGB"`
}

// ListContractsResponse defines the response for ListContracts handler
type ListContractsResponse struct {
	ID              int64     `json:"id"`
	ClientID        int64     `json:"client_id"`
	Status          string    `json:"status"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Price           float64   `json:"price"`
	PriceTimeUnit   string    `json:"price_time_unit"`
	CareName        string    `json:"care_name"`
	CareType        string    `json:"care_type"`
	FinancingAct    string    `json:"financing_act"`
	FinancingOption string    `json:"financing_option"`
	CreatedAt       time.Time `json:"created_at"`
	SenderID        *int64    `json:"sender_id"`
	SenderName      *string   `json:"sender_name"`
	ClientFirstName string    `json:"client_first_name"`
	ClientLastName  string    `json:"client_last_name"`
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
// @Success 200 {object} pagination.Response[ListContractsResponse]
// @Router /contracts [get]
func (server *Server) ListContractsApi(ctx *gin.Context) {
	var req ListContractsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	contracts, err := server.store.ListContracts(ctx, db.ListContractsParams{
		Limit:           params.Limit,
		Offset:          params.Offset,
		Search:          req.Search,
		Status:          req.Status,
		CareType:        req.CareType,
		FinancingAct:    req.FinancingAct,
		FinancingOption: req.FinancingOption,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(contracts) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListContractsResponse{}, 0)
		res := SuccessResponse(pag, "No contracts found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	totalCount := contracts[0].TotalCount

	contractsRes := make([]ListContractsResponse, len(contracts))
	for i, contract := range contracts {
		contractsRes[i] = ListContractsResponse{
			ID:              contract.ID,
			ClientID:        contract.ClientID,
			Status:          contract.Status,
			StartDate:       contract.StartDate.Time,
			EndDate:         contract.EndDate.Time,
			Price:           contract.Price,
			PriceTimeUnit:   contract.PriceTimeUnit,
			CareName:        contract.CareName,
			CareType:        contract.CareType,
			FinancingAct:    contract.FinancingAct,
			FinancingOption: contract.FinancingOption,
			SenderID:        contract.SenderID,
			SenderName:      contract.SenderName,
			ClientFirstName: contract.ClientFirstName,
			ClientLastName:  contract.ClientLastName,
			CreatedAt:       contract.CreatedAt.Time,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, contractsRes, totalCount)
	res := SuccessResponse(pag, "Contracts retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// func (server *Server) GenerateContractPdfApi(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	contractID, err := strconv.ParseInt(id, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	contract, err := server.store.GetClientContract(ctx, contractID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

//     documentData := pdf.ContractData{
// 		ID:              contract.ID,
// 		Status: 		contract.Status,

// 	}
// 	pdf, err := pdf.GenerateAndUploadContractPDF(ctx)
// }

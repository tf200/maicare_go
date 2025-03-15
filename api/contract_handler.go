package api

import (
	db "maicare_go/db/sqlc"
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
func (s *Server) CreateContractTypeApi(ctx *gin.Context) {
	var req CreateContractTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contractType, err := s.store.CreateContractType(ctx, req.Name)
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
func (s *Server) ListContractTypesApi(ctx *gin.Context) {
	contractTypes, err := s.store.ListContractTypes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	contractTypesRes := make([]ListContractTypesResponse, len(contractTypes))
	for i, contractType := range contractTypes {
		contractTypesRes[i] = ListContractTypesResponse{
			ID:   contractType.ID,
			Name: contractType.Name,
		}
	}

	res := SuccessResponse(contractTypesRes, "Contract types retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// CreateContractRequest defines the request for CreateContract handler
type CreateContractRequest struct {
	TypeID          *int64      `json:"type_id"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  int32       `json:"reminder_period"`
	Tax             *int32      `json:"tax"`
	Price           float64     `json:"price"`
	PriceFrequency  string      `json:"price_frequency"`
	Hours           *int32      `json:"hours"`
	HoursType       string      `json:"hours_type"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type"`
	SenderID        *int64      `json:"sender_id"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
}

// CreateContractResponse defines the response for CreateContract handler
type CreateContractResponse struct {
	ID              int64              `json:"id"`
	TypeID          *int64             `json:"type_id"`
	Status          string             `json:"status"`
	StartDate       time.Time          `json:"start_date"`
	EndDate         time.Time          `json:"end_date"`
	ReminderPeriod  int32              `json:"reminder_period"`
	Tax             *int32             `json:"tax"`
	Price           float64            `json:"price"`
	PriceFrequency  string             `json:"price_frequency"`
	Hours           *int32             `json:"hours"`
	HoursType       string             `json:"hours_type"`
	CareName        string             `json:"care_name"`
	CareType        string             `json:"care_type"`
	ClientID        int64              `json:"client_id"`
	SenderID        *int64             `json:"sender_id"`
	AttachmentIds   []uuid.UUID        `json:"attachment_ids"`
	FinancingAct    string             `json:"financing_act"`
	FinancingOption string             `json:"financing_option"`
	DepartureReason *string            `json:"departure_reason"`
	DepartureReport *string            `json:"departure_report"`
	Updated         pgtype.Timestamptz `json:"updated"`
	Created         pgtype.Timestamptz `json:"created"`
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
		Tax:             req.Tax,
		Price:           req.Price,
		PriceFrequency:  req.PriceFrequency,
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
		Tax:             contract.Tax,
		Price:           contract.Price,
		PriceFrequency:  contract.PriceFrequency,
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
		Updated:         contract.Updated,
		Created:         contract.Created,
	}, "Contract created successfully")
	ctx.JSON(http.StatusOK, res)

}

// GetContractResponse defines the response for GetContract handler
type GetClientContractResponse struct {
	ID              int64              `json:"id"`
	TypeID          *int64             `json:"type_id"`
	Status          string             `json:"status"`
	StartDate       time.Time          `json:"start_date"`
	EndDate         time.Time          `json:"end_date"`
	ReminderPeriod  int32              `json:"reminder_period"`
	Tax             *int32             `json:"tax"`
	Price           float64            `json:"price"`
	PriceFrequency  string             `json:"price_frequency"`
	Hours           *int32             `json:"hours"`
	HoursType       string             `json:"hours_type"`
	CareName        string             `json:"care_name"`
	CareType        string             `json:"care_type"`
	ClientID        int64              `json:"client_id"`
	SenderID        *int64             `json:"sender_id"`
	AttachmentIds   []uuid.UUID        `json:"attachment_ids"`
	FinancingAct    string             `json:"financing_act"`
	FinancingOption string             `json:"financing_option"`
	DepartureReason *string            `json:"departure_reason"`
	DepartureReport *string            `json:"departure_report"`
	Updated         pgtype.Timestamptz `json:"updated"`
	Created         pgtype.Timestamptz `json:"created"`
}

// GetClientContractApi returns a contract by ID
// @Summary Get a contract by ID
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Contract ID"
// @Success 200 {object} GetContractResponse
// @Router /clients/{id}/contracts [get]

func (server *Server) GetClientContractApi(ctx *gin.Context) {
	id := ctx.Param("id")
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
		Status:          contract.Status,
		StartDate:       contract.StartDate.Time,
		EndDate:         contract.EndDate.Time,
		ReminderPeriod:  contract.ReminderPeriod,
		Tax:             contract.Tax,
		Price:           contract.Price,
		PriceFrequency:  contract.PriceFrequency,
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
		Updated:         contract.Updated,
		Created:         contract.Created,
	}, "Contract retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateClientAllergyRequest defines the request for creating a client allergy
type CreateClientAllergyRequest struct {
	AllergyTypeID int64   `json:"allergy_id" binding:"required"`
	Severity      string  `json:"severity" binding:"required"`
	Reaction      string  `json:"reaction" binding:"required"`
	Notes         *string `json:"notes"`
}

// CreateClientAllergyResponse defines the response for creating a client allergy
type CreateClientAllergyResponse struct {
	ID            int64     `json:"id"`
	ClientID      int64     `json:"client_id"`
	AllergyTypeID int64     `json:"allergy_type_id"`
	Severity      string    `json:"severity"`
	Reaction      string    `json:"reaction"`
	Notes         *string   `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateClientAllergyApi creates a client allergy
// @Summary Create a client allergy
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body CreateClientAllergyRequest true "Client allergy data"
// @Success 200 {object} Response[CreateClientAllergyResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/client_allergies [post]
func (server *Server) CreateClientAllergyApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateClientAllergyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateClientAllergyParams{
		ClientID:      clientID,
		AllergyTypeID: req.AllergyTypeID,
		Severity:      req.Severity,
		Reaction:      req.Reaction,
		Notes:         req.Notes,
	}
	clientAllergy, err := server.store.CreateClientAllergy(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateClientAllergyResponse{
		ID:            clientAllergy.ID,
		ClientID:      clientAllergy.ClientID,
		AllergyTypeID: clientAllergy.AllergyTypeID,
		Severity:      clientAllergy.Severity,
		Reaction:      clientAllergy.Reaction,
		Notes:         clientAllergy.Notes,
		CreatedAt:     clientAllergy.CreatedAt.Time,
	}, "Client allergy created successfully")
	ctx.JSON(http.StatusCreated, res)
}

// ListClientAllergiesRequest defines the request for listing client allergies
type ListClientAllergiesRequest struct {
	pagination.Request
}

// ListClientAllergiesResponse defines the response for listing client allergies
type ListClientAllergiesResponse struct {
	ID            int64     `json:"id"`
	ClientID      int64     `json:"client_id"`
	AllergyTypeID int64     `json:"allergy_type_id"`
	Severity      string    `json:"severity"`
	Reaction      string    `json:"reaction"`
	Notes         *string   `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	AllergyType   string    `json:"allergy_type"`
}

// ListClientAllergiesApi lists all client allergies
// @Summary List all client allergies
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param client_id query int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListClientAllergiesResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/client_allergies [get]
func (server *Server) ListClientAllergiesApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req ListClientAllergiesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	params := req.GetParams()
	arg := db.ListClientAllergiesParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}
	clientAllergies, err := server.store.ListClientAllergies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	totalCount := clientAllergies[0].TotalAllergies

	allergies := make([]ListClientAllergiesResponse, 0)
	for _, allergy := range clientAllergies {
		allergies = append(allergies, ListClientAllergiesResponse{
			ID:            allergy.ID,
			ClientID:      allergy.ClientID,
			AllergyTypeID: allergy.AllergyTypeID,
			Severity:      allergy.Severity,
			Reaction:      allergy.Reaction,
			Notes:         allergy.Notes,
			CreatedAt:     allergy.CreatedAt.Time,
			AllergyType:   allergy.AllergyType,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, allergies, totalCount)
	res := SuccessResponse(pag, "Client allergies fetched successfully")

	ctx.JSON(http.StatusOK, res)
}

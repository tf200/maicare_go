package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListAllergyTypesRequest defines the request for listing allergy types
type ListAllergyTypesRequest struct {
	pagination.Request
	Search *string `form:"search"`
}

// ListAllergyTypesResponse defines the response for listing allergy types
type ListAllergyTypesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ListAllergyTypesApi lists all allergy types
// @Summary List all allergy types
// @Tags client_Medical
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Produce json
// @Success 200 {object} Response[ListAllergyTypesResponse]
// @Failure 400,404 {object} Response[any]
// @Router /allergy_types [get]
func (server *Server) ListAllergyTypesApi(ctx *gin.Context) {
	var req ListAllergyTypesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	params := req.GetParams()
	arg := db.ListAllergiesParams{
		Limit:  params.Limit,
		Offset: params.Offset,
		Search: req.Search,
	}
	allergyTypes, err := server.store.ListAllergies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	allergyTypesResponse := make([]ListAllergyTypesResponse, 0)
	for _, allergyType := range allergyTypes {
		allergyTypesResponse = append(allergyTypesResponse, ListAllergyTypesResponse{
			ID:   allergyType.ID,
			Name: allergyType.Name,
		})
	}
	res := SuccessResponse(allergyTypesResponse, "Allergy types fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

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
// @Success 201 {object} Response[CreateClientAllergyResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/allergies [post]
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
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListClientAllergiesResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/allergies [get]
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

	allergies := make([]ListClientAllergiesResponse, len(clientAllergies))
	for i, allergy := range clientAllergies {
		allergies[i] = ListClientAllergiesResponse{
			ID:            allergy.ID,
			ClientID:      allergy.ClientID,
			AllergyTypeID: allergy.AllergyTypeID,
			Severity:      allergy.Severity,
			Reaction:      allergy.Reaction,
			Notes:         allergy.Notes,
			CreatedAt:     allergy.CreatedAt.Time,
			AllergyType:   allergy.AllergyType,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, allergies, totalCount)
	res := SuccessResponse(pag, "Client allergies fetched successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetClientAllergyResponse defines the request for geting a client allergy
type GetClientAllergyResponse struct {
	ID            int64     `json:"id"`
	ClientID      int64     `json:"client_id"`
	AllergyTypeID int64     `json:"allergy_type_id"`
	Severity      string    `json:"severity"`
	Reaction      string    `json:"reaction"`
	Notes         *string   `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetClientAllergyApi gets a client allergy
// @Summary Get a client allergy
// @Tags client_Medical
// @Produce json
// @Param id path int true "Client ID"
// @Param allergy_id path int true "Allergy ID"
// @Success 200 {object} Response[GetClientAllergyResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/allergies/{allergy_id} [get]
func (server *Server) GetClientAllergyApi(ctx *gin.Context) {
	id := ctx.Param("allergy_id")
	allergyID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	allergy, err := server.store.GetClientAllergy(ctx, allergyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetClientAllergyResponse{
		ID:            allergy.ID,
		ClientID:      allergy.ClientID,
		AllergyTypeID: allergy.AllergyTypeID,
		Severity:      allergy.Severity,
		Reaction:      allergy.Reaction,
		Notes:         allergy.Notes,
		CreatedAt:     allergy.CreatedAt.Time,
	}, "Client allergy fetched successfully")
	ctx.JSON(http.StatusOK, res)

}

// UpdateClientAllergyRequest defines the request for updating a client allergy
type UpdateClientAllergyRequest struct {
	AllergyTypeID *int64  `json:"allergy_type_id"`
	Severity      *string `json:"severity"`
	Reaction      *string `json:"reaction"`
	Notes         *string `json:"notes"`
}

// UpdateClientAllergyResponse defines the response for updating a client allergy
type UpdateClientAllergyResponse struct {
	ID            int64     `json:"id"`
	ClientID      int64     `json:"client_id"`
	AllergyTypeID int64     `json:"allergy_type_id"`
	Severity      string    `json:"severity"`
	Reaction      string    `json:"reaction"`
	Notes         *string   `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// UpdateClientAllergyApi updates a client allergy
// @Summary Update a client allergy
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param allergy_id path int true "Allergy ID"
// @Param request body UpdateClientAllergyRequest true "Client allergy data"
// @Success 200 {object} Response[UpdateClientAllergyResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/allergies/{allergy_id} [put]
func (server *Server) UpdateClientAllergyApi(ctx *gin.Context) {
	id := ctx.Param("allergy_id")
	allergyID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateClientAllergyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateClientAllergyParams{
		ID:            allergyID,
		AllergyTypeID: req.AllergyTypeID,
		Severity:      req.Severity,
		Reaction:      req.Reaction,
		Notes:         req.Notes,
	}

	allergy, err := server.store.UpdateClientAllergy(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateClientAllergyResponse{
		ID:            allergy.ID,
		ClientID:      allergy.ClientID,
		AllergyTypeID: allergy.AllergyTypeID,
		Severity:      allergy.Severity,
		Reaction:      allergy.Reaction,
		Notes:         allergy.Notes,
		CreatedAt:     allergy.CreatedAt.Time,
	}, "Client allergy updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteClientAllergyResponse defines the response for deleting a client allergy
type DeleteClientAllergyResponse struct {
	ID int64 `json:"id"`
}

// DeleteClientAllergyApi deletes a client allergy
// @Summary Delete a client allergy
// @Tags client_Medical
// @Produce json
// @Param id path int true "Client ID"
// @Param allergy_id path int true "Allergy ID"
// @Success 200 {object} Response[any]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/allergies/{allergy_id} [delete]
func (server *Server) DeleteClientAllergyApi(ctx *gin.Context) {
	id := ctx.Param("allergy_id")
	allergyID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	allergy, err := server.store.DeleteClientAllergy(ctx, allergyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(DeleteClientAllergyResponse{
		ID: allergy.ID,
	}, "Client allergy deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateClientDiagnosisRequest defines the request for creating a client diagnosis
type CreateClientDiagnosisRequest struct {
	ClientID            int64     `json:"client_id"`
	Title               *string   `json:"title"`
	DiagnosisCode       string    `json:"diagnosis_code"`
	Description         string    `json:"description"`
	DateOfDiagnosis     time.Time `json:"date_of_diagnosis"`
	Severity            *string   `json:"severity"`
	Status              string    `json:"status"`
	DiagnosingClinician string    `json:"diagnosing_clinician"`
	Notes               *string   `json:"notes"`
}

// CreateClientDiagnosisResponse defines the response for creating a client diagnosis
type CreateClientDiagnosisResponse struct {
	ID                  int64     `json:"id"`
	Title               *string   `json:"title"`
	ClientID            int64     `json:"client_id"`
	DiagnosisCode       string    `json:"diagnosis_code"`
	Description         string    `json:"description"`
	DateOfDiagnosis     time.Time `json:"date_of_diagnosis"`
	Severity            *string   `json:"severity"`
	Status              string    `json:"status"`
	DiagnosingClinician string    `json:"diagnosing_clinician"`
	Notes               *string   `json:"notes"`
	CreatedAt           time.Time `json:"created_at"`
}

// CreateClientDiagnosisApi creates a client diagnosis
// @Summary Create a client diagnosis
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body CreateClientDiagnosisRequest true "Client diagnosis data"
// @Success 201 {object} Response[CreateClientDiagnosisResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/diagnosis [post]
func (server *Server) CreateClientDiagnosisApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateClientDiagnosisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateClientDiagnosisParams{
		ClientID:            clientID,
		Title:               req.Title,
		DiagnosisCode:       req.DiagnosisCode,
		Description:         req.Description,
		DateOfDiagnosis:     pgtype.Timestamptz{Time: req.DateOfDiagnosis, Valid: true},
		Severity:            req.Severity,
		Status:              req.Status,
		DiagnosingClinician: req.DiagnosingClinician,
		Notes:               req.Notes,
	}

	clientDiagnosis, err := server.store.CreateClientDiagnosis(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateClientDiagnosisResponse{
		ID:                  clientDiagnosis.ID,
		Title:               clientDiagnosis.Title,
		ClientID:            clientDiagnosis.ClientID,
		DiagnosisCode:       clientDiagnosis.DiagnosisCode,
		Description:         clientDiagnosis.Description,
		DateOfDiagnosis:     clientDiagnosis.DateOfDiagnosis.Time,
		Severity:            clientDiagnosis.Severity,
		Status:              clientDiagnosis.Status,
		DiagnosingClinician: clientDiagnosis.DiagnosingClinician,
		Notes:               clientDiagnosis.Notes,
		CreatedAt:           clientDiagnosis.CreatedAt.Time,
	}, "Client diagnosis created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListClientDiagnosesRequest defines the request for listing client diagnoses
type ListClientDiagnosesRequest struct {
	pagination.Request
}

// ListClientDiagnosesResponse defines the response for listing client diagnoses
type ListClientDiagnosesResponse struct {
	ID                  int64     `json:"id"`
	Title               *string   `json:"title"`
	ClientID            int64     `json:"client_id"`
	DiagnosisCode       string    `json:"diagnosis_code"`
	Description         string    `json:"description"`
	DateOfDiagnosis     time.Time `json:"date_of_diagnosis"`
	Severity            *string   `json:"severity"`
	Status              string    `json:"status"`
	DiagnosingClinician string    `json:"diagnosing_clinician"`
	Notes               *string   `json:"notes"`
	CreatedAt           time.Time `json:"created_at"`
}

// ListClientDiagnosesApi lists all client diagnoses
// @Summary List all client diagnoses
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListClientDiagnosesResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/diagnosis [get]
func (server *Server) ListClientDiagnosesApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListClientDiagnosesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()
	arg := db.ListClientDiagnosesParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}

	clientDiagnoses, err := server.store.ListClientDiagnoses(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalCount := clientDiagnoses[0].TotalDiagnoses

	diagnoses := make([]ListClientDiagnosesResponse, len(clientDiagnoses))
	for i, diagnosis := range clientDiagnoses {
		diagnoses[i] = ListClientDiagnosesResponse{
			ID:                  diagnosis.ID,
			Title:               diagnosis.Title,
			ClientID:            diagnosis.ClientID,
			DiagnosisCode:       diagnosis.DiagnosisCode,
			Description:         diagnosis.Description,
			DateOfDiagnosis:     diagnosis.DateOfDiagnosis.Time,
			Severity:            diagnosis.Severity,
			Status:              diagnosis.Status,
			DiagnosingClinician: diagnosis.DiagnosingClinician,
			Notes:               diagnosis.Notes,
			CreatedAt:           diagnosis.CreatedAt.Time,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, diagnoses, totalCount)
	res := SuccessResponse(pag, "Client diagnoses fetched successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetClientDiagnosisResponse defines the response for getting a client diagnosis
type GetClientDiagnosisResponse struct {
	ID                  int64     `json:"id"`
	Title               *string   `json:"title"`
	ClientID            int64     `json:"client_id"`
	DiagnosisCode       string    `json:"diagnosis_code"`
	Description         string    `json:"description"`
	DateOfDiagnosis     time.Time `json:"date_of_diagnosis"`
	Severity            *string   `json:"severity"`
	Status              string    `json:"status"`
	DiagnosingClinician string    `json:"diagnosing_clinician"`
	Notes               *string   `json:"notes"`
	CreatedAt           time.Time `json:"created_at"`
}

// GetClientDiagnosisApi gets a client diagnosis
// @Summary Get a client diagnosis
// @Tags client_Medical
// @Produce json
// @Param id path int true "Client ID"
// @Param diagnosis_id path int true "Diagnosis ID"
// @Success 200 {object} Response[GetClientDiagnosisResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/diagnosis/{diagnosis_id} [get]
func (server *Server) GetClientDiagnosisApi(ctx *gin.Context) {
	id := ctx.Param("diagnosis_id")
	diagnosisID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	diagnosis, err := server.store.GetClientDiagnosis(ctx, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetClientDiagnosisResponse{
		ID:                  diagnosis.ID,
		Title:               diagnosis.Title,
		ClientID:            diagnosis.ClientID,
		DiagnosisCode:       diagnosis.DiagnosisCode,
		Description:         diagnosis.Description,
		DateOfDiagnosis:     diagnosis.DateOfDiagnosis.Time,
		Severity:            diagnosis.Severity,
		Status:              diagnosis.Status,
		DiagnosingClinician: diagnosis.DiagnosingClinician,
		Notes:               diagnosis.Notes,
		CreatedAt:           diagnosis.CreatedAt.Time,
	}, "Client diagnosis fetched successfully")
	ctx.JSON(http.StatusOK, res)

}

// UpdateClientDiagnosisApi updates a client diagnosis
type UpdateClientDiagnosisRequest struct {
	Title               *string            `json:"title"`
	DiagnosisCode       *string            `json:"diagnosis_code"`
	Description         *string            `json:"description"`
	DateOfDiagnosis     pgtype.Timestamptz `json:"date_of_diagnosis"`
	Severity            *string            `json:"severity"`
	Status              *string            `json:"status"`
	DiagnosingClinician *string            `json:"diagnosing_clinician"`
	Notes               *string            `json:"notes"`
}

// UpdateClientDiagnosisApi updates a client diagnosis
type UpdateClientDiagnosisResponse struct {
	ID                  int64              `json:"id"`
	Title               *string            `json:"title"`
	ClientID            int64              `json:"client_id"`
	DiagnosisCode       string             `json:"diagnosis_code"`
	Description         string             `json:"description"`
	DateOfDiagnosis     pgtype.Timestamptz `json:"date_of_diagnosis"`
	Severity            *string            `json:"severity"`
	Status              string             `json:"status"`
	DiagnosingClinician string             `json:"diagnosing_clinician"`
	Notes               *string            `json:"notes"`
	CreatedAt           pgtype.Timestamptz `json:"created_at"`
}

// UpdateClientDiagnosisApi updates a client diagnosis
// @Summary Update a client diagnosis
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param diagnosis_id path int true "Diagnosis ID"
// @Param request body UpdateClientDiagnosisRequest true "Client diagnosis data"
// @Success 200 {object} Response[UpdateClientDiagnosisResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/diagnosis/{diagnosis_id} [put]
func (server *Server) UpdateClientDiagnosisApi(ctx *gin.Context) {
	id := ctx.Param("diagnosis_id")
	diagnosisID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateClientDiagnosisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateClientDiagnosisParams{
		ID:                  diagnosisID,
		Title:               req.Title,
		DiagnosisCode:       req.DiagnosisCode,
		Description:         req.Description,
		DateOfDiagnosis:     req.DateOfDiagnosis,
		Severity:            req.Severity,
		Status:              req.Status,
		DiagnosingClinician: req.DiagnosingClinician,
		Notes:               req.Notes,
	}

	diagnosis, err := server.store.UpdateClientDiagnosis(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateClientDiagnosisResponse{
		ID:                  diagnosis.ID,
		Title:               diagnosis.Title,
		ClientID:            diagnosis.ClientID,
		DiagnosisCode:       diagnosis.DiagnosisCode,
		Description:         diagnosis.Description,
		DateOfDiagnosis:     diagnosis.DateOfDiagnosis,
		Severity:            diagnosis.Severity,
		Status:              diagnosis.Status,
		DiagnosingClinician: diagnosis.DiagnosingClinician,
		Notes:               diagnosis.Notes,
		CreatedAt:           diagnosis.CreatedAt,
	}, "Client diagnosis updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteClientDiagnosisResponse defines the response for deleting a client diagnosis
type DeleteClientDiagnosisResponse struct {
	ID int64 `json:"id"`
}

// DeleteClientDiagnosisApi deletes a client diagnosis
// @Summary Delete a client diagnosis
// @Tags client_Medical
// @Produce json
// @Param id path int true "Client ID"
// @Param diagnosis_id path int true "Diagnosis ID"
// @Success 200 {object} Response[any]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/diagnosis/{diagnosis_id} [delete]
func (server *Server) DeleteClientDiagnosisApi(ctx *gin.Context) {
	id := ctx.Param("diagnosis_id")
	diagnosisID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	diagnosis, err := server.store.DeleteClientDiagnosis(ctx, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(DeleteClientDiagnosisResponse{
		ID: diagnosis.ID,
	}, "Client diagnosis deleted successfully")

	ctx.JSON(http.StatusOK, res)
}

package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	clientp "maicare_go/service/client"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

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

	res := SuccessResponse(result, "Client diagnosis created successfully")

	ctx.JSON(http.StatusCreated, res)
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

	res := SuccessResponse(pag, "Client diagnoses fetched successfully")
	ctx.JSON(http.StatusOK, res)

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

	diagnosis, err := server.businessService.ClientService.GetClientDiagnosis(ctx, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(diagnosis, "Client diagnosis fetched successfully")
	ctx.JSON(http.StatusOK, res)

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

	var req clientp.UpdateClientDiagnosisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	diagnosis, err := server.businessService.ClientService.UpdateClientDiagnosis(ctx, req, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(diagnosis, "Client diagnosis updated successfully")

	ctx.JSON(http.StatusOK, res)
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

	diag, err := server.businessService.ClientService.DeleteClientDiagnosis(ctx, diagnosisID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(diag, "Client diagnosis deleted successfully")

	ctx.JSON(http.StatusOK, res)
}



// CreateClientMedicationApi creates a client medication
// @Summary Create a client medication
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param diagnosis_id path int true "Diagnosis ID"
// @Param request body CreateclientMedicationRequest true "Client medication data"
// @Success 201 {object} Response[CreateClientMedicationResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/diagnosis/{diagnosis_id}/medications [post]
func (server *Server) CreateClientMedicationApi(ctx *gin.Context) {
	id := ctx.Param("diagnosis_id")
	diagnosisID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateclientMedicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateClientMedicationParams{
		DiagnosisID:      &diagnosisID,
		Name:             req.Name,
		Dosage:           req.Dosage,
		StartDate:        pgtype.Date{Time: req.StartDate, Valid: true},
		EndDate:          pgtype.Date{Time: req.EndDate, Valid: true},
		Notes:            req.Notes,
		SelfAdministered: req.SelfAdministered,
		AdministeredByID: req.AdministeredByID,
		IsCritical:       req.IsCritical,
	}

	clientMedication, err := server.store.CreateClientMedication(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateClientMedicationResponse{
		ID:               clientMedication.ID,
		DiagnosisID:      clientMedication.DiagnosisID,
		Name:             clientMedication.Name,
		Dosage:           clientMedication.Dosage,
		StartDate:        clientMedication.StartDate.Time,
		EndDate:          clientMedication.EndDate.Time,
		Notes:            clientMedication.Notes,
		SelfAdministered: clientMedication.SelfAdministered,
		AdministeredByID: clientMedication.AdministeredByID,
		IsCritical:       clientMedication.IsCritical,
		UpdatedAt:        clientMedication.UpdatedAt.Time,
		CreatedAt:        clientMedication.CreatedAt.Time,
	}, "Client medication created successfully")

	ctx.JSON(http.StatusCreated, res)

}

type ListClientMedicationsRequest struct {
	pagination.Request
}

// ListClientMedicationsResponse defines the response for listing client medications
type ListClientMedicationsResponse struct {
	ID               int64     `json:"id"`
	DiagnosisID      *int64    `json:"diagnosis_id"`
	Name             string    `json:"name"`
	Dosage           string    `json:"dosage"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Notes            *string   `json:"notes"`
	SelfAdministered bool      `json:"self_administered"`
	AdministeredByID *int64    `json:"administered_by_id"`
	IsCritical       bool      `json:"is_critical"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedAt        time.Time `json:"created_at"`
}

// ListClientMedicationsApi lists all client medications
// @Summary List all client medications
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param diagnosis_id path int true "Diagnosis ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListClientMedicationsResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/diagnosis/{diagnosis_id}/medications [get]
func (server *Server) ListClientMedicationsApi(ctx *gin.Context) {
	id := ctx.Param("diagnosis_id")
	diagnosisID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListClientMedicationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()
	arg := db.ListMedicationsByDiagnosisIDParams{
		DiagnosisID: &diagnosisID,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}

	clientMedications, err := server.store.ListMedicationsByDiagnosisID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(clientMedications) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListClientMedicationsResponse{}, 0)
		res := SuccessResponse(pag, "No client medications found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	totalCount := clientMedications[0].TotalMedications

	clientMedicationList := make([]ListClientMedicationsResponse, len(clientMedications))
	for i, medication := range clientMedications {
		clientMedicationList[i] = ListClientMedicationsResponse{
			ID:               medication.ID,
			Name:             medication.Name,
			Dosage:           medication.Dosage,
			StartDate:        medication.StartDate.Time,
			EndDate:          medication.EndDate.Time,
			Notes:            medication.Notes,
			SelfAdministered: medication.SelfAdministered,
			IsCritical:       medication.IsCritical,
			DiagnosisID:      medication.DiagnosisID,
			AdministeredByID: medication.AdministeredByID,
			UpdatedAt:        medication.UpdatedAt.Time,
			CreatedAt:        medication.CreatedAt.Time,
		}
	}
	pag := pagination.NewResponse(ctx, req.Request, clientMedicationList, totalCount)
	res := SuccessResponse(pag, "Client medications fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientMedicationResponse defines the response for getting a client medication
type GetClientMedicationResponse struct {
	ID                      int64     `json:"id"`
	Name                    string    `json:"name"`
	Dosage                  string    `json:"dosage"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	Notes                   *string   `json:"notes"`
	SelfAdministered        bool      `json:"self_administered"`
	DiagnosisID             *int64    `json:"diagnosis_id"`
	AdministeredByID        *int64    `json:"administered_by_id"`
	IsCritical              bool      `json:"is_critical"`
	UpdatedAt               time.Time `json:"updated_at"`
	CreatedAt               time.Time `json:"created_at"`
	AdministeredByFirstName string    `json:"administered_by_first_name"`
	AdministeredByLastName  string    `json:"administered_by_last_name"`
}

// GetClientMedicationApi gets a client medication
// @Summary Get a client medication
// @Tags client_Medical
// @Produce json
// @Param id path int true "Client ID"
// @Param medication_id path int true "Medication ID"
// @Param diagnosis_id path int true "Diagnosis ID"
// @Success 200 {object} Response[GetClientMedicationResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/diagnosis/{diagnosis_id}/medications/{medication_id} [get]
func (server *Server) GetClientMedicationApi(ctx *gin.Context) {
	id := ctx.Param("medication_id")
	medicationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	medication, err := server.store.GetMedication(ctx, medicationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetClientMedicationResponse{
		ID:                      medication.ID,
		Name:                    medication.Name,
		Dosage:                  medication.Dosage,
		StartDate:               medication.StartDate.Time,
		EndDate:                 medication.EndDate.Time,
		Notes:                   medication.Notes,
		SelfAdministered:        medication.SelfAdministered,
		DiagnosisID:             medication.DiagnosisID,
		AdministeredByID:        medication.AdministeredByID,
		IsCritical:              medication.IsCritical,
		UpdatedAt:               medication.UpdatedAt.Time,
		CreatedAt:               medication.CreatedAt.Time,
		AdministeredByFirstName: medication.AdministeredByFirstName,
		AdministeredByLastName:  medication.AdministeredByLastName,
	}, "Client medication fetched successfully")
	ctx.JSON(http.StatusOK, res)

}

// UpdateClientMedicationRequest defines the request for updating a client medication
type UpdateClientMedicationRequest struct {
	Name             *string   `json:"name"`
	Dosage           *string   `json:"dosage"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Notes            *string   `json:"notes"`
	SelfAdministered *bool     `json:"self_administered"`
	AdministeredByID *int64    `json:"administered_by_id"`
	IsCritical       *bool     `json:"is_critical"`
}

// UpdateClientMedicationResponse defines the response for updating a client medication
type UpdateClientMedicationResponse struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Dosage           string    `json:"dosage"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Notes            *string   `json:"notes"`
	SelfAdministered bool      `json:"self_administered"`
	DiagnosisID      *int64    `json:"diagnosis_id"`
	AdministeredByID *int64    `json:"administered_by_id"`
	IsCritical       bool      `json:"is_critical"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedAt        time.Time `json:"created_at"`
}

// UpdateClientMedicationApi updates a client medication
// @Summary Update a client medication
// @Tags client_Medical
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param diagnosis_id path int true "Diagnosis ID"
// @Param medication_id path int true "Medication ID"
// @Param request body UpdateClientMedicationRequest true "Client medication data"
// @Success 200 {object} Response[UpdateClientMedicationResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/diagnosis/{diagnosis_id}/medications/{medication_id} [put]
func (server *Server) UpdateClientMedicationApi(ctx *gin.Context) {
	id := ctx.Param("medication_id")
	medicationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateClientMedicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateClientMedicationParams{
		ID:               medicationID,
		Name:             req.Name,
		Dosage:           req.Dosage,
		StartDate:        pgtype.Date{Time: req.StartDate, Valid: true},
		EndDate:          pgtype.Date{Time: req.EndDate, Valid: true},
		Notes:            req.Notes,
		SelfAdministered: req.SelfAdministered,
		AdministeredByID: req.AdministeredByID,
		IsCritical:       req.IsCritical,
	}

	medication, err := server.store.UpdateClientMedication(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateClientMedicationResponse{
		ID:               medication.ID,
		Name:             medication.Name,
		Dosage:           medication.Dosage,
		StartDate:        medication.StartDate.Time,
		EndDate:          medication.EndDate.Time,
		Notes:            medication.Notes,
		SelfAdministered: medication.SelfAdministered,
		DiagnosisID:      medication.DiagnosisID,
		AdministeredByID: medication.AdministeredByID,
		IsCritical:       medication.IsCritical,
		UpdatedAt:        medication.UpdatedAt.Time,
		CreatedAt:        medication.CreatedAt.Time,
	}, "Client medication updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteClientMedicationApi deletes a client medication
// @Summary Delete a client medication
// @Tags client_Medical
// @Produce json
// @Param id path int true "Client ID"
// @Param diagnosis_id path int true "Diagnosis ID"
// @Param medication_id path int true "Medication ID"
// @Success 200 {object} Response[any]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/medications/{medication_id} [delete]
func (server *Server) DeleteClientMedicationApi(ctx *gin.Context) {
	id := ctx.Param("medication_id")
	medicationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteClientMedication(ctx, medicationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Client medication deleted successfully")

	ctx.JSON(http.StatusOK, res)
}

package handler

import (
	"errors"
	"io"
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"
	"maicare_go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RegistrationFormHandler struct {
	service domain.RegistrationFormService
}

func NewRegistrationFormHandler(service domain.RegistrationFormService) *RegistrationFormHandler {
	return &RegistrationFormHandler{service: service}
}

func (h *RegistrationFormHandler) StartUploadSession(ctx *gin.Context) {
	token, err := h.service.StartUploadSession(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to start registration", ""))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK(gin.H{"registration_token": token}, "Registration started successfully"))
}

func (h *RegistrationFormHandler) InitRegistrationUpload(ctx *gin.Context) {
	token := ctx.GetHeader("X-Registration-Token")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("registration token is required", ""))
		return
	}
	var req initAttachmentUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}
	result, err := h.service.InitRegistrationUpload(ctx.Request.Context(), token, domain.InitAttachmentUploadParams{
		Filename: req.Filename, ContentType: req.ContentType, Size: req.Size,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK(result, "Upload initiated successfully"))
}

func (h *RegistrationFormHandler) CreateRegistrationForm(ctx *gin.Context) {
	var req createRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	params, err := toCreateRegistrationFormParams(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}
	params.RegistrationUploadToken = ctx.GetHeader("X-Registration-Token")
	if params.RegistrationUploadToken == "" {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("registration token is required", ""))
		return
	}
	result, err := h.service.CreateRegistrationForm(ctx.Request.Context(), params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create registration form", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toRegistrationFormResponse(*result), "Registration form created successfully"))
}

func (h *RegistrationFormHandler) ListRegistrationForms(ctx *gin.Context) {
	var req listRegistrationFormsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.ListRegistrationForms(ctx.Request.Context(), toListRegistrationFormsParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list registration forms", ""))
		return
	}

	items := make([]registrationFormListItemResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = toRegistrationFormListItemResponse(item)
	}

	ctx.JSON(http.StatusOK, httpapi.OK(httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount), "Registration forms retrieved successfully"))
}

func (h *RegistrationFormHandler) GetRegistrationFormCounts(ctx *gin.Context) {
	result, err := h.service.GetRegistrationFormCounts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get registration form counts", ""))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK(gin.H{
		"total":          result.Total,
		"pending_review": result.PendingReview,
		"processed":      result.Processed,
		"high_risk":      result.HighRisk,
	}, "Registration form counts retrieved successfully"))
}

func (h *RegistrationFormHandler) GetRegistrationForm(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid registration form ID", ""))
		return
	}

	result, err := h.service.GetRegistrationForm(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, domain.ErrRegistrationFormNotFound) {
			ctx.JSON(http.StatusNotFound, httpapi.Fail("registration form not found", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get registration form", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toRegistrationFormResponse(*result), "Registration form retrieved successfully"))
}

func (h *RegistrationFormHandler) UpdateRegistrationForm(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid registration form ID", ""))
		return
	}

	var req updateRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	params, err := toUpdateRegistrationFormParams(id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.UpdateRegistrationForm(ctx.Request.Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, domain.ErrRegistrationFormNotFound) {
			ctx.JSON(http.StatusNotFound, httpapi.Fail("registration form not found", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update registration form", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toRegistrationFormResponse(*result), "Registration form updated successfully"))
}

func (h *RegistrationFormHandler) ReplaceRegistrationFormDocument(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid registration form ID", ""))
		return
	}

	var req replaceRegistrationFormDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}
	if !supportedRegistrationDocumentTypes[req.DocumentType] {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid document_type", "unsupported registration document type"))
		return
	}

	if _, err := h.service.GetAttachment(ctx.Request.Context(), req.FileID); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("attachment not found", ""))
		return
	}
	result, err := h.service.ReplaceRegistrationFormDocument(ctx.Request.Context(), domain.ReplaceRegistrationFormDocumentParams{
		ID: id, DocumentType: req.DocumentType, FileID: req.FileID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrRegistrationFormNotFound) || errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, httpapi.Fail("registration form not found", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to replace registration document", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toRegistrationFormResponse(*result), "Registration document replaced successfully"))
}

func (h *RegistrationFormHandler) DeleteRegistrationForm(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid registration form ID", ""))
		return
	}

	if err := h.service.DeleteRegistrationForm(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete registration form", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Registration form deleted successfully"))
}

func (h *RegistrationFormHandler) UpdateRegistrationFormStatus(ctx *gin.Context) {
	payload, ok := middleware.AuthPayloadFromContext(ctx.Request.Context())
	if !ok || payload == nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized access", ""))
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid registration form ID", ""))
		return
	}

	var req updateRegistrationFormStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	if err := h.service.UpdateRegistrationFormStatus(ctx.Request.Context(), toUpdateRegistrationFormStatusParams(id, req, payload.EmployeeID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update registration form status", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Registration form status updated successfully"))
}

func (h *RegistrationFormHandler) ProcessRegistrationForm(ctx *gin.Context) {
	payload, ok := middleware.AuthPayloadFromContext(ctx.Request.Context())
	if !ok || payload == nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized access", ""))
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid registration form ID", ""))
		return
	}

	var req processRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	if err := h.service.ProcessRegistrationForm(ctx.Request.Context(), domain.ProcessRegistrationFormParams{
		ID:                        id,
		EmployeeID:                payload.EmployeeID,
		IntakeAppointmentLocation: req.IntakeAppointmentLocation,
		AddmissionType:            req.AddmissionType,
		ProposedDates:             req.ProposedDates,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to process registration form", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Registration form processed and intake options sent"))
}

func (h *RegistrationFormHandler) GetPublicIntakeOptions(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("token is required", ""))
		return
	}

	result, err := h.service.GetPublicIntakeOptions(ctx.Request.Context(), token)
	if err != nil {
		ctx.JSON(http.StatusNotFound, httpapi.Fail("invalid token or form not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toPublicIntakeOptionsResponse(*result), "Intake options retrieved successfully"))
}

func (h *RegistrationFormHandler) SelectIntakeDate(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("token is required", ""))
		return
	}

	var req selectIntakeDateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	if err := h.service.SelectIntakeDate(ctx.Request.Context(), domain.SelectIntakeDateParams{
		Token:        token,
		SelectedDate: req.SelectedDate,
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Intake date selected successfully"))
}

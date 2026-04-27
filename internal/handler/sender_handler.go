package handler

import (
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SenderHandler struct {
	service domain.SenderService
}

func NewSenderHandler(service domain.SenderService) *SenderHandler {
	return &SenderHandler{service: service}
}

// CreateSender creates a new sender.
// @Summary Create a new sender
// @Tags senders
// @Accept json
// @Produce json
// @Param request body createSenderRequest true "Sender data"
// @Success 201 {object} httpapi.Envelope[senderResponse]
// @Router /senders [post]
func (h *SenderHandler) CreateSender(ctx *gin.Context) {
	var req createSenderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.CreateSender(ctx.Request.Context(), toCreateSenderParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create sender", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toSenderResponse(result), "Sender created successfully"))
}

// ListSenders returns a paginated list of senders.
// @Summary List senders
// @Tags senders
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search"
// @Param include_archived query bool false "Include archived"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[senderResponse]]
// @Router /senders [get]
func (h *SenderHandler) ListSenders(ctx *gin.Context) {
	var req listSendersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid query params", err.Error()))
		return
	}

	result, err := h.service.ListSenders(ctx.Request.Context(), toListSendersParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list senders", err.Error()))
		return
	}

	items := make([]senderResponse, len(result.Items))
	for i, s := range result.Items {
		items[i] = toSenderResponse(&s)
	}

	ctx.JSON(http.StatusOK, httpapi.OK(httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount), "Senders retrieved successfully"))
}

// GetSenderByID returns a sender by ID with invoice template items.
// @Summary Get a sender
// @Tags senders
// @Produce json
// @Param id path string true "Sender ID"
// @Success 200 {object} httpapi.Envelope[senderDetailResponse]
// @Router /senders/{id} [get]
func (h *SenderHandler) GetSenderByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid sender ID", err.Error()))
		return
	}

	result, err := h.service.GetSenderByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get sender", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toSenderDetailResponse(result), "Sender retrieved successfully"))
}

// UpdateSender updates an existing sender.
// @Summary Update a sender
// @Tags senders
// @Accept json
// @Produce json
// @Param id path string true "Sender ID"
// @Param request body updateSenderRequest true "Sender data"
// @Success 200 {object} httpapi.Envelope[senderResponse]
// @Router /senders/{id} [put]
func (h *SenderHandler) UpdateSender(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid sender ID", err.Error()))
		return
	}

	var req updateSenderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.UpdateSender(ctx.Request.Context(), toUpdateSenderParams(id, req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update sender", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toSenderResponse(result), "Sender updated successfully"))
}

// DeleteSender deletes (archives) a sender.
// @Summary Delete a sender
// @Tags senders
// @Produce json
// @Param id path string true "Sender ID"
// @Success 200 {object} httpapi.Envelope[any]
// @Router /senders/{id} [delete]
func (h *SenderHandler) DeleteSender(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid sender ID", err.Error()))
		return
	}

	if err := h.service.DeleteSender(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete sender", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Sender deleted successfully"))
}

// CreateSenderInvoiceTemplate creates a new invoice template for a sender.
// @Summary Create a sender invoice template
// @Tags senders
// @Accept json
// @Produce json
// @Param id path string true "Sender ID"
// @Param request body createSenderInvoiceTemplateRequest true "Invoice template IDs"
// @Success 200 {object} httpapi.Envelope[any]
// @Router /senders/{id}/invoice_template [post]
func (h *SenderHandler) CreateSenderInvoiceTemplate(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid sender ID", err.Error()))
		return
	}

	var req createSenderInvoiceTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	if err := h.service.CreateSenderInvoiceTemplate(ctx.Request.Context(), id, req.InvoiceTemplateIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create sender invoice template", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Sender invoice template created successfully"))
}

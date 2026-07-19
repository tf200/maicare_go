package handler

import (
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttachmentHandler struct {
	service domain.AttachmentService
}

func NewAttachmentHandler(service domain.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: service}
}

func RegisterAttachmentRoutes(rg *gin.RouterGroup, handler *AttachmentHandler, auth gin.HandlerFunc) {
	attachments := rg.Group("/attachments")
	attachments.POST("/upload/init", auth, handler.InitUpload)
	attachments.POST("/upload/confirm", auth, handler.ConfirmUpload)
	attachments.GET("/:id", auth, handler.GetAttachment)
	attachments.DELETE("/:id", auth, handler.DeleteAttachment)
}

type initAttachmentUploadRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	Size        int64  `json:"size" binding:"required"`
}

type confirmAttachmentUploadRequest struct {
	FileID uuid.UUID `json:"file_id" binding:"required"`
}

func (h *AttachmentHandler) InitUpload(ctx *gin.Context) {
	var req initAttachmentUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}
	result, err := h.service.InitUpload(ctx.Request.Context(), domain.InitAttachmentUploadParams{
		Filename: req.Filename, ContentType: req.ContentType, Size: req.Size,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK(result, "Upload initiated successfully"))
}

func (h *AttachmentHandler) ConfirmUpload(ctx *gin.Context) {
	var req confirmAttachmentUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}
	result, err := h.service.ConfirmUpload(ctx.Request.Context(), req.FileID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to confirm upload", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK(result, "Upload confirmed successfully"))
}

func (h *AttachmentHandler) GetAttachment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid attachment ID format", ""))
		return
	}
	result, err := h.service.GetAttachment(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get attachment", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK(result, "Attachment retrieved successfully"))
}

func (h *AttachmentHandler) DeleteAttachment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid attachment ID format", ""))
		return
	}
	result, err := h.service.DeleteAttachment(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete attachment", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK(result, "Attachment deleted successfully"))
}

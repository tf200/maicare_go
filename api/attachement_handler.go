package api

import (
	"fmt"
	"net/http"

	"maicare_go/service/attachment"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// InitUploadHandlerApi initiates a file upload
// @Summary Initiate file upload
// @Description Initiate a file upload to get a presigned URL
// @Tags attachments
// @Accept json
// @Produce json
// @Param request body attachment.InitUploadRequest true "Init Upload Request"
// @Success 200 {object} Response[attachment.InitUploadResponse]
// @Router /attachments/upload/init [post]
func (server *Server) InitUploadHandlerApi(ctx *gin.Context) {
	var req attachment.InitUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "InitUploadHandlerApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.AttachmentService.InitUpload(ctx, &req)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "InitUploadHandlerApi", "Failed to initiate upload", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Upload initiated successfully")
	ctx.JSON(http.StatusOK, res)
}

// ConfirmUploadHandlerApi confirms a file upload
// @Summary Confirm file upload
// @Description Confirm that a file has been uploaded to the storage
// @Tags attachments
// @Accept json
// @Produce json
// @Param request body attachment.ConfirmUploadRequest true "Confirm Upload Request"
// @Success 200 {object} Response[attachment.ConfirmUploadResponse]
// @Router /attachments/upload/confirm [post]
func (server *Server) ConfirmUploadHandlerApi(ctx *gin.Context) {
	var req attachment.ConfirmUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "ConfirmUploadHandlerApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.AttachmentService.ConfirmUpload(ctx, &req)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ConfirmUploadHandlerApi", "Failed to confirm upload", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Upload confirmed successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetAttachmentByIdApi retrieves an attachment by its ID
// @Summary Get an attachment by ID
// @Description Get an attachment by its ID
// @Tags attachments
// @Produce json
// @Param id path uuid true "Attachment ID"
// @Success 200 {object} Response[attachment.GetAttachmentByIdResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /attachments/{id} [get]
func (server *Server) GetAttachmentByIdApi(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetAttachmentByIdApi", "Invalid UUID format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid attachment ID format")))
		return
	}
	attachment, err := server.businessService.AttachmentService.GetAttachmentById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(attachment, "Attachment retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteAttachment deletes an attachment by its ID
// @Summary Delete an attachment by ID
// @Description Delete an attachment by its ID
// @Tags attachments
// @Produce json
// @Param id path uuid true "Attachment ID"
// @Success 200 {object} Response[attachment.DeleteAttachmentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /attachments/{id} [delete]
func (server *Server) DeleteAttachment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	attachment, err := server.businessService.AttachmentService.DeleteAttachment(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(attachment, "Attachment deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

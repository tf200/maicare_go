package api

import (
	"fmt"
	"maicare_go/service/attachment"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UploadHandler uploads a file to the server
// @Summary Upload a file
// @Description Upload a file to the server
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} Response[attachment.UploadHandlerResponse]
// @Router /attachments/upload [post]
func (server *Server) UploadHandlerApi(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, attachment.MaxFileSize)
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			ctx.JSON(http.StatusRequestEntityTooLarge, errorResponse(fmt.Errorf("file size exceeds maximum limit of 10MB")))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := server.businessService.AttachmentService.UploadAttachment(ctx, file, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to upload attachment: %v", err)))
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

// GetAttachmentByIdApi retrieves an attachment by its ID
// @Summary Get an attachment by ID
// @Description Get an attachment by its ID
// @Tags attachments
// @Produce json
// @Param id path string true "Attachment ID"
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
// @Param id path string true "Attachment ID"
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

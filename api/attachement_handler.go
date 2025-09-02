package api

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const maxFileSize int64 = 10 << 20 // 10 MB

var allowedMimeTypes = map[string]bool{
	// Images
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,

	// PDF
	"application/pdf": true,

	// Microsoft Office Documents
	"application/msword": true, // .doc
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
	"application/vnd.ms-excel": true, // .xls
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true, // .xlsx

	// OpenDocument Formats
	"application/vnd.oasis.opendocument.text":        true, // .odt
	"application/vnd.oasis.opendocument.spreadsheet": true, // .ods

	// Allow ZIP files (for .docx, etc.)
	"application/zip": true,
}

// UploadHandler handles file uploads
type UploadHandlerResponse struct {
	FileURL   string    `json:"file_url"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

// UploadHandler uploads a file to the server
// @Summary Upload a file
// @Description Upload a file to the server
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} Response[UploadHandlerResponse]
// @Router /attachments/upload [post]
func (server *Server) UploadHandlerApi(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxFileSize)
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			server.logBusinessEvent(LogLevelError, "UploadHandlerApi", "File size exceeds maximum limit", zap.Error(err))
			ctx.JSON(http.StatusRequestEntityTooLarge, errorResponse(fmt.Errorf("file size exceeds maximum limit of 10MB")))
			return
		}
		server.logBusinessEvent(LogLevelError, "UploadHandlerApi", "Failed to get form file", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Basic validations
	if err := bucket.ValidateFile(header, maxFileSize); err != nil {
		server.logBusinessEvent(LogLevelError, "UploadHandlerApi", "File validation failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("file validation error")))
		return
	}

	filename := bucket.GenerateUniqueFilename(header.Filename)

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil && err != io.EOF {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error reading file: %v", err)))
		return
	}

	// Reset file pointer after reading
	if _, err := file.Seek(0, 0); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error resetting file: %v", err)))
		return
	}

	// Verify content type
	contentType := http.DetectContentType(buff)
	if !allowedMimeTypes[contentType] {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("unsupported file type: %s", contentType)))
		return
	}

	key, size, err := server.b2Client.Upload(ctx, file, filename, contentType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateAttachmentParams{
		Name: filename,
		File: key,
		Size: int32(size),
		Tag:  util.StringPtr(""),
	}
	attachment, err := server.store.CreateAttachment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UploadHandlerResponse{
		FileURL:   key,
		FileID:    attachment.Uuid,
		CreatedAt: attachment.Created.Time,
		Size:      int64(attachment.Size),
	}, "File uploaded successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetAttachmentByIdResponse represents the response for GetAttachmentByIdApi
type GetAttachmentByIdResponse struct {
	FileURL   string    `json:"file_url"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

// GetAttachmentByIdApi retrieves an attachment by its ID
// @Summary Get an attachment by ID
// @Description Get an attachment by its ID
// @Tags attachments
// @Produce json
// @Param id path string true "Attachment ID"
// @Success 200 {object} Response[GetAttachmentByIdResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /attachments/{id} [get]
func (server *Server) GetAttachmentByIdApi(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetAttachmentByIdApi", "Invalid UUID format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid attachment ID format")))
		return
	}

	attachment, err := server.store.GetAttachmentById(ctx, id)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetAttachmentByIdApi", "Failed to get attachment by ID", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get attachment by ID")))
		return
	}

	presignedUrl, err := server.b2Client.GeneratePresignedURL(ctx, attachment.File, time.Minute*15)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetAttachmentByIdApi", "Failed to generate presigned URL", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to generate presigned URL")))
		return
	}

	res := SuccessResponse(GetAttachmentByIdResponse{
		FileURL:   presignedUrl,
		FileID:    attachment.Uuid,
		CreatedAt: attachment.Created.Time,
		Size:      int64(attachment.Size),
	}, "Attachment retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteAttachmentResponse represents the response for DeleteAttachment
type DeleteAttachmentResponse struct {
	FileURL   string    `json:"file_url"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

// DeleteAttachment deletes an attachment by its ID
// @Summary Delete an attachment by ID
// @Description Delete an attachment by its ID
// @Tags attachments
// @Produce json
// @Param id path string true "Attachment ID"
// @Success 200 {object} Response[DeleteAttachmentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /attachments/{id} [delete]
func (server *Server) DeleteAttachment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	attachment, err := server.store.GetAttachmentById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			server.logBusinessEvent(LogLevelError, "DeleteAttachment", "Failed to rollback transaction", zap.Error(rollbackErr))
		}
	}()
	qtx := server.store.WithTx(tx)

	attachment, err = qtx.DeleteAttachment(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.b2Client.Delete(ctx, attachment.Name)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteAttachment", "Failed to delete attachment from B2", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete attachment from B2: %w", err)))
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteAttachment", "Failed to commit transaction", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to commit transaction: %w", err)))
		return
	}

	res := SuccessResponse(DeleteAttachmentResponse{
		FileURL:   attachment.File,
		FileID:    attachment.Uuid,
		CreatedAt: attachment.Created.Time,
		Size:      int64(attachment.Size),
	}, "Attachment deleted successfully")
	ctx.JSON(http.StatusOK, res)

}

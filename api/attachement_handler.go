package api

import (
	"fmt"
	"net/http"

	db "maicare_go/db/sqlc"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
)

type UploadHandlerResponse struct {
	Message   string `json:"message"`
	FileURL   string `json:"file_url"`
	FileID    string `json:"file_id"`
	CreatedAt string `json:"created_at"`
}

func (server *Server) UploadHandler(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer file.Close()

	// Check if file is empty
	if header.Size == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("file cannot be empty")))
		return
	}

	filename := header.Filename

	err = server.b2Client.UploadToB2(ctx.Request.Context(), file, filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fileURL := fmt.Sprintf("%s/file/%s/%s",
		server.b2Client.Bucket.BaseURL(),
		server.b2Client.Bucket.Name(),
		filename)

	arg := db.CreateAttachmentParams{
		Name: filename,
		File: fileURL,
		Size: int32(header.Size),
		Tag:  util.StringPtr(""),
	}
	attachment, err := server.store.CreateAttachment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := UploadHandlerResponse{
		Message:   "File uploaded successfully",
		FileURL:   fileURL,
		FileID:    string(attachment.Uuid.Bytes[:]),
		CreatedAt: attachment.Created.Time.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, res)
}

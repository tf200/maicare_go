package api

import (
	"fmt"
	"io"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/util"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// IntakeFormUploadHandlerResponse represents a response from the intake form upload handler
type IntakeFormUploadHandlerResponse struct {
	FileURL   string    `json:"file_url"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

// @Summary Upload a file for intake form
// @Description Upload a file for intake form
// @Tags intake_form
// @Accept mpfd
// @Produce json
// @Param token path string true "Intake form token"
// @Param file formData file true "File to upload"
// @Success 201 {object} Response[IntakeFormUploadHandlerResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 413 {object} Response[any] "Request entity too large"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /intake_form/upload [post]
// @Security -
func (server *Server) IntakeFormUploadHandlerApi(ctx *gin.Context) {
	token := ctx.Param("token")

	dbToken, err := server.store.GetIntakeFormToken(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if dbToken.IsRevoked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if dbToken.ExpiresAt.Time.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxFileSize)
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			ctx.JSON(http.StatusRequestEntityTooLarge, errorResponse(fmt.Errorf("file size exceeds maximum limit of 10MB")))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Basic validations
	if err := bucket.ValidateFile(header, maxFileSize); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

	res := SuccessResponse(UploadHandlerResponse{
		FileURL:   fileURL,
		FileID:    attachment.Uuid,
		CreatedAt: attachment.Created.Time,
		Size:      int64(attachment.Size),
	}, "File uploaded successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GenerateIntakeFormTokenResponse represents a response from the generate intake form token handler
type GenerateIntakeFormTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// @Summary Generate an intake form token
// @Description Generate an intake form token
// @Tags intake_form
// @Produce json
// @Success 201 {object} Response[GenerateIntakeFormTokenResponse]
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /intake_form/token [post]
// @Security -
func (server *Server) GenerateIntakeFormToken(ctx *gin.Context) {
	arg := db.CreateIntakeFormTokenParams{
		Token:     uuid.New().String(),
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24), Valid: true},
	}

	token, err := server.store.CreateIntakeFormToken(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GenerateIntakeFormTokenResponse{
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt.Time,
	}, "Intake form token created successfully")
	ctx.JSON(http.StatusCreated, res)

}

// @Summary Verify an intake form token
// @Description Verify an intake form token
// @Tags intake_form
// @Produce json
// @Param token path string true "Intake form token"
// @Success 200 {object} Response[any]
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /intake_form/token/{token} [get]
// @Security -
func (server *Server) VerifyIntakeFormToken(ctx *gin.Context) {
	token := ctx.Param("token")

	dbToken, err := server.store.GetIntakeFormToken(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if dbToken.IsRevoked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if dbToken.ExpiresAt.Time.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Intake form token verified"})

}

// CreateIntakeFormRequest represents a request to create an intake form
type CreateIntakeFormRequest struct {
	IntakeFormToken            string      `json:"intake_form_token" binding:"required"`
	FirstName                  string      `json:"first_name" binding:"required"`
	LastName                   string      `json:"last_name" binding:"required"`
	DateOfBirth                time.Time   `json:"date_of_birth" binding:"required"`
	Phonenumber                string      `json:"phonenumber" binding:"required"`
	Gender                     string      `json:"gender" binding:"required"`
	PlaceOfBirth               string      `json:"place_of_birth" binding:"required"`
	RepresentativeFirstName    string      `json:"representative_first_name" binding:"required"`
	RepresentativeLastName     string      `json:"representative_last_name" binding:"required"`
	RepresentativePhoneNumber  string      `json:"representative_phone_number" binding:"required"`
	RepresentativeEmail        string      `json:"representative_email" binding:"required"`
	RepresentativeRelationship string      `json:"representative_relationship" binding:"required"`
	RepresentativeAddress      string      `json:"representative_address" binding:"required"`
	AttachementIds             []uuid.UUID `json:"attachement_ids"`
}

// CreateIntakeFormResponse represents a response from the create intake form handler
type CreateIntakeFormResponse struct {
	ID                         int64       `json:"id"`
	IntakeFormToken            string      `json:"intake_form_token"`
	FirstName                  string      `json:"first_name"`
	LastName                   string      `json:"last_name"`
	DateOfBirth                time.Time   `json:"date_of_birth"`
	Phonenumber                string      `json:"phonenumber"`
	Gender                     string      `json:"gender"`
	PlaceOfBirth               string      `json:"place_of_birth"`
	RepresentativeFirstName    string      `json:"representative_first_name"`
	RepresentativeLastName     string      `json:"representative_last_name"`
	RepresentativePhoneNumber  string      `json:"representative_phone_number"`
	RepresentativeEmail        string      `json:"representative_email"`
	RepresentativeRelationship string      `json:"representative_relationship"`
	RepresentativeAddress      string      `json:"representative_address"`
	AttachementIds             []uuid.UUID `json:"attachement_ids"`
}

// @Summary Create an intake form
// @Description Create an intake form
// @Tags intake_form
// @Accept json
// @Produce json
// @Param token path string true "Intake form token"
// @Param request body CreateIntakeFormRequest true "Intake form request"
// @Success 201 {object} Response[CreateIntakeFormResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /intake_form [post]
// @Security -
func (server *Server) CreateIntakeFormApi(ctx *gin.Context) {
	token := ctx.Param("token")

	dbToken, err := server.store.GetIntakeFormToken(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if dbToken.IsRevoked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if dbToken.ExpiresAt.Time.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	var req CreateIntakeFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := server.store.WithTx(tx)

	arg := db.CreateIntakeFormParams{
		IntakeFormToken:            req.IntakeFormToken,
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		DateOfBirth:                pgtype.Date{Time: req.DateOfBirth, Valid: true},
		PhoneNumber:                req.Phonenumber,
		Gender:                     req.Gender,
		PlaceOfBirth:               req.PlaceOfBirth,
		RepresentativeFirstName:    req.RepresentativeFirstName,
		RepresentativeLastName:     req.RepresentativeLastName,
		RepresentativePhoneNumber:  req.RepresentativePhoneNumber,
		RepresentativeEmail:        req.RepresentativeEmail,
		RepresentativeRelationship: req.RepresentativeRelationship,
		RepresentativeAddress:      req.RepresentativeAddress,
		AttachementIds:             req.AttachementIds,
	}

	form, err := qtx.CreateIntakeForm(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = qtx.RevokedIntakeFormToken(ctx, req.IntakeFormToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, attachmentID := range req.AttachementIds {
		_, err = qtx.SetAttachmentAsUsedorUnused(ctx, db.SetAttachmentAsUsedorUnusedParams{
			Uuid:   attachmentID,
			IsUsed: true,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateIntakeFormResponse{
		ID:                         form.ID,
		IntakeFormToken:            form.IntakeFormToken,
		FirstName:                  form.FirstName,
		LastName:                   form.LastName,
		DateOfBirth:                form.DateOfBirth.Time,
		Phonenumber:                form.PhoneNumber,
		Gender:                     form.Gender,
		PlaceOfBirth:               form.PlaceOfBirth,
		RepresentativeFirstName:    form.RepresentativeFirstName,
		RepresentativeLastName:     form.RepresentativeLastName,
		RepresentativePhoneNumber:  form.RepresentativePhoneNumber,
		RepresentativeEmail:        form.RepresentativeEmail,
		RepresentativeRelationship: form.RepresentativeRelationship,
		RepresentativeAddress:      form.RepresentativeAddress,
		AttachementIds:             form.AttachementIds,
	}, "Intake form created successfully")

	ctx.JSON(http.StatusCreated, res)
}

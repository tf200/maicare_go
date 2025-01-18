package api

import (
	"encoding/json"
	db "maicare_go/db/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateClientDetailsRequest represents a request to create a new client
type CreateClientDetailsRequest struct {
	FirstName             string      `json:"first_name" binding:"required"`
	LastName              string      `json:"last_name" binding:"required"`
	Email                 string      `json:"email" binding:"required,email"`
	Organisation          *string     `json:"organisation" binding:"required"`
	LocationID            *int64      `json:"location" binding:"required"`
	LegalMeasure          *string     `json:"legal_measure"`
	Birthplace            *string     `json:"birthplace" binding:"required"`
	Departement           *string     `json:"departement" binding:"required"`
	Gender                string      `json:"gender" binding:"required"`
	Filenumber            string      `json:"filenumber" binding:"required"`
	DateOfBirth           string      `json:"date_of_birth" binding:"required" time_format:"2006-01-02"`
	PhoneNumber           *string     `json:"phone_number" binding:"required"`
	Infix                 *string     `json:"infix"`
	Source                *string     `json:"source" binding:"required"`
	Bsn                   *string     `json:"bsn"`
	Addresses             []Address   `json:"addresses"`
	IdentityAttachmentIds []uuid.UUID `json:"identity_attachment_ids"`
}

// Address represents a client address
type Address struct {
	BelongsTo   *string `json:"belongs_to"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	ZipCode     *string `json:"zip_code"`
	PhoneNumber *string `json:"phone_number"`
}

// CreateClientDetailsResponse represents a response to a create client request
type CreateClientDetailsResponse struct {
	ID                    int64              `json:"id"`
	FirstName             string             `json:"first_name"`
	LastName              string             `json:"last_name"`
	DateOfBirth           pgtype.Date        `json:"date_of_birth"`
	Identity              bool               `json:"identity"`
	Status                *string            `json:"status"`
	Bsn                   *string            `json:"bsn"`
	Source                *string            `json:"source"`
	Birthplace            *string            `json:"birthplace"`
	Email                 string             `json:"email"`
	PhoneNumber           *string            `json:"phone_number"`
	Organisation          *string            `json:"organisation"`
	Departement           *string            `json:"departement"`
	Gender                string             `json:"gender"`
	Filenumber            string             `json:"filenumber"`
	ProfilePicture        *string            `json:"profile_picture"`
	Infix                 *string            `json:"infix"`
	Created               pgtype.Timestamptz `json:"created"`
	SenderID              *int64             `json:"sender_id"`
	LocationID            *int64             `json:"location_id"`
	IdentityAttachmentIds []uuid.UUID        `json:"identity_attachment_ids"`
	DepartureReason       *string            `json:"departure_reason"`
	DepartureReport       *string            `json:"departure_report"`
	Addresses             []Address          `json:"addresses"`
	LegalMeasure          *string            `json:"legal_measure"`
	HasUntakenMedications bool               `json:"has_untaken_medications"`
}

// CreateClientApi creates a new client
// @Summary Create a new client
// @Tags clients
// @Accept json
// @Produce json
// @Param request body CreateClientDetailsRequest true "Client details"
// @Success 200 {object} Response[CreateClientDetailsResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients [post]
func (server *Server) CreateClientApi(ctx *gin.Context) {
	var req CreateClientDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	parsedDateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	identityAttachmentIdsJSON, err := json.Marshal(req.IdentityAttachmentIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	AddressesJSON, err := json.Marshal(req.Addresses)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	client, err := server.store.CreateClientDetailsTx(ctx, db.CreateClientDetailsTxParams{
		CreateClientParams: db.CreateClientDetailsParams{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			DateOfBirth: pgtype.Date{Time: parsedDateOfBirth, Valid: true},
			Identity:    true, // needs to be checked
			//Status:                nil,  // needs to be checked
			Bsn:          req.Bsn,
			Source:       req.Source,
			Birthplace:   req.Birthplace,
			Email:        req.Email,
			PhoneNumber:  req.PhoneNumber,
			Organisation: req.Organisation,
			Departement:  req.Departement,
			Infix:        req.Infix,
			Gender:       req.Gender,
			Filenumber:   req.Filenumber,
			//ProfilePicture:        nil, // needs to be checked
			//SenderID:              nil, // needs to be checked
			LocationID:            req.LocationID,
			IdentityAttachmentIds: identityAttachmentIdsJSON,
			//DepartureReason:       nil, // needs to be checked
			//DepartureReport:       nil, // needs to be checked
			//GpsPosition:           nil,  needs to be checked
			//MaturityDomains:       nil, // needs to be checked
			Addresses:    AddressesJSON,
			LegalMeasure: req.LegalMeasure,
			//HasUntakenMedications: false, // needs to be checked
		},
		IdentityAttachments: req.IdentityAttachmentIds,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var addresses []Address
	err = json.Unmarshal(client.Client.Addresses, &addresses)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var identityAttachmentIds []uuid.UUID
	err = json.Unmarshal(client.Client.IdentityAttachmentIds, &identityAttachmentIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateClientDetailsResponse{
		ID:                    client.Client.ID,
		FirstName:             client.Client.FirstName,
		LastName:              client.Client.LastName,
		DateOfBirth:           client.Client.DateOfBirth,
		Identity:              client.Client.Identity,
		Status:                client.Client.Status,
		Bsn:                   client.Client.Bsn,
		Source:                client.Client.Source,
		Birthplace:            client.Client.Birthplace,
		Email:                 client.Client.Email,
		PhoneNumber:           client.Client.PhoneNumber,
		Organisation:          client.Client.Organisation,
		Departement:           client.Client.Departement,
		Gender:                client.Client.Gender,
		Filenumber:            client.Client.Filenumber,
		ProfilePicture:        client.Client.ProfilePicture,
		Infix:                 client.Client.Infix,
		Created:               client.Client.Created,
		SenderID:              client.Client.SenderID,
		LocationID:            client.Client.LocationID,
		IdentityAttachmentIds: identityAttachmentIds,
		DepartureReason:       client.Client.DepartureReason,
		DepartureReport:       client.Client.DepartureReport,
		Addresses:             addresses,
		LegalMeasure:          client.Client.LegalMeasure,
		HasUntakenMedications: client.Client.HasUntakenMedications,
	}, "Client created successfully")
	ctx.JSON(http.StatusOK, res)

}

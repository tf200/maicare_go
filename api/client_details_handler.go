package api

import (
	"encoding/json"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
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
	LocationID            *int64      `json:"location_id" binding:"required"`
	LegalMeasure          *string     `json:"legal_measure"`
	Birthplace            *string     `json:"birthplace" binding:"required"`
	Departement           *string     `json:"departement" binding:"required"`
	Gender                string      `json:"gender" binding:"required"`
	Filenumber            string      `json:"filenumber" binding:"required"`
	DateOfBirth           string      `json:"date_of_birth" binding:"required" time_format:"2006-01-02"`
	PhoneNumber           *string     `json:"phone_number" binding:"required"`
	SenderID              int64       `json:"sender_id" binding:"required"`
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
	ID                    int64       `json:"id"`
	FirstName             string      `json:"first_name"`
	LastName              string      `json:"last_name"`
	DateOfBirth           time.Time   `json:"date_of_birth"`
	Identity              bool        `json:"identity"`
	Status                *string     `json:"status"`
	Bsn                   *string     `json:"bsn"`
	Source                *string     `json:"source"`
	Birthplace            *string     `json:"birthplace"`
	Email                 string      `json:"email"`
	PhoneNumber           *string     `json:"phone_number"`
	Organisation          *string     `json:"organisation"`
	Departement           *string     `json:"departement"`
	Gender                string      `json:"gender"`
	Filenumber            string      `json:"filenumber"`
	ProfilePicture        *string     `json:"profile_picture"`
	Infix                 *string     `json:"infix"`
	Created               time.Time   `json:"created"`
	SenderID              int64       `json:"sender_id"`
	LocationID            *int64      `json:"location_id"`
	IdentityAttachmentIds []uuid.UUID `json:"identity_attachment_ids"`
	DepartureReason       *string     `json:"departure_reason"`
	DepartureReport       *string     `json:"departure_report"`
	Addresses             []Address   `json:"addresses"`
	LegalMeasure          *string     `json:"legal_measure"`
	HasUntakenMedications bool        `json:"has_untaken_medications"`
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
			SenderID:              req.SenderID, // needs to be checked
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
		DateOfBirth:           client.Client.DateOfBirth.Time,
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
		Created:               client.Client.CreatedAt.Time,
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

// ListClientsApiParams represents a request to list clients
type ListClientsApiParams struct {
	pagination.Request
	Status     *string `form:"status"`
	LocationID *int64  `form:"location_id"`
	Search     *string `form:"search"`
}

// ListClientsApiResponse represents a response to a list clients request
type ListClientsApiResponse struct {
	ID                    int64     `json:"id"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	DateOfBirth           time.Time `json:"date_of_birth"`
	Identity              bool      `json:"identity"`
	Status                *string   `json:"status"`
	Bsn                   *string   `json:"bsn"`
	Source                *string   `json:"source"`
	Birthplace            *string   `json:"birthplace"`
	Email                 string    `json:"email"`
	PhoneNumber           *string   `json:"phone_number"`
	Organisation          *string   `json:"organisation"`
	Departement           *string   `json:"departement"`
	Gender                string    `json:"gender"`
	Filenumber            string    `json:"filenumber"`
	ProfilePicture        *string   `json:"profile_picture"`
	Infix                 *string   `json:"infix"`
	CreatedAt             time.Time `json:"created_at"`
	SenderID              int64     `json:"sender_id"`
	LocationID            *int64    `json:"location_id"`
	DepartureReason       *string   `json:"departure_reason"`
	DepartureReport       *string   `json:"departure_report"`
	Addresses             []Address `json:"addresses"`
	LegalMeasure          *string   `json:"legal_measure"`
	HasUntakenMedications bool      `json:"has_untaken_medications"`
}

// ListClientsApi lists clients
// @Summary List clients
// @Tags clients
// @Produce json
// @Param status query string false "Client status"
// @Param location_id query int false "Location ID"
// @Param search query string false "Search query"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListClientsApiResponse]]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients [get]
func (server *Server) ListClientsApi(ctx *gin.Context) {
	var req ListClientsApiParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	params := req.GetParams()

	clients, err := server.store.ListClientDetails(ctx, db.ListClientDetailsParams{
		Status: req.Status,
		Search: req.Search,
		Offset: params.Offset,
		Limit:  params.Limit,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalCount := clients[0].TotalCount

	clientList := make([]ListClientsApiResponse, len(clients))
	for i, client := range clients {
		var addresses []Address
		err = json.Unmarshal(client.Addresses, &addresses)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		clientList[i] = ListClientsApiResponse{
			ID:                    client.ID,
			FirstName:             client.FirstName,
			LastName:              client.LastName,
			DateOfBirth:           client.DateOfBirth.Time,
			Identity:              client.Identity,
			Status:                client.Status,
			Bsn:                   client.Bsn,
			Source:                client.Source,
			Birthplace:            client.Birthplace,
			Email:                 client.Email,
			PhoneNumber:           client.PhoneNumber,
			Organisation:          client.Organisation,
			Departement:           client.Departement,
			Gender:                client.Gender,
			Filenumber:            client.Filenumber,
			ProfilePicture:        client.ProfilePicture,
			Infix:                 client.Infix,
			CreatedAt:             client.CreatedAt.Time,
			SenderID:              client.SenderID,
			LocationID:            client.LocationID,
			DepartureReason:       client.DepartureReason,
			DepartureReport:       client.DepartureReport,
			Addresses:             addresses,
			LegalMeasure:          client.LegalMeasure,
			HasUntakenMedications: client.HasUntakenMedications,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, clientList, totalCount)

	res := SuccessResponse(pag, "Clients fetched successfully")
	ctx.JSON(http.StatusOK, res)

}

// GetClientApiResponse represents a response to a get client request
type GetClientApiResponse struct {
	ID                    int64       `json:"id"`
	FirstName             string      `json:"first_name"`
	LastName              string      `json:"last_name"`
	DateOfBirth           time.Time   `json:"date_of_birth"`
	Identity              bool        `json:"identity"`
	Status                *string     `json:"status"`
	Bsn                   *string     `json:"bsn"`
	Source                *string     `json:"source"`
	Birthplace            *string     `json:"birthplace"`
	Email                 string      `json:"email"`
	PhoneNumber           *string     `json:"phone_number"`
	Organisation          *string     `json:"organisation"`
	Departement           *string     `json:"departement"`
	Gender                string      `json:"gender"`
	Filenumber            string      `json:"filenumber"`
	ProfilePicture        *string     `json:"profile_picture"`
	Infix                 *string     `json:"infix"`
	CreatedAt             time.Time   `json:"created_at"`
	SenderID              int64       `json:"sender_id"`
	LocationID            *int64      `json:"location_id"`
	DepartureReason       *string     `json:"departure_reason"`
	DepartureReport       *string     `json:"departure_report"`
	Addresses             []Address   `json:"addresses"`
	IdentityAttachmentIds []uuid.UUID `json:"identity_attachment_ids"`
	LegalMeasure          *string     `json:"legal_measure"`
	HasUntakenMedications bool        `json:"has_untaken_medications"`
}

// GetClientApi gets a client
// @Summary Get a client
// @Tags clients
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[GetClientApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id} [get]
func (server *Server) GetClientApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	client, err := server.store.GetClientDetails(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var addresses []Address
	err = json.Unmarshal(client.Addresses, &addresses)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var identityAttachmentIds []uuid.UUID
	err = json.Unmarshal(client.IdentityAttachmentIds, &identityAttachmentIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetClientApiResponse{
		ID:                    client.ID,
		FirstName:             client.FirstName,
		LastName:              client.LastName,
		DateOfBirth:           client.DateOfBirth.Time,
		Identity:              client.Identity,
		Status:                client.Status,
		Bsn:                   client.Bsn,
		Source:                client.Source,
		Birthplace:            client.Birthplace,
		Email:                 client.Email,
		PhoneNumber:           client.PhoneNumber,
		Organisation:          client.Organisation,
		Departement:           client.Departement,
		Gender:                client.Gender,
		Filenumber:            client.Filenumber,
		ProfilePicture:        client.ProfilePicture,
		Infix:                 client.Infix,
		CreatedAt:             client.CreatedAt.Time,
		SenderID:              client.SenderID,
		LocationID:            client.LocationID,
		DepartureReason:       client.DepartureReason,
		DepartureReport:       client.DepartureReport,
		Addresses:             addresses,
		IdentityAttachmentIds: identityAttachmentIds,
		LegalMeasure:          client.LegalMeasure,
		HasUntakenMedications: client.HasUntakenMedications,
	}, "Client fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// SetClientProfilePictureRequest represents a request to update a client
type SetClientProfilePictureRequest struct {
	AttachmentID uuid.UUID `json:"attachement_id" binding:"required"`
}

// SetClientProfilePictureResponse represents a response to a set client profile picture request
type SetClientProfilePictureResponse struct {
	ID             int64   `json:"id"`
	ProfilePicture *string `json:"profile_picture"`
}

// SetClientProfilePictureApi sets a client profile picture
// @Summary Set a client profile picture
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body SetClientProfilePictureRequest true "Client profile picture"
// @Success 200 {object} Response[SetClientProfilePictureResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/profile_picture [put]
func (server *Server) SetClientProfilePictureApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req SetClientProfilePictureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.SetClientProfilePictureTxParams{
		ClientID:     clientID,
		AttachmentID: req.AttachmentID,
	}
	client, err := server.store.SetClientProfilePictureTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(SetClientProfilePictureResponse{
		ID:             client.User.ID,
		ProfilePicture: client.User.ProfilePicture,
	}, "Profile picture set successfully")
	ctx.JSON(http.StatusOK, res)

}

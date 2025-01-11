package api

import (
	"encoding/json"
	"log"
	"net/http"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
)

// Contact represents a contact information.
type Contact struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"email"`
	PhoneNumber *string `json:"phone_number"`
}

// CreateSenderRequest represents a request to create a new sender.
type CreateSenderRequest struct {
	Types        string    `json:"types" binding:"required,oneof=main_provider local_authority particular_party healthcare_institution"`
	Name         string    `json:"name" binding:"required"`
	Address      *string   `json:"address"`
	PostalCode   *string   `json:"postal_code"`
	Place        *string   `json:"place"`
	Land         *string   `json:"land"`
	KVKNumber    *string   `json:"KVKnumber"`
	BTWNumber    *string   `json:"BTWnumber"`
	PhoneNumber  *string   `json:"phone_number"`
	ClientNumber *string   `json:"client_number"`
	Contacts     []Contact `json:"contacts" binding:"dive"`
}

// CreateSenderResponse represents a response to a request to create a new sender.
type CreateSenderResponse struct {
	ID           int64     `json:"id"`
	Types        string    `json:"types"`
	Name         string    `json:"name"`
	Address      *string   `json:"address"`
	PostalCode   *string   `json:"postal_code"`
	Place        *string   `json:"place"`
	Land         *string   `json:"land"`
	KVKNumber    *string   `json:"KVKnumber"`
	BTWNumber    *string   `json:"BTWnumber"`
	PhoneNumber  *string   `json:"phone_number"`
	ClientNumber *string   `json:"client_number"`
	Contacts     []Contact `json:"contacts"`
}

// CreateSenderApi creates a new sender.
// @Summary Create a new sender
// @Description Create a new sender
// @Tags senders
// @Accept json
// @Produce json
// @Param request body CreateSenderRequest true "Sender data"
// @Success 201 {object} Response[CreateSenderResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders [post]
func (server *Server) CreateSenderApi(ctx *gin.Context) {
	var req CreateSenderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	contacts, err := json.Marshal(req.Contacts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sender, err := server.store.CreateSender(ctx, db.CreateSenderParams{
		Types:        req.Types,
		Name:         req.Name,
		Address:      req.Address,
		PostalCode:   req.PostalCode,
		Place:        req.Place,
		Land:         req.Land,
		Kvknumber:    req.KVKNumber,
		Btwnumber:    req.BTWNumber,
		PhoneNumber:  req.PhoneNumber,
		ClientNumber: req.ClientNumber,
		Contacts:     contacts,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := SuccessResponse(CreateSenderResponse{
		ID:           sender.ID,
		Types:        sender.Types,
		Name:         sender.Name,
		Address:      sender.Address,
		PostalCode:   sender.PostalCode,
		Place:        sender.Place,
		Land:         sender.Land,
		KVKNumber:    sender.Kvknumber,
		BTWNumber:    sender.Btwnumber,
		PhoneNumber:  sender.PhoneNumber,
		ClientNumber: sender.ClientNumber,
		Contacts:     req.Contacts,
	}, "Sender created successfully")

	ctx.JSON(http.StatusCreated, rsp)
}

// GetSenderRequest represents a request to get a sender by ID.
type ListSendersRequest struct {
	pagination.Request
	IncludeArchived *bool   `form:"include_archived"`
	Search          *string `form:"search"`
}

// GetSenderResponse represents a response to a request to get a sender by ID.
type ListSendersResponse struct {
	ID           int64     `json:"id"`
	Types        string    `json:"types"`
	Name         string    `json:"name"`
	Address      *string   `json:"address"`
	PostalCode   *string   `json:"postal_code"`
	Place        *string   `json:"place"`
	Land         *string   `json:"land"`
	KVKNumber    *string   `json:"KVKnumber"`
	BTWNumber    *string   `json:"BTWnumber"`
	PhoneNumber  *string   `json:"phone_number"`
	ClientNumber *string   `json:"client_number"`
	Contacts     []Contact `json:"contacts"`
}

// ListSendersAPI returns a list of senders.
// @Summary List senders
// @Description List senders
// @Tags senders
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param search query string false "Search"
// @Param include_archived query bool false "Include archived"
// @Success 200 {object} Response[pagination.Response[ListSendersResponse]]
// @Failure 400,404,500 {object} Response[any]
func (server *Server) ListSendersAPI(ctx *gin.Context) {
	var req ListSendersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	arg := db.ListSendersParams{
		Limit:           params.Limit,
		Offset:          params.Offset,
		Search:          req.Search,
		IncludeArchived: req.IncludeArchived,
	}
	log.Printf("arg: %v", arg)

	// Fetch senders from the database
	senders, err := server.store.ListSenders(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	responseSenders := make([]ListSendersResponse, len(senders))
	for i, sender := range senders {
		contacts := make([]Contact, 0)
		if err := json.Unmarshal(sender.Contacts, &contacts); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		responseSenders[i] = ListSendersResponse{
			ID:           sender.ID,
			Types:        sender.Types,
			Name:         sender.Name,
			Address:      sender.Address,
			PostalCode:   sender.PostalCode,
			Place:        sender.Place,
			Land:         sender.Land,
			KVKNumber:    sender.Kvknumber,
			BTWNumber:    sender.Btwnumber,
			PhoneNumber:  sender.PhoneNumber,
			ClientNumber: sender.ClientNumber,
			Contacts:     contacts,
		}
	}

	totalCount, err := server.store.CountSenders(ctx, req.IncludeArchived)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := pagination.NewResponse(ctx, req.Request, responseSenders, totalCount)
	res := SuccessResponse(response, "Senders retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

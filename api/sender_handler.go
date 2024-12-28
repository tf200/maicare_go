package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
)

type Contact struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"email"`
	PhoneNumber *string `json:"phone_number"`
}

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

	rsp := CreateSenderResponse{
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
	}

	ctx.JSON(http.StatusCreated, rsp)
}

type ListSendersRequest struct {
	pagination.Request
	IncludeArchived *bool   `form:"include_archived"`
	Search          *string `form:"search"`
}

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

	// Fetch senders from the database
	senders, err := server.store.ListSenders(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var responseSenders []ListSendersResponse
	for _, sender := range senders {
		var contacts []Contact

		if sender.Contacts != nil {
			if err := json.Unmarshal(sender.Contacts, &contacts); err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to parse contacts: %w", err)))
				return
			}
		}
		responseSenders = append(responseSenders, ListSendersResponse{
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
		})

	}

	totalCount, err := server.store.CountSenders(ctx, req.IncludeArchived)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := pagination.NewResponse(ctx, req.Request, responseSenders, totalCount)

	ctx.JSON(http.StatusOK, response)
}

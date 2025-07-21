package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
)

// Contact represents a contact information.
type SenderContact struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"email"`
	PhoneNumber *string `json:"phone_number"`
}

// CreateSenderRequest represents a request to create a new sender.
type CreateSenderRequest struct {
	Types        string          `json:"types" binding:"required,oneof=main_provider local_authority particular_party healthcare_institution"`
	Name         string          `json:"name" binding:"required"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	KVKNumber    *string         `json:"KVKnumber"`
	BTWNumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	Contacts     []SenderContact `json:"contacts" binding:"dive"`
}

// CreateSenderResponse represents a response to a request to create a new sender.
type CreateSenderResponse struct {
	ID           int64           `json:"id"`
	Types        string          `json:"types"`
	Name         string          `json:"name"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	KVKNumber    *string         `json:"KVKnumber"`
	BTWNumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	Contacts     []SenderContact `json:"contacts"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
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
	contactsParam, err := json.Marshal(req.Contacts)
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
		Contacts:     contactsParam,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	contactsResp := make([]SenderContact, 0)
	if err := json.Unmarshal(sender.Contacts, &contactsResp); err != nil {
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
		Contacts:     contactsResp,
		CreatedAt:    sender.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:    sender.UpdatedAt.Time.Format(time.RFC3339),
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
	ID           int64           `json:"id"`
	Types        string          `json:"types"`
	Name         string          `json:"name"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	KVKNumber    *string         `json:"KVKnumber"`
	BTWNumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	Contacts     []SenderContact `json:"contacts"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

// ListSendersAPI returns a list of senders.
// @Summary List senders
// @Description List senders
// @Tags senders
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search"
// @Param include_archived query bool false "Include archived"
// @Success 200 {object} Response[pagination.Response[ListSendersResponse]]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders [get]
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
		contacts := make([]SenderContact, 0)
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
			CreatedAt:    sender.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:    sender.UpdatedAt.Time.Format(time.RFC3339),
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

type TemplateItem struct {
	ID           int64  `json:"id"`
	ItemTag      string `json:"item_tag"`
	Description  string `json:"description"`
	SourceTable  string `json:"source_table"`
	SourceColumn string `json:"source_column"`
}

// GetSenderByIdResponse represents a response to a request to get a sender by ID.
type GetSenderByIdResponse struct {
	ID                   int64           `json:"id"`
	Types                string          `json:"types"`
	Name                 string          `json:"name"`
	Address              *string         `json:"address"`
	PostalCode           *string         `json:"postal_code"`
	Place                *string         `json:"place"`
	Land                 *string         `json:"land"`
	Kvknumber            *string         `json:"KVKnumber"`
	Btwnumber            *string         `json:"BTWnumber"`
	PhoneNumber          *string         `json:"phone_number"`
	ClientNumber         *string         `json:"client_number"`
	EmailAddress         *string         `json:"email_address"`
	Contacts             []SenderContact `json:"contacts"`
	InvoiceTemplateItems []TemplateItem  `json:"invoice_template_items"`
	IsArchived           bool            `json:"is_archived"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

// GetSenderAPI returns a sender by ID.
// @Summary Get a sender
// @Description Get a sender
// @Tags senders
// @Produce json
// @Param id path int true "Sender ID"
// @Success 200 {object} Response[GetSenderByIdResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders/{id} [get]
func (server *Server) GetSenderByIdAPI(ctx *gin.Context) {
	id := ctx.Param("id")
	senderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sender, err := server.store.GetSenderById(ctx, senderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var contactsResp []SenderContact
	if err := json.Unmarshal(sender.Contacts, &contactsResp); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	tmplids, err := server.store.GetTemplateItemsBySourceTable(ctx, sender.InvoiceTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var invoiceTemplateItems []TemplateItem
	if len(tmplids) != 0 {
		for _, tmpl := range tmplids {
			invoiceTemplateItems = append(invoiceTemplateItems, TemplateItem{
				ID:           tmpl.ID,
				ItemTag:      tmpl.ItemTag,
				Description:  tmpl.Description,
				SourceTable:  tmpl.SourceTable,
				SourceColumn: tmpl.SourceColumn,
			})
		}
	}

	res := SuccessResponse(GetSenderByIdResponse{
		ID:                   sender.ID,
		Types:                sender.Types,
		Name:                 sender.Name,
		Address:              sender.Address,
		PostalCode:           sender.PostalCode,
		Place:                sender.Place,
		Land:                 sender.Land,
		Kvknumber:            sender.Kvknumber,
		Btwnumber:            sender.Btwnumber,
		PhoneNumber:          sender.PhoneNumber,
		ClientNumber:         sender.ClientNumber,
		EmailAddress:         sender.EmailAddress,
		Contacts:             contactsResp,
		IsArchived:           sender.IsArchived,
		InvoiceTemplateItems: invoiceTemplateItems,
		CreatedAt:            sender.CreatedAt.Time,
		UpdatedAt:            sender.UpdatedAt.Time,
	}, "Sender retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// UpdateSenderRequest represents a request to update a sender.
type UpdateSenderRequest struct {
	Name         *string         `json:"name"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	Kvknumber    *string         `json:"KVKnumber"`
	Btwnumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	EmailAddress *string         `json:"email_address"`
	Contacts     []SenderContact `json:"contacts"`
	IsArchived   *bool           `json:"is_archived"`
	Types        *string         `json:"types" binding:"omitempty,oneof=main_provider local_authority particular_party healthcare_institution"`
}

// UpdateSenderResponse represents a response to a request to update a sender.
type UpdateSenderResponse struct {
	ID           int64           `json:"id"`
	Types        string          `json:"types"`
	Name         string          `json:"name"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	Kvknumber    *string         `json:"KVKnumber"`
	Btwnumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	EmailAddress *string         `json:"email_address"`
	Contacts     []SenderContact `json:"contacts"`
	IsArchived   bool            `json:"is_archived"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

// @Summary Update a sender
// @Description Update a sender
// @Tags senders
// @Accept json
// @Produce json
// @Param id path int true "Sender ID"
// @Param request body UpdateSenderRequest true "Sender data"
// @Success 200 {object} Response[UpdateSenderResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders/{id} [put]
func (server *Server) UpdateSenderApi(ctx *gin.Context) {
	id := ctx.Param("id")
	senderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateSenderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Initialize update params with the ID
	params := db.UpdateSenderParams{
		ID:           senderID,
		Name:         req.Name,
		Address:      req.Address,
		PostalCode:   req.PostalCode,
		Place:        req.Place,
		Land:         req.Land,
		Kvknumber:    req.Kvknumber,
		Btwnumber:    req.Btwnumber,
		PhoneNumber:  req.PhoneNumber,
		ClientNumber: req.ClientNumber,
		EmailAddress: req.EmailAddress,
		IsArchived:   req.IsArchived,
		Types:        req.Types,
	}

	// Only include contacts in the update if it's provided in the request
	if req.Contacts != nil {
		contactsParam, err := json.Marshal(req.Contacts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		params.Contacts = contactsParam
	}

	updatedSender, err := server.store.UpdateSender(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var contactsResp []SenderContact
	if err := json.Unmarshal(updatedSender.Contacts, &contactsResp); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateSenderResponse{
		ID:           updatedSender.ID,
		Types:        updatedSender.Types,
		Name:         updatedSender.Name,
		Address:      updatedSender.Address,
		PostalCode:   updatedSender.PostalCode,
		Place:        updatedSender.Place,
		Land:         updatedSender.Land,
		Kvknumber:    updatedSender.Kvknumber,
		Btwnumber:    updatedSender.Btwnumber,
		PhoneNumber:  updatedSender.PhoneNumber,
		ClientNumber: updatedSender.ClientNumber,
		EmailAddress: updatedSender.EmailAddress,
		Contacts:     contactsResp,
		IsArchived:   updatedSender.IsArchived,
		CreatedAt:    updatedSender.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:    updatedSender.UpdatedAt.Time.Format(time.RFC3339),
	}, "Sender updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete a sender
// @Description Delete a sender
// @Tags senders
// @Accept json
// @Produce json
// @Param id path int true "Sender ID"
// @Success 200 {object} Response[any]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders/{id} [delete]
func (server *Server) DeleteSenderApi(ctx *gin.Context) {
	id := ctx.Param("id")
	senderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteSender(ctx, senderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Sender deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateSenderInvoiceTemplateRequest represents a request to create a new sender invoice template.
type CreateSenderInvoiceTemplateRequest struct {
	InvoiceTemplateIDs []int64 `json:"invoice_template" binding:"required"`
}

// CreateSenderInvoiceTemplateApi handles the creation of a new sender invoice template.
// @Summary Create a new sender invoice template
// @Description Create a new sender invoice template
// @Tags senders
// @Accept json
// @Produce json
// @Param id path int true "Sender ID"
// @Param request body CreateSenderInvoiceTemplateRequest true "Invoice template IDs"
// @Success 201 {object} Response[any]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders/{id}/invoice_template [post]
func (server *Server) CreateSenderInvoiceTemplateApi(ctx *gin.Context) {
	id := ctx.Param("id")
	senderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateSenderInvoiceTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// verify that the ids exist in the template table
	tmplids, err := server.store.GetTemplateItemsByIds(ctx, req.InvoiceTemplateIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.CreateSenderInvoiceTemplate(ctx, db.CreateSenderInvoiceTemplateParams{
		ID:              senderID,
		InvoiceTemplate: tmplids,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Sender invoice template created successfully")
	ctx.JSON(http.StatusOK, res)
}

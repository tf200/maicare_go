package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/invoice"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// GenerateInvoiceRequest represents the request body for generating an invoice.
type GenerateInvoiceRequest struct {
	ClientID  int64     `json:"client_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

// GenerateInvoiceResponse represents the response body for generating an invoice.
type GenerateInvoiceResponse struct {
	ID              int64                    `json:"id"`
	InvoiceNumber   string                   `json:"invoice_number"`
	IssueDate       time.Time                `json:"issue_date"`
	DueDate         pgtype.Date              `json:"due_date"`
	Status          string                   `json:"status"`
	InvoiceDetails  []invoice.InvoiceDetails `json:"invoice_details"`
	TotalAmount     float64                  `json:"total_amount"`
	PdfAttachmentID *uuid.UUID               `json:"pdf_attachment_id"`
	ExtraContent    *string                  `json:"extra_content"`
	ClientID        int64                    `json:"client_id"`
	SenderID        *int64                   `json:"sender_id"`
	UpdatedAt       time.Time                `json:"updated_at"`
	CreatedAt       time.Time                `json:"created_at"`
}

// @Summary Generate Invoice
// @Description Generate an invoice based on the provided client ID and date range.
// @Tags Invoice
// @Accept json
// @Produce json
// @Param request body GenerateInvoiceRequest true "Generate Invoice Request"
// @Success 200 {object} Response[GenerateInvoiceResponse] "Successful response with invoice details"
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /invoices/generate [post]
func (server *Server) GenerateInvoiceApi(ctx *gin.Context) {
	var req GenerateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	invoiceData := invoice.InvoiceParams{
		ClientID:  req.ClientID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	clientSender, err := server.store.GetClientSender(ctx.Request.Context(), invoiceData.ClientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var invoiceResult *invoice.InvoiceData

	invoiceResult, warningCount, err := invoice.GenerateInvoice(server.store, invoiceData, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	invoiceDetailsBytes, err := json.Marshal(invoiceResult.InvoiceDetails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	invoice, err := server.store.CreateInvoice(ctx.Request.Context(), db.CreateInvoiceParams{
		ClientID:       invoiceResult.ClientID,
		SenderID:       &clientSender.ID,
		DueDate:        pgtype.Date{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		TotalAmount:    invoiceResult.TotalAmount,
		InvoiceDetails: invoiceDetailsBytes,
		InvoiceNumber:  invoice.GenerateInvoiceNumber(invoiceResult.ClientID, time.Now()),
		ExtraContent:   nil, // Assuming no extra content for simplicity
		WarningCount:   int32(warningCount),
		IssueDate:      pgtype.Date{Time: invoiceResult.InvoiceDate, Valid: true},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := GenerateInvoiceResponse{
		ID:              invoice.ID,
		InvoiceNumber:   invoice.InvoiceNumber,
		IssueDate:       invoice.IssueDate.Time,
		DueDate:         invoice.DueDate,
		Status:          invoice.Status,
		InvoiceDetails:  invoiceResult.InvoiceDetails,
		TotalAmount:     invoice.TotalAmount,
		PdfAttachmentID: invoice.PdfAttachmentID,
		ExtraContent:    invoice.ExtraContent,
		ClientID:        invoice.ClientID,
		SenderID:        invoice.SenderID,
		UpdatedAt:       invoice.UpdatedAt.Time,
		CreatedAt:       invoice.CreatedAt.Time,
	}
	res := SuccessResponse(response, "Invoice generated successfully")
	ctx.JSON(http.StatusOK, res)

}

// ListInvoicesRequest represents the request parameters for listing invoices.
type ListInvoicesRequest struct {
	ClientID  *int64    `form:"client_id"`
	SenderID  *int64    `form:"sender_id"`
	Status    *string   `form:"status" binding:"omitempty,oneof=outstanding partially_paid paid expired overpaid imported concept"`
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
	pagination.Request
}

// ListInvoicesResponse represents the response body for listing invoices.
type ListInvoicesResponse struct {
	ID              int64                    `json:"id"`
	InvoiceNumber   string                   `json:"invoice_number"`
	IssueDate       time.Time                `json:"issue_date"`
	DueDate         time.Time                `json:"due_date"`
	Status          string                   `json:"status"`
	InvoiceDetails  []invoice.InvoiceDetails `json:"invoice_details"`
	TotalAmount     float64                  `json:"total_amount"`
	PdfAttachmentID *uuid.UUID               `json:"pdf_attachment_id"`
	ExtraContent    *string                  `json:"extra_content"`
	ClientID        int64                    `json:"client_id"`
	SenderID        *int64                   `json:"sender_id"`
	UpdatedAt       time.Time                `json:"updated_at"`
	CreatedAt       time.Time                `json:"created_at"`
	SenderName      *string                  `json:"sender_name"`
	ClientFirstName string                   `json:"client_first_name"`
	ClientLastName  string                   `json:"client_last_name"`
	WarningCount    int32                    `json:"warning_count"`
}

// @Summary List Invoices
// @Description List invoices based on optional filters like client ID, sender ID, status, and date range.
// @Tags Invoice
// @Produce json
// @Param client_id query int64 false "Client ID"
// @Param sender_id query int64 false "Sender ID"
// @Param status query string false "Invoice status (outstanding, partially_paid, paid, expired, overpaid, imported, concept)"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number for pagination"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} Response[pagination.Response[ListInvoicesResponse]] "Successful response with paginated invoices"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices [get]
func (server *Server) ListInvoicesApi(ctx *gin.Context) {
	var req ListInvoicesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	invoices, err := server.store.ListInvoices(ctx.Request.Context(), db.ListInvoicesParams{
		ClientID:  req.ClientID,
		SenderID:  req.SenderID,
		Status:    req.Status,
		StartDate: pgtype.Date{Time: req.StartDate},
		EndDate:   pgtype.Date{Time: req.EndDate},
		Offset:    params.Offset,
		Limit:     params.Limit,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	invoicesRes := make([]ListInvoicesResponse, len(invoices))

	if len(invoices) == 0 {
		ctx.JSON(http.StatusOK, SuccessResponse(invoicesRes, "No invoices found"))
		return
	}

	for i, items := range invoices {
		var invoiceDetails []invoice.InvoiceDetails
		if err := json.Unmarshal(items.InvoiceDetails, &invoiceDetails); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		invoicesRes[i] = ListInvoicesResponse{
			ID:              items.ID,
			InvoiceNumber:   items.InvoiceNumber,
			IssueDate:       items.IssueDate.Time,
			DueDate:         items.DueDate.Time,
			Status:          items.Status,
			InvoiceDetails:  invoiceDetails,
			TotalAmount:     items.TotalAmount,
			PdfAttachmentID: items.PdfAttachmentID,
			ExtraContent:    items.ExtraContent,
			ClientID:        items.ClientID,
			SenderID:        items.SenderID,
			UpdatedAt:       items.UpdatedAt.Time,
			CreatedAt:       items.CreatedAt.Time,
			SenderName:      items.SenderName,
			ClientFirstName: items.ClientFirstName,
			ClientLastName:  items.ClientLastName,
			WarningCount:    items.WarningCount,
		}
	}
	pag := pagination.NewResponse(ctx, req.Request, invoicesRes, invoices[0].TotalCount)
	res := SuccessResponse(pag, "Invoices retrieved successfully")
	ctx.JSON(http.StatusOK, res)

}

// GetInvoiceByIDResponse represents the response body for getting an invoice by ID.
type GetInvoiceByIDResponse struct {
	ID              int64                    `json:"id"`
	InvoiceNumber   string                   `json:"invoice_number"`
	IssueDate       time.Time                `json:"issue_date"`
	DueDate         time.Time                `json:"due_date"`
	Status          string                   `json:"status"`
	InvoiceDetails  []invoice.InvoiceDetails `json:"invoice_details"`
	TotalAmount     float64                  `json:"total_amount"`
	PdfAttachmentID *uuid.UUID               `json:"pdf_attachment_id"`
	ExtraContent    *string                  `json:"extra_content"`
	ClientID        int64                    `json:"client_id"`
	SenderID        *int64                   `json:"sender_id"`
	UpdatedAt       time.Time                `json:"updated_at"`
	CreatedAt       time.Time                `json:"created_at"`
	SenderName      *string                  `json:"sender_name"`
	ClientFirstName string                   `json:"client_first_name"`
	ClientLastName  string                   `json:"client_last_name"`
}

// @Summary Get Invoice by ID
// @Description Get an invoice by its ID.
// @Tags Invoice
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Success 200 {object} Response[GetInvoiceByIDResponse] "Successful response with invoice details"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id} [get]
func (server *Server) GetInvoiceByIDApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	invoiceItem, err := server.store.GetInvoice(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var invoiceDetails []invoice.InvoiceDetails
	if err := json.Unmarshal(invoiceItem.InvoiceDetails, &invoiceDetails); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := GetInvoiceByIDResponse{
		ID:              invoiceItem.ID,
		InvoiceNumber:   invoiceItem.InvoiceNumber,
		IssueDate:       invoiceItem.IssueDate.Time,
		DueDate:         invoiceItem.DueDate.Time,
		Status:          invoiceItem.Status,
		InvoiceDetails:  invoiceDetails,
		TotalAmount:     invoiceItem.TotalAmount,
		PdfAttachmentID: invoiceItem.PdfAttachmentID,
		ExtraContent:    invoiceItem.ExtraContent,
		ClientID:        invoiceItem.ClientID,
		SenderID:        invoiceItem.SenderID,
		UpdatedAt:       invoiceItem.UpdatedAt.Time,
		CreatedAt:       invoiceItem.CreatedAt.Time,
		SenderName:      invoiceItem.SenderName,
		ClientFirstName: invoiceItem.ClientFirstName,
		ClientLastName:  invoiceItem.ClientLastName,
	}
	ctx.JSON(http.StatusOK, response)
}

// UpdateInvoiceRequest represents the request body for updating an invoice.
type UpdateInvoiceRequest struct {
	IssueDate      time.Time `json:"issue_date"`
	DueDate        time.Time `json:"due_date"`
	InvoiceDetails []byte    `json:"invoice_details"`
	TotalAmount    float64   `json:"total_amount"`
	ExtraContent   *string   `json:"extra_content"`
	Status         string    `json:"status"`
	WarningCount   int32     `json:"warning_count"`
}

// UpdateInvoiceResponse represents the response body for updating an invoice.
type UpdateInvoiceResponse struct {
	ID              int64                    `json:"id"`
	InvoiceNumber   string                   `json:"invoice_number"`
	IssueDate       time.Time                `json:"issue_date"`
	DueDate         pgtype.Date              `json:"due_date"`
	Status          string                   `json:"status"`
	InvoiceDetails  []invoice.InvoiceDetails `json:"invoice_details"`
	TotalAmount     float64                  `json:"total_amount"`
	PdfAttachmentID *uuid.UUID               `json:"pdf_attachment_id"`
	ExtraContent    *string                  `json:"extra_content"`
	ClientID        int64                    `json:"client_id"`
	SenderID        *int64                   `json:"sender_id"`
	UpdatedAt       time.Time                `json:"updated_at"`
	CreatedAt       time.Time                `json:"created_at"`
}

// @Summary Update Invoice
// @Description Update an existing invoice by its ID.
// @Tags Invoice
// @Accept json
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Param request body UpdateInvoiceRequest true "Update Invoice Request"
// @Success 200 {object} Response[UpdateInvoiceResponse] "Successful response with updated invoice details"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id} [put]
func (server *Server) UpdateInvoiceApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var invoiceDetails []invoice.InvoiceDetails
	if err := json.Unmarshal(req.InvoiceDetails, &invoiceDetails); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if invoice.VerifyTotalAmount(invoiceDetails, req.TotalAmount) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("total amount does not match the sum of invoice details")))
		return
	}

	updatedInvoice, err := server.store.UpdateInvoice(ctx.Request.Context(), db.UpdateInvoiceParams{
		ID:             invoiceID,
		IssueDate:      pgtype.Date{Time: req.IssueDate},
		DueDate:        pgtype.Date{Time: req.DueDate},
		InvoiceDetails: req.InvoiceDetails,
		TotalAmount:    req.TotalAmount,
		ExtraContent:   req.ExtraContent,
		Status:         req.Status,
		WarningCount:   req.WarningCount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = json.Unmarshal(updatedInvoice.InvoiceDetails, &invoiceDetails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := UpdateInvoiceResponse{
		ID:              updatedInvoice.ID,
		InvoiceNumber:   updatedInvoice.InvoiceNumber,
		IssueDate:       updatedInvoice.IssueDate.Time,
		DueDate:         updatedInvoice.DueDate,
		Status:          updatedInvoice.Status,
		InvoiceDetails:  invoiceDetails,
		TotalAmount:     updatedInvoice.TotalAmount,
		PdfAttachmentID: updatedInvoice.PdfAttachmentID,
		ExtraContent:    updatedInvoice.ExtraContent,
		ClientID:        updatedInvoice.ClientID,
		SenderID:        updatedInvoice.SenderID,
		UpdatedAt:       updatedInvoice.UpdatedAt.Time,
		CreatedAt:       updatedInvoice.CreatedAt.Time,
	}
	ctx.JSON(http.StatusOK, response)

}

// DeleteInvoiceApi handles the deletion of an invoice by its ID.
// @Summary Delete Invoice
// @Description Delete an invoice by its ID.
// @Tags Invoice
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Success 200 {object} Response[any] "Successful response indicating deletion"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id} [delete]
func (server *Server) DeleteInvoiceApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteInvoice(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse[any](nil, "Invoice deleted successfully"))
}

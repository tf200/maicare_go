package api

import (
	"database/sql"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/invoice"
	"maicare_go/pagination"
	"maicare_go/util"
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
		DueDate:         invoice.DueDate.Time,
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
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Invoice retrieved successfully"))

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

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := server.store.WithTx(tx)

	employeeID, err := qtx.GetEmployeeIDByUserID(ctx.Request.Context(), payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL myapp.current_employee_id = %d", employeeID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	updatedInvoice, err := qtx.UpdateInvoice(ctx.Request.Context(), db.UpdateInvoiceParams{
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

	if err := tx.Commit(ctx); err != nil {
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
		DueDate:         updatedInvoice.DueDate.Time,
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
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Invoice updated successfully"))

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

type GetInvoiceAuditLogResponse struct {
	AuditID            int64           `json:"audit_id"`
	InvoiceID          int64           `json:"invoice_id"`
	Operation          string          `json:"operation"`
	ChangedBy          *int64          `json:"changed_by"`
	ChangedAt          time.Time       `json:"changed_at"`
	OldValues          util.JSONObject `json:"old_values"`
	NewValues          util.JSONObject `json:"new_values"`
	ChangedFields      []string        `json:"changed_fields"`
	ChangedByFirstName *string         `json:"changed_by_first_name"`
	ChangedByLastName  *string         `json:"changed_by_last_name"`
}

func (server *Server) GetInvoiceAuditLogApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	logs, err := server.store.GetInvoiceAuditLogs(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(logs) == 0 {
		ctx.JSON(http.StatusNotFound, SuccessResponse[any](nil, "No audit logs found for this invoice"))
		return
	}

	response := make([]GetInvoiceAuditLogResponse, len(logs))
	for i, log := range logs {
		response[i] = GetInvoiceAuditLogResponse{
			AuditID:            log.AuditID,
			InvoiceID:          log.InvoiceID,
			Operation:          log.Operation,
			ChangedBy:          log.ChangedBy,
			ChangedAt:          log.ChangedAt.Time,
			OldValues:          util.ParseJSONToObject(log.OldValues),
			NewValues:          util.ParseJSONToObject(log.NewValues),
			ChangedFields:      log.ChangedFields,
			ChangedByFirstName: log.ChangedByFirstName,
			ChangedByLastName:  log.ChangedByLastName,
		}
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Audit logs retrieved successfully"))

}

// ================== Invoice Logs ==================

// GetInvoiceAuditLogsResponse represents the response body for getting invoice audit logs.
type GetInvoiceAuditLogsResponse struct {
	AuditID            int64           `json:"audit_id"`
	InvoiceID          int64           `json:"invoice_id"`
	Operation          string          `json:"operation"`
	ChangedBy          *int64          `json:"changed_by"`
	ChangedAt          time.Time       `json:"changed_at"`
	OldValues          util.JSONObject `json:"old_values"`
	NewValues          util.JSONObject `json:"new_values"`
	ChangedFields      []string        `json:"changed_fields"`
	ChangedByFirstName *string         `json:"changed_by_first_name"`
	ChangedByLastName  *string         `json:"changed_by_last_name"`
}

// GetInvoiceAuditLogsApi handles the retrieval of audit logs for a specific invoice.
// @Summary Get Invoice Audit Logs
// @Description Retrieve audit logs for a specific invoice by its ID.
// @Tags Invoice
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Success 200 {object} Response[[]GetInvoiceAuditLogsResponse] "Successful
// response with audit logs"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id}/audit [get]
func (server *Server) GetInvoiceAuditLogsApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	logs, err := server.store.GetInvoiceAuditLogs(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(logs) == 0 {
		ctx.JSON(http.StatusNotFound, SuccessResponse[any](nil, "No audit logs found for this invoice"))
		return
	}

	response := make([]GetInvoiceAuditLogsResponse, len(logs))
	for i, log := range logs {
		response[i] = GetInvoiceAuditLogsResponse{
			AuditID:            log.AuditID,
			InvoiceID:          log.InvoiceID,
			Operation:          log.Operation,
			ChangedBy:          log.ChangedBy,
			ChangedAt:          log.ChangedAt.Time,
			OldValues:          util.ParseJSONToObject(log.OldValues),
			NewValues:          util.ParseJSONToObject(log.NewValues),
			ChangedFields:      log.ChangedFields,
			ChangedByFirstName: log.ChangedByFirstName,
			ChangedByLastName:  log.ChangedByLastName,
		}
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Audit logs retrieved successfully"))

}

// ================== Payment Api ==================

// CreatePaymentRequest represents the request body for creating a payment.
type CreatePaymentRequest struct {
	PaymentMethod    *string   `json:"payment_method" binding:"oneof=credit_card bank_transfer cash check other"`
	PaymentStatus    string    `json:"payment_status" binding:"required,oneof=pending completed failed refunded reversed"`
	Amount           float64   `json:"amount" binding:"required,min=0"`
	PaymentDate      time.Time `json:"payment_date" binding:"required" example:"2023-10-01T00:00:00Z"`
	PaymentReference *string   `json:"payment_reference"`
	Notes            *string   `json:"notes"`
}

// CreatePaymentResponse represents the response body for creating a payment.
type CreatePaymentResponse struct {
	PaymentID            int64     `json:"payment_id"`
	InvoiceID            int64     `json:"invoice_id"`
	PaymentMethod        *string   `json:"payment_method"`
	PaymentStatus        string    `json:"payment_status"`
	Amount               float64   `json:"amount"`
	PaymentDate          time.Time `json:"payment_date"`
	PaymentReference     *string   `json:"payment_reference"`
	Notes                *string   `json:"notes"`
	InvoiceStatusChanged bool      `json:"invoice_status_changed"`
	CurrentInvoiceStatus string    `json:"current_invoice_status"`
	RecordedBy           *int64    `json:"recorded_by"`
}

// @Summary Create Payment
// @Description Create a payment for an invoice.
// @Tags Invoice
// @Accept json
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Param request body CreatePaymentRequest true "Create Payment Request"
// @Success 200 {object} Response[CreatePaymentResponse] "Successful response with payment
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id}/payments [post]
func (server *Server) CreatePaymentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	invoiceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid invoice ID: %s", id)))
		return
	}
	var req CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)
	qtx := server.store.WithTx(tx)

	employeeID, err := qtx.GetEmployeeIDByUserID(ctx.Request.Context(), payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL myapp.current_employee_id = %d", employeeID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	getInvoice, err := qtx.GetInvoice(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payment, err := qtx.CreatePayment(ctx, db.CreatePaymentParams{
		InvoiceID:        invoiceID,
		PaymentMethod:    req.PaymentMethod,
		PaymentStatus:    req.PaymentStatus,
		Amount:           req.Amount,
		PaymentDate:      pgtype.Date{Time: req.PaymentDate, Valid: true},
		PaymentReference: req.PaymentReference,
		Notes:            req.Notes,
		RecordedBy:       &employeeID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var newInvoiceStatus string
	var statusChanged bool = false

	if req.PaymentStatus == string(invoice.PaymentStatusCompleted) {
		totalPaid, err := qtx.GetTotalPaidAmountByInvoice(ctx.Request.Context(), invoiceID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newStatus, err := invoice.DetermineInvoiceStatus(getInvoice.TotalAmount, totalPaid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if string(newStatus) != getInvoice.Status {
			updatedInvoice, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
				ID:     invoiceID,
				Status: string(newStatus)})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			newInvoiceStatus = updatedInvoice.Status
			statusChanged = true
		} else {
			newInvoiceStatus = getInvoice.Status
		}
	}
	if err := tx.Commit(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := CreatePaymentResponse{
		PaymentID:            payment.ID,
		InvoiceID:            payment.InvoiceID,
		PaymentMethod:        payment.PaymentMethod,
		PaymentStatus:        payment.PaymentStatus,
		Amount:               payment.Amount,
		PaymentDate:          payment.PaymentDate.Time,
		PaymentReference:     payment.PaymentReference,
		Notes:                payment.Notes,
		RecordedBy:           payment.RecordedBy,
		InvoiceStatusChanged: statusChanged,
		CurrentInvoiceStatus: newInvoiceStatus,
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Payment created successfully"))
}

// ListPaymentsResponse represents the response body for listing payments.
type ListPaymentsResponse struct {
	PaymentID           int64              `json:"payment_id"`
	InvoiceID           int64              `json:"invoice_id"`
	PaymentMethod       *string            `json:"payment_method"`
	PaymentStatus       string             `json:"payment_status"`
	Amount              float64            `json:"amount"`
	PaymentDate         pgtype.Date        `json:"payment_date"`
	PaymentReference    *string            `json:"payment_reference"`
	Notes               *string            `json:"notes"`
	RecordedBy          *int64             `json:"recorded_by"`
	CreatedAt           pgtype.Timestamptz `json:"created_at"`
	UpdatedAt           pgtype.Timestamptz `json:"updated_at"`
	RecordedByFirstName *string            `json:"recorded_by_first_name"`
	RecordedByLastName  *string            `json:"recorded_by_last_name"`
}

// @Summary List Payments
// @Description List all payments for a specific invoice.
// @Tags Invoice
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Success 200 {object} Response[[]ListPaymentsResponse] "Successful response with
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id}/payments [get]
func (server *Server) ListPaymentsApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid invoice ID: %s", ctx.Param("id"))))
		return
	}

	payments, err := server.store.ListPayments(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(payments) == 0 {
		ctx.JSON(http.StatusOK, SuccessResponse([]any{}, "No payments found for this invoice"))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(payments, "Payments retrieved successfully"))
}

type GetPaymentByIDResponse struct {
	PaymentID           int64     `json:"payment_id"`
	InvoiceID           int64     `json:"invoice_id"`
	PaymentMethod       *string   `json:"payment_method"`
	PaymentStatus       string    `json:"payment_status"`
	Amount              float64   `json:"amount"`
	PaymentDate         time.Time `json:"payment_date"`
	PaymentReference    *string   `json:"payment_reference"`
	Notes               *string   `json:"notes"`
	RecordedBy          *int64    `json:"recorded_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	RecordedByFirstName *string   `json:"recorded_by_first_name"`
	RecordedByLastName  *string   `json:"recorded_by_last_name"`
}

func (server *Server) GetPaymentByIDApi(ctx *gin.Context) {
	paymentID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid payment ID: %s", ctx.Param("id"))))
		return
	}

	payment, err := server.store.GetPayment(ctx.Request.Context(), paymentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := GetPaymentByIDResponse{
		PaymentID:           payment.ID,
		InvoiceID:           payment.InvoiceID,
		PaymentMethod:       payment.PaymentMethod,
		PaymentStatus:       payment.PaymentStatus,
		Amount:              payment.Amount,
		PaymentDate:         payment.PaymentDate.Time,
		PaymentReference:    payment.PaymentReference,
		Notes:               payment.Notes,
		RecordedBy:          payment.RecordedBy,
		CreatedAt:           payment.CreatedAt.Time,
		UpdatedAt:           payment.UpdatedAt.Time,
		RecordedByFirstName: payment.RecordedByFirstName,
		RecordedByLastName:  payment.RecordedByLastName,
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Payment retrieved successfully"))

}

// UpdatePaymentRequest represents the request body for updating a payment.
type UpdatePaymentRequest struct {
	PaymentMethod    *string    `json:"payment_method"`
	PaymentStatus    *string    `json:"payment_status"`
	Amount           *float64   `json:"amount"`
	PaymentDate      *time.Time `json:"payment_date"`
	PaymentReference *string    `json:"payment_reference"`
	Notes            *string    `json:"notes"`
	RecordedBy       *int64     `json:"recorded_by"`
}

// UpdatePaymentResponse represents the response body for updating a payment.
type UpdatePaymentResponse struct {
	PaymentID             int64     `json:"payment_id"`
	InvoiceID             int64     `json:"invoice_id"`
	PaymentMethod         *string   `json:"payment_method"`
	PaymentStatus         string    `json:"payment_status"`
	Amount                float64   `json:"amount"`
	PaymentDate           time.Time `json:"payment_date"`
	PaymentReference      *string   `json:"payment_reference"`
	Notes                 *string   `json:"notes"`
	RecordedBy            *int64    `json:"recorded_by"`
	InvoiceStatusChanged  bool      `json:"invoice_status_changed"`
	CurrentInvoiceStatus  string    `json:"current_invoice_status"`
	PreviousInvoiceStatus string    `json:"previous_invoice_status"`
}

// @Summary Update Payment
// @Description Update a payment for an invoice.
// @Tags Invoice
// @Accept json
// @Produce json
// @Param invoice_id path int64 true "Invoice ID"
// @Param payment_id path int64 true "Payment ID"
// @Param request body UpdatePaymentRequest true "Update Payment Request"
// @Success 200 {object} Response[UpdatePaymentResponse] "Successful response with updated payment"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{invoice_id}/payments/{payment_id} [put]
func (server *Server) UpdatePaymentApi(ctx *gin.Context) {
	invoiceIDStr := ctx.Param("invoice_id")
	paymentIDStr := ctx.Param("payment_id")

	invoiceID, err := strconv.ParseInt(invoiceIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid invoice ID: %s", invoiceIDStr)))
		return
	}

	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid payment ID: %s", paymentIDStr)))
		return
	}

	var req UpdatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)
	qtx := server.store.WithTx(tx)

	employeeID, err := qtx.GetEmployeeIDByUserID(ctx.Request.Context(), payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = tx.Exec(ctx, "SET LOCAL myapp.current_employee_id = $1", employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get the current payment with invoice info to validate and track changes
	currentPayment, err := qtx.GetPaymentWithInvoice(ctx.Request.Context(), paymentID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("payment not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Validate that payment belongs to the specified invoice
	if currentPayment.InvoiceID != invoiceID {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("payment does not belong to specified invoice")))
		return
	}

	// Store original invoice status for comparison
	originalInvoiceStatus := currentPayment.InvoiceStatus

	arg := db.UpdatePaymentParams{
		ID:               paymentID,
		PaymentMethod:    req.PaymentMethod,
		PaymentStatus:    req.PaymentStatus,
		Amount:           req.Amount,
		PaymentReference: req.PaymentReference,
		Notes:            req.Notes,
		RecordedBy:       req.RecordedBy,
	}

	if req.PaymentDate != nil {
		arg.PaymentDate = pgtype.Date{Time: *req.PaymentDate, Valid: true}
	}

	// Update the payment
	updatedPayment, err := qtx.UpdatePayment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var newInvoiceStatus string
	var statusChanged bool = false

	// Recalculate invoice status if payment status is completed
	// This handles cases where:
	// 1. Payment status changed to completed
	// 2. Payment amount changed
	// 3. Payment status changed from completed to something else
	if updatedPayment.PaymentStatus == string(invoice.PaymentStatusCompleted) ||
		currentPayment.PaymentStatus == string(invoice.PaymentStatusCompleted) {

		// Get fresh total paid amount after the update
		totalPaid, err := qtx.GetTotalPaidAmountByInvoice(ctx.Request.Context(), invoiceID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newStatus, err := invoice.DetermineInvoiceStatus(currentPayment.InvoiceTotalAmount, totalPaid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if string(newStatus) != originalInvoiceStatus {
			updatedInvoice, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
				ID:     invoiceID,
				Status: string(newStatus),
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			newInvoiceStatus = updatedInvoice.Status
			statusChanged = true
		} else {
			newInvoiceStatus = originalInvoiceStatus
		}
	} else {
		newInvoiceStatus = originalInvoiceStatus
	}

	if err := tx.Commit(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := UpdatePaymentResponse{
		PaymentID:             updatedPayment.ID,
		InvoiceID:             updatedPayment.InvoiceID,
		PaymentMethod:         updatedPayment.PaymentMethod,
		PaymentStatus:         updatedPayment.PaymentStatus,
		Amount:                updatedPayment.Amount,
		PaymentDate:           updatedPayment.PaymentDate.Time,
		PaymentReference:      updatedPayment.PaymentReference,
		Notes:                 updatedPayment.Notes,
		RecordedBy:            updatedPayment.RecordedBy,
		InvoiceStatusChanged:  statusChanged,
		CurrentInvoiceStatus:  newInvoiceStatus,
		PreviousInvoiceStatus: originalInvoiceStatus,
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Payment updated successfully"))
}

// DeletePaymentResponse represents the response body for deleting a payment.
type DeletePaymentResponse struct {
	DeletedPaymentID      int64   `json:"deleted_payment_id"`
	InvoiceID             int64   `json:"invoice_id"`
	DeletedAmount         float64 `json:"deleted_amount"`
	DeletedPaymentStatus  string  `json:"deleted_payment_status"`
	InvoiceStatusChanged  bool    `json:"invoice_status_changed"`
	CurrentInvoiceStatus  string  `json:"current_invoice_status"`
	PreviousInvoiceStatus string  `json:"previous_invoice_status"`
}

// @Summary Delete Payment
// @Description Delete a payment from an invoice.
// @Tags Invoice
// @Accept json
// @Produce json
// @Param invoice_id path int64 true "Invoice ID"
// @Param payment_id path int64 true "Payment ID"
// @Success 200 {object} Response[DeletePaymentResponse] "Successful response with deletion details"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{invoice_id}/payments/{payment_id} [delete]
func (server *Server) DeletePaymentApi(ctx *gin.Context) {
	invoiceIDStr := ctx.Param("invoice_id")
	paymentIDStr := ctx.Param("payment_id")

	invoiceID, err := strconv.ParseInt(invoiceIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid invoice ID: %s", invoiceIDStr)))
		return
	}

	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid payment ID: %s", paymentIDStr)))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)
	qtx := server.store.WithTx(tx)

	employeeID, err := qtx.GetEmployeeIDByUserID(ctx.Request.Context(), payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = tx.Exec(ctx, "SET LOCAL myapp.current_employee_id = $1", employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get the payment with invoice info before deletion to validate and track changes
	paymentToDelete, err := qtx.GetPaymentWithInvoice(ctx.Request.Context(), paymentID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("payment not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Validate that payment belongs to the specified invoice
	if paymentToDelete.InvoiceID != invoiceID {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("payment does not belong to specified invoice")))
		return
	}

	// Store original invoice status for comparison
	originalInvoiceStatus := paymentToDelete.InvoiceStatus

	// Delete the payment
	deletedPayment, err := qtx.DeletePayment(ctx.Request.Context(), paymentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var newInvoiceStatus string
	var statusChanged bool = false

	// Recalculate invoice status only if the deleted payment was completed
	// This ensures we only recalculate when the deletion actually affects the paid amount
	if deletedPayment.PaymentStatus == string(invoice.PaymentStatusCompleted) {
		// Get fresh total paid amount after the deletion
		totalPaid, err := qtx.GetTotalPaidAmountByInvoice(ctx.Request.Context(), invoiceID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newStatus, err := invoice.DetermineInvoiceStatus(paymentToDelete.InvoiceTotalAmount, totalPaid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if string(newStatus) != originalInvoiceStatus {
			updatedInvoice, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
				ID:     invoiceID,
				Status: string(newStatus),
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			newInvoiceStatus = updatedInvoice.Status
			statusChanged = true
		} else {
			newInvoiceStatus = originalInvoiceStatus
		}
	} else {
		// If deleted payment wasn't completed, invoice status shouldn't change
		newInvoiceStatus = originalInvoiceStatus
	}

	if err := tx.Commit(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := DeletePaymentResponse{
		DeletedPaymentID:      deletedPayment.ID,
		InvoiceID:             deletedPayment.InvoiceID,
		DeletedAmount:         deletedPayment.Amount,
		DeletedPaymentStatus:  deletedPayment.PaymentStatus,
		InvoiceStatusChanged:  statusChanged,
		CurrentInvoiceStatus:  newInvoiceStatus,
		PreviousInvoiceStatus: originalInvoiceStatus,
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Payment deleted successfully"))

}

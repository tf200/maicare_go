package handler

import (
	"errors"
	"fmt"
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InvoiceHandler struct {
	service domain.InvoiceService
}

func NewInvoiceHandler(service domain.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

// CreateInvoice handles POST /invoices
func (h *InvoiceHandler) CreateInvoice(ctx *gin.Context) {
	var req CreateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), "VALIDATION_ERROR"))
		return
	}

	dparams := domain.CreateInvoiceParams{
		ClientID:    req.ClientID,
		InvoiceType: req.InvoiceType,
		IssueDate:   req.IssueDate,
		DueDate:     req.DueDate,
		Lines:       make([]domain.CreateInvoiceLineInput, len(req.Lines)),
	}
	for i, l := range req.Lines {
		dparams.Lines[i] = domain.CreateInvoiceLineInput{
			LineType:    l.LineType,
			ContractID:  l.ContractID,
			ServiceType: l.ServiceType,
			Description: l.Description,
			PeriodStart: l.PeriodStart,
			PeriodEnd:   l.PeriodEnd,
			Quantity:    l.Quantity,
			Unit:        l.Unit,
			UnitPrice:   l.UnitPrice,
			VatRate:     l.VatRate,
		}
	}

	inv, lines, err := h.service.CreateInvoice(ctx.Request.Context(), dparams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "CREATE_ERROR"))
		return
	}

	lineDTOs := make([]InvoiceLineDTO, len(lines))
	for i, l := range lines {
		lineDTOs[i] = InvoiceLineDTO{
			ID: l.ID, LineNo: l.LineNo, LineType: l.LineType, ContractID: l.ContractID,
			ServiceType: l.ServiceType, Description: l.Description, PeriodStart: l.PeriodStart,
			PeriodEnd: l.PeriodEnd, Quantity: l.Quantity, Unit: l.Unit, UnitPrice: l.UnitPrice,
			NetAmount: l.NetAmount, VatRate: l.VatRate, VatAmount: l.VatAmount, GrossAmount: l.GrossAmount,
		}
	}

	resp := CreateInvoiceResponseDTO{
		ID: inv.ID, InvoiceNumber: inv.InvoiceNumber, IssueDate: inv.IssueDate, DueDate: inv.DueDate,
		Status: inv.Status, Source: inv.Source, InvoiceType: inv.InvoiceType, Currency: inv.Currency,
		NetTotal: inv.NetTotal, VatTotal: inv.VatTotal, GrossTotal: inv.GrossTotal,
		PdfAttachmentID: inv.PdfAttachmentID, ClientID: inv.ClientID, SenderID: inv.SenderID,
		Lines: lineDTOs, UpdatedAt: inv.UpdatedAt, CreatedAt: inv.CreatedAt,
	}
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Invoice created successfully"))
}

// GenerateInvoice handles POST /invoices/generate
func (h *InvoiceHandler) GenerateInvoice(ctx *gin.Context) {
	var req GenerateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), "VALIDATION_ERROR"))
		return
	}

	dparams := domain.GenerateInvoiceParams{
		ClientID:        req.ClientID,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		BillingTimezone: req.BillingTimezone,
		BillingCycle:    req.BillingCycle,
	}

	inv, _, err := h.service.GenerateInvoice(ctx.Request.Context(), dparams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "GENERATE_ERROR"))
		return
	}

	lineDTOs := make([]InvoiceLineDTO, len(inv.Lines))
	for i, l := range inv.Lines {
		lineDTOs[i] = InvoiceLineDTO{
			ID: l.ID, LineNo: l.LineNo, LineType: l.LineType, ContractID: l.ContractID,
			ServiceType: l.ServiceType, Description: l.Description, PeriodStart: l.PeriodStart,
			PeriodEnd: l.PeriodEnd, Quantity: l.Quantity, Unit: l.Unit, UnitPrice: l.UnitPrice,
			NetAmount: l.NetAmount, VatRate: l.VatRate, VatAmount: l.VatAmount, GrossAmount: l.GrossAmount,
		}
	}

	resp := GenerateInvoiceResponseDTO{
		ID: inv.ID, InvoiceNumber: inv.InvoiceNumber, IssueDate: inv.IssueDate, DueDate: inv.DueDate,
		Status: inv.Status, Source: inv.Source, InvoiceType: inv.InvoiceType,
		PeriodStart: inv.PeriodStart, PeriodEnd: inv.PeriodEnd, Currency: inv.Currency,
		NetTotal: inv.NetTotal, VatTotal: inv.VatTotal, GrossTotal: inv.GrossTotal,
		PdfAttachmentID: inv.PdfAttachmentID, ClientID: inv.ClientID, SenderID: inv.SenderID,
		Lines: lineDTOs, Warnings: inv.Warnings, UpdatedAt: inv.UpdatedAt, CreatedAt: inv.CreatedAt,
	}
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Invoice generated successfully"))
}

// CreditInvoice handles POST /invoices/:id/credit
func (h *InvoiceHandler) CreditInvoice(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}

	result, err := h.service.CreditInvoice(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "CREDIT_ERROR"))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK(CreditInvoiceResponseDTO{ID: result.ID}, "Credit note created successfully"))
}

// ListInvoices handles GET /invoices
func (h *InvoiceHandler) ListInvoices(ctx *gin.Context) {
	var req ListInvoicesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), "VALIDATION_ERROR"))
		return
	}

	pageReq := httpapi.PageRequest{Page: 1, PageSize: 20}
	if req.PageSize > 0 {
		pageReq.PageSize = req.PageSize
	}

	params := domain.ListInvoicesParams{
		ClientID:        req.ClientID,
		SenderID:        req.SenderID,
		Status:          req.Status,
		Statuses:        req.Statuses,
		Source:          req.Source,
		InvoiceType:     req.InvoiceType,
		RunID:           req.RunID,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		PeriodStart:     req.PeriodStart,
		PeriodEnd:       req.PeriodEnd,
		Locked:          req.Locked,
		MinWarningCount: req.MinWarningCount,
		Q:               req.Q,
		SortBy:          req.SortBy,
		SortDir:         req.SortDir,
	}
	if req.Page > 0 {
		pageReq.Page = req.Page
	}
	if req.PageSize > 0 {
		pageReq.PageSize = req.PageSize
	}
	p := pageReq.Params()
	params.Limit = p.Limit
	params.Offset = p.Offset

	result, err := h.service.ListInvoices(ctx.Request.Context(), params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "LIST_ERROR"))
		return
	}

	items := make([]ListInvoicesItemDTO, len(result.Items))
	for i, item := range result.Items {
		items[i] = ListInvoicesItemDTO{
			ID: item.ID, InvoiceNumber: item.InvoiceNumber, SenderName: item.SenderName,
			IsOverdue: item.IsOverdue, ClientFirstName: item.ClientFirstName, ClientLastName: item.ClientLastName,
			ClientFilenumber: item.ClientFilenumber, Currency: item.Currency, GrossTotal: item.GrossTotal,
			BalanceDue: item.BalanceDue, PaidTotal: item.PaidTotal, Status: item.Status,
			IssueDate: item.IssueDate, DueDate: item.DueDate, ClientID: item.ClientID, SenderID: item.SenderID,
		}
	}

	ctx.JSON(http.StatusOK, httpapi.OK(httpapi.NewPageResponse(ctx, pageReq, items, result.TotalCount), "Invoices retrieved successfully"))
}

// GetInvoiceByID handles GET /invoices/:id
func (h *InvoiceHandler) GetInvoiceByID(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}

	inv, lines, paymentPct, err := h.service.GetInvoiceByID(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "GET_ERROR"))
		return
	}

	lineDTOs := make([]InvoiceLineDTO, len(lines))
	for i, l := range lines {
		lineDTOs[i] = InvoiceLineDTO{
			ID: l.ID, LineNo: l.LineNo, LineType: l.LineType, ContractID: l.ContractID,
			ServiceType: l.ServiceType, Description: l.Description, PeriodStart: l.PeriodStart,
			PeriodEnd: l.PeriodEnd, Quantity: l.Quantity, Unit: l.Unit, UnitPrice: l.UnitPrice,
			NetAmount: l.NetAmount, VatRate: l.VatRate, VatAmount: l.VatAmount, GrossAmount: l.GrossAmount,
		}
	}

	resp := GetInvoiceByIDResponseDTO{
		ID: inv.ID, InvoiceNumber: inv.InvoiceNumber, IssueDate: inv.IssueDate, DueDate: inv.DueDate,
		Status: inv.Status, Source: inv.Source, InvoiceType: inv.InvoiceType,
		OriginalInvoiceID: inv.OriginalInvoiceID, ReplacesInvoiceID: inv.ReplacesInvoiceID,
		PeriodStart: inv.PeriodStart, PeriodEnd: inv.PeriodEnd, BillingTimezone: inv.BillingTimezone,
		Currency: inv.Currency, NetTotal: inv.NetTotal, VatTotal: inv.VatTotal, GrossTotal: inv.GrossTotal,
		ClientID: inv.ClientID, SenderID: inv.SenderID, Lines: lineDTOs, UpdatedAt: inv.UpdatedAt, CreatedAt: inv.CreatedAt,
		SenderName: inv.SenderName, ClientFirstName: inv.ClientFirstName, ClientLastName: inv.ClientLastName,
		PaymentCompletionPrc: paymentPct,
	}
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Invoice retrieved successfully"))
}

// UpdateInvoice handles PUT /invoices/:id
func (h *InvoiceHandler) UpdateInvoice(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}

	var req UpdateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), "VALIDATION_ERROR"))
		return
	}

	dparams := domain.CreateInvoiceParams{
		IssueDate: req.IssueDate,
		DueDate:   req.DueDate,
		Lines:     make([]domain.CreateInvoiceLineInput, len(req.Lines)),
	}
	for i, l := range req.Lines {
		dparams.Lines[i] = domain.CreateInvoiceLineInput{
			LineType:    l.LineType,
			ContractID:  l.ContractID,
			ServiceType: l.ServiceType,
			Description: l.Description,
			PeriodStart: l.PeriodStart,
			PeriodEnd:   l.PeriodEnd,
			Quantity:    l.Quantity,
			Unit:        l.Unit,
			UnitPrice:   l.UnitPrice,
			VatRate:     l.VatRate,
		}
	}

	inv, err := h.service.UpdateInvoice(ctx.Request.Context(), invoiceID, dparams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "UPDATE_ERROR"))
		return
	}

	resp := UpdateInvoiceResponseDTO{
		ID: inv.ID, InvoiceNumber: inv.InvoiceNumber, IssueDate: inv.IssueDate, DueDate: inv.DueDate,
		Status: inv.Status, Currency: inv.Currency, NetTotal: inv.NetTotal, VatTotal: inv.VatTotal,
		GrossTotal: inv.GrossTotal, PdfAttachmentID: inv.PdfAttachmentID, ClientID: inv.ClientID,
		SenderID: inv.SenderID, UpdatedAt: inv.UpdatedAt, CreatedAt: inv.CreatedAt,
	}
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Invoice updated successfully"))
}

// DeleteInvoice handles DELETE /invoices/:id
func (h *InvoiceHandler) DeleteInvoice(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}

	if err := h.service.DeleteInvoice(ctx.Request.Context(), invoiceID); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "DELETE_ERROR"))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Invoice deleted successfully"))
}

// GenerateInvoicePDF handles GET /invoices/:id/generate_pdf
func (h *InvoiceHandler) GenerateInvoicePDF(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}

	result, err := h.service.GenerateInvoicePDF(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "PDF_ERROR"))
		return
	}
	ctx.JSON(http.StatusCreated, httpapi.OK(GeneratePDFResponseDTO{FileURL: result.FileURL}, "Invoice PDF generated"))
}

// GetInvoiceTemplateItems handles GET /invoices/template_items
func (h *InvoiceHandler) GetInvoiceTemplateItems(ctx *gin.Context) {
	items, err := h.service.GetInvoiceTemplateItems(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "TEMPLATE_ERROR"))
		return
	}

	dtoItems := make([]InvoiceTemplateItemDTO, len(items))
	for i, it := range items {
		dtoItems[i] = InvoiceTemplateItemDTO{
			ID: it.ID, ItemTag: it.ItemTag, Description: it.Description,
			SourceTable: it.SourceTable, SourceColumn: it.SourceColumn,
		}
	}
	ctx.JSON(http.StatusOK, httpapi.OK(dtoItems, "Template items retrieved successfully"))
}

// SendInvoiceReminder handles POST /invoices/:id/send_reminder
func (h *InvoiceHandler) SendInvoiceReminder(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}

	if err := h.service.SendInvoiceReminder(ctx.Request.Context(), invoiceID); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "REMINDER_ERROR"))
		return
	}
	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Invoice reminder sent successfully"))
}

// GetInvoiceAuditLog handles GET /invoices/:id/audit
func (h *InvoiceHandler) GetInvoiceAuditLog(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}

	logs, err := h.service.GetInvoiceAuditLogs(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "AUDIT_ERROR"))
		return
	}

	dtoLogs := make([]InvoiceAuditLogDTO, len(logs))
	for i, l := range logs {
		dtoLogs[i] = InvoiceAuditLogDTO{
			AuditID: l.AuditID, InvoiceID: l.InvoiceID, Operation: l.Operation,
			ChangedBy: l.ChangedBy, ChangedAt: l.ChangedAt, ChangedFields: l.ChangedFields,
			ChangedByFirstName: l.ChangedByFirstName, ChangedByLastName: l.ChangedByLastName,
		}
	}
	ctx.JSON(http.StatusOK, httpapi.OK(dtoLogs, "Audit logs retrieved successfully"))
}

// ==================== Payment Handlers ====================

// CreatePayment handles POST /invoices/:id/payments
func (h *InvoiceHandler) CreatePayment(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}

	var req CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), "VALIDATION_ERROR"))
		return
	}

	dparams := domain.CreatePaymentParams{
		PaymentMethod:    req.PaymentMethod,
		PaymentStatus:    req.PaymentStatus,
		Amount:           req.Amount,
		PaymentDate:      req.PaymentDate,
		PaymentReference: req.PaymentReference,
		Notes:            req.Notes,
	}

	result, err := h.service.CreatePayment(ctx.Request.Context(), invoiceID, dparams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "PAYMENT_ERROR"))
		return
	}

	resp := CreatePaymentResponseDTO{
		PaymentID: result.PaymentID, InvoiceID: result.InvoiceID, PaymentMethod: result.PaymentMethod,
		PaymentStatus: result.PaymentStatus, Amount: result.Amount, PaymentDate: result.PaymentDate,
		PaymentReference: result.PaymentReference, Notes: result.Notes,
		InvoiceStatusChanged: result.InvoiceStatusChanged, CurrentInvoiceStatus: result.CurrentInvoiceStatus,
		RecordedBy: result.RecordedBy,
	}
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Payment created successfully"))
}

// ListPayments handles GET /invoices/:id/payments
func (h *InvoiceHandler) ListPayments(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}

	payments, err := h.service.ListPayments(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "LIST_ERROR"))
		return
	}

	items := make([]ListPaymentsItemDTO, len(payments))
	for i, p := range payments {
		items[i] = ListPaymentsItemDTO{
			PaymentID: p.ID, InvoiceID: p.InvoiceID, PaymentMethod: p.PaymentMethod,
			PaymentStatus: p.PaymentStatus, Amount: p.Amount, PaymentDate: p.PaymentDate,
			PaymentReference: p.PaymentReference, Notes: p.Notes, RecordedBy: p.RecordedBy,
			CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
			RecordedByFirstName: p.RecordedByFirstName, RecordedByLastName: p.RecordedByLastName,
		}
	}
	ctx.JSON(http.StatusOK, httpapi.OK(items, "Payments retrieved successfully"))
}

// GetPaymentByID handles GET /invoices/:id/payments/:payment_id
func (h *InvoiceHandler) GetPaymentByID(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid invoice ID", "INVALID_ID"))
		return
	}
	paymentID, err := uuid.Parse(ctx.Param("payment_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid payment ID", "INVALID_ID"))
		return
	}

	payment, err := h.service.GetPaymentByID(ctx.Request.Context(), invoiceID, paymentID)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			ctx.JSON(http.StatusNotFound, httpapi.Fail(err.Error(), "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "GET_ERROR"))
		return
	}

	resp := GetPaymentByIDResponseDTO{
		PaymentID: payment.ID, InvoiceID: payment.InvoiceID, PaymentMethod: payment.PaymentMethod,
		PaymentStatus: payment.PaymentStatus, Amount: payment.Amount, PaymentDate: payment.PaymentDate,
		PaymentReference: payment.PaymentReference, Notes: payment.Notes, RecordedBy: payment.RecordedBy,
		CreatedAt: payment.CreatedAt, UpdatedAt: payment.UpdatedAt,
		RecordedByFirstName: payment.RecordedByFirstName, RecordedByLastName: payment.RecordedByLastName,
	}
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Payment retrieved successfully"))
}

// UpdatePayment handles PUT /invoices/:id/payments/:payment_id
func (h *InvoiceHandler) UpdatePayment(ctx *gin.Context) {
	invoiceIDStr := ctx.Param("id")
	paymentIDStr := ctx.Param("payment_id")

	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(fmt.Sprintf("invalid invoice ID: %s", invoiceIDStr), "INVALID_ID"))
		return
	}

	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(fmt.Sprintf("invalid payment ID: %s", paymentIDStr), "INVALID_ID"))
		return
	}

	var req UpdatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), "VALIDATION_ERROR"))
		return
	}

	dparams := domain.UpdatePaymentParams{
		PaymentMethod:    req.PaymentMethod,
		PaymentStatus:    req.PaymentStatus,
		Amount:           req.Amount,
		PaymentDate:      req.PaymentDate,
		PaymentReference: req.PaymentReference,
		Notes:            req.Notes,
	}

	result, err := h.service.UpdatePayment(ctx.Request.Context(), invoiceID, paymentID, dparams)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			ctx.JSON(http.StatusNotFound, httpapi.Fail(err.Error(), "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "UPDATE_ERROR"))
		return
	}

	resp := UpdatePaymentResponseDTO{
		PaymentID: result.PaymentID, InvoiceID: result.InvoiceID, PaymentMethod: result.PaymentMethod,
		PaymentStatus: result.PaymentStatus, Amount: result.Amount, PaymentDate: result.PaymentDate,
		PaymentReference: result.PaymentReference, Notes: result.Notes, RecordedBy: result.RecordedBy,
		InvoiceStatusChanged: result.InvoiceStatusChanged, CurrentInvoiceStatus: result.CurrentInvoiceStatus,
		PreviousInvoiceStatus: result.PreviousInvoiceStatus,
	}
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Payment updated successfully"))
}

// DeletePayment handles DELETE /invoices/:id/payments/:payment_id
func (h *InvoiceHandler) DeletePayment(ctx *gin.Context) {
	invoiceIDStr := ctx.Param("id")
	paymentIDStr := ctx.Param("payment_id")

	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(fmt.Sprintf("invalid invoice ID: %s", invoiceIDStr), "INVALID_ID"))
		return
	}

	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(fmt.Sprintf("invalid payment ID: %s", paymentIDStr), "INVALID_ID"))
		return
	}

	result, err := h.service.DeletePayment(ctx.Request.Context(), invoiceID, paymentID)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			ctx.JSON(http.StatusNotFound, httpapi.Fail(err.Error(), "NOT_FOUND"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail(err.Error(), "DELETE_ERROR"))
		return
	}

	resp := DeletePaymentResponseDTO{
		DeletedPaymentID: result.DeletedPaymentID, InvoiceID: result.InvoiceID,
		DeletedAmount: result.DeletedAmount, DeletedPaymentStatus: result.DeletedPaymentStatus,
		InvoiceStatusChanged: result.InvoiceStatusChanged, CurrentInvoiceStatus: result.CurrentInvoiceStatus,
		PreviousInvoiceStatus: result.PreviousInvoiceStatus,
	}
	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Payment deleted successfully"))
}

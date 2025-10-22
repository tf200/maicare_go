package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/invoice"
	"maicare_go/pdf"
	invserv "maicare_go/service/invoice"
	"maicare_go/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

// @Summary Create Invoice
// @Description Create a new invoice with the provided details.
// @Tags Invoice
// @Accept json
// @Produce json
// @Param request body CreateInvoiceRequest true "Create Invoice Request"
// @Success 200 {object} Response[CreateInvoiceResponse] "Successful response with invoice details"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices [post]
func (server *Server) CreateInvoiceApi(ctx *gin.Context) {
	var req invserv.CreateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	response, err := server.businessService.InvoiceService.CreateInvoice(ctx.Request.Context(), req, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(response, "Invoice created successfully")
	ctx.JSON(http.StatusOK, res)
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
	var req invserv.GenerateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelWarn, "GenerateInvoiceApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	inv, _, err := server.businessService.InvoiceService.GenerateInvoice(req, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(inv, "Invoice generated successfully")
	ctx.JSON(http.StatusOK, res)

}

// @Summary Credit Invoice
// @Description Create a credit note for an existing invoice.
// @Tags Invoice
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Success 200 {object} Response[CreditInvoiceResponse] "Successful response with credit note details"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id}/credit [post]
func (server *Server) CreditInvoiceApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid invoice ID: %s", ctx.Param("id"))))
		return
	}
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	creditNoteInvoice, err := server.businessService.InvoiceService.CreditInvoice(ctx.Request.Context(), invoiceID, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(creditNoteInvoice, "Credit note created successfully"))
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
	var req invserv.ListInvoicesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pag, err := server.businessService.InvoiceService.ListInvoices(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(pag, "Invoices retrieved successfully")
	ctx.JSON(http.StatusOK, res)

}

// @Summary Get Invoice by ID
// @Description Get an invoice by its ID.
// @Tags Invoice
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Success 200 {object} Response[invoice.GetInvoiceByIDResponse] "Successful response with invoice details"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id} [get]
func (server *Server) GetInvoiceByIDApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid invoice ID")))
		return
	}

	response, err := server.businessService.InvoiceService.GetInvoiceByID(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Invoice retrieved successfully"))

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

	var req invserv.UpdateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	response, err := server.businessService.InvoiceService.UpdateInvoice(ctx, invoiceID, req, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
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

	err = server.businessService.InvoiceService.DeleteInvoice(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse[any](nil, "Invoice deleted successfully"))
}

// GenerateInvoicePDFResponse represents the response body for generating an invoice PDF.
type GenerateInvoicePDFResponse struct {
	FileUrl string `json:"file_url"`
}

// GenerateInvoicePdfApi handles generation of invoice in pdf format
// @Summary Generate InvoicePdf
// @Description Generate an invoice in PDF format by its ID.
// @Tags Invoice
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Success 201 {object} Response[GenerateInvoicePDFResponse] "Successful response indicating generation"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id}/generate_pdf [get]
func (server *Server) GenerateInvoicePdfApi(ctx *gin.Context) {
	id := ctx.Param("id")
	invoiceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid invoice ID: %s", id)))
		return
	}

	invoiceData, err := server.store.GetInvoice(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var invoiceDetails []invoice.InvoiceDetails
	if err := json.Unmarshal(invoiceData.InvoiceDetails, &invoiceDetails); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var senderContacts []SenderContact

	err = json.Unmarshal(invoiceData.SenderContacts, &senderContacts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var pdfInvoiceDetails []pdf.InvoiceDetail

	for _, value := range invoiceDetails {
		var pdfInvoicePeriods []pdf.InvoicePeriod
		for _, period := range value.Periods {
			pdfInvoicePeriods = append(pdfInvoicePeriods, pdf.InvoicePeriod{
				StartDate:             period.StartDate,
				EndDate:               period.EndDate,
				AcommodationTimeFrame: util.DerefString(period.AcommodationTimeFrame),
				AmbulanteTotalMinutes: util.DerefFloat64(period.AmbulanteTotalMinutes),
			})

		}
		pdfInvoiceDetails = append(pdfInvoiceDetails, pdf.InvoiceDetail{
			CareType:      value.ContractType,
			Price:         value.Price,
			PriceTimeUnit: value.PriceTimeUnit,
			PreVatTotal:   value.PreVatTotal,
			Total:         value.Total,
			Periods:       pdfInvoicePeriods,
		})
	}

	var extraItems map[string]string
	if invoiceData.ExtraContent != nil {
		if err := json.Unmarshal(invoiceData.ExtraContent, &extraItems); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	arg := pdf.InvoicePDFData{
		ID:                   invoiceData.ID,
		SenderName:           util.DerefString(invoiceData.SenderName),
		SenderContactPerson:  util.DerefString(senderContacts[0].Name),
		SenderAddressLine1:   util.DerefString(invoiceData.SenderAddress),
		SenderPostalCodeCity: util.DerefString(invoiceData.SenderPostalCode),
		InvoiceNumber:        invoiceData.InvoiceNumber,
		InvoiceDate:          invoiceData.IssueDate.Time,
		DueDate:              invoiceData.DueDate.Time,
		InvoiceDetails:       pdfInvoiceDetails,
		ExtraItems:           extraItems,
	}

	fileKey, filesize, err := pdf.GenerateAndUploadInvoicePDF(ctx, arg, server.b2Client)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to generate invoice PDF: %w", err)))
		return
	}

	fileArgs := db.CreateAttachmentParams{
		Name: "Invoice_" + invoiceData.InvoiceNumber + ".pdf",
		File: fileKey,
		Size: int32(filesize),
		Tag:  util.StringPtr(""),
	}
	attachment, err := server.store.CreateAttachment(ctx, fileArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	url := server.generateResponsePresignedURL(&attachment.File)
	if url == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to generate presigned url")))
		return
	}
	res := SuccessResponse(GenerateInvoicePDFResponse{
		FileUrl: *url,
	}, "Invoice Pdf generated")
	ctx.JSON(http.StatusCreated, res)

}

// GetInvoiceTemplateItemsApi handles the retrieval of all invoice template items.
// @Summary Get Invoice Template Items
// @Description Retrieve all invoice template items.
// @Tags Invoice
// @Produce json
// @Success 200 {object} Response[[]GetInvoiceTemplateItemsResponse] "Successful
// response with template items"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/template_items [get]
func (server *Server) GetInvoiceTemplateItemsApi(ctx *gin.Context) {
	response, err := server.businessService.InvoiceService.GetInvoiceTemplateItemsApi(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Template items retrieved successfully"))

}

// SendInvoiceReminderApi handles sending a reminder for a specific invoice.
// @Summary Send Invoice Reminder
// @Description Send a reminder for a specific invoice by its ID.
// @Tags Invoice
// @Produce json
// @Param id path int64 true "Invoice ID"
// @Success 200 {object} Response[any] "Successful response indicating reminder sent"
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id}/send_reminder [post]
func (server *Server) SendInvoiceReminderApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid invoice ID: %s", ctx.Param("id"))))
		return
	}

	err = server.businessService.InvoiceService.SendInvoiceReminder(ctx, invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to send invoice reminder: %w", err)))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse[any](nil, "Invoice reminder sent successfully"))
}

// ================== Invoice Logs ==================

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
func (server *Server) GetInvoiceAuditLogApi(ctx *gin.Context) {
	invoiceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	logs, err := server.businessService.InvoiceService.GetInvoiceAuditLogs(ctx.Request.Context(), invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(logs) == 0 {
		ctx.JSON(http.StatusNotFound, SuccessResponse[any](nil, "No audit logs found for this invoice"))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(logs, "Audit logs retrieved successfully"))

}

// ================== Payment Api ==================

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
	var req invserv.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	response, err := server.businessService.InvoiceService.CreatePayment(ctx, invoiceID, req, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Payment created successfully"))
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

	response, err := server.businessService.InvoiceService.ListPayments(ctx, invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Payments retrieved successfully"))
}

// @Summary Get Payment by ID
// @Description Get a payment by its ID.
// @Tags Invoice
// @Produce json
// @Param id path int64 true "Payment ID"
// @Success 200 {object} Response[GetPaymentByIDResponse] "Successful response
// @Failure 400,401,404,500 {object} Response[any]
// @Router /invoices/{id}/payments/{payment_id} [get]
func (server *Server) GetPaymentByIDApi(ctx *gin.Context) {
	paymentID, err := strconv.ParseInt(ctx.Param("payment_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid payment ID: %s", ctx.Param("id"))))
		return
	}

	response, err := server.businessService.InvoiceService.GetPaymentByID(ctx, paymentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Payment retrieved successfully"))

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
	invoiceIDStr := ctx.Param("id")
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

	var req invserv.UpdatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	response, err := server.businessService.InvoiceService.UpdatePayment(ctx, invoiceID, payload.EmployeeID, paymentID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Payment updated successfully"))
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

	response, err := server.businessService.InvoiceService.DeletePayment(ctx, invoiceID, paymentID, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Payment deleted successfully"))

}

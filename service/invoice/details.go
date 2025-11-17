package invoice

import (
	"context"
	"encoding/json"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/pdf"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *invoiceService) CreateInvoice(
	ctx context.Context,
	req CreateInvoiceRequest,
	employeeID uuid.UUID,
) (*CreateInvoiceResponse, error) {
	_, err := VerifyTotalAmount(req.InvoiceDetails, req.TotalAmount)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"CreateInvoice",
			"Total amount verification failed",
			zap.Error(err),
			zap.String("client_id", req.ClientID.String()),
		)
		return nil, fmt.Errorf("total amount verification failed: %v", err)
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"CreateInvoice",
			"Failed to begin transaction",
			zap.Error(err),
			zap.String("client_id", req.ClientID.String()),
		)
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.Store.WithTx(tx)
	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL myapp.current_employee_id = %d", employeeID))
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"CreateInvoice",
			"Failed to set current employee ID",
			zap.Error(err),
			zap.String("client_id", req.ClientID.String()),
		)
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}

	sender, err := qtx.GetClientSender(ctx, req.ClientID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"CreateInvoice",
			"Failed to get client sender",
			zap.Error(err),
			zap.String("client_id", req.ClientID.String()),
		)
		return nil, fmt.Errorf("failed to get client sender: %v", err)
	}

	invoiceNumber, invoiceSequence, err := s.GenerateInvoiceNumber(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"CreateInvoice",
			"Failed to generate invoice number",
			zap.Error(err),
			zap.String("client_id", req.ClientID.String()),
		)
		return nil, fmt.Errorf("failed to generate invoice number: %v", err)
	}

	invoiceDetailsBytes, err := json.Marshal(req.InvoiceDetails)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"CreateInvoice",
			"Failed to marshal invoice details",
			zap.Error(err),
			zap.String("client_id", req.ClientID.String()),
		)
		return nil, fmt.Errorf("failed to marshal invoice details: %v", err)
	}

	arg := db.CreateInvoiceParams{
		InvoiceNumber:   invoiceNumber,
		InvoiceSequence: invoiceSequence,
		DueDate:         pgtype.Date{Time: req.DueDate, Valid: true},
		IssueDate:       pgtype.Date{Time: req.IssueDate, Valid: true},
		InvoiceDetails:  invoiceDetailsBytes,
		TotalAmount:     req.TotalAmount,
		ExtraContent:    util.ParseObjectToJSON(req.ExtraContent),
		ClientID:        req.ClientID,
		SenderID:        &sender.ID,
		WarningCount:    0,
		InvoiceType:     db.InvoiceTypeEnum(req.InvoiceType),
	}
	invoice, err := qtx.CreateInvoice(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"CreateInvoice",
			"Failed to create invoice",
			zap.Error(err),
			zap.String("client_id", req.ClientID.String()),
		)
		return nil, fmt.Errorf("failed to create invoice: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"CreateInvoice",
			"Failed to commit transaction",
			zap.Error(err),
			zap.String("client_id", req.ClientID.String()),
		)
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &CreateInvoiceResponse{
		ID:              invoice.ID,
		InvoiceNumber:   invoice.InvoiceNumber,
		IssueDate:       invoice.IssueDate.Time,
		DueDate:         invoice.DueDate.Time,
		Status:          string(invoice.Status),
		InvoiceDetails:  req.InvoiceDetails,
		TotalAmount:     req.TotalAmount,
		PdfAttachmentID: invoice.PdfAttachmentID,
		ExtraContent:    util.ParseJSONToObject(invoice.ExtraContent),
		ClientID:        invoice.ClientID,
		SenderID:        invoice.SenderID,
		UpdatedAt:       invoice.UpdatedAt.Time,
		CreatedAt:       invoice.CreatedAt.Time,
	}, nil
}

func (s *invoiceService) GetInvoiceByID(
	ctx context.Context,
	invoiceID int64,
) (*GetInvoiceByIDResponse, error) {
	inv, err := s.Store.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"GetInvoiceByID",
			"Failed to get invoice",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, err
	}

	var invoiceDetails []InvoiceDetails
	if err := json.Unmarshal(inv.InvoiceDetails, &invoiceDetails); err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"GetInvoiceByID",
			"Failed to unmarshal invoice details",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to unmarshal invoice details")
	}

	paymentCompletionPrc := s.calculatePaymentCompletionPercentage(ctx, inv.TotalAmount, invoiceID)

	return &GetInvoiceByIDResponse{
		ID:                   inv.ID,
		InvoiceNumber:        inv.InvoiceNumber,
		IssueDate:            inv.IssueDate.Time,
		DueDate:              inv.DueDate.Time,
		Status:               string(inv.Status),
		InvoiceDetails:       invoiceDetails,
		TotalAmount:          inv.TotalAmount,
		PdfAttachmentID:      inv.PdfAttachmentID,
		ExtraContent:         util.ParseJSONToObject(inv.ExtraContent),
		ClientID:             inv.ClientID,
		SenderID:             inv.SenderID,
		InvoiceType:          string(inv.InvoiceType),
		OriginalInvoiceID:    inv.OriginalInvoiceID,
		UpdatedAt:            inv.UpdatedAt.Time,
		CreatedAt:            inv.CreatedAt.Time,
		SenderName:           inv.SenderName,
		SenderKvknumber:      inv.SenderKvknumber,
		SenderBtwnumber:      inv.SenderBtwnumber,
		ClientFirstName:      inv.ClientFirstName,
		ClientLastName:       inv.ClientLastName,
		PaymentCompletionPrc: paymentCompletionPrc,
	}, nil
}

func (s *invoiceService) ListInvoices(
	ctx *gin.Context,
	req ListInvoicesRequest,
) (*pagination.Response[ListInvoicesResponse], error) {
	params := req.GetParams()

	invoices, err := s.Store.ListInvoices(ctx, db.ListInvoicesParams{
		ClientID: req.ClientID,
		SenderID: req.SenderID,
		Status: func() db.NullInvoiceStatusEnum {
			if req.Status != nil {
				return db.NullInvoiceStatusEnum{Valid: true, InvoiceStatusEnum: db.InvoiceStatusEnum(*req.Status)}
			}
			return db.NullInvoiceStatusEnum{Valid: false}
		}(),

		StartDate: pgtype.Date{Time: req.StartDate, Valid: !req.StartDate.IsZero()},
		EndDate:   pgtype.Date{Time: req.EndDate, Valid: !req.EndDate.IsZero()},
		Limit:     params.Limit,
		Offset:    params.Offset,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"ListInvoices",
			"Failed to list invoices",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list invoices: %v", err)
	}

	var invoiceResponses []ListInvoicesResponse
	for _, inv := range invoices {
		var invoiceDetails []InvoiceDetails
		if err := json.Unmarshal(inv.InvoiceDetails, &invoiceDetails); err != nil {
			s.Logger.LogBusinessEvent(
				logger.LogLevelError,
				"ListInvoices",
				"Failed to unmarshal invoice details",
				zap.Error(err),
				zap.Int64("invoice_id", inv.ID),
			)
			return nil, fmt.Errorf(
				"failed to unmarshal invoice details for invoice ID %d: %v",
				inv.ID,
				err,
			)
		}

		invoiceResponses = append(invoiceResponses, ListInvoicesResponse{
			ID:                inv.ID,
			InvoiceNumber:     inv.InvoiceNumber,
			IssueDate:         inv.IssueDate.Time,
			DueDate:           inv.DueDate.Time,
			Status:            string(inv.Status),
			InvoiceDetails:    invoiceDetails,
			TotalAmount:       inv.TotalAmount,
			PdfAttachmentID:   inv.PdfAttachmentID,
			ExtraContent:      util.ParseJSONToObject(inv.ExtraContent),
			ClientID:          inv.ClientID,
			SenderID:          inv.SenderID,
			InvoiceType:       string(inv.InvoiceType),
			OriginalInvoiceID: inv.OriginalInvoiceID,
			UpdatedAt:         inv.UpdatedAt.Time,
			CreatedAt:         inv.CreatedAt.Time,
			SenderName:        inv.SenderName,
			ClientFirstName:   inv.ClientFirstName,
			ClientLastName:    inv.ClientLastName,
			WarningCount:      inv.WarningCount,
		})
	}
	pag := pagination.NewResponse(ctx, req.Request, invoiceResponses, invoices[0].TotalCount)
	return &pag, nil
}

func (s *invoiceService) UpdateInvoice(
	ctx context.Context,
	invoiceID int64,
	req UpdateInvoiceRequest,
	employeeID uuid.UUID,
) (*UpdateInvoiceResponse, error) {
	_, err := VerifyTotalAmount(req.InvoiceDetails, req.TotalAmount)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"UpdateInvoice",
			"Total amount verification failed",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("total amount verification failed: %v", err)
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"UpdateInvoice",
			"Failed to begin transaction",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.Store.WithTx(tx)
	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL myapp.current_employee_id = %d", employeeID))
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"UpdateInvoice",
			"Failed to set current employee ID",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}

	invoiceDetailsBytes, err := json.Marshal(req.InvoiceDetails)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"UpdateInvoice",
			"Failed to marshal invoice details",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to marshal invoice details: %v", err)
	}

	arg := db.UpdateInvoiceParams{
		ID:             invoiceID,
		IssueDate:      pgtype.Date{Time: req.IssueDate},
		DueDate:        pgtype.Date{Time: req.DueDate},
		InvoiceDetails: invoiceDetailsBytes,
		TotalAmount:    &req.TotalAmount,
		ExtraContent:   util.ParseObjectToJSON(req.ExtraContent),
		Status: db.NullInvoiceStatusEnum{
			Valid:             true,
			InvoiceStatusEnum: db.InvoiceStatusEnum(req.Status),
		},
		WarningCount: &req.WarningCount,
	}
	updatedInv, err := qtx.UpdateInvoice(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"UpdateInvoice",
			"Failed to update invoice",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to update invoice: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"UpdateInvoice",
			"Failed to commit transaction",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &UpdateInvoiceResponse{
		ID:              updatedInv.ID,
		InvoiceNumber:   updatedInv.InvoiceNumber,
		IssueDate:       updatedInv.IssueDate.Time,
		DueDate:         updatedInv.DueDate.Time,
		Status:          string(updatedInv.Status),
		InvoiceDetails:  req.InvoiceDetails,
		TotalAmount:     req.TotalAmount,
		PdfAttachmentID: updatedInv.PdfAttachmentID,
		ExtraContent:    util.ParseJSONToObject(updatedInv.ExtraContent),
		ClientID:        updatedInv.ClientID,
		SenderID:        updatedInv.SenderID,
		UpdatedAt:       updatedInv.UpdatedAt.Time,
		CreatedAt:       updatedInv.CreatedAt.Time,
	}, nil
}

func (s *invoiceService) DeleteInvoice(ctx context.Context, invoiceID int64) error {
	err := s.Store.DeleteInvoice(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"DeleteInvoice",
			"Failed to delete invoice",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return fmt.Errorf("failed to delete invoice: %v", err)
	}
	return nil
}

func (s *invoiceService) GetInvoiceAuditLogs(
	ctx context.Context,
	invoiceID int64,
) ([]GetInvoiceAuditLogResponse, error) {
	logs, err := s.Store.GetInvoiceAuditLogs(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"GetInvoiceAuditLog",
			"Failed to get invoice audit logs",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to get invoice audit logs: %v", err)
	}

	response := []GetInvoiceAuditLogResponse{}

	for _, log := range logs {
		response = append(response, GetInvoiceAuditLogResponse{
			AuditID:            log.AuditID,
			InvoiceID:          log.InvoiceID,
			Operation:          string(log.Operation),
			ChangedBy:          log.ChangedBy,
			ChangedAt:          log.ChangedAt.Time,
			OldValues:          util.ParseJSONToObject(log.OldValues),
			NewValues:          util.ParseJSONToObject(log.NewValues),
			ChangedFields:      log.ChangedFields,
			ChangedByFirstName: log.ChangedByFirstName,
			ChangedByLastName:  log.ChangedByLastName,
		})
	}
	return response, nil
}

// To do : Implement pdf attachment handling

func (s *invoiceService) GetInvoiceTemplateItemsApi(
	ctx context.Context,
) ([]GetInvoiceTemplateItemsResponse, error) {
	templateItems, err := s.Store.GetAllTemplateItems(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"GetInvoiceTemplateItems",
			"Failed to get invoice template items",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get invoice template items: %v", err)
	}

	response := []GetInvoiceTemplateItemsResponse{}

	for _, item := range templateItems {
		response = append(response, GetInvoiceTemplateItemsResponse{
			ID:           item.ID,
			ItemTag:      item.ItemTag,
			Description:  item.Description,
			SourceTable:  item.SourceTable,
			SourceColumn: item.SourceColumn,
		})
	}
	return response, nil
}

func (s *invoiceService) GenerateInvoicePdf(
	ctx context.Context,
	invoiceID int64,
) (*GenerateInvoicePDFResponse, error) {
	invoiceData, err := s.Store.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"GenerateInvoicePdf",
			"Failed to get invoice data",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to get invoice data: %v", err)
	}

	var invoiceDetails []InvoiceDetails
	if err := json.Unmarshal(invoiceData.InvoiceDetails, &invoiceDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal invoice details: %v", err)
	}

	var senderContacts []SenderContact
	err = json.Unmarshal(invoiceData.SenderContacts, &senderContacts)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"GenerateInvoicePdf",
			"Failed to unmarshal sender contacts",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to unmarshal sender contacts: %v", err)
	}

	var pdfInvoiceDetails []pdf.InvoiceDetail
	for _, detail := range invoiceDetails {
		var pdfInvoicePeriods []pdf.InvoicePeriod
		for _, period := range detail.Periods {
			pdfInvoicePeriods = append(pdfInvoicePeriods, pdf.InvoicePeriod{
				StartDate:             period.StartDate,
				EndDate:               period.EndDate,
				AcommodationTimeFrame: util.DerefString(period.AcommodationTimeFrame),
				AmbulanteTotalMinutes: util.DerefFloat64(period.AmbulanteTotalMinutes),
			})
		}
		pdfInvoiceDetails = append(pdfInvoiceDetails, pdf.InvoiceDetail{
			CareType:      detail.ContractType,
			Periods:       pdfInvoicePeriods,
			Price:         detail.Price,
			PriceTimeUnit: detail.PriceTimeUnit,
			PreVatTotal:   detail.PreVatTotal,
			Total:         detail.Total,
		})
	}

	var extraItems map[string]string
	if err := json.Unmarshal(invoiceData.ExtraContent, &extraItems); err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"GenerateInvoicePdf",
			"Failed to unmarshal extra content",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to unmarshal extra content: %v", err)
	}

	pdfData := pdf.InvoicePDFData{
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
	key, size, err := s.PDFService.GenerateAndUploadInvoicePDF(ctx, pdfData)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"GenerateInvoicePdf",
			"Failed to generate and upload invoice PDF",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to generate and upload invoice PDF: %v", err)
	}

	fileArg := db.CreateAttachmentParams{
		Name: "Invoice_" + invoiceData.InvoiceNumber + ".pdf",
		File: key,
		Size: int32(size),
		Tag:  util.StringPtr("invoice_pdf"),
	}

	attachment, err := s.Store.CreateAttachment(ctx, fileArg)
	if err != nil {
		s.Logger.LogBusinessEvent(
			logger.LogLevelError,
			"GenerateInvoicePdf",
			"Failed to create attachment",
			zap.Error(err),
			zap.Int64("invoice_id", invoiceID),
		)
		return nil, fmt.Errorf("failed to create attachment: %v", err)
	}

	url := s.GenerateResponsePresignedURL(&attachment.File, ctx)

	return &GenerateInvoicePDFResponse{
		FileUrl: *url,
	}, nil
}

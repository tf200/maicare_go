package invoice

import (
	"context"
	"github.com/goccy/go-json"
	"fmt"
	"math"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/pdf"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *invoiceService) CreateInvoice(ctx context.Context, req CreateInvoiceRequest, employeeID uuid.UUID) (*CreateInvoiceResponse, error) {
	if req.ClientID == uuid.Nil {
		return nil, fmt.Errorf("client_id is required")
	}
	if len(req.Lines) == 0 {
		return nil, fmt.Errorf("lines must not be empty")
	}

	client, err := s.Store.GetClientDetails(ctx, req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to load client details: %w", err)
	}
	if client.SenderID == nil || *client.SenderID == uuid.Nil {
		return nil, fmt.Errorf("sender_id is not set for the client")
	}
	senderID := *client.SenderID

	invoiceNumber, invoiceSequence, err := s.GenerateInvoiceNumber(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	if employeeID != uuid.Nil {
		_, err = tx.Exec(ctx, "SELECT set_config('myapp.current_employee_id', $1, true)", employeeID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to set current employee id: %w", err)
		}
	}
	qtx := s.Store.WithTx(tx)

	billingTz := DefaultBillingTimezone
	currency := "EUR"

	inv, err := qtx.CreateInvoice(ctx, db.CreateInvoiceParams{
		InvoiceNumber:     invoiceNumber,
		InvoiceSequence:   invoiceSequence,
		DueDate:           pgtype.Date{Time: req.DueDate, Valid: true},
		IssueDate:         pgtype.Date{Time: req.IssueDate, Valid: true},
		Status:            db.InvoiceStatusEnumConcept,
		InvoiceType:       db.InvoiceTypeEnum(req.InvoiceType),
		Source:            db.InvoiceSourceEnumManual,
		OriginalInvoiceID: nil,
		ReplacesInvoiceID: nil,
		PeriodStart:       pgtype.Timestamptz{Valid: false},
		PeriodEnd:         pgtype.Timestamptz{Valid: false},
		BillingCycle:      nil,
		BillingTimezone:   billingTz,
		Currency:          currency,
		ExtraContent:      util.ParseObjectToJSON(req.ExtraContent),
		ClientID:          req.ClientID,
		SenderID:          senderID,
		WarningCount:      0,
		RunID:             nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	lineNo := int32(1)
	var netTotalCents, vatTotalCents, grossTotalCents int64
	lines := make([]InvoiceLine, 0, len(req.Lines))

	for _, in := range req.Lines {
		netCents := centsFromAmount(in.UnitPrice * in.Quantity)
		vatRatePct := int64(in.VatRate)
		amts := computeVat(netCents, vatRatePct)

		metaBytes, _ := json.Marshal(map[string]any{"manual": true})
		line, err := qtx.CreateInvoiceLine(ctx, db.CreateInvoiceLineParams{
			InvoiceID:   inv.ID,
			ClientID:    inv.ClientID,
			SenderID:    inv.SenderID,
			LineNo:      lineNo,
			LineType:    db.InvoiceLineTypeEnum(in.LineType),
			ContractID:  in.ContractID,
			ServiceType: in.ServiceType,
			Description: in.Description,
			PeriodStart: pgtype.Timestamptz{Time: in.PeriodStart, Valid: true},
			PeriodEnd:   pgtype.Timestamptz{Time: in.PeriodEnd, Valid: true},
			Quantity:    in.Quantity,
			Unit:        in.Unit,
			UnitPrice:   roundTo(in.UnitPrice, 4),
			NetAmount:   amountFromCents(amts.netCents),
			VatRate:     in.VatRate,
			VatAmount:   amountFromCents(amts.vatCents),
			GrossAmount: amountFromCents(amts.grossCents),
			Metadata:    metaBytes,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create invoice line: %w", err)
		}

		lines = append(lines, InvoiceLine{
			ID:          line.ID,
			LineNo:      line.LineNo,
			LineType:    string(line.LineType),
			ContractID:  line.ContractID,
			ServiceType: line.ServiceType,
			Description: line.Description,
			PeriodStart: line.PeriodStart.Time,
			PeriodEnd:   line.PeriodEnd.Time,
			Quantity:    line.Quantity,
			Unit:        line.Unit,
			UnitPrice:   line.UnitPrice,
			NetAmount:   line.NetAmount,
			VatRate:     line.VatRate,
			VatAmount:   line.VatAmount,
			GrossAmount: line.GrossAmount,
		})

		lineNo++
		netTotalCents += amts.netCents
		vatTotalCents += amts.vatCents
		grossTotalCents += amts.grossCents
	}

	netTotal := amountFromCents(netTotalCents)
	vatTotal := amountFromCents(vatTotalCents)
	grossTotal := amountFromCents(grossTotalCents)
	snapshotBytes, _ := json.Marshal(lines)

	updated, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:                inv.ID,
		IssueDate:         pgtype.Date{Valid: false},
		DueDate:           pgtype.Date{Valid: false},
		PeriodStart:       pgtype.Timestamptz{Valid: false},
		PeriodEnd:         pgtype.Timestamptz{Valid: false},
		BillingCycle:      nil,
		BillingTimezone:   &billingTz,
		Source:            db.NullInvoiceSourceEnum{Valid: true, InvoiceSourceEnum: db.InvoiceSourceEnumManual},
		OriginalInvoiceID: nil,
		ReplacesInvoiceID: nil,
		BillToSnapshot:    nil,
		ClientSnapshot:    nil,
		DetailsSnapshot:   snapshotBytes,
		NetTotalAmount:    &netTotal,
		VatTotalAmount:    &vatTotal,
		GrossTotalAmount:  &grossTotal,
		Currency:          &currency,
		ExtraContent:      nil,
		Status:            db.NullInvoiceStatusEnum{Valid: true, InvoiceStatusEnum: db.InvoiceStatusEnumConcept},
		WarningCount:      nil,
		RunID:             nil,
		LockedAt:          pgtype.Timestamptz{Valid: false},
		CalcVersion:       nil,
		CalcMetadata:      nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice totals: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &CreateInvoiceResponse{
		ID:              updated.ID,
		InvoiceNumber:   updated.InvoiceNumber,
		IssueDate:       updated.IssueDate.Time,
		DueDate:         updated.DueDate.Time,
		Status:          string(updated.Status),
		Source:          string(updated.Source),
		InvoiceType:     string(updated.InvoiceType),
		Currency:        updated.Currency,
		NetTotal:        updated.NetTotalAmount,
		VatTotal:        updated.VatTotalAmount,
		GrossTotal:      updated.GrossTotalAmount,
		PdfAttachmentID: updated.PdfAttachmentID,
		ExtraContent:    util.ParseJSONToObject(updated.ExtraContent),
		ClientID:        updated.ClientID,
		SenderID:        updated.SenderID,
		Lines:           lines,
		UpdatedAt:       updated.UpdatedAt.Time,
		CreatedAt:       updated.CreatedAt.Time,
	}, nil
}

func (s *invoiceService) GetInvoiceByID(ctx context.Context, invoiceID uuid.UUID) (*GetInvoiceByIDResponse, error) {
	inv, err := s.Store.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetInvoiceByID", "Failed to get invoice", zap.Error(err))
		return nil, err
	}

	lines, err := s.Store.ListInvoiceLinesByInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoice lines: %w", err)
	}
	lineDTOs := make([]InvoiceLine, 0, len(lines))
	for _, l := range lines {
		lineDTOs = append(lineDTOs, InvoiceLine{
			ID:          l.ID,
			LineNo:      l.LineNo,
			LineType:    string(l.LineType),
			ContractID:  l.ContractID,
			ServiceType: l.ServiceType,
			Description: l.Description,
			PeriodStart: l.PeriodStart.Time,
			PeriodEnd:   l.PeriodEnd.Time,
			Quantity:    l.Quantity,
			Unit:        l.Unit,
			UnitPrice:   l.UnitPrice,
			NetAmount:   l.NetAmount,
			VatRate:     l.VatRate,
			VatAmount:   l.VatAmount,
			GrossAmount: l.GrossAmount,
		})
	}

	paymentCompletionPrc := s.calculatePaymentCompletionPercentage(ctx, inv.GrossTotalAmount, invoiceID)

	var ps, pe *time.Time
	if inv.PeriodStart.Valid {
		t := inv.PeriodStart.Time
		ps = &t
	}
	if inv.PeriodEnd.Valid {
		t := inv.PeriodEnd.Time
		pe = &t
	}

	return &GetInvoiceByIDResponse{
		ID:                   inv.ID,
		InvoiceNumber:        inv.InvoiceNumber,
		IssueDate:            inv.IssueDate.Time,
		DueDate:              inv.DueDate.Time,
		Status:               string(inv.Status),
		Source:               string(inv.Source),
		InvoiceType:          string(inv.InvoiceType),
		OriginalInvoiceID:    inv.OriginalInvoiceID,
		ReplacesInvoiceID:    inv.ReplacesInvoiceID,
		PeriodStart:          ps,
		PeriodEnd:            pe,
		BillingTimezone:      inv.BillingTimezone,
		Currency:             inv.Currency,
		NetTotal:             inv.NetTotalAmount,
		VatTotal:             inv.VatTotalAmount,
		GrossTotal:           inv.GrossTotalAmount,
		ExtraContent:         util.ParseJSONToObject(inv.ExtraContent),
		ClientID:             inv.ClientID,
		SenderID:             inv.SenderID,
		Lines:                lineDTOs,
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

func (s *invoiceService) ListInvoices(ctx *gin.Context, req ListInvoicesRequest) (*pagination.Response[ListInvoicesResponse], error) {
	params := req.GetParams()

	issueStart, issueStartValid, err := util.ParseYYYYMMDD(req.StartDate)
	if err != nil {
		return nil, err
	}
	issueEnd, issueEndValid, err := util.ParseYYYYMMDD(req.EndDate)
	if err != nil {
		return nil, err
	}

	periodStart, periodStartValid, err := util.ParseRFC3339OrYYYYMMDD(req.PeriodStart)
	if err != nil {
		return nil, err
	}
	periodEnd, periodEndValid, err := util.ParseRFC3339OrYYYYMMDD(req.PeriodEnd)
	if err != nil {
		return nil, err
	}

	statuses := []db.InvoiceStatusEnum(nil)
	if len(req.Statuses) > 0 {
		statuses = make([]db.InvoiceStatusEnum, 0, len(req.Statuses))
		for _, st := range req.Statuses {
			statuses = append(statuses, db.InvoiceStatusEnum(st))
		}
	}

	source := db.NullInvoiceSourceEnum{Valid: false}
	if req.Source != nil {
		source = db.NullInvoiceSourceEnum{InvoiceSourceEnum: db.InvoiceSourceEnum(*req.Source), Valid: true}
	}

	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "updated_at"
	}
	sortDir := req.SortDir
	if sortDir == "" {
		sortDir = "desc"
	}

	listParams := db.ListInvoicesParams{
		ClientID:        req.ClientID,
		SenderID:        req.SenderID,
		Status:          db.NullInvoiceStatusFromPtr(req.Status),
		Statuses:        statuses,
		Source:          source,
		InvoiceType:     db.NullInvoiceTypeFromPtr(req.InvoiceType),
		RunID:           req.RunID,
		StartDate:       pgtype.Date{Time: issueStart, Valid: issueStartValid},
		EndDate:         pgtype.Date{Time: issueEnd, Valid: issueEndValid},
		PeriodStart:     pgtype.Timestamptz{Time: periodStart, Valid: periodStartValid},
		PeriodEnd:       pgtype.Timestamptz{Time: periodEnd, Valid: periodEndValid},
		MinWarningCount: req.MinWarningCount,
		Locked:          req.Locked,
		Q:               req.Q,
		SortBy:          sortBy,
		SortDir:         sortDir,
		Limit:           params.Limit,
		Offset:          params.Offset,
	}
	var invoices []db.ListInvoicesRow
	if err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		invoices, err = q.ListInvoices(ctx, listParams)
		return err
	}); err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	if len(invoices) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListInvoicesResponse{}, 0)
		return &pag, nil
	}

	resp := make([]ListInvoicesResponse, 0, len(invoices))
	for _, inv := range invoices {
		resp = append(resp, ListInvoicesResponse{
			ID:            inv.ID,
			InvoiceNumber: inv.InvoiceNumber,
			SenderName:    inv.SenderName,
			IsOverdue:     inv.IsOverdue,

			ClientFirstName:  inv.ClientFirstName,
			ClientLastName:   inv.ClientLastName,
			ClientFilenumber: inv.ClientFilenumber,

			Currency:   inv.Currency,
			GrossTotal: inv.GrossTotalAmount,
			BalanceDue: inv.BalanceDueAmount,
			PaidTotal:  inv.PaidTotalAmount,

			Status:    string(inv.Status),
			IssueDate: inv.IssueDate.Time,
			DueDate:   inv.DueDate.Time,

			ClientID: inv.ClientID,
			SenderID: inv.SenderID,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, resp, invoices[0].TotalCount)
	return &pag, nil
}

func (s *invoiceService) UpdateInvoice(ctx context.Context, invoiceID uuid.UUID, req UpdateInvoiceRequest, employeeID uuid.UUID) (*UpdateInvoiceResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	if employeeID != uuid.Nil {
		_, err = tx.Exec(ctx, "SELECT set_config('myapp.current_employee_id', $1, true)", employeeID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to set current employee id: %w", err)
		}
	}
	qtx := s.Store.WithTx(tx)

	inv, err := qtx.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	var detailsSnapshot []byte
	var netTotalAmount, vatTotalAmount, grossTotalAmount *float64

	if req.Lines != nil {
		if len(req.Lines) == 0 {
			return nil, fmt.Errorf("lines must not be empty")
		}
		if inv.Source == db.InvoiceSourceEnumImported {
			return nil, fmt.Errorf("cannot update invoice lines for imported invoices")
		}
		if inv.InvoiceType == db.InvoiceTypeEnumCreditNote {
			return nil, fmt.Errorf("cannot update invoice lines for credit notes")
		}
		if inv.Status == db.InvoiceStatusEnumCanceled {
			return nil, fmt.Errorf("cannot update invoice lines for canceled invoices")
		}
		if inv.LockedAt.Valid {
			return nil, fmt.Errorf("cannot update invoice lines on a locked invoice")
		}

		paid, err := qtx.GetTotalPaidAmountByInvoice(ctx, invoiceID)
		if err != nil {
			return nil, fmt.Errorf("failed to check payments: %w", err)
		}
		if paid > 0 {
			return nil, fmt.Errorf("cannot update invoice lines when payments exist")
		}

		billedCount, err := qtx.CountBilledCalendarEventsByInvoice(ctx, invoiceID)
		if err != nil {
			return nil, fmt.Errorf("failed to check billed appointments: %w", err)
		}

		linkedCount, err := qtx.CountInvoiceLineCalendarEventsByInvoice(ctx, invoiceID)
		if err != nil {
			return nil, fmt.Errorf("failed to check invoice appointment links: %w", err)
		}
		existingLines, err := qtx.ListInvoiceLinesByInvoice(ctx, invoiceID)
		if err != nil {
			return nil, fmt.Errorf("failed to list existing invoice lines: %w", err)
		}
		if len(existingLines) == 0 {
			return nil, fmt.Errorf("invoice has no existing lines to update")
		}

		sameUUIDPtr := func(a, b *uuid.UUID) bool {
			if a == nil && b == nil {
				return true
			}
			if a == nil || b == nil {
				return false
			}
			return *a == *b
		}

		var netTotalCents, vatTotalCents, grossTotalCents int64
		snapshotLines := make([]InvoiceLine, 0, len(req.Lines))

		// If the invoice has appointment-linked lines, we must keep the existing invoice_line IDs.
		// That means: no add/remove/reorder, only edit values in-place.
		if billedCount > 0 || linkedCount > 0 {
			if len(req.Lines) != len(existingLines) {
				return nil, fmt.Errorf("cannot add/remove/reorder lines for invoices with appointment links")
			}

			for i, in := range req.Lines {
				cur := existingLines[i]
				if in.LineType != string(cur.LineType) {
					return nil, fmt.Errorf("line %d: line_type cannot be changed", i+1)
				}
				if !sameUUIDPtr(in.ContractID, cur.ContractID) {
					return nil, fmt.Errorf("line %d: contract_id cannot be changed", i+1)
				}
				if in.ServiceType != cur.ServiceType {
					return nil, fmt.Errorf("line %d: service_type cannot be changed", i+1)
				}

				netCents := centsFromAmount(in.UnitPrice * in.Quantity)
				vatRatePct := int64(in.VatRate)
				amts := computeVat(netCents, vatRatePct)

				unitPrice := roundTo(in.UnitPrice, 4)
				periodStart := pgtype.Timestamptz{Time: in.PeriodStart, Valid: !in.PeriodStart.IsZero()}
				periodEnd := pgtype.Timestamptz{Time: in.PeriodEnd, Valid: !in.PeriodEnd.IsZero()}
				description := in.Description
				unit := in.Unit
				qty := in.Quantity
				netAmt := amountFromCents(amts.netCents)
				vatRate := in.VatRate
				vatAmt := amountFromCents(amts.vatCents)
				grossAmt := amountFromCents(amts.grossCents)

				updatedLine, err := qtx.UpdateInvoiceLine(ctx, db.UpdateInvoiceLineParams{
					LineType:    db.NullInvoiceLineTypeEnum{Valid: false},
					ContractID:  nil,
					ServiceType: nil,
					Description: &description,
					PeriodStart: periodStart,
					PeriodEnd:   periodEnd,
					Quantity:    &qty,
					Unit:        &unit,
					UnitPrice:   &unitPrice,
					NetAmount:   &netAmt,
					VatRate:     &vatRate,
					VatAmount:   &vatAmt,
					GrossAmount: &grossAmt,
					Metadata:    nil,
					ID:          cur.ID,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to update invoice line: %w", err)
				}

				snapshotLines = append(snapshotLines, InvoiceLine{
					ID:          updatedLine.ID,
					LineNo:      updatedLine.LineNo,
					LineType:    string(updatedLine.LineType),
					ContractID:  updatedLine.ContractID,
					ServiceType: updatedLine.ServiceType,
					Description: updatedLine.Description,
					PeriodStart: updatedLine.PeriodStart.Time,
					PeriodEnd:   updatedLine.PeriodEnd.Time,
					Quantity:    updatedLine.Quantity,
					Unit:        updatedLine.Unit,
					UnitPrice:   updatedLine.UnitPrice,
					NetAmount:   updatedLine.NetAmount,
					VatRate:     updatedLine.VatRate,
					VatAmount:   updatedLine.VatAmount,
					GrossAmount: updatedLine.GrossAmount,
				})

				netTotalCents += amts.netCents
				vatTotalCents += amts.vatCents
				grossTotalCents += amts.grossCents
			}
		} else {
			// No appointment linkage: safe to fully replace the invoice lines.
			if err := qtx.DeleteInvoiceLinesByInvoice(ctx, invoiceID); err != nil {
				return nil, fmt.Errorf("failed to delete invoice lines: %w", err)
			}

			lineNo := int32(1)
			for _, in := range req.Lines {
				netCents := centsFromAmount(in.UnitPrice * in.Quantity)
				vatRatePct := int64(in.VatRate)
				amts := computeVat(netCents, vatRatePct)

				metaBytes := []byte("{}")
				if inv.Source == db.InvoiceSourceEnumManual {
					metaBytes, _ = json.Marshal(map[string]any{"manual": true, "updated": true})
				}
				periodStart := pgtype.Timestamptz{Time: in.PeriodStart, Valid: !in.PeriodStart.IsZero()}
				periodEnd := pgtype.Timestamptz{Time: in.PeriodEnd, Valid: !in.PeriodEnd.IsZero()}
				line, err := qtx.CreateInvoiceLine(ctx, db.CreateInvoiceLineParams{
					InvoiceID:   inv.ID,
					ClientID:    inv.ClientID,
					SenderID:    inv.SenderID,
					LineNo:      lineNo,
					LineType:    db.InvoiceLineTypeEnum(in.LineType),
					ContractID:  in.ContractID,
					ServiceType: in.ServiceType,
					Description: in.Description,
					PeriodStart: periodStart,
					PeriodEnd:   periodEnd,
					Quantity:    in.Quantity,
					Unit:        in.Unit,
					UnitPrice:   roundTo(in.UnitPrice, 4),
					NetAmount:   amountFromCents(amts.netCents),
					VatRate:     in.VatRate,
					VatAmount:   amountFromCents(amts.vatCents),
					GrossAmount: amountFromCents(amts.grossCents),
					Metadata:    metaBytes,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to create invoice line: %w", err)
				}

				snapshotLines = append(snapshotLines, InvoiceLine{
					ID:          line.ID,
					LineNo:      line.LineNo,
					LineType:    string(line.LineType),
					ContractID:  line.ContractID,
					ServiceType: line.ServiceType,
					Description: line.Description,
					PeriodStart: line.PeriodStart.Time,
					PeriodEnd:   line.PeriodEnd.Time,
					Quantity:    line.Quantity,
					Unit:        line.Unit,
					UnitPrice:   line.UnitPrice,
					NetAmount:   line.NetAmount,
					VatRate:     line.VatRate,
					VatAmount:   line.VatAmount,
					GrossAmount: line.GrossAmount,
				})

				lineNo++
				netTotalCents += amts.netCents
				vatTotalCents += amts.vatCents
				grossTotalCents += amts.grossCents
			}
		}

		netTotal := amountFromCents(netTotalCents)
		vatTotal := amountFromCents(vatTotalCents)
		grossTotal := amountFromCents(grossTotalCents)
		detailsSnapshot, _ = json.Marshal(snapshotLines)
		netTotalAmount = &netTotal
		vatTotalAmount = &vatTotal
		grossTotalAmount = &grossTotal
	}

	var locked pgtype.Timestamptz
	if req.LockedAt != nil {
		locked = pgtype.Timestamptz{Time: *req.LockedAt, Valid: true}
	} else {
		locked = pgtype.Timestamptz{Valid: false}
	}

	updated, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:               invoiceID,
		IssueDate:        pgtype.Date{Time: req.IssueDate, Valid: !req.IssueDate.IsZero()},
		DueDate:          pgtype.Date{Time: req.DueDate, Valid: !req.DueDate.IsZero()},
		PeriodStart:      pgtype.Timestamptz{Valid: false},
		PeriodEnd:        pgtype.Timestamptz{Valid: false},
		BillingCycle:     nil,
		BillingTimezone:  nil,
		Source:           db.NullInvoiceSourceEnum{Valid: false},
		BillToSnapshot:   nil,
		ClientSnapshot:   nil,
		DetailsSnapshot:  detailsSnapshot,
		NetTotalAmount:   netTotalAmount,
		VatTotalAmount:   vatTotalAmount,
		GrossTotalAmount: grossTotalAmount,
		Currency:         nil,
		ExtraContent:     util.ParseObjectToJSON(req.ExtraContent),
		Status:           db.NullInvoiceStatusEnum{Valid: req.Status != "", InvoiceStatusEnum: db.InvoiceStatusEnum(req.Status)},
		WarningCount:     req.WarningCount,
		RunID:            nil,
		LockedAt:         locked,
		CalcVersion:      nil,
		CalcMetadata:     nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return &UpdateInvoiceResponse{
		ID:              updated.ID,
		InvoiceNumber:   updated.InvoiceNumber,
		IssueDate:       updated.IssueDate.Time,
		DueDate:         updated.DueDate.Time,
		Status:          string(updated.Status),
		Currency:        updated.Currency,
		NetTotal:        updated.NetTotalAmount,
		VatTotal:        updated.VatTotalAmount,
		GrossTotal:      updated.GrossTotalAmount,
		PdfAttachmentID: updated.PdfAttachmentID,
		ExtraContent:    util.ParseJSONToObject(updated.ExtraContent),
		ClientID:        updated.ClientID,
		SenderID:        updated.SenderID,
		UpdatedAt:       updated.UpdatedAt.Time,
		CreatedAt:       updated.CreatedAt.Time,
	}, nil
}

func (s *invoiceService) DeleteInvoice(ctx context.Context, invoiceID uuid.UUID) error {
	if err := s.Store.DeleteInvoice(ctx, invoiceID); err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}
	return nil
}

func (s *invoiceService) GetInvoiceAuditLogs(ctx context.Context, invoiceID uuid.UUID) ([]GetInvoiceAuditLogResponse, error) {
	logs, err := s.Store.GetInvoiceAuditLogs(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice audit logs: %w", err)
	}

	out := make([]GetInvoiceAuditLogResponse, 0, len(logs))
	for _, l := range logs {
		out = append(out, GetInvoiceAuditLogResponse{
			AuditID:            l.AuditID,
			InvoiceID:          l.InvoiceID,
			Operation:          string(l.Operation),
			ChangedBy:          l.ChangedBy,
			ChangedAt:          l.ChangedAt.Time,
			OldValues:          util.ParseJSONToObject(l.OldValues),
			NewValues:          util.ParseJSONToObject(l.NewValues),
			ChangedFields:      l.ChangedFields,
			ChangedByFirstName: l.ChangedByFirstName,
			ChangedByLastName:  l.ChangedByLastName,
		})
	}
	return out, nil
}

func (s *invoiceService) GetInvoiceTemplateItemsApi(ctx context.Context) ([]GetInvoiceTemplateItemsResponse, error) {
	items, err := s.Store.GetAllTemplateItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get template items: %w", err)
	}

	out := make([]GetInvoiceTemplateItemsResponse, 0, len(items))
	for _, it := range items {
		out = append(out, GetInvoiceTemplateItemsResponse{
			ID:           it.ID,
			ItemTag:      it.ItemTag,
			Description:  it.Description,
			SourceTable:  it.SourceTable,
			SourceColumn: it.SourceColumn,
		})
	}
	return out, nil
}

func (s *invoiceService) GenerateInvoicePdf(ctx context.Context, invoiceID uuid.UUID) (*GenerateInvoicePDFResponse, error) {
	inv, err := s.Store.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// If we already have a stored PDF attachment, just return a presigned URL.
	if inv.PdfAttachmentID != nil {
		att, err := s.Store.GetAttachmentById(ctx, *inv.PdfAttachmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get pdf attachment: %w", err)
		}
		url, err := s.B2Client.GeneratePresignedURL(ctx, att.File, 15*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned url: %w", err)
		}
		return &GenerateInvoicePDFResponse{FileUrl: url}, nil
	}

	lines, err := s.Store.ListInvoiceLinesByInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoice lines: %w", err)
	}

	details := make([]pdf.InvoiceDetail, 0, len(lines))
	for _, l := range lines {
		details = append(details, pdf.InvoiceDetail{
			CareType:      l.ServiceType,
			Periods:       []pdf.InvoicePeriod{{StartDate: l.PeriodStart.Time, EndDate: l.PeriodEnd.Time}},
			Price:         l.UnitPrice,
			PriceTimeUnit: l.Unit,
			PreVatTotal:   l.NetAmount,
			Total:         l.GrossAmount,
		})
	}

	extra := map[string]string{}
	if inv.ExtraContent != nil {
		_ = json.Unmarshal(inv.ExtraContent, &extra)
	}

	pdfData := pdf.InvoicePDFData{
		ID: inv.ID,
		SenderName: func() string {
			if inv.SenderName != nil {
				return *inv.SenderName
			}
			return ""
		}(),
		SenderStreet: func() string {
			if inv.SenderStreet != nil {
				return *inv.SenderStreet
			}
			return ""
		}(),
		SenderHouseNumber: func() string {
			if inv.SenderHouseNumber != nil {
				return *inv.SenderHouseNumber
			}
			return ""
		}(),
		SenderPostalCode: func() string {
			if inv.SenderPostalCode != nil {
				return *inv.SenderPostalCode
			}
			return ""
		}(),
		SenderCity: func() string {
			if inv.SenderCity != nil {
				return *inv.SenderCity
			}
			return ""
		}(),
		InvoiceNumber:  inv.InvoiceNumber,
		InvoiceDate:    inv.IssueDate.Time,
		DueDate:        inv.DueDate.Time,
		InvoiceDetails: details,
		TotalAmount:    inv.GrossTotalAmount,
		ExtraItems:     extra,
	}

	// Generate and upload PDF.
	// Note: pdf service returns an object key, not a ready-to-use URL.
	objKey, size, err := s.PDFService.GenerateAndUploadInvoicePDF(ctx, pdfData)
	if err != nil {
		return nil, err
	}

	// If invoice is immutable, persist linkage so future requests don't regenerate.
	immutable := inv.LockedAt.Valid || inv.Status != db.InvoiceStatusEnumConcept
	if immutable {
		if size > math.MaxInt32 {
			return nil, fmt.Errorf("pdf too large to store in attachment_file: %d bytes", size)
		}

		attID := uuid.New()
		name := fmt.Sprintf("invoice_%s.pdf", inv.InvoiceNumber)
		tag := "invoice_pdf"

		tx, err := s.Store.ConnPool.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)
		qtx := s.Store.WithTx(tx)

		_, err = qtx.CreateAttachment(ctx, db.CreateAttachmentParams{
			Uuid: attID,
			Name: name,
			File: objKey,
			Size: int32(size),
			Tag:  &tag,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create attachment: %w", err)
		}

		_, err = qtx.SetAttachmentAsUsedorUnused(ctx, db.SetAttachmentAsUsedorUnusedParams{
			Uuid:   attID,
			IsUsed: true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to mark attachment as used: %w", err)
		}

		_, err = qtx.InsertIncoicePdfUrl(ctx, db.InsertIncoicePdfUrlParams{
			ID:              invoiceID,
			PdfAttachmentID: &attID,
		})
		if err != nil {
			// Another request may have stored a PDF concurrently. In that case, return the stored one.
			if err == pgx.ErrNoRows {
				latest, gErr := qtx.GetInvoice(ctx, invoiceID)
				if gErr == nil && latest.PdfAttachmentID != nil {
					att, aErr := qtx.GetAttachmentById(ctx, *latest.PdfAttachmentID)
					if aErr == nil {
						_ = tx.Commit(ctx)
						url, uErr := s.B2Client.GeneratePresignedURL(ctx, att.File, 15*time.Minute)
						if uErr == nil {
							return &GenerateInvoicePDFResponse{FileUrl: url}, nil
						}
					}
				}
			}
			return nil, fmt.Errorf("failed to set invoice pdf attachment id: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	url, err := s.B2Client.GeneratePresignedURL(ctx, objKey, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned url: %w", err)
	}
	return &GenerateInvoicePDFResponse{FileUrl: url}, nil
}

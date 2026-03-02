package invoice

import (
	"context"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// CreditInvoiceResponse represents the response body for crediting an invoice.
type CreditInvoiceResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *invoiceService) CreditInvoice(ctx context.Context, invoiceID uuid.UUID, employeeID uuid.UUID) (*CreditInvoiceResponse, error) {
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

	original, err := qtx.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreditInvoice", "Failed to get original invoice", zap.Error(err))
		return nil, fmt.Errorf("failed to get original invoice: %w", err)
	}
	if original.InvoiceType == db.InvoiceTypeEnumCreditNote {
		return nil, fmt.Errorf("cannot credit a credit note")
	}

	origLines, err := qtx.ListInvoiceLinesByInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list original invoice lines: %w", err)
	}
	if len(origLines) == 0 {
		return nil, fmt.Errorf("original invoice has no lines to credit")
	}

	creditNumber, creditSeq, err := s.GenerateInvoiceNumber(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	creditInv, err := qtx.CreateInvoice(ctx, db.CreateInvoiceParams{
		InvoiceNumber:     creditNumber,
		InvoiceSequence:   creditSeq,
		DueDate:           pgtype.Date{Time: now.Add(30 * 24 * time.Hour), Valid: true},
		IssueDate:         pgtype.Date{Time: now, Valid: true},
		Status:            db.InvoiceStatusEnumConcept,
		InvoiceType:       db.InvoiceTypeEnumCreditNote,
		Source:            original.Source,
		OriginalInvoiceID: &original.ID,
		ReplacesInvoiceID: nil,
		PeriodStart:       original.PeriodStart,
		PeriodEnd:         original.PeriodEnd,
		BillingCycle:      original.BillingCycle,
		BillingTimezone:   original.BillingTimezone,
		Currency:          original.Currency,
		ExtraContent:      []byte("{}"),
		ClientID:          original.ClientID,
		SenderID:          original.SenderID,
		WarningCount:      0,
		RunID:             nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create credit invoice: %w", err)
	}

	var netTotalCents, vatTotalCents, grossTotalCents int64
	lineNo := int32(1)

	for _, l := range origLines {
		netCents := -centsFromAmount(l.NetAmount)
		vatCents := -centsFromAmount(l.VatAmount)
		grossCents := -centsFromAmount(l.GrossAmount)

		_, err := qtx.CreateInvoiceLine(ctx, db.CreateInvoiceLineParams{
			InvoiceID:   creditInv.ID,
			ClientID:    creditInv.ClientID,
			SenderID:    creditInv.SenderID,
			LineNo:      lineNo,
			LineType:    l.LineType,
			ContractID:  l.ContractID,
			ServiceType: l.ServiceType,
			Description: fmt.Sprintf("CREDIT: %s", l.Description),
			PeriodStart: l.PeriodStart,
			PeriodEnd:   l.PeriodEnd,
			Quantity:    -l.Quantity,
			Unit:        l.Unit,
			UnitPrice:   l.UnitPrice,
			NetAmount:   amountFromCents(netCents),
			VatRate:     l.VatRate,
			VatAmount:   amountFromCents(vatCents),
			GrossAmount: amountFromCents(grossCents),
			Metadata:    l.Metadata,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create credit invoice line: %w", err)
		}

		lineNo++
		netTotalCents += netCents
		vatTotalCents += vatCents
		grossTotalCents += grossCents
	}

	netTotal := amountFromCents(netTotalCents)
	vatTotal := amountFromCents(vatTotalCents)
	grossTotal := amountFromCents(grossTotalCents)

	// Void billed appointments from the original invoice to allow re-billing if needed.
	if err := qtx.VoidBilledCalendarEventsByInvoice(ctx, original.ID); err != nil {
		return nil, fmt.Errorf("failed to void billed appointments: %w", err)
	}

	// Mark the original invoice as canceled and link the credit note.
	if _, err := qtx.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{ID: original.ID, Status: db.InvoiceStatusEnumCanceled}); err != nil {
		return nil, fmt.Errorf("failed to cancel original invoice: %w", err)
	}

	// Set credit note totals.
	_, err = qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:                creditInv.ID,
		IssueDate:         pgtype.Date{Valid: false},
		DueDate:           pgtype.Date{Valid: false},
		PeriodStart:       pgtype.Timestamptz{Valid: false},
		PeriodEnd:         pgtype.Timestamptz{Valid: false},
		BillingCycle:      nil,
		BillingTimezone:   nil,
		Source:            db.NullInvoiceSourceEnum{Valid: false},
		OriginalInvoiceID: nil,
		ReplacesInvoiceID: nil,
		BillToSnapshot:    nil,
		ClientSnapshot:    nil,
		DetailsSnapshot:   nil,
		NetTotalAmount:    &netTotal,
		VatTotalAmount:    &vatTotal,
		GrossTotalAmount:  &grossTotal,
		Currency:          nil,
		ExtraContent:      nil,
		Status:            db.NullInvoiceStatusEnum{Valid: false},
		WarningCount:      nil,
		RunID:             nil,
		LockedAt:          pgtype.Timestamptz{Valid: false},
		CalcVersion:       nil,
		CalcMetadata:      nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update credit invoice totals: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit credit transaction: %w", err)
	}

	return &CreditInvoiceResponse{ID: creditInv.ID}, nil
}

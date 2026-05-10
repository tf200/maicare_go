package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

const (
	defaultBillingTimezone = "Europe/Amsterdam"
	defaultBillingCycle    = "iso_4_week"
	paymentTolerance       = 50.0
)

var (
	ErrNoBillableItems          = errors.New("no billable items")
	ErrAutoInvoiceAlreadyExists = errors.New("auto invoice already exists for period")
)

type InvoiceService struct {
	store      *db.Store
	logger     domain.Logger
	storage    domain.Storage
	pdfService domain.PDFService
}

func NewInvoiceService(store *db.Store, logger domain.Logger, storage domain.Storage, pdfService domain.PDFService) *InvoiceService {
	return &InvoiceService{
		store:      store,
		logger:     logger,
		storage:    storage,
		pdfService: pdfService,
	}
}

// ==================== Helpers ====================

func centsFromAmount(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

func amountFromCents(cents int64) float64 {
	return float64(cents) / 100
}

func roundHalfUpDiv(n, d int64) int64 {
	if d == 0 {
		panic("division by zero")
	}
	if n == 0 {
		return 0
	}
	sign := int64(1)
	if n < 0 {
		sign = -1
		n = -n
	}
	q := n / d
	r := n % d
	if 2*r >= d {
		q++
	}
	return sign * q
}

func clampTime(t, min, max time.Time) time.Time {
	if t.Before(min) {
		return min
	}
	if t.After(max) {
		return max
	}
	return t
}

func countCalendarDays(start, end time.Time, tz string) (int64, error) {
	if end.Before(start) || end.Equal(start) {
		return 0, nil
	}
	if tz == "" {
		tz = defaultBillingTimezone
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return 0, fmt.Errorf("invalid billing timezone %q: %w", tz, err)
	}
	sy, sm, sd := start.In(loc).Date()
	ey, em, ed := end.In(loc).Date()
	sUTC := time.Date(sy, sm, sd, 0, 0, 0, 0, time.UTC)
	eUTC := time.Date(ey, em, ed, 0, 0, 0, 0, time.UTC)
	if eUTC.Before(sUTC) {
		return 0, nil
	}
	return int64(eUTC.Sub(sUTC).Hours() / 24), nil
}

type lineAmounts struct {
	netCents   int64
	vatCents   int64
	grossCents int64
}

func computeVat(netCents int64, vatRatePct int64) lineAmounts {
	vatCents := roundHalfUpDiv(netCents*vatRatePct, 100)
	return lineAmounts{
		netCents:   netCents,
		vatCents:   vatCents,
		grossCents: netCents + vatCents,
	}
}

func roundTo(amount float64, decimals int) float64 {
	pow := math.Pow10(decimals)
	return math.Round(amount*pow) / pow
}

func toInvoiceLine(l db.InvoiceLine) domain.InvoiceLine {
	return domain.InvoiceLine{
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
	}
}

func toInvoiceLines(lines []db.InvoiceLine) []domain.InvoiceLine {
	out := make([]domain.InvoiceLine, 0, len(lines))
	for _, l := range lines {
		out = append(out, toInvoiceLine(l))
	}
	return out
}

func invoiceFromRow(inv db.Invoice) domain.Invoice {
	var ps, pe *time.Time
	if inv.PeriodStart.Valid {
		t := inv.PeriodStart.Time
		ps = &t
	}
	if inv.PeriodEnd.Valid {
		t := inv.PeriodEnd.Time
		pe = &t
	}
	var locked *time.Time
	if inv.LockedAt.Valid {
		t := inv.LockedAt.Time
		locked = &t
	}
	return domain.Invoice{
		ID:                inv.ID,
		InvoiceNumber:     inv.InvoiceNumber,
		IssueDate:         inv.IssueDate.Time,
		DueDate:           inv.DueDate.Time,
		Status:            string(inv.Status),
		Source:            string(inv.Source),
		InvoiceType:       string(inv.InvoiceType),
		OriginalInvoiceID: inv.OriginalInvoiceID,
		ReplacesInvoiceID: inv.ReplacesInvoiceID,
		PeriodStart:       ps,
		PeriodEnd:         pe,
		BillingTimezone:   inv.BillingTimezone,
		Currency:          inv.Currency,
		NetTotal:          inv.NetTotalAmount,
		VatTotal:          inv.VatTotalAmount,
		GrossTotal:        inv.GrossTotalAmount,
		PdfAttachmentID:   inv.PdfAttachmentID,
		ExtraContent:      inv.ExtraContent,
		ClientID:          inv.ClientID,
		SenderID:          inv.SenderID,
		WarningCount:      inv.WarningCount,
		LockedAt:          locked,
		RunID:             inv.RunID,
		UpdatedAt:         inv.UpdatedAt.Time,
		CreatedAt:         inv.CreatedAt.Time,
	}
}

func invoiceFromGetRow(inv db.GetInvoiceRow) domain.Invoice {
	var ps, pe *time.Time
	if inv.PeriodStart.Valid {
		t := inv.PeriodStart.Time
		ps = &t
	}
	if inv.PeriodEnd.Valid {
		t := inv.PeriodEnd.Time
		pe = &t
	}
	var locked *time.Time
	if inv.LockedAt.Valid {
		t := inv.LockedAt.Time
		locked = &t
	}
	return domain.Invoice{
		ID:                inv.ID,
		InvoiceNumber:     inv.InvoiceNumber,
		IssueDate:         inv.IssueDate.Time,
		DueDate:           inv.DueDate.Time,
		Status:            string(inv.Status),
		Source:            string(inv.Source),
		InvoiceType:       string(inv.InvoiceType),
		OriginalInvoiceID: inv.OriginalInvoiceID,
		ReplacesInvoiceID: inv.ReplacesInvoiceID,
		PeriodStart:       ps,
		PeriodEnd:         pe,
		BillingTimezone:   inv.BillingTimezone,
		Currency:          inv.Currency,
		NetTotal:          inv.NetTotalAmount,
		VatTotal:          inv.VatTotalAmount,
		GrossTotal:        inv.GrossTotalAmount,
		PdfAttachmentID:   inv.PdfAttachmentID,
		ExtraContent:      inv.ExtraContent,
		ClientID:          inv.ClientID,
		SenderID:          inv.SenderID,
		WarningCount:      inv.WarningCount,
		LockedAt:          locked,
		RunID:             inv.RunID,
		SenderName:        inv.SenderName,
		SenderStreet:      inv.SenderStreet,
		SenderHouseNumber: inv.SenderHouseNumber,
		SenderPostalCode:  inv.SenderPostalCode,
		SenderCity:        inv.SenderCity,
		ClientFirstName:   inv.ClientFirstName,
		ClientLastName:    inv.ClientLastName,
		UpdatedAt:         inv.UpdatedAt.Time,
		CreatedAt:         inv.CreatedAt.Time,
	}
}

// ==================== Invoice Number Generation ====================

func (s *InvoiceService) generateInvoiceNumber(ctx context.Context) (string, int64, error) {
	now := time.Now()
	datePart := now.Format("20060102")
	maxSeq, err := s.store.GetMaxInvoiceSequenceForDate(ctx, now)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get max invoice sequence: %w", err)
	}
	nextSeq := maxSeq + 1
	return fmt.Sprintf("INV-%s-%04d", datePart, nextSeq), nextSeq, nil
}

// ==================== Create Invoice ====================

func (s *InvoiceService) CreateInvoice(ctx context.Context, params domain.CreateInvoiceParams, employeeID uuid.UUID) (*domain.Invoice, []domain.InvoiceLine, error) {
	if params.ClientID == uuid.Nil {
		return nil, nil, fmt.Errorf("client_id is required")
	}
	if len(params.Lines) == 0 {
		return nil, nil, fmt.Errorf("lines must not be empty")
	}

	client, err := s.store.GetClientDetails(ctx, params.ClientID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load client details: %w", err)
	}
	if client.SenderID == nil || *client.SenderID == uuid.Nil {
		return nil, nil, fmt.Errorf("sender_id is not set for the client")
	}
	senderID := *client.SenderID

	invoiceNumber, invoiceSequence, err := s.generateInvoiceNumber(ctx)
	if err != nil {
		return nil, nil, err
	}

	tx, err := s.store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	if employeeID != uuid.Nil {
		_, err = tx.Exec(ctx, "SELECT set_config('myapp.current_employee_id', $1, true)", employeeID.String())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to set current employee id: %w", err)
		}
	}
	qtx := s.store.WithTx(tx)

	billingTz := defaultBillingTimezone
	currency := "EUR"
	extraContent := []byte("{}")
	if len(params.ExtraContent) > 0 {
		extraContent = params.ExtraContent
	}

	inv, err := qtx.CreateInvoice(ctx, db.CreateInvoiceParams{
		InvoiceNumber:     invoiceNumber,
		InvoiceSequence:   invoiceSequence,
		DueDate:           pgtype.Date{Time: params.DueDate, Valid: true},
		IssueDate:         pgtype.Date{Time: params.IssueDate, Valid: true},
		Status:            db.InvoiceStatusEnumConcept,
		InvoiceType:       db.InvoiceTypeEnum(params.InvoiceType),
		Source:            db.InvoiceSourceEnumManual,
		OriginalInvoiceID: nil,
		ReplacesInvoiceID: nil,
		PeriodStart:       pgtype.Timestamptz{Valid: false},
		PeriodEnd:         pgtype.Timestamptz{Valid: false},
		BillingCycle:      nil,
		BillingTimezone:   billingTz,
		Currency:          currency,
		ExtraContent:      extraContent,
		ClientID:          params.ClientID,
		SenderID:          senderID,
		WarningCount:      0,
		RunID:             nil,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	lineNo := int32(1)
	var netTotalCents, vatTotalCents, grossTotalCents int64
	lines := make([]domain.InvoiceLine, 0, len(params.Lines))

	for _, in := range params.Lines {
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
			return nil, nil, fmt.Errorf("failed to create invoice line: %w", err)
		}

		lines = append(lines, toInvoiceLine(line))
		lineNo++
		netTotalCents += amts.netCents
		vatTotalCents += amts.vatCents
		grossTotalCents += amts.grossCents
	}

	netTotal := amountFromCents(netTotalCents)
	vatTotal := amountFromCents(vatTotalCents)
	grossTotal := amountFromCents(grossTotalCents)
	snapshotBytes, _ := json.Marshal(lines)

	manualSource := db.InvoiceSourceEnumManual
	conceptStatus := db.InvoiceStatusEnumConcept
	updated, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:                inv.ID,
		IssueDate:         pgtype.Date{Valid: false},
		DueDate:           pgtype.Date{Valid: false},
		PeriodStart:       pgtype.Timestamptz{Valid: false},
		PeriodEnd:         pgtype.Timestamptz{Valid: false},
		BillingCycle:      nil,
		BillingTimezone:   &billingTz,
		Source:            &manualSource,
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
		Status:            &conceptStatus,
		WarningCount:      nil,
		RunID:             nil,
		LockedAt:          pgtype.Timestamptz{Valid: false},
		CalcVersion:       nil,
		CalcMetadata:      nil,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update invoice totals: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	d := invoiceFromRow(updated)
	return &d, lines, nil
}

// ==================== Get Invoice By ID ====================

func (s *InvoiceService) GetInvoiceByID(ctx context.Context, invoiceID uuid.UUID) (*domain.Invoice, []domain.InvoiceLine, float64, error) {
	inv, err := s.store.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.GetInvoiceByID", "failed to get invoice", err, zap.String("invoice_id", invoiceID.String()))
		return nil, nil, 0, err
	}

	lines, err := s.store.ListInvoiceLinesByInvoice(ctx, invoiceID)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to list invoice lines: %w", err)
	}
	lineDTOs := toInvoiceLines(lines)

	paymentCompletionPrc := s.calculatePaymentCompletionPercentage(ctx, inv.GrossTotalAmount, invoiceID)

	d := invoiceFromGetRow(inv)
	return &d, lineDTOs, paymentCompletionPrc, nil
}

// ==================== List Invoices ====================

func (s *InvoiceService) ListInvoices(ctx context.Context, params domain.ListInvoicesParams) (*domain.ListResult[domain.InvoiceListItem], error) {
	issueDate := pgtype.Date{Valid: false}
	if params.StartDate != "" {
		t, err := time.Parse("2006-01-02", params.StartDate)
		if err == nil {
			issueDate = pgtype.Date{Time: t, Valid: true}
		}
	}
	endDate := pgtype.Date{Valid: false}
	if params.EndDate != "" {
		t, err := time.Parse("2006-01-02", params.EndDate)
		if err == nil {
			endDate = pgtype.Date{Time: t, Valid: true}
		}
	}

	periodStart := pgtype.Timestamptz{Valid: false}
	if params.PeriodStart != "" {
		t, err := time.Parse(time.RFC3339, params.PeriodStart)
		if err == nil {
			periodStart = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}
	periodEnd := pgtype.Timestamptz{Valid: false}
	if params.PeriodEnd != "" {
		t, err := time.Parse(time.RFC3339, params.PeriodEnd)
		if err == nil {
			periodEnd = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}

	statuses := []db.InvoiceStatusEnum(nil)
	if len(params.Statuses) > 0 {
		statuses = make([]db.InvoiceStatusEnum, 0, len(params.Statuses))
		for _, st := range params.Statuses {
			statuses = append(statuses, db.InvoiceStatusEnum(st))
		}
	}

	var source *db.InvoiceSourceEnum
	if params.Source != nil {
		value := db.InvoiceSourceEnum(*params.Source)
		source = &value
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "updated_at"
	}
	sortDir := params.SortDir
	if sortDir == "" {
		sortDir = "desc"
	}

	listParams := db.ListInvoicesParams{
		ClientID:        params.ClientID,
		SenderID:        params.SenderID,
		Status:          db.NullInvoiceStatusFromPtr(params.Status),
		Statuses:        statuses,
		Source:          source,
		InvoiceType:     db.NullInvoiceTypeFromPtr(params.InvoiceType),
		RunID:           params.RunID,
		StartDate:       issueDate,
		EndDate:         endDate,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		MinWarningCount: params.MinWarningCount,
		Locked:          params.Locked,
		Q:               params.Q,
		SortBy:          sortBy,
		SortDir:         sortDir,
		Limit:           params.Limit,
		Offset:          params.Offset,
	}

	invoices, err := s.store.ListInvoices(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	if len(invoices) == 0 {
		return &domain.ListResult[domain.InvoiceListItem]{Items: []domain.InvoiceListItem{}, TotalCount: 0}, nil
	}

	resp := make([]domain.InvoiceListItem, 0, len(invoices))
	for _, inv := range invoices {
		resp = append(resp, domain.InvoiceListItem{
			ID:               inv.ID,
			InvoiceNumber:    inv.InvoiceNumber,
			SenderName:       inv.SenderName,
			IsOverdue:        inv.IsOverdue,
			ClientFirstName:  inv.ClientFirstName,
			ClientLastName:   inv.ClientLastName,
			ClientFilenumber: inv.ClientFilenumber,
			Currency:         inv.Currency,
			GrossTotal:       inv.GrossTotalAmount,
			BalanceDue:       inv.BalanceDueAmount,
			PaidTotal:        inv.PaidTotalAmount,
			Status:           string(inv.Status),
			IssueDate:        inv.IssueDate.Time,
			DueDate:          inv.DueDate.Time,
			ClientID:         inv.ClientID,
			SenderID:         inv.SenderID,
		})
	}

	return &domain.ListResult[domain.InvoiceListItem]{Items: resp, TotalCount: invoices[0].TotalCount}, nil
}

// ==================== Update Invoice ====================

func (s *InvoiceService) UpdateInvoice(ctx context.Context, invoiceID uuid.UUID, employeeID uuid.UUID, params domain.CreateInvoiceParams) (*domain.Invoice, error) {
	tx, err := s.store.ConnPool.Begin(ctx)
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
	qtx := s.store.WithTx(tx)

	inv, err := qtx.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	var detailsSnapshot []byte
	var netTotalAmount, vatTotalAmount, grossTotalAmount *float64

	if params.Lines != nil {
		if len(params.Lines) == 0 {
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
		snapshotLines := make([]domain.InvoiceLine, 0, len(params.Lines))

		if billedCount > 0 || linkedCount > 0 {
			if len(params.Lines) != len(existingLines) {
				return nil, fmt.Errorf("cannot add/remove/reorder lines for invoices with appointment links")
			}

			for i, in := range params.Lines {
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
					LineType:    nil,
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

				snapshotLines = append(snapshotLines, toInvoiceLine(updatedLine))
				netTotalCents += amts.netCents
				vatTotalCents += amts.vatCents
				grossTotalCents += amts.grossCents
			}
		} else {
			if err := qtx.DeleteInvoiceLinesByInvoice(ctx, invoiceID); err != nil {
				return nil, fmt.Errorf("failed to delete invoice lines: %w", err)
			}

			lineNo := int32(1)
			for _, in := range params.Lines {
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

				snapshotLines = append(snapshotLines, toInvoiceLine(line))
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

	issueDate := pgtype.Date{Time: params.IssueDate, Valid: !params.IssueDate.IsZero()}
	dueDate := pgtype.Date{Time: params.DueDate, Valid: !params.DueDate.IsZero()}

	updated, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:               invoiceID,
		IssueDate:        issueDate,
		DueDate:          dueDate,
		PeriodStart:      pgtype.Timestamptz{Valid: false},
		PeriodEnd:        pgtype.Timestamptz{Valid: false},
		BillingCycle:     nil,
		BillingTimezone:  nil,
		Source:           nil,
		BillToSnapshot:   nil,
		ClientSnapshot:   nil,
		DetailsSnapshot:  detailsSnapshot,
		NetTotalAmount:   netTotalAmount,
		VatTotalAmount:   vatTotalAmount,
		GrossTotalAmount: grossTotalAmount,
		Currency:         nil,
		ExtraContent:     params.ExtraContent,
		Status:           nil,
		WarningCount:     nil,
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

	d := invoiceFromRow(updated)
	return &d, nil
}

// ==================== Delete Invoice ====================

func (s *InvoiceService) DeleteInvoice(ctx context.Context, invoiceID uuid.UUID) error {
	return s.store.DeleteInvoice(ctx, invoiceID)
}

// ==================== Generate Invoice (Auto) ====================

func (s *InvoiceService) GenerateInvoice(ctx context.Context, params domain.GenerateInvoiceParams) (*domain.GenerateInvoiceResult, int64, error) {
	if params.ClientID == uuid.Nil {
		return nil, 0, fmt.Errorf("client_id is required")
	}
	if params.StartDate.IsZero() || params.EndDate.IsZero() {
		return nil, 0, fmt.Errorf("start_date and end_date are required")
	}
	if !params.EndDate.After(params.StartDate) {
		return nil, 0, fmt.Errorf("end_date must be after start_date")
	}

	billingTz := params.BillingTimezone
	if billingTz == "" {
		billingTz = defaultBillingTimezone
	}
	billingCycle := params.BillingCycle
	if billingCycle == "" {
		billingCycle = defaultBillingCycle
	}

	senderIDs, err := s.store.ListClientSendersForPeriod(ctx, db.ListClientSendersForPeriodParams{
		ClientID:    params.ClientID,
		PeriodStart: pgtype.Timestamptz{Time: params.StartDate, Valid: true},
		PeriodEnd:   pgtype.Timestamptz{Time: params.EndDate, Valid: true},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list client senders for period: %w", err)
	}
	if len(senderIDs) == 0 {
		return nil, 0, fmt.Errorf("no approved contracts found for client in the specified period")
	}
	if len(senderIDs) > 1 {
		return nil, 0, fmt.Errorf("multiple senders found for client in the specified period; generate per sender is required")
	}

	return s.generateInvoiceForTarget(ctx, generateTargetParams{
		ClientID:        params.ClientID,
		SenderID:        senderIDs[0],
		PeriodStart:     params.StartDate,
		PeriodEnd:       params.EndDate,
		BillingTimezone: billingTz,
		BillingCycle:    billingCycle,
		Source:          db.InvoiceSourceEnumAuto,
		RunID:           nil,
	})
}

type generateTargetParams struct {
	ClientID        uuid.UUID
	SenderID        uuid.UUID
	PeriodStart     time.Time
	PeriodEnd       time.Time
	BillingTimezone string
	BillingCycle    string
	Source          db.InvoiceSourceEnum
	RunID           *uuid.UUID
}

func (s *InvoiceService) generateInvoiceForTarget(ctx context.Context, p generateTargetParams) (*domain.GenerateInvoiceResult, int64, error) {
	var warningCount int64
	warnings := []string{}

	contracts, err := s.store.ListApprovedContractsForClientSenderInPeriod(ctx, db.ListApprovedContractsForClientSenderInPeriodParams{
		ClientID:    p.ClientID,
		SenderID:    p.SenderID,
		PeriodStart: pgtype.Timestamptz{Time: p.PeriodStart, Valid: true},
		PeriodEnd:   pgtype.Timestamptz{Time: p.PeriodEnd, Valid: true},
	})
	if err != nil {
		return nil, warningCount, fmt.Errorf("failed to list approved contracts: %w", err)
	}
	if len(contracts) == 0 {
		return nil, warningCount, fmt.Errorf("no approved contracts found for client/sender in period")
	}

	byCareType := map[db.CareTypeEnum]db.ListApprovedContractsForClientSenderInPeriodRow{}
	for _, c := range contracts {
		if _, ok := byCareType[c.CareType]; ok {
			return nil, warningCount, fmt.Errorf("multiple approved contracts found for care_type=%s; generation aborted", c.CareType)
		}
		byCareType[c.CareType] = c
	}

	invoiceNumber, invoiceSequence, err := s.generateInvoiceNumber(ctx)
	if err != nil {
		return nil, warningCount, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	issueDate := time.Now()
	dueDate := issueDate.Add(30 * 24 * time.Hour)

	tx, err := s.store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, warningCount, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.store.WithTx(tx)

	inv, err := qtx.CreateInvoice(ctx, db.CreateInvoiceParams{
		InvoiceNumber:     invoiceNumber,
		InvoiceSequence:   invoiceSequence,
		DueDate:           pgtype.Date{Time: dueDate, Valid: true},
		IssueDate:         pgtype.Date{Time: issueDate, Valid: true},
		Status:            db.InvoiceStatusEnumConcept,
		InvoiceType:       db.InvoiceTypeEnumStandard,
		Source:            p.Source,
		OriginalInvoiceID: nil,
		ReplacesInvoiceID: nil,
		PeriodStart:       pgtype.Timestamptz{Time: p.PeriodStart, Valid: true},
		PeriodEnd:         pgtype.Timestamptz{Time: p.PeriodEnd, Valid: true},
		BillingCycle:      &p.BillingCycle,
		BillingTimezone:   p.BillingTimezone,
		Currency:          "EUR",
		ExtraContent:      []byte("{}"),
		ClientID:          p.ClientID,
		SenderID:          p.SenderID,
		WarningCount:      0,
		RunID:             p.RunID,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && p.Source == db.InvoiceSourceEnumAuto {
			return nil, warningCount, fmt.Errorf("%w: %v", ErrAutoInvoiceAlreadyExists, err)
		}
		return nil, warningCount, fmt.Errorf("failed to create invoice: %w", err)
	}

	lineNo := int32(1)
	var netTotalCents, vatTotalCents, grossTotalCents int64

	for _, contract := range byCareType {
		segments, err := qtx.GetBillablePeriodsForContract(ctx, db.GetBillablePeriodsForContractParams{
			InvoiceStartDate: pgtype.Timestamptz{Time: p.PeriodStart, Valid: true},
			InvoiceEndDate:   pgtype.Timestamptz{Time: p.PeriodEnd, Valid: true},
			ContractID:       contract.ID,
		})
		if err != nil {
			warningCount++
			warnings = append(warnings, fmt.Sprintf("failed to compute approved periods for contract %s", contract.ID.String()))
			continue
		}

		for _, seg := range segments {
			segStart := clampTime(seg.BillableStart.Time, contract.StartDate.Time, contract.EndDate.Time)
			segStart = clampTime(segStart, p.PeriodStart, p.PeriodEnd)
			segEnd := clampTime(seg.BillableEnd.Time, contract.StartDate.Time, contract.EndDate.Time)
			segEnd = clampTime(segEnd, p.PeriodStart, p.PeriodEnd)
			if !segEnd.After(segStart) {
				continue
			}

			vatRatePct := int64(0)
			if contract.Vat != nil {
				vatRatePct = int64(*contract.Vat)
			}

			switch contract.CareType {
			case db.CareTypeEnumAccommodation:
				days, err := countCalendarDays(segStart, segEnd, p.BillingTimezone)
				if err != nil {
					warningCount++
					warnings = append(warnings, fmt.Sprintf("failed to count accommodation days for contract %s: %v", contract.ID.String(), err))
					continue
				}
				if days <= 0 {
					continue
				}

				priceCents := centsFromAmount(contract.Price)
				var unitPrice float64
				var netCents int64

				if contract.PriceTimeUnit == db.PriceTimeUnitEnumDaily {
					unitPrice = roundTo(contract.Price, 4)
					netCents = priceCents * days
				} else if contract.PriceTimeUnit == db.PriceTimeUnitEnumWeekly {
					unitPrice = roundTo(contract.Price/7, 4)
					netCents = roundHalfUpDiv(priceCents*days, 7)
				} else {
					warningCount++
					warnings = append(warnings, fmt.Sprintf("unsupported accommodation price_time_unit=%s for contract %s", contract.PriceTimeUnit, contract.ID.String()))
					continue
				}

				amts := computeVat(netCents, vatRatePct)
				_, err = qtx.CreateInvoiceLine(ctx, db.CreateInvoiceLineParams{
					InvoiceID:   inv.ID,
					ClientID:    inv.ClientID,
					SenderID:    inv.SenderID,
					LineNo:      lineNo,
					LineType:    db.InvoiceLineTypeEnumContract,
					ContractID:  &contract.ID,
					ServiceType: string(contract.CareType),
					Description: contract.CareName,
					PeriodStart: pgtype.Timestamptz{Time: segStart, Valid: true},
					PeriodEnd:   pgtype.Timestamptz{Time: segEnd, Valid: true},
					Quantity:    float64(days),
					Unit:        "day",
					UnitPrice:   unitPrice,
					NetAmount:   amountFromCents(amts.netCents),
					VatRate:     float64(vatRatePct),
					VatAmount:   amountFromCents(amts.vatCents),
					GrossAmount: amountFromCents(amts.grossCents),
					Metadata:    []byte("{}"),
				})
				if err != nil {
					return nil, warningCount, fmt.Errorf("failed to create invoice_line: %w", err)
				}

				lineNo++
				netTotalCents += amts.netCents
				vatTotalCents += amts.vatCents
				grossTotalCents += amts.grossCents

			case db.CareTypeEnumAmbulante:
				appts, err := qtx.ListUnbilledClientAppointmentsStartingInRange(ctx, db.ListUnbilledClientAppointmentsStartingInRangeParams{
					ClientID:  &inv.ClientID,
					StartDate: pgtype.Timestamptz{Time: segStart, Valid: true},
					EndDate:   pgtype.Timestamptz{Time: segEnd, Valid: true},
				})
				if err != nil {
					warningCount++
					warnings = append(warnings, fmt.Sprintf("failed to list appointments for contract %s", contract.ID.String()))
					continue
				}
				if len(appts) == 0 {
					continue
				}

				var totalSeconds int64
				for _, a := range appts {
					d := a.EndTime.Time.Sub(a.StartTime.Time)
					if d > 0 {
						totalSeconds += int64(d.Seconds())
					}
				}
				if totalSeconds <= 0 {
					continue
				}

				priceCents := centsFromAmount(contract.Price)
				var quantity float64
				var unit string
				var unitPrice float64
				var netCents int64

				if contract.PriceTimeUnit == db.PriceTimeUnitEnumMinute {
					unit = "minute"
					quantity = roundTo(float64(totalSeconds)/60, 4)
					unitPrice = roundTo(contract.Price, 4)
					netCents = roundHalfUpDiv(priceCents*totalSeconds, 60)
				} else if contract.PriceTimeUnit == db.PriceTimeUnitEnumHourly {
					unit = "hour"
					quantity = roundTo(float64(totalSeconds)/3600, 4)
					unitPrice = roundTo(contract.Price, 4)
					netCents = roundHalfUpDiv(priceCents*totalSeconds, 3600)
				} else {
					warningCount++
					warnings = append(warnings, fmt.Sprintf("unsupported ambulante price_time_unit=%s for contract %s", contract.PriceTimeUnit, contract.ID.String()))
					continue
				}

				amts := computeVat(netCents, vatRatePct)
				invoiceLine, err := qtx.CreateInvoiceLine(ctx, db.CreateInvoiceLineParams{
					InvoiceID:   inv.ID,
					ClientID:    inv.ClientID,
					SenderID:    inv.SenderID,
					LineNo:      lineNo,
					LineType:    db.InvoiceLineTypeEnumContract,
					ContractID:  &contract.ID,
					ServiceType: string(contract.CareType),
					Description: contract.CareName,
					PeriodStart: pgtype.Timestamptz{Time: segStart, Valid: true},
					PeriodEnd:   pgtype.Timestamptz{Time: segEnd, Valid: true},
					Quantity:    quantity,
					Unit:        unit,
					UnitPrice:   unitPrice,
					NetAmount:   amountFromCents(amts.netCents),
					VatRate:     float64(vatRatePct),
					VatAmount:   amountFromCents(amts.vatCents),
					GrossAmount: amountFromCents(amts.grossCents),
					Metadata:    []byte("{}"),
				})
				if err != nil {
					return nil, warningCount, fmt.Errorf("failed to create invoice_line: %w", err)
				}

				for _, a := range appts {
					dur := a.EndTime.Time.Sub(a.StartTime.Time)
					if dur <= 0 {
						continue
					}
					minutes := dur.Minutes()
					_, err = qtx.CreateInvoiceLineCalendarEvent(ctx, db.CreateInvoiceLineCalendarEventParams{
						InvoiceLineID:   invoiceLine.ID,
						CalendarEventID: a.AppointmentID,
						ClientID:        inv.ClientID,
						StartAt:         a.StartTime,
						EndAt:           a.EndTime,
						MinutesBilled:   minutes,
						Metadata:        []byte("{}"),
					})
					if err != nil {
						return nil, warningCount, fmt.Errorf("failed to create invoice_line_calendar_event: %w", err)
					}

					_, err = qtx.InsertBilledCalendarEvent(ctx, db.InsertBilledCalendarEventParams{
						CalendarEventID: a.AppointmentID,
						ClientID:        inv.ClientID,
						InvoiceID:       inv.ID,
						InvoiceLineID:   invoiceLine.ID,
					})
					if err != nil {
						return nil, warningCount, fmt.Errorf("failed to insert billed_calendar_event: %w", err)
					}
				}

				lineNo++
				netTotalCents += amts.netCents
				vatTotalCents += amts.vatCents
				grossTotalCents += amts.grossCents
			}
		}
	}

	if lineNo == 1 {
		return nil, warningCount, fmt.Errorf("%w for client/sender in the specified period", ErrNoBillableItems)
	}

	lines, err := qtx.ListInvoiceLinesByInvoice(ctx, inv.ID)
	if err != nil {
		return nil, warningCount, fmt.Errorf("failed to list invoice lines: %w", err)
	}
	lineDTOs := toInvoiceLines(lines)
	snapshotBytes, _ := json.Marshal(lineDTOs)

	netTotal := amountFromCents(netTotalCents)
	vatTotal := amountFromCents(vatTotalCents)
	grossTotal := amountFromCents(grossTotalCents)

	wc := int32(warningCount)
	calcMeta, _ := json.Marshal(map[string]any{
		"warnings": warnings,
	})

	updated, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:               inv.ID,
		IssueDate:        pgtype.Date{Valid: false},
		DueDate:          pgtype.Date{Valid: false},
		PeriodStart:      pgtype.Timestamptz{Valid: false},
		PeriodEnd:        pgtype.Timestamptz{Valid: false},
		BillingCycle:     nil,
		BillingTimezone:  nil,
		Source:           nil,
		BillToSnapshot:   nil,
		ClientSnapshot:   nil,
		DetailsSnapshot:  snapshotBytes,
		NetTotalAmount:   &netTotal,
		VatTotalAmount:   &vatTotal,
		GrossTotalAmount: &grossTotal,
		Currency:         nil,
		ExtraContent:     nil,
		Status:           nil,
		WarningCount:     &wc,
		RunID:            nil,
		LockedAt:         pgtype.Timestamptz{Valid: false},
		CalcVersion:      nil,
		CalcMetadata:     calcMeta,
	})
	if err != nil {
		return nil, warningCount, fmt.Errorf("failed to update invoice totals: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, warningCount, fmt.Errorf("failed to commit invoice transaction: %w", err)
	}

	var ps, pe *time.Time
	if updated.PeriodStart.Valid {
		t := updated.PeriodStart.Time
		ps = &t
	}
	if updated.PeriodEnd.Valid {
		t := updated.PeriodEnd.Time
		pe = &t
	}

	result := &domain.GenerateInvoiceResult{
		Invoice: domain.Invoice{
			ID:              updated.ID,
			InvoiceNumber:   updated.InvoiceNumber,
			IssueDate:       updated.IssueDate.Time,
			DueDate:         updated.DueDate.Time,
			Status:          string(updated.Status),
			Source:          string(updated.Source),
			InvoiceType:     string(updated.InvoiceType),
			PeriodStart:     ps,
			PeriodEnd:       pe,
			Currency:        updated.Currency,
			NetTotal:        updated.NetTotalAmount,
			VatTotal:        updated.VatTotalAmount,
			GrossTotal:      updated.GrossTotalAmount,
			PdfAttachmentID: updated.PdfAttachmentID,
			ExtraContent:    updated.ExtraContent,
			ClientID:        updated.ClientID,
			SenderID:        updated.SenderID,
			UpdatedAt:       updated.UpdatedAt.Time,
			CreatedAt:       updated.CreatedAt.Time,
		},
		Lines:    lineDTOs,
		Warnings: warnings,
	}
	return result, warningCount, nil
}

// ==================== Credit Invoice ====================

func (s *InvoiceService) CreditInvoice(ctx context.Context, invoiceID uuid.UUID, employeeID uuid.UUID) (*domain.CreditInvoiceResult, error) {
	tx, err := s.store.ConnPool.Begin(ctx)
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
	qtx := s.store.WithTx(tx)

	original, err := qtx.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.CreditInvoice", "failed to get original invoice", err, zap.String("invoice_id", invoiceID.String()))
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

	creditNumber, creditSeq, err := s.generateInvoiceNumber(ctx)
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
	cl := int32(1)

	for _, l := range origLines {
		netCents := -centsFromAmount(l.NetAmount)
		vatC := -centsFromAmount(l.VatAmount)
		grossC := -centsFromAmount(l.GrossAmount)

		_, err = qtx.CreateInvoiceLine(ctx, db.CreateInvoiceLineParams{
			InvoiceID:   creditInv.ID,
			ClientID:    creditInv.ClientID,
			SenderID:    creditInv.SenderID,
			LineNo:      cl,
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
			VatAmount:   amountFromCents(vatC),
			GrossAmount: amountFromCents(grossC),
			Metadata:    l.Metadata,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create credit invoice line: %w", err)
		}

		cl++
		netTotalCents += netCents
		vatTotalCents += vatC
		grossTotalCents += grossC
	}

	netTotal := amountFromCents(netTotalCents)
	vatTotal := amountFromCents(vatTotalCents)
	grossTotal := amountFromCents(grossTotalCents)

	if err := qtx.VoidBilledCalendarEventsByInvoice(ctx, original.ID); err != nil {
		return nil, fmt.Errorf("failed to void billed appointments: %w", err)
	}

	if _, err := qtx.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{ID: original.ID, Status: db.InvoiceStatusEnumCanceled}); err != nil {
		return nil, fmt.Errorf("failed to cancel original invoice: %w", err)
	}

	_, err = qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:                creditInv.ID,
		IssueDate:         pgtype.Date{Valid: false},
		DueDate:           pgtype.Date{Valid: false},
		PeriodStart:       pgtype.Timestamptz{Valid: false},
		PeriodEnd:         pgtype.Timestamptz{Valid: false},
		BillingCycle:      nil,
		BillingTimezone:   nil,
		Source:            nil,
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
		Status:            nil,
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

	return &domain.CreditInvoiceResult{ID: creditInv.ID}, nil
}

// ==================== Get Invoice Audit Logs ====================

func (s *InvoiceService) GetInvoiceAuditLogs(ctx context.Context, invoiceID uuid.UUID) ([]domain.InvoiceAuditLog, error) {
	logs, err := s.store.GetInvoiceAuditLogs(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice audit logs: %w", err)
	}

	out := make([]domain.InvoiceAuditLog, 0, len(logs))
	for _, l := range logs {
		out = append(out, domain.InvoiceAuditLog{
			AuditID:            l.AuditID,
			InvoiceID:          l.InvoiceID,
			Operation:          string(l.Operation),
			ChangedBy:          l.ChangedBy,
			ChangedAt:          l.ChangedAt.Time,
			OldValues:          l.OldValues,
			NewValues:          l.NewValues,
			ChangedFields:      l.ChangedFields,
			ChangedByFirstName: l.ChangedByFirstName,
			ChangedByLastName:  l.ChangedByLastName,
		})
	}
	return out, nil
}

// ==================== Invoice Template Items ====================

func (s *InvoiceService) GetInvoiceTemplateItems(ctx context.Context) ([]domain.InvoiceTemplateItemData, error) {
	items, err := s.store.GetAllTemplateItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get template items: %w", err)
	}

	out := make([]domain.InvoiceTemplateItemData, 0, len(items))
	for _, it := range items {
		out = append(out, domain.InvoiceTemplateItemData{
			ID:           it.ID,
			ItemTag:      it.ItemTag,
			Description:  it.Description,
			SourceTable:  it.SourceTable,
			SourceColumn: it.SourceColumn,
		})
	}
	return out, nil
}

// ==================== Generate Invoice PDF ====================

func (s *InvoiceService) GenerateInvoicePDF(ctx context.Context, invoiceID uuid.UUID) (*domain.GeneratePDFResult, error) {
	inv, err := s.store.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	if inv.PdfAttachmentID != nil {
		att, err := s.store.GetAttachmentById(ctx, *inv.PdfAttachmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get pdf attachment: %w", err)
		}
		url, err := s.storage.GeneratePresignedURL(ctx, att.File, 15*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned url: %w", err)
		}
		return &domain.GeneratePDFResult{FileURL: url}, nil
	}

	lines, err := s.store.ListInvoiceLinesByInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoice lines: %w", err)
	}

	details := make([]domain.InvoiceDetailPDF, 0, len(lines))
	for _, l := range lines {
		details = append(details, domain.InvoiceDetailPDF{
			CareType:      l.ServiceType,
			Periods:       []domain.InvoicePeriodPDF{{StartDate: l.PeriodStart.Time, EndDate: l.PeriodEnd.Time}},
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

	pdfData := domain.InvoicePDF{
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

	objKey, size, err := s.pdfService.GenerateAndUploadInvoicePDF(ctx, pdfData)
	if err != nil {
		return nil, err
	}

	immutable := inv.LockedAt.Valid || inv.Status != db.InvoiceStatusEnumConcept
	if immutable {
		if size > math.MaxInt32 {
			return nil, fmt.Errorf("pdf too large to store in attachment_file: %d bytes", size)
		}

		attID := uuid.New()
		name := fmt.Sprintf("invoice_%s.pdf", inv.InvoiceNumber)
		tag := "invoice_pdf"

		tx, err := s.store.ConnPool.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)
		qtx := s.store.WithTx(tx)

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
			if err == pgx.ErrNoRows {
				latest, gErr := qtx.GetInvoice(ctx, invoiceID)
				if gErr == nil && latest.PdfAttachmentID != nil {
					att, aErr := qtx.GetAttachmentById(ctx, *latest.PdfAttachmentID)
					if aErr == nil {
						_ = tx.Commit(ctx)
						url, uErr := s.storage.GeneratePresignedURL(ctx, att.File, 15*time.Minute)
						if uErr == nil {
							return &domain.GeneratePDFResult{FileURL: url}, nil
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

	url, err := s.storage.GeneratePresignedURL(ctx, objKey, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned url: %w", err)
	}
	return &domain.GeneratePDFResult{FileURL: url}, nil
}

// ==================== Send Invoice Reminder ====================

func (s *InvoiceService) SendInvoiceReminder(ctx context.Context, invoiceID uuid.UUID) error {
	senderID, err := s.store.GetInvoiceSenderID(ctx, invoiceID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.SendInvoiceReminder", "failed to get sender ID for invoice", err, zap.String("invoice_id", invoiceID.String()))
		return err
	}
	s.logger.LogInfo(ctx, "InvoiceService.SendInvoiceReminder", "sending reminder email",
		zap.String("invoice_id", invoiceID.String()), zap.String("sender_id", senderID.String()))
	return nil
}

// ==================== Payments ====================

func (s *InvoiceService) determineInvoiceStatus(invoiceTotal, totalPaid float64) (db.InvoiceStatusEnum, error) {
	diff := totalPaid - invoiceTotal

	if totalPaid <= paymentTolerance {
		return db.InvoiceStatusEnumOutstanding, nil
	}

	if diff < -paymentTolerance {
		return db.InvoiceStatusEnumPartiallyPaid, nil
	}

	if diff >= -paymentTolerance && diff <= paymentTolerance {
		return db.InvoiceStatusEnumPaid, nil
	}
	if diff > paymentTolerance {
		return db.InvoiceStatusEnumOverpaid, nil
	}

	return "", fmt.Errorf("could not determine invoice status for totalPaid: %f, invoiceTotal: %f", totalPaid, invoiceTotal)
}

func (s *InvoiceService) calculatePaymentCompletionPercentage(ctx context.Context, totalAmount float64, invoiceID uuid.UUID) float64 {
	if totalAmount == 0 {
		return 0
	}
	totalPaid, err := s.store.GetCompletedPaymentSum(ctx, invoiceID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.calculatePaymentCompletionPercentage", "failed to get total completed payment", err, zap.String("invoice_id", invoiceID.String()))
		return 0
	}
	return (totalPaid / totalAmount) * 100
}

func (s *InvoiceService) CreatePayment(ctx context.Context, invoiceID uuid.UUID, employeeID uuid.UUID, params domain.CreatePaymentParams) (*domain.CreatePaymentResult, error) {
	tx, err := s.store.ConnPool.Begin(ctx)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.CreatePayment", "failed to begin transaction", err, zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "SELECT set_config('myapp.current_employee_id', $1, true)", employeeID.String())
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.CreatePayment", "failed to set current employee ID", err, zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}
	qtx := s.store.WithTx(tx)

	getInvoice, err := qtx.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.CreatePayment", "failed to get invoice", err, zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to get invoice: %v", err)
	}

	paymentParams := db.CreatePaymentParams{
		InvoiceID:        invoiceID,
		PaymentMethod:    db.PaymentMethodEnum(params.PaymentMethod),
		PaymentStatus:    db.PaymentStatusEnum(params.PaymentStatus),
		Amount:           params.Amount,
		PaymentDate:      pgtype.Date{Time: params.PaymentDate, Valid: true},
		PaymentReference: params.PaymentReference,
		Notes:            params.Notes,
		RecordedBy:       &employeeID,
	}

	payment, err := qtx.CreatePayment(ctx, paymentParams)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.CreatePayment", "failed to create payment", err, zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	var newInvoiceStatus db.InvoiceStatusEnum
	invoiceStatusChanged := false

	if params.PaymentStatus == string(domain.PaymentStatusCompleted) {
		totalPaid, err := qtx.GetCompletedPaymentSum(ctx, invoiceID)
		if err != nil {
			s.logger.LogError(ctx, "InvoiceService.CreatePayment", "failed to get total completed payment", err, zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to get total completed payment: %v", err)
		}

		newInvoiceStatus, err = s.determineInvoiceStatus(getInvoice.GrossTotalAmount, totalPaid)
		if err != nil {
			s.logger.LogError(ctx, "InvoiceService.CreatePayment", "failed to determine invoice status", err, zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to determine invoice status: %v", err)
		}

		if newInvoiceStatus != db.InvoiceStatusEnum(getInvoice.Status) {
			invoiceStatusChanged = true
			_, err = qtx.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{
				ID:     invoiceID,
				Status: newInvoiceStatus,
			})
			if err != nil {
				s.logger.LogError(ctx, "InvoiceService.CreatePayment", "failed to update invoice status", err, zap.String("invoice_id", invoiceID.String()))
				return nil, fmt.Errorf("failed to update invoice status: %v", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.CreatePayment", "failed to commit transaction", err, zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &domain.CreatePaymentResult{
		PaymentID:            payment.ID,
		InvoiceID:            payment.InvoiceID,
		PaymentMethod:        string(payment.PaymentMethod),
		PaymentStatus:        string(payment.PaymentStatus),
		Amount:               payment.Amount,
		PaymentDate:          payment.PaymentDate.Time,
		PaymentReference:     payment.PaymentReference,
		Notes:                payment.Notes,
		RecordedBy:           payment.RecordedBy,
		InvoiceStatusChanged: invoiceStatusChanged,
		CurrentInvoiceStatus: string(newInvoiceStatus),
	}, nil
}

func (s *InvoiceService) ListPayments(ctx context.Context, invoiceID uuid.UUID) ([]domain.Payment, error) {
	payments, err := s.store.ListPayments(ctx, invoiceID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.ListPayments", "failed to list payments", err, zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to list payments: %v", err)
	}

	response := make([]domain.Payment, 0, len(payments))
	for _, payment := range payments {
		response = append(response, domain.Payment{
			ID:                  payment.ID,
			InvoiceID:           payment.InvoiceID,
			PaymentMethod:       string(payment.PaymentMethod),
			PaymentStatus:       string(payment.PaymentStatus),
			Amount:              payment.Amount,
			PaymentDate:         payment.PaymentDate.Time,
			PaymentReference:    payment.PaymentReference,
			Notes:               payment.Notes,
			RecordedBy:          payment.RecordedBy,
			CreatedAt:           payment.CreatedAt.Time,
			UpdatedAt:           payment.UpdatedAt.Time,
			RecordedByFirstName: payment.RecordedByFirstName,
			RecordedByLastName:  payment.RecordedByLastName,
		})
	}
	return response, nil
}

func (s *InvoiceService) GetPaymentByID(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
	payment, err := s.store.GetPayment(ctx, paymentID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.GetPaymentByID", "failed to get payment by ID", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to get payment by ID: %v", err)
	}

	return &domain.Payment{
		ID:                  payment.ID,
		InvoiceID:           payment.InvoiceID,
		PaymentMethod:       string(payment.PaymentMethod),
		PaymentStatus:       string(payment.PaymentStatus),
		Amount:              payment.Amount,
		PaymentDate:         payment.PaymentDate.Time,
		PaymentReference:    payment.PaymentReference,
		Notes:               payment.Notes,
		RecordedBy:          payment.RecordedBy,
		CreatedAt:           payment.CreatedAt.Time,
		UpdatedAt:           payment.UpdatedAt.Time,
		RecordedByFirstName: payment.RecordedByFirstName,
		RecordedByLastName:  payment.RecordedByLastName,
	}, nil
}

func (s *InvoiceService) UpdatePayment(ctx context.Context, invoiceID uuid.UUID, paymentID uuid.UUID, employeeID uuid.UUID, params domain.UpdatePaymentParams) (*domain.UpdatePaymentResult, error) {
	tx, err := s.store.ConnPool.Begin(ctx)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.UpdatePayment", "failed to begin transaction", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "SELECT set_config('myapp.current_employee_id', $1, true)", employeeID.String())
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.UpdatePayment", "failed to set current employee ID", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}
	qtx := s.store.WithTx(tx)

	currentPayment, err := qtx.GetPaymentWithInvoice(ctx, paymentID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.UpdatePayment", "failed to get current payment", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to get current payment: %v", err)
	}

	if currentPayment.InvoiceID != invoiceID {
		return nil, fmt.Errorf("payment does not belong to the specified invoice")
	}

	originalInvoiceStatus := currentPayment.InvoiceStatus
	updateParams := db.UpdatePaymentParams{
		ID:               paymentID,
		PaymentMethod:    db.NullPaymentMethodFromPtr(params.PaymentMethod),
		PaymentStatus:    db.NullPaymentStatusFromPtr(params.PaymentStatus),
		Amount:           params.Amount,
		PaymentReference: params.PaymentReference,
		Notes:            params.Notes,
		RecordedBy:       &employeeID,
	}
	if params.PaymentDate != nil {
		updateParams.PaymentDate = pgtype.Date{Time: *params.PaymentDate, Valid: true}
	}

	updatedPayment, err := qtx.UpdatePayment(ctx, updateParams)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.UpdatePayment", "failed to update payment", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to update payment: %v", err)
	}

	var newInvoiceStatus db.InvoiceStatusEnum
	var statusChanged bool = false

	if updatedPayment.PaymentStatus == db.PaymentStatusEnumCompleted ||
		currentPayment.PaymentStatus == db.PaymentStatusEnumCompleted {
		totalPaid, err := qtx.GetTotalPaidAmountByInvoice(ctx, invoiceID)
		if err != nil {
			s.logger.LogError(ctx, "InvoiceService.UpdatePayment", "failed to get total paid amount", err, zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to get total paid amount: %v", err)
		}

		newStatus, err := s.determineInvoiceStatus(currentPayment.InvoiceTotalAmount, totalPaid)
		if err != nil {
			s.logger.LogError(ctx, "InvoiceService.UpdatePayment", "failed to determine invoice status", err, zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to determine invoice status: %v", err)
		}

		if newStatus != originalInvoiceStatus {
			updatedInvoice, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
				ID:     invoiceID,
				Status: &newStatus,
			})
			if err != nil {
				s.logger.LogError(ctx, "InvoiceService.UpdatePayment", "failed to update invoice status", err, zap.String("invoice_id", invoiceID.String()))
				return nil, fmt.Errorf("failed to update invoice status: %v", err)
			}
			newInvoiceStatus = updatedInvoice.Status
			statusChanged = true
		} else {
			newInvoiceStatus = originalInvoiceStatus
		}
	} else {
		newInvoiceStatus = originalInvoiceStatus
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.UpdatePayment", "failed to commit transaction", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &domain.UpdatePaymentResult{
		PaymentID:             updatedPayment.ID,
		InvoiceID:             updatedPayment.InvoiceID,
		PaymentMethod:         string(updatedPayment.PaymentMethod),
		PaymentStatus:         string(updatedPayment.PaymentStatus),
		Amount:                updatedPayment.Amount,
		PaymentDate:           updatedPayment.PaymentDate.Time,
		PaymentReference:      updatedPayment.PaymentReference,
		Notes:                 updatedPayment.Notes,
		RecordedBy:            updatedPayment.RecordedBy,
		InvoiceStatusChanged:  statusChanged,
		CurrentInvoiceStatus:  string(newInvoiceStatus),
		PreviousInvoiceStatus: string(originalInvoiceStatus),
	}, nil
}

func (s *InvoiceService) DeletePayment(ctx context.Context, invoiceID uuid.UUID, paymentID uuid.UUID, employeeID uuid.UUID) (*domain.DeletePaymentResult, error) {
	tx, err := s.store.ConnPool.Begin(ctx)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.DeletePayment", "failed to begin transaction", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "SELECT set_config('myapp.current_employee_id', $1, true)", employeeID.String())
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.DeletePayment", "failed to set current employee ID", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}
	qtx := s.store.WithTx(tx)

	paymentToDelete, err := qtx.GetPaymentWithInvoice(ctx, paymentID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.DeletePayment", "failed to get payment to delete", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to get payment to delete: %v", err)
	}

	if paymentToDelete.InvoiceID != invoiceID {
		return nil, fmt.Errorf("payment does not belong to the specified invoice")
	}

	originalInvoiceStatus := paymentToDelete.InvoiceStatus

	deletedPayment, err := qtx.DeletePayment(ctx, paymentID)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.DeletePayment", "failed to delete payment", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to delete payment: %v", err)
	}

	var newInvoiceStatus db.InvoiceStatusEnum
	var statusChanged bool = false

	if deletedPayment.PaymentStatus == db.PaymentStatusEnumCompleted {
		totalPaid, err := qtx.GetTotalPaidAmountByInvoice(ctx, invoiceID)
		if err != nil {
			s.logger.LogError(ctx, "InvoiceService.DeletePayment", "failed to get total paid amount", err, zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to get total paid amount: %v", err)
		}

		newStatus, err := s.determineInvoiceStatus(paymentToDelete.InvoiceTotalAmount, totalPaid)
		if err != nil {
			s.logger.LogError(ctx, "InvoiceService.DeletePayment", "failed to determine invoice status", err, zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to determine invoice status: %v", err)
		}

		if newStatus != originalInvoiceStatus {
			updatedInvoice, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
				ID:     invoiceID,
				Status: &newStatus,
			})
			if err != nil {
				s.logger.LogError(ctx, "InvoiceService.DeletePayment", "failed to update invoice status", err, zap.String("invoice_id", invoiceID.String()))
				return nil, fmt.Errorf("failed to update invoice status: %v", err)
			}
			newInvoiceStatus = updatedInvoice.Status
			statusChanged = true
		} else {
			newInvoiceStatus = originalInvoiceStatus
		}
	} else {
		newInvoiceStatus = originalInvoiceStatus
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.LogError(ctx, "InvoiceService.DeletePayment", "failed to commit transaction", err, zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &domain.DeletePaymentResult{
		DeletedPaymentID:      deletedPayment.ID,
		InvoiceID:             deletedPayment.InvoiceID,
		DeletedAmount:         deletedPayment.Amount,
		DeletedPaymentStatus:  string(deletedPayment.PaymentStatus),
		InvoiceStatusChanged:  statusChanged,
		CurrentInvoiceStatus:  string(newInvoiceStatus),
		PreviousInvoiceStatus: string(originalInvoiceStatus),
	}, nil
}

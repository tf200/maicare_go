package invoice

import (
	"context"
	"errors"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/util"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

var (
	ErrNoBillableItems          = errors.New("no billable items")
	ErrAutoInvoiceAlreadyExists = errors.New("auto invoice already exists for period")
)

func (s *invoiceService) GenerateInvoiceNumber(ctx context.Context) (string, int64, error) {
	now := time.Now()
	datePart := now.Format("20060102") // YYYYMMDD

	maxSeq, err := s.Store.GetMaxInvoiceSequenceForDate(ctx, now)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get max invoice sequence: %w", err)
	}

	nextSeq := maxSeq + 1
	return fmt.Sprintf("INV-%s-%04d", datePart, nextSeq), nextSeq, nil
}

func (s *invoiceService) GenerateInvoice(req GenerateInvoiceRequest, ctx context.Context) (*GenerateInvoiceResponse, int64, error) {
	if req.ClientID == uuid.Nil {
		return nil, 0, fmt.Errorf("client_id is required")
	}
	if req.StartDate.IsZero() || req.EndDate.IsZero() {
		return nil, 0, fmt.Errorf("start_date and end_date are required")
	}
	if !req.EndDate.After(req.StartDate) {
		return nil, 0, fmt.Errorf("end_date must be after start_date")
	}

	billingTz := req.BillingTimezone
	if billingTz == "" {
		billingTz = DefaultBillingTimezone
	}
	billingCycle := req.BillingCycle
	if billingCycle == "" {
		billingCycle = DefaultBillingCycle
	}

	senderIDs, err := s.Store.ListClientSendersForPeriod(ctx, db.ListClientSendersForPeriodParams{
		ClientID:    req.ClientID,
		PeriodStart: pgtype.Timestamptz{Time: req.StartDate, Valid: true},
		PeriodEnd:   pgtype.Timestamptz{Time: req.EndDate, Valid: true},
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

	result, warnings, err := s.generateInvoiceForTarget(ctx, generateTargetParams{
		ClientID:        req.ClientID,
		SenderID:        senderIDs[0],
		PeriodStart:     req.StartDate,
		PeriodEnd:       req.EndDate,
		BillingTimezone: billingTz,
		BillingCycle:    billingCycle,
		Source:          db.InvoiceSourceEnumAuto,
		RunID:           nil,
	})
	return result, warnings, err
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

func (s *invoiceService) generateInvoiceForTarget(ctx context.Context, p generateTargetParams) (*GenerateInvoiceResponse, int64, error) {
	var warningCount int64
	warnings := []string{}

	contracts, err := s.Store.ListApprovedContractsForClientSenderInPeriod(ctx, db.ListApprovedContractsForClientSenderInPeriodParams{
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

	// Enforce: at most 1 contract per care_type.
	byCareType := map[db.CareTypeEnum]db.ListApprovedContractsForClientSenderInPeriodRow{}
	for _, c := range contracts {
		if _, ok := byCareType[c.CareType]; ok {
			return nil, warningCount, fmt.Errorf("multiple approved contracts found for care_type=%s; generation aborted", c.CareType)
		}
		byCareType[c.CareType] = c
	}

	invoiceNumber, invoiceSequence, err := s.GenerateInvoiceNumber(ctx)
	if err != nil {
		return nil, warningCount, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	issueDate := time.Now()
	dueDate := issueDate.Add(30 * 24 * time.Hour)

	// Create in a single transaction to keep lines + billed locks consistent.
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, warningCount, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.Store.WithTx(tx)

	// Create the invoice header with minimal fields; then UpdateInvoice to set the generation metadata/totals.
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
		// Auto invoices are idempotent via a partial unique index; treat duplicates as a "skip".
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
				line, err := qtx.CreateInvoiceLine(ctx, db.CreateInvoiceLineParams{
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

				_ = line // reserved for future source linkage (e.g., accommodation stay sources)
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

				// Link billed appointments and lock them to prevent double-billing.
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

	// Store a simple snapshot for now (backward-friendly), derived from the canonical lines.
	lines, err := qtx.ListInvoiceLinesByInvoice(ctx, inv.ID)
	if err != nil {
		return nil, warningCount, fmt.Errorf("failed to list invoice lines: %w", err)
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
		Source:           db.NullInvoiceSourceEnum{Valid: false},
		BillToSnapshot:   nil,
		ClientSnapshot:   nil,
		DetailsSnapshot:  snapshotBytes,
		NetTotalAmount:   &netTotal,
		VatTotalAmount:   &vatTotal,
		GrossTotalAmount: &grossTotal,
		Currency:         nil,
		ExtraContent:     nil,
		Status:           db.NullInvoiceStatusEnum{Valid: false},
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

	return &GenerateInvoiceResponse{
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
		ExtraContent:    util.ParseJSONToObject(updated.ExtraContent),
		ClientID:        updated.ClientID,
		SenderID:        updated.SenderID,
		Lines:           lineDTOs,
		Warnings:        warnings,
		UpdatedAt:       updated.UpdatedAt.Time,
		CreatedAt:       updated.CreatedAt.Time,
	}, warningCount, nil
}

func (s *invoiceService) BatchGenerateInvoices(ctx context.Context) error {
	// Keep existing scheduling logic, but make period boundaries stable.
	now := time.Now().UTC()
	isoWeek, isoYear := now.ISOWeek()
	if isoWeek%4 != 1 || now.Weekday() != time.Monday {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "BatchGenerateInvoices", "Not the scheduled time for batch invoice generation",
			zap.Int("current_week", isoWeek), zap.Int("current_year", isoYear))
		return nil
	}

	periodEnd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	periodStart := periodEnd.AddDate(0, 0, -28)

	run, err := s.Store.CreateInvoiceRun(ctx, db.CreateInvoiceRunParams{
		BillingCycle: DefaultBillingCycle,
		PeriodStart:  pgtype.Timestamptz{Time: periodStart, Valid: true},
		PeriodEnd:    pgtype.Timestamptz{Time: periodEnd, Valid: true},
		Timezone:     DefaultBillingTimezone,
		DryRun:       false,
		Status:       db.InvoiceRunStatusEnumRunning,
		Params:       []byte("{}"),
		CreatedBy:    nil,
	})
	if err != nil {
		return fmt.Errorf("failed to create invoice_run: %w", err)
	}

	targets, err := s.Store.ListInvoiceTargetsForPeriod(ctx, db.ListInvoiceTargetsForPeriodParams{
		PeriodStart: pgtype.Timestamptz{Time: periodStart, Valid: true},
		PeriodEnd:   pgtype.Timestamptz{Time: periodEnd, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to list invoice targets: %w", err)
	}

	var hadErrors bool
	for _, t := range targets {
		warningsBytes, _ := json.Marshal([]string{})
		runItem, err := s.Store.CreateInvoiceRunItem(ctx, db.CreateInvoiceRunItemParams{
			RunID:     run.ID,
			ClientID:  t.ClientID,
			SenderID:  t.SenderID,
			Status:    db.InvoiceRunItemStatusEnumCreated,
			InvoiceID: nil,
			Error:     nil,
			Warnings:  warningsBytes,
		})
		if err != nil {
			hadErrors = true
			continue
		}

		inv, _, genErr := s.generateInvoiceForTarget(ctx, generateTargetParams{
			ClientID:        t.ClientID,
			SenderID:        t.SenderID,
			PeriodStart:     periodStart,
			PeriodEnd:       periodEnd,
			BillingTimezone: DefaultBillingTimezone,
			BillingCycle:    DefaultBillingCycle,
			Source:          db.InvoiceSourceEnumAuto,
			RunID:           &run.ID,
		})
		if genErr != nil {
			msg := genErr.Error()
			status := db.InvoiceRunItemStatusEnumFailed
			if errors.Is(genErr, ErrNoBillableItems) || errors.Is(genErr, ErrAutoInvoiceAlreadyExists) {
				// "No billable items" and "already exists" are expected outcomes for a batch run.
				status = db.InvoiceRunItemStatusEnumSkipped
			} else {
				hadErrors = true
			}
			_, _ = s.Store.UpdateInvoiceRunItem(ctx, db.UpdateInvoiceRunItemParams{
				ID:        runItem.ID,
				Status:    db.NullInvoiceRunItemStatusEnum{Valid: true, InvoiceRunItemStatusEnum: status},
				InvoiceID: nil,
				Error:     &msg,
				Warnings:  nil,
			})
			continue
		}

		// Persist generation warnings (if any) on the run item for audit/debugging.
		wb, _ := json.Marshal(inv.Warnings)
		_, _ = s.Store.UpdateInvoiceRunItem(ctx, db.UpdateInvoiceRunItemParams{
			ID:        runItem.ID,
			Status:    db.NullInvoiceRunItemStatusEnum{Valid: true, InvoiceRunItemStatusEnum: db.InvoiceRunItemStatusEnumCreated},
			InvoiceID: &inv.ID,
			Error:     nil,
			Warnings:  wb,
		})
	}

	finalStatus := db.InvoiceRunStatusEnumCompleted
	if hadErrors {
		finalStatus = db.InvoiceRunStatusEnumCompletedWithErrors
	}
	_, err = s.Store.UpdateInvoiceRun(ctx, db.UpdateInvoiceRunParams{
		ID:         run.ID,
		Status:     db.NullInvoiceRunStatusEnum{Valid: true, InvoiceRunStatusEnum: finalStatus},
		FinishedAt: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update invoice_run: %w", err)
	}
	return nil
}

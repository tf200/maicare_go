package invoice

import (
	"context"
	"database/sql"

	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/util"
	"time"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type InvoiceData struct {
	ClientID          int64            `json:"client_id"`
	SenderID          int64            `json:"sender_id"`
	InvoiceDate       time.Time        `json:"invoice_date"`
	InvoiceNumber     string           `json:"invoice_number"`
	InvoiceSequence   int64            `json:"invoice_sequence"`
	PreVatTotalAmount float64          `json:"pre_vat_total"` // Total before VAT
	TotalAmount       float64          `json:"total_amount"`  // Total amount for the invoice
	InvoiceDetails    []InvoiceDetails `json:"invoice_details"`
}

type InvoicePeriod struct {
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
	AcommodationTimeFrame *string   `json:"accommodation_time_frame,omitempty"`
	AmbulanteTotalMinutes *float64  `json:"ambulante_total_minutes,omitempty"`
}

func (s *invoiceService) GenerateInvoiceNumber(ctx context.Context) (string, int64, error) {
	now := time.Now()
	datePart := now.Format("20060102") // YYYYMMDD

	// Get the maximum sequence number for today
	maxSeq, err := s.Store.GetMaxInvoiceSequenceForDate(ctx, now)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get max sequence: %w", err)
	}

	nextSeq := maxSeq + 1
	invoiceNumber := fmt.Sprintf("INV-%s-%04d", datePart, nextSeq)

	return invoiceNumber, nextSeq, nil
}

func (s *invoiceService) GenerateInvoice(req GenerateInvoiceRequest, ctx context.Context) (*GenerateInvoiceResponse, int64, error) {

	clientSender, err := s.Store.GetClientSender(ctx, req.ClientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "GenerateInvoice", "No sender found for client",
				zap.Int64("client_id", req.ClientID))
			return nil, 0, fmt.Errorf("no sender found for client %d", req.ClientID)
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateInvoice", "Database error during sender retrieval",
			zap.Int64("client_id", req.ClientID), zap.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to get client sender: %v", err)
	}

	var warningCount int64
	if req.ClientID <= 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "GenerateInvoice", "Invalid client ID",
			zap.Int64("client_id", req.ClientID))
		return nil, warningCount, fmt.Errorf("invalid client ID: %d", req.ClientID)
	}
	if req.StartDate.IsZero() || req.EndDate.IsZero() {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "GenerateInvoice", "Start date and end date must be specified",
			zap.Int64("client_id", req.ClientID))
		return nil, warningCount, fmt.Errorf("start date and end date must be specified")
	}
	if req.EndDate.Before(req.StartDate) {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "GenerateInvoice", "End date cannot be before start date",
			zap.Int64("client_id", req.ClientID))
		return nil, warningCount, fmt.Errorf("end date cannot be before start date")
	}

	// Get all client contracts
	contracts, err := s.Store.ListClientContracts(ctx, db.ListClientContractsParams{
		ClientID: req.ClientID,
		Limit:    1000,
		Offset:   0,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateInvoice", "Database error during contract retrieval",
			zap.Int64("client_id", req.ClientID), zap.String("error", err.Error()))
		return nil, warningCount, fmt.Errorf("failed to get client contracts for client %d: %w", req.ClientID, err)
	}
	if len(contracts) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "GenerateInvoice", "No contracts found for client",
			zap.Int64("client_id", req.ClientID))
		return nil, warningCount, fmt.Errorf("no contracts found for client %d", req.ClientID)
	}
	var totalInvoiceItems int

	var totalAmount float64
	var totalPreVat float64
	var invoice = make([]InvoiceDetails, len(contracts))

	for i, contract := range contracts {
		billablePeriods, err := s.Store.GetBillablePeriodsForContract(ctx, db.GetBillablePeriodsForContractParams{
			ContractID:       contract.ID,
			InvoiceStartDate: pgtype.Timestamptz{Time: req.StartDate, Valid: true},
			InvoiceEndDate:   pgtype.Timestamptz{Time: req.EndDate, Valid: true},
		})
		if err != nil {
			warningCount++
			continue
		}
		if len(billablePeriods) == 0 {
			continue
		}
		totalInvoiceItems++

		invoice[i] = InvoiceDetails{
			ContractID:    contract.ID,
			Price:         contract.Price,
			ContractType:  contract.CareType,
			PriceTimeUnit: contract.PriceTimeUnit,
			Vat:           float64(*contract.Vat),
			Warnings:      []string{},
			Periods:       []InvoicePeriod{},
		}

		if len(billablePeriods) > 1 {
			invoice[i].Warnings = append(invoice[i].Warnings,
				fmt.Sprintf("multiple billable periods found for contract %d, make sure to verify contract details", contract.ID))
			warningCount++
		}

		for _, period := range billablePeriods {
			var periodItem InvoicePeriod
			periodItem.StartDate = period.BillableStart.Time
			periodItem.EndDate = period.BillableEnd.Time

			if contract.CareType == "accommodation" {
				totals, err := CalculateAccomodationInvoiceTotal(AccommodationInvoiceParams{
					Price:               contract.Price,
					PriceTimeUnit:       contract.PriceTimeUnit,
					VAT:                 float64(*contract.Vat),
					BillablePeriodStart: period.BillableStart.Time,
					BillablePeriodEnd:   period.BillableEnd.Time,
				})
				if err != nil {
					invoice[i].Warnings = append(invoice[i].Warnings,
						fmt.Sprintf("failed to calculate accommodation invoice total for contract %d: %v", contract.ID, err))
					warningCount++
					continue
				}

				invoice[i].PreVatTotal += totals.PreVatTotal
				invoice[i].Total += totals.Total
				periodItem.AcommodationTimeFrame = &totals.TimeFrame
				totalAmount += totals.Total
				totalPreVat += totals.PreVatTotal

			} else if contract.CareType == "ambulante" {
				appointments, err := s.Store.ListClientAppointmentsStartingInRange(ctx, db.ListClientAppointmentsStartingInRangeParams{
					ClientID: contract.ClientID,
					StartDate: pgtype.Timestamp{
						Time:  period.BillableStart.Time,
						Valid: true,
					},
					EndDate: pgtype.Timestamp{
						Time:  period.BillableEnd.Time,
						Valid: true,
					},
				})
				if err != nil {
					invoice[i].Warnings = append(invoice[i].Warnings,
						fmt.Sprintf("failed to get appointments for contract %d: %v", contract.ID, err))
					warningCount++
					continue
				}
				if len(appointments) == 0 {
					invoice[i].Warnings = append(invoice[i].Warnings,
						fmt.Sprintf("no appointments found for ambulante contract %d in the specified date range", contract.ID))
					warningCount++
					continue
				}

				totalMinutes := 0.0
				for _, appointment := range appointments {
					duration := appointment.EndTime.Time.Sub(appointment.StartTime.Time)
					totalMinutes += duration.Minutes()
				}

				totals, err := CalculateAmbulanteInvoiceTotal(AmbulanteInvoiceParams{
					Price:         contract.Price,
					PriceTimeUnit: contract.PriceTimeUnit,
					VAT:           float64(*contract.Vat),
					TotalMinutes:  totalMinutes,
				})
				if err != nil {
					invoice[i].Warnings = append(invoice[i].Warnings,
						fmt.Sprintf("failed to calculate ambulante invoice total for contract %d: %v", contract.ID, err))
					warningCount++
					continue
				}

				invoice[i].PreVatTotal += totals.PreVatTotal
				invoice[i].Total += totals.Total
				periodItem.AmbulanteTotalMinutes = &totals.TotalMinutes
				totalAmount += totals.Total
				totalPreVat += totals.PreVatTotal
			}

			invoice[i].Periods = append(invoice[i].Periods, periodItem)
		}
	}

	if totalInvoiceItems == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "GenerateInvoice", "No billable items found for client",
			zap.Int64("client_id", req.ClientID))
		return nil, warningCount, fmt.Errorf("no billable items found for client %d in the specified date range", req.ClientID)
	}

	invoiceDate := time.Now()
	invoiceNumber, invoiceSequence, err := s.GenerateInvoiceNumber(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateInvoice", "Failed to generate invoice number",
			zap.Int64("client_id", req.ClientID), zap.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	finalInvoice := InvoiceData{
		ClientID: req.ClientID,
		// SenderID:          *contracts[0].SenderID,
		InvoiceNumber:   invoiceNumber,
		InvoiceSequence: invoiceSequence,
		// InvoiceDate:       invoiceDate,
		InvoiceDetails: invoice,
		TotalAmount:    totalAmount,
		// PreVatTotalAmount: totalPreVat,
	}
	invoiceDetailsBytes, err := json.Marshal(finalInvoice.InvoiceDetails)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateInvoice", "Failed to marshal invoice details",
			zap.Int64("client_id", req.ClientID), zap.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to marshal invoice details: %v", err)
	}

	extraContent, err := s.Store.FetchInvoiceTemplateItems(ctx, db.FetchQueryData{
		ClientID:   req.ClientID,
		ContractID: contracts[0].ID,
		SenderID:   clientSender.ID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateInvoice", "Failed to fetch invoice template items",
			zap.Int64("client_id", req.ClientID), zap.String("error", err.Error()))

	}
	if len(extraContent) == 0 {
		extraContent = map[string]string{}
	}

	extraContentBytes, err := json.Marshal(extraContent)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateInvoice", "Failed to marshal extra content",
			zap.Int64("client_id", req.ClientID), zap.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to marshal extra content: %v", err)
	}

	createdInv, err := s.Store.CreateInvoice(ctx, db.CreateInvoiceParams{
		ClientID:        finalInvoice.ClientID,
		SenderID:        &clientSender.ID,
		DueDate:         pgtype.Date{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		TotalAmount:     finalInvoice.TotalAmount,
		InvoiceDetails:  invoiceDetailsBytes,
		InvoiceNumber:   finalInvoice.InvoiceNumber,
		ExtraContent:    extraContentBytes,
		WarningCount:    int32(warningCount),
		IssueDate:       pgtype.Date{Time: invoiceDate, Valid: true},
		InvoiceType:     "standard",
		InvoiceSequence: finalInvoice.InvoiceSequence,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateInvoice", "Failed to create invoice in database",
			zap.Int64("client_id", req.ClientID), zap.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to create invoice in database: %v", err)
	}

	result := &GenerateInvoiceResponse{
		ID:              createdInv.ID,
		InvoiceNumber:   createdInv.InvoiceNumber,
		IssueDate:       createdInv.IssueDate.Time,
		DueDate:         createdInv.DueDate.Time,
		Status:          createdInv.Status,
		InvoiceDetails:  finalInvoice.InvoiceDetails,
		TotalAmount:     createdInv.TotalAmount,
		PdfAttachmentID: createdInv.PdfAttachmentID,
		ExtraContent:    util.ParseJSONToObject(createdInv.ExtraContent),
		ClientID:        createdInv.ClientID,
		SenderID:        createdInv.SenderID,
		UpdatedAt:       createdInv.UpdatedAt.Time,
		CreatedAt:       createdInv.CreatedAt.Time,
	}

	return result, warningCount, nil
}

func (s *invoiceService) BatchGenerateInvoices(ctx context.Context) error {
	// sleep for 10 seconds to allow other services to start
	time.Sleep(10 * time.Second)
	now := time.Now()
	currentIsoWeek, currentYear := now.ISOWeek()

	// Quarterly invoices
	if currentIsoWeek%4 == 1 && now.Weekday() == time.Monday {

		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "BatchGenerateInvoices", "Starting batch invoice generation",
			zap.Int("current_week", currentIsoWeek), zap.Int("current_year", currentYear))

		startDate := time.Date(currentYear, now.Month(), now.Day()-28, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(currentYear, now.Month(), now.Day()-1, 23, 59, 59, 0, time.UTC)
		// Generate invoices for all clients
		clientIDs, err := s.Store.GetAllClientsIDs(ctx)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "BatchGenerateInvoices", "Failed to fetch client IDs",
				zap.String("error", err.Error()))
			return fmt.Errorf("failed to fetch client IDs: %v", err)
		}
		for _, clientID := range clientIDs {
			req := GenerateInvoiceRequest{
				ClientID:  clientID,
				StartDate: startDate,
				EndDate:   endDate,
			}
			_, warnings, err := s.GenerateInvoice(req, ctx)
			if err != nil {
				s.Logger.LogBusinessEvent(logger.LogLevelError, "BatchGenerateInvoices", "Failed to generate invoice",
					zap.Int64("client_id", clientID), zap.String("error", err.Error()))
				continue
			}
			s.Logger.LogBusinessEvent(logger.LogLevelInfo, "BatchGenerateInvoices", "Successfully generated invoice",
				zap.Int64("client_id", clientID), zap.Int64("warnings", warnings))
		}

	} else {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "BatchGenerateInvoices", "Not the scheduled time for batch invoice generation",
			zap.Int("current_week", currentIsoWeek), zap.Int("current_year", currentYear))
	}
	return nil
}

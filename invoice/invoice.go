package invoice

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type InvoiceParams struct {
	ClientID  int64     `json:"client_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

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

// InvoiceDetails contains details for each contract in the invoice
type InvoiceDetails struct {
	ContractID    int64           `json:"contract_id"`
	ContractType  string          `json:"contract_name"`
	Periods       []InvoicePeriod `json:"periods"`
	PreVatTotal   float64         `json:"pre_vat_total_price"`
	Total         float64         `json:"total_price"`
	Vat           float64         `json:"vat"`
	Price         float64         `json:"price"`
	PriceTimeUnit string          `json:"price_time_unit"`
	Warnings      []string        `json:"warnings"`
}

type InvoicePeriod struct {
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
	AcommodationTimeFrame *string   `json:"accommodation_time_frame,omitempty"`
	AmbulanteTotalMinutes *float64  `json:"ambulante_total_minutes,omitempty"`
}

func GenerateInvoiceNumber(ctx context.Context, now time.Time, s *db.Store) (string, int64, error) {
	datePart := now.Format("20060102") // YYYYMMDD

	// Get the maximum sequence number for today
	maxSeq, err := s.GetMaxInvoiceSequenceForDate(ctx, now)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get max sequence: %w", err)
	}

	nextSeq := maxSeq + 1
	invoiceNumber := fmt.Sprintf("INV-%s-%04d", datePart, nextSeq)

	return invoiceNumber, nextSeq, nil
}

func GenerateInvoice(store *db.Store, invoiceData InvoiceParams, ctx context.Context) (*InvoiceData, int64, error) {
	// Validate invoice data
	var warningCount int64
	if invoiceData.ClientID <= 0 {
		return nil, warningCount, fmt.Errorf("invalid client ID: %d", invoiceData.ClientID)
	}
	if invoiceData.StartDate.IsZero() || invoiceData.EndDate.IsZero() {
		return nil, warningCount, fmt.Errorf("start date and end date must be specified")
	}
	if invoiceData.EndDate.Before(invoiceData.StartDate) {
		return nil, warningCount, fmt.Errorf("end date cannot be before start date")
	}

	// Get all client contracts
	contracts, err := store.ListClientContracts(ctx, db.ListClientContractsParams{
		ClientID: invoiceData.ClientID,
		Limit:    1000,
		Offset:   0,
	})
	if err != nil {
		return nil, warningCount, fmt.Errorf("failed to get client contracts for client %d: %w", invoiceData.ClientID, err)
	}
	if len(contracts) == 0 {
		return nil, warningCount, fmt.Errorf("no contracts found for client %d", invoiceData.ClientID)
	}

	var totalAmount float64
	var totalPreVat float64
	var invoice []InvoiceDetails = make([]InvoiceDetails, len(contracts))

	for i, contract := range contracts {
		billablePeriods, err := store.GetBillablePeriodsForContract(ctx, db.GetBillablePeriodsForContractParams{
			ContractID:       contract.ID,
			InvoiceStartDate: pgtype.Timestamptz{Time: invoiceData.StartDate, Valid: true},
			InvoiceEndDate:   pgtype.Timestamptz{Time: invoiceData.EndDate, Valid: true},
		})
		if err != nil {
			invoice[i].Warnings = append(invoice[i].Warnings,
				fmt.Sprintf("failed to get billable periods for contract %d: %v", contract.ID, err))
			warningCount++
			continue
		}
		if len(billablePeriods) == 0 {
			continue
		}

		invoice[i] = InvoiceDetails{
			ContractID:    contract.ID,
			Price:         contract.Price,
			ContractType:  contract.CareType,
			PriceTimeUnit: contract.PriceTimeUnit,
			Vat:           float64(*contract.Vat),
			Warnings:      []string{},
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
				appointments, err := store.ListClientAppointmentsStartingInRange(ctx, db.ListClientAppointmentsStartingInRangeParams{
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

	invoiceDate := time.Now()
	invoiceNumber, invoiceSequence, err := GenerateInvoiceNumber(ctx, invoiceDate, store)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	finalInvoice := InvoiceData{
		ClientID:          invoiceData.ClientID,
		SenderID:          *contracts[0].SenderID,
		InvoiceNumber:     invoiceNumber,
		InvoiceSequence:   invoiceSequence,
		InvoiceDate:       invoiceDate,
		InvoiceDetails:    invoice,
		TotalAmount:       totalAmount,
		PreVatTotalAmount: totalPreVat,
	}

	return &finalInvoice, warningCount, nil
}

func VerifyTotalAmount(invoiceDetails []InvoiceDetails, totalAmount float64) (bool, error) {
	var calculatedTotal float64

	for _, detail := range invoiceDetails {
		calculatedTotal += detail.Total
	}

	if calculatedTotal != totalAmount {
		return false, fmt.Errorf("total amount does not match the sum of invoice details: expected %.2f, got %.2f", totalAmount, calculatedTotal)
	}
	return true, nil

}

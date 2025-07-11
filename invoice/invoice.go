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
	ClientID       int64            `json:"client_id"`
	InvoiceDate    time.Time        `json:"invoice_date"`
	InvoiceDetails []InvoiceDetails `json:"invoice_details"`
}

type InvoiceDetails struct {
	ContractID            int64     `json:"contract_id"`
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
	PreVatTotal           float64   `json:"pre_vat_total_price"`
	Total                 float64   `json:"total_price"`
	Vat                   float64   `json:"vat"`
	AcommodationTimeFrame *string   `json:"accommodation_time_frame,omitempty"`
	AmbulanteTotalMinutes *float64  `json:"ambulante_total_minutes,omitempty"`
	Price                 float64   `json:"price"`
	PriceTimeUnit         string    `json:"price_time_unit"`
	Warnings              []string  `json:"warnings"` // Optional message for errors
}

func GenerateInvoice(store *db.Store, invoiceData InvoiceParams, ctx context.Context) (*InvoiceData, error) {
	// Validate invoice data
	if invoiceData.ClientID <= 0 {
		return nil, fmt.Errorf("invalid client ID: %d", invoiceData.ClientID)
	}
	if invoiceData.StartDate.IsZero() || invoiceData.EndDate.IsZero() {
		return nil, fmt.Errorf("start date and end date must be specified")
	}
	if invoiceData.EndDate.Before(invoiceData.StartDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	// set the invoice date to now

	// get all client contracts
	contracts, err := store.ListClientContracts(ctx, db.ListClientContractsParams{
		ClientID: invoiceData.ClientID})
	if err != nil {
		return nil, fmt.Errorf("failed to get client contracts for client %d: %w", invoiceData.ClientID, err)
	}
	if len(contracts) == 0 {
		return nil, fmt.Errorf("no contracts found for client %d", invoiceData.ClientID)
	}

	var invoice []InvoiceDetails = make([]InvoiceDetails, len(contracts))
	// Iterate through each contract and
	for i, contract := range contracts {
		invoice[i] = InvoiceDetails{
			ContractID:    contract.ID,
			Price:         contract.Price,
			PriceTimeUnit: contract.PriceTimeUnit,
			Vat:           float64(*contract.Vat),
		}

		billablePeriod, err := store.GetBillablePeriodsForContract(ctx, db.GetBillablePeriodsForContractParams{
			ContractID:       contract.ID,
			InvoiceStartDate: pgtype.Timestamptz{Time: invoiceData.StartDate, Valid: true},
			InvoiceEndDate:   pgtype.Timestamptz{Time: invoiceData.EndDate, Valid: true},
		})
		if err != nil {
			invoice[i].Warnings = append(invoice[i].Warnings, fmt.Sprintf("failed to get billable periods for contract %d: %v", contract.ID, err))
			continue // Skip this contract if billable periods cannot be retrieved
		}
		if len(billablePeriod) == 0 {
			continue
		}

		// itterate through the billable periods of the contract
		for _, period := range billablePeriod {
			invoice[i].StartDate = period.BillableStart.Time
			invoice[i].EndDate = period.BillableEnd.Time
			if contract.CareType == "accommodation" {
				accomodationTotals, err := CalculateAccomodationInvoiceTotal(AccommodationInvoiceParams{
					Price:               contract.Price,
					PriceTimeUnit:       contract.PriceTimeUnit,
					VAT:                 float64(*contract.Vat),
					BillablePeriodStart: period.BillableStart.Time,
					BillablePeriodEnd:   period.BillableEnd.Time,
				})
				if err != nil {
					invoice[i].Warnings = append(invoice[i].Warnings, fmt.Sprintf("failed to calculate accommodation invoice total for contract %d: %v", contract.ID, err))
					continue // Skip this contract if calculation fails
				}
				invoice[i].PreVatTotal = accomodationTotals.PreVatTotal
				invoice[i].Total = accomodationTotals.Total
				invoice[i].AcommodationTimeFrame = &accomodationTotals.TimeFrame

			} else if contract.CareType == "ambulante" {
				appointments, err := store.ListClientAppointmentsStartingInRange(ctx, db.ListClientAppointmentsStartingInRangeParams{
					ClientID:  contract.ClientID,
					StartDate: pgtype.Timestamp{Time: period.BillableStart.Time, Valid: true},
					EndDate:   pgtype.Timestamp{Time: period.BillableEnd.Time, Valid: true},
				})
				if err != nil {
					invoice[i].Warnings = append(invoice[i].Warnings, fmt.Sprintf("failed to get appointments for contract %d: %v", contract.ID, err))
					continue // Skip this contract if appointments cannot be retrieved
				}
				if len(appointments) == 0 {
					continue // Skip this contract if no appointments are found
				}
				// Calculate total minutes from appointments
				totalMinutes := 0.0
				for _, appointment := range appointments {
					duration := appointment.EndTime.Time.Sub(appointment.StartTime.Time)
					totalMinutes += duration.Minutes()
				}
				ambulanteTotals, err := CalculateAmbulanteInvoiceTotal(AmbulanteInvoiceParams{
					Price:         contract.Price,
					PriceTimeUnit: contract.PriceTimeUnit,
					VAT:           float64(*contract.Vat),
					TotalMinutes:  totalMinutes,
				})
				if err != nil {
					invoice[i].Warnings = append(invoice[i].Warnings, fmt.Sprintf("failed to calculate ambulante invoice total for contract %d: %v", contract.ID, err))
					continue // Skip this contract if calculation fails
				}
				invoice[i].PreVatTotal = ambulanteTotals.PreVatTotal
				invoice[i].Total = ambulanteTotals.Total
				invoice[i].AmbulanteTotalMinutes = &totalMinutes

			}
		}
	}
	finalInvoice := InvoiceData{
		ClientID:       invoiceData.ClientID,
		InvoiceDate:    time.Now(),
		InvoiceDetails: invoice,
	}
	return &finalInvoice, nil
}

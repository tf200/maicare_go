package invoice

import (
	"fmt"
	"math"
	"time"
)

type AccommodationInvoiceParams struct {
	Price               float64   `json:"price"`
	PriceTimeUnit       string    `json:"price_time_unit"`
	VAT                 float64   `json:"vat"`
	BillablePeriodStart time.Time `json:"billable_period_start"`
	BillablePeriodEnd   time.Time `json:"billable_period_end"`
}

type AccomodationInvoiceTotals struct {
	PreVatTotal float64 `json:"pre_vat_total_price"`
	Total       float64 `json:"total_price"`
	Vat         float64 `json:"vat"`
	TimeFrame   string  `json:"time_frame"`
}

func CalculateAccomodationInvoiceTotal(params AccommodationInvoiceParams) (*AccomodationInvoiceTotals, error) {
	if params.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than zero")
	}
	if params.PriceTimeUnit == "" {
		return nil, fmt.Errorf("price time unit must be specified")
	}
	if params.BillablePeriodStart.IsZero() || params.BillablePeriodEnd.IsZero() {
		return nil, fmt.Errorf("billable period start and end must be specified")
	}
	if params.BillablePeriodEnd.Before(params.BillablePeriodStart) {
		return nil, fmt.Errorf("billable period end cannot be before start")
	}

	// Calculate the number of days in the billable period
	daysInPeriod := params.BillablePeriodEnd.Sub(params.BillablePeriodStart)
	days := int(math.Ceil(daysInPeriod.Hours() / 24))
	if days <= 0 {
		return nil, fmt.Errorf("billable period must be at least one day")
	}

	if params.PriceTimeUnit == "daily" {
		preVatTotal := params.Price * float64(days)
		vat := preVatTotal * (params.VAT / 100)
		total := preVatTotal + vat

		return &AccomodationInvoiceTotals{
			PreVatTotal: preVatTotal,
			Total:       total,
			Vat:         vat,
			TimeFrame:   fmt.Sprintf("%d days", days),
		}, nil
	}
	if params.PriceTimeUnit == "weekly" {
		dailyRate := params.Price / 7
		weeks := days / 7
		preVatTotal := dailyRate * float64(days)
		vat := preVatTotal * (params.VAT / 100)
		total := preVatTotal + vat

		return &AccomodationInvoiceTotals{
			PreVatTotal: preVatTotal,
			Total:       total,
			Vat:         vat,
			TimeFrame:   fmt.Sprintf("%d weeks", weeks),
		}, nil
	}

	return nil, fmt.Errorf("unsupported price time unit: %s", params.PriceTimeUnit)

}

type AmbulanteInvoiceParams struct {
	Price         float64 `json:"price"`
	PriceTimeUnit string  `json:"price_time_unit"`
	VAT           float64 `json:"vat"`
	TotalMinutes  float64 `json:"total_minutes"` // Total minutes of care provided in the billable period
}
type AmbulanteInvoiceTotals struct {
	PreVatTotal  float64 `json:"pre_vat_total_price"`
	Total        float64 `json:"total_price"`
	Vat          float64 `json:"vat"`
	TotalMinutes float64 `json:"total_minutes"` // Total minutes of care provided in the billable period
}

func CalculateAmbulanteInvoiceTotal(params AmbulanteInvoiceParams) (*AmbulanteInvoiceTotals, error) {
	if params.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than zero")
	}
	if params.PriceTimeUnit == "" {
		return nil, fmt.Errorf("price time unit must be specified")
	}
	if params.TotalMinutes <= 0 {
		return nil, fmt.Errorf("total minutes must be greater than zero")
	}

	if params.PriceTimeUnit == "minute" {
		preVatTotal := params.Price * params.TotalMinutes
		vat := preVatTotal * (params.VAT / 100)
		total := preVatTotal + vat

		return &AmbulanteInvoiceTotals{
			PreVatTotal:  preVatTotal,
			Total:        total,
			Vat:          vat,
			TotalMinutes: params.TotalMinutes,
		}, nil
	}
	if params.PriceTimeUnit == "hourly" {
		hours := params.TotalMinutes / 60
		preVatTotal := params.Price * hours
		vat := preVatTotal * (params.VAT / 100)
		total := preVatTotal + vat

		return &AmbulanteInvoiceTotals{
			PreVatTotal:  preVatTotal,
			Total:        total,
			Vat:          vat,
			TotalMinutes: params.TotalMinutes,
		}, nil
	}
	return nil, fmt.Errorf("unsupported price time unit: %s", params.PriceTimeUnit)
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

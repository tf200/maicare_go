package pdf

import (
	"context"
	"log"
	"maicare_go/bucket"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateInvoicePDF(t *testing.T) {

	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load conf %v", err)
	}

	testb2Client, err := bucket.NewB2Client(config)
	if err != nil {
		log.Fatal("cannot create b2 client:", err)
	}

	data := InvoicePDFData{
		SenderName:           "Stichting Zorgklant",
		SenderContactPerson:  "Jan Jansen",
		SenderAddressLine1:   "Voorbeeldstraat 12",
		SenderPostalCodeCity: "1234 AB Amsterdam",
		InvoiceNumber:        "INV-20250718-001",
		InvoiceDate:          time.Now(),
		DueDate:              time.Now().AddDate(0, 0, 30),
		InvoiceDetails: []InvoiceDetail{
			{
				CareType:      "Accommodatie",
				Periods:       []InvoicePeriod{{StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7), AcommodationTimeFrame: "30 days", AmbulanteTotalMinutes: 0}},
				Price:         75.00,
				PriceTimeUnit: "uur",
				PreVatTotal:   150.00,
				Total:         181.50,
			},
			{
				CareType:      "Ambulante",
				Periods:       []InvoicePeriod{{StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7), AcommodationTimeFrame: "", AmbulanteTotalMinutes: 120}},
				Price:         60.00,
				PriceTimeUnit: "uur",
				PreVatTotal:   120.00,
				Total:         145.20,
			},
		},
		TotalAmount: 326.70,
		ExtraItems: map[string]string{
			"Client Geboortedatum": "01-01-1990",
			"Financieringsoptie":   "PGB",
			"Opmerkingen":          "Geen bijzonderheden.",
		},
	}

	pdfBytes, _, _, err := GenerateAndUploadInvoicePDF(context.Background(), data, testb2Client)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)

}

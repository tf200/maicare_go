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

func TestGenerateContractPDF(t *testing.T) {
	// Mock data for the contract
	contractData := ContractData{
		ID:                2,
		Status:            "approved",
		StartDate:         "2024-01-01",
		EndDate:           "2025-01-01",
		ReminderPeriod:    30,
		CareType:          "ambulante",
		SenderName:        "Maicare BV",
		SenderAddress:     "123 Care St, Health City",
		SenderContactInfo: "info@maicare.com, +31 20 123 4567",
		ClientFirstName:   "Jane",
		ClientLastName:    "Doe",
		ClientAddress:     "456 Client Ave, Wellness Town",
		ClientContactInfo: "jane.doe@example.com, +31 6 9876 5432",

		CareName:         "Standard Care Package",
		FinancingAct:     "WMO",
		FinancingOption:  "ZIN",
		Hours:            20.0,
		HoursType:        "weekly",
		AmbulanteDisplay: "block",
		Price:            1000.0,
		PriceTimeUnit:    "monthly",
		Vat:              21.0,
		TypeName:         "Standard Care Package",
		GenerationDate:   time.Now().Format("2006-01-02"),
	}

	pdfBytes, err := GenerateContractPDF(contractData)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)

	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load config %v", err)
	}

	testb2Client, err := bucket.NewObjectStorageClient(context.Background(), config)
	if err != nil {
		log.Fatal("cannot create b2 client:", err)
	}

	ctx := context.Background()
	filename, err := UploadContractPDF(ctx, pdfBytes, contractData.ID, testb2Client)
	require.NoError(t, err)
	require.NotEmpty(t, filename)

	log.Println("Contract PDF uploaded successfully:", filename)
}

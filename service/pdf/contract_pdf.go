package pdf

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"
)

// ContractData represents the data structure for the care contract template.
type ContractData struct {
	// Header Information
	ID     int64  `json:"ContractID"` // Corresponds to {.ContractID}
	Status string `json:"Status"`     // Corresponds to {.Status} and used for status-badge class

	// Contract Periods
	StartDate      string `json:"StartDate"`      // Corresponds to {.StartDate} and [START_DATE] in terms
	EndDate        string `json:"EndDate"`        // Corresponds to {.EndDate} and [END_DATE] in terms
	ReminderPeriod int    `json:"ReminderPeriod"` // Corresponds to {.ReminderPeriod}

	// Parties
	SenderName        string `json:"SenderName"`
	SenderStreet      string `json:"SenderStreet"`
	SenderHouseNumber string `json:"SenderHouseNumber"`
	SenderPostalCode  string `json:"SenderPostalCode"`
	SenderCity        string `json:"SenderCity"`
	SenderContactInfo string `json:"SenderContactInfo"`

	ClientFirstName   string `json:"ClientFirstName"`   // Corresponds to {.ClientFirstName}
	ClientLastName    string `json:"ClientLastName"`    // Corresponds to {.ClientLastName}
	ClientAddress     string `json:"ClientAddress"`     // Corresponds to {.ClientAddress}
	ClientContactInfo string `json:"ClientContactInfo"` // Corresponds to {.ClientContactInfo}

	// Care Specifications
	CareType        string `json:"CareType"`        // Corresponds to {.CareType}
	CareName        string `json:"CareName"`        // Corresponds to {.CareName}
	FinancingAct    string `json:"FinancingAct"`    // Corresponds to {.FinancingAct}
	FinancingOption string `json:"FinancingOption"` // Corresponds to {.FinancingOption}

	// Ambulante Care Hours (conditional display)
	Hours            float64 `json:"Hours"`            // Corresponds to {.Hours}
	HoursType        string  `json:"HoursType"`        // Corresponds to {.HoursType}
	AmbulanteDisplay string  `json:"AmbulanteDisplay"` // "block" or "none" for [AMBULANTE_DISPLAY]

	// Financial Terms
	Price          float64 `json:"Price"`          // Corresponds to {.Price}
	PriceTimeUnit  string  `json:"PriceTimeUnit"`  // Corresponds to {.PriceTimeUnit}
	Vat            float64 `json:"Vat"`            // Corresponds to {.Vat}
	TypeName       string  `json:"TypeName"`       // Corresponds to {.TypeName}
	GenerationDate string  `json:"GenerationDate"` // Date when the contract was generated
}

// GenerateIncidentPDF generates a PDF from incident data and returns the PDF bytes
func (s *pdfService) generateContractPDF(contractData ContractData) (multipart.File, error) {
	headerLines := []string{
		fmt.Sprintf("Contract ID: %d", contractData.ID),
		fmt.Sprintf("Status: %s", contractData.Status),
		fmt.Sprintf("Generation date: %s", contractData.GenerationDate),
		fmt.Sprintf("Period: %s to %s", contractData.StartDate, contractData.EndDate),
		fmt.Sprintf("Reminder period (days): %d", contractData.ReminderPeriod),
	}

	sections := []documentSection{
		{
			Title: "Sender",
			Lines: []string{
				fmt.Sprintf("Name: %s", contractData.SenderName),
				fmt.Sprintf("Address: %s %s, %s %s", contractData.SenderStreet, contractData.SenderHouseNumber, contractData.SenderPostalCode, contractData.SenderCity),
				fmt.Sprintf("Contact info: %s", contractData.SenderContactInfo),
			},
		},
		{
			Title: "Client",
			Lines: []string{
				fmt.Sprintf("Name: %s %s", contractData.ClientFirstName, contractData.ClientLastName),
				fmt.Sprintf("Address: %s", contractData.ClientAddress),
				fmt.Sprintf("Contact info: %s", contractData.ClientContactInfo),
			},
		},
		{
			Title: "Care specification",
			Lines: []string{
				fmt.Sprintf("Care type: %s", contractData.CareType),
				fmt.Sprintf("Care name: %s", contractData.CareName),
				fmt.Sprintf("Financing act: %s", contractData.FinancingAct),
				fmt.Sprintf("Financing option: %s", contractData.FinancingOption),
				fmt.Sprintf("Hours: %.2f (%s)", contractData.Hours, contractData.HoursType),
			},
		},
		{
			Title: "Financial terms",
			Lines: []string{
				fmt.Sprintf("Price: EUR %.2f per %s", contractData.Price, contractData.PriceTimeUnit),
				fmt.Sprintf("VAT: %.2f%%", contractData.Vat),
				fmt.Sprintf("Contract type: %s", contractData.TypeName),
			},
		},
	}

	pdfBytes, err := buildSectionsPDF("Care Contract", headerLines, sections)
	if err != nil {
		return nil, fmt.Errorf("failed to generate contract pdf: %w", err)
	}

	return toMultipartFile(pdfBytes), nil
}

// UploadIncidentPDF uploads a PDF to B2 with a generated filename
func (s *pdfService) uploadContractPDF(ctx context.Context, pdfFile multipart.File, contractID int64) (string, error) {
	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("contract/%s/contract-%d.pdf", timestamp, contractID)

	// Upload to B2
	key, _, err := s.bucketClient.Upload(ctx, pdfFile, filename, "application/pdf")
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF to B2: %w", err)
	}
	return key, nil
}

// Helper function to do both operations if needed
func (s *pdfService) GenerateAndUploadContractPDF(ctx context.Context, contractData ContractData) (string, error) {
	// Generate PDF
	pdfFile, err := s.generateContractPDF(contractData)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload PDF
	fileURL, err := s.uploadContractPDF(ctx, pdfFile, contractData.ID)
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF: %w", err)
	}

	return fileURL, nil
}

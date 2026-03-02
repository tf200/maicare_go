package pdf

import (
	"context"

	"maicare_go/bucket"
)

type PdfService interface {
	GenerateAppointmentCardPDF(ctx context.Context, cardData AppointmentCard) ([]byte, error)
	GenerateAndUploadAppointmentCardPDF(ctx context.Context, cardData AppointmentCard) (string, error)
	GenerateAndUploadInvoicePDF(ctx context.Context, invoiceData InvoicePDFData) (string, int64, error)
	GenerateAndUploadContractPDF(ctx context.Context, contractData ContractData) (string, error)
	GenerateIncidentPDF(ctx context.Context, incidentData IncidentReportData) ([]byte, error)
	GenerateAndUploadIncidentPDF(ctx context.Context, incidentData IncidentReportData) (string, error)
}

type pdfService struct {
	bucketClient bucket.ObjectStorageInterface
}

func NewPdfService(bucketClient bucket.ObjectStorageInterface) PdfService {
	return &pdfService{
		bucketClient: bucketClient,
	}
}

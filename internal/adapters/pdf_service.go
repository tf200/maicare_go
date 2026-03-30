package adapters

import (
	"context"

	"maicare_go/internal/domain"
	pkgpdf "maicare_go/pkg/pdf"
)

type PDFServiceAdapter struct {
	service interface {
		GenerateAppointmentCardPDF(ctx context.Context, cardData pkgpdf.AppointmentCard) ([]byte, error)
		GenerateAndUploadAppointmentCardPDF(ctx context.Context, cardData pkgpdf.AppointmentCard) (string, error)
		GenerateAndUploadInvoicePDF(ctx context.Context, invoiceData pkgpdf.InvoicePDFData) (string, int64, error)
		GenerateAndUploadContractPDF(ctx context.Context, contractData pkgpdf.ContractData) (string, error)
		GenerateIncidentPDF(ctx context.Context, incidentData pkgpdf.IncidentReportData) ([]byte, error)
		GenerateAndUploadIncidentPDF(ctx context.Context, incidentData pkgpdf.IncidentReportData) (string, error)
	}
}

func NewPDFServiceAdapter(service interface {
	GenerateAppointmentCardPDF(ctx context.Context, cardData pkgpdf.AppointmentCard) ([]byte, error)
	GenerateAndUploadAppointmentCardPDF(ctx context.Context, cardData pkgpdf.AppointmentCard) (string, error)
	GenerateAndUploadInvoicePDF(ctx context.Context, invoiceData pkgpdf.InvoicePDFData) (string, int64, error)
	GenerateAndUploadContractPDF(ctx context.Context, contractData pkgpdf.ContractData) (string, error)
	GenerateIncidentPDF(ctx context.Context, incidentData pkgpdf.IncidentReportData) ([]byte, error)
	GenerateAndUploadIncidentPDF(ctx context.Context, incidentData pkgpdf.IncidentReportData) (string, error)
}) domain.PDFService {
	return &PDFServiceAdapter{service: service}
}

func (a *PDFServiceAdapter) GenerateAppointmentCardPDF(ctx context.Context, cardData domain.AppointmentCardPDF) ([]byte, error) {
	return a.service.GenerateAppointmentCardPDF(ctx, pkgpdf.AppointmentCard{
		ID:                     cardData.ID,
		ClientName:             cardData.ClientName,
		Date:                   cardData.Date,
		Mentor:                 cardData.Mentor,
		GeneralInformation:     cardData.GeneralInformation,
		ImportantContacts:      cardData.ImportantContacts,
		HouseholdInfo:          cardData.HouseholdInfo,
		OrganizationAgreements: cardData.OrganizationAgreements,
		YouthOfficerAgreements: cardData.YouthOfficerAgreements,
		TreatmentAgreements:    cardData.TreatmentAgreements,
		SmokingRules:           cardData.SmokingRules,
		Work:                   cardData.Work,
		SchoolInternship:       cardData.SchoolInternship,
		Travel:                 cardData.Travel,
		Leave:                  cardData.Leave,
	})
}

func (a *PDFServiceAdapter) GenerateAndUploadAppointmentCardPDF(ctx context.Context, cardData domain.AppointmentCardPDF) (string, error) {
	return a.service.GenerateAndUploadAppointmentCardPDF(ctx, pkgpdf.AppointmentCard{
		ID:                     cardData.ID,
		ClientName:             cardData.ClientName,
		Date:                   cardData.Date,
		Mentor:                 cardData.Mentor,
		GeneralInformation:     cardData.GeneralInformation,
		ImportantContacts:      cardData.ImportantContacts,
		HouseholdInfo:          cardData.HouseholdInfo,
		OrganizationAgreements: cardData.OrganizationAgreements,
		YouthOfficerAgreements: cardData.YouthOfficerAgreements,
		TreatmentAgreements:    cardData.TreatmentAgreements,
		SmokingRules:           cardData.SmokingRules,
		Work:                   cardData.Work,
		SchoolInternship:       cardData.SchoolInternship,
		Travel:                 cardData.Travel,
		Leave:                  cardData.Leave,
	})
}

func (a *PDFServiceAdapter) GenerateAndUploadInvoicePDF(ctx context.Context, invoiceData domain.InvoicePDF) (string, int64, error) {
	return a.service.GenerateAndUploadInvoicePDF(ctx, pkgpdf.InvoicePDFData{
		ID:                  invoiceData.ID,
		SenderName:          invoiceData.SenderName,
		SenderContactPerson: invoiceData.SenderContactPerson,
		SenderStreet:        invoiceData.SenderStreet,
		SenderHouseNumber:   invoiceData.SenderHouseNumber,
		SenderPostalCode:    invoiceData.SenderPostalCode,
		SenderCity:          invoiceData.SenderCity,
		InvoiceNumber:       invoiceData.InvoiceNumber,
		InvoiceDate:         invoiceData.InvoiceDate,
		DueDate:             invoiceData.DueDate,
		InvoiceDetails:      toPkgInvoiceDetails(invoiceData.InvoiceDetails),
		TotalAmount:         invoiceData.TotalAmount,
		ExtraItems:          invoiceData.ExtraItems,
	})
}

func (a *PDFServiceAdapter) GenerateAndUploadContractPDF(ctx context.Context, contractData domain.ContractPDF) (string, error) {
	return a.service.GenerateAndUploadContractPDF(ctx, pkgpdf.ContractData{
		ID:                contractData.ID,
		Status:            contractData.Status,
		StartDate:         contractData.StartDate,
		EndDate:           contractData.EndDate,
		ReminderPeriod:    contractData.ReminderPeriod,
		SenderName:        contractData.SenderName,
		SenderStreet:      contractData.SenderStreet,
		SenderHouseNumber: contractData.SenderHouseNumber,
		SenderPostalCode:  contractData.SenderPostalCode,
		SenderCity:        contractData.SenderCity,
		SenderContactInfo: contractData.SenderContactInfo,
		ClientFirstName:   contractData.ClientFirstName,
		ClientLastName:    contractData.ClientLastName,
		ClientAddress:     contractData.ClientAddress,
		ClientContactInfo: contractData.ClientContactInfo,
		CareType:          contractData.CareType,
		CareName:          contractData.CareName,
		FinancingAct:      contractData.FinancingAct,
		FinancingOption:   contractData.FinancingOption,
		Hours:             contractData.Hours,
		HoursType:         contractData.HoursType,
		AmbulanteDisplay:  contractData.AmbulanteDisplay,
		Price:             contractData.Price,
		PriceTimeUnit:     contractData.PriceTimeUnit,
		Vat:               contractData.Vat,
		TypeName:          contractData.TypeName,
		GenerationDate:    contractData.GenerationDate,
	})
}

func (a *PDFServiceAdapter) GenerateIncidentPDF(ctx context.Context, incidentData domain.IncidentReportPDF) ([]byte, error) {
	return a.service.GenerateIncidentPDF(ctx, toPkgIncidentReport(incidentData))
}

func (a *PDFServiceAdapter) GenerateAndUploadIncidentPDF(ctx context.Context, incidentData domain.IncidentReportPDF) (string, error) {
	return a.service.GenerateAndUploadIncidentPDF(ctx, toPkgIncidentReport(incidentData))
}

func toPkgIncidentReport(data domain.IncidentReportPDF) pkgpdf.IncidentReportData {
	return pkgpdf.IncidentReportData{
		ID:                      data.ID,
		EmployeeID:              data.EmployeeID,
		EmployeeFirstName:       data.EmployeeFirstName,
		EmployeeLastName:        data.EmployeeLastName,
		LocationID:              data.LocationID,
		ReporterInvolvement:     data.ReporterInvolvement,
		InformedParties:         data.InformedParties,
		OccurredAt:              data.OccurredAt,
		IncidentType:            data.IncidentType,
		SeverityOfIncident:      data.SeverityOfIncident,
		IncidentExplanation:     data.IncidentExplanation,
		RecurrenceRisk:          data.RecurrenceRisk,
		IncidentPreventSteps:    data.IncidentPreventSteps,
		IncidentTakenMeasures:   data.IncidentTakenMeasures,
		CauseCategories:         data.CauseCategories,
		CauseExplanation:        data.CauseExplanation,
		PhysicalInjury:          data.PhysicalInjury,
		PhysicalInjuryDesc:      data.PhysicalInjuryDesc,
		PsychologicalDamage:     data.PsychologicalDamage,
		PsychologicalDamageDesc: data.PsychologicalDamageDesc,
		NeededConsultation:      data.NeededConsultation,
		FollowUpActions:         data.FollowUpActions,
		FollowUpNotes:           data.FollowUpNotes,
		IsEmployeeAbsent:        data.IsEmployeeAbsent,
		AdditionalDetails:       data.AdditionalDetails,
		ClientID:                data.ClientID,
		ClientFirstName:         data.ClientFirstName,
		ClientLastName:          data.ClientLastName,
		LocationName:            data.LocationName,
	}
}

func toPkgInvoiceDetails(details []domain.InvoiceDetailPDF) []pkgpdf.InvoiceDetail {
	result := make([]pkgpdf.InvoiceDetail, len(details))
	for i, detail := range details {
		result[i] = pkgpdf.InvoiceDetail{
			CareType:      detail.CareType,
			Periods:       toPkgInvoicePeriods(detail.Periods),
			Price:         detail.Price,
			PriceTimeUnit: detail.PriceTimeUnit,
			PreVatTotal:   detail.PreVatTotal,
			Total:         detail.Total,
		}
	}
	return result
}

func toPkgInvoicePeriods(periods []domain.InvoicePeriodPDF) []pkgpdf.InvoicePeriod {
	result := make([]pkgpdf.InvoicePeriod, len(periods))
	for i, period := range periods {
		result[i] = pkgpdf.InvoicePeriod{
			StartDate:             period.StartDate,
			EndDate:               period.EndDate,
			AcommodationTimeFrame: period.AcommodationTimeFrame,
			AmbulanteTotalMinutes: period.AmbulanteTotalMinutes,
		}
	}
	return result
}

var _ domain.PDFService = (*PDFServiceAdapter)(nil)

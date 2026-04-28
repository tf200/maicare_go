package adapters

import (
	"context"

	"maicare_go/internal/domain"
	pkgpdf "maicare_go/pkg/pdf"
)

type IncidentPDFGeneratorAdapter struct {
	service interface {
		GenerateIncidentPDF(ctx context.Context, incidentData pkgpdf.IncidentReportData) ([]byte, error)
	}
}

func NewIncidentPDFGeneratorAdapter(service interface {
	GenerateIncidentPDF(ctx context.Context, incidentData pkgpdf.IncidentReportData) ([]byte, error)
}) domain.IncidentPDFGenerator {
	return &IncidentPDFGeneratorAdapter{service: service}
}

func (a *IncidentPDFGeneratorAdapter) GenerateIncidentPDF(ctx context.Context, data domain.IncidentPDFData) ([]byte, error) {
	return a.service.GenerateIncidentPDF(ctx, pkgpdf.IncidentReportData{
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
	})
}

var _ domain.IncidentPDFGenerator = (*IncidentPDFGeneratorAdapter)(nil)

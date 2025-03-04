package tasks

import (
	"context"
	"fmt"
	"log"
	"maicare_go/email"
	"maicare_go/pdf"

	"github.com/goccy/go-json"

	"github.com/hibiken/asynq"
)

func (processor *AsynqServer) ProcessEmailTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal email task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if p.To == "" || p.UserEmail == "" || p.UserPassword == "" {
		return fmt.Errorf("invalid email payload: missing required fields: %w", asynq.SkipRetry)

	}

	log.Printf("Sending email to %s", p.To)

	err := processor.smtp.SendCredentials(ctx, []string{p.To}, email.Credentials{Email: p.UserEmail, Password: p.UserPassword})
	if err != nil {
		log.Printf("Failed to send email to %s: %v", p.To, err)
		return fmt.Errorf("failed to send email to %s: %v: %w", p.To, err, asynq.SkipRetry)
	}

	return nil
}

func (processor *AsynqServer) ProcessIncidentTask(ctx context.Context, t *asynq.Task) error {
	var p IncidentPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal incident task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	incidentData := pdf.IncidentReportData{
		ID:                      p.ID,
		EmployeeID:              p.EmployeeID,
		EmployeeFirstName:       p.EmployeeFirstName,
		EmployeeLastName:        p.EmployeeLastName,
		LocationID:              p.LocationID,
		ReporterInvolvement:     p.ReporterInvolvement,
		InformWho:               p.InformWho,
		IncidentDate:            p.IncidentDate,
		IncidentType:            p.IncidentType,
		PassingAway:             p.PassingAway,
		SelfHarm:                p.SelfHarm,
		Violence:                p.Violence,
		FireWaterDamage:         p.FireWaterDamage,
		Accident:                p.Accident,
		ClientAbsence:           p.ClientAbsence,
		Medicines:               p.Medicines,
		Organization:            p.Organization,
		UseProhibitedSubstances: p.UseProhibitedSubstances,
		OtherNotifications:      p.OtherNotifications,
		SeverityOfIncident:      p.SeverityOfIncident,
		IncidentExplanation:     p.IncidentExplanation,
		RecurrenceRisk:          p.RecurrenceRisk,
		IncidentPreventSteps:    p.IncidentPreventSteps,
		IncidentTakenMeasures:   p.IncidentTakenMeasures,
		Technical:               p.Technical,
		Organizational:          p.Organizational,
		MeseWorker:              p.MeseWorker,
		ClientOptions:           p.ClientOptions,
		OtherCause:              p.OtherCause,
		CauseExplanation:        p.CauseExplanation,
		PhysicalInjury:          p.PhysicalInjury,
		PhysicalInjuryDesc:      p.PhysicalInjuryDesc,
		PsychologicalDamage:     p.PsychologicalDamage,
		PsychologicalDamageDesc: p.PsychologicalDamageDesc,
		NeededConsultation:      p.NeededConsultation,
		Succession:              p.Succession,
		SuccessionDesc:          p.SuccessionDesc,
		Other:                   p.Other,
		OtherDesc:               p.OtherDesc,
		AdditionalAppointments:  p.AdditionalAppointments,
		EmployeeAbsenteeism:     p.EmployeeAbsenteeism,
		ClientID:                p.ClientID,
		LocationName:            p.LocationName,
	}

	pdfName, err := pdf.GenerateAndUploadIncidentPDF(ctx, incidentData, processor.b2Bucket)
	if err != nil {
		log.Printf("Failed to generate and upload incident PDF: %v", err)
		return fmt.Errorf("failed to generate and upload incident PDF: %v: %w", err, asynq.SkipRetry)
	}

	err = processor.smtp.SendIncident(ctx, p.To, email.Incident{
		IncidentID:   p.ID,
		IncidentType: p.IncidentType,
		Severity:     p.SeverityOfIncident,
		Location:     p.LocationName,
		ReportedBy:   fmt.Sprintf("%s %s", p.EmployeeFirstName, p.EmployeeLastName),
		DocumentLink: pdfName,
	})
	if err != nil {
		log.Printf("Failed to send incident email to %s: %v", p.To, err)
		return fmt.Errorf("failed to send incident email to %s: %v: %w", p.To, err, asynq.SkipRetry)
	}

	return nil
}

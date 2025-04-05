package async

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

	err := processor.smtp.SendCredentials(ctx, []string{p.To}, email.Credentials{Email: p.UserEmail, Password: p.UserPassword, Name: p.Name})
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

// ProcessNotificationTask handles tasks of type TypeNotificationSend.
// It decodes the payload and delegates to the NotificationService.
func (a *AsynqServer) ProcessNotificationTask(ctx context.Context, t *asynq.Task) error {
	var payload NotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		// Return a non-nil error to indicate failure, but don't retry if payload is invalid
		return fmt.Errorf("failed to unmarshal notification payload: %w: %v", asynq.SkipRetry, err)
	}

	log.Printf("Received notification task: %+v", payload) // Log received payload

	// Ensure the notification service is available
	if a.notificationService == nil {
		// Don't retry if the fundamental dependency is missing
		return fmt.Errorf("notification service not initialized on AsynqServer: %w", asynq.SkipRetry)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		// Log the error and return it
		log.Printf("Failed to marshal notification payload: %v", err)
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	// Delegate the actual work to the notification service
	err = a.notificationService.CreateAndDeliver(ctx, payloadBytes)
	if err != nil {
		// Log the error from the service
		log.Printf("Error processing notification task (ID: %s, Type: %s): %v", t.ResultWriter().TaskID(), payload.Type, err)
		// Return the error so Asynq can handle retries based on its configuration
		return fmt.Errorf("notification service failed to process task: %w", err)
	}

	log.Printf("Successfully processed notification task (ID: %s, Type: %s)", t.ResultWriter().TaskID(), payload.Type)
	return nil
}

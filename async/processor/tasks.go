package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"maicare_go/async/aclient"
	db "maicare_go/db/sqlc"
	"maicare_go/email"
	"maicare_go/service/notification"
	"maicare_go/service/pdf"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

func (processor *AsynqServer) ProcessEmailTask(ctx context.Context, t *asynq.Task) error {
	var p aclient.EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal email task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if p.To == "" || p.UserEmail == "" || p.UserPassword == "" {
		return fmt.Errorf("invalid email payload: missing required fields: %w", asynq.SkipRetry)
	}

	log.Printf("Sending email to %s", p.To)

	err := processor.brevoConf.SendCredentials(ctx, []string{p.To}, email.Credentials{Email: p.UserEmail, Password: p.UserPassword, Name: p.Name})
	if err != nil {
		log.Printf("Failed to send email to %s: %v", p.To, err)
		return fmt.Errorf("failed to send email to %s: %v: %w", p.To, err, asynq.SkipRetry)
	}

	return nil
}

func (processor *AsynqServer) ProcessIncidentTask(ctx context.Context, t *asynq.Task) error {
	var p aclient.IncidentPayload
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

	pdfName, err := processor.service.PDFService.GenerateAndUploadIncidentPDF(ctx, incidentData)
	if err != nil {
		log.Printf("Failed to generate and upload incident PDF: %v", err)
		return fmt.Errorf("failed to generate and upload incident PDF: %v: %w", err, asynq.SkipRetry)
	}

	err = processor.brevoConf.SendIncident(ctx, p.To, email.Incident{
		IncidentID:   p.ID.String(),
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
	var payload notification.NotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		// Return a non-nil error to indicate failure, but don't retry if payload is invalid
		return fmt.Errorf("failed to unmarshal notification payload: %w: %v", asynq.SkipRetry, err)
	}

	log.Printf("Received notification task: %+v", payload) // Log received payload

	// Ensure the notification service is available
	if a.service.NotificationService == nil {
		// Don't retry if the fundamental dependency is missing
		return fmt.Errorf("notification service not initialized on AsynqServer: %w", asynq.SkipRetry)
	}

	// Delegate the actual work to the notification service
	err := a.service.NotificationService.CreateAndDeliver(ctx, payload)
	if err != nil {
		// Log the error from the service
		log.Printf("Error processing notification task (ID: %s, Type: %s): %v", t.ResultWriter().TaskID(), payload.Type, err)
		// Return the error so Asynq can handle retries based on its configuration
		return fmt.Errorf("notification service failed to process task: %w", err)
	}

	log.Printf("Successfully processed notification task (ID: %s, Type: %s)", t.ResultWriter().TaskID(), payload.Type)
	return nil
}

func (processor *AsynqServer) ProcessRegistrationFormTask(ctx context.Context, t *asynq.Task) error {
	var p aclient.AcceptedRegistrationFormPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal incident task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	formData := email.AcceptedRegitrationForm{
		ReferrerName:        p.ReferrerName,
		ChildName:           p.ChildName,
		ChildBSN:            p.ChildBSN,
		AppointmentDate:     p.AppointmentDate,
		AppointmentLocation: p.AppointmentLocation,
	}

	err := processor.brevoConf.SendAcceptedRegistrationForm(ctx, []string{p.To}, formData)
	if err != nil {
		log.Printf("Failed to send registration Form Email %s: %v", p.To, err)
		return fmt.Errorf("failed to send incident email to %s: %v: %w", p.To, err, asynq.SkipRetry)
	}

	return nil
}

func (processor *AsynqServer) ProcessProcessRegistrationFormEmailTask(ctx context.Context, t *asynq.Task) error {
	var p aclient.ProcessRegistrationFormEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal process registration form email task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	emailData := email.ProcessRegistrationForm{
		RecipientName: p.ReferrerName,
		ClientName:    p.ClientName,
		Location:      p.Location,
		Link:          p.Link,
	}

	err := processor.brevoConf.SendProcessRegistrationForm(ctx, p.To, emailData)
	if err != nil {
		log.Printf("Failed to send process registration form email to %v: %v", p.To, err)
		return fmt.Errorf("failed to send email: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}

func (c *AsynqServer) ProcessContractRemiderTask(ctx context.Context, t *asynq.Task) error {
	contractsToBeReminded, err := c.store.ListContractsTobeReminded(ctx)
	if err != nil {
		log.Printf("Failed to list contracts to be reminded: %v", err)
		return fmt.Errorf("failed to list contracts to be reminded: %v: %w", err, asynq.SkipRetry)
	}

	if len(contractsToBeReminded) == 0 {
		log.Println("No contracts to be reminded")
		return nil
	}

	for _, contract := range contractsToBeReminded {
		log.Printf("Processing reminder for contract ID: %d", contract.ID)

		reminder, err := c.store.CreateContractReminder(ctx, db.CreateContractReminderParams{
			ContractID:     contract.ID,
			ReminderSentAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		)
		if err != nil {
			log.Printf("Failed to create contract reminder for contract ID %d: %v", contract.ID, err)
			return fmt.Errorf("failed to create contract reminder for contract ID %d: %v: %w", contract.ID, err, asynq.SkipRetry)
		}

		log.Printf("Created contract reminder with ID: %d for contract ID: %d", reminder.ID, contract.ID)

		notificationData := notification.ClientContractReminderData{
			ClientID:           contract.ClientID,
			ClientFirstName:    contract.ClientFirstName,
			ClientLastName:     contract.ClientLastName,
			ContractID:         contract.ID,
			CareType:           string(contract.CareType),
			ContractStart:      contract.StartDate.Time,
			ContractEnd:        contract.EndDate.Time,
			ReminderType:       string(reminder.ReminderType),
			LastReminderSentAt: &reminder.ReminderSentAt.Time,
		}

		adminUsers, err := c.store.GetAllAdminUsers(ctx)
		if err != nil {
			log.Printf("Failed to get admin users: %v", err)
			return fmt.Errorf("failed to get admin users: %v: %w", err, asynq.SkipRetry)
		}

		if len(adminUsers) == 0 {
			log.Println("No admin users found to notify")
			return nil // No admin users to notify, but we can still create the reminder
		}

		notificationPayload := notification.NotificationPayload{
			RecipientUserIDs: make([]uuid.UUID, len(adminUsers)),
			Type:             notification.TypeClientContractReminder,
			Data: notification.NotificationData{
				ClientContractReminder: &notificationData,
			},
			CreatedAt: time.Now(),
		}
		for i, user := range adminUsers {
			notificationPayload.RecipientUserIDs[i] = user.ID
		}

		err = c.service.NotificationService.CreateAndDeliver(ctx, notificationPayload)
		if err != nil {
			log.Printf("Failed to deliver notification for contract ID %d: %v", contract.ID, err)
			return fmt.Errorf("failed to deliver notification for contract ID %d: %v: %w", contract.ID, err, asynq.SkipRetry)
		}
		log.Printf("Notification for contract ID %d delivered successfully", contract.ID)

	}

	log.Println("All contract reminders processed successfully")
	return nil
}

func (c *AsynqServer) ProcessClientCareStatusSyncTask(ctx context.Context, t *asynq.Task) error {
	activatedClientIDs, err := c.store.ActivateDueScheduledInCareClients(ctx)
	if err != nil {
		log.Printf("Failed to activate due scheduled-in-care clients: %v", err)
		return fmt.Errorf("failed to activate due scheduled-in-care clients: %v: %w", err, asynq.SkipRetry)
	}

	if len(activatedClientIDs) == 0 {
		log.Println("No scheduled_in_care clients due for activation")
		return nil
	}

	log.Printf("Activated %d scheduled_in_care clients", len(activatedClientIDs))
	return nil
}

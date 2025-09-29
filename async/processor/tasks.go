package processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maicare_go/async/aclient"
	db "maicare_go/db/sqlc"
	"maicare_go/email"
	"maicare_go/notification"
	"maicare_go/pdf"
	"time"

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

	pdfName, err := pdf.GenerateAndUploadIncidentPDF(ctx, incidentData, processor.b2Bucket)
	if err != nil {
		log.Printf("Failed to generate and upload incident PDF: %v", err)
		return fmt.Errorf("failed to generate and upload incident PDF: %v: %w", err, asynq.SkipRetry)
	}

	err = processor.brevoConf.SendIncident(ctx, p.To, email.Incident{
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
	var payload notification.NotificationPayload
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

	// Delegate the actual work to the notification service
	err := a.notificationService.CreateAndDeliver(ctx, payload)
	if err != nil {
		// Log the error from the service
		log.Printf("Error processing notification task (ID: %s, Type: %s): %v", t.ResultWriter().TaskID(), payload.Type, err)
		// Return the error so Asynq can handle retries based on its configuration
		return fmt.Errorf("notification service failed to process task: %w", err)
	}

	log.Printf("Successfully processed notification task (ID: %s, Type: %s)", t.ResultWriter().TaskID(), payload.Type)
	return nil
}

func (c *AsynqServer) ProcessAppointmentTask(ctx context.Context, t *asynq.Task) error {
	log.Printf("Processing appointment task: %s", t.Type())
	var p aclient.AppointmentPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal appointment task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	if p.AppointmentTemplateID == (uuid.UUID{}) {
		return fmt.Errorf("invalid appointment payload: missing required fields: %w", asynq.SkipRetry)
	}

	appointemntTemplate, err := c.store.GetAppointmentTemplate(ctx, p.AppointmentTemplateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Appointment template not found: %v", err)
			return fmt.Errorf("appointment template not found: %v: %w", err, asynq.SkipRetry)
		}
		log.Printf("Failed to get appointment template: %v", err)
		return fmt.Errorf("failed to get appointment template: %v", err)
	}

	if appointemntTemplate.RecurrenceType == nil {
		return fmt.Errorf("invalid appointment template: missing recurrence type: %w", asynq.SkipRetry)
	}

	if !appointemntTemplate.StartTime.Valid || !appointemntTemplate.EndTime.Valid {
		return fmt.Errorf("invalid appointment template: missing start or end time: %w", asynq.SkipRetry)
	}

	startTime := appointemntTemplate.StartTime.Time
	endTime := appointemntTemplate.EndTime.Time
	duration := endTime.Sub(startTime)
	recurenceType := *appointemntTemplate.RecurrenceType
	interval := int(*appointemntTemplate.RecurrenceInterval)
	endDate := time.Time{}
	hasEndDate := appointemntTemplate.RecurrenceEndDate.Valid
	if hasEndDate {
		endDate = appointemntTemplate.RecurrenceEndDate.Time
	}

	safetyHorizonDate := time.Now().AddDate(2, 0, 0)

	maxOccurrences := 750

	currentStartTime := startTime
	occurrenceCount := 0

	log.Printf("Starting generation for template %d. EndDate: %v, SafetyHorizon: %v",
		appointemntTemplate.ID, appointemntTemplate.RecurrenceEndDate, safetyHorizonDate)

	for occurrenceCount < maxOccurrences {
		if hasEndDate && currentStartTime.After(endDate) {
			break
		}
		if currentStartTime.After(safetyHorizonDate) {
			break
		}

		currentEndTime := currentStartTime.Add(duration)

		scheduledAppt, err := c.store.CreateAppointment(ctx, db.CreateAppointmentParams{
			CreatorEmployeeID: &appointemntTemplate.CreatorEmployeeID,
			StartTime:         pgtype.Timestamp{Time: currentStartTime, Valid: true},
			EndTime:           pgtype.Timestamp{Time: currentEndTime, Valid: true},
			Location:          appointemntTemplate.Location,
			Description:       appointemntTemplate.Description,
			Color:             appointemntTemplate.Color,
		})
		if err != nil {
			log.Printf("Failed to create appointment: %v", err)
			return fmt.Errorf("failed to create appointment: %v: %w", err, asynq.SkipRetry)
		}
		occurrenceCount++
		log.Printf("Created appointment %d for template %d", scheduledAppt.ID, appointemntTemplate.ID)

		if len(p.ClientIDs) > 0 {
			err = c.store.BulkAddAppointmentClients(ctx, db.BulkAddAppointmentClientsParams{
				AppointmentID: scheduledAppt.ID,
				ClientIds:     p.ClientIDs,
			})
			if err != nil {
				log.Printf("Failed to add clients to appointment: %v", err)
				return fmt.Errorf("failed to add clients to appointment: %v: %w", err, asynq.SkipRetry)
			}
		}
		if len(p.ParticipantEmployeeIDs) > 0 {
			err = c.store.BulkAddAppointmentParticipants(ctx, db.BulkAddAppointmentParticipantsParams{
				AppointmentID: scheduledAppt.ID,
				EmployeeIds:   p.ParticipantEmployeeIDs,
			})
			if err != nil {
				log.Printf("Failed to add participants to appointment: %v", err)
				return fmt.Errorf("failed to add participants to appointment: %v: %w", err, asynq.SkipRetry)
			}
		}

		switch recurenceType {
		case "DAILY":
			currentStartTime = currentStartTime.AddDate(0, 0, interval)
		case "WEEKLY":
			currentStartTime = currentStartTime.AddDate(0, 0, interval*7)
		case "MONTHLY":
			currentStartTime = currentStartTime.AddDate(0, interval, 0)
		default:
			log.Printf("Unknown recurrence type: %s", recurenceType)
			return fmt.Errorf("unknown recurrence type: %s: %w", recurenceType, asynq.SkipRetry)
		}

	}
	if occurrenceCount >= maxOccurrences {
		log.Printf("Max occurrences reached: %d", maxOccurrences)
	}

	log.Printf("Finished generating appointments for template %d. Created %d occurrences", appointemntTemplate.ID, occurrenceCount)
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
			CareType:           contract.CareType,
			ContractStart:      contract.StartDate.Time,
			ContractEnd:        contract.EndDate.Time,
			ReminderType:       reminder.ReminderType,
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
			RecipientUserIDs: make([]int64, len(adminUsers)),
			Type:             notification.TypeClientContractReminder,
			Data: notification.NotificationData{
				ClientContractReminder: &notificationData,
			},
			CreatedAt: time.Now(),
		}
		for i, user := range adminUsers {
			notificationPayload.RecipientUserIDs[i] = user.ID
		}

		err = c.notificationService.CreateAndDeliver(ctx, notificationPayload)
		if err != nil {
			log.Printf("Failed to deliver notification for contract ID %d: %v", contract.ID, err)
			return fmt.Errorf("failed to deliver notification for contract ID %d: %v: %w", contract.ID, err, asynq.SkipRetry)
		}
		log.Printf("Notification for contract ID %d delivered successfully", contract.ID)

	}

	log.Println("All contract reminders processed successfully")
	return nil

}

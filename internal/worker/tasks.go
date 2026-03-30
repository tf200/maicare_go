package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	db "maicare_go/db/sqlc"
	pkgasynq "maicare_go/pkg/asynq"
	pkgemail "maicare_go/pkg/email"
	pkgpdf "maicare_go/pkg/pdf"
	"maicare_go/service/notification"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	hibikenasynq "github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

func (processor *AsynqServer) ProcessEmailTask(ctx context.Context, t *hibikenasynq.Task) error {
	var p pkgasynq.EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal email task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, hibikenasynq.SkipRetry)
	}

	if p.To == "" || p.UserEmail == "" || p.UserPassword == "" {
		return fmt.Errorf("invalid email payload: missing required fields: %w", hibikenasynq.SkipRetry)
	}

	log.Printf("Sending email to %s", p.To)

	err := processor.brevoConf.SendCredentials(ctx, []string{p.To}, pkgemail.Credentials{
		Email:    p.UserEmail,
		Password: p.UserPassword,
		Name:     p.Name,
	})
	if err != nil {
		log.Printf("Failed to send email to %s: %v", p.To, err)
		return fmt.Errorf("failed to send email to %s: %v: %w", p.To, err, hibikenasynq.SkipRetry)
	}

	return nil
}

func (processor *AsynqServer) ProcessIncidentTask(ctx context.Context, t *hibikenasynq.Task) error {
	var p pkgasynq.IncidentPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal incident task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, hibikenasynq.SkipRetry)
	}

	incidentData := pkgpdf.IncidentReportData{
		ID:                      p.ID,
		EmployeeID:              p.EmployeeID,
		EmployeeFirstName:       p.EmployeeFirstName,
		EmployeeLastName:        p.EmployeeLastName,
		LocationID:              p.LocationID,
		ReporterInvolvement:     p.ReporterInvolvement,
		InformedParties:         p.InformedParties,
		OccurredAt:              p.OccurredAt,
		IncidentType:            p.IncidentType,
		SeverityOfIncident:      p.SeverityOfIncident,
		IncidentExplanation:     p.IncidentExplanation,
		RecurrenceRisk:          p.RecurrenceRisk,
		IncidentPreventSteps:    p.IncidentPreventSteps,
		IncidentTakenMeasures:   p.IncidentTakenMeasures,
		CauseCategories:         p.CauseCategories,
		CauseExplanation:        p.CauseExplanation,
		PhysicalInjury:          p.PhysicalInjury,
		PhysicalInjuryDesc:      p.PhysicalInjuryDesc,
		PsychologicalDamage:     p.PsychologicalDamage,
		PsychologicalDamageDesc: p.PsychologicalDamageDesc,
		NeededConsultation:      p.NeededConsultation,
		FollowUpActions:         p.FollowUpActions,
		FollowUpNotes:           p.FollowUpNotes,
		IsEmployeeAbsent:        p.IsEmployeeAbsent,
		AdditionalDetails:       p.AdditionalDetails,
		ClientID:                p.ClientID,
		LocationName:            p.LocationName,
	}

	pdfName, err := processor.service.PDFService.GenerateAndUploadIncidentPDF(ctx, incidentData)
	if err != nil {
		log.Printf("Failed to generate and upload incident PDF: %v", err)
		return fmt.Errorf("failed to generate and upload incident PDF: %v: %w", err, hibikenasynq.SkipRetry)
	}

	err = processor.brevoConf.SendIncident(ctx, p.Emails, pkgemail.Incident{
		IncidentID:   p.ID.String(),
		IncidentType: p.IncidentType,
		Severity:     p.SeverityOfIncident,
		Location:     p.LocationName,
		ReportedBy:   fmt.Sprintf("%s %s", p.EmployeeFirstName, p.EmployeeLastName),
		DocumentLink: pdfName,
	})
	if err != nil {
		log.Printf("Failed to send incident email to %v: %v", p.Emails, err)
		return fmt.Errorf("failed to send incident email to %v: %v: %w", p.Emails, err, hibikenasynq.SkipRetry)
	}

	return nil
}

func (processor *AsynqServer) ProcessIncidentConfirmedEmailTask(ctx context.Context, t *hibikenasynq.Task) error {
	var p pkgasynq.IncidentConfirmedEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal incident confirmed email task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, hibikenasynq.SkipRetry)
	}

	incident, err := processor.store.GetIncident(ctx, p.IncidentID)
	if err != nil {
		log.Printf("Failed to load incident %s: %v", p.IncidentID.String(), err)
		return fmt.Errorf("failed to load incident: %w", err)
	}

	if !incident.IsConfirmed {
		log.Printf("Incident %s is not confirmed; skipping email", incident.ID.String())
		return nil
	}

	if incident.ConfirmationEmailSentAt.Valid {
		log.Printf("Incident %s confirmation email already sent at %s", incident.ID.String(), incident.ConfirmationEmailSentAt.Time.Format(time.RFC3339))
		return nil
	}

	recipientsPtr, err := processor.store.ListIncidentReportRecipientEmails(ctx, incident.ClientID)
	if err != nil {
		log.Printf("Failed to resolve incident email recipients for client %s: %v", incident.ClientID.String(), err)
		return fmt.Errorf("failed to resolve recipients: %w", err)
	}

	recipients := make([]string, 0, len(recipientsPtr))
	for _, e := range recipientsPtr {
		if e == nil || *e == "" {
			continue
		}
		recipients = append(recipients, *e)
	}
	if len(recipients) == 0 {
		log.Printf("No incident report recipients for client %s; marking as sent", incident.ClientID.String())
		_, _ = processor.store.MarkIncidentConfirmationEmailSent(ctx, incident.ID)
		return nil
	}

	informedParties := make([]string, 0, len(incident.InformedParties))
	for _, v := range incident.InformedParties {
		informedParties = append(informedParties, string(v))
	}

	causeCategories := make([]string, 0, len(incident.CauseCategories))
	for _, v := range incident.CauseCategories {
		causeCategories = append(causeCategories, string(v))
	}

	followUpActions := make([]string, 0, len(incident.FollowUpActions))
	for _, v := range incident.FollowUpActions {
		followUpActions = append(followUpActions, string(v))
	}

	incidentData := pkgpdf.IncidentReportData{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       incident.EmployeeFirstName,
		EmployeeLastName:        incident.EmployeeLastName,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     string(incident.ReporterInvolvement),
		InformedParties:         informedParties,
		OccurredAt:              incident.OccurredAt.Time,
		IncidentType:            string(incident.IncidentType),
		SeverityOfIncident:      string(incident.SeverityOfIncident),
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          string(incident.RecurrenceRisk),
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		CauseCategories:         causeCategories,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          string(incident.PhysicalInjury),
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     string(incident.PsychologicalDamage),
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      string(incident.NeededConsultation),
		FollowUpActions:         followUpActions,
		FollowUpNotes:           incident.FollowUpNotes,
		IsEmployeeAbsent:        incident.IsEmployeeAbsent,
		AdditionalDetails:       incident.AdditionalDetails,
		ClientID:                incident.ClientID,
		ClientFirstName:         incident.ClientFirstName,
		ClientLastName:          incident.ClientLastName,
		LocationName:            incident.LocationName,
	}

	pdfBytes, err := processor.service.PDFService.GenerateIncidentPDF(ctx, incidentData)
	if err != nil {
		log.Printf("Failed to generate incident PDF for %s: %v", incident.ID.String(), err)
		return fmt.Errorf("failed to generate incident pdf: %w", err)
	}

	clientName := fmt.Sprintf("%s %s", incident.ClientFirstName, incident.ClientLastName)
	reportedBy := fmt.Sprintf("%s %s", incident.EmployeeFirstName, incident.EmployeeLastName)
	attachmentName := fmt.Sprintf("incident_report_%s.pdf", incident.ID.String())

	err = processor.brevoConf.SendIncidentWithAttachment(ctx, recipients, pkgemail.Incident{
		IncidentID:   incident.ID.String(),
		ReportedBy:   reportedBy,
		ClientName:   clientName,
		IncidentType: string(incident.IncidentType),
		Severity:     string(incident.SeverityOfIncident),
		Location:     incident.LocationName,
		DocumentLink: "",
	}, attachmentName, pdfBytes)
	if err != nil {
		log.Printf("Failed to send confirmed incident email to %v: %v", recipients, err)
		return fmt.Errorf("failed to send incident email: %w", err)
	}

	_, err = processor.store.MarkIncidentConfirmationEmailSent(ctx, incident.ID)
	if err != nil {
		log.Printf("Failed to mark incident confirmation email sent for %s: %v", incident.ID.String(), err)
		return fmt.Errorf("failed to mark email sent: %w", err)
	}

	return nil
}

func (a *AsynqServer) ProcessNotificationTask(ctx context.Context, t *hibikenasynq.Task) error {
	var payload notification.NotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal notification payload: %w: %v", hibikenasynq.SkipRetry, err)
	}

	log.Printf("Received notification task: %+v", payload)

	if a.service.NotificationService == nil {
		return fmt.Errorf("notification service not initialized on AsynqServer: %w", hibikenasynq.SkipRetry)
	}

	err := a.service.NotificationService.CreateAndDeliver(ctx, payload)
	if err != nil {
		log.Printf("Error processing notification task (ID: %s, Type: %s): %v", t.ResultWriter().TaskID(), payload.Type, err)
		return fmt.Errorf("notification service failed to process task: %w", err)
	}

	log.Printf("Successfully processed notification task (ID: %s, Type: %s)", t.ResultWriter().TaskID(), payload.Type)
	return nil
}

func (processor *AsynqServer) ProcessRegistrationFormTask(ctx context.Context, t *hibikenasynq.Task) error {
	var p pkgasynq.AcceptedRegistrationFormPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal incident task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, hibikenasynq.SkipRetry)
	}

	formData := pkgemail.AcceptedRegitrationForm{
		ReferrerName:        p.ReferrerName,
		ChildName:           p.ChildName,
		ChildBSN:            p.ChildBSN,
		AppointmentDate:     p.AppointmentDate,
		AppointmentLocation: p.AppointmentLocation,
	}

	err := processor.brevoConf.SendAcceptedRegistrationForm(ctx, []string{p.To}, formData)
	if err != nil {
		log.Printf("Failed to send registration Form Email %s: %v", p.To, err)
		return fmt.Errorf("failed to send incident email to %s: %v: %w", p.To, err, hibikenasynq.SkipRetry)
	}

	return nil
}

func (processor *AsynqServer) ProcessProcessRegistrationFormEmailTask(ctx context.Context, t *hibikenasynq.Task) error {
	var p pkgasynq.ProcessRegistrationFormEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal process registration form email task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, hibikenasynq.SkipRetry)
	}

	emailData := pkgemail.ProcessRegistrationForm{
		RecipientName: p.ReferrerName,
		ClientName:    p.ClientName,
		Location:      p.Location,
		Link:          p.Link,
	}

	err := processor.brevoConf.SendProcessRegistrationForm(ctx, p.To, emailData)
	if err != nil {
		log.Printf("Failed to send process registration form email to %v: %v", p.To, err)
		return fmt.Errorf("failed to send email: %v: %w", err, hibikenasynq.SkipRetry)
	}

	return nil
}

func (c *AsynqServer) ProcessContractRemiderTask(ctx context.Context, t *hibikenasynq.Task) error {
	contractsToBeReminded, err := c.store.ListContractsTobeReminded(ctx)
	if err != nil {
		log.Printf("Failed to list contracts to be reminded: %v", err)
		return fmt.Errorf("failed to list contracts to be reminded: %v: %w", err, hibikenasynq.SkipRetry)
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
		})
		if err != nil {
			log.Printf("Failed to create contract reminder for contract ID %d: %v", contract.ID, err)
			return fmt.Errorf("failed to create contract reminder for contract ID %d: %v: %w", contract.ID, err, hibikenasynq.SkipRetry)
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
			return fmt.Errorf("failed to get admin users: %v: %w", err, hibikenasynq.SkipRetry)
		}

		if len(adminUsers) == 0 {
			log.Println("No admin users found to notify")
			return nil
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
			return fmt.Errorf("failed to deliver notification for contract ID %d: %v: %w", contract.ID, err, hibikenasynq.SkipRetry)
		}
		log.Printf("Notification for contract ID %d delivered successfully", contract.ID)
	}

	log.Println("All contract reminders processed successfully")
	return nil
}

func (c *AsynqServer) ProcessClientCareStatusSyncTask(ctx context.Context, t *hibikenasynq.Task) error {
	activatedClientIDs, err := c.store.ActivateDueScheduledInCareClients(ctx)
	if err != nil {
		log.Printf("Failed to activate due scheduled-in-care clients: %v", err)
		return fmt.Errorf("failed to activate due scheduled-in-care clients: %v: %w", err, hibikenasynq.SkipRetry)
	}

	if len(activatedClientIDs) == 0 {
		log.Println("No scheduled_in_care clients due for activation")
		return nil
	}

	log.Printf("Activated %d scheduled_in_care clients", len(activatedClientIDs))

	outOfCareClientIDs, err := c.store.ActivateDueScheduledOutOfCareClients(ctx)
	if err != nil {
		log.Printf("Failed to activate due scheduled-out-of-care clients: %v", err)
		return fmt.Errorf("failed to activate due scheduled-out-of-care clients: %v: %w", err, hibikenasynq.SkipRetry)
	}
	if len(outOfCareClientIDs) > 0 {
		log.Printf("Activated %d scheduled_out_of_care clients", len(outOfCareClientIDs))
	}

	missingFinalEvaluationClientIDs, err := c.store.ListDueScheduledOutOfCareMissingFinalEvaluation(ctx)
	if err != nil {
		log.Printf("Failed to list due scheduled-out-of-care clients missing final evaluation: %v", err)
		return fmt.Errorf("failed to list due scheduled-out-of-care clients missing final evaluation: %v: %w", err, hibikenasynq.SkipRetry)
	}
	if len(missingFinalEvaluationClientIDs) > 0 {
		log.Printf("Skipped %d scheduled_out_of_care clients due to missing final_evaluation", len(missingFinalEvaluationClientIDs))
	}

	return nil
}

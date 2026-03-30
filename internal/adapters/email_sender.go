package adapters

import (
	"context"

	"maicare_go/internal/domain"
	pkgemail "maicare_go/pkg/email"
)

type EmailSenderAdapter struct {
	sender *pkgemail.BrevoConf
}

func NewEmailSenderAdapter(sender *pkgemail.BrevoConf) domain.EmailSender {
	return &EmailSenderAdapter{sender: sender}
}

func (a *EmailSenderAdapter) SendCredentials(ctx context.Context, to []string, data domain.EmailCredentials) error {
	return a.sender.SendCredentials(ctx, to, pkgemail.Credentials{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	})
}

func (a *EmailSenderAdapter) SendIncident(ctx context.Context, to []string, data domain.IncidentEmail) error {
	return a.sender.SendIncident(ctx, to, toPkgIncident(data))
}

func (a *EmailSenderAdapter) SendIncidentWithAttachment(ctx context.Context, to []string, data domain.IncidentEmail, attachmentName string, attachmentBytes []byte) error {
	return a.sender.SendIncidentWithAttachment(ctx, to, toPkgIncident(data), attachmentName, attachmentBytes)
}

func (a *EmailSenderAdapter) SendAcceptedRegistrationForm(ctx context.Context, to []string, data domain.AcceptedRegistrationFormEmail) error {
	return a.sender.SendAcceptedRegistrationForm(ctx, to, pkgemail.AcceptedRegitrationForm{
		ReferrerName:        data.ReferrerName,
		ChildName:           data.ChildName,
		ChildBSN:            data.ChildBSN,
		AppointmentDate:     data.AppointmentDate,
		AppointmentLocation: data.AppointmentLocation,
	})
}

func (a *EmailSenderAdapter) SendProcessRegistrationForm(ctx context.Context, to []string, data domain.ProcessRegistrationFormEmail) error {
	return a.sender.SendProcessRegistrationForm(ctx, to, pkgemail.ProcessRegistrationForm{
		RecipientName: data.RecipientName,
		ClientName:    data.ClientName,
		Location:      data.Location,
		Link:          data.Link,
	})
}

func (a *EmailSenderAdapter) SendClientContractReminder(ctx context.Context, to []string, data domain.ClientContractReminderEmail) error {
	return a.sender.SendClientContractReminder(ctx, to, pkgemail.ClientContractReminder{
		ClientID:           data.ClientID,
		ClientFirstName:    data.ClientFirstName,
		ClientLastName:     data.ClientLastName,
		ContractID:         data.ContractID,
		CareType:           data.CareType,
		ContractStartDate:  data.ContractStartDate,
		ContractEndDate:    data.ContractEndDate,
		ContractStatus:     data.ContractStatus,
		ReminderType:       data.ReminderType,
		LastReminderSentAt: data.LastReminderSentAt,
		CurrentDate:        data.CurrentDate,
		CurrentYear:        data.CurrentYear,
	})
}

func toPkgIncident(data domain.IncidentEmail) pkgemail.Incident {
	return pkgemail.Incident{
		IncidentID:   data.IncidentID,
		ReportedBy:   data.ReportedBy,
		ClientName:   data.ClientName,
		IncidentType: data.IncidentType,
		Severity:     data.Severity,
		Location:     data.Location,
		DocumentLink: data.DocumentLink,
	}
}

var _ domain.EmailSender = (*EmailSenderAdapter)(nil)

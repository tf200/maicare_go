package adapters

import (
	"context"

	"maicare_go/internal/domain"
	pkgasynq "maicare_go/pkg/asynq"

	hibikenasynq "github.com/hibiken/asynq"
)

type TaskQueueAdapter struct {
	client *pkgasynq.AsynqClient
}

func NewTaskQueueAdapter(client *pkgasynq.AsynqClient) domain.TaskQueue {
	return &TaskQueueAdapter{client: client}
}

func (a *TaskQueueAdapter) EnqueueEmailDelivery(ctx context.Context, payload domain.EmailDeliveryTaskPayload, opts *domain.TaskEnqueueOptions) error {
	return a.client.EnqueueEmailDelivery(pkgasynq.EmailDeliveryPayload{
		To:           payload.To,
		Name:         payload.Name,
		UserEmail:    payload.UserEmail,
		UserPassword: payload.UserPassword,
	}, ctx, toAsynqOptions(opts)...)
}

func (a *TaskQueueAdapter) EnqueueIncident(ctx context.Context, payload domain.IncidentTaskPayload, opts *domain.TaskEnqueueOptions) error {
	return a.client.EnqueueIncident(pkgasynq.IncidentPayload{
		ID:                      payload.ID,
		EmployeeID:              payload.EmployeeID,
		EmployeeFirstName:       payload.EmployeeFirstName,
		EmployeeLastName:        payload.EmployeeLastName,
		LocationID:              payload.LocationID,
		ReporterInvolvement:     payload.ReporterInvolvement,
		InformedParties:         payload.InformedParties,
		OccurredAt:              payload.OccurredAt,
		IncidentType:            payload.IncidentType,
		SeverityOfIncident:      payload.SeverityOfIncident,
		IncidentExplanation:     payload.IncidentExplanation,
		RecurrenceRisk:          payload.RecurrenceRisk,
		IncidentPreventSteps:    payload.IncidentPreventSteps,
		IncidentTakenMeasures:   payload.IncidentTakenMeasures,
		CauseCategories:         payload.CauseCategories,
		CauseExplanation:        payload.CauseExplanation,
		PhysicalInjury:          payload.PhysicalInjury,
		PhysicalInjuryDesc:      payload.PhysicalInjuryDesc,
		PsychologicalDamage:     payload.PsychologicalDamage,
		PsychologicalDamageDesc: payload.PsychologicalDamageDesc,
		NeededConsultation:      payload.NeededConsultation,
		FollowUpActions:         payload.FollowUpActions,
		FollowUpNotes:           payload.FollowUpNotes,
		IsEmployeeAbsent:        payload.IsEmployeeAbsent,
		AdditionalDetails:       payload.AdditionalDetails,
		ClientID:                payload.ClientID,
		LocationName:            payload.LocationName,
		Emails:                  payload.Emails,
	}, ctx, toAsynqOptions(opts)...)
}

func (a *TaskQueueAdapter) EnqueueIncidentConfirmedEmail(ctx context.Context, payload domain.IncidentConfirmedEmailTaskPayload, opts *domain.TaskEnqueueOptions) error {
	return a.client.EnqueueIncidentConfirmedEmail(ctx, pkgasynq.IncidentConfirmedEmailPayload{
		IncidentID: payload.IncidentID,
	}, toAsynqOptions(opts)...)
}

func (a *TaskQueueAdapter) EnqueueNotificationTask(ctx context.Context, payload domain.NotificationTaskPayload, opts *domain.TaskEnqueueOptions) error {
	return a.client.EnqueueNotificationTask(ctx, pkgasynq.NotificationPayload{
		RecipientUserIDs: payload.RecipientUserIDs,
		Type:             payload.Type,
		Data:             toAsynqNotificationData(payload.Data),
		CreatedAt:        payload.CreatedAt,
		Message:          payload.Message,
	}, toAsynqOptions(opts)...)
}

func (a *TaskQueueAdapter) EnqueueAcceptedRegistration(ctx context.Context, payload domain.AcceptedRegistrationFormTaskPayload, opts *domain.TaskEnqueueOptions) error {
	return a.client.EnqueueAcceptedRegistration(ctx, pkgasynq.AcceptedRegistrationFormPayload{
		ReferrerName:        payload.ReferrerName,
		ChildName:           payload.ChildName,
		ChildBSN:            payload.ChildBSN,
		AppointmentDate:     payload.AppointmentDate,
		AppointmentLocation: payload.AppointmentLocation,
		To:                  payload.To,
	}, toAsynqOptions(opts)...)
}

func (a *TaskQueueAdapter) EnqueueProcessRegistrationFormEmail(ctx context.Context, payload domain.ProcessRegistrationFormEmailTaskPayload, opts *domain.TaskEnqueueOptions) error {
	return a.client.EnqueueProcessRegistrationFormEmail(ctx, pkgasynq.ProcessRegistrationFormEmailPayload{
		ReferrerName: payload.ReferrerName,
		ClientName:   payload.ClientName,
		Location:     payload.Location,
		Link:         payload.Link,
		To:           payload.To,
	}, toAsynqOptions(opts)...)
}

func (a *TaskQueueAdapter) Close() error {
	return a.client.Close()
}

func toAsynqOptions(opts *domain.TaskEnqueueOptions) []hibikenasynq.Option {
	if opts == nil {
		return nil
	}

	result := make([]hibikenasynq.Option, 0, 3)
	if opts.Queue != "" {
		result = append(result, hibikenasynq.Queue(opts.Queue))
	}
	if opts.MaxRetry > 0 {
		result = append(result, hibikenasynq.MaxRetry(opts.MaxRetry))
	}
	if opts.ProcessAt != nil && !opts.ProcessAt.IsZero() {
		result = append(result, hibikenasynq.ProcessAt(*opts.ProcessAt))
	}

	return result
}

func toAsynqNotificationData(data domain.NotificationTaskData) pkgasynq.NotificationData {
	return pkgasynq.NotificationData{
		NewAppointment:          toAsynqNewAppointmentData(data.NewAppointment),
		NewClientAssignment:     toAsynqNewClientAssignmentData(data.NewClientAssignment),
		ClientContractReminder:  toAsynqClientContractReminderData(data.ClientContractReminder),
		NewIncidentReport:       toAsynqNewIncidentReportData(data.NewIncidentReport),
		NewScheduleNotification: toAsynqNewScheduleNotificationData(data.NewScheduleNotification),
	}
}

func toAsynqNewAppointmentData(data *domain.NewAppointmentTaskData) *pkgasynq.NewAppointmentData {
	if data == nil {
		return nil
	}

	return &pkgasynq.NewAppointmentData{
		AppointmentID: data.AppointmentID,
		CreatedBy:     data.CreatedBy,
		StartTime:     data.StartTime,
		EndTime:       data.EndTime,
		Location:      data.Location,
	}
}

func toAsynqNewClientAssignmentData(data *domain.NewClientAssignmentTaskData) *pkgasynq.NewClientAssignmentData {
	if data == nil {
		return nil
	}

	return &pkgasynq.NewClientAssignmentData{
		ClientID:        data.ClientID,
		ClientFirstName: data.ClientFirstName,
		ClientLastName:  data.ClientLastName,
		ClientLocation:  data.ClientLocation,
	}
}

func toAsynqClientContractReminderData(data *domain.ClientContractReminderTaskData) *pkgasynq.ClientContractReminderData {
	if data == nil {
		return nil
	}

	return &pkgasynq.ClientContractReminderData{
		ClientID:           data.ClientID,
		ClientFirstName:    data.ClientFirstName,
		ClientLastName:     data.ClientLastName,
		ContractID:         data.ContractID,
		CareType:           data.CareType,
		ContractStart:      data.ContractStart,
		ContractEnd:        data.ContractEnd,
		ReminderType:       data.ReminderType,
		LastReminderSentAt: data.LastReminderSentAt,
	}
}

func toAsynqNewIncidentReportData(data *domain.NewIncidentReportTaskData) *pkgasynq.NewIncidentReportData {
	if data == nil {
		return nil
	}

	return &pkgasynq.NewIncidentReportData{
		ID:                 data.ID,
		EmployeeID:         data.EmployeeID,
		EmployeeFirstName:  data.EmployeeFirstName,
		EmployeeLastName:   data.EmployeeLastName,
		LocationID:         data.LocationID,
		LocationName:       data.LocationName,
		ClientID:           data.ClientID,
		ClientFirstName:    data.ClientFirstName,
		ClientLastName:     data.ClientLastName,
		SeverityOfIncident: data.SeverityOfIncident,
	}
}

func toAsynqNewScheduleNotificationData(data *domain.NewScheduleNotificationTaskData) *pkgasynq.NewScheduleNotificationData {
	if data == nil {
		return nil
	}

	return &pkgasynq.NewScheduleNotificationData{
		ScheduleID: data.ScheduleID,
		CreatedBy:  data.CreatedBy,
		StartTime:  data.StartTime,
		EndTime:    data.EndTime,
		Location:   data.Location,
	}
}

var _ domain.TaskQueue = (*TaskQueueAdapter)(nil)

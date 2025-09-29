package appointment

import (
	"context"
	"database/sql"
	"fmt"
	"maicare_go/async/aclient"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/notification"
	"maicare_go/util"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *appointmentService) CreateAppointment(req *CreateAppointmentRequest, userID int64, ctx context.Context) (*CreateAppointmentResponse, error) {
	if req.StartTime.After(req.EndTime) {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Start time is after end time")
		return nil, fmt.Errorf("start time must be before end time")
	}

	employee, err := s.Store.GetEmployeeProfileByUserID(ctx, userID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to get employee profile", zap.Error(err))
		return nil, fmt.Errorf("failed to get employee profile")
	}

	filteredParticipants := req.ParticipantEmployeeIDs[:0] // Reuse the slice's backing array
	for _, id := range req.ParticipantEmployeeIDs {
		if id != employee.EmployeeID {
			filteredParticipants = append(filteredParticipants, id)
		}
	}
	req.ParticipantEmployeeIDs = filteredParticipants
	if req.RecurrenceType == "NONE" {
		return s.createNormalAppointment(req, employee.EmployeeID, employee.FirstName, employee.LastName, ctx)
	} else {
		return s.createRecurringAppointment(req, employee.EmployeeID, employee.FirstName, employee.LastName, ctx)
	}

}

func (s *appointmentService) AddParticipantToAppointment(
	ctx context.Context,
	appointmentID uuid.UUID,
	req AddParticipantToAppointmentRequest) error {
	err := s.Store.BulkAddAppointmentParticipants(ctx, db.BulkAddAppointmentParticipantsParams{
		AppointmentID: appointmentID,
		EmployeeIds:   req.ParticipantEmployeeIDs,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddParticipantToAppointmentApi", "Failed to add appointment participants", zap.Error(err))
		return fmt.Errorf("failed to add participants to appointment")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "AddParticipantToAppointmentApi", "Participants added to appointment successfully", zap.String("appointment_id", appointmentID.String()))
	return nil
}

func (s *appointmentService) AddClientToAppointment(
	ctx context.Context,
	appointmentID uuid.UUID,
	req AddClientToAppointmentRequest) error {
	err := s.Store.BulkAddAppointmentClients(ctx, db.BulkAddAppointmentClientsParams{
		AppointmentID: appointmentID,
		ClientIds:     req.ClientIDs,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddClientToAppointmentApi", "Failed to add appointment clients", zap.Error(err))
		return fmt.Errorf("failed to add clients to appointment")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "AddClientToAppointmentApi", "Clients added to appointment successfully", zap.String("appointment_id", appointmentID.String()))
	return nil
}

func (s *appointmentService) ListAppointmentsForEmployeeInRange(
	ctx context.Context,
	employeeID int64,
	req ListAppointmentsForEmployeeInRangeRequest) ([]ListAppointmentsForEmployeeInRangeResponse, error) {
	if req.StartDate.After(req.EndDate) {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAppointmentsForEmployeeInRangeApi", "Start date is after end date")
		return nil, fmt.Errorf("start date must be before end date")
	}

	appointments, err := s.Store.ListEmployeeAppointmentsInRange(ctx, db.ListEmployeeAppointmentsInRangeParams{
		EmployeeID: &employeeID,
		StartDate:  pgtype.Timestamp{Time: req.StartDate, Valid: true},
		EndDate:    pgtype.Timestamp{Time: req.EndDate, Valid: true},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAppointmentsForEmployeeInRangeApi", "Failed to list appointments", zap.Error(err))
		return nil, fmt.Errorf("failed to list appointments")
	}

	if len(appointments) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListAppointmentsForEmployeeInRangeApi", "No appointments found")
		return []ListAppointmentsForEmployeeInRangeResponse{}, nil
	}

	var appointmentIDs []uuid.UUID
	for _, appt := range appointments {
		appointmentIDs = append(appointmentIDs, appt.AppointmentID)
	}

	participantsMap := make(map[uuid.UUID][]ParticipantsDetails)
	participants, err := s.Store.GetAppointmentParticipants(ctx, appointmentIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAppointmentsForEmployeeInRangeApi", "Failed to get appointment participants", zap.Error(err))
		return nil, fmt.Errorf("failed to list appointments")
	}
	for _, p := range participants {
		participantsMap[p.AppointmentID] = append(participantsMap[p.AppointmentID], ParticipantsDetails{
			EmployeeID: p.EmployeeID,
			FirstName:  p.FirstName,
			LastName:   p.LastName,
		})
	}

	clientsMap := make(map[uuid.UUID][]ClientsDetails)
	clients, err := s.Store.GetAppointmentClients(ctx, appointmentIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAppointmentsForEmployeeInRangeApi", "Failed to get appointment clients", zap.Error(err))
		return nil, fmt.Errorf("failed to list appointments")
	}
	for _, c := range clients {
		clientsMap[c.AppointmentID] = append(clientsMap[c.AppointmentID], ClientsDetails{
			ClientID:  c.ClientID,
			FirstName: c.FirstName,
			LastName:  c.LastName,
		})
	}

	var resp []ListAppointmentsForEmployeeInRangeResponse
	for _, appt := range appointments {
		resp = append(resp, ListAppointmentsForEmployeeInRangeResponse{
			ID:                  appt.AppointmentID,
			CreatorEmployeeID:   appt.CreatorEmployeeID,
			StartTime:           appt.StartTime.Time,
			EndTime:             appt.EndTime.Time,
			Color:               appt.Color,
			Location:            appt.Location,
			Description:         appt.Description,
			ParticipantsDetails: participantsMap[appt.AppointmentID],
			ClientsDetails:      clientsMap[appt.AppointmentID],
		})
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListAppointmentsForEmployeeInRangeApi", "Appointments listed successfully", zap.Int("count", len(resp)))
	return resp, nil
}

func (s *appointmentService) ListAppointmentsForClientInRange(
	ctx context.Context,
	clientID int64,
	req ListAppointmentsForClientRequest) ([]ListAppointmentsForClientResponse, error) {
	if req.StartDate.After(req.EndDate) {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAppointmentsForClientInRangeApi", "Start date is after end date")
		return nil, fmt.Errorf("start date must be before end date")
	}

	appointments, err := s.Store.ListClientAppointmentsInRange(ctx, db.ListClientAppointmentsInRangeParams{
		ClientID:  clientID,
		StartDate: pgtype.Timestamp{Time: req.StartDate, Valid: true},
		EndDate:   pgtype.Timestamp{Time: req.EndDate, Valid: true},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAppointmentsForClientInRangeApi", "Failed to list appointments", zap.Error(err))
		return nil, fmt.Errorf("failed to list appointments")
	}

	if len(appointments) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListAppointmentsForClientInRangeApi", "No appointments found")
		return []ListAppointmentsForClientResponse{}, nil
	}

	var appointmentIDs []uuid.UUID
	for _, appt := range appointments {
		appointmentIDs = append(appointmentIDs, appt.AppointmentID)
	}

	participantsMap := make(map[uuid.UUID][]ParticipantsDetails)
	participants, err := s.Store.GetAppointmentParticipants(ctx, appointmentIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAppointmentsForClientInRangeApi", "Failed to get appointment participants", zap.Error(err))
		return nil, fmt.Errorf("failed to list appointments")
	}
	for _, p := range participants {
		participantsMap[p.AppointmentID] = append(participantsMap[p.AppointmentID], ParticipantsDetails{
			EmployeeID: p.EmployeeID,
			FirstName:  p.FirstName,
			LastName:   p.LastName,
		})
	}

	clientsMap := make(map[uuid.UUID][]ClientsDetails)
	clients, err := s.Store.GetAppointmentClients(ctx, appointmentIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAppointmentsForClientInRangeApi", "Failed to get appointment clients", zap.Error(err))
		return nil, fmt.Errorf("failed to list appointments")
	}
	for _, c := range clients {
		clientsMap[c.AppointmentID] = append(clientsMap[c.AppointmentID], ClientsDetails{
			ClientID:  c.ClientID,
			FirstName: c.FirstName,
			LastName:  c.LastName,
		})
	}

	var resp []ListAppointmentsForClientResponse
	for _, appt := range appointments {
		resp = append(resp, ListAppointmentsForClientResponse{
			ID:                  appt.AppointmentID,
			CreatorEmployeeID:   appt.CreatorEmployeeID,
			StartTime:           appt.StartTime.Time,
			EndTime:             appt.EndTime.Time,
			Color:               appt.Color,
			Location:            appt.Location,
			Description:         appt.Description,
			ParticipantsDetails: participantsMap[appt.AppointmentID],
			ClientsDetails:      clientsMap[appt.AppointmentID],
		})
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListAppointmentsForClientInRangeApi", "Appointments listed successfully", zap.Int("count", len(resp)))
	return resp, nil
}

func (s *appointmentService) GetAppointment(
	ctx context.Context,
	appointmentID uuid.UUID) (*GetAppointmentResponse, error) {
	appointment, err := s.Store.GetScheduledAppointmentByID(ctx, appointmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "GetAppointmentApi", "Appointment not found", zap.String("appointment_id", appointmentID.String()))
			return nil, fmt.Errorf("appointment not found")
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetAppointmentApi", "Failed to get appointment", zap.Error(err))
		return nil, fmt.Errorf("failed to get appointment")
	}

	participants, err := s.Store.GetAppointmentParticipants(ctx, []uuid.UUID{appointmentID})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetAppointmentApi", "Failed to get appointment participants", zap.Error(err))
		return nil, fmt.Errorf("failed to get appointment")
	}
	var participantDetails []ParticipantsDetails
	for _, p := range participants {
		participantDetails = append(participantDetails, ParticipantsDetails{
			EmployeeID: p.EmployeeID,
			FirstName:  p.FirstName,
			LastName:   p.LastName,
		})
	}

	clients, err := s.Store.GetAppointmentClients(ctx, []uuid.UUID{appointmentID})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetAppointmentApi", "Failed to get appointment clients", zap.Error(err))
		return nil, fmt.Errorf("failed to get appointment")
	}
	var clientDetails []ClientsDetails
	for _, c := range clients {
		clientDetails = append(clientDetails, ClientsDetails{
			ClientID:  c.ClientID,
			FirstName: c.FirstName,
			LastName:  c.LastName,
		})
	}

	resp := &GetAppointmentResponse{
		ID:                  appointment.ID,
		CreatorEmployeeID:   appointment.CreatorEmployeeID,
		StartTime:           appointment.StartTime.Time,
		EndTime:             appointment.EndTime.Time,
		Color:               appointment.Color,
		Location:            appointment.Location,
		Description:         appointment.Description,
		ParticipantsDetails: participantDetails,
		ClientsDetails:      clientDetails,
	}

	return resp, nil
}

func (s *appointmentService) UpdateAppointment(
	ctx context.Context,
	appointmentID uuid.UUID,
	req *UpdateAppointmentRequest) (*UpdateAppointmentResponse, error) {

	if req.StartTime.After(req.EndTime) {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentApi", "Start time is after end time")
		return nil, fmt.Errorf("start time must be before end time")
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentApi", "Failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to update appointment")
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err != sql.ErrTxDone {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentApi", "Failed to rollback transaction", zap.Error(err))
		}
	}()
	qtx := s.Store.WithTx(tx)

	appointment, err := qtx.UpdateAppointment(ctx, db.UpdateAppointmentParams{
		ID:          appointmentID,
		StartTime:   pgtype.Timestamp{Time: req.StartTime, Valid: true},
		EndTime:     pgtype.Timestamp{Time: req.EndTime, Valid: true},
		Location:    req.Location,
		Color:       req.Color,
		Description: req.Description,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentApi", "Failed to update appointment", zap.Error(err))
		return nil, fmt.Errorf("failed to update appointment")
	}

	if req.ParticipantEmployeeIDs != nil {
		err = qtx.DeleteAppointmentParticipants(ctx, appointment.ID)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentApi", "Failed to delete appointment participants", zap.Error(err))
			return nil, fmt.Errorf("failed to update appointment")
		}
		if len(*req.ParticipantEmployeeIDs) > 0 {
			err = qtx.BulkAddAppointmentParticipants(ctx, db.BulkAddAppointmentParticipantsParams{
				AppointmentID: appointment.ID,
				EmployeeIds:   *req.ParticipantEmployeeIDs,
			})
			if err != nil {
				s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentApi", "Failed to add appointment participants", zap.Error(err))
				return nil, fmt.Errorf("failed to update appointment")
			}
		}
	}

	if req.ClientIDs != nil {
		err = qtx.DeleteAppointmentClients(ctx, appointment.ID)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentApi", "Failed to delete appointment clients", zap.Error(err))
			return nil, fmt.Errorf("failed to update appointment")
		}
		if len(*req.ClientIDs) > 0 {
			err = qtx.BulkAddAppointmentClients(ctx, db.BulkAddAppointmentClientsParams{
				AppointmentID: appointment.ID,
				ClientIds:     *req.ClientIDs,
			})
			if err != nil {
				s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentApi", "Failed to add appointment clients", zap.Error(err))
				return nil, fmt.Errorf("failed to update appointment")
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentApi", "Failed to commit transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to update appointment")
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateAppointmentApi", "Appointment updated successfully", zap.String("appointment_id", appointment.ID.String()))

	return &UpdateAppointmentResponse{
		ID:                     appointment.ID,
		AppointmentTemplatesID: appointment.AppointmentTemplatesID,
		CreatorEmployeeID:      appointment.CreatorEmployeeID,
		StartTime:              appointment.StartTime,
		EndTime:                appointment.EndTime,
		Location:               appointment.Location,
		Description:            appointment.Description,
		Color:                  appointment.Color,
		Status:                 appointment.Status,
		IsConfirmed:            appointment.IsConfirmed,
		ConfirmedByEmployeeID:  appointment.ConfirmedByEmployeeID,
		ConfirmedAt:            appointment.ConfirmedAt,
		CreatedAt:              appointment.CreatedAt,
		UpdatedAt:              appointment.UpdatedAt,
	}, nil

}

func (s *appointmentService) DeleteAppointment(
	ctx context.Context,
	appointmentID uuid.UUID) error {
	err := s.Store.DeleteAppointment(ctx, appointmentID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteAppointmentApi", "Failed to delete appointment", zap.Error(err))
		return fmt.Errorf("failed to delete appointment")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "DeleteAppointmentApi", "Appointment deleted successfully", zap.String("appointment_id", appointmentID.String()))
	return nil
}

func (s *appointmentService) ConfirmAppointment(
	ctx context.Context,
	appointmentID uuid.UUID,
	employeeID int64) error {
	err := s.Store.ConfirmAppointment(ctx, db.ConfirmAppointmentParams{
		ID:         appointmentID,
		EmployeeID: &employeeID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ConfirmAppointmentApi", "Failed to confirm appointment", zap.Error(err))
		return fmt.Errorf("failed to confirm appointment")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ConfirmAppointmentApi", "Appointment confirmed successfully", zap.String("appointment_id", appointmentID.String()))
	return nil
}

func (s *appointmentService) createNormalAppointment(
	req *CreateAppointmentRequest,
	employeeID int64,
	employeeFirstName string,
	employeeLastName string,
	ctx context.Context) (*CreateAppointmentResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to create appointment")
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err != sql.ErrTxDone {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to rollback transaction", zap.Error(err))
		}
	}()

	qtx := s.Store.WithTx(tx)

	appointment, err := qtx.CreateAppointment(ctx, db.CreateAppointmentParams{
		CreatorEmployeeID: &employeeID,
		StartTime:         pgtype.Timestamp{Time: req.StartTime, Valid: true},
		EndTime:           pgtype.Timestamp{Time: req.EndTime, Valid: true},
		Location:          req.Location,
		Description:       req.Description,
		Color:             req.Color,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to create appointment", zap.Error(err))
		return nil, fmt.Errorf("failed to create appointment")
	}

	if len(req.ParticipantEmployeeIDs) > 0 {
		err = qtx.BulkAddAppointmentParticipants(ctx, db.BulkAddAppointmentParticipantsParams{
			AppointmentID: appointment.ID,
			EmployeeIds:   req.ParticipantEmployeeIDs,
		})
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to add appointment participants", zap.Error(err))
			return nil, fmt.Errorf("failed to create appointment")
		}
	}

	if len(req.ClientIDs) > 0 {
		err = qtx.BulkAddAppointmentClients(ctx, db.BulkAddAppointmentClientsParams{
			AppointmentID: appointment.ID,
			ClientIds:     req.ClientIDs,
		})
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to add appointment clients", zap.Error(err))
			return nil, fmt.Errorf("failed to create appointment")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to commit transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to create appointment")
	}

	if len(req.ParticipantEmployeeIDs) > 0 {

		data := notification.NewAppointmentData{
			AppointmentID: appointment.ID,
			CreatedBy:     employeeFirstName + " " + employeeLastName,
			StartTime:     appointment.StartTime.Time,
			EndTime:       appointment.EndTime.Time,
			Location:      util.DerefString(req.Location),
		}

		err = s.AsynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
			RecipientUserIDs: req.ParticipantEmployeeIDs,
			Type:             notification.TypeNewAppointment,
			Data: notification.NotificationData{
				NewAppointment: &data,
			},
			Message:   data.NewAppointmentMessage(),
			CreatedAt: time.Now(),
		})
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to enqueue notification task", zap.Error(err))
		}
	}

	resp := &CreateAppointmentResponse{
		ID:                appointment.ID,
		CreatorEmployeeID: appointment.CreatorEmployeeID,
		StartTime:         appointment.StartTime.Time,
		EndTime:           appointment.EndTime.Time,
		Color:             appointment.Color,
		Location:          appointment.Location,
		Description:       appointment.Description,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "CreateAppointmentApi", "Appointment created successfully", zap.String("appointment_id", appointment.ID.String()))
	return resp, nil
}

func (s *appointmentService) createRecurringAppointment(
	req *CreateAppointmentRequest,
	employeeID int64,
	employeeFirstName string,
	employeeLastName string,
	ctx context.Context) (*CreateAppointmentResponse, error) {

	appointmentTemp, err := s.Store.CreateAppointmentTemplate(ctx, db.CreateAppointmentTemplateParams{
		CreatorEmployeeID:  employeeID,
		StartTime:          pgtype.Timestamp{Time: req.StartTime, Valid: true},
		EndTime:            pgtype.Timestamp{Time: req.EndTime, Valid: true},
		Location:           req.Location,
		Description:        req.Description,
		Color:              req.Color,
		RecurrenceType:     &req.RecurrenceType,
		RecurrenceInterval: req.RecurrenceInterval,
		RecurrenceEndDate:  pgtype.Date{Time: req.RecurrenceEndDate, Valid: true},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to create appointment template", zap.Error(err))
		return nil, fmt.Errorf("failed to create appointment")
	}

	err = s.AsynqClient.EnqueueAppointmentTask(ctx, aclient.AppointmentPayload{
		AppointmentTemplateID:  appointmentTemp.ID,
		ParticipantEmployeeIDs: req.ParticipantEmployeeIDs,
		ClientIDs:              req.ClientIDs,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to enqueue appointment creation task", zap.Error(err))
		return nil, fmt.Errorf("failed to create appointment")
	}

	if len(req.ParticipantEmployeeIDs) > 0 {

		data := notification.NewAppointmentData{
			AppointmentID: appointmentTemp.ID,
			CreatedBy:     employeeFirstName + " " + employeeLastName,
			StartTime:     appointmentTemp.StartTime.Time,
			EndTime:       appointmentTemp.EndTime.Time,
			Location:      util.DerefString(req.Location),
		}

		err = s.AsynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
			RecipientUserIDs: req.ParticipantEmployeeIDs,
			Type:             notification.TypeNewAppointment,
			Data: notification.NotificationData{
				NewAppointment: &data,
			},
			Message:   data.NewAppointmentMessage(),
			CreatedAt: time.Now(),
		})
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to enqueue notification task", zap.Error(err))
		}
	}

	resp := &CreateAppointmentResponse{
		ID:                appointmentTemp.ID,
		CreatorEmployeeID: &appointmentTemp.CreatorEmployeeID,
		StartTime:         appointmentTemp.StartTime.Time,
		EndTime:           appointmentTemp.EndTime.Time,
		Color:             appointmentTemp.Color,
		Location:          appointmentTemp.Location,
		Description:       appointmentTemp.Description,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "CreateAppointmentApi", "Recurring appointment template created successfully", zap.String("appointment_template_id", appointmentTemp.ID.String()))
	return resp, nil
}

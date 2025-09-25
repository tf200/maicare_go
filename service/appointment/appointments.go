package appointment

import (
	"context"
	"database/sql"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *appointmentService) CreateAppointment(req *CreateAppointmentRequest, userID int64, ctx context.Context) (*CreateAppointmentResponse, error) {
	employee, err := s.Store.GetEmployeeProfileByUserID(ctx, userID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Failed to get employee profile", zap.Error(err))
		return nil, fmt.Errorf("failed to get employee profile")
	}

	if req.StartTime.After(req.EndTime) {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Start time is after end time")
		return nil, fmt.Errorf("start time must be before end time")
	}

	filteredParticipants := req.ParticipantEmployeeIDs[:0] // Reuse the slice's backing array
	for _, id := range req.ParticipantEmployeeIDs {
		if id != employee.EmployeeID {
			filteredParticipants = append(filteredParticipants, id)
		}
	}
	req.ParticipantEmployeeIDs = filteredParticipants

}

func (s *appointmentService) CreateNormalAppointment(req *CreateAppointmentRequest, employeeID int64, ctx context.Context) (*CreateAppointmentResponse, error) {
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

func (s *appointmentService) CreateRecurringAppointment(req *CreateAppointmentRequest, employeeID int64, ctx context.Context) (*CreateAppointmentResponse, error) {

	appointmentTemp, err := server.store.CreateAppointmentTemplate(ctx, db.CreateAppointmentTemplateParams{
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

	err = 

	s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentApi", "Recurring appointments are not yet implemented")
	return nil, fmt.Errorf("recurring appointments are not yet implemented")
}

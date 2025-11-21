package employees

import (
	"context"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/util"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *employeeService) ListWorkingHours(
	ctx context.Context,
	employeeID uuid.UUID,
	req *ListWorkingHoursRequest,
) (*ListWorkingHoursResponse, error) {
	// Get period start and end dates
	periodStart, periodEnd, err := util.GetStartAndEndOfISOWeek(int(req.Year), int(req.Week))
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"ListWorkingHours",
			"Failed to get ISO week dates",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get ISO week dates: %w", err)
	}

	// Fetch appointments and schedules
	appointments, schedules, err := s.fetchWorkingHoursData(ctx, employeeID, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	// Build working hours items
	workingHours, summary := s.buildWorkingHoursItems(appointments, schedules)

	// Calculate overtime
	err = s.calculateOvertime(ctx, employeeID, &summary)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"ListWorkingHours",
			"Failed to calculate overtime",
			zap.Error(err),
		)
		// Don't fail the entire request, just log the error
	}

	// Build period information
	period := s.buildPeriodInfo(req.Year, req.Week, periodStart, periodEnd)

	return &ListWorkingHoursResponse{
		EmployeeID:   employeeID,
		Period:       period,
		Summary:      summary,
		WorkingHours: workingHours,
	}, nil
}

func (s *employeeService) fetchWorkingHoursData(
	ctx context.Context,
	employeeID uuid.UUID,
	periodStart, periodEnd time.Time,
) ([]db.ListEmployeeAppointmentsInRangeRow, []db.GetEmployeeSchedulesRow, error) {
	// Fetch employee appointments
	appointments, err := s.Store.ListEmployeeAppointmentsInRange(
		ctx,
		db.ListEmployeeAppointmentsInRangeParams{
			StartDate: pgtype.Timestamp{
				Time:  periodStart,
				Valid: true,
			},
			EndDate: pgtype.Timestamp{
				Time:  periodEnd,
				Valid: true,
			},
			EmployeeID: &employeeID,
		},
	)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"fetchWorkingHoursData",
			"Failed to fetch employee appointments",
			zap.String("employee_id", employeeID.String()),
			zap.Error(err),
		)
		return nil, nil, fmt.Errorf("failed to fetch employee appointments: %w", err)
	}

	// Fetch employee schedules
	schedules, err := s.Store.GetEmployeeSchedules(ctx, db.GetEmployeeSchedulesParams{
		PeriodStart: pgtype.Timestamp{
			Time:  periodStart,
			Valid: true,
		},
		PeriodEnd: pgtype.Timestamp{
			Time:  periodEnd,
			Valid: true,
		},
		EmployeeID: employeeID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"fetchWorkingHoursData",
			"Failed to fetch employee schedules",
			zap.String("employee_id", employeeID.String()),
			zap.Error(err),
		)
		return nil, nil, fmt.Errorf("failed to fetch employee schedules: %w", err)
	}

	return appointments, schedules, nil
}

func (s *employeeService) buildWorkingHoursItems(
	appointments []db.ListEmployeeAppointmentsInRangeRow,
	schedules []db.GetEmployeeSchedulesRow,
) ([]WorkingHourItem, Summary) {
	workingHours := make([]WorkingHourItem, len(schedules)+len(appointments))
	var appointmentHours, shiftHours float64
	uniqueDays := make(map[string]bool)

	// Process schedules
	for i, schedule := range schedules {
		duration := schedule.EndDatetime.Time.Sub(schedule.StartDatetime.Time).Hours()
		shiftHours += duration
		dayKey := schedule.StartDatetime.Time.Format("2006-01-02")
		uniqueDays[dayKey] = true

		workingHours[i] = WorkingHourItem{
			ID:            schedule.ID,
			Type:          "schedule",
			StartTime:     schedule.StartDatetime.Time,
			EndTime:       schedule.EndDatetime.Time,
			DurationHours: duration,
			Location:      schedule.LocationName,
			LocationID:    &schedule.LocationID,
			Description:   nil,
			Status:        nil,
		}
	}

	// Process appointments
	for i, appointment := range appointments {
		duration := appointment.EndTime.Time.Sub(appointment.StartTime.Time).Hours()
		appointmentHours += duration
		dayKey := appointment.StartTime.Time.Format("2006-01-02")
		uniqueDays[dayKey] = true

		workingHours[len(schedules)+i] = WorkingHourItem{
			ID:            appointment.AppointmentID,
			Type:          "appointment",
			StartTime:     appointment.StartTime.Time,
			EndTime:       appointment.EndTime.Time,
			DurationHours: duration,
			Location:      *appointment.Location,
			LocationID:    nil,
			Description:   appointment.Description,
			Status:        util.StringPtr(string(appointment.Status)),
		}
	}

	summary := Summary{
		TotalHours:       appointmentHours + shiftHours,
		AppointmentHours: appointmentHours,
		ShiftHours:       shiftHours,
		TotalDaysWorked:  len(uniqueDays),
		OverTime:         0.0,
	}

	return workingHours, summary
}

func (s *employeeService) calculateOvertime(
	ctx context.Context,
	employeeID uuid.UUID,
	summary *Summary,
) error {
	// Fetch employee contract details
	contractDetails, err := s.Store.GetEmployeeContractDetails(ctx, employeeID)
	if err != nil {
		return fmt.Errorf("failed to fetch employee contract details: %w", err)
	}

	// Calculate overtime if applicable
	if contractDetails.ContractHours != nil {
		standardHours := *contractDetails.ContractHours
		if standardHours > 0 && summary.TotalHours > standardHours {
			summary.OverTime = summary.TotalHours - standardHours
		}
	}

	return nil
}

func (s *employeeService) buildPeriodInfo(
	year, week int32,
	periodStart, periodEnd time.Time,
) Period {
	currentYear, currentWeek := time.Now().ISOWeek()

	return Period{
		Year:          year,
		Week:          week,
		MonthName:     periodStart.Month().String(),
		IsCurrentWeek: year == int32(currentYear) && week == int32(currentWeek),
		DateRange: struct {
			Start string `json:"start"`
			End   string `json:"end"`
		}{
			Start: periodStart.Format("2006-01-02"),
			End:   periodEnd.Format("2006-01-02"),
		},
	}
}

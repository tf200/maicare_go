package employees

import (
	"context"
	"fmt"
	"sort"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/appointment"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *employeeService) GetMyScheduleTimeline(
	ctx context.Context,
	employeeID uuid.UUID,
	req *GetMyScheduleTimelineRequest,
) ([]GetMyScheduleTimelineDayResponse, error) {
	startDate, endDate, err := parseTimelineDateRange(req)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetMyScheduleTimeline",
			"Invalid date range",
			zap.String("start_date", req.StartDate),
			zap.String("end_date", req.EndDate),
			zap.Error(err),
		)
		return nil, err
	}

	endExclusive := endDate.AddDate(0, 0, 1)

	schedules, err := s.Store.GetEmployeeSchedules(ctx, db.GetEmployeeSchedulesParams{
		EmployeeID: employeeID,
		PeriodStart: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		PeriodEnd: pgtype.Timestamptz{
			Time:  endExclusive,
			Valid: true,
		},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetMyScheduleTimeline",
			"Failed to fetch employee schedules",
			zap.String("employee_id", employeeID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch employee schedules: %w", err)
	}

	events, err := s.appointmentService.ListEvents(ctx, appointment.ListEventsRequest{
		StartAt: startDate,
		EndAt:   endExclusive,
	}, employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetMyScheduleTimeline",
			"Failed to fetch employee events",
			zap.String("employee_id", employeeID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch employee events: %w", err)
	}

	dayMap := make(map[string][]GetMyScheduleTimelineItem)

	for _, sched := range schedules {
		dayKey := sched.StartDatetime.Time.UTC().Format("2006-01-02")
		dayMap[dayKey] = append(dayMap[dayKey], GetMyScheduleTimelineItem{
			ItemType:  TimelineItemTypeShift,
			StartTime: sched.StartDatetime.Time,
			EndTime:   sched.EndDatetime.Time,
			Shift: &GetMyScheduleShiftItemInfo{
				ScheduleID:   sched.ID,
				LocationID:   sched.LocationID,
				LocationName: sched.LocationName,
			},
		})
	}

	for _, event := range events {
		if event.Kind == appointment.EventKindReminder {
			continue
		}

		dayKey := event.StartAt.UTC().Format("2006-01-02")
		dayMap[dayKey] = append(dayMap[dayKey], GetMyScheduleTimelineItem{
			ItemType:  TimelineItemTypeEvent,
			StartTime: event.StartAt,
			EndTime:   event.EndAt,
			Event: &GetMyScheduleEventItemInfo{
				EventID:            event.ID,
				MasterEventID:      event.MasterEventID,
				Title:              event.Title,
				Description:        event.Description,
				Location:           event.Location,
				Color:              event.Color,
				WorkApprovalStatus: event.WorkApprovalStatus,
				RecurrenceID:       event.RecurrenceID,
			},
		})
	}

	response := make([]GetMyScheduleTimelineDayResponse, 0, int(endDate.Sub(startDate).Hours()/24)+1)
	for day := startDate; !day.After(endDate); day = day.AddDate(0, 0, 1) {
		dayKey := day.Format("2006-01-02")
		items := append([]GetMyScheduleTimelineItem{}, dayMap[dayKey]...)
		sort.Slice(items, func(i, j int) bool {
			if items[i].StartTime.Equal(items[j].StartTime) {
				return items[i].ItemType < items[j].ItemType
			}
			return items[i].StartTime.Before(items[j].StartTime)
		})

		response = append(response, GetMyScheduleTimelineDayResponse{
			Date:  dayKey,
			Items: items,
		})
	}

	return response, nil
}

func parseTimelineDateRange(req *GetMyScheduleTimelineRequest) (time.Time, time.Time, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date format, expected YYYY-MM-DD")
	}

	if endDate.Before(startDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("end_date must be on or after start_date")
	}

	return startDate, endDate, nil
}

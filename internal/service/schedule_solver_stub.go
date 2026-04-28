//go:build !ortools

package service

import (
	"context"
	"fmt"
	"time"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

func (s *ScheduleService) generateSchedulesWithORTools(ctx context.Context, locationID uuid.UUID, timezone string, locationTZ *time.Location, employees []domain.ScheduleEmployeeContractHours, locationShifts []domain.ScheduleLocationShift, week int32, year int32) (*domain.AutoGenerateSchedulesResponse, error) {
	weekStart, err := isoWeekStartDate(int(year), int(week), locationTZ)
	if err != nil {
		return nil, err
	}

	inputsEmp, inputsShifts, _, err := buildAutoGenInputs(employees, locationShifts)
	if err != nil {
		return nil, err
	}

	response := buildAutoGenResponseSkeleton(locationID, timezone, week, year, weekStart, inputsEmp, inputsShifts)
	response.Warnings = []domain.SchedulePlanWarning{{
		Code:    "SOLVER_UNAVAILABLE",
		Message: "Auto-generation requires an ortools build. Rebuild with -tags ortools to enable this feature.",
	}}

	return response, fmt.Errorf("auto schedule generation requires an ortools build tag")
}

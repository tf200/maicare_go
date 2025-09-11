package employees

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"time"

	"go.uber.org/zap"
)

type UpdateEmployeeIsSubcontractorRequest struct {
	EmployeeID      int64
	IsSubcontractor *bool
}

type UpdateEmployeeIsSubcontractorResult struct {
	ID                int64
	IsSubcontractor   *bool
	ContractType      *string
	ContractHours     *float64
	ContractStartDate time.Time
	ContractEndDate   time.Time
	ContractRate      *float64
}

func (s *employeeService) UpdateEmployeeIsSubcontractor(
	req UpdateEmployeeIsSubcontractorRequest,
	ctx context.Context) (*UpdateEmployeeIsSubcontractorResult, error) {
	contractType := "loondienst"
	if req.IsSubcontractor != nil && *req.IsSubcontractor {
		contractType = "ZZP"
	}

	emp, err := s.Store.UpdateEmployeeIsSubcontractor(ctx, db.UpdateEmployeeIsSubcontractorParams{
		ID:              req.EmployeeID,
		IsSubcontractor: req.IsSubcontractor,
		ContractType:    &contractType,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateEmployeeIsSubcontractor", "Failed to update employee subcontractor status", zap.Error(err), zap.Int64("EmployeeID", req.EmployeeID))
		return nil, fmt.Errorf("failed to update employee subcontractor status")
	}
	res := &UpdateEmployeeIsSubcontractorResult{
		ID:                emp.ID,
		IsSubcontractor:   emp.IsSubcontractor,
		ContractType:      emp.ContractType,
		ContractHours:     emp.ContractHours,
		ContractStartDate: emp.ContractStartDate.Time,
		ContractEndDate:   emp.ContractEndDate.Time,
		ContractRate:      emp.ContractRate,
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateEmployeeIsSubcontractor", "Successfully updated employee subcontractor status", zap.Int64("EmployeeID", req.EmployeeID), zap.Bool("IsSubcontractor", *req.IsSubcontractor))
	return res, nil
}

package employees

import (
	"context"
	"fmt"
	"maicare_go/async/aclient"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/util"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type UpdateContractHoursRequest struct{}

type UpdateContractHoursResponse struct{}

// func (server *employeeService) UpdateContractHours(req UpdateContractHoursRequest) (*UpdateContractHoursResponse, error) {
// }

func (s *employeeService) AddEmployeeContractDetails(req AddEmployeeContractDetailsRequest, employeeID int64, ctx context.Context) (*AddEmployeeContractDetailsResponse, error) {
	arg := db.AddEmployeeContractDetailsParams{
		ID:                employeeID,
		ContractHours:     req.ContractHours,
		ContractStartDate: pgtype.Date{Time: req.ContractStartDate, Valid: true},
		ContractEndDate:   pgtype.Date{Time: req.ContractEndDate, Valid: true},
		ContractRate:      req.ContractRate, // Optional field, can be set later if needed
	}
	contractDetails, err := s.Store.AddEmployeeContractDetails(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeContractDetails", "Failed to add contract details to employee profile", zap.Error(err), zap.Int64("EmployeeID", employeeID))
		return nil, fmt.Errorf("failed to add contract details: %w", err)
	}

	user, err := s.Store.GetUserByID(ctx, contractDetails.UserID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeContractDetails", "Failed to get user by ID", zap.Error(err), zap.Int64("UserID", contractDetails.UserID))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	password := util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeContractDetails", "Failed to hash password", zap.Error(err), zap.Int64("UserID", user.ID))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.Store.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:       user.ID,
		Password: hashedPassword,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeContractDetails", "Failed to update password", zap.Error(err), zap.Int64("UserID", user.ID))
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	err = s.AsynqClient.EnqueueEmailDelivery(aclient.EmailDeliveryPayload{
		Name:         contractDetails.FirstName + " " + contractDetails.LastName,
		To:           contractDetails.Email,
		UserEmail:    user.Email,
		UserPassword: password,
	}, ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeContractDetails", "Failed to enqueue email delivery", zap.Error(err))
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	res := &AddEmployeeContractDetailsResponse{
		ID:                contractDetails.ID,
		ContractHours:     contractDetails.ContractHours,
		ContractStartDate: contractDetails.ContractStartDate.Time,
		ContractEndDate:   contractDetails.ContractEndDate.Time,
		ContractRate:      contractDetails.ContractRate,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "AddEmployeeContractDetails", "Successfully added contract details to employee profile", zap.Int64("EmployeeID", employeeID))
	return res, nil
}

func (s *employeeService) GetEmployeeContractDetails(employeeID int64, ctx context.Context) (*GetEmployeeContractDetailsResponse, error) {
	contractDetails, err := s.Store.GetEmployeeContractDetails(ctx, employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetEmployeeContractDetails", "Failed to get employee contract details", zap.Error(err), zap.Int64("EmployeeID", employeeID))
		return nil, fmt.Errorf("failed to get contract details: %w", err)
	}

	res := &GetEmployeeContractDetailsResponse{
		ContractHours:     contractDetails.ContractHours,
		ContractStartDate: contractDetails.ContractStartDate.Time,
		ContractEndDate:   contractDetails.ContractEndDate.Time,
		ContractType:      contractDetails.ContractType,
		ContractRate:      contractDetails.ContractRate, // Optional field for contract rate
		IsSubcontractor:   contractDetails.IsSubcontractor,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "GetEmployeeContractDetails", "Successfully retrieved employee contract details", zap.Int64("EmployeeID", employeeID))
	return res, nil
}

package employees

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *employeeService) AddEmployeeCertification(req AddEmployeeCertificationRequest, employeeID int64, ctx context.Context) (*AddEmployeeCertificationResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", req.DateIssued)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeCertification", "Failed to parse date issued", zap.Error(err), zap.Int64("EmployeeID", employeeID))
		return nil, fmt.Errorf("invalid date issued format: %w", err)
	}

	arg := db.AddEmployeeCertificationParams{
		EmployeeID: employeeID,
		Name:       req.Name,
		IssuedBy:   req.IssuedBy,
		DateIssued: pgtype.Date{Time: parsedDate, Valid: true},
	}
	certification, err := s.Store.AddEmployeeCertification(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeCertification", "Failed to add certification to employee profile", zap.Error(err), zap.Int64("EmployeeID", employeeID))
		return nil, fmt.Errorf("failed to add certification: %w", err)
	}

	res := &AddEmployeeCertificationResponse{
		ID:         certification.ID,
		EmployeeID: certification.EmployeeID,
		Name:       certification.Name,
		IssuedBy:   certification.IssuedBy,
		DateIssued: certification.DateIssued.Time,
		CreatedAt:  certification.CreatedAt.Time,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "AddEmployeeCertification", "Successfully added certification to employee profile", zap.Int64("EmployeeID", employeeID), zap.Int64("CertificationID", certification.ID))
	return res, nil
}

func (s *employeeService) ListEmployeeCertification(employeeID int64, ctx context.Context) ([]ListEmployeeCertificationResponse, error) {
	certifications, err := s.Store.ListEmployeeCertifications(ctx, employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListEmployeeCertification", "Failed to list employee certifications", zap.Error(err), zap.Int64("EmployeeID", employeeID))
		return nil, fmt.Errorf("failed to list certifications: %w", err)
	}

	responseCertifications := make([]ListEmployeeCertificationResponse, len(certifications))
	for i, certification := range certifications {
		responseCertifications[i] = ListEmployeeCertificationResponse{
			ID:         certification.ID,
			EmployeeID: certification.EmployeeID,
			Name:       certification.Name,
			IssuedBy:   certification.IssuedBy,
			DateIssued: certification.DateIssued.Time,
		}
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListEmployeeCertification", "Successfully listed employee certifications", zap.Int64("EmployeeID", employeeID), zap.Int("Count", len(certifications)))
	return responseCertifications, nil
}

func (s *employeeService) UpdateEmployeeCertification(req UpdateEmployeeCertificationRequest, certificationID int64, ctx context.Context) (*UpdateEmployeeCertificationResponse, error) {
	var parsedDate time.Time
	var err error
	if req.DateIssued != nil {
		parsedDate, err = time.Parse("2006-01-02", *req.DateIssued)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateEmployeeCertification", "Failed to parse date issued", zap.Error(err), zap.Int64("CertificationID", certificationID))
			return nil, fmt.Errorf("invalid date issued format: %w", err)
		}
	}

	certification, err := s.Store.UpdateEmployeeCertification(ctx, db.UpdateEmployeeCertificationParams{
		ID:         certificationID,
		Name:       req.Name,
		IssuedBy:   req.IssuedBy,
		DateIssued: pgtype.Date{Time: parsedDate, Valid: true},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateEmployeeCertification", "Failed to update employee certification", zap.Error(err), zap.Int64("CertificationID", certificationID))
		return nil, fmt.Errorf("failed to update certification: %w", err)
	}

	res := &UpdateEmployeeCertificationResponse{
		ID:         certification.ID,
		EmployeeID: certification.EmployeeID,
		Name:       certification.Name,
		IssuedBy:   certification.IssuedBy,
		DateIssued: certification.DateIssued.Time,
		CreatedAt:  certification.CreatedAt.Time,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateEmployeeCertification", "Successfully updated employee certification", zap.Int64("CertificationID", certificationID))
	return res, nil
}

func (s *employeeService) DeleteEmployeeCertification(certificationID int64, ctx context.Context) (*DeleteEmployeeCertificationResponse, error) {
	certification, err := s.Store.DeleteEmployeeCertification(ctx, certificationID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteEmployeeCertification", "Failed to delete employee certification", zap.Error(err), zap.Int64("CertificationID", certificationID))
		return nil, fmt.Errorf("failed to delete certification: %w", err)
	}

	res := &DeleteEmployeeCertificationResponse{
		ID:         certification.ID,
		EmployeeID: certification.EmployeeID,
		Name:       certification.Name,
		IssuedBy:   certification.IssuedBy,
		DateIssued: certification.DateIssued.Time,
		CreatedAt:  certification.CreatedAt.Time,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "DeleteEmployeeCertification", "Successfully deleted employee certification", zap.Int64("CertificationID", certificationID))
	return res, nil
}
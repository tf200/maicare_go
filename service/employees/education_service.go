package employees

import (
	"context"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *employeeService) AddEducationToEmployeeProfile(
	req AddEducationToEmployeeProfileRequest,
	employeeID uuid.UUID,
	ctx context.Context,
) (*AddEducationToEmployeeProfileResponse, error) {
	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"AddEducationToEmployeeProfile",
			"Failed to parse start date",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
		)
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"AddEducationToEmployeeProfile",
			"Failed to parse end date",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
		)
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	arg := db.AddEducationToEmployeeProfileParams{
		EmployeeID:      employeeID,
		InstitutionName: req.InstitutionName,
		Degree:          req.Degree,
		FieldOfStudy:    req.FieldOfStudy,
		StartDate:       pgtype.Date{Time: parsedStartDate, Valid: true},
		EndDate:         pgtype.Date{Time: parsedEndDate, Valid: true},
	}
	education, err := s.Store.AddEducationToEmployeeProfile(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"AddEducationToEmployeeProfile",
			"Failed to add education to employee profile",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
		)
		return nil, fmt.Errorf("failed to add education: %w", err)
	}

	res := &AddEducationToEmployeeProfileResponse{
		ID:              education.ID,
		EmployeeID:      education.EmployeeID,
		InstitutionName: education.InstitutionName,
		Degree:          education.Degree,
		FieldOfStudy:    education.FieldOfStudy,
		StartDate:       education.StartDate.Time,
		EndDate:         education.EndDate.Time,
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"AddEducationToEmployeeProfile",
		"Successfully added education to employee profile",
		zap.String("EmployeeID", employeeID.String()),
		zap.Int64("EducationID", education.ID),
	)
	return res, nil
}

func (s *employeeService) ListEmployeeEducation(
	employeeID uuid.UUID,
	ctx context.Context,
) ([]ListEmployeeEducationResponse, error) {
	educations, err := s.Store.ListEducations(ctx, employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"ListEmployeeEducation",
			"Failed to list employee education",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
		)
		return nil, fmt.Errorf("failed to list education: %w", err)
	}

	responseEducations := make([]ListEmployeeEducationResponse, len(educations))
	for i, education := range educations {
		responseEducations[i] = ListEmployeeEducationResponse{
			ID:              education.ID,
			EmployeeID:      education.EmployeeID,
			InstitutionName: education.InstitutionName,
			Degree:          education.Degree,
			FieldOfStudy:    education.FieldOfStudy,
			StartDate:       education.StartDate.Time,
			EndDate:         education.EndDate.Time,
		}
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"ListEmployeeEducation",
		"Successfully listed employee education",
		zap.String("EmployeeID", employeeID.String()),
		zap.Int("Count", len(educations)),
	)
	return responseEducations, nil
}

func (s *employeeService) UpdateEmployeeEducation(
	req UpdateEmployeeEducationRequest,
	educationID int64,
	ctx context.Context,
) (*UpdateEmployeeEducationResponse, error) {
	var parsedStartDate time.Time
	var err error
	if req.StartDate != nil {
		parsedStartDate, err = time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			s.Logger.LogBusinessEvent(
				ctx,
				logger.LogLevelError,
				"UpdateEmployeeEducation",
				"Failed to parse start date",
				zap.Error(err),
				zap.Int64("EducationID", educationID),
			)
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
	}
	var parsedEndDate time.Time
	if req.EndDate != nil {
		parsedEndDate, err = time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			s.Logger.LogBusinessEvent(
				ctx,
				logger.LogLevelError,
				"UpdateEmployeeEducation",
				"Failed to parse end date",
				zap.Error(err),
				zap.Int64("EducationID", educationID),
			)
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
	}

	education, err := s.Store.UpdateEmployeeEducation(ctx, db.UpdateEmployeeEducationParams{
		ID:              educationID,
		InstitutionName: req.InstitutionName,
		Degree:          req.Degree,
		FieldOfStudy:    req.FieldOfStudy,
		StartDate:       pgtype.Date{Time: parsedStartDate, Valid: true},
		EndDate:         pgtype.Date{Time: parsedEndDate, Valid: true},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"UpdateEmployeeEducation",
			"Failed to update employee education",
			zap.Error(err),
			zap.Int64("EducationID", educationID),
		)
		return nil, fmt.Errorf("failed to update education: %w", err)
	}

	res := &UpdateEmployeeEducationResponse{
		ID:              education.ID,
		InstitutionName: education.InstitutionName,
		Degree:          education.Degree,
		FieldOfStudy:    education.FieldOfStudy,
		StartDate:       education.StartDate.Time,
		EndDate:         education.EndDate.Time,
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"UpdateEmployeeEducation",
		"Successfully updated employee education",
		zap.Int64("EducationID", educationID),
	)
	return res, nil
}

func (s *employeeService) DeleteEmployeeEducation(
	educationID int64,
	ctx context.Context,
) (*DeleteEmployeeEducationResponse, error) {
	education, err := s.Store.DeleteEmployeeEducation(ctx, educationID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"DeleteEmployeeEducation",
			"Failed to delete employee education",
			zap.Error(err),
			zap.Int64("EducationID", educationID),
		)
		return nil, fmt.Errorf("failed to delete education: %w", err)
	}

	res := &DeleteEmployeeEducationResponse{
		ID:              education.ID,
		InstitutionName: education.InstitutionName,
		Degree:          education.Degree,
		FieldOfStudy:    education.FieldOfStudy,
		StartDate:       education.StartDate.Time,
		EndDate:         education.EndDate.Time,
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"DeleteEmployeeEducation",
		"Successfully deleted employee education",
		zap.Int64("EducationID", educationID),
	)
	return res, nil
}

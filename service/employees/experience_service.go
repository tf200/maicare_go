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

func (s *employeeService) AddEmployeeExperience(req AddEmployeeExperienceRequest, employeeID int64, ctx context.Context) (*AddEmployeeExperienceResponse, error) {
	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeExperience", "Failed to parse start date", zap.Error(err), zap.Int64("EmployeeID", employeeID))
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeExperience", "Failed to parse end date", zap.Error(err), zap.Int64("EmployeeID", employeeID))
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	arg := db.AddEmployeeExperienceParams{
		EmployeeID:  employeeID,
		JobTitle:    req.JobTitle,
		CompanyName: req.CompanyName,
		StartDate:   pgtype.Date{Time: parsedStartDate, Valid: true},
		EndDate:     pgtype.Date{Time: parsedEndDate, Valid: true},
		Description: req.Description,
	}
	experience, err := s.Store.AddEmployeeExperience(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddEmployeeExperience", "Failed to add experience to employee profile", zap.Error(err), zap.Int64("EmployeeID", employeeID))
		return nil, fmt.Errorf("failed to add experience: %w", err)
	}

	res := &AddEmployeeExperienceResponse{
		ID:          experience.ID,
		EmployeeID:  experience.EmployeeID,
		JobTitle:    experience.JobTitle,
		CompanyName: experience.CompanyName,
		StartDate:   experience.StartDate.Time.Format(time.RFC3339),
		EndDate:     experience.EndDate.Time.Format(time.RFC3339),
		Description: experience.Description,
		CreatedAt:   experience.CreatedAt.Time.Format(time.RFC3339),
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "AddEmployeeExperience", "Successfully added experience to employee profile", zap.Int64("EmployeeID", employeeID), zap.Int64("ExperienceID", experience.ID))
	return res, nil
}

func (s *employeeService) ListEmployeeExperience(employeeID int64, ctx context.Context) ([]ListEmployeeExperienceResponse, error) {
	experiences, err := s.Store.ListEmployeeExperience(ctx, employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListEmployeeExperience", "Failed to list employee experience", zap.Error(err), zap.Int64("EmployeeID", employeeID))
		return nil, fmt.Errorf("failed to list experience: %w", err)
	}

	responseExperiences := make([]ListEmployeeExperienceResponse, len(experiences))
	for i, experience := range experiences {
		responseExperiences[i] = ListEmployeeExperienceResponse{
			ID:          experience.ID,
			EmployeeID:  experience.EmployeeID,
			JobTitle:    experience.JobTitle,
			CompanyName: experience.CompanyName,
			StartDate:   experience.StartDate.Time.Format(time.RFC3339),
			EndDate:     experience.EndDate.Time.Format(time.RFC3339),
			Description: experience.Description,
			CreatedAt:   experience.CreatedAt.Time.Format(time.RFC3339),
		}
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListEmployeeExperience", "Successfully listed employee experience", zap.Int64("EmployeeID", employeeID), zap.Int("Count", len(experiences)))
	return responseExperiences, nil
}

func (s *employeeService) UpdateEmployeeExperience(req UpdateEmployeeExperienceRequest, experienceID int64, ctx context.Context) (*UpdateEmployeeExperienceResponse, error) {
	var parsedStartDate time.Time
	var err error
	if req.StartDate != nil {
		parsedStartDate, err = time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateEmployeeExperience", "Failed to parse start date", zap.Error(err), zap.Int64("ExperienceID", experienceID))
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
	}
	var parsedEndDate time.Time
	if req.EndDate != nil {
		parsedEndDate, err = time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateEmployeeExperience", "Failed to parse end date", zap.Error(err), zap.Int64("ExperienceID", experienceID))
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
	}

	experience, err := s.Store.UpdateEmployeeExperience(ctx, db.UpdateEmployeeExperienceParams{
		ID:          experienceID,
		JobTitle:    req.JobTitle,
		CompanyName: req.CompanyName,
		StartDate:   pgtype.Date{Time: parsedStartDate, Valid: true},
		EndDate:     pgtype.Date{Time: parsedEndDate, Valid: true},
		Description: req.Description,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateEmployeeExperience", "Failed to update employee experience", zap.Error(err), zap.Int64("ExperienceID", experienceID))
		return nil, fmt.Errorf("failed to update experience: %w", err)
	}

	res := &UpdateEmployeeExperienceResponse{
		ID:          experience.ID,
		EmployeeID:  experience.EmployeeID,
		JobTitle:    experience.JobTitle,
		CompanyName: experience.CompanyName,
		StartDate:   experience.StartDate.Time,
		EndDate:     experience.EndDate.Time,
		Description: experience.Description,
		CreatedAt:   experience.CreatedAt.Time,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateEmployeeExperience", "Successfully updated employee experience", zap.Int64("ExperienceID", experienceID))
	return res, nil
}

func (s *employeeService) DeleteEmployeeExperience(experienceID int64, ctx context.Context) (*DeleteEmployeeExperienceResponse, error) {
	experience, err := s.Store.DeleteEmployeeExperience(ctx, experienceID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteEmployeeExperience", "Failed to delete employee experience", zap.Error(err), zap.Int64("ExperienceID", experienceID))
		return nil, fmt.Errorf("failed to delete experience: %w", err)
	}

	res := &DeleteEmployeeExperienceResponse{
		ID:          experience.ID,
		EmployeeID:  experience.EmployeeID,
		JobTitle:    experience.JobTitle,
		CompanyName: experience.CompanyName,
		StartDate:   experience.StartDate.Time,
		EndDate:     experience.EndDate.Time,
		Description: experience.Description,
		CreatedAt:   experience.CreatedAt.Time,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "DeleteEmployeeExperience", "Successfully deleted employee experience", zap.Int64("ExperienceID", experienceID))
	return res, nil
}
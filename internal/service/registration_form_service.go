package service

import (
	"context"
	"fmt"

	"maicare_go/internal/domain"
	"maicare_go/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RegistrationFormService struct {
	repo      *repository.RegistrationFormRepository
	logger    domain.Logger
	taskQueue domain.TaskQueue
}

func NewRegistrationFormService(repo *repository.RegistrationFormRepository, logger domain.Logger, taskQueue domain.TaskQueue) domain.RegistrationFormService {
	return &RegistrationFormService{
		repo:      repo,
		logger:    logger,
		taskQueue: taskQueue,
	}
}

func (s *RegistrationFormService) CreateRegistrationForm(ctx context.Context, params domain.CreateRegistrationFormParams) (*domain.RegistrationForm, error) {
	form, err := s.repo.CreateRegistrationForm(ctx, params)
	if err != nil {
		s.logError(ctx, "CreateRegistrationForm", err, zap.String("client_name", params.ClientFirstName+" "+params.ClientLastName))
		return nil, err
	}
	return form, nil
}

func (s *RegistrationFormService) ListRegistrationForms(ctx context.Context, params domain.ListRegistrationFormsParams) (*domain.ListResult[domain.RegistrationFormListItem], error) {
	result, err := s.repo.ListRegistrationForms(ctx, params)
	if err != nil {
		s.logError(ctx, "ListRegistrationForms", err)
		return nil, err
	}
	return result, nil
}

func (s *RegistrationFormService) GetRegistrationForm(ctx context.Context, id uuid.UUID) (*domain.RegistrationForm, error) {
	form, err := s.repo.GetRegistrationForm(ctx, id)
	if err != nil {
		s.logError(ctx, "GetRegistrationForm", err, zap.String("form_id", id.String()))
		return nil, err
	}
	return form, nil
}

func (s *RegistrationFormService) UpdateRegistrationForm(ctx context.Context, params domain.UpdateRegistrationFormParams) (*domain.RegistrationForm, error) {
	form, err := s.repo.UpdateRegistrationForm(ctx, params)
	if err != nil {
		s.logError(ctx, "UpdateRegistrationForm", err, zap.String("form_id", params.ID.String()))
		return nil, err
	}
	return form, nil
}

func (s *RegistrationFormService) DeleteRegistrationForm(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteRegistrationForm(ctx, id)
	if err != nil {
		s.logError(ctx, "DeleteRegistrationForm", err, zap.String("form_id", id.String()))
		return err
	}
	return nil
}

func (s *RegistrationFormService) UpdateRegistrationFormStatus(ctx context.Context, params domain.UpdateRegistrationFormStatusParams) error {
	err := s.repo.UpdateRegistrationFormStatus(ctx, params)
	if err != nil {
		s.logError(ctx, "UpdateRegistrationFormStatus", err, zap.String("form_id", params.ID.String()))
		return err
	}
	return nil
}

func (s *RegistrationFormService) ProcessRegistrationForm(ctx context.Context, params domain.ProcessRegistrationFormParams) error {
	form, token, err := s.repo.ProcessRegistrationForm(ctx, params)
	if err != nil {
		s.logError(ctx, "ProcessRegistrationForm", err, zap.String("form_id", params.ID.String()))
		return err
	}

	// Build email recipients from referrer + guardian emails
	recipients := []string{}
	if form.ReferrerEmail != "" {
		recipients = append(recipients, form.ReferrerEmail)
	}
	if form.Guardian1Email != "" {
		recipients = append(recipients, form.Guardian1Email)
	}
	if form.Guardian2Email != "" {
		recipients = append(recipients, form.Guardian2Email)
	}

	// Construct link
	link := fmt.Sprintf("https://maicare.online/intake/schedule/%s", token)

	// Build email payload
	emailData := domain.ProcessRegistrationFormEmailTaskPayload{
		ReferrerName: form.ReferrerFirstName,
		ClientName:   form.ClientFirstName + " " + form.ClientLastName,
		Location:     derefString(form.IntakeAppointmentLocation),
		Link:         link,
		To:           recipients,
	}

	// Enqueue email task - log error if fails but don't fail the request
	if err := s.taskQueue.EnqueueProcessRegistrationFormEmail(ctx, emailData, nil); err != nil {
		s.logError(ctx, "ProcessRegistrationForm::EnqueueEmail", err,
			zap.String("form_id", params.ID.String()),
			zap.String("link", link),
		)
	}

	return nil
}

func (s *RegistrationFormService) GetPublicIntakeOptions(ctx context.Context, token string) (*domain.PublicIntakeOptions, error) {
	options, err := s.repo.GetPublicIntakeOptions(ctx, token)
	if err != nil {
		s.logError(ctx, "GetPublicIntakeOptions", err, zap.String("token", token))
		return nil, err
	}
	return options, nil
}

func (s *RegistrationFormService) SelectIntakeDate(ctx context.Context, params domain.SelectIntakeDateParams) error {
	err := s.repo.SelectIntakeDate(ctx, params)
	if err != nil {
		s.logError(ctx, "SelectIntakeDate", err, zap.String("token", params.Token))
		return err
	}
	return nil
}

func (s *RegistrationFormService) logError(ctx context.Context, method string, err error, fields ...zap.Field) {
	if s.logger != nil {
		s.logger.LogError(ctx, "RegistrationFormService."+method, err.Error(), err, fields...)
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to calculate risk count for list operations
func calculateRiskCount(
	riskAggressiveBehavior *bool,
	riskSuicidalSelfharm *bool,
	riskSubstanceAbuse *bool,
	riskPsychiatricIssues *bool,
	riskCriminalHistory *bool,
	riskFlightBehavior *bool,
	riskWeaponPossession *bool,
	riskSexualBehavior *bool,
	riskDayNightRhythm *bool,
	riskOther *bool,
) int {
	count := 0
	isTruePtr := func(b *bool) bool {
		return b != nil && *b
	}

	if isTruePtr(riskAggressiveBehavior) {
		count++
	}
	if isTruePtr(riskSuicidalSelfharm) {
		count++
	}
	if isTruePtr(riskSubstanceAbuse) {
		count++
	}
	if isTruePtr(riskPsychiatricIssues) {
		count++
	}
	if isTruePtr(riskCriminalHistory) {
		count++
	}
	if isTruePtr(riskFlightBehavior) {
		count++
	}
	if isTruePtr(riskWeaponPossession) {
		count++
	}
	if isTruePtr(riskSexualBehavior) {
		count++
	}
	if isTruePtr(riskDayNightRhythm) {
		count++
	}
	if isTruePtr(riskOther) {
		count++
	}

	return count
}

var _ domain.RegistrationFormService = (*RegistrationFormService)(nil)

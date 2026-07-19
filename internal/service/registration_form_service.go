package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RegistrationFormService struct {
	repo        domain.RegistrationFormRepository
	uploads     domain.RegistrationUploadSessionRepository
	attachments domain.AttachmentService
	logger      domain.Logger
	audit       domain.AuditLogger
	taskQueue   domain.TaskQueue
}

func NewRegistrationFormService(repo domain.RegistrationFormRepository, uploads domain.RegistrationUploadSessionRepository, attachments domain.AttachmentService, logger domain.Logger, taskQueue domain.TaskQueue, audit domain.AuditLogger) domain.RegistrationFormService {
	return &RegistrationFormService{
		repo:        repo,
		uploads:     uploads,
		attachments: attachments,
		logger:      logger,
		taskQueue:   taskQueue,
		audit:       audit,
	}
}

func (s *RegistrationFormService) StartUploadSession(ctx context.Context) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)
	hash := sha256.Sum256([]byte(token))
	if _, err := s.uploads.Create(ctx, hex.EncodeToString(hash[:]), time.Now().Add(24*time.Hour)); err != nil {
		return "", err
	}
	return token, nil
}

func (s *RegistrationFormService) GetRegistrationFormCounts(ctx context.Context) (domain.RegistrationFormCounts, error) {
	return s.repo.GetRegistrationFormCounts(ctx)
}

func (s *RegistrationFormService) InitRegistrationUpload(ctx context.Context, token string, params domain.InitAttachmentUploadParams) (*domain.InitAttachmentUploadResult, error) {
	if params.Size > 20<<20 {
		return nil, fmt.Errorf("file size exceeds maximum limit of 20MB")
	}
	if params.ContentType != "application/pdf" && params.ContentType != "image/jpeg" && params.ContentType != "image/png" {
		return nil, fmt.Errorf("unsupported file type: %s", params.ContentType)
	}
	session, err := s.uploads.GetActive(ctx, hashUploadToken(token))
	if err != nil {
		return nil, fmt.Errorf("invalid registration upload session")
	}
	result, err := s.attachments.InitUpload(ctx, params)
	if err != nil {
		return nil, err
	}
	if err := s.uploads.AddAttachment(ctx, session.ID, result.FileID); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *RegistrationFormService) CreateRegistrationForm(ctx context.Context, params domain.CreateRegistrationFormParams) (*domain.RegistrationForm, error) {
	ids := registrationAttachmentIDs(params)
	session, err := s.uploads.GetActive(ctx, hashUploadToken(params.RegistrationUploadToken))
	if err != nil {
		return nil, fmt.Errorf("invalid registration upload session")
	}
	valid, err := s.uploads.HasAttachments(ctx, session.ID, ids)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid registration attachments")
	}
	for _, id := range ids {
		if _, err := s.attachments.ConfirmUpload(ctx, id); err != nil {
			return nil, err
		}
	}
	form, err := s.repo.CreateRegistrationForm(ctx, params)
	if err != nil {
		s.logError(ctx, "CreateRegistrationForm", err, zap.String("client_name", params.ClientFirstName+" "+params.ClientLastName))
		return nil, err
	}
	if err := s.uploads.Consume(ctx, session.ID); err != nil {
		return nil, err
	}
	return form, nil
}

func hashUploadToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func registrationAttachmentIDs(params domain.CreateRegistrationFormParams) []uuid.UUID {
	ids := make([]uuid.UUID, 0, 6)
	for _, id := range []*uuid.UUID{params.DocumentReferral, params.DocumentEducationReport, params.DocumentPsychiatricReport, params.DocumentDiagnosis, params.DocumentSafetyPlan, params.DocumentIDCopy} {
		if id != nil {
			ids = append(ids, *id)
		}
	}
	return ids
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

	if s.audit != nil {
		fid := id.String()
		if auditErr := s.audit.Log(ctx, domain.AuditEvent{
			EventType:   "record_access",
			Action:      "read",
			Result:      "success",
			SubjectType: "registration_form",
			SubjectID:   fid,
			AccessRule:  strPtr("REGISTRATION_FORM.VIEW"),
		}); auditErr != nil {
			s.logError(ctx, "GetRegistrationForm::Audit", auditErr, zap.String("form_id", fid))
		}
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

func (s *RegistrationFormService) ReplaceRegistrationFormDocument(ctx context.Context, params domain.ReplaceRegistrationFormDocumentParams) (*domain.RegistrationForm, error) {
	if _, err := s.attachments.ConfirmUpload(ctx, params.FileID); err != nil {
		s.logError(ctx, "ReplaceRegistrationFormDocument::ConfirmUpload", err, zap.String("file_id", params.FileID.String()))
		return nil, err
	}
	form, err := s.repo.ReplaceRegistrationFormDocument(ctx, params)
	if err != nil {
		s.logError(ctx, "ReplaceRegistrationFormDocument", err, zap.String("form_id", params.ID.String()))
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

var _ domain.RegistrationFormService = (*RegistrationFormService)(nil)

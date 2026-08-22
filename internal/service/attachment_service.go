package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var attachmentCategories = map[string]domain.AttachmentCategory{
	"image/jpeg": domain.AttachmentImage, "image/jpg": domain.AttachmentImage,
	"image/png": domain.AttachmentImage, "image/gif": domain.AttachmentImage,
	"image/webp": domain.AttachmentImage, "image/svg+xml": domain.AttachmentImage,
	"application/pdf": domain.AttachmentDocument, "application/msword": domain.AttachmentDocument,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": domain.AttachmentDocument,
	"application/vnd.ms-excel": domain.AttachmentDocument,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": domain.AttachmentDocument,
	"text/plain": domain.AttachmentDocument, "text/csv": domain.AttachmentDocument,
	"application/zip": domain.AttachmentArchive, "application/x-rar": domain.AttachmentArchive,
	"application/x-7z-compressed": domain.AttachmentArchive,
}

type AttachmentService struct {
	repository domain.AttachmentRepository
	storage    domain.Storage
	logger     domain.Logger
}

func NewAttachmentService(repository domain.AttachmentRepository, storage domain.Storage, logger domain.Logger) domain.AttachmentService {
	return &AttachmentService{repository: repository, storage: storage, logger: logger}
}

func (s *AttachmentService) InitUpload(ctx context.Context, params domain.InitAttachmentUploadParams) (*domain.InitAttachmentUploadResult, error) {
	if params.Size <= 0 {
		return nil, errors.New("file size must be greater than 0")
	}
	if params.Size > domain.MaxAttachmentSize {
		return nil, fmt.Errorf("file size exceeds maximum limit of %dMB", domain.MaxAttachmentSize>>20)
	}
	category, ok := attachmentCategories[params.ContentType]
	if !ok {
		return nil, fmt.Errorf("unsupported file type: %s", params.ContentType)
	}

	id := uuid.New()
	key := attachmentObjectKey(id, params.Filename, category)
	name := sanitizeAttachmentFilename(params.Filename)
	if _, err := s.repository.CreateAttachment(ctx, domain.CreateAttachmentParams{
		ID: id, Name: name, File: key, Size: int32(params.Size), Tag: &name,
	}); err != nil {
		s.logError(ctx, "AttachmentService.InitUpload", "failed to create attachment record", err, zap.String("file_id", id.String()))
		return nil, fmt.Errorf("failed to create attachment record: %w", err)
	}

	uploadURL, err := s.storage.GeneratePresignedUploadURL(ctx, key, 15*time.Minute)
	if err != nil {
		s.logError(ctx, "AttachmentService.InitUpload", "failed to generate upload URL", err, zap.String("file_id", id.String()))
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}
	return &domain.InitAttachmentUploadResult{UploadURL: uploadURL, FileID: id, Key: key}, nil
}

func (s *AttachmentService) ConfirmUpload(ctx context.Context, id uuid.UUID) (*domain.ConfirmAttachmentUploadResult, error) {
	attachment, err := s.repository.GetAttachmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("attachment not found: %w", err)
	}
	size, err := s.storage.GetFileInfo(ctx, attachment.File)
	if err != nil {
		return nil, fmt.Errorf("file verification failed: %w", err)
	}
	if size <= 0 {
		return nil, errors.New("file is empty")
	}
	if size > domain.MaxAttachmentSize {
		return nil, fmt.Errorf("file size exceeds maximum limit of %dMB", domain.MaxAttachmentSize>>20)
	}
	if attachment.Size > 0 && size != int64(attachment.Size) {
		s.logWarn(ctx, "AttachmentService.ConfirmUpload", "uploaded file size differs from initialized size", zap.String("file_id", id.String()))
	}

	updated, err := s.repository.SetAttachmentUsed(ctx, id, true)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm upload: %w", err)
	}
	fileURL, err := s.storage.GeneratePresignedURL(ctx, updated.File, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate download URL: %w", err)
	}
	return &domain.ConfirmAttachmentUploadResult{FileURL: fileURL, FileID: updated.UUID, CreatedAt: updated.CreatedAt, Size: size}, nil
}

func (s *AttachmentService) GetAttachment(ctx context.Context, id uuid.UUID) (*domain.AttachmentResult, error) {
	attachment, err := s.repository.GetAttachmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}
	fileURL, err := s.storage.GeneratePresignedURL(ctx, attachment.File, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return &domain.AttachmentResult{FileURL: fileURL, FileID: attachment.UUID, CreatedAt: attachment.CreatedAt, Size: int64(attachment.Size)}, nil
}

func (s *AttachmentService) DeleteAttachment(ctx context.Context, id uuid.UUID) (*domain.AttachmentResult, error) {
	deleted, err := s.repository.DeleteAttachment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete attachment: %w", err)
	}
	if err := s.storage.Delete(ctx, deleted.File); err != nil {
		s.logError(ctx, "AttachmentService.DeleteAttachment", "attachment record deleted but storage cleanup failed", err, zap.String("file_id", id.String()))
		if _, restoreErr := s.repository.CreateAttachment(ctx, domain.CreateAttachmentParams{
			ID: deleted.UUID, Name: deleted.Name, File: deleted.File, Size: deleted.Size, Tag: deleted.Tag,
		}); restoreErr != nil {
			s.logError(ctx, "AttachmentService.DeleteAttachment", "failed to restore attachment after storage cleanup failure", restoreErr, zap.String("file_id", id.String()))
		} else if deleted.IsUsed {
			if _, restoreErr := s.repository.SetAttachmentUsed(ctx, id, true); restoreErr != nil {
				s.logError(ctx, "AttachmentService.DeleteAttachment", "failed to restore attachment usage state", restoreErr, zap.String("file_id", id.String()))
			}
		}
		return nil, fmt.Errorf("failed to delete file from storage: %w", err)
	}
	return &domain.AttachmentResult{FileURL: deleted.File, FileID: deleted.UUID, CreatedAt: deleted.CreatedAt, Size: int64(deleted.Size)}, nil
}

func attachmentObjectKey(id uuid.UUID, filename string, category domain.AttachmentCategory) string {
	now := time.Now().UTC()
	return fmt.Sprintf("%s/%d/%02d/%s_%s", category, now.Year(), now.Month(), id, sanitizeAttachmentFilename(filename))
}

func sanitizeAttachmentFilename(filename string) string {
	name := filepath.Base(strings.TrimSpace(filename))
	name = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r < 32 {
			return '_'
		}
		return r
	}, name)
	name = strings.TrimSpace(name)
	if name == "" || name == "." {
		return "file"
	}
	if len(name) > 100 {
		name = name[:100]
	}
	return name
}

func (s *AttachmentService) logError(ctx context.Context, operation, message string, err error, fields ...zap.Field) {
	if s.logger != nil {
		s.logger.LogError(ctx, operation, message, err, fields...)
	}
}

func (s *AttachmentService) logWarn(ctx context.Context, operation, message string, fields ...zap.Field) {
	if s.logger != nil {
		s.logger.LogWarn(ctx, operation, message, fields...)
	}
}

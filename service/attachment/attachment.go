package attachment

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *attachmentService) InitUpload(ctx context.Context, req *InitUploadRequest) (*InitUploadResponse, error) {
	if req.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum limit of 100MB")
	}

	category, allowed := allowedMimeTypes[req.ContentType]
	if !allowed {
		return nil, fmt.Errorf("unsupported file type: %s", req.ContentType)
	}

	// Calculate secure key
	key, fileUUID := s.generateSecureKey(req.Filename, category)

	// Create initial record in DB (unused)
	arg := db.CreateAttachmentParams{
		Uuid: fileUUID,
		Name: req.Filename,
		File: key,
		Size: int32(req.Size),
		Tag:  &req.Filename,
	}

	_, err := s.Store.CreateAttachment(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "InitUpload", "Failed to create attachment record", zap.Error(err))
		return nil, fmt.Errorf("failed to create attachment record: %v", err)
	}

	// Generate presigned upload URL
	uploadURL, err := s.B2Client.GeneratePresignedUploadURL(ctx, key, 15*time.Minute)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "InitUpload", "Failed to generate presigned upload URL", zap.Error(err))
		return nil, fmt.Errorf("failed to generate upload URL: %v", err)
	}

	return &InitUploadResponse{
		UploadURL: uploadURL,
		FileID:    fileUUID,
		Key:       key,
	}, nil
}

func (s *attachmentService) ConfirmUpload(ctx context.Context, req *ConfirmUploadRequest) (*ConfirmUploadResponse, error) {
	attachment, err := s.Store.GetAttachmentById(ctx, req.FileID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ConfirmUpload", "Attachment not found", zap.Error(err))
		return nil, fmt.Errorf("attachment not found")
	}

	// Check if file exists in storage
	size, err := s.B2Client.GetFileInfo(ctx, attachment.File)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ConfirmUpload", "File verification failed", zap.Error(err))
		return nil, fmt.Errorf("file verification failed: %v", err)
	}

	if size == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	// Mark as used
	updatedAttachment, err := s.Store.SetAttachmentAsUsedorUnused(ctx, db.SetAttachmentAsUsedorUnusedParams{
		Uuid:   req.FileID,
		IsUsed: true,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ConfirmUpload", "Failed to update attachment status", zap.Error(err))
		return nil, fmt.Errorf("failed to confirm upload: %v", err)
	}

	// Generate download URL for response
	url, err := s.B2Client.GeneratePresignedURL(ctx, updatedAttachment.File, 15*time.Minute)
	if err != nil {
		// Log error but don't fail the confirmation? Or fail?
		// Better to return what we have, URL generation is secondary here?
		// But response expects URL.
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ConfirmUpload", "Failed to generate download URL", zap.Error(err))
		return nil, fmt.Errorf("failed to generate download URL")
	}

	return &ConfirmUploadResponse{
		FileURL:   url,
		FileID:    updatedAttachment.Uuid,
		CreatedAt: updatedAttachment.Created.Time,
		Size:      size,
	}, nil
}

func (s *attachmentService) GetAttachmentById(ctx context.Context, id uuid.UUID) (*GetAttachmentByIdResponse, error) {
	attachment, err := s.Store.GetAttachmentById(ctx, id)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetAttachmentById", "Failed to get attachment by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get attachment")
	}
	url := s.GenerateResponsePresignedURL(&attachment.File, ctx)
	if url == nil {
		return nil, fmt.Errorf("failed to generate presigned URL")
	}
	return &GetAttachmentByIdResponse{
		FileURL:   *url,
		FileID:    attachment.Uuid,
		CreatedAt: attachment.Created.Time,
		Size:      int64(attachment.Size),
	}, nil
}

func (s *attachmentService) DeleteAttachment(ctx context.Context, id uuid.UUID) (*DeleteAttachmentResponse, error) {
	attachment, err := s.Store.GetAttachmentById(ctx, id)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteAttachment", "Failed to get attachment by ID", zap.Error(err))

		return nil, fmt.Errorf("failed to get attachment")
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteAttachment", "Failed to begin transaction", zap.Error(err))

		return nil, fmt.Errorf("failed to begin transaction")
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteAttachment", "Failed to rollback transaction", zap.Error(rollbackErr))
		}
	}()
	qtx := s.Store.WithTx(tx)

	attachment, err = qtx.DeleteAttachment(ctx, id)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteAttachment", "Failed to delete attachment record", zap.Error(err))

		return nil, fmt.Errorf("failed to delete attachment record")
	}

	err = s.B2Client.Delete(ctx, attachment.File)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteAttachment", "Failed to delete attachment from B2", zap.Error(err))
		return nil, fmt.Errorf("failed to delete attachment from B2: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteAttachment", "Failed to commit transaction", zap.Error(err))

		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &DeleteAttachmentResponse{
		FileURL:   attachment.File,
		FileID:    attachment.Uuid,
		CreatedAt: attachment.Created.Time,
		Size:      int64(attachment.Size),
	}, nil
}

func (s *attachmentService) generateSecureKey(filename string, category FileCategory) (string, uuid.UUID) {
	now := time.Now().UTC()

	cleanFilename := s.sanitizeFilename(filename)

	uuid := uuid.New()
	key := fmt.Sprintf("%s/%d/%02d/%s_%s",
		string(category),
		now.Year(),
		now.Month(),
		uuid,
		cleanFilename,
	)
	return key, uuid
}

func (s *attachmentService) sanitizeFilename(filename string) string {
	// Remove or replace problematic characters
	clean := strings.ReplaceAll(filename, " ", "_")
	clean = strings.ReplaceAll(clean, "+", "_")
	clean = strings.ReplaceAll(clean, "&", "_")
	clean = strings.ReplaceAll(clean, "=", "_")

	// Limit filename length
	if len(clean) > 100 {
		ext := filepath.Ext(clean)
		name := clean[:100-len(ext)] + ext
		clean = name
	}

	return clean
}

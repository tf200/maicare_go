package attachment

import (
	"context"
	"errors"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (s *attachmentService) InitUpload(ctx context.Context, req *InitUploadRequest) (*InitUploadResponse, error) {
	if req.Size <= 0 {
		return nil, fmt.Errorf("file size must be greater than 0")
	}
	if req.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum limit of %dMB", MaxFileSize>>20)
	}

	category, allowed := allowedMimeTypes[req.ContentType]
	if !allowed {
		return nil, fmt.Errorf("unsupported file type: %s", req.ContentType)
	}

	// Calculate secure key
	key, fileUUID := s.generateSecureKey(req.Filename, category)

	arg := db.CreateAttachmentParams{
		Uuid: fileUUID,
		Name: req.Filename,
		File: key,
		Size: int32(req.Size),
		Tag:  &req.Filename,
	}

	_, err := s.Store.CreateAttachment(ctx, arg)
	if err != nil {
		s.Logger.LogError(ctx, "InitUpload", "Failed to create attachment record", err, zap.String("file_id", fileUUID.String()))
		return nil, fmt.Errorf("failed to create attachment record: %w", err)
	}

	// Generate presigned upload URL
	uploadURL, err := s.B2Client.GeneratePresignedUploadURL(ctx, key, 15*time.Minute)
	if err != nil {
		s.Logger.LogError(ctx, "InitUpload", "Failed to generate presigned upload URL", err, zap.String("file_id", fileUUID.String()))
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
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
		s.Logger.LogError(ctx, "ConfirmUpload", "Attachment not found", err, zap.String("file_id", req.FileID.String()))
		return nil, fmt.Errorf("attachment not found")
	}

	// Check if file exists in storage
	size, err := s.B2Client.GetFileInfo(ctx, attachment.File)
	if err != nil {
		s.Logger.LogError(ctx, "ConfirmUpload", "File verification failed", err, zap.String("file_id", req.FileID.String()))
		return nil, fmt.Errorf("file verification failed: %w", err)
	}

	if size == 0 {
		return nil, fmt.Errorf("file is empty")
	}
	if size > MaxFileSize {
		s.Logger.LogWarn(ctx, "ConfirmUpload", "Uploaded file exceeds maximum size", zap.String("file_id", req.FileID.String()), zap.Int64("size", size), zap.Int64("max_size", MaxFileSize))
		return nil, fmt.Errorf("file size exceeds maximum limit of %dMB", MaxFileSize>>20)
	}
	if attachment.Size > 0 && size != int64(attachment.Size) {
		s.Logger.LogWarn(ctx, "ConfirmUpload", "Uploaded file size differs from initialized size", zap.String("file_id", req.FileID.String()), zap.Int64("initialized_size", int64(attachment.Size)), zap.Int64("uploaded_size", size))
	}

	// Mark as used
	updatedAttachment, err := s.Store.SetAttachmentAsUsedorUnused(ctx, db.SetAttachmentAsUsedorUnusedParams{
		Uuid:   req.FileID,
		IsUsed: true,
	})
	if err != nil {
		s.Logger.LogError(ctx, "ConfirmUpload", "Failed to update attachment status", err, zap.String("file_id", req.FileID.String()))
		return nil, fmt.Errorf("failed to confirm upload: %w", err)
	}

	// Generate download URL for response
	urlPtr := s.GenerateResponsePresignedURL(&updatedAttachment.File, ctx)
	if urlPtr == nil {
		s.Logger.LogError(ctx, "ConfirmUpload", "Failed to generate download URL", errors.New("presigned url is nil"), zap.String("file_id", req.FileID.String()))
		return nil, fmt.Errorf("failed to generate download URL")
	}

	return &ConfirmUploadResponse{
		FileURL:   *urlPtr,
		FileID:    updatedAttachment.Uuid,
		CreatedAt: updatedAttachment.Created.Time,
		Size:      size,
	}, nil
}

func (s *attachmentService) GetAttachmentById(ctx context.Context, id uuid.UUID) (*GetAttachmentByIdResponse, error) {
	attachment, err := s.Store.GetAttachmentById(ctx, id)
	if err != nil {
		s.Logger.LogError(ctx, "GetAttachmentById", "Failed to get attachment by ID", err, zap.String("file_id", id.String()))
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}
	url := s.GenerateResponsePresignedURL(&attachment.File, ctx)
	if url == nil {
		s.Logger.LogError(ctx, "GetAttachmentById", "Failed to generate presigned URL", errors.New("presigned url is nil"), zap.String("file_id", id.String()))
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
		s.Logger.LogError(ctx, "DeleteAttachment", "Failed to get attachment by ID", err, zap.String("file_id", id.String()))
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "DeleteAttachment", "Failed to begin transaction", err, zap.String("file_id", id.String()))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			s.Logger.LogError(ctx, "DeleteAttachment", "Failed to rollback transaction", rollbackErr, zap.String("file_id", id.String()))
		}
	}()
	qtx := s.Store.WithTx(tx)

	attachment, err = qtx.DeleteAttachment(ctx, id)
	if err != nil {
		s.Logger.LogError(ctx, "DeleteAttachment", "Failed to delete attachment record", err, zap.String("file_id", id.String()))
		return nil, fmt.Errorf("failed to delete attachment record: %w", err)
	}

	err = s.B2Client.Delete(ctx, attachment.File)
	if err != nil {
		s.Logger.LogError(ctx, "DeleteAttachment", "Failed to delete attachment from B2", err, zap.String("file_id", id.String()))
		return nil, fmt.Errorf("failed to delete attachment from B2: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "DeleteAttachment", "Failed to commit transaction", err, zap.String("file_id", id.String()))
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

	cleanFilename := util.SanitizeFilename(filename, 100)

	fileUUID := uuid.New()
	key := fmt.Sprintf("%s/%d/%02d/%s_%s",
		string(category),
		now.Year(),
		now.Month(),
		fileUUID,
		cleanFilename,
	)
	return key, fileUUID
}

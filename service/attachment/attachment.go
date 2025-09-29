package attachment

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *attachmentService) UploadAttachment(ctx context.Context,
	file multipart.File,
	header *multipart.FileHeader,
) (*UploadHandlerResponse, error) {
	fileInfo, err := s.validateAndFinalizeFile(file, header)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UploadAttachment", "File validation failed", zap.Error(err))
		return nil, err
	}

	key, uuid := s.generateSecureKey(header.Filename, fileInfo)

	objectKey, size, err := s.B2Client.Upload(ctx, file, key, fileInfo.ContentType)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UploadAttachment", "File upload failed", zap.Error(err))
		return nil, fmt.Errorf("file upload error: %v", err)
	}

	arg := db.CreateAttachmentParams{
		Uuid: uuid,
		File: objectKey,
		Size: int32(size),
		Tag:  &header.Filename,
	}
	attachment, err := s.Store.CreateAttachment(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UploadAttachment", "Failed to create attachment record", zap.Error(err))
		return nil, fmt.Errorf("failed to create attachment record: %v", err)
	}

	return &UploadHandlerResponse{
		FileURL:   key,
		FileID:    attachment.Uuid,
		CreatedAt: attachment.Created.Time,
		Size:      int64(attachment.Size),
	}, nil
}

func (s *attachmentService) GetAttachmentById(ctx context.Context, id uuid.UUID) (*GetAttachmentByIdResponse, error) {
	attachment, err := s.Store.GetAttachmentById(ctx, id)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetAttachmentById", "Failed to get attachment by ID", zap.Error(err))
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
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteAttachment", "Failed to get attachment by ID", zap.Error(err))

		return nil, fmt.Errorf("failed to get attachment")
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteAttachment", "Failed to begin transaction", zap.Error(err))

		return nil, fmt.Errorf("failed to begin transaction")
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteAttachment", "Failed to rollback transaction", zap.Error(rollbackErr))
		}
	}()
	qtx := s.Store.WithTx(tx)

	attachment, err = qtx.DeleteAttachment(ctx, id)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteAttachment", "Failed to delete attachment record", zap.Error(err))

		return nil, fmt.Errorf("failed to delete attachment record")
	}

	err = s.B2Client.Delete(ctx, attachment.File)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteAttachment", "Failed to delete attachment from B2", zap.Error(err))
		return nil, fmt.Errorf("failed to delete attachment from B2: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteAttachment", "Failed to commit transaction", zap.Error(err))

		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &DeleteAttachmentResponse{
		FileURL:   attachment.File,
		FileID:    attachment.Uuid,
		CreatedAt: attachment.Created.Time,
		Size:      int64(attachment.Size),
	}, nil
}

func (s *attachmentService) validateAndFinalizeFile(file multipart.File, header *multipart.FileHeader) (*FileInfo, error) {
	if header.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum limit of 100MB")
	}

	if header.Size == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	buffer := make([]byte, min(512, header.Size))
	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error resetting file: %v", err)
	}

	detectedType := http.DetectContentType(buffer[:bytesRead])

	if detectedType == "application/octet-stream" {
		if extType := mime.TypeByExtension(filepath.Ext(header.Filename)); extType != "" {
			detectedType = extType
		}
	}

	category, allowed := allowedMimeTypes[detectedType]
	if !allowed {
		return nil, fmt.Errorf("unsupported file type: %s", detectedType)
	}

	checksum, md5Hash, err := s.generateChecksums(file)
	if err != nil {
		return nil, fmt.Errorf("error generating file checksums: %v", err)
	}

	return &FileInfo{
		Size:        header.Size,
		ContentType: detectedType,
		Checksum:    checksum,
		MD5Hash:     md5Hash,
		Category:    category,
		Extension:   filepath.Ext(header.Filename),
	}, nil
}

func (s *attachmentService) generateChecksums(file multipart.File) (string, string, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return "", "", fmt.Errorf("error resetting file for checksum calculation: %v", err)
	}

	hasherSHA256 := sha256.New()
	hasherMD5 := md5.New()

	multiWriter := io.MultiWriter(hasherSHA256, hasherMD5)
	if _, err := io.Copy(multiWriter, file); err != nil {
		return "", "", fmt.Errorf("error calculating file checksums: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", "", fmt.Errorf("error resetting file after checksum calculation: %v", err)
	}

	sha256hash := hex.EncodeToString(hasherSHA256.Sum(nil))
	md5hash := hex.EncodeToString(hasherMD5.Sum(nil))

	return sha256hash, md5hash, nil
}

func (s *attachmentService) generateSecureKey(filename string, fileInfo *FileInfo) (string, uuid.UUID) {
	now := time.Now().UTC()

	cleanFilename := s.sanitizeFilename(filename)

	uuid := uuid.New()
	key := fmt.Sprintf("%s/%d/%02d/%s_%s",
		string(fileInfo.Category),
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

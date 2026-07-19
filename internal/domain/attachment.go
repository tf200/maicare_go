package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const MaxAttachmentSize int64 = 100 << 20

type AttachmentCategory string

const (
	AttachmentImage    AttachmentCategory = "image"
	AttachmentDocument AttachmentCategory = "document"
	AttachmentArchive  AttachmentCategory = "archive"
)

type InitAttachmentUploadParams struct {
	Filename    string
	ContentType string
	Size        int64
}

type InitAttachmentUploadResult struct {
	UploadURL string    `json:"upload_url"`
	FileID    uuid.UUID `json:"file_id"`
	Key       string    `json:"key"`
}

type ConfirmAttachmentUploadResult struct {
	FileURL   string    `json:"file_url"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

type AttachmentResult struct {
	FileURL   string    `json:"file_url"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

type AttachmentRepository interface {
	CreateAttachment(ctx context.Context, params CreateAttachmentParams) (*AttachmentFile, error)
	GetAttachmentByID(ctx context.Context, id uuid.UUID) (*AttachmentFile, error)
	SetAttachmentUsed(ctx context.Context, id uuid.UUID, used bool) (*AttachmentFile, error)
	DeleteAttachment(ctx context.Context, id uuid.UUID) (*AttachmentFile, error)
}

type CreateAttachmentParams struct {
	ID   uuid.UUID
	Name string
	File string
	Size int32
	Tag  *string
}

type AttachmentService interface {
	InitUpload(ctx context.Context, params InitAttachmentUploadParams) (*InitAttachmentUploadResult, error)
	ConfirmUpload(ctx context.Context, id uuid.UUID) (*ConfirmAttachmentUploadResult, error)
	GetAttachment(ctx context.Context, id uuid.UUID) (*AttachmentResult, error)
	DeleteAttachment(ctx context.Context, id uuid.UUID) (*AttachmentResult, error)
}

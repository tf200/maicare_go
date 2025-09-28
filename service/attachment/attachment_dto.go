package attachment

import (
	"time"

	"github.com/google/uuid"
)

const (
	MaxFileSize     = 100 << 20 // 100MB - increased from 10MB
	MinPartSize     = 5 << 20   // 5MB - minimum part size for multipart upload
	MaxMemoryBuffer = 32 << 20  // 32MB - max memory buffer before streaming
	UploadTimeout   = 30 * time.Minute
	ChecksumTimeout = 5 * time.Minute
)

var allowedMimeTypes = map[string]FileCategory{
	// Images
	"image/jpeg":    ImageFile,
	"image/jpg":     ImageFile,
	"image/png":     ImageFile,
	"image/gif":     ImageFile,
	"image/webp":    ImageFile,
	"image/svg+xml": ImageFile,

	// Documents
	"application/pdf":    DocumentFile,
	"application/msword": DocumentFile,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": DocumentFile,
	"application/vnd.ms-excel": DocumentFile,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": DocumentFile,
	"text/plain": DocumentFile,
	"text/csv":   DocumentFile,

	// Archives
	"application/zip":             ArchiveFile,
	"application/x-rar":           ArchiveFile,
	"application/x-7z-compressed": ArchiveFile,
}

type FileCategory string

const (
	ImageFile    FileCategory = "image"
	DocumentFile FileCategory = "document"
	ArchiveFile  FileCategory = "archive"
)

type FileInfo struct {
	Size        int64
	ContentType string
	Category    FileCategory
	Extension   string
	Checksum    string
	MD5Hash     string
}

// UploadHandler handles file uploads
type UploadHandlerResponse struct {
	FileURL   string    `json:"file_url"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

// GetAttachmentByIdResponse represents the response for GetAttachmentByIdApi
type GetAttachmentByIdResponse struct {
	FileURL   string    `json:"file_url"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

// DeleteAttachmentResponse represents the response for DeleteAttachment
type DeleteAttachmentResponse struct {
	FileURL   string    `json:"file_url"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

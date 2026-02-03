package attachment

import (
	"context"

	"maicare_go/service/deps"

	"github.com/google/uuid"
)

type AttachmentService interface {
	InitUpload(ctx context.Context, req *InitUploadRequest) (*InitUploadResponse, error)
	ConfirmUpload(ctx context.Context, req *ConfirmUploadRequest) (*ConfirmUploadResponse, error)
	GetAttachmentById(ctx context.Context, id uuid.UUID) (*GetAttachmentByIdResponse, error)
	DeleteAttachment(ctx context.Context, id uuid.UUID) (*DeleteAttachmentResponse, error)
}

type attachmentService struct {
	*deps.ServiceDependencies
}

func NewAttachmentService(deps *deps.ServiceDependencies) AttachmentService {
	return &attachmentService{
		ServiceDependencies: deps,
	}
}

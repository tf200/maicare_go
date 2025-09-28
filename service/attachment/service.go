package attachment

import (
	"context"
	"maicare_go/service/deps"
	"mime/multipart"

	"github.com/google/uuid"
)

type AttachmentService interface {
	UploadAttachment(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*UploadHandlerResponse, error)
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

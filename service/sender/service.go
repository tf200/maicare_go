package sender

import (
	"context"

	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SenderService interface {
	CreateSender(ctx context.Context, req *CreateSenderRequest) (*CreateSenderResponse, error)
	ListSenders(ctx *gin.Context, req *ListSendersRequest) (*pagination.Response[ListSendersResponse], error)
	GetSenderByID(ctx context.Context, senderID uuid.UUID) (*GetSenderByIdResponse, error)
	UpdateSender(ctx context.Context, senderID uuid.UUID, req *UpdateSenderRequest) (*UpdateSenderResponse, error)
	DeleteSender(ctx context.Context, senderID uuid.UUID) error
	CreateSenderInvoiceTemplate(ctx context.Context, senderID uuid.UUID, req *CreateSenderInvoiceTemplateRequest) error
}

type senderService struct {
	*deps.ServiceDependencies
}

func NewSenderService(deps *deps.ServiceDependencies) SenderService {
	return &senderService{
		ServiceDependencies: deps,
	}
}

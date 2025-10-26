package sender

import (
	"context"
	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
)

type SenderService interface {
	CreateSender(ctx context.Context, req *CreateSenderRequest) (*CreateSenderResponse, error)
	ListSenders(ctx *gin.Context, req *ListSendersRequest) (*pagination.Response[ListSendersResponse], error)
	GetSenderByID(ctx context.Context, senderID int64) (*GetSenderByIdResponse, error)
	UpdateSender(ctx context.Context, senderID int64, req *UpdateSenderRequest) (*UpdateSenderResponse, error)
	DeleteSender(ctx context.Context, senderID int64) error
	CreateSenderInvoiceTemplate(ctx context.Context, senderID int64, req *CreateSenderInvoiceTemplateRequest) error
}

type senderService struct {
	*deps.ServiceDependencies
}

func NewSenderService(deps *deps.ServiceDependencies) SenderService {
	return &senderService{
		ServiceDependencies: deps,
	}
}

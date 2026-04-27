package deps

import (
	"context"
	"time"

	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/internal/ws"
	"maicare_go/logger"
	"maicare_go/service/ai"
	"maicare_go/service/pdf"
	"maicare_go/token"
	"maicare_go/util"
)

type ServiceDependencies struct {
	Store      *db.Store
	TokenMaker token.Maker
	Logger     logger.Logger
	Config     *util.Config
	B2Client   bucket.ObjectStorageInterface
	WsHub      *ws.Hub
	PDFService pdf.PdfService
	AIService  ai.AIService
	GrpcClient grpclient.GrpcClientInterface
}

func NewServiceDependencies(store *db.Store, tokenMaker token.Maker, logger logger.Logger, config *util.Config, b2Client bucket.ObjectStorageInterface, grpcClient grpclient.GrpcClientInterface, wsHub *ws.Hub, aiService ai.AIService) *ServiceDependencies {
	return &ServiceDependencies{
		Store:      store,
		TokenMaker: tokenMaker,
		Logger:     logger,
		Config:     config,
		B2Client:   b2Client,
		WsHub:      wsHub,
		PDFService: pdf.NewPdfService(b2Client),
		AIService:  aiService,
		GrpcClient: grpcClient,
	}
}

func (d *ServiceDependencies) GenerateResponsePresignedURL(fileKey *string, ctx context.Context) *string {
	if fileKey == nil {
		return nil
	}

	url, err := d.B2Client.GeneratePresignedURL(ctx, *fileKey, time.Minute*15)
	if err != nil {
		return nil
	}
	return &url
}

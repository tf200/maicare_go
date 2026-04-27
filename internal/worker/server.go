package worker

import (
	"context"
	"crypto/tls"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	pkgasynq "maicare_go/pkg/asynq"
	"maicare_go/pkg/bucket"
	pkgemail "maicare_go/pkg/email"
	pkgpdf "maicare_go/pkg/pdf"

	hibikenasynq "github.com/hibiken/asynq"
)

// IncidentPDFService covers the incident PDF methods the worker needs.
// Implemented by *pkgpdf.pdfService (returned by pkgpdf.NewPdfService).
type IncidentPDFService interface {
	GenerateIncidentPDF(ctx context.Context, incidentData pkgpdf.IncidentReportData) ([]byte, error)
	GenerateAndUploadIncidentPDF(ctx context.Context, incidentData pkgpdf.IncidentReportData) (string, error)
}

type AsynqServer struct {
	notifSvc  domain.NotificationService
	pdfSvc    IncidentPDFService
	server    *hibikenasynq.Server
	store     *db.Store
	brevoConf *pkgemail.BrevoConf
	b2Bucket  bucket.ObjectStorageClient
}

func NewAsynqServer(
	redisHost,
	redisUser,
	redisPassword string,
	store *db.Store,
	tlsConfig *tls.Config,
	brevoConf *pkgemail.BrevoConf,
	b2Bucket bucket.ObjectStorageClient,
	notifSvc domain.NotificationService,
	pdfSvc IncidentPDFService,
) *AsynqServer {
	srv := hibikenasynq.NewServer(
		hibikenasynq.RedisClientOpt{
			Addr:         redisHost,
			Username:     redisUser,
			Password:     redisPassword,
			TLSConfig:    tlsConfig,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		hibikenasynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				pkgasynq.QueueCritical: 6,
				pkgasynq.QueueDefault:  3,
				pkgasynq.QueueLow:      1,
			},
			RetryDelayFunc: func(n int, err error, t *hibikenasynq.Task) time.Duration {
				return time.Duration(n*n) * time.Second
			},
		},
	)

	return &AsynqServer{
		server:    srv,
		store:     store,
		brevoConf: brevoConf,
		b2Bucket:  b2Bucket,
		notifSvc:  notifSvc,
		pdfSvc:    pdfSvc,
	}
}

func (a *AsynqServer) Start() error {
	mux := hibikenasynq.NewServeMux()
	mux.HandleFunc(pkgasynq.TypeEmailDelivery, a.ProcessEmailTask)
	mux.HandleFunc(pkgasynq.TypeIncidentProcess, a.ProcessIncidentTask)
	mux.HandleFunc(pkgasynq.TypeIncidentConfirmedEmail, a.ProcessIncidentConfirmedEmailTask)
	mux.HandleFunc(pkgasynq.TypeNotificationSend, a.ProcessNotificationTask)
	mux.HandleFunc(pkgasynq.TypeAcceptedRegistration, a.ProcessRegistrationFormTask)
	mux.HandleFunc(pkgasynq.TypeProcessRegistrationFormEmail, a.ProcessProcessRegistrationFormEmailTask)
	mux.HandleFunc(TypeContractReminder, a.ProcessContractRemiderTask)
	mux.HandleFunc(TypeClientCareStatusSync, a.ProcessClientCareStatusSyncTask)

	return a.server.Start(mux)
}

func (a *AsynqServer) Shutdown() {
	a.server.Shutdown()
}

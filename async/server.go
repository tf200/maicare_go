package async

import (
	"crypto/tls"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/email"
	"maicare_go/notification"
	"maicare_go/service"

	"time"

	"github.com/hibiken/asynq"
)

type AsynqServer struct {
	businessService     *service.BusinessService
	server              *asynq.Server
	store               *db.Store
	brevoConf           *email.BrevoConf
	b2Bucket            *bucket.ObjectStorageClient
	notificationService *notification.Service
}

func NewAsynqServer(redisHost, redisUser, redisPassword string,
	store *db.Store, tls *tls.Config,
	brevoConf *email.BrevoConf, b2Bucket *bucket.ObjectStorageClient,
	notificationService *notification.Service,
	businessService *service.BusinessService) *AsynqServer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:         redisHost,
			Username:     redisUser,
			Password:     redisPassword,
			TLSConfig:    tls,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			RetryDelayFunc: func(n int, err error, t *asynq.Task) time.Duration {
				return time.Duration(n*n) * time.Second // Exponential backoff
			},
		},
	)
	return &AsynqServer{server: srv,
		store:               store,
		brevoConf:           brevoConf,
		b2Bucket:            b2Bucket,
		notificationService: notificationService,
		businessService:     businessService}
}

func (a *AsynqServer) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, a.ProcessEmailTask)
	mux.HandleFunc(TypeIncidentProcess, a.ProcessIncidentTask)
	mux.HandleFunc(TypeNotificationSend, a.ProcessNotificationTask)
	mux.HandleFunc(TypeAppointmentCreate, a.ProcessAppointmentTask)
	mux.HandleFunc(TypeAcceptedRegistration, a.ProcessRegistrationFormTask)
	mux.HandleFunc(TypeContractReminder, a.ProcessContractRemiderTask)

	return a.server.Start(mux)
}

func (a *AsynqServer) Shutdown() {
	a.server.Shutdown()
}

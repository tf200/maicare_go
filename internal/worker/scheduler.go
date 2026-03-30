package worker

import (
	"crypto/tls"
	"log"
	"time"

	hibikenasynq "github.com/hibiken/asynq"
)

const (
	TypeContractReminder     = "contract:reminder"
	TypeClientCareStatusSync = "client:care_status_sync"
)

type Scheduler struct {
	Scheduler *hibikenasynq.Scheduler
}

func NewScheduler(redisHost, redisUser, redisPassword string, tlsConfig *tls.Config) *Scheduler {
	scheduler := hibikenasynq.NewScheduler(
		hibikenasynq.RedisClientOpt{
			Addr:         redisHost,
			Username:     redisUser,
			Password:     redisPassword,
			TLSConfig:    tlsConfig,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		nil,
	)

	return &Scheduler{Scheduler: scheduler}
}

func (s *Scheduler) ScheduleContractReminder() error {
	task := hibikenasynq.NewTask(TypeContractReminder, nil)

	entryID, err := s.Scheduler.Register("0 0 * * *", task)
	if err != nil {
		return err
	}

	log.Printf("Scheduled contract reminder with entry ID: %s", entryID)
	return nil
}

func (s *Scheduler) ScheduleClientCareStatusSync() error {
	task := hibikenasynq.NewTask(TypeClientCareStatusSync, nil)

	entryID, err := s.Scheduler.Register("0 * * * *", task)
	if err != nil {
		return err
	}

	log.Printf("Scheduled client care status sync with entry ID: %s", entryID)
	return nil
}

func (s *Scheduler) Start() error {
	if err := s.ScheduleContractReminder(); err != nil {
		return err
	}
	if err := s.ScheduleClientCareStatusSync(); err != nil {
		return err
	}

	if err := s.Scheduler.Run(); err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) Shutdown() {
	if s == nil || s.Scheduler == nil {
		return
	}

	s.Scheduler.Shutdown()
}

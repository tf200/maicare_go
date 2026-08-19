package worker

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"maicare_go/internal/ctxkeys"
	pkgasynq "maicare_go/pkg/asynq"

	"github.com/goccy/go-json"
	hibikenasynq "github.com/hibiken/asynq"
)

const (
	TypeContractReminder     = "contract:reminder"
	TypeClientCareStatusSync = "client:care_status_sync"
)

type Scheduler struct {
	Scheduler *hibikenasynq.Scheduler
	actor     ctxkeys.ActorIdentity
}

func NewScheduler(redisHost, redisUser, redisPassword string, tlsConfig *tls.Config, actor ctxkeys.ActorIdentity) *Scheduler {
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

	return &Scheduler{Scheduler: scheduler, actor: actor}
}

func (s *Scheduler) task(taskType string) (*hibikenasynq.Task, error) {
	if !s.actor.IsValid() {
		return nil, fmt.Errorf("scheduled task %s requires a valid system actor", taskType)
	}
	payload, err := json.Marshal(pkgasynq.ScheduledWorkerPayload{Actor: pkgasynq.ActorPayload{
		UserID: s.actor.UserID, EmployeeID: s.actor.EmployeeID,
	}})
	if err != nil {
		return nil, err
	}
	return hibikenasynq.NewTask(taskType, payload), nil
}

func (s *Scheduler) ScheduleContractReminder() error {
	task, err := s.task(TypeContractReminder)
	if err != nil {
		return err
	}

	entryID, err := s.Scheduler.Register("0 0 * * *", task)
	if err != nil {
		return err
	}

	log.Printf("Scheduled contract reminder with entry ID: %s", entryID)
	return nil
}

func (s *Scheduler) ScheduleClientCareStatusSync() error {
	task, err := s.task(TypeClientCareStatusSync)
	if err != nil {
		return err
	}

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

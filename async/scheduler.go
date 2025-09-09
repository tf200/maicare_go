package async

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeContractReminder = "contract:reminder"
)

type Scheduler struct {
	Scheduler *asynq.Scheduler
}

func NewScheduler(redisHost, redisUser, redisPassword string, tls *tls.Config) *Scheduler {
	sch := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:         redisHost,
			Username:     redisUser,
			Password:     redisPassword,
			TLSConfig:    tls,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		}, nil)
	return &Scheduler{Scheduler: sch}

}

func (s *Scheduler) ScheduleContractReminder() error {
	task := asynq.NewTask(TypeContractReminder, nil)

	entryID, err := s.Scheduler.Register("0 0 * * *", task)
	if err != nil {
		return err
	}
	// Optional: log the entry ID for debugging
	log.Printf("Scheduled contract reminder with entry ID: %s", entryID)

	return nil
}

// func (s *)

func (s *Scheduler) Start() error {

	if err := s.ScheduleContractReminder(); err != nil {
		return err
	}

	if err := s.Scheduler.Run(); err != nil {
		return err
	}
	return nil
}

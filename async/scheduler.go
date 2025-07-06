package async

import (
	"crypto/tls"
	"time"

	"github.com/hibiken/asynq"
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

func (s *Scheduler) Start() error {
	if err := s.Scheduler.Run(); err != nil {
		return err
	}
	return nil
}

package storage

import (
	"context"
	"github.com/Duke1616/etools/job"
	"github.com/Duke1616/etools/job/storage/mysql"
)

type Storager interface {
	Preempt(ctx context.Context) (job.CronJob, error)
}

type storage struct {
	db mysql.GormStoragerCronJob
}

func (s *storage) Preempt(ctx context.Context) (job.CronJob, error) {
	j, err := s.db.Preempt(ctx)

	return job.CronJob{
		Id: j.Id,
	}, err
}

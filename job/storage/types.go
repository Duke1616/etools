package storage

import (
	"context"
	"github.com/Duke1616/etools/job"
	"github.com/Duke1616/etools/job/storage/mysql"
)

type Storager interface {
	Preempt(ctx context.Context) (job.CronJob, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
}

type storage struct {
	db mysql.GormStoragerCronJob
}

func NewStorage(db mysql.GormStoragerCronJob) Storager {
	return &storage{
		db: db,
	}
}

func (s *storage) Preempt(ctx context.Context) (job.CronJob, error) {
	j, err := s.db.Preempt(ctx)

	return job.CronJob{
		Id:         j.Id,
		Expression: j.Expression,
		Executor:   j.Executor,
		Name:       j.Name,
	}, err
}

func (s *storage) Release(ctx context.Context, jid int64) error {
	return s.db.Release(ctx, jid)
}

func (s *storage) UpdateUtime(ctx context.Context, id int64) error {
	return s.db.UpdateUtime(ctx, id)
}

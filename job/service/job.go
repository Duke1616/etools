package service

import (
	"context"
	"github.com/Duke1616/etools/job"
)

type Service interface {
	Preempt(ctx context.Context) (job.CronJob, error)
}

type service struct {
}

func (s *service) Preempt(ctx context.Context) (job.CronJob, error) {
	return job.CronJob{}, nil
}

package service

import (
	"context"
	"fmt"
	"github.com/Duke1616/etools/job"
	"github.com/Duke1616/etools/job/storage"
	"time"
)

type Service interface {
	Preempt(ctx context.Context) (job.CronJob, error)
}

type service struct {
	storage storage.Storager

	refreshInterval time.Duration
}

func NewService(storage storage.Storager) Service {
	return &service{
		storage:         storage,
		refreshInterval: time.Minute}
}

func (s *service) Preempt(ctx context.Context) (job.CronJob, error) {
	j, err := s.storage.Preempt(ctx)
	if err != nil {
		return job.CronJob{}, err
	}

	// 任务执行超时退出
	ticker := time.NewTicker(s.refreshInterval)
	go func() {
		for range ticker.C {
			s.refresh(j.Id)
		}
	}()
	j.CancelFunc = func() {
		ticker.Stop()
		iCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := s.storage.Release(iCtx, j.Id)
		if er != nil {
			fmt.Print("释放 job 失败")
		}
	}
	return j, err
}

func (s *service) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := s.storage.UpdateUtime(ctx, id)
	if err != nil {
		fmt.Print("续约失败")
	}
}

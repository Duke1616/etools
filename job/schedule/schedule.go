package schedule

import (
	"context"
	"github.com/Duke1616/etools/job/executor"
	"github.com/Duke1616/etools/job/service"
	"time"
)

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()

		if err != nil {
			continue
		}

		// 获取到了可执行 j
		exec, ok := s.executors[j.Executor]
		if !ok {
			continue
		}

		go func() {
			defer func() {
				j.CancelFunc()
			}()
			er := exec.Exec(ctx, j)
			if er != nil {
				return
			}
		}()

	}
}

func NewScheduler(svc service.Service) *Scheduler {
	return &Scheduler{
		svc:       svc,
		dbTimeout: time.Second,
		executors: map[string]executor.Executor{},
	}
}

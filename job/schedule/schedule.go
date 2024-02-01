package schedule

import (
	"context"
	"fmt"
	"github.com/Duke1616/etools/job/executor"
	"github.com/Duke1616/etools/job/service"
	"golang.org/x/sync/semaphore"
	"time"
)

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			return err
		}

		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()

		if err != nil {
			continue
		}

		// 获取到了可执行 j
		exec, ok := s.executors[j.Executor]
		if !ok {
			fmt.Println("找不到执行器")
			continue
		}

		go func() {
			defer func() {
				s.limiter.Release(1)
				j.CancelFunc()
			}()
			er := exec.Exec(ctx, j)
			if er != nil {
				return
			}
			er = s.svc.ResetNextTime(ctx, j)
			if er != nil {
				fmt.Println("重置下次执行时间失败")
			}
		}()

	}
}

func (s *Scheduler) RegisterExecutor(exec executor.Executor) {
	s.executors[exec.Name()] = exec
}

func NewScheduler(svc service.Service) *Scheduler {
	return &Scheduler{
		svc:       svc,
		dbTimeout: time.Second,
		limiter:   semaphore.NewWeighted(100),
		executors: map[string]executor.Executor{},
	}
}

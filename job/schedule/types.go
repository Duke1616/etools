package schedule

import (
	"github.com/Duke1616/etools/job/executor"
	"github.com/Duke1616/etools/job/service"
	"golang.org/x/sync/semaphore"
	"time"
)

type Scheduler struct {
	dbTimeout time.Duration

	svc       service.Service
	executors map[string]executor.Executor

	limiter *semaphore.Weighted
}

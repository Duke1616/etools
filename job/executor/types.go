package executor

import (
	"context"
	"github.com/Duke1616/etools/job"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, j job.CronJob) error
}

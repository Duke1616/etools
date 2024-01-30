package executor

import (
	"context"
	"fmt"
	"github.com/Duke1616/etools/job"
)

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j job.CronJob) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: make(map[string]func(ctx context.Context, j job.CronJob) error),
	}
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, j job.CronJob) error {
	fn, ok := l.funcs[j.Name]
	if !ok {
		return fmt.Errorf("未注册本地方法 %s", j.Name)
	}
	return fn(ctx, j)
}

func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, j job.CronJob) error) {
	l.funcs[name] = fn
}

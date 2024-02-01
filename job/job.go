package job

import (
	"github.com/robfig/cron/v3"
	"time"
)

type CronJob struct {
	Id   int64
	Name string

	// Cron 表达式
	Expression string
	Executor   string

	CancelFunc func()
}

func (j CronJob) NextTime() time.Time {
	c := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := c.Parse(j.Expression)
	return s.Next(time.Now())
}

package job

type CronJob struct {
	Id   int64
	Name string

	// Cron 表达式
	Expression string
	Executor   string

	CancelFunc func()
}

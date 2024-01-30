package schedule

import (
	"github.com/Duke1616/etools/job/service"
	"time"
)

type Scheduler struct {
	dbTimeout time.Duration

	svc service.Service
}

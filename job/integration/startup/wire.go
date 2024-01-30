//go:build wireinject

package startup

import (
	"github.com/Duke1616/etools/job/schedule"
	"github.com/Duke1616/etools/job/service"
	"github.com/Duke1616/etools/job/storage"
	"github.com/Duke1616/etools/job/storage/mysql"
	"github.com/google/wire"
)

func InitServer() *schedule.Scheduler {
	wire.Build(
		InitDB,
		service.NewService,
		storage.NewStorage,
		mysql.NewGormJob,
		schedule.NewScheduler,
	)

	return &schedule.Scheduler{}
}

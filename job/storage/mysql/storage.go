package mysql

import (
	"context"
	"github.com/Duke1616/etools/job"
	"gorm.io/gorm"
)

type GormStoragerCronJob interface {
	Preempt(ctx context.Context) (job.CronJob, error)
}

type gormJob struct {
	db *gorm.DB
}

func (g *gormJob) Preempt(ctx context.Context) (CronJob, error) {
	return CronJob{}, nil
}

type CronJob struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
}

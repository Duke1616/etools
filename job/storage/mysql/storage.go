package mysql

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type GormStoragerCronJob interface {
	Preempt(ctx context.Context) (CronJob, error)
	Release(ctx context.Context, jid int64) error
	UpdateUtime(ctx context.Context, jid int64) error
}

type gormJob struct {
	db *gorm.DB
}

func NewGormJob(db *gorm.DB) GormStoragerCronJob {
	return &gormJob{
		db: db,
	}
}

func (g *gormJob) Preempt(ctx context.Context) (CronJob, error) {
	db := g.db.WithContext(ctx)
	for {
		var j CronJob
		now := time.Now().UnixMilli()
		// 取出小于当前时间，可执行任务
		// TODO 这里是缺少找到续约失败的 JOB 出来执行
		err := db.Where("status = ? AND next_time <?",
			jobStatusWaiting, now).
			First(&j).Error
		if err != nil {
			return j, err
		}

		// 修改版本号，防止重复抢占任务
		res := db.WithContext(ctx).Model(&CronJob{}).
			Where("id = ? AND version = ?", j.Id, j.Version).
			Updates(map[string]any{
				"status":  jobStatusRunning,
				"version": j.Version + 1,
				"utime":   now,
			})
		if res.Error != nil {
			return CronJob{}, res.Error
		}

		// 没有抢到，重新进入循环
		if res.RowsAffected == 0 {
			continue
		}
		return j, err
	}
}

func (g *gormJob) Release(ctx context.Context, jid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&CronJob{}).
		Where("id = ?", jid).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  now,
	}).Error
}

func (g *gormJob) UpdateUtime(ctx context.Context, jid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&CronJob{}).
		Where("id = ?", jid).Updates(map[string]any{
		"utime": now,
	}).Error
}

type CronJob struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Name       string `gorm:"type:varchar(128);unique"`
	Executor   string
	Expression string

	Status  int
	Version int64

	NextTime int64 `gorm:"index"`
	Utime    int64
	Ctime    int64
}

const (
	// jobStatusWaiting 没人抢
	jobStatusWaiting = iota
	// jobStatusRunning 已经被人抢了
	jobStatusRunning
	// jobStatusPaused 不再需要调度了
	jobStatusPaused
)

package integration

import (
	"context"
	"github.com/Duke1616/etools/job"
	"github.com/Duke1616/etools/job/executor"
	"github.com/Duke1616/etools/job/integration/startup"
	"github.com/Duke1616/etools/job/schedule"
	"github.com/Duke1616/etools/job/storage/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type SchedulerTestSuite struct {
	suite.Suite
	scheduler *schedule.Scheduler
	db        *gorm.DB
}

func (s *SchedulerTestSuite) SetupSuite() {
	s.db = startup.InitDB()
	s.scheduler = startup.InitServer()
}

func (s *SchedulerTestSuite) TearDownSuite() {
	err := s.db.Exec("TRUNCATE TABLE `cron_jobs`").Error
	assert.NoError(s.T(), err)
}

// TestSchedule 测试调度
func (s *SchedulerTestSuite) TestSchedule() {
	testCases := []struct {
		name string

		// 准备数据
		before func(t *testing.T)
		// 验证数据
		after func(t *testing.T)

		// Start 运行时间
		interval time.Duration
		wantErr  error
		wantJob  *testJob
	}{
		{
			name: "测试JOB",
			before: func(t *testing.T) {
				// 在 DB 里面插入一个 JOB，等待被调度
				j := mysql.CronJob{
					Id:       1,
					Name:     "job",
					Executor: "local",
					// 每五秒执行一次
					Expression: "*/5 * * * * ?",
					NextTime:   time.Now().UnixMilli(),
					Ctime:      123456,
					Utime:      654321,
				}
				err := s.db.Create(&j).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var j mysql.CronJob
				err := s.db.Where("id=?", 1).First(&j).Error
				assert.NoError(t, err)
				// 应该大于当下
				assert.True(t, j.NextTime > time.Now().UnixMilli())
				j.NextTime = 0
				assert.True(t, j.Ctime > 0)
				j.Ctime = 0
				assert.True(t, j.Utime > 0)
				j.Utime = 0
				assert.Equal(t, mysql.CronJob{
					Id:       1,
					Name:     "job",
					Executor: "local",
					//Cfg:      "这是一个配置",
					// 每五秒执行一次
					Expression: "*/5 * * * * ?",
					// 必然还原回来了
					Status: 0,
					// 抢占会导致版本升高
					Version: 1,
				}, j)
			},
			wantErr: context.DeadlineExceeded,
			// 运行了一次
			wantJob: &testJob{cnt: 1},
			// 开始调度一秒钟
			interval: time.Second,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			exec := executor.NewLocalFuncExecutor()
			j := &testJob{}
			exec.RegisterFunc("job", j.Do)
			s.scheduler.RegisterExecutor(exec)
			ctx, cancel := context.WithTimeout(context.Background(), tc.interval)
			defer cancel()
			err := s.scheduler.Schedule(ctx)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantJob, j)
		})
	}
}

func TestScheduler(t *testing.T) {
	suite.Run(t, &SchedulerTestSuite{})
}

type testJob struct {
	cnt int
}

func (t *testJob) Do(ctx context.Context, j job.CronJob) error {
	t.cnt++
	return nil
}

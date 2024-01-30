package integration

import (
	"github.com/Duke1616/etools/job/schedule"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SchedulerTestSuite struct {
	suite.Suite
	scheduler *schedule.Scheduler
	db        *gorm.DB
}

package startup

import (
	mysqlModel "github.com/Duke1616/etools/job/storage/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/etools"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&mysqlModel.CronJob{})
	if err != nil {
		panic(err)
	}

	return db
}

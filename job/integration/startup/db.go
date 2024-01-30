package startup

import (
	mysqlModel "github.com/Duke1616/etools/job/storage/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/etools"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&mysqlModel.CronJob{})
	if err != nil {
		panic(err)
	}

	return db
}

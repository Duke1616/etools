package startup

import (
	"fmt"
	"github.com/Duke1616/etools/gormx/connpool"
	"github.com/Duke1616/etools/logger"
	"github.com/Duke1616/etools/migrator/example"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SrcDB *gorm.DB
type DstDB *gorm.DB

func InitSrcDB() SrcDB {
	return initDB("src")
}

func InitDstDB() DstDB {
	return initDB("dst")
}

func InitDoubleWritePool(src SrcDB, dst DstDB, l logger.Logger) *connpool.DoubleWritePool {
	return connpool.NewDoubleWritePool(src, dst, l)
}

func InitPoolDB(p *connpool.DoubleWritePool) *gorm.DB {
	doubleWrite, err := gorm.Open(mysql.New(mysql.Config{
		Conn: p,
	}))
	if err != nil {
		panic(err)
	}
	return doubleWrite
}

func initDB(key string) *gorm.DB {
	dsn := fmt.Sprintf("root:123456@tcp(localhost:3306)/demo_%s", key)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitTables(db *gorm.DB) error {
	// 严格来说，这个不是优秀实践
	return db.AutoMigrate(
		&example.User{},
	)
}

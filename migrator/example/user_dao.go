package example

import (
	"context"
	"github.com/Duke1616/etools/migrator"
	"gorm.io/gorm"
	"time"
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewGORMUserDAO(db *gorm.DB) *GORMUserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"index"`
	Password string
	Phone    string
	Ctime    int64
	Utime    int64
}

func (u User) ID() int64 {
	return u.Id
}

func (u User) CompareTo(dst migrator.Entity) bool {
	dstVal, ok := dst.(User)
	return ok && u == dstVal
}

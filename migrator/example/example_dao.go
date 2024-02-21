package example

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"index"`
	Password string
	Phone    string
	Ctime    int64
	Utime    int64
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

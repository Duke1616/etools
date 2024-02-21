package example

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type DoubleWriteTestSuite struct {
	suite.Suite
	src *gorm.DB
	dst *gorm.DB
	dao *DoubleWriteDAO
}

func (s *DoubleWriteTestSuite) SetupSuite() {
	t := s.T()
	src, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/demo"))
	require.NoError(t, err)
	err = src.AutoMigrate(&User{})
	require.NoError(t, err)
	dst, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/demo_dst"))
	require.NoError(t, err)
	err = dst.AutoMigrate(&User{})
	require.NoError(t, err)
	s.dao = NewDoubleWriteDAO(NewGORMUserDAO(src), NewGORMUserDAO(dst))
	s.src = src
	s.dst = dst
	s.dao.UpdatePattern(patternSrcFirst)
}

func (s *DoubleWriteTestSuite) TearDownTest() {
	s.src.Exec("TRUNCATE TABLE users")
	s.dst.Exec("TRUNCATE TABLE users")
}

func (s *DoubleWriteTestSuite) TestDoubleWriteTest() {
	t := s.T()
	user := User{
		Email:    "1234567890@163.com",
		Password: "1234567890",
		Phone:    "1234567890",
	}
	err := s.dao.Insert(context.Background(), user)
	assert.NoError(t, err)

	// 数据校验比对
	srcFirst := &User{}
	s.src.Where("id = ?", 1).First(&srcFirst)
	dstFirst := &User{}
	s.dst.Where("id = ?", 1).First(&dstFirst)
	resetTime(dstFirst)
	resetTime(srcFirst)
	assert.Equal(t, srcFirst, dstFirst)
}

func TestDoubleWrite(t *testing.T) {
	suite.Run(t, new(DoubleWriteTestSuite))
}

func resetTime(user *User) *User {
	user.Utime = 1708490429590
	user.Ctime = 1708490429590
	return user
}

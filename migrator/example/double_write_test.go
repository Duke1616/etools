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

	s.dao = NewDoubleWriteDAOV1(src, dst)
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

	err := s.dao.Insert(context.Background(), User{
		Email:    "1234567890@163.com",
		Password: "1234567890",
		Phone:    "1234567890",
	})
	assert.NoError(t, err)
}

func TestDoubleWrite(t *testing.T) {
	suite.Run(t, new(DoubleWriteTestSuite))
}

package connpool

import (
	"github.com/Duke1616/etools/logger"
	"github.com/Duke1616/etools/migrator/example"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type DoubleWriteTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DoubleWriteTestSuite) SetupSuite() {
	t := s.T()
	src, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/demo_src"))
	require.NoError(t, err)
	err = src.AutoMigrate(&example.User{})
	require.NoError(t, err)
	dst, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/demo_dst"))
	require.NoError(t, err)
	err = dst.AutoMigrate(&example.User{})
	require.NoError(t, err)
	doubleWrite, err := gorm.Open(mysql.New(mysql.Config{
		Conn: &DoubleWritePool{
			src:     src.ConnPool,
			dst:     dst.ConnPool,
			pattern: atomicx.NewValueOf(PatternSrcFirst),
			l:       logger.NewNopLogger(),
		},
	}))
	require.NoError(t, err)
	s.db = doubleWrite
}

func (s *DoubleWriteTestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE TABLE users")
}

// 集成测试，需要启动数据库
func (s *DoubleWriteTestSuite) TestDoubleWriteTest() {
	t := s.T()
	err := s.db.Create(&example.User{
		Email:    "17682333",
		Password: "1234567",
	}).Error
	assert.NoError(t, err)
	// 查询数据库就可以看到对应的数据
}

func (s *DoubleWriteTestSuite) TestDoubleWriteTransaction() {
	t := s.T()
	err := s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&example.User{
			Email:    "17682333",
			Password: "1234567",
		}).Error
	})
	require.NoError(t, err)
}

func TestDoubleWrite(t *testing.T) {
	suite.Run(t, new(DoubleWriteTestSuite))
}

package scheduler

import (
	"context"
	"fmt"
	"github.com/Duke1616/etools/gormx/connpool"
	"github.com/Duke1616/etools/httpx"
	"github.com/Duke1616/etools/httpx/ginx"
	"github.com/Duke1616/etools/logger"
	"github.com/Duke1616/etools/migrator"
	"github.com/Duke1616/etools/migrator/events"
	"github.com/Duke1616/etools/migrator/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync"
)

type Scheduler[T migrator.Entity] struct {
	lock sync.Mutex
	src  *gorm.DB
	dst  *gorm.DB
	pool *connpool.DoubleWritePool

	l          logger.Logger
	cancelFull func()
	cancelIncr func()

	pattern  string
	producer events.Producer
}

func NewScheduler[T migrator.Entity](l logger.Logger, src *gorm.DB, dst *gorm.DB, pool *connpool.DoubleWritePool,
	producer events.Producer) *Scheduler[T] {
	return &Scheduler[T]{
		l:       l,
		src:     src,
		dst:     dst,
		pattern: connpool.PatternSrcOnly,
		cancelFull: func() {
			// 初始的时候，啥也不用做
		},
		cancelIncr: func() {
			// 初始的时候，啥也不用做
		},
		pool:     pool,
		producer: producer,
	}
}

func (s *Scheduler[T]) RegisterRoutes(server *gin.RouterGroup) {
	server.POST("/src_only", ginx.Wrap(s.SrcOnly))
	server.POST("/src_first", ginx.Wrap(s.SrcFirst))
	server.POST("/dst_first", ginx.Wrap(s.DstFirst))
	server.POST("/dst_only", ginx.Wrap(s.DstOnly))
	server.POST("/full/start", ginx.Wrap(s.StartFullValidation))
	server.POST("/full/stop", ginx.Wrap(s.StopFullValidation))
}

func (s *Scheduler[T]) SrcOnly(c *gin.Context) (httpx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternSrcOnly
	_ = s.pool.UpdatePattern(connpool.PatternSrcOnly)
	return httpx.Result{
		Msg: "OK",
	}, nil
}

func (s *Scheduler[T]) SrcFirst(c *gin.Context) (httpx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternSrcFirst
	_ = s.pool.UpdatePattern(connpool.PatternSrcFirst)
	return httpx.Result{
		Msg: "OK",
	}, nil
}

func (s *Scheduler[T]) DstFirst(c *gin.Context) (httpx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternDstFirst
	_ = s.pool.UpdatePattern(connpool.PatternDstFirst)
	return httpx.Result{
		Msg: "OK",
	}, nil
}

func (s *Scheduler[T]) DstOnly(c *gin.Context) (httpx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternDstOnly
	_ = s.pool.UpdatePattern(connpool.PatternDstOnly)
	return httpx.Result{
		Msg: "OK",
	}, nil
}

func (s *Scheduler[T]) StopFullValidation(c *gin.Context) (httpx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cancelFull()
	return httpx.Result{
		Msg: "停止全量校验成功",
	}, nil
}

func (s *Scheduler[T]) StartFullValidation(c *gin.Context) (httpx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 取消上一次的
	cancel := s.cancelFull

	v, err := s.newValidator()
	if err != nil {
		return httpx.Result{}, err
	}
	var ctx context.Context
	ctx, s.cancelFull = context.WithCancel(context.Background())

	go func() {
		cancel()
		err = v.Validate(ctx)
		if err != nil {
			s.l.Warn("退出全量校验", logger.Error(err))
		}
	}()

	return httpx.Result{
		Msg: "启动全量校验成功",
	}, nil
}

func (s *Scheduler[T]) newValidator() (*validator.Validator[T], error) {
	switch s.pattern {
	case connpool.PatternDstOnly, connpool.PatternSrcOnly:
		return validator.NewValidator[T](s.src, s.dst, s.pattern, s.l, s.producer), nil
	case connpool.PatternSrcFirst, connpool.PatternDstFirst:
		return validator.NewValidator[T](s.dst, s.src, s.pattern, s.l, s.producer), nil
	default:
		return nil, fmt.Errorf("未知的 pattern %s", s.pattern)
	}
}

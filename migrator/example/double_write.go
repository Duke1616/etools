package example

import (
	"context"
	"errors"
	"github.com/Duke1616/etools/logger"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

var errUnknownPattern = errors.New("未知的双写模式")

const (
	patternDstOnly  = "DST_ONLY"
	patternSrcOnly  = "SRC_ONLY"
	patternDstFirst = "DST_FIRST"
	patternSrcFirst = "SRC_FIRST"
)

type DoubleWriteDAO struct {
	src     UserDAO
	dst     UserDAO
	pattern *atomicx.Value[string]
	l       logger.Logger
}

func (d *DoubleWriteDAO) UpdatePattern(pattern string) {
	d.pattern.Store(pattern)
}

func NewDoubleWriteDAO(src UserDAO, dst UserDAO) *DoubleWriteDAO {
	return &DoubleWriteDAO{
		src:     src,
		dst:     dst,
		pattern: atomicx.NewValueOf(patternSrcOnly),
	}
}

func NewDoubleWriteDAOV1(src *gorm.DB, dst *gorm.DB) *DoubleWriteDAO {
	return &DoubleWriteDAO{
		src:     NewGORMUserDAO(src),
		dst:     NewGORMUserDAO(dst),
		pattern: atomicx.NewValueOf(patternSrcOnly),
	}
}

func (d *DoubleWriteDAO) Insert(ctx context.Context, u User) error {
	switch d.pattern.Load() {
	case patternSrcOnly:
		return d.src.Insert(ctx, u)
	case patternSrcFirst:
		err := d.src.Insert(ctx, u)
		if err != nil {
			return err
		}
		err = d.dst.Insert(ctx, u)
		if err != nil {
			d.l.Error("dst 写入失败", logger.Error(err))
		}
		return nil
	case patternDstOnly:
		return d.dst.Insert(ctx, u)
	case patternDstFirst:
		err := d.dst.Insert(ctx, u)
		if err != nil {
			return err
		}
		err = d.src.Insert(ctx, u)
		if err != nil {
			d.l.Error("src 写入失败", logger.Error(err))
		}
		return nil
	default:
		return errUnknownPattern
	}
}

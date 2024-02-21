package validator

import (
	"context"
	"github.com/Duke1616/etools/logger"
	"github.com/Duke1616/etools/migrator"
	"github.com/Duke1616/etools/migrator/events"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	base      *gorm.DB
	target    *gorm.DB
	direction string
	batchSize int

	l        logger.Logger
	producer events.Producer

	fromBase func(ctx context.Context, offset int) (T, error)
}

func NewValidator[T migrator.Entity](base *gorm.DB, target *gorm.DB, direction string, l logger.Logger,
	p events.Producer) *Validator[T] {
	res := &Validator[T]{
		base:      base,
		target:    target,
		l:         l,
		producer:  p,
		direction: direction,
		batchSize: 100,
	}
	res.fromBase = res.fullFromBase
	return res
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		return v.validateBaseToTarget(ctx)
	})
	eg.Go(func() error {
		return v.validateTargetToBase(ctx)
	})
	return eg.Wait()
}

func (v *Validator[T]) validateBaseToTarget(ctx context.Context) error {
	offset := 1
	for {
		src, err := v.fullFromBase(ctx, offset)
		if err == context.DeadlineExceeded || err == context.Canceled {
			return nil
		}

		if err == gorm.ErrRecordNotFound {
			// TODO 增量模式 、没有数据了
			continue
		}

		if err != nil {
			v.l.Error("base -> target 查询 base 失败", logger.Error(err))
			offset++
			continue
		}

		var dst T
		err = v.target.Where("id = ?", src.ID()).First(&dst).Error
		switch err {
		case gorm.ErrRecordNotFound:
			v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
		case nil:
			//var srcAny any = src
			//if c, ok := srcAny.(interface {
			//	CompareTo(entity migrator.Entity) bool
			//}); ok {
			//	// 有，我就用它的
			//	if !c.CompareTo(dst) {
			//		v.notify(src.ID(), events.InconsistentEventTypeNEQ)
			//	}
			//} else {
			//	// 没有，我就用反射
			//	if !reflect.DeepEqual(src, dst) {
			//		v.notify(src.ID(), events.InconsistentEventTypeNEQ)
			//	}
			//}

			equal := src.CompareTo(dst)
			if !equal {
				// 要丢一条消息到 Kafka 上
				v.notify(src.ID(), events.InconsistentEventTypeNEQ)
			}
		default:
			v.l.Error("base -> target 查询 target 失败",
				logger.Int64("id", src.ID()),
				logger.Error(err))
		}
		offset++
	}
}

func (v *Validator[T]) validateTargetToBase(ctx context.Context) error {
	offset := 0
	for {
		var ts []T
		err := v.target.WithContext(ctx).Select("id").
			Order("id").Limit(v.batchSize).Offset(offset).Find(&ts).Error
		if err == context.DeadlineExceeded || err == context.Canceled {
			return nil
		}

		if err == gorm.ErrRecordNotFound || len(ts) == 0 {
			continue
		}

		if err != nil {
			v.l.Error("target => base 查询 target 失败", logger.Error(err))
			offset += len(ts)
			continue
		}

		var srcTs []T
		ids := slice.Map(ts, func(idx int, t T) int64 {
			return t.ID()
		})
		err = v.base.WithContext(ctx).Select("id").Where("id IN ?", ids).Find(&srcTs).Error
		if err == gorm.ErrRecordNotFound || len(srcTs) == 0 {
			// 都代表。base 里面一条对应的数据都没有
			v.notifyBaseMissing(ts)
			offset += len(ts)
			continue
		}
		if err != nil {
			v.l.Error("target => base 查询 base 失败", logger.Error(err))
			// 保守起见，我都认为 base 里面没有数据
			offset += len(ts)
			continue
		}
		// 找差集，diff 里面的，就是 target 有，但是 base 没有的
		diff := slice.DiffSetFunc(ts, srcTs, func(src, dst T) bool {
			return src.ID() == dst.ID()
		})

		v.notifyBaseMissing(diff)
		if len(ts) < v.batchSize {
			// 说明比对完成了
		}
		offset += len(ts)
	}
}

func (v *Validator[T]) fullFromBase(ctx context.Context, offset int) (T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src T
	err := v.base.WithContext(dbCtx).Order("id").
		Offset(offset).First(&src).Error
	return src, err
}

func (v *Validator[T]) notifyBaseMissing(ts []T) {
	for _, val := range ts {
		v.notify(val.ID(), events.InconsistentEventTypeBaseMissing)
	}
}

// notify 通知kafka
func (v *Validator[T]) notify(id int64, typ string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := v.producer.ProduceInconsistentEvent(ctx, events.InconsistentEvent{
		ID:        id,
		Type:      typ,
		Direction: v.direction,
	})
	if err != nil {
		v.l.Error("发送不一致消息失败",
			logger.Error(err),
			logger.String("type", typ),
			logger.Int64("id", id))
	}
}

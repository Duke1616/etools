package validator

import (
	"context"
	"github.com/Duke1616/etools/logger"
	"github.com/Duke1616/etools/migrator"
	"github.com/Duke1616/etools/migrator/events"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	base     *gorm.DB
	target   *gorm.DB
	l        logger.Logger
	producer events.Producer
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	return v.validateBaseToTarget(ctx)
}

func (v *Validator[T]) validateBaseToTarget(ctx context.Context) error {
	offset := 1
	for {
		src, err := v.fullFromBase(ctx, offset)
		if err == context.DeadlineExceeded || err == context.Canceled {
			return nil
		}

		if err == gorm.ErrRecordNotFound {
			// 没有数据了
			continue
		}

		if err != nil {
			// 查询出错了
			v.l.Error("base -> target 查询 base 失败", logger.Error(err))
			// 在这里，
			offset++
			continue
		}

		var dst T
		err = v.target.Where("id = ?", src.ID()).First(&dst).Error
		switch err {
		case gorm.ErrRecordNotFound:
			v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
		case nil:
		default:

		}

		offset++

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

func (v *Validator[T]) notify(id int64, typ string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := v.producer.ProduceInconsistentEvent(ctx, events.InconsistentEvent{
		ID:   id,
		Type: typ,
	})
	if err != nil {
		v.l.Error("发送不一致消息失败",
			logger.Error(err),
			logger.String("type", typ),
			logger.Int64("id", id))
	}
}

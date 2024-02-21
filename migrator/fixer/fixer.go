package fixer

import (
	"context"
	"github.com/Duke1616/etools/migrator"
	"github.com/Duke1616/etools/migrator/events"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OverrideFixer[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB

	columns []string
}

func (o OverrideFixer[T]) Fix(ctx context.Context, evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeNEQ, events.InconsistentEventTypeTargetMissing:
		var t T
		err := o.base.WithContext(ctx).Where("id=?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return o.target.Model(&t).Delete("id = ?", evt.ID).Error
		case nil:
			return o.target.WithContext(ctx).Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns(o.columns),
			}).Create(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return o.target.Model(new(T)).Delete("id = ?", evt.ID).Error
	}
	return nil
}

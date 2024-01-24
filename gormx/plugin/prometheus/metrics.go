package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type MetricsPlugin struct {
	vector *prometheus.SummaryVec
}

type gormHookFunc func(tx *gorm.DB)

type gormRegister interface {
	Register(name string, fn func(*gorm.DB)) error
}

func NewPlugin(opts prometheus.SummaryOpts, register *prometheus.Registry) *MetricsPlugin {
	vector := prometheus.NewSummaryVec(opts,
		[]string{"type", "table"})

	register.MustRegister(vector)
	return &MetricsPlugin{
		vector: vector,
	}
}

func (p *MetricsPlugin) Name() string {
	return "prometheus"
}

func (p *MetricsPlugin) Initialize(db *gorm.DB) error {
	cb := db.Callback()
	hooks := []struct {
		callback gormRegister
		hook     gormHookFunc
		name     string
	}{
		{cb.Create().Before("gorm:create"), p.Before(), "gorm_create_before"},
		{cb.Create().After("gorm:create"), p.After("CREATE"), "gorm_create_after"},

		{cb.Query().Before("gorm:query"), p.Before(), "gorm_query_before"},
		{cb.Query().After("gorm:query"), p.After("QUERY"), "gorm_query_after"},

		{cb.Delete().Before("gorm:delete"), p.Before(), "gorm_delete_before"},
		{cb.Delete().After("gorm:delete"), p.After("DELETE"), "gorm_delete_after"},

		{cb.Update().Before("gorm:update"), p.Before(), "gorm_update_before"},
		{cb.Update().After("gorm:update"), p.After("UPDATE"), "gorm_update_after"},

		{cb.Row().Before("gorm:row"), p.Before(), "gorm_row_before"},
		{cb.Row().After("gorm:row"), p.After("ROW"), "gorm_row_after"},

		{cb.Raw().Before("gorm:raw"), p.Before(), "gorm_raw_before"},
		{cb.Raw().After("gorm:raw"), p.After("RAW"), "gorm_rwa_after"},
	}

	var firstErr error

	for _, h := range hooks {
		if err := h.callback.Register(h.name, h.hook); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("callback register %s failed: %w", h.name, err)
		}
	}

	return firstErr
}

func (p *MetricsPlugin) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start := time.Now()
		db.Set("start_time", start)
	}
}

func (p *MetricsPlugin) After(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		start, ok := val.(time.Time)
		if ok {
			duration := time.Since(start).Milliseconds()
			p.vector.WithLabelValues(typ, db.Statement.Table).
				Observe(float64(duration))
		}
	}
}

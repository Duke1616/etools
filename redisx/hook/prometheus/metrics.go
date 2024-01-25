package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

type MetricsHook struct {
	vector *prometheus.SummaryVec
}

//go:generate mockgen -package=redismocks -destination=mocks/redis_cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
func NewMetricsHook(opts prometheus.SummaryOpts, register *prometheus.Registry) *MetricsHook {
	vector := prometheus.NewSummaryVec(opts,
		[]string{"cmd", "key_exist"})

	register.MustRegister(vector)
	return &MetricsHook{
		vector: vector,
	}
}

func (m *MetricsHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (m *MetricsHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		var err error
		defer func() {
			duration := time.Since(start).Milliseconds()
			keyExists := err == redis.Nil
			m.vector.WithLabelValues(cmd.Name(), strconv.FormatBool(keyExists)).
				Observe(float64(duration))
		}()
		err = next(ctx, cmd)
		return err
	}
}

func (m *MetricsHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

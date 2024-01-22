package victoriametrics

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
	"time"
)

type Builder struct {
	Namespace string
}

func NewBuilder(namespace string) *Builder {
	return &Builder{
		Namespace: namespace,
	}
}

var metricsMap = make(map[ginMetrics]*metrics.Summary)

type ginMetrics struct {
	method  string
	pattern string
	status  int
}

func NewSummaryExt(hm ginMetrics, namespace string) *metrics.Summary {
	metrics.ExposeMetadata(true)

	name := fmt.Sprintf("%s{pattern=\"%s\", method=\"%s\", status=\"%d\"}",
		namespace, hm.pattern, hm.method, hm.status)

	return metrics.NewSummaryExt(name, time.Millisecond*20, []float64{0.5, 0.75, 0.9, 0.99, 0.999})
}

func (b *Builder) BuildResponseTime() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			pattern := ctx.FullPath()
			if ctx.FullPath() == "" {
				pattern = "empty"
			}

			gm := ginMetrics{
				method:  ctx.Request.Method,
				pattern: pattern,
				status:  ctx.Writer.Status(),
			}

			vector, ok := metricsMap[gm]
			if !ok {
				vector = NewSummaryExt(gm, b.Namespace)
				metricsMap[gm] = vector
			}

			// 记录请求开始时间
			start := time.Now()

			// 计算并更新指标
			vector.UpdateDuration(start)
		}()

		ctx.Next()
	}
}

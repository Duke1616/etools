package prometheus

import (
	"context"
	_ "embed"
	redismock "github.com/Duke1616/etools/redisx/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPrometheusHook(t *testing.T) {
	testCases := []struct {
		name string

		wantCode string
	}{
		{
			name:     "验证成功",
			wantCode: "200 OK",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建客户端
			client := http.DefaultClient
			client.Timeout = 1 * time.Second

			// 注册prometheus
			reg := prometheus.NewRegistry()
			hook := NewMetricsHook(prometheus.SummaryOpts{
				Namespace: "mgr",
				Subsystem: "etools",
				Name:      "redis_db",
				Help:      "统计 REDIS 的数据库查询",
				ConstLabels: map[string]string{
					"instance_id": "my_instance",
				},
				Objectives: map[float64]float64{
					0.5:   0.01,
					0.75:  0.01,
					0.9:   0.01,
					0.99:  0.001,
					0.999: 0.0001,
				},
			}, reg)

			// 添加hook
			db, mock := redismock.NewClientMock()
			db.AddHook(hook)

			// 发送信息
			mock.MatchExpectationsInOrder(true)
			_ = db.Eval(context.TODO(), "key", []string{"field"}, time.Now().Unix())

			// 发送http请求
			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer backend.Close()

			resp, err := client.Get(backend.URL)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			// 判断信息
			assert.Equal(t, tc.wantCode, resp.Status)
			gather, err := reg.Gather()
			require.NoError(t, err)

			assert.Equal(t, len(gather), 1)
		})
	}
}

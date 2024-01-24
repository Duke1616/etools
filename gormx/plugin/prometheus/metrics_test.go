package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPrometheusPlugin(t *testing.T) {
	testCases := []struct {
		name string
		do   func(ctx context.Context, db *gorm.DB)

		wantCode string
	}{
		{
			name: "验证成功",
			do: func(ctx context.Context, db *gorm.DB) {
				err := db.Exec("CREATE TABLE foo (id int)").Error
				require.NoError(t, err)
				var num int
				param := 1
				err = db.WithContext(ctx).Table("foo").Select("id", param).Where("id = ?", param).Scan(&num).Error
				require.NoError(t, err)
			},
			wantCode: "200 OK",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建客户端
			client := http.DefaultClient
			client.Timeout = 1 * time.Second

			// sqlite数据
			db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
			require.NoError(t, err)

			// 注册prometheus
			reg := prometheus.NewRegistry()
			cb := NewPlugin(prometheus.SummaryOpts{
				Namespace: "mgr",
				Subsystem: "etools",
				Name:      "gorm_db",
				Help:      "统计 GORM 的数据库查询",
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

			// 添加插件
			err = db.Use(cb)
			require.NoError(t, err)

			// 运行SQL、执行plugin插件操作
			tc.do(context.Background(), db)

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

			assert.Equal(t, tc.wantCode, resp.Status)

			// 判断信息
			gather, err := reg.Gather()
			require.NoError(t, err)

			assert.Equal(t, len(gather), 1)
		})
	}
}

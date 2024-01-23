package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestPrometheusPlugin(t *testing.T) {
//	testCases := []struct {
//		name string
//		do   func(ctx context.Context, db *gorm.DB)
//
//		// 构造请求，预期中输入
//		reqBody string
//	}{
//		{
//			name: "验证成功",
//			do: func(ctx context.Context, db *gorm.DB) {
//				err := db.Exec("CREATE TABLE foo (id int)").Error
//				require.NoError(t, err)
//				var num int
//				param := 42
//				err = db.WithContext(ctx).Table("foo").Select("id", param).Where("id = ?", param).Scan(&num).Error
//				fmt.Print(num)
//				require.NoError(t, err)
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			//// 初始化http
//			//mux := http.NewServeMux()
//			//
//			//// 将 /metrics 路由与 Prometheus handler 关联
//			//mux.Handle("/metrics", promhttp.Handler())
//			//
//			//// 使用 httptest.NewServer 创建模拟 HTTP 服务器
//			//server := httptest.NewServer(mux)
//			//defer server.Close() // 在测试结束后关闭服务器
//			//mux.ServeHTTP(recorder, req)
//			// 初始化DB
//			db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
//			require.NoError(t, err)
//			//initPrometheus()
//			// 加载gorm plugin
//			cb := NewPlugin(prometheus.SummaryOpts{
//				Namespace: "mgr",
//				Subsystem: "etools",
//				Name:      "gorm_db",
//				Help:      "统计 GORM 的数据库查询",
//				ConstLabels: map[string]string{
//					"instance_id": "my_instance",
//				},
//				Objectives: map[float64]float64{
//					0.5:   0.01,
//					0.75:  0.01,
//					0.9:   0.01,
//					0.99:  0.001,
//					0.999: 0.0001,
//				},
//			})
//			err = db.Use(cb)
//			require.NoError(t, err)
//
//			tc.do(context.TODO(), db)
//
//			// 准备Req和记录的 recorder
//			promhttp.
//			recorder := httptest.NewRecorder()
//			req, err := http.NewRequest(http.MethodGet,
//				"/", bytes.NewReader([]byte(tc.reqBody)))
//			req.Header.Set("Content-Type", "application/json")
//
//			assert.NoError(t, err)
//			promhttp.Handler().ServeHTTP(recorder, req)
//
//		})
//	}
//}

func initPrometheus() {
	go func() {
		// 专门给 prometheus 用的端口
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

func HandlerFor(reg prometheus.Gatherer, opts promhttp.HandlerOpts) http.Handler {
	return promhttp.HandlerForTransactional(prometheus.ToTransactionalGatherer(reg), opts)
}

type blockingCollector struct {
	CollectStarted, Block chan struct{}
}

func (b blockingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy_desc", "not helpful", nil, nil)
}

func (b blockingCollector) Collect(ch chan<- prometheus.Metric) {
	select {
	case b.CollectStarted <- struct{}{}:
	default:
	}
	// Collects nothing, just waits for a channel receive.
	<-b.Block
}

func TestHandlerMaxRequestsInFlight(t *testing.T) {
	reg := prometheus.NewRegistry()
	handler := HandlerFor(reg, promhttp.HandlerOpts{MaxRequestsInFlight: 1})
	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	w3 := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Accept", "test/plain")

	c := blockingCollector{Block: make(chan struct{}), CollectStarted: make(chan struct{}, 1)}
	reg.MustRegister(c)

	rq1Done := make(chan struct{})
	go func() {
		handler.ServeHTTP(w1, request)
		close(rq1Done)
	}()
	<-c.CollectStarted

	handler.ServeHTTP(w2, request)

	if got, want := w2.Code, http.StatusServiceUnavailable; got != want {
		t.Errorf("got HTTP status code %d, want %d", got, want)
	}
	if got, want := w2.Body.String(), "Limit of concurrent requests reached (1), try again later.\n"; got != want {
		t.Errorf("got body %q, want %q", got, want)
	}

	close(c.Block)
	<-rq1Done

	handler.ServeHTTP(w3, request)

	if got, want := w3.Code, http.StatusOK; got != want {
		t.Errorf("got HTTP status code %d, want %d", got, want)
	}
}

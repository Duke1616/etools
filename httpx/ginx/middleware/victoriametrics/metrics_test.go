package victoriametrics

import (
	"github.com/Duke1616/etools/httpx"
	"github.com/Duke1616/etools/httpx/ginx"
	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestMetrics(t *testing.T) {
	// Handler
	h := &Handler{}

	// 准备服务器，注册路由
	server := gin.Default()

	b := NewBuilder("etools_gin_http_resp_time")
	middleware := b.BuildResponseTime()

	server.Use(middleware)

	initPrometheus()

	h.RegisterRoutes(server)

	server.Run(":8080")
}

func initPrometheus() {
	go func() {
		// 专门给 prometheus 用的端口
		http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
			metrics.WritePrometheus(w, true)
		})

		//http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

type Handler struct {
}

func (h *Handler) RegisterRoutes(server *gin.Engine) {
	server.GET("/example", ginx.WrapBody[User](h.Example))
}

func (h *Handler) Example(ctx *gin.Context, req User) (httpx.Result, error) {
	data := &User{
		Username: "张三",
		Email:    "123456@qq.com",
		Password: "123456",
	}

	return httpx.Success(data), nil
}

type User struct {
	Username string
	Email    string
	Password string
}

func (u *User) Desensitization() {
	u.Password = ""
}

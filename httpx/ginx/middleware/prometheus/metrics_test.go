package prometheus

import (
	"github.com/Duke1616/etools/httpx"
	"github.com/Duke1616/etools/httpx/ginx"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"testing"
)

func TestMetrics(t *testing.T) {
	// Handler
	h := &Handler{}

	// 准备服务器，注册路由
	server := gin.Default()

	pb := &Builder{
		Namespace: "mer",
		Subsystem: "etools",
		Name:      "gin_http",
		Help:      "统计 GIN 的HTTP接口数据",
	}

	middleware := pb.BuildResponseTime()

	server.Use(middleware)

	initPrometheus()

	h.RegisterRoutes(server)

	server.Run(":8080")
}

func initPrometheus() {
	go func() {
		// 专门给 prometheus 用的端口
		http.Handle("/metrics", promhttp.Handler())
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

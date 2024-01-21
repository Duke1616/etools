package ginx

import (
	"bytes"
	"encoding/json"
	"github.com/Duke1616/etools/httpx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrapBody(t *testing.T) {
	testCases := []struct {
		name string
		// 构造请求，预期中输入
		reqBody string

		// 预期中的输出
		wantCode int
		wantRes  httpx.Result
	}{
		{
			name:     "脱敏成功",
			wantCode: http.StatusOK,
			wantRes: httpx.Result{
				Code: 10000,
				Msg:  "OK",
				Data: map[string]interface{}{
					"Email":    "123456@qq.com",
					"Username": "张三",
					"Password": "",
				},
			},
			reqBody: `{}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Handler
			h := &Handler{}

			// 准备服务器，注册路由
			server := gin.Default()
			h.RegisterRoutes(server)

			// 准备Req和记录的 recorder
			req, err := http.NewRequest(http.MethodPost,
				"/example", bytes.NewReader([]byte(tc.reqBody)))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)

			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if recorder.Code != http.StatusOK {
				return
			}
			var res httpx.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)

		})
	}
}

type Handler struct {
}

func (h *Handler) RegisterRoutes(server *gin.Engine) {
	server.POST("/example", WrapBody[User](h.Example))
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

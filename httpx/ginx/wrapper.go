package ginx

import (
	"github.com/Duke1616/etools/httpx"
	"github.com/Duke1616/etools/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var L logger.Logger = logger.NewNopLogger()

func WrapBody[Req any](
	bizFn func(ctx *gin.Context, req Req) (httpx.Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			L.Error("输入错误", logger.Error(err))
			return
		}
		L.Debug("输入参数", logger.Field{Key: "req", Val: req})
		res, err := bizFn(ctx, req)
		if err != nil {
			L.Error("执行业务逻辑失败", logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

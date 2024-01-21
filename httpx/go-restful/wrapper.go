package restfulx

import (
	"github.com/Duke1616/etools/httpx"
	"github.com/Duke1616/etools/logger"
	"github.com/emicklei/go-restful/v3"
	"net/http"
)

var L logger.Logger = logger.NewNopLogger()

func WrapBody[Req any](
	bizFn func(r *restful.Request, w *restful.Response, req Req) (httpx.Result, error),
) restful.RouteFunction {
	return func(r *restful.Request, w *restful.Response) {
		var req Req
		if err := r.ReadEntity(&req); err != nil {
			L.Error("输入错误", logger.Error(err))
			return
		}
		L.Debug("输入参数", logger.Field{Key: "req", Val: req})

		res, err := bizFn(r, w, req)
		if err != nil {
			L.Error("执行业务逻辑失败", logger.Error(err))
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteAsJson(res)
	}
}

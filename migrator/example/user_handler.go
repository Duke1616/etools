package example

import (
	"github.com/Duke1616/etools/httpx"
	"github.com/Duke1616/etools/httpx/ginx"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	dao UserDAO
}

func NewUserHandler(dao UserDAO) *UserHandler {
	return &UserHandler{
		dao: dao,
	}
}

func (s *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/create", ginx.WrapBody[User](s.MockUsers))
}

func (s *UserHandler) MockUsers(ctx *gin.Context, req User) (httpx.Result, error) {
	err := s.dao.Insert(ctx, req)
	if err != nil {
		return httpx.Result{
			Msg: "创建用户失败",
		}, err
	}

	return httpx.Result{
		Msg: "创建用户成功",
	}, nil
}

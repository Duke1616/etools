package startup

import (
	"github.com/Duke1616/etools/gormx/connpool"
	"github.com/Duke1616/etools/httpx/ginx"
	"github.com/Duke1616/etools/logger"
	"github.com/Duke1616/etools/migrator/events"
	"github.com/Duke1616/etools/migrator/example"
	"github.com/Duke1616/etools/migrator/scheduler"
	"github.com/gin-gonic/gin"
)

func InitGinServer(l logger.Logger, src SrcDB, dst DstDB, pool *connpool.DoubleWritePool,
	producer events.Producer, userHdl *example.UserHandler) *ginx.Server {
	engine := gin.Default()

	sch := scheduler.NewScheduler[example.User](l, src, dst, pool, producer)
	sch.RegisterRoutes(engine)
	userHdl.RegisterRoutes(engine)
	return &ginx.Server{
		Engine: engine,
		Addr:   ":8082",
	}
}

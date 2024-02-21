package startup

import (
	"github.com/Duke1616/etools/gormx/connpool"
	"github.com/Duke1616/etools/httpx/ginx"
	"github.com/Duke1616/etools/logger"
	"github.com/Duke1616/etools/migrator/events"
	"github.com/Duke1616/etools/migrator/events/fixer"
	"github.com/Duke1616/etools/migrator/example"
	"github.com/Duke1616/etools/migrator/scheduler"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func InitGinServer(l logger.Logger, src SrcDB, dst DstDB, pool *connpool.DoubleWritePool,
	producer events.Producer) *ginx.Server {
	engine := gin.Default()
	group := engine.Group("/migrator")

	sch := scheduler.NewScheduler[example.User](l, src, dst, pool, producer)
	sch.RegisterRoutes(group)
	return &ginx.Server{
		Engine: engine,
		Addr:   ":8082",
	}
}

func InitFixerProducer(p sarama.SyncProducer) events.Producer {
	return events.NewSaramaProducer("inconsistent_users", p)
}

func InitFixerConsumer(client sarama.Client,
	l logger.Logger,
	src SrcDB,
	dst DstDB) *fixer.Consumer[example.User] {
	res, err := fixer.NewConsumer[example.User](client, l, src, dst, "inconsistent_users")
	if err != nil {
		panic(err)
	}
	return res
}

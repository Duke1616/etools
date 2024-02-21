//go:build wireinject

package startup

import (
	"github.com/Duke1616/etools/migrator/example"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	InitSrcDB,
	InitDstDB,
	InitDoubleWritePool,
	InitPoolDB,
	InitLogger,
	InitSaramaClient,
	InitSaramaSyncProducer,
)

var UsersSet = wire.NewSet(
	example.NewGORMUserDAO,
	example.NewUserHandler,
)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		UsersSet,
		InitFixerConsumer,
		InitFixerProducer,
		InitConsumers,
		InitGinServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}

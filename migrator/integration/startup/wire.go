//go:build wireinject

package startup

import (
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

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		InitFixerConsumer,
		InitFixerProducer,
		InitConsumers,
		InitGinServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}

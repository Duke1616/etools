// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/Duke1616/etools/migrator/example"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitApp() *App {
	client := InitSaramaClient()
	logger := InitLogger()
	srcDB := InitSrcDB()
	dstDB := InitDstDB()
	consumer := InitFixerConsumer(client, logger, srcDB, dstDB)
	v := InitConsumers(consumer)
	doubleWritePool := InitDoubleWritePool(srcDB, dstDB, logger)
	syncProducer := InitSaramaSyncProducer(client)
	producer := InitFixerProducer(syncProducer)
	db := InitPoolDB(doubleWritePool)
	userDAO := example.NewGORMUserDAO(db)
	userHandler := example.NewUserHandler(userDAO)
	server := InitGinServer(logger, srcDB, dstDB, doubleWritePool, producer, userHandler)
	app := &App{
		Consumers:   v,
		AdminServer: server,
	}
	return app
}

// wire.go:

var thirdPartySet = wire.NewSet(
	InitSrcDB,
	InitDstDB,
	InitDoubleWritePool,
	InitPoolDB,
	InitLogger,
	InitSaramaClient,
	InitSaramaSyncProducer,
)

var UsersSet = wire.NewSet(example.NewGORMUserDAO, example.NewUserHandler)

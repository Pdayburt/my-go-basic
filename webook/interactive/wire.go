//go:build wireinject

package main

import (
	"example.com/mod/webook/interactive/events"
	"example.com/mod/webook/interactive/ioc"
	"example.com/mod/webook/interactive/repository"
	"example.com/mod/webook/interactive/repository/cache"
	"example.com/mod/webook/interactive/repository/dao"
	"example.com/mod/webook/interactive/service"
	"example.com/mod/webook/interactive/webookgrpc"
	"github.com/google/wire"
)

var interactiveService = wire.NewSet(dao.NewInteractiveDao, cache.NewInteractiveCache,
	repository.NewInteractiveRepository, service.NewInteractiveService,
)

// ioc.InitKafka 暂时
var thirdProvider = wire.NewSet(ioc.InitRedis, ioc.InitDbB,
	ioc.InitKafka, ioc.NewConsumer,
	ioc.InitGRPCxServer)

func InitApp() *App {
	wire.Build(thirdProvider, interactiveService,
		webookgrpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		wire.Struct(new(App), "*"))

	return new(App)
}

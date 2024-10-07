//go:build wireinject

package main

import (
	"example.com/mod/webook/internal/events/article"
	ioc2 "example.com/mod/webook/internal/ioc"
	"example.com/mod/webook/internal/repository"
	"example.com/mod/webook/internal/repository/cache"
	"example.com/mod/webook/internal/repository/dao"
	"example.com/mod/webook/internal/service"
	"example.com/mod/webook/internal/web"
	"github.com/google/wire"
)

func InitWebServerByWire() *App {

	wire.Build(
		ioc2.InitDbB, ioc2.InitRedis,
		ioc2.InitSMService,
		ioc2.InitKafka,
		ioc2.NewSyncProducer,

		article.NewInteractiveReadEventConsumer,
		article.NewKafkaProducer,
		ioc2.NewConsumer,

		//dao.NewUserDao,
		dao.NewUserDao, dao.NewArticleDao, dao.NewInteractiveDao,
		cache.NewUserCache, cache.NewCodeCache, cache.NewInteractiveCache,
		repository.NewCodeRepository, repository.NewUserRepository,
		repository.NewArticleRepository,
		repository.NewInteractiveRepository,

		service.NewUserService, service.NewCodeService,
		service.NewArticleService,
		service.NewInteractiveService,

		web.NewUserHandler,
		web.NewArticleHandler,

		//gin.Default,
		//还需要中间件和注册路由
		ioc2.InitMiddlewares, ioc2.InitGin,
		wire.Struct(new(App), "*"),
	)

	return new(App)
}

//go:build wireinject

package main

import (
	ioc2 "example.com/mod/webook/internal/ioc"
	"example.com/mod/webook/internal/repository"
	"example.com/mod/webook/internal/repository/cache"
	"example.com/mod/webook/internal/repository/dao"
	"example.com/mod/webook/internal/service"
	"example.com/mod/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServerByWire() *gin.Engine {

	wire.Build(
		ioc2.InitDbB, ioc2.InitRedis,
		ioc2.InitSMService,

		dao.NewUserDao,
		cache.NewUserCache, cache.NewCodeCache,
		repository.NewCodeRepository, repository.NewUserRepository,

		service.NewUserService, service.NewCodeService,
		service.NewArticleService,

		web.NewUserHandler,
		web.NewArticleHandler,

		//gin.Default,
		//还需要中间件和注册路由
		ioc2.InitMiddlewares, ioc2.InitGin,
	)

	return new(gin.Engine)
}

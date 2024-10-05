//go:build wireinject

package startup

import (
	"example.com/mod/webook/internal/ioc"
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
		ioc.InitDbB, ioc.InitRedis,
		ioc.InitSMService,

		dao.NewUserDao, dao.NewArticleDao,
		cache.NewUserCache, cache.NewCodeCache,
		repository.NewCodeRepository, repository.NewUserRepository,
		repository.NewArticleRepository,

		service.NewUserService, service.NewCodeService,
		service.NewArticleService,

		web.NewUserHandler,
		web.NewArticleHandler,

		//gin.Default,
		//还需要中间件和注册路由
		ioc.InitMiddlewares, ioc.InitGin,
	)

	return new(gin.Engine)
}

var thirdPart = wire.NewSet(
	ioc.InitDbB, ioc.InitRedis,
	ioc.InitSMService,
)

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		ioc.InitDbB,
		dao.NewArticleDao,
		repository.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}

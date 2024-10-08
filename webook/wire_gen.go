// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"example.com/mod/webook/internal/events/article"
	"example.com/mod/webook/internal/ioc"
	"example.com/mod/webook/internal/repository"
	"example.com/mod/webook/internal/repository/cache"
	"example.com/mod/webook/internal/repository/dao"
	"example.com/mod/webook/internal/service"
	"example.com/mod/webook/internal/web"
	"github.com/google/wire"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebServerByWire() *App {
	v := ioc.InitMiddlewares()
	db := ioc.InitDbB()
	userDao := dao.NewUserDao(db)
	cmdable := ioc.InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDao, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMService()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService)
	articleDao := dao.NewArticleDao(db)
	articleRepository := repository.NewArticleRepository(articleDao)
	client := ioc.InitKafka()
	syncProducer := ioc.NewSyncProducer(client)
	producer := article.NewKafkaProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, producer)
	interactiveDao := dao.NewInteractiveDao(db)
	interactiveCache := cache.NewInteractiveCache(cmdable)
	interactiveRepository := repository.NewInteractiveRepository(interactiveDao, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	articleHandler := web.NewArticleHandler(articleService, interactiveService)
	engine := ioc.InitGin(v, userHandler, articleHandler)
	interactiveReadEventConsumer := article.NewInteractiveReadEventConsumer(client, interactiveRepository)
	v2 := ioc.NewConsumer(interactiveReadEventConsumer)
	rankingService := service.NewBatchRankingService(articleService, interactiveService)
	job := ioc.InitRankingJob(rankingService)
	rankingJobAdapter := ioc.InitRankingJobAdapter(job)
	cron := ioc.InitJobs(rankingJobAdapter)
	app := &App{
		server:   engine,
		consumer: v2,
		cron:     cron,
	}
	return app
}

// wire.go:

var rankingServiceSet = wire.NewSet(service.NewBatchRankingService)

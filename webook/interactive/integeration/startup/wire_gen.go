// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"example.com/mod/webook/api/proto/gen/intr/v1"
	"example.com/mod/webook/interactive/repository"
	"example.com/mod/webook/interactive/repository/cache"
	"example.com/mod/webook/interactive/repository/dao"
	"example.com/mod/webook/interactive/service"
	"example.com/mod/webook/interactive/webookgrpc"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitInteractiveService() service.InteractiveService {
	gormDB := InitTestDB()
	interactiveDao := dao.NewInteractiveDao(gormDB)
	cmdable := InitRedis()
	interactiveCache := cache.NewInteractiveCache(cmdable)
	interactiveRepository := repository.NewInteractiveRepository(interactiveDao, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	return interactiveService
}

func InitInteractiveGRPCServer() intrv1.InteractiveServiceServer {
	gormDB := InitTestDB()
	interactiveDao := dao.NewInteractiveDao(gormDB)
	cmdable := InitRedis()
	interactiveCache := cache.NewInteractiveCache(cmdable)
	interactiveRepository := repository.NewInteractiveRepository(interactiveDao, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	return interactiveServiceServer
}

// wire.go:

var thirdProvider = wire.NewSet(InitRedis,
	InitTestDB, InitKafka)

var interactiveServer = wire.NewSet(cache.NewInteractiveCache, dao.NewInteractiveDao, repository.NewInteractiveRepository, service.NewInteractiveService)

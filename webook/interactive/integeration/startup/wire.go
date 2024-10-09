//go:build wireinject

package startup

import (
	intrv1 "example.com/mod/webook/api/proto/gen/intr/v1"
	"example.com/mod/webook/interactive/repository"
	"example.com/mod/webook/interactive/repository/cache"
	"example.com/mod/webook/interactive/repository/dao"
	"example.com/mod/webook/interactive/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis,
	InitTestDB, InitKafka)

var interactiveServer = wire.NewSet(cache.NewInteractiveCache, dao.NewInteractiveDao,
	repository.NewInteractiveRepository, service.NewInteractiveService)

func InitInteractiveService() service.InteractiveService {
	wire.Build(thirdProvider, interactiveServer)
	return service.NewInteractiveService(nil)
}

func InitInteractiveGRPCServer() intrv1.InteractiveServiceServer {
	wire.Build(thirdProvider, interactiveServer, grpc.NewInteractiveServiceServer)
	return grpc.NewInteractiveServiceServer(nil)
}

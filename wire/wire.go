//go:build wireinject

package wire

import (
	"example.com/mod/wire/repository"
	"example.com/mod/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	//传入各个组件的初始化方法
	wire.Build(repository.NewUserRepository, dao.NewUserDao, InitDB)
	return new(repository.UserRepository)
}

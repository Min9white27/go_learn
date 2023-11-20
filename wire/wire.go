//go:build wireinject

//让 wire 来注入这里的代码

package wire

import (
	"gitee.com/geekbang/basic-go/wire/repository"
	"gitee.com/geekbang/basic-go/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	// 只需要在这里声明要用的各种东西，具体如何构造如何编排顺序直接交给 wire 去管就行
	// 这个方法里面传入各个组件的初始化方法
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	return new(repository.UserRepository)
}

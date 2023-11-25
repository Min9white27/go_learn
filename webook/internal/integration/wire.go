//go:build wireinject

package integration

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基本的第三方依赖
		ioc.InitDB, ioc.InitRedis,

		// 初始化 DAO
		dao.NewUserDAO,
		// 初始化 cache
		cache.NewUserCache, cache.NewCodeCache,

		// 初始化 repository
		repository.NewUserRepository, repository.NewCodeRepository,

		// 初始化 service
		service.NewUserService, service.NewCodeService, ioc.InitSMSService,

		// 初始化 UserHandler
		web.NewUserHandler,

		// 中间件呢？
		// 注册路由呢？
		// unused provider 没有用到前面的任何东西
		//gin.Default,

		// 初始化 gin.Engine
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}

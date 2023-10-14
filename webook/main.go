package main

import (
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/repository"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/repository/dao"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/service"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/web"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {

	db := initDB()
	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)
	//分散式注册路由写法，优点是比较有条理，缺点是找路由时不太好找
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	//server.Use(func(ctx *gin.Context) {
	//	println("这是第一个 middleware")
	//})
	//
	//server.Use(func(ctx *gin.Context) {
	//	println("这是第二个 middleware")
	//})

	server.Use(cors.New(cors.Config{
		// AllowOrigins 用这种的话就需要加入所有需要的地址，也可以用 AllowOriginFunc ，这样写可以不用刻意写全部地址
		//AllowOrigins: []string{"http://localhost:3000"},
		// AllowMethods 不写表示全都包括在内
		//AllowMethods:  []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//ExposeHeaders: []string{"Origin"},
		// AllowCredentials 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("^dg#Wx%vaw8D$OBIRT4AXolVhiM103fN"),
	//	[]byte("!YwTy&F^mBG@4gZ1WDEna9je#chOslQz"))

	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("^dg#Wx%vaw8D$OBIRT4AXolVhiM103fN"), []byte("!YwTy&F^mBG@4gZ1WDEna9je#chOslQz"))
	if err != nil {
		panic(err)
	}

	server.Use(sessions.Sessions("mysession", store))
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("users/signup").IgnorePaths("users/login").Build())
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不用启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

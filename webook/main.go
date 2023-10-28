package main

import (
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/repository"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/repository/dao"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/service"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/web"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/web/middleware"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {

	db := initDB()
	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)
	//分散式注册路由写法，优点是比较有条理，缺点是找路由时不太好找
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，阿橙")
	})
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

	redisClient := redis.NewClient(&redis.Options{
		Addr: "webook-live-redis:11479",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		// AllowOrigins 用这种的话就需要加入所有需要的地址，也可以用 AllowOriginFunc ，这样写可以不用刻意写全部地址
		//AllowOrigins: []string{"http://localhost:3000"},
		// AllowMethods 不写表示全都包括在内
		//AllowMethods:  []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 不加 ExposeHeaders ，前端是拿不到你的token
		ExposeHeaders: []string{"x-jwt-token"},
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
	store := memstore.NewStore([]byte("^dg#Wx%vaw8D$OBIRT4AXolVhiM103fN"),
		[]byte("!YwTy&F^mBG@4gZ1WDEna9je#chOslQz"))

	//第一个参数是最大空闲连接数量
	//第二个就是 tcp ，不太可能用udp
	//第三、四个就是连接信息和密码
	//第五、六个就是两个key:Authentication 是指身份认证，Encryption 是指数据加密
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("^dg#Wx%vaw8D$OBIRT4AXolVhiM103fN"), []byte("!YwTy&F^mBG@4gZ1WDEna9je#chOslQz"))
	//if err != nil {
	//	panic(err)
	//}

	//myStore := &sqlx_store.Store{}

	server.Use(sessions.Sessions("mysession", store))
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
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
	db, err := gorm.Open(mysql.Open("root:root@tcp(lwebook-live-mysql:11309)/webook"))
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

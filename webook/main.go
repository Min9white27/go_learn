package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/integration"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//db := initDB()
	//rdb := initRedis()
	//server := initWebServer()
	//u := initUser(db, rdb)
	//u.RegisterRoutes(server)
	server := integration.InitWebServer()
	//分散式注册路由写法，优点是比较有条理，缺点是找路由时不太好找
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，阿橙")
	})
	//server.Run(":8081")
	server.Run(":8080")
}

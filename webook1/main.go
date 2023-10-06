package main

import (
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook1/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	server.Run()

	//server.Use(cors.New(cors.Config{
	//	AllowOrigins: []string{"http://localhost:3000"},
	//	//AllowMethods:     []string{"POST"},
	//	AllowHeaders: []string{"Origin"},
	//	//ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	AllowOriginFunc: func(origin string) bool {
	//		return origin == "https://github.com"
	//	},
	//	MaxAge: 12 * time.Hour,
	//}))
	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.Run("8080")
}

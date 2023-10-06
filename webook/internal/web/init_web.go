package web

import "github.com/gin-gonic/gin"

//不同注册路由的方法

func RegisterRoutes() *gin.Engine {
	server := gin.Default()

	registerUsersRoutes(server)

	return server
}

func registerUsersRoutes(server *gin.Engine) {
	u := &UserHandler{}

	server.POST("/users/signup", u.SignUp)
	// REST 风格
	//server.PUT("/users", func(context *gin.Context) {
	//
	//})

	server.POST("/users/login", u.Login)

	server.POST("/users/edit", u.Edit)
	// REST 风格
	//server.POST("/users/:id", func(context *gin.Context) {
	//
	//})

	server.GET("/users/profile", u.Profile)
	// REST 风格
	//server.GET("/users/profile", func(context *gin.Context) {
	//
	//})
}

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// get /users/DaMing 查询
// delete /users/DaMing 删除
// put /users/DaMing 注册
// post /users/DaMing 修改

func main() {
	server := gin.Default()
	//当一个HTTP请求，用GET方法访问的时候，如果访问路径是/hello,就执行下面这段代码
	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello,go!")
	}) //这里就是路由注册

	server.POST("/POST", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,POST 方法")
	})
	//上面的就是静态路由

	//参数路由,不能不带冒号 : ，否则会被理解为静态路由
	server.GET("/users/:name", func(ctx *gin.Context) {
		//通过Gin里面的Param方法来获得参数
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "这是你传过来的名字 %s", name)
	})

	//通配符路由，其里面的*不能单独出现，比如/views/*,或者/views/*/a
	server.GET("/views/*.html", func(ctx *gin.Context) {
		path := ctx.Param(".html")
		ctx.String(http.StatusOK, "匹配上的值是 %s", path)
	})

	//查询参数
	server.GET("order", func(ctx *gin.Context) {
		oid := ctx.Query("id")
		ctx.String(http.StatusOK, "查询参数的值为"+oid)
	})

	server.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}

//在Gin里面,一个Web服务器被抽象为Engine,Engine承担了路由注册、接入middleware的核心职责
//gin.Context是Gin里面的核心类型，它的字面意思就是“上下文”，它的职责是处理请求（Request）和返回响应（ResponseWriter）
//Gin为每一个HTTP方法都提供了一个注册路由的方法，其基本上都是两个参数：1.路由规则：比如说/hello这种静态路由。2.处理函数：也就是注册的，返回hello,go的方法
//Gin支持很多类型的路由（常用的有）：
//1.静态路由：完全匹配的路由，前面注册的/hello和/POST都属于。2.参数路由：在路径中带上了参数的路由。3.通配符路由：任意匹配的路由
//适用初学者的两条原则：1.用户是查询数据的，用GET，参数放到查询参数里面。即?a=123这样。2.用户是提交数据的，用POST，参数全部放到Body里面。

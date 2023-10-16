package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		//不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//不需要登录校验的
		//if ctx.Request.URL.Path == "/users/signup" || ctx.Request.URL.Path == "/users/login" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")

		if id == nil {
			//	没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now()
		if updateTime == nil {
			//说明还没有刷新过，刚登录，还没有刷新过
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}
		// updateTime 已经有了
		updateTimeVal, _ := updateTime.(time.Time)
		//if !ok {
		//	ctx.AbortWithStatus(http.StatusInternalServerError)
		//	return
		//}
		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}
	}

}

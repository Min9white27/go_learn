package middleware

import (
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

// LoginJWTMiddlewareBuilder jWT 登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//	用 JWT 来检验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			//	没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			//	没有登录,有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			//	没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			//ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		//	err 为 nil ，token 不为 nil
		if token == nil || !token.Valid {
			//	没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}

}

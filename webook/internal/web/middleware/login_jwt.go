package middleware

import (
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

// LoginJWTMiddlewareBuilder jWT 登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}
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
		tokenStr := l.ExtractToken(ctx)
		claims := &ijwt.UserClaims{}
		// ParseWithClaims 会修改 claims 的数据，所以上面需要传指针
		// ParseWithClaims 里面，一定要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			//	没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			//ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		//正常不需要验证 token 的过期时间，因为 Valid 会验证
		//claims.ExpiresAt.Time.Before(time.Now()){
		////	过期了
		//}
		//	err 为 nil ，token 不为 nil

		if token == nil || !token.Valid || claims.Uid == 0 {
			//	没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 理论上是要加监控的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			// 要么 redis 有问题，要么已经退出登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 这样做退出登录是无效的，因为这样做，tokenStr 并没有改变
		//token.Valid = false

		// 短的 token 过期了，但长的 token 还在，搞个新的短 token
		// 自动刷新机制
		//now := time.Now()
		//// 每十秒刷新一次
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	tokenStr, err := token.SignedString(web.JWTKey)
		//	if err != nil {
		//		//记录日志
		//		log.Println("jwt 续约失败", err)
		//	}
		//	ctx.Header("x-jwt-token", tokenStr)
		//}

		ctx.Set("claims", claims)
		//ctx.Set("userId", claims)
	}

}

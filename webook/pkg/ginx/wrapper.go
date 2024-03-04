package ginx

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var L logger.LoggerV1

// WrapLog 弊端：无法根据这个err，精准找到是哪段业务代码出错
// 可以在主业务代码返回 err的时候，打印出具体错误
// fmt.Errorf("系统错误")
func WrapLog[T any](fn func(ctx context.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			L.Error("业务逻辑出错",
				// 通过打印路径去判断错误在哪
				logger.String("path", ctx.Request.URL.Path),
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

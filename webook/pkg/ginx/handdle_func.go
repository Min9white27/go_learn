package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WrapReq[T any](fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//, uc ijwt.UserClaims
		// 顺便把 userClaims 也取出来
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		if err != nil {

		}
		ctx.JSON(http.StatusOK, res)
	}
}

type Result struct {
	// 这个叫做业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

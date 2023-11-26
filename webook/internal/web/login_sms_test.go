package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_LoginSMS(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		reqBody string

		wantCode int
		wantBody Result
	}{
		{
			name: "验证码校验通过",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "15812345678", "123456").
					Return(true, nil)

				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "15812345678").Return(domain.User{
					Phone: "15812345678",
				}, nil)

				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "15812345678",
	"code": "123456"
}
`,
			wantCode: 200,
			wantBody: Result{
				Msg: "验证码校验通过",
			},
		},
		{
			name: "参数不对，bind 失败",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "15812345678",
	"code": "123456"
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "redis 错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "15812345678", "123456").
					Return(false, errors.New("mock redis 错误"))

				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "15812345678",
	"code": "123456"
}
`,
			wantCode: 200,
			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "验证码错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "15812345678", "123456").
					Return(false, nil)

				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "15812345678",
	"code": "123456"
}
`,
			wantCode: 200,
			wantBody: Result{
				Code: 4,
				Msg:  "验证码错误",
			},
		},
		{
			name: "数据库错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "15812345678", "123456").
					Return(true, nil)

				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "15812345678").Return(domain.User{},
					errors.New("mock db 错误"))

				return userSvc, codeSvc
			},
			reqBody: `
{
	"phone": "15812345678",
	"code": "123456"
}
`,
			wantCode: 200,
			wantBody: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 注册路由
			server := gin.Default()
			us, cs := tc.mock(ctrl)
			h := NewUserHandler(us, cs)
			h.RegisterRoutes(server)

			//	构造 http 请求，因为是单元测试，没有必要真的经过网络发起 HTTP 请求
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewBuffer([]byte(tc.reqBody)))
			//	自己写的代码，应该可以保证这里不会 err，所以这里可以断言一下
			require.NoError(t, err)
			//	将请求塞到 header 里面
			req.Header.Set("Content-Type", "application/json")

			//	构造 HTTP 响应
			resp := httptest.NewRecorder()
			//server.Use(func(ctx *gin.Context) {
			//	ctx.Set("user", UserClaims{
			//		Uid: 123,
			//	})
			//})
			//	HTTP 请求 进入到 gin 框架，这样调用，gin 就会处理这个请求
			//	请求后的响应写会到 resp
			server.ServeHTTP(resp, req)

			//server.GET("user", func(ctx *gin.Context) {
			//
			//})

			//err := h.setJWTToken(server)

			//	比较 http 响应码
			assert.Equal(t, tc.wantCode, resp.Code)
			//	如果响应码不等于 200，req.Body 没有数据，所以就没有必要进行反序列化，直接返回
			if resp.Code != 200 {
				return
			}
			var webRes Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webRes)
		})
	}
}

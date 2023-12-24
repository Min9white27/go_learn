package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/internal/integration/startup"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := startup.InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string

		// 准备数据
		before func(t *testing.T)
		// 验证数据
		after   func(t *testing.T)
		reqBody string

		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				//	redis 什么数据都没有
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//	要清理数据
				val, err := rdb.GetDel(ctx, "phone_code:login:15812345678").Result()
				cancel()
				assert.NoError(t, err)
				// 验证码是 6 位
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phone": "15812345678"
}
`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 这个手机号码，已经又一个验证码了（先给它塞一个）
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15812345678", "123456",
					time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//	要清理数据
				val, err := rdb.GetDel(ctx, "phone_code:login:15812345678").Result()
				cancel()
				assert.NoError(t, err)
				// 验证码是 6 位，清理并没有被覆盖，还是123456
				assert.Equal(t, "123456", val)
			},
			reqBody: `
{
	"phone": "15812345678"
}
`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送太频繁，请稍后再试",
			},
		},
		{
			name: "系统错误",
			// 没有过期时间
			before: func(t *testing.T) {
				// 这个手机号码，已经又一个验证码了，没有过期时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15812345678", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				//	要清理数据
				val, err := rdb.GetDel(ctx, "phone_code:login:15812345678").Result()
				cancel()
				assert.NoError(t, err)
				// 验证码是 6 位
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phone": "15812345678"
}
`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "手机号码输入错误",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
{
	"phone": "14212345678"
}
`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 4,
				Msg:  "手机号码输入错误",
			},
		},
		{
			name: "数据格式错误",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
{
	"phone": ，
}
`,
			wantCode: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send",
				bytes.NewBuffer([]byte(tc.reqBody)))
			// 自己写的代码，可以保证一定不会返回 error, 所以不用判断，要求一定不能 NoError
			// 如果 require 条件不成立，就不会往下执行
			require.NoError(t, err)
			// 数据是 json 格式
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口
			// 当你这样调用的时候， GIN 就会处理这个请求
			// 响应写回到 resp 里
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes web.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webRes)
			tc.after(t)
		})
	}
}

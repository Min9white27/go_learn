package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	repomocks "gitee.com/geekbang/basic-go/webook/internal/repository/mocks"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func Test_userService_Login(t *testing.T) {
	// 做个测试用例都用到的时间
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		//	输入
		email    string
		password string

		//	输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登陆成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$L4q056N2sqw2JX3xGsL40eimrgxvfa3quV./GAUIUrai4MkZVl4t2",
						Phone:    "15812345678",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$L4q056N2sqw2JX3xGsL40eimrgxvfa3quV./GAUIUrai4MkZVl4t2",
				Phone:    "15812345678",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB 错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
		{
			name: "密码不正确",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$L4q056N2sqw2JX3xGsL40eimrgxvfa3quV./GAUIUrai4MkZVl4t2",
						Phone:    "15812345678",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world124",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//	具体的测试代码
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl), logger.L())
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}

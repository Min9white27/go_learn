package repository

import (
	"context"
	"database/sql"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	cachemocks "gitee.com/geekbang/basic-go/webook/internal/repository/cache/mocks"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	daomocks "gitee.com/geekbang/basic-go/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	// 111ms.1111ns
	now := time.Now()
	// 要去掉毫秒以外的部分
	// 111ms
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx context.Context
		id  int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存未命中，查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)

				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
					Id: 123,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "this is password",
					Phone: sql.NullString{
						String: "15812345678",
						Valid:  true,
					},
					Nickname:        "abc",
					Birthday:        now.UnixMilli(),
					PersonalProfile: "hello",
					Ctime:           now.UnixMilli(),
					Utime:           now.UnixMilli(),
				}, nil)

				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:              123,
					Email:           "123@qq.com",
					Password:        "this is password",
					Phone:           "15812345678",
					Nickname:        "abc",
					PersonalProfile: "hello",
					Ctime:           now,
					Birthday:        now,
				}).Return(nil)

				return d, c
			},

			ctx: context.Background(),
			id:  123,

			wantUser: domain.User{
				Id:              123,
				Email:           "123@qq.com",
				Password:        "this is password",
				Phone:           "15812345678",
				Nickname:        "abc",
				PersonalProfile: "hello",
				Ctime:           now,
				Birthday:        now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					Id:              123,
					Email:           "123@qq.com",
					Password:        "this is password",
					Phone:           "15812345678",
					Nickname:        "abc",
					PersonalProfile: "hello",
					Ctime:           now,
					Birthday:        now,
				}, nil)

				d := daomocks.NewMockUserDAO(ctrl)

				return d, c
			},

			ctx: context.Background(),
			id:  123,

			wantUser: domain.User{
				Id:              123,
				Email:           "123@qq.com",
				Password:        "this is password",
				Phone:           "15812345678",
				Nickname:        "abc",
				PersonalProfile: "hello",
				Ctime:           now,
				Birthday:        now,
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中，查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存未命中，查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)

				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{}, errors.New("mock db 错误"))

				return d, c
			},

			ctx: context.Background(),
			id:  123,

			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
			time.Sleep(time.Second)
			//	检测 testSignal
		})
	}
}

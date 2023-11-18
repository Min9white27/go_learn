package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeRedisCache struct {
	client redis.Cmdable
}

func NewCodeRedisCache(client redis.Cmdable) CodeCache {
	return &CodeRedisCache{
		client: client,
	}
}

func (c *CodeRedisCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		//	毫无问题
		return nil
	case -1:
		//	发送太频繁
		return ErrCodeSendTooMany
	//case -2:
	default:
		//	系统错误
		return errors.New("系统错误")
	}
}

func (c *CodeRedisCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		//	毫无问题
		return true, nil
	case -1:
		//	如果频繁出现这个错误，就需要警惕，系统可能在被入侵
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		//	系统错误
		return false, nil
		//default:
		//return false, ErrUnknownForCode
	}
	return false, ErrUnknownForCode
}

//func (c *CodeRedisCache) Verify (ctx context.Context, biz, phone, code string) error {
//
//}

func (c *CodeRedisCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
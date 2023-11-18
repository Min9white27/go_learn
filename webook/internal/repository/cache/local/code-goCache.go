package local

import (
	"context"
	"errors"
	"fmt"
	cache2 "gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type CodeGoCache struct {
	GoCache *cache.Cache
	lock    sync.Mutex
	//Item    *cache.Item
}

func NewCodeFreeCache(GoCache *cache.Cache) *CodeGoCache {
	return &CodeGoCache{
		GoCache: GoCache,
		//Item:    Item,
	}
}

func (c *CodeGoCache) Set(ctx context.Context, biz, phone, code string) error {
	// 设置普通锁
	c.lock.Lock()
	// 解锁普通锁
	defer c.lock.Unlock()

	key := c.key(biz, phone)

	val, ok := c.GoCache.Get(key)
	if !ok {
		//	没有验证码
		return c.GoCache.Add(key, codeItems{
			code: code,
			cnt:  3,
		}, time.Minute*10)
	}
	// 有验证码，检查过期时间
	val, ok = c.GoCache.Get(key)
	itm, ok := val.(codeItems)
	if !ok {
		// 基本上不可能进到这里
		return errors.New("系统错误")
	}
	now := time.Now()
	if itm.expire.Sub(now) > time.Minute*9 {
		// 系统发送验证码太频繁
		return cache2.ErrCodeSendTooMany
	}
	// 没有必要，用户不需要知道
	//if itm.Expired() {
	//	return errors.New("验证码已过期")
	//}
	return c.GoCache.Add(key, codeItems{
		code: code,
		cnt:  3,
	}, time.Minute*10)
}

func (c *CodeGoCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	// 设置普通锁
	c.lock.Lock()
	// 释放普通锁
	defer c.lock.Unlock()

	key := c.key(biz, phone)
	val, ok := c.GoCache.Get(key)
	if !ok {
		//	没有发送验证码
		return false, cache2.ErrKeyNotExist
	}
	itm, ok := val.(codeItems)
	if !ok {
		// 基本上不可能进入到这里
		return false, errors.New("系统错误")
	}
	if itm.cnt <= 0 {
		// 用户验证码验证太多次，可能有人在搞你
		return false, cache2.ErrCodeVerifyTooManyTimes
	}
	itm.cnt--
	return itm.code == inputCode, nil
}

func (c *CodeGoCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type codeItems struct {
	code   string
	cnt    int
	expire time.Time
}

func (ci *codeItems) Expired() bool {
	return ci.expire.Before(time.Now())
}

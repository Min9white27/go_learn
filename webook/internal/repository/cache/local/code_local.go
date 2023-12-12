package local

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"github.com/hashicorp/golang-lru/v2"
	"sync"
	"time"
)

// 技术选型考虑的点
//  1. 功能性：功能是否能够完全覆盖你的需求。
//  2. 社区和支持度：社区是否活跃，文档是否齐全，
//     以及百度（搜索引擎）能不能搜索到你需要的各种信息，有没有帮你踩过坑
//  3. 非功能性：易用性（用户友好度，学习曲线要平滑），
//     扩展性（如果开源软件的某些功能需要定制，框架是否支持定制，以及定制的难度高不高）
//     性能（追求性能的公司，往往有能力自研）

type CodeLocalCache struct {
	lruCache *lru.Cache[string, any]
	lock     sync.Mutex
	//rwLock     sync.RWMutex
	//maps       sync.Map
	// 过期时间
	expiration time.Duration
}

func NewLocalCodeCache(LruCache *lru.Cache[string, any], expiration time.Duration) *CodeLocalCache {
	return &CodeLocalCache{
		lruCache:   LruCache,
		expiration: expiration,
	}
}

func (l *CodeLocalCache) Set(ctx context.Context, biz, phone, code string) error {
	// 设置普通锁
	l.lock.Lock()
	// 释放普通锁
	defer l.lock.Unlock()

	key := l.key(biz, phone)

	// 如果你的 key 非常多，这个 maps 本身就占据了很多内存
	//lock, _ := l.maps.LoadOrStore(key, &sync.Mutex{})
	//lock.(*sync.Mutex).Lock()
	//defer lock.(*sync.Mutex).Unlock()

	//lock, _ := l.maps.LoadOrStore(key, &sync.Mutex{})
	//lock.(*sync.Mutex).Lock()
	//defer func() {
	//	l.maps.Delete(key)
	//	lock.(*sync.Mutex).Unlock()
	//}()

	now := time.Now()

	val, ok := l.lruCache.Get(key)
	if !ok {
		// 没有验证码
		l.lruCache.Add(key, codeItem{
			code:   code,
			cnt:    3,
			expire: now.Add(l.expiration),
		})
		return nil
	}
	// 有验证码，检查过期时间
	itm, ok := val.(codeItem)
	if !ok {
		// 基本上不可能进入这里
		return errors.New("系统错误")
	}
	if itm.expire.Sub(now) > time.Minute*9 {
		// 发送验证码太频繁
		return cache.ErrCodeSendTooMany
	}
	// 验证码过期，重新发送
	l.lruCache.Add(key, codeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expiration),
	})
	return nil
}

func (l *CodeLocalCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	// 设置普通锁
	l.lock.Lock()
	// 释放普通锁
	defer l.lock.Unlock()

	key := l.key(biz, phone)
	val, ok := l.lruCache.Get(key)
	if !ok {
		// 没有发送验证码
		return false, cache.ErrKeyNotExist
	}
	// 有验证码，减少可验证次数并比较用户输入的验证码和系统发送的验证码是否相同
	itm, ok := val.(codeItem)
	if !ok {
		// 基本上不可能进入这里
		return false, errors.New("系统错误")
	}
	if itm.cnt <= 0 {
		// 验证码验证太多次，可能有人在搞你
		return false, cache.ErrCodeVerifyTooManyTimes
	}
	itm.cnt--
	return itm.code == inputCode, nil
}

func (l *CodeLocalCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}

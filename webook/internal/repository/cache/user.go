package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	// 可以传单机 Redis
	// 也可以传 cluster 的 Redis
	client redis.Cmdable
	// 过期时间
	expiration time.Duration
}

// A 用到了 B，B 一定是接口 => 保证面向接口
// A 用到了 B，B 一定是 A 的字段 => 规避包变量、包方法，这两种都缺乏扩展性
// A 用到了 B，A 绝对不初始化 B，而是外面注入 => 保持依赖注入（DI, Dependency Injection）和依赖反转（IOC）
// expiration 1s, 1m

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// Get 只要 error 为 nil, 就认为缓存里面有数据
// 如果没有数据，返回一个特定的 error
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	// 数据不存在，err = redis.Nil
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	//if err!=nil{
	//	return domain.User{},err
	//}
	//return u,nil
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

// 不好的写法
// main 函数里面初始化好
//var RedisClient *redis.Client
//
//func GetUser(ctx context.Context,id int64){
//	RedisClient.Get()
//}

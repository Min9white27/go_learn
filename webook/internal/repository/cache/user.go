package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	// 可以传单机 Redis
	// 也可以传 cluster 的 Redis
	client redis.Cmdable
	// 过期时间
	expiration time.Duration
}

// A 用到了 B，B 一定是借口
// A 用到了 B，B 一定是 A 的字段
// A 用到了 B，A 绝对不初始化 B，而是外面注入

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// Get 只要 error 为 nil, 就认为缓存里面有数据
// 如果没有数据，返回一个特定的 error
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
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

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

// 不好的写法
// main 函数里面初始化好
//var RedisClient *redis.Client
//
//func GetUser(ctx context.Context,id int64){
//	RedisClient.Get()
//}
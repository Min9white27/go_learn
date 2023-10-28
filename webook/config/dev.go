//go:build !k8s

// dd go:build dev
// dd go:build test
// dd go:build e2e

// 没有 k8s 这个编译标签
package config

var Config = config{
	DB: DBConfig{
		//本地连接
		DSN: "localhost:13316",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}

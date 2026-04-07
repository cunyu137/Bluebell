package redis

import (
	"Bluebell/settings"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var (
	client *redis.Client
	Nil    = redis.Nil
)

// Init 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 20,               // ✅ 保持20个空闲连接
		IdleTimeout:  5 * time.Minute,  // ✅ 空闲连接5分钟后关闭
		MaxConnAge:   10 * time.Minute, // ✅ 连接最大生命周期10分钟
		DialTimeout:  5 * time.Second,  // ✅ 连接超时5秒
		ReadTimeout:  3 * time.Second,  // ✅ 读超时3秒
		WriteTimeout: 3 * time.Second,  // ✅ 写超时3秒
		PoolTimeout:  4 * time.Second,  // ✅ 获取连接超时4秒
	})

	_, err = client.Ping().Result()
	return
}

func Close() {
	_ = client.Close()
}

func GetPoolStats() redis.PoolStats {
	return *client.PoolStats()
}

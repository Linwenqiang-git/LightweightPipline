package redis_proxy

import (
	. "lightweightpipline/configs"

	"github.com/go-redis/redis/v8"
)

// redis客户端（根据情况是否对go-redis包一层）
type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(configure Configure) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     configure.Redis.Addr,     // Redis 地址
		Password: configure.Redis.Password, // Redis 密码
		DB:       configure.Redis.DB,       // 默认数据库
	})
	redisClient := &RedisClient{
		Client: client,
	}
	return redisClient, nil
}

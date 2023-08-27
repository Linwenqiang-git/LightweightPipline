package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type RedisOption struct {
	addr            string
	password        string
	DB              int
	IsOpenSentinel  bool
	RedisSentinelIp []string
	ConnectTimeout  int
	RedisPrefix     string
}

func (r *RedisOption) CreateConnect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     r.addr,     // Redis 地址
		Password: r.password, // Redis 密码
		DB:       r.DB,       // 默认数据库
	})
	return client
}

func NewRedisOption() *RedisOption {
	return &RedisOption{
		addr:        viper.GetString("redis.Addr"),
		DB:          viper.GetInt("redis.DB"),
		password:    viper.GetString("redis.Password"),
		RedisPrefix: viper.GetString("redis.RedisPrefix"),
	}
}

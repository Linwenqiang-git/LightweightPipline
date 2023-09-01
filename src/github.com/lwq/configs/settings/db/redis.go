package db

import (
	"github.com/spf13/viper"
)

type RedisOption struct {
	Addr            string
	Password        string
	DB              int
	IsOpenSentinel  bool
	RedisSentinelIp []string
	ConnectTimeout  int
	RedisPrefix     string
}

func NewRedisOption() *RedisOption {
	return &RedisOption{
		Addr:        viper.GetString("redis.Addr"),
		DB:          viper.GetInt("redis.DB"),
		Password:    viper.GetString("redis.Password"),
		RedisPrefix: viper.GetString("redis.RedisPrefix"),
	}
}

//go:build wireinject
// +build wireinject

package wire

import (
	conf "lightweightpipline/configs"
	irepo "lightweightpipline/internal/biz/repo"
	repo "lightweightpipline/internal/data/repo"
	redis "lightweightpipline/third_party/redis_proxy"

	"github.com/google/wire"
)

var configureSet = wire.NewSet(conf.ProvideConfigure)

var dbContextSet = wire.NewSet(configureSet, irepo.ProvideDbContext)

var redisSet = wire.NewSet(configureSet, redis.NewRedisClient)

var commandRepoSet = wire.NewSet(dbContextSet, repo.NewCommandRepo)

// =========配置信息=========
func GetConfigure() (conf.Configure, error) {
	panic(wire.Build(configureSet))
}

// =========数据库上下文=========
//database
func GetDbContext() (*irepo.DbContext, error) {
	panic(wire.Build(dbContextSet))
}
//redis
func NewRedis() (*redis.RedisClient, error) {
	panic(wire.Build(redisSet))
}

// =========仓储层=========
func NewCommandRepo() (irepo.ICommandRepo, error) {
	panic(wire.Build(commandRepoSet))
}

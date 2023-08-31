//go:build wireinject
// +build wireinject

package wire

import (
	conf "lightweightpipline/configs"
	irepo "lightweightpipline/internal/biz/repo"
	data "lightweightpipline/internal/data/repo"

	"github.com/google/wire"
)

var configureSet = wire.NewSet(conf.ProvideConfigure)

var dbContextSet = wire.NewSet(configureSet, irepo.ProvideDbContext)

var commandRepoSet = wire.NewSet(dbContextSet, data.NewCommandRepo)

// 配置信息
func GetConfigure() (conf.Configure, error) {
	panic(wire.Build(configureSet))
}

// 数据库上下文
func GetDbContext() (*irepo.DbContext, error) {
	panic(wire.Build(dbContextSet))
}

// 仓储层
func NewCommandRepo() (irepo.ICommandRepo, error) {
	panic(wire.Build(commandRepoSet))
}

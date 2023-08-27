//go:build wireinject
// +build wireinject

package wire

import (
	conf "lightweightpipline/configs"
	dbContext "lightweightpipline/internal/biz/repo"

	"github.com/google/wire"
)

var configureSet = wire.NewSet(conf.ProvideConfigure)

var dbContextSet = wire.NewSet(configureSet, dbContext.ProvideDbContext)

func GetConfigure() (conf.Configure, error) {
	wire.Build(configureSet)
	return conf.Configure{}, nil
}

func GetDbContext() (dbContext.DbContext, error) {
	wire.Build(dbContextSet)
	return dbContext.DbContext{}, nil
}

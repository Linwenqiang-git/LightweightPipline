package settings

import (
	"github.com/spf13/viper"
)

type SystemOption struct {
	Env          string
	Addr         string
	DbType       string
	BuildRootDir string
}

func NewSystemConfig() *SystemOption {
	systemConfig := &SystemOption{
		Env:          viper.GetString("system.env"),
		Addr:         viper.GetString("system.addr"),
		DbType:       viper.GetString("system.db-type"),
		BuildRootDir: viper.GetString("system.build-root-dir"),
	}
	return systemConfig
}

package settings

import "github.com/spf13/viper"

type SystemOption struct {
	Env    string
	Addr   string
	DbType string
}

func NewSystemConfig() *SystemOption {
	return &SystemOption{
		Env:    viper.GetString("system.env"),
		Addr:   viper.GetString("system.addr"),
		DbType: viper.GetString("system.db-type"),
	}
}

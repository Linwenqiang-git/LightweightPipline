package settings

import "github.com/spf13/viper"

type System struct {
	Env    string
	Addr   int
	DbType string
}

func NewSystemConfig() *System {
	return &System{
		Env:    viper.GetString("system.env"),
		Addr:   viper.GetInt("system.addr"),
		DbType: viper.GetString("system.db-type"),
	}
}

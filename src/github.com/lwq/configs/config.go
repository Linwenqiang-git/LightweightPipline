package config

import (
	"fmt"
	"sync"

	. "lightweightpipline/configs/settings"
	. "lightweightpipline/configs/settings/db"

	"github.com/spf13/viper"
)

var configure Configure
var once sync.Once

type Configure struct {
	System *SystemOption
	Mysql  *Mysql
	Pgsql  *Pgsql
	Redis  *RedisOption
}

func ProvideConfigure() (Configure, error) {
	once.Do(func() {
		viper.SetConfigName("appsettings") // 配置文件名
		viper.SetConfigType("yaml")        // 配置文件类型
		viper.AddConfigPath("../configs")  // 配置文件路径(需配置本地机密文件地址)
		// 加载配置文件
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
		//init setting
		configure = Configure{
			System: NewSystemConfig(),
			Mysql:  NewMysql(),
			Pgsql:  NewPgSql(),
			Redis:  NewRedisOption(),
		}
	})
	return configure, nil
}

package db

import (
	"fmt"

	"github.com/spf13/viper"
)

type Pgsql struct {
	Addr         string `mapstructure:"addr" json:"addr" yaml:"addr"`                             // 服务器地址:端口
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`                             //:端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                       // 高级配置
	Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`                     // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                 // 数据库用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                 // 数据库密码
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`                  // 是否开启Gorm全局日志
	LogZap       bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"`                     // 是否通过zap写入日志文件
}

func (s *Pgsql) GetConnectStr() string {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		s.Username,
		s.Password,
		s.Addr,
		s.Port,
		s.Dbname)

	return dsn
}

func NewPgSql() *Pgsql {
	return &Pgsql{
		Addr:     viper.GetString("pgsql.addr"),
		Port:     viper.GetInt("pgsql.port"),
		Dbname:   viper.GetString("pgsql.db-name"),
		Username: viper.GetString("pgsql.username"),
		Password: viper.GetString("pgsql.password"),
	}
}

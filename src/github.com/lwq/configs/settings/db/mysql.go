package db

import (
	"fmt"

	"github.com/spf13/viper"
)

type Mysql struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`             // 服务器地址
	Port         string `mapstructure:"port" json:"port" yaml:"port"`             // 端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`       // 高级配置
	Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`     // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"` // 数据库用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Charset      string `mapstructure:"charset" json:"charset" yaml:"charset"`                    // 数据库密码
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`                  // 是否开启Gorm全局日志
	LogZap       bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"`                     // 是否通过zap写入日志文件
}

func (s *Mysql) GetConnectStr() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&allowNativePasswords=true",
		s.Username,
		s.Password,
		s.Path,
		s.Port,
		s.Dbname,
		s.Charset)
	return dsn
}

func NewMysql() *Mysql {
	return &Mysql{
		Path:     viper.GetString("mysql.path"),
		Port:     viper.GetString("mysql.port"),
		Dbname:   viper.GetString("mysql.db-name"),
		Username: viper.GetString("mysql.username"),
		Password: viper.GetString("mysql.password"),
		Charset:  viper.GetString("mysql.charset"),
	}
}

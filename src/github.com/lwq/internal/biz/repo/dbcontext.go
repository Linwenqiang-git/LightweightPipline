package repo

import (
	. "lightweightpipline/configs"
	. "lightweightpipline/configs/settings/log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbContext struct {
	db *gorm.DB
}

func (c *DbContext) GetDb() *gorm.DB {
	return c.db
}

func ProvideDbContext(configure Configure) (dbContext *DbContext, err error) {
	configType := configure.System.DbType
	var db *gorm.DB
	switch configType {
	case "mysql":
		connectStr := configure.Mysql.GetConnectStr()
		db, err = gorm.Open(mysql.Open(connectStr), &gorm.Config{})
	case "pgsql":
		connectStr := configure.Pgsql.GetConnectStr()
		db, err = gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	default:
		Logger.Fatal("不支持的数据库类型")
	}

	return &DbContext{db: db}, err
}

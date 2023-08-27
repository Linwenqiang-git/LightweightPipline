package main

import (
	provider "lightweightpipline/cmd/wire"
	entity "lightweightpipline/internal/biz/aggregates"

	"gorm.io/gorm"
)

func MigrateDb() *gorm.DB {

	dbContext, err := provider.GetDbContext()
	if err != nil {
		panic("connect error：" + err.Error())
	}
	return dbContext.GetDb()
}

func main() {
	db := MigrateDb()
	//db.AutoMigrate(&entity.Application{})
	db.AutoMigrate(&entity.AppBuildStep{})
	//db.AutoMigrate(&entity.AppBuildRecord{})
	//db.AutoMigrate(&entity.AppBuildRecordDetail{})
	// db.AutoMigrate(&entity.Template{})
	// db.AutoMigrate(&entity.Stage{})
	// db.AutoMigrate(&entity.StageCommand{})
	// db.AutoMigrate(&entity.Command{})
}

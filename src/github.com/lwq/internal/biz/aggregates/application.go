package aggregates

import "gorm.io/gorm"

//应用信息（界面维护）
type Application struct {
	gorm.Model
	Name string `sql:"COMMENT:'应用名称'"`
	Path string `sql:"COMMENT:'应用路径'"`
}

func (Application) TableName() string {
	return "application"
}

//应用构建阶段信息（界面维护）
type AppBuildStep struct {
	gorm.Model
	AppId     uint
	StageId   uint
	StageName string
	Sort      int  `sql:"COMMENT:'阶段顺序'"`
	IsOpen    bool `sql:"COMMENT:'是否启用'"`
}

func (AppBuildStep) TableName() string {
	return "app_build_step"
}

//应用构建记录信息（系统生成）
type AppBuildRecord struct {
	gorm.Model
	AppId     uint
	ImageName string
	Branch    string
	Progress  uint
	Result    uint
}

func (AppBuildRecord) TableName() string {
	return "appbuildrecord"
}

//构建记录明细（系统生成）
type AppBuildRecordDetail struct {
	gorm.Model
	BuildRecordId uint
	StageId       uint
	StageName     string
	CommandId     uint
	Command       string `sql:"COMMENT:'命令'"`
	Output        string `sql:"COMMENT:'命令行输出'"`
	Result        uint   `sql:"COMMENT:'运行结果'"`
}

func (AppBuildRecordDetail) TableName() string {
	return "appbuildrecorddetail"
}

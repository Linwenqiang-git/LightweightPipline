package aggregates

import "gorm.io/gorm"

//应用信息（界面维护）
type Application struct {
	gorm.Model
	Name             string `sql:"COMMENT:'应用名称'"`
	CodeRepository   string `sql:"COMMENT:'代码仓库地址'"`
	ServiceName      string `sql:"COMMENT:'服务名称'"`
	ClusterNameSpace string `sql:"COMMENT:'集群命名空间'"`
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
	Command       string  `sql:"COMMENT:'命令'"`
	SecondCommand string  `sql:"COMMENT:'第二命令'"`
	Output        string  `sql:"COMMENT:'命令行输出'"`
	RunTime       float64 `sql:"COMMENT:'运行时间（秒）'"`
	Result        uint    `sql:"COMMENT:'运行结果'"`
}

func (AppBuildRecordDetail) TableName() string {
	return "appbuildrecorddetail"
}

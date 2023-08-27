package aggregates

import "gorm.io/gorm"

//构建阶段信息（界面维护）
type Stage struct {
	gorm.Model
	Name      string
	Explain   string `sql:"COMMENT:'阶段说明'"`
	TempletId uint   `sql:"COMMENT:'模板Id'"`
	Sort      int    `sql:"COMMENT:'阶段排序'"`
}

func (Stage) TableName() string {
	return "stage"
}

//构建阶段命令信息（界面维护）
type StageCommand struct {
	gorm.Model
	StageId   uint
	CommandId uint
	Sort      int `sql:"COMMENT:'命令排序'"`
}

func (StageCommand) TableName() string {
	return "stagecommand"
}

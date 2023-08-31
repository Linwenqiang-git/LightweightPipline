package aggregates

import "gorm.io/gorm"

//命令信息（界面维护）
type Command struct {
	gorm.Model
	Detail        string `sql:"COMMENT:'命令细节'"`
	ErrorKeyword  string `sql:"COMMENT:'命令执行出错关键字'"`
	SecondCommand string `sql:"COMMENT:'替代命令'"`
	Explain       string `sql:"COMMENT:'命令说明'"`
}

func (Command) TableName() string {
	return "command"
}

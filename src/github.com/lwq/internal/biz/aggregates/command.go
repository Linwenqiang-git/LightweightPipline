package aggregates

import "gorm.io/gorm"

type Command struct {
	gorm.Model
	Detail    string `sql:"COMMENT:'命令细节'"`
	Parameter string `sql:"COMMENT:'命令参数'"`
	Explain   string `sql:"COMMENT:'命令说明'"`
}

func (Command) TableName() string {
	return "command"
}

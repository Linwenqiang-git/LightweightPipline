package aggregates

import "gorm.io/gorm"

//构建模板信息（界面维护）
type Template struct {
	gorm.Model
	Name    string
	Explain string `sql:"COMMENT:'模板说明'"`
}

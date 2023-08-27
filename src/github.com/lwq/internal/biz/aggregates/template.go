package aggregates

import "gorm.io/gorm"

type Template struct {
	gorm.Model
	Name    string
	Explain string `sql:"COMMENT:'模板说明'"`
}

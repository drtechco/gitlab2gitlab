// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameLang = "lang"

// Lang mapped from table <lang>
type Lang struct {
	Key  string `gorm:"column:key;type:text(255)" json:"key"`
	Memo string `gorm:"column:memo;type:text(255)" json:"memo"`
}

// TableName Lang's table name
func (*Lang) TableName() string {
	return TableNameLang
}
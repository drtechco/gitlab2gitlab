// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameI18n = "i18n"

// I18n mapped from table <i18n>
type I18n struct {
	LangKey string `gorm:"column:lang_key;type:varchar(255)" json:"lang_key"`
	Key     string `gorm:"column:key;type:varchar(255)" json:"key"`
	Value   string `gorm:"column:value;type:varchar(3000)" json:"value"`
	Memo    string `gorm:"column:memo;type:varchar(255)" json:"memo"`
}

// TableName I18n's table name
func (*I18n) TableName() string {
	return TableNameI18n
}

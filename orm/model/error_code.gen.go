// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameErrorCode = "error_code"

// ErrorCode mapped from table <error_code>
type ErrorCode struct {
	Code    int32  `gorm:"column:code;type:INTEGER" json:"code"`
	I18nKey string `gorm:"column:i18n_key;type:text(255)" json:"i18n_key"`
	Memo    string `gorm:"column:memo;type:text(255)" json:"memo"`
}

// TableName ErrorCode's table name
func (*ErrorCode) TableName() string {
	return TableNameErrorCode
}

package entity

import "gorm.io/gorm"

type BaseEntity struct {
	*gorm.DB    `gorm:"-" json:"-"`
	Id          int64  `gorm:"primaryKey" json:"id"`
	GmtCreate   string `json:"gmt_create"` //日期时间字段统一设置为字符串即可
	GmtModified string `json:"gmt_modified"`
	//DeletedAt gorm.DeletedAt `json:"deleted_at"`   // 如果开发者需要使用软删除功能，打开本行注释掉的代码即可，同时需要在数据库的所有表增加字段deleted_at 类型为 datetime
}

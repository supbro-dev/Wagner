/*
* @Author: supbro
* @Date:   2025/6/3 11:19
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/3 11:19
 */
package entity

type ScriptEntity struct {
	BaseEntity
	Name    string `gorm:"column:name" json:"name"`
	Type    string `gorm:"column:type" json:"type"`
	Desc    string `gorm:"column:desc" json:"desc"`
	Content string `gorm:"column:content" json:"content"`
	Version int    `gorm:"column:version" json:"version"`
}

func (u *ScriptEntity) TableName() string {
	return "script" // 自定义表名
}

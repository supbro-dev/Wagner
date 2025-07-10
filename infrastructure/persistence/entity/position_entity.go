/*
* @Author: supbro
* @Date:   2025/7/9 11:23
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/9 11:23
 */
package entity

type PositionEntity struct {
	BaseEntity
	Name string `gorm:"column:name" json:"name"`
	Code string `gorm:"column:code" json:"code"`
}

func (u *PositionEntity) TableName() string {
	return "position" // 自定义表名
}

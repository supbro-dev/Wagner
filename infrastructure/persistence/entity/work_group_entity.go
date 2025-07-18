/*
* @Author: supbro
* @Date:   2025/7/17 14:29
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/17 14:29
 */
package entity

// 工作组
type WorkGroupEntity struct {
	BaseEntity
	Name          string `gorm:"column:name" json:"name"`
	Code          string `gorm:"column:code" json:"code"`
	WorkplaceCode string `gorm:"column:workplace_code" json:"workplaceCode"`
	PositionCode  string `gorm:"column:position_code" json:"positionCode"`
	Desc          string `gorm:"column:desc" json:"desc"`
}

func (u *WorkGroupEntity) TableName() string {
	return "work_group" // 自定义表名
}

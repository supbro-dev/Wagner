/*
* @Author: supbro
* @Date:   2025/6/2 11:26
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 11:26
 */
package entity

type WorkplaceEntity struct {
	BaseEntity
	Code            string `gorm:"column:code" json:"code"`
	Name            string `gorm:"column:name" json:"name"`
	RegionCode      string `gorm:"column:region_code" json:"region_code"`
	IndustryCode    string `gorm:"column:industry_code" json:"industry_code"`
	SubIndustryCode string `gorm:"column:sub_industry_code" json:"sub_industry_code"`
}

func (u *WorkplaceEntity) TableName() string {
	return "workplace" // 自定义表名
}

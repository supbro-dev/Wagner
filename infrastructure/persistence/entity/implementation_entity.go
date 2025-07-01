/*
* @Author: supbro
* @Date:   2025/7/1 12:59
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/1 12:59
 */
package entity

type Implementation struct {
	BaseEntity
	Name                   string               `gorm:"column:name" json:"name"`
	TargetCode             string               `gorm:"column:target_code" json:"targetCode"`
	TargetType             TargetType           `gorm:"column:target_type" json:"targetType"`
	StandardPositionRootId int64                `gorm:"column:standard_position_root_id" json:"standardPositionRootId"`
	Status                 ImplementationStatus `gorm:"column:status" json:"status"`
}

type TargetType string

var (
	Workplace   TargetType = "workplace"
	Industry    TargetType = "industry"
	SubIndustry TargetType = "subIndustry"
)

type ImplementationStatus string

var (
	// 准备中
	Preparing ImplementationStatus = "preparing"
	// 上线
	Online ImplementationStatus = "online"
	// 下线
	Offline ImplementationStatus = "offline"
)

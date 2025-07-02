/*
* @Author: supbro
* @Date:   2025/7/1 12:59
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/1 12:59
 */
package entity

type ProcessImplementationEntity struct {
	BaseEntity
	Name       string     `gorm:"column:name" json:"name"`
	TargetCode string     `gorm:"column:target_code" json:"targetCode"`
	TargetType TargetType `gorm:"column:target_type" json:"targetType"`
	// rootId同时也是ProcessPosition的版本号
	ProcessPositionRootId int64                `gorm:"column:process_position_root_id" json:"processPositionRootId"`
	Status                ImplementationStatus `gorm:"column:status" json:"status"`
}

func (u *ProcessImplementationEntity) TableName() string {
	return "process_implementation" // 自定义表名
}

type TargetType string

var (
	Workplace   TargetType = "workplace"
	Industry    TargetType = "industry"
	SubIndustry TargetType = "subIndustry"
)

type ImplementationStatus string

var (
	// 上线
	Online ImplementationStatus = "online"
	// 下线
	Offline ImplementationStatus = "offline"
)

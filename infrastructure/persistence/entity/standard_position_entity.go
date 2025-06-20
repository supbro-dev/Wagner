package entity

// 标准岗位
type StandardPositionEntity struct {
	BaseEntity
	Code            string               `gorm:"column:code" json:"code"`
	Name            string               `gorm:"column:name" json:"name"`
	ParentCode      string               `gorm:"column:parent_code" json:"parent_code"`
	Type            StandardPositionType `gorm:"column:type" json:"type"`
	Level           int                  `gorm:"column:level" json:"level"`
	Version         int                  `gorm:"column:version" json:"version"`
	IndustryCode    string               `gorm:"column:industry_code" json:"industry_code"`
	SubIndustryCode string               `gorm:"column:sub_industry_code" json:"sub_industry_code"`
	Script          string               `gorm:"column:script" json:"script"`
	Properties      string               `gorm:"column:properties" json:"properties"`
	Order           int                  `gorm:"column:order" json:"order"`
}

func (u *StandardPositionEntity) TableName() string {
	return "standard_position" // 自定义表名
}

type StandardPositionType string

var (
	// 部门
	DEPT StandardPositionType = "DEPT"
	// 岗位
	POSITION StandardPositionType = "POSITION"
	// 直接环节
	DIRECT_PROCESS StandardPositionType = "DIRECT_PROCESS"
	// 间接环节
	INDIRECT_PROCESS StandardPositionType = "INDIRECT_PROCESS"
)

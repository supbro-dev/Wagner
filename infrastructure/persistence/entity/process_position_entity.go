package entity

// 标准岗位
type ProcessPositionEntity struct {
	BaseEntity
	Code            string              `gorm:"column:code" json:"code"`
	Name            string              `gorm:"column:name" json:"name"`
	ParentCode      string              `gorm:"column:parent_code" json:"parent_code"`
	Type            ProcessPositionType `gorm:"column:type" json:"type"`
	Level           int                 `gorm:"column:level" json:"level"`
	Version         int                 `gorm:"column:version" json:"version"`
	IndustryCode    string              `gorm:"column:industry_code" json:"industry_code"`
	SubIndustryCode string              `gorm:"column:sub_industry_code" json:"sub_industry_code"`
	Script          string              `gorm:"column:script" json:"script"`
	Properties      string              `gorm:"column:properties" json:"properties"`
	Order           int                 `gorm:"column:order" json:"order"`
}

func (u *ProcessPositionEntity) TableName() string {
	return "process_position" // 自定义表名
}

type ProcessPositionType string

var (
	ROOT ProcessPositionType = "ROOT"
	// 部门
	DEPT ProcessPositionType = "DEPT"
	// 岗位
	POSITION ProcessPositionType = "POSITION"
	// 直接环节
	DIRECT_PROCESS ProcessPositionType = "DIRECT_PROCESS"
	// 间接环节
	INDIRECT_PROCESS ProcessPositionType = "INDIRECT_PROCESS"
)

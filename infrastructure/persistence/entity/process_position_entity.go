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
	Properties      *string             `gorm:"column:properties" json:"properties"`
	SortIndex       int                 `gorm:"column:sort_index" json:"sortIndex"`
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
	// 工作组
	WORK_GROUP ProcessPositionType = "WORK_GROUP"
	// 直接环节
	DIRECT_PROCESS ProcessPositionType = "DIRECT_PROCESS"
	// 间接环节
	INDIRECT_PROCESS ProcessPositionType = "INDIRECT_PROCESS"
)

func ProcessPositionType2Desc(processPositionType ProcessPositionType) string {
	switch processPositionType {
	case ROOT:
		return "根节点"
	case DEPT:
		return "部门"
	case POSITION:
		return "岗位"
	case DIRECT_PROCESS:
		return "直接环节"
	case INDIRECT_PROCESS:
		return "间接环节"
	default:
		return "未知"
	}
}

var (
	MaxTimeInMinuteKey = "maxTimeInMinute"
	MinIdleTimeKey     = "minIdleTimeInMinute"
	WorkLoadRollUpKey  = "workLoadRollUp"
)

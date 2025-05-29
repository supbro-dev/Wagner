package entity

type StandardPositionEntity struct {
	BaseEntity
	Code            string `gorm:"column:code" json:"code"`
	Name            string `gorm:"column:name" json:"name"`
	ParentCode      string `gorm:"column:parent_code" json:"parent_code"`
	Type            string `gorm:"column:type" json:"type"`
	Level           string `gorm:"column:level" json:"level"`
	Version         string `gorm:"column:version" json:"version"`
	IndustryCode    string `gorm:"column:industry_code" json:"industry_code"`
	SubIndustryCode string `gorm:"column:sub_industry_code" json:"sub_industry_code"`
}

func (u *StandardPositionEntity) TableName() string {
	return "standard_position" // 自定义表名
}

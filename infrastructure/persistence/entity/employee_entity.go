package entity

type EmployeeEntity struct {
	BaseEntity
	Name          string `gorm:"column:name" json:"name"`
	Number        string `gorm:"column:number" json:"number"`
	Identity      string `gorm:"column:identity" json:"identity"`
	SensitiveInfo string `gorm:"column:sensitive_info" json:"sensitive_info"`
	WorkplaceCode string `gorm:"column:workplace_code" json:"workplace_code"`
	PositionCode  string `gorm:"column:position_code" json:"position_code"`
	WorkGroupCode string `gorm:"column:work_group_code" json:"work_group_code"`
	Status        string `gorm:"column:status" json:"status"`
}

func (u *EmployeeEntity) TableName() string {
	return "employee" // 自定义表名
}

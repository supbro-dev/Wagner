package entity

type ActionEntity struct {
	BaseEntity
	EmployeeNumber string `gorm:"column:employee_number" json:"employee_number"`
	WorkplaceCode  string `gorm:"column:workplace_code" json:"workplace_code"`
	OperateDay     string `gorm:"column:operate_day;type:date" json:"operate_day"`
	StartTime      string `gorm:"column:start_time" json:"start_time"`
	EndTime        string `gorm:"column:end_time" json:"end_time"`
	ActionCode     string `gorm:"column:action_code" json:"action_code"`
	ActionType     string `gorm:"column:action_type" json:"action_type"`
	WorkLoad       string `gorm:"column:work_load" json:"work_load"`
	Properties     string `gorm:"column:properties" json:"properties"`
}

func (u *ActionEntity) TableName() string {
	return "action" // 自定义表名
}

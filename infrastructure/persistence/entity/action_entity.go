package entity

import "time"

type ActionEntity struct {
	BaseEntity
	EmployeeNumber string                 `gorm:"column:employee_number" json:"employee_number"`
	WorkplaceCode  string                 `gorm:"column:workplace_code" json:"workplace_code"`
	OperateDay     time.Time              `gorm:"column:operate_day;type:date" json:"operate_day"`
	StartTime      time.Time              `gorm:"column:start_time" json:"start_time"`
	EndTime        time.Time              `gorm:"column:end_time" json:"end_time"`
	ActionCode     string                 `gorm:"column:action_code" json:"action_code"`
	ActionType     string                 `gorm:"column:action_type" json:"action_type"`
	WorkLoad       map[string]interface{} `gorm:"column:work_load:type:json" json:"work_load"`
	Properties     map[string]interface{} `gorm:"column:properties;type:json" json:"properties"`
}

func (u *ActionEntity) TableName() string {
	return "action" // 自定义表名
}

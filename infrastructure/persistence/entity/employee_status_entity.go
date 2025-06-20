/*
* @Author: supbro
* @Date:   2025/6/18 14:06
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/18 14:06
 */
package entity

import "time"

type EmployeeStatusEntity struct {
	BaseEntity
	EmployeeNumber string         `gorm:"column:employee_number" json:"employeeNumber"`
	EmployeeName   string         `gorm:"column:employee_name" json:"employeeName"`
	OperateDay     *time.Time     `gorm:"column:operate_day" json:"operateDay"`
	WorkplaceCode  string         `gorm:"column:workplace_code" json:"workplaceCode"`
	Status         EmployeeStatus `gorm:"column:status" json:"status"`
	LastActionTime *time.Time     `gorm:"column:last_action_time" json:"lastActionTime"`
	LastActionCode string         `gorm:"column:last_action_code" json:"lastActionCode"`
	WorkGroupCode  string         `gorm:"column:work_group_code" json:"workGroupCode"`
}

type EmployeeStatus string

var (
	DirectWorking         EmployeeStatus = "DIRECT_WORKING"
	IndirectWorking       EmployeeStatus = "INDIRECT_WORKING"
	Idle                  EmployeeStatus = "IDLE"
	Rest                  EmployeeStatus = "REST"
	OffDuty               EmployeeStatus = "OFF_DUTY"
	OffDutyWithoutEndTime EmployeeStatus = "OFF_DUTY_WITHOUT_END_TIME"
)

func (u *EmployeeStatusEntity) TableName() string {
	return "employee_status" // 自定义表名
}

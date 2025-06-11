/*
* @Author: supbro
* @Date:   2025/6/11 09:56
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:56
 */
package entity

import "time"

type EmployeeSummaryEntity struct {
	EmployeeNumber       string    `gorm:"column:employee_number" json:"employeeNumber"`
	EmployeeName         string    `gorm:"column:employee_name" json:"employeeName"`
	OperateDay           time.Time `gorm:"column:operate_day" json:"operateDay"`
	ProcessCode          string    `gorm:"column:process_code" json:"processCode"`
	PositionCode         string    `gorm:"column:position_code" json:"positionCode"` // 作业岗位编码
	WorkplaceCode        string    `gorm:"column:workplace_code" json:"workplaceCode"`
	WorkplaceName        string    `gorm:"column:workplace_name" json:"workplaceName"`
	EmployeePositionCode string    `gorm:"column:employee_position_code" json:"EmployeePositionCode"`
	WorkGroupCode        string    `gorm:"column:work_group_code" json:"workGroupCode"`       // 员工工作组编码
	RegionCode           string    `gorm:"column:region_code" json:"regionCode"`              // 工作点所属区域
	IndustryCode         string    `gorm:"column:industry_code" json:"industryCode"`          // 工作点所属行业
	SubIndustryCode      string    `gorm:"column:sub_industry_code" json:"subIndustryCode"`   // 工作点所属子行业
	DirectWorkTime       int       `gorm:"column:direct_work_time" json:"directWorkTime"`     // 直接作业时长（秒）
	IndirectWorkTime     int       `gorm:"column:indirect_work_time" json:"indirectWorkTime"` // 间接作业时长
	IdleTime             int       `gorm:"column:idle_time" json:"idleTime"`                  // 闲置时长
	RestTime             int       `gorm:"column:rest_time" json:"restTime"`
	AttendanceTime       int       `gorm:"column:attendance_time" json:"attendanceTime"`   // 出勤时长
	ProcessProperty      string    `gorm:"column:process_property" json:"processProperty"` // 环节额外属性
}

type WorkLoadWithSummaryEntity struct {
	EmployeeSummary *EmployeeSummaryEntity
	WorkLoad        map[string]float64
}

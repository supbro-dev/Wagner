/*
* @Author: supbro
* @Date:   2025/6/5 17:27
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/5 17:27
 */
package entity

import (
	"time"
)

// 小时人效聚合表
type HourSummaryResultEntity struct {
	BaseEntity
	OperateTime          time.Time `gorm:"column:operate_time;not null" json:"operateTime"`                            // 作业时间（小时）
	OperateDay           time.Time `gorm:"column:operate_day;not null" json:"operateDay"`                              // 作业日期
	ProcessCode          string    `gorm:"column:process_code;type:varchar(45);not null" json:"processCode"`           // 环节编码
	PositionCode         string    `gorm:"column:position_code;type:varchar(45)" json:"positionCode"`                  // 作业岗位编码
	WorkplaceCode        string    `gorm:"column:workplace_code;type:varchar(45);not null" json:"workplaceCode"`       // 工作点编码
	EmployeeNumber       string    `gorm:"column:employee_number;type:varchar(45);not null" json:"employeeNumber"`     // 员工工号
	EmployeeName         string    `gorm:"column:employee_name;type:varchar(45);not null" json:"employeeName"`         // 员工姓名
	EmployeePositionCode string    `gorm:"column:employee_position_code;type:varchar(45)" json:"employeePositionCode"` // 员工归属岗位
	WorkGroupCode        string    `gorm:"column:work_group_code;type:varchar(45)" json:"workGroupCode"`               // 员工工作组编码
	RegionCode           string    `gorm:"column:region_code" json:"regionCode"`                                       // 工作点所属区域
	IndustryCode         string    `gorm:"column:industry_code" json:"industryCode"`                                   // 工作点所属行业
	SubIndustryCode      string    `gorm:"column:sub_industry_code" json:"subIndustryCode"`                            // 工作点所属子行业
	WorkLoad             string    `gorm:"column:work_load;type:json" json:"workLoad"`                                 // 员工工作量
	DirectWorkTime       int       `gorm:"column:direct_work_time;not null;default:0" json:"directWorkTime"`           // 直接作业时长（秒）
	IndirectWorkTime     int       `gorm:"column:indirect_work_time;not null;default:0" json:"indirectWorkTime"`       // 间接作业时长
	IdleTime             int       `gorm:"column:idle_time;not null;default:0" json:"idleTime"`                        // 闲置时长
	AttendanceTime       int       `gorm:"column:attendance_time;not null;default:0" json:"attendanceTime"`            // 出勤时长
	ProcessProperty      string    `gorm:"column:process_property;type:json;not null" json:"processProperty"`          // 环节属性
	Properties           string    `gorm:"column:properties;type:json;" json:"properties"`                             // 其他属性
	IsDeleted            int8      `gorm:"column:is_deleted;not null;default:0" json:"isDeleted"`                      // 是否删除 (0-未删除 1-已删除)
	UniqueKey            string    `gorm:"column:unique_key;not null" json:"uniqueKey"`
}

// 设置表名
func (HourSummaryResultEntity) TableName() string {
	return "hour_summary_result"
}

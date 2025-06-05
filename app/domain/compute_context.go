/*
* @Author: supbro
* @Date:   2025/6/2 22:50
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 22:50
 */
package domain

import (
	"time"
	"wagner/app/service/calc_dynamic_param"
)

// 加工节点上下文
type ComputeContext struct {
	// 员工快照
	Employee EmployeeSnapshot
	// 日期
	OperateDay time.Time
	// 计算参数
	CalcOtherParam calc_dynamic_param.CalcOtherParam
	// 计算开始时间
	CalcStartTime time.Time
	// 计算结束时间
	CalcEndTime time.Time

	// 最近三天工作列表
	YesterdayWorkList []Work
	TodayWorkList     []Work
	TomorrowWorkList  []Work

	// 最近三天考勤列表
	YesterdayAttendance Attendance
	TodayAttendance     Attendance
	TomorrowAttendance  Attendance
	// 最近两天的排班信息
	YesterdayScheduling Scheduling
	TodayScheduling     Scheduling

	// 上下班信息
	YesterdayAttendanceStartTime time.Time
	YesterdayAttendanceEndTime   time.Time
	TodayAttendanceStartTime     time.Time
	TodayAttendanceEndTime       time.Time
	TomorrowAttendanceStartTime  time.Time

	TodayAttendanceNoMissing bool
}

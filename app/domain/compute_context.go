/*
* @Author: supbro
* @Date:   2025/6/2 22:50
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 22:50
 */
package domain

import "time"

// 加工节点上下文
type ComputeContext struct {
	// 员工快照
	Employee EmployeeSnapshot
	// 日期
	OperateDay time.Time

	// 最近三天工作列表
	YesterdayWorkList []Work
	TodayWorkList     []Work
	TomorrowWorkList  []Work

	// 最近三天考勤列表
	YesterdayAttendance Attendance
	TodayAttendance     Attendance
	TomorrowAttendance  Attendance
	// 当天排班信息
	TodayScheduling Scheduling
}

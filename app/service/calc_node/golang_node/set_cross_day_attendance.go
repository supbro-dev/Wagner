/*
* @Author: supbro
* @Date:   2025/6/5 14:03
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/5 14:03
 */
package golang_node

import (
	"wagner/app/domain"
)

// 设置跨天考勤上下班信息
func SetCrossDayAttendance(ctx *domain.ComputeContext) *domain.ComputeContext {
	// 未设置过当天考勤下班时间
	if &ctx.TodayAttendance != nil && &ctx.TodayAttendanceEndTime == nil {
		ctx.TodayAttendanceStartTime = ctx.TodayAttendance.StartTime
		ctx.TodayAttendanceEndTime = ctx.TodayAttendance.EndTime
	}
	// 未设置过前一天考勤下班时间
	if &ctx.YesterdayAttendance != nil && ctx.YesterdayAttendanceEndTime == nil {
		ctx.YesterdayAttendanceStartTime = ctx.YesterdayAttendance.StartTime
		ctx.YesterdayAttendanceEndTime = ctx.YesterdayAttendance.EndTime
	}
	// 未设置过第二天考勤上班时间
	if &ctx.TomorrowAttendance != nil && ctx.TomorrowAttendanceStartTime == nil {
		ctx.TomorrowAttendanceStartTime = ctx.TomorrowAttendance.StartTime
	}

	return ctx
}

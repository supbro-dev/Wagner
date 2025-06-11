/*
* @Author: supbro
* @Date:   2025/6/9 09:05
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/9 09:05
 */
package golang_node

import (
	"time"
	"wagner/app/domain"
)

// 计算处理当天默认考勤开始时间
func ComputeAttendanceDefaultStartTime(ctx *domain.ComputeContext) *domain.ComputeContext {
	if ctx.TodayAttendance != nil {
		if ctx.TodayAttendance.ComputedStartTime == nil && ctx.TodayAttendance.ComputedEndTime != nil {
			computedStartTime := computeAttendanceStartTime(ctx.TodayAttendance.EndTime, ctx.TodayScheduling, ctx.CalcOtherParam.Attendance.AttendanceAbsencePenaltyHour)
			ctx.TodayAttendance.ComputedStartTime = &computedStartTime
			ctx.TodayAttendanceStartTime = &computedStartTime
		}
	}
	return ctx
}

func computeAttendanceStartTime(todayAttendanceEndTime *time.Time, todayScheduling *domain.Scheduling, attendanceAbsencePenaltyHour int) time.Time {
	var defaultAttendanceStartTime time.Time
	if todayScheduling != nil {
		defaultAttendanceStartTime = *todayScheduling.StartTime
	} else if attendanceAbsencePenaltyHour > 0 {
		defaultAttendanceStartTime = todayAttendanceEndTime.Add(time.Duration(-attendanceAbsencePenaltyHour) * time.Hour)
	} else {
		if todayAttendanceEndTime.Hour() < 12 {
			defaultAttendanceStartTime = time.Date(
				todayAttendanceEndTime.Year(),
				todayAttendanceEndTime.Month(),
				todayAttendanceEndTime.Day(),
				12, // 小时设置为12
				todayAttendanceEndTime.Minute(),
				todayAttendanceEndTime.Second(),
				todayAttendanceEndTime.Nanosecond(),
				todayAttendanceEndTime.Location(),
			).AddDate(0, 0, -1)
		} else {
			defaultAttendanceStartTime = time.Date(
				todayAttendanceEndTime.Year(),
				todayAttendanceEndTime.Month(),
				todayAttendanceEndTime.Day(),
				0, // 小时设置为12
				todayAttendanceEndTime.Minute(),
				todayAttendanceEndTime.Second(),
				todayAttendanceEndTime.Nanosecond(),
				todayAttendanceEndTime.Location())
		}
	}

	return defaultAttendanceStartTime
}

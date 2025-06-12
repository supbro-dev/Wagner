/*
* @Author: supbro
* @Date:   2025/6/5 14:12
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/5 14:12
 */
package golang_node

import (
	"time"
	"wagner/app/domain"
	"wagner/app/utils/datetime_util"
)

// 缺卡情况设置默认下班时间
func ComputeAttendanceDefaultEndTime(ctx *domain.ComputeContext) *domain.ComputeContext {
	now := ctx.CalcStartTime

	// 今天下班卡缺卡
	if ctx.TodayAttendanceEndTime == nil && ctx.TodayAttendanceStartTime != nil {
		defaultEndTime := computeAttendanceEndTime(now, *ctx.TodayAttendanceStartTime, ctx.TodayScheduling, ctx.CalcOtherParam.Attendance.AttendanceAbsencePenaltyHour)
		ctx.TodayAttendanceEndTime = &defaultEndTime
	} else {
		ctx.TodayAttendanceNoMissing = true
	}

	// 前一天下班卡缺卡
	if ctx.YesterdayAttendanceEndTime == nil && ctx.YesterdayAttendanceStartTime != nil {
		defaultEndTime := computeAttendanceEndTime(now, *ctx.YesterdayAttendanceStartTime, ctx.YesterdayScheduling, ctx.CalcOtherParam.Attendance.AttendanceAbsencePenaltyHour)
		ctx.YesterdayAttendanceEndTime = &defaultEndTime
	}

	return ctx
}

// 计算缺卡情况下的考勤下班时间
// Parameters: now 当前时间, todayAttendanceStartTime 当天考勤上班还见, tomorrowAttendanceStartTime 第二天考勤上班时间，todayScheduling 排班信息， attendanceAbsencePenaltyHour 考勤缺卡惩罚时间(H)
// Returns:
func computeAttendanceEndTime(now, todayAttendanceStartTime time.Time, todayScheduling *domain.Scheduling, attendanceAbsencePenaltyHour int) time.Time {
	// 计算下班缺卡惩罚时间
	penaltyDefaultAttendanceEndTime := todayAttendanceStartTime.Add(time.Duration(attendanceAbsencePenaltyHour) * time.Hour)

	var defaultAttendanceEndTime time.Time
	// 如果有排班，使用排班下班时间
	if todayScheduling != nil && todayScheduling.EndTime != nil {
		defaultAttendanceEndTime = *todayScheduling.EndTime
	} else {
		// 上班时间在当日12点后, 默认下班时间为上班日期后一日12点
		if todayAttendanceStartTime.Hour() > 12 {
			defaultAttendanceEndTime = time.Date(
				todayAttendanceStartTime.Year(),
				todayAttendanceStartTime.Month(),
				todayAttendanceStartTime.Day(),
				12, // 小时设置为12
				todayAttendanceStartTime.Minute(),
				todayAttendanceStartTime.Second(),
				todayAttendanceStartTime.Nanosecond(),
				todayAttendanceStartTime.Location(),
			).AddDate(0, 0, 1)
		} else {
			// 上班时间在当日12点前, 默认下班时间为上班日期后一日0点
			defaultAttendanceEndTime = time.Date(
				todayAttendanceStartTime.Year(),
				todayAttendanceStartTime.Month(),
				todayAttendanceStartTime.Day(),
				0, // 小时设置为0
				todayAttendanceStartTime.Minute(),
				todayAttendanceStartTime.Second(),
				todayAttendanceStartTime.Nanosecond(),
				todayAttendanceStartTime.Location(),
			).AddDate(0, 0, 1)
		}
	}

	// 如果是当天，且当前时间处于惩罚结束时间和系统判定的结束时间之间，使用当前时间作为默认结束时间，为了处理当天员工未下班且已经工作到惩罚结束时间之后的场景
	if penaltyDefaultAttendanceEndTime.Before(now) && now.Before(defaultAttendanceEndTime) {
		defaultAttendanceEndTime = now
	} else {
		// 默认下班时间取计算下班时间和惩罚下班时间里最早的那个
		defaultAttendanceEndTime = datetime_util.Min(penaltyDefaultAttendanceEndTime, defaultAttendanceEndTime)
	}

	// 有节点单独处理，这里先不处理
	//// 考虑第二天上班时间的情况
	//if &tomorrowAttendanceStartTime != nil {
	//	// 如果默认下班时间超过第二天的上班时间，自动给下班时间设置为第二天上班时间减一秒，防止两天的考勤交叉
	//	if defaultAttendanceEndTime.Equal(tomorrowAttendanceStartTime) || defaultAttendanceEndTime.After(tomorrowAttendanceStartTime) {
	//		defaultAttendanceEndTime = tomorrowAttendanceStartTime.Add(-time.Second)
	//	}
	//}

	// 默认下班时间不超过当前时间
	if defaultAttendanceEndTime.After(now) {
		defaultAttendanceEndTime = now
	}

	return defaultAttendanceEndTime
}

/*
* @Author: supbro
* @Date:   2025/6/7 17:01
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/7 17:01
 */
package golang

import (
	"fmt"
	"time"
	"wagner/app/domain"
	"wagner/app/utils/datetime_util"
)

// 处理昨天/当天/第二天，3天考勤计算出来有重叠的场景
func CutOffAttendanceTime(ctx *domain.ComputeContext) *domain.ComputeContext {
	todayAttendance := ctx.TodayAttendance

	// 处理昨天下班时间大于当天上班时间
	if timeOverlap(*ctx.YesterdayAttendanceEndTime, *ctx.TodayAttendanceStartTime) {
		todayAttendanceStartTimeMinus1Sec := ctx.TodayAttendanceStartTime.Add(-time.Second)
		ctx.YesterdayAttendanceEndTime = &todayAttendanceStartTimeMinus1Sec
	}

	// 处理昨天上班时间大于昨天下班时间
	if timeOverlap(*ctx.YesterdayAttendanceStartTime, *ctx.YesterdayAttendanceEndTime) {
		ctx.YesterdayAttendanceStartTime = ctx.YesterdayAttendanceEndTime
	}

	// 处理当天下班时间大于第二天上班时间
	if timeOverlap(*ctx.TodayAttendanceEndTime, *ctx.TomorrowAttendanceStartTime) {
		tomorrowAttendanceStartTimeMinus1Sec := ctx.TomorrowAttendanceStartTime.Add(-time.Second)
		ctx.TodayAttendanceEndTime = &tomorrowAttendanceStartTimeMinus1Sec

		if todayAttendance != nil {
			todayAttendance.ComputedEndTime = &tomorrowAttendanceStartTimeMinus1Sec
			todayAttendance.AppendOperationMsg(
				fmt.Sprint(`今日原结束时间: %v, 超过第二天考勤开始时间: %v, 调整后: %v`,
					datetime_util.FormatDatetime(*todayAttendance.EndTime),
					datetime_util.FormatDatetime(*ctx.TomorrowAttendanceStartTime),
					datetime_util.FormatDatetime(*todayAttendance.ComputedEndTime)))
		}
	}

	// 处理当天上班时间大于当天下班时间
	if timeOverlap(*ctx.TodayAttendanceStartTime, *ctx.TodayAttendanceEndTime) {
		ctx.TodayAttendanceStartTime = ctx.TodayAttendanceEndTime

		if todayAttendance != nil {
			todayAttendance.ComputedStartTime = ctx.TodayAttendanceEndTime
			todayAttendance.AppendOperationMsg(
				fmt.Sprint(`今日原开始时间: %v, 超过今天考勤结束时间: %v, 调整后: %v`,
					datetime_util.FormatDatetime(*todayAttendance.StartTime),
					datetime_util.FormatDatetime(*todayAttendance.EndTime),
					datetime_util.FormatDatetime(*todayAttendance.ComputedStartTime)))
		}
	}

	return ctx
}

// 判断左边的时间是否与右边的时间重叠
func timeOverlap(leftTime, rightTime time.Time) bool {
	return &leftTime != nil && &rightTime != nil && (leftTime.Equal(rightTime) || leftTime.After(rightTime))
}

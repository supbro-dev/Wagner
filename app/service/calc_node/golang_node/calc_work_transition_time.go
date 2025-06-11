/*
* @Author: supbro
* @Date:   2025/6/10 12:39
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/10 12:39
 */
package golang_node

import (
	"fmt"
	"wagner/app/domain"
	"wagner/app/utils/datetime_util"
)

// 处理上班打卡到第一次作业开始的时间，处理最后一次作业完成到下班打卡的时间
func CalcWorkTransitionTime(ctx *domain.ComputeContext) *domain.ComputeContext {
	if ctx.TodayWorkList == nil || len(ctx.TodayWorkList) == 0 {
		return ctx
	}
	// 处理上班打卡到第一次作业开始的时间
	if ctx.TodayAttendanceStartTime != nil {
		maxRunUpTimeInMinute := ctx.CalcOtherParam.Attendance.MaxRunUpTimeInMinute
		firstWork := ctx.TodayWorkList[0]
		if firstWork.GetAction().ComputedStartTime.After(*ctx.TodayAttendanceStartTime) {
			diff := firstWork.GetAction().ComputedStartTime.Sub(*ctx.TodayAttendanceStartTime)
			// 小于最大开班时间，设置为工作，大于最大开班时间这里不处理，后续变为闲置
			if diff.Minutes() <= float64(maxRunUpTimeInMinute) {
				originalStartTime := firstWork.GetAction().ComputedStartTime
				firstWork.GetAction().ComputedStartTime = ctx.TodayAttendanceStartTime

				firstWork.GetAction().AppendOperationMsg(fmt.Sprintf("作业加入正常开班时间, 原开始时间: %v, 调整后: %v", datetime_util.FormatDatetime(*originalStartTime),
					datetime_util.FormatDatetime(*ctx.TodayAttendanceStartTime)))
			}
		}
	}

	// 处理最后一次作业完成到下班打卡的时间
	if ctx.TodayAttendanceEndTime != nil {
		lastWork := ctx.TodayWorkList[len(ctx.TodayWorkList)-1]

		if lastWork.GetAction().ComputedEndTime.Before(*ctx.TodayAttendanceEndTime) {
			minIdleTime := getOrDefaultMinIdleTime(lastWork.GetAction().Process, ctx.CalcOtherParam.Work.DefaultMinIdleTimeInMinute)

			diff := ctx.TodayAttendanceEndTime.Sub(*lastWork.GetAction().ComputedEndTime)
			if diff.Minutes() <= float64(minIdleTime) {
				originalEndTime := lastWork.GetAction().ComputedEndTime
				lastWork.GetAction().ComputedEndTime = ctx.TodayAttendanceEndTime

				lastWork.GetAction().AppendOperationMsg(fmt.Sprintf("作业完成到下班正常间隙调整, 原开始时间: %v, 调整后: %v", datetime_util.FormatDatetime(*originalEndTime),
					datetime_util.FormatDatetime(*ctx.TodayAttendanceEndTime)))
			}
		}
	}
	return ctx
}

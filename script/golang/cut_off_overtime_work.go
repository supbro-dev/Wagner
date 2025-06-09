/*
* @Author: supbro
* @Date:   2025/6/9 11:27
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/9 11:27
 */
package golang

import (
	"fmt"
	"time"
	"wagner/app/domain"
	"wagner/app/utils/datetime_util"
)

// 超长的Work进行截断
// 1.根据环节超长配置截断
// 2.根据下班卡截断（间接工作）
func CutOffOvertimeWork(ctx *domain.ComputeContext) *domain.ComputeContext {
	for _, work := range ctx.TodayWorkList {
		// 先进性超长截断
		maxDurationInMinute := getOrDefaultMaxTime(work.GetAction().Process, ctx.CalcOtherParam.Work.DefaultMaxTimeInMinute)

		diff := work.GetAction().ComputedEndTime.Sub(*work.GetAction().ComputedStartTime)
		if diff.Minutes() > float64(maxDurationInMinute) {
			originalEndTime := work.GetAction().ComputedEndTime
			computedEndTime := work.GetAction().ComputedStartTime.Add(time.Duration(maxDurationInMinute) * time.Minute)
			work.SetComputedEndTime(computedEndTime)
			work.GetAction().AppendOperationMsg(fmt.Sprintf("持续时长过长被系统截断, 持续时长: %v分, 最长时长: %v分, 原结束时间: %v, 调整后: %v",
				diff.Minutes(), maxDurationInMinute, datetime_util.FormatDatetime(*originalEndTime), datetime_util.FormatDatetime(computedEndTime)))
		}

		// 再进行考勤截断
		if ctx.TodayAttendance != nil {
			isWorkCrossAttendanceEnd := work.GetAction().ComputedStartTime.Before(*ctx.TodayAttendance.ComputedEndTime) && work.GetAction().ComputedEndTime.After(*ctx.TodayAttendance.ComputedEndTime)

			if isWorkCrossAttendanceEnd {
				originalEndTime := work.GetAction().ComputedEndTime
				computedEndTime := ctx.TodayAttendance.ComputedEndTime
				work.SetComputedEndTime(*computedEndTime)

				work.GetAction().AppendOperationMsg(fmt.Sprintf("持续时长过长被下班卡截断, 原结束时间: %v, 调整后: %v", datetime_util.FormatDatetime(*originalEndTime), datetime_util.FormatDatetime(*computedEndTime)))
			}
		}
	}

	return ctx
}

const MaxTimeKey = "maxTime"

func getOrDefaultMaxTime(process domain.StandardPosition, defaultMaxTime int) int {
	if process.Properties != nil {
		if value, exists := process.Properties[MaxTimeKey]; exists {
			return value.(int)
		}
	}

	return defaultMaxTime
}

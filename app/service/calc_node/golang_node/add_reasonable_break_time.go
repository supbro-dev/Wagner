/*
* @Author: supbro
* @Date:   2025/6/9 13:38
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/9 13:38
 */
package golang_node

import (
	"fmt"
	"sort"
	"wagner/app/domain"
	"wagner/app/utils/datetime_util"
)

// 在作业之间添加合理休息
func AddReasonableBreakTime(ctx *domain.ComputeContext) *domain.ComputeContext {
	if ctx.TodayWorkList == nil || len(ctx.TodayWorkList) == 0 {
		return ctx
	}

	actionList := make([]*domain.Action, 0)

	for _, work := range ctx.TodayWorkList {
		actionList = append(actionList, work.GetAction())
	}
	if ctx.TodayRestList != nil && len(ctx.TodayRestList) > 0 {
		for _, rest := range ctx.TodayRestList {
			actionList = append(actionList, &rest.Action)
		}
	}

	// 排序
	sort.Slice(actionList, func(i, j int) bool {
		return actionList[i].ComputedStartTime.Before(*actionList[j].ComputedStartTime)
	})

	for i, action := range actionList {
		var nextWorkOrRest *domain.Action

		if i < len(actionList)-1 {
			nextWorkOrRest = actionList[i+1]
		} else {
			nextWorkOrRest = nil
		}

		// 休息不会添加break时长
		if action.ActionType != domain.REST && nextWorkOrRest != nil {
			minIdleTime := getOrDefaultMinIdleTime(action.Process, ctx.CalcOtherParam.Work.DefaultMinIdleTimeInMinute)

			nextStartTime := nextWorkOrRest.ComputedStartTime
			diff := nextWorkOrRest.ComputedStartTime.Sub(*action.ComputedEndTime)

			if diff.Minutes() > 0 && diff.Minutes() <= float64(minIdleTime) {
				originalEndTime := action.ComputedEndTime
				action.ComputedEndTime = nextStartTime
				action.AppendOperationMsg(fmt.Sprintf("作业正常间隙调整, 原结束时间: %v, 调整后: %v",
					datetime_util.FormatDatetime(*originalEndTime), datetime_util.FormatDatetime(*nextStartTime)))
			}
		}
	}

	return ctx
}

const MinIdleTimeKey = "minIdleTimeInMinute"

func getOrDefaultMinIdleTime(process *domain.StandardPosition, defaultMinIdleTime int) int {
	if process.Properties != nil {
		if value, exists := process.Properties[MinIdleTimeKey]; exists {
			return value.(int)
		}
	}

	return defaultMinIdleTime
}

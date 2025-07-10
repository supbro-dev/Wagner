/*
* @Author: supbro
* @Date:   2025/6/9 13:21
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/9 13:21
 */
package golang_node

import (
	"fmt"
	"wagner/app/domain"
	"wagner/app/utils/datetime_util"
)

// 处理作业交叉截断
func CutOffCrossWork(ctx *domain.ComputeContext) *domain.ComputeContext {
	for i, work := range ctx.TodayWorkList {
		var nextWork domain.Actionable
		if i < len(ctx.TodayWorkList)-1 {
			nextWork = ctx.TodayWorkList[i+1]
		} else {
			nextWork = nil
		}

		if nextWork != nil {
			// 如果当前作业未结束，下一个作业已经开始，视为作业交叉，需要进行截断
			if work.GetAction().ComputedEndTime.After(*nextWork.GetAction().ComputedStartTime) {
				originalEndTime := work.GetAction().ComputedEndTime
				computedEndTime := nextWork.GetAction().ComputedStartTime

				work.(domain.Work).SetCutOffWorkCode(nextWork.GetAction().ActionCode)
				work.GetAction().ComputedEndTime = computedEndTime

				work.GetAction().AppendOperationMsg(fmt.Sprintf("作业间交叉被截断, 与%v存在交叉, 原结束时间: %v, 调整后: %v", nextWork.GetAction().ActionCode,
					datetime_util.FormatDatetime(*originalEndTime),
					datetime_util.FormatDatetime(*computedEndTime)))
			}
		}
	}
	return ctx
}

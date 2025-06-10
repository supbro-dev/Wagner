/*
* @Author: supbro
* @Date:   2025/6/9 11:11
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/9 11:11
 */
package golang

import (
	"fmt"
	"wagner/app/domain"
	"wagner/app/utils/datetime_util"
)

// 把未结束的作业结束时间设置为当前时间
func PaddingUnfinishedWorkEndTime(ctx *domain.ComputeContext) *domain.ComputeContext {
	now := ctx.CalcStartTime
	for _, work := range ctx.TodayWorkList {
		if work.GetAction().EndTime == nil {
			work.GetAction().ComputedEndTime = &now
			work.GetAction().AppendOperationMsg(fmt.Sprintf(`工作未结束, 结束时间调整到当前时间, 调整后: %v`, datetime_util.FormatDatetime(now)))
		}
	}
	return ctx
}

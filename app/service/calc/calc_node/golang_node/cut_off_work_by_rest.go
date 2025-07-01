/*
* @Author: supbro
* @Date:   2025/6/10 09:06
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/10 09:06
 */
package golang_node

import (
	"fmt"
	"github.com/jinzhu/copier"
	"sort"
	"wagner/app/domain"
	"wagner/app/utils/datetime_util"
)

// 作业被休息截断
func CutOffWorkByRest(ctx *domain.ComputeContext) *domain.ComputeContext {
	if ctx.TodayRestList == nil || len(ctx.TodayRestList) == 0 {
		return ctx
	}

	newWorkList := make([]domain.Actionable, 0)
	for _, work := range ctx.TodayWorkList {
		restNum := len(ctx.TodayRestList)
		var cutOffRestNum int
		newWork := cutOffWork(work, ctx.TodayRestList)
		for newWork != nil && cutOffRestNum < restNum {
			cutOffRestNum++
			newWorkList = append(newWorkList, newWork)
			newWork = cutOffWork(newWork, ctx.TodayRestList)
		}
	}

	ctx.TodayWorkList = append(ctx.TodayWorkList, newWorkList...)

	// 每次操作完workList，进行排序
	sort.Slice(ctx.TodayWorkList, func(i, j int) bool {
		if ctx.TodayWorkList[i].GetAction().ComputedStartTime.Before(*ctx.TodayWorkList[j].GetAction().ComputedStartTime) {
			return true
		} else if ctx.TodayWorkList[i].GetAction().ComputedStartTime.After(*ctx.TodayWorkList[j].GetAction().ComputedStartTime) {
			return false
		} else {
			return ctx.TodayWorkList[i].GetAction().ComputedEndTime.Before(*ctx.TodayWorkList[j].GetAction().ComputedEndTime)
		}
	})

	return ctx
}

// 考虑同一个作业可能被多个休息截断多次，所以传restList
func cutOffWork(work domain.Actionable, restList []*domain.Rest) domain.Work {
	for _, rest := range restList {
		if rest.ComputedStartTime.Equal(*rest.ComputedEndTime) {
			continue
		}
		// work |s         |e
		// rest    |s
		startBeforeEqualRestEndAfterRest := datetime_util.LeftBeforeOrEqualRight(*work.GetAction().ComputedStartTime, *rest.ComputedStartTime) &&
			work.GetAction().ComputedEndTime.After(*rest.ComputedStartTime)

		if startBeforeEqualRestEndAfterRest {
			originalEndTime := work.GetAction().ComputedEndTime
			// work |s          |e
			// rest    |s           |e
			if datetime_util.LeftBeforeOrEqualRight(*work.GetAction().ComputedEndTime, *rest.ComputedEndTime) {
				work.GetAction().ComputedEndTime = rest.ComputedStartTime
				work.GetAction().AppendOperationMsg(fmt.Sprintf("被休息开始截断, 原结束时间: %v, 调整后: %v",
					datetime_util.FormatDatetime(*originalEndTime), datetime_util.FormatDatetime(*work.GetAction().ComputedEndTime)))
			} else if work.GetAction().ComputedEndTime.After(*rest.ComputedEndTime) {
				// work |s        |e
				// rest     |s  |e

				// 设置新生成的作业
				var newWork domain.Work
				if directWork, ok := work.(*domain.DirectWork); ok {
					newDirectWork := &domain.DirectWork{}
					copier.Copy(&newDirectWork, directWork)
					newWork = newDirectWork
				} else if indirectWork, ok := work.(*domain.IndirectWork); ok {
					newIndirectWork := &domain.IndirectWork{}
					copier.Copy(&newIndirectWork, indirectWork)
					newWork = newIndirectWork
				}
				newWork.GetAction().Process = work.GetAction().Process
				newWork.GetAction().OperationMsgList = make([]string, 0)
				// 清空所有工作量
				if newDirectWork, ok := newWork.(*domain.DirectWork); ok {
					newDirectWork.WorkLoad = make(map[string]float64)
				}
				newWork.GetAction().ComputedStartTime = rest.ComputedEndTime
				newWork.GetAction().ComputedEndTime = work.GetAction().ComputedEndTime
				newWork.GetAction().AppendOperationMsg(fmt.Sprintf("因休息开始创建, 开始时间：%v，结束时间: %v", datetime_util.FormatDatetime(*rest.ComputedEndTime),
					datetime_util.FormatDatetime(*work.GetAction().ComputedEndTime)))

				// 处理原作业
				work.GetAction().ComputedEndTime = rest.ComputedStartTime
				work.GetAction().AppendOperationMsg(fmt.Sprintf("被休息开始截断, 原结束时间: %v, 调整后: %v",
					datetime_util.FormatDatetime(*originalEndTime), datetime_util.FormatDatetime(*work.GetAction().ComputedEndTime)))

				originalStartTime := work.GetAction().ComputedStartTime
				newWork.GetAction().AppendOperationMsg(fmt.Sprintf("被休息截断, 原开始时间: %v, 调整后: %v", datetime_util.FormatDatetime(*originalStartTime),
					datetime_util.FormatDatetime(*rest.ComputedEndTime)))

				return newWork
			}

		}
	}

	return nil
}

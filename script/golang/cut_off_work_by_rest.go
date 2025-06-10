/*
* @Author: supbro
* @Date:   2025/6/10 09:06
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/10 09:06
 */
package golang

import (
	"fmt"
	"github.com/jinzhu/copier"
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
		startBeforeEqualRestEndAfterRest := (work.GetAction().ComputedStartTime.Before(*rest.ComputedStartTime) || work.GetAction().ComputedStartTime.Equal(*rest.ComputedStartTime)) &&
			work.GetAction().ComputedEndTime.After(*rest.ComputedStartTime)

		if startBeforeEqualRestEndAfterRest {
			originalEndTime := work.GetAction().ComputedEndTime
			// work |s          |e
			// rest    |s           |e
			if work.GetAction().ComputedEndTime.Before(*rest.ComputedEndTime) || work.GetAction().ComputedEndTime.Equal(*rest.ComputedEndTime) {
				work.GetAction().ComputedEndTime = rest.ComputedStartTime
				work.GetAction().AppendOperationMsg(fmt.Sprintf("被休息开始截断, 原结束时间: %v, 调整后: %v",
					datetime_util.FormatDatetime(*originalEndTime), datetime_util.FormatDatetime(*work.GetAction().ComputedEndTime)))
			} else if work.GetAction().ComputedEndTime.After(*rest.ComputedEndTime) {
				// work |s        |e
				// rest     |s  |e
				var newWork domain.Work
				if _, ok := work.(*domain.DirectWork); ok {
					newWork = &domain.DirectWork{}
				} else if _, ok := work.(*domain.IndirectWork); ok {
					newWork = &domain.IndirectWork{}
				}

				copyErr := copier.Copy(&newWork, work)
				if copyErr != nil {
					panic(copyErr)
				}

				newWork.GetAction().Process = work.GetAction().Process
				newWork.GetAction().OperationMsgList = make([]string, 0)
				// 清空所有工作量
				if newDirectWork, ok := newWork.(*domain.DirectWork); ok {
					newDirectWork.WorkLoad = make(map[string]float64)
				}
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

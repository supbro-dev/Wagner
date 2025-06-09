/*
* @Author: supbro
* @Date:   2025/6/8 22:19
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/8 22:19
 */
package golang

import (
	"sort"
	"time"
	"wagner/app/domain"
)

// 根据参数配置的每个operateDay只计算x天内的数据，把早于这天的数据丢弃
func FilterExpiredData(ctx *domain.ComputeContext) *domain.ComputeContext {
	// 每一个operateDay只计算x天之内的数据
	lookBackDays := ctx.CalcOtherParam.Work.LookBackDays

	theMinDate := ctx.OperateDay.AddDate(0, 0, -(lookBackDays))

	if ctx.YesterdayWorkList != nil && len(ctx.YesterdayWorkList) > 0 {
		afterFilterYesterdayWorkList := make([]domain.Work, 0)
		for _, work := range ctx.YesterdayWorkList {
			if !isWorkBeforeTheMinDate(work, theMinDate) {
				afterFilterYesterdayWorkList = append(afterFilterYesterdayWorkList, work)
			}
		}
		ctx.YesterdayWorkList = afterFilterYesterdayWorkList
	}

	if ctx.TodayWorkList != nil && len(ctx.TodayWorkList) > 0 {
		afterFilterTodayWorkList := make([]domain.Work, 0)
		for _, work := range ctx.TodayWorkList {
			if !isWorkBeforeTheMinDate(work, theMinDate) {
				afterFilterTodayWorkList = append(afterFilterTodayWorkList, work)
			}
		}
		ctx.TodayWorkList = afterFilterTodayWorkList
	}

	if ctx.TomorrowWorkList != nil && len(ctx.TomorrowWorkList) > 0 {
		afterFilterTomorrowWorkList := make([]domain.Work, 0)
		for _, work := range ctx.TomorrowWorkList {
			if !isWorkBeforeTheMinDate(work, theMinDate) {
				afterFilterTomorrowWorkList = append(afterFilterTomorrowWorkList, work)
			}
		}
		ctx.TomorrowWorkList = afterFilterTomorrowWorkList
	}

	// 每次操作完workList，进行排序
	sort.Slice(ctx.TodayWorkList, func(i, j int) bool {
		return ctx.TodayWorkList[i].GetAction().ComputedStartTime.Before(*ctx.TodayWorkList[j].GetAction().ComputedStartTime)
	})

	return ctx
}

func isWorkBeforeTheMinDate(work domain.Work, theMinDate time.Time) bool {
	startTime := work.GetAction().ComputedStartTime
	endTime := work.GetAction().ComputedEndTime
	isStartTimeBefore := startTime != nil && (startTime.Equal(theMinDate) || startTime.Before(theMinDate))
	isEndTimeBefore := endTime != nil && (endTime.Equal(theMinDate) || endTime.Before(theMinDate))

	return isStartTimeBefore || isEndTimeBefore
}

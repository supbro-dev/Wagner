/*
* @Author: supbro
* @Date:   2025/6/7 22:12
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/7 22:12
 */
package golang

import (
	"sort"
	"wagner/app/domain"
)

// 把第二天/昨天两天的数据，在今天考勤上下班时间范围内的数据，归属到今天
func AddCrossDayData(ctx *domain.ComputeContext) *domain.ComputeContext {
	// 处理第二天的数据
	if ctx.TodayAttendanceEndTime != nil {
		tomorrowWorksBelongsToday := make([]domain.Actionable, 0)
		for _, tomorrowWork := range ctx.TomorrowWorkList {
			workComputedStartTime := (tomorrowWork).GetAction().ComputedStartTime
			if workComputedStartTime.Before(*ctx.TodayAttendanceEndTime) || workComputedStartTime.Equal(*ctx.TodayAttendanceEndTime) {
				tomorrowWorksBelongsToday = append(tomorrowWorksBelongsToday, tomorrowWork)
			}
		}
		if len(tomorrowWorksBelongsToday) > 0 {
			if len(ctx.TodayWorkList) == 0 {
				ctx.TodayWorkList = make([]domain.Actionable, 0)
			}
			for _, work := range tomorrowWorksBelongsToday {
				ctx.TodayWorkList = append(ctx.TodayWorkList, work)
			}
		}
	}

	// 处理昨天的数据
	if ctx.TomorrowAttendanceStartTime != nil {
		yesterdayWorksBelongsToday := make([]domain.Actionable, 0)
		for _, yesterdayWork := range ctx.YesterdayWorkList {
			workComputedStartTime := yesterdayWork.GetAction().ComputedStartTime
			if workComputedStartTime.Equal(*ctx.TodayAttendanceStartTime) || workComputedStartTime.After(*ctx.TodayAttendanceStartTime) {
				yesterdayWorksBelongsToday = append(yesterdayWorksBelongsToday, yesterdayWork)
			}
		}

		if len(yesterdayWorksBelongsToday) > 0 {
			if len(ctx.TomorrowWorkList) == 0 {
				ctx.TodayWorkList = make([]domain.Actionable, 0)
			}

			for _, work := range yesterdayWorksBelongsToday {
				ctx.TodayWorkList = append(ctx.TodayWorkList, work)
			}
		}
	}

	// 每次操作完workList，进行排序
	sort.Slice(ctx.TodayWorkList, func(i, j int) bool {
		return ctx.TodayWorkList[i].GetAction().ComputedStartTime.Before(*ctx.TodayWorkList[j].GetAction().ComputedStartTime)
	})

	return ctx
}

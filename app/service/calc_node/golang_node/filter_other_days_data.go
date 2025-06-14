/*
* @Author: supbro
* @Date:   2025/6/8 22:04
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/8 22:04
 */
package golang_node

import (
	"sort"
	"wagner/app/domain"
)

// 把今天数据中开始时间小于昨天考勤下班时间，或大于第二天考勤上班时间的数据过滤掉丢弃
func FilterOtherDaysData(ctx *domain.ComputeContext) *domain.ComputeContext {
	finalTodayWorkList := ctx.TodayWorkList

	// 排除属于昨天的数据
	if ctx.YesterdayAttendanceEndTime != nil {
		afterFilterYesterdayWorks := make([]domain.Actionable, 0)
		for _, work := range finalTodayWorkList {
			if work.GetAction().ComputedStartTime == nil || work.GetAction().ComputedStartTime.After(*ctx.YesterdayAttendanceEndTime) {
				afterFilterYesterdayWorks = append(afterFilterYesterdayWorks, work)
			}
		}

		finalTodayWorkList = afterFilterYesterdayWorks
	}

	// 测试排除属于第二天的数据
	if ctx.TomorrowAttendanceStartTime != nil {
		afterFilterTomorrowWorks := make([]domain.Actionable, 0)
		for _, work := range finalTodayWorkList {
			if work.GetAction().ComputedStartTime == nil || work.GetAction().ComputedStartTime.Before(*ctx.TomorrowAttendanceStartTime) {
				afterFilterTomorrowWorks = append(afterFilterTomorrowWorks, work)
			}
		}
		finalTodayWorkList = afterFilterTomorrowWorks
	}

	ctx.TodayWorkList = finalTodayWorkList
	ctx.YesterdayWorkList = nil
	ctx.TomorrowWorkList = nil

	// 每次操作完workList，进行排序
	sort.Slice(ctx.TodayWorkList, func(i, j int) bool {
		return ctx.TodayWorkList[i].GetAction().ComputedStartTime.Before(*ctx.TodayWorkList[j].GetAction().ComputedStartTime)
	})

	return ctx
}

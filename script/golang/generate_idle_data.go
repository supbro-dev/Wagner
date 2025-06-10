/*
* @Author: supbro
* @Date:   2025/6/10 13:05
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/10 13:05
 */
package golang

import (
	"time"
	"wagner/app/domain"
)

// 生成闲置工时
func GenerateIdleDataList(ctx *domain.ComputeContext) *domain.ComputeContext {
	idleList := make([]*domain.Idle, 0)
	var nextWork domain.Work
	for i, work := range ctx.TodayWorkList {
		if i < len(ctx.TodayWorkList)-1 {
			nextWork = ctx.TodayWorkList[i+1]
		} else {
			nextWork = nil
		}
		// 考虑考勤上班时间
		if i == 0 && ctx.TodayAttendanceStartTime != nil {
			if work.GetAction().ComputedStartTime.After(*ctx.TodayAttendanceStartTime) {
				idle := generateIdle(*ctx.TodayAttendanceStartTime, *work.GetAction().ComputedStartTime, work.GetAction().Process)
				idleList = append(idleList, idle)
			}
		} else if i == len(ctx.TodayWorkList)-1 && ctx.TodayAttendanceEndTime != nil {
			// 考虑考勤下班时间
			if work.GetAction().ComputedEndTime.Before(*ctx.TodayAttendanceEndTime) {
				idle := generateIdle(*work.GetAction().ComputedEndTime, *ctx.TodayAttendanceEndTime, work.GetAction().Process)
				idleList = append(idleList, idle)
			}
		}

		if nextWork != nil {

		}

	}
	return ctx
}

func generateIdle(startTime, endTime time.Time, process domain.StandardPosition) *domain.Idle {

}

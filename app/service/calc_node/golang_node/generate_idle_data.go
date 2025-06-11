/*
* @Author: supbro
* @Date:   2025/6/10 13:05
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/10 13:05
 */
package golang_node

import (
	"sort"
	"time"
	"wagner/app/domain"
	"wagner/app/service"
)

// 生成闲置工时
func GenerateIdleDataList(ctx *domain.ComputeContext) *domain.ComputeContext {
	// 把休息时间段放到切片中并进行排序
	todayActionList := ctx.TodayWorkList
	for _, rest := range ctx.TodayRestList {
		todayActionList = append(todayActionList, rest)
	}
	sort.Slice(todayActionList, func(i, j int) bool {
		return todayActionList[i].GetAction().ComputedStartTime.Before(*todayActionList[j].GetAction().ComputedStartTime)
	})

	idleList := make([]domain.Actionable, 0)
	var nextAction domain.Actionable
	for i, action := range todayActionList {
		if i < len(todayActionList)-1 {
			nextAction = todayActionList[i+1]
		} else {
			nextAction = nil
		}
		// 考虑考勤上班时间
		if i == 0 && ctx.TodayAttendanceStartTime != nil {
			if action.GetAction().ComputedStartTime.After(*ctx.TodayAttendanceStartTime) {
				idle := generateIdle(*ctx.TodayAttendanceStartTime, *action.GetAction().ComputedStartTime, action.GetAction().Process, action.GetAction().Properties)
				idleList = append(idleList, idle)
			}
		} else if i == len(todayActionList)-1 && ctx.TodayAttendanceEndTime != nil {
			// 考虑考勤下班时间
			if action.GetAction().ComputedEndTime.Before(*ctx.TodayAttendanceEndTime) {
				idle := generateIdle(*action.GetAction().ComputedEndTime, *ctx.TodayAttendanceEndTime, action.GetAction().Process, action.GetAction().Properties)
				idleList = append(idleList, idle)
			}
		}

		if nextAction != nil && action.GetAction().ComputedEndTime.Before(*nextAction.GetAction().ComputedStartTime) {
			idle := generateIdle(*action.GetAction().ComputedEndTime, *nextAction.GetAction().ComputedStartTime, action.GetAction().Process, action.GetAction().Properties)
			idleList = append(idleList, idle)
		}
	}

	// 处理全天闲置
	if ctx.TodayAttendanceStartTime != nil && ctx.TodayAttendanceEndTime != nil && (ctx.TodayRestList == nil || len(ctx.TodayRestList) == 0) {
		// 如果前后都没有环节，使用员工所属岗位下第一个环节
		standardPositionService := service.DomainHolder.StandardPositionService
		firstProcess := standardPositionService.FindPositionFirstProcess(ctx.Employee.PositionCode, ctx.Workplace.IndustryCode, ctx.Workplace.SubIndustryCode)

		idle := generateIdle(*ctx.TodayAttendanceStartTime, *ctx.TodayAttendanceEndTime, *firstProcess, make(map[string]interface{}))
		idleList = append(idleList, idle)
	}

	ctx.TodayIdleList = idleList
	// 最终的结果中加入休息、闲置
	todayActionList = append(todayActionList, idleList...)
	sort.Slice(todayActionList, func(i, j int) bool {
		return todayActionList[i].GetAction().ComputedStartTime.Before(*todayActionList[j].GetAction().ComputedStartTime)
	})

	ctx.TodayWorkList = todayActionList
	return ctx
}

func generateIdle(startTime, endTime time.Time, process domain.StandardPosition, properties map[string]interface{}) *domain.Idle {
	idle := &domain.Idle{
		Action: domain.Action{
			ComputedStartTime: &startTime,
			ComputedEndTime:   &endTime,
			Process:           process,
			Properties:        properties,
		},
	}
	return idle
}

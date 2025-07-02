/*
* @Author: supbro
* @Date:   2025/6/10 13:57
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/10 13:57
 */
package golang_node

import (
	"sort"
	"wagner/app/domain"
	"wagner/app/service"
)

// 匹配休息的归属环节（休息归属于他上一个环节， 上一个环节不存在，则归属于休息的下一个环节）
func MatchRestProcess(ctx *domain.ComputeContext) *domain.ComputeContext {
	// 把休息时间段放到切片中并进行排序
	todayActionList := ctx.TodayWorkList
	for _, rest := range ctx.TodayRestList {
		todayActionList = append(todayActionList, rest)
	}
	sort.Slice(todayActionList, func(i, j int) bool {
		if todayActionList[i].GetAction().ComputedStartTime.Before(*todayActionList[j].GetAction().ComputedStartTime) {
			return true
		} else if todayActionList[i].GetAction().ComputedStartTime.After(*todayActionList[j].GetAction().ComputedStartTime) {
			return false
		} else {
			return todayActionList[i].GetAction().ComputedEndTime.Before(*todayActionList[j].GetAction().ComputedEndTime)
		}
	})

	for i, action := range todayActionList {
		if rest, ok := action.(*domain.Rest); ok {
			process := findPreviousActionProcess(i-1, todayActionList)
			if process == nil {
				process = findNextActionProcess(i+1, todayActionList)
			}
			if process != nil {
				rest.GetAction().Process = process
				rest.GetAction().ProcessCode = process.Code
			} else {
				// 如果前后都没有环节，使用员工所属岗位下第一个环节
				processService := service.DomainHolder.ProcessService
				firstProcess := processService.FindFirstProcess(ctx.Employee.PositionCode, ctx.Workplace)

				rest.GetAction().Process = firstProcess
			}
		}
	}

	ctx.TodayWorkList = todayActionList

	return ctx
}

func findPreviousActionProcess(i int, actionList []domain.Actionable) *domain.ProcessPosition {
	if i < 0 {
		return nil
	}

	if &actionList[i].GetAction().Process != nil {
		return actionList[i].GetAction().Process
	} else {
		return findPreviousActionProcess(i-1, actionList)
	}
}

func findNextActionProcess(i int, actionList []domain.Actionable) *domain.ProcessPosition {
	if i > len(actionList)-1 {
		return nil
	}
	if &actionList[i].GetAction().Process != nil {
		return actionList[i].GetAction().Process
	} else {
		return findNextActionProcess(i+1, actionList)
	}
}

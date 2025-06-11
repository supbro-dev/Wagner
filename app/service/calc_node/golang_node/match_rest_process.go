/*
* @Author: supbro
* @Date:   2025/6/10 13:57
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/10 13:57
 */
package golang_node

import (
	"wagner/app/domain"
	"wagner/app/service"
)

// 匹配休息的归属环节（休息归属于他上一个环节， 上一个环节不存在，则归属于休息的下一个环节）
func MatchRestProcess(ctx *domain.ComputeContext) *domain.ComputeContext {
	for i, action := range ctx.TodayWorkList {
		if rest, ok := action.(*domain.Rest); ok {
			process := findLastActionProcess(i, ctx.TodayWorkList)
			if process == nil {
				process = findNextActionProcess(i, ctx.TodayWorkList)
			}
			if process != nil {
				rest.GetAction().Process = *process
			} else {
				// 如果前后都没有环节，使用员工所属岗位下第一个环节
				standardPositionService := service.DomainHolder.StandardPositionService
				firstProcess := standardPositionService.FindPositionFirstProcess(ctx.Employee.PositionCode, ctx.Workplace.IndustryCode, ctx.Workplace.SubIndustryCode)

				rest.GetAction().Process = *firstProcess
			}
		}
	}

	return ctx
}

func findLastActionProcess(i int, actionList []domain.Actionable) *domain.StandardPosition {
	if i <= 0 {
		return nil
	}

	if &actionList[i].GetAction().Process != nil {
		return &actionList[i].GetAction().Process
	} else {
		return findLastActionProcess(i-1, actionList)
	}
}

func findNextActionProcess(i int, actionList []domain.Actionable) *domain.StandardPosition {
	if i >= len(actionList)-1 {
		return nil
	}
	if &actionList[i].GetAction().Process != nil {
		return &actionList[i].GetAction().Process
	} else {
		return findNextActionProcess(i+1, actionList)
	}
}

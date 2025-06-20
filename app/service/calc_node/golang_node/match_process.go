/*
* @Author: supbro
* @Date:   2025/6/6 15:53
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 15:53
 */
package golang_node

import (
	"wagner/app/domain"
	"wagner/app/service/standard_position"
	"wagner/app/utils/script_util"
)

// 为所有作业匹配对应的环节
func MarchProcess(ctx *domain.ComputeContext) *domain.ComputeContext {
	for _, work := range ctx.TodayWorkList {
		process := findFirstProcess(work, ctx.ProcessList)

		if process != nil {
			work.GetAction().Process = process
			work.GetAction().ProcessCode = process.Code
		}
	}

	return ctx
}

// 遍历所有环节节点，根据表达式匹配到第一个环节
func findFirstProcess(work domain.Actionable, processList []*domain.StandardPosition) *domain.StandardPosition {
	for _, process := range processList {
		if process.Script == "" {
			continue
		}
		isThisProcess, err := script_util.Run[map[string]interface{}, bool]("", process.Script, script_util.EL, work.GetAction().Properties)

		if isThisProcess && err == nil {
			return process
		}
	}

	return standard_position.OtherProcess
}

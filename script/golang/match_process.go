/*
* @Author: supbro
* @Date:   2025/6/6 15:53
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 15:53
 */
package golang

import (
	"wagner/app/domain"
	"wagner/app/utils/script_util"
)

// 为所有作业匹配对应的环节
func MarchProcess(ctx *domain.ComputeContext) *domain.ComputeContext {
	for _, work := range ctx.TodayWorkList {
		process := findFirstProcess(&work, ctx.ProcessList)

		if process != nil {
			work.SetProcess(*process)
		}
	}
	return ctx
}

// 遍历所有环节节点，根据表达式匹配到第一个环节
func findFirstProcess(work *domain.Work, processList *[]domain.StandardPosition) *domain.StandardPosition {
	for _, process := range *processList {
		if isThisProcess, err := script_util.Run[map[string]interface{}, bool](process.ScriptName, process.Script, script_util.EL, (*work).GetProperties()); err != nil && isThisProcess {
			return &process
		}
	}

	return nil
}

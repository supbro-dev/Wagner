/*
* @Author: supbro
* @Date:   2025/6/4 20:19
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/4 20:19
 */
package golang

import "wagner/app/domain"

func ProcessChangeOperator(ctx *domain.ComputeContext) *domain.ComputeContext {

}

func handleChangeOperatorTime(actionList *[]domain.Action) {
	for _, action := range actionList {
		if action.IsChangeOperator() {

		}
	}
}

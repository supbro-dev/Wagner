/*
* @Author: supbro
* @Date:   2025/6/12 08:18
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 08:18
 */
package node

import (
	"wagner/app/domain"
	"wagner/app/service/calc/calc_dynamic_param"
	"wagner/app/utils/datetime_util"
)

func BuildTestCtx() domain.ComputeContext {
	ctx := domain.ComputeContext{}
	employee := domain.EmployeeSnapshot{
		Name:          "大兄弟",
		Number:        "A1001",
		WorkplaceCode: "workplace1",
		PositionCode:  "picker",
		WorkGroupCode: "workGroupCode1",
	}
	workplace := domain.Workplace{
		Code:            "workplace1",
		Name:            "1号工作点",
		RegionCode:      "NORTH",
		IndustryCode:    "FOOD",
		SubIndustryCode: "ConvenientFood",
	}
	operateDay, _ := datetime_util.ParseDate("2025-06-12")
	calcParam := calc_dynamic_param.DefaultCalcOtherParam
	calcStartTime, _ := datetime_util.ParseDatetime("2025-06-12 18:00:00")

	ctx.Employee = &employee
	ctx.Workplace = &workplace
	ctx.OperateDay = operateDay
	ctx.CalcOtherParam = calcParam
	ctx.CalcStartTime = calcStartTime

	return ctx
}

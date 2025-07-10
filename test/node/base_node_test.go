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
	calcParam := calc_dynamic_param.CalcOtherParam{
		Attendance: calc_dynamic_param.AttendanceParam{
			// 默认惩罚8小时
			AttendanceAbsencePenaltyHour: 8,
			MaxRunUpTimeInMinute:         20,
		},
		HourSummary: calc_dynamic_param.HourSummaryParam{
			// 默认聚合到结束的那个小时里
			WorkLoadAggregateType: calc_dynamic_param.AggregateEndHour,
		},
		Work: calc_dynamic_param.WorkParam{
			WorkLoadUnits: []calc_dynamic_param.WorkLoadUnit{
				{"件数", "itemNum"},
				{"SKU数", "skuNum"},
				{"包裹数", "packageNum"},
			},
			LookBackDays:               2,
			DefaultMaxTimeInMinute:     60,
			DefaultMinIdleTimeInMinute: 10,
		},
	}
	calcStartTime, _ := datetime_util.ParseDatetime("2025-06-12 18:00:00")

	ctx.Employee = &employee
	ctx.Workplace = &workplace
	ctx.OperateDay = operateDay
	ctx.CalcOtherParam = calcParam
	ctx.CalcStartTime = calcStartTime

	return ctx
}

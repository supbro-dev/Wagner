/*
* @Author: supbro
* @Date:   2025/6/11 09:52
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:52
 */
package query

import (
	"time"
	"wagner/app/domain"
	"wagner/app/service/calc_dynamic_param"
)

type HourSummaryResultQuery struct {
	WorkplaceCode      string
	EmployeeNumber     string
	DateRange          []*time.Time
	AggregateDimension domain.AggregateDimension
	IsCrossPosition    domain.IsCrossPosition
	WorkLoadUnit       []calc_dynamic_param.WorkLoadUnit
	CurrentPage        int
	PageSize           int
}

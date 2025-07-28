/*
* @Author: supbro
* @Date:   2025/6/11 09:52
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:52
 */
package query

import (
	"time"
)

type HourSummaryResultQuery struct {
	WorkplaceCode      string
	EmployeeNumber     string
	EmployeeNumberList []string
	DateRange          []*time.Time
	AggregateDimension string
	IsCrossPosition    string
	CurrentPage        int
	PageSize           int
}

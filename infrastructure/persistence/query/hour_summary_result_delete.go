/*
* @Author: supbro
* @Date:   2025/6/12 21:09
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 21:09
 */
package query

import "time"

type HourSummaryResultDelete struct {
	EmployeeNumber string
	WorkplaceCode  string
	OperateDay     time.Time
	UniqueKeyList  []string
}

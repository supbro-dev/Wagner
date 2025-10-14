/*
* @Author: supbro
* @Date:   2025/10/14 09:56
* @Last Modified by:   supbro
* @Last Modified time: 2025/10/14 09:56
 */
package qo

type EmployeeEfficiencyPageQueryQo struct {
	EmployeeNumber     string `json:"employeeNumber"`
	WorkplaceCode      string `json:"workplaceCode"`
	AggregateDimension string `json:"aggregateDimension"`
	IsCrossPosition    string `json:"isCrossPosition"`
	StartDateStr       string `json:"startDate"`
	EndDateStr         string `json:"endDate"`
	CurrentPage        string `json:"currentPage"`
	PageSize           string `json:"pageSize"`
}

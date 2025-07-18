/*
* @Author: supbro
* @Date:   2025/7/18 15:15
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/18 15:15
 */
package vo

type EmployeeVO struct {
	Name          string `json:"name"`
	Number        string `json:"number"`
	WorkplaceCode string `json:"workplaceCode"`
	WorkGroupCode string `json:"workGroupCode"`
}

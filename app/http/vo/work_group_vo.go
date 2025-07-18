/*
* @Author: supbro
* @Date:   2025/7/18 08:44
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/18 08:44
 */
package vo

type WorkGroupVO struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	PositionCode  string `json:"positionCode"`
	PositionName  string `json:"positionName"`
	WorkplaceCode string `json:"workplaceCode"`
	Desc          string `json:"desc"`
}

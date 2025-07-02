/*
* @Author: supbro
* @Date:   2025/7/2 16:43
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/2 16:43
 */
package vo

type ProcessImplementationPageVO struct {
	TableDataList []*ProcessImplementationVO `json:"tableDataList"`
	Page          *Page                      `json:"page"`
}

type ProcessImplementationVO struct {
	Key            string `json:"key"`
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	TargetTypeDesc string `json:"targetTypeDesc"`
	TargetName     string `json:"targetName"`
	Status         string `json:"status"`
	StatusDesc     string `json:"statusDesc"`
}

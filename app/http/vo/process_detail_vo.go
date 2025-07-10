/*
* @Author: supbro
* @Date:   2025/7/8 09:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/8 09:48
 */
package vo

type ProcessDetailVo struct {
	Id                      int64  `json:"id"`
	ParentCode              string `json:"parentCode"`
	Name                    string `json:"name"`
	Code                    string `json:"code"`
	Type                    string `json:"type"`
	TypeDesc                string `json:"typeDesc"`
	MaxTimeInMinute         string `json:"maxTimeInMinute"`
	MaxTimeInMinuteDesc     string `json:"maxTimeInMinuteDesc"`
	MinIdleTimeInMinute     string `json:"minIdleTimeInMinute"`
	MinIdleTimeInMinuteDesc string `json:"minIdleTimeInMinuteDesc"`
	Script                  string `json:"script"`
	WorkLoadRollUp          string `json:"workLoadRollUp"`
	WorkLoadRollUpDesc      string `json:"workLoadRollUpDesc"`
	SortIndex               int    `json:"sortIndex"`
}

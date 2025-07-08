/*
* @Author: supbro
* @Date:   2025/7/8 09:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/8 09:48
 */
package vo

type ProcessDetailVo struct {
	Id                  int64  `json:"id"`
	ProcessName         string `json:"processName"`
	ProcessCode         string `json:"processCode"`
	TypeDesc            string `json:"typeDesc"`
	MaxTimeInMinute     string `json:"maxTimeInMinute"`
	MinIdleTimeInMinute string `json:"minIdleTimeInMinute"`
	Script              string `json:"script"`
	WorkLoadRollUpDesc  string `json:"workLoadRollUpDesc"`
}

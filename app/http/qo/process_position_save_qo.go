/*
* @Author: supbro
* @Date:   2025/7/8 21:51
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/8 21:51
 */
package qo

type ProcessPositionSaveQo struct {
	ProcessImplId     string `json:"processImplId"`
	ParentProcessCode string `json:"parentProcessCode"`
	AddLevelType      string `json:"addLevelType"`
	ProcessName       string `json:"processName"`
	ProcessCode       string `json:"processCode"`
	Type              string `json:"type"`
	WorkLoadRollUp    string `json:"workLoadRollUp"`
}

type AddLevelType string

var (
	SameLevel AddLevelType = "sameLevel"
	NextLevel AddLevelType = "nextLevel"
)

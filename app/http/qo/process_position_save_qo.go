/*
* @Author: supbro
* @Date:   2025/7/8 21:51
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/8 21:51
 */
package qo

type ProcessPositionSaveQo struct {
	Id                  int64  `json:"id"`
	ProcessImplId       int64  `json:"processImplId"`
	ParentPositionCode  string `json:"parentPositionCode"`
	AddLevelType        string `json:"addLevelType"`
	Name                string `json:"name"`
	Code                string `json:"code"`
	Type                string `json:"type"`
	WorkLoadRollUp      string `json:"workLoadRollUp"`
	MaxTimeInMinute     string `json:"maxTimeInMinute"`
	MinIdleTimeInMinute string `json:"minIdleTimeInMinute"`
	SortIndex           int    `json:"sortIndex"`
	Script              string `json:"script"`
}

type AddLevelType string

var (
	SameLevel AddLevelType = "sameLevel"
	NextLevel AddLevelType = "nextLevel"
)

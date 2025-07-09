/*
* @Author: supbro
* @Date:   2025/7/7 20:18
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/7 20:18
 */
package vo

type ProcessPositionTreeNodeVo struct {
	Id             int64                        `json:"id"`
	Title          string                       `json:"title"`
	Key            string                       `json:"key"`
	Type           string                       `json:"type"`
	ParentName     string                       `json:"parentName"`
	ParentCode     string                       `json:"parentCode"`
	WorkLoadRollUp string                       `json:"workLoadRollUp"`
	SortIndex      int                          `json:"sortIndex"`
	Children       []*ProcessPositionTreeNodeVo `json:"children"`
}

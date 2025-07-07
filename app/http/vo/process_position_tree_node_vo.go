/*
* @Author: supbro
* @Date:   2025/7/7 20:18
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/7 20:18
 */
package vo

type ProcessPositionTreeNodeVo struct {
	Title    string                       `json:"title"`
	Key      string                       `json:"key"`
	Type     string                       `json:"type"`
	Children []*ProcessPositionTreeNodeVo `json:"children"`
}

/*
* @Author: supbro
* @Date:   2025/7/7 16:53
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/7 16:53
 */
package domain

import "wagner/infrastructure/persistence/entity"

type ProcessPositionTreeNode struct {
	Id             int64
	Name           string
	Code           string
	ParentName     string
	ParentCode     string
	Type           entity.ProcessPositionType
	Children       []*ProcessPositionTreeNode
	WorkLoadRollUp string
	SortIndex      int
}

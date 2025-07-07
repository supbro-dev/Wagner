/*
* @Author: supbro
* @Date:   2025/7/7 16:53
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/7 16:53
 */
package domain

import "wagner/infrastructure/persistence/entity"

type ProcessPositionTreeNode struct {
	Name     string
	Code     string
	Type     entity.ProcessPositionType
	Children []*ProcessPositionTreeNode
}

/*
* @Author: supbro
* @Date:   2025/6/5 10:30
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/5 10:50
 */
package domain

import "wagner/infrastructure/persistence/entity"

type ProcessPosition struct {
	Id         int64
	Name       string
	Code       string
	ParentCode string
	Type       entity.ProcessPositionType
	// 层级（1代表一级部门、2代表2级部门，最后一级为环节，倒数第二级为岗位）
	Level int
	// 最大部门层级
	MaxDeptLevel int
	// 环节的属性
	Properties map[string]interface{}
	// 路径
	Path []*ProcessPosition
	// 环节匹配执行脚本
	Script string
	// 标识一套StandardPosition的唯一版本
	Version   int
	SortIndex int
}

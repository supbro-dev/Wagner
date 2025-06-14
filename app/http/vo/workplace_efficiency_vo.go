/*
* @Author: supbro
* @Date:   2025/6/14 10:09
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/14 10:09
 */
package vo

type WorkplaceEfficiencyVO struct {
	Root    *WorkplaceStructureVO `json:"root"`
	Columns []*TableColumnVO      `json:"columns"`
}

type WorkplaceStructureVO struct {
	Name string `json:"name"`
	Code string `json:"code"`
	// 层级（1代表一级部门、2代表2级部门，最后一级为环节，倒数第二级为岗位）
	Level int `json:"level"`
	// 最大部门层级
	MaxDeptLevel int `json:"maxDeptLevel"`
	// 工作量
	WorkLoad map[string]float64 `json:"workLoad"`
	// 工时
	DirectWorkTime   float64 `json:"directWorkTime"`
	IndirectWorkTime float64 `json:"indirectWorkTime"` // 间接作业时长
	IdleTime         float64 `json:"idleTime"`         // 闲置时长
	RestTime         float64 `json:"restTime"`
	AttendanceTime   float64 `json:"attendanceTime"` // 出勤时长

	Children []*WorkplaceStructureVO `json:"children"`

	WorkLoadRollUp bool // 工作量是否向上汇总
}

/*
* @Author: supbro
* @Date:   2025/6/11 09:24
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:24
 */
package vo

type EmployeeEfficiencyVO struct {
	TableDataList []*EmployeeSummaryVO `json:"tableDataList"`
	Columns       []*TableColumnVO     `json:"columns"`
}

type EmployeeSummaryVO struct {
	EmployeeNumber       string             `json:"employeeNumber"`
	EmployeeName         string             `json:"employeeName"`
	OperateDay           string             `json:"operateDay"`
	ProcessCode          string             `json:"processCode"`
	PositionCode         string             `json:"positionCode"` // 作业岗位编码
	WorkplaceCode        string             `json:"workplaceCode"`
	WorkplaceName        string             `json:"workplaceName"`
	EmployeePositionCode string             `json:"EmployeePositionCode"`
	WorkGroupCode        string             `json:"workGroupCode"`    // 员工工作组编码
	RegionCode           string             `json:"regionCode"`       // 工作点所属区域
	IndustryCode         string             `json:"industryCode"`     // 工作点所属行业
	SubIndustryCode      string             `json:"subIndustryCode"`  // 工作点所属子行业
	DirectWorkTime       int                `json:"directWorkTime"`   // 直接作业时长（秒）
	IndirectWorkTime     int                `json:"indirectWorkTime"` // 间接作业时长
	IdleTime             int                `json:"idleTime"`         // 闲置时长
	RestTime             int                `json:"restTime"`
	AttendanceTime       int                `json:"attendanceTime"`  // 出勤时长
	ProcessProperty      string             `json:"processProperty"` // 环节额外属性
	WorkLoad             map[string]float64 `json:"workLoad"`
}

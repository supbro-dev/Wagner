/*
* @Author: supbro
* @Date:   2025/6/18 11:31
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/18 11:31
 */
package vo

type EmployeeStatusVO struct {
	GroupStatus []*GroupStatusVO `json:"groupStatus"`
}

type GroupStatusVO struct {
	GroupCode          string                  `json:"groupCode"`
	GroupName          string                  `json:"groupName"`
	EmployeeStatusList []*EachEmployeeStatusVO `json:"employeeStatusList"`
	GroupStatusNum     *GroupStatusNumVO       `json:"groupStatusNum"`
}

type EachEmployeeStatusVO struct {
	EmployeeNumber string `json:"employeeNumber"`
	EmployeeName   string `json:"employeeName"`
	StatusDesc     string `json:"statusDesc"`
	LastActionDesc string `json:"lastActionDesc"`
}

type GroupStatusNumVO struct {
	DirectWorkingNum         int `json:"directWorkingNum"`
	IndirectWorkingNum       int `json:"indirectWorkingNum"`
	IdleNum                  int `json:"idleNum"`
	RestNum                  int `json:"restNum"`
	OffDutyNum               int `json:"offDutyNum"`
	OffDutyWithoutEndTimeNum int `json:"offDutyWithoutEndTimeNum"`
}

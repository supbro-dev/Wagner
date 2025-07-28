/*
* @Author: supbro
* @Date:   2025/6/13 08:02
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/13 08:02
 */
package vo

import (
	"time"
	"wagner/app/domain"
)

type TimeOnTaskVO struct {
	EmployeeNumber string `json:"employeeNumber"`
	EmployeeName   string `json:"employeeName"`
	WorkplaceName  string `json:"workplaceName"`
	RegionCode     string `json:"regionCode"`
	OperateDay     string `json:"operateDay"`

	Attendance          *AttendanceVO        `json:"attendance"`
	Scheduling          *SchedulingVO        `json:"scheduling"`
	ProcessDurationList []*ProcessDurationVO `json:"processDurationList"`
}

type AttendanceVO struct {
	ActionType domain.ActionType `json:"actionType"`
	StartTime  time.Time         `json:"startTime"`
	EndTime    time.Time         `json:"endTime"`
}

type SchedulingVO struct {
	ActionType domain.ActionType `json:"actionType"`
	StartTime  time.Time         `json:"startTime"`
	EndTime    time.Time         `json:"endTime"`
	RestList   []RestVO          `json:"restList"`
}

type RestVO struct {
	ActionType domain.ActionType `json:"actionType"`
	StartTime  time.Time         `json:"startTime"`
	EndTime    time.Time         `json:"endTime"`
}

type ProcessDurationVO struct {
	Id            string                    `json:"id"`
	StartTime     time.Time                 `json:"startTime"`
	EndTime       time.Time                 `json:"endTime"`
	ActionType    domain.ActionType         `json:"actionType"`
	ProcessCode   string                    `json:"processCode"`
	ProcessName   string                    `json:"processName"`
	WorkplaceName string                    `json:"workplaceName"`
	WorkLoad      map[string]float64        // 不对外透传
	WorkLoadDesc  string                    `json:"workLoadDesc"`
	Duration      float64                   `json:"duration"`
	Details       []ProcessDurationDetailVO `json:"details"`
}

type ProcessDurationDetailVO struct {
	StartTime        string  `json:"startTime"`
	EndTime          string  `json:"endTime"`
	Duration         float64 `json:"duration"`
	OperationMessage string  `json:"operationMessage"`
}

package domain

import "time"

type Action struct {
	EmployeeNumber string
	WorkplaceCode  string
	OperateDay     time.Time
	ActionCode     string
	ActionType     string
	StartTime      time.Time
	EndTime        time.Time
	Properties     string
}

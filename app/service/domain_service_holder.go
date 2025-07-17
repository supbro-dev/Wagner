package service

import (
	"wagner/app/service/action"
	"wagner/app/service/calc/calc_dynamic_param"
	"wagner/app/service/employee_snapshot"
	"wagner/app/service/position"
	"wagner/app/service/process"
	"wagner/app/service/work_group"
	"wagner/app/service/workplace"
)

type DomainServiceHolder struct {
	EmployeeSnapshotService *employee_snapshot.EmployeeSnapshotService
	ActionService           *action.ActionService
	ProcessService          process.ProcessService
	CalcDynamicParamService *calc_dynamic_param.CalcDynamicParamService
	WorkplaceService        *workplace.WorkplaceService
	PositionService         *position.PositionService
	WorkGroupService        *work_group.WorkGroupService
}

var (
	DomainHolder DomainServiceHolder
)

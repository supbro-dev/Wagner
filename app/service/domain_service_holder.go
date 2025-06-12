package service

import (
	"wagner/app/service/action"
	"wagner/app/service/calc_dynamic_param"
	"wagner/app/service/employee_snapshot"
	"wagner/app/service/standard_position"
	"wagner/app/service/workplace"
)

type DomainServiceHolder struct {
	EmployeeSnapshotService *employee_snapshot.EmployeeSnapshotService
	ActionService           *action.ActionService
	StandardPositionService standard_position.StandardPositionItf
	CalcDynamicParamService *calc_dynamic_param.CalcDynamicParamService
	WorkplaceService        *workplace.WorkplaceService
}

var (
	DomainHolder DomainServiceHolder
)

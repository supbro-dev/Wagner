package service

import (
	"wagner/app/service/action"
	"wagner/app/service/employee_snapshot"
	"wagner/app/service/standard_position"
)

type DomainServiceHolder struct {
	EmployeeSnapshotService *employee_snapshot.EmployeeSnapshotService
	ActionService           *action.ActionService
	StandardPositionService *standard_position.StandardPositionService
}

var (
	DomainHolder DomainServiceHolder
)
